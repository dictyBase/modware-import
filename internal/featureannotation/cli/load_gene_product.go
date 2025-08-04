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

func LoadGeneProduct(c *cli.Context) error {
	logger := registry.GetLogger()
	client := registry.GetFeatureAnnotationAPIClient()
	user := c.String("user")

	workerFunc := func(
		ctx context.Context,
		job concurrent.Job[GeneProduct],
	) (*pb.FeatureAnnotation, error) {
		return manageGeneProduct(ctx, client, job.Payload, user, logger)
	}

	processor := concurrent.NewBatchProcessor(
		workerFunc,
		c.Int("batch-size"),
		concurrent.WithWorkers[GeneProduct, *pb.FeatureAnnotation](
			c.Int("workers"),
		),
	)

	file, err := os.Open(c.String("input"))
	if err != nil {
		return fmt.Errorf(
			"error opening input file %s: %w",
			c.String("input"),
			err,
		)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return fmt.Errorf("empty csv file %s", c.String("input"))
		}
		return fmt.Errorf("error reading header from csv: %w", err)
	}

	processor.Start()
	go readCSVAndProcess(reader, processor, logger)
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

func manageGeneProduct(
	ctx context.Context,
	client feature.FeatureAnnotationServiceClient,
	product GeneProduct,
	user string,
	logger *logrus.Entry,
) (*pb.FeatureAnnotation, error) {
	existing, err := client.GetFeatureAnnotation(
		ctx,
		&pb.FeatureAnnotationId{Id: product.GeneID},
	)
	if err != nil {
		return handleNewGenePdtFromCsv(
			ctx,
			client,
			product,
			user,
			err,
			logger,
		)
	}

	ok := slices.ContainsFunc(
		existing.Attributes.Properties,
		func(tag *pb.TagProperty) bool {
			return tag.Tag == productTag &&
				tag.Value == product.Product
		},
	)
	if ok {
		logger.Debugf(
			"gene product %s already exists for %s",
			product.Product,
			product.GeneID,
		)
		return existing, nil
	}

	updated, err := client.AddTag(ctx, &pb.AddTagRequest{
		Id: product.GeneID,
		Tag: &pb.TagPropertyCreate{
			Tag:       productTag,
			Value:     product.Product,
			CreatedBy: user,
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error creating tags for %s: %w",
			product.GeneID,
			err,
		)
	}
	logger.Debugf("successfully updated gene product for %s", updated.Id)
	return updated, nil
}

func handleNewGenePdtFromCsv(
	ctx context.Context,
	client feature.FeatureAnnotationServiceClient,
	product GeneProduct,
	user string,
	grpcErr error,
	logger *logrus.Entry,
) (*pb.FeatureAnnotation, error) {
	if status.Code(grpcErr) != codes.NotFound {
		return nil, fmt.Errorf(
			"error finding feature annotation for %s: %w",
			product.GeneID,
			grpcErr,
		)
	}
	logger.Debugf(
		"creating new feature annotation for %s",
		product.GeneID,
	)
	nfa := &pb.NewFeatureAnnotation{
		Id:        product.GeneID,
		CreatedBy: user,
		CreatedAt: timestamppb.Now(),
		Attributes: &pb.FeatureAnnotationAttributes{
			Name: product.GeneID,
			Properties: []*pb.TagProperty{{
				Tag:       productTag,
				Value:     product.Product,
				CreatedBy: user,
				CreatedAt: timestamppb.Now(),
			}},
		},
	}
	created, err := client.CreateFeatureAnnotation(ctx, nfa)
	if err != nil {
		return nil, fmt.Errorf(
			"error creating annotation for %s: %w",
			product.GeneID,
			err,
		)
	}
	logger.Debugf("created new feature annotation %s", created.Id)
	return created, nil
}

func readCSVAndProcess(
	reader *csv.Reader,
	processor *concurrent.BatchProcessor[GeneProduct, *pb.FeatureAnnotation],
	logger *logrus.Entry,
) {
	defer processor.Close()
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
		processor.Add(GeneProduct{GeneID: record[0], Product: record[1]})
	}
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
