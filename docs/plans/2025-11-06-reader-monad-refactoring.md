# Reader Monad Refactoring Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Refactor plasmid CSV parsing to use fp-go Reader monad for point-free composition and clean dependency injection.

**Architecture:** Replace manual currying with Reader monad pattern. Parse functions lift pure Either computations into ReaderEither. Enrichment functions use named helpers for business logic, wrapped in Reader for dependency access. All steps share uniform signature: `ParseContext → ReaderEither[Dependencies, error, ParseContext]`.

**Tech Stack:** Go 1.21+, github.com/IBM/fp-go (Reader, Either, Option), github.com/stretchr/testify

---

## Prerequisites

**Verify existing tests pass:**
```bash
gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...
```

Expected: All tests PASS

**File to modify:** `internal/datasource/csv/stockcenter/plasmid.go`

---

## Task 1: Add Core Types and Import

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:1-17`

**Step 1: Add ReaderEither import**

Add to imports section (line 10-14):
```go
import (
    A "github.com/IBM/fp-go/array"
    E "github.com/IBM/fp-go/either"
    F "github.com/IBM/fp-go/function"
    O "github.com/IBM/fp-go/option"
    RE "github.com/IBM/fp-go/readereither"  // NEW

    "github.com/dictyBase/modware-import/internal/datasource"
    csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)
```

**Step 2: Add ParseContext type**

Add after line 64 (after Plasmid struct):
```go
// ParseContext carries data through the parsing pipeline
type ParseContext struct {
    Record  []string
    Plasmid *Plasmid
}
```

**Step 3: Add Dependencies type**

Add after ParseContext:
```go
// Dependencies holds lookup services for enrichment
type Dependencies struct {
    Alookup StockAnnotatorLookup
    Plookup StockPubLookup
    Glookup StockGeneLookup
}
```

**Step 4: Add PlasmidParser type alias**

Add after Dependencies:
```go
// PlasmidParser represents a parsing step in the pipeline
type PlasmidParser = RE.ReaderEither[Dependencies, error, ParseContext]
```

**Step 5: Run tests to verify no breakage**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: All tests still PASS

**Step 6: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): add Reader monad types for plasmid parsing

Add ParseContext, Dependencies, and PlasmidParser type alias to support
Reader monad refactoring. Import fp-go readereither package.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 2: Convert parseId to Reader Pattern

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:106-121`

**Step 1: Replace parseId function**

Replace lines 106-121 with:
```go
// parseId extracts and sets the plasmid ID from record
func parseId(ctx ParseContext) PlasmidParser {
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
```

**Step 2: Run tests to verify no breakage**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: May FAIL because parseId signature changed. We'll fix call sites in later tasks.

**Step 3: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): convert parseId to Reader pattern

