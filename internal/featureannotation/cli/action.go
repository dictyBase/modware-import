package cli

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AQL query to update documents based on featureprop_id
const updateAQLQuery = `
	FOR row IN @data
		FOR prop IN @@collection
			FILTER prop.featureprop_id == row.featureprop_id
			UPDATE prop WITH { value: row.value } IN @@collection
			RETURN NEW
`

// DefaultUserName is the default creator/updater for annotations
const DefaultUserName = "dcr@dictycr.org"

var annMap = map[string]string{
	"CGM_DDB_PASC": "pgaudet@northwestern.edu",
	"CGM_DDB_PFEY": "pfey@northwestern.edu",
	"CGM_DDB_BOBD": "robert-dodson@northwestern.edu",
	"CGM_DDB_KPIL": "kpilchar@northwestern.edu",
	"CGM_DDB":      "dictybase@northwestern.edu",
}

// processGeneEntryParams holds the parameters for the processGeneEntry
// function.
type processGeneEntryParams struct {
	entry  *Gene
	dbh    *arangomanager.Database
	client feature.FeatureAnnotationServiceClient
	logger *logrus.Entry
}

// submitBatchParams holds the parameters for the submitBatch function.
type submitBatchParams struct {
	dbh            *arangomanager.Database
	collectionName string
	docs           []map[string]interface{}
	logger         *logrus.Entry
}

// processRecordParams holds the parameters for the processRecord function.
type processRecordParams struct {
	record             []string
	featurePropIDIndex int
	valueIndex         int
	logger             *logrus.Entry
}

