# Modware-Stock Mock gRPC Server Design

**Date:** 2025-11-07
**Status:** Design Complete
**Purpose:** Integration testing and local development mock for stock service

## Overview

This design document describes a mock gRPC server implementation for the modware-stock service. The mock server will support integration testing of modware-import tools and provide a lightweight development environment without requiring the full production stack (ArangoDB, NATS, etc.).

### Goals

1. **Integration Testing** - Test modware-import against stock service without external dependencies
2. **Development Environment** - Provide lightweight stock service for local development
3. **Production Fidelity** - Match production behavior including validation, errors, and filter syntax
4. **Deep fp-go Integration** - Use functional patterns throughout: IOEither, Reader, Either, Option

### Key Design Decisions

- **Storage:** Pebble key-value store with in-memory mode
- **Serialization:** Hybrid protobuf (storage) + JSON (indices)
- **ID Generation:** Sequential starting from 1 (DBS0000001, DBP0000001)
- **Validation:** Full go-proto-validators validation
- **Filter Parser:** Complete implementation of all operators
- **Pagination:** Offset-based with snapshot isolation
- **CLI:** Subcommand in existing mock-grpc-server

## Architecture

### Core fp-go Patterns

**1. Progressive Context Building**

Each workflow stage enriches typed context structs:

```go
// Example: GetStrain workflow
getStrainContext → withValidatedStrainID → withStrainDoc
```

**2. IOEither Pipelines**

```go
func (s *StockService) GetStrain(ctx context.Context, req *stock.StockId) (*stock.Strain, error) {
    result := F.Pipe4(
        IOE.Of[error](getStrainContext{ctx: ctx, request: req, storage: s.storage}),
        IOE.Bind(setValidatedStrainID, validateStrainIDRequest),
        IOE.Bind(setStrainDoc, retrieveStrainFromStorage),
        IOE.Map[error](extractStrainResponse),
        toStrainResult(ctx),
    )
    return result.F1, result.F2
}
```

**3. Curried Setters**

```go
var (
    setValidatedStrainID = F.Curry2(
        func(id string, ctx getStrainContext) withValidatedStrainID {
            return withValidatedStrainID{
                getStrainContext: ctx,
                validatedID:      id,
            }
        },
    )

    setStrainDoc = F.Curry2(
        func(strain *stock.Strain, ctx withValidatedStrainID) withStrainDoc {
            return withStrainDoc{
                withValidatedStrainID: ctx,
                strain:                strain,
            }
        },
    )
)
```

**4. Reader Monad for Dependencies**

```go
type PaginationDependencies struct {
    DB         *pebble.DB
    Snapshot   *pebble.Snapshot
    FilterExpr FilterExpression
    Params     *stock.StockParameters
}

type PaginationContext struct {
    Results    []*stock.Strain
    Sliced     []*stock.Strain
    NextCursor int64
    Limit      int64
}

type PaginationReader = RE.ReaderEither[PaginationDependencies, error, PaginationContext]
```

**5. Named Helper Functions**

All functions are named (not var), except curried setters which use `var` with `F.Curry2`. Functions are univariate, receiving a single typed struct parameter.

### Component Overview

```
mock-grpc-server/
├── cmd/mock-grpc-server/
│   └── main.go (add "stock" subcommand)
├── internal/mock-grpc-server/
│   ├── cli/
│   │   └── stock.go (RunStockServer)
│   ├── service/stock/
│   │   ├── service.go (gRPC service implementation)
│   │   ├── strain.go (strain workflows)
│   │   ├── plasmid.go (plasmid workflows)
│   │   ├── strain_types.go (context types for strain)
│   │   ├── plasmid_types.go (context types for plasmid)
│   │   └── result_converters.go (tuple result conversion)
│   ├── storage/pebble/
│   │   ├── storage.go (interface implementation)
│   │   ├── strain_create.go (create strain workflow)
│   │   ├── strain_read.go (get/list strain workflows)
│   │   ├── strain_update.go (update strain workflow)
│   │   ├── plasmid_create.go (create plasmid workflow)
│   │   ├── plasmid_read.go (get/list plasmid workflows)
│   │   ├── plasmid_update.go (update plasmid workflow)
│   │   ├── common.go (shared helpers)
│   │   └── namespaces.go (key namespace operations)
│   ├── filter/
│   │   ├── parser.go (filter parser)
│   │   ├── tokenizer.go (lexical tokenizer)
│   │   ├── evaluator.go (predicate evaluation)
│   │   └── types.go (filter expression types)
│   └── pagination/
│       ├── cursor.go (cursor encoding/decoding)
│       └── slice.go (pagination helpers)
```

