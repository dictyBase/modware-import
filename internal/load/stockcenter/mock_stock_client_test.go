package stockcenter

import (
	"context"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockStockClient struct {
	mock.Mock
}

func (m *MockStockClient) GetPlasmid(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) != nil {
		return args.Get(0).(*stock.Plasmid), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStockClient) UpdatePlasmid(
	ctx context.Context,
	in *stock.PlasmidUpdate,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) != nil {
		return args.Get(0).(*stock.Plasmid), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStockClient) CreatePlasmid(
	ctx context.Context,
	in *stock.NewPlasmid,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Plasmid), args.Error(1)
}

func (m *MockStockClient) GetStrain(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Strain), args.Error(1)
}

func (m *MockStockClient) CreateStrain(
	ctx context.Context,
	in *stock.NewStrain,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Strain), args.Error(1)
}

func (m *MockStockClient) UpdateStrain(
	ctx context.Context,
	in *stock.StrainUpdate,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Strain), args.Error(1)
}

func (m *MockStockClient) ListStrains(
	ctx context.Context,
	in *stock.StockParameters,
	opts ...grpc.CallOption,
) (*stock.StrainCollection, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.StrainCollection), args.Error(1)
}

func (m *MockStockClient) ListPlasmids(
	ctx context.Context,
	in *stock.StockParameters,
	opts ...grpc.CallOption,
) (*stock.PlasmidCollection, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.PlasmidCollection), args.Error(1)
}

func (m *MockStockClient) LoadPlasmid(
	ctx context.Context,
	in *stock.ExistingPlasmid,
	opts ...grpc.CallOption,
) (*stock.Plasmid, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Plasmid), args.Error(1)
}

func (m *MockStockClient) LoadStrain(
	ctx context.Context,
	in *stock.ExistingStrain,
	opts ...grpc.CallOption,
) (*stock.Strain, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.Strain), args.Error(1)
}

//nolint:revive // method name matches interface
func (m *MockStockClient) ListStrainsByIds(
	ctx context.Context,
	in *stock.StockIdList,
	opts ...grpc.CallOption,
) (*stock.StrainList, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*stock.StrainList), args.Error(1)
}

func (m *MockStockClient) RemoveStock(
	ctx context.Context,
	in *stock.StockId,
	opts ...grpc.CallOption,
) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockStockClient) OboJSONFileUpload(
	ctx context.Context,
	opts ...grpc.CallOption,
) (stock.StockService_OboJSONFileUploadClient, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(stock.StockService_OboJSONFileUploadClient), args.Error(1)
}
