# NewStockStorage fp-go Integration Design

**Date:** 2025-11-07
**Status:** Approved Design
**Pattern:** Hybrid approach - functional internals with traditional Go signature

## Overview

Refactor `NewStockStorage` to use fp-go internally while maintaining the traditional `(*pebbleStorage, error)` signature. This matches the established codebase pattern where:
- Storage methods (CreateStrain, CreatePlasmid) use full IOEither pipelines
- Constructor functions maintain traditional Go signatures for easier consumption
- Business logic is functional, setup/teardown is imperative

## Design Principles

1. **Point-free style with named functions** - No inline lambdas, all pipeline steps are named
2. **Option for config presence** - Handle nil config with Option monad
3. **Either for storage mode discrimination** - Type-safe mode selection
4. **Separate functions per mode** - In-memory and persistent paths are distinct
5. **Same return types** - Both modes return FilesystemSetup for consistency
6. **IOEither for side effects** - Filesystem and database operations wrapped in IOEither

## Type Definitions

### Configuration Types

```go
// ValidatedConfig represents a validated storage configuration
type ValidatedConfig struct {
    mode StorageMode
}

// StorageMode discriminates between in-memory and persistent storage
type StorageMode = E.Either[InMemoryConfig, PersistentConfig]

type InMemoryConfig struct{}

type PersistentConfig struct {
    dataDir string
}
```

### Setup Result Types

```go
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
```

## Pipeline Functions

### Layer 1: Configuration Validation (Pure)

```go
// validateConfig converts optional config into validated config
func validateConfig(config *Config) E.Either[error, ValidatedConfig] {
    return F.Pipe2(
        O.FromNullable(config),
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
```

### Layer 2: Filesystem Setup (Side Effects)

```go
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
```

### Layer 3: Database Opening (Side Effects)

```go
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
```

### Layer 4: Storage Construction (Pure)

```go
// buildStorage constructs the final pebbleStorage from database setup
func buildStorage(dbSetup DatabaseSetup) *pebbleStorage {
    return &pebbleStorage{
        db:      dbSetup.db,
        keys:    newKeyBuilder(),
        memFS:   dbSetup.memFS,
        tempDir: dbSetup.tempDir,
    }
}
```

## Main Pipeline

```go
// NewStockStorage creates a new Pebble-backed stock storage
func NewStockStorage(config *Config) (*pebbleStorage, error) {
    pipeline := F.Pipe3(
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

    return E.Fold(
        func(err error) (*pebbleStorage, error) {
            return nil, err
        },
        func(storage *pebbleStorage) (*pebbleStorage, error) {
            return storage, nil
        },
    )(pipeline)
}
```

## Data Flow

1. **validateConfig** → `Either[error, ValidatedConfig]` (pure)
   - Handle nil config with Option.FromNullable
   - Provide default empty Config
   - Determine storage mode based on DataDir

2. **Extract mode** → `Either[error, StorageMode]` (pure)
   - Map ValidatedConfig to StorageMode
   - StorageMode is Either[InMemoryConfig, PersistentConfig]

3. **E.Chain into IOEither pipeline**:
   - **setupFilesystem** → `IOEither[error, FilesystemSetup]` (side effect)
     - Fold on mode to dispatch to correct setup function
     - In-memory: create virtual FS
     - Persistent: create real directory

   - **openDatabase** → `IOEither[error, DatabaseSetup]` (side effect)
     - Check fsSetup.memFS to determine mode
     - Open Pebble with appropriate options

   - **buildStorage** → `IOEither[error, *pebbleStorage]` (pure)
     - Construct final pebbleStorage struct
     - Initialize keyBuilder

4. **Execute IOEither** by calling `()` → `Either[error, *pebbleStorage]`

5. **Fold to traditional Go** → `(*pebbleStorage, error)`

## Benefits

1. **Composable** - Each step is a named, testable function
2. **Type-safe** - StorageMode discriminated union prevents invalid states
3. **Point-free** - No inline lambdas except for extraction/folding
4. **Error handling** - All errors flow through Either/IOEither
5. **Testable** - Each layer can be tested independently
6. **Backward compatible** - Same signature, existing callers unaffected
7. **Consistent** - Matches CreateStrain pattern for internal operations

## Testing Strategy

Each layer can be tested independently:

```go
// Test pure validation
func TestValidateConfig(t *testing.T) {
    result := validateConfig(nil)
    // Assert Either is Right with in-memory mode
}

// Test filesystem setup
func TestSetupInMemoryFilesystem(t *testing.T) {
    result := setupInMemoryFilesystem(InMemoryConfig{})()
    // Assert Either is Right with valid FilesystemSetup
}

// Test database opening (requires test filesystem)
func TestOpenDatabase(t *testing.T) {
    fsSetup := FilesystemSetup{memFS: vfs.NewMem(), tempDir: "/mem"}
    result := openDatabase(fsSetup)()
    // Assert Either is Right with open database
}

// Test end-to-end
func TestNewStockStorage(t *testing.T) {
    storage, err := NewStockStorage(&Config{})
    require.NoError(t, err)
    require.NotNil(t, storage)
    defer storage.Close()
}
```

## Migration Notes

- Replace existing NewStockStorage implementation
- No changes required to callers (CLI code)
- Add imports: `O "github.com/IBM/fp-go/option"`
- Existing E, F, IOE imports already present
- Consider extracting mode discrimination logic if reused elsewhere
