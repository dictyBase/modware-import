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

// typeAssertString performs type assertion to string, returning (value, success)
func typeAssertString(val any) (string, bool) {
	str, ok := val.(string)
	return str, ok
}

// extractString converts any value to Option[string] using Optionize1
var extractString = O.Optionize1(typeAssertString)

// tryExtractFloat64 attempts direct float64 type assertion
func tryExtractFloat64(val any) O.Option[float64] {
	if num, ok := val.(float64); ok {
		return O.Some(num)
	}
	return O.None[float64]()
}

// tryExtractFloat32 attempts float32 type assertion and converts to float64
func tryExtractFloat32(val any) O.Option[float64] {
	if num, ok := val.(float32); ok {
		return O.Some(float64(num))
	}
	return O.None[float64]()
}

// tryExtractInt attempts int type assertion and converts to float64
func tryExtractInt(val any) O.Option[float64] {
	if num, ok := val.(int); ok {
		return O.Some(float64(num))
	}
	return O.None[float64]()
}

// tryExtractInt64 attempts int64 type assertion and converts to float64
func tryExtractInt64(val any) O.Option[float64] {
	if num, ok := val.(int64); ok {
		return O.Some(float64(num))
	}
	return O.None[float64]()
}

// extractFloat64 tries multiple numeric type conversions using O.Alt chain
func extractFloat64(val any) O.Option[float64] {
	return F.Pipe3(
		tryExtractFloat64(val),
		O.Alt(func() O.Option[float64] { return tryExtractFloat32(val) }),
		O.Alt(func() O.Option[float64] { return tryExtractInt(val) }),
		O.Alt(func() O.Option[float64] { return tryExtractInt64(val) }),
	)
}

// typeAssertArray performs type assertion to []any, returning (value, success)
func typeAssertArray(val any) ([]any, bool) {
	arr, ok := val.([]any)
	return arr, ok
}

// extractArray converts any value to Option[[]any] using Optionize1
var extractArray = O.Optionize1(typeAssertArray)

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

// parseFloat converts string to float64 using O.TryCatch for error handling
func parseFloat(str string) O.Option[float64] {
	return O.TryCatch(func() (float64, error) {
		return strconv.ParseFloat(str, 64)
	})
}

// tryParseRFC3339 attempts to parse time in RFC3339 format (ISO 8601)
func tryParseRFC3339(str string) O.Option[time.Time] {
	return O.TryCatch(func() (time.Time, error) {
		return time.Parse(time.RFC3339, str)
	})
}

// tryParseDateOnly attempts to parse time in date-only format
func tryParseDateOnly(str string) O.Option[time.Time] {
	return O.TryCatch(func() (time.Time, error) {
		return time.Parse("2006-01-02", str)
	})
}

// parseTime converts string to time.Time using O.Alt for fallback logic
func parseTime(str string) O.Option[time.Time] {
	// Try RFC3339 first, fall back to date-only format
	return F.Pipe1(
		tryParseRFC3339(str),
		O.Alt(func() O.Option[time.Time] {
			return tryParseDateOnly(str)
		}),
	)
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
