package filter

import (
	"strings"

	A "github.com/IBM/fp-go/array"
	EQ "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	ORD "github.com/IBM/fp-go/ord"
)

// Ord and Eq instances for functional comparison
var (
	intOrd   = ORD.FromStrictCompare[int]()
	stringEq = EQ.FromStrictEquals[string]()
	leqOrd   = ORD.Leq(intOrd)
)

// operatorSymbols defines all valid operator symbols in order of precedence (longest first)
var operatorSymbols = []string{
	// Three character operators
	"===", "!==", "@=~", "@!~", "@==", "@!=", "=@=",
	// Two character operators
	"$>=", "$<=", "$>", "$<", "$=",
	"#>=", "#<=", "#>", "#<", "#=",
	"=~", "!~",
}

// operatorContext holds the context needed to check if an operator matches
type operatorContext struct {
	filter   string
	position int
	opSymbol string
}

// Tokenize converts a filter string into a sequence of tokens
func Tokenize(filter string) []Token {
	if filter == "" {
		return []Token{}
	}

	tokens := make([]Token, 0)
	position := 0

	for position < len(filter) {
		position = processNextToken(filter, position, &tokens)
	}

	return tokens
}

// processNextToken processes the next token in the filter string and returns the new position
func processNextToken(filter string, position int, tokens *[]Token) int {
	// Check for logical operators first
	if newPos, handled := tryLogicalOperator(filter, position, tokens); handled {
		return newPos
	}

	// Try to match a comparison operator
	if newPos, handled := tryComparisonOperator(filter, position, tokens); handled {
		return newPos
	}

	// Extract field or value
	return extractFieldOrValue(filter, position, tokens)
}

// tryLogicalOperator checks for AND (;) or OR (,) operators
func tryLogicalOperator(
	filter string,
	position int,
	tokens *[]Token,
) (int, bool) {
	switch filter[position] {
	case ';':
		*tokens = append(*tokens, Token{Type: TokenAnd, Value: ";"})
		return position + 1, true
	case ',':
		*tokens = append(*tokens, Token{Type: TokenOr, Value: ","})
		return position + 1, true
	}
	return position, false
}

// tryComparisonOperator attempts to match a comparison operator at the current position
func tryComparisonOperator(
	filter string,
	position int,
	tokens *[]Token,
) (int, bool) {
	for _, opSymbol := range operatorSymbols {
		if position+len(opSymbol) <= len(filter) &&
			filter[position:position+len(opSymbol)] == opSymbol {
			*tokens = append(
				*tokens,
				Token{Type: TokenOperator, Value: opSymbol},
			)
			return position + len(opSymbol), true
		}
	}
	return position, false
}

// extractFieldOrValue extracts either a field name or a value based on context
func extractFieldOrValue(filter string, position int, tokens *[]Token) int {
	startPos := position
	isField := isFieldToken(*tokens)

	if isField {
		return extractField(filter, startPos, position, tokens)
	}
	return extractValue(filter, startPos, position, tokens)
}

// isFieldToken determines if the next token should be a field based on the previous token
func isFieldToken(tokens []Token) bool {
	if len(tokens) == 0 {
		return true
	}
	return tokens[len(tokens)-1].Type != TokenOperator
}

// extractField extracts a field name up to the next operator
func extractField(filter string, startPos, position int, tokens *[]Token) int {
	for position < len(filter) {
		if isAtOperator(filter, position) {
			break
		}
		position++
	}

	fieldValue := strings.TrimSpace(filter[startPos:position])
	if fieldValue != "" {
		*tokens = append(*tokens, Token{Type: TokenField, Value: fieldValue})
	}
	return position
}

// extractValue extracts a value up to the next logical operator (AND/OR) or end of string
func extractValue(filter string, startPos, position int, tokens *[]Token) int {
	for position < len(filter) {
		if filter[position] == ';' || filter[position] == ',' {
			break
		}
		position++
	}

	valueText := strings.TrimSpace(filter[startPos:position])
	if valueText != "" {
		*tokens = append(*tokens, Token{Type: TokenValue, Value: valueText})
	}
	return position
}

// operatorMatches checks if operator matches at position using a unified pipeline
func operatorMatches(ctx operatorContext) bool {
	endPos := ctx.position + len(ctx.opSymbol)

	return F.Pipe4(
		endPos,
		// 1. Check boundary: endPos <= len(filter)
		O.FromPredicate(leqOrd(len(ctx.filter))),
		// 2. Extract substring if valid
		O.Map(func(_ int) string {
			return ctx.filter[ctx.position:endPos]
		}),
		// 3. Check if substring equals operator
		O.Map(func(substring string) bool {
			return stringEq.Equals(substring, ctx.opSymbol)
		}),
		// 4. Return false if any step failed, otherwise return the result
		O.GetOrElse(F.Constant(false)),
	)
}

// toOperatorContext creates an operatorContext from a symbol
func toOperatorContext(
	filter string,
	position int,
) func(string) operatorContext {
	return func(opSymbol string) operatorContext {
		return operatorContext{
			filter:   filter,
			position: position,
			opSymbol: opSymbol,
		}
	}
}

// isAtOperator checks if the current position is at the start of an operator
// using point-free functional composition
func isAtOperator(filter string, position int) bool {
	addContext := toOperatorContext(filter, position)
	return F.Pipe2(
		operatorSymbols,
		A.Map(addContext),
		A.Any(operatorMatches),
	)
}
