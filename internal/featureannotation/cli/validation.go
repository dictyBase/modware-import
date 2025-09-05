package cli

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator instance for the CLI package
var validate = validator.New()

// LoadGeneDescriptionParams holds validated parameters for loading gene descriptions
type LoadGeneDescriptionParams struct {
	InputFile string `validate:"required,file"`
	Workers   int    `validate:"gte=1,lte=50"`
	BatchSize int    `validate:"gte=1,lte=1000"`
	User      string `validate:"required,email"`
}

// LoadGeneProductParams holds validated parameters for loading gene products
type LoadGeneProductParams struct {
	InputFiles []string `validate:"required,dive,file"`
	Workers    int      `validate:"gte=1,lte=50"`
	BatchSize  int      `validate:"gte=1,lte=1000"`
	User       string   `validate:"required,email"`
}

// GeneDescriptionRequest holds validated data for a single gene description
type GeneDescriptionRequest struct {
	GeneID      string `validate:"required,min=1"`
	Description string `validate:"required,min=1"`
	User        string `validate:"required,email"`
}

// CLIValidationError wraps multiple validation errors from go-playground/validator
type CLIValidationError struct {
	Errors validator.ValidationErrors
}

func (e *CLIValidationError) Error() string {
	var messages []string
	for _, err := range e.Errors {
		switch err.Tag() {
		case "required":
			messages = append(
				messages,
				fmt.Sprintf("%s is required", err.Field()),
			)
		case "email":
			messages = append(
				messages,
				fmt.Sprintf("%s must be a valid email address", err.Field()),
			)
		case "min":
			messages = append(
				messages,
				fmt.Sprintf(
					"%s must be at least %s characters long",
					err.Field(),
					err.Param(),
				),
			)
		case "max":
			messages = append(
				messages,
				fmt.Sprintf(
					"%s must be at most %s characters long",
					err.Field(),
					err.Param(),
				),
			)
		case "gte":
			messages = append(
				messages,
				fmt.Sprintf(
					"%s must be greater than or equal to %s",
					err.Field(),
					err.Param(),
				),
			)
		case "lte":
			messages = append(
				messages,
				fmt.Sprintf(
					"%s must be less than or equal to %s",
					err.Field(),
					err.Param(),
				),
			)
		case "file":
			messages = append(
				messages,
				fmt.Sprintf(
					"%s must be a valid file path that exists",
					err.Field(),
				),
			)
		default:
			messages = append(
				messages,
				fmt.Sprintf(
					"%s failed validation rule: %s",
					err.Field(),
					err.Tag(),
				),
			)
		}
	}
	return strings.Join(messages, "; ")
}

// ValidateStruct validates a struct and returns a formatted error if validation fails
func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return &CLIValidationError{Errors: validationErrors}
		}
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
