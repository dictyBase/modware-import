package stockcenter

import (
	"fmt"
	"strings"
	"time"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	Pred "github.com/IBM/fp-go/predicate"
	S "github.com/IBM/fp-go/string"
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

// parseCommaSeparatedField parses a comma-separated string into Option[[]string]
// Empty string → None, populated → Some([]string)
func parseCommaSeparatedField(field string) O.Option[[]string] {
	trimmed := strings.TrimSpace(field)

	// Empty string → None
	if S.IsEmpty(trimmed) {
		return O.None[[]string]()
	}

	// Split, trim, filter non-empty using predicate
	return F.Pipe3(
		strings.Split(trimmed, ","),
		A.Map(strings.TrimSpace),
		A.Filter(Pred.Not(S.IsEmpty)),
		O.Some[[]string],
	)
}

// parsePublicationField parses a single PMID into Option[[]string]
// Empty → None, populated → Some([]string{pmid})
func parsePublicationField(field string) O.Option[[]string] {
	trimmed := strings.TrimSpace(field)

	if S.IsEmpty(trimmed) {
		return O.None[[]string]()
	}

	return O.Some([]string{trimmed})
}

// HasValidRecordLength checks if CSV record has exactly 7 fields
func HasValidRecordLength(r []string) bool {
	return len(r) == 7
}

// RecordLengthError creates error for invalid record length
func RecordLengthError(r []string) error {
	return fmt.Errorf("invalid CSV record: expected 7 fields, got %d", len(r))
}

// BuildPlasmid constructs GoldenBraidPlasmid immutably from CSV record (curried)
func BuildPlasmid(userEmail string, plasmidCVTerm string) func([]string) *GoldenBraidPlasmid {
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

// hasValidName checks if plasmid name starts with 'p' and is non-empty
var hasValidName = func(p *GoldenBraidPlasmid) bool {
	return !S.IsEmpty(p.Name) && strings.HasPrefix(p.Name, "p")
}

// nameError creates error for invalid name
var nameError = func(p *GoldenBraidPlasmid) error {
	return fmt.Errorf("invalid plasmid name '%s': must start with 'p'", p.Name)
}

// hasValidSummary checks if summary is non-empty
var hasValidSummary = func(p *GoldenBraidPlasmid) bool {
	return !S.IsEmpty(p.Summary)
}

// summaryError creates error for empty summary
var summaryError = func(p *GoldenBraidPlasmid) error {
	return fmt.Errorf("plasmid '%s' has empty summary", p.Name)
}

// hasValidUser checks if user email is non-empty
var hasValidUser = func(p *GoldenBraidPlasmid) bool {
	return !S.IsEmpty(p.User)
}

// userError creates error for missing user
var userError = func(p *GoldenBraidPlasmid) error {
	return fmt.Errorf("plasmid '%s' has no user email", p.Name)
}

// ValidatePlasmid validates a GoldenBraidPlasmid using predicate chain
func ValidatePlasmid(
	p *GoldenBraidPlasmid,
) E.Either[error, *GoldenBraidPlasmid] {
	return F.Pipe3(
		E.Of[error](p),
		E.Chain(E.FromPredicate(hasValidName, nameError)),
		E.Chain(E.FromPredicate(hasValidSummary, summaryError)),
		E.Chain(E.FromPredicate(hasValidUser, userError)),
	)
}
