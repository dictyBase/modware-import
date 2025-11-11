package filter

import (
	"fmt"

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

// ParseFilter parses a filter string into a FilterExpression using Either composition
func ParseFilter(filterStr string) E.Either[error, FilterExpression] {
	if filterStr == "" {
		return E.Right[error](FilterExpression(AlwaysTrueFilter{}))
	}

	return F.Pipe2(
		filterStr,
		Tokenize,
		parseTokens,
	)
}

// parseTokens converts tokens into FilterExpression
func parseTokens(tokens []Token) E.Either[error, FilterExpression] {
	if len(tokens) == 0 {
		return E.Right[error](FilterExpression(AlwaysTrueFilter{}))
	}

	// Split by OR operators (lowest precedence)
	orGroups := splitByOperator(tokens, TokenOr)
	if len(orGroups) > 1 {
		return buildOrTree(orGroups)
	}

	// Split by AND operators (higher precedence)
	andGroups := splitByOperator(tokens, TokenAnd)
	if len(andGroups) > 1 {
		return buildAndTree(andGroups)
	}

	// Single predicate
	return parsePredicate(tokens)
}

// splitByOperator splits tokens by a specific operator type
func splitByOperator(tokens []Token, opType TokenType) [][]Token {
	groups := [][]Token{}
	current := []Token{}

	for _, token := range tokens {
		if token.Type == opType {
			if len(current) > 0 {
				groups = append(groups, current)
				current = []Token{}
			}
		} else {
			current = append(current, token)
		}
	}

	if len(current) > 0 {
		groups = append(groups, current)
	}

	return groups
}

// buildOrTree builds an OR expression tree from token groups
func buildOrTree(groups [][]Token) E.Either[error, FilterExpression] {
	if len(groups) == 0 {
		return E.Left[FilterExpression](fmt.Errorf("empty OR groups"))
	}

	if len(groups) == 1 {
		return parseTokens(groups[0])
	}

	return F.Pipe1(
		parseTokens(groups[0]),
		E.Chain(func(left FilterExpression) E.Either[error, FilterExpression] {
			return F.Pipe1(
				buildOrTree(groups[1:]),
				E.Map[error](func(right FilterExpression) FilterExpression {
					return OrExpression{Left: left, Right: right}
				}),
			)
		}),
	)
}

// buildAndTree builds an AND expression tree from token groups
func buildAndTree(groups [][]Token) E.Either[error, FilterExpression] {
	if len(groups) == 0 {
		return E.Left[FilterExpression](fmt.Errorf("empty AND groups"))
	}

	if len(groups) == 1 {
		return parseTokens(groups[0])
	}

	return F.Pipe1(
		parseTokens(groups[0]),
		E.Chain(func(left FilterExpression) E.Either[error, FilterExpression] {
			return F.Pipe1(
				buildAndTree(groups[1:]),
				E.Map[error](func(right FilterExpression) FilterExpression {
					return AndExpression{Left: left, Right: right}
				}),
			)
		}),
	)
}

// Predicates for token validation
// hasLengthThree checks if token slice has exactly 3 elements
func hasLengthThree(t []Token) bool {
	return len(t) == 3
}

// lengthError creates an error for invalid token length
func lengthError(t []Token) error {
	return fmt.Errorf(
		"invalid predicate: expected 3 tokens, got %d",
		len(t),
	)
}

// hasFieldTokenType checks if first token is a Field type
func hasFieldTokenType(t []Token) bool {
	return t[0].Type == TokenField
}

// fieldTokenError creates an error for invalid field token
func fieldTokenError(t []Token) error {
	return fmt.Errorf(
		"invalid predicate: first token must be a Field, got %v",
		t[0].Type,
	)
}

// hasOperatorTokenType checks if second token is an Operator type
func hasOperatorTokenType(t []Token) bool {
	return t[1].Type == TokenOperator
}

// operatorTokenError creates an error for invalid operator token
func operatorTokenError(t []Token) error {
	return fmt.Errorf(
		"invalid predicate: second token must be an Operator, got %v",
		t[1].Type,
	)
}

// hasValueTokenType checks if third token is a Value type
func hasValueTokenType(t []Token) bool {
	return t[2].Type == TokenValue
}

// valueTokenError creates an error for invalid value token
func valueTokenError(t []Token) error {
	return fmt.Errorf(
		"invalid predicate: third token must be a Value, got %v",
		t[2].Type,
	)
}

// validateTokenStructure executes the full validation pipeline.
// It takes tokens and returns the result of the composed validation checks.
func validateTokenStructure(tokens []Token) E.Either[error, []Token] {
	return F.Pipe4(
		tokens,
		E.FromPredicate(hasLengthThree, lengthError),
		E.Chain(
			E.FromPredicate(hasFieldTokenType, fieldTokenError),
		),
		E.Chain(
			E.FromPredicate(
				hasOperatorTokenType,
				operatorTokenError,
			),
		),
		E.Chain(
			E.FromPredicate(hasValueTokenType, valueTokenError),
		),
	)
}

// operatorExists checks if operator symbol exists in operatorMap
func operatorExists(operatorSymbol string) bool {
	_, ok := operatorMap[operatorSymbol]
	return ok
}

// validateOperatorExists validates operator symbol using predicate
var validateOperatorExists = E.FromPredicate(
	operatorExists,
	func(symbol string) error {
		return fmt.Errorf("unknown operator: %s", symbol)
	},
