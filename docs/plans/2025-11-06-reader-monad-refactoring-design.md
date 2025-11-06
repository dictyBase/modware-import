# Reader Monad Refactoring for Plasmid CSV Parsing

**Date:** 2025-11-06
**Status:** Design Complete
**File:** `internal/datasource/csv/stockcenter/plasmid.go`

## Overview

This refactoring replaces manual currying with fp-go's Reader monad pattern, enabling point-free composition and clean dependency injection.

## Goals

1. **Enable point-free composition** - Compose pipeline steps without passing arguments explicitly
2. **Improve readability** - Use named functions instead of inline closures
3. **Maximize composability** - Reorder, add, or remove pipeline steps easily
4. **Clean dependency injection** - Thread lookup services through enrichment functions

## Current State

Parse functions use manual currying:
```go
func parseId(record []string) func(*Plasmid) E.Either[error, *Plasmid]
```

Enrichment functions use `F.Curry2`:
```go
var enrichWithAnnotator = F.Curry2(
    func(alookup StockAnnotatorLookup, plasmid *Plasmid) E.Either[error, *Plasmid]
)
```

Pipeline requires explicit argument passing:
```go
F.Pipe3(
    parsePlasmidFields(plr.Record)(plasmid),
    E.Chain(enrichWithAnnotator(plr.alookup)),
    E.Chain(enrichWithPublications(plr.plookup)),
    E.Chain(enrichWithGenes(plr.glookup)),
)
```

## Proposed Design

### Core Types

**ParseContext** - Data flowing through pipeline:
```go
type ParseContext struct {
    Record  []string
    Plasmid *Plasmid
}
```

**Dependencies** - Lookup services injected at runtime:
```go
type Dependencies struct {
    Alookup StockAnnotatorLookup
    Plookup StockPubLookup
    Glookup StockGeneLookup
}
```

**Type alias** for clarity:
```go
type PlasmidParser = ReaderEither[Dependencies, error, ParseContext]
```

### Parse Functions

Parse functions return `ReaderEither` for uniform composition even though they need no dependencies:

```go
func parseId(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return RE.FromEither[Dependencies](
        F.Pipe1(
            getRecordField(ctx.Record, 0),
            O.Fold(
                func() E.Either[error, ParseContext] {
                    return E.Left[ParseContext](fmt.Errorf("missing id at index 0"))
                },
                func(id string) E.Either[error, ParseContext] {
                    ctx.Plasmid.Id = id
                    return E.Right[error](ctx)
                },
            ),
        ),
    )
}

func parseName(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return RE.FromEither[Dependencies](
        F.Pipe1(
            getRecordField(ctx.Record, 1),
            O.Fold(
                func() E.Either[error, ParseContext] {
                    return E.Left[ParseContext](fmt.Errorf("missing name at index 1"))
                },
                func(name string) E.Either[error, ParseContext] {
                    ctx.Plasmid.Name = F.Pipe1(name, ensurePlasmidPrefix)
                    return E.Right[error](ctx)
                },
            ),
        ),
    )
}

func parseSummary(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return RE.FromEither[Dependencies](
        F.Pipe1(
            getRecordField(ctx.Record, 2),
            O.Fold(
                func() E.Either[error, ParseContext] {
                    return E.Left[ParseContext](fmt.Errorf("missing summary at index 2"))
                },
                func(summary string) E.Either[error, ParseContext] {
                    ctx.Plasmid.Summary = summary
                    return E.Right[error](ctx)
                },
            ),
        ),
    )
}
```

### Enrichment Functions

**Named helper functions** separate business logic from Reader plumbing:

```go
func applyAnnotator(ctx ParseContext, deps Dependencies) ParseContext {
    user, createdOn, updatedOn, ok := deps.Alookup.StockAnnotator(ctx.Plasmid.Id)
    if ok {
        ctx.Plasmid.User = user
        ctx.Plasmid.CreatedOn = createdOn
        ctx.Plasmid.UpdatedOn = updatedOn
    }
    return ctx
}

func applyPublications(ctx ParseContext, deps Dependencies) ParseContext {
    return F.Pipe2(
        deps.Plookup.StockPub(ctx.Plasmid.Id),
        O.FromPredicate(isNonEmptySlice),
        O.Fold(
            func() ParseContext { return ctx },
            func(pubs []string) ParseContext {
                filteredPubs := A.Filter(func(pub string) bool { return pub != "" })(pubs)
                ctx.Plasmid.Publications = append(ctx.Plasmid.Publications, filteredPubs...)
                return ctx
            },
        ),
    )
}

func applyGenes(ctx ParseContext, deps Dependencies) ParseContext {
    return F.Pipe2(
        deps.Glookup.StockGene(ctx.Plasmid.Id),
        O.FromPredicate(isNonEmptySlice),
        O.Fold(
            func() ParseContext { return ctx },
            func(genes []string) ParseContext {
                filteredGenes := A.Filter(func(gene string) bool { return gene != "" })(genes)
                ctx.Plasmid.Genes = append(ctx.Plasmid.Genes, filteredGenes...)
                return ctx
            },
        ),
    )
}
```

