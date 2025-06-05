package cli

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HTML Stripping Utilities
var (
	spaceNormalizerRegexp = regexp.MustCompile(`\s+`)
)

type processSinglePropertyUpdateParams struct {
	ctx        context.Context
	grpcClient feature_annotation.FeatureAnnotationServiceClient
	featAnno   *feature_annotation.FeatureAnnotation
	prop       StrippedProperty
	arangoUser string
	logger     *logrus.Entry
	jobID      string
}

type handleCreateFeatureAnnotationParams struct {
	ctx           context.Context
	grpcClient    feature_annotation.FeatureAnnotationServiceClient
	processedData ProcessedGeneData
	jobID         string
	logger        *logrus.Entry
}

type handleUpdateExistingFeatureAnnotationParams struct {
	ctx           context.Context
	grpcClient    feature_annotation.FeatureAnnotationServiceClient
	featAnno      *feature_annotation.FeatureAnnotation
	processedData ProcessedGeneData
	config        AppConfig
	jobID         string
	logger        *logrus.Entry
}

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

func processSinglePropertyUpdate(
	params *processSinglePropertyUpdateParams,
) error {
	if slices.ContainsFunc(
		params.featAnno.Attributes.Properties,
		func(existingTag *feature_annotation.TagProperty) bool {
			return existingTag.Tag == params.prop.OriginalName
		},
	) {
		_, err := params.grpcClient.UpdateTag(
			params.ctx,
			&feature_annotation.UpdateTagRequest{
				Id: params.featAnno.Id,
				Tag: &feature_annotation.TagPropertyUpdate{
					Tag:       params.prop.OriginalName,
					Value:     params.prop.StrippedText,
					UpdatedBy: DefaultUserName,
				},
			},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to UpdateTag for property %s: %v",
				params.prop.OriginalName,
				err,
			)
		}
		params.logger.Debugf(
			"gRPC Worker (Job %s): successfully updated tag %s for gene ID %s",
			params.jobID,
			params.prop.OriginalName,
			params.featAnno.Id,
		)
	} else {
		// Add new tag
		_, err := params.grpcClient.AddTag(params.ctx,
			&feature_annotation.AddTagRequest{
				Id: params.featAnno.Id,
				Tag: &feature_annotation.TagPropertyCreate{
					Tag:       params.prop.OriginalName,
					Value:     params.prop.StrippedText,
					CreatedBy: DefaultUserName,
				},
			})
		if err != nil {
			return fmt.Errorf(
				"failed to AddTag for property %s: %v",
				params.prop.OriginalName,
				err,
			)
		}
		params.logger.Debugf(
			"gRPC Worker (Job %s): successfully added tag %s for gene ID %s",
			params.jobID,
			params.prop.OriginalName,
			params.featAnno.Id,
		)
	}
	return nil
}

func strippedPropertyToTagProperty(
	prop StrippedProperty,
) *feature_annotation.TagProperty {
	return &feature_annotation.TagProperty{
		Tag:       prop.OriginalName,
		Value:     prop.StrippedText,
		CreatedBy: DefaultUserName,
	}
}

func handleCreateFeatureAnnotation(
	params *handleCreateFeatureAnnotationParams,
) (GrpcUpdateResult, error) {
	result := GrpcUpdateResult{
		GeneID:  params.processedData.GeneID,
		Success: false,
	}
	_, createErr := params.grpcClient.CreateFeatureAnnotation(
		params.ctx,
		&feature_annotation.NewFeatureAnnotation{
			Id:        params.processedData.GeneID,
			CreatedBy: DefaultUserName,
			Attributes: &feature_annotation.FeatureAnnotationAttributes{
				Name: params.processedData.GeneID, // Use GeneID as name if not otherwise available
				Properties: collection.Map(
					params.processedData.StrippedPropsText,
					strippedPropertyToTagProperty,
				),
			},
		})
	if createErr != nil {
		result.Message = fmt.Sprintf(
			"failed to CreateFeatureAnnotation after not found: %v",
			createErr,
		)
		result.Error = createErr
		return result, createErr
	}
	result.Success = true
	result.Message = fmt.Sprintf(
		"Successfully created FeatureAnnotation with %d tags",
		len(params.processedData.StrippedPropsText),
	)
	return result, nil
}

func handleUpdateExistingFeatureAnnotation(
	params *handleUpdateExistingFeatureAnnotationParams,
) (GrpcUpdateResult, error) {
	result := GrpcUpdateResult{
		GeneID:  params.processedData.GeneID,
		Success: false,
	}
	for _, prop := range params.processedData.StrippedPropsText {
		err := processSinglePropertyUpdate(
			&processSinglePropertyUpdateParams{
				ctx:        params.ctx,
				grpcClient: params.grpcClient,
				featAnno:   params.featAnno,
				prop:       prop,
				arangoUser: params.config.ArangoUser,
				logger:     params.logger,
				jobID:      params.jobID,
			},
		)
		if err != nil {
			result.Message = err.Error()
			result.Error = err
			return result, err // Return on the first error encountered
		}
	}
	result.Success = true
	result.Message = fmt.Sprintf(
		"Successfully processed %d tags",
		len(params.processedData.StrippedPropsText),
	)
	return result, nil
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
			if status.Code(err) == codes.NotFound {
				return handleCreateFeatureAnnotation(
					&handleCreateFeatureAnnotationParams{
						ctx:           ctx,
						grpcClient:    grpcClient,
						processedData: processedData,
						jobID:         job.ID,
						logger:        logger,
					},
				)
			}
			result.Message = fmt.Sprintf(
				"failed to GetFeatureAnnotation: %v",
				err,
			)
			result.Error = err
			return result, err
		}
		return handleUpdateExistingFeatureAnnotation(
			&handleUpdateExistingFeatureAnnotationParams{
				ctx:           ctx,
				grpcClient:    grpcClient,
				featAnno:      featAnno,
				processedData: processedData,
				config:        config,
				jobID:         job.ID,
				logger:        logger,
			},
		)
	}
}
