package cli

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValidateFileExists(t *testing.T) {
	// Create a test file
	tmpFile, err := os.CreateTemp("", "test_*.csv")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Test with existing file
	fileInfo, err := os.Stat(tmpFile.Name())
	require.NoError(t, err)
	require.False(t, fileInfo.IsDir())

	// Test with non-existent file
	_, err = os.Stat("/non/existent/file.csv")
	require.Error(t, err)
}

func TestReadCSVFile(t *testing.T) {
	tests := []struct {
		name          string
		csvContent    string
		expectedError bool
		expectedRows  int
	}{
		{
			name:          "valid CSV file",
			csvContent:    "GeneID,gene_description\nDDB_G0269114,test description\nDDB_G0278243,another test",
			expectedError: false,
			expectedRows:  3,
		},
		{
			name:          "empty CSV file",
			csvContent:    "",
			expectedError: true, // Now expects error due to stricter validation
			expectedRows:  0,
		},
		{
			name:          "header only CSV",
			csvContent:    "GeneID,gene_description",
			expectedError: true, // Now expects error due to stricter validation
			expectedRows:  1,
		},
		{
			name:          "CSV with quoted fields",
			csvContent:    "GeneID,gene_description\nDDB_G0269114,\"test, description\"\nDDB_G0278243,\"another \"\"test\"\"\"",
			expectedError: false,
			expectedRows:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary CSV file
			tmpFile, err := os.CreateTemp("", "test_*.csv")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.csvContent)
			require.NoError(t, err)
			tmpFile.Close()

			// Test readCSVFile
			records, err := readCSVFile(tmpFile.Name())

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, records, tt.expectedRows)
			}
		})
	}
}

func TestQueryGraphQL(t *testing.T) {
	tests := []struct {
		name           string
		geneID         string
		serverResponse string
		serverStatus   int
		expectedError  bool
		expectedDesc   string
	}{
		{
			name:   "successful query",
			geneID: "DDB_G0269114",
			serverResponse: `{
				"data": {
					"geneGeneralInformation": {
						"id": "DDB_G0269114",
						"description": "test gene description"
					}
				}
			}`,
			serverStatus:  200,
			expectedError: false,
			expectedDesc:  "test gene description",
		},
		{
			name:   "gene not found",
			geneID: "INVALID_GENE",
			serverResponse: `{
				"data": {
					"geneGeneralInformation": null
				}
			}`,
			serverStatus:  200,
			expectedError: true,
			expectedDesc:  "",
		},
		{
			name:   "GraphQL error",
			geneID: "DDB_G0269114",
			serverResponse: `{
				"errors": [
					{"message": "Gene not found"}
				]
			}`,
			serverStatus:  200,
			expectedError: true,
			expectedDesc:  "",
		},
		{
			name:          "server error",
			geneID:        "DDB_G0269114",
			serverStatus:  500,
			expectedError: true,
			expectedDesc:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				if tt.serverStatus == 200 {
					w.Write([]byte(tt.serverResponse))
				}
			}))
			defer server.Close()

			client := &http.Client{Timeout: 5 * time.Second}
			description, err := queryGraphQLBackwardCompat(client, server.URL, tt.geneID)

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedDesc, description)
			}
		})
	}
}

