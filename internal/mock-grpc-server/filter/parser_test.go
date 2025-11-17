package filter

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/require"
)

func TestParseSimplePredicate(t *testing.T) {
	filter := "depositor===Costanza"
	result := ParseFilter(filter)

	require.True(t, E.IsRight(result))

	expr := E.GetOrElse(func(error) Expression { return AlwaysTrueFilter{} })(result)
	pred, ok := expr.(Predicate)
	require.True(t, ok, "expected Predicate type")
	require.Equal(t, "depositor", pred.Field)
	require.Equal(t, Equals, pred.Operator)
	require.Equal(t, "Costanza", pred.Value)
}

func TestParseAndExpression(t *testing.T) {
	filter := "depositor===Costanza;species===Dictyostelium"
	result := ParseFilter(filter)

	require.True(t, E.IsRight(result))

	expr := E.GetOrElse(func(error) Expression { return AlwaysTrueFilter{} })(result)
	andExpr, ok := expr.(AndExpression)
	require.True(t, ok, "expected AndExpression type")

	// Check left predicate
	leftPred, ok := andExpr.Left.(Predicate)
	require.True(t, ok)
	require.Equal(t, "depositor", leftPred.Field)
	require.Equal(t, Equals, leftPred.Operator)
	require.Equal(t, "Costanza", leftPred.Value)

	// Check right predicate
	rightPred, ok := andExpr.Right.(Predicate)
	require.True(t, ok)
	require.Equal(t, "species", rightPred.Field)
	require.Equal(t, Equals, rightPred.Operator)
	require.Equal(t, "Dictyostelium", rightPred.Value)
}

func TestParseOrExpression(t *testing.T) {
	filter := "depositor===Costanza,depositor===Seinfeld"
	result := ParseFilter(filter)

	require.True(t, E.IsRight(result))

	expr := E.GetOrElse(func(error) Expression { return AlwaysTrueFilter{} })(result)
	orExpr, ok := expr.(OrExpression)
	require.True(t, ok, "expected OrExpression type")

	// Check left predicate
	leftPred, ok := orExpr.Left.(Predicate)
	require.True(t, ok)
	require.Equal(t, "depositor", leftPred.Field)

	// Check right predicate
	rightPred, ok := orExpr.Right.(Predicate)
	require.True(t, ok)
	require.Equal(t, "depositor", rightPred.Field)
	require.Equal(t, "Seinfeld", rightPred.Value)
}

func TestParseComplexExpression(t *testing.T) {
	// depositor===Costanza AND created_at$>=2018-12-01 OR species===Dictyostelium
	// AND has higher precedence than OR, so this should parse as:
	// (depositor===Costanza AND created_at$>=2018-12-01) OR species===Dictyostelium
	filter := "depositor===Costanza;created_at$>=2018-12-01,species===Dictyostelium"
	result := ParseFilter(filter)

	require.True(t, E.IsRight(result))

	expr := E.GetOrElse(func(error) Expression { return AlwaysTrueFilter{} })(result)
	orExpr, ok := expr.(OrExpression)
	require.True(t, ok, "expected OrExpression at root")

	// Left side should be AND expression
	_, ok = orExpr.Left.(AndExpression)
	require.True(t, ok, "expected AndExpression on left side of OR")

	// Right side should be simple predicate
	rightPred, ok := orExpr.Right.(Predicate)
	require.True(t, ok)
	require.Equal(t, "species", rightPred.Field)
}

func TestParseOperatorMapping(t *testing.T) {
	tests := []struct {
		name             string
		filter           string
		expectedOperator Operator
	}{
		{"equals", "field===value", Equals},
		{"not equals", "field!==value", NotEquals},
		{"contains", "field=~value", Contains},
		{"not contains", "field!~value", NotContains},
		{"numeric equals", "field#=10", NumEquals},
		{"greater than", "field#>10", GreaterThan},
		{"less than", "field#<10", LessThan},
		{"greater or equal", "field#>=10", GreaterOrEqual},
		{"less or equal", "field#<=10", LessOrEqual},
		{"date equals", "field$=2020-01-01", DateEquals},
		{"date greater", "field$>2020-01-01", DateGreater},
		{"date less", "field$<2020-01-01", DateLess},
		{"date greater or equal", "field$>=2020-01-01", DateGreaterOrEqual},
		{"date less or equal", "field$<=2020-01-01", DateLessOrEqual},
		{"array contains", "field@=~value", ArrayContains},
		{"array not contains", "field@!~value", ArrayNotContains},
		{"array equals", "field@==value", ArrayEquals},
		{"array not equals", "field@!=value", ArrayNotEquals},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseFilter(tt.filter)
			require.True(t, E.IsRight(result))

			expr := E.GetOrElse(
				func(error) Expression { return AlwaysTrueFilter{} },
			)(
				result,
			)
			pred, ok := expr.(Predicate)
			require.True(t, ok)
			require.Equal(t, tt.expectedOperator, pred.Operator)
		})
	}
}

func TestParseEmptyFilter(t *testing.T) {
	result := ParseFilter("")
	require.True(t, E.IsRight(result))

	expr := E.GetOrElse(func(error) Expression { return nil })(result)
	_, ok := expr.(AlwaysTrueFilter)
	require.True(t, ok, "empty filter should return AlwaysTrueFilter")
}

func TestParseInvalidFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter string
	}{
		{"missing operator", "depositor"},
		{"missing value", "depositor==="},
		{"missing field", "===value"},
		{"invalid operator", "depositor<>value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseFilter(tt.filter)
			require.True(t, E.IsLeft(result), "invalid filter should return Left")
		})
	}
}

func TestSplitByOperator_EdgeCases(t *testing.T) {
	fieldA := Token{Type: TokenField, Value: "A"}
	fieldB := Token{Type: TokenField, Value: "B"}
	fieldC := Token{Type: TokenField, Value: "C"}
	orToken := Token{Type: TokenOr}

	tests := []struct {
		name     string
		tokens   []Token
		opType   TokenType
		expected [][]Token
	}{
		{"empty input", []Token{}, TokenOr, [][]Token{}},
		{"no separators", []Token{fieldA, fieldB}, TokenOr, [][]Token{{fieldA, fieldB}}},
		{
			"consecutive separators",
			[]Token{fieldA, orToken, orToken, fieldB},
			TokenOr,
			[][]Token{{fieldA}, {fieldB}},
		},
		{"trailing separator", []Token{fieldA, orToken}, TokenOr, [][]Token{{fieldA}}},
		{"leading separator", []Token{orToken, fieldA}, TokenOr, [][]Token{{fieldA}}},
		{
			"multiple consecutive separators",
			[]Token{fieldA, orToken, orToken, orToken, fieldB},
			TokenOr,
			[][]Token{{fieldA}, {fieldB}},
		},
		{
			"normal case with single separators",
			[]Token{fieldA, orToken, fieldB, orToken, fieldC},
			TokenOr,
			[][]Token{{fieldA}, {fieldB}, {fieldC}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitByOperator(tt.tokens, tt.opType)
			require.Equal(t, tt.expected, result)
		})
	}
}
