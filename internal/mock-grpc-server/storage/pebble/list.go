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
func (storage *Storage) ListStrains(
	params *stock.StockParameters,
) IOE.IOEither[error, *stock.StrainCollection] {
	return IOE.TryCatchError(func() (*stock.StrainCollection, error) {
		filterExpr := parseFilterOrDefault(params.Filter)
		strains, err := storage.collectFilteredStrains(filterExpr)
		if err != nil {
			return nil, err
		}

		limit, cursor := normalizePaginationParams(params.Limit, params.Cursor)
		paginatedStrains, sliceResult := applyPagination(strains, limit, cursor)
		collectionData := convertToStrainCollectionData(paginatedStrains)

		return &stock.StrainCollection{
			Data: collectionData,
			Meta: createMeta(limit, int64(len(strains)), sliceResult.NextCursor),
		}, nil
	})
}

// collectFilteredStrains collects all strains that match the filter expression
func (storage *Storage) collectFilteredStrains(
	filterExpr filter.Expression,
) ([]*stock.Strain, error) {
	iter, err := createIndexIterator(storage.db)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	strains := make([]*stock.Strain, 0)
	for iter.First(); iter.Valid(); iter.Next() {
		strain := storage.processStrainIndexEntry(iter, filterExpr)
		if strain != nil {
			strains = append(strains, strain)
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	return strains, nil
}

// processStrainIndexEntry processes a single index entry and returns the strain if it matches the filter
func (storage *Storage) processStrainIndexEntry(
	iter *pebble.Iterator,
	filterExpr filter.Expression,
) *stock.Strain {
	indexMap, stockID := extractIndexData(iter, storage.db)
	if indexMap == nil {
		return nil
	}

	if !filterExpr.Evaluate(indexMap) {
		return nil
	}

	return extractStrainFromResult(storage.GetStrain(stockID)())
}

// ListPlasmids lists plasmids with filtering and pagination
func (storage *Storage) ListPlasmids(
	params *stock.StockParameters,
) IOE.IOEither[error, *stock.PlasmidCollection] {
	return IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
		filterExpr := parseFilterOrDefault(params.Filter)
		plasmids, err := storage.collectFilteredPlasmids(filterExpr)
		if err != nil {
			return nil, err
		}

		limit, cursor := normalizePaginationParams(params.Limit, params.Cursor)
		paginatedPlasmids, sliceResult := applyPagination(plasmids, limit, cursor)
		collectionData := convertToPlasmidCollectionData(paginatedPlasmids)

		return &stock.PlasmidCollection{
			Data: collectionData,
			Meta: createMeta(limit, int64(len(plasmids)), sliceResult.NextCursor),
		}, nil
	})
}

// collectFilteredPlasmids collects all plasmids that match the filter expression
func (storage *Storage) collectFilteredPlasmids(
	filterExpr filter.Expression,
) ([]*stock.Plasmid, error) {
	iter, err := createIndexIterator(storage.db)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	plasmids := make([]*stock.Plasmid, 0)
	for iter.First(); iter.Valid(); iter.Next() {
		plasmid := storage.processPlasmidIndexEntry(iter, filterExpr)
		if plasmid != nil {
			plasmids = append(plasmids, plasmid)
		}
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("iterator error: %w", err)
	}

	return plasmids, nil
}

// processPlasmidIndexEntry processes a single index entry and returns the plasmid if it matches the filter
func (storage *Storage) processPlasmidIndexEntry(
	iter *pebble.Iterator,
	filterExpr filter.Expression,
) *stock.Plasmid {
	indexMap, stockID := extractIndexData(iter, storage.db)
	if indexMap == nil {
		return nil
	}

	if !filterExpr.Evaluate(indexMap) {
		return nil
	}

	return extractPlasmidFromResult(storage.GetPlasmid(stockID)())
}

// Helper functions shared between ListStrains and ListPlasmids

