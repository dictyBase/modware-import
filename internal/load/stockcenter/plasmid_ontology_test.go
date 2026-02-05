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

// TestShouldSkip tests the shouldSkip predicate
func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		property string
		want     bool
	}{
		{"GB vector should skip", "vector", "GB vector", true},
		{"target term should skip", "vector", "vector", true},
		{"case insensitive GB vector", "vector", "gb VECTOR", true},
		{"case insensitive target", "vector", "VECTOR", true},
		{"other value should process", "vector", "expression vector", false},
		{"empty should process", "vector", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pctx := ProcessContext{
				Term:    tt.term,
				Plasmid: createTestPlasmidData("DBP001", tt.property),
			}
			require.Equal(t, tt.want, shouldSkip(pctx))
		})
	}
}

// TestHasGBVector tests the hasGBVector predicate
func TestHasGBVector(t *testing.T) {
	tests := []struct {
		name     string
		property string
		want     bool
	}{
		{"exact match", "GB vector", true},
		{"case insensitive", "gb VECTOR", true},
		{"different value", "expression vector", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pctx := ProcessContext{
				Plasmid: createTestPlasmidData("DBP001", tt.property),
			}
			require.Equal(t, tt.want, hasGBVector(pctx))
		})
	}
}

// TestHasProperty tests the hasProperty function
func TestHasProperty(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		property string
		want     bool
	}{
		{"exact match", "vector", "vector", true},
		{"case insensitive match", "vector", "VECTOR", true},
		{"no match", "vector", "expression vector", false},
		{"empty no match", "vector", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pctx := ProcessContext{
				Term:    tt.target,
				Plasmid: createTestPlasmidData("DBP001", tt.property),
			}
			require.Equal(t, tt.want, hasProperty(pctx))
		})
	}
}

// TestProcessBatch tests the processBatch function
func TestProcessBatch(t *testing.T) {
	mockClient := new(MockStockClient)

	plasmids := []*stockpb.PlasmidCollection_Data{
		createTestPlasmidData("DBP001", "GB vector"),         // Skip (has GB vector)
		createTestPlasmidData("DBP002", "vector"),            // Skip (already has target term)
		createTestPlasmidData("DBP003", "expression vector"), // Update this one
	}

	// Mock: UpdatePlasmid called only for DBP003
	mockClient.On("UpdatePlasmid", mock.Anything, mock.MatchedBy(
		func(req *stockpb.PlasmidUpdate) bool {
			return req.Data.Id == "DBP003"
		},
	), mock.Anything).Return(&stockpb.Plasmid{}, nil).Once()

	stats := processBatch(plasmids, mockClient, "vector")

	require.Equal(t, 3, stats.ProcessedCount)
	require.Equal(t, 1, stats.UpdatedCount)
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
