package pebble

import (
	"fmt"
	"os"
	"sync"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/dictyBase/modware-import/internal/config"
	"github.com/dictyBase/modware-import/internal/fputil"
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

// Storage implements the StockStorage interface using Pebble KV store
type Storage struct {
	db      *pebble.DB
	keys    keyBuilder
	memFS   vfs.FS     // Only used for in-memory mode
	tempDir string     // Only used for in-memory mode
	mu      sync.Mutex // Protects ID generation
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
func setupInMemoryFilesystem(
	_ InMemoryConfig,
) IOE.IOEither[error, FilesystemSetup] {
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
func setupPersistentFilesystem(
	cfg PersistentConfig,
) IOE.IOEither[error, FilesystemSetup] {
	return IOE.TryCatchError(func() (FilesystemSetup, error) {
		if err := os.MkdirAll(cfg.dataDir, config.SharedDirectoryPermission); err != nil {
			return FilesystemSetup{}, fmt.Errorf(
				"failed to create data directory: %w",
				err,
			)
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
				return DatabaseSetup{}, fmt.Errorf(
					"failed to open in-memory pebble database: %w",
					err,
				)
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

// buildStorage constructs the final Storage from database setup
func buildStorage(dbSetup DatabaseSetup) *Storage {
	return &Storage{
		db:      dbSetup.db,
		keys:    newKeyBuilder(),
		memFS:   dbSetup.memFS,
		tempDir: dbSetup.tempDir,
	}
}

// ToTuple converts Either to Go-style tuple (value, error)
func ToTuple[A any](e E.Either[error, A]) T.Tuple2[A, error] {
	return F.Pipe1(
		e,
		E.Fold(
			func(err error) T.Tuple2[A, error] {
				var zero A
				return T.MakeTuple2(zero, err)
			},
			func(val A) T.Tuple2[A, error] {
				return T.MakeTuple2[A, error](val, nil)
			},
		))
}

// extractStorageMode extracts StorageMode from ValidatedConfig
func extractStorageMode(cfg ValidatedConfig) StorageMode {
	return cfg.mode
}

// executeStorageSetup executes the complete storage setup pipeline
func executeStorageSetup(
	mode StorageMode,
) E.Either[error, *Storage] {
	return F.Pipe3(
		setupFilesystem(mode),
		IOE.Chain(openDatabase),
		IOE.Map[error](buildStorage),
		fputil.ToEither[error, *Storage],
	)
}

// NewStockStorage creates a new Pebble-backed stock storage
func NewStockStorage(config *Config) (*Storage, error) {
	result := F.Pipe4(
		config,
		validateConfig,
		E.Map[error](extractStorageMode),
		E.Chain(executeStorageSetup),
		ToTuple[*Storage],
	)

	return result.F1, result.F2
}

// Close closes the Pebble database
func (storage *Storage) Close() error {
	if storage.db != nil {
		if err := storage.db.Close(); err != nil {
			return fmt.Errorf("failed to close pebble database: %w", err)
		}
	}
	return nil
}
