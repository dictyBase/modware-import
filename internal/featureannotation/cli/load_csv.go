package cli

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/dictyBase/arangomanager"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Stage 1: Setup
type SetupConfig struct {
	Logger         *logrus.Entry
	DBH            *arangomanager.Database
	CSVFilePath    string
	CollectionName string
	BatchSize      int
	Delimiter      string
	Workers        int
}

// Stage 2: File Processing
type FileContext struct {
	Setup  SetupConfig
	File   *os.File
	Reader *csv.Reader
	Error  error // To propagate errors
}

// Stage 3: Header Validation
type ProcessingContext struct {
	FileContext
	FeaturePropIDIndex int
	ValueIndex         int
	// Error field is inherited from FileContext
}

// Stage 4: Data Processing / Final Pipeline Result
type PipelineResult struct {
	File        *os.File
	Setup       SetupConfig
	UpdateCount int
	Error       error
}

// ProcessSingleRecordParams holds the parameters for the processSingleRecordAndValidate function.
type ProcessSingleRecordParams struct {
	Record             []string
	FeaturePropIDIndex int
	ValueIndex         int
	RowNumForLogging   int // Actual row number in the CSV file for logging
	Logger             *logrus.Entry
}

// SubmitBatchAndLogParams holds the parameters for the submitBatchAndLog function.
type SubmitBatchAndLogParams struct {
	Setup            *SetupConfig
	Docs             []map[string]interface{}
	Logger           *logrus.Entry
	BatchDescription string
}

// Stage 1: Initialize setup
func setupPipeline(cltx *cli.Context) SetupConfig {
	return SetupConfig{
		Logger:         registry.GetLogger(),
		DBH:            registry.GetArangodbConnection(),
		CSVFilePath:    registry.GetCSVFilePath(),
		CollectionName: cltx.String("collection"),
		BatchSize:      cltx.Int("batch-size"),
		Delimiter:      cltx.String("delimiter"),
		Workers:        cltx.Int("workers"),
	}
}

// Stage 2: Open file and create CSV reader
func setupFileProcessing(config SetupConfig) FileContext {
	file, err := os.Open(config.CSVFilePath)
	if err != nil {
		// Return context with error set
		return FileContext{
			Setup: config,
			File:  nil,
			Error: fmt.Errorf("error opening CSV file: %w", err),
		}
	}

	reader := csv.NewReader(file)
	if config.Delimiter != "," && len(config.Delimiter) == 1 {
		reader.Comma = rune(config.Delimiter[0])
	} else if config.Delimiter != "," {
		config.Logger.Warnf("delimiter '%s' is not a single character, using default ','", config.Delimiter)
	}

	config.Logger.Infof(
		"starting update for collection %s from %s (batch size: %d)",
		config.CollectionName, config.CSVFilePath, config.BatchSize,
	)

	return FileContext{
		Setup:  config,
		File:   file,
		Reader: reader,
	}
}

// Stage 3: Validate headers
func validateHeaders(fileCtx FileContext) ProcessingContext {
	// Check for errors from previous stage
	if fileCtx.Error != nil {
		return ProcessingContext{FileContext: fileCtx}
	}

	headers, err := fileCtx.Reader.Read()
	if err != nil {
		if err == io.EOF {
			fileCtx.Error = fmt.Errorf("CSV file is empty")
		} else {
			fileCtx.Error = fmt.Errorf("error reading CSV header: %w", err)
		}
		return ProcessingContext{FileContext: fileCtx}
	}

	fileCtx.Setup.Logger.Debugf("CSV headers: %v", headers)

	if len(headers) == 0 {
		fileCtx.Error = fmt.Errorf("CSV file has no headers")
		return ProcessingContext{FileContext: fileCtx}
	}

	featurePropIDIndex := slices.Index(headers, "FEATUREPROP_ID")
	valueIndex := slices.Index(headers, "VALUE")

	if featurePropIDIndex == -1 {
		fileCtx.Error = fmt.Errorf(
			"CSV file must contain a FEATUREPROP_ID column",
		)
		return ProcessingContext{FileContext: fileCtx}
	}

	if valueIndex == -1 {
		fileCtx.Error = fmt.Errorf("CSV file must contain a VALUE column")
		return ProcessingContext{FileContext: fileCtx}
	}

	return ProcessingContext{
		FileContext:        fileCtx,
		FeaturePropIDIndex: featurePropIDIndex,
		ValueIndex:         valueIndex,
	}
}

