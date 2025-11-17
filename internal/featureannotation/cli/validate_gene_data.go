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
	"slices"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/hasura/go-graphql-client"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	defaultHTTPTimeoutSeconds  = 30       // Default HTTP client timeout
	maxIdleConns               = 10       // Maximum idle connections
	maxIdleConnsPerHost        = 10       // Maximum idle connections per host
	idleConnTimeoutSeconds     = 30       // Idle connection timeout
	maxMismatchExamples        = 3        // Maximum mismatch examples to display
	fullValidationPercentage   = 100      // Full validation percentage
	bytesPerKilobyte           = 1024     // Bytes per kilobyte
	maxFileSizeBytes           = 50       // Maximum file size in MB
	maxRecordCount             = 100000   // Maximum number of records allowed
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
	Client  *graphql.Client
	Record  []string
	Context context.Context
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
	MaxFileSizeBytes: maxFileSizeBytes * bytesPerKilobyte * bytesPerKilobyte, // 50MB
	MaxRecords:       maxRecordCount,                                          // 100k records
}

// GeneGeneralInformationQuery represents the GraphQL query structure for retrieving
// gene information including ID and description fields.
type GeneGeneralInformationQuery struct {
	GeneGeneralInformation *struct {
		ID          string `graphql:"id"`
		Description string `graphql:"description"`
	} `graphql:"geneGeneralInformation(gene: $gene)"`
}

// GraphQLClientConfig holds configuration options for creating a GraphQL client,
// including URL, HTTP client, and custom headers.
type GraphQLClientConfig struct {
	URL        string
	HTTPClient *http.Client
	Headers    map[string]string
}

// GraphQLClientOption is a functional option for configuring GraphQL client
type GraphQLClientOption func(*GraphQLClientConfig)

// WithHTTPClient sets the HTTP client for GraphQL operations
func WithHTTPClient(client *http.Client) GraphQLClientOption {
	return func(cfg *GraphQLClientConfig) {
		cfg.HTTPClient = client
	}
}

// WithHeaders sets custom headers for GraphQL requests
func WithHeaders(headers map[string]string) GraphQLClientOption {
	return func(cfg *GraphQLClientConfig) {
		cfg.Headers = headers
	}
}

