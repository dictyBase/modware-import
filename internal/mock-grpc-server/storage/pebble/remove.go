package pebble

import (
	"fmt"

	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/cockroachdb/pebble"
)

// RemoveStock removes a stock by ID
func (storage *Storage) RemoveStock(stockID string) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) {
		batch := storage.db.NewBatch()
		keys := storage.keys

		// Delete stock document
		if err := batch.Delete(keys.stockKey(stockID), pebble.Sync); err != nil {
			return struct{}{}, fmt.Errorf("failed to delete stock: %w", err)
		}

		// Delete JSON index
		if err := batch.Delete(keys.indexKey(stockID), pebble.Sync); err != nil {
			return struct{}{}, fmt.Errorf("failed to delete index: %w", err)
		}

		// Delete type classification
		if err := batch.Delete(keys.typeKey(stockID), pebble.Sync); err != nil {
			return struct{}{}, fmt.Errorf("failed to delete type: %w", err)
		}

		// Commit batch
		if err := batch.Commit(pebble.Sync); err != nil {
			return struct{}{}, fmt.Errorf("failed to commit delete batch: %w", err)
		}

		return struct{}{}, nil
	})
}
