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

// parsePredicate parses a single field-operator-value predicate
func parsePredicate(tokens []Token) E.Either[error, FilterExpression] {
	if len(tokens) != 3 {
		return E.Left[FilterExpression](
			fmt.Errorf("invalid predicate: expected 3 tokens, got %d", len(tokens)),
		)
	}

	if tokens[0].Type != TokenField {
		return E.Left[FilterExpression](
			fmt.Errorf("invalid predicate: first token must be field, got %v", tokens[0].Type),
		)
	}

	if tokens[1].Type != TokenOperator {
		return E.Left[FilterExpression](
			fmt.Errorf("invalid predicate: second token must be operator, got %v", tokens[1].Type),
		)
	}

	if tokens[2].Type != TokenValue {
		return E.Left[FilterExpression](
			fmt.Errorf("invalid predicate: third token must be value, got %v", tokens[2].Type),
		)
	}

	operator, ok := operatorMap[tokens[1].Value]
	if !ok {
		return E.Left[FilterExpression](
			fmt.Errorf("unknown operator: %s", tokens[1].Value),
		)
	}

	return E.Right[error](FilterExpression(Predicate{
		Field:    tokens[0].Value,
		Operator: operator,
		Value:    tokens[2].Value,
	}))
}
