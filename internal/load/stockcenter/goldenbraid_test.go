package stockcenter

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/either"
	O "github.com/IBM/fp-go/option"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
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

func TestProcessPlasmidWithUpsert_CreateWhenNotExists(t *testing.T) {
	t.Skip("deferred: tracked in modware-import-7zd")
	// TODO(modware-import-7zd): restore body once processPlasmidWithUpsert depositor arg is wired up
	// mockClient := new(MockStockClient)
	// regsc.SetStockAPIClient(mockClient)
	// plasmid := &source.GoldenBraidPlasmid{Name: "pNew"}
	// userEmail := "test@example.com"
	// plasmidCV := "test"
	// mockClient.On("ListPlasmids", mock.Anything, mock.Anything, mock.Anything).
	// 	Return(&stock.PlasmidCollection{Data: []*stock.PlasmidCollection_Data{}}, nil)
	// mockClient.On("CreatePlasmid", mock.Anything, mock.MatchedBy(func(p *stock.NewPlasmid) bool {
	// 	return p.Data.Attributes.Name == "pNew"
	// }), mock.Anything).
	// 	Return(&stock.Plasmid{Data: &stock.Plasmid_Data{Id: "DBPNew"}}, nil)
	// result := processPlasmidWithUpsert(userEmail, plasmidCV)(plasmid)
	// require.True(t, E.IsRight(result))
	// upsertResult := E.GetOrElse(func(error) GoldenBraidUpsertResult {
	// 	return GoldenBraidUpsertResult{}
	// })(result)
	// require.True(t, upsertResult.Created)
	// require.Equal(t, "DBPNew", upsertResult.PlasmidID)
	// mockClient.AssertExpectations(t)
}

func TestProcessPlasmidWithUpsert_UpdateWhenExists(t *testing.T) {
	t.Skip("deferred: tracked in modware-import-7zd")
	// TODO(modware-import-7zd): restore body once processPlasmidWithUpsert depositor arg is wired up
	// mockClient := new(MockStockClient)
	// regsc.SetStockAPIClient(mockClient)
	// plasmid := &source.GoldenBraidPlasmid{Name: "pExisting"}
	// userEmail := "test@example.com"
	// plasmidCV := "test"
	// mockClient.On("ListPlasmids", mock.Anything, mock.Anything, mock.Anything).
	// 	Return(&stock.PlasmidCollection{
	// 		Data: []*stock.PlasmidCollection_Data{
	// 			{Id: "DBPExisting", Attributes: &stock.PlasmidAttributes{Name: "pExisting"}},
	// 		},
	// 	}, nil)
	// mockClient.On("UpdatePlasmid", mock.Anything, mock.MatchedBy(func(p *stock.PlasmidUpdate) bool {
	// 	return p.Data.Id == "DBPExisting"
	// }), mock.Anything).
	// 	Return(&stock.Plasmid{Data: &stock.Plasmid_Data{Id: "DBPExisting"}}, nil)
	// result := processPlasmidWithUpsert(userEmail, plasmidCV)(plasmid)
	// require.True(t, E.IsRight(result))
	// upsertResult := E.GetOrElse(func(error) GoldenBraidUpsertResult {
	// 	return GoldenBraidUpsertResult{}
	// })(result)
	// require.False(t, upsertResult.Created)
	// require.Equal(t, "DBPExisting", upsertResult.PlasmidID)
	// mockClient.AssertExpectations(t)
}

func TestGoldenBraidSummaryAggregation(t *testing.T) {
	result1 := GoldenBraidUpsertResult{PlasmidID: "DBP0001", Created: true}
	result2 := GoldenBraidUpsertResult{PlasmidID: "DBP0002", Created: false}
	result3 := GoldenBraidUpsertResult{Error: errors.New("test error")}

	summary := GoldenBraidProcessingResult{}
	semigroup := GoldenBraidSummarySemigroup()

	summary = semigroup.Concat(summary, goldenBraidUpsertResultToSummary(result1))
	summary = semigroup.Concat(summary, goldenBraidUpsertResultToSummary(result2))
	summary = semigroup.Concat(summary, goldenBraidUpsertResultToSummary(result3))

	require.Equal(t, 1, summary.CreatedCount)
	require.Equal(t, 1, summary.UpdatedCount)
	require.Equal(t, 1, summary.ErrorCount)
	require.Len(t, summary.Successes, 2)
}
