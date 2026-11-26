// Package normalize converts create workbooks for phenotype annotation into
// reviewable CSV files. It does not call a stock or an annotation service.
package normalize

import (
	"context"
	"errors"

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
	errSpec      = errors.New("--folder and --out must not be blank")
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
	Folder string
	Sheet  string
	Out    string
}

// Output pairs the path of one CSV file with its rows, header first.
type Output = PR.Pair[string, [][]string]

// State carries the options, the dependencies, and the results of one run.
// Batch grows row by row while the workbooks are read.
// fp-go:Lens
type State struct {
	Context context.Context
	Options Options
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
