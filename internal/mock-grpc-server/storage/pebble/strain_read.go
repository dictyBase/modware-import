package pebble

import (
	"errors"
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/cockroachdb/pebble"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/fperrors"
)

// GetStrain retrieves a strain by ID
func (storage *pebbleStorage) GetStrain(stockID string) IOE.IOEither[error, *stock.Strain] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*stock.Strain, error) {
			data, closer, err := storage.db.Get(storage.keys.stockKey(stockID))
			if err != nil {
				if errors.Is(err, pebble.ErrNotFound) {
					return nil, fmt.Errorf("strain %s not found", stockID)
				}
				return nil, fmt.Errorf("failed to get strain: %w", err)
			}
			defer closer.Close()

			// Make a copy since data is only valid until closer.Close()
			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)

			strain, err := deserializeStrain(dataCopy)
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize strain: %w", err)
			}

			return strain, nil
		}),
		IOE.MapLeft[*stock.Strain](
			fperrors.OnError("storage operation failed"),
		),
	)
}

// GetPlasmid retrieves a plasmid by ID
func (storage *pebbleStorage) GetPlasmid(stockID string) IOE.IOEither[error, *stock.Plasmid] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*stock.Plasmid, error) {
			data, closer, err := storage.db.Get(storage.keys.stockKey(stockID))
			if err != nil {
				if errors.Is(err, pebble.ErrNotFound) {
					return nil, fmt.Errorf("plasmid %s not found", stockID)
				}
				return nil, fmt.Errorf("failed to get plasmid: %w", err)
			}
			defer closer.Close()

			// Make a copy since data is only valid until closer.Close()
			dataCopy := make([]byte, len(data))
			copy(dataCopy, data)

			plasmid, err := deserializePlasmid(dataCopy)
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize plasmid: %w", err)
			}

			return plasmid, nil
		}),
		IOE.MapLeft[*stock.Plasmid](
			fperrors.OnError("storage operation failed"),
		),
	)
}

// ListStrainsByIds retrieves multiple strains by their IDs
func (storage *pebbleStorage) ListStrainsByIds(ids []string) IOE.IOEither[error, *stock.StrainList] {
	return IOE.TryCatchError(func() (*stock.StrainList, error) {
		strains := make([]*stock.StrainList_Data, 0, len(ids))

		for _, stockID := range ids {
			result := storage.GetStrain(stockID)()

			// Extract strain from Either, skip errors
			F.Pipe1(
				result,
				E.Fold(
					func(err error) *stock.Strain {
						// Skip errors silently
						return nil
					},
					func(strain *stock.Strain) *stock.Strain {
						if strain != nil {
							strains = append(strains, &stock.StrainList_Data{
								Type:       strain.Data.Type,
								Id:         strain.Data.Id,
								Attributes: strain.Data.Attributes,
							})
						}
						return strain
					},
				),
			)
		}

		return &stock.StrainList{
			Data: strains,
		}, nil
	})
}
