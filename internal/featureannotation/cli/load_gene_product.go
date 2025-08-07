package cli

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	productTag = "gene_product"
)

type GeneProduct struct {
	GeneID  string
	Product string
}

type manageGeneProductParams struct {
	ctx     context.Context
	client  feature.FeatureAnnotationServiceClient
	product GeneProduct
	user    string
	logger  *logrus.Entry
}

type handleNewGenePdtFromCsvParams struct {
	ctx     context.Context
	client  feature.FeatureAnnotationServiceClient
	product GeneProduct
	user    string
	grpcErr error
	logger  *logrus.Entry
}

func LoadGeneProduct(c *cli.Context) error {
	logger := registry.GetLogger()
	client := registry.GetFeatureAnnotationAPIClient()
	user := c.String("user")

	workerFunc := func(
		ctx context.Context,
		job concurrent.Job[GeneProduct],
	) (*pb.FeatureAnnotation, error) {
		return manageGeneProduct(&manageGeneProductParams{
			ctx:     ctx,
			client:  client,
			product: job.Payload,
			user:    user,
			logger:  logger,
		})
	}

	processor := concurrent.NewBatchProcessor(
		workerFunc,
		c.Int("batch-size"),
		concurrent.WithWorkers[GeneProduct, *pb.FeatureAnnotation](
			c.Int("workers"),
		),
	)

	products, err := readAndDeduplicate(c.StringSlice("input"), logger)
	if err != nil {
		return err
	}

	processor.Start()
	processor.AddBatch(products)
	processor.Close()

	successCount, errorCount := processResults(processor, logger)

	logger.Infof(
		"Finished loading gene products. Success: %d, Errors: %d",
		successCount,
		errorCount,
	)
	if errorCount > 0 {
		return fmt.Errorf("encountered %d errors during loading", errorCount)
	}

	return nil
}

func readAndDeduplicate(
	files []string,
	logger *logrus.Entry,
) ([]GeneProduct, error) {
	seen := make(map[string]bool)
	var products []GeneProduct
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("error opening input file %s: %w", file, err)
		}
		defer f.Close()
		reader := csv.NewReader(f)
		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				logger.Warnf("empty csv file %s", file)
				continue
			}
			return nil, fmt.Errorf("error reading header from csv: %w", err)
		}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Errorf("error reading record from csv: %s", err)
				continue
			}
			if len(record) < 2 {
				logger.Warnf("skipping malformed record: %v", record)
				continue
			}
			if seen[record[0]] {
				continue
			}
			products = append(
				products,
				GeneProduct{GeneID: record[0], Product: record[1]},
			)
			seen[record[0]] = true
		}
	}
	return products, nil
}

func manageGeneProduct(
	params *manageGeneProductParams,
) (*pb.FeatureAnnotation, error) {
	existing, err := params.client.GetFeatureAnnotation(
		params.ctx,
		&pb.FeatureAnnotationId{Id: params.product.GeneID},
	)
	if err != nil {
		return handleNewGenePdtFromCsv(
			&handleNewGenePdtFromCsvParams{
				ctx:     params.ctx,
				client:  params.client,
				product: params.product,
				user:    params.user,
				grpcErr: err,
				logger:  params.logger,
			},
		)
	}

	ok := slices.ContainsFunc(
		existing.Attributes.Properties,
		func(tag *pb.TagProperty) bool {
			return tag.Tag == productTag &&
				tag.Value == params.product.Product
		},
	)
	if ok {
		params.logger.Debugf(
			"gene product %s already exists for %s",
			params.product.Product,
			params.product.GeneID,
		)
		return existing, nil
	}

	updated, err := params.client.AddTag(params.ctx, &pb.AddTagRequest{
		Id: params.product.GeneID,
		Tag: &pb.TagPropertyCreate{
			Tag:       productTag,
			Value:     params.product.Product,
			CreatedBy: params.user,
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error creating tags for %s: %w",
			params.product.GeneID,
			err,
		)
	}
	params.logger.Debugf(
		"successfully updated gene product for %s",
		updated.Id,
	)
	return updated, nil
}

func handleNewGenePdtFromCsv(
	params *handleNewGenePdtFromCsvParams,
) (*pb.FeatureAnnotation, error) {
	if status.Code(params.grpcErr) != codes.NotFound {
		return nil, fmt.Errorf(
			"error finding feature annotation for %s: %w",
			params.product.GeneID,
			params.grpcErr,
		)
	}
	nfa := &pb.NewFeatureAnnotation{
		Id:        params.product.GeneID,
		CreatedBy: params.user,
		CreatedAt: timestamppb.Now(),
		Attributes: &pb.FeatureAnnotationAttributes{
			Name: params.product.GeneID,
			Properties: []*pb.TagProperty{{
				Tag:       productTag,
				Value:     params.product.Product,
				CreatedBy: params.user,
				CreatedAt: timestamppb.Now(),
			}},
		},
	}
	created, err := params.client.CreateFeatureAnnotation(params.ctx, nfa)
	if err != nil {
		return nil, fmt.Errorf(
			"error creating annotation for %s: %w",
			params.product.GeneID,
			err,
		)
	}
	params.logger.Debugf("created new feature annotation %s", created.Id)
	return created, nil
}

func processResults(
	processor *concurrent.BatchProcessor[GeneProduct, *pb.FeatureAnnotation],
	logger *logrus.Entry,
) (int, int) {
	var successCount, errorCount int
	for result := range processor.Results() {
		if result.Error != nil {
			logger.WithFields(logrus.Fields{
				"job-id": result.JobID,
				"error":  result.Error,
			}).Error("error loading gene product")
			errorCount++
		} else {
			successCount++
		}
	}
	for err := range processor.Errors() {
		logger.Errorf("processor error: %s", err)
		errorCount++
	}
	return successCount, errorCount
}
