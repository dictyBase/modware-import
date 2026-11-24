# Phenotype workbook normalization: implementation plan

> **For agentic workers:** Use superpowers:executing-plans or superpowers:subagent-driven-development to do this plan one task at a time. Each step has a checkbox (`- [ ]`).

## Table of Contents

- [Global constraints](#global-constraints)
- [Review focus](#review-focus)
- [Command and input contract](#command-and-input-contract)
- [Worksheet contract](#worksheet-contract)
- [Canonical outputs](#canonical-outputs)
- [Failure and publication contract](#failure-and-publication-contract)
- [File map](#file-map)
- [Implementation tasks](#implementation-tasks)
  - [Task 0: Add the dependencies and linter](#task-0-add-the-dependencies-and-linter)
  - [Task 1: Types, generated lenses, schema, and the row classifier](#task-1-types-generated-lenses-schema-and-the-row-classifier)
  - [Task 2: File selection, workbook reads, and batch accumulation](#task-2-file-selection-workbook-reads-and-batch-accumulation)
  - [Task 3: Atomic CSV publication](#task-3-atomic-csv-publication)
  - [Task 4: Command, process status, and functional gates](#task-4-command-process-status-and-functional-gates)
  - [Task 5: Corpus check and README](#task-5-corpus-check-and-readme)
- [Reference implementation](#reference-implementation)
  - [types.go](#typesgo)
  - [optics.go](#opticsgo)
  - [schema.go](#schemago)
  - [row.go](#rowgo)
  - [source.go](#sourcego)
  - [workbook.go](#workbookgo)
  - [batch.go](#batchgo)
  - [output.go](#outputgo)
  - [rename_linux.go, rename_darwin.go, and rename_other.go](#rename_linuxgo-rename_darwingo-and-rename_othergo)
  - [cli/arguments.go](#cliargumentsgo)
  - [cli/normalize.go](#clinormalizego)
  - [cmd/phenoannotation/main.go](#cmdphenoannotationmaingo)
  - [normalize_test.go](#normalize_testgo)
  - [cli/normalize_test.go](#clinormalize_testgo)
  - [cli/gates_test.go](#cligates_testgo)
- [fp-go v2 rules for this package](#fp-go-v2-rules-for-this-package)
- [Verification record](#verification-record)

**Goal:** Make reviewable CSV files for strain creation and phenotype import. Do not change a stock or an annotation.

**Architecture:** One `phenoannotation` binary will have three subcommands. This plan covers only `normalize`. The `create-strains` and `load-phenotypes` subcommands will read its files later. The command uses no service credentials, no ontology download, and no cluster access.

**Stack:** Go 1.26.10, `urfave/cli/v3` v3.14.0, `excelize/v2` v2.11.1, `encoding/csv`, `fp-go/v2` v2.13.2, `fp-go-loom` v0.10.0, and `golang.org/x/sys` v0.49.0.

**Basis:** The decisions of this session and the create workbooks in `kosher-files/`. The code below compiled and passed all gates in a scratch copy on 2026-11-23. The tests are the contract.

## Global constraints

- Keep the current `baserow`, `import`, and other binaries on their present CLI versions.
- Use only `Option`, `Either`, `IO`, and `IOEither` for effects. Use the alias `IOE` for `ioeither`.
- Use term IDs only. Ignore all term-name columns.
- Do not create a strain characteristic.
- Each row without a DBS ID requests one new strain. Do not use the descriptor as a key.
- Do not read the delete files or the update files.
- Do not call a network service.

## Review focus

These are the primary failures that a person can see. Each line identifies the test that covers it.

1. A cell has a nonbreaking space (U+00A0) at an edge. The value must match. `TestProcess` covers it.
2. A workbook has a blank row between two records. The row loop in `readSheet` counts each physical row, so row numbers stay physical and blank rows are counted without storing raw rows. `TestWorkbookFailures` and `TestProcess` cover it.
3. A note has a tab, a quote, or a line break. Standard CSV quoting from `encoding/csv` keeps it. `TestProcess` compares the decoded field.
4. A second run uses the same `--out`. The command must not replace the first batch. `TestPublish` and `TestProcess` cover it.
5. The operator stops the command during the write. No partial directory must stay. `TestCancelDuringPublish` and `TestCancelledRun` cover it.
6. A resource close or cleanup fails after a read or write error. The release error must join the error of the work step and must not replace it. `TestReleaseJoinsErrors` covers it.
7. A file exists at the target path of `WriteFile`. `WriteFile` must fail with `os.ErrExist` and must not overwrite it. `TestWriteFileExclusive` covers it.

## Command and input contract

```sh
phenoannotation normalize --folder ./kosher-files --out ./batch-01
phenoannotation normalize --input ./kosher-files/Phenotype_annotation_APR_06_2023.xlsx --out ./sample
```

- Give only one of `--input` and `--folder`. Give a new destination directory with `--out`.
- In folder mode, the command selects regular files with the pattern `Phenotype_annotation_*.xlsx` directly in the folder. It sorts them by name.
- The command ignores update files, delete files, lock files, and other files. It rejects a file from `--input` that does not match the create pattern. It also rejects a symbolic link. An empty selection is an error.
- The command reads `.xlsx` files, not binary `.xls` files. It uses the only worksheet by default. A workbook with more than one worksheet must have `--sheet NAME`.
- Each output row keeps the file name and the physical row number in Excel. Row 1 is the header. A blank row does not change the number of the next row.
- Keep each source workbook read only. Use one input folder with different file names for each batch.
- The command does not read a date from the file name. It does not correct the `2043` file name.

## Worksheet contract

Row 1 sets the layout. The cell count of a data row does not set it. Excel does not always include empty cells at the end of a row. The two observed headers follow:

```text
Strain ID | ddpheno ID | Phenotype term | Notes | Assay ID | Assay term | Envirionment ID | Environment term | Reference(s) | Assigned By
Strain ID | Strain descriptor | ddpheno ID | Phenotype term | Notes | Assay ID | Assay term | Envirionment ID | Environment term | Reference(s) | Deleted | Assigned By
```

The command accepts `Environment ID` as an alias for `Envirionment ID`. It trims each header cell. It rejects an unknown header or a duplicate header before it writes a file. It adds empty cells to a short data row. A nonempty cell after the last header column is an error. The command counts a fully blank row as blank and writes no output for it.

All 77 create workbooks use the worksheet name `Sheet 1`. Do not use the Baserow default `Strain_Annotations`. Do not use the old assay predicate in `internal/datasource/xls/phenotype/annotation.go`. It uses an underscore, but the workbooks use a colon.

## Canonical outputs

Each CSV file is UTF-8, uses commas, and uses standard CSV quoting from `encoding/csv`. Each file has a header, also when it has no records. The two accepted files have these fields in this order:

```text
source_file,source_row,strain_id,descriptor,pheno_id,assay_id,env_id,note,pmid,assigned_by
```

| Field | Rule |
|---|---|
| `source_file`, `source_row` | The file name and the 1-based row number in Excel. Together they identify one source row. |
| `strain_id` | `DBS` and seven digits. A malformed ID is an error. It is not a request for a new strain. |
| `descriptor` | Kept as text. It is necessary only when the strain ID is empty. It is not a key. |
| `pheno_id` | `DDPHENO:` and seven digits. The command does not use a term name. |
| `assay_id` | Empty, or `DDASSAY:` and seven digits. The command clears a `DDSTRAINCHAR:` ID and writes a warning. |
| `env_id` | Empty, or `DDENVIR:` and seven digits. The command changes `DDDENVIR:` to `DDENVIR:` and writes a warning. |
| `note` | Kept with its tabs, quotes, and line breaks. An empty note is valid. |
| `pmid` | One positive decimal number. The command removes a `PMID:` prefix. It does not query PubMed. |
| `assigned_by` | One email address with no display name. Case stays the same. |

The command trims the edges of each field with `strings.TrimSpace`, which also removes U+00A0. It does not change internal spaces. It does not repair an unknown prefix. A term name without its ID gives a warning. A nonempty `Deleted` cell in a create workbook gives an error.

- **`phenotype.csv`** has each valid row with a DBS ID. Repeated observations stay.
- **`new-strain.csv`** has each valid row without a DBS ID. Two rows with the same descriptor are two strains.
- **`rejects.csv`** has the columns `source_file`, `source_row`, `severity`, `code`, `field`, and `raw_value`. The severity is `error` or `warning`. A warning does not reject the row. The diagnostics of one row are in the column order.

**Downstream boundary:** `create-strains` will map `(source_file, source_row)` to a new DBS ID. It will use the stock center depositor `dictystocks@northwestern.edu`. `load-phenotypes` will send term IDs and skip groups that exist. This plan does not include those steps. The `normalize` command does not remove duplicates, does not set ranks, and does not get IDs from a service.

## Failure and publication contract

`readSheet` loops over each worksheet row by row and calls the row pipeline `addRow` for each data row. The batch grows row by row, and `ReadFiles` joins the per-file batches in order. Only accepted records and diagnostics stay in memory, not raw rows. The command checks all rows before `Publish` writes a file. `Publish` writes the three files into a new temporary directory next to the destination. Then it renames that directory to the destination with one no-replace rename.

The rename is the only commit point. The command does not replace a destination that exists. After a failure or a cancellation, the command removes the temporary directory.

`WriteFile` is an `IOEither` effect and matches `FileWriter = func(out Output) IOE.IOEither[error, string]`. Each `IOE.Bracket` resource scope (workbooks, row iterators, output files, and the temporary directory) joins a release error with the error of the work step. A release error does not replace the primary error.

| Exit | Output directory | Meaning |
|---|---|---|
| `0` | Complete | All nonblank rows are valid. Warnings can occur. |
| `2` | Complete | One or more rows are not valid. Examine `rejects.csv` before you use the batch. |
| `1` | None | The options, a workbook, a header, a read, or a write failed. |

The command prints one summary line after the commit. That line shows the counts of files, blank rows, rows for known strains, rows for new strains, rejected rows, and warnings. The command does not change a source file.

## File map

| File | Responsibility |
|---|---|
| `internal/phenoannotation/normalize/types.go` | Typed enums, typed codes and field names, named records, `FileWriter`, and the `go:generate` directive. |
| `internal/phenoannotation/normalize/gen_lens.go` | Generated lenses from `go generate ./internal/phenoannotation/normalize`. Commit this file. |
| `internal/phenoannotation/normalize/optics.go` | Generated lens instances and composed lenses for nested fields. |
| `internal/phenoannotation/normalize/schema.go` | Headers, layouts, `toCells`, `recordHeader`, `diagnosticHeader`, `recordFields`, `diagnosticFields`, `recordRows`, and `diagnosticRows`. |
| `internal/phenoannotation/normalize/row.go` | Predicates, the `rules` list, `newRule`, `diagnose`, `checkRule`, `kindOf`, `toRecord`, and `classify`. |
| `internal/phenoannotation/normalize/source.go` | Option validation (`Valid`), `workbookPath`, and file selection (`SelectFiles`). |
| `internal/phenoannotation/normalize/workbook.go` | Row alignment (`newLine`, `toSourceRow`), `pickSheet`, `readHeader`, the `readSheet` row loop, `IOE.Bracket` resources with joined errors, and `ReadFiles`. |
| `internal/phenoannotation/normalize/batch.go` | The `addRow` row pipeline, `outcomeBatch`, `joinBatch`, and `PrintSummary`. |
| `internal/phenoannotation/normalize/output.go` | `batchFiles`, `writeCSV`, `WriteFile` with `IOE.Bracket`, `Publish` with `IOE.Bracket`, and `renameInto`. |
| `internal/phenoannotation/normalize/rename_linux.go` | Linux atomic no-replace rename (`Renameat2`). |
| `internal/phenoannotation/normalize/rename_darwin.go` | Darwin atomic no-replace rename (`RenamexNp`). |
| `internal/phenoannotation/normalize/rename_other.go` | Error stub for platforms without atomic no-replace rename. |
| `internal/phenoannotation/normalize/normalize_test.go` | Unit, adapter, resource, and process tests. |
| `internal/phenoannotation/cli/arguments.go` | CLI boundary record and the `go:generate` directive. |
| `internal/phenoannotation/cli/gen_lens.go` | Generated lenses from `go generate ./internal/phenoannotation/cli`. Commit this file. |
| `internal/phenoannotation/cli/normalize.go` | `MakeLensRef` CLI builders, `BuildFrom` state setup, and the `runNormalize` pipeline. |
| `internal/phenoannotation/cli/normalize_test.go` | CLI unit tests for cancellation, subcommand, and review status. |
| `internal/phenoannotation/cli/gates_test.go` | `pipelinecheck` and `loomcheck` functional gates over all three packages. |
| `cmd/phenoannotation/main.go` | Signal setup, stderr error output, and process exit codes. |

## Implementation tasks

For each task, write the test first. Run the test. Make sure that it fails. Then add the code. Run the test again. Copy code from the [reference implementation](#reference-implementation).

### Task 0: Add the dependencies and linter

- [ ] Run this command in the module root:

```sh
go get github.com/urfave/cli/v3@v3.14.0 github.com/dictyBase/fp-go-loom@v0.10.0
```

- [ ] Run `golangci-lint version`. Make sure that the version is v2.14.0. The local tool must be v2.13.2 or a higher version, because v2.12.2 has an `unparam` defect. The CI pin in `.github/workflows/lint-test-build.yml` is v2.13.2.
- [ ] After Task 3 adds the rename adapters, run `go mod tidy`. Make sure that `golang.org/x/sys` moves to the direct requirements. Make sure that no other module version changes.

### Task 1: Types, generated lenses, schema, and the row classifier

**Files:** `types.go`, `gen_lens.go` (generated), `optics.go`, `schema.go`, `row.go`, and `normalize_test.go` (with `TestCSVColumns` and `TestHeader`).

**Produces:** Typed enums and strings (`Layout`, `RowKind`, `Severity`, `Code`, `FieldName`), named records (`Origin`, `Cells`, `SourceRow`, `Record`, `Diagnostic`, `RowOutcome`, `Batch`, `Options`, `State`, `rule`, `line`, `sheetRequest`), `Output`, `FileWriter`, generated and composed lenses, `decodeHeader`, `toCells`, `cellAt`, `layoutColumns`, `layoutWidth`, `recordHeader`, `diagnosticHeader`, `recordFields`, `diagnosticFields`, `recordRows`, `diagnosticRows`, `nonblank`, `firstNonblank`, `isError`, `isWarning`, `kindEq`, `rules`, `newRule`, `diagnose`, `checkRule`, `kindOf`, `toRecord`, and `classify`.

- [ ] Add `check`, `TestCSVColumns`, and `TestHeader` to `normalize_test.go`. Run `go test ./internal/phenoannotation/normalize -run '^Test(CSVColumns|Header)$'`. The build fails because the package does not exist.
- [ ] Add `types.go`.
- [ ] Run `go generate ./internal/phenoannotation/normalize`. Make sure that `gen_lens.go` appears in `internal/phenoannotation/normalize/`.
- [ ] Add `optics.go`, `schema.go`, and `row.go`.
- [ ] Run `go test ./internal/phenoannotation/normalize -run '^Test(CSVColumns|Header)$'`. Make sure that the two tests pass.
- [ ] Commit `types.go`, `gen_lens.go`, `optics.go`, `schema.go`, `row.go`, and `normalize_test.go` with the message `feat(phenoannotation): classify phenotype workbook rows`.

### Task 2: File selection, workbook reads, and batch accumulation

**Files:** `source.go`, `workbook.go`, `batch.go`, and the tests `TestAlignShortRow`, `TestValidTerms`, `TestRowFaults`, `TestAccounting`, `TestSelectFiles`, and `TestWorkbookFailures` in `normalize_test.go`.

**Consumes:** `types.go`, `optics.go`, `decodeHeader`, `toCells`, `cellAt`, `layoutColumns`, `layoutWidth`, `nonblank`, `firstNonblank`, `isWarning`, `kindEq`, and `classify`.

**Produces:** `Valid`, `liveState`, `workbookPath`, `SelectFiles`, `newLine`, `toSourceRow`, `pickSheet`, `readHeader`, `readSheet`, `joinUseError`, `closeJoined`, `ReadFiles`, `outcomeBatch`, `joinBatch`, `addRow`, and `PrintSummary`.

- [ ] Add `force`, `dataRow`, `fixture`, `codes`, and `baseState` (without `WriteFile` at this step) to `normalize_test.go`. Add `TestAlignShortRow`, `TestValidTerms`, `TestRowFaults`, `TestAccounting`, `TestSelectFiles`, and `TestWorkbookFailures`.
- [ ] Run `go test ./internal/phenoannotation/normalize -run '^Test(CSVColumns|Header|AlignShortRow|ValidTerms|RowFaults|Accounting|SelectFiles|WorkbookFailures)$'`. Make sure that the build fails.
- [ ] Add `source.go`, `workbook.go`, and `batch.go`. `readSheet` loops over the rows of each worksheet and calls `addRow` for each data row so the batch grows row by row. `ReadFiles` joins the per-file batches. `IOE.Bracket` closes each workbook and row iterator with `closeJoined` on all paths.
- [ ] Run the same test command. Make sure that all eight tests pass.
- [ ] Commit with the message `feat(phenoannotation): read create workbooks into batches`.

### Task 3: Atomic CSV publication

**Files:** `output.go`, `rename_linux.go`, `rename_darwin.go`, `rename_other.go`, and the tests `TestPublish`, `TestCancelDuringPublish`, `TestReleaseJoinsErrors`, and `TestWriteFileExclusive` in `normalize_test.go`.

**Consumes:** `types.go`, `optics.go`, `recordRows`, `diagnosticRows`, `liveState`, `joinUseError`, and `closeJoined`.

**Produces:** `batchFiles`, `writeCSV`, `WriteFile`, `removeJoined`, `renameInto`, and `Publish`.

- [ ] Set `s.Writer = WriteFile` in `baseState`. Add `outputNames`, `readCSV`, `emptyDir`, `failWriter`, `badCloser`, `failUse`, `failDir`, `TestPublish`, `TestCancelDuringPublish`, `TestReleaseJoinsErrors`, and `TestWriteFileExclusive` to `normalize_test.go`.
- [ ] Run `go test ./internal/phenoannotation/normalize -run '^Test(Publish|CancelDuringPublish|ReleaseJoinsErrors|WriteFileExclusive)$'`. Make sure that the build fails.
- [ ] Add `output.go`, `rename_linux.go`, `rename_darwin.go`, and `rename_other.go`.
- [ ] Run `go mod tidy`. Make sure that `golang.org/x/sys` moves to the direct requirements.
- [ ] Run `go test ./internal/phenoannotation/normalize -count=1`. Make sure that all 12 tests pass.
- [ ] Run `CGO_ENABLED=0 GOOS=linux go vet ./internal/phenoannotation/...`. Run `CGO_ENABLED=0 GOOS=windows go build ./internal/phenoannotation/...`. Make sure that the two commands pass.
- [ ] Commit with the message `feat(phenoannotation): publish normalized CSV batches`.

### Task 4: Command, process status, and functional gates

**Files:** `cli/arguments.go`, `cli/gen_lens.go` (generated), `cli/normalize.go`, `cmd/phenoannotation/main.go`, `cli/normalize_test.go` (`TestCancelledRun`, `TestNeedsSubcommand`, `TestReviewStatus`), `cli/gates_test.go` (`TestFunctionalGates`), and `TestProcess` in `normalize_test.go`.

**Consumes:** `State`, `Options`, `MakeStateLenses`, `MakeOptionsLenses`, `MakeBatchLenses`, `Valid`, `SelectFiles`, `ReadFiles`, `Publish`, `PrintSummary`, and `WriteFile`.

**Produces:** `Command(io.Writer) *cli.Command`, `ErrReview`, and the `phenoannotation` binary.

- [ ] Add `cli/normalize_test.go`, `cli/gates_test.go`, and `TestProcess` in `normalize_test.go`. Run `go test ./internal/phenoannotation/... -count=1`. Make sure that the build fails.
- [ ] Add `cli/arguments.go`.
- [ ] Run `go generate ./internal/phenoannotation/cli`. Make sure that `gen_lens.go` appears in `internal/phenoannotation/cli/`.
- [ ] Add `cli/normalize.go` and `cmd/phenoannotation/main.go`.
- [ ] Run `go test ./internal/phenoannotation/... -count=1`. Make sure that all 17 tests pass across the two packages (13 in `normalize` and 4 in `cli`).
- [ ] Run the gates in the [verification record](#verification-record), plus repo-wide `golangci-lint run ./...` and `go test ./...`.
- [ ] Commit `cli/arguments.go`, `cli/gen_lens.go`, `cli/normalize.go`, `cli/normalize_test.go`, `cli/gates_test.go`, `cmd/phenoannotation/main.go`, and `normalize_test.go` with the message `feat(phenoannotation): add normalize command`.

### Task 5: Corpus check and README

- [ ] Run the local corpus check. CI must not use these files.

```sh
go build -o /tmp/phenoannotation ./cmd/phenoannotation
/tmp/phenoannotation normalize --folder ../../kosher-files --out /tmp/pheno-review-01
python3 -c 'import csv,sys,collections; print(collections.Counter(r[3] for r in list(csv.reader(open(sys.argv[1])))[1:]))' /tmp/pheno-review-01/rejects.csv
```

- [ ] Make sure that the result agrees with the [verification record](#verification-record). If a count is different, examine `rejects.csv` before you continue.
- [ ] Add a `phenoannotation` section to `README.md`. Show the two commands, the three output files, the exit codes, and the two later subcommands. Tell the operator to examine `rejects.csv` before the next step.
- [ ] Commit with the message `docs(phenoannotation): document normalize`.

## Reference implementation

This code compiled and passed all gates in a scratch copy on 2026-11-23. It is a reference, not a change to the repository.

`go generate ./internal/phenoannotation/...` generates `internal/phenoannotation/normalize/gen_lens.go` and `internal/phenoannotation/cli/gen_lens.go` from the `// fp-go:Lens` markers (`fp-go` v2.13.2 lens generator). Commit the two generated files. Do not write them by hand. This plan does not include them.

### types.go

`internal/phenoannotation/normalize/types.go`. This file defines the typed domain enums, typed strings, named records, `FileWriter`, and the `go:generate` line for `gen_lens.go`.

```go
// Package normalize converts create workbooks for phenotype annotation into
// reviewable CSV files. It does not call a stock or an annotation service.
package normalize

import (
	"context"
	"errors"
	"io"

	IOE "github.com/IBM/fp-go/v2/ioeither"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
)

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

// Layout identifies one of the two accepted worksheet headers.
type Layout int

// The accepted layouts. LayoutUnknown marks a header that is not accepted.
const (
	LayoutUnknown Layout = iota
	LayoutShort
	LayoutExtended
)

// RowKind is the classification of one physical worksheet row.
type RowKind int

// The row kinds. A rejected row has at least one error diagnostic.
const (
	KindBlank RowKind = iota + 1
	KindRejected
	KindExisting
	KindNew
)

// Severity is the severity of one diagnostic.
type Severity string

// The diagnostic severities. A warning does not reject a row.
const (
	SevError   Severity = "error"
	SevWarning Severity = "warning"
)

// Code identifies the rule that made a diagnostic.
type Code string

// The diagnostic codes.
const (
	CodeInvalidStrainID   Code = "invalid_strain_id"
	CodeMissingDescriptor Code = "missing_descriptor"
	CodeInvalidPhenoID    Code = "invalid_pheno_id"
	CodeStrainChar        Code = "unsupported_strain_characteristic"
	CodeInvalidAssayID    Code = "invalid_assay_id"
	CodeLabelWithoutID    Code = "label_without_id"
	CodeEnvPrefix         Code = "corrected_env_prefix"
	CodeInvalidEnvID      Code = "invalid_env_id"
	CodeInvalidPMID       Code = "invalid_pmid"
	CodeDeletedInCreate   Code = "deleted_in_create"
	CodeInvalidCurator    Code = "invalid_assigned_by"
	CodeExtraCell         Code = "extra_cell"
)

// FieldName names the output field or source column of a diagnostic.
type FieldName string

// The field names that diagnostics use.
const (
	FieldStrainID   FieldName = "strain_id"
	FieldDescriptor FieldName = "descriptor"
	FieldPhenoID    FieldName = "pheno_id"
	FieldAssayID    FieldName = "assay_id"
	FieldAssayTerm  FieldName = "assay_term"
	FieldEnvID      FieldName = "env_id"
	FieldEnvTerm    FieldName = "env_term"
	FieldPMID       FieldName = "pmid"
	FieldDeleted    FieldName = "deleted"
	FieldAssignedBy FieldName = "assigned_by"
	FieldOverflow   FieldName = "overflow"
)

const (
	shortWidth    = 10
	extendedWidth = 12
	headerRow     = 1
	firstDataRow  = 2
)

var (
	errBadLayout = errors.New("unknown or duplicate worksheet header")
	errNoHeader  = errors.New("missing worksheet header")
	errNoSheet   = errors.New("workbook has no worksheet")
	errSpec      = errors.New(
		"specify exactly one of --input and --folder and a new --out directory",
	)
	errDisplayName    = errors.New("assigned_by has a display name")
	errChangedAddress = errors.New("assigned_by is not a bare address")
)

// Origin identifies one physical source row: the workbook base name and
// the 1-based Excel row number.
// fp-go:Lens
type Origin struct {
	File string
	Row  int
}

// Cells holds the twelve worksheet columns in the extended order. A short
// row has empty Descriptor and Deleted cells.
// fp-go:Lens
type Cells struct {
	Strain     string
	Descriptor string
	Pheno      string
	PhenoTerm  string
	Note       string
	Assay      string
	AssayTerm  string
	Env        string
	EnvTerm    string
	Reference  string
	Deleted    string
	Curator    string
}

// SourceRow is one physical worksheet row. Raw keeps each cell as read.
// Clean has the same cells with the edges trimmed. Extra is the first
// nonblank cell after the last header column.
// fp-go:Lens
type SourceRow struct {
	Origin Origin
	Raw    Cells
	Clean  Cells
	Extra  string
}

// Record is one accepted row. It follows the accepted CSV columns.
// fp-go:Lens
type Record struct {
	Origin     Origin
	StrainID   string
	Descriptor string
	PhenoID    string
	AssayID    string
	EnvID      string
	Note       string
	PMID       string
	AssignedBy string
}

// Diagnostic is one error or warning for one source row.
// fp-go:Lens
type Diagnostic struct {
	Origin   Origin
	Severity Severity
	Code     Code
	Field    FieldName
	Raw      string
}

// RowOutcome is the classification of one SourceRow.
// fp-go:Lens
type RowOutcome struct {
	Kind   RowKind
	Record Record
	Issues []Diagnostic
}

// Batch holds every accepted row and every diagnostic of one run.
// fp-go:Lens
type Batch struct {
	Existing []Record
	New      []Record
	Issues   []Diagnostic
	Blank    int
	Rejected int
	Warnings int
}

// Options holds the command options of one run.
// fp-go:Lens
type Options struct {
	Input  string
	Folder string
	Sheet  string
	Out    string
}

// Output pairs the path of one CSV file with its rows, header first.
type Output = PR.Pair[string, [][]string]

// FileWriter returns an effect that creates the file at the path of out,
// writes the rows as CSV, and returns the path. Tests replace it with a
// writer that fails.
type FileWriter = func(out Output) IOE.IOEither[error, string]

// State carries the options, the dependencies, and the results of one run.
// Batch grows row by row while the workbooks are read.
// fp-go:Lens
type State struct {
	Context context.Context
	Options Options
	Writer  FileWriter
	Log     io.Writer
	Files   []string
	Batch   Batch
}

// rule describes one row check. Broken is true when the row breaks the
// rule. Raw reads the raw cell that the diagnostic reports.
// fp-go:Lens
type rule struct {
	Broken   P.Predicate[SourceRow]
	Severity Severity
	Code     Code
	Field    FieldName
	Raw      func(SourceRow) string
}

// line is one data row of a worksheet as read: the workbook base name,
// the header layout, the 1-based row number, and the cells.
// fp-go:Lens
type line struct {
	File   string
	Layout Layout
	Row    int
	Cells  []string
}

// sheetRequest pairs the worksheet names of a workbook with the --sheet
// option.
// fp-go:Lens
type sheetRequest struct {
	Names []string
	Want  string
}
```

### optics.go

`internal/phenoannotation/normalize/optics.go`. This file instantiates the generated lenses and composes lenses for nested fields.

```go
package normalize

import (
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

// Generated lenses for each record, and composed lenses that reach a
// nested field from the record that the pipeline carries.
var (
	originL  = MakeOriginLenses()
	cellsL   = MakeCellsLenses()
	rowL     = MakeSourceRowLenses()
	recordL  = MakeRecordLenses()
	diagL    = MakeDiagnosticLenses()
	outcomeL = MakeRowOutcomeLenses()
	batchL   = MakeBatchLenses()
	optionsL = MakeOptionsLenses()
	stateL   = MakeStateLenses()
	ruleL    = MakeruleLenses()
	lineL    = MakelineLenses()
	requestL = MakesheetRequestLenses()

	// State to options.
	stateInput  = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Input))
	stateFolder = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Folder))
	stateSheet  = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Sheet))
	stateOut    = F.Pipe1(stateL.Options, L.Compose[State](optionsL.Out))

	// Source row to origin.
	recordRow = F.Pipe1(recordL.Origin, L.Compose[Record](originL.Row))
	diagFile  = F.Pipe1(diagL.Origin, L.Compose[Diagnostic](originL.File))
	diagRow   = F.Pipe1(diagL.Origin, L.Compose[Diagnostic](originL.Row))
	recFile   = F.Pipe1(recordL.Origin, L.Compose[Record](originL.File))

	// Source row to trimmed cells. Rules compare these values.
	cleanStrain     = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Strain))
	cleanDescriptor = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Descriptor))
	cleanPheno      = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Pheno))
	cleanNote       = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Note))
	cleanAssay      = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Assay))
	cleanAssayTerm  = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.AssayTerm))
	cleanEnv        = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Env))
	cleanEnvTerm    = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.EnvTerm))
	cleanReference  = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Reference))
	cleanDeleted    = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Deleted))
	cleanCurator    = F.Pipe1(rowL.Clean, L.Compose[SourceRow](cellsL.Curator))

	// Source row to raw cells. Diagnostics report these values.
	rawStrain     = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Strain))
	rawDescriptor = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Descriptor))
	rawPheno      = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Pheno))
	rawAssay      = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Assay))
	rawAssayTerm  = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.AssayTerm))
	rawEnv        = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Env))
	rawEnvTerm    = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.EnvTerm))
	rawReference  = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Reference))
	rawDeleted    = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Deleted))
	rawCurator    = F.Pipe1(rowL.Raw, L.Compose[SourceRow](cellsL.Curator))
)
```

### schema.go

`internal/phenoannotation/normalize/schema.go`. The full header comparison also finds a duplicate header. `recordHeader` and `recordFields` define the columns of `phenotype.csv` and `new-strain.csv`. `diagnosticHeader` and `diagnosticFields` define the columns of `rejects.csv`. `TestCSVColumns` makes sure that each header and field list have the same width.

```go
package normalize

import (
	"strconv"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
)

// Internal column positions. The workbook adapter expands a short row to
// these positions before it builds Cells.
const (
	colStrain = iota
	colDescriptor
	colPheno
	colPhenoTerm
	colNote
	colAssay
	colAssayTerm
	colEnv
	colEnvTerm
	colReference
	colDeleted
	colCurator
)

var (
	header10 = []string{
		"Strain ID", "ddpheno ID", "Phenotype term", "Notes", "Assay ID",
		"Assay term", "Environment ID", "Environment term", "Reference(s)",
		"Assigned By",
	}
	header12 = []string{
		"Strain ID", "Strain descriptor", "ddpheno ID", "Phenotype term",
		"Notes", "Assay ID", "Assay term", "Environment ID",
		"Environment term", "Reference(s)", "Deleted", "Assigned By",
	}
	headerAlias = strings.NewReplacer("Envirionment ID", "Environment ID")
	headerEq    = A.StrictEquals[string]()
	layoutEq    = EQ.FromStrictEquals[Layout]()

	isShortHeader    = F.Pipe1(header10, EQ.Equals(headerEq))
	isExtendedHeader = F.Pipe1(header12, EQ.Equals(headerEq))
	knownLayout      = P.Not(F.Pipe1(LayoutUnknown, EQ.Equals(layoutEq)))

	// chooseLayout compares the full header, so a duplicate header fails.
	chooseLayout = MO.Match(LayoutUnknown, []O.Kleisli[[]string, Layout]{
		MO.Const(isShortHeader, LayoutShort),
		MO.Const(isExtendedHeader, LayoutExtended),
	})

	// layoutColumns maps each internal column to its source column. The
	// value -1 marks a column that the short layout does not have.
	layoutColumns = map[Layout][]int{
		LayoutShort:    {0, -1, 1, 2, 3, 4, 5, 6, 7, 8, -1, 9},
		LayoutExtended: {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
	}
	layoutWidth = map[Layout]int{
		LayoutShort:    shortWidth,
		LayoutExtended: extendedWidth,
	}

	// toCells builds Cells from a row in internal column order.
	toCells = endoutils.BuildFrom([]func([]string) En.Endomorphism[Cells]{
		F.Flow2(F.Bind2nd(cellAt, colStrain), cellsL.Strain.Set),
		F.Flow2(F.Bind2nd(cellAt, colDescriptor), cellsL.Descriptor.Set),
		F.Flow2(F.Bind2nd(cellAt, colPheno), cellsL.Pheno.Set),
		F.Flow2(F.Bind2nd(cellAt, colPhenoTerm), cellsL.PhenoTerm.Set),
		F.Flow2(F.Bind2nd(cellAt, colNote), cellsL.Note.Set),
		F.Flow2(F.Bind2nd(cellAt, colAssay), cellsL.Assay.Set),
		F.Flow2(F.Bind2nd(cellAt, colAssayTerm), cellsL.AssayTerm.Set),
		F.Flow2(F.Bind2nd(cellAt, colEnv), cellsL.Env.Set),
		F.Flow2(F.Bind2nd(cellAt, colEnvTerm), cellsL.EnvTerm.Set),
		F.Flow2(F.Bind2nd(cellAt, colReference), cellsL.Reference.Set),
		F.Flow2(F.Bind2nd(cellAt, colDeleted), cellsL.Deleted.Set),
		F.Flow2(F.Bind2nd(cellAt, colCurator), cellsL.Curator.Set),
	})

	// recordHeader is the header of phenotype.csv and new-strain.csv. It
	// follows the order of recordFields.
	recordHeader = []string{
		"source_file", "source_row", "strain_id", "descriptor", "pheno_id",
		"assay_id", "env_id", "note", "pmid", "assigned_by",
	}

	// diagnosticHeader is the header of rejects.csv. It follows the order of
	// diagnosticFields.
	diagnosticHeader = []string{
		"source_file", "source_row", "severity", "code", "field", "raw_value",
	}
)

// cellAt returns the cell at index i, or an empty cell when the row is
// short or i is negative.
func cellAt(cells []string, i int) string {
	return F.Pipe2(cells, A.Lookup[string](i), O.GetOrElse(F.Constant("")))
}

// decodeHeader trims each header cell, accepts the alias for the
// misspelled environment column, and returns the layout.
func decodeHeader(cells []string) E.Either[error, Layout] {
	return F.Pipe3(
		cells,
		A.Map(F.Flow2(strings.TrimSpace, headerAlias.Replace)),
		chooseLayout,
		E.FromPredicate(knownLayout, F.Constant1[Layout](errBadLayout)),
	)
}

// recordFields returns one CSV line of an accepted record.
func recordFields(r Record) []string {
	return []string{
		recFile.Get(r),
		strconv.Itoa(recordRow.Get(r)),
		recordL.StrainID.Get(r),
		recordL.Descriptor.Get(r),
		recordL.PhenoID.Get(r),
		recordL.AssayID.Get(r),
		recordL.EnvID.Get(r),
		recordL.Note.Get(r),
		recordL.PMID.Get(r),
		recordL.AssignedBy.Get(r),
	}
}

// diagnosticFields returns one CSV line of a diagnostic.
func diagnosticFields(d Diagnostic) []string {
	return []string{
		diagFile.Get(d),
		strconv.Itoa(diagRow.Get(d)),
		string(diagL.Severity.Get(d)),
		string(diagL.Code.Get(d)),
		string(diagL.Field.Get(d)),
		diagL.Raw.Get(d),
	}
}

// recordRows returns the lines of a record file, header first.
func recordRows(records []Record) [][]string {
	return F.Pipe2(records, A.Map(recordFields), A.Prepend(recordHeader))
}

// diagnosticRows returns the lines of rejects.csv, header first.
func diagnosticRows(issues []Diagnostic) [][]string {
	return F.Pipe2(issues, A.Map(diagnosticFields), A.Prepend(diagnosticHeader))
}
```

### row.go

`internal/phenoannotation/normalize/row.go`. Predicates compare trimmed cells through composed lenses, and diagnostics show the raw cells. Each rule lives in `rules` in column order. `checkRule` applies one rule with `diagnose`, and `kindOf` selects the row kind with `MO.First`.

```go
package normalize

import (
	"net/mail"
	"regexp"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"
)

// Identifier formats and text predicates.
var (
	dbsID     = regexp.MustCompile(`^DBS[0-9]{7}$`)
	phenoID   = regexp.MustCompile(`^DDPHENO:[0-9]{7}$`)
	assayID   = regexp.MustCompile(`^DDASSAY:[0-9]{7}$`)
	charID    = regexp.MustCompile(`^DDSTRAINCHAR:[0-9]{7}$`)
	envID     = regexp.MustCompile(`^DDENVIR:[0-9]{7}$`)
	envTypoID = regexp.MustCompile(`^DDDENVIR:[0-9]{7}$`)
	positive  = regexp.MustCompile(`^0*[1-9][0-9]*$`)

	nonblank  = F.Flow2(strings.TrimSpace, Str.IsNonEmpty)
	stripPMID = F.Bind2nd(strings.TrimPrefix, "PMID:")
	dropTypoD = F.Bind2nd(strings.TrimPrefix, "D")

	// correctEnv changes the observed prefix DDDENVIR: to DDENVIR:.
	correctEnv = F.Pipe1(envTypoID.MatchString, P.Fold(F.Identity[string], dropTypoD))
	// supportedAssay clears a strain characteristic ID in the assay column.
	supportedAssay = F.Pipe1(
		charID.MatchString,
		P.Fold(F.Identity[string], F.Constant1[string]("")),
	)

	unnamedAddress = F.Flow2(addressName, Str.IsEmpty)
	// isBareEmail accepts one address with no display name.
	isBareEmail = F.Flow2(bareAddress, E.IsRight[error, string])
)

// Diagnostic and outcome predicates.
var (
	severityEq = EQ.FromStrictEquals[Severity]()
	kindEq     = EQ.FromStrictEquals[RowKind]()
	isError    = F.Pipe2(SevError, EQ.Equals(severityEq), P.ContraMap(diagL.Severity.Get))
	isWarning  = F.Pipe2(SevWarning, EQ.Equals(severityEq), P.ContraMap(diagL.Severity.Get))
)

// Row predicates. Each one reads the trimmed cells through a composed lens.
var (
	isBlankRow = F.Pipe1(
		F.Flow2(rowL.Clean.Get, P.IsZero[Cells]()),
		P.And(F.Flow2(rowL.Extra.Get, Str.IsEmpty)),
	)

	strainGiven   = F.Flow2(cleanStrain.Get, Str.IsNonEmpty)
	strainBlank   = P.Not(strainGiven)
	validStrainID = F.Flow2(cleanStrain.Get, dbsID.MatchString)
	badStrainID   = F.Pipe1(strainGiven, P.And(P.Not(validStrainID)))
	noDescriptor  = F.Pipe1(strainBlank, P.And(F.Flow2(cleanDescriptor.Get, Str.IsEmpty)))
	badPhenoID    = P.Not(F.Flow2(cleanPheno.Get, phenoID.MatchString))

	assayGiven     = F.Flow2(cleanAssay.Get, Str.IsNonEmpty)
	assayIsChar    = F.Flow2(cleanAssay.Get, charID.MatchString)
	assayIsTerm    = F.Flow2(cleanAssay.Get, assayID.MatchString)
	badAssayID     = F.Pipe2(assayGiven, P.And(P.Not(assayIsChar)), P.And(P.Not(assayIsTerm)))
	assayLabelOnly = F.Pipe1(
		P.Not(assayGiven),
		P.And(F.Flow2(cleanAssayTerm.Get, Str.IsNonEmpty)),
	)

	envGiven     = F.Flow2(cleanEnv.Get, Str.IsNonEmpty)
	envHasTypo   = F.Flow2(cleanEnv.Get, envTypoID.MatchString)
	envValid     = F.Flow3(cleanEnv.Get, correctEnv, envID.MatchString)
	badEnvID     = F.Pipe1(envGiven, P.And(P.Not(envValid)))
	envLabelOnly = F.Pipe1(P.Not(envGiven), P.And(F.Flow2(cleanEnvTerm.Get, Str.IsNonEmpty)))

	badPMID      = P.Not(F.Flow3(cleanReference.Get, stripPMID, positive.MatchString))
	deletedGiven = F.Flow2(cleanDeleted.Get, Str.IsNonEmpty)
	badCurator   = P.Not(F.Flow2(cleanCurator.Get, isBareEmail))
	extraGiven   = F.Flow2(rowL.Extra.Get, Str.IsNonEmpty)
)

// Rules and builders.
var (
	// rules lists the row checks in column order. Each broken rule gives one
	// diagnostic, so the diagnostics of a row are in this order.
	rules = []rule{
		newRule(badStrainID, SevError, CodeInvalidStrainID, FieldStrainID, rawStrain.Get),
		newRule(noDescriptor, SevError, CodeMissingDescriptor, FieldDescriptor, rawDescriptor.Get),
		newRule(badPhenoID, SevError, CodeInvalidPhenoID, FieldPhenoID, rawPheno.Get),
		newRule(assayIsChar, SevWarning, CodeStrainChar, FieldAssayID, rawAssay.Get),
		newRule(badAssayID, SevError, CodeInvalidAssayID, FieldAssayID, rawAssay.Get),
		newRule(assayLabelOnly, SevWarning, CodeLabelWithoutID, FieldAssayTerm, rawAssayTerm.Get),
		newRule(envHasTypo, SevWarning, CodeEnvPrefix, FieldEnvID, rawEnv.Get),
		newRule(badEnvID, SevError, CodeInvalidEnvID, FieldEnvID, rawEnv.Get),
		newRule(envLabelOnly, SevWarning, CodeLabelWithoutID, FieldEnvTerm, rawEnvTerm.Get),
		newRule(badPMID, SevError, CodeInvalidPMID, FieldPMID, rawReference.Get),
		newRule(deletedGiven, SevError, CodeDeletedInCreate, FieldDeleted, rawDeleted.Get),
		newRule(badCurator, SevError, CodeInvalidCurator, FieldAssignedBy, rawCurator.Get),
		newRule(extraGiven, SevError, CodeExtraCell, FieldOverflow, rowL.Extra.Get),
	}

	// toRecord builds the canonical record from the trimmed cells.
	toRecord = endoutils.BuildFrom([]func(SourceRow) En.Endomorphism[Record]{
		F.Flow2(rowL.Origin.Get, recordL.Origin.Set),
		F.Flow2(cleanStrain.Get, recordL.StrainID.Set),
		F.Flow2(cleanDescriptor.Get, recordL.Descriptor.Set),
		F.Flow2(cleanPheno.Get, recordL.PhenoID.Set),
		F.Flow3(cleanAssay.Get, supportedAssay, recordL.AssayID.Set),
		F.Flow3(cleanEnv.Get, correctEnv, recordL.EnvID.Set),
		F.Flow2(cleanNote.Get, recordL.Note.Set),
		F.Flow3(cleanReference.Get, stripPMID, recordL.PMID.Set),
		F.Flow2(cleanCurator.Get, recordL.AssignedBy.Set),
	})

	hasErrorIssue = A.Any(isError)
	blankOutcome  = F.Pipe1(RowOutcome{}, outcomeL.Kind.Set(KindBlank))

	// classify returns one outcome for each physical row.
	classify = F.Pipe1(
		isBlankRow,
		P.Fold(classifyData, F.Constant1[SourceRow](blankOutcome)),
	)
)

func addressName(a *mail.Address) string  { return a.Name }
func addressValue(a *mail.Address) string { return a.Address }

// bareAddress accepts one email address with no display name. The parsed
// address must equal the input, so a list or a comment fails.
func bareAddress(s string) E.Either[error, string] {
	return F.Pipe4(
		s,
		E.Eitherize1(mail.ParseAddress),
		E.Chain(E.FromPredicate(unnamedAddress, F.Constant1[*mail.Address](errDisplayName))),
		E.Map[error](addressValue),
		E.Chain(E.FromPredicate(F.Pipe1(s, predord.StrEq), F.Constant1[string](errChangedAddress))),
	)
}

// firstNonblank returns the first cell that is not blank, or "".
func firstNonblank(cells []string) string {
	return F.Pipe2(cells, A.FindFirst(nonblank), O.GetOrElse(F.Constant("")))
}

// newRule describes one row check.
func newRule(
	broken P.Predicate[SourceRow],
	sev Severity,
	code Code,
	field FieldName,
	raw func(SourceRow) string,
) rule {
	return F.Pipe5(
		rule{},
		ruleL.Broken.Set(broken),
		ruleL.Severity.Set(sev),
		ruleL.Code.Set(code),
		ruleL.Field.Set(field),
		ruleL.Raw.Set(raw),
	)
}

// diagnose makes the diagnostic of rule r for row.
func diagnose(r rule, row SourceRow) Diagnostic {
	return F.Pipe5(
		Diagnostic{},
		diagL.Origin.Set(rowL.Origin.Get(row)),
		diagL.Severity.Set(ruleL.Severity.Get(r)),
		diagL.Code.Set(ruleL.Code.Get(r)),
		diagL.Field.Set(ruleL.Field.Get(r)),
		diagL.Raw.Set(F.Pipe1(row, ruleL.Raw.Get(r))),
	)
}

// checkRule returns the diagnostic of rule r when row breaks r.
func checkRule(row SourceRow, r rule) O.Option[Diagnostic] {
	return F.Pipe2(
		row,
		O.FromPredicate(ruleL.Broken.Get(r)),
		O.Map(F.Bind1st(diagnose, r)),
	)
}

// kindOf rejects a row with an error before it looks at the strain ID.
func kindOf(row SourceRow, issues []Diagnostic) RowKind {
	return MO.First(KindExisting, []O.Option[RowKind]{
		F.Pipe1(issues, MO.Const(hasErrorIssue, KindRejected)),
		F.Pipe1(row, MO.Const(strainBlank, KindNew)),
	})
}

// classifyData checks a row that is not blank against every rule and
// classifies it.
func classifyData(row SourceRow) RowOutcome {
	issues := F.Pipe1(rules, A.FilterMap(F.Bind1st(checkRule, row)))

	return F.Pipe3(
		RowOutcome{},
		outcomeL.Kind.Set(kindOf(row, issues)),
		outcomeL.Record.Set(toRecord(row)),
		outcomeL.Issues.Set(issues),
	)
}
```

### source.go

`internal/phenoannotation/normalize/source.go`. `os.ReadDir` sorts by name. `fs.DirEntry.Info` returns `Lstat` data. Thus, a symbolic link is not a regular file. `workbookPath` joins the folder and file name. `selectSource` uses `P.Fold` on `hasInput` to select `--folder` or `--input`.

```go
package normalize

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predfs "github.com/dictyBase/fp-go-loom/predicate/fs"
	predstrings "github.com/dictyBase/fp-go-loom/predicate/strings"
)

var (
	hasCreatePrefix = F.Bind2nd(strings.HasPrefix, "Phenotype_annotation_")
	hasXlsxSuffix   = predstrings.HasSuffix(".xlsx")
	createName      = F.Pipe1(hasCreatePrefix, P.And(hasXlsxSuffix))
	createEntry     = F.Pipe1(createName, P.ContraMap(fs.DirEntry.Name))
	createFile      = F.Pipe1(
		predfs.IsRegularInfo,
		P.And(F.Pipe1(createName, P.ContraMap(fs.FileInfo.Name))),
	)

	hasInput   = F.Flow2(stateInput.Get, Str.IsNonEmpty)
	hasFolder  = F.Flow2(stateFolder.Get, Str.IsNonEmpty)
	hasOut     = F.Flow2(stateOut.Get, nonblank)
	inputOnly  = F.Pipe1(hasInput, P.And(P.Not(hasFolder)))
	folderOnly = F.Pipe1(hasFolder, P.And(P.Not(hasInput)))
	validSpec  = F.Pipe2(inputOnly, P.Or(folderOnly), P.And(hasOut))

	// Valid checks the option combination before any file access.
	Valid = E.FromPredicate(validSpec, F.Constant1[State](errSpec))

	// liveState fails when the run context is done.
	liveState = E.FromError(contextError)

	// selectSource reads one explicit file or lists one folder.
	selectSource = F.Pipe1(hasInput, P.Fold(folderFiles, explicitFile))
)

func contextError(s State) error { return stateL.Context.Get(s).Err() }

func notCreateFile(info fs.FileInfo) error {
	return fmt.Errorf("not a regular create workbook: %s", info.Name())
}

func noWorkbooks(dir string, _ []string) error {
	return fmt.Errorf("no create workbooks in %s", dir)
}

// workbookPath returns the path of a listed workbook in dir.
func workbookPath(dir string, info fs.FileInfo) string {
	return filepath.Join(dir, info.Name())
}

// lstatCreateFile checks one path without following a symbolic link.
func lstatCreateFile(path string) IOE.IOEither[error, fs.FileInfo] {
	return F.Pipe2(
		path,
		IOE.Eitherize1(os.Lstat),
		IOE.ChainEitherK(E.FromPredicate(createFile, notCreateFile)),
	)
}

// explicitFile accepts one regular create workbook.
func explicitFile(s State) IOE.IOEither[error, []string] {
	return F.Pipe4(
		s,
		stateInput.Get,
		IOE.Of[error, string],
		IOE.ChainFirst(lstatCreateFile),
		IOE.Map[error](A.Of[string]),
	)
}

// folderFiles lists create workbooks directly under the folder. ReadDir
// sorts by name. DirEntry.Info returns Lstat data, so a symbolic link is
// not a regular file.
func folderFiles(s State) IOE.IOEither[error, []string] {
	dir := stateFolder.Get(s)

	return F.Pipe6(
		dir,
		IOE.Eitherize1(os.ReadDir),
		IOE.Map[error](A.Filter(createEntry)),
		IOE.Chain(IOE.TraverseArraySeq(IOE.Eitherize1(fs.DirEntry.Info))),
		IOE.Map[error](A.Filter(predfs.IsRegularInfo)),
		IOE.Map[error](A.Map(F.Bind1st(workbookPath, dir))),
		IOE.ChainEitherK(E.FromPredicate(
			predarrays.IsNonEmpty[string](),
			F.Bind1st(noWorkbooks, dir),
		)),
	)
}

// SelectFiles returns the workbook paths in a stable order.
func SelectFiles(s State) IOE.IOEither[error, []string] {
	return F.Pipe3(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Chain(selectSource),
	)
}
```

### workbook.go

`internal/phenoannotation/normalize/workbook.go`. `readHeader` reads row 1 and decodes the layout. `readSheet` loops over the data rows of one worksheet and calls the row pipeline `addRow` for each row so the batch grows row by row. `ReadFiles` joins the per-file batches with `A.Reduce(joinBatch, Batch{})`. `IOE.Bracket` and `closeJoined` manage each workbook and row iterator so a close error joins the error of the work step. `pickSheet` chooses the worksheet from `--sheet` and the workbook sheet list. `openWorkbook`, `readHeader`, and `readSheet` are the only functions with allow directives.

```go
package normalize

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	"github.com/xuri/excelize/v2"
)

var (
	oneSheet = predarrays.LenEq[string](1)

	// lineOrigin is the source file and the row number of a line.
	lineOrigin = endoutils.BuildFrom([]func(line) En.Endomorphism[Origin]{
		F.Flow2(lineL.File.Get, originL.File.Set),
		F.Flow2(lineL.Row.Get, originL.Row.Set),
	})

	// toSourceRow keeps the raw cells and a trimmed copy. The trimmed copy
	// goes into the checks and the records. The raw copy goes into the
	// diagnostics.
	toSourceRow = endoutils.BuildFrom([]func(line) En.Endomorphism[SourceRow]{
		F.Flow2(lineOrigin, rowL.Origin.Set),
		F.Flow3(expandCells, toCells, rowL.Raw.Set),
		F.Flow4(expandCells, A.Map(strings.TrimSpace), toCells, rowL.Clean.Set),
		F.Flow2(overflow, rowL.Extra.Set),
	})

	wantsSheet = F.Pipe1(Str.IsNonEmpty, P.ContraMap(requestL.Want.Get))
	// chooseSheet uses the sole worksheet when the option is empty.
	chooseSheet = F.Pipe1(wantsSheet, P.Fold(onlySheet, namedSheet))
)

// expandCells puts the cells of one line in internal column order.
func expandCells(l line) []string {
	return F.Pipe1(
		layoutColumns[lineL.Layout.Get(l)],
		A.Map(F.Bind1st(cellAt, lineL.Cells.Get(l))),
	)
}

// overflow returns the first nonblank cell after the last header column.
func overflow(l line) string {
	return F.Pipe3(
		l,
		lineL.Cells.Get,
		A.SliceRight[string](layoutWidth[lineL.Layout.Get(l)]),
		firstNonblank,
	)
}

// newLine records one data row as read.
func newLine(file string, layout Layout, row int, cells []string) line {
	return F.Pipe4(
		line{},
		lineL.File.Set(file),
		lineL.Layout.Set(layout),
		lineL.Row.Set(row),
		lineL.Cells.Set(cells),
	)
}

func sheetCountError(names []string) error {
	return fmt.Errorf("%d worksheets: specify --sheet", len(names))
}

func missingSheet(want string) error { return fmt.Errorf("missing worksheet %q", want) }

func onlySheet(r sheetRequest) E.Either[error, string] {
	return F.Pipe4(
		r,
		requestL.Names.Get,
		E.FromPredicate(oneSheet, sheetCountError),
		E.Map[error](A.Head[string]),
		E.Chain(E.FromOption[string](F.Constant(errNoSheet))),
	)
}

func namedSheet(r sheetRequest) E.Either[error, string] {
	return F.Pipe2(
		r,
		requestL.Want.Get,
		E.FromPredicate(F.Bind1st(slices.Contains[[]string], requestL.Names.Get(r)), missingSheet),
	)
}

// openWorkbook closes the partial handle that excelize can return
// together with an open error.
//
// fp-go:allow-branch excelize can return a handle together with an error
func openWorkbook(path string) (*excelize.File, error) {
	book, err := excelize.OpenFile(path)
	if err != nil && book != nil {
		return nil, errors.Join(err, book.Close())
	}

	return book, err
}

// readHeader reads row 1 and returns the layout of the worksheet.
//
// fp-go:allow-branch excelize exposes only a stateful row iterator
func readHeader(rows *excelize.Rows) (Layout, error) {
	if !rows.Next() {
		return LayoutUnknown, errors.Join(errNoHeader, rows.Error())
	}
	cells, err := rows.Columns()
	if err != nil {
		return LayoutUnknown, fmt.Errorf("row %d: %w", headerRow, err)
	}

	return E.UnwrapError(decodeHeader(cells))
}

// readSheet reads one worksheet row by row. Each Next advances exactly
// one physical row, so the counter is the Excel row number. Each data
// row goes through the row pipeline addRow and into the batch.
//
// fp-go:allow-loop excelize exposes only a stateful row iterator
// fp-go:allow-branch the iterator reports an error for each row
func (s State) readSheet(file string, rows *excelize.Rows) (Batch, error) {
	layout, err := readHeader(rows)
	if err != nil {
		return Batch{}, err
	}
	ctx := stateL.Context.Get(s)
	batch := Batch{}
	for row := firstDataRow; rows.Next(); row++ {
		if err := ctx.Err(); err != nil {
			return Batch{}, err
		}
		cells, err := rows.Columns()
		if err != nil {
			return Batch{}, fmt.Errorf("row %d: %w", row, err)
		}
		batch = addRow(batch, newLine(file, layout, row, cells))
	}

	return batch, rows.Error()
}

// closeEffect closes a workbook, a row iterator, or a file.
func closeEffect[R io.Closer](r R) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) { return struct{}{}, r.Close() })
}

// joinUseError joins the error of the use step, if any, with a release
// error. IOE.Bracket otherwise returns only the release error.
func joinUseError[A any](used E.Either[error, A], releaseErr error) error {
	return errors.Join(E.ToError(used), releaseErr)
}

// closeJoined is the release step of IOE.Bracket for a closer. A close
// error never hides the error of the use step.
func closeJoined[R io.Closer, A any](r R, used E.Either[error, A]) IOE.IOEither[error, struct{}] {
	return F.Pipe2(
		r,
		closeEffect[R],
		IOE.MapLeft[struct{}](F.Bind1st(joinUseError[A], used)),
	)
}

func workbookError(path string, err error) error {
	return fmt.Errorf("workbook %s: %w", path, err)
}

// pickSheet chooses a worksheet from the --sheet option and the
// worksheet names of a workbook.
func pickSheet(want string, names []string) E.Either[error, string] {
	return F.Pipe3(
		sheetRequest{},
		requestL.Want.Set(want),
		requestL.Names.Set(names),
		chooseSheet,
	)
}

// scanSheet reads one worksheet and closes the row iterator. Excelize
// keeps the workbook path in book.Path.
func (s State) scanSheet(book *excelize.File, name string) IOE.IOEither[error, Batch] {
	return IOE.Bracket(
		F.Pipe1(name, IOE.Eitherize1(book.Rows)),
		F.Bind1st(IOE.Eitherize2(s.readSheet), filepath.Base(book.Path)),
		closeJoined[*excelize.Rows, Batch],
	)
}

// readBook chooses the worksheet of an open workbook and reads it.
func (s State) readBook(book *excelize.File) IOE.IOEither[error, Batch] {
	return F.Pipe3(
		book.GetSheetList(),
		IOE.Of[error, []string],
		IOE.ChainEitherK(F.Bind1st(pickSheet, stateSheet.Get(s))),
		IOE.Chain(F.Bind1st(s.scanSheet, book)),
	)
}

// bookBatch opens the workbook at path, reads it, and closes it on every
// path.
func (s State) bookBatch(path string) IOE.IOEither[error, Batch] {
	return IOE.Bracket(
		F.Pipe1(path, IOE.Eitherize1(openWorkbook)),
		s.readBook,
		closeJoined[*excelize.File, Batch],
	)
}

// readWorkbook checks the context, reads one workbook, and names the
// workbook in an error.
func (s State) readWorkbook(path string) IOE.IOEither[error, Batch] {
	return F.Pipe4(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Chain(F.Bind2nd(State.bookBatch, path)),
		IOE.MapLeft[Batch](F.Bind1st(workbookError, path)),
	)
}

// ReadFiles reads the selected workbooks one at a time, in order, and
// joins their batches.
func ReadFiles(s State) IOE.IOEither[error, Batch] {
	return F.Pipe3(
		s,
		stateL.Files.Get,
		IOE.TraverseArraySeq(s.readWorkbook),
		IOE.Map[error](A.Reduce(joinBatch, Batch{})),
	)
}
```

Excelize does not show all XML decoder errors through `Rows.Error`. A successful read is not proof of OOXML conformance.

### batch.go

`internal/phenoannotation/normalize/batch.go`. `addRow` is the row pipeline: it runs `toSourceRow`, `classify`, `outcomeBatch`, and `joinBatch` to add one row to the batch. `PrintSummary` writes one summary line to the state log writer as an `IO` effect.

```go
package normalize

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
)

const summaryFormat = "files=%d blank=%d existing=%d new=%d rejected=%d warnings=%d\n"

var (
	stateExisting = F.Pipe1(stateL.Batch, L.Compose[State](batchL.Existing))
	stateNew      = F.Pipe1(stateL.Batch, L.Compose[State](batchL.New))
	stateBlank    = F.Pipe1(stateL.Batch, L.Compose[State](batchL.Blank))
	stateRejected = F.Pipe1(stateL.Batch, L.Compose[State](batchL.Rejected))
	stateWarnings = F.Pipe1(stateL.Batch, L.Compose[State](batchL.Warnings))
)

// kindIs matches an outcome of the given kind.
func kindIs(kind RowKind) P.Predicate[RowOutcome] {
	return F.Pipe2(kind, EQ.Equals(kindEq), P.ContraMap(outcomeL.Kind.Get))
}

// recordIf returns the record of o when o has the given kind.
func recordIf(o RowOutcome, kind RowKind) []Record {
	return F.Pipe3(
		o,
		O.FromPredicate(kindIs(kind)),
		O.Map(outcomeL.Record.Get),
		O.Fold(A.Empty[Record], A.Of[Record]),
	)
}

// countIf returns 1 when o has the given kind, else 0.
func countIf(o RowOutcome, kind RowKind) int {
	return F.Pipe2(
		o,
		O.FromPredicate(kindIs(kind)),
		O.Fold(F.Constant(0), F.Constant1[RowOutcome](1)),
	)
}

// outcomeBatch turns one classified row into a batch of one row. A
// warning of a rejected row stays in the issues. A diagnostic is never
// counted as a rejected row.
func outcomeBatch(o RowOutcome) Batch {
	issues := outcomeL.Issues.Get(o)

	return F.Pipe6(
		Batch{},
		batchL.Existing.Set(recordIf(o, KindExisting)),
		batchL.New.Set(recordIf(o, KindNew)),
		batchL.Issues.Set(issues),
		batchL.Warnings.Set(F.Pipe2(issues, A.Filter(isWarning), A.Size[Diagnostic])),
		batchL.Blank.Set(countIf(o, KindBlank)),
		batchL.Rejected.Set(countIf(o, KindRejected)),
	)
}

// joinBatch appends batch b after batch a. It keeps the row order.
func joinBatch(a, b Batch) Batch {
	return F.Pipe6(
		a,
		batchL.Existing.Set(A.ArrayConcatAll(batchL.Existing.Get(a), batchL.Existing.Get(b))),
		batchL.New.Set(A.ArrayConcatAll(batchL.New.Get(a), batchL.New.Get(b))),
		batchL.Issues.Set(A.ArrayConcatAll(batchL.Issues.Get(a), batchL.Issues.Get(b))),
		batchL.Warnings.Set(batchL.Warnings.Get(a)+batchL.Warnings.Get(b)),
		batchL.Blank.Set(batchL.Blank.Get(a)+batchL.Blank.Get(b)),
		batchL.Rejected.Set(batchL.Rejected.Get(a)+batchL.Rejected.Get(b)),
	)
}

// addRow is the row pipeline: align the cells, classify the row, and
// add it to the batch.
func addRow(b Batch, l line) Batch {
	return F.Pipe4(l, toSourceRow, classify, outcomeBatch, F.Bind1st(joinBatch, b))
}

// PrintSummary writes one summary line to the log writer of the state.
func PrintSummary(s State) IO.IO[IO.Void] {
	return IO.FromImpure(func() {
		_, _ = fmt.Fprintf(
			stateL.Log.Get(s), summaryFormat,
			len(stateL.Files.Get(s)), stateBlank.Get(s),
			len(stateExisting.Get(s)), len(stateNew.Get(s)),
			stateRejected.Get(s), stateWarnings.Get(s),
		)
	})
}
```

### output.go

`internal/phenoannotation/normalize/output.go`. `batchFiles` builds the three CSV files from a batch. `writeCSV` writes rows with `encoding/csv` and default commas. `WriteFile` and `Publish` use `IOE.Bracket`. As a result, file close errors and temporary directory removal errors join the error of the work step. `renameInto` reads the context and does the atomic no-replace rename.

```go
package normalize

import (
	"encoding/csv"
	"os"
	"path/filepath"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	PR "github.com/IBM/fp-go/v2/pair"
)

const (
	filePerm       = 0o600
	createFlags    = os.O_CREATE | os.O_EXCL | os.O_WRONLY
	stagingPattern = ".pheno-*"
)

// batchFiles returns the three output files of a batch in dir.
func batchFiles(dir string, b Batch) []Output {
	return []Output{
		PR.MakePair(filepath.Join(dir, "phenotype.csv"), recordRows(batchL.Existing.Get(b))),
		PR.MakePair(filepath.Join(dir, "new-strain.csv"), recordRows(batchL.New.Get(b))),
		PR.MakePair(filepath.Join(dir, "rejects.csv"), diagnosticRows(batchL.Issues.Get(b))),
	}
}

// createExclusive creates path and fails when path exists.
func createExclusive(path string) IOE.IOEither[error, *os.File] {
	return IOE.TryCatchError(func() (*os.File, error) {
		return os.OpenFile(path, createFlags, filePerm)
	})
}

// writeCSV writes all rows to f. WriteAll flushes the writer and returns
// the flush error.
func writeCSV(f *os.File, rows [][]string) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) {
		return struct{}{}, csv.NewWriter(f).WriteAll(rows)
	})
}

// WriteFile creates the file of out exclusively, writes its rows as CSV,
// and closes the file. A close error joins a write error.
func WriteFile(out Output) IOE.IOEither[error, string] {
	return F.Pipe1(
		IOE.Bracket(
			createExclusive(PR.Head(out)),
			F.Bind2nd(writeCSV, PR.Tail(out)),
			closeJoined[*os.File, struct{}],
		),
		IOE.MapTo[error, struct{}](PR.Head(out)),
	)
}

// removeAll removes the staging directory. After a successful rename the
// directory does not exist, and RemoveAll returns nil.
func removeAll(dir string) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) { return struct{}{}, os.RemoveAll(dir) })
}

// removeJoined is the release step of IOE.Bracket for the staging
// directory. A removal error never hides the error of the commit.
func removeJoined[A any](dir string, used E.Either[error, A]) IOE.IOEither[error, struct{}] {
	return F.Pipe2(dir, removeAll, IOE.MapLeft[struct{}](F.Bind1st(joinUseError[A], used)))
}

// stagingDir checks the context and creates a new directory next to the
// destination.
func stagingDir(s State) IOE.IOEither[error, string] {
	return F.Pipe4(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Map[error](F.Flow3(stateOut.Get, filepath.Clean, filepath.Dir)),
		IOE.Chain(F.Bind2nd(IOE.Eitherize2(os.MkdirTemp), stagingPattern)),
	)
}

// writeOutput checks the context, then writes one file with the writer of
// the state.
func (s State) writeOutput(out Output) IOE.IOEither[error, string] {
	return F.Pipe3(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Chain(F.Constant1[State](F.Pipe1(out, stateL.Writer.Get(s)))),
	)
}

// renameInto checks the context and renames dir to the destination. It
// never replaces an existing destination. The rename is the only commit
// point.
func renameInto(dir string, s State) IOE.IOEither[error, State] {
	return F.Pipe5(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Map[error](F.Flow2(stateOut.Get, filepath.Clean)),
		IOE.Chain(F.Bind1st(IOE.Eitherize2(renameNew), dir)),
		IOE.MapTo[error, string](s),
	)
}

// commit writes the three files of the batch into dir and renames dir to
// the destination.
func (s State) commit(dir string) IOE.IOEither[error, State] {
	return F.Pipe4(
		s,
		stateL.Batch.Get,
		F.Bind1st(batchFiles, dir),
		IOE.TraverseArraySeq(s.writeOutput),
		IOE.Chain(F.Constant1[[]string](renameInto(dir, s))),
	)
}

// Publish writes the batch into a new sibling directory and renames it to
// the destination. It removes the staging directory on every path.
func Publish(s State) IOE.IOEither[error, State] {
	return IOE.Bracket(stagingDir(s), s.commit, removeJoined[State])
}
```

### rename_linux.go, rename_darwin.go, and rename_other.go

`internal/phenoannotation/normalize/rename_linux.go`:

```go
package normalize

import "golang.org/x/sys/unix"

// renameNew renames from to to and fails when to exists.
func renameNew(from, to string) (string, error) {
	return to, unix.Renameat2(unix.AT_FDCWD, from, unix.AT_FDCWD, to, unix.RENAME_NOREPLACE)
}
```

`internal/phenoannotation/normalize/rename_darwin.go`:

```go
package normalize

import "golang.org/x/sys/unix"

// renameNew renames from to to and fails when to exists.
func renameNew(from, to string) (string, error) {
	return to, unix.RenamexNp(from, to, unix.RENAME_EXCL)
}
```

`internal/phenoannotation/normalize/rename_other.go`. Do not use `os.Lstat` and then `os.Rename` as a fallback, because two runs can race.

```go
//go:build !linux && !darwin

package normalize

import "errors"

var errNoAtomicRename = errors.New(
	"atomic no-replace directory rename unavailable on this platform",
)

// renameNew fails, because this platform has no no-replace rename. Do not
// use os.Lstat followed by os.Rename: two runs can race.
func renameNew(string, string) (string, error) { return "", errNoAtomicRename }
```

The rename prevents a partial directory. It does not guarantee durability after a power failure.

### cli/arguments.go

`internal/phenoannotation/cli/arguments.go`. This file defines the boundary record and the `go:generate` directive for `internal/phenoannotation/cli/gen_lens.go`.

```go
// Package cli registers the phenoannotation commands.
package cli

import (
	"context"

	cli "github.com/urfave/cli/v3"
)

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

// arguments is the boundary record of one CLI invocation.
// fp-go:Lens
type arguments struct {
	Ctx context.Context
	Cmd *cli.Command
}
```

### cli/normalize.go

`internal/phenoannotation/cli/normalize.go`. `MakeLensRef` lenses build `*cli.Command` and `*cli.StringFlag`. `runNormalize` runs one `F.Pipe13` pipeline, sets `norm.WriteFile` on the state, binds `norm.ReadFiles` into `stateL.Batch`, and prints the summary one time after `Publish` through `IOE.ChainFirstIOK`.

```go
package cli

import (
	"context"
	"errors"
	"io"

	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	P "github.com/IBM/fp-go/v2/predicate"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	cli "github.com/urfave/cli/v3"

	norm "github.com/dictyBase/modware-import/internal/phenoannotation/normalize"
)

var (
	// ErrReview marks a complete batch that contains rejected rows.
	ErrReview = errors.New("batch contains rejected rows: review rejects.csv")
	errFlags  = errors.New("use one input flag and no positional arguments")
	errNoSub  = errors.New("expected normalize subcommand")

	argsL    = MakeargumentsLenses()
	stateL   = norm.MakeStateLenses()
	optionsL = norm.MakeOptionsLenses()
	batchL   = norm.MakeBatchLenses()

	stateRejected = F.Pipe1(stateL.Batch, L.Compose[norm.State](batchL.Rejected))

	cmdName     = L.MakeLensRef(getCmdName, setCmdName)
	cmdUsage    = L.MakeLensRef(getCmdUsage, setCmdUsage)
	cmdAction   = L.MakeLensRef(getCmdAction, setCmdAction)
	cmdFlags    = L.MakeLensRef(getCmdFlags, setCmdFlags)
	cmdWriter   = L.MakeLensRef(getCmdWriter, setCmdWriter)
	cmdCommands = L.MakeLensRef(getCmdCommands, setCmdCommands)

	flagName     = L.MakeLensRef(getFlagName, setFlagName)
	flagUsage    = L.MakeLensRef(getFlagUsage, setFlagUsage)
	flagRequired = L.MakeLensRef(getFlagRequired, setFlagRequired)

	// Flag shape: no positional argument and not both input flags. An
	// explicit empty value still counts as a set flag.
	flagSet      = (*cli.Command).IsSet
	noPositional = F.Flow2((*cli.Command).NArg, P.IsZero[int]())
	bothSources  = F.Pipe1(F.Bind2nd(flagSet, "folder"), P.And(F.Bind2nd(flagSet, "input")))
	plainCall    = F.Pipe1(noPositional, P.And(P.Not(bothSources)))

	requireFlagShape = E.FromPredicate(
		F.Pipe1(plainCall, P.ContraMap(argsL.Cmd.Get)),
		F.Constant1[arguments](errFlags),
	)

	// requireNoRejects turns a batch with rejected rows into ErrReview.
	// The batch is already published when this step runs.
	requireNoRejects = E.FromPredicate(
		F.Flow2(stateRejected.Get, P.IsZero[int]()),
		F.Constant1[norm.State](ErrReview),
	)

	flagText = (*cli.Command).String

	// toOptions copies the flag values into the options.
	toOptions = endoutils.BuildFrom([]func(*cli.Command) En.Endomorphism[norm.Options]{
		F.Flow2(F.Bind2nd(flagText, "input"), optionsL.Input.Set),
		F.Flow2(F.Bind2nd(flagText, "folder"), optionsL.Folder.Set),
		F.Flow2(F.Bind2nd(flagText, "sheet"), optionsL.Sheet.Set),
		F.Flow2(F.Bind2nd(flagText, "out"), optionsL.Out.Set),
	})

	// toState builds a new normalize state from one invocation.
	toState = endoutils.BuildFrom([]func(arguments) En.Endomorphism[norm.State]{
		F.Flow2(argsL.Ctx.Get, stateL.Context.Set),
		F.Flow3(argsL.Cmd.Get, toOptions, stateL.Options.Set),
		F.Flow3(argsL.Cmd.Get, rootWriter, stateL.Log.Set),
	})
)

func getCmdName(c *cli.Command) string                 { return c.Name }
func setCmdName(c *cli.Command, v string) *cli.Command { c.Name = v; return c }

func getCmdUsage(c *cli.Command) string                 { return c.Usage }
func setCmdUsage(c *cli.Command, v string) *cli.Command { c.Usage = v; return c }

func getCmdAction(c *cli.Command) cli.ActionFunc { return c.Action }

func setCmdAction(c *cli.Command, v cli.ActionFunc) *cli.Command { c.Action = v; return c }

func getCmdFlags(c *cli.Command) []cli.Flag                 { return c.Flags }
func setCmdFlags(c *cli.Command, v []cli.Flag) *cli.Command { c.Flags = v; return c }

func getCmdWriter(c *cli.Command) io.Writer                 { return c.Writer }
func setCmdWriter(c *cli.Command, v io.Writer) *cli.Command { c.Writer = v; return c }

func getCmdCommands(c *cli.Command) []*cli.Command { return c.Commands }

func setCmdCommands(c *cli.Command, v []*cli.Command) *cli.Command {
	c.Commands = v
	return c
}

func getFlagName(f *cli.StringFlag) string                     { return f.Name }
func setFlagName(f *cli.StringFlag, v string) *cli.StringFlag  { f.Name = v; return f }
func getFlagUsage(f *cli.StringFlag) string                    { return f.Usage }
func setFlagUsage(f *cli.StringFlag, v string) *cli.StringFlag { f.Usage = v; return f }
func getFlagRequired(f *cli.StringFlag) bool                   { return f.Required }

func setFlagRequired(f *cli.StringFlag, v bool) *cli.StringFlag {
	f.Required = v
	return f
}

// rootWriter returns the writer of the root command.
func rootWriter(c *cli.Command) io.Writer { return c.Root().Writer }

// runNormalize validates the flags, selects files, reads and classifies
// every row, publishes it, prints the summary, and returns the
// review status. It forces the effect once.
func runNormalize(ctx context.Context, c *cli.Command) error {
	return F.Pipe13(
		arguments{},
		argsL.Ctx.Set(ctx),
		argsL.Cmd.Set(c),
		IOE.Of[error, arguments],
		IOE.ChainEitherK(requireFlagShape),
		IOE.Map[error](F.Flow2(toState, stateL.Writer.Set(norm.WriteFile))),
		IOE.ChainEitherK(norm.Valid),
		IOE.Bind(stateL.Files.Set, norm.SelectFiles),
		IOE.Bind(stateL.Batch.Set, norm.ReadFiles),
		IOE.Chain(norm.Publish),
		IOE.ChainFirstIOK[error](norm.PrintSummary),
		IOE.ChainEitherK(requireNoRejects),
		ioeutils.ToEither[error, norm.State],
		E.ToError[norm.State],
	)
}

func needCommand(context.Context, *cli.Command) error { return errNoSub }

func textFlag(name, usage string, required bool) cli.Flag {
	return F.Pipe3(
		&cli.StringFlag{},
		flagName.Set(name),
		flagUsage.Set(usage),
		flagRequired.Set(required),
	)
}

func normalizeCommand() *cli.Command {
	return F.Pipe4(
		&cli.Command{},
		cmdName.Set("normalize"),
		cmdUsage.Set("convert create workbooks into reviewable CSV files"),
		cmdAction.Set(runNormalize),
		cmdFlags.Set([]cli.Flag{
			textFlag("input", "one create workbook", false),
			textFlag("folder", "folder of create workbooks", false),
			textFlag("out", "new output directory", true),
			textFlag("sheet", "exact worksheet name", false),
		}),
	)
}

// Command returns the root command. It writes the summary to w.
func Command(w io.Writer) *cli.Command {
	return F.Pipe4(
		&cli.Command{},
		cmdName.Set("phenoannotation"),
		cmdAction.Set(needCommand),
		cmdWriter.Set(w),
		cmdCommands.Set([]*cli.Command{normalizeCommand()}),
	)
}
```

### cmd/phenoannotation/main.go

`cmd/phenoannotation/main.go`. `exitCode` maps the error with `MO.Match`, and `reportFailure` prints to stderr with `P.Fold`. `main` calls `stop()` directly before `os.Exit` because `os.Exit` skips deferred calls.

```go
// Command phenoannotation prepares phenotype annotation batches.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
	MO "github.com/dictyBase/fp-go-loom/matchopt"

	app "github.com/dictyBase/modware-import/internal/phenoannotation/cli"
)

// Exit codes of the process.
const (
	exitOK     = 0
	exitFailed = 1
	exitReview = 2
)

var (
	isReview = F.Bind2nd(errors.Is, app.ErrReview)

	// exitCode maps the command error to the process exit code.
	exitCode = MO.Match(exitFailed, []O.Kleisli[error, int]{
		MO.Const(P.IsZero[error](), exitOK),
		MO.Const(isReview, exitReview),
	})

	// reportFailure prints a command error to stderr and keeps it.
	reportFailure = F.Pipe1(P.IsNonZero[error](), P.Fold(F.Identity[error], printFailure))
)

// printFailure is the process edge: it writes err to stderr.
func printFailure(err error) error {
	_, _ = fmt.Fprintln(os.Stderr, err)
	return err
}

// main runs the command. A command action never calls os.Exit, so the
// cleanup of the action runs before main stops the signal handler.
func main() {
	background := context.Background()
	ctx, stop := signal.NotifyContext(background, os.Interrupt, syscall.SIGTERM)
	err := app.Command(os.Stdout).Run(ctx, os.Args)
	stop()
	os.Exit(F.Pipe2(err, reportFailure, exitCode))
}
```

### normalize_test.go

`internal/phenoannotation/normalize/normalize_test.go`. These 13 tests cover CSV column widths, headers, short-row alignment, validation, batch counts, file selection, and workbook failures. They also cover publication, cancellation, joined release errors, exclusive file creation, and process exit codes.

```go
package normalize

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	PR "github.com/IBM/fp-go/v2/pair"
	"github.com/xuri/excelize/v2"
)

var outputNames = []string{"phenotype.csv", "new-strain.csv", "rejects.csv"}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func force[A any](ioe IOE.IOEither[error, A]) (A, error) { return E.UnwrapError(ioe()) }

// dataRow builds a row in the extended layout with the production aligner.
func dataRow(cells ...string) SourceRow {
	return toSourceRow(newLine("Phenotype_annotation_T.xlsx", LayoutExtended, firstDataRow, cells))
}

func fixture(t *testing.T, path string, rows ...[]string) {
	t.Helper()
	f := excelize.NewFile()
	for i, r := range rows {
		if len(r) == 0 {
			continue
		}
		axis, err := excelize.CoordinatesToCellName(1, i+1)
		check(t, err)
		check(t, f.SetSheetRow("Sheet1", axis, &r))
	}
	check(t, f.SaveAs(path))
	check(t, f.Close())
}

func readCSV(t *testing.T, path string) [][]string {
	t.Helper()
	f, err := os.Open(path)
	check(t, err)
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	check(t, err)
	check(t, f.Close())
	return records
}

func codes(o RowOutcome) []Code {
	var out []Code
	for _, d := range o.Issues {
		out = append(out, d.Code)
	}
	return out
}

func emptyDir(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	check(t, err)
	if len(entries) != 0 {
		t.Fatal(entries)
	}
}

func baseState(out string) State {
	var s State
	s.Context = context.Background()
	s.Options.Out = out
	s.Writer = WriteFile
	return s
}

func invoke(t *testing.T, bin string, want int, args ...string) string {
	t.Helper()
	buf, err := exec.Command(bin, args...).CombinedOutput()
	code := 0
	if exit, ok := errors.AsType[*exec.ExitError](err); ok {
		code = exit.ExitCode()
	} else {
		check(t, err)
	}
	if code != want {
		t.Fatalf("exit=%d want=%d: %s", code, want, buf)
	}
	return string(buf)
}

// TestCSVColumns keeps each header and its field function the same width.
func TestCSVColumns(t *testing.T) {
	if len(recordHeader) != len(recordFields(Record{})) {
		t.Fatal("record header and fields differ")
	}
	if len(diagnosticHeader) != len(diagnosticFields(Diagnostic{})) {
		t.Fatal("diagnostic header and fields differ")
	}
}

func TestHeader(t *testing.T) {
	alias := slices.Clone(header10)
	alias[6] = " Envirionment ID "
	layout, err := E.UnwrapError(decodeHeader(alias))
	check(t, err)
	if layout != LayoutShort {
		t.Fatal(layout)
	}
	alias[6] = alias[7]
	if _, err = E.UnwrapError(decodeHeader(alias)); !errors.Is(err, errBadLayout) {
		t.Fatal("duplicate header accepted")
	}
}

func TestAlignShortRow(t *testing.T) {
	r := toSourceRow(newLine("a.xlsx", LayoutShort, 4, []string{
		" DBS0000001 ", "DDPHENO:0000001", "", "", "", "", "", "", "1", "a@b.org", "", "x",
	}))
	if r.Origin.Row != 4 || r.Raw.Strain != " DBS0000001 " || r.Clean.Strain != "DBS0000001" {
		t.Fatal(r)
	}
	if r.Raw.Descriptor != "" || r.Clean.Curator != "a@b.org" || r.Extra != "x" {
		t.Fatal(r)
	}
}

func TestValidTerms(t *testing.T) {
	o := classify(dataRow(
		"DBS0000001", "", "DDPHENO:0000001", "", "", "DDASSAY:0000057", "",
		"DDENVIR:0000026", "", "36982997", "", "a@b.org",
	))
	if o.Kind != KindExisting || len(o.Issues) != 0 {
		t.Fatal(o)
	}
	if o.Record.AssayID != "DDASSAY:0000057" || o.Record.EnvID != "DDENVIR:0000026" {
		t.Fatal(o.Record)
	}
}

func TestRowFaults(t *testing.T) {
	cells := []string{
		"DBS0000001", "", "DDPHENO:0000001", "", "", "", "assay label", "",
		"environment label", "PMID:0", "", "Display Name <curator@example.org>", "outside",
	}
	o := classify(dataRow(cells...))
	want := []Code{
		CodeLabelWithoutID, CodeLabelWithoutID, CodeInvalidPMID,
		CodeInvalidCurator, CodeExtraCell,
	}
	if o.Kind != KindRejected || !slices.Equal(codes(o), want) {
		t.Fatal(o)
	}
	cells = cells[:12]
	cells[5], cells[7], cells[9], cells[11] = "BAD:0000001", "DDDENVIR:xyz", "1", "c@a.org,c@b.org"
	want = []Code{CodeInvalidAssayID, CodeInvalidEnvID, CodeInvalidCurator}
	if got := codes(classify(dataRow(cells...))); !slices.Equal(got, want) {
		t.Fatal(got)
	}
	got := codes(classify(dataRow(
		"", "", "DDPHENO:0000001", "", "", "", "", "", "", "1", "", "a@b.org",
	)))
	if !slices.Equal(got, []Code{CodeMissingDescriptor}) {
		t.Fatal(got)
	}
}

func TestAccounting(t *testing.T) {
	cells := []string{
		"", "same", "DDPHENO:0000001", "", "", "DDSTRAINCHAR:0000002", "", "",
		"", "1", "", "x@y.org",
	}
	fresh := dataRow(cells...)
	cells[0] = "DBS0000001"
	known := dataRow(cells...)
	var b Batch
	for _, r := range []SourceRow{fresh, fresh, known, {}} {
		b = joinBatch(b, outcomeBatch(classify(r)))
	}
	got := []int{len(b.Existing), len(b.New), b.Blank, b.Warnings, b.Rejected}
	if !slices.Equal(got, []int{1, 2, 1, 3, 0}) {
		t.Fatal(b)
	}
	if b.New[0].Descriptor != b.New[1].Descriptor || b.New[0].AssayID != "" {
		t.Fatal(b.New)
	}
}

func TestSelectFiles(t *testing.T) {
	dir := t.TempDir()
	s := stateFolder.Set(dir)(baseState(""))
	if _, err := force(SelectFiles(s)); err == nil {
		t.Fatal("empty selection accepted")
	}
	good := filepath.Join(dir, "Phenotype_annotation_B.xlsx")
	fixture(t, good, header10)
	update := filepath.Join(dir, "Phenotypes_update_JUN_11_2023.xlsx")
	check(t, os.WriteFile(update, nil, 0o600))
	check(t, os.WriteFile(filepath.Join(dir, "~$Phenotype_annotation_B.xlsx"), nil, 0o600))
	link := filepath.Join(dir, "Phenotype_annotation_A.xlsx")
	check(t, os.Symlink(good, link))
	files, err := force(SelectFiles(s))
	check(t, err)
	if !slices.Equal(files, []string{good}) {
		t.Fatal(files)
	}
	files, err = force(SelectFiles(stateInput.Set(good)(baseState(""))))
	check(t, err)
	if !slices.Equal(files, []string{good}) {
		t.Fatal(files)
	}
	for _, input := range []string{link, update} {
		if _, err := force(SelectFiles(stateInput.Set(input)(baseState("")))); err == nil {
			t.Fatalf("accepted %s", input)
		}
	}
}

func TestWorkbookFailures(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Phenotype_annotation_A.xlsx")
	s := baseState("")
	check(t, os.WriteFile(path, []byte("not a zip archive"), 0o600))
	if _, err := force(s.readWorkbook(path)); err == nil {
		t.Fatal("corrupt archive accepted")
	}
	fixture(t, path)
	if _, err := force(s.readWorkbook(path)); !errors.Is(err, errNoHeader) {
		t.Fatal("missing header accepted")
	}
	f := excelize.NewFile()
	_, err := f.NewSheet("Curated")
	check(t, err)
	check(t, f.SetSheetRow("Curated", "A1", &header10))
	row := []string{"DBS0000001"}
	check(t, f.SetSheetRow("Curated", "A4", &row))
	check(t, f.SaveAs(path))
	check(t, f.Close())
	if _, err = force(s.readWorkbook(path)); err == nil {
		t.Fatal("ambiguous sheet accepted")
	}
	if _, err = force(stateSheet.Set("absent")(s).readWorkbook(path)); err == nil {
		t.Fatal("missing sheet accepted")
	}
	b, err := force(stateSheet.Set("Curated")(s).readWorkbook(path))
	check(t, err)
	if b.Blank != 2 || b.Rejected != 1 || b.Issues[0].Origin.Row != 4 {
		t.Fatal(b)
	}
}

func failWriter(Output) IOE.IOEither[error, string] { return IOE.Left[string](io.ErrClosedPipe) }

func TestPublish(t *testing.T) {
	d := t.TempDir()
	dst := filepath.Join(d, "out")
	_, err := force(Publish(stateL.Writer.Set(failWriter)(baseState(dst))))
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatal(err)
	}
	emptyDir(t, d)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = force(Publish(stateL.Context.Set(ctx)(baseState(dst))))
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	emptyDir(t, d)
	_, err = force(Publish(baseState(dst)))
	check(t, err)
	for _, name := range outputNames {
		if len(readCSV(t, filepath.Join(dst, name))) != 1 {
			t.Fatal(name)
		}
	}
	if _, err = force(Publish(baseState(dst))); err == nil {
		t.Fatal("replaced an existing destination")
	}
}

func TestCancelDuringPublish(t *testing.T) {
	d := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cancelAfter := func(o Output) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			defer cancel()
			return WriteFile(o)()
		}
	}
	s := baseState(filepath.Join(d, "out"))
	s.Context = ctx
	s.Writer = cancelAfter
	if _, err := force(Publish(s)); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	emptyDir(t, d)
}

type badCloser struct{}

func (badCloser) Close() error { return io.ErrClosedPipe }

func failUse(badCloser) IOE.IOEither[error, int] { return IOE.Left[int](io.ErrUnexpectedEOF) }

func TestReleaseJoinsErrors(t *testing.T) {
	_, err := force(IOE.Bracket(IOE.Of[error](badCloser{}), failUse, closeJoined[badCloser, int]))
	if !errors.Is(err, io.ErrUnexpectedEOF) || !errors.Is(err, io.ErrClosedPipe) {
		t.Fatal(err)
	}
	dir := t.TempDir()
	_, err = force(IOE.Bracket(IOE.Of[error](dir), failDir, removeJoined[int]))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(dir); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatal(statErr)
	}
}

func failDir(string) IOE.IOEither[error, int] { return IOE.Left[int](io.ErrUnexpectedEOF) }

func TestWriteFileExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.csv")
	out := PR.MakePair(path, [][]string{{"a"}})
	if _, err := force(WriteFile(out)); err != nil {
		t.Fatal(err)
	}
	if _, err := force(WriteFile(out)); !errors.Is(err, os.ErrExist) {
		t.Fatal(err)
	}
}

func buildBinary(t *testing.T, dir string) string {
	t.Helper()
	bin := filepath.Join(dir, "phenoannotation")
	out, err := exec.Command("go", "build", "-o", bin, "../../../cmd/phenoannotation").
		CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v: %s", err, out)
	}
	return bin
}

const processNote = "first\t\"quoted\"\nlast"

// writeSources writes one short and one extended workbook and an update file.
func writeSources(t *testing.T, src string) (string, string) {
	t.Helper()
	a := filepath.Join(src, "Phenotype_annotation_A.xlsx")
	z := filepath.Join(src, "Phenotype_annotation_Z.xlsx")
	valid := []string{
		"DBS0000001", "DDPHENO:0000001", "", processNote, "DDSTRAINCHAR:0000002", "",
		"DDDENVIR:0000001", "", "PMID:123", " a@b.org\u00a0",
	}
	fixture(t, z, header10, valid, nil, valid)
	fresh := []string{"", "same", "DDPHENO:0000001", "", "", "", "", "", "", "42", "", "c@x.org"}
	bad := slices.Clone(fresh)
	bad[0], bad[10], bad[11] = "DBS1", "x", ""
	fixture(t, a, header12, fresh, fresh, bad)
	update := filepath.Join(src, "Phenotypes_update_JUN_11_2023.xlsx")
	check(t, os.WriteFile(update, []byte("ignore"), 0o600))
	return a, z
}

func checkClean(t *testing.T, dir string) {
	t.Helper()
	p := readCSV(t, filepath.Join(dir, "phenotype.csv"))
	header := recordHeader
	want := []string{
		"Phenotype_annotation_Z.xlsx", "2", "DBS0000001", "", "DDPHENO:0000001", "",
		"DDENVIR:0000001", processNote, "123", "a@b.org",
	}
	if len(p) != 3 || !slices.Equal(p[0], header) || !slices.Equal(p[1], want) || p[2][1] != "4" {
		t.Fatal(p)
	}
}

func checkFolder(t *testing.T, first, second string) {
	t.Helper()
	for _, name := range outputNames {
		x, err := os.ReadFile(filepath.Join(first, name))
		check(t, err)
		y, err := os.ReadFile(filepath.Join(second, name))
		check(t, err)
		if string(x) != string(y) {
			t.Fatal(name)
		}
	}
	n := readCSV(t, filepath.Join(first, "new-strain.csv"))
	if len(n) != 3 || n[1][3] != n[2][3] || n[2][1] != "3" {
		t.Fatal(n)
	}
	var got []string
	for _, d := range readCSV(t, filepath.Join(first, "rejects.csv"))[1:] {
		got = append(got, d[3])
	}
	want := []string{
		"invalid_strain_id", "deleted_in_create", "invalid_assigned_by",
		"unsupported_strain_characteristic", "corrected_env_prefix",
		"unsupported_strain_characteristic", "corrected_env_prefix",
	}
	if !slices.Equal(got, want) {
		t.Fatal(got)
	}
}

func TestProcess(t *testing.T) {
	d := t.TempDir()
	src := filepath.Join(d, "src")
	check(t, os.Mkdir(src, 0o700))
	bin := buildBinary(t, d)
	a, z := writeSources(t, src)

	clean := filepath.Join(d, "clean")
	invoke(t, bin, 0, "normalize", "--input", z, "--out", clean)
	checkClean(t, clean)

	first, second := filepath.Join(d, "first"), filepath.Join(d, "second")
	for _, dest := range []string{first, second} {
		msg := invoke(t, bin, 2, "normalize", "--folder", src, "--out", dest)
		if !strings.Contains(msg, "files=2 blank=1 existing=2 new=2 rejected=1 warnings=4") {
			t.Fatal(msg)
		}
	}
	checkFolder(t, first, second)

	unused := filepath.Join(d, "unused")
	invoke(t, bin, 1, "normalize", "--input", a, "--out", first)
	invoke(t, bin, 1, "normalize", "--input", a, "--folder", src, "--out", unused)
	invoke(t, bin, 1, "normalize", "--input", z, "--sheet", "absent", "--out", unused)
	invoke(t, bin, 1, "normalize", "--typo")
	invoke(t, bin, 1)
	if _, err := os.Stat(unused); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}
```

### cli/normalize_test.go

`internal/phenoannotation/cli/normalize_test.go`. These tests cover cancellation, the subcommand requirement, and `requireNoRejects`.

```go
package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	E "github.com/IBM/fp-go/v2/either"

	norm "github.com/dictyBase/modware-import/internal/phenoannotation/normalize"
)

func TestCancelledRun(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := Command(io.Discard).Run(ctx, []string{
		"phenoannotation", "normalize", "--folder", dir, "--out", out,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

func TestNeedsSubcommand(t *testing.T) {
	err := Command(io.Discard).Run(context.Background(), []string{"phenoannotation"})
	if !errors.Is(err, errNoSub) {
		t.Fatal(err)
	}
}

func TestReviewStatus(t *testing.T) {
	var s norm.State
	if E.IsLeft(requireNoRejects(s)) {
		t.Fatal("clean batch rejected")
	}
	_, err := E.UnwrapError(requireNoRejects(stateRejected.Set(1)(s)))
	if !errors.Is(err, ErrReview) {
		t.Fatal(err)
	}
}
```

### cli/gates_test.go

`internal/phenoannotation/cli/gates_test.go`. The gate runs all `pipelinecheck` v0.10.0 rules and `loomcheck` on `cli`, `normalize`, and `cmd/phenoannotation`.

```go
package cli

import (
	"strings"
	"testing"

	"github.com/dictyBase/fp-go-loom/loomcheck"
	pc "github.com/dictyBase/fp-go-loom/pipelinecheck"
)

// roots lists every package of the phenoannotation command.
var roots = []string{".", "../normalize", "../../../cmd/phenoannotation"}

func isEntrypoint(name string) bool {
	return strings.HasPrefix(name, "run") || strings.HasPrefix(name, "Run")
}

// TestFunctionalGates runs every pipelinecheck rule and loomcheck on the
// three packages. The checks are syntax only. A pass is not a proof.
func TestFunctionalGates(t *testing.T) {
	var cfg pc.Config
	cfg.Roots = roots
	cfg.IsEntrypoint = isEntrypoint
	cfg.RequirePointFreeSeed = true
	cfg.RequirePointFreeBranching = true
	cfg.RequireTrivialDelegate = true
	cfg.TerminalFunctions = []string{"Publish", "PrintSummary", "WriteFile"}
	pc.Require(t, cfg)
	pc.RequireBareIf(t, cfg)
	pc.RequireNoIfErrInTryCatch(t, cfg)
	pc.RequireNoNonRawTryCatchCallback(t, cfg)
	pc.RequireSafePrintf(t, cfg)
	pc.RequireNoRedundantFold(t, cfg)
	pc.RequireNoHandRolledFileJoinInArrayMap(t, cfg)
	pc.RequireNoImperativeLoop(t, cfg)
	pc.RequireNoImperativeBranch(t, cfg)
	pc.RequireNoDirectApply(t, cfg)
	pc.RequireNoNewPipe(t, cfg)
	pc.RequireNoBoolExpr(t, cfg)
	pc.RequireNoPopulatedStruct(t, cfg)
	pc.RequireNoLocalPipeClosure(t, cfg)
	pc.RequireNoFlap(t, cfg)
	pc.RequireTerminalEffect(t, cfg)
	var lc loomcheck.Config
	lc.Roots = roots
	loomcheck.Require(t, lc)
}
```

## fp-go v2 rules for this package

- Load `fp-go`, `fp-go-pipe-flow`, `fp-go-predicates`, `fp-go-lens`, and `fp-go-loom` before you change the code.
- Use a named plain function, not a point-free combinator chain. Do not use `Flap`, `Flip`, or `Pair` carriers only to keep a value without a name. `RequireNoFlap` in `TestFunctionalGates` enforces this rule for `Flap`, `MonadFlap`, and `F.Flip`.
- A lens `Set` can accept a getter call as its argument. That pattern in `diagnose`, `classifyData`, `outcomeBatch`, and `joinBatch` is not a nested-call violation.
- Use typed domain values (`Layout`, `RowKind`, `Severity`, `Code`, `FieldName`), not interchangeable strings or integers.
- Use named record structs (`Origin`, `Cells`, `SourceRow`, `Record`, `Diagnostic`, `RowOutcome`, `Batch`, `Options`, `State`, `line`), not positional arrays. Convert to string slices only at the CSV output boundary with `recordHeader`, `recordFields`, `diagnosticHeader`, and `diagnosticFields`.
- Mark record structs with `// fp-go:Lens`. Run `go generate ./internal/phenoannotation/...` to create `gen_lens.go`. Commit the generated files. Use `L.Compose` to reach nested fields through `State`, `SourceRow`, `Record`, and `Diagnostic`. Use `L.MakeLensRef` for external pointer structs (`*cli.Command`, `*cli.StringFlag`).
- In non-test code, write boolean expressions (`==`, `!=`, `&&`, `||`, `!`) only in named raw adapters (`openWorkbook`, `readHeader`, `readSheet`). Each raw adapter must have an `fp-go:allow-branch` or `fp-go:allow-loop` directive. Build named predicate variables from lens getters, `Str.IsNonEmpty`, `Str.IsEmpty`, and `regexp.Regexp.MatchString`. Combine them with `P.And`, `P.Or`, `P.Not`, `P.ContraMap`, `P.IsZero`, `P.IsNonZero`, and `EQ.Equals` over typed `Eq` instances.
- Put row checks as data in one `rules` list of `rule` records with a `Broken` predicate, in column order. Run them in `classifyData` through `checkRule` and `diagnose`.
- Do not keep an unexported helper with one caller when its body is one `F.Pipe1`. Move its steps into the caller. The `TryCatch*` and `MapLeft` steps are not inside this rule.
- Use `endoutils.BuildFrom` when each field of a record comes independently from the same input. Examples: `toCells`, `toSourceRow`, `lineOrigin`, `toRecord`, `toOptions`, and `toState`.
- Use `IOE.Bind` for dependent effect steps (`SelectFiles` into `stateL.Files`, `ReadFiles` into `stateL.Batch`).
- Use `P.Fold` for two-way dispatch (`classify`, `selectSource`, `chooseSheet`, `correctEnv`, `supportedAssay`, `reportFailure`). Use `MO.Match` (`chooseLayout`, `exitCode`) or `MO.First` (`kindOf`) for multi-way dispatch.
- Return `E.Either` from pure fallible steps. Lift them into `IOEither` pipelines with `IOE.ChainEitherK`. Do not bridge pure helpers with `IOE.FromEither`.
- Use `IOE.Bracket` for each resource with a release step (workbooks, row iterators, output files, and temporary directories). Join the error of the work step and the release error with `errors.Join` (`joinUseError`, `closeJoined`, `removeJoined`).
- Keep side effects out of `E.Fold` arms. Print the summary one time after `Publish` through `IOE.ChainFirstIOK[error](norm.PrintSummary)`. Force the `IOEither` one time at the CLI boundary with `ioeutils.ToEither` and `E.ToError`.
- Use `F.Pipe1(value, comb(args))`, not direct application `comb(args)(value)`. Do not write populated struct literals, local pipe closures, or trivial wrapper functions with one call in non-test code.
- Put an imperative loop or branch only in a named adapter for a stateful external API. Keep the loop body as a call to a pure row pipeline (`readSheet` calls `addRow` for each row so the batch grows row by row, and `ReadFiles` joins per-file batches with `A.Reduce(joinBatch, Batch{})`). Only three adapters carry directives: `openWorkbook` (`fp-go:allow-branch`), `readHeader` (`fp-go:allow-branch`), and `readSheet` (`fp-go:allow-loop` and `fp-go:allow-branch`). Give each directive a clear reason.
- Run `TestFunctionalGates` (`pipelinecheck` and `loomcheck`) on `cli`, `normalize`, and `cmd/phenoannotation`.

## Verification record

The scratch copy passed these gates on 2026-11-23 on Darwin/arm64 with Go 1.26.10:

| Gate | Result |
|---|---|
| `go test ./internal/phenoannotation/... -count=1` | 17 tests passed in 2 packages (13 in `normalize`, 4 in `cli`). Includes `TestCSVColumns` and `TestFunctionalGates`. |
| `golangci-lint` v2.14.0, the local tool (`fmt` and `run` on the 3 packages) | 0 issues. `fmt` made no change. |
| `gopls check -severity=hint` on each non-generated file | No diagnostics. |
| `CGO_ENABLED=0 GOOS=linux go vet ./internal/phenoannotation/... ./cmd/phenoannotation/` | Passed. |
| `CGO_ENABLED=0 GOOS=windows go build ./internal/phenoannotation/...` | Passed (`rename_other.go`). |
| `go build ./...` (full module) | Passed. |
| `rg` check for `I.Flap`, `A.Flap`, `F.Flip`, `PR.MapTail`, `PR.MapHead`, `PR.Of` | No hits in `internal/phenoannotation` or `cmd/phenoannotation`. |
| `TestFunctionalGates` (`pipelinecheck` v0.10.0 and `loomcheck`) | Passed on `.`, `../normalize`, and `../../../cmd/phenoannotation`. |

With Go 1.26.10 on Darwin, `GOOS=linux go vet` and `GOOS=windows go build` try to compile `runtime/cgo` with the macOS SDK and fail. Set `CGO_ENABLED=0` for the two cross-platform commands.

The CI pin in `.github/workflows/lint-test-build.yml` is `golangci-lint` v2.13.2. `golangci-lint` v2.12.2 bundles `mvdan.cc/unparam` v0.0.0-20251027182757. That `unparam` version does not unalias the generic type alias `E.Either[error, A]`. Thus it measures the size of a type parameter and panics with `go/types/gcsizes.go:17: assertion failed` on `closeJoined` and `removeJoined`. `golangci-lint` v2.13.2 and higher versions bundle `unparam` v0.0.0-20260823230713, which has the fix. Do not use v2.12.2.

`TestFunctionalGates` runs `loomcheck` and all `pipelinecheck` v0.10.0 rules on all three roots:

- `Require`, `BareIf`, `NoIfErrInTryCatch`, `NoNonRawTryCatchCallback`, `SafePrintf`, `NoRedundantFold`, `NoHandRolledFileJoinInArrayMap`
- `NoImperativeLoop`, `NoImperativeBranch`, `NoDirectApply`, `NoNewPipe`, `NoBoolExpr`, `NoPopulatedStruct`, `NoLocalPipeClosure`, `NoFlap`
- `TerminalEffect` with `TerminalFunctions` `Publish`, `PrintSummary`, and `WriteFile`

A probe added a temporary file with a boolean-expression predicate, an `if` statement, and a populated struct literal. `TestFunctionalGates` failed 3 times, which showed that the gate is live.

We ran the corpus check with the revised code and no service connection:

```sh
/tmp/phenoannotation normalize --folder ../../kosher-files --out /tmp/pheno-review-01
```

It returned exit `0` and this summary:

```text
files=77 blank=683 existing=368 new=64 rejected=0 warnings=13
```

A second probe added a file with `A.Flap`. `RequireNoFlap` failed one time, which confirmed that the rule is live.

`rejects.csv` contained 10 `corrected_env_prefix` records and 3 `unsupported_strain_characteristic` records. The decoded CSV content equals the earlier TSV output on the full corpus. The files hold 369 rows in `phenotype.csv`, 65 rows in `new-strain.csv`, and 14 rows in `rejects.csv`, with headers.

After the row-loop revision (`readSheet`, `addRow`, and `joinBatch`), all gates above (`go test`, `golangci-lint` v2.14.0, `gopls`, `CGO_ENABLED=0` Linux `go vet`, `go build ./...`, and `TestFunctionalGates` with `RequireNoFlap` and `fp-go-loom` v0.10.0) passed again, and the corpus check gave the same summary with byte-identical CSV files.

Checker blind spots:

- `pipelinecheck` and `loomcheck` examine syntax only. A pass is not proof of correctness.
- The manual nested-call regular expression gives false positives on function literals in `IOE.TryCatchError` and `IO.FromImpure`.
- The gates do not scan test files.
- Excelize does not show all XML decoder errors through `Rows.Error`.

We did not run repo-wide `golangci-lint run ./...` or `go test ./...` on the scratch copy. Run the two commands as implementer gates:

```sh
golangci-lint run ./...
go test ./...
```