// WithTimeout sets the timeout for HTTP requests
func WithTimeout(timeout time.Duration) GraphQLClientOption {
	return func(cfg *GraphQLClientConfig) {
		if cfg.HTTPClient == nil {
			cfg.HTTPClient = &http.Client{}
		}
		cfg.HTTPClient.Timeout = timeout
	}
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

	records, err := readCSVFileWithValidation(params.InputFile)
	if err != nil {
		return fmt.Errorf("validation process failed: %w", err)
	}

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

// readCSVFileWithValidation reads CSV file with security constraints and
// validation
func readCSVFileWithValidation(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file %s: %w", filePath, err)
	}
	defer file.Close()

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

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	if len(records) == 1 {
		return nil, fmt.Errorf("CSV file contains only header row")
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
	if len(records) == 0 {
		return ValidationReport{}, fmt.Errorf("CSV file is empty")
	}

	// Skip header row and filter out empty records using slices.DeleteFunc
	dataRecords := slices.DeleteFunc(
		records[1:],
		func(record []string) bool {
			return len(record) == 0 ||
				(len(record) == 1 &&
					strings.TrimSpace(record[0]) == "")
		},
	)

	if len(dataRecords) == 0 {
		return ValidationReport{}, fmt.Errorf(
			"no data records to process after skipping header",
		)
	}

	// Process records concurrently with errgroup
	results, err := processRecordsConcurrently(
		context.Background(),
		params,
		dataRecords,
	)
	if err != nil {
		return ValidationReport{}, err
	}

	// Generate report using functional programming
	report := generateReport(results)
	return report, nil
}

func processRecordsConcurrently(
	ctx context.Context,
	params GeneValidationParams,
	records [][]string,
) ([]ValidationResult, error) {
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(params.Workers)

	results := make([]ValidationResult, len(records))
	sharedClient := graphql.NewClient(params.GraphQLURL, nil)

	for i, record := range records {
		i, record := i, record // Capture loop variables
		g.Go(func() error {
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
				results[i] = validateSingleGene(
					SingleGeneValidationParams{
						Client:  sharedClient,
						Record:  record,
						Context: gctx,
					})
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("concurrent validation failed: %w", err)
	}

	return results, nil
}

// NewGraphQLClient creates a new GraphQL client with the given options
func NewGraphQLClient(
	graphqlURL string,
	opts ...GraphQLClientOption,
) *graphql.Client {
	cfg := &GraphQLClientConfig{
		URL:     graphqlURL,
		Headers: make(map[string]string),
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Set default HTTP client if none provided
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: defaultHTTPTimeoutSeconds * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        maxIdleConns,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
				IdleConnTimeout:     idleConnTimeoutSeconds * time.Second,
			},
		}
	}

	client := graphql.NewClient(cfg.URL, cfg.HTTPClient)

	// Add custom headers if provided
	if len(cfg.Headers) > 0 {
		client = client.WithRequestModifier(func(req *http.Request) {
			for key, value := range cfg.Headers {
				req.Header.Set(key, value)
			}
		})
	}

	return client
}

// validateSingleGene validates a single gene record against GraphQL endpoint
func validateSingleGene(params SingleGeneValidationParams) ValidationResult {
	if len(params.Record) < config.MinimumFieldCount {
		geneID := ""
		if len(params.Record) > 0 {
			geneID = params.Record[0]
		}
		return ValidationResult{
			GeneID: geneID,
			Error:  "invalid record: insufficient columns (expected at least 2)",
		}
	}

	geneID := params.Record[0]
	csvDesc := params.Record[1]

	// Query GraphQL for gene description
	var query GeneGeneralInformationQuery
	err := params.Client.Query(
		params.Context,
		&query,
		map[string]any{"gene": geneID},
	)
	if err != nil {
		return ValidationResult{
			GeneID:         geneID,
			CSVDescription: csvDesc,
			Error: fmt.Errorf("failed to execute GraphQL query: %w", err).
				Error(),
		}
	}
	if query.GeneGeneralInformation == nil {
		return ValidationResult{
			GeneID:         geneID,
			CSVDescription: csvDesc,
			Error:          fmt.Sprintf("gene not found: %s", geneID),
		}
	}

	gqlDescription := strings.TrimSpace(
		query.GeneGeneralInformation.Description,
	)

	return ValidationResult{
		GeneID:         geneID,
		CSVDescription: csvDesc,
		GQLDescription: gqlDescription,
		Match: strings.EqualFold(
			csvDesc,
			gqlDescription,
		),
	}
}

// generateReport creates a validation report from the results
func generateReport(results []ValidationResult) ValidationReport {
	var matchCount, mismatchCount, errorCount int

	for _, result := range results {
		switch {
		case result.Error != "":
			errorCount++
		case result.Match:
			matchCount++
		default:
			mismatchCount++
		}
	}

	return ValidationReport{
		TotalGenes:    len(results),
		MatchCount:    matchCount,
		MismatchCount: mismatchCount,
		ErrorCount:    errorCount,
		Results:       results,
		GeneratedAt:   time.Now(),
	}
}

// saveReport saves the validation report to a JSON file with secure permissions
func saveReport(outputPath string, report ValidationReport) error {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	// Use secure file permissions (owner read/write only)
	err = os.WriteFile(outputPath, jsonData, config.DefaultFilePermission)
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
	return float64(numerator) / float64(denominator) * fullValidationPercentage
}

// printMismatchExamples prints the first few mismatches
func printMismatchExamples(results []ValidationResult) {
	var mismatches []ValidationResult
	for _, result := range results {
		if !result.Match && result.Error == "" {
			mismatches = append(mismatches, result)
		}
	}

	maxExamples := min(len(mismatches), maxMismatchExamples)
	if maxExamples == 0 {
		return
	}
	examples := mismatches[:maxExamples]

	if len(examples) > 0 {
		fmt.Printf("First few mismatches:\n")
		for _, result := range examples {
			fmt.Printf("  Gene: %s\n", result.GeneID)
			fmt.Printf(
				"    CSV: %s\n",
				truncateString(result.CSVDescription, config.DefaultLineWidth),
			)
			fmt.Printf(
				"    GQL: %s\n\n",
				truncateString(result.GQLDescription, config.DefaultLineWidth),
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
