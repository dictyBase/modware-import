package cli

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	E "github.com/IBM/fp-go/either"
	O "github.com/IBM/fp-go/option"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

// ValidationResult represents the result of validating a single gene
type ValidationResult struct {
	GeneID         string `json:"gene_id"`
	CSVDescription string `json:"csv_description"`
	GQLDescription string `json:"gql_description"`
	Match          bool   `json:"match"`
	Error          string `json:"error,omitempty"`
}

// ValidationReport contains the complete validation report
type ValidationReport struct {
	TotalGenes    int                `json:"total_genes"`
	MatchCount    int                `json:"match_count"`
	MismatchCount int                `json:"mismatch_count"`
	ErrorCount    int                `json:"error_count"`
	Results       []ValidationResult `json:"results"`
	GeneratedAt   time.Time          `json:"generated_at"`
}

// GeneValidationParams contains all parameters needed for gene validation
type GeneValidationParams struct {
	InputFile    string `validate:"required,file_exists,max_file_size"`
	OutputReport string `validate:"required,endswith=.json"`
	GraphQLURL   string `validate:"required,url,startswith=https"`
	Timeout      int    `validate:"gte=10,lte=300"`
	Workers      int    `validate:"gte=1,lte=20"`
}

// SingleGeneValidationParams contains parameters for validating a single gene
type SingleGeneValidationParams struct {
	Client     *http.Client
	GraphQLURL string
	Record     []string
}

// ConcurrentValidationParams contains parameters for concurrent validation
type ConcurrentValidationParams struct {
	ValidationParams GeneValidationParams
	Records          [][]string
	Context          context.Context
}

// CSVFileConstraints defines file size and content limits
type CSVFileConstraints struct {
	MaxFileSizeBytes int64
	MaxRecords       int
}

var csvConstraints = CSVFileConstraints{
	MaxFileSizeBytes: 50 * 1024 * 1024, // 50MB
	MaxRecords:       100000,           // 100k records
}

