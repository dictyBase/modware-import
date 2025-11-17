# Plasmid Loading fp-go Migration Design

**Date:** 2025-11-06
**Status:** Approved
**Scope:** Maximum fp-go adoption for plasmid loading functionality

## Executive Summary

This design migrates the plasmid loading system to use functional programming patterns from the `github.com/IBM/fp-go` library. The refactoring achieves:

- **IOEither** for side effects (database, API, file I/O)
- **Reader monad** for dependency injection
- **Either** for validation with error accumulation
- **Option** for nullable values
- **Streaming recursion** for memory-efficient processing
- **Hybrid architecture** maintaining registry for setup, fp-go for business logic

## Goals

1. **Maximum fp-go adoption** - Transform all business logic to functional patterns
2. **Error accumulation** - Collect all validation errors instead of failing fast
3. **Memory efficiency** - Stream process CSV records without loading all into memory
4. **Type safety** - Leverage fp-go's type system for compile-time guarantees
5. **Maintainability** - Clear separation of pure and impure code

## Architecture Overview

### Three-Layer Functional Architecture

```
┌─────────────────────────────────────────────────────┐
│  CLI Layer (internal/cli/stockcenter/plasmid.go)   │
│  - Registry initialization (existing pattern)       │
│  - Build Reader context from registry               │
│  - Execute IOEither workflow                        │
│  - Convert result to (result, error) tuple          │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│  Business Logic (internal/load/stockcenter/)        │
│  - Reader monad context: PlasmidEnv                 │
│  - IOEither-based pipeline using Do/Bind            │
│  - Streaming processor with recursive IOEither      │
│  - Validation accumulator using Either[[]error]     │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│  Data Source (internal/datasource/csv/stockcenter/) │
│  - Functional CSV reader wrapping iterator          │
│  - Option monad for nullable fields                 │
│  - Either for field parsing validation              │
│  - Pure transformation functions                    │
└─────────────────────────────────────────────────────┘
```

## Type System Design

### Reader Context

```go
// PlasmidEnv contains all dependencies for plasmid loading
type PlasmidEnv struct {
    Logger          *logrus.Entry
    APIClient       pb.StockServiceClient
    PlasmidReader   io.Reader
    AnnotatorReader io.Reader
    PubReader       io.Reader
    GeneReader      io.Reader
}

// Type aliases for cleaner signatures
type (
    PlasmidIO[A any]     = CRIOE.ContextReaderIOEither[PlasmidEnv, error, A]
    PlasmidEither[A any] = E.Either[error, A]
)
```

### Error Accumulation

```go
// ValidationError represents a single validation failure
type ValidationError struct {
    PlasmidID string
    Field     string
    Message   string
    Err       error
}

// ProcessingResult accumulates successes and errors
type ProcessingResult struct {
    Successes []string          // Successfully processed plasmid IDs
    Errors    []ValidationError // All validation/processing errors
}

// Semigroup instance for combining results
var ProcessingResultSemigroup = SG.MakeSemigroup(
    func(a, b ProcessingResult) ProcessingResult {
        return ProcessingResult{
            Successes: append(a.Successes, b.Successes...),
            Errors:    append(a.Errors, b.Errors...),
        }
    },
)
```

### Progressive Context Types (Do/Bind Pattern)

```go
// Stage 1: Lookups initialized
type WithLookups struct {
    Env        PlasmidEnv
    Annotator  source.StockAnnotatorLookup
    PubLookup  source.StockPubLookup
    GeneLookup source.StockGeneLookup
}

// Stage 2: CSV reader constructed
type WithReader struct {
    WithLookups
    CSVReader source.PlasmidReader
}

// Stage 3: Single plasmid being processed
type WithPlasmid struct {
    Env     PlasmidEnv
    Plasmid *source.Plasmid
}
```

### Curried Context Builders

```go
var (
    setLookups = F.Curry2(
        func(lookups LookupTuple, env PlasmidEnv) WithLookups {
            return WithLookups{
                Env:        env,
                Annotator:  lookups.Annotator,
                PubLookup:  lookups.PubLookup,
                GeneLookup: lookups.GeneLookup,
            }
        },
    )

    setReader = F.Curry2(
        func(reader source.PlasmidReader, ctx WithLookups) WithReader {
            return WithReader{WithLookups: ctx, CSVReader: reader}
        },
    )
)
```

