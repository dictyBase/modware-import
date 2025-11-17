package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	E "github.com/IBM/fp-go/either"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/stretchr/testify/require"
)

// eitherToTuple converts an Either to a (value, error) tuple
// Test helper for asserting Either results
func eitherToTuple[A any](either E.Either[error, A]) (A, error) {
	var zero A
	if E.IsLeft(either) {
		// Extract error from Left side
		leftEither := E.Swap(either)
		err := E.GetOrElse(
			func(A) error { return fmt.Errorf("unknown error") },
		)(
			leftEither,
		)
		return zero, err
	}
	// Extract value from Right side
	value := E.GetOrElse(func(error) A { return zero })(either)
	return value, nil
}

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
			expected:    []string{"DDB_G0288781", "DDB_G0276393", "DDB_G0274173"},
			shouldError: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			expected:    []string{},
			shouldError: false,
		},
		{
			name:        "only comments and empty lines",
			fileContent: "# Comment\n\n# Another comment\n",
			expected:    []string{},
			shouldError: false,
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
			result := fputil.ToEither(readGeneIDsFromFile(tmpFile))
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
