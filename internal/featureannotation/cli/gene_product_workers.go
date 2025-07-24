package cli

import (
	"context"
	"fmt"

	fanno "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// legacyDBQueryWorkerFunc creates worker for querying legacy database
func legacyDBQueryWorkerFunc(
	config GeneProductAppConfig,
) concurrent.WorkerFunc[GeneInfo, []ProcessedGeneProduct] {
	return func(
		ctx context.Context,
		job concurrent.Job[GeneInfo],
	) ([]ProcessedGeneProduct, error) {
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
			return nil, fmt.Errorf(
				"failed to query gene product for gene %s: %w",
				gene.GeneID,
				err,
			)
		}
		defer cursor.Close()

		// Check if we have results
		if cursor.IsEmpty() {
			logger.Debugf("No gene product found for gene %s", gene.GeneID)
			return []ProcessedGeneProduct{{
				GeneID:   gene.GeneID,
				GeneName: gene.Name,
			}}, nil
		}

		var processedGeneProducts []ProcessedGeneProduct

		// Loop through all cursor results to get all gene products for
		// this gene
		for cursor.Scan() {
			result := GeneProductResult{}
			if err := cursor.Read(&result); err != nil {
				return nil, fmt.Errorf(
					"failed to read gene product result: %w",
					err,
				)
			}

			processedGeneProducts = append(
				processedGeneProducts,
				ProcessedGeneProduct{
					GeneID:      gene.GeneID,
					GeneName:    gene.Name,
					GeneProduct: result.GeneProduct,
					CreatedBy:   result.CreatedBy,
					CreatedOn:   result.CreatedOn.Time,
				},
			)

			logger.Debugf(
				"Found gene product '%s' for gene %s",
				result.GeneProduct,
				gene.GeneID,
			)
		}
		logger.Debugf(
			"Found %d gene products for gene %s",
			len(processedGeneProducts),
			gene.GeneID,
		)

		return processedGeneProducts, nil
	}
}

