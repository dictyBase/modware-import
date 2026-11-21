package gbexport_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/require"

	"github.com/dictyBase/modware-import/internal/goldenbraid/gbexport"
	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
)

const (
	simpleStem = "pDGB_A1_act15p_mNeonGreen_act8t"
	denseStem  = "pDGB_A1_act15p_GST_H2Bv3_mNeonGreen_act8t"
	gapStem    = "pDGB_A1_act15p_DDB_G0278295_mNeonGreen_act8t"
)

var allStems = []string{simpleStem, denseStem, gapStem}

type logRecord map[string]any

type capture struct {
	buf bytes.Buffer
}

func (c *capture) logger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(&c.buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func (c *capture) records(t *testing.T) []logRecord {
	t.Helper()
	var out []logRecord
	for line := range strings.SplitSeq(strings.TrimSpace(c.buf.String()), "\n") {
		if line == "" {
			continue
		}
		var rec logRecord
		require.NoError(t, json.Unmarshal([]byte(line), &rec))
		out = append(out, rec)
	}
	return out
}

func (c *capture) withMessage(t *testing.T, msg string) []logRecord {
	t.Helper()
	var out []logRecord
	for _, rec := range c.records(t) {
		if rec["msg"] == msg {
			out = append(out, rec)
		}
	}
	return out
}

func fixture(t *testing.T, stem string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "snapgene", "testdata", stem+".dna"))
	require.NoError(t, err)
	return data
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, data, 0o600))
}

func run(
	ctx context.Context,
	t *testing.T,
	logs *capture,
	params gbexport.Params,
) (gbexport.State, error) {
	t.Helper()
	state := gbexport.State{Ctx: ctx, Logger: logs.logger(), Params: params}
	return E.UnwrapError(gbexport.Pipeline(state)())
}

func artifactPaths(out, stem string) []string {
	return []string{
		filepath.Join(out, gbexport.FastaDir, stem+".fasta"),
		filepath.Join(out, gbexport.GenBankDir, stem+".gbk"),
		filepath.Join(out, gbexport.SVGDir, stem+".svg"),
	}
}

func requireArtifacts(t *testing.T, out string, stems ...string) {
	t.Helper()
	for _, stem := range stems {
		for _, path := range artifactPaths(out, stem) {
			info, err := os.Stat(path)
			require.NoError(t, err, path)
			require.Positive(t, info.Size(), path)
		}
	}
}

func requireNoArtifacts(t *testing.T, out string, stems ...string) {
	t.Helper()
	for _, stem := range stems {
		for _, path := range artifactPaths(out, stem) {
			_, err := os.Stat(path)
			require.ErrorIs(t, err, os.ErrNotExist, path)
		}
	}
}

func requireSummary(t *testing.T, logs *capture, total, succeeded, failed float64) {
	t.Helper()
	summaries := logs.withMessage(t, "export summary")
	require.Len(t, summaries, 1)
	require.InDelta(t, total, summaries[0]["total"], 0)
	require.InDelta(t, succeeded, summaries[0]["succeeded"], 0)
	require.InDelta(t, failed, summaries[0]["failed"], 0)
}