Change parseId from manual currying to Reader monad pattern.
Now accepts ParseContext and returns PlasmidParser.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 3: Convert parseName to Reader Pattern

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:135-150`

**Step 1: Replace parseName function**

Replace lines 135-150 with:
```go
// parseName extracts and normalizes plasmid name with prefix
func parseName(ctx ParseContext) PlasmidParser {
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
```

**Step 2: Run tests**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: Still may FAIL - we're converting incrementally

**Step 3: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): convert parseName to Reader pattern

Change parseName from manual currying to Reader monad pattern.
Now accepts ParseContext and returns PlasmidParser.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 4: Convert parseSummary to Reader Pattern

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:153-168`

**Step 1: Replace parseSummary function**

Replace lines 153-168 with:
```go
// parseSummary extracts plasmid summary from record
func parseSummary(ctx ParseContext) PlasmidParser {
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

**Step 2: Run tests**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: Still may FAIL - we're converting incrementally

**Step 3: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): convert parseSummary to Reader pattern

Change parseSummary from manual currying to Reader monad pattern.
Now accepts ParseContext and returns PlasmidParser.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 5: Add Named Helper Functions for Enrichment

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go` (add before line 171)

**Step 1: Add applyAnnotator helper**

Add before line 171 (before enrichWithAnnotator):
```go
// applyAnnotator applies annotator lookup to plasmid context
func applyAnnotator(ctx ParseContext, deps Dependencies) ParseContext {
    user, createdOn, updatedOn, ok := deps.Alookup.StockAnnotator(ctx.Plasmid.Id)
    if ok {
        ctx.Plasmid.User = user
        ctx.Plasmid.CreatedOn = createdOn
        ctx.Plasmid.UpdatedOn = updatedOn
    }
    return ctx
}
```

**Step 2: Add applyPublications helper**

Add after applyAnnotator:
```go
// applyPublications applies publication lookup to plasmid context
func applyPublications(ctx ParseContext, deps Dependencies) ParseContext {
    return F.Pipe2(
        deps.Plookup.StockPub(ctx.Plasmid.Id),
        O.FromPredicate(isNonEmptySlice),
        O.Fold(
            func() ParseContext { return ctx },
            func(pubs []string) ParseContext {
                filteredPubs := F.Pipe1(
                    pubs,
                    A.Filter(func(pub string) bool {
                        return pub != ""
                    }),
                )
                ctx.Plasmid.Publications = append(ctx.Plasmid.Publications, filteredPubs...)
                return ctx
            },
        ),
    )
}
```

**Step 3: Add applyGenes helper**

Add after applyPublications:
```go
// applyGenes applies gene lookup to plasmid context
func applyGenes(ctx ParseContext, deps Dependencies) ParseContext {
    return F.Pipe2(
        deps.Glookup.StockGene(ctx.Plasmid.Id),
        O.FromPredicate(isNonEmptySlice),
        O.Fold(
            func() ParseContext { return ctx },
            func(genes []string) ParseContext {
                filteredGenes := F.Pipe1(
                    genes,
                    A.Filter(func(gene string) bool {
                        return gene != ""
                    }),
                )
                ctx.Plasmid.Genes = append(ctx.Plasmid.Genes, filteredGenes...)
                return ctx
            },
        ),
    )
}
```

**Step 4: Run tests**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: Still may FAIL - helpers added but not yet used

**Step 5: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): add named helper functions for enrichment

Add applyAnnotator, applyPublications, and applyGenes helpers to separate
business logic from Reader plumbing. These will be wrapped in Reader functions
in the next step.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 6: Convert Enrichment Functions to Reader Pattern

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:171-240`

**Step 1: Replace enrichWithAnnotator**

Replace lines 171-181 (the curried enrichWithAnnotator) with:
```go
// enrichWithAnnotator enriches plasmid with annotator data using Reader
func enrichWithAnnotator(ctx ParseContext) PlasmidParser {
    return F.Pipe1(
        RE.Ask[Dependencies](),
        RE.Map(func(deps Dependencies) ParseContext {
            return applyAnnotator(ctx, deps)
        }),
    )
}
```

**Step 2: Replace enrichWithPublications**

Replace lines 189-213 (the curried enrichWithPublications) with:
```go
// enrichWithPublications enriches plasmid with publications using Reader
func enrichWithPublications(ctx ParseContext) PlasmidParser {
    return F.Pipe1(
        RE.Ask[Dependencies](),
        RE.Map(func(deps Dependencies) ParseContext {
            return applyPublications(ctx, deps)
        }),
    )
}
```

**Step 3: Replace enrichWithGenes**

Replace lines 216-240 (the curried enrichWithGenes) with:
```go
// enrichWithGenes enriches plasmid with genes using Reader
func enrichWithGenes(ctx ParseContext) PlasmidParser {
    return F.Pipe1(
        RE.Ask[Dependencies](),
        RE.Map(func(deps Dependencies) ParseContext {
            return applyGenes(ctx, deps)
        }),
    )
}
```

**Step 4: Run tests**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: Still may FAIL - functions converted but not yet integrated

**Step 5: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): convert enrichment functions to Reader pattern

Replace curried enrichWith* functions with Reader monad versions.
Each uses RE.Ask to access dependencies and calls named helper.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 7: Add parseAllFields Composition Function

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go` (add before Value method)

**Step 1: Add parseAllFields function**

Add before the Value() method (around line 254):
```go
// parseAllFields composes all field parsing functions
func parseAllFields(ctx ParseContext) PlasmidParser {
    return F.Pipe3(
        parseId(ctx),
        RE.Chain(parseName),
        RE.Chain(parseSummary),
    )
}
```

**Step 2: Run tests**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: Still may FAIL - function added but not yet used

**Step 3: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): add parseAllFields composition

