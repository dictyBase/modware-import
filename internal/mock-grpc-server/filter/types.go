package filter

// TokenType represents the type of a token in the filter expression
type TokenType int

const (
	// TokenField represents a field name (e.g., "depositor", "created_at")
	TokenField TokenType = iota
	// TokenOperator represents a comparison operator (e.g., "===", "$>=")
	TokenOperator
	// TokenValue represents a value to compare against
	TokenValue
	// TokenAnd represents the AND logical operator (;)
	TokenAnd
	// TokenOr represents the OR logical operator (,)
	TokenOr
)

// Token represents a single token in the filter expression
type Token struct {
	Type  TokenType
	Value string
}

// Operator represents a filter comparison operator
type Operator int

const (
	// String operators
	Contains Operator = iota
	NotContains
	Equals
	NotEquals

	// Numeric operators
	NumEquals
	GreaterThan
	LessThan
	GreaterOrEqual
	LessOrEqual

	// Date operators
	DateEquals
	DateGreater
	DateLess
	DateGreaterOrEqual
	DateLessOrEqual

	// Array operators
	ArrayContains
	ArrayNotContains
	ArrayEquals
	ArrayNotEquals
)

// FilterExpression represents an evaluatable filter expression
type FilterExpression interface {
	Evaluate(data map[string]any) bool
}

// Predicate represents a single comparison operation
type Predicate struct {
	Field    string
	Operator Operator
	Value    string
}

// Evaluate method is implemented in evaluator.go

// AndExpression represents a logical AND of two expressions
type AndExpression struct {
	Left  FilterExpression
	Right FilterExpression
}

// Evaluate evaluates the AND expression
func (expr AndExpression) Evaluate(data map[string]interface{}) bool {
	return expr.Left.Evaluate(data) && expr.Right.Evaluate(data)
}

// OrExpression represents a logical OR of two expressions
type OrExpression struct {
	Left  FilterExpression
	Right FilterExpression
}

// Evaluate evaluates the OR expression
func (expr OrExpression) Evaluate(data map[string]interface{}) bool {
	return expr.Left.Evaluate(data) || expr.Right.Evaluate(data)
}

// AlwaysTrueFilter is a filter that always returns true
type AlwaysTrueFilter struct{}

// Evaluate always returns true
func (AlwaysTrueFilter) Evaluate(_ map[string]interface{}) bool {
	return true
}
