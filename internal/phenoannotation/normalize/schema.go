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