## Core Patterns

### Pattern 1: Streaming Recursive Processor

**Problem:** Process potentially large CSV files without loading everything into memory.

**Solution:** Tail-recursive IOEither processor with accumulator.

```go
// streamProcessRecords processes CSV records recursively using IOEither
func streamProcessRecords(
    reader source.PlasmidReader,
) PlasmidIO[ProcessingResult] {
    return CRIOE.Map[PlasmidEnv](
        processNextRecord(reader, ProcessingResult{}),
    )
}

// Tail-recursive processor with accumulator
func processNextRecord(
    reader source.PlasmidReader,
    acc ProcessingResult,
) PlasmidIO[ProcessingResult] {
    return F.Pipe3(
        hasNextRecord(reader),
        CRIOE.Chain[PlasmidEnv](func(hasNext bool) PlasmidIO[ProcessingResult] {
            if !hasNext {
                // Base case: no more records, return accumulated result
                return CRIOE.Of[PlasmidEnv](acc)
            }
            // Recursive case: process one record and continue
            return processSingleRecord(reader, acc)
        }),
        CRIOE.MapLeft[PlasmidEnv, ProcessingResult](
            fperrors.OnError("stream processing failed"),
        ),
    )
}

// Process one record and recurse
func processSingleRecord(
    reader source.PlasmidReader,
    acc ProcessingResult,
) PlasmidIO[ProcessingResult] {
    return F.Pipe4(
        readPlasmidRecord(reader),
        CRIOE.Chain[PlasmidEnv](validateAndProcessPlasmid),
        CRIOE.Map[PlasmidEnv](mergeResult(acc)),
        CRIOE.Chain[PlasmidEnv](func(newAcc ProcessingResult) PlasmidIO[ProcessingResult] {
            // Tail recursion: continue with updated accumulator
            return processNextRecord(reader, newAcc)
        }),
    )
}
```

**Key Benefits:**
- Constant memory usage regardless of file size
- Purely functional recursion
- Composable with other IOEither operations

### Pattern 2: Validation with Error Accumulation

**Problem:** Traditional Go error handling stops at first error, losing information about subsequent failures.

**Solution:** Use `Either[[]ValidationError, T]` to accumulate all validation errors.

```go
// validatePlasmid performs all validation checks, accumulating errors
func validatePlasmid(plasmid *source.Plasmid) E.Either[[]ValidationError, *source.Plasmid] {
    return F.Pipe4(
        E.Of[[]ValidationError](plasmid),
        E.Chain(validateUserAssignment),
        E.Chain(validateRequiredFields),
        E.Chain(validatePublications),
        E.MapLeft[*source.Plasmid](func(errs []ValidationError) []ValidationError {
            // Enrich errors with plasmid ID context
            return A.Map(errs, func(err ValidationError) ValidationError {
                err.PlasmidID = plasmid.Id
                return err
            })
        }),
    )
}

// Individual validators return Either with error arrays
func validateUserAssignment(p *source.Plasmid) E.Either[[]ValidationError, *source.Plasmid] {
    if len(p.User) == 0 {
        return E.Left[*source.Plasmid]([]ValidationError{{
            Field:   "User",
            Message: "user assignment required",
        }})
    }
    return E.Right[[]ValidationError](p)
}

func validateRequiredFields(p *source.Plasmid) E.Either[[]ValidationError, *source.Plasmid] {
    errors := []ValidationError{}

    if p.Id == "" {
        errors = append(errors, ValidationError{Field: "Id", Message: "required"})
    }
    if p.Name == "" {
        errors = append(errors, ValidationError{Field: "Name", Message: "required"})
    }

    if len(errors) > 0 {
        return E.Left[*source.Plasmid](errors)
    }
    return E.Right[[]ValidationError](p)
}
```

**Key Benefits:**
- Collect all validation errors in single pass
- Better user experience (see all issues at once)
- Composable validators using Either chain