type Gene struct {
	FeatureID int    `json:"feature_id"`
	GeneID    string `json:"gene_id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

// Stage 1: Setup
type SetupConfig struct {
	Logger         *logrus.Entry
	DBH            *arangomanager.Database
	CSVFilePath    string
	CollectionName string
	BatchSize      int
	Delimiter      string
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

// Stage 4: Data Processing
type DataProcessingResult struct {
	ProcessingContext
	Documents   []map[string]interface{}
	UpdateCount int
}

// Stage 5: Final Result
type FinalResult struct {
	TotalUpdated int
	Success      bool
	Error        error
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

	featurePropIDIndex := slices.Index(headers, "featureprop_id")
	valueIndex := slices.Index(headers, "value")

	if featurePropIDIndex == -1 {
		fileCtx.Error = fmt.Errorf(
			"CSV file must contain a 'featureprop_id' column",
		)
		return ProcessingContext{FileContext: fileCtx}
	}

	if valueIndex == -1 {
		fileCtx.Error = fmt.Errorf("CSV file must contain a 'value' column")
		return ProcessingContext{FileContext: fileCtx}
	}

	return ProcessingContext{
		FileContext:        fileCtx,
		FeaturePropIDIndex: featurePropIDIndex,
		ValueIndex:         valueIndex,
	}
}

// Stage 4: Process CSV records into batches
func processCSVRecords(procCtx ProcessingContext) DataProcessingResult {
	// Check for errors from previous stage
	if procCtx.Error != nil {
		return DataProcessingResult{ProcessingContext: procCtx}
	}

	result := DataProcessingResult{
		ProcessingContext: procCtx,
		Documents: make(
			[]map[string]interface{},
			0,
			procCtx.Setup.BatchSize,
		),
		UpdateCount: 0,
	}

	rowNum := 1 // Start with first data row (header already read)

	for {
		record, err := procCtx.Reader.Read()
		if err != nil {
			if err == io.EOF {
				break // End of file, normal termination
			}
			// Critical read error
			result.Error = fmt.Errorf(
				"error reading CSV data row %d: %w",
				rowNum,
				err,
			)
			return result // Return immediately with the error
		}

		// Process the record using the existing helper
		doc := processRecord(&processRecordParams{
			record:             record,
			featurePropIDIndex: procCtx.FeaturePropIDIndex,
			valueIndex:         procCtx.ValueIndex,
			logger:             procCtx.Setup.Logger,
		})

		if doc == nil {
			rowNum++ // Skip this row but continue
			continue
		}

		result.Documents = append(result.Documents, doc)

		// If batch is full, process it
		if len(result.Documents) >= procCtx.Setup.BatchSize {
			batchResult, batchErr := submitBatch(&submitBatchParams{
				dbh:            procCtx.Setup.DBH,
				collectionName: procCtx.Setup.CollectionName,
				docs:           result.Documents,
				logger:         procCtx.Setup.Logger,
			})

			result.UpdateCount += batchResult // Accumulate count

			if batchErr != nil {
				result.Error = fmt.Errorf(
					"error submitting batch after row %d: %w",
					rowNum,
					batchErr,
				)
				// Clear documents as the batch failed, but keep the count processed so far
				result.Documents = nil
				return result // Return immediately with the error
			}

			procCtx.Setup.Logger.Infof(
				"Processed batch ending at row %d. Documents updated in this batch: %d. Total updated so far: %d",
				rowNum,
				batchResult,
				result.UpdateCount,
			)

			// Reset batch documents
			result.Documents = make(
				[]map[string]interface{},
				0,
				procCtx.Setup.BatchSize,
			)
		}

		rowNum++
	}

	// Return the result containing any remaining documents for the final batch
	return result
}

// Stage 5: Submit batches and finalize
func finalizeBatchProcessing(result DataProcessingResult) FinalResult {
	// Always ensure the file is closed if it was opened
	if result.File != nil {
		defer result.File.Close()
	}

	// Check for errors propagated from previous stages
	if result.Error != nil {
		return FinalResult{
			TotalUpdated: result.UpdateCount, // Report count processed before error
			Success:      false,
			Error:        result.Error,
		}
	}

	// Process final batch if any documents remain
	finalCount := result.UpdateCount
	if len(result.Documents) > 0 {
		batchResult, err := submitBatch(&submitBatchParams{
			dbh:            result.Setup.DBH,
			collectionName: result.Setup.CollectionName,
			docs:           result.Documents,
			logger:         result.Setup.Logger,
		})

		finalCount += batchResult // Add final batch count

		if err != nil {
			// Error occurred during the final batch submission
			return FinalResult{
				TotalUpdated: finalCount, // Report total including this partial/failed batch
				Success:      false,
				Error: fmt.Errorf(
					"error submitting final batch: %w",
					err,
				),
			}
		}

		result.Setup.Logger.Infof(
			"Processed final batch. Documents updated in this batch: %d. Total updated: %d",
			batchResult,
			finalCount,
		)
	}

	// Log success only if no errors occurred throughout the pipeline
	result.Setup.Logger.Infof(
		"successfully finished processing CSV for collection %s. Total documents updated: %d",
		result.Setup.CollectionName,
		finalCount,
	)

	return FinalResult{
		TotalUpdated: finalCount,
		Success:      true,
		Error:        nil,
	}
}

func LoadCSVToArangodb(cltx *cli.Context) error {
	result := collection.Pipe5(
		cltx,
		setupPipeline,
		setupFileProcessing,
		validateHeaders,
		processCSVRecords,
		finalizeBatchProcessing,
	)

	// Check the final result for errors
	if result.Error != nil {
		registry.GetLogger().
			Errorf("Pipeline failed: %s", result.Error)
		return cli.Exit(
			fmt.Sprintf("Pipeline execution failed: %s", result.Error.Error()),
			2,
		)
	}
	// If successful, log the final count (already logged in
	// finalizeBatchProcessing, but good for confirmation)
	registry.GetLogger().
		Infof(
			"Pipeline completed successfully. Total documents updated: %d",
			result.TotalUpdated,
		)
	return nil
}

func LoadFeatureAnnotation(cltx *cli.Context) error {
	logger := registry.GetLogger()
	dbh := registry.GetArangodbConnection()
	client := registry.GetFeatureAnnotationAPIClient()

	result, err := dbh.Search(ListActiveGenesQ)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	defer result.Close()

	if result.IsEmpty() {
		logger.Error("No PubMed references found in feature annotations")
		return nil
	}

	for result.Scan() {
		entry := &Gene{}
		if err := result.Read(&entry); err != nil {
			return cli.Exit(
				fmt.Sprintf("error reading query result: %s", err),
				2,
			)
		}
		// Call helper function to process the entry
		params := &processGeneEntryParams{
			entry:  entry,
			dbh:    dbh,
			client: client,
			logger: logger,
		}
		if err := processGeneEntry(params); err != nil {
			// Exit if the helper function encounters an error
			return cli.Exit(err.Error(), 2)
		}
	}

	return nil
}

// fetchPubmedIDs queries and retrieves PubMed IDs associated with a given gene feature.
func fetchPubmedIDs(
	entry *Gene,
	dbh *arangomanager.Database,
	logger *logrus.Entry,
) ([]string, error) {
	pubmedResult, err := dbh.SearchRows(
		ListPubmedsByFeature,
		map[string]interface{}{
			"feature_id": entry.FeatureID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error querying PubMed IDs for feature %d: %w",
			entry.FeatureID,
			err,
		)
	}
	defer pubmedResult.Close()

	if pubmedResult.IsEmpty() {
		logger.Infof(
			"No PubMed references found for feature %d",
			entry.FeatureID,
		)
		return []string{}, nil // Return empty slice instead of nil
	}

	pubmedIDs := make([]string, 0)
	for pubmedResult.Scan() {
		var pubmed string
		if err := pubmedResult.Read(&pubmed); err != nil {
			return nil, fmt.Errorf(
				"error reading PubMed ID for feature %d: %w",
				entry.FeatureID,
				err,
			)
		}
		pubmedIDs = append(pubmedIDs, pubmed)
	}
	logger.Infof("Feature %d has PubMed reference: %v",
		entry.FeatureID,
		pubmedIDs,
	)
	return pubmedIDs, nil
}

// processGeneEntry handles fetching PubMed IDs and creating the annotation for a single gene entry.
func processGeneEntry(params *processGeneEntryParams) error {
	logger := params.logger
	entry := params.entry
	dbh := params.dbh
	client := params.client

	logger.Debugf("Feature has geneid %s", entry.GeneID)

	// Fetch PubMed IDs
	pubmedIDs, err := fetchPubmedIDs(entry, dbh, logger)
	if err != nil {
		return err // Propagate error from fetching IDs
	}

	// If no PubMed IDs were found or pubmedIDs is nil, skip creating the annotation for this entry
	if len(pubmedIDs) == 0 {
		logger.Infof("Skipping feature %s with no PubMed IDs", entry.GeneID)
		return nil
	}

	createdBy := DefaultUserName
	if val, ok := annMap[entry.CreatedBy]; ok {
		createdBy = val
	}

	// Set up gRPC call
	res, err := client.CreateFeatureAnnotation(
		context.Background(),
		&feature.NewFeatureAnnotation{
			Type:       "gene",
			Id:         entry.GeneID,
			IsObsolete: false,
			CreatedBy:  createdBy,
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:   entry.Name,
				Pubmed: pubmedIDs,
			},
		})
	if err != nil {
		return fmt.Errorf(
			"failed to create feature annotation for gene %s: %w",
			entry.GeneID,
			err,
		)
	}

	logger.Infof(
		"Created new feature annotation record %s for feature name %s",
		res.Attributes.Name,
		res.Id,
	)
	return nil
}

// processRecord validates a CSV record and converts it to a document map.
// Returns the document map, or nil if the record should be skipped.
func processRecord(
	params *processRecordParams,
) map[string]interface{} {
	maxIndex := params.featurePropIDIndex
	if params.valueIndex > maxIndex {
		maxIndex = params.valueIndex
	}
	// Check if record has enough columns up to the highest required index
	if len(params.record) <= maxIndex {
		params.logger.Warnf(
			"skipping row with insufficient fields (expected at least %d): %v",
			maxIndex+1,
			params.record,
		)
		return nil // Indicate skip
	}
	return map[string]interface{}{
		"featureprop_id": params.record[params.featurePropIDIndex],
		"value":          params.record[params.valueIndex],
	}
}

// submitBatch updates a batch of documents in ArangoDB using the updateAQLQuery.
// Returns the number of documents successfully updated, or an error.
func submitBatch(
	params *submitBatchParams,
) (int, error) {
	if len(params.docs) == 0 {
		return 0, nil
	}
	count, err := params.dbh.CountWithParams(updateAQLQuery,
		map[string]interface{}{
			"data":        params.docs,
			"@collection": params.collectionName,
		})
	if err != nil {
		return 0, fmt.Errorf("AQL query execution failed: %w", err)
	}
	return int(count), nil
}
