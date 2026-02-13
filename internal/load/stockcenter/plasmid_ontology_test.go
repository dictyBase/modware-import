package stockcenter

import (
	"errors"
	"testing"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test helpers
func createTestPlasmidData(id string, property string) *stockpb.PlasmidCollection_Data {
	return &stockpb.PlasmidCollection_Data{
		Id:   id,
		Type: "plasmid",
		Attributes: &stockpb.PlasmidAttributes{
			DictyPlasmidProperty: property,
		},
	}
}

// TestProcessBatch tests the processBatch function
// Note: Server-side filtering now excludes plasmids with target term or "GB vector"
// before they reach processBatch, so all plasmids in the batch should be updated
func TestProcessBatch(t *testing.T) {
	mockClient := new(MockStockClient)

	// These plasmids already passed server-side filtering
	plasmids := []*stockpb.PlasmidCollection_Data{
		createTestPlasmidData("DBP001", "expression vector"),
		createTestPlasmidData("DBP002", "cloning vector"),
	}

	// Mock: UpdatePlasmid called for both plasmids
	mockClient.On("UpdatePlasmid", mock.Anything, mock.Anything, mock.Anything).
		Return(&stockpb.Plasmid{}, nil).Twice()

	stats := processBatch(plasmids, mockClient, "vector")

	require.Equal(t, 2, stats.ProcessedCount)
	require.Equal(t, 2, stats.UpdatedCount)
	require.Equal(t, 0, stats.ErrorCount)

	mockClient.AssertExpectations(t)
}

// TestProcessBatchWithErrors tests processBatch when updates fail
func TestProcessBatchWithErrors(t *testing.T) {
	mockClient := new(MockStockClient)

	plasmids := []*stockpb.PlasmidCollection_Data{
		createTestPlasmidData("DBP001", "expression vector"), // Update - will succeed
		createTestPlasmidData("DBP002", "cloning vector"),    // Update - will fail
	}

	// Mock: DBP001 succeeds
	mockClient.On("UpdatePlasmid", mock.Anything, mock.MatchedBy(
		func(req *stockpb.PlasmidUpdate) bool {
			return req.Data.Id == "DBP001"
		},
	), mock.Anything).Return(&stockpb.Plasmid{}, nil).Once()

	// Mock: DBP002 fails
	mockClient.On("UpdatePlasmid", mock.Anything, mock.MatchedBy(
		func(req *stockpb.PlasmidUpdate) bool {
			return req.Data.Id == "DBP002"
		},
	), mock.Anything).Return(nil, errors.New("update failed")).Once()

	stats := processBatch(plasmids, mockClient, "vector")

	require.Equal(t, 2, stats.ProcessedCount)
	require.Equal(t, 1, stats.UpdatedCount)
	require.Equal(t, 1, stats.ErrorCount)
	require.Len(t, stats.Errors, 1)

	mockClient.AssertExpectations(t)
}

// TestStatsSemigroup tests the StatsSemigroup combinator
func TestStatsSemigroup(t *testing.T) {
	semigroup := StatsSemigroup()

	stats1 := OntologyUpdateStats{
		ProcessedCount: 5,
		UpdatedCount:   3,
		ErrorCount:     2,
		Errors:         []error{errors.New("error 1")},
	}

	stats2 := OntologyUpdateStats{
		ProcessedCount: 10,
		UpdatedCount:   7,
		ErrorCount:     3,
		Errors:         []error{errors.New("error 2")},
	}

	combined := semigroup.Concat(stats1, stats2)

	require.Equal(t, 15, combined.ProcessedCount)
	require.Equal(t, 10, combined.UpdatedCount)
	require.Equal(t, 5, combined.ErrorCount)
	require.Len(t, combined.Errors, 2)
}

// TestStatsSemigroupErrorLimit tests that error sampling is limited
func TestStatsSemigroupErrorLimit(t *testing.T) {
	semigroup := StatsSemigroup()

	// Create stats with more than maxErrorSamples errors
	stats1 := OntologyUpdateStats{
		Errors: []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		},
	}

	stats2 := OntologyUpdateStats{
		Errors: []error{
			errors.New("error 4"),
			errors.New("error 5"),
			errors.New("error 6"),
		},
	}

	combined := semigroup.Concat(stats1, stats2)

	// Should be limited to maxErrorSamples (5)
	require.LessOrEqual(t, len(combined.Errors), maxErrorSamples)
	require.Equal(t, 5, len(combined.Errors))
}

// TestCountPlasmids tests the countPlasmids function
func TestCountPlasmids(t *testing.T) {
	plasmids := []*stockpb.PlasmidCollection_Data{
		createTestPlasmidData("DBP001", "vector"),
		createTestPlasmidData("DBP002", "expression vector"),
		createTestPlasmidData("DBP003", "GB vector"),
	}

	count := countPlasmids(plasmids)
	require.Equal(t, 3, count)
}

// TestCreateInitialStats tests the createInitialStats function
func TestCreateInitialStats(t *testing.T) {
	stats := createInitialStats(10)
	require.Equal(t, 10, stats.ProcessedCount)
	require.Equal(t, 0, stats.UpdatedCount)
	require.Equal(t, 0, stats.ErrorCount)
	require.Empty(t, stats.Errors)
}
