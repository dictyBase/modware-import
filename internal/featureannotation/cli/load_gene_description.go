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
	// descriptionTag is the tag name used for gene descriptions in the feature annotation system
	descriptionTag = "gene_description"
)

// GeneDescription represents a gene with its description for loading operations.
type GeneDescription struct {
	GeneID      string `validate:"required,min=1"`
	Description string `validate:"required,min=1"`
}

// manageGeneDescriptionParams holds parameters for managing gene descriptions.
type manageGeneDescriptionParams struct {
	ctx         context.Context
	client      feature.FeatureAnnotationServiceClient
	description GeneDescription `validate:"required"`
	user        string          `validate:"required,email"`
	logger      *logrus.Entry   `validate:"required"`
}

// handleNewGeneDescFromCsvParams holds parameters for creating new gene descriptions.
type handleNewGeneDescFromCsvParams struct {
	ctx         context.Context
	client      feature.FeatureAnnotationServiceClient
	description GeneDescription `validate:"required"`
	user        string          `validate:"required,email"`
	grpcErr     error
	logger      *logrus.Entry `validate:"required"`
}

// LoadGeneDescription loads gene descriptions from a CSV file into the feature annotation service.
// It validates input parameters, processes the CSV file concurrently, and reports results.
func LoadGeneDescription(c *cli.Context) error {
	// Validate CLI parameters
	params := LoadGeneDescriptionParams{
		InputFile: c.String("input"),
		Workers:   c.Int("workers"),
		BatchSize: c.Int("batch-size"),
		User:      c.String("user"),
	}

	if err := ValidateStruct(params); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	logger := registry.GetLogger()
	client := registry.GetFeatureAnnotationAPIClient()
	user := params.User

	workerFunc := func(
		ctx context.Context,
		job concurrent.Job[GeneDescription],
	) (*pb.FeatureAnnotation, error) {
		return manageGeneDescription(&manageGeneDescriptionParams{
			ctx:         ctx,
			client:      client,
			description: job.Payload,
			user:        user,
			logger:      logger,
		})
	}

	processor := concurrent.NewBatchProcessor(
		workerFunc,
		params.BatchSize,
		concurrent.WithWorkers[GeneDescription, *pb.FeatureAnnotation](
			params.Workers,
		),
	)

	processor.Start()

	// Stream CSV records directly to processor in a separate goroutine
	go func() {
		defer processor.Close()
		if err := streamGeneDescriptionsFromCSV(params.InputFile, processor, logger); err != nil {
			logger.Errorf("error streaming CSV file: %s", err)
		}
	}()

	successCount, errorCount := processGeneDescriptionResults(processor, logger)

	logger.Infof(
		"Finished loading gene descriptions. Success: %d, Errors: %d",
		successCount,
		errorCount,
	)
	if errorCount > 0 {
		return fmt.Errorf("encountered %d errors during loading", errorCount)
	}

	return nil
}

// streamGeneDescriptionsFromCSV parses a CSV file and streams gene descriptions
// directly to the processor. It expects a CSV format with gene ID in the first
// column and description in the second column. Duplicate gene IDs are detected
// and skipped with a warning. Records are processed as they are read, providing
// constant memory usage regardless of file size.
func streamGeneDescriptionsFromCSV(
	file string,
	processor *concurrent.BatchProcessor[GeneDescription, *pb.FeatureAnnotation],
	logger *logrus.Entry,
) error {
	reader, cleanup, err := openCSVFileWithHeader(file, logger)
	if err != nil {
		return err
	}
	defer cleanup()

	processGeneDescriptionRecords(reader, processor, logger)
	return nil
}

// openCSVFileWithHeader opens a CSV file and validates the header
func openCSVFileWithHeader(file string, logger *logrus.Entry) (*csv.Reader, func(), error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening input file %s: %w", file, err)
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

// processGeneDescriptionRecords handles the main CSV processing loop
func processGeneDescriptionRecords(
	reader *csv.Reader,
	processor *concurrent.BatchProcessor[GeneDescription, *pb.FeatureAnnotation],
	logger *logrus.Entry,
) {
	seen := make(map[string]struct{})
	lineNumber := 2 // Start at 2 since we skipped header
	processedCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Errorf("error reading record from csv at line %d: %s", lineNumber, err)
			lineNumber++
			continue
		}

		if processRecord(record, seen, processor, logger, lineNumber) {
			processedCount++
		}
		lineNumber++

		reportProgressIfNeeded(processedCount, logger)
	}

	logger.Infof("completed streaming %d gene descriptions", processedCount)
}

