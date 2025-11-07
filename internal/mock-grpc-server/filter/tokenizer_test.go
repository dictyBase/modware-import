package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenizeSimpleEquality(t *testing.T) {
	filter := "depositor===Costanza"
	tokens := Tokenize(filter)

	require.Len(t, tokens, 3)
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "depositor", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "===", tokens[1].Value)
	require.Equal(t, TokenValue, tokens[2].Type)
	require.Equal(t, "Costanza", tokens[2].Value)
}

func TestTokenizeWithAnd(t *testing.T) {
	filter := "depositor===Costanza;species===Dictyostelium"
	tokens := Tokenize(filter)

	require.Len(t, tokens, 7)
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "depositor", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "===", tokens[1].Value)
	require.Equal(t, TokenValue, tokens[2].Type)
	require.Equal(t, "Costanza", tokens[2].Value)
	require.Equal(t, TokenAnd, tokens[3].Type)
	require.Equal(t, ";", tokens[3].Value)
	require.Equal(t, TokenField, tokens[4].Type)
	require.Equal(t, "species", tokens[4].Value)
	require.Equal(t, TokenOperator, tokens[5].Type)
	require.Equal(t, "===", tokens[5].Value)
	require.Equal(t, TokenValue, tokens[6].Type)
	require.Equal(t, "Dictyostelium", tokens[6].Value)
}

func TestTokenizeWithOr(t *testing.T) {
	filter := "depositor===Costanza,depositor===Seinfeld"
	tokens := Tokenize(filter)

	require.Len(t, tokens, 7)
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "depositor", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "===", tokens[1].Value)
	require.Equal(t, TokenValue, tokens[2].Type)
	require.Equal(t, "Costanza", tokens[2].Value)
	require.Equal(t, TokenOr, tokens[3].Type)
	require.Equal(t, ",", tokens[3].Value)
	require.Equal(t, TokenField, tokens[4].Type)
	require.Equal(t, "depositor", tokens[4].Value)
}

func TestTokenizeDateOperator(t *testing.T) {
	filter := "created_at$>=2018-12-01"
	tokens := Tokenize(filter)

	require.Len(t, tokens, 3)
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "created_at", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "$>=", tokens[1].Value)
	require.Equal(t, TokenValue, tokens[2].Type)
	require.Equal(t, "2018-12-01", tokens[2].Value)
}

func TestTokenizeNumericOperator(t *testing.T) {
	filter := "quantity#>10"
	tokens := Tokenize(filter)

	require.Len(t, tokens, 3)
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "quantity", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "#>", tokens[1].Value)
	require.Equal(t, TokenValue, tokens[2].Type)
	require.Equal(t, "10", tokens[2].Value)
}

func TestTokenizeComplexFilter(t *testing.T) {
	filter := "depositor===Costanza;created_at$>=2018-12-01,species===Dictyostelium"
	tokens := Tokenize(filter)

	// depositor===Costanza ; created_at$>=2018-12-01 , species===Dictyostelium
	// Should have: field, op, value, AND, field, op, value, OR, field, op, value = 11 tokens
	require.Len(t, tokens, 11)

	// First predicate
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "depositor", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "===", tokens[1].Value)

	// AND operator
	require.Equal(t, TokenAnd, tokens[3].Type)

	// Second predicate
	require.Equal(t, TokenField, tokens[4].Type)
	require.Equal(t, "created_at", tokens[4].Value)

	// OR operator
	require.Equal(t, TokenOr, tokens[7].Type)

	// Third predicate
	require.Equal(t, TokenField, tokens[8].Type)
	require.Equal(t, "species", tokens[8].Value)
}

func TestTokenizeEmptyString(t *testing.T) {
	filter := ""
	tokens := Tokenize(filter)

	require.Len(t, tokens, 0)
}

func TestTokenizeArrayOperators(t *testing.T) {
	filter := "tags@=~important"
	tokens := Tokenize(filter)

	require.Len(t, tokens, 3)
	require.Equal(t, TokenField, tokens[0].Type)
	require.Equal(t, "tags", tokens[0].Value)
	require.Equal(t, TokenOperator, tokens[1].Type)
	require.Equal(t, "@=~", tokens[1].Value)
	require.Equal(t, TokenValue, tokens[2].Type)
	require.Equal(t, "important", tokens[2].Value)
}