### Pattern 3: Option for Nullable Values

**Problem:** Dealing with optional CSV fields and potential missing data.

**Solution:** Use `Option[T]` instead of nil pointers or empty strings.

```go
// Safe array access using Option
func getRecordField(record []string, index int) O.Option[string] {
    if index < 0 || index >= len(record) {
        return O.None[string]()
    }
    return O.Some(record[index])
}

// Parse name with Option-based transformation
func parseName(record []string) func(*Plasmid) E.Either[error, *Plasmid] {
    return func(p *Plasmid) E.Either[error, *Plasmid] {
        return F.Pipe2(
            getRecordField(record, 1),
            O.Match(
                func() E.Either[error, *Plasmid] {
                    return E.Left[*Plasmid](fmt.Errorf("missing name at index 1"))
                },
                func(name string) E.Either[error, *Plasmid] {
                    p.Name = F.Pipe1(name, ensurePlasmidPrefix)
                    return E.Right[error](p)
                },
            ),
        )
    }
}

// Pure transformation using Option
func ensurePlasmidPrefix(name string) string {
    return F.Pipe2(
        O.Some(name),
        O.Filter(func(n string) bool {
            return strings.HasPrefix(n, "p")
        }),
        O.GetOrElse(F.Constant(fmt.Sprintf("p%s", name))),
    )
}
```

**Key Benefits:**
- Explicit handling of missing values
- No nil pointer dereferences
- Composable with functional pipelines

### Pattern 4: IOEither for API Interactions

**Problem:** gRPC API calls have side effects and can fail.

**Solution:** Wrap all API calls in `IOEither` for explicit effect tracking.

```go
// checkPlasmidExists wraps gRPC call in IOEither
func checkPlasmidExists(id string) PlasmidIO[bool] {
    return CRIOE.FromIOEither[PlasmidEnv](
        IOE.TryCatchError(func() (bool, error) {
            return F.Pipe1(
                CRIOE.Ask[PlasmidEnv](),
                CRIOE.Map[PlasmidEnv](func(env PlasmidEnv) bool {
                    _, err := env.APIClient.GetPlasmid(
                        context.Background(),
                        &pb.StockId{Id: id},
                    )
                    if err != nil && status.Code(err) == codes.NotFound {
                        return false
                    }
                    return err == nil
                }),
            )
        }),
    )
}

// createPlasmidAPI wraps create operation
func createPlasmidAPI(plasmid *source.Plasmid) PlasmidIO[string] {
    return F.Pipe3(
        CRIOE.Ask[PlasmidEnv](),
        CRIOE.Chain[PlasmidEnv](func(env PlasmidEnv) PlasmidIO[string] {
            return CRIOE.FromIOEither[PlasmidEnv](
                IOE.TryCatchError(func() (string, error) {
                    attr := populateExistingPlasmidAttributes(env.Logger, plasmid)
                    _, err := env.APIClient.LoadPlasmid(
                        context.Background(),
                        &pb.ExistingPlasmid{
                            Data: &pb.ExistingPlasmid_Data{
                                Type:       "plasmid",
                                Id:         plasmid.Id,
                                Attributes: attr,
                            },
                        },
                    )
                    if err != nil {
                        return "", fmt.Errorf("failed to create plasmid %s: %w", plasmid.Id, err)
                    }
                    env.Logger.Debugf("created plasmid %s", plasmid.Id)
                    return plasmid.Id, nil
                }),
            )
        }),
        CRIOE.MapLeft[PlasmidEnv, string](
            fperrors.OnError("create plasmid API call failed"),
        ),
    )
}
```

**Key Benefits:**
- Explicit side effect tracking
- Composable with pure transformations
- Automatic error propagation

### Pattern 5: Reader Monad for Dependency Injection

**Problem:** Need to pass multiple dependencies through many layers.

**Solution:** Use `ContextReaderIOEither` to inject dependencies functionally.

