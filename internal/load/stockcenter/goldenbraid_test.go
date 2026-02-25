package stockcenter

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/either"
	O "github.com/IBM/fp-go/option"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchPlasmidByName_ReturnsNoneWhenNotFound(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)

	// Mock: ListPlasmids returns empty collection
	mockClient.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Limit == 1 && p.Filter == "plasmid_name===nonexistent"
	}), mock.Anything).
		Return(&stock.PlasmidCollection{Data: []*stock.PlasmidCollection_Data{}}, nil)

	result := fetchPlasmidByName("nonexistent")()

	require.True(t, E.IsRight(result))
	opt := E.GetOrElse(func(error) O.Option[*stock.Plasmid] {
		return O.None[*stock.Plasmid]()
	})(result)
	require.True(t, O.IsNone(opt))
	mockClient.AssertExpectations(t)
}

func TestFetchPlasmidByName_ReturnsSomeWhenFound(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)

	// Mock: ListPlasmids returns collection with one plasmid
	mockClient.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Limit == 1 && p.Filter == "plasmid_name===pTest1"
	}), mock.Anything).
		Return(&stock.PlasmidCollection{
			Data: []*stock.PlasmidCollection_Data{
				{
					Type: "plasmid",
					Id:   "DBP0001",
					Attributes: &stock.PlasmidAttributes{
						Name: "pTest1",
					},
				},
			},
		}, nil)

	result := fetchPlasmidByName("pTest1")()

	require.True(t, E.IsRight(result))
	opt := E.GetOrElse(func(error) O.Option[*stock.Plasmid] {
		return O.None[*stock.Plasmid]()
	})(result)
	require.True(t, O.IsSome(opt))
	plasmid := O.GetOrElse(func() *stock.Plasmid { return nil })(opt)
	require.Equal(t, "DBP0001", plasmid.Data.Id)
	mockClient.AssertExpectations(t)
}

func TestProcessPlasmid_CreateWhenNotExists(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)

	ctx := &source.GoldenBraidContext{
		Name:         "pNew",
		Summary:      "A test plasmid",
		User:         "test@example.com",
		PlasmidType:  "plasmid",
		Depositor:    "depositor@example.com",
		Genes:        O.None[[]string](),
		Publications: O.None[[]string](),
	}

	mockClient.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Limit == 1 && p.Filter == "plasmid_name===pNew"
	}), mock.Anything).
		Return(&stock.PlasmidCollection{Data: []*stock.PlasmidCollection_Data{}}, nil)
	mockClient.On("CreatePlasmid", mock.Anything, mock.MatchedBy(func(p *stock.NewPlasmid) bool {
		return p.Data.Attributes.Name == "pNew"
	}), mock.Anything).
		Return(&stock.Plasmid{Data: &stock.Plasmid_Data{Id: "DBPNew"}}, nil)

	result := processPlasmid(ctx)

	require.True(t, E.IsRight(result))
	processResult := E.GetOrElse(func(error) GoldenBraidResult {
		return GoldenBraidResult{}
	})(result)
	require.True(t, processResult.Created)
	require.Equal(t, "DBPNew", processResult.PlasmidID)
	mockClient.AssertExpectations(t)
}

func TestProcessPlasmid_SkipWhenExists(t *testing.T) {
	mockClient := new(MockStockClient)
	regsc.SetStockAPIClient(mockClient)

	ctx := &source.GoldenBraidContext{
		Name:         "pExisting",
		Summary:      "An existing plasmid",
		User:         "test@example.com",
		PlasmidType:  "plasmid",
		Depositor:    "depositor@example.com",
		Genes:        O.None[[]string](),
		Publications: O.None[[]string](),
	}

	mockClient.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Limit == 1 && p.Filter == "plasmid_name===pExisting"
	}), mock.Anything).
		Return(&stock.PlasmidCollection{
			Data: []*stock.PlasmidCollection_Data{
				{
					Type: "plasmid",
					Id:   "DBPExisting",
					Attributes: &stock.PlasmidAttributes{
						Name: "pExisting",
					},
				},
			},
		}, nil)

	result := processPlasmid(ctx)

	require.True(t, E.IsRight(result))
	processResult := E.GetOrElse(func(error) GoldenBraidResult {
		return GoldenBraidResult{}
	})(result)
	require.False(t, processResult.Created)
	require.Equal(t, "DBPExisting", processResult.PlasmidID)
	mockClient.AssertExpectations(t)
}

func TestGoldenBraidSummaryAggregation(t *testing.T) {
	result1 := GoldenBraidResult{PlasmidID: "DBP0001", Created: true}
	result2 := GoldenBraidResult{PlasmidID: "DBP0002", Created: false}
	result3 := GoldenBraidResult{Error: errors.New("test error")}

	summary := GoldenBraidProcessingResult{}
	semigroup := GoldenBraidSummarySemigroup()

	summary = semigroup.Concat(summary, goldenBraidResultToSummary(result1))
	summary = semigroup.Concat(summary, goldenBraidResultToSummary(result2))
	summary = semigroup.Concat(summary, goldenBraidResultToSummary(result3))

	require.Equal(t, 1, summary.CreatedCount)
	require.Equal(t, 1, summary.SkippedCount)
	require.Equal(t, 1, summary.ErrorCount)
	require.Len(t, summary.Successes, 2)
}
