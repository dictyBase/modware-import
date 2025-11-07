package pebble

import (
	"encoding/json"
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/cockroachdb/pebble"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/filter"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/pagination"
)

// ListStrains lists strains with filtering and pagination
func (storage *pebbleStorage) ListStrains(
	params *stock.StockParameters,
) IOE.IOEither[error, *stock.StrainCollection] {
	return IOE.TryCatchError(func() (*stock.StrainCollection, error) {
		// Parse filter or use always-true filter
		filterExpr := F.Pipe1(
			filter.ParseFilter(params.Filter),
			E.GetOrElse(func(err error) filter.FilterExpression {
				return filter.AlwaysTrueFilter{}
			}),
		)

		strains := make([]*stock.Strain, 0)

		// Create iterator with index prefix to only get strain indices
		iter, err := storage.db.NewIter(&pebble.IterOptions{
			LowerBound: []byte(indexPrefix),
			UpperBound: []byte(indexPrefix + "\xff"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create iterator: %w", err)
		}
		defer iter.Close()

		// Collect all strains and apply filter
		for iter.First(); iter.Valid(); iter.Next() {
			// Extract stock ID from index key
			indexKey := string(iter.Key())
			stockID := indexKey[len(indexPrefix):]

			// Get JSON index for filtering
			jsonData, closer, err := storage.db.Get(iter.Key())
			if err != nil {
				continue // Skip on error
			}
			defer closer.Close()

			var indexMap map[string]interface{}
			if err := json.Unmarshal(jsonData, &indexMap); err != nil {
				continue // Skip invalid JSON
			}

			// Apply filter
			if !filterExpr.Evaluate(indexMap) {
				continue
			}

			// Get the actual strain document
			result := storage.GetStrain(stockID)()

			// Extract strain from Either
			F.Pipe1(
				result,
				E.Fold(
					func(err error) *stock.Strain {
						return nil
					},
					func(strain *stock.Strain) *stock.Strain {
						if strain != nil {
							strains = append(strains, strain)
						}
						return strain
					},
				),
			)
		}

		if err := iter.Error(); err != nil {
			return nil, fmt.Errorf("iterator error: %w", err)
		}

		// Apply pagination
		limit := params.Limit
		if limit <= 0 {
			limit = 10
		}

		cursor := params.Cursor
		if cursor < 0 {
			cursor = 0
		}

		total := int64(len(strains))
		sliceResult := pagination.CalculateSlice(pagination.SliceParams{
			Cursor: cursor,
			Limit:  limit,
			Total:  total,
		})

		// Slice the results
		paginatedStrains := strains[sliceResult.Start:sliceResult.End]

		// Convert to collection data format
		collectionData := make([]*stock.StrainCollection_Data, len(paginatedStrains))
		for idx, strain := range paginatedStrains {
			collectionData[idx] = &stock.StrainCollection_Data{
				Type:       strain.Data.Type,
				Id:         strain.Data.Id,
				Attributes: strain.Data.Attributes,
			}
		}

		return &stock.StrainCollection{
			Data: collectionData,
			Meta: &stock.Meta{
				Limit:      limit,
				Total:      total,
				NextCursor: sliceResult.NextCursor,
			},
		}, nil
	})
}

// ListPlasmids lists plasmids with filtering and pagination
func (storage *pebbleStorage) ListPlasmids(
	params *stock.StockParameters,
) IOE.IOEither[error, *stock.PlasmidCollection] {
	return IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
		// Parse filter or use always-true filter
		filterExpr := F.Pipe1(
			filter.ParseFilter(params.Filter),
			E.GetOrElse(func(err error) filter.FilterExpression {
				return filter.AlwaysTrueFilter{}
			}),
		)

		plasmids := make([]*stock.Plasmid, 0)

		// Create iterator with index prefix to only get plasmid indices
		iter, err := storage.db.NewIter(&pebble.IterOptions{
			LowerBound: []byte(indexPrefix),
			UpperBound: []byte(indexPrefix + "\xff"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create iterator: %w", err)
		}
		defer iter.Close()

		// Collect all plasmids and apply filter
		for iter.First(); iter.Valid(); iter.Next() {
			// Extract stock ID from index key
			indexKey := string(iter.Key())
			stockID := indexKey[len(indexPrefix):]

			// Get JSON index for filtering
			jsonData, closer, err := storage.db.Get(iter.Key())
			if err != nil {
				continue
			}
			defer closer.Close()

			var indexMap map[string]interface{}
			if err := json.Unmarshal(jsonData, &indexMap); err != nil {
				continue
			}

			// Apply filter
			if !filterExpr.Evaluate(indexMap) {
				continue
			}

			// Get the actual plasmid document
			result := storage.GetPlasmid(stockID)()

			// Extract plasmid from Either
			F.Pipe1(
				result,
				E.Fold(
					func(err error) *stock.Plasmid {
						return nil
					},
					func(plasmid *stock.Plasmid) *stock.Plasmid {
						if plasmid != nil {
							plasmids = append(plasmids, plasmid)
						}
						return plasmid
					},
				),
			)
		}

		if err := iter.Error(); err != nil {
			return nil, fmt.Errorf("iterator error: %w", err)
		}

		// Apply pagination
		limit := params.Limit
		if limit <= 0 {
			limit = 10
		}

		cursor := params.Cursor
		if cursor < 0 {
			cursor = 0
		}

		total := int64(len(plasmids))
		sliceResult := pagination.CalculateSlice(pagination.SliceParams{
			Cursor: cursor,
			Limit:  limit,
			Total:  total,
		})

		// Slice the results
		paginatedPlasmids := plasmids[sliceResult.Start:sliceResult.End]

		// Convert to collection data format
		collectionData := make([]*stock.PlasmidCollection_Data, len(paginatedPlasmids))
		for idx, plasmid := range paginatedPlasmids {
			collectionData[idx] = &stock.PlasmidCollection_Data{
				Type:       plasmid.Data.Type,
				Id:         plasmid.Data.Id,
				Attributes: plasmid.Data.Attributes,
			}
		}

		return &stock.PlasmidCollection{
			Data: collectionData,
			Meta: &stock.Meta{
				Limit:      limit,
				Total:      total,
				NextCursor: sliceResult.NextCursor,
			},
		}, nil
	})
}
