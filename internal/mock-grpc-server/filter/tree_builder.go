package filter

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
)

// combineAsOr creates an OrExpression from left and right expressions
func combineAsOr(
	left Expression,
) func(Expression) Expression {
	return func(right Expression) Expression {
		return OrExpression{Left: left, Right: right}
	}
}

// buildOrTreeRight builds right side of OR tree and combines with left
func buildOrTreeRight(
	groups [][]Token,
) func(Expression) E.Either[error, Expression] {
	return func(left Expression) E.Either[error, Expression] {
		return F.Pipe1(
			buildOrTree(groups),
			E.Map[error](combineAsOr(left)),
		)
	}
}

// buildOrTree builds an OR expression tree from token groups using Either composition
func buildOrTree(groups [][]Token) E.Either[error, Expression] {
	if len(groups) == 0 {
		return E.Left[Expression](fmt.Errorf("empty OR groups"))
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
	left Expression,
) func(Expression) Expression {
	return func(right Expression) Expression {
		return AndExpression{Left: left, Right: right}
	}
}

// buildAndTreeRight builds right side of AND tree and combines with left
func buildAndTreeRight(
	groups [][]Token,
) func(Expression) E.Either[error, Expression] {
	return func(left Expression) E.Either[error, Expression] {
		return F.Pipe1(
			buildAndTree(groups),
			E.Map[error](combineAsAnd(left)),
		)
	}
}

// buildAndTree builds an AND expression tree from token groups using Either composition
func buildAndTree(groups [][]Token) E.Either[error, Expression] {
	if len(groups) == 0 {
		return E.Left[Expression](fmt.Errorf("empty AND groups"))
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
