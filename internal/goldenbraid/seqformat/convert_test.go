package seqformat_test

import (
	"bytes"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/bebop/poly/io/genbank"
	"github.com/stretchr/testify/require"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
	"github.com/dictyBase/modware-import/internal/goldenbraid/seqformat"
)

func TestToGenBankBytesDeterministicSortedQualifiers(t *testing.T) {
	record := testRecord()
	first, err := E.UnwrapError(seqformat.ToGenBankBytes(record))
	require.NoError(t, err)
	second, err := E.UnwrapError(seqformat.ToGenBankBytes(record))
	require.NoError(t, err)
	require.Equal(t, first, second)
	text := string(first)
	require.Equal(t, 1, strings.Count(text, "FEATURES             Location/Qualifiers\n"))
	require.Contains(t, text, "     CDS             complement(join(5..20,40..60))\n")
	require.Contains(
		t,
		text,
		"                     /ApEinfo_fwdcolor=\"#0a0b0c\"\n                     /label=\"mNeonGreen\"\n",
	)
	require.Equal(
		t,
		[]string{"CDS", "primer_bind"},
		[]string{record.Features[0].Type, record.Features[1].Type},
	)
	require.Equal(t, "#0a0b0c", record.Features[0].Attributes["ApEinfo_fwdcolor"])
}

func TestToGenBankBytesRoundTripsSharedRecord(t *testing.T) {
	record := testRecord()
	data, err := E.UnwrapError(seqformat.ToGenBankBytes(record))
	require.NoError(t, err)
	parsed, err := genbank.Parse(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, record.Sequence, parsed.Sequence)
	require.Equal(t, record.Meta.Locus.Circular, parsed.Meta.Locus.Circular)
	require.Equal(t, record.Meta.Locus.Name, parsed.Meta.Locus.Name)
	require.Equal(t, record.Meta.Definition, parsed.Meta.Definition)
	require.Equal(t, record.Meta.Other["COMMENT"], parsed.Meta.Other["COMMENT"])
	require.Len(t, parsed.Features, 2)
	require.Equal(t, "CDS", parsed.Features[0].Type)
	require.Equal(
		t,
		"complement(join(5..20,40..60))",
		genbank.BuildLocationString(parsed.Features[0].Location),
	)
	require.True(t, parsed.Features[0].Location.Complement)
	require.True(t, parsed.Features[0].Location.Join)
	require.Equal(t, [][2]int{{4, 20}, {39, 60}}, [][2]int{
		{
			parsed.Features[0].Location.SubLocations[0].Start,
			parsed.Features[0].Location.SubLocations[0].End,
		},
		{
			parsed.Features[0].Location.SubLocations[1].Start,
			parsed.Features[0].Location.SubLocations[1].End,
		},
	})
	require.Equal(t, "mNeonGreen", parsed.Features[0].Attributes["label"])
	require.Equal(t, "#0a0b0c", parsed.Features[0].Attributes["ApEinfo_fwdcolor"])
	require.Equal(t, "primer_bind", parsed.Features[1].Type)
	require.Equal(t, "80..90", genbank.BuildLocationString(parsed.Features[1].Location))
	require.False(t, parsed.Features[1].Location.Complement)
	require.Equal(t, "primer_F", parsed.Features[1].Attributes["label"])
	require.NotContains(t, parsed.Features[1].Attributes, "ApEinfo_fwdcolor")
}

func TestToFASTABytesUsesFullNameAndSequence(t *testing.T) {
	record := testRecord()
	data, err := E.UnwrapError(seqformat.ToFASTABytes(record))
	require.NoError(t, err)
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	require.Equal(t, ">pDGB_A1_act15p_mNeonGreen_act8t", lines[0])
	require.NotEqual(t, ">"+record.Meta.Locus.Name, lines[0])
	require.Equal(t, record.Sequence, strings.Join(lines[1:], ""))
	require.Len(t, lines[1], 80)
}

