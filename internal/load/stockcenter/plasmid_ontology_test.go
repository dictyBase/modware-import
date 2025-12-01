package stockcenter

import (
	"errors"
	"os"
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSummarySemigroup(t *testing.T) {
	// Create two summaries with errors
	s1 := KeywordProcessingSummary{
		ErrorCount: 1,
		Errors:     []string{"error 1"},
	}
	s2 := KeywordProcessingSummary{
		ErrorCount: 1,
		Errors:     []string{"error 2"},
	}

	sg := SummarySemigroup()
	result := sg.Concat(s1, s2)

	require.Equal(t, 2, result.ErrorCount)
	require.Len(t, result.Errors, 2)
	require.Contains(t, result.Errors, "error 1")
	require.Contains(t, result.Errors, "error 2")
}

func TestLoadPlasmidOntology(t *testing.T) {
	t.Run("Success", testLoadPlasmidOntologySuccess)
	t.Run("Skip Existing", testLoadPlasmidOntologySkipExisting)
	t.Run("Validation Error", testLoadPlasmidOntologyValidationError)
	t.Run("API Error", testLoadPlasmidOntologyAPIError)
}

func testLoadPlasmidOntologySuccess(t *testing.T) {
	// Mock setup
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)
	defer func() { regsc.SetStockAPIClient(nil) }()

	// Test data
	csvContent := "DBP0000001\tkeyword\tvector\n"
	filePath := createTempCSV(t, csvContent)
	defer os.Remove(filePath)

	// Mock expectations
	mockClient.On(
		"GetPlasmid",
		mock.Anything,
		&stock.StockId{Id: "DBP0000001"},
		mock.Anything,
	).
		Return(&stock.Plasmid{
			Data: &stock.Plasmid_Data{
				Attributes: &stock.PlasmidAttributes{
					DictyPlasmidProperty: "", // Current property is empty
				},
			},
		}, nil)

	mockClient.On(
		"UpdatePlasmid",
		mock.Anything,
		mock.MatchedBy(func(req *stock.PlasmidUpdate) bool {
			return req.Data.Id == "DBP0000001" &&
				req.Data.Attributes.DictyPlasmidProperty == "vector"
		}),
		mock.Anything,
	).Return(&stock.Plasmid{}, nil)

	// Command setup
	cmd := &cobra.Command{}
	cmd.Flags().String("input", filePath, "")
	cmd.Flags().String("property", "keyword", "")

	// Execute
	err := LoadPlasmidOntology(cmd, []string{})
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func testLoadPlasmidOntologySkipExisting(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)
	defer func() { regsc.SetStockAPIClient(nil) }()

	csvContent := "DBP0000002\tkeyword\tvector\n"
	filePath := createTempCSV(t, csvContent)
	defer os.Remove(filePath)

	// Should return existing matching property
	mockClient.On(
		"GetPlasmid",
		mock.Anything,
		&stock.StockId{Id: "DBP0000002"},
		mock.Anything,
	).
		Return(&stock.Plasmid{
			Data: &stock.Plasmid_Data{
				Attributes: &stock.PlasmidAttributes{
					DictyPlasmidProperty: "vector",
				},
			},
		}, nil)

	// UpdatePlasmid should NOT be called

	cmd := &cobra.Command{}
	cmd.Flags().String("input", filePath, "")
	cmd.Flags().String("property", "keyword", "")

	err := LoadPlasmidOntology(cmd, []string{})
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}

func testLoadPlasmidOntologyValidationError(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)
	defer func() { regsc.SetStockAPIClient(nil) }()

	// Invalid CSV (only 2 columns)
	csvContent := "DBP0000003\tkeyword\n"
	filePath := createTempCSV(t, csvContent)
	defer os.Remove(filePath)

	cmd := &cobra.Command{}
	cmd.Flags().String("input", filePath, "")
	cmd.Flags().String("property", "keyword", "")

	// The function logs errors but returns nil for partial successes/failures
	err := LoadPlasmidOntology(cmd, []string{})
	require.NoError(t, err)

	// No API calls expected
	mockClient.AssertExpectations(t)
}

func testLoadPlasmidOntologyAPIError(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)
	defer func() { regsc.SetStockAPIClient(nil) }()

	csvContent := "DBP0000004\tkeyword\tvector\n"
	filePath := createTempCSV(t, csvContent)
	defer os.Remove(filePath)

	mockClient.On(
		"GetPlasmid",
		mock.Anything,
		&stock.StockId{Id: "DBP0000004"},
		mock.Anything,
	).
		Return(nil, errors.New("network error"))

	cmd := &cobra.Command{}
	cmd.Flags().String("input", filePath, "")
	cmd.Flags().String("property", "keyword", "")

	err := LoadPlasmidOntology(cmd, []string{})
	require.NoError(t, err)

	mockClient.AssertExpectations(t)
}

// Helper to create temp file
func createTempCSV(t *testing.T, content string) string {
	f, err := os.CreateTemp("", "plasmid_ontology_*.tsv")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
	return f.Name()
}
