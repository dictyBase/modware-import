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

	parseAddress   = E.Eitherize1(mail.ParseAddress)
	unnamedAddress = F.Flow2(addressName, Str.IsEmpty)
	// isBareEmail accepts one address with no display name.
	isBareEmail = F.Flow2(bareAddress, O.IsSome[string])
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
func bareAddress(s string) O.Option[string] {
	return F.Pipe5(
		s,
		parseAddress,
		E.ToOption[error, *mail.Address],
		O.Filter(unnamedAddress),
		O.Map(addressValue),
		O.Filter(F.Pipe1(s, predord.StrEq)),
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
