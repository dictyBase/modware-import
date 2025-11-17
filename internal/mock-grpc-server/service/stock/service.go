package stock

import (
	"context"

	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ServiceConfig holds configuration for the stock service
type ServiceConfig struct {
	StrainOntology  string
	StrainTerm      string
	PlasmidOntology string
	PlasmidTerm     string
	Logger          *logrus.Logger
}

// Service implements the gRPC stock service
type Service struct {
	stock.UnimplementedStockServiceServer
	storage storage.StockStorage
	config  *ServiceConfig
}

// NewService creates a new stock service instance
func NewService(storage storage.StockStorage, config *ServiceConfig) *Service {
	return &Service{
		storage: storage,
		config:  config,
	}
}

// GetStrain retrieves a strain by ID
func (service *Service) GetStrain(
	ctx context.Context,
	req *stock.StockId,
) (*stock.Strain, error) {
	result := getStrain(getStrainParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// CreateStrain creates a new strain
func (service *Service) CreateStrain(
	ctx context.Context,
	req *stock.NewStrain,
) (*stock.Strain, error) {
	result := createStrain(createStrainParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
		config:  service.config,
	})
	return result.F1, result.F2
}

// UpdateStrain updates an existing strain
func (service *Service) UpdateStrain(
	ctx context.Context,
	req *stock.StrainUpdate,
) (*stock.Strain, error) {
	result := updateStrain(updateStrainParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// GetPlasmid retrieves a plasmid by ID
func (service *Service) GetPlasmid(
	ctx context.Context,
	req *stock.StockId,
) (*stock.Plasmid, error) {
	result := getPlasmid(getPlasmidParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// CreatePlasmid creates a new plasmid
func (service *Service) CreatePlasmid(
	ctx context.Context,
	req *stock.NewPlasmid,
) (*stock.Plasmid, error) {
	result := createPlasmid(createPlasmidParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
		config:  service.config,
	})
	return result.F1, result.F2
}

// UpdatePlasmid updates an existing plasmid
func (service *Service) UpdatePlasmid(
	ctx context.Context,
	req *stock.PlasmidUpdate,
) (*stock.Plasmid, error) {
	result := updatePlasmid(updatePlasmidParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// RemoveStock removes a stock by ID
func (service *Service) RemoveStock(
	ctx context.Context,
	req *stock.StockId,
) (*emptypb.Empty, error) {
	result := removeStock(removeStockParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// ListStrains lists strains with filtering and pagination
func (service *Service) ListStrains(
	ctx context.Context,
	params *stock.StockParameters,
) (*stock.StrainCollection, error) {
	result := listStrains(listStrainsParams{
		ctx:     ctx,
		params:  params,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// ListPlasmids lists plasmids with filtering and pagination
func (service *Service) ListPlasmids(
	ctx context.Context,
	params *stock.StockParameters,
) (*stock.PlasmidCollection, error) {
	result := listPlasmids(listPlasmidsParams{
		ctx:     ctx,
		params:  params,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// LoadStrain loads an existing strain with a specific ID
func (service *Service) LoadStrain(
	ctx context.Context,
	req *stock.ExistingStrain,
) (*stock.Strain, error) {
	result := loadStrain(loadStrainParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// LoadPlasmid loads an existing plasmid with a specific ID
func (service *Service) LoadPlasmid(
	ctx context.Context,
	req *stock.ExistingPlasmid,
) (*stock.Plasmid, error) {
	result := loadPlasmid(loadPlasmidParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}

// ListStrainsByIds retrieves multiple strains by their IDs
//
//nolint:revive // Method name must match protobuf interface
func (service *Service) ListStrainsByIds(
	ctx context.Context,
	req *stock.StockIdList,
) (*stock.StrainList, error) {
	result := listStrainsByIDs(listStrainsByIDsParams{
		ctx:     ctx,
		request: req,
		storage: service.storage,
	})
	return result.F1, result.F2
}