// GraphQLResponse represents the response from the GraphQL endpoint
type GraphQLResponse struct {
	Data struct {
		GeneGeneralInformation *struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		} `json:"geneGeneralInformation"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

var geneValidator = validator.New()

// Custom validation functions
func init() {
	geneValidator.RegisterValidation("file_exists", validateFileExists)
	geneValidator.RegisterValidation("startswith", validateStartsWith)
	geneValidator.RegisterValidation("endswith", validateEndsWith)
	geneValidator.RegisterValidation("max_file_size", validateMaxFileSize)
}

func validateFileExists(fl validator.FieldLevel) bool {
	filePath := fl.Field().String()
	if filePath == "" {
		return false
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	// Ensure it's actually a file, not a directory
	return !fileInfo.IsDir()
}

func validateStartsWith(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()
	return strings.HasPrefix(value, param)
}

func validateEndsWith(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	param := fl.Param()
	return strings.HasSuffix(value, param)
}

func validateMaxFileSize(fl validator.FieldLevel) bool {
	filePath := fl.Field().String()
	if filePath == "" {
		return false
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}

	return fileInfo.Size() <= csvConstraints.MaxFileSizeBytes
}

// ValidateGeneData is the main action function for gene data validation
func ValidateGeneData(c *cli.Context) error {
	params := GeneValidationParams{
		InputFile:    c.String("input"),
		OutputReport: c.String("output-report"),
		GraphQLURL:   c.String("graphql-url"),
		Timeout:      c.Int("timeout"),
		Workers:      c.Int("workers"),
	}

	// Enhanced validation with better error context
	if err := geneValidator.Struct(params); err != nil {
		return fmt.Errorf(
			"gene validation parameters invalid (file: %s, url: %s): %w",
			params.InputFile,
			params.GraphQLURL,
			err,
		)
	}

	// Additional URL scheme validation for security
	if err := validateGraphQLURLScheme(params.GraphQLURL); err != nil {
		return fmt.Errorf("GraphQL URL security validation failed: %w", err)
	}

	// Read CSV file with security constraints
	records, err := readCSVFileWithValidation(params.InputFile)
	if err != nil {
		return fmt.Errorf("validation process failed: %w", err)
	}

	// Validate gene descriptions with functional programming patterns
	report, err := validateGeneDescriptionsWithFP(params, records)
	if err != nil {
		return fmt.Errorf("validation process failed: %w", err)
	}

	return saveReport(params.OutputReport, report)
}

// validateGraphQLURLScheme ensures GraphQL URL uses HTTPS for security
func validateGraphQLURLScheme(graphqlURL string) error {
	parsedURL, err := url.Parse(graphqlURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme != "https" {
		return fmt.Errorf(
			"GraphQL URL must use HTTPS scheme, got: %s",
			parsedURL.Scheme,
		)
	}

	return nil
}

// readCSVFile reads and parses the CSV file containing gene descriptions (legacy function)
func readCSVFile(filePath string) ([][]string, error) {
	return readCSVFileWithValidation(filePath)
}

// readCSVFileWithValidation reads CSV file with security constraints and validation
func readCSVFileWithValidation(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file %s: %w", filePath, err)
	}
	defer file.Close()

	records, err := readAndValidateCSVRecordsTraditional(file)
	if err != nil {
		return nil, fmt.Errorf(
			"CSV validation failed for file %s: %w",
			filePath,
			err,
		)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	if len(records) == 1 {
		return nil, fmt.Errorf("CSV file contains only header row")
	}

	return records, nil
}

// readAndValidateCSVRecordsTraditional reads CSV records with size validation (traditional error handling)
func readAndValidateCSVRecordsTraditional(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	reader.ReuseRecord = true // Memory optimization for large files

	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Create a copy since ReuseRecord is true
		recordCopy := make([]string, len(record))
		copy(recordCopy, record)
		records = append(records, recordCopy)

		// Prevent memory exhaustion
		if len(records) > csvConstraints.MaxRecords {
			return nil, fmt.Errorf(
				"CSV file exceeds maximum record limit of %d",
				csvConstraints.MaxRecords,
			)
		}
	}

	return records, nil
}

// validateGeneDescriptionsWithFP validates gene descriptions with functional programming patterns
func validateGeneDescriptionsWithFP(
	params GeneValidationParams,
	records [][]string,
) (
	ValidationReport,
	error,
) {
	// Check if records is not empty
	if len(records) == 0 {
		return ValidationReport{}, fmt.Errorf("CSV file is empty")
	}

	// Skip header row using functional programming
	dataRecords := collection.Filter(
		records[1:],
		func(record []string) bool {
			return len(record) > 0 // Remove empty records
		})

	if len(dataRecords) == 0 {
		return ValidationReport{}, fmt.Errorf(
			"no data records to process after skipping header",
		)
	}

	// Process records concurrently with errgroup
	results, err := processRecordsConcurrentlyTraditional(
		context.Background(),
		params,
		dataRecords,
	)
	if err != nil {
		return ValidationReport{}, err
	}

	// Generate report using functional programming
	report := generateReportWithFP(results)
	return report, nil
}

// processRecordsConcurrentlyTraditional processes validation records using errgroup for safe concurrency
func processRecordsConcurrentlyTraditional(
	ctx context.Context,
	params GeneValidationParams,
	records [][]string,
) ([]ValidationResult, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf(
			"no data records to process after skipping header",
		)
	}

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(params.Workers)

	results := make([]ValidationResult, len(records))
	sharedClient := createSharedHTTPClient(params.Timeout)

	for i, record := range records {
		i, record := i, record // Capture loop variables
		g.Go(func() error {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
				validationParams := SingleGeneValidationParams{
					Client:     sharedClient,
					GraphQLURL: params.GraphQLURL,
					Record:     record,
				}
				results[i] = validateSingleGene(validationParams)
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("concurrent validation failed: %w", err)
	}

	return results, nil
}

// createSharedHTTPClient creates a reusable HTTP client with timeout and rate limiting
func createSharedHTTPClient(timeoutSeconds int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
	}
}

// validateSingleGene validates a single gene record against GraphQL endpoint (refactored to use parameter struct)
func validateSingleGene(params SingleGeneValidationParams) ValidationResult {
	recordOpt := validateRecordStructure(params.Record)

	return O.Fold(
		func() ValidationResult {
			return ValidationResult{
				GeneID: safeGetGeneID(params.Record),
				Error:  "CSV record validation failed",
			}
		},
		func(geneData struct{ geneID, csvDesc string }) ValidationResult {
			queryResult := queryGraphQL(
				params.Client,
				params.GraphQLURL,
				geneData.geneID,
			)
			return E.Fold(
				func(err error) ValidationResult {
					return ValidationResult{
						GeneID:         geneData.geneID,
						CSVDescription: geneData.csvDesc,
						Error:          err.Error(),
					}
				},
				func(gqlDescription string) ValidationResult {
					match := strings.EqualFold(
						strings.TrimSpace(geneData.csvDesc),
						strings.TrimSpace(gqlDescription),
					)
					return ValidationResult{
						GeneID:         geneData.geneID,
						CSVDescription: geneData.csvDesc,
						GQLDescription: gqlDescription,
						Match:          match,
					}
				},
			)(queryResult)
		},
	)(recordOpt)
}

// validateRecordStructure validates CSV record structure and extracts gene data
func validateRecordStructure(
	record []string,
) O.Option[struct{ geneID, csvDesc string }] {
	if len(record) < 2 {
		return O.None[struct{ geneID, csvDesc string }]()
	}

	geneID := strings.TrimSpace(record[0])
	csvDescription := strings.TrimSpace(record[1])

	if geneID == "" {
		return O.None[struct{ geneID, csvDesc string }]()
	}

	return O.Some(struct{ geneID, csvDesc string }{
		geneID:  geneID,
		csvDesc: csvDescription,
	})
}

// safeGetGeneID safely extracts gene ID from record
func safeGetGeneID(record []string) string {
	if len(record) > 0 {
		return strings.TrimSpace(record[0])
	}
	return ""
}

// queryGraphQL queries the GraphQL endpoint for gene information using Either monad
func queryGraphQL(
	client *http.Client,
	graphqlURL, geneID string,
) E.Either[error, string] {
	// Use a traditional approach to handle HTTP response properly
	result, err := queryGraphQLTraditional(client, graphqlURL, geneID)
	if err != nil {
		return E.Left[string](
			fmt.Errorf("GraphQL query failed for gene %s: %w", geneID, err),
		)
	}
	return E.Right[error](result)
}

// queryGraphQLTraditional queries GraphQL using traditional error handling for proper resource management
func queryGraphQLTraditional(
	client *http.Client,
	graphqlURL, geneID string,
) (string, error) {
	request := GraphQLRequest{
		Query: `
			query GeneGeneralInformationSummary($gene: String!) {
				geneGeneralInformation(gene: $gene) {
					id
					description
				}
			}
		`,
		Variables: map[string]any{
			"gene": geneID,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	resp, err := client.Post(
		graphqlURL,
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return "", fmt.Errorf("failed to make GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"GraphQL request failed with status: %d",
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var response GraphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	if len(response.Errors) > 0 {
		return "", fmt.Errorf("GraphQL errors: %s", response.Errors[0].Message)
	}

	if response.Data.GeneGeneralInformation == nil {
		return "", fmt.Errorf("gene not found")
	}

	return response.Data.GeneGeneralInformation.Description, nil
}

// generateReportWithFP creates a validation report using functional programming patterns
func generateReportWithFP(results []ValidationResult) ValidationReport {
	return collection.Pipe2(
		results,
		calculateReportCountsWithFP,
		func(counts ReportCounts) ValidationReport {
			return ValidationReport{
				TotalGenes:    len(results),
				MatchCount:    counts.MatchCount,
				MismatchCount: counts.MismatchCount,
				ErrorCount:    counts.ErrorCount,
				Results:       results,
				GeneratedAt:   time.Now(),
			}
		},
	)
}

// generateReport creates a validation report from the results using functional composition (backward compatibility)
func generateReport(results []ValidationResult) ValidationReport {
	return generateReportWithFP(results)
}

// ReportCounts holds the count statistics for a validation report
type ReportCounts struct {
	MatchCount    int
	MismatchCount int
	ErrorCount    int
}

// calculateReportCountsWithFP calculates statistics using functional programming patterns
func calculateReportCountsWithFP(results []ValidationResult) ReportCounts {
	matchCount := len(collection.Filter(results, func(r ValidationResult) bool {
		return r.Error == "" && r.Match
	}))

	errorCount := len(collection.Filter(results, func(r ValidationResult) bool {
		return r.Error != ""
	}))

	mismatchCount := len(
		collection.Filter(results, func(r ValidationResult) bool {
			return r.Error == "" && !r.Match
		}),
	)

	return ReportCounts{
		MatchCount:    matchCount,
		MismatchCount: mismatchCount,
		ErrorCount:    errorCount,
	}
}

// saveReport saves the validation report to a JSON file with secure permissions
func saveReport(outputPath string, report ValidationReport) error {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	// Fix security issue: use 0600 for secure file permissions
	err = os.WriteFile(outputPath, jsonData, 0o600)
	if err != nil {
		return fmt.Errorf("failed to save report to %s: %w", outputPath, err)
	}

	fmt.Printf("Validation report saved to: %s\n", outputPath)
	printReportSummary(report)
	return nil
}

// printReportSummary prints a summary of the validation report to stdout (removed emojis per guidelines)
func printReportSummary(report ValidationReport) {
	fmt.Printf("\nValidation Summary:\n")
	fmt.Printf("  Total genes processed: %d\n", report.TotalGenes)
	fmt.Printf("  Matches: %d (%.1f%%)\n",
		report.MatchCount,
		calculatePercentage(report.MatchCount, report.TotalGenes))
	fmt.Printf("  Mismatches: %d (%.1f%%)\n",
		report.MismatchCount,
		calculatePercentage(report.MismatchCount, report.TotalGenes))
	fmt.Printf("  Errors: %d (%.1f%%)\n",
		report.ErrorCount,
		calculatePercentage(report.ErrorCount, report.TotalGenes))
	fmt.Printf(
		"  Generated at: %s\n\n",
		report.GeneratedAt.Format(time.RFC3339),
	)

	if report.MismatchCount > 0 {
		printMismatchExamples(report.Results)
	}
}

// calculatePercentage safely calculates percentage avoiding division by zero
func calculatePercentage(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0.0
	}
	return float64(numerator) / float64(denominator) * 100
}

// printMismatchExamples prints the first few mismatches using functional programming
func printMismatchExamples(results []ValidationResult) {
	mismatchResults := collection.Filter(
		results,
		func(r ValidationResult) bool {
			return !r.Match && r.Error == ""
		},
	)

	// Take first 3 mismatches
	exampleCount := min(len(mismatchResults), 3)

	if exampleCount > 0 {
		fmt.Printf("First few mismatches:\n")
		for i := range exampleCount {
			result := mismatchResults[i]
			fmt.Printf("  Gene: %s\n", result.GeneID)
			fmt.Printf(
				"    CSV: %s\n",
				truncateString(result.CSVDescription, 80),
			)
			fmt.Printf(
				"    GQL: %s\n\n",
				truncateString(result.GQLDescription, 80),
			)
		}
	}
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Backward compatibility functions for tests

// queryGraphQLBackwardCompat provides backward compatibility for tests
func queryGraphQLBackwardCompat(
	client *http.Client,
	graphqlURL, geneID string,
) (string, error) {
	// Use traditional approach for simplicity in tests
	request := GraphQLRequest{
		Query: `
			query GeneGeneralInformationSummary($gene: String!) {
				geneGeneralInformation(gene: $gene) {
					id
					description
				}
			}
		`,
		Variables: map[string]any{
			"gene": geneID,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	resp, err := client.Post(
		graphqlURL,
		"application/json",
		strings.NewReader(string(jsonData)),
	)
	if err != nil {
		return "", fmt.Errorf("failed to make GraphQL request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"GraphQL request failed with status: %d",
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var response GraphQLResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse GraphQL response: %w", err)
	}

	if len(response.Errors) > 0 {
		return "", fmt.Errorf("GraphQL errors: %s", response.Errors[0].Message)
	}

	if response.Data.GeneGeneralInformation == nil {
		return "", fmt.Errorf("gene not found: %s", geneID)
	}

	return response.Data.GeneGeneralInformation.Description, nil
}
