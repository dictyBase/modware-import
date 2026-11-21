# GoldenBraid Export Subcommand (`goldenbraid export`) Implementation Plan

## Table of Contents

- [Global Constraints](#global-constraints)
- [Review Focus](#review-focus)
- [File Structure](#file-structure)
- [Tasks](#tasks)
  - [Task 1: Dependencies, AST Gates, and `seqmodel` Package](#task-1-dependencies-ast-gates-and-seqmodel-package)
  - [Task 2: SnapGene Binary TLV and XML Parser](#task-2-snapgene-binary-tlv-and-xml-parser)
  - [Task 3: GenBank and FASTA Converters](#task-3-genbank-and-fasta-converters)
  - [Task 4: Circular Plasmid SVG Renderer](#task-4-circular-plasmid-svg-renderer)
  - [Task 5: Input Resolution and Local Export Pipeline](#task-5-input-resolution-and-local-export-pipeline)
  - [Task 6: S3 Input and Output with `s3blob`](#task-6-s3-input-and-output-with-s3blob)
  - [Task 7: CLI Subcommand and Main Integration](#task-7-cli-subcommand-and-main-integration)
  - [Task 8: Kubernetes Job and Dagu Workflow](#task-8-kubernetes-job-and-dagu-workflow)

> **For agentic workers:** Use `superpowers:subagent-driven-development` or `superpowers:executing-plans` to implement this plan task by task. Each step uses a checkbox to record progress.

**Goal:** Export GoldenBraid assemblies (finished constructs) from [assemblies.zip](https://dicty.labs.duke.edu/assets/goldenbraid/assemblies.zip) as FASTA, GenBank, and SVG. Backbones, parts, and legacy groups are out of scope.

**Architecture:** Resolve `--input` and `--output` independently. Each location has a `*blob.Bucket`: `fileblob` for local paths and `file:` URLs, or `s3blob` for `s3:` URLs. Read one construct, open an archive in memory, or list a directory or prefix. For each selected construct, write three objects under the output root. Log counts and per-file errors before the exit decision.

The native Go `snapgene` parser decodes TLV packets and XML directly into one `genbank.Genbank`. Each encoder reads the same record. Its feature list includes primer sites. Do not add a Python parser or a second record model. The `seqformat` and `svgmap` packages produce GenBank bytes, FASTA bytes, and circular SVG maps in pure `Either` pipelines. The `gbexport` package runs lazy `IOEither` effects over state structs, and `gbcli` exposes a point-free `urfave/cli/v2` action. One test file in `internal/goldenbraid/gates_test.go` enforces `pipelinecheck` and `loomcheck` across all feature packages.

**Tech Stack:** Go `1.26.9`, `github.com/IBM/fp-go/v2` (`v2.13.2`), `github.com/dictyBase/fp-go-loom` (`v0.5.0`), `github.com/bebop/poly` (`v0.31.2`), `github.com/urfave/cli/v2` (`v2.27.7`), `gocloud.dev/blob` (`v0.49.0`, `fileblob` and `s3blob` drivers), AWS SDK for Go v2 (`service/s3`, `credentials`), `github.com/stretchr/testify/require`, `github.com/testcontainers/testcontainers-go/modules/minio`, `log/slog`.

**Spec:** Visual architecture guide in [goldenbraid_visual_guide.html](file:///Users/sba964/.gemini/antigravity/brain/6d122c78-c754-4f9f-8b1d-a28641c1e4ca/goldenbraid_visual_guide.html) and the Duke assembly archive in the Goal paragraph. The visual guide is not available in this environment. Coverage against the guide is unknown.

## Global Constraints

- Work inside the current module `github.com/dictyBase/modware-import` (`go 1.26.9` in `go.mod`). Do not run `go mod init`.
- Use `github.com/urfave/cli/v2` (`v2.27.7`).
- Keep all `internal/goldenbraid/` packages free of `github.com/urfave/cli/v2` imports except `internal/goldenbraid/gbcli`.
- Import only the house monad allowlist: `O` (`option`), `E` (`either`), `IO` (`io`), and `IOE` (`ioeither`). Alias `IOE` to `ioeither` only. Do not import `Reader`, `ReaderIOEither`, `Result`, or `IOResult`.
- Import `github.com/dictyBase/fp-go-loom` (`v0.5.0`) for predicates, `eithererr.ConstErr`, `eitherparse.ParseInt`, `matchopt` (`MO`), `ioeutils.ToEither`, `pipelinecheck`, and `loomcheck`.
- Use `github.com/bebop/poly@v0.31.2` as the only bio library. The alternatives evaluated (`biogo/hts`, `biogo/biogo`, and `gonomics`) do not meet this export contract: BAM/SAM-only read paths, missing GenBank support, or `log.Fatal`-based error handling. Use `encoding/xml` for SVG. Do not add a graphics library.
- Use package-level `L.MakeLens` values for external poly structs. A setter changes a copy and returns it. Do not generate lenses for external types.
- Annotate local domain and state structs with `// fp-go:Lens`. Generate `gen_lens.go` with `//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go`.
- Construct structs from `T{}` through `F.PipeN(T{}, lens.Set(...))` and read fields with `lens.Get`. Write zero populated struct literals in any file that imports `fp-go` or `fp-go-loom`.
- Write no local closures and no local `F.Flow` or `F.Pipe` setup variables inside functions. Write each function body as one flat `F.PipeN`.
- Make helpers accept plain arguments and bind them with `F.Bind1st` or `F.Bind2nd`. Never return a closure from a helper. Put reusable combinators in package-level `var` blocks.
- Audit unexported helpers before each commit. Inline single-use `F.PipeN` helpers, replace ignored-state adapters with `F.Constant1[State](val)`, inline single-use pointer helpers, and remove pass-through `P.Fold` wrappers.
- Lift pure `E.Either` functions into `IOE.IOEither` pipes with `IOE.ChainEitherK`. Never write `IOE.FromEither` wrapper helpers.
- Put only one raw call inside each `IOE.TryCatchError` callback. Wrap errors in `IOE.MapLeft` with `ER.OnError` (`ER "github.com/IBM/fp-go/v2/errors"`), and project fields in `IOE.Map`.
- Describe terminal effects as data, converge all branches before the terminal sink, and execute the sink once. Encode artifact bytes purely before any bucket write.
- Combine slices with `A.ArrayConcatAll` or `A.Flatten`. Never call `slices.Concat` or `append(a, b...)`.
- Write no imperative `if`, `for`, or `switch` outside `TryCatch*` callbacks. Use `P.Fold`, `O.FromPredicate`, `MO.Match`, `E.FromPredicate`, `A.Reduce`, and recursion.
- Derive fixed-root path mappers once with `root := F.Pipe1(dir, F.Flip(FILE.Join))`, then apply `A.Map(root)`.
- Declare package sentinel errors and match them with `errors.Is`. Never match `err.Error()` text.
- Force `IOEither` at the CLI edge with `ioeutils.ToEither` and `E.ToError`. Do not write a redundant `E.Fold` that discards the `Right` value.
- Define named constants for numeric literals (`32 << 20`, `16`, `640`). The `.golangci.yml` file enables `mnd`, `funlen` (80 lines, 50 statements), `lll`, `cyclop`, and `revive`.
- Use `IOE.WithResource` for both buckets. Write only `fasta/<stem>.fasta`, `genbank/<stem>.gbk`, and `svg/<stem>.svg` under the output root.
- Run all Go verification checks after each Go task: `go generate`, `golangci-lint fmt`, `golangci-lint run ./...`, `go test ./...`, `go build ./...`, `gopls`, and four `rg` checks.

## Review Focus

- **Circular features across base 1:** A single XML segment with `start > end` crosses the origin. Make sure that the parser splits it at base 1. The renderer must keep both parts on one ring in source order.
- **DNA flags:** Examine bits `0x01` and `0x02` independently. Both `0x0f` and `0x1f` denote circular, double-stranded DNA. Ignore the other bits.
- **Oversized TLV packet lengths:** Packet lengths greater than `32 MiB` (`MaxPacketBytes`) must fail with `ErrPacketTooLarge` before slice allocation.
- **History-heavy files:** Use synthetic repeated history packets (`0x0b`). Apply the size limit per packet, not per file. Skip each unused packet even if its ID occurs more than once. Do not add a Go stack frame for each packet.
- **Dense short features on SVG rings:** The adjacent short tags (`GGGSx4 Linker`, `H2Bv3`) in the dense fixture must stay in the `viewBox`. Their Y coordinates must differ by at least `16 px`.
- **Partial failure:** Continue after a corrupt file. Log each failed key and its error once, log the counts once, then return `Left(ErrPartialExportFailure)`.
- **Duplicate stems:** Reserve the stem before export. Fail the second source with `ErrDuplicateStem`. Do not replace the first source's artifacts.
- **Mixed schemes:** Export S3 input to local output. The output opener must use the output scheme, not the input scheme.
- **Bucket lifetime:** Make sure that `Close` runs on failure as well as success. Keep the list operation, reads, and output writes inside `IOE.WithResource`.
- **Depth and names:** Include root files and direct children of first-level folders such as `GoldenBraid Assemblies/`. Preserve the space. Skip deeper entries, archive directory records, `__MACOSX/` metadata, and `._` names. Do not list a sibling prefix such as `runs/ab/` for `runs/a/`.

---

## File Structure

- `internal/goldenbraid/doc.go` — Declares package `goldenbraid` with a package doc comment only.
- `internal/goldenbraid/gates_test.go` — Enforces all 12 `pipelinecheck` rules and `loomcheck.Require` over `Roots: []string{"."}`.
- `internal/goldenbraid/seqmodel/types.go` — Declares `// fp-go:Lens` structs (`Color`, `ExportSummary`) and sentinel errors.
- `internal/goldenbraid/seqmodel/gen_lens.go` — Holds generated lenses for `seqmodel`.
- `internal/goldenbraid/seqmodel/validate.go` — Checks upper-case bases against `ACGT`. Returns the validated sequence.
- `internal/goldenbraid/seqmodel/monoid.go` — Defines `ExportSummaryMonoid` with `A.ArrayConcatAll` and `M.ConcatAll`.
- `internal/goldenbraid/polylens/lenses.go` — Defines shared external lenses for poly record, metadata, location, feature, and FASTA fields.
- `internal/goldenbraid/snapgene/tlv.go` — Scans binary SnapGene TLV packets with an indexed fold and applies `MaxPacketBytes` to each packet.
- `internal/goldenbraid/snapgene/xml.go` — Unmarshals `<Features>`, `<Notes>`, and `<Primers>` XML packets with `eitherparse.ParseInt` and `MO.Match`.
- `internal/goldenbraid/snapgene/parser.go` — Defines `ParseInput` (`// fp-go:Lens`) and exposes pure `Parse` (`E.Either`) for source bytes.
- `internal/goldenbraid/snapgene/gen_lens.go` — Holds generated lenses for `snapgene`.
- `internal/goldenbraid/seqformat/convert.go` — Writes GenBank and FASTA bytes from `genbank.Genbank`.
- `internal/goldenbraid/svgmap/config.go` — Defines `// fp-go:Lens` `Config`, `RenderInput`, named layout constants, and `En.Endomorphism[Config]` options.
- `internal/goldenbraid/svgmap/gen_lens.go` — Holds generated lenses for `svgmap`.
- `internal/goldenbraid/svgmap/layout.go` — Assigns concentric feature rings and relaxes label Y positions with `A.Reduce` over sorted slices.
- `internal/goldenbraid/svgmap/render.go` — Emits accessible circular `<svg>` bytes as `E.Either[error, []byte]`.
- `internal/goldenbraid/gbexport/state.go` — Declares `// fp-go:Lens` `Params`, `FileBundle`, and `State` structs plus sentinel errors and file constants.
- `internal/goldenbraid/gbexport/gen_lens.go` — Holds generated lenses for `gbexport`.
- `internal/goldenbraid/gbexport/effects.go` — Writes each encoded artifact with `bucket.WriteAll` and adds context to errors.
- `internal/goldenbraid/gbexport/source.go` — Resolves URLs and input types, opens buckets, reads archive entries, and lists two levels of keys.
- `internal/goldenbraid/gbexport/source_integration_test.go` — Contains MinIO tests for the S3 opener and shared export pipeline.
- `internal/goldenbraid/gbexport/pipeline.go` — Composes `Pipeline(seed State) IOE.IOEither[error, State]` for all input types, collects per-file results with `IO.TraverseArraySeq`, and logs failures and counts.
- `internal/goldenbraid/gbcli/flags.go` — Defines `exportFlags()` with `urfave/cli/v2` flag struct literals and zero `fp-go` imports.
- `internal/goldenbraid/gbcli/export.go` — Defines `ExportFlags()`, `toExportState`, and point-free `RunExport(c *cli.Context) error` with `logger.GetCliSlogHandler` and `stockcenter` flag helpers.
- `cmd/goldenbraid/main.go` (Modify) — Registers the `export` subcommand and converts `slices.Concat` calls to `A.ArrayConcatAll`.
- `just-modules/goldenbraid.just` (Modify) — Submits the export Job.
- `dagu-goldenbraid-export.yaml` — Runs the export Job and retrieves its logs.
- `schemas/goldenbraid-export-params.json` — Defines the inputs for the workflow.

---

## Tasks

### Task 1: Dependencies, AST Gates, and `seqmodel` Package

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Create: `internal/goldenbraid/doc.go`
- Create: `internal/goldenbraid/gates_test.go`
- Create: `internal/goldenbraid/polylens/lenses.go`
- Create: `internal/goldenbraid/seqmodel/types.go`
- Create: `internal/goldenbraid/seqmodel/gen_lens.go`
- Create: `internal/goldenbraid/seqmodel/validate.go`
- Create: `internal/goldenbraid/seqmodel/monoid.go`
- Create: `internal/goldenbraid/seqmodel/seqmodel_test.go`

**Interfaces:**
- Consumes: `github.com/bebop/poly/io/genbank` and `github.com/bebop/poly/io/fasta` (`v0.31.2`), `github.com/IBM/fp-go/v2` (`v2.13.2`), and `github.com/dictyBase/fp-go-loom` (`v0.5.0`).
- Produces:
  - Sentinels: `ErrInvalidColor`, `ErrEmptyName`, `ErrEmptySequence`, `ErrInvalidDNABase`, `ErrSegmentOutOfBounds`
  - Local types (`// fp-go:Lens`):
    - `type Color struct{ R, G, B uint8 }`
    - `ExportSummary`: `Total int`, `Succeeded int`, `Failed int`, `Errors []error`. Do not add JSON tags.
  - `func ParseColor(s string) E.Either[error, Color]` and `func (c Color) Hex() string`
  - `func ValidateSequence(sequence string) E.Either[error, string]`
  - `func ReduceExportResults(results []E.Either[error, string]) ExportSummary`

- [ ] **Step 1: Pin poly and `fp-go-loom` and create the AST gates**

Run:

```bash
go get github.com/bebop/poly@v0.31.2
go get github.com/dictyBase/fp-go-loom@v0.5.0
```

Create `internal/goldenbraid/doc.go` with package `goldenbraid` and a doc comment only.

Create `internal/goldenbraid/gates_test.go` in package `goldenbraid_test` with `Roots: []string{"."}`. Do not use `pipelinecheck.ModuleRoot(t)` because current repository code outside `internal/goldenbraid` will fail these gates. Set `IsEntrypoint` to match prefixes `"Run"` or `"run"` so the gate includes `gbcli.RunExport`. Enable `RequirePointFreeSeed`, `RequirePointFreeBranching`, and `RequireTrivialDelegate`. Then call `loomcheck.Require(t, loomcheck.Config{Roots: []string{"."}})` and all 12 `pipelinecheck` rules:
- `Require`, `RequireBareIf`, `RequireNoIfErrInTryCatch`, `RequireNoNonRawTryCatchCallback`
- `RequireSafePrintf`, `RequireNoRedundantFold`, `RequireNoHandRolledFileJoinInArrayMap`
- `RequireNoImperativeLoop`, `RequireNoImperativeBranch`, `RequireNoDirectApply`, `RequireNoPopulatedStruct`, `RequireNoLocalPipeClosure`

Later tasks inherit this file and add no gate files.

- [ ] **Step 2: Write unit tests in `internal/goldenbraid/seqmodel/seqmodel_test.go`**

```go
package seqmodel_test

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/require"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
)

func TestValidateSequence(t *testing.T) {
	require.Equal(t, E.Right[error]("ACGT"), seqmodel.ValidateSequence("acGt"))
	require.ErrorIs(t, E.ToError[string](seqmodel.ValidateSequence("ATGCN")), seqmodel.ErrInvalidDNABase)
	require.ErrorIs(t, E.ToError[string](seqmodel.ValidateSequence("")), seqmodel.ErrEmptySequence)
}

func TestReduceExportResultsUsesMonoidConcat(t *testing.T) {
	corruptErr := errors.New("corrupt input")
	results := []E.Either[error, string]{
		E.Right[error]("valid.dna"),
		E.Left[string](corruptErr),
	}

	got := seqmodel.ReduceExportResults(results)
	sumL := seqmodel.MakeExportSummaryLenses()
	require.Equal(t, 2, sumL.Total.Get(got))
	require.Equal(t, 1, sumL.Succeeded.Get(got))
	require.Equal(t, 1, sumL.Failed.Get(got))
}
```

In the same test file, pass `#0A0B0C` to `ParseColor`. Assert that `Hex` on the parsed color returns `#0a0b0c` (lower-case digits with leading zeros). Test a short value and a value with a non-hexadecimal digit. Both must return `ErrInvalidColor`.

- [ ] **Step 3: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/seqmodel/...`
Expected: FAIL because `internal/goldenbraid/seqmodel` files do not exist.

- [ ] **Step 4: Implement `types.go`, `gen_lens.go`, `validate.go`, and `monoid.go`**

Generate lenses only for `Color` and `ExportSummary` in `seqmodel`. Put the house generator directive in `types.go`. Keep the summary fields free of JSON tags. In `validate.go`, reject an empty sequence. Accept only ASCII `ACGTacgt` and return the upper-case sequence. Use `E.FromPredicate` and sentinel error factories.

In `polylens/lenses.go`, use `L.MakeLens` for the fields used on `genbank.Genbank`, `genbank.Meta`, `genbank.Locus`, `genbank.Feature`, `genbank.Location`, and `fasta.Fasta`. Use fresh maps and slices for qualifiers and locations. Do not mutate an input through a shared map or slice.

In `ParseColor`, accept only `#` followed by six hexadecimal digits. Return `ErrInvalidColor` for any other input. Return lower-case `#rrggbb` from `Hex`.

In `monoid.go`, define `ExportSummaryMonoid` with `A.ArrayConcatAll` for `Errors` and addition for each count, then fold results with `M.ConcatAll`.

- [ ] **Step 5: Run verification gates and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go'
```

Expected: All tests, AST gates, linter, build, and `gopls` pass with zero output. Under `internal/goldenbraid`, accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Expect no other matches there. CLI struct literals in `cmd/goldenbraid/main.go` are an accepted exception. Task 7 removes `slices.Concat` from `cmd/goldenbraid/main.go`.

```bash
git add go.mod go.sum internal/goldenbraid/doc.go internal/goldenbraid/gates_test.go internal/goldenbraid/polylens internal/goldenbraid/seqmodel
git commit -m "feat(goldenbraid): add AST gates and seqmodel domain package"
```

---

### Task 2: SnapGene Binary TLV and XML Parser

**Files:**
- Create: `internal/goldenbraid/snapgene/tlv.go`
- Create: `internal/goldenbraid/snapgene/xml.go`
- Create: `internal/goldenbraid/snapgene/parser.go`
- Create: `internal/goldenbraid/snapgene/gen_lens.go`
- Create: `internal/goldenbraid/snapgene/parser_test.go`
- Existing: `internal/goldenbraid/snapgene/testdata/README.md`
- Existing: `internal/goldenbraid/snapgene/testdata/pDGB_A1_act15p_mNeonGreen_act8t.dna`
- Existing: `internal/goldenbraid/snapgene/testdata/pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t.dna`
- Existing: `internal/goldenbraid/snapgene/testdata/pDGB_A1_act15p_DDB_G0278295_mNeonGreen_act8t.dna`

**Interfaces:**
- Consumes: `genbank.Genbank`, `polylens`, `seqmodel.ValidateSequence`, `seqmodel.ParseColor`, `eitherparse.ParseInt`, `MO.Match`, `A.ReduceWithIndex`, `predord.StrEq`
- Produces:
  - `const MaxPacketBytes uint32 = 32 << 20`
  - Sentinels: `ErrInvalidMagicCookie`, `ErrPacketTooLarge`, `ErrTruncatedPacket`, `ErrMissingDNAPacket`
  - `// fp-go:Lens` `type ParseInput struct { SourceFile string; Name string; Data []byte }`
  - `MakeParseInputLenses()` for `SourceFile`, `Name`, and `Data`
  - `// fp-go:Lens` `type Packet struct { ID byte; Body []byte }`
  - `func Packets(data []byte) E.Either[error, []Packet]`
  - `const LocusNameMax = 16`
  - `// fp-go:Lens` `type ParseResult struct { Record genbank.Genbank; Warnings []error }`
  - `func Parse(in ParseInput) E.Either[error, ParseResult]`
  - Warning sentinels: `ErrInvalidFeature`, `ErrInvalidPrimer`

- [ ] **Step 1: Write unit tests in `internal/goldenbraid/snapgene/parser_test.go`**

Put the tests in package `snapgene_test`. Keep the synthetic tests. Build packets with `encodePacket(id byte, body []byte)`. Write one ID byte, four bytes for the big-endian body length, then the body.
1. `TestParseValidCircularSnapGeneDNA` — Build cookie `0x09` (`"SnapGene"` plus 6 version bytes). Add DNA block `0x00` (`flags = 0x0f`, `"GGAGCAAAATTTTGTATTTT"`). Add XML packets `0x0A`, `0x06`, and `0x05`. Set `directionality="1"` on `Feature` and `color="#31849b"` on its first standard `Segment`. Assert `Meta.Locus.Circular`, `MoleculeType == "DNA"`, a decimal `SequenceLength`, and an upper-case sequence. Check the label, `ApEinfo_fwdcolor == "#31849b"`, and the non-complemented location. Use a stem longer than 16 characters. Check the short locus name, full `Meta.Name`, definition, and comment.
2. `TestParseOriginWrapRange` — On a 20-base circular record, encode one `<Segment range="15-4" type="standard" color="noColor"/>`. Assert `Location.Join` with leaves `[14,20)` and `[0,4)` in order. Assert no color qualifier.
3. `TestParseSkipsUnknownPacketIDs` — Include packets `0x02`, `0x07`, `0x08`, and `0x11`. Assert `E.IsRight(got)`.
4. `TestParseRejectsOversizedPacketWithSentinel` — Encode packet `0x00` with big-endian length `0x7FFFFFFF`. Assert `require.ErrorIs(t, err, snapgene.ErrPacketTooLarge)`.
5. `TestParseRejectsInvalidCookieWithSentinel` — Start stream with `0x00`. Assert `require.ErrorIs(t, err, snapgene.ErrInvalidMagicCookie)`.
6. `TestParseManyHistoryPackets` — Insert 1000 `0x0b` packets between a valid cookie and DNA packet. Assert that `Parse` returns the expected sequence.
7. `TestParseSkipsInvalidFeature` — On a 20-base sequence, put XML range `1-21` between two valid features. Its converted location `[0,21)` exceeds the sequence length. Match both `ErrInvalidFeature` and `seqmodel.ErrSegmentOutOfBounds` on the warning. Expect `Right` with both valid features in source order. Expect exactly one warning. Make sure that `errors.Is(warning, ErrInvalidFeature)` is true.
8. Test two distinct binding sites on one primer, plus a simplified duplicate. Expect two `primer_bind` features with the expected half-open bounds and `Complement` flags.
9. Include one invalid primer and one valid primer. Expect `Right` with the valid primer and one warning that wraps `ErrInvalidPrimer`.

Read each golden input from `testdata/`. Assert the counts in the table. Examine the flags byte with `Packets`. Assert an upper-case sequence, circular locus, and one `primer_bind` site in `Record.Features`. The Features column counts source annotations only: assert total feature counts of 9, 14, and 12, respectively. Compare the leaves of each annotation with the converted standard XML ranges. Exclude both gaps in the last fixture.

| Test | Fixture | Bases | Flags | Features | Source gaps | Primers |
|---|---|---:|---|---:|---:|---:|
| `TestParseRealAssembly` | `pDGB_A1_act15p_mNeonGreen_act8t.dna` | 3666 | `0x0f` | 8 | 0 | 1 |
| `TestParseRealDenseTags` | `pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t.dna` | 4932 | `0x0f` | 13 | 0 | 1 |
| `TestParseRealGapSegments` | `pDGB_A1_act15p_DDB_G0278295_mNeonGreen_act8t.dna` | 4458 | `0x0f` | 11 | 2 | 1 |

The archive contains 156 `.dna` entries, all under `GoldenBraid Assemblies/`. All 156 assemblies are circular. Fourteen have gap segments. None have an origin wrap. Use `encodePacket` for the wrap and history tests.

- [ ] **Step 2: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/snapgene/...`
Expected: FAIL because the parser is not implemented.

- [ ] **Step 3: Implement `tlv.go`, `xml.go`, `parser.go`, and `gen_lens.go`**

In `tlv.go`, use `A.ReduceWithIndex` over the input bytes as a fold over offsets. Carry the next header offset and decoded packet state in an `Either`. Use a named reducer and `P.Fold`: keep the state at offsets inside a body. Read a header only at the next header offset. Keep the first `Left` unchanged. This scans the input once without recursion, a separate offset array, or an imperative loop in the parser.

Each header must contain five bytes. Read the body length as a big-endian `uint32`. If the length is more than `MaxPacketBytes`, return `ErrPacketTooLarge` before you allocate memory or take a slice. Reject an incomplete header or body with `ErrTruncatedPacket`. Advance past the entire body even for an unused ID. Do not limit the file size. Accept repeated packet IDs.

| Packet ID | Body | Action |
|---|---|---|
| `0x09` | `SnapGene` plus six version bytes (examples: `0001000d000e`, `0001000e000e` in hex) | The first packet must be this cookie. Examine the signature and size, not one fixed version. |
| `0x00` | Flags byte, then sequence bytes | Reject a body without a flags byte. Set circular when `flags & 0x01 != 0`. Bit `0x02` identifies double-stranded input. Do not add a strand field to the record. Pass the bases to `ValidateSequence`. The base count is the body length minus one. |
| `0x0a`, `0x05`, `0x06` | Features, primers, notes XML respectively | Use the XML rules in this step. Use the DNA length to resolve feature ranges, independent of packet order. |
| `0x02`, `0x03`, `0x07`, `0x08`, `0x0b`, `0x0d`, `0x0e`, `0x10`, `0x11`, `0x1b` | Observed but unused | Skip each occurrence. Apply the same length checks and skip rule to unknown IDs. |

In `xml.go`, read `name`, `type`, and optional `directionality` from `Feature`. Preserve UTF-8 names and the source type (for example, `CDS`, `gene`, `mRNA`, `tRNA`, `misc_feature`, `primer_bind`, `promoter`, `rep_origin`, `terminator`). Set `Feature.Type` from `type` and `Attributes["label"]` from `name`. Use `MO.Match` on direction: set `Location.Complement` only for `2`.

- Read `range`, `type`, and `color` from each `Segment`. Skip `type="gap"`.
- Read color only from the first standard segment. For absent color or `noColor`, omit `ApEinfo_fwdcolor`. For `Some` parsed color, set the qualifier to `Color.Hex()`.
- Parse the one-based inclusive `s-e` with `eitherparse.ParseInt`. Convert to a `genbank.Location` leaf `[s-1,e)`. For a circular wrap, split into `[s-1,len)` and `[0,e)`. Reject a wrap on linear DNA. Check each source endpoint before the split, then require `0 <= Start < End <= len` on each leaf. Wrap `ErrSegmentOutOfBounds` for a bounds failure.
- Use one leaf as the location. For two or more leaves, set `Join` and `SubLocations` on the parent. Put `Complement` on the whole location, not on each child.
- Keep XML order even with `swappedSegmentNumbering`.

The SVG renderer does not need the segment attributes `name` and `translated` or feature qualifiers such as `<Q name="note"><V text="…"/></Q>`.

In packet `0x05`, skip `BindingSite` elements with `simplified="1"`. For each remaining site, create a `genbank.Feature` of type `primer_bind` with the primer name as its label. Convert zero-based inclusive `location="a-b"` to `[a,b+1)`. Check the same leaf bounds. Set `Complement` for `boundStrand="1"`, not `"0"`. Place the site features after the XML features in the same slice, in source order.

Decode each feature and primer separately. If an entry has invalid XML, range, or color, skip that entry. Record one warning that wraps `ErrInvalidFeature` or `ErrInvalidPrimer`. Keep the valid entries. Return `Right` if the file has no fatal error. Do not convert DNA validation or packet framing failures to warnings.

Build metadata in `Parse`, not in the writers. Reject an empty input stem with `ErrEmptyName`. Use these fields:

| Field | Source and rule |
|---|---|
| `Meta.Name`, `Meta.Locus.Name` | Keep the full stem in `Meta.Name`. Limit the locus name to its first 16 characters. |
| `Meta.Locus.SequenceLength`, `MoleculeType`, `Circular` | Use the decimal base count as a string, `"DNA"`, and the circular flag from the DNA packet. |
| `Meta.Definition` | Use `<full stem>, <lower-case Type> construct.` from notes packet `0x06`. Without `Type`, use `<full stem>.` |
| `Meta.Other["COMMENT"]` | Use `SnapGene UUID <UUID>; created <Created>; modified <LastModified>; transformed into <TransformedInto>.` Remove `™`. Omit parts with no value. Omit the key if no parts remain. Do not add other keys. |

In `parser.go`, generate lenses for local `ParseInput`, `Packet`, and `ParseResult`. Use `polylens` to set metadata, sequence, and features on a zero `genbank.Genbank`. Return the record and warnings. Keep packet framing and DNA errors fatal.

- [ ] **Step 4: Run verification gates and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid -g '!**/*_test.go'
```

Expected: All tests, AST gates, linter, build, and `gopls` pass. Accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Expect no other matches.

```bash
git add internal/goldenbraid/snapgene
git commit -m "feat(goldenbraid): add pure-Go SnapGene .dna TLV and XML parser"
```

---

### Task 3: GenBank and FASTA Converters

**Files:**
- Create: `internal/goldenbraid/seqformat/convert.go`
- Create: `internal/goldenbraid/seqformat/convert_test.go`

**Interfaces:**
- Consumes: `genbank.Genbank`, `polylens`, and `github.com/bebop/poly/io/fasta`. Do not import the parser.
- Produces: `func ToGenBankBytes(r genbank.Genbank) E.Either[error, []byte]` and `func ToFASTABytes(r genbank.Genbank) E.Either[error, []byte]`.

- [ ] **Step 1: Write converter tests in `convert_test.go`**

1. Build a record with the shared lenses. Include a short locus name, the full stem, a definition, and a COMMENT. Add a reverse joined feature with a label and color. Add a primer site without color. Call `ToGenBankBytes` twice. Assert identical bytes, exactly one FEATURES heading, and `ApEinfo_fwdcolor` before `label` on the colored feature. The order assertion must not depend on Go map iteration.
2. Read the bytes with `genbank.Parse(bytes.NewReader(data))`. Compare sequence, circular flag, locus name, definition, COMMENT, and feature order. Check coordinates, reverse orientation, labels, optional color, and primer type. Compare `genbank.BuildLocationString` on each parsed location with the expected expression. For the reverse joined feature, expect `complement(join(...))`.
3. Call `ToFASTABytes` on the same record. Assert a header with the full stem, not the short locus name. Join the sequence lines and compare with `Record.Sequence`.

- [ ] **Step 2: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/seqformat/...`
Expected: FAIL because `internal/goldenbraid/seqformat` does not exist.

- [ ] **Step 3: Implement `convert.go`**

Use the shared Features lens to clear `Features` on a value copy. Lift `genbank.Build` with `E.Eitherize1`. The pinned writer emits an empty FEATURES heading even on this copy. Replace the exact line pair `FEATURES             Location/Qualifiers\nORIGIN\n` with the new FEATURES block followed by `ORIGIN\n`. Return a descriptive `Left` if the pair is missing. Keep the header and sequence bytes unchanged. The parser sets no `BaseCount` and at most one `Meta.Other` entry (COMMENT), so the pair is adjacent and the header has no map-order variation.

Start the replacement with the same FEATURES heading. Map features in slice order. For each line, write five spaces, the feature type left-aligned in 16 columns, then `genbank.BuildLocationString(location)` and a newline. Sort qualifier keys in ascending string order. Write one line per key in that order. Indent each `/key="value"` by 21 spaces. Double embedded quotes in values. Do not use `genbank.BuildFeatureString`: it ranges over a map. Combine the lines with pure folds or joins.

For FASTA, set `fasta.Fasta.Name` to `Meta.Name` and `Sequence` to the validated sequence. Call `fasta.Build` on a one-element slice. Use `E.Eitherize1` for the builder and `E.MapLeft` with `ER.OnError` for build errors.

- [ ] **Step 4: Run verification gates and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid -g '!**/*_test.go'
```

Expected: All tests, AST gates, linter, build, and `gopls` pass. Accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Expect no other matches.

```bash
git add internal/goldenbraid/seqformat
git commit -m "feat(goldenbraid): add GenBank and FASTA converters"
```

---

### Task 4: Circular Plasmid SVG Renderer

**Files:**
- Create: `internal/goldenbraid/svgmap/config.go`
- Create: `internal/goldenbraid/svgmap/gen_lens.go`
- Create: `internal/goldenbraid/svgmap/layout.go`
- Create: `internal/goldenbraid/svgmap/render.go`
- Create: `internal/goldenbraid/svgmap/render_test.go`

**Interfaces:**
- Consumes: `genbank.Genbank`, `polylens`, `encoding/xml`, `En "github.com/IBM/fp-go/v2/endomorphism"`, `M "github.com/IBM/fp-go/v2/monoid"`
- Produces:
  - Named constants: `DefaultSizePx = 640`, `MinLabelGapPx = 16`, `LabelMarginPx = 16`, `ArrowAngleDeg = 2.5`
  - `// fp-go:Lens` `type Config struct { Size int; ShowPrimers bool; MaxLabels int }`
  - `type Option = En.Endomorphism[Config]`
  - `func WithSize(px int) Option`
  - `func WithShowPrimers(show bool) Option`
  - `func WithMaxLabels(limit int) Option`
  - `// fp-go:Lens` `type RenderInput struct { Record genbank.Genbank; Options []Option }`
  - `func Render(in RenderInput) E.Either[error, []byte]`

- [ ] **Step 1: Write unit tests in `internal/goldenbraid/svgmap/render_test.go`**

Add these unit tests in `internal/goldenbraid/svgmap/render_test.go`:
1. `TestRenderProducesAccessibleValidSVG` — Unmarshal output bytes with `encoding/xml`. Assert `<svg role="img" viewBox="0 0 640 640">`, a `<title>` with plasmid name and bp count, and `<path>` elements with `data-feature-id`.
2. `TestAllocateRingsPacksOverlappingFeatures` — Assert an outer `gene` (`[99,900)`) lands on ring `0`, and `CDS` (`[99,850)`) plus `HA` (`[850,900)`) share ring `1`.
3. `TestRenderHandlesCircularOriginFeature` — Assert a circular feature with joined leaves `[11899,12000)` and `[0,150)` renders valid finite SVG coordinates.
4. `TestLayoutLabelsSpacesDenseClusterInsideViewBox` — Parse `../snapgene/testdata/pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t.dna` and render its labels. Make sure that the dense `GGGSx4 Linker` and `H2Bv3` labels are present. Assert `16 <= y <= 624` and at least `16 px` between adjacent labels on each side.
5. Render a record with a labeled `primer_bind` twice. Assert its path and label are absent by default and present with `WithShowPrimers(true)`. Check the fill for both a colored and an uncolored feature.

- [ ] **Step 2: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/svgmap/...`
Expected: FAIL because `internal/goldenbraid/svgmap` does not exist.

- [ ] **Step 3: Implement `config.go`, `gen_lens.go`, `layout.go`, and `render.go`**

Add `//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go` in `config.go` and run `go generate ./internal/goldenbraid/...`. Define `WithSize`, `WithShowPrimers`, and `WithMaxLabels` as lens setters (`En.Endomorphism[Config]`) folded onto `defaultConfig()` with `M.ConcatAll(En.Monoid[Config]())`. Use named constants for all numeric geometry values to satisfy `mnd`.

- Render circular maps only. Read `Meta.Name` for the title and `Sequence` for the base count. Derive `data-feature-id` from the original feature index before you filter the features.
- Filter `Type == "primer_bind"` before layout unless `ShowPrimers` is true. Use the half-open leaves of `Location` for arcs. Keep joined leaves on one ring. Do not draw across gaps. Use `Complement` for reverse orientation.
- Read label text from `Attributes["label"]`. Use `Attributes["ApEinfo_fwdcolor"]` for the fill if the attribute exists. Otherwise, use a named default fill.
- In `layout.go`, avoid imperative loops and branches. Sort angular spans largest to smallest, then assign concentric rings with `A.Reduce`. Next, sort left and right labels by ideal Y, then relax Y coordinates with `A.Reduce` to enforce `MinLabelGapPx`.
- In `render.go`, define typed SVG elements with XML tags. Build them with local lenses and encode with `xml.Marshal` from `encoding/xml` lifted into `Either`. Do not concatenate XML strings. Emit clockwise polar arcs from `-math.Pi/2` with `ArrowAngleDeg` arrowheads and `<polyline>` leader lines.

- [ ] **Step 4: Run verification gates and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid -g '!**/*_test.go'
```

Expected: All tests, AST gates, linter, build, and `gopls` pass. Accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Expect no other matches.

```bash
git add internal/goldenbraid/svgmap
git commit -m "feat(goldenbraid): add pure-Go circular plasmid SVG renderer"
```

- [ ] **Step 5: Get a visual review before Task 5**

Render the three Task 2 fixtures to a temporary directory. Present all three SVG files to the user. Ask the user to examine the arcs, labels, and colors. Wait for the review before you start Task 5.

---

### Task 5: Input Resolution and Local Export Pipeline

**Files:**
- Modify: `go.mod` and `go.sum`
- Create: `internal/goldenbraid/gbexport/source.go`
- Create: `internal/goldenbraid/gbexport/state.go`
- Create: `internal/goldenbraid/gbexport/gen_lens.go`
- Create: `internal/goldenbraid/gbexport/effects.go`
- Create: `internal/goldenbraid/gbexport/pipeline.go`
- Create: `internal/goldenbraid/gbexport/pipeline_test.go`

**Interfaces:**
- Consumes: `snapgene.Parse`, `seqformat.ToFASTABytes`, `seqformat.ToGenBankBytes`, `svgmap.Render`, `seqmodel.ReduceExportResults`, `archive/zip`, `net/url`, `blob`, and `fileblob`.
- Produces: `Pipeline(seed State) IOE.IOEither[error, State]` and the sentinels `ErrNoDNAFiles`, `ErrPartialExportFailure`, `ErrDuplicateStem`, `ErrUnknownInputType`, and `ErrUnsupportedScheme`.
- Constants: `FastaDir = "fasta"`, `GenBankDir = "genbank"`, `SVGDir = "svg"`.
- Lens types: `Params` holds `Input`, `InputType`, `Output`, `S3Endpoint`, `S3Region`, `AccessKey`, `SecretKey` (strings), and `ShowPrimers bool`. `State` holds `Ctx context.Context`, `Logger *slog.Logger`, `Params Params`, and `Summary seqmodel.ExportSummary`. Keep open bucket handles inside the resource scopes, not in the returned state.
- `FileBundle` holds `Stem string` and the encoded `FASTA`, `GenBank`, and `SVG` byte slices. A per-file error carries the source key and wraps the cause for `errors.Is`.

- [ ] **Step 1: Write pipeline tests in `pipeline_test.go`**

Run `go get gocloud.dev@v0.49.0`. Use `fileblob`, `t.TempDir()`, the three assembly fixtures, and a captured `slog` handler. These tests do not need Docker.

| Case | Input or setup | Assertions |
|---|---|---|
| Single construct | Pass a `.dna` file directly. Put another fixture beside it. | Export only the named stem. Read `fasta/<stem>.fasta`, `genbank/<stem>.gbk`, and `svg/<stem>.svg`. Assert one success. Check the full stem in FASTA and the primer site in GenBank, even with `ShowPrimers` false. |
| Partial failure | Put two valid files and one corrupt file in a directory. | Both valid bundles exist. Assert `ErrPartialExportFailure`, one summary `(3, 2, 1)`, and one failure log with the corrupt key and error. |
| Parser warning | Use a valid file with one invalid feature. | Expect one warning log and one success. |
| Partial write | Make the SVG write fail after the FASTA write. | The FASTA stays in place. Count one file failure. |
| Directory depth | Put files at the root, one folder deep, and two folders deep. | Export the first two only. |
| Assembly archive | Build a zip with `archive/zip`. Put all three fixtures under `GoldenBraid Assemblies/`. Include a directory entry, `__MACOSX/._x.dna`, and a two-deep `.dna` entry. | Export exactly three bundles. In a separate case, include a root `.dna` entry and make sure that it exports. |
| Explicit type | Give the zip a name without `.zip`. Select `InputType = "zip"`. | Export the same three stems. |
| Empty input | Use an empty folder. | Return `ErrNoDNAFiles`. Do not write files. |
| Duplicate stem | Use different fixtures at `a/same.dna` and `b/same.dna`. | Count one success and one failure. Match the recorded cause with `errors.Is(err, ErrDuplicateStem)`. Make sure that the FASTA sequence matches the first file. |
| Overwrite | Seed `genbank/<stem>.gbk` with stale bytes. | The GenBank output replaces the stale data. |
| Cancellation | Cancel `State.Ctx` before a single-file read. | The per-file error wraps `context.Canceled`. The batch returns `ErrPartialExportFailure`. |
| Dispatch | Use an unknown input type or an `https:` input or output. | Match `ErrUnknownInputType` or `ErrUnsupportedScheme`, respectively. |

- [ ] **Step 2: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/gbexport/...`
Expected: FAIL because the pipeline does not exist.

- [ ] **Step 3: Resolve inputs and run the batch inside the bucket lifetimes**

Annotate the state and domain types with `// fp-go:Lens`. Add the lens generator to `state.go`. Run `go generate ./internal/goldenbraid/...`. Keep CLI types out of this package.

- Parse both locations with `url.Parse`. Use `MO.Match` on each scheme. For an empty scheme or `file`, select `fileblob`. Until Task 6 adds `s3`, reject other schemes with `ErrUnsupportedScheme`.
- Use `MO.Match` on the input type. Accept `dna`, `zip`, `dir`, and `auto` (default). Reject other values with `ErrUnknownInputType`. For `auto`, inspect the extension of the last path element. Give an explicit type precedence over the extension.

| Type | Auto selection | Local input bucket root | Selection inside bucket |
|---|---|---|---|
| `dna` | `.dna` | Parent directory | Read the base name only. |
| `zip` | `.zip` | Parent directory | Read the named object once. Open it with `zip.NewReader(bytes.NewReader(data), int64(len(data)))`. Do not extract files to disk. |
| `dir` | Other | Input directory | List from the root with an empty prefix. |

- Open the input with `IOE.WithResource`. Use `fileblob.OpenBucket(dir, nil)` for local input. Resolve candidate entries before you open the output.
- For a directory or S3 prefix, list the root with `bucket.List` and `Delimiter: "/"`. Use `Prefix: ""` for the bucket root. For a nonempty prefix, put a slash at the end.
- List each first-level sub-prefix once with `Prefix: sub` and `Delimiter: "/"`. Ignore all sub-prefixes from the second list.
- Build `blob.ListOptions` with external lenses. Drain the iterator with `iter.Next(ctx)` until `io.EOF`. Keep the raw drain exception `// fp-go:allow-loop blob iterator drain is one raw effect`. End the list at `io.EOF`.

A list failure stops the run. It is not a per-file failure.

- Select `.dna` files at the root and one folder deep. Apply the same depth limit to archive member names.
- Skip directories as files. Skip `__MACOSX/` trees and names that start with `._`. Preserve spaces in folder names.
- Keep each selected archive member's `*zip.File` so that its read does not depend on a name lookup. Close each member reader after use.
- Use `predarrays.IsNonEmpty` in `IOE.ChainEitherK` to return `ErrNoDNAFiles` for an empty selection.

In the input scope, acquire the output with a second `IOE.WithResource`. For local output use `fileblob.OpenBucket(dir, &fileblob.Options{CreateDir: true})`. Build these options with a lens in fp-go code. Release each acquired bucket with `Close` on success and failure. Wrap acquisition, reads, writes, and release with the house raw-effect adapters and `IOE.MapLeft`.

Sort selected sources by key (stable for equal archive names). Derive the stem from the base name without its extension. Reserve each stem in order before any write. Mark each later occurrence as a per-file `ErrDuplicateStem`, even if the first export fails. Do not use an output existence check: old output from an earlier run must not block export.

- For a bucket object, use `bucket.ReadAll(ctx, key)` with `State.Ctx`. For a zip member, read from its open member reader.
- Examine `context.Cause(ctx)` before you parse the bytes, even after a successful local read.
- Set the `ParseInput` lenses for source name, stem, and bytes. Use `IOE.ChainEitherK` for `Parse` and the pure converters. Encode all three artifacts before the first write. Use `WithShowPrimers` for the SVG option.

- Use `path.Join(prefix, formatDir, stem + extension)` for each output key. For local output, use an empty prefix.
- Call `bucket.WriteAll(ctx, key, data, nil)` for each artifact. Do not implement a second replacement policy.
- Count a success only after all three writes succeed.

`WriteAll` replaces the object if the key already exists. `fileblob` uses a temporary file and a rename. Write FASTA, GenBank, and SVG in that order. If a later write fails, keep earlier artifacts. Do not roll back a partial write.

- Immediately after `Parse`, use `IOE.ChainFirstIOK[error](s.logWarnings)`. Log each warning with its source key. Extract `ParseResult.Record` once. Pass it to both converters and set the same record on `svgmap.RenderInput`. Use `ShowPrimers` only for SVG. Always include primer sites in GenBank.
- Return the key on success. On failure, use `IOE.MapLeft[string, error](ER.OnError("export %q", key))`, then `IOE.TapLeft[string, error](s.logFailure)`. Log the enriched error once. Keep the `Left` for the batch result. Do not log it again during reduction.
- Process one source at a time. Use `IO.TraverseArraySeq`, the sequential form of `IO.TraverseArray`, to collect `[]E.Either[error, string]`. Do not start concurrent exports. Lift the results with `IOE.FromIO[error]`. Reduce them into `State.Summary`.
- Use `IOE.ChainFirstIOK[error](State.logSummary)` once. Log only `total`, `succeeded`, and `failed`. Build each log effect with `IO.FromImpure` and the logger from `State`. For `logFailure`, lift the effect with `IOE.FromIO[error]` for `TapLeft`. Pass `s.logWarnings` and `s.logFailure` as method values and `State.logSummary` as a method expression. Do not bind the state with `F.Bind1st` or return a closure from a log helper.
- Then use `E.FromPredicate` in `IOE.ChainEitherK` to return `Left(ErrPartialExportFailure)` when `Failed > 0`. Examine wrapped causes at the per-file result or summary boundary before this final decision.

- [ ] **Step 4: Run verification gates and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid -g '!**/*_test.go'
```

Expected: All tests, AST gates, linter, build, and `gopls` pass. Accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Expect no other matches.

```bash
git add go.mod go.sum internal/goldenbraid/gbexport
git commit -m "feat(goldenbraid): export assembly files and archives to format folders"
```

---

### Task 6: S3 Input and Output with `s3blob`

**Files:**
- Modify: `go.mod` and `go.sum`
- Modify: `internal/goldenbraid/gbexport/source.go`
- Create: `internal/goldenbraid/gbexport/source_test.go`
- Create: `internal/goldenbraid/gbexport/source_integration_test.go`

**Interfaces:**
- Consumes: `gocloud.dev/blob/s3blob`, AWS SDK v2 `service/s3` and `credentials`, `Params`, and `State.Ctx`.
- Produces: an `s3` opener for either resource scope in Task 5. The archive reader, two-level list operation, writer, and error policy do not change.

- [ ] **Step 1: Write S3 unit and integration tests**

Run `go get gocloud.dev/blob/s3blob@v0.49.0`. In `source_test.go`, examine the options without a server. Assert `BaseEndpoint` points to `http://localhost:9000`, `UsePathStyle` is true, and the configured region and credentials reach the client. Make sure that `s3://bucket/runs/a` resolves to bucket `bucket` and key `runs/a`.

In `source_integration_test.go`, set `//go:build integration`. Start MinIO with `tcminio.Run` as in `internal/load/stockcenter/plasmid_ontology_integration_test.go`. Create a bucket and seed objects with `WriteAll`. Use a fresh bucket or clean output prefix per case.

| Input | Output | Assertions |
|---|---|---|
| `s3://bucket/prefix/assemblies.zip` with three assembly fixtures | `s3://bucket/out` | Read the three stems under each of `out/fasta/`, `out/genbank/`, and `out/svg/`. |
| `s3://bucket/prefix` with a valid root file, a valid first-level file, and one corrupt file | `s3://bucket/out` | Assert six valid artifacts, counts `(3, 2, 1)`, the corrupt key in the failure log, and `ErrPartialExportFailure`. Ignore two-deep and sibling-prefix objects. |
| Same S3 zip | Local temporary directory | Read all three local format folders. No S3 output is necessary. |

- [ ] **Step 2: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/gbexport/...`
Expected: FAIL because the S3 options builder does not exist.

- [ ] **Step 3: Add the S3 opener to `source.go`**

Define package-level `L.MakeLens` lenses for the external `s3.Options` fields. Start with `s3.Options{}` and apply the four setters with `F.Pipe4`:

| Field | Value |
|---|---|
| `BaseEndpoint` | Pointer to `Params.S3Endpoint` (`http://<s3-server>:<s3-server-port>`) |
| `Region` | `Params.S3Region` |
| `Credentials` | `credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")` from `Params` |
| `UsePathStyle` | `true` so MinIO receives the bucket name in the URL path, not in the host name |

- Use only `--access-key` and `--secret-key` for credentials. Do not use the AWS default credential chain. Always use HTTP. Do not add an HTTPS option.
- In Kubernetes, `stockcenter.S3Flags` reads `MINIO_SERVICE_HOST` and `MINIO_SERVICE_PORT` for the endpoint. The `minio` secret supplies `user` and `pass` through the credential flags' `EnvVars`.
- Use `s3.New(opts)` to build the client before the raw callback. In the `s3` arm of `MO.Match`, take the bucket name from the URL host. Remove the `/` at the start of the URL path to get the object key or prefix.
- Wrap only `s3blob.OpenBucket(ctx, client, bucketName, nil)` in `IOE.TryCatchError`. Use `State.Ctx`. Add context with `IOE.MapLeft`. Never include credentials in an error or log.
- Use the shared input resolver and output writer from Task 5.

- [ ] **Step 4: Run verification gates, integration tests, and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go test -tags integration ./internal/goldenbraid/...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid -g '!**/*_test.go'
```

Docker must be available for `go test -tags integration ./internal/goldenbraid/...`. Expected: All unit tests, integration tests, AST gates, linter, build, and `gopls` pass. Accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Expect no other matches.

```bash
git add go.mod go.sum internal/goldenbraid/gbexport
git commit -m "feat(goldenbraid): support S3 assembly input and artifact output with s3blob"
```

---

### Task 7: CLI Subcommand and Main Integration

**Files:**
- Create: `internal/goldenbraid/gbcli/flags.go`
- Create: `internal/goldenbraid/gbcli/export.go`
- Create: `internal/goldenbraid/gbcli/export_test.go`
- Modify: `cmd/goldenbraid/main.go`

**Interfaces:**
- Consumes: `github.com/urfave/cli/v2`, `internal/cli/stockcenter` (`S3Flags`, `LoggingFlags`), `internal/logger` (`GetCliSlogHandler`), `internal/goldenbraid/gbexport` (`Pipeline`, `State`, `Params`), `ioeutils.ToEither`
- Produces:
  - `func ExportFlags() []cli.Flag`
  - `func RunExport(c *cli.Context) error`
  - `export` subcommand registered in `cmd/goldenbraid/main.go`

- [ ] **Step 1: Write CLI unit tests in `internal/goldenbraid/gbcli/export_test.go`**

1. `TestExportFlagsMatchesExpectedFlagSet` — Use `require.ElementsMatch` for `input`, `input-type`, `output`, `show-primers`, `s3-region`, `s3-server`, `s3-server-port`, `access-key`, `secret-key`, `log-level`, and `log-format`. Make sure that both `input` and `output` have `Required: true`. Compare the flag defaults with Step 3.
2. `TestRunExportLocalFolderWritesArtifacts` — Build a `cli.Context` with `--input` (a folder with one assembly) and `--output` (a second directory). Call `gbcli.RunExport(c)`. Assert success and make sure that all three format paths exist.

- [ ] **Step 2: Run the unit tests to confirm failure**

Run: `go test ./internal/goldenbraid/gbcli/...`
Expected: FAIL because `internal/goldenbraid/gbcli` does not exist.

- [ ] **Step 3: Implement `flags.go`, `export.go`, and update `cmd/goldenbraid/main.go`**

In `internal/goldenbraid/gbcli/flags.go`, do not import `fp-go` or `fp-go-loom`. The `RequireNoPopulatedStruct` AST gate scans only files that import `fp-go` or `fp-go-loom`, so `cli.StringFlag` and `cli.BoolFlag` literals stay valid in `flags.go`. Define `exportFlags() []cli.Flag` with these flags:
- `--input` (`Required: true`): local path or `s3://bucket/key` for an assembly file, zip archive, or directory/prefix.
- `--input-type` (`Value: "auto"`): `auto`, `dna`, `zip`, or `dir`. Use an explicit type to override extension detection.
- `--output` (`Required: true`): local directory or `s3://bucket/prefix` for the three format folders.
- `--s3-region` (`Value: "us-east-1"`): AWS region for the S3 client.
- `--show-primers` (`cli.BoolFlag`): draw `primer_bind` features on the SVG map. This option does not change the GenBank output.

In `internal/goldenbraid/gbcli/export.go`:
- Define `ExportFlags()` as `A.ArrayConcatAll(exportFlags(), stockcenter.S3Flags(), stockcenter.LoggingFlags())`.
- Build the logger with `F.Pipe2(c, logger.GetCliSlogHandler, slog.New)` (`internal/logger`). Do not add a custom logger package.
- Implement `toExportState(c *cli.Context) gbexport.State` with generated lenses (`Ctx` from `c.Context`, `Logger`, and `Params` from flags). Set `S3Endpoint` to `http://<s3-server>:<s3-server-port>`. Set `S3Region`, `AccessKey`, and `SecretKey` from `s3-region`, `access-key`, and `secret-key`. Keep the server, port, and credential flags in `stockcenter.S3Flags()`.
- Use `c` as the seed in `RunExport` to satisfy `RequirePointFreeSeed`:

```go
func RunExport(c *cli.Context) error {
	return F.Pipe4(
		c,
		toExportState,
		gbexport.Pipeline,
		ioeutils.ToEither[error, gbexport.State],
		E.ToError[gbexport.State],
	)
}
```

In `cmd/goldenbraid/main.go`:
- Register the `export` subcommand with `Name: "export"`, `Usage: "Export GoldenBraid assemblies to FASTA, GenBank, and SVG"`, `Action: gbcli.RunExport`, `Flags: gbcli.ExportFlags()`, and no `Before` hook.
- Convert every `slices.Concat` call in `cmd/goldenbraid/main.go` to `A.ArrayConcatAll` (`A "github.com/IBM/fp-go/v2/array"`) and remove the `"slices"` import.

- [ ] **Step 4: Run all verification gates and commit**

```bash
go generate ./internal/goldenbraid/...
golangci-lint fmt
golangci-lint run ./...
go test ./...
go build ./...
git ls-files --cached --others --exclude-standard '*.go' | rg -v '^vendor/' | while read -r f; do
  head -5 "$f" | rg -q 'Code generated .* DO NOT EDIT' && continue
  out=$(gopls check -severity=hint "$f") || exit 1
  [ -z "$out" ] || { echo "$out"; exit 1; }
done
rg '[A-Z][A-Za-z0-9_.]*\([a-z][A-Za-z0-9_.]*\(.*\)' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg '\)\([a-z]' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg -P -U '(?<![\w\]])[A-Za-z_][A-Za-z0-9_.]*\{\s*\n?\s*[A-Za-z_][A-Za-z0-9_]*\s*:' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go' -g '!**/gen_lens.go'
rg 'slices\.Concat' internal/goldenbraid cmd/goldenbraid -g '!**/*_test.go'
```

Expected: All tests, AST gates, linter, build, and `gopls` pass with zero diagnostics. Accept `A.ArrayConcatAll(...Flags())` in the nested-call `rg` check. Accept CLI flag and command struct literals in `internal/goldenbraid/gbcli/flags.go` and `cmd/goldenbraid/main.go`. Expect no other matches.

```bash
git add internal/goldenbraid/gbcli cmd/goldenbraid/main.go
git commit -m "feat(goldenbraid): register the assembly export command with URL-based input and output"
```

---

### Task 8: Kubernetes Job and Dagu Workflow

**Files:**
- Modify: `just-modules/goldenbraid.just`
- Create: `dagu-goldenbraid-export.yaml`
- Create: `schemas/goldenbraid-export-params.json`
- Reference: `dagu-goldenbraid-inventory.yaml`, `dagu-goldenbraid-plasmid.yaml`, `schemas/goldenbraid-plasmid-params.json`, and the `dagu-params` recipe in `justfile`.

**Interfaces:** The recipe returns the Job name on stdout. The workflow passes this name to the wait, log, and failure steps. Local paths refer to the container file system, not the submitter's machine. Use S3 paths for artifacts that must survive Job deletion.

- [ ] **Step 1: Add the Job recipe**

Copy `run-goldenbraid-inventory` in `just-modules/goldenbraid.just`. Use this signature:

```just
run-goldenbraid-export tag k8s_config input output k8s_namespace="dev" debug="false":
```

| Job attribute | Normal | Debug (`debug="true"`) |
|---|---|---|
| `ttlSecondsAfterFinished` | `120` | `300` |
| `generateName` | `goldenbraid-export-` | `goldenbraid-export-debug-` |
| Container name | `goldenbraid-export` | `goldenbraid-export-debug` |
| Extra arguments | None | `--log-level debug --log-format text` |

Keep `restartPolicy: Never`, `namespace: {{ k8s_namespace }}`, and `image: {{ ghcr_image }}:{{ tag }}`. Keep `envFrom.secretRef.name: minio`. Replace the container arguments with this list:

```yaml
args:
  - export
  - --input
  - "{{ input }}"
  - --output
  - "{{ output }}"
  ${debug_args}
```

Keep the debug arguments at the same indentation as the other arguments in the recipe heredoc. Submit with `kubectl create -f - --kubeconfig {{ k8s_config }} -o jsonpath='{.metadata.name}'`.

- [ ] **Step 2: Define the workflow parameters**

Copy `schemas/goldenbraid-plasmid-params.json` to `schemas/goldenbraid-export-params.json`. Remove `EMAIL`. Keep the draft, object type, and string property types. Set `required` to `["TAG", "INPUT", "OUTPUT"]`.

| Property | Description or constraint | Default |
|---|---|---|
| `TAG` | The Docker image tag to use | No default |
| `INPUT` | Input .dna, .zip, or folder. local path or s3://bucket/key | No default |
| `OUTPUT` | Output root. local path or s3://bucket/prefix | No default |
| `ENV` | Enum: `dev`, `staging`, `production` | `dev` |
| `NAMESPACE` | Kubernetes namespace | `dev` |
| `KUBE_CFG` | Keep the source description. Supply a kubeconfig path outside `dev`. | `""` |

- [ ] **Step 3: Connect the Job to Dagu**

Copy `dagu-goldenbraid-inventory.yaml` to `dagu-goldenbraid-export.yaml`. Set `name: goldenbraid-export` and describe the export Job. Keep `type: graph`. Set `params.schema` to `./schemas/goldenbraid-export-params.json`. Keep the values `ENV: dev`, `NAMESPACE: dev`, and `KUBE_CFG: ""`.

Keep `k8s-generate` unchanged. In `submit`, keep `depends: [k8s-generate]` and `output: JOB_NAME`. Use this command body:

```bash
set -euo pipefail
JOB_OUTPUT=$(just goldenbraid::run-goldenbraid-export "${TAG}" "${KPATH}" "${INPUT}" "${OUTPUT}" "${NAMESPACE}")
echo "$JOB_OUTPUT"
```

Copy `wait`, `get-logs`, and `handler_on.failure` without changes. The failure handler calls `just job-debug` only if `JOB_NAME` exists.

- [ ] **Step 4: Examine the files and commit**

```bash
just dagu-params dagu-goldenbraid-export.yaml
just --list goldenbraid
yq . dagu-goldenbraid-export.yaml
jq . schemas/goldenbraid-export-params.json
```

Expected: The parameter list shows `TAG`, `INPUT`, `OUTPUT`, `ENV`, `NAMESPACE`, and `KUBE_CFG`. The recipe list includes `run-goldenbraid-export`. Both parsers exit with status 0. These checks do not submit a Job.

```bash
git add just-modules/goldenbraid.just dagu-goldenbraid-export.yaml schemas/goldenbraid-export-params.json
git commit -m "feat(goldenbraid): add export job recipe and Dagu workflow"
```
