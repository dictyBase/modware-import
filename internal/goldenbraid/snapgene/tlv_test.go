package snapgene_test

import (
	"bytes"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/require"

	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
)

func TestPacketsSkipsRepeatedUnusedBodies(t *testing.T) {
	data := cookie()
	for range 1000 {
		data = append(data, encodePacket(11, []byte("history"))...)
	}
	data = append(data, encodePacket(255, []byte("unknown"))...)
	data = append(data, dna(15, bases)...)
	packets, err := E.UnwrapError(snapgene.Packets(data))
	require.NoError(t, err)
	require.Equal(t, 2, len(packets), "only recognized payloads need retained packet state")
	require.Equal(t, byte(9), packets[0].ID)
	require.Equal(t, byte(0), packets[1].ID)
}

func TestPacketsLimitAppliesPerBodyNotFile(t *testing.T) {
	body := bytes.Repeat([]byte{'x'}, int(snapgene.MaxPacketBytes))
	data := stream(cookie(), encodePacket(11, body), dna(15, bases))
	require.Greater(t, len(data), int(snapgene.MaxPacketBytes))
	packets, err := E.UnwrapError(snapgene.Packets(data))
	require.NoError(t, err)
	require.Equal(t, 2, len(packets))
	require.Equal(t, bases, string(packets[1].Body[1:]))
}

func TestParseRejectsLeadingUnusedPacket(t *testing.T) {
	data := stream(encodePacket(11, []byte("history")), cookie(), dna(15, bases))
	require.ErrorIs(
		t,
		E.ToError(snapgene.Parse(snapgene.ParseInput{Name: "leading", Data: data})),
		snapgene.ErrInvalidMagicCookie,
	)
}
