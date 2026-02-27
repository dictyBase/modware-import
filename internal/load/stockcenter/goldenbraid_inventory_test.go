package stockcenter

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/nao1215/filesql"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// InventoryMockStockClient
type InventoryMockStockClient struct {
	mock.Mock
}

func (m *InventoryMockStockClient) ListPlasmids(
	ctx context.Context,
	in *stock.StockParameters,
	opts ...grpc.CallOption,
) (*stock.PlasmidCollection, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.PlasmidCollection), args.Error(1)
}

func (m *InventoryMockStockClient) GetPlasmid(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Plasmid), args.Error(1)
}

func (m *InventoryMockStockClient) CreatePlasmid(
	_ context.Context,
	_ *stock.NewPlasmid,
	_ ...grpc.CallOption,
) (*stock.Plasmid, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) UpdatePlasmid(
	_ context.Context,
	_ *stock.PlasmidUpdate,
	_ ...grpc.CallOption,
) (*stock.Plasmid, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) DeletePlasmid(
	_ context.Context,
	_ *stock.StockId,
	_ ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) LoadPlasmid(
	_ context.Context,
	_ *stock.ExistingPlasmid,
	_ ...grpc.CallOption,
) (*stock.Plasmid, error) {
	return nil, nil
}

// Removed OboJSON as it caused undefined type issues and is not used
func (m *InventoryMockStockClient) CreateStrain(
	_ context.Context,
	_ *stock.NewStrain,
	_ ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) ListStrains(
	_ context.Context,
	_ *stock.StockParameters,
	_ ...grpc.CallOption,
) (*stock.StrainCollection, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) GetStrain(
	_ context.Context,
	_ *stock.StockId,
	_ ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) UpdateStrain(
	_ context.Context,
	_ *stock.StrainUpdate,
	_ ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) DeleteStrain(
	_ context.Context,
	_ *stock.StockId,
	_ ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) LoadStrain(
	_ context.Context,
	_ *stock.ExistingStrain,
	_ ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) ListStrainsByIds( //nolint:revive // API-defined method name
	_ context.Context,
	_ *stock.StockIdList,
	_ ...grpc.CallOption,
) (*stock.StrainList, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) OboJSONFileUpload(
	_ context.Context,
	_ ...grpc.CallOption,
) (stock.StockService_OboJSONFileUploadClient, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) RemoveStock(
	_ context.Context,
	_ *stock.StockId,
	_ ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

// InventoryMockAnnotationClient
type InventoryMockAnnotationClient struct {
	mock.Mock
}

func (m *InventoryMockAnnotationClient) ListAnnotationGroups(
	ctx context.Context,
	in *annotation.ListGroupParameters,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotationGroupCollection, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*annotation.TaggedAnnotationGroupCollection), args.Error(1)
}

func (m *InventoryMockAnnotationClient) DeleteAnnotationGroup(
	ctx context.Context,
	in *annotation.GroupEntryId,
	opts ...grpc.CallOption,
) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *InventoryMockAnnotationClient) CreateAnnotationGroup(
	ctx context.Context,
	in *annotation.AnnotationIdList,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotationGroup, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*annotation.TaggedAnnotationGroup), args.Error(1)
}

func (m *InventoryMockAnnotationClient) CreateAnnotation(
	ctx context.Context,
	in *annotation.NewTaggedAnnotation,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*annotation.TaggedAnnotation), args.Error(1)
}

func (m *InventoryMockAnnotationClient) AddToAnnotationGroup(
	_ context.Context,
	_ *annotation.AnnotationGroupId,
	_ ...grpc.CallOption,
) (*annotation.TaggedAnnotationGroup, error) {
	return nil, nil
}

// Stub other methods
func (m *InventoryMockAnnotationClient) GetAnnotation(
	_ context.Context,
	_ *annotation.AnnotationId,
	_ ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) GetEntryAnnotation(
	_ context.Context,
	_ *annotation.EntryAnnotationRequest,
	_ ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) GetAnnotationGroup(
	_ context.Context,
	_ *annotation.GroupEntryId,
	_ ...grpc.CallOption,
) (*annotation.TaggedAnnotationGroup, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) UpdateAnnotation(
	_ context.Context,
	_ *annotation.TaggedAnnotationUpdate,
	_ ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) DeleteAnnotation(
	_ context.Context,
	_ *annotation.DeleteAnnotationRequest,
	_ ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) GetAnnotationTag(
	_ context.Context,
	_ *annotation.TagRequest,
	_ ...grpc.CallOption,
) (*annotation.AnnotationTag, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) ListAnnotations(
	_ context.Context,
	_ *annotation.ListParameters,
	_ ...grpc.CallOption,
) (*annotation.TaggedAnnotationCollection, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) OboJSONFileUpload(
	_ context.Context,
	_ ...grpc.CallOption,
) (annotation.TaggedAnnotationService_OboJSONFileUploadClient, error) {
	return nil, nil
}