## Storage Layer Design

### Interface Definition

```go
type StockStorage interface {
    GetStrain(id string) IOE.IOEither[error, *stock.Strain]
    GetPlasmid(id string) IOE.IOEither[error, *stock.Plasmid]

    CreateStrain(req *stock.NewStrain) IOE.IOEither[error, *stock.Strain]
    CreatePlasmid(req *stock.NewPlasmid) IOE.IOEither[error, *stock.Plasmid]

    UpdateStrain(req *stock.StrainUpdate) IOE.IOEither[error, *stock.Strain]
    UpdatePlasmid(req *stock.PlasmidUpdate) IOE.IOEither[error, *stock.Plasmid]

    LoadStrain(id string, req *stock.ExistingStrain) IOE.IOEither[error, *stock.Strain]
    LoadPlasmid(id string, req *stock.ExistingPlasmid) IOE.IOEither[error, *stock.Plasmid]

    RemoveStock(id string) IOE.IOEither[error, struct{}]

    ListStrains(params *stock.StockParameters) IOE.IOEither[error, *stock.StrainCollection]
    ListPlasmids(params *stock.StockParameters) IOE.IOEither[error, *stock.PlasmidCollection]
    ListStrainsByIds(ids []string) IOE.IOEither[error, *stock.StrainList]
}
```

### Key Namespace Schema

Pebble uses namespace prefixes to organize different data types:

```go
// Stock documents - protobuf serialized
stock:{id} → binary protobuf data
  Example: stock:DBS0000001 → Strain protobuf bytes

// JSON indices - for filtering/queries
index:{id} → JSON representation
  Example: index:DBS0000001 → {"id":"DBS0000001","depositor":"Costanza",...}

// Type edges - stock classification
type:{stock_id} → "strain" | "plasmid"
  Example: type:DBS0000001 → "strain"

// Parent relationships - strain hierarchy
parent:{strain_id} → parent strain ID
  Example: parent:DBS0000042 → "DBS0000001"

// Ontology terms - stock classification
term:{stock_id} → ontology term
  Example: term:DBS0000001 → "general strain"

// ID counters - sequential ID generation
counter:strain → int64 next ID
counter:plasmid → int64 next ID
  Example: counter:strain → 43

// Reverse indices - for common queries
depositor:{name}:{id} → empty (existence check)
species:{species}:{id} → empty (existence check)
  Example: depositor:Costanza:DBS0000001 → ""
```

### CreateStrain Workflow

**Context Types:**

```go
type createStrainContext struct {
    req *stock.NewStrain
    db  *pebble.DB
}

type withGeneratedID struct {
    createStrainContext
    generatedID string
}

type withTimestamps struct {
    withGeneratedID
    createdAt time.Time
    updatedAt time.Time
}

type withBuiltStrain struct {
    withTimestamps
    strain *stock.Strain
}

type withSerializedData struct {
    withBuiltStrain
    protoBytes []byte
    jsonBytes  []byte
}

type withBatchKeys struct {
    withSerializedData
    batch *pebble.Batch
}
```

**Curried Setters:**

```go
var (
    setGeneratedID = F.Curry2(
        func(id string, ctx createStrainContext) withGeneratedID {
            return withGeneratedID{
                createStrainContext: ctx,
                generatedID:         id,
            }
        },
    )

    setTimestamps = F.Curry2(
        func(ts T.Tuple2[time.Time, time.Time], ctx withGeneratedID) withTimestamps {
            return withTimestamps{
                withGeneratedID: ctx,
                createdAt:       ts.F1,
                updatedAt:       ts.F2,
            }
        },
    )

    // ... additional setters
)
```

**Workflow Functions:**

