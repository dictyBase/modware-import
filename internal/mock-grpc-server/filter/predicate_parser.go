package filter

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
)

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

// unknownOperatorError creates an error for unknown operator symbol
func unknownOperatorError(symbol string) error {
	return fmt.Errorf("unknown operator: %s", symbol)
}

// validateOperatorExists validates operator symbol using predicate
var validateOperatorExists = E.FromPredicate(
	operatorExists,
	unknownOperatorError,
)

// lookupOperator retrieves the Operator enum value from operatorMap
func lookupOperator(symbol string) Operator {
	return operatorMap[symbol]
}

// parseOperator extracts and validates operator from tokens using Either composition
func parseOperator(tokens []Token) E.Either[error, Operator] {
	return F.Pipe1(
		validateOperatorExists(tokens[1].Value),
		E.Map[error](lookupOperator),
	)
}

// constructPredicate creates a FilterExpression from operator and tokens
func constructPredicate(op Operator, tokens []Token) FilterExpression {
	return FilterExpression(Predicate{
		Field:    tokens[0].Value,
		Operator: op,
		Value:    tokens[2].Value,
	})
}

// buildPredicate constructs FilterExpression from validated tokens and operator
var buildPredicate = F.Curry2(constructPredicate)

// applyOperatorToPredicate applies parsed operator to build the final predicate
func applyOperatorToPredicate(
	validTokens []Token,
) func(Operator) FilterExpression {
	return func(op Operator) FilterExpression {
		return buildPredicate(op)(validTokens)
	}
}

// buildPredicateFromTokens chains operator parsing and predicate construction
func buildPredicateFromTokens(
	validTokens []Token,
) E.Either[error, FilterExpression] {
	return F.Pipe1(
		parseOperator(validTokens),
		E.Map[error](applyOperatorToPredicate(validTokens)),
	)
}

// parsePredicate parses a single field-operator-value predicate using Either composition
func parsePredicate(tokens []Token) E.Either[error, FilterExpression] {
	return F.Pipe1(
		validateTokenStructure(tokens),
		E.Chain(buildPredicateFromTokens),
	)
}
