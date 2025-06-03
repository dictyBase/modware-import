package cli

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// HTML Stripping Utilities
var (
	spaceNormalizerRegexp = regexp.MustCompile(`\s+`)
)

func stripHTMLWithParser(htmlString string) (string, error) {
	doc, err := html.Parse(
		strings.NewReader(strings.ReplaceAll(htmlString, "\\n", "")),
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}
	var b strings.Builder
	extractTextFromNode(doc, &b)
	return strings.TrimSpace(
		spaceNormalizerRegexp.ReplaceAllString(b.String(), " "),
	), nil
}

func extractTextFromNode(node *html.Node, b *strings.Builder) {
	if node == nil {
		return
	}
	switch node.Type {
	case html.ElementNode:
		// Skip content of certain tags
		switch node.Data {
		case "script", "style", "noscript", "iframe", "noembed":
			return
		}
	case html.CommentNode:
		return
	case html.TextNode:
		b.WriteString(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		extractTextFromNode(child, b)
	}
}

func htmlProcessingWorkerFunc(
	logger *logrus.Entry,
) concurrent.WorkerFunc[ArangoResultDoc, ProcessedGeneData] {
	return func(
		ctx context.Context,
		job concurrent.Job[ArangoResultDoc],
	) (ProcessedGeneData, error) {
		arangoDoc := job.Payload
		var strippedProps []StrippedProperty
		for _, prop := range arangoDoc.Props {
			strippedText, err := stripHTMLWithParser(prop.Value)
			if err != nil {
				return ProcessedGeneData{
						GeneID: arangoDoc.ID,
					}, fmt.Errorf(
						"failed to strip HTML for gene %s, property %s (Job %s): %w",
						arangoDoc.ID,
						prop.Name,
						job.ID,
						err,
					)
			}
			strippedProps = append(
				strippedProps,
				StrippedProperty{
					OriginalName: prop.Name,
					StrippedText: strippedText,
				},
			)
		}
		logger.Debugf(
			"HTML Worker (Job %s): finished processing gene ID: %s, found %d stripped props.",
			job.ID,
			arangoDoc.ID,
			len(strippedProps),
		)
		return ProcessedGeneData{
			GeneID:            arangoDoc.ID,
			StrippedPropsText: strippedProps,
		}, nil
	}
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
						CreatedBy: DefaultUserName,
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
