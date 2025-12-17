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
	ctx context.Context,
	in *stock.NewPlasmid,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) UpdatePlasmid(
	ctx context.Context,
	in *stock.PlasmidUpdate,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) DeletePlasmid(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) LoadPlasmid(
	ctx context.Context,
	in *stock.ExistingPlasmid,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	return nil, nil
}

// Removed OboJSON as it caused undefined type issues and is not used
func (m *InventoryMockStockClient) CreateStrain(
	ctx context.Context,
	in *stock.NewStrain,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) ListStrains(
	ctx context.Context,
	in *stock.StockParameters,
	opts ...grpc.CallOption,
) (*stock.StrainCollection, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) GetStrain(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) UpdateStrain(
	ctx context.Context,
	in *stock.StrainUpdate,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) DeleteStrain(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) LoadStrain(
	ctx context.Context,
	in *stock.ExistingStrain,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) ListStrainsByIds(
	ctx context.Context,
	in *stock.StockIdList,
	opts ...grpc.CallOption,
) (*stock.StrainList, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) OboJSONFileUpload(
	ctx context.Context,
	opts ...grpc.CallOption,
) (stock.StockService_OboJSONFileUploadClient, error) {
	return nil, nil
}

func (m *InventoryMockStockClient) RemoveStock(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
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
	ctx context.Context,
	in *annotation.AnnotationGroupId,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotationGroup, error) {
	return nil, nil
}

// Stub other methods
func (m *InventoryMockAnnotationClient) GetAnnotation(
	ctx context.Context,
	in *annotation.AnnotationId,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) GetEntryAnnotation(
	ctx context.Context,
	in *annotation.EntryAnnotationRequest,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) GetAnnotationGroup(
	ctx context.Context,
	in *annotation.GroupEntryId,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotationGroup, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) UpdateAnnotation(
	ctx context.Context,
	in *annotation.TaggedAnnotationUpdate,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotation, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) DeleteAnnotation(
	ctx context.Context,
	in *annotation.DeleteAnnotationRequest,
	opts ...grpc.CallOption,
) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) GetAnnotationTag(
	ctx context.Context,
	in *annotation.TagRequest,
	opts ...grpc.CallOption,
) (*annotation.AnnotationTag, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) ListAnnotations(
	ctx context.Context,
	in *annotation.ListParameters,
	opts ...grpc.CallOption,
) (*annotation.TaggedAnnotationCollection, error) {
	return nil, nil
}

func (m *InventoryMockAnnotationClient) OboJSONFileUpload(
	ctx context.Context,
	opts ...grpc.CallOption,
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
	require.Contains(t, err.Error(), "plasmid not found")
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
