package pagination

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
)

// EncodeCursor encodes an offset into a cursor value
// For simplicity, we use the offset directly as the cursor
func EncodeCursor(offset int64) int64 {
	return offset
}

// DecodeCursor decodes a cursor value into an offset
// Returns Either[error, int64] for safe error handling
func DecodeCursor(cursor int64) E.Either[error, int64] {
	if cursor < 0 {
		return E.Left[int64](fmt.Errorf("invalid cursor: must be non-negative, got %d", cursor))
	}
	return E.Right[error](cursor)
}
