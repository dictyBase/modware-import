package cli

import (
	"encoding/csv"
	"fmt"
	"iter"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/urfave/cli/v2"
)

// GeneDataRecord represents a single gene data record
type GeneDataRecord struct {
	GeneID      string `validate:"required,min=1" json:"gene_id"`
	GeneProduct string `                          json:"gene_product"`
	Description string `                          json:"description"`
}

// ParseUnknowmeDataParams contains all parameters for the ParseUnknowmeData function
type ParseUnknowmeDataParams struct {
	InputFiles            []string `validate:"required,min=1,dive,min=1" json:"input_files"`
	GeneProductOutput     string   `validate:"required,min=1" json:"gene_product_output"`
	GeneDescriptionOutput string   `validate:"required,min=1" json:"gene_description_output"`
}

// ParsingConfig holds configuration for parsing operations
type ParsingConfig struct {
	ddbGeneRegex           *regexp.Regexp
	geneProductStartColumn int
	geneProductEndColumn   int
	skipEmptyProduct       bool
	skipEmptyDescription   bool
}

// ParsingOption is a function that configures ParsingConfig
type ParsingOption func(*ParsingConfig)

// WithGeneProductColumnRange sets the column range for searching gene products
func WithGeneProductColumnRange(startCol, endCol int) ParsingOption {
	return func(config *ParsingConfig) {
		config.geneProductStartColumn = startCol
		config.geneProductEndColumn = endCol
	}
}

// WithSkipEmptyProduct controls whether to skip records with empty gene products
func WithSkipEmptyProduct(skip bool) ParsingOption {
	return func(config *ParsingConfig) {
		config.skipEmptyProduct = skip
	}
}

// WithSkipEmptyDescription controls whether to skip records with empty descriptions
func WithSkipEmptyDescription(skip bool) ParsingOption {
	return func(config *ParsingConfig) {
		config.skipEmptyDescription = skip
	}
}

