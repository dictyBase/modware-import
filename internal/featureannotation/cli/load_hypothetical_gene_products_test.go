package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadGeneIDsFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expected    []string
		shouldError bool
	}{
		{
			name: "valid gene IDs",
			fileContent: `DDB_G0288781
DDB_G0276393
DDB_G0274173`,
			expected:    []string{"DDB_G0288781", "DDB_G0276393", "DDB_G0274173"},
			shouldError: false,
		},
		{
			name: "gene IDs with empty lines",
			fileContent: `DDB_G0288781

DDB_G0276393

DDB_G0274173`,
			expected:    []string{"DDB_G0288781", "DDB_G0276393", "DDB_G0274173"},
			shouldError: false,
		},
		{
			name: "gene IDs with comments",
			fileContent: `# This is a comment
DDB_G0288781
# Another comment
DDB_G0276393`,
			expected:    []string{"DDB_G0288781", "DDB_G0276393"},
			shouldError: false,
		},
		{
			name:        "gene IDs with BOM and whitespace",
			fileContent: "\ufeffDDB_G0288781\n  DDB_G0276393  \n\tDDB_G0274173\t",
			expected:    []string{"\ufeffDDB_G0288781", "DDB_G0276393", "DDB_G0274173"},
			shouldError: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			expected:    nil,
			shouldError: true,
		},
		{
			name:        "only comments and empty lines",
			fileContent: "# Comment\n\n# Another comment\n",
			expected:    nil,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "genes.txt")
			err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0o600)
			require.NoError(t, err)

			// Execute using the IOEither pipeline and convert to tuple
			result := toEither(readGeneIDsFromFile(tmpFile))
			geneIDs, err := eitherToTuple(result)

			// Assert
			if tt.shouldError {
				require.Error(t, err)
				require.Nil(t, geneIDs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, geneIDs)
			}
		})
	}
}

func TestHypotheticalProteinProduct(t *testing.T) {
	t.Run("constant is correct", func(t *testing.T) {
		require.Equal(t, "conserved hypothetical protein", HypotheticalProteinProduct)
	})
}
