package filter

import (
	"strconv"
	"strings"
	"time"

	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

// OperatorEvaluator is a function that evaluates an operator against actual and expected values
type OperatorEvaluator func(actual any, expected string) bool

// operatorEvaluators maps operators to their evaluation functions
var operatorEvaluators = map[Operator]OperatorEvaluator{
	// String operators
	Contains:    evalStringContains,
	NotContains: evalStringNotContains,
	Equals:      evalStringEquals,
	NotEquals:   evalStringNotEquals,

	// Numeric operators
	NumEquals:      evalNumericEquals,
	GreaterThan:    evalNumericGreaterThan,
	LessThan:       evalNumericLessThan,
	GreaterOrEqual: evalNumericGreaterOrEqual,
	LessOrEqual:    evalNumericLessOrEqual,

	// Date operators
	DateEquals:         evalDateEquals,
	DateGreater:        evalDateGreater,
	DateLess:           evalDateLess,
	DateGreaterOrEqual: evalDateGreaterOrEqual,
	DateLessOrEqual:    evalDateLessOrEqual,

	// Array operators
	ArrayContains:    evalArrayContains,
	ArrayNotContains: evalArrayNotContains,
	ArrayEquals:      evalArrayEquals,
	ArrayNotEquals:   evalArrayNotEquals,
}

// UpdatePredicateEvaluate updates the Predicate's Evaluate method
func (pred Predicate) Evaluate(data map[string]any) bool {
	// Get field value from data
	fieldValue, exists := data[pred.Field]
	if !exists {
		return false
	}

	// Get evaluator for operator
	evaluator, ok := operatorEvaluators[pred.Operator]
	if !ok {
		return false
	}

	return evaluator(fieldValue, pred.Value)
}

// String evaluators using Option for type safety

func evalStringContains(actual any, expected string) bool {
	return F.Pipe2(
		extractString(actual),
		O.Map(func(str string) bool {
			return strings.Contains(str, expected)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

func evalStringNotContains(actual any, expected string) bool {
	return F.Pipe2(
		extractString(actual),
		O.Map(func(str string) bool {
			return !strings.Contains(str, expected)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

func evalStringEquals(actual any, expected string) bool {
	return F.Pipe2(
		extractString(actual),
		O.Map(func(str string) bool {
			return str == expected
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

func evalStringNotEquals(actual any, expected string) bool {
	return F.Pipe2(
		extractString(actual),
		O.Map(func(str string) bool {
			return str != expected
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

// Numeric evaluators using Option chains

func evalNumericEquals(actual any, expected string) bool {
	return evaluateNumeric(
		actual,
		expected,
		func(a, b float64) bool { return a == b },
	)
}

func evalNumericGreaterThan(actual any, expected string) bool {
	return evaluateNumeric(
		actual,
		expected,
		func(a, b float64) bool { return a > b },
	)
}

func evalNumericLessThan(actual any, expected string) bool {
	return evaluateNumeric(
		actual,
		expected,
		func(a, b float64) bool { return a < b },
	)
}

func evalNumericGreaterOrEqual(actual any, expected string) bool {
	return evaluateNumeric(
		actual,
		expected,
		func(a, b float64) bool { return a >= b },
	)
}

func evalNumericLessOrEqual(actual any, expected string) bool {
	return evaluateNumeric(
		actual,
		expected,
		func(a, b float64) bool { return a <= b },
	)
}

func evaluateNumeric(
	actual any,
	expected string,
	comparator func(float64, float64) bool,
) bool {
	return F.Pipe2(
		extractFloat64(actual),
		O.Chain(func(actualNum float64) O.Option[bool] {
			return F.Pipe1(
				parseFloat(expected),
				O.Map(func(expectedNum float64) bool {
					return comparator(actualNum, expectedNum)
				}),
			)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

// Date evaluators

func evalDateEquals(actual any, expected string) bool {
	return evaluateDate(
		actual,
		expected,
		func(a, b time.Time) bool { return a.Equal(b) },
	)
}

func evalDateGreater(actual any, expected string) bool {
	return evaluateDate(
		actual,
		expected,
		func(a, b time.Time) bool { return a.After(b) },
	)
}

func evalDateLess(actual any, expected string) bool {
	return evaluateDate(
		actual,
		expected,
		func(a, b time.Time) bool { return a.Before(b) },
	)
}

func evalDateGreaterOrEqual(actual any, expected string) bool {
	return evaluateDate(actual, expected, func(a, b time.Time) bool {
		return a.After(b) || a.Equal(b)
	})
}

func evalDateLessOrEqual(actual any, expected string) bool {
	return evaluateDate(actual, expected, func(a, b time.Time) bool {
		return a.Before(b) || a.Equal(b)
	})
}

func evaluateDate(
	actual any,
	expected string,
	comparator func(time.Time, time.Time) bool,
) bool {
	return F.Pipe2(
		extractTime(actual),
		O.Chain(func(actualTime time.Time) O.Option[bool] {
			return F.Pipe1(
				parseTime(expected),
				O.Map(func(expectedTime time.Time) bool {
					return comparator(actualTime, expectedTime)
				}),
			)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

// Array evaluators

func evalArrayContains(actual any, expected string) bool {
	return F.Pipe2(
		extractArray(actual),
		O.Map(func(arr []any) bool {
			for _, item := range arr {
				if itemMatches(item, expected) {
					return true
				}
			}
			return false
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

func evalArrayNotContains(actual any, expected string) bool {
	return F.Pipe2(
		extractArray(actual),
		O.Map(func(arr []any) bool {
			for _, item := range arr {
				if itemMatches(item, expected) {
					return false
				}
			}
			return true
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

func evalArrayEquals(actual any, expected string) bool {
	return F.Pipe2(
		extractArray(actual),
		O.Map(func(arr []any) bool {
			if len(arr) != 1 {
				return false
			}
			return itemMatches(arr[0], expected)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

func evalArrayNotEquals(actual any, expected string) bool {
	return F.Pipe2(
		extractArray(actual),
		O.Map(func(arr []any) bool {
			if len(arr) != 1 {
				return true
			}
			return !itemMatches(arr[0], expected)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}

// Helper functions using Option for type safety

func extractString(val any) O.Option[string] {
	str, ok := val.(string)
	if !ok {
		return O.None[string]()
	}
	return O.Some(str)
}

func extractFloat64(val any) O.Option[float64] {
	switch num := val.(type) {
	case float64:
		return O.Some(num)
	case float32:
		return O.Some(float64(num))
	case int:
		return O.Some(float64(num))
	case int64:
		return O.Some(float64(num))
	default:
		return O.None[float64]()
	}
}

func extractArray(val any) O.Option[[]any] {
	arr, ok := val.([]any)
	if !ok {
		return O.None[[]any]()
	}
	return O.Some(arr)
}

func extractTime(val any) O.Option[time.Time] {
	// Try parsing as string first (common case from JSON)
	if str, ok := val.(string); ok {
		return parseTime(str)
	}

	// Try direct time.Time
	if timeVal, ok := val.(time.Time); ok {
		return O.Some(timeVal)
	}

	return O.None[time.Time]()
}

func parseFloat(str string) O.Option[float64] {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return O.None[float64]()
	}
	return O.Some(num)
}

func parseTime(str string) O.Option[time.Time] {
	// Try RFC3339 format first (ISO 8601)
	timeVal, err := time.Parse(time.RFC3339, str)
	if err == nil {
		return O.Some(timeVal)
	}

	// Try date-only format
	timeVal, err = time.Parse("2006-01-02", str)
	if err == nil {
		return O.Some(timeVal)
	}

	return O.None[time.Time]()
}

func itemMatches(item any, expected string) bool {
	return F.Pipe2(
		extractString(item),
		O.Map(func(str string) bool {
			return strings.Contains(str, expected)
		}),
		O.GetOrElse(func() bool { return false }),
	)
}
