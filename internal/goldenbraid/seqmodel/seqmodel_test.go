package seqmodel_test

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/require"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
)

func TestValidateSequence(t *testing.T) {
	require.Equal(t, E.Right[error]("ACGT"), seqmodel.ValidateSequence("acGt"))
	require.ErrorIs(
		t,
		E.ToError(seqmodel.ValidateSequence("ATGCN")),
		seqmodel.ErrInvalidDNABase,
	)
	require.ErrorIs(t, E.ToError(seqmodel.ValidateSequence("")), seqmodel.ErrEmptySequence)
}

func TestReduceExportResultsUsesMonoidConcat(t *testing.T) {
	corruptErr := errors.New("corrupt input")
	results := []E.Either[error, string]{
		E.Right[error]("valid.dna"),
		E.Left[string](corruptErr),
	}

	got := seqmodel.ReduceExportResults(results)
	sumL := seqmodel.MakeExportSummaryLenses()
	require.Equal(t, 2, sumL.Total.Get(got))
	require.Equal(t, 1, sumL.Succeeded.Get(got))
	require.Equal(t, 1, sumL.Failed.Get(got))
}

func TestParseColor(t *testing.T) {
	got := E.Map[error](seqmodel.Color.Hex)(seqmodel.ParseColor("#0A0B0C"))
	require.Equal(t, E.Right[error]("#0a0b0c"), got)
	for _, input := range []string{"#abc", "#0A0B0Z", "0A0B0C", "#0000000", "#000000\n"} {
		require.ErrorIs(
			t,
			E.ToError(seqmodel.ParseColor(input)),
			seqmodel.ErrInvalidColor,
		)
	}
}

func TestReduceEmptyAndOrderedErrors(t *testing.T) {
	first, second := errors.New("first"), errors.New("second")
	l := seqmodel.MakeExportSummaryLenses()
	empty := seqmodel.ReduceExportResults(nil)
	require.Zero(t, l.Total.Get(empty))
	require.Empty(t, l.Errors.Get(empty))
	got := seqmodel.ReduceExportResults(
		[]E.Either[error, string]{
			E.Left[string](first),
			E.Right[error]("ok"),
			E.Left[string](second),
		},
	)
	require.Equal(t, []error{first, second}, l.Errors.Get(got))
	require.Equal(t, 3, l.Total.Get(got))
	require.Equal(t, 1, l.Succeeded.Get(got))
	require.Equal(t, 2, l.Failed.Get(got))
}

func TestValidateSequenceRejectsNonASCIIAndWhitespace(t *testing.T) {
	for _, input := range []string{" ACGT", "ACGT\n", "ＡＣＧＴ", "ACGU", "acgtn"} {
		require.ErrorIs(
			t,
			E.ToError(seqmodel.ValidateSequence(input)),
			seqmodel.ErrInvalidDNABase,
		)
	}
}