```go
// generateStrainID generates next sequential ID
func generateStrainID(ctx createStrainContext) IOE.IOEither[error, string] {
    return IOE.TryCatchError(func() (string, error) {
        counter, err := ctx.db.Get([]byte("counter:strain"))
        if err != nil && !errors.Is(err, pebble.ErrNotFound) {
            return "", fmt.Errorf("failed to read strain counter: %w", err)
        }

        var nextID int64 = 1
        if err == nil {
            nextID = int64(binary.BigEndian.Uint64(counter)) + 1
        }

        return fmt.Sprintf("DBS%07d", nextID), nil
    })
}

// generateTimestamps creates current timestamps
func generateTimestamps(ctx withGeneratedID) IOE.IOEither[error, T.Tuple2[time.Time, time.Time]] {
    return func() E.Either[error, T.Tuple2[time.Time, time.Time]] {
        now := time.Now()
        return E.Right[error](T.MakeTuple2(now, now))
    }
}

// buildStrainFromRequest constructs strain - pure function
func buildStrainFromRequest(ctx withTimestamps) *stock.Strain {
    return F.Pipe2(
        buildStrainAttributes(ctx),
        func(attrs *stock.StrainAttributes) *stock.Strain_Data {
            return buildStrainData(ctx.generatedID, attrs)
        },
        buildCompleteStrain,
    )
}

// serializeStrain serializes to protobuf and JSON
func serializeStrain(ctx withBuiltStrain) IOE.IOEither[error, T.Tuple2[[]byte, []byte]] {
    return IOE.TryCatchError(func() (T.Tuple2[[]byte, []byte], error) {
        protoBytes, err := proto.Marshal(ctx.strain)
        if err != nil {
            return T.MakeTuple2[[]byte, []byte](nil, nil),
                fmt.Errorf("failed to marshal strain: %w", err)
        }

        jsonIndex := buildJSONIndex(ctx.strain)
        jsonBytes, err := json.Marshal(jsonIndex)
        if err != nil {
            return T.MakeTuple2[[]byte, []byte](nil, nil),
                fmt.Errorf("failed to marshal JSON index: %w", err)
        }

        return T.MakeTuple2(protoBytes, jsonBytes), nil
    })
}

// buildBatchWrite creates pebble batch with all keys
func buildBatchWrite(ctx withSerializedData) IOE.IOEither[error, *pebble.Batch] {
    return IOE.TryCatchError(func() (*pebble.Batch, error) {
        batch := ctx.db.NewBatch()

        // Write stock document
        if err := batch.Set(
            []byte("stock:"+ctx.generatedID),
            ctx.protoBytes,
            pebble.Sync,
        ); err != nil {
            return nil, fmt.Errorf("failed to set stock: %w", err)
        }

        // Write JSON index
        if err := batch.Set(
            []byte("index:"+ctx.generatedID),
            ctx.jsonBytes,
            pebble.Sync,
        ); err != nil {
            return nil, fmt.Errorf("failed to set index: %w", err)
        }

        // Write type, counter, etc.
        // ...

        return batch, nil
    })
}

// commitBatch commits the pebble batch
func commitBatch(ctx withBatchKeys) IOE.IOEither[error, withBatchKeys] {
    return IOE.TryCatchError(func() (withBatchKeys, error) {
        if err := ctx.batch.Commit(pebble.Sync); err != nil {
            return ctx, fmt.Errorf("failed to commit batch: %w", err)
        }
        return ctx, nil
    })
}
```

**Main Pipeline:**

```go
func (s *pebbleStorage) CreateStrain(req *stock.NewStrain) IOE.IOEither[error, *stock.Strain] {
    return F.Pipe7(
        IOE.Of[error](createStrainContext{req: req, db: s.db}),
        IOE.Bind(setGeneratedID, generateStrainID),
        IOE.Bind(setTimestamps, generateTimestamps),
        IOE.Let[error](setBuiltStrain, buildStrainFromRequest),
        IOE.Bind(setSerializedData, serializeStrain),
        IOE.Bind(setBatchKeys, buildBatchWrite),
        IOE.Chain(commitBatch),
        IOE.Map[error](extractCreatedStrain),
    )
}
```

### GetStrain Workflow

Simpler workflow for retrieval:

```go
func (s *pebbleStorage) GetStrain(id string) IOE.IOEither[error, *stock.Strain] {
    return F.Pipe1(
        IOE.TryCatchError(func() (*stock.Strain, error) {
            data, err := s.db.Get([]byte("stock:" + id))
            if err != nil {
                if errors.Is(err, pebble.ErrNotFound) {
                    return nil, fmt.Errorf("strain %s not found", id)
                }
                return nil, fmt.Errorf("failed to get strain: %w", err)
            }

            var strain stock.Strain
            if err := proto.Unmarshal(data, &strain); err != nil {
                return nil, fmt.Errorf("failed to unmarshal strain: %w", err)
            }

            return &strain, nil
        }),
        IOE.MapLeft[*stock.Strain](
            fperrors.OnError("storage operation failed"),
        ),
    )
}
```

## Service Layer Design

### GetStrain Workflow