Add parseAllFields to compose parse functions in point-free style.
Chains parseId, parseName, and parseSummary using RE.Chain.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 8: Update Value() Method to Use Reader Pipeline

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:255-270`

**Step 1: Replace Value() implementation**

Replace lines 255-270 (the current Value method) with:
```go
// Value gets a new Plasmid instance using fp-go Reader pattern
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
        RE.Chain(enrichWithAnnotator),
        RE.Chain(enrichWithPublications),
        RE.Chain(enrichWithGenes),
        RE.Map(func(ctx ParseContext) *Plasmid {
            return ctx.Plasmid
        }),
    )(deps) // Run with dependencies

    return result
}
```

**Step 2: Run tests to verify refactoring works**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: All tests should PASS now (refactoring complete, behavior preserved)

**Step 3: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): update Value() to use Reader pipeline

Replace manual pipeline with point-free Reader composition.
Create ParseContext and Dependencies, run pipeline with deps.

All tests pass - refactoring preserves existing behavior.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 9: Remove Old parsePlasmidFields Function

**Files:**
- Modify: `internal/datasource/csv/stockcenter/plasmid.go:243-252`

**Step 1: Delete parsePlasmidFields**

Delete lines 243-252 (the old parsePlasmidFields curried function):
```go
// DELETE THIS ENTIRE BLOCK:
// parsePlasmidFields is a curried function that composes all field parsers
var parsePlasmidFields = F.Curry2(
    func(record []string, plasmid *Plasmid) E.Either[error, *Plasmid] {
        return F.Pipe3(
            E.Of[error](plasmid),
            E.Chain(parseId(record)),
            E.Chain(parseName(record)),
            E.Chain(parseSummary(record)),
        )
    },
)
```

**Step 2: Run tests to verify removal didn't break anything**

Run: `gotestsum --format-hide-empty-pkg --format dots -- ./internal/datasource/csv/stockcenter/...`
Expected: All tests PASS (function no longer used)

**Step 3: Commit**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "refactor(stockcenter): remove old parsePlasmidFields function

Delete parsePlasmidFields as it's replaced by parseAllFields.
All tests pass - old code safely removed.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Task 10: Format Code and Run Final Verification

**Files:**
- `internal/datasource/csv/stockcenter/plasmid.go`

**Step 1: Format code**

Run: `gofumpt -w internal/datasource/csv/stockcenter/plasmid.go`
Expected: Code formatted according to gofumpt standards

**Step 2: Run linter**

Run: `golangci-lint run internal/datasource/csv/stockcenter/...`
Expected: No linting issues

**Step 3: Run all tests with verbose output**

Run: `gotestsum --format-hide-empty-pkg --format standard-verbose -- ./internal/datasource/csv/stockcenter/...`
Expected: All tests PASS with detailed output

**Step 4: Run all project tests to verify no side effects**

Run: `gotestsum --format-hide-empty-pkg --format dots`
Expected: All tests PASS across entire project

**Step 5: Commit formatting if changes made**

```bash
git add internal/datasource/csv/stockcenter/plasmid.go
git commit -m "style(stockcenter): format code after Reader refactoring

Run gofumpt to ensure consistent formatting.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## Verification Checklist

After completing all tasks, verify:

- [ ] All tests pass: `gotestsum --format-hide-empty-pkg --format dots`
- [ ] No lint issues: `golangci-lint run`
- [ ] Code formatted: `gofumpt -w .`
- [ ] parseId, parseName, parseSummary use Reader pattern
- [ ] Enrichment functions use named helpers with Reader
- [ ] parseAllFields composes parse functions point-free
- [ ] Value() uses Reader pipeline with ParseContext and Dependencies
- [ ] Old parsePlasmidFields removed
- [ ] All commits follow conventional commit format
- [ ] Behavior unchanged (tests prove this)

---

## Notes

**Point-free style achieved:** After initial `parseAllFields(ctx)`, all subsequent steps compose without explicit arguments:
```go
RE.Chain(enrichWithAnnotator),
RE.Chain(enrichWithPublications),
RE.Chain(enrichWithGenes),
```

**Named functions preferred:** Business logic in `applyAnnotator`, `applyPublications`, `applyGenes` - testable independently.

**Dependencies injected once:** Call `(deps)` at end to run entire pipeline with lookups.

**Backward compatible:** `Value()` signature unchanged - returns `E.Either[error, *Plasmid]` as before.