// NewParsingConfig creates a new ParsingConfig with default values and applies options
func NewParsingConfig(opts ...ParsingOption) *ParsingConfig {
	config := &ParsingConfig{
		ddbGeneRegex:           regexp.MustCompile(`^DDB_G\d+`),
		geneProductStartColumn: 3,
		geneProductEndColumn:   7,
		skipEmptyProduct:       true,
		skipEmptyDescription:   true,
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

// ParseUnknowmeData processes HTML files containing gene data tables and extracts
// DDB_G entries with their corresponding gene product and description information
func ParseUnknowmeData(cliCtx *cli.Context) error {
	params := ParseUnknowmeDataParams{
		InputFiles:            cliCtx.StringSlice("input"),
		GeneProductOutput:     cliCtx.String("gene-product-output"),
		GeneDescriptionOutput: cliCtx.String("gene-description-output"),
	}

	// Validate parameters
	if err := ValidateStruct(params); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return processUnknowmeDataWithParams(params)
}

// processUnknowmeDataWithParams performs the actual processing with validated parameters
func processUnknowmeDataWithParams(params ParseUnknowmeDataParams) error {
	// Validate all input files exist
	if err := validateInputFilesExist(params.InputFiles); err != nil {
		return fmt.Errorf("input file validation failed: %w", err)
	}

	// Create parsing configuration with default options
	parsingConfig := NewParsingConfig()

	// Create consolidated iterator for processing gene data from multiple files
	geneDataIterator, err := createMultiFileGeneDataIterator(
		params.InputFiles,
		parsingConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to create gene data iterator: %w", err)
	}

	// Create CSV writers using functional approach
	csvWriterParams := CSVWriterCreationParams{
		GeneProductOutputFile:     params.GeneProductOutput,
		GeneDescriptionOutputFile: params.GeneDescriptionOutput,
	}

	writers, err := createCSVWriters(csvWriterParams)
	if err != nil {
		return fmt.Errorf("failed to create CSV writers: %w", err)
	}
	defer closeAllWriters(writers)

	// Process records with unique gene ID collection and count for reporting
	processingResult, err := processUniqueGeneDataRecords(geneDataIterator, writers)
	if err != nil {
		return fmt.Errorf("failed to process gene data records: %w", err)
	}

	// Validate processing results
	if err := validateProcessingResults(processingResult); err != nil {
		return fmt.Errorf("processing validation failed: %w", err)
	}

	// Report successful processing
	reportProcessingResults(processingResult, params)

	return nil
}

// CSVWriterCreationParams contains parameters for creating CSV writers
type CSVWriterCreationParams struct {
	GeneProductOutputFile     string `validate:"required,min=1" json:"gene_product_output_file"`
	GeneDescriptionOutputFile string `validate:"required,min=1" json:"gene_description_output_file"`
}

// ProcessingResult contains the results of processing gene data records
type ProcessingResult struct {
	TotalRecordsProcessed     int
	GeneProductRecordsWritten int
	DescriptionRecordsWritten int
}

// validateInputFileExists checks if the input file exists and is readable
func validateInputFileExists(inputFile string) error {
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}
	return nil
}

// validateInputFilesExist checks if all input files exist and are readable
func validateInputFilesExist(inputFiles []string) error {
	for _, inputFile := range inputFiles {
		if err := validateInputFileExists(inputFile); err != nil {
			return fmt.Errorf("validation failed for file %s: %w", inputFile, err)
		}
	}
	return nil
}

// createMultiFileGeneDataIterator creates an iterator that processes gene data from multiple HTML files
func createMultiFileGeneDataIterator(
	filenames []string,
	config *ParsingConfig,
) (iter.Seq[GeneDataRecord], error) {
	// Load all HTML documents first to validate them
	htmlDocuments := make([]*goquery.Document, 0, len(filenames))
	for _, filename := range filenames {
		htmlDocument, err := loadHTMLDocument(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to load HTML document %s: %w", filename, err)
		}
		htmlDocuments = append(htmlDocuments, htmlDocument)
	}

	return createIteratorFromMultipleDocuments(htmlDocuments, config), nil
}

// createIteratorFromMultipleDocuments creates an iterator from multiple parsed HTML documents
func createIteratorFromMultipleDocuments(
	htmlDocuments []*goquery.Document,
	config *ParsingConfig,
) iter.Seq[GeneDataRecord] {
	return func(yield func(GeneDataRecord) bool) {
		for _, htmlDocument := range htmlDocuments {
			htmlDocument.Find("table tr").
				Each(func(i int, row *goquery.Selection) {
					record, shouldProcess := processTableRow(row, config)
					if shouldProcess {
						if !yield(record) {
							return // Early termination requested by consumer
						}
					}
				})
		}
	}
}

// createCSVWriters creates all required CSV writers
func createCSVWriters(
	params CSVWriterCreationParams,
) ([]CSVRecordWriter, error) {
	if err := ValidateStruct(params); err != nil {
		return nil, fmt.Errorf(
			"CSV writer parameters validation failed: %w",
			err,
		)
	}

	geneProductWriter, err := NewGeneProductCSVWriter(
		params.GeneProductOutputFile,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gene product writer: %w", err)
	}

	geneDescriptionWriter, err := NewGeneDescriptionCSVWriter(
		params.GeneDescriptionOutputFile,
	)
	if err != nil {
		geneProductWriter.Close()
		return nil, fmt.Errorf(
			"failed to create gene description writer: %w",
			err,
		)
	}

	return []CSVRecordWriter{geneProductWriter, geneDescriptionWriter}, nil
}

// closeAllWriters closes all CSV writers and handles errors
func closeAllWriters(writers []CSVRecordWriter) {
	for _, writer := range writers {
		if err := writer.Close(); err != nil {
			// Log error but don't fail since this is cleanup
			fmt.Printf("Warning: failed to close writer: %v\n", err)
		}
	}
}

// processUniqueGeneDataRecords processes gene data records with unique gene ID collection
func processUniqueGeneDataRecords(
	geneDataIterator iter.Seq[GeneDataRecord],
	writers []CSVRecordWriter,
) (*ProcessingResult, error) {
	result := &ProcessingResult{}
	uniqueGenes := make(map[string]GeneDataRecord)

	// First pass: collect unique gene records
	for record := range geneDataIterator {
		result.TotalRecordsProcessed++

		// Validate individual record
		if err := ValidateStruct(record); err != nil {
			return nil, fmt.Errorf(
				"record validation failed for gene %s: %w",
				record.GeneID,
				err,
			)
		}

		// Store record if gene ID is not already present, or merge/update as needed
		if existingRecord, exists := uniqueGenes[record.GeneID]; exists {
			// Apply merging logic: prefer non-empty values
			mergedRecord := mergeGeneDataRecords(existingRecord, record)
			uniqueGenes[record.GeneID] = mergedRecord
		} else {
			uniqueGenes[record.GeneID] = record
		}
	}

	// Second pass: process unique records with writers
	for _, uniqueRecord := range uniqueGenes {
		err := processRecordWithWriters(uniqueRecord, writers, result)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to process unique record %s: %w",
				uniqueRecord.GeneID,
				err,
			)
		}
	}

	return result, nil
}

// mergeGeneDataRecords merges two gene data records, preferring non-empty values
func mergeGeneDataRecords(existing, newRecord GeneDataRecord) GeneDataRecord {
	merged := existing

	// Prefer non-empty gene product
	if strings.TrimSpace(newRecord.GeneProduct) != "" && strings.TrimSpace(existing.GeneProduct) == "" {
		merged.GeneProduct = newRecord.GeneProduct
	}

	// Prefer non-empty description
	if strings.TrimSpace(newRecord.Description) != "" && strings.TrimSpace(existing.Description) == "" {
		merged.Description = newRecord.Description
	}

	return merged
}

