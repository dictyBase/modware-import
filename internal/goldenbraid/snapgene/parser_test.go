package snapgene_test

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/bebop/poly/io/genbank"
	"github.com/stretchr/testify/require"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
)

func encodePacket(id byte, body []byte) []byte {
	var out bytes.Buffer
	out.WriteByte(id)
	size := len(body)
	if size > int(snapgene.MaxPacketBytes) {
		panic("oversized test body")
	}
	_ = binary.Write(&out, binary.BigEndian, uint32(size))
	out.Write(body)
	return out.Bytes()
}

func cookie() []byte { return encodePacket(9, append([]byte("SnapGene"), 0, 1, 0, 13, 0, 14)) }

func dna(flags byte, seq string) []byte {
	return encodePacket(0, append([]byte{flags}, []byte(seq)...))
}

func stream(packets ...[]byte) []byte { return bytes.Join(packets, nil) }

func parse(t *testing.T, data []byte, name string) snapgene.ParseResult {
	t.Helper()
	result, err := E.UnwrapError(
		snapgene.Parse(snapgene.ParseInput{Name: name, SourceFile: name + ".dna", Data: data}),
	)
	require.NoError(t, err)
	return result
}

func featurePacket(entries string) []byte {
	return encodePacket(10, []byte("<Features>"+entries+"</Features>"))
}

func feature(name, span, color, direction string) string {
	return fmt.Sprintf(
		`<Feature name="%s" type="CDS" directionality="%s"><Segment range="%s" type="standard" color="%s"/></Feature>`,
		name,
		direction,
		span,
		color,
	)
}

const bases = "GGAGCAAAATTTTGTATTTT"

func TestParseValidCircularSnapGeneDNA(t *testing.T) {
	name := "a_full_name_over_sixteen_characters"
	data := stream(
		cookie(),
		featurePacket(feature("UTF8 β", "2-8", "#31849B", "1")),
		encodePacket(
			6,
			[]byte(
				`<Notes><Type>Synthetic</Type><UUID>uuid</UUID><Created>2019.8.29</Created>`+
					`<LastModified>2020.1.2</LastModified><TransformedInto>DH10B™</TransformedInto></Notes>`,
			),
		),
		encodePacket(
			5,
			[]byte(
				`<Primers><Primer name="p"><BindingSite location="1-3" boundStrand="0"/></Primer></Primers>`,
			),
		),
		dna(15, strings.ToLower(bases)),
	)
	got := parse(t, data, name)
	r := got.Record
	require.Empty(t, got.Warnings)
	require.Equal(t, bases, r.Sequence)
	require.True(t, r.Meta.Locus.Circular)
	require.Equal(t, "DNA", r.Meta.Locus.MoleculeType)
	require.Equal(t, "20", r.Meta.Locus.SequenceLength)
	require.Equal(t, name, r.Meta.Name)
	require.Equal(t, name[:16], r.Meta.Locus.Name)
	require.Equal(t, name+", synthetic construct.", r.Meta.Definition)
	require.Equal(
		t,
		map[string]string{
			"COMMENT": "SnapGene UUID uuid; created 2019.8.29; modified 2020.1.2; transformed into DH10B.",
		},
		r.Meta.Other,
	)
	require.Len(t, r.Features, 2)
	require.Equal(t, "CDS", r.Features[0].Type)
	require.Equal(
		t,
		map[string]string{"label": "UTF8 β", "ApEinfo_fwdcolor": "#31849b"},
		r.Features[0].Attributes,
	)
	require.Equal(t, genbank.Location{Start: 1, End: 8}, r.Features[0].Location)
	require.Equal(t, "primer_bind", r.Features[1].Type)
}

func TestParseOriginWrapRange(t *testing.T) {
	got := parse(
		t,
		stream(cookie(), dna(15, bases), featurePacket(feature("wrap", "15-4", "noColor", "2"))),
		"wrap",
	)
	require.Empty(t, got.Warnings)
	require.Equal(
		t,
		genbank.Location{
			Join:         true,
			Complement:   true,
			SubLocations: []genbank.Location{{Start: 14, End: 20}, {Start: 0, End: 4}},
		},
		got.Record.Features[0].Location,
	)
	require.NotContains(t, got.Record.Features[0].Attributes, "ApEinfo_fwdcolor")
}