// createMockGraphQLServer creates a test server that responds to GraphQL gene queries
func createMockGraphQLServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request GraphQLRequest
		json.NewDecoder(r.Body).Decode(&request)

		geneID := request.Variables["gene"].(string)

		var response string
		if geneID == "DDB_G0269114" {
			response = `{
				"data": {
					"geneGeneralInformation": {
						"id": "DDB_G0269114",
						"description": "test description"
					}
				}
			}`
		} else {
			response = `{
				"data": {
					"geneGeneralInformation": null
				}
			}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

// createValidationTestCases returns test cases for single gene validation
func createValidationTestCases() []struct {
	name          string
	record        []string
	expectedMatch bool
	expectedError bool
} {
	return []struct {
		name          string
		record        []string
		expectedMatch bool
		expectedError bool
	}{
		{
			name:          "exact match",
			record:        []string{"DDB_G0269114", "test description"},
			expectedMatch: true,
			expectedError: false,
		},
		{
			name:          "case insensitive match",
			record:        []string{"DDB_G0269114", "TEST DESCRIPTION"},
			expectedMatch: true,
			expectedError: false,
		},
		{
			name:          "no match",
			record:        []string{"DDB_G0269114", "different description"},
			expectedMatch: false,
			expectedError: false,
		},
		{
			name:          "gene not found",
			record:        []string{"INVALID_GENE", "some description"},
			expectedMatch: false,
			expectedError: true,
		},
		{
			name:          "invalid record - too few columns",
			record:        []string{"DDB_G0269114"},
			expectedMatch: false,
			expectedError: true,
		},
		{
			name:          "empty gene ID",
			record:        []string{"", "some description"},
			expectedMatch: false,
			expectedError: true,
		},
	}
}

// runSingleValidationTest executes a single validation test case
func runSingleValidationTest(t *testing.T, client *http.Client, serverURL string, testCase struct {
	name          string
	record        []string
	expectedMatch bool
	expectedError bool
},
) {
	params := SingleGeneValidationParams{
		Client:     client,
		GraphQLURL: serverURL,
		Record:     testCase.record,
	}
	result := validateSingleGene(params)

	if testCase.expectedError {
		require.NotEmpty(t, result.Error)
	} else {
		require.Empty(t, result.Error)
		require.Equal(t, testCase.expectedMatch, result.Match)
	}

	if len(testCase.record) > 0 {
		require.Equal(t, testCase.record[0], result.GeneID)
	}
}

func TestValidateSingleGene(t *testing.T) {
	server := createMockGraphQLServer()
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	tests := createValidationTestCases()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runSingleValidationTest(t, client, server.URL, tt)
		})
	}
}

func TestGenerateReport(t *testing.T) {
	results := []ValidationResult{
		{GeneID: "GENE1", Match: true},
		{GeneID: "GENE2", Match: false},
		{GeneID: "GENE3", Error: "some error"},
		{GeneID: "GENE4", Match: true},
		{GeneID: "GENE5", Match: false},
	}

	report := generateReport(results)

	require.Equal(t, 5, report.TotalGenes)
	require.Equal(t, 2, report.MatchCount)
	require.Equal(t, 2, report.MismatchCount)
	require.Equal(t, 1, report.ErrorCount)
	require.Len(t, report.Results, 5)
	require.False(t, report.GeneratedAt.IsZero())
}

func TestSaveReport(t *testing.T) {
	report := ValidationReport{
		TotalGenes:    3,
		MatchCount:    2,
		MismatchCount: 1,
		ErrorCount:    0,
		Results: []ValidationResult{
			{GeneID: "GENE1", Match: true},
			{GeneID: "GENE2", Match: true},
			{GeneID: "GENE3", Match: false},
		},
		GeneratedAt: time.Now(),
	}

	tmpFile, err := os.CreateTemp("", "test_report_*.json")
	require.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = saveReport(tmpFile.Name(), report)
	require.NoError(t, err)

	// Verify file contents
	data, err := os.ReadFile(tmpFile.Name())
	require.NoError(t, err)

	var savedReport ValidationReport
	err = json.Unmarshal(data, &savedReport)
	require.NoError(t, err)

	require.Equal(t, report.TotalGenes, savedReport.TotalGenes)
	require.Equal(t, report.MatchCount, savedReport.MatchCount)
	require.Equal(t, report.MismatchCount, savedReport.MismatchCount)
	require.Equal(t, report.ErrorCount, savedReport.ErrorCount)
	require.Len(t, savedReport.Results, 3)
}

// createTestCSVFile creates a temporary CSV file with test gene data
func createTestCSVFile(t *testing.T) string {
	tmpCSV, err := os.CreateTemp("", "test_genes_*.csv")
	require.NoError(t, err)

	writer := csv.NewWriter(tmpCSV)
	err = writer.WriteAll([][]string{
		{"GeneID", "gene_description"},
		{"DDB_G0269114", "test description"},
		{"DDB_G0278243", "another description"},
		{"INVALID_GENE", "invalid gene description"},
	})
	require.NoError(t, err)
	tmpCSV.Close()

	return tmpCSV.Name()
}

// createIntegrationGraphQLServer creates a test server for integration tests
func createIntegrationGraphQLServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request GraphQLRequest
		json.NewDecoder(r.Body).Decode(&request)

		geneID := request.Variables["gene"].(string)

		var response string
		switch geneID {
		case "DDB_G0269114":
			response = `{
				"data": {
					"geneGeneralInformation": {
						"id": "DDB_G0269114",
						"description": "test description"
					}
				}
			}`
		case "DDB_G0278243":
			response = `{
				"data": {
					"geneGeneralInformation": {
						"id": "DDB_G0278243",
						"description": "different description"
					}
				}
			}`
		default:
			response = `{
				"data": {
					"geneGeneralInformation": null
				}
			}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

// verifyValidationReport checks the generated validation report
func verifyValidationReport(t *testing.T, reportPath string) {
	data, err := os.ReadFile(reportPath)
	require.NoError(t, err)

	var report ValidationReport
	err = json.Unmarshal(data, &report)
	require.NoError(t, err)

	require.Equal(t, 3, report.TotalGenes)
	require.Equal(t, 1, report.MatchCount)    // DDB_G0269114 matches
	require.Equal(t, 1, report.MismatchCount) // DDB_G0278243 doesn't match
	require.Equal(t, 1, report.ErrorCount)    // INVALID_GENE causes error

	// Check specific results
	geneResults := make(map[string]ValidationResult)
	for _, result := range report.Results {
		geneResults[result.GeneID] = result
	}

	require.True(t, geneResults["DDB_G0269114"].Match)
	require.Empty(t, geneResults["DDB_G0269114"].Error)

	require.False(t, geneResults["DDB_G0278243"].Match)
	require.Empty(t, geneResults["DDB_G0278243"].Error)

	require.False(t, geneResults["INVALID_GENE"].Match)
	require.NotEmpty(t, geneResults["INVALID_GENE"].Error)
}

func TestValidateGeneDataIntegration(t *testing.T) {
	csvPath := createTestCSVFile(t)
	defer os.Remove(csvPath)

	server := createIntegrationGraphQLServer()
	defer server.Close()

	tmpReport, err := os.CreateTemp("", "test_report_*.json")
	require.NoError(t, err)
	tmpReport.Close()
	defer os.Remove(tmpReport.Name())

	// Test the validation function directly with HTTP URL for integration test
	params := GeneValidationParams{
		InputFile:    csvPath,
		OutputReport: tmpReport.Name(),
		GraphQLURL:   server.URL, // This will be HTTP for test server
		Timeout:      10,
		Workers:      2,
	}

	// Read CSV file
	records, err := readCSVFile(csvPath)
	require.NoError(t, err)

	// Validate gene descriptions directly (bypassing CLI validation)
	report, err := validateGeneDescriptionsWithFP(params, records)
	require.NoError(t, err)

	// Save report
	err = saveReport(tmpReport.Name(), report)
	require.NoError(t, err)

	verifyValidationReport(t, tmpReport.Name())
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "short string unchanged",
			input:    "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "exact length unchanged",
			input:    "exactly10c",
			maxLen:   10,
			expected: "exactly10c",
		},
		{
			name:     "long string truncated",
			input:    "this is a very long string that should be truncated",
			maxLen:   20,
			expected: "this is a very lo...",
		},
		{
			name:     "empty string",
			input:    "",
			maxLen:   5,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateString(tt.input, tt.maxLen)
			require.Equal(t, tt.expected, result)
		})
	}
}
