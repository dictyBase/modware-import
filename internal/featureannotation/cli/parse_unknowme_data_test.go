package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseUnknowmeDataParams_Validation(t *testing.T) {
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "test_validation")
	require.NoError(t, err, "failed to create temp directory")
	defer os.RemoveAll(tempDir)

	// Create test files
	validFile1 := filepath.Join(tempDir, "file1.html")
	validFile2 := filepath.Join(tempDir, "file2.html")

	err = os.WriteFile(validFile1, []byte("<html></html>"), 0o600)
	require.NoError(t, err, "failed to create test file 1")

	err = os.WriteFile(validFile2, []byte("<html></html>"), 0o600)
	require.NoError(t, err, "failed to create test file 2")

	tests := []struct {
		name        string
		params      ParseUnknowmeDataParams
		expectError bool
	}{
		{
			name: "valid parameters with multiple files",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{validFile1, validFile2},
				GeneProductOutput:     "products.csv",
				GeneDescriptionOutput: "descriptions.csv",
			},
			expectError: false,
		},
		{
			name: "invalid parameters with empty input files",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{},
				GeneProductOutput:     "products.csv",
				GeneDescriptionOutput: "descriptions.csv",
			},
			expectError: true,
		},
		{
			name: "invalid parameters with empty output file",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{"file1.html"},
				GeneProductOutput:     "",
				GeneDescriptionOutput: "descriptions.csv",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.params)
			if tt.expectError {
				require.Error(t, err, "expected validation error but got none")
			} else {
				require.NoError(t, err, "unexpected validation error")
			}
		})
	}
}

func TestMergeGeneDataRecords(t *testing.T) {
	tests := []struct {
		name     string
		existing GeneDataRecord
		newRec   GeneDataRecord
		expected GeneDataRecord
	}{
		{
			name: "merge with empty existing gene product",
			existing: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "",
				Description: "existing description",
			},
			newRec: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "new gene product",
				Description: "",
			},
			expected: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "new gene product",
				Description: "existing description",
			},
		},
		{
			name: "prefer existing non-empty values",
			existing: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "existing product",
				Description: "existing description",
			},
			newRec: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "new product",
				Description: "new description",
			},
			expected: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "existing product",
				Description: "existing description",
			},
		},
		{
			name: "fill empty description from new record",
			existing: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "existing product",
				Description: "",
			},
			newRec: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "",
				Description: "new description",
			},
			expected: GeneDataRecord{
				GeneID:      "DDB_G0123456",
				GeneProduct: "existing product",
				Description: "new description",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeGeneDataRecords(tt.existing, tt.newRec)

			require.Equal(t, tt.expected.GeneID, result.GeneID, "GeneID should match expected value")
			require.Equal(t, tt.expected.GeneProduct, result.GeneProduct, "GeneProduct should match expected value")
			require.Equal(t, tt.expected.Description, result.Description, "Description should match expected value")
		})
	}
}

func TestValidateInputFilesExist(t *testing.T) {
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "test_parse_unknowme")
	require.NoError(t, err, "failed to create temp directory")
	defer os.RemoveAll(tempDir)

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.html")
	testFile2 := filepath.Join(tempDir, "test2.html")

	err = os.WriteFile(testFile1, []byte("<html></html>"), 0o600)
	require.NoError(t, err, "failed to create test file 1")

	err = os.WriteFile(testFile2, []byte("<html></html>"), 0o600)
	require.NoError(t, err, "failed to create test file 2")

	tests := []struct {
		name        string
		params      ParseUnknowmeDataParams
		expectError bool
	}{
		{
			name: "all files exist",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{testFile1, testFile2},
				GeneProductOutput:     "products.csv",
				GeneDescriptionOutput: "descriptions.csv",
			},
			expectError: false,
		},
		{
			name: "one file missing",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{testFile1, "nonexistent.html"},
				GeneProductOutput:     "products.csv",
				GeneDescriptionOutput: "descriptions.csv",
			},
			expectError: true,
		},
		{
			name: "empty file list",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{},
				GeneProductOutput:     "products.csv",
				GeneDescriptionOutput: "descriptions.csv",
			},
			expectError: true, // empty file list should fail validation
		},
		// TODO(human): Add more test cases for file validation
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.params)
			if tt.expectError {
				require.Error(t, err, "expected validation error but got none")
			} else {
				require.NoError(t, err, "unexpected validation error")
			}
		})
	}
}