func TestParseSkipsUnknownPacketIDs(t *testing.T) {
	data := cookie()
	for _, id := range []byte{2, 7, 8, 17, 255, 2} {
		data = append(data, encodePacket(id, []byte{1, 2, 3})...)
	}
	got := parse(t, append(data, dna(31, bases)...), "unknown")
	require.Equal(t, bases, got.Record.Sequence)
	require.True(t, got.Record.Meta.Locus.Circular)
}

func TestParseRejectsOversizedPacketWithSentinel(t *testing.T) {
	data := append(cookie(), 0, 127, 255, 255, 255)
	require.ErrorIs(
		t,
		E.ToError(snapgene.Parse(snapgene.ParseInput{Name: "huge", Data: data})),
		snapgene.ErrPacketTooLarge,
	)
}

func TestParseRejectsInvalidCookieWithSentinel(t *testing.T) {
	for _, data := range [][]byte{
		nil, dna(15, bases),
		encodePacket(9, []byte("bad cookie")),
		encodePacket(9, []byte("SnapGene")),
	} {
		require.ErrorIs(
			t,
			E.ToError(snapgene.Parse(snapgene.ParseInput{Name: "bad", Data: data})),
			snapgene.ErrInvalidMagicCookie,
		)
	}
}

func TestParseManyHistoryPackets(t *testing.T) {
	data := cookie()
	for range 1000 {
		data = append(data, encodePacket(11, []byte("history"))...)
	}
	got := parse(t, append(data, dna(15, bases)...), "history")
	require.Equal(t, bases, got.Record.Sequence)
}

func TestPacketsRejectsTruncatedFrames(t *testing.T) {
	for _, tail := range [][]byte{{10}, {10, 0, 0, 0}, {10, 0, 0, 0, 2, 1}} {
		require.ErrorIs(
			t,
			E.ToError(snapgene.Packets(append(cookie(), tail...))),
			snapgene.ErrTruncatedPacket,
		)
	}
}

func TestParseFatalDNANeverBecomesWarning(t *testing.T) {
	cases := []struct {
		data []byte
		err  error
	}{
		{cookie(), snapgene.ErrMissingDNAPacket},
		{stream(cookie(), encodePacket(0, nil)), seqmodel.ErrEmptySequence},
		{stream(cookie(), dna(15, "")), seqmodel.ErrEmptySequence},
		{stream(cookie(), dna(15, "ACGTN")), seqmodel.ErrInvalidDNABase},
	}
	for _, tc := range cases {
		require.ErrorIs(
			t,
			E.ToError(snapgene.Parse(snapgene.ParseInput{Name: "bad", Data: tc.data})),
			tc.err,
		)
	}
	require.ErrorIs(
		t,
		E.ToError(snapgene.Parse(snapgene.ParseInput{Data: stream(cookie(), dna(15, bases))})),
		seqmodel.ErrEmptyName,
	)
}

func TestDNAFlagsIndependent(t *testing.T) {
	for _, flags := range []byte{0, 1, 2, 3, 15, 31} {
		got := parse(t, stream(cookie(), dna(flags, bases)), "flags")
		require.Equal(t, flags&1 != 0, got.Record.Meta.Locus.Circular)
	}
}

func TestParseSkipsInvalidFeature(t *testing.T) {
	entries := feature(
		"one",
		"1-3",
		"",
		"0",
	) + feature(
		"bad",
		"1-21",
		"",
		"0",
	) + feature(
		"two",
		"4-7",
		"",
		"2",
	)
	got := parse(t, stream(cookie(), dna(15, bases), featurePacket(entries)), "features")
	require.Len(t, got.Warnings, 1)
	require.ErrorIs(t, got.Warnings[0], snapgene.ErrInvalidFeature)
	require.ErrorIs(t, got.Warnings[0], seqmodel.ErrSegmentOutOfBounds)
	require.Len(t, got.Record.Features, 2)
	require.Equal(t, "one", got.Record.Features[0].Attributes["label"])
	require.Equal(t, "two", got.Record.Features[1].Attributes["label"])
	require.True(t, got.Record.Features[1].Location.Complement)
}

