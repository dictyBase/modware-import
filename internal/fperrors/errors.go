package fperrors

import "fmt"

// OnError returns a function that wraps an error with additional context
func OnError(message string) func(error) error {
	return func(err error) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("%s: %w", message, err)
	}
}
