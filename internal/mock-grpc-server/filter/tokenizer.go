package filter

import (
	"strings"

	A "github.com/IBM/fp-go/array"
	EQ "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	ORD "github.com/IBM/fp-go/ord"
)

// Ord and Eq instances for functional comparison
var (
	intOrd   = ORD.FromStrictCompare[int]()
	stringEq = EQ.FromEquals(func(a, b string) bool { return a == b })
)

// operatorSymbols defines all valid operator symbols in order of precedence (longest first)
var operatorSymbols = []string{
	// Three character operators
	"===", "!==", "@=~", "@!~", "@==", "@!=",
	// Two character operators
	"$>=", "$<=", "$>", "$<", "$=",
	"#>=", "#<=", "#>", "#<", "#=",
	"=~", "!~",
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
func tryLogicalOperator(filter string, position int, tokens *[]Token) (int, bool) {
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
func tryComparisonOperator(filter string, position int, tokens *[]Token) (int, bool) {
	for _, opSymbol := range operatorSymbols {
		if position+len(opSymbol) <= len(filter) && filter[position:position+len(opSymbol)] == opSymbol {
			*tokens = append(*tokens, Token{Type: TokenOperator, Value: opSymbol})
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

// hasEnoughSpaceImpl checks if there's enough space for an operator
func hasEnoughSpaceImpl(filterLen, position, opLen int) bool {
	return ORD.Leq(intOrd)(filterLen)(position + opLen)
}

// hasEnoughSpace is the curried version using F.Bind1of3
var hasEnoughSpace = F.Bind1of3(hasEnoughSpaceImpl)

// extractSubstringImpl safely extracts substring from filter
func extractSubstringImpl(filter string, position, opLen int) string {
	endPos := position + opLen
	if endPos > len(filter) {
		return ""
	}
	return filter[position:endPos]
}

// extractSubstring is the curried version
var extractSubstring = F.Bind1of3(extractSubstringImpl)

// substringMatchesImpl checks if substring equals operator symbol
func substringMatchesImpl(filter string, position int, opSymbol string) bool {
	substring := extractSubstring(filter)(position, len(opSymbol))
	return stringEq.Equals(substring, opSymbol)
}

// substringMatches is the curried version
var substringMatches = F.Bind1of3(substringMatchesImpl)

// operatorMatchesAtPositionImpl combines boundary and match checks
func operatorMatchesAtPositionImpl(filter string, position int, opSymbol string) bool {
	// Boundary check first
	if !hasEnoughSpace(len(filter))(position, len(opSymbol)) {
		return false
	}
	// Then string match
	return substringMatches(filter)(position, opSymbol)
}

// isAtOperator checks if the current position is at the start of an operator
// using functional predicate composition with Ord and Eq
func isAtOperator(filter string, position int) bool {
	return A.Any(func(opSymbol string) bool {
		return operatorMatchesAtPositionImpl(filter, position, opSymbol)
	})(operatorSymbols)
}
