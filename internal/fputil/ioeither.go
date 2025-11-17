// Package fputil provides utility functions for working with fp-go types
package fputil

import (
	E "github.com/IBM/fp-go/either"
	IOE "github.com/IBM/fp-go/ioeither"
)

// ToEither executes an IOEither to get an Either result.
// This is useful for converting deferred side effects into immediate Either values.
//
// Example:
//
//	result := ToEither(loadFromAPI(id))
//	E.Fold(handleError, handleSuccess)(result)
func ToEither[ER, A any](ioe IOE.IOEither[ER, A]) E.Either[ER, A] {
	return ioe()
}