func handleFeatureAnnotationLookup(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	processedGene ProcessedGeneProduct,
	grpcErr error,
) (GrpcUpdateResult, error) {
	result := GrpcUpdateResult{
		GeneID:  processedGene.GeneID,
		Success: false,
	}
	if status.Code(grpcErr) != codes.NotFound {
		result.Error = grpcErr
		result.Message = fmt.Sprintf(
			"Failed to get feature annotation: %v",
			grpcErr,
		)
		return result, grpcErr
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

// handleNoGeneProduct processes cases where gene product is empty
func handleNoGeneProduct(
	processedGene ProcessedGeneProduct,
) (bool, GrpcUpdateResult) {
	if processedGene.GeneProduct != "" {
		return false, GrpcUpdateResult{}
	}
	return true, GrpcUpdateResult{
		GeneID:  processedGene.GeneID,
		Success: true,
		Message: "No gene product to update",
	}
}

// handleExistingGeneProduct checks if gene product tag already exists using slices.ContainsFunc
func handleExistingGeneProduct(
	featAnno *fanno.FeatureAnnotation,
	processedGene ProcessedGeneProduct,
	logger *logrus.Entry,
) (bool, GrpcUpdateResult) {
	hasGeneProduct := slices.ContainsFunc(
		featAnno.Attributes.Properties,
		func(prop *fanno.TagProperty) bool {
			return prop.Tag == GeneProductTag
		},
	)
	if !hasGeneProduct {
		return false, GrpcUpdateResult{}
	}
	logger.Debugf(
		"Gene product tag already exists for gene %s, skipping",
		processedGene.GeneID,
	)
	return true, GrpcUpdateResult{
		GeneID:  processedGene.GeneID,
		Success: true,
		Message: "Gene product tag already exists",
	}
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

// batchGeneProductGrpcWorkerFunc creates worker for batch processing gene products
func batchGeneProductGrpcWorkerFunc(
	config GeneProductAppConfig,
	grpcClient fanno.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[BatchGeneProductJob, BatchGeneProductResult] {
	return func(
		ctx context.Context,
		job concurrent.Job[BatchGeneProductJob],
	) (BatchGeneProductResult, error) {
		batchJob := job.Payload
		logger := config.Logger
		geneProducts := batchJob.GeneProducts

		result := BatchGeneProductResult{
			Success:        false,
			ProcessedCount: 0,
			SkippedCount:   0,
		}

		// Validate and filter products
		validProducts, skippedCount := filterValidGeneProducts(geneProducts, logger)
		result.SkippedCount += skippedCount

		if len(geneProducts) == 0 || len(validProducts) == 0 {
			result.Success = true
			result.Message = "No valid gene products to process"
			return result, nil
		}

		// All gene products should have the same GeneID
		geneID := geneProducts[0].GeneID
		result.GeneID = geneID

		// Get existing feature annotation
		featAnno, err := grpcClient.GetFeatureAnnotation(
			ctx,
			&fanno.FeatureAnnotationId{Id: geneID},
		)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// Create new feature annotation with all gene products
				return createFeatureAnnotationWithMultipleProducts(ctx, grpcClient, validProducts)
			}
			result.Error = err
			result.Message = fmt.Sprintf("Failed to get feature annotation: %v", err)
			return result, err
		}

		// Process existing annotation
		return processBatchUpdateForExistingAnnotation(
			ctx,
			grpcClient,
			featAnno,
			validProducts,
			result,
			logger,
		)
	}
}

// filterValidGeneProducts filters out empty gene products and returns valid ones with skip count
func filterValidGeneProducts(
	geneProducts []ProcessedGeneProduct,
	logger *logrus.Entry,
) ([]ProcessedGeneProduct, int) {
	var validProducts []ProcessedGeneProduct
	skippedCount := 0

	for _, product := range geneProducts {
		if product.GeneProduct != "" {
			validProducts = append(validProducts, product)
		} else {
			skippedCount++
			logger.Debugf("Skipping empty gene product for gene %s", product.GeneID)
		}
	}

	return validProducts, skippedCount
}

// processBatchUpdateForExistingAnnotation handles batch update for existing annotations
func processBatchUpdateForExistingAnnotation(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	featAnno *fanno.FeatureAnnotation,
	validProducts []ProcessedGeneProduct,
	result BatchGeneProductResult,
	logger *logrus.Entry,
) (BatchGeneProductResult, error) {
	geneID := validProducts[0].GeneID

	// Collect existing gene product values for comparison
	existingValues := make(map[string]bool)
	for _, prop := range featAnno.Attributes.Properties {
		if prop.Tag == GeneProductTag {
			existingValues[prop.Value] = true
		}
	}

	// Filter new products that don't already exist
	var newProducts []ProcessedGeneProduct
	for _, product := range validProducts {
		if !existingValues[product.GeneProduct] {
			newProducts = append(newProducts, product)
		} else {
			result.SkippedCount++
			logger.Debugf(
				"Gene product '%s' already exists for gene %s, skipping",
				product.GeneProduct,
				geneID,
			)
		}
	}

	if len(newProducts) == 0 {
		result.Success = true
		result.Message = "All gene products already exist"
		return result, nil
	}

	// Add new gene products using batch update
	return updateFeatureAnnotationWithMultipleProducts(ctx, grpcClient, featAnno, newProducts)
}

// createFeatureAnnotationWithMultipleProducts creates new feature annotation with multiple gene products
func createFeatureAnnotationWithMultipleProducts(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	geneProducts []ProcessedGeneProduct,
) (BatchGeneProductResult, error) {
	if len(geneProducts) == 0 {
		return BatchGeneProductResult{
			Success: true,
			Message: "No gene products to create",
		}, nil
	}

	geneID := geneProducts[0].GeneID
	result := BatchGeneProductResult{
		GeneID: geneID,
	}

	// Create properties for all gene products
	var properties []*fanno.TagProperty
	for _, product := range geneProducts {
		properties = append(properties, &fanno.TagProperty{
			Tag:       GeneProductTag,
			Value:     product.GeneProduct,
			CreatedBy: product.CreatedBy,
			CreatedAt: timestamppb.New(product.CreatedOn),
		})
	}

	// Create new feature annotation
	req := &fanno.NewFeatureAnnotation{
		Id:        geneID,
		CreatedBy: DefaultUserName,
		Attributes: &fanno.FeatureAnnotationAttributes{
			Properties: properties,
		},
	}

	_, err := grpcClient.CreateFeatureAnnotation(ctx, req)
	if err != nil {
		result.Error = err
		result.Message = fmt.Sprintf("Failed to create feature annotation: %v", err)
		return result, err
	}

	result.Success = true
	result.ProcessedCount = len(geneProducts)
	result.Message = fmt.Sprintf("Successfully created feature annotation with %d gene products", len(geneProducts))
	return result, nil
}

// updateFeatureAnnotationWithMultipleProducts updates existing feature annotation with multiple gene products
func updateFeatureAnnotationWithMultipleProducts(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	existingAnnotation *fanno.FeatureAnnotation,
	newGeneProducts []ProcessedGeneProduct,
) (BatchGeneProductResult, error) {
	if len(newGeneProducts) == 0 {
		return BatchGeneProductResult{
			Success: true,
			Message: "No new gene products to add",
		}, nil
	}

	geneID := newGeneProducts[0].GeneID
	result := BatchGeneProductResult{
		GeneID: geneID,
	}

	// Create properties for new gene products
	var newProperties []*fanno.TagProperty
	for _, product := range newGeneProducts {
		newProperties = append(newProperties, &fanno.TagProperty{
			Tag:       GeneProductTag,
			Value:     product.GeneProduct,
			CreatedBy: product.CreatedBy,
			CreatedAt: timestamppb.New(product.CreatedOn),
		})
	}

	// Combine existing properties with new ones
	allProperties := make([]*fanno.TagProperty, 0, len(existingAnnotation.Attributes.Properties)+len(newProperties))
	allProperties = append(allProperties, existingAnnotation.Attributes.Properties...)
	allProperties = append(allProperties, newProperties...)

	// Update feature annotation
	req := &fanno.FeatureAnnotationUpdate{
		Id: geneID,
		Attributes: &fanno.FeatureAnnotationAttributes{
			Properties: allProperties,
		},
		UpdatedBy: DefaultUserName,
	}

	_, err := grpcClient.UpdateFeatureAnnotation(ctx, req)
	if err != nil {
		result.Error = err
		result.Message = fmt.Sprintf("Failed to update feature annotation: %v", err)
		return result, err
	}

	result.Success = true
	result.ProcessedCount = len(newGeneProducts)
	result.Message = fmt.Sprintf("Successfully added %d gene products", len(newGeneProducts))
	return result, nil
}