```go
// Main workflow uses Reader context
func loadPlasmidWorkflow() PlasmidIO[ProcessingResult] {
    return F.Pipe1(
        CRIOE.Do[PlasmidEnv, ProcessingResult](
            // Step 1: Initialize all lookups
            CRIOE.Bind(
                "lookups",
                func() PlasmidIO[LookupTuple] {
                    return CRIOE.FromIOEither[PlasmidEnv](
                        IOE.Chain(
                            IOE.Ask[PlasmidEnv](),
                            func(env PlasmidEnv) IOE.IOEither[error, LookupTuple] {
                                return initAllLookups(env)
                            },
                        ),
                    )
                },
            ),
            // Step 2: Create CSV reader with lookups
            CRIOE.Bind(
                "reader",
                func(ctx map[string]interface{}) PlasmidIO[source.PlasmidReader] {
                    lookups := ctx["lookups"].(LookupTuple)
                    return createPlasmidReader(lookups)
                },
            ),
            // Step 3: Stream process all records
            CRIOE.Chain(func(ctx map[string]interface{}) PlasmidIO[ProcessingResult] {
                reader := ctx["reader"].(source.PlasmidReader)
                return streamProcessRecords(reader)
            }),
        ),
        CRIOE.MapLeft[PlasmidEnv, ProcessingResult](
            fperrors.OnError("plasmid loading workflow failed"),
        ),
    )
}

// Execute workflow by providing environment
func LoadPlasmid(cmd *cobra.Command, args []string) error {
    env := buildPlasmidEnv()

    result := F.Pipe3(
        loadPlasmidWorkflow(),
        CRIOE.RunReaderIOEither(env),  // Inject dependencies
        convertToGoResult,
        logFinalStats(env.Logger),
    )

    return result.F2
}
```

**Key Benefits:**
- No global state access in business logic
- Testable without mocking global registry
- Dependencies explicit in type signature

## Implementation Plan

### Phase 1: Data Source Layer
**File:** `internal/datasource/csv/stockcenter/plasmid.go`

1. Add fp-go imports (E, O, F, A)
2. Refactor `Value()` method to return `E.Either[error, *Plasmid]`
3. Implement field parsers using `Option` for safe array access
4. Add `ensurePlasmidPrefix` as pure function
5. Implement lookup enrichment using curried functions
6. Update tests to verify Either-based parsing

### Phase 2: Business Logic Layer
**File:** `internal/load/stockcenter/plasmid.go`

1. Add fp-go imports (IOE, CRIOE, E, O, F, A, SG, fperrors)
2. Define `PlasmidEnv` and type aliases
3. Define `ValidationError` and `ProcessingResult` types
4. Implement `ProcessingResultSemigroup`
5. Implement progressive context types (`WithLookups`, `WithReader`)
6. Implement curried context builders
7. Implement `initAllLookups` with IOEither
8. Implement streaming processor:
   - `streamProcessRecords`
   - `processNextRecord` (recursive)
   - `processSingleRecord`
   - `hasNextRecord`
   - `readPlasmidRecord`
9. Implement validation pipeline:
   - `validatePlasmid`
   - `validateUserAssignment`
   - `validateRequiredFields`
   - `validatePublications`
10. Implement API interactions:
    - `checkPlasmidExists`
    - `createPlasmidAPI`
    - `updatePlasmidAPI`
11. Implement `processPlasmidWithEnv` with Do/Bind
12. Implement `loadPlasmidWorkflow` orchestration
13. Update tests using IOEither test patterns

### Phase 3: CLI Layer
**File:** `internal/cli/stockcenter/plasmid.go`

1. Add fp-go imports (CRIOE, T, F)
2. Implement `buildPlasmidEnv` from registry
3. Update `LoadPlasmid` to use functional workflow
4. Implement `convertToGoResult` (IOEither → Tuple2)
5. Implement `formatValidationErrors`
6. Implement `logFinalStats` with curried logger
7. Keep `setPlasmidPreRun` unchanged (registry pattern)
8. Verify integration with existing Cobra commands

### Phase 4: Testing & Validation

1. Run `gotestsum --format-hide-empty-pkg --format dots`
2. Verify all tests pass
3. Run `golangci-lint run` for code quality
4. Test with sample CSV files
5. Verify error accumulation behavior
6. Check memory usage with large files
7. Validate logging output