**Context Types:**

```go
type getStrainContext struct {
    ctx     context.Context
    request *stock.StockId
    storage StockStorage
}

type withValidatedStrainID struct {
    getStrainContext
    validatedID string
}

type withStrainDoc struct {
    withValidatedStrainID
    strain *stock.Strain
}
```

**Workflow Functions:**

```go
// validateStrainIDRequest validates the stock ID
func validateStrainIDRequest(ctx getStrainContext) IOE.IOEither[error, string] {
    return F.Pipe1(
        validationToEither(ctx),
        IOE.FromEither[error, string],
    )
}

// validationToEither converts validation to Either
func validationToEither(ctx getStrainContext) E.Either[error, string] {
    return F.Pipe2(
        validateRequest(ctx.request),
        O.FromNillable[error],
        O.Fold(
            func() E.Either[error, string] {
                return F.Pipe1(
                    extractRequestID(ctx),
                    validateIDFormat,
                )
            },
            func(err error) E.Either[error, string] {
                return E.Left[string](
                    fmt.Errorf("invalid request parameters: %w", err),
                )
            },
        ),
    )
}

// retrieveStrainFromStorage retrieves strain with error enrichment
func retrieveStrainFromStorage(ctx withValidatedStrainID) IOE.IOEither[error, *stock.Strain] {
    return F.Pipe1(
        ctx.storage.GetStrain(ctx.validatedID),
        IOE.MapLeft[*stock.Strain](enrichStrainError(ctx.validatedID)),
    )
}
```

**Result Conversion:**

```go
type (
    StrainResult    = T.Tuple2[*stock.Strain, error]
    StrainEither    = E.Either[error, *stock.Strain]
    StrainIO        = IOE.IOEither[error, *stock.Strain]
    StrainConverter = func(StrainIO) StrainResult
)

// toStrainResult converts IOEither to tuple result
func toStrainResult(ctx context.Context) StrainConverter {
    return F.Flow2(
        IOE.ToEither[error, *stock.Strain],
        E.Fold(
            createErrorResult(ctx),
            createSuccessResult,
        ),
    )
}

// createErrorResult creates error result tuple
func createErrorResult(ctx context.Context) func(error) StrainResult {
    return func(err error) StrainResult {
        return T.MakeTuple2(
            &stock.Strain{},
            errorToGRPCStatus(ctx)(err),
        )
    }
}

// createSuccessResult creates success result tuple
func createSuccessResult(strain *stock.Strain) StrainResult {
    return T.MakeTuple2[*stock.Strain, error](strain, nil)
}
```

**Main Method:**

```go
func (s *StockService) GetStrain(
    ctx context.Context,
    req *stock.StockId,
) (*stock.Strain, error) {
    result := F.Pipe4(
        IOE.Of[error](getStrainContext{
            ctx:     ctx,
            request: req,
            storage: s.storage,
        }),
        IOE.Bind(setValidatedStrainID, validateStrainIDRequest),
        IOE.Bind(setStrainDoc, retrieveStrainFromStorage),
        IOE.Map[error](extractStrainResponse),
        toStrainResult(ctx),
    )
    return result.F1, result.F2
}
```

### CreateStrain Workflow

Similar pattern with validation and default term application:

```go
func (s *StockService) CreateStrain(
    ctx context.Context,
    req *stock.NewStrain,
) (*stock.Strain, error) {
    result := F.Pipe4(
        IOE.Of[error](createStrainContext{
            ctx:     ctx,
            request: req,
            storage: s.storage,
            config:  s.config,
        }),
        IOE.Bind(setValidatedNewStrain, validateNewStrainRequest),
        IOE.Bind(setCreatedStrain, createStrainInStorage),
        IOE.Map[error](extractCreatedStrainResponse),
        toStrainResult(ctx),
    )
    return result.F1, result.F2
}
```

### ListStrains Workflow

Uses Reader monad for complex pagination with dependencies:

```go
func (s *StockService) ListStrains(
    ctx context.Context,
    params *stock.StockParameters,
) (*stock.StrainCollection, error) {
    limit := defaultLimit(params.Limit)

    result := F.Pipe4(
        IOE.Of[error](listStrainsContext{
            ctx:     ctx,
            params:  params,
            limit:   limit,
            storage: s.storage,
        }),
        IOE.Bind(setValidatedFilter, validateStrainFilter),
        IOE.Bind(setStrainCollection, retrieveStrainsFromStorage),
        IOE.Map[error](extractStrainCollectionResponse),
        toStrainCollectionResult(ctx, limit),
    )
    return result.F1, result.F2
}
```

