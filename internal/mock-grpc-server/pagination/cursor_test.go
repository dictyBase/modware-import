package pagination

import (
	"testing"

	E "github.com/IBM/fp-go/either"
	"github.com/stretchr/testify/require"
)

func TestEncodeCursor(t *testing.T) {
	offset := int64(42)
	cursor := EncodeCursor(offset)

	require.Equal(t, int64(42), cursor)
}

func TestDecodeCursor(t *testing.T) {
	cursor := int64(42)
	result := DecodeCursor(cursor)

	require.True(t, E.IsRight(result))
	offset := E.GetOrElse(func(error) int64 { return -1 })(result)
	require.Equal(t, int64(42), offset)
}

func TestDecodeNegativeCursor(t *testing.T) {
	cursor := int64(-1)
	result := DecodeCursor(cursor)

	require.True(t, E.IsLeft(result), "negative cursor should return Left")
}

func TestDecodeZeroCursor(t *testing.T) {
	cursor := int64(0)
	result := DecodeCursor(cursor)

	require.True(t, E.IsRight(result))
	offset := E.GetOrElse(func(error) int64 { return -1 })(result)
	require.Equal(t, int64(0), offset)
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []int64{0, 1, 10, 100, 1000, 9999}

	for _, original := range tests {
		cursor := EncodeCursor(original)
		result := DecodeCursor(cursor)

		require.True(t, E.IsRight(result))
		decoded := E.GetOrElse(func(error) int64 { return -1 })(result)
		require.Equal(t, original, decoded)
	}
}
