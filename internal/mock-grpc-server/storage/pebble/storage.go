package pebble

import (
	"fmt"
	"os"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
)

// Config holds configuration for Pebble storage
type Config struct {
	// DataDir is the directory for persistent storage.
	// If empty, uses in-memory storage.
	DataDir string
}

// pebbleStorage implements the StockStorage interface using Pebble KV store
type pebbleStorage struct {
	db      *pebble.DB
	keys    keyBuilder
	memFS   vfs.FS // Only used for in-memory mode
	tempDir string // Only used for in-memory mode
}

// NewStockStorage creates a new Pebble-backed stock storage
func NewStockStorage(config *Config) (*pebbleStorage, error) {
	if config == nil {
		config = &Config{}
	}

	var (
		db      *pebble.DB
		err     error
		memFS   vfs.FS
		tempDir string
	)

	if config.DataDir == "" {
		// In-memory mode using virtual filesystem
		memFS = vfs.NewMem()
		tempDir = "/mem"

		db, err = pebble.Open(tempDir, &pebble.Options{
			FS: memFS,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to open in-memory pebble database: %w", err)
		}
	} else {
		// Persistent mode using actual filesystem
		if err := os.MkdirAll(config.DataDir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}

		db, err = pebble.Open(config.DataDir, &pebble.Options{})
		if err != nil {
			return nil, fmt.Errorf("failed to open pebble database: %w", err)
		}
	}

	return &pebbleStorage{
		db:      db,
		keys:    newKeyBuilder(),
		memFS:   memFS,
		tempDir: tempDir,
	}, nil
}

// Close closes the Pebble database
func (storage *pebbleStorage) Close() error {
	if storage.db != nil {
		if err := storage.db.Close(); err != nil {
			return fmt.Errorf("failed to close pebble database: %w", err)
		}
	}
	return nil
}