## Filter Engine Design

### Parser Architecture

**Tokenizer → Parser → Evaluator**

1. **Tokenizer** - Splits filter string into tokens
2. **Parser** - Builds expression AST using Either
3. **Evaluator** - Evaluates predicates against JSON data

### Filter Expression Types

```go
type FilterExpression interface {
    Evaluate(stockJSON map[string]interface{}) bool
}

type Operator int

const (
    // String operators
    Contains Operator = iota
    NotContains
    Equals
    NotEquals

    // Numeric operators
    NumEquals
    GreaterThan
    LessThan
    GreaterOrEqual
    LessOrEqual

    // Date operators
    DateEquals
    DateGreater
    DateLess
    DateGreaterOrEqual
    DateLessOrEqual

    // Array operators
    ArrayContains
    ArrayNotContains
    ArrayEquals
    ArrayNotEquals
)

type Predicate struct {
    Field    string
    Operator Operator
    Value    string
}

type AndExpression struct {
    Left  FilterExpression
    Right FilterExpression
}

type OrExpression struct {
    Left  FilterExpression
    Right FilterExpression
}

type AlwaysTrueFilter struct{}
```

### Evaluator Using fp-go Patterns

**Function Dispatch via Map:**

```go
type OperatorEvaluator func(actual interface{}, expected string) bool

var operatorEvaluators = map[Operator]OperatorEvaluator{
    Contains:    evalStringContains,
    NotContains: evalStringNotContains,
    Equals:      evalStringEquals,
    NotEquals:   evalStringNotEquals,
    // ... all operators
}

func getEvaluator(op Operator) O.Option[OperatorEvaluator] {
    eval, ok := operatorEvaluators[op]
    if !ok {
        return O.None[OperatorEvaluator]()
    }
    return O.Some(eval)
}

func evaluateValue(op Operator, actual interface{}, expected string) bool {
    return F.Pipe1(
        getEvaluator(op),
        O.Fold(
            func() bool { return false },
            func(evaluator OperatorEvaluator) bool {
                return evaluator(actual, expected)
            },
        ),
    )
}
```

**String Evaluators Using Option:**

```go
func evalStringContains(actual interface{}, expected string) bool {
    return F.Pipe1(
        O.FromPredicate(isString)(actual),
        O.Map[string](func(str string) bool {
            return strings.Contains(str, expected)
        }),
        O.GetOrElse(func() bool { return false }),
    )
}

func isString(val interface{}) bool {
    _, ok := val.(string)
    return ok
}

func extractString(val interface{}) O.Option[string] {
    return O.FromPredicate(isString)(val)
}
```

**Numeric Evaluators Using Option Chain:**

```go
func evaluateNumeric(
    actual interface{},
    expected string,
    comparator func(float64, float64) bool,
) bool {
    return F.Pipe2(
        extractFloat64(actual),
        O.Chain(func(actualNum float64) O.Option[bool] {
            return F.Pipe1(
                parseFloat(expected),
                O.Map[float64](func(expectedNum float64) bool {
                    return comparator(actualNum, expectedNum)
                }),
            )
        }),
        O.GetOrElse(func() bool { return false }),
    )
}

func parseFloat(str string) O.Option[float64] {
    num, err := strconv.ParseFloat(str, 64)
    if err != nil {
        return O.None[float64]()
    }
    return O.Some(num)
}
```

**Array Evaluators Using Array Functions:**

```go
func evalArrayContains(actual interface{}, expected string) bool {
    return F.Pipe2(
        extractArray(actual),
        O.Map[[]interface{}](func(arr []interface{}) bool {
            return A.Some(arrayItemContains(expected))(arr)
        }),
        O.GetOrElse(func() bool { return false }),
    )
}

func arrayItemContains(expected string) func(interface{}) bool {
    return func(item interface{}) bool {
        return F.Pipe1(
            extractString(item),
            O.Map[string](func(str string) bool {
                return strings.Contains(str, expected)
            }),
            O.GetOrElse(func() bool { return false }),
        )
    }
}
```

### Parser Using Either Composition