// Helper function to submit a batch and log results
func submitBatchAndLog(
	params *SubmitBatchAndLogParams,
) (int, error) {
	if len(params.Docs) == 0 {
		return 0, nil
	}
	// Core logic from the former submitBatch function
	dbCount, dbErr := params.Setup.DBH.CountWithParams(updateAQLQuery,
		map[string]interface{}{
			"data":        params.Docs,
			"@collection": params.Setup.CollectionName,
		})

	// Error handling and logging (adapted from original submitBatchAndLog)
	if dbErr != nil {
		// Construct the error message that submitBatch would have returned
		returnedErr := fmt.Errorf("AQL query execution failed: %w", dbErr)
		// Log this error
		params.Logger.Errorf(
			"Error submitting %s batch: %v",
			params.BatchDescription,
			returnedErr,
		)
		return 0, returnedErr // Return the constructed error
	}

	// Success logging
	params.Logger.Infof(
		"Processed %s batch. Documents updated: %d.",
		params.BatchDescription,
		dbCount,
	)
	return int(dbCount), nil
}

// processSingleRecordAndValidate processes a single CSV record, validates it,
// and transforms it into a map.
// It returns the processed document and a boolean indicating if processing was successful.
func processSingleRecordAndValidate(
	params *ProcessSingleRecordParams,
) (map[string]interface{}, bool) {
	maxIndex := params.FeaturePropIDIndex
	if params.ValueIndex > maxIndex {
		maxIndex = params.ValueIndex
	}
	if len(params.Record) <= maxIndex {
		params.Logger.Warnf(
			"skipping row %d with insufficient fields (expected at least %d): %v",
			params.RowNumForLogging,
			maxIndex+1,
			params.Record,
		)
		return nil, false
	}
	processedDoc := map[string]interface{}{
		"featureprop_id": params.Record[params.FeaturePropIDIndex],
		"value":          params.Record[params.ValueIndex],
	}
	return processedDoc, true
}

// Stage 4: Process CSV records sequentially in batches.
func processCSVRecords(procCtx ProcessingContext) PipelineResult {
	currentFile := procCtx.FileContext.File
	currentSetup := procCtx.FileContext.Setup
	logger := currentSetup.Logger
	if procCtx.Error != nil { // Error from a previous stage
		return PipelineResult{
			File:  currentFile,
			Setup: currentSetup,
			Error: procCtx.Error,
		}
	}
	batchSize := currentSetup.BatchSize
	totalUpdatedCount := 0
	documentsForBatch := make([]map[string]interface{}, 0, batchSize)
	rowNum := 1 // Start with first data row (header already read)
	for {
		record, err := procCtx.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break // End of file, normal termination
			}
			return PipelineResult{
				File:  currentFile,
				Setup: currentSetup,
				Error: fmt.Errorf(
					"error reading CSV data row %d: %w",
					rowNum,
					err,
				),
			}
		}
		rowNum++
		processedDoc, isValid := processSingleRecordAndValidate(
			&ProcessSingleRecordParams{
				Record:             record,
				FeaturePropIDIndex: procCtx.FeaturePropIDIndex,
				ValueIndex:         procCtx.ValueIndex,
				RowNumForLogging:   rowNum - 1, // rowNum is 1-based for data rows, and incremented before this call
				Logger:             logger,
			})
		if !isValid {
			continue
		}
		documentsForBatch = append(documentsForBatch, processedDoc)
		if len(documentsForBatch) >= batchSize {
			count, batchErr := submitBatchAndLog(&SubmitBatchAndLogParams{
				Setup:            &currentSetup,
				Docs:             documentsForBatch,
				Logger:           logger,
				BatchDescription: "intermediate",
			})
			// batchErr is logged by submitBatchAndLog; current logic continues on error.
			totalUpdatedCount += count
			if batchErr == nil { // Log total only on success of batch
				logger.Infof("Total updated so far: %d", totalUpdatedCount)
			}
			documentsForBatch = slices.Delete(
				documentsForBatch,
				0,
				len(documentsForBatch),
			)
		}
	}
	if len(documentsForBatch) > 0 {
		count, batchErr := submitBatchAndLog(&SubmitBatchAndLogParams{
			Setup:            &currentSetup,
			Docs:             documentsForBatch,
			Logger:           logger,
			BatchDescription: "final",
		})
		totalUpdatedCount += count
		if batchErr == nil { // Log total only on success of batch
			logger.Infof(
				"Total updated after final batch: %d",
				totalUpdatedCount,
			)
		}
	}
	logger.Infof(
		"Successfully processed all records. Total documents updated: %d",
		totalUpdatedCount,
	)
	return PipelineResult{
		File:        currentFile,
		Setup:       currentSetup,
		UpdateCount: totalUpdatedCount,
	}
}