**Enrichment functions** use Reader to access dependencies:

```go
func enrichWithAnnotator(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return F.Pipe1(
        RE.Ask[Dependencies](),
        RE.Map(func(deps Dependencies) ParseContext {
            return applyAnnotator(ctx, deps)
        }),
    )
}

func enrichWithPublications(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return F.Pipe1(
        RE.Ask[Dependencies](),
        RE.Map(func(deps Dependencies) ParseContext {
            return applyPublications(ctx, deps)
        }),
    )
}

func enrichWithGenes(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return F.Pipe1(
        RE.Ask[Dependencies](),
        RE.Map(func(deps Dependencies) ParseContext {
            return applyGenes(ctx, deps)
        }),
    )
}
```

### Point-Free Composition

**Compose parse functions:**
```go
func parseAllFields(ctx ParseContext) ReaderEither[Dependencies, error, ParseContext] {
    return F.Pipe3(
        parseId(ctx),
        RE.Chain(parseName),
        RE.Chain(parseSummary),
    )
}
```

**Main pipeline in `Value()` method:**
```go
func (plr *csvPlasmidReader) Value() E.Either[error, *Plasmid] {
    // Check for CSV reader errors first
    if plr.Err != nil {
        return E.Left[*Plasmid](plr.Err)
    }

    // Create initial context
    ctx := ParseContext{
        Record:  plr.Record,
        Plasmid: new(Plasmid),
    }

    // Create dependencies
    deps := Dependencies{
        Alookup: plr.alookup,
        Plookup: plr.plookup,
        Glookup: plr.glookup,
    }

    // Run the point-free pipeline
    result := F.Pipe4(
        parseAllFields(ctx),
        RE.Chain(enrichWithAnnotator),      // pure point-free!
        RE.Chain(enrichWithPublications),   // no arguments!
        RE.Chain(enrichWithGenes),          // beautiful!
        RE.Map(func(ctx ParseContext) *Plasmid {
            return ctx.Plasmid
        }),
    )(deps) // Run with dependencies

    return result
}
```

## Migration Strategy

1. Add new types: `ParseContext`, `Dependencies`, `PlasmidParser` type alias
2. Convert parse functions: `parseId`, `parseName`, `parseSummary` to return `ReaderEither`
3. Convert enrichment functions: Add named helpers (`applyAnnotator`, etc.), then wrap in Reader
4. Update `Value()` method: Use the new pipeline with `.Run(deps)`
5. Remove old code: Delete `parsePlasmidFields` curried function

## Error Handling

The `ReaderEither` type handles errors through the `Either` layer. When any step returns `E.Left[error]`, the pipeline short-circuits automatically. Existing error handling requires no changes.

## Testing Strategy

- **Parse functions:** Test with mock `ParseContext`, ignore dependencies
- **Enrichment helpers:** Test helpers (`applyAnnotator`) directly with mock lookups
- **Full pipeline:** Create `Dependencies` with test doubles

## Benefits

1. **Point-free composition** - Pipeline reads left-to-right without argument passing
2. **Dependency injection** - Provide lookups once at the end
3. **Testability** - Test named helpers independently
4. **Consistency** - All steps share one signature: `ParseContext → ReaderEither`
5. **Composability** - Reorder, add, or remove steps easily

## Trade-offs

- **Learning curve** - Reader monad is more abstract than currying
- **Indirection** - One more layer (Reader) between code and dependencies
- **Debugging** - Stack traces might be slightly harder to read with Reader abstraction

## Required Imports

```go
import (
    A "github.com/IBM/fp-go/array"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
    O "github.com/IBM/fp-go/option"
    RE "github.com/IBM/fp-go/readereither"  // NEW
)
```

## What Gets Removed

- `parsePlasmidFields` curried function (replaced by `parseAllFields`)
- `enrichWithAnnotator` curried variable (becomes named function)
- `enrichWithPublications` curried variable (becomes named function)
- `enrichWithGenes` curried variable (becomes named function)

## What Stays Unchanged

- `getRecordField` helper
- `ensurePlasmidPrefix` helper
- `isNonEmptySlice` helper
- `PlasmidReader` interface
- `Value()` method signature (still returns `E.Either[error, *Plasmid]`)