func TestToGenBankBytesEscapesEmbeddedQualifierQuotes(t *testing.T) {
	record := testRecord()
	record.Features[0].Attributes["note"] = "protein \"tag\""
	data, err := E.UnwrapError(seqformat.ToGenBankBytes(record))
	require.NoError(t, err)
	require.Contains(t, string(data), "                     /note=\"protein \"\"tag\"\"\"\n")
}

func TestToGenBankBytesPreservesHeaderAndSequence(t *testing.T) {
	record := testRecord()
	emptyFeatures := F.Pipe1(record, PL.GenbankFeatures.Set(nil))
	baseline, err := genbank.Build(emptyFeatures)
	require.NoError(t, err)
	data, err := E.UnwrapError(seqformat.ToGenBankBytes(record))
	require.NoError(t, err)
	before, after, found := strings.Cut(
		string(baseline),
		"FEATURES             Location/Qualifiers\nORIGIN\n",
	)
	require.True(t, found)
	require.True(t, strings.HasPrefix(string(data), before))
	require.True(t, strings.HasSuffix(string(data), "ORIGIN\n"+after))
}

func TestToGenBankBytesRejectsMissingFeatureBoundary(t *testing.T) {
	record := testRecord()
	record.Meta.BaseCount = []genbank.BaseCount{{Base: "a", Count: 30}}
	_, err := E.UnwrapError(seqformat.ToGenBankBytes(record))
	require.Error(t, err)
	require.ErrorIs(t, err, seqformat.ErrMissingFeatureBoundary)
}

func testRecord() genbank.Genbank {
	first := F.Pipe2(genbank.Location{}, PL.LocationStart.Set(4), PL.LocationEnd.Set(20))
	second := F.Pipe2(genbank.Location{}, PL.LocationStart.Set(39), PL.LocationEnd.Set(60))
	joined := F.Pipe3(
		genbank.Location{},
		PL.LocationComplement.Set(true),
		PL.LocationJoin.Set(true),
		PL.LocationSubLocations.Set([]genbank.Location{first, second}),
	)
	coding := F.Pipe3(
		genbank.Feature{},
		PL.FeatureType.Set("CDS"),
		PL.FeatureLocation.Set(joined),
		PL.FeatureAttributes.Set(
			map[string]string{"label": "mNeonGreen", "ApEinfo_fwdcolor": "#0a0b0c"},
		),
	)
	primerLocation := F.Pipe2(genbank.Location{}, PL.LocationStart.Set(79), PL.LocationEnd.Set(90))
	primer := F.Pipe3(
		genbank.Feature{},
		PL.FeatureType.Set("primer_bind"),
		PL.FeatureLocation.Set(primerLocation),
		PL.FeatureAttributes.Set(map[string]string{"label": "primer_F"}),
	)
	locus := F.Pipe7(
		genbank.Locus{},
		PL.LocusName.Set("pDGB_A1_act15p"),
		PL.LocusSequenceLength.Set("120"),
		PL.LocusMoleculeType.Set(
			"DNA",
		),
		PL.LocusCircular.Set(true),
		PL.LocusGenbankDivision.Set("SYN"),
		PL.LocusModificationDate.Set("01-JAN-2000"),
		PL.LocusSequenceCoding.Set("bp"),
	)
	meta := F.Pipe4(
		genbank.Meta{},
		PL.MetaLocus.Set(locus),
		PL.MetaName.Set("pDGB_A1_act15p_mNeonGreen_act8t"),
		PL.MetaDefinition.Set(
			"Test construct.",
		),
		PL.MetaOther.Set(map[string]string{"COMMENT": "SnapGene UUID example."}),
	)
	return F.Pipe3(
		genbank.Genbank{},
		PL.GenbankMeta.Set(meta),
		PL.GenbankSequence.Set(strings.Repeat("ACGT", 30)),
		PL.GenbankFeatures.Set([]genbank.Feature{coding, primer}),
	)
}
