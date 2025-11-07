package filter

import (
	"strings"
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
		// Check for AND operator (;)
		if filter[position] == ';' {
			tokens = append(tokens, Token{Type: TokenAnd, Value: ";"})
			position++
			continue
		}

		// Check for OR operator (,)
		if filter[position] == ',' {
			tokens = append(tokens, Token{Type: TokenOr, Value: ","})
			position++
			continue
		}

		// Try to match an operator
		operatorMatched := false
		for _, opSymbol := range operatorSymbols {
			if position+len(opSymbol) <= len(filter) && filter[position:position+len(opSymbol)] == opSymbol {
				tokens = append(tokens, Token{Type: TokenOperator, Value: opSymbol})
				position += len(opSymbol)
				operatorMatched = true
				break
			}
		}

		if operatorMatched {
			continue
		}

		// Extract field or value
		// Field: everything before an operator
		// Value: everything after an operator until AND/OR or end
		startPos := position

		// Determine if we're looking at a field (before operator) or value (after operator)
		// Look ahead to find what comes next
		isField := true
		if len(tokens) > 0 && tokens[len(tokens)-1].Type == TokenOperator {
			isField = false
		}

		if isField {
			// Extract field name (up to operator)
			for position < len(filter) {
				// Check if we're at an operator
				atOperator := false
				for _, opSymbol := range operatorSymbols {
					if position+len(opSymbol) <= len(filter) && filter[position:position+len(opSymbol)] == opSymbol {
						atOperator = true
						break
					}
				}
				if atOperator {
					break
				}
				position++
			}
			fieldValue := strings.TrimSpace(filter[startPos:position])
			if fieldValue != "" {
				tokens = append(tokens, Token{Type: TokenField, Value: fieldValue})
			}
		} else {
			// Extract value (up to AND/OR or end)
			for position < len(filter) {
				if filter[position] == ';' || filter[position] == ',' {
					break
				}
				position++
			}
			valueText := strings.TrimSpace(filter[startPos:position])
			if valueText != "" {
				tokens = append(tokens, Token{Type: TokenValue, Value: valueText})
			}
		}
	}

	return tokens
}