func TestParseFeatureValidationAndGaps(t *testing.T) {
	for _, bad := range []string{
		feature("bad", "0-4", "", "0"),
		feature("bad", "21-4", "", "0"),
		feature("bad", "x-4", "", "0"),
		feature("bad", "1-4", "#oops", "0"),
		feature("bad", "1-4", "", "bad"),
		`<Feature name="gaps" type="CDS"><Segment range="1-3" type="gap"/></Feature>`,
	} {
		got := parse(
			t,
			stream(cookie(), dna(15, bases), featurePacket(bad+feature("ok", "1-4", "", "0"))),
			"badentry",
		)
		require.Len(t, got.Warnings, 1)
		require.ErrorIs(t, got.Warnings[0], snapgene.ErrInvalidFeature)
		require.Len(t, got.Record.Features, 1)
	}
	got := parse(
		t,
		stream(cookie(), dna(2, bases), featurePacket(feature("linearwrap", "15-4", "", "0"))),
		"linear",
	)
	require.Len(t, got.Warnings, 1)
	require.ErrorIs(t, got.Warnings[0], seqmodel.ErrSegmentOutOfBounds)
}

func TestParsePreservesSegmentOrderAndFirstStandardColor(t *testing.T) {
	entry := `<Feature name="joined" type="gene" directionality="2" swappedSegmentNumbering="1">` +
		`<Segment type="gap" range="4-7" color="#bad"/>` +
		`<Segment type="standard" range="12-15" color="noColor"/>` +
		`<Segment type="standard" range="1-3" color="#ff0000"/></Feature>`
	got := parse(t, stream(cookie(), dna(15, bases), featurePacket(entry)), "join")
	require.Empty(t, got.Warnings)
	require.NotContains(t, got.Record.Features[0].Attributes, "ApEinfo_fwdcolor")
	require.Equal(
		t,
		genbank.Location{
			Join:         true,
			Complement:   true,
			SubLocations: []genbank.Location{{Start: 11, End: 15}, {Start: 0, End: 3}},
		},
		got.Record.Features[0].Location,
	)
}

func TestParsePrimerSitesAndSimplifiedDuplicate(t *testing.T) {
	entry := `<Primers><Primer name="primer"><BindingSite location="0-3" boundStrand="0"/>` +
		`<BindingSite location="8-10" boundStrand="1"/>` +
		`<BindingSite location="8-10" boundStrand="1" simplified="1"/></Primer></Primers>`
	got := parse(t, stream(cookie(), encodePacket(5, []byte(entry)), dna(15, bases)), "primer")
	require.Empty(t, got.Warnings)
	require.Len(t, got.Record.Features, 2)
	require.Equal(t, genbank.Location{Start: 0, End: 4}, got.Record.Features[0].Location)
	require.Equal(
		t,
		genbank.Location{Start: 8, End: 11, Complement: true},
		got.Record.Features[1].Location,
	)
}

func TestParseSkipsInvalidPrimer(t *testing.T) {
	entry := `<Primers><Primer name="bad"><BindingSite location="0-20" boundStrand="0"/></Primer>` +
		`<Primer name="good"><BindingSite location="2-4" boundStrand="1"/></Primer></Primers>`
	got := parse(t, stream(cookie(), dna(15, bases), encodePacket(5, []byte(entry))), "primer")
	require.Len(t, got.Warnings, 1)
	require.ErrorIs(t, got.Warnings[0], snapgene.ErrInvalidPrimer)
	require.ErrorIs(t, got.Warnings[0], seqmodel.ErrSegmentOutOfBounds)
	require.Len(t, got.Record.Features, 1)
	require.Equal(t, "good", got.Record.Features[0].Attributes["label"])
}

func TestMalformedAnnotationPacketsPreserveSiblings(t *testing.T) {
	for _, tc := range []struct {
		id       byte
		data     string
		sentinel error
	}{
		{10, "<Features><Feature></Features>", snapgene.ErrInvalidFeature},
		{5, "<Primers><Primer></Primers>", snapgene.ErrInvalidPrimer},
	} {
		got := parse(
			t,
			stream(
				cookie(),
				dna(15, bases),
				encodePacket(tc.id, []byte(tc.data)),
				featurePacket(feature("valid", "1-5", "", "0")),
				encodePacket(
					5,
					[]byte(
						`<Primers><Primer name="p"><BindingSite location="0-2" boundStrand="0"/></Primer></Primers>`,
					),
				),
			),
			"broken",
		)
		require.Len(t, got.Warnings, 1)
		require.ErrorIs(t, got.Warnings[0], tc.sentinel)
		require.Len(t, got.Record.Features, 2)
	}
}

