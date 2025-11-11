package filter

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

// The accumulator state for our reduce operation.
type splitAccumulator struct {
	groups    [][]Token
	lastIndex int
}

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

// parseTokens converts tokens into FilterExpression using Either composition
func parseTokens(tokens []Token) E.Either[error, FilterExpression] {
	// Early return for empty tokens
	if len(tokens) == 0 {
		return E.Right[error](FilterExpression(AlwaysTrueFilter{}))
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

// combineAsOr creates an OrExpression from left and right expressions
func combineAsOr(
	left FilterExpression,
) func(FilterExpression) FilterExpression {
	return func(right FilterExpression) FilterExpression {
		return OrExpression{Left: left, Right: right}
	}
}

// buildOrTreeRight builds right side of OR tree and combines with left
func buildOrTreeRight(
	groups [][]Token,
) func(FilterExpression) E.Either[error, FilterExpression] {
	return func(left FilterExpression) E.Either[error, FilterExpression] {
		return F.Pipe1(
			buildOrTree(groups),
			E.Map[error](combineAsOr(left)),
		)
	}
}

// buildOrTree builds an OR expression tree from token groups using Either composition
func buildOrTree(groups [][]Token) E.Either[error, FilterExpression] {
	if len(groups) == 0 {
		return E.Left[FilterExpression](fmt.Errorf("empty OR groups"))
	}

	if len(groups) == 1 {
		return parseTokens(groups[0])
	}

	// Parse left side, then recursively build right side, combining into OrExpression
	return F.Pipe1(
		parseTokens(groups[0]),
		E.Chain(buildOrTreeRight(groups[1:])),
	)
}

// combineAsAnd creates an AndExpression from left and right expressions
func combineAsAnd(
	left FilterExpression,
) func(FilterExpression) FilterExpression {
	return func(right FilterExpression) FilterExpression {
		return AndExpression{Left: left, Right: right}
	}
}

// buildAndTreeRight builds right side of AND tree and combines with left
func buildAndTreeRight(
	groups [][]Token,
) func(FilterExpression) E.Either[error, FilterExpression] {
	return func(left FilterExpression) E.Either[error, FilterExpression] {
		return F.Pipe1(
			buildAndTree(groups),
			E.Map[error](combineAsAnd(left)),
		)
	}
}

// buildAndTree builds an AND expression tree from token groups using Either composition
func buildAndTree(groups [][]Token) E.Either[error, FilterExpression] {
	if len(groups) == 0 {
		return E.Left[FilterExpression](fmt.Errorf("empty AND groups"))
	}

	if len(groups) == 1 {
		return parseTokens(groups[0])
	}

	// Parse left side, then recursively build right side, combining into AndExpression
	return F.Pipe1(
		parseTokens(groups[0]),
		E.Chain(buildAndTreeRight(groups[1:])),
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

// checkIfSeparatorImpl checks if a token at index i matches the separator type.
// Returns Some(index) if it's a separator, None otherwise.
func checkIfSeparatorImpl(
	opType TokenType,
	i int,
	token Token,
) O.Option[int] {
	// Use FromPredicate idiomatically: test token, return index if match
	return F.Pipe1(
		i,
		O.FromPredicate(func(idx int) bool {
			return token.Type == opType
		}),
	)
}

// buildTokenGroup creates a new accumulator state with the next group of tokens.
func buildTokenGroupImpl(
	tokens []Token,
	acc splitAccumulator,
	sepIndex int,
) splitAccumulator {
	newGroup := tokens[acc.lastIndex:sepIndex]
	return splitAccumulator{
		groups:    append(acc.groups, newGroup),
		lastIndex: sepIndex + 1,
	}
}

// checkIfSeparator is the partially applicable version using F.Bind1of3.
var checkIfSeparator = F.Bind1of3(checkIfSeparatorImpl)

var buildTokenGroup = F.Bind1of3(buildTokenGroupImpl)

// extractGroups is a simple accessor function that pulls the final result
// out of our accumulator state.
func extractGroups(state splitAccumulator) [][]Token {
	return state.groups
}

// splitByOperator splits tokens by a specific operator type using functional composition.
// It handles edge cases like consecutive separators by filtering out empty groups.
func splitByOperator(tokens []Token, opType TokenType) [][]Token {
	// Handle the edge case of an empty input slice.
	if len(tokens) == 0 {
		return [][]Token{}
	}

	// Define the initial state for our reduce operation.
	initialState := splitAccumulator{
		groups:    make([][]Token, 0),
		lastIndex: 0,
	}

	// The composed data transformation pipeline.
	// It reads as a clear sequence of steps, using our named functions.
	return F.Pipe5(
		tokens, // 1. Start with the original tokens.

		// 2. Find the indices of all separator tokens.
		// F.Bind1of3 partially applies opType to checkIfSeparator
		A.FilterMapWithIndex(checkIfSeparator(opType)),

		// 3. Append the final index to mark the end of the last group.
		func(indices []int) []int {
			return append(indices, len(tokens))
		},

		// 4. Reduce the list of indices into groups of tokens.
		// F.Bind1of3 partially applies tokens to buildTokenGroup
		A.Reduce(buildTokenGroup(tokens), initialState),

		// 5. Extract the final list of groups from the accumulator state.
		extractGroups,

		// 6. Filter out empty groups to match imperative behavior.
		// This handles edge cases like consecutive separators (e.g., [A, OR, OR, B])
		A.Filter(func(group []Token) bool { return len(group) > 0 }),
	)
}