```go
func ParseFilter(filterStr string) E.Either[error, FilterExpression] {
    return F.Pipe2(
        filterStr,
        tokenizeFilter,
        parseOrExpression,
    )
}

func buildAndTreeRecursive(groups [][]Token) E.Either[error, FilterExpression] {
    return F.Pipe2(
        parseExpression(groups[0]),
        E.Chain(func(left FilterExpression) E.Either[error, FilterExpression] {
            return F.Pipe1(
                buildAndTree(groups[1:]),
                E.Map[error](func(right FilterExpression) FilterExpression {
                    return AndExpression{Left: left, Right: right}
                }),
            )
        }),
    )
}
```

## Pagination Design

### Reader Monad Pattern

**Dependencies and Context:**

```go
type PaginationDependencies struct {
    DB         *pebble.DB
    Snapshot   *pebble.Snapshot
    FilterExpr FilterExpression
    Params     *stock.StockParameters
}

type PaginationContext struct {
    Results    []*stock.Strain
    Sliced     []*stock.Strain
    NextCursor int64
    Limit      int64
}

type PaginationReader = RE.ReaderEither[PaginationDependencies, error, PaginationContext]
```

**Univariate Functions with Typed Structs:**

```go
// SliceContext holds slicing parameters
type SliceContext struct {
    Results []*stock.Strain
    Cursor  int64
    Limit   int64
}

func sliceResults(ctx PaginationContext) PaginationReader {
    return F.Pipe1(
        RE.Ask[PaginationDependencies](),
        RE.Map[PaginationDependencies, error](func(deps PaginationDependencies) PaginationContext {
            sliceCtx := SliceContext{
                Results: ctx.Results,
                Cursor:  deps.Params.Cursor,
                Limit:   ctx.Limit,
            }

            sliced, nextCursor := applyPaginationSlice(sliceCtx)

            ctx.Sliced = sliced
            ctx.NextCursor = nextCursor
            return ctx
        }),
    )
}

// applyPaginationSlice performs the actual slicing
func applyPaginationSlice(ctx SliceContext) ([]*stock.Strain, int64) {
    window := calculateSliceWindow(SliceWindow{
        Cursor: ctx.Cursor,
        Limit:  ctx.Limit,
        Total:  int64(len(ctx.Results)),
    })

    start := window.Start
    end := window.End

    sliced := ctx.Results[start:end]
    nextCursor := computeNextCursorValue(NextCursorParams{
        TotalResults: int64(len(ctx.Results)),
        CurrentEnd:   end,
    })

    return sliced, nextCursor
}
```

**Main Pipeline:**

```go
func (s *pebbleStorage) ListStrains(
    params *stock.StockParameters,
) IOE.IOEither[error, *stock.StrainCollection] {
    filterExpr := F.Pipe1(
        params.Filter,
        parseFilterOrDefault,
        E.GetOrElse(func() FilterExpression { return AlwaysTrueFilter{} }),
    )

    limit := defaultLimit(params.Limit)

    deps := PaginationDependencies{
        DB:         s.db,
        FilterExpr: filterExpr,
        Params:     params,
    }

    initialCtx := PaginationContext{Limit: limit}

    return F.Pipe1(
        F.Pipe4(
            createSnapshotWithDeps(initialCtx),
            RE.Chain[PaginationDependencies, error](iterateAndFilter),
            RE.Chain[PaginationDependencies, error](applySortIfRequested),
            RE.Chain[PaginationDependencies, error](sliceResults),
            RE.Map[PaginationDependencies, error](buildStrainCollectionFromContext),
        )(deps),
        IOE.FromEither[error, *stock.StrainCollection],
    )
}
```

### Cursor Encoding

Simple offset-based encoding:

```go
type CursorToken struct {
    QueryTimestamp int64
    Offset         int64
}

func encodeCursor(token CursorToken) int64 {
    return token.Offset
}

func decodeCursor(cursor int64) E.Either[error, CursorToken] {
    if cursor < 0 {
        return E.Left[CursorToken](
            fmt.Errorf("invalid cursor: must be non-negative"),
        )
    }

    return E.Right[error](CursorToken{
        QueryTimestamp: time.Now().UnixNano(),
        Offset:         cursor,
    })
}
```

## CLI Integration

### Command Structure

Add new subcommand to existing mock-grpc-server:

```go
// In cmd/mock-grpc-server/main.go
app := &cli.App{
    Name:  "mock-grpc-server",
    Usage: "Mock gRPC servers for integration testing",
    Commands: []*cli.Command{
        {
            Name:   "stock",
            Usage:  "Mock gRPC server for stock service",
            Action: stockcli.RunStockServer,
            Flags:  stockFlags(),
        },
        // ... existing commands
    },
}
```

### Flags (Matching Production)