func TestMalformedNotesPreservesSequenceAndValidMetadata(t *testing.T) {
	broken := encodePacket(6, []byte("<Notes><Type></Notes>"))
	got := parse(t, stream(cookie(), dna(15, bases), broken), "fallback")
	require.Equal(t, bases, got.Record.Sequence)
	require.Equal(t, "fallback.", got.Record.Meta.Definition)
	require.Empty(t, got.Record.Meta.Other)
	require.Len(t, got.Warnings, 1)
	require.ErrorIs(t, got.Warnings[0], snapgene.ErrInvalidNotes)
	got = parse(
		t,
		stream(
			cookie(),
			dna(15, bases),
			encodePacket(6, []byte(`<Notes><Type>Synthetic</Type><UUID>valid</UUID></Notes>`)),
			broken,
		),
		"valid",
	)
	require.Equal(t, "valid, synthetic construct.", got.Record.Meta.Definition)
	require.Equal(t, "SnapGene UUID valid.", got.Record.Meta.Other["COMMENT"])
	require.Len(t, got.Warnings, 1)
}

func TestParseRealAssembly(
	t *testing.T,
) {
	checkFixture(t, "pDGB_A1_act15p_mNeonGreen_act8t", 3666, 8)
}

func TestParseRealDenseTags(t *testing.T) {
	checkFixture(t, "pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t", 4932, 13)
}

func TestParseRealGapSegments(t *testing.T) {
	checkFixture(t, "pDGB_A1_act15p_DDB_G0278295_mNeonGreen_act8t", 4458, 11)
}

func checkFixture(t *testing.T, name string, size, count int) {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name + ".dna")
	require.NoError(t, err)
	got := parse(t, data, name)
	require.Empty(t, got.Warnings)
	require.Len(t, got.Record.Sequence, size)
	require.Equal(t, strings.ToUpper(got.Record.Sequence), got.Record.Sequence)
	require.True(t, got.Record.Meta.Locus.Circular)
	require.Len(t, got.Record.Features, count+1)
	require.Equal(t, "primer_bind", got.Record.Features[count].Type)
	packets, err := E.UnwrapError(snapgene.Packets(data))
	require.NoError(t, err)
	for _, packet := range packets {
		if packet.ID == 0 {
			require.Equal(t, byte(15), packet.Body[0])
		}
		if packet.ID != 10 {
			continue
		}
		var root struct {
			Features []struct {
				Segments []struct {
					Range string `xml:"range,attr"`
					Type  string `xml:"type,attr"`
				} `xml:"Segment"`
			} `xml:"Feature"`
		}
		require.NoError(t, xml.Unmarshal(packet.Body, &root))
		for i, f := range root.Features {
			var want [][2]int
			for _, segment := range f.Segments {
				if segment.Type == "gap" {
					continue
				}
				var start, end int
				_, err := fmt.Sscanf(segment.Range, "%d-%d", &start, &end)
				require.NoError(t, err)
				want = append(want, [2]int{start - 1, end})
			}
			require.Equal(t, want, locationLeaves(got.Record.Features[i].Location))
		}
	}
}

func locationLeaves(location genbank.Location) [][2]int {
	if !location.Join {
		return [][2]int{{location.Start, location.End}}
	}
	var out [][2]int
	for _, child := range location.SubLocations {
		out = append(out, locationLeaves(child)...)
	}
	return out
}

func TestParseNonGapSegments(t *testing.T) {
	entry := `<Feature name="other" type="misc_feature"><Segment range="1-3" type="other"/></Feature>`
	got := parse(t, stream(cookie(), dna(15, bases), featurePacket(entry)), "other")
	require.Empty(t, got.Warnings)
	require.Equal(t, genbank.Location{Start: 0, End: 3}, got.Record.Features[0].Location)
}

func TestParseRemovesTrademarksFromFullComment(t *testing.T) {
	notes := `<Notes><Created>date™</Created><LastModified>later™</LastModified></Notes>`
	got := parse(t, stream(cookie(), dna(15, bases), encodePacket(6, []byte(notes))), "notes")
	require.Equal(t, "created date; modified later.", got.Record.Meta.Other["COMMENT"])
}

func TestParseOmitsEmptyCleanedCommentParts(t *testing.T) {
	notes := `<Notes><UUID>™</UUID><Created>™</Created>` +
		`<LastModified>™</LastModified><TransformedInto>™</TransformedInto></Notes>`
	got := parse(t, stream(cookie(), dna(15, bases), encodePacket(6, []byte(notes))), "empty")
	require.Empty(t, got.Warnings)
	require.Empty(t, got.Record.Meta.Other)
}
