package filter

import (
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
)

// operatorMap maps operator symbols to Operator enum values
var operatorMap = map[string]Operator{
	// String operators
	"===": Equals,
	"!==": NotEquals,
	"=~":  Contains,
	"!~":  NotContains,
	"=@=": Contains, // Alternative syntax for contains

	// Numeric operators
	"#=":  NumEquals,
	"#>":  GreaterThan,
	"#<":  LessThan,
	"#>=": GreaterOrEqual,
	"#<=": LessOrEqual,

	// Date operators
	"$=":  DateEquals,
	"$>":  DateGreater,
	"$<":  DateLess,
	"$>=": DateGreaterOrEqual,
	"$<=": DateLessOrEqual,

	// Array operators
	"@=~": ArrayContains,
	"@!~": ArrayNotContains,
	"@==": ArrayEquals,
	"@!=": ArrayNotEquals,
}

// ParseFilter parses a filter string into a Expression using Either composition
func ParseFilter(filterStr string) E.Either[error, Expression] {
	if filterStr == "" {
		return E.Right[error](Expression(AlwaysTrueFilter{}))
	}

	return F.Pipe2(
		filterStr,
		Tokenize,
		parseTokens,
	)
}

// parseTokens converts tokens into Expression using Either composition
func parseTokens(tokens []Token) E.Either[error, Expression] {
	// Early return for empty tokens
	if len(tokens) == 0 {
		return E.Right[error](Expression(AlwaysTrueFilter{}))
	}

	// Try parsing as OR expression first (lowest precedence)
	orGroups := splitByOperator(tokens, TokenOr)
	if len(orGroups) > 1 {
		return buildOrTree(orGroups)
	}

	// Try parsing as AND expression (higher precedence)
	andGroups := splitByOperator(tokens, TokenAnd)
	if len(andGroups) > 1 {
		return buildAndTree(andGroups)
	}

	// Single predicate
	return parsePredicate(tokens)
}
