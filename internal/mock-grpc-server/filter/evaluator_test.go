package filter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEvaluateStringEquals(t *testing.T) {
	pred := Predicate{
		Field:    "depositor",
		Operator: Equals,
		Value:    "Costanza",
	}

	data := map[string]interface{}{
		"depositor": "Costanza",
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateStringNotEquals(t *testing.T) {
	pred := Predicate{
		Field:    "depositor",
		Operator: NotEquals,
		Value:    "Seinfeld",
	}

	data := map[string]interface{}{
		"depositor": "Costanza",
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateStringContains(t *testing.T) {
	pred := Predicate{
		Field:    "label",
		Operator: Contains,
		Value:    "axe",
	}

	data := map[string]interface{}{
		"label": "axeA2 axeB2 axeC2",
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateStringNotContains(t *testing.T) {
	pred := Predicate{
		Field:    "label",
		Operator: NotContains,
		Value:    "xyz",
	}

	data := map[string]interface{}{
		"label": "axeA2 axeB2 axeC2",
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateNumericEquals(t *testing.T) {
	pred := Predicate{
		Field:    "quantity",
		Operator: NumEquals,
		Value:    "42",
	}

	data := map[string]interface{}{
		"quantity": float64(42),
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateNumericGreaterThan(t *testing.T) {
	pred := Predicate{
		Field:    "quantity",
		Operator: GreaterThan,
		Value:    "10",
	}

	data := map[string]interface{}{
		"quantity": float64(42),
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateNumericLessThan(t *testing.T) {
	pred := Predicate{
		Field:    "quantity",
		Operator: LessThan,
		Value:    "100",
	}

	data := map[string]interface{}{
		"quantity": float64(42),
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateDateGreaterOrEqual(t *testing.T) {
	pred := Predicate{
		Field:    "created_at",
		Operator: DateGreaterOrEqual,
		Value:    "2018-12-01",
	}

	// Simulate timestamp from protobuf conversion
	createdTime := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	timestamp := timestamppb.New(createdTime)

	data := map[string]interface{}{
		"created_at": timestamp.AsTime().Format(time.RFC3339),
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateDateLess(t *testing.T) {
	pred := Predicate{
		Field:    "created_at",
		Operator: DateLess,
		Value:    "2020-01-01",
	}

	createdTime := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	timestamp := timestamppb.New(createdTime)

	data := map[string]interface{}{
		"created_at": timestamp.AsTime().Format(time.RFC3339),
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateArrayContains(t *testing.T) {
	pred := Predicate{
		Field:    "tags",
		Operator: ArrayContains,
		Value:    "important",
	}

	data := map[string]interface{}{
		"tags": []interface{}{"urgent", "important", "flagged"},
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateArrayNotContains(t *testing.T) {
	pred := Predicate{
		Field:    "tags",
		Operator: ArrayNotContains,
		Value:    "spam",
	}

	data := map[string]interface{}{
		"tags": []interface{}{"urgent", "important", "flagged"},
	}

	result := pred.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateMissingField(t *testing.T) {
	pred := Predicate{
		Field:    "nonexistent",
		Operator: Equals,
		Value:    "value",
	}

	data := map[string]interface{}{
		"depositor": "Costanza",
	}

	result := pred.Evaluate(data)
	require.False(t, result, "missing field should return false")
}

func TestEvaluateTypeMismatch(t *testing.T) {
	pred := Predicate{
		Field:    "depositor",
		Operator: GreaterThan,
		Value:    "10",
	}

	data := map[string]interface{}{
		"depositor": "Costanza",
	}

	result := pred.Evaluate(data)
	require.False(t, result, "type mismatch should return false")
}

func TestEvaluateAndExpression(t *testing.T) {
	expr := AndExpression{
		Left: Predicate{
			Field:    "depositor",
			Operator: Equals,
			Value:    "Costanza",
		},
		Right: Predicate{
			Field:    "species",
			Operator: Equals,
			Value:    "Dictyostelium",
		},
	}

	data := map[string]interface{}{
		"depositor": "Costanza",
		"species":   "Dictyostelium",
	}

	result := expr.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateOrExpression(t *testing.T) {
	expr := OrExpression{
		Left: Predicate{
			Field:    "depositor",
			Operator: Equals,
			Value:    "Costanza",
		},
		Right: Predicate{
			Field:    "depositor",
			Operator: Equals,
			Value:    "Seinfeld",
		},
	}

	data := map[string]interface{}{
		"depositor": "Costanza",
	}

	result := expr.Evaluate(data)
	require.True(t, result)
}

func TestEvaluateComplexExpression(t *testing.T) {
	// (depositor === Costanza AND species === Dictyostelium) OR label =~ axe
	expr := OrExpression{
		Left: AndExpression{
			Left: Predicate{
				Field:    "depositor",
				Operator: Equals,
				Value:    "Costanza",
			},
			Right: Predicate{
				Field:    "species",
				Operator: Equals,
				Value:    "Dictyostelium",
			},
		},
		Right: Predicate{
			Field:    "label",
			Operator: Contains,
			Value:    "axe",
		},
	}

	// Test case 1: Matches left side of OR
	data1 := map[string]interface{}{
		"depositor": "Costanza",
		"species":   "Dictyostelium",
		"label":     "wild type",
	}
	require.True(t, expr.Evaluate(data1))

	// Test case 2: Matches right side of OR
	data2 := map[string]interface{}{
		"depositor": "Seinfeld",
		"species":   "Other",
		"label":     "axeA2",
	}
	require.True(t, expr.Evaluate(data2))

	// Test case 3: Matches neither
	data3 := map[string]interface{}{
		"depositor": "Seinfeld",
		"species":   "Other",
		"label":     "wild type",
	}
	require.False(t, expr.Evaluate(data3))
}