```go
func stockFlags() []cli.Flag {
    return []cli.Flag{
        &cli.IntFlag{
            Name:    "port",
            Aliases: []string{"p"},
            Value:   9560,
            Usage:   "Server port",
            EnvVars: []string{"STOCK_GRPC_PORT"},
        },
        &cli.StringFlag{
            Name:    "log-level",
            Aliases: []string{"l"},
            Value:   "info",
            Usage:   "Log level (debug, info, warn, error)",
            EnvVars: []string{"LOG_LEVEL"},
        },
        &cli.StringFlag{
            Name:    "log-format",
            Value:   "json",
            Usage:   "Log format (json, text)",
            EnvVars: []string{"LOG_FORMAT"},
        },
        &cli.StringFlag{
            Name:    "data-dir",
            Value:   "",
            Usage:   "Pebble data directory (empty for in-memory)",
            EnvVars: []string{"PEBBLE_DATA_DIR"},
        },
        &cli.BoolFlag{
            Name:  "reflection",
            Value: true,
            Usage: "Enable gRPC server reflection",
        },
        &cli.StringFlag{
            Name:  "strain-ontology",
            Value: "dicty_strain_property",
            Usage: "Ontology for strain grouping terms",
        },
        &cli.StringFlag{
            Name:  "strain-term",
            Value: "general strain",
            Usage: "Default ontology term for strains",
        },
        &cli.StringFlag{
            Name:  "plasmid-ontology",
            Value: "plasmid_keywords",
            Usage: "Ontology for plasmid grouping terms",
        },
        &cli.StringFlag{
            Name:  "plasmid-term",
            Value: "cloning vector",
            Usage: "Default ontology term for plasmids",
        },
    }
}
```

### Server Initialization

```go
func RunStockServer(c *cli.Context) error {
    // Initialize Pebble storage
    storage, err := pebble.NewStockStorage(&pebble.Config{
        DataDir: c.String("data-dir"),
    })
    if err != nil {
        return cli.Exit(fmt.Sprintf("failed to create storage: %v", err), 1)
    }
    defer storage.Close()

    // Create gRPC server
    grpcServer := grpc.NewServer(
        grpc.ChainUnaryInterceptor(
            grpc_ctxtags.UnaryServerInterceptor(),
            grpc_logrus.UnaryServerInterceptor(getLogger(c)),
        ),
    )

    // Register stock service
    stock.RegisterStockServiceServer(
        grpcServer,
        stockservice.NewStockService(storage, &stockservice.ServiceConfig{
            StrainOntology:  c.String("strain-ontology"),
            StrainTerm:      c.String("strain-term"),
            PlasmidOntology: c.String("plasmid-ontology"),
            PlasmidTerm:     c.String("plasmid-term"),
        }),
    )

    // Enable reflection
    if c.Bool("reflection") {
        reflection.Register(grpcServer)
    }

    // Start server
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Int("port")))
    if err != nil {
        return cli.Exit(fmt.Sprintf("failed to listen: %v", err), 1)
    }

    log.Printf("Starting stock mock server on port %d", c.Int("port"))
    return grpcServer.Serve(lis)
}
```

## Error Handling

### Production-Matching gRPC Errors

```go
func errorToGRPCStatus(ctx context.Context) func(error) error {
    return func(err error) error {
        if err == nil {
            return nil
        }

        errMsg := err.Error()

        switch {
        case strings.Contains(errMsg, "not found"):
            return status.Error(codes.NotFound, errMsg)
        case strings.Contains(errMsg, "invalid"):
            return status.Error(codes.InvalidArgument, errMsg)
        case strings.Contains(errMsg, "already exists"):
            return status.Error(codes.AlreadyExists, errMsg)
        default:
            return status.Error(codes.Internal, errMsg)
        }
    }
}
```

### Validation Errors

Full go-proto-validators support:

```go
func validateNewStrainRequest(ctx createStrainContext) IOE.IOEither[error, *stock.NewStrain] {
    return func() E.Either[error, *stock.NewStrain] {
        if err := ctx.request.Validate(); err != nil {
            return E.Left[*stock.NewStrain](
                fmt.Errorf("invalid request parameters: %w", err),
            )
        }

        // Apply defaults
        if len(ctx.request.Data.Attributes.DictyStrainProperty) == 0 {
            ctx.request.Data.Attributes.DictyStrainProperty = ctx.config.StrainTerm
        }

        return E.Right[error](ctx.request)
    }
}
```

## Testing Strategy

### Unit Tests

