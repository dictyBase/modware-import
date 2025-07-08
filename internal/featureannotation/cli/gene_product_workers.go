package cli

import (
	"context"
	"fmt"

	fanno "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// legacyDBQueryWorkerFunc creates worker for querying legacy database
func legacyDBQueryWorkerFunc(
	config GeneProductAppConfig,
) concurrent.WorkerFunc[GeneInfo, ProcessedGeneProduct] {
	return func(
		ctx context.Context,
		job concurrent.Job[GeneInfo],
	) (ProcessedGeneProduct, error) {
		gene := job.Payload
		logger := config.Logger

		// Get legacy database connection
		legacyDB := registry.GetLegacyArangodbConnection()
		cursor, err := legacyDB.SearchRows(
			GeneProductQuery,
			map[string]interface{}{
				"feature_id": gene.FeatureID,
			},
		)
		if err != nil {
			return ProcessedGeneProduct{}, fmt.Errorf(
				"failed to query gene product for gene %s: %w",
				gene.GeneID,
				err,
			)
		}
		defer cursor.Close()

		// Check if we have results
		if cursor.IsEmpty() {
			logger.Debugf("No gene product found for gene %s", gene.GeneID)
			return ProcessedGeneProduct{
				GeneID:   gene.GeneID,
				GeneName: gene.Name,
			}, nil
		}

		result := GeneProductResult{}
		if cursor.Scan() {
			if err := cursor.Read(&result); err != nil {
				return ProcessedGeneProduct{}, fmt.Errorf(
					"failed to read gene product result: %w",
					err,
				)
			}
		}

		logger.Debugf(
			"Found gene product '%s' for gene %s",
			result.GeneProduct,
			gene.GeneID,
		)

		return ProcessedGeneProduct{
			GeneID:      gene.GeneID,
			GeneName:    gene.Name,
			GeneProduct: result.GeneProduct,
			CreatedBy:   result.CreatedBy,
			CreatedOn:   result.CreatedOn.Time,
		}, nil
	}
}

func createFeatureAnnotationWithProduct(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	processedGene ProcessedGeneProduct,
) (GrpcUpdateResult, error) {
	result := GrpcUpdateResult{
		GeneID:  processedGene.GeneID,
		Success: false,
	}
	_, createErr := grpcClient.CreateFeatureAnnotation(
		ctx,
		&fanno.NewFeatureAnnotation{
			Id:        processedGene.GeneID,
			CreatedBy: DefaultUserName,
			Attributes: &fanno.FeatureAnnotationAttributes{
				Name: processedGene.GeneName,
				Properties: []*fanno.TagProperty{
					{
						Tag:       GeneProductTag,
						Value:     processedGene.GeneProduct,
						CreatedBy: processedGene.CreatedBy,
						CreatedAt: timestamppb.New(processedGene.CreatedOn),
					},
				},
			},
		},
	)
	if createErr != nil {
		result.Error = createErr
		result.Message = fmt.Sprintf(
			"Failed to create feature annotation: %v",
			createErr,
		)
		return result, createErr
	}
	result.Success = true
	result.Message = "Created new feature annotation with gene product"
	return result, nil
}

func addGeneProductTag(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	processedGene ProcessedGeneProduct,
) (GrpcUpdateResult, error) {
	result := GrpcUpdateResult{
		GeneID:  processedGene.GeneID,
		Success: false,
	}
	_, err := grpcClient.AddTag(
		ctx,
		&fanno.AddTagRequest{
			Id: processedGene.GeneID,
			Tag: &fanno.TagPropertyCreate{
				Tag:       GeneProductTag,
				Value:     processedGene.GeneProduct,
				CreatedBy: processedGene.CreatedBy,
			},
		},
	)
	if err != nil {
		result.Error = err
		result.Message = fmt.Sprintf(
			"Failed to add gene product tag: %v",
			err,
		)
		return result, err
	}

	result.Success = true
	result.Message = "Successfully added gene product tag"
	return result, nil
}

// geneProductGrpcWorkerFunc creates worker for updating via gRPC
func geneProductGrpcWorkerFunc(
	config GeneProductAppConfig,
	grpcClient fanno.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[ProcessedGeneProduct, GrpcUpdateResult] {
	return func(
		ctx context.Context,
		job concurrent.Job[ProcessedGeneProduct],
	) (GrpcUpdateResult, error) {
		processedGene := job.Payload
		logger := config.Logger

		result := GrpcUpdateResult{
			GeneID:  processedGene.GeneID,
			Success: false,
		}

		// Skip if no gene product
		if processedGene.GeneProduct == "" {
			result.Success = true
			result.Message = "No gene product to update"
			return result, nil
		}

		// Get existing feature annotation
		featAnno, err := grpcClient.GetFeatureAnnotation(
			ctx,
			&fanno.FeatureAnnotationId{
				Id: processedGene.GeneID,
			},
		)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// Create new feature annotation
				return createFeatureAnnotationWithProduct(
					ctx,
					grpcClient,
					processedGene,
				)
			}
			result.Error = err
			result.Message = fmt.Sprintf(
				"Failed to get feature annotation: %v",
				err,
			)
			return result, err
		}

		// Check if gene product tag already exists
		for _, prop := range featAnno.Attributes.Properties {
			if prop.Tag == GeneProductTag {
				logger.Debugf(
					"Gene product tag already exists for gene %s, skipping",
					processedGene.GeneID,
				)
				result.Success = true
				result.Message = "Gene product tag already exists"
				return result, nil
			}
		}

		// Add gene product tag
		return addGeneProductTag(ctx, grpcClient, processedGene)
	}
}