func packet(id byte, body []byte) []byte {
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

func syntheticWithInvalidFeature() []byte {
	cookie := packet(9, append([]byte("SnapGene"), 0, 1, 0, 13, 0, 14))
	dna := packet(0, append([]byte{15}, []byte("GGAGCAAAATTTTGTATTTT")...))
	features := packet(10, []byte(
		`<Features>`+
			`<Feature name="ok" type="CDS"><Segment range="2-8" type="standard"/></Feature>`+
			`<Feature name="bad" type="CDS"><Segment range="1-21" type="standard"/></Feature>`+
			`</Features>`,
	))
	return bytes.Join([][]byte{cookie, dna, features}, nil)
}

func zipArchive(t *testing.T, entries map[string][]byte, dirs ...string) []byte {
	t.Helper()
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	for _, dir := range dirs {
		_, err := writer.Create(dir)
		require.NoError(t, err)
	}
	for name, data := range entries {
		entry, err := writer.Create(name)
		require.NoError(t, err)
		_, err = entry.Write(data)
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())
	return buf.Bytes()
}

func assemblyArchive(t *testing.T) []byte {
	t.Helper()
	entries := map[string][]byte{
		"__MACOSX/GoldenBraid Assemblies/._x.dna": []byte("apple double"),
		"GoldenBraid Assemblies/deep/nested.dna":  fixture(t, simpleStem),
	}
	for _, stem := range allStems {
		entries["GoldenBraid Assemblies/"+stem+".dna"] = fixture(t, stem)
	}
	return zipArchive(t, entries, "GoldenBraid Assemblies/")
}

func TestPipelineExportsSingleConstruct(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, simpleStem+".dna"), fixture(t, simpleStem))
	writeFile(t, filepath.Join(in, denseStem+".dna"), fixture(t, denseStem))
	logs := &capture{}
	state, err := run(context.Background(), t, logs, gbexport.Params{
		Input:  filepath.Join(in, simpleStem+".dna"),
		Output: out,
	})
	require.NoError(t, err)
	requireArtifacts(t, out, simpleStem)
	requireNoArtifacts(t, out, denseStem)
	requireSummary(t, logs, 1, 1, 0)
	require.Equal(t, 1, state.Summary.Succeeded)
	fasta, err := os.ReadFile(artifactPaths(out, simpleStem)[0])
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(fasta), ">"+simpleStem+"\n"))
	gbk, err := os.ReadFile(artifactPaths(out, simpleStem)[1])
	require.NoError(t, err)
	require.Contains(t, string(gbk), "primer_bind")
}

func TestPipelineContinuesAfterCorruptFile(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, simpleStem+".dna"), fixture(t, simpleStem))
	writeFile(t, filepath.Join(in, denseStem+".dna"), fixture(t, denseStem))
	writeFile(t, filepath.Join(in, "corrupt.dna"), []byte("not a snapgene file"))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.ErrorIs(t, err, gbexport.ErrPartialExportFailure)
	requireArtifacts(t, out, simpleStem, denseStem)
	requireNoArtifacts(t, out, "corrupt")
	requireSummary(t, logs, 3, 2, 1)
	failures := logs.withMessage(t, "export failed")
	require.Len(t, failures, 1)
	require.Equal(t, "corrupt.dna", failures[0]["key"])
	require.NotEmpty(t, failures[0]["error"])
}

func TestPipelineLogsParserWarnings(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, "warned.dna"), syntheticWithInvalidFeature())
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.NoError(t, err)
	warnings := logs.withMessage(t, "parse warning")
	require.Len(t, warnings, 1)
	require.Equal(t, "warned.dna", warnings[0]["key"])
	requireSummary(t, logs, 1, 1, 0)
	requireArtifacts(t, out, "warned")
}

func TestPipelineKeepsEarlierArtifactsOnLaterWriteFailure(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, simpleStem+".dna"), fixture(t, simpleStem))
	// A directory at the SVG key makes the final write fail.
	require.NoError(t, os.MkdirAll(artifactPaths(out, simpleStem)[2], 0o750))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.ErrorIs(t, err, gbexport.ErrPartialExportFailure)
	_, err = os.Stat(artifactPaths(out, simpleStem)[0])
	require.NoError(t, err)
	requireSummary(t, logs, 1, 0, 1)
}

func TestPipelineLimitsDirectoryDepth(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, simpleStem+".dna"), fixture(t, simpleStem))
	writeFile(
		t,
		filepath.Join(in, "GoldenBraid Assemblies", denseStem+".dna"),
		fixture(t, denseStem),
	)
	writeFile(t, filepath.Join(in, "a", "b", gapStem+".dna"), fixture(t, gapStem))
	writeFile(t, filepath.Join(in, "._hidden.dna"), []byte("apple double"))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.NoError(t, err)
	requireArtifacts(t, out, simpleStem, denseStem)
	requireNoArtifacts(t, out, gapStem, "._hidden")
	requireSummary(t, logs, 2, 2, 0)
}

