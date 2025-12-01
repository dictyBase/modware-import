package stockcenter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSummarySemigroup(t *testing.T) {
	// Create two summaries with errors
	s1 := KeywordProcessingSummary{
		ErrorCount: 1,
		Errors:     []string{"error 1"},
	}
	s2 := KeywordProcessingSummary{
		ErrorCount: 1,
		Errors:     []string{"error 2"},
	}

	sg := SummarySemigroup()
	result := sg.Concat(s1, s2)

	require.Equal(t, 2, result.ErrorCount)
	require.Len(t, result.Errors, 2)
	require.Contains(t, result.Errors, "error 1")
	require.Contains(t, result.Errors, "error 2")
}