// processRecordWithWriters processes a single record with all writers
func processRecordWithWriters(
	record GeneDataRecord,
	writers []CSVRecordWriter,
	result *ProcessingResult,
) error {
	for i, writer := range writers {
		if !writer.ShouldSkip(record) {
			if err := writer.WriteRecord(record); err != nil {
				return fmt.Errorf(
					"failed to write record with writer %d: %w",
					i,
					err,
				)
			}

			// Update counters based on writer type
			switch writer.(type) {
			case *GeneProductCSVWriter:
				result.GeneProductRecordsWritten++
			case *GeneDescriptionCSVWriter:
				result.DescriptionRecordsWritten++
			}
		}
	}
	return nil
}

// validateProcessingResults validates the processing results
func validateProcessingResults(result *ProcessingResult) error {
	if result.TotalRecordsProcessed == 0 {
		return fmt.Errorf("no gene data records found in the HTML file")
	}
	return nil
}

// reportProcessingResults reports the results of processing
func reportProcessingResults(
	result *ProcessingResult,
	params ParseUnknowmeDataParams,
) {
	fmt.Printf(
		"Successfully processed %d gene records from %d input file(s)\n",
		result.TotalRecordsProcessed,
		len(params.InputFiles),
	)
	fmt.Printf("Input files processed:\n")
	for i, inputFile := range params.InputFiles {
		fmt.Printf("  %d. %s\n", i+1, inputFile)
	}
	fmt.Printf(
		"Gene products written to: %s (%d unique records)\n",
		params.GeneProductOutput,
		result.GeneProductRecordsWritten,
	)
	fmt.Printf("Gene descriptions written to: %s (%d unique records)\n",
		params.GeneDescriptionOutput, result.DescriptionRecordsWritten)
}

// processTableRow processes a single table row and returns a gene data record
func processTableRow(
	row *goquery.Selection,
	config *ParsingConfig,
) (GeneDataRecord, bool) {
	cells := row.Find("td")
	cellCount := cells.Length()

	// Skip rows that don't have at least 3 cells
	if cellCount < 3 {
		return GeneDataRecord{}, false
	}

	// Extract and validate gene ID
	geneID := extractGeneID(cells)
	if !isValidGeneID(geneID, config.ddbGeneRegex) {
		return GeneDataRecord{}, false
	}

	// Extract gene product and description using functional approach
	geneProduct, geneProductColumn := extractGeneProductWithIndex(
		cells,
		config,
	)
	description := extractGeneDescription(
		cells,
		geneProductColumn,
		cellCount,
	)

	record := GeneDataRecord{
		GeneID:      geneID,
		GeneProduct: geneProduct,
		Description: description,
	}

	return record, true
}

// extractGeneID extracts the gene ID from the first cell
func extractGeneID(cells *goquery.Selection) string {
	return strings.TrimSpace(cells.Eq(0).Text())
}

// isValidGeneID checks if the gene ID matches the expected pattern
func isValidGeneID(geneID string, regex *regexp.Regexp) bool {
	return regex.MatchString(geneID)
}

// extractGeneProductWithIndex extracts gene product and returns its column
// index
func extractGeneProductWithIndex(
	cells *goquery.Selection,
	config *ParsingConfig,
) (string, int) {
	return extractTextFromColumnsWithIndex(
		cells,
		config.geneProductStartColumn,
		config.geneProductEndColumn,
	)
}

// extractGeneDescription extracts description based on gene product location
func extractGeneDescription(
	cells *goquery.Selection,
	geneProductColumn, cellCount int,
) string {
	if geneProductColumn != -1 {
		// Gene product found: start description scanning after gene product column
		return extractTextFromColumns(cells, geneProductColumn+1, cellCount-1)
	}
	// Gene product not found: scan all remaining cells for description
	return extractTextFromColumns(cells, 1, cellCount-1)
}

// loadHTMLDocument loads and parses an HTML file
func loadHTMLDocument(filename string) (*goquery.Document, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return doc, nil
}

// skipGeneProduct determines if a gene product record should be skipped
func skipGeneProduct(record GeneDataRecord) bool {
	return isEmptyGeneProduct(record.GeneProduct) ||
		containsExcludedGeneProductTerms(record.GeneProduct)
}

// isEmptyGeneProduct checks if the gene product is empty
func isEmptyGeneProduct(geneProduct string) bool {
	return strings.TrimSpace(geneProduct) == ""
}