func TestPipelineExportsAssemblyArchive(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, "assemblies.zip"), assemblyArchive(t))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{
		Input:  filepath.Join(in, "assemblies.zip"),
		Output: out,
	})
	require.NoError(t, err)
	requireArtifacts(t, out, allStems...)
	requireNoArtifacts(t, out, "nested", "._x")
	requireSummary(t, logs, 3, 3, 0)
}

func TestPipelineExportsArchiveRootEntry(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	archive := zipArchive(t, map[string][]byte{simpleStem + ".dna": fixture(t, simpleStem)})
	writeFile(t, filepath.Join(in, "root.zip"), archive)
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{
		Input:  filepath.Join(in, "root.zip"),
		Output: out,
	})
	require.NoError(t, err)
	requireArtifacts(t, out, simpleStem)
}

func TestPipelineExplicitInputTypeOverridesExtension(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, "bundle.bin"), assemblyArchive(t))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{
		Input:     filepath.Join(in, "bundle.bin"),
		InputType: "zip",
		Output:    out,
	})
	require.NoError(t, err)
	requireArtifacts(t, out, allStems...)
}

func TestPipelineRejectsEmptyInput(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.ErrorIs(t, err, gbexport.ErrNoDNAFiles)
	entries, readErr := os.ReadDir(out)
	require.NoError(t, readErr)
	require.Empty(t, entries)
}

func TestPipelineRejectsDuplicateStem(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, "a", "same.dna"), fixture(t, simpleStem))
	writeFile(t, filepath.Join(in, "b", "same.dna"), fixture(t, denseStem))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.ErrorIs(t, err, gbexport.ErrPartialExportFailure)
	require.ErrorIs(t, err, gbexport.ErrDuplicateStem)
	requireSummary(t, logs, 2, 1, 1)
	fasta, readErr := os.ReadFile(filepath.Join(out, gbexport.FastaDir, "same.fasta"))
	require.NoError(t, readErr)
	expected, readErr := os.ReadFile(filepath.Join(in, "a", "same.dna"))
	require.NoError(t, readErr)
	require.NotEmpty(t, expected)
	require.True(t, strings.HasPrefix(string(fasta), ">same\n"))
	gbk, readErr := os.ReadFile(filepath.Join(out, gbexport.GenBankDir, "same.gbk"))
	require.NoError(t, readErr)
	require.NotContains(t, string(gbk), "H2Bv3", "first source must win")
}

func TestPipelineReplacesStaleOutput(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, simpleStem+".dna"), fixture(t, simpleStem))
	stale := artifactPaths(out, simpleStem)[1]
	writeFile(t, stale, []byte("stale"))
	logs := &capture{}
	_, err := run(context.Background(), t, logs, gbexport.Params{Input: in, Output: out})
	require.NoError(t, err)
	data, err := os.ReadFile(stale)
	require.NoError(t, err)
	require.NotEqual(t, "stale", string(data))
	require.Contains(t, string(data), "LOCUS")
}

func TestPipelineReportsCancellationPerFile(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(in, simpleStem+".dna"), fixture(t, simpleStem))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	logs := &capture{}
	_, err := run(ctx, t, logs, gbexport.Params{
		Input:  filepath.Join(in, simpleStem+".dna"),
		Output: out,
	})
	require.ErrorIs(t, err, gbexport.ErrPartialExportFailure)
	require.ErrorIs(t, err, context.Canceled)
	requireSummary(t, logs, 1, 0, 1)
}

func TestPipelineRejectsUnknownTypeAndScheme(t *testing.T) {
	in, out := t.TempDir(), t.TempDir()
	cases := []struct {
		name   string
		params gbexport.Params
		want   error
	}{
		{
			"type",
			gbexport.Params{Input: in, InputType: "tarball", Output: out},
			gbexport.ErrUnknownInputType,
		},
		{
			"input scheme",
			gbexport.Params{Input: "https://example.org/a.dna", Output: out},
			gbexport.ErrUnsupportedScheme,
		},
		{
			"output scheme",
			gbexport.Params{Input: in, Output: "https://example.org/out"},
			gbexport.ErrUnsupportedScheme,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := run(context.Background(), t, &capture{}, tc.params)
			require.True(t, errors.Is(err, tc.want), "got %v", err)
		})
	}
}
