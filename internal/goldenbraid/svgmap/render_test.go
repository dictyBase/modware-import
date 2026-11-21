package svgmap_test

import (
	"encoding/xml"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/bebop/poly/io/genbank"
	"github.com/stretchr/testify/require"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
	"github.com/dictyBase/modware-import/internal/goldenbraid/svgmap"
)

type svgDocument struct {
	XMLName     xml.Name
	Role        string     `xml:"role,attr"`
	ViewBox     string     `xml:"viewBox,attr"`
	Title       string     `xml:"title"`
	Description string     `xml:"desc"`
	Paths       []svgPath  `xml:"path"`
	Labels      []svgLabel `xml:"text"`
	Leaders     []struct {
		Points string `xml:"points,attr"`
	} `xml:"polyline"`
}
type svgPath struct {
	ID   string `xml:"data-feature-id,attr"`
	Ring int    `xml:"data-ring,attr"`
	Fill string `xml:"fill,attr"`
	D    string `xml:"d,attr"`
}
type svgLabel struct {
	ID     string  `xml:"data-feature-id,attr"`
	Side   string  `xml:"data-side,attr"`
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Length float64 `xml:"textLength,attr"`
	Text   string  `xml:",chardata"`
}

func makeRecord(length int, features ...genbank.Feature) genbank.Genbank {
	locus := F.Pipe1(genbank.Locus{}, PL.LocusCircular.Set(true))
	meta := F.Pipe2(genbank.Meta{}, PL.MetaName.Set("test_construct"), PL.MetaLocus.Set(locus))
	return F.Pipe3(
		genbank.Genbank{},
		PL.GenbankMeta.Set(meta),
		PL.GenbankSequence.Set(strings.Repeat("A", length)),
		PL.GenbankFeatures.Set(features),
	)
}

func feature(label, kind string, start, end int) genbank.Feature {
	location := F.Pipe2(genbank.Location{}, PL.LocationStart.Set(start), PL.LocationEnd.Set(end))
	return F.Pipe3(
		genbank.Feature{},
		PL.FeatureType.Set(kind),
		PL.FeatureLocation.Set(location),
		PL.FeatureAttributes.Set(map[string]string{"label": label}),
	)
}

func render(t *testing.T, record genbank.Genbank, options ...svgmap.Option) (svgDocument, []byte) {
	t.Helper()
	lenses := svgmap.MakeRenderInputLenses()
	input := F.Pipe2(svgmap.RenderInput{}, lenses.Record.Set(record), lenses.Options.Set(options))
	data, err := E.UnwrapError(svgmap.Render(input))
	require.NoError(t, err)
	var doc svgDocument
	require.NoError(t, xml.Unmarshal(data, &doc))
	return doc, data
}