// containsExcludedGeneProductTerms checks if gene product contains excluded terms
func containsExcludedGeneProductTerms(geneProduct string) bool {
	excludedTerms := []string{"no gp", "unknown"}
	lowerGeneProduct := strings.ToLower(geneProduct)

	_, found := collection.Find(excludedTerms, func(term string) bool {
		return strings.Contains(lowerGeneProduct, term)
	})
	return found
}

// CSVRecordWriter defines a common interface for writing gene data records to CSV
type CSVRecordWriter interface {
	WriteRecord(record GeneDataRecord) error
	ShouldSkip(record GeneDataRecord) bool
	Close() error
}

// GeneProductCSVWriter writes gene product data to CSV format
type GeneProductCSVWriter struct {
	writer *csv.Writer
	file   *os.File
}

// GeneProductCSVWriterParams contains parameters for creating a gene product CSV writer
type GeneProductCSVWriterParams struct {
	Filename string `validate:"required,min=1" json:"filename"`
}

// NewGeneProductCSVWriter creates a new CSV writer for gene products
func NewGeneProductCSVWriter(filename string) (*GeneProductCSVWriter, error) {
	params := GeneProductCSVWriterParams{Filename: filename}
	return NewGeneProductCSVWriterWithParams(params)
}

// NewGeneProductCSVWriterWithParams creates a new CSV writer for gene products with validated parameters
func NewGeneProductCSVWriterWithParams(
	params GeneProductCSVWriterParams,
) (*GeneProductCSVWriter, error) {
	if err := ValidateStruct(params); err != nil {
		return nil, fmt.Errorf(
			"gene product CSV writer parameters validation failed: %w",
			err,
		)
	}

	file, err := os.Create(params.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create gene product file: %w", err)
	}

	csvWriter := csv.NewWriter(file)

	// Write header using predefined structure
	header := []string{"GeneID", "gene_product"}
	if err := csvWriter.Write(header); err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to write gene product header: %w", err)
	}

	return &GeneProductCSVWriter{
		writer: csvWriter,
		file:   file,
	}, nil
}

func (w *GeneProductCSVWriter) WriteRecord(record GeneDataRecord) error {
	return w.writer.Write([]string{record.GeneID, record.GeneProduct})
}

func (w *GeneProductCSVWriter) ShouldSkip(record GeneDataRecord) bool {
	return skipGeneProduct(record)
}

func (w *GeneProductCSVWriter) Close() error {
	w.writer.Flush()
	return w.file.Close()
}

// GeneDescriptionCSVWriter writes gene description data to CSV format
type GeneDescriptionCSVWriter struct {
	writer *csv.Writer
	file   *os.File
}

// GeneDescriptionCSVWriterParams contains parameters for creating a gene description CSV writer
type GeneDescriptionCSVWriterParams struct {
	Filename string `validate:"required,min=1" json:"filename"`
}

// NewGeneDescriptionCSVWriter creates a new CSV writer for gene descriptions
func NewGeneDescriptionCSVWriter(
	filename string,
) (*GeneDescriptionCSVWriter, error) {
	params := GeneDescriptionCSVWriterParams{Filename: filename}
	return NewGeneDescriptionCSVWriterWithParams(params)
}

// NewGeneDescriptionCSVWriterWithParams creates a new CSV writer for gene descriptions with validated parameters
func NewGeneDescriptionCSVWriterWithParams(
	params GeneDescriptionCSVWriterParams,
) (*GeneDescriptionCSVWriter, error) {
	if err := ValidateStruct(params); err != nil {
		return nil, fmt.Errorf(
			"gene description CSV writer parameters validation failed: %w",
			err,
		)
	}

	file, err := os.Create(params.Filename)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create gene description file: %w",
			err,
		)
	}

	csvWriter := csv.NewWriter(file)

	// Write header using predefined structure
	header := []string{"GeneID", "gene_description"}
	if err := csvWriter.Write(header); err != nil {
		file.Close()
		return nil, fmt.Errorf(
			"failed to write gene description header: %w",
			err,
		)
	}

	return &GeneDescriptionCSVWriter{
		writer: csvWriter,
		file:   file,
	}, nil
}

func (w *GeneDescriptionCSVWriter) WriteRecord(record GeneDataRecord) error {
	return w.writer.Write([]string{record.GeneID, record.Description})
}

func (w *GeneDescriptionCSVWriter) ShouldSkip(record GeneDataRecord) bool {
	return shouldSkipGeneDescription(record)
}

// shouldSkipGeneDescription determines if a gene description record should be skipped
func shouldSkipGeneDescription(record GeneDataRecord) bool {
	return isEmptyGeneDescription(record.Description)
}

// isEmptyGeneDescription checks if the gene description is empty
func isEmptyGeneDescription(description string) bool {
	return strings.TrimSpace(description) == ""
}

func (w *GeneDescriptionCSVWriter) Close() error {
	w.writer.Flush()
	return w.file.Close()
}
