package seqmodel_test

import (
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

func TestValidateSequenceRejectsNonASCIIAndWhitespace(t *testing.T) {
	for _, input := range []string{" ACGT", "ACGT\n", "ＡＣＧＴ", "ACGU", "acgtn"} {
		require.ErrorIs(
			t,
			E.ToError(seqmodel.ValidateSequence(input)),
			seqmodel.ErrInvalidDNABase,
		)
	}
}
