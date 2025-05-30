package cli
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"

	"github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/html"
)

// DefaultAQLQuery is the default query to fetch gene data from ArangoDB.
// Exported for use in flag.go
const DefaultAQLQuery = `
FOR ftype IN cvterm
 FOR feat IN feature
    FOR dbx IN dbxref
        FILTER ftype.name == 'gene'
        FILTER feat.type_id == ftype.cvterm_id
        FILTER feat.dbxref_id == dbx.dbxref_id
        LET props = (
            FOR fprop IN featureprop
                FOR cvt IN cvterm
                    FILTER cvt.name IN ['description','name description']
                    FILTER feat.feature_id == fprop.feature_id
                    FILTER fprop.type_id == cvt.cvterm_id
                    RETURN {
                        name: cvt.name,
                        value: fprop.value
                    }
        )
        FILTER LENGTH(props) > 0
        RETURN {
            id: dbx.accession,
            props: props
        }
`

// queryArangoParams holds the parameters for the queryArango function.
type queryArangoParams struct {
	ctx            context.Context
	config         AppConfig // Contains AQLQuery and Logger
	arangoDocsChan chan<- ArangoResultDoc
	mainCancel     context.CancelFunc
}

// AppConfig holds all configuration for the application.
type AppConfig struct {
	AQLQuery             string
	ArangoUser           string // For authorship in gRPC updates
	NumProcessingWorkers int
	NumGrpcWorkers       int
	Logger               *logrus.Entry
}

// ArangoResultDoc represents the structure of a document from ArangoDB.
type ArangoResultDoc struct {
	ID    string           `json:"id"` // This is dbx.accession, likely the feature_id
	Props []ArangoProperty `json:"props"`
}

// ArangoProperty represents a single property object from ArangoDB.
type ArangoProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProcessedGeneData holds the gene ID and its list of HTML-stripped property values.
type ProcessedGeneData struct {
	GeneID            string
	StrippedPropsText []StrippedProperty
}

// StrippedProperty holds the original property name and its stripped text.
type StrippedProperty struct {
	OriginalName string
	StrippedText string
}

// GrpcUpdateResult holds the result of a gRPC update operation.
type GrpcUpdateResult struct {
	GeneID  string
	Success bool
	Message string
	Error   error
}

// HTML Stripping Utilities
var (
	spaceNormalizerRegexp = regexp.MustCompile(`\s+`)
)
func stripHTMLWithParser(htmlString string) (string, error) {
	processedHTMLString := strings.ReplaceAll(htmlString, "\\n", "")
	doc, err := html.Parse(strings.NewReader(processedHTMLString))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}
	var b strings.Builder
	extractTextFromNode(doc, &b)
	normalizedText := spaceNormalizerRegexp.ReplaceAllString(b.String(), " ")
	return strings.TrimSpace(normalizedText), nil
}

func extractTextFromNode(node *html.Node, b *strings.Builder) {
	if node == nil {
		return
	}
	if node.Type == html.ElementNode {
		switch node.Data {
		case "script", "style", "noscript", "iframe", "noembed":
			return
		}
	}
	if node.Type == html.CommentNode {
		return
	}
	if node.Type == html.TextNode {
		b.WriteString(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		extractTextFromNode(child, b)
	}
}

// queryArango is responsible for querying ArangoDB and sending documents to a channel.
// It calls params.mainCancel if critical errors occur or the context is done.
func queryArango(params *queryArangoParams) {
	dbh := registry.GetArangodbConnection()
	cursor, err := dbh.Search(params.config.AQLQuery)
	if err != nil {
		params.config.Logger.Errorf("error in arangodb query %v", err)
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	defer cursor.Close()
	if cursor.IsEmpty() {
		params.config.Logger.Warn("ArangoDB querier finished (no data).")
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	docCount := 0
	for cursor.Scan() {
		select {
		case <-params.ctx.Done():
			return
		default:
			var doc ArangoResultDoc
			if errRead := cursor.Read(&doc); errRead != nil {
				params.config.Logger.Errorf(
					"Failed to read document from ArangoDB cursor: %v",
					errRead,
				)
				params.mainCancel() // Signal other goroutines to stop
			}
			// Send the document to the channel, handling context
			// cancellation
			select {
			case params.arangoDocsChan <- doc:
				docCount++
			case <-params.ctx.Done():
				return
			}
		}
	}
	params.config.Logger.Infof(
		"Successfully fetched %d documents from ArangoDB.",
		docCount,
	)
}

func grpcUpdateWorkerFunc(
	config AppConfig,
	grpcClient feature_annotation.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[ProcessedGeneData, GrpcUpdateResult] {
	return func(ctx context.Context, job concurrent.Job[ProcessedGeneData]) (GrpcUpdateResult, error) {
		logger := config.Logger
		processedData := job.Payload
		logger.Debugf(
			"gRPC Worker (Job %s): updating gene ID: %s",
			job.ID,
			processedData.GeneID,
		)
		result := GrpcUpdateResult{GeneID: processedData.GeneID, Success: false}
		featAnno, err := grpcClient.GetFeatureAnnotation(
			ctx, &feature_annotation.FeatureAnnotationId{
				Id: processedData.GeneID,
			})
		if err != nil {
			result.Message = fmt.Sprintf(
				"failed to GetFeatureAnnotation: %v",
				err,
			)
			result.Error = err
			return result, err
		}
		for _, prop := range processedData.StrippedPropsText {
			_, err := grpcClient.AddTag(ctx,
				&feature_annotation.AddTagRequest{
					Id: featAnno.Id,
					Tag: &feature_annotation.TagPropertyCreate{
						Tag:       prop.OriginalName,
						Value:     prop.StrippedText,
						CreatedBy: config.ArangoUser,
					},
				})
			if err != nil {
				errMsg := fmt.Sprintf(
					"failed to AddTag for property %s: %v",
					prop.OriginalName,
					err,
				)
				result.Message = errMsg
				result.Error = err
				logger.Errorf(
					"gRPC Worker (Job %s): %s for gene ID %s",
					job.ID, errMsg, processedData.GeneID,
				)
				return result, err
			}
			logger.Debugf(
				"gRPC Worker (Job %s): successfully added tag %s for gene ID %s",
				job.ID,
				prop.OriginalName,
				processedData.GeneID,
			)
		}
		result.Success = true
		result.Message = fmt.Sprintf(
			"Successfully added %d tags",
			len(processedData.StrippedPropsText),
		)
		return result, nil
	}
}
