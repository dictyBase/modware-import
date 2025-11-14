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
	Name         string              // Column 0: "Plasmid  Name"
	Summary      string              // Column 5: "Description"
	Genes        O.Option[[]string]  // Column 3: parsed comma-separated genes
	Publications O.Option[[]string]  // Column 6: "PMID" (stored as single-element slice)
	User         string              // From CLI --user-email
	PlasmidType  string              // From CLI --plasmid-cvterm
	CreatedOn    time.Time           // Default timestamp
	UpdatedOn    time.Time           // Default timestamp
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

// setPlasmidName is a curried setter for the Name field
var setPlasmidName = F.Curry2(
	func(record []string, p *GoldenBraidPlasmid) *GoldenBraidPlasmid {
		p.Name = strings.TrimSpace(record[0])
		return p
	},
)

// setPlasmidSummary is a curried setter for the Summary field
var setPlasmidSummary = F.Curry2(
	func(record []string, p *GoldenBraidPlasmid) *GoldenBraidPlasmid {
		p.Summary = strings.TrimSpace(record[5])
		return p
	},
)

// setPlasmidGenes is a curried setter for the Genes field
var setPlasmidGenes = F.Curry2(
	func(record []string, p *GoldenBraidPlasmid) E.Either[error, *GoldenBraidPlasmid] {
		p.Genes = parseCommaSeparatedField(record[3])
		return E.Of[error](p)
	},
)

// setPlasmidPublications is a curried setter for the Publications field
var setPlasmidPublications = F.Curry2(
	func(record []string, p *GoldenBraidPlasmid) E.Either[error, *GoldenBraidPlasmid] {
		p.Publications = parsePublicationField(record[6])
		return E.Of[error](p)
	},
)

// setMetadata is a curried setter for user, plasmid type, and timestamps
var setMetadata = F.Curry3(
	func(
		userEmail string,
		plasmidCVTerm string,
		p *GoldenBraidPlasmid,
	) *GoldenBraidPlasmid {
		now := time.Now()
		p.User = userEmail
		p.PlasmidType = plasmidCVTerm
		p.CreatedOn = now
		p.UpdatedOn = now
		return p
	},
)

// ParseRecord parses a single CSV record into a GoldenBraidPlasmid
func ParseRecord(
	record []string,
	userEmail string,
	plasmidCVTerm string,
) E.Either[error, *GoldenBraidPlasmid] {
	// Validate record length
	if len(record) != 7 {
		return E.Left[*GoldenBraidPlasmid](
			fmt.Errorf("invalid CSV record: expected 7 fields, got %d", len(record)),
		)
	}

	return F.Pipe5(
		E.Of[error](&GoldenBraidPlasmid{}),
		E.Map[error](setPlasmidName(record)),
		E.Map[error](setPlasmidSummary(record)),
		E.Chain[error](setPlasmidGenes(record)),
		E.Chain[error](setPlasmidPublications(record)),
		E.Map[error](setMetadata(userEmail)(plasmidCVTerm)),
	)
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
func ValidatePlasmid(p *GoldenBraidPlasmid) E.Either[error, *GoldenBraidPlasmid] {
	return F.Pipe3(
		E.Of[error](p),
		E.Chain(E.FromPredicate(hasValidName, nameError)),
		E.Chain(E.FromPredicate(hasValidSummary, summaryError)),
		E.Chain(E.FromPredicate(hasValidUser, userError)),
	)
}