// processRecord validates and processes a single CSV record
func processRecord(
	record []string,
	seen map[string]struct{},
	processor *concurrent.BatchProcessor[GeneDescription, *pb.FeatureAnnotation],
	logger *logrus.Entry,
	lineNumber int,
) bool {
	if !isValidRecord(record, logger, lineNumber) {
		return false
	}

	if isDuplicateGene(record[0], seen, logger, lineNumber) {
		return false
	}

	geneDesc := GeneDescription{GeneID: record[0], Description: record[1]}
	if err := ValidateStruct(geneDesc); err != nil {
		logger.Warnf("skipping invalid gene description at line %d: %s", lineNumber, err)
		return false
	}

	processor.Add(geneDesc)
	seen[record[0]] = struct{}{}
	return true
}

// isValidRecord checks if the CSV record has the required format
func isValidRecord(record []string, logger *logrus.Entry, lineNumber int) bool {
	if len(record) < 2 {
		logger.Warnf("skipping malformed record at line %d: %v", lineNumber, record)
		return false
	}
	return true
}

// isDuplicateGene checks if the gene ID has already been processed
func isDuplicateGene(geneID string, seen map[string]struct{}, logger *logrus.Entry, lineNumber int) bool {
	if _, exists := seen[geneID]; exists {
		logger.Warnf("duplicate gene ID %s found at line %d, skipping", geneID, lineNumber)
		return true
	}
	return false
}

// reportProgressIfNeeded logs progress every 1000 records
func reportProgressIfNeeded(processedCount int, logger *logrus.Entry) {
	if processedCount%1000 == 0 {
		logger.Infof("processed %d gene descriptions from CSV", processedCount)
	}
}

// manageGeneDescription handles the creation or update of gene descriptions.
// It checks if the description already exists and either skips or adds the new
// description.
func manageGeneDescription(
	params *manageGeneDescriptionParams,
) (*pb.FeatureAnnotation, error) {
	existing, err := params.client.GetFeatureAnnotation(
		params.ctx,
		&pb.FeatureAnnotationId{Id: params.description.GeneID},
	)
	if err != nil {
		return handleNewGeneDescFromCsv(
			&handleNewGeneDescFromCsvParams{
				ctx:         params.ctx,
				client:      params.client,
				description: params.description,
				user:        params.user,
				grpcErr:     err,
				logger:      params.logger,
			},
		)
	}

	ok := slices.ContainsFunc(
		existing.Attributes.Properties,
		func(tag *pb.TagProperty) bool {
			return tag.Tag == descriptionTag &&
				tag.Value == params.description.Description
		},
	)
	if ok {
		params.logger.Debugf(
			"gene description %s already exists for %s",
			params.description.Description,
			params.description.GeneID,
		)
		return existing, nil
	}

	updated, err := params.client.AddTag(params.ctx, &pb.AddTagRequest{
		Id: params.description.GeneID,
		Tag: &pb.TagPropertyCreate{
			Tag:       descriptionTag,
			Value:     params.description.Description,
			CreatedBy: params.user,
			CreatedAt: timestamppb.Now(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf(
			"error creating tags for %s: %w",
			params.description.GeneID,
			err,
		)
	}
	params.logger.Debugf(
		"successfully updated gene description for %s",
		updated.Id,
	)
	return updated, nil
}

// handleNewGeneDescFromCsv creates a new feature annotation when the gene
// doesn't exist. It handles the NotFound error case from the initial lookup.
func handleNewGeneDescFromCsv(
	params *handleNewGeneDescFromCsvParams,
) (*pb.FeatureAnnotation, error) {
	if status.Code(params.grpcErr) != codes.NotFound {
		return nil, fmt.Errorf(
			"error finding feature annotation for %s: %w",
			params.description.GeneID,
			params.grpcErr,
		)
	}
	nfa := &pb.NewFeatureAnnotation{
		Id:        params.description.GeneID,
		CreatedBy: params.user,
		CreatedAt: timestamppb.Now(),
		Attributes: &pb.FeatureAnnotationAttributes{
			Name: params.description.GeneID,
			Properties: []*pb.TagProperty{{
				Tag:       descriptionTag,
				Value:     params.description.Description,
				CreatedBy: params.user,
				CreatedAt: timestamppb.Now(),
			}},
		},
	}
	created, err := params.client.CreateFeatureAnnotation(params.ctx, nfa)
	if err != nil {
		return nil, fmt.Errorf(
			"error creating annotation for %s: %w",
			params.description.GeneID,
			err,
		)
	}
	params.logger.Debugf("created new feature annotation %s", created.Id)
	return created, nil
}

// processGeneDescriptionResults processes the results from the batch processor.
// It counts successes and errors, logging detailed information for each result.
func processGeneDescriptionResults(
	processor *concurrent.BatchProcessor[GeneDescription, *pb.FeatureAnnotation],
	logger *logrus.Entry,
) (int, int) {
	var successCount, errorCount int
	for result := range processor.Results() {
		if result.Error != nil {
			logger.WithFields(logrus.Fields{
				"job-id": result.JobID,
				"error":  result.Error,
			}).Error("error loading gene description")
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