func visibleLabels(doc svgDocument) []svgLabel {
	var labels []svgLabel
	for _, label := range doc.Labels {
		if label.ID != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

func TestRenderProducesAccessibleValidSVG(t *testing.T) {
	record := makeRecord(1000, feature("promoter <&>", "promoter", 10, 100))
	doc, data := render(t, record)
	require.Equal(t, "svg", doc.XMLName.Local)
	require.Equal(t, "http://www.w3.org/2000/svg", doc.XMLName.Space)
	require.Equal(t, "img", doc.Role)
	require.Equal(t, "0 0 640 640", doc.ViewBox)
	require.Contains(t, doc.Title, "test_construct")
	require.Contains(t, doc.Title, "1000 bp")
	require.Len(t, doc.Paths, 1)
	require.Equal(t, "0", doc.Paths[0].ID)
	require.Equal(t, "promoter <&>", visibleLabels(doc)[0].Text)
	require.NotEmpty(t, doc.Leaders)
	_, again := render(t, record)
	require.Equal(t, data, again)
}

func TestAllocateRingsPacksOverlappingFeatures(t *testing.T) {
	doc, _ := render(
		t,
		makeRecord(
			1000,
			feature("outer", "gene", 99, 900),
			feature("inner", "CDS", 99, 850),
			feature("HA", "misc_feature", 850, 900),
		),
	)
	rings := map[string]int{}
	for _, p := range doc.Paths {
		rings[p.ID] = p.Ring
	}
	require.Equal(t, map[string]int{"0": 0, "1": 1, "2": 1}, rings)
}

func TestRenderHandlesCircularOriginFeature(t *testing.T) {
	leaves := []genbank.Location{{Start: 11899, End: 12000}, {Start: 0, End: 150}}
	location := F.Pipe2(
		genbank.Location{},
		PL.LocationJoin.Set(true),
		PL.LocationSubLocations.Set(leaves),
	)
	f := F.Pipe1(feature("origin", "gene", 0, 1), PL.FeatureLocation.Set(location))
	doc, data := render(t, makeRecord(12000, f))
	require.Len(t, doc.Paths, 2)
	require.Equal(t, doc.Paths[0].ID, doc.Paths[1].ID)
	require.Equal(t, doc.Paths[0].Ring, doc.Paths[1].Ring)
	require.NotContains(t, string(data), "NaN")
	require.NotContains(t, string(data), "Inf")
	require.NotEqual(t, doc.Paths[0].D, doc.Paths[1].D)
}

func TestLayoutLabelsSpacesDenseClusterInsideViewBox(t *testing.T) {
	data, err := os.ReadFile("../snapgene/testdata/pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t.dna")
	require.NoError(t, err)
	lenses := snapgene.MakeParseInputLenses()
	input := F.Pipe3(
		snapgene.ParseInput{},
		lenses.Data.Set(data),
		lenses.Name.Set("dense"),
		lenses.SourceFile.Set("dense.dna"),
	)
	parsed, err := E.UnwrapError(snapgene.Parse(input))
	require.NoError(t, err)
	doc, _ := render(t, parsed.Record)
	labels := visibleLabels(doc)
	names := []string{}
	for _, label := range labels {
		names = append(names, label.Text)
	}
	require.Contains(t, names, "GGGSx4 Linker")
	require.Contains(t, names, "H2Bv3")
	assertLabelBounds(t, labels, 640)
}

func assertLabelBounds(t *testing.T, labels []svgLabel, size float64) {
	t.Helper()
	sides := map[string][]float64{}
	for _, label := range labels {
		require.GreaterOrEqual(t, label.Y, 16.0)
		require.LessOrEqual(t, label.Y, size-16)
		require.GreaterOrEqual(t, label.X, 16.0)
		require.LessOrEqual(t, label.X+label.Length, size-16+0.001)
		sides[label.Side] = append(sides[label.Side], label.Y)
	}
	for _, ys := range sides {
		sort.Float64s(ys)
		for i := 1; i < len(ys); i++ {
			require.GreaterOrEqual(t, ys[i]-ys[i-1], 16.0)
		}
	}
}

func TestRenderPrimerVisibilityAndFeatureColors(t *testing.T) {
	colored := F.Pipe1(
		feature("colored", "CDS", 50, 200),
		PL.FeatureAttributes.Set(
			map[string]string{"label": "colored", "ApEinfo_fwdcolor": "#31849b"},
		),
	)
	primer := feature("primer", "primer_bind", 300, 320)
	record := makeRecord(1000, primer, colored, feature("plain", "gene", 400, 500))
	doc, _ := render(t, record)
	require.Len(t, doc.Paths, 2)
	require.Equal(t, "1", doc.Paths[0].ID)
	require.Equal(t, "#31849b", doc.Paths[0].Fill)
	require.NotEmpty(t, doc.Paths[1].Fill)
	require.NotEqual(t, doc.Paths[0].Fill, doc.Paths[1].Fill)
	require.NotContains(t, fmt.Sprint(visibleLabels(doc)), "primer")
	with, _ := render(t, record, svgmap.WithShowPrimers(true))
	require.Len(t, with.Paths, 3)
	require.Equal(t, "0", with.Paths[0].ID)
	require.Contains(t, fmt.Sprint(visibleLabels(with)), "primer")
}

func TestRenderReversesArrowOrientationAndKeepsGaps(t *testing.T) {
	forward := feature("joined", "CDS", 10, 100)
	location := genbank.Location{
		Join:         true,
		SubLocations: []genbank.Location{{Start: 10, End: 100}, {Start: 200, End: 250}},
	}
	forward = F.Pipe1(forward, PL.FeatureLocation.Set(location))
	a, _ := render(t, makeRecord(1000, forward))
	reverse := F.Pipe1(forward, PL.FeatureLocation.Modify(PL.LocationComplement.Set(true)))
	b, _ := render(t, makeRecord(1000, reverse))
	require.Len(t, a.Paths, 2)
	require.Len(t, b.Paths, 2)
	require.NotEqual(t, a.Paths[0].D, b.Paths[0].D)
	require.NotEqual(t, a.Paths[1].D, b.Paths[1].D)
}

func TestRenderBoundsLabelCapacityAndReportsOmissions(t *testing.T) {
	var features []genbank.Feature
	for i := range 80 {
		features = append(features, feature(fmt.Sprintf("label-%d", i), "gene", i*2, i*2+1))
	}
	doc, _ := render(t, makeRecord(1000, features...), svgmap.WithMaxLabels(100))
	require.Len(t, doc.Paths, 80)
	defaultDoc, _ := render(t, makeRecord(1000, features...))
	require.Len(t, visibleLabels(defaultDoc), 32)
	require.Len(t, visibleLabels(doc), 39)
	require.Contains(t, doc.Description, "41")
	assertLabelBounds(t, visibleLabels(doc), 640)
	zero, _ := render(t, makeRecord(1000, features...), svgmap.WithMaxLabels(0))
	require.Empty(t, visibleLabels(zero))
	require.Len(t, zero.Paths, 80)
	require.Contains(t, zero.Description, "80")
	limited, _ := render(t, makeRecord(1000, features...), svgmap.WithMaxLabels(3))
	require.Len(t, visibleLabels(limited), 3)
	require.ElementsMatch(
		t,
		[]string{"0", "1", "2"},
		[]string{
			visibleLabels(limited)[0].ID,
			visibleLabels(limited)[1].ID,
			visibleLabels(limited)[2].ID,
		},
	)
}

func TestRenderOptionsAndInvalidGeometry(t *testing.T) {
	record := makeRecord(1000, feature("x", "gene", 10, 30))
	doc, _ := render(t, record, svgmap.WithSize(800), svgmap.WithSize(480))
	require.Equal(t, "0 0 480 480", doc.ViewBox)
	assertLabelBounds(t, visibleLabels(doc), 480)
	inputL := svgmap.MakeRenderInputLenses()
	invalidOptions := [][]svgmap.Option{
		{
			svgmap.WithSize(32),
		},
		{svgmap.WithSize(33)},
		{svgmap.WithSize(-1)},
		{svgmap.WithMaxLabels(-1)},
	}
	for _, options := range invalidOptions {
		input := F.Pipe2(
			svgmap.RenderInput{},
			inputL.Record.Set(record),
			inputL.Options.Set(options),
		)
		require.ErrorIs(t, E.ToError(svgmap.Render(input)), svgmap.ErrInvalidConfig)
	}
	tiny, _ := render(t, record, svgmap.WithSize(110))
	assertLabelBounds(t, visibleLabels(tiny), 110)
	var overlapping []genbank.Feature
	for i := range 80 {
		overlapping = append(overlapping, feature(fmt.Sprint(i), "gene", 1, 999))
	}
	input := F.Pipe1(svgmap.RenderInput{}, inputL.Record.Set(makeRecord(1000, overlapping...)))
	require.ErrorIs(t, E.ToError(svgmap.Render(input)), svgmap.ErrInvalidConfig)
}

func TestRenderRejectsNonCircularEmptyAndInvalidLocations(t *testing.T) {
	inputL := svgmap.MakeRenderInputLenses()
	record := makeRecord(1000, feature("x", "gene", 0, 1001))
	input := F.Pipe1(svgmap.RenderInput{}, inputL.Record.Set(record))
	require.ErrorIs(t, E.ToError(svgmap.Render(input)), seqmodel.ErrSegmentOutOfBounds)
	linear := makeRecord(1000)
	linear.Meta.Locus.Circular = false
	input = inputL.Record.Set(linear)(input)
	require.ErrorIs(t, E.ToError(svgmap.Render(input)), svgmap.ErrNotCircular)
	input = inputL.Record.Set(makeRecord(0))(input)
	require.ErrorIs(t, E.ToError(svgmap.Render(input)), seqmodel.ErrEmptySequence)
	emptyName := makeRecord(1000)
	emptyName.Meta.Name = ""
	input = inputL.Record.Set(emptyName)(input)
	require.ErrorIs(t, E.ToError(svgmap.Render(input)), seqmodel.ErrEmptyName)
}

// TestWriteVisualFixtures checks all real maps and optionally retains review artifacts.
func TestWriteVisualFixtures(t *testing.T) {
	output := t.TempDir()
	if os.Getenv("GOLDENBRAID_SVG_REVIEW") == "1" {
		output = "../../../.superpowers/sdd/2026-11-20-goldenbraid-export-cli-plan/visual"
		require.NoError(t, os.MkdirAll(output, 0o750))
	}
	files, err := filepath.Glob("../snapgene/testdata/*.dna")
	require.NoError(t, err)
	require.Len(t, files, 3)
	t.Logf("visual artifacts: %s", output)
	for _, name := range files {
		data, readErr := os.ReadFile(name)
		require.NoError(t, readErr)
		stem := strings.TrimSuffix(filepath.Base(name), ".dna")
		lenses := snapgene.MakeParseInputLenses()
		input := F.Pipe3(
			snapgene.ParseInput{},
			lenses.Data.Set(data),
			lenses.Name.Set(stem),
			lenses.SourceFile.Set(name),
		)
		parsed, parseErr := E.UnwrapError(snapgene.Parse(input))
		require.NoError(t, parseErr)
		doc, svg := render(t, parsed.Record)
		assertLabelBounds(t, visibleLabels(doc), 640)
		require.NotContains(t, string(svg), "NaN")
		require.NotContains(t, string(svg), "Inf")
		require.NoError(t, os.WriteFile(filepath.Join(output, stem+".svg"), svg, 0o600))
	}
}

func TestLabelLeaderAnchorsOnOccupiedJoinedSegment(t *testing.T) {
	location := genbank.Location{
		Join:         true,
		SubLocations: []genbank.Location{{Start: 10, End: 50}, {Start: 400, End: 500}},
	}
	joined := F.Pipe1(feature("gapped", "CDS", 10, 50), PL.FeatureLocation.Set(location))
	doc, _ := render(t, makeRecord(1000, joined))
	require.Len(t, doc.Leaders, 1)
	var x, y float64
	_, err := fmt.Sscanf(doc.Leaders[0].Points, "%f,%f", &x, &y)
	require.NoError(t, err)
	// Half of 140 occupied bases is 70: 40 in the first leaf, then 30 at base 430.
	angle := math.Atan2(y-320, x-320) + math.Pi/2
	require.InDelta(t, 2*math.Pi*0.43, angle, 0.00001)
}