// parseFilterOrDefault parses the filter string or returns an always-true filter
func parseFilterOrDefault(filterStr string) filter.Expression {
	return F.Pipe1(
		filter.ParseFilter(filterStr),
		E.GetOrElse(func(error) filter.Expression {
			return filter.AlwaysTrueFilter{}
		}),
	)
}

// createIndexIterator creates a pebble iterator for the index prefix
func createIndexIterator(db *pebble.DB) (*pebble.Iterator, error) {
	iter, err := db.NewIter(&pebble.IterOptions{
		LowerBound: []byte(indexPrefix),
		UpperBound: []byte(indexPrefix + "\xff"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create iterator: %w", err)
	}
	return iter, nil
}

// extractIndexData extracts and parses the JSON index data from an iterator entry
func extractIndexData(iter *pebble.Iterator, _ *pebble.DB) (map[string]interface{}, string) {
	indexKey := string(iter.Key())
	stockID := indexKey[len(indexPrefix):]

	// Use iterator's value directly instead of Get
	jsonData := iter.Value()
	if jsonData == nil {
		return nil, ""
	}

	var indexMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &indexMap); err != nil {
		return nil, ""
	}

	return indexMap, stockID
}

// extractStrainFromResult extracts a strain from an Either result
func extractStrainFromResult(result E.Either[error, *stock.Strain]) *stock.Strain {
	return F.Pipe1(
		result,
		E.Fold(
			func(error) *stock.Strain {
				return nil
			},
			func(strain *stock.Strain) *stock.Strain {
				return strain
			},
		),
	)
}

// extractPlasmidFromResult extracts a plasmid from an Either result
func extractPlasmidFromResult(result E.Either[error, *stock.Plasmid]) *stock.Plasmid {
	return F.Pipe1(
		result,
		E.Fold(
			func(error) *stock.Plasmid {
				return nil
			},
			func(plasmid *stock.Plasmid) *stock.Plasmid {
				return plasmid
			},
		),
	)
}

// normalizePaginationParams ensures limit and cursor have valid values
func normalizePaginationParams(limit, cursor int64) (int64, int64) {
	if limit <= 0 {
		limit = 10
	}
	if cursor < 0 {
		cursor = 0
	}
	return limit, cursor
}

// applyPagination applies pagination to a slice and returns the paginated result
func applyPagination[T any](items []T, limit, cursor int64) ([]T, pagination.SliceResult) {
	total := int64(len(items))
	sliceResult := pagination.CalculateSlice(pagination.SliceParams{
		Cursor: cursor,
		Limit:  limit,
		Total:  total,
	})
	return items[sliceResult.Start:sliceResult.End], sliceResult
}

// convertToStrainCollectionData converts strains to collection data format
func convertToStrainCollectionData(strains []*stock.Strain) []*stock.StrainCollection_Data {
	collectionData := make([]*stock.StrainCollection_Data, len(strains))
	for idx, strain := range strains {
		collectionData[idx] = &stock.StrainCollection_Data{
			Type:       strain.Data.Type,
			Id:         strain.Data.Id,
			Attributes: strain.Data.Attributes,
		}
	}
	return collectionData
}

// convertToPlasmidCollectionData converts plasmids to collection data format
func convertToPlasmidCollectionData(plasmids []*stock.Plasmid) []*stock.PlasmidCollection_Data {
	collectionData := make([]*stock.PlasmidCollection_Data, len(plasmids))
	for idx, plasmid := range plasmids {
		collectionData[idx] = &stock.PlasmidCollection_Data{
			Type:       plasmid.Data.Type,
			Id:         plasmid.Data.Id,
			Attributes: plasmid.Data.Attributes,
		}
	}
	return collectionData
}

// createMeta creates a Meta object with pagination information
func createMeta(limit, total, nextCursor int64) *stock.Meta {
	return &stock.Meta{
		Limit:      limit,
		Total:      total,
		NextCursor: nextCursor,
	}
}