func TestProcessRow(t *testing.T) {
	mockStock := new(InventoryMockStockClient)
	mockAnno := new(InventoryMockAnnotationClient)

	deps := Deps{
		StockClient:      mockStock,
		AnnotationClient: mockAnno,
	}
	ctx := PipelineContext{
		Deps:        deps,
		PlasmidName: "pTest",
		Location:    "Box1",
	}
	plasmidID := "123"

	// Mock ListPlasmids (Found)
	mockStock.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Filter == "plasmid_name===pTest"
	}), mock.Anything).Return(&stock.PlasmidCollection{
		Data: []*stock.PlasmidCollection_Data{
			{Id: plasmidID},
		},
	}, nil).Once()

	// Mock Check Inventory (Found)
	expectedFilter := fmt.Sprintf(
		"entry_id===%s;tag===%s;ontology===%s",
		plasmidID,
		regsc.InvLocationTag,
		regsc.PlasmidInvOntO,
	)

	mockAnno.On("ListAnnotationGroups", mock.Anything, mock.MatchedBy(func(p *annotation.ListGroupParameters) bool {
		return p.Filter == expectedFilter
	}), mock.Anything).
		Return(&annotation.TaggedAnnotationGroupCollection{
			Data: []*annotation.TaggedAnnotationGroupCollection_Data{
				{
					Group: &annotation.TaggedAnnotationGroup{GroupId: "grp1"},
				},
			},
		}, nil).
		Once()

	// Mock Delete
	mockAnno.On("DeleteAnnotationGroup", mock.Anything, mock.Anything, mock.Anything).
		Return(&emptypb.Empty{}, nil).Once()

	// Mock Create Annotations
	mockAnno.On("CreateAnnotation", mock.Anything, mock.Anything, mock.Anything).
		Return(&annotation.TaggedAnnotation{Data: &annotation.TaggedAnnotation_Data{Id: "anno1"}}, nil).
		Times(3)

	// Mock Create Group
	mockAnno.On("CreateAnnotationGroup", mock.Anything, mock.Anything, mock.Anything).
		Return(&annotation.TaggedAnnotationGroup{}, nil).Once()

	summary := processRowToSummary(ctx)
	require.Equal(t, 0, summary.ErrorCount, "unexpected errors: %v", summary.Errors)
	require.Equal(t, 1, summary.SuccessCount)

	mockStock.AssertExpectations(t)
}

func TestProcessRowSkipsWhenPlasmidNotFound(t *testing.T) {
	mockStock := new(InventoryMockStockClient)
	mockAnno := new(InventoryMockAnnotationClient)
	ctx := PipelineContext{
		Deps: Deps{
			StockClient:      mockStock,
			AnnotationClient: mockAnno,
		},
		PlasmidName: "pMissing",
		Location:    "Box99",
	}

	// ListPlasmids returns empty collection — plasmid not in stock
	mockStock.On("ListPlasmids", mock.Anything, mock.Anything, mock.Anything).
		Return(&stock.PlasmidCollection{}, nil).Once()

	summary := processRowToSummary(ctx)

	// Silent skip: neither success nor error
	require.Equal(t, 0, summary.ErrorCount)
	require.Equal(t, 0, summary.SuccessCount)
	mockStock.AssertExpectations(t)
	// Annotation client must never be reached on the None branch
	mockAnno.AssertNotCalled(t, "ListAnnotationGroups")
}

func TestFileSQLQuery(t *testing.T) {
	// pTest1 matches by Plasmid Name directly
	// pTest2 matches via Synonym (pAlias2 in inventory)
	// pTest3 has no inventory entry → must NOT appear in results
	// pTest1 appears twice with different locations → both distinct rows
	gbCSV := `Plasmid Name,Synonym,Depositor,Genes,Keywords,Description,PMID
pTest1,syn1,Depositor A,geneA,keywords,desc1,123456
pTest2,pAlias2,Depositor B,geneB,keywords,desc2,789012
pTest3,syn3,Depositor C,geneC,keywords,desc3,345678`

	invCSV := `Name,Location
pTest1,Box1
pTest1,Box1A
pAlias2,Box2`

	gbReader := io.NopCloser(strings.NewReader(gbCSV))
	invReader := io.NopCloser(strings.NewReader(invCSV))
	ctx := context.Background()

	builder, err := filesql.NewBuilder().
		AddReader(invReader, "goldenbraid_inventory", filesql.FileTypeCSV).
		AddReader(gbReader, "goldenbraid", filesql.FileTypeCSV).
		Build(ctx)
	require.NoError(t, err)

	db, err := builder.Open(ctx)
	require.NoError(t, err)
	defer db.Close()

	// Test JOIN query returns correct (name, location) pairs
	rows, err := db.Query(inventoryJoinQuery)
	require.NoError(t, err)
	defer rows.Close()

	type result struct{ name, loc string }
	var results []result
	for rows.Next() {
		var r result
		require.NoError(t, rows.Scan(&r.name, &r.loc))
		results = append(results, r)
	}
	require.NoError(t, rows.Err())
	// pTest1/Box1 + pTest1/Box1A + pTest2/Box2 = 3 distinct rows; pTest3 absent
	require.Equal(t, 3, len(results))

	// Test COUNT query matches row count
	var count int
	err = db.QueryRow(inventoryCountQuery).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 3, count)
}
