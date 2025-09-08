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

	processor.Start()

	// Stream CSV records directly to processor in a separate goroutine
	go func() {
		defer processor.Close()
		streamGeneProductsFromCSVFiles(
			c.StringSlice("input"),
			processor,
			logger,
		)
	}()

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

// streamGeneProductsFromCSVFiles processes multiple CSV files and streams gene products directly to the processor.
// It handles multiple input files sequentially, maintaining duplicate detection across all files.
func streamGeneProductsFromCSVFiles(
	files []string,
	processor *concurrent.BatchProcessor[GeneProduct, *pb.FeatureAnnotation],
	logger *logrus.Entry,
) {
	seen := make(map[string]struct{})
	totalProcessed := 0

	for fileIndex, file := range files {
		fileProcessed := processGeneProductFile(
			file,
			fileIndex,
			len(files),
			seen,
			processor,
			logger,
			&totalProcessed,
		)

		logger.Infof(
			"completed file %s: processed %d gene products",
			file,
			fileProcessed,
		)
	}

	logger.Infof(
		"completed streaming %d gene products from %d files",
		totalProcessed,
		len(files),
	)
}

// processGeneProductFile handles processing of a single CSV file
func processGeneProductFile(
	file string,
	fileIndex, totalFiles int,
	seen map[string]struct{},
	processor *concurrent.BatchProcessor[GeneProduct, *pb.FeatureAnnotation],
	logger *logrus.Entry,
	totalProcessed *int,
) int {
	logger.Infof("processing file %d of %d: %s", fileIndex+1, totalFiles, file)

	reader, cleanup, err := openGeneProductCSV(file, logger)
	if err != nil {
		logger.Errorf("Error opening file %s: %v", file, err)
		return 0
	}
	if reader == nil {
		return 0 // Empty file
	}
	defer cleanup()

	return processGeneProductRecords(
		reader,
		file,
		seen,
		processor,
		logger,
		totalProcessed,
	)
}

// openGeneProductCSV opens and validates a CSV file for gene products
func openGeneProductCSV(
	file string,
	logger *logrus.Entry,
) (*csv.Reader, func(), error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"error opening input file %s: %w",
			file,
			err,
		)
	}

	reader := csv.NewReader(f)
	if _, err := reader.Read(); err != nil {
		f.Close()
		if err == io.EOF {
			logger.Warnf("empty csv file %s", file)
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("error reading header from csv: %w", err)
	}

	cleanup := func() { f.Close() }
	return reader, cleanup, nil
}

// processGeneProductRecords processes all records in a CSV file
func processGeneProductRecords(
	reader *csv.Reader,
	filename string,
	seen map[string]struct{},
	processor *concurrent.BatchProcessor[GeneProduct, *pb.FeatureAnnotation],
	logger *logrus.Entry,
	totalProcessed *int,
) int {
	lineNumber := 2 // Start at 2 since we skipped header
	fileProcessed := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf(
				"error reading record from csv at line %d in file %s: %s",
				lineNumber,
				filename,
				err,
			)
			lineNumber++
			continue
		}

		if processGeneProductRecord(
			record,
			filename,
			seen,
			processor,
			logger,
			lineNumber,
		) {
			fileProcessed++
			*totalProcessed++
			reportCSVProgress(*totalProcessed, logger)
		}
		lineNumber++
	}

	return fileProcessed
}

// processGeneProductRecord validates and processes a single gene product record
func processGeneProductRecord(
	record []string,
	filename string,
	seen map[string]struct{},
	processor *concurrent.BatchProcessor[GeneProduct, *pb.FeatureAnnotation],
	logger *logrus.Entry,
	lineNumber int,
) bool {
	if !isValidGeneProductRecord(record, filename, logger, lineNumber) {
		return false
	}

	if isDuplicateGeneProduct(record[0], filename, seen, logger, lineNumber) {
		return false
	}

	product := GeneProduct{GeneID: record[0], Product: record[1]}
	processor.Add(product)
	seen[record[0]] = struct{}{}
	return true
}

// isValidGeneProductRecord checks if the CSV record has the required format
func isValidGeneProductRecord(
	record []string,
	filename string,
	logger *logrus.Entry,
	lineNumber int,
) bool {
	if len(record) < 2 {
		logger.Warnf(
			"skipping malformed record at line %d in file %s: %v",
			lineNumber,
			filename,
			record,
		)
		return false
	}
	return true
}

// isDuplicateGeneProduct checks if the gene ID has already been processed
func isDuplicateGeneProduct(
	geneID, filename string,
	seen map[string]struct{},
	logger *logrus.Entry,
	lineNumber int,
) bool {
	if _, exists := seen[geneID]; exists {
		logger.Warnf(
			"duplicate gene ID %s found at line %d in file %s, skipping",
			geneID,
			lineNumber,
			filename,
		)
		return true
	}
	return false
}

// reportCSVProgress logs progress every 1000 records
func reportCSVProgress(totalProcessed int, logger *logrus.Entry) {
	if totalProcessed%1000 == 0 {
		logger.Infof(
			"processed %d gene products across all files",
			totalProcessed,
		)
	}
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
			return tag.Tag == GeneProductTag &&
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
			Tag:       GeneProductTag,
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
				Tag:       GeneProductTag,
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
