package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseUnknowmeDataParams_Validation(t *testing.T) {
	tests := []struct {
		name        string
		params      ParseUnknowmeDataParams
		expectError bool
	}{
		{
			name: "valid parameters with multiple files",
			params: ParseUnknowmeDataParams{
				InputFiles:            []string{"file1.html", "file2.html"},
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
			if tt.expectError && err == nil {
				t.Errorf("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
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

			if result.GeneID != tt.expected.GeneID {
				t.Errorf("expected GeneID %s, got %s", tt.expected.GeneID, result.GeneID)
			}
			if result.GeneProduct != tt.expected.GeneProduct {
				t.Errorf("expected GeneProduct %s, got %s", tt.expected.GeneProduct, result.GeneProduct)
			}
			if result.Description != tt.expected.Description {
				t.Errorf("expected Description %s, got %s", tt.expected.Description, result.Description)
			}
		})
	}
}

func TestValidateInputFilesExist(t *testing.T) {
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "test_parse_unknowme")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFile1 := filepath.Join(tempDir, "test1.html")
	testFile2 := filepath.Join(tempDir, "test2.html")

	if err := os.WriteFile(testFile1, []byte("<html></html>"), 0o600); err != nil {
		t.Fatalf("failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("<html></html>"), 0o600); err != nil {
		t.Fatalf("failed to create test file 2: %v", err)
	}

	tests := []struct {
		name        string
		inputFiles  []string
		expectError bool
	}{
		{
			name:        "all files exist",
			inputFiles:  []string{testFile1, testFile2},
			expectError: false,
		},
		{
			name:        "one file missing",
			inputFiles:  []string{testFile1, "nonexistent.html"},
			expectError: true,
		},
		{
			name:        "empty file list",
			inputFiles:  []string{},
			expectError: false, // validateInputFilesExist doesn't validate empty lists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputFilesExist(tt.inputFiles)
			if tt.expectError && err == nil {
				t.Errorf("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}
