package stockcenter

import (
	"fmt"
	"strings"
	"time"

	A "github.com/IBM/fp-go/array"
	Eq "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	Pred "github.com/IBM/fp-go/predicate"
	S "github.com/IBM/fp-go/string"
)

const (
	// GoldenBraidFieldCount is the expected number of fields in a GoldenBraid CSV record
	GoldenBraidFieldCount = 7
)

// GoldenBraidPlasmid represents a plasmid from the GoldenBraid CSV file
type GoldenBraidPlasmid struct {
	Name         string             // Column 0: "Plasmid  Name"
	Summary      string             // Column 5: "Description"
	Genes        O.Option[[]string] // Column 3: parsed comma-separated genes
	Publications O.Option[[]string] // Column 6: "PMID" (stored as single-element slice)
	User         string             // From CLI --user-email
	PlasmidType  string             // From CLI --plasmid-cvterm
	CreatedOn    time.Time          // Default timestamp
	UpdatedOn    time.Time          // Default timestamp
}

var (
	splitWithComma = F.Bind2nd(strings.Split, ",")

	hasPrefix = F.Bind1st(strings.HasPrefix, "p")

	HasValidUser = F.Flow2(
		func(p *GoldenBraidPlasmid) string { return p.User },
		S.IsNonEmpty,
	)

	HasValidSummary = F.Flow2(
		func(p *GoldenBraidPlasmid) string { return p.Summary },
		S.IsNonEmpty,
	)

	// HasValidName checks if plasmid name starts with 'p' and is non-empty
	HasValidName = F.Flow2(
		func(p *GoldenBraidPlasmid) string { return p.Name },
		Pred.MonoidAll[string]().Concat(
			S.IsNonEmpty,
			hasPrefix,
		),
	)

	// Eq.Of creates an Eq instance for comparable types.
	// We then create a predicate that checks for equality with GoldenBraidFieldCount.
	HasValidRecordLength = F.Flow2(
		A.Size[string],
		Eq.Equals(Eq.FromStrictEquals[int]())(GoldenBraidFieldCount),
	)
)

// parseCommaSeparatedField parses a comma-separated string into Option[[]string]
// Empty string → None, populated → Some([]string)
func parseCommaSeparatedField(field string) O.Option[[]string] {
	return F.Pipe3(
		field,
		splitWithComma,
		A.Map(F.Flow2(
			strings.TrimSpace,
			O.FromPredicate(S.IsNonEmpty),
		)),
		O.SequenceArray[string],
	)
}

// parsePublicationField parses a single PMID into Option[[]string]
// Empty → None, populated → Some([]string{pmid})
func parsePublicationField(field string) O.Option[[]string] {
	return F.Pipe3(
		field,
		strings.TrimSpace,
		O.FromPredicate(S.IsNonEmpty),
		O.Map(A.Of[string]),
	)
}

// RecordLengthError creates error for invalid record length
func RecordLengthError(r []string) error {
	return fmt.Errorf(
		"invalid CSV record: expected %d fields, got %d",
		GoldenBraidFieldCount,
		len(r),
	)
}

// BuildPlasmid constructs GoldenBraidPlasmid immutably from CSV record (curried)
func BuildPlasmid(
	userEmail string,
	plasmidCVTerm string,
) func([]string) *GoldenBraidPlasmid {
	return func(r []string) *GoldenBraidPlasmid {
		now := time.Now()
		return &GoldenBraidPlasmid{
			Name:         strings.TrimSpace(r[0]),
			Summary:      strings.TrimSpace(r[5]),
			Genes:        parseCommaSeparatedField(r[3]),
			Publications: parsePublicationField(r[6]),
			User:         userEmail,
			PlasmidType:  plasmidCVTerm,
			CreatedOn:    now,
			UpdatedOn:    now,
		}
	}
}

// Validation predicates

// NameError creates error for invalid name
func NameError(p *GoldenBraidPlasmid) error {
	return fmt.Errorf("invalid plasmid name '%s': must start with 'p'", p.Name)
}

// SummaryError creates error for empty summary
func SummaryError(p *GoldenBraidPlasmid) error {
	return fmt.Errorf("plasmid '%s' has empty summary", p.Name)
}

// UserError creates error for missing user
func UserError(p *GoldenBraidPlasmid) error {
	return fmt.Errorf("plasmid '%s' has no user email", p.Name)
}
