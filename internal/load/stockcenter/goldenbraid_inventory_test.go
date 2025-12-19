package stockcenter

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/either"
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

// Test ResolvePlasmidID
func TestResolvePlasmidID(t *testing.T) {
	mockStock := new(InventoryMockStockClient)
	regsc.SetStockAPIClient(mockStock)

	// Case 1: Found
	mockStock.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Filter == "name==pTest"
	}), mock.Anything).Return(&stock.PlasmidCollection{
		Data: []*stock.PlasmidCollection_Data{
			{Id: "12345"}, // Assumed field based on usage
		},
	}, nil).Once()

	id, err := E.Unwrap(resolvePlasmidID("pTest")())
	require.NoError(t, err)
	require.Equal(t, "12345", id)

	// Case 2: Not Found
	mockStock.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Filter == "name==pUnknown"
	}), mock.Anything).Return(&stock.PlasmidCollection{
		Data: []*stock.PlasmidCollection_Data{},
	}, nil).Once()

	_, err = E.Unwrap(resolvePlasmidID("pUnknown")())
	require.Error(t, err)
	require.Contains(t, err.Error(), "expected 1 plasmid")
}

func TestResolvePlasmidIDPointFree(t *testing.T) {
	mockStock := new(InventoryMockStockClient)
	regsc.SetStockAPIClient(mockStock)

	// Case 1: Found
	name := "pDV-CFPC-GB0000"
	mockStock.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Filter == "name=="+name
	}), mock.Anything).Return(&stock.PlasmidCollection{
		Data: []*stock.PlasmidCollection_Data{
			{Id: "test-id-123"},
		},
	}, nil).Once()

	id, err := E.Unwrap(resolvePlasmidID(name)())
	require.NoError(t, err)
	require.Equal(t, "test-id-123", id)

	// Case 2: Not Found
	unknownName := "pUnknown"
	mockStock.On("ListPlasmids", mock.Anything, mock.MatchedBy(func(p *stock.StockParameters) bool {
		return p.Filter == "name=="+unknownName
	}), mock.Anything).Return(&stock.PlasmidCollection{
		Data: []*stock.PlasmidCollection_Data{},
	}, nil).Once()

	_, err = E.Unwrap(resolvePlasmidID(unknownName)())
	require.Error(t, err)
	require.Contains(t, err.Error(), "expected 1 plasmid")
}

// Test SyncInventory
func TestSyncInventory(t *testing.T) {
	mockAnno := new(InventoryMockAnnotationClient)
	regsc.SetAnnotationAPIClient(mockAnno)

	plasmidID := "123"
	location := "Box1"

	// Mock Check Inventory (Found)
	expectedFilter := fmt.Sprintf(
		"entry_id==%s;tag==%s;ontology==%s",
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

	_, err := E.Unwrap(syncInventory(plasmidID, location)())
	require.NoError(t, err)
	mockAnno.AssertExpectations(t)
}

func TestFileSQLQuery(t *testing.T) {
	csvContent := `Name,Location
pTest1,Loc1
pTest2,Loc2`

	reader := io.NopCloser(strings.NewReader(csvContent))
	ctx := context.Background()

	builder, err := filesql.NewBuilder().
		AddReader(reader, "inventory", filesql.FileTypeCSV).
		Build(ctx)
	require.NoError(t, err)

	db, err := builder.Open(ctx)
	require.NoError(t, err)
	defer db.Close()

	rows, err := db.Query("SELECT Name, Location FROM inventory")
	require.NoError(t, err)
	defer rows.Close()

	count := 0
	for rows.Next() {
		var name, loc string
		err := rows.Scan(&name, &loc)
		require.NoError(t, err)
		if count == 0 {
			require.Equal(t, "pTest1", name)
			require.Equal(t, "Loc1", loc)
		}
		count++
	}
	require.Equal(t, 2, count)
}
