package pebble

import (
	"fmt"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
)

// Config holds configuration for Pebble storage
type Config struct {
	// DataDir is the directory for persistent storage.
	// If empty, uses in-memory storage.
	DataDir string
}

// ValidatedConfig represents a validated storage configuration
type ValidatedConfig struct {
	mode StorageMode
}

// StorageMode discriminates between in-memory and persistent storage
type StorageMode = E.Either[InMemoryConfig, PersistentConfig]

// InMemoryConfig represents in-memory storage configuration
type InMemoryConfig struct{}

// PersistentConfig represents persistent storage configuration
type PersistentConfig struct {
	dataDir string
}

// FilesystemSetup holds the result of filesystem initialization
type FilesystemSetup struct {
	memFS   vfs.FS // Only set for in-memory mode
	tempDir string // Only set for in-memory mode
	dataDir string // Only set for persistent mode
}

// DatabaseSetup holds the opened database and filesystem info
type DatabaseSetup struct {
	db      *pebble.DB
	memFS   vfs.FS
	tempDir string
}

// pebbleStorage implements the StockStorage interface using Pebble KV store
type pebbleStorage struct {
	db      *pebble.DB
	keys    keyBuilder
	memFS   vfs.FS // Only used for in-memory mode
	tempDir string // Only used for in-memory mode
}

// validateConfig converts optional config into validated config
func validateConfig(config *Config) E.Either[error, ValidatedConfig] {
	return F.Pipe2(
		O.FromNillable(config),
		O.GetOrElse(func() *Config { return &Config{} }),
		determineStorageMode,
	)
}

// determineStorageMode classifies config into storage mode
func determineStorageMode(config *Config) E.Either[error, ValidatedConfig] {
	if config.DataDir == "" {
		return E.Right[error](ValidatedConfig{
			mode: E.Left[PersistentConfig](InMemoryConfig{}),
		})
	}
	return E.Right[error](ValidatedConfig{
		mode: E.Right[InMemoryConfig](PersistentConfig{
			dataDir: config.DataDir,
		}),
	})
}

// setupInMemoryFilesystem creates virtual filesystem for in-memory mode
func setupInMemoryFilesystem(cfg InMemoryConfig) IOE.IOEither[error, FilesystemSetup] {
	return func() E.Either[error, FilesystemSetup] {
		memFS := vfs.NewMem()
		return E.Right[error](FilesystemSetup{
			memFS:   memFS,
			tempDir: "/mem",
			dataDir: "",
		})
	}
}

// setupPersistentFilesystem creates directory for persistent mode
func setupPersistentFilesystem(cfg PersistentConfig) IOE.IOEither[error, FilesystemSetup] {
	return IOE.TryCatchError(func() (FilesystemSetup, error) {
		if err := os.MkdirAll(cfg.dataDir, 0o755); err != nil {
			return FilesystemSetup{}, fmt.Errorf("failed to create data directory: %w", err)
		}
		return FilesystemSetup{
			memFS:   nil,
			tempDir: "",
			dataDir: cfg.dataDir,
		}, nil
	})
}

// setupFilesystem dispatches to appropriate setup function based on mode
func setupFilesystem(mode StorageMode) IOE.IOEither[error, FilesystemSetup] {
	return E.Fold(
		setupInMemoryFilesystem,
		setupPersistentFilesystem,
	)(mode)
}

// openDatabase opens Pebble database with appropriate configuration
func openDatabase(fsSetup FilesystemSetup) IOE.IOEither[error, DatabaseSetup] {
	return IOE.TryCatchError(func() (DatabaseSetup, error) {
		var db *pebble.DB
		var err error

		if fsSetup.memFS != nil {
			// In-memory mode
			db, err = pebble.Open(fsSetup.tempDir, &pebble.Options{
				FS: fsSetup.memFS,
			})
			if err != nil {
				return DatabaseSetup{}, fmt.Errorf("failed to open in-memory pebble database: %w", err)
			}
		} else {
			// Persistent mode
			db, err = pebble.Open(fsSetup.dataDir, &pebble.Options{})
			if err != nil {
				return DatabaseSetup{}, fmt.Errorf("failed to open pebble database: %w", err)
			}
		}

		return DatabaseSetup{
			db:      db,
			memFS:   fsSetup.memFS,
			tempDir: fsSetup.tempDir,
		}, nil
	})
}

// buildStorage constructs the final pebbleStorage from database setup
func buildStorage(dbSetup DatabaseSetup) *pebbleStorage {
	return &pebbleStorage{
		db:      dbSetup.db,
		keys:    newKeyBuilder(),
		memFS:   dbSetup.memFS,
		tempDir: dbSetup.tempDir,
	}
}

// NewStockStorage creates a new Pebble-backed stock storage
func NewStockStorage(config *Config) (*pebbleStorage, error) {
	result := F.Pipe2(
		validateConfig(config),
		E.Map[error](func(cfg ValidatedConfig) StorageMode {
			return cfg.mode
		}),
		E.Chain(func(mode StorageMode) E.Either[error, *pebbleStorage] {
			ioePipeline := F.Pipe2(
				setupFilesystem(mode),
				IOE.Chain(openDatabase),
				IOE.Map[error](buildStorage),
			)
			return ioePipeline()
		}),
	)

	if E.IsLeft(result) {
		return nil, E.Fold(
			func(err error) error { return err },
			func(_ *pebbleStorage) error { return nil },
		)(result)
	}

	return E.Fold(
		func(_ error) *pebbleStorage { return nil },
		func(storage *pebbleStorage) *pebbleStorage { return storage },
	)(result), nil
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
