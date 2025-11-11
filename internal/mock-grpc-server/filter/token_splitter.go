package filter

import (
	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

// splitAccumulator is the accumulator state for the reduce operation.
type splitAccumulator struct {
	groups    [][]Token
	lastIndex int
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

// buildTokenGroupImpl creates a new accumulator state with the next group of tokens.
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

// buildTokenGroup is the partially applicable version using F.Bind1of3.
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