## Testing Strategy

### Unit Tests

```go
func TestValidatePlasmid(t *testing.T) {
    tests := []struct {
        name          string
        plasmid       *source.Plasmid
        expectedRight bool
        expectedErrs  int
    }{
        {
            name: "valid plasmid",
            plasmid: &source.Plasmid{
                Id:   "DBP001",
                Name: "pDV10",
                User: "curator@example.com",
            },
            expectedRight: true,
            expectedErrs:  0,
        },
        {
            name: "missing user",
            plasmid: &source.Plasmid{
                Id:   "DBP001",
                Name: "pDV10",
            },
            expectedRight: false,
            expectedErrs:  1,
        },
        {
            name: "multiple errors",
            plasmid: &source.Plasmid{},
            expectedRight: false,
            expectedErrs:  3,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validatePlasmid(tt.plasmid)

            if tt.expectedRight {
                require.True(t, E.IsRight(result))
            } else {
                require.True(t, E.IsLeft(result))
                errors := E.GetLeft(result)
                require.Len(t, errors, tt.expectedErrs)
            }
        })
    }
}
```

### Integration Tests

```go
func TestLoadPlasmidWorkflow(t *testing.T) {
    // Setup test environment
    env := PlasmidEnv{
        Logger:    logrus.NewEntry(logrus.New()),
        APIClient: mockAPIClient,
        // ... test readers
    }

    // Execute workflow
    result := F.Pipe2(
        loadPlasmidWorkflow(),
        CRIOE.RunReaderIOEither(env),
    )

    // Verify result
    F.Pipe1(
        result(),
        E.Match(
            func(err error) {
                t.Fatalf("workflow failed: %v", err)
            },
            func(procResult ProcessingResult) {
                require.GreaterOrEqual(t, len(procResult.Successes), 1)
                // Verify error accumulation worked
            },
        ),
    )
}
```

## Migration Risks & Mitigation

### Risk 1: Learning Curve
**Impact:** Team unfamiliar with functional patterns
**Mitigation:**
- Comprehensive code comments
- Reference existing fp-go examples in `.claude/patterns/`
- Pair programming sessions

### Risk 2: Performance
**Impact:** Functional overhead vs imperative code
**Mitigation:**
- Benchmark before/after
- Streaming processor prevents memory issues
- Tail-call recursion optimized by compiler

### Risk 3: Debugging Complexity
**Impact:** Functional stack traces harder to read
**Mitigation:**
- Extensive logging at each pipeline stage
- Use `fperrors.OnError` for context enrichment
- Test each function independently

### Risk 4: Breaking Changes
**Impact:** Changes to public APIs
**Mitigation:**
- Keep CLI interface unchanged
- Internal refactoring only
- Comprehensive test coverage

## Success Criteria

1. ✅ All existing tests pass
2. ✅ No regression in functionality
3. ✅ Error accumulation working (batch validation)
4. ✅ Memory usage stable with large CSV files
5. ✅ Code passes `golangci-lint run`
6. ✅ Logging provides detailed processing statistics
7. ✅ Type safety verified at compile time

## References

- [fp-go library](https://github.com/IBM/fp-go)
- [fp-go patterns](.claude/patterns/fp-go-patterns.md)
- [fp-go examples](.claude/patterns/fp-go-concepts/)
- [CLAUDE.md coding conventions](.claude/CLAUDE.md)

## Appendix: Import Structure

```go
import (
    "context"
    "fmt"
    "strings"
    "time"

    A "github.com/IBM/fp-go/array"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
    O "github.com/IBM/fp-go/option"
    T "github.com/IBM/fp-go/tuple"
    IOE "github.com/IBM/fp-go/ioeither"
    CRIOE "github.com/IBM/fp-go/context/readerioeither"
    SG "github.com/IBM/fp-go/semigroup"
    fperrors "github.com/IBM/fp-go/errors"

    pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
    source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
    "github.com/dictyBase/modware-import/internal/registry"
    regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"

    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)
```
