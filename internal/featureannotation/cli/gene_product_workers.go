package cli

import (
	"context"
	"fmt"
	"slices"

	fanno "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/collection"
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
		// i) No gene product
		if ok, result := handleNoGeneProduct(processedGene); ok {
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
			return handleFeatureAnnotationLookup(
				ctx,
				grpcClient,
				processedGene,
				err,
			)
		}
		// iii) Gene product already exists
		if ok, result := handleExistingGeneProduct(
			featAnno,
			processedGene,
			config.Logger,
		); ok {
			return result, nil
		}

		// iv) then goes to addgeneProductTag
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

		// Validate and filter products
		validProducts := collection.Filter(
			geneProducts,
			func(gp ProcessedGeneProduct) bool {
				return gp.GeneProduct != ""
			})
		result := BatchGeneProductResult{
			GeneID:         validProducts[0].GeneID,
			Success:        false,
			ProcessedCount: 0,
			SkippedCount:   0,
		}
		if len(validProducts) == 0 {
			result.Success = true
			result.Message = "No valid gene products to process"
			return result, nil
		}
		geneID := validProducts[0].GeneID

		// Get existing feature annotation
		featAnno, err := grpcClient.GetFeatureAnnotation(
			ctx,
			&fanno.FeatureAnnotationId{Id: geneID},
		)
		if err != nil {
			return handleNewFeatAnnoWithMulti(
				ctx,
				grpcClient,
				validProducts,
				err,
			)
		}

		// Process existing annotation
		hasNew, newProducts := handleExistingFeatAnnoWithMulti(
			featAnno,
			validProducts,
		)

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

func handleExistingFeatAnnoWithMulti(
	featAnno *fanno.FeatureAnnotation,
	validProducts []ProcessedGeneProduct,
) (bool, []ProcessedGeneProduct) {
	existingSet := mapset.NewSet(
		collection.Pipe2(
			featAnno.Attributes.Properties,
			collection.CurriedFilter(
				func(prop *fanno.TagProperty) bool {
					return prop.Tag == GeneProductTag
				}),
			collection.CurriedMap(
				func(prop *fanno.TagProperty) string {
					return prop.Value
				}),
		)...)

	// Filter new products that don't already exist
	newProducts := collection.Filter(
		validProducts,
		func(product ProcessedGeneProduct) bool {
			return !existingSet.Contains(product.GeneProduct)
		})

	return len(newProducts) > 0, newProducts
}

// createFeatureAnnotationWithMultipleProducts creates new feature annotation with multiple gene products
func handleNewFeatAnnoWithMulti(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	geneProducts []ProcessedGeneProduct,
	grpcErr error,
) (BatchGeneProductResult, error) {
	geneID := geneProducts[0].GeneID
	result := BatchGeneProductResult{
		GeneID: geneID,
	}
	if status.Code(grpcErr) != codes.NotFound {
		result.Error = grpcErr
		result.Message = fmt.Sprintf(
			"Failed to get feature annotation: %v",
			grpcErr,
		)
		return result, grpcErr
	}
	properties := collection.Map(
		geneProducts,
		func(product ProcessedGeneProduct) *fanno.TagProperty {
			return &fanno.TagProperty{
				Tag:       GeneProductTag,
				Value:     product.GeneProduct,
				CreatedBy: product.CreatedBy,
				CreatedAt: timestamppb.New(product.CreatedOn),
			}
		},
	)
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
		result.Message = fmt.Sprintf(
			"Failed to create feature annotation: %v",
			err,
		)
		return result, err
	}

	result.Success = true
	result.ProcessedCount = len(geneProducts)
	result.Message = fmt.Sprintf(
		"Successfully created feature annotation with %d gene products",
		len(geneProducts),
	)
	return result, nil
}

// updateFeatureAnnotationWithMultipleProducts updates existing feature annotation with multiple gene products
func handleUpdateFeatAnnoWithMulti(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	existingAnnotation *fanno.FeatureAnnotation,
	newGeneProducts []ProcessedGeneProduct,
	count int,
) (BatchGeneProductResult, error) {
	geneID := newGeneProducts[0].GeneID
	result := BatchGeneProductResult{
		GeneID:       geneID,
		SkippedCount: count,
	}
	// Create properties for new gene products
	newProperties := collection.Map(
		newGeneProducts,
		func(product ProcessedGeneProduct) *fanno.TagPropertyCreate {
			return &fanno.TagPropertyCreate{
				Tag:       GeneProductTag,
				Value:     product.GeneProduct,
				CreatedBy: product.CreatedBy,
				CreatedAt: timestamppb.New(product.CreatedOn),
			}
		})
	_, err := grpcClient.AddTags(ctx, &fanno.AddTagsRequest{
		Id:   existingAnnotation.Id,
		Tags: newProperties,
	})
	if err != nil {
		result.Error = err
		result.Message = fmt.Sprintf(
			"Failed to add gene products to feature annotation: %v",
			err,
		)
		return result, err
	}

	result.Success = true
	result.ProcessedCount = len(newGeneProducts)
	result.Message = fmt.Sprintf(
		"Successfully added %d gene products",
		len(newGeneProducts),
	)
	return result, nil
}