Test each layer independently using testify/require:

```go
func TestPebbleStorage_CreateStrain(t *testing.T) {
    storage := setupTestStorage(t)
    defer storage.Close()

    req := &stock.NewStrain{
        Data: &stock.NewStrain_Data{
            Type: "strain",
            Attributes: &stock.NewStrainAttributes{
                CreatedBy: "test@example.com",
                UpdatedBy: "test@example.com",
                Depositor: "John Doe",
                Label:     "test strain",
                Species:   "Dictyostelium discoideum",
            },
        },
    }

    result := storage.CreateStrain(req)()

    F.Pipe1(
        result,
        E.Fold(
            func(err error) {
                require.Fail(t, "Expected success but got error: %v", err)
            },
            func(strain *stock.Strain) {
                require.Regexp(t, `^DBS\d{7}$`, strain.Data.Id)
                require.Equal(t, "test strain", strain.Data.Attributes.Label)
            },
        ),
    )
}
```

### Integration Tests

Test full gRPC service workflows:

```go
func TestStockService_CreateAndGetStrain(t *testing.T) {
    fix := setupServiceTest(t)
    defer fix.cleanup()

    ctx := context.Background()

    createReq := &stock.NewStrain{
        Data: &stock.NewStrain_Data{
            Type: "strain",
            Attributes: &stock.NewStrainAttributes{
                CreatedBy: "curator@dictybase.org",
                UpdatedBy: "curator@dictybase.org",
                Depositor: "John Doe",
                Label:     "axeA2 axeB2 axeC2",
                Species:   "Dictyostelium discoideum",
            },
        },
    }

    created, err := fix.service.CreateStrain(ctx, createReq)
    require.NoError(t, err)
    require.NotEmpty(t, created.Data.Id)

    getReq := &stock.StockId{Id: created.Data.Id}
    retrieved, err := fix.service.GetStrain(ctx, getReq)
    require.NoError(t, err)
    require.Equal(t, created.Data.Id, retrieved.Data.Id)
}
```

### Filter Engine Tests

Test all operators and combinations:

```go
func TestFilterParser_ComplexQuery(t *testing.T) {
    filter := "depositor===Costanza;created_at$>=2018-12-01"
    result := ParseFilter(filter)

    require.True(t, E.IsRight(result))

    strain := &stock.Strain{
        Data: &stock.Strain_Data{
            Attributes: &stock.StrainAttributes{
                Depositor: "Costanza",
                CreatedAt: timestamppb.New(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
            },
        },
    }

    jsonData := buildJSONIndex(strain)

    matches := F.Pipe2(
        result,
        E.Map[error](func(expr FilterExpression) bool {
            return expr.Evaluate(jsonData)
        }),
        E.GetOrElse(func() bool { return false }),
    )

    require.True(t, matches)
}
```

## Implementation Scope

### Phase 1: Core Infrastructure

1. Pebble storage interface and implementation
2. Basic CRUD operations (Create, Get, Update, Delete)
3. Sequential ID generation
4. CLI integration with stock subcommand

### Phase 2: Service Layer

1. gRPC service with IOEither pipelines
2. Progressive context building for all workflows
3. Result converters (Either → Tuple)
4. Production-matching error codes

### Phase 3: Filter Engine

1. Tokenizer with Option-based safe access
2. Parser with Either composition
3. Evaluator with function dispatch
4. All operators (string, numeric, date, array)

### Phase 4: Pagination

1. Reader monad for dependencies
2. Snapshot-based iteration
3. Cursor encoding/decoding
4. Sorting support

### Phase 5: Testing & Documentation

1. Unit tests for all layers
2. Integration tests for workflows
3. Filter parser comprehensive tests
4. Usage documentation

## Success Criteria

1. **All 9 RPC methods implemented** - CRUD + listing functionality
2. **Full filter parser** - All operators working correctly
3. **Production-matching behavior** - Same validation, errors, responses
4. **Deep fp-go integration** - IOEither, Reader, Either, Option throughout
5. **Comprehensive tests** - >90% coverage on core logic
6. **CLI integration** - Seamless subcommand in mock-grpc-server
7. **Documentation** - Clear usage examples and architecture docs

## Future Enhancements

1. **Persistent mode** - Support for on-disk Pebble storage
2. **Snapshot export/import** - Save/restore test data
3. **Performance metrics** - Request timing and throughput stats
4. **Health checks** - gRPC health checking protocol
5. **Multiple simultaneous servers** - Run stock + annotation mocks together
