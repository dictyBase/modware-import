package storage

import (
	IOE "github.com/IBM/fp-go/ioeither"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

// StockStorage defines the interface for stock data persistence operations.
// All methods return IOEither to represent side effects with error handling.
type StockStorage interface {
	// GetStrain retrieves a strain by ID
	GetStrain(stockID string) IOE.IOEither[error, *stock.Strain]

	// GetPlasmid retrieves a plasmid by ID
	GetPlasmid(stockID string) IOE.IOEither[error, *stock.Plasmid]

	// CreateStrain creates a new strain and returns the created strain with generated ID
	CreateStrain(req *stock.NewStrain) IOE.IOEither[error, *stock.Strain]

	// CreatePlasmid creates a new plasmid and returns the created plasmid with generated ID
	CreatePlasmid(req *stock.NewPlasmid) IOE.IOEither[error, *stock.Plasmid]

	// UpdateStrain updates an existing strain
	UpdateStrain(req *stock.StrainUpdate) IOE.IOEither[error, *stock.Strain]

	// UpdatePlasmid updates an existing plasmid
	UpdatePlasmid(req *stock.PlasmidUpdate) IOE.IOEither[error, *stock.Plasmid]

	// LoadStrain loads an existing strain with a specific ID (for bulk loading)
	LoadStrain(stockID string, req *stock.ExistingStrain) IOE.IOEither[error, *stock.Strain]

	// LoadPlasmid loads an existing plasmid with a specific ID (for bulk loading)
	LoadPlasmid(stockID string, req *stock.ExistingPlasmid) IOE.IOEither[error, *stock.Plasmid]

	// RemoveStock removes a stock by ID
	RemoveStock(stockID string) IOE.IOEither[error, struct{}]

	// ListStrains lists strains with filtering and pagination
	ListStrains(params *stock.StockParameters) IOE.IOEither[error, *stock.StrainCollection]

	// ListPlasmids lists plasmids with filtering and pagination
	ListPlasmids(params *stock.StockParameters) IOE.IOEither[error, *stock.PlasmidCollection]

	// ListStrainsByIDs retrieves multiple strains by their IDs
	ListStrainsByIDs(ids []string) IOE.IOEither[error, *stock.StrainList]

	// Close closes the storage backend
	Close() error
}
