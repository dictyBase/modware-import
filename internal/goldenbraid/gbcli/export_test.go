package gbcli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/dictyBase/modware-import/internal/goldenbraid/gbcli"
	"github.com/dictyBase/modware-import/internal/goldenbraid/gbexport"
)

const stem = "pDGB_A1_act15p_mNeonGreen_act8t"

func flagsByName(flags []cli.Flag) map[string]cli.Flag {
	out := make(map[string]cli.Flag, len(flags))
	for _, flag := range flags {
		out[flag.Names()[0]] = flag
	}
	return out
}

func TestExportFlagsMatchesExpectedFlagSet(t *testing.T) {
	flags := flagsByName(gbcli.ExportFlags())
	names := make([]string, 0, len(flags))
	for name := range flags {
		names = append(names, name)
	}
	require.ElementsMatch(t, []string{
		"input", "output", "show-primers", "s3-region",
		"s3-server", "s3-server-port", "access-key", "secret-key",
		"log-level", "log-format",
	}, names)
	require.True(t, flags["input"].(*cli.StringFlag).Required)
	require.True(t, flags["output"].(*cli.StringFlag).Required)
	require.Equal(t, "us-east-1", flags["s3-region"].(*cli.StringFlag).Value)
	require.False(t, flags["show-primers"].(*cli.BoolFlag).Value)
}

func TestRunExportLocalFolderWritesArtifacts(t *testing.T) {
	out := t.TempDir()
	app := &cli.App{Name: "goldenbraid", Flags: gbcli.ExportFlags(), Action: gbcli.RunExport}
	require.NoError(t, app.Run([]string{
		"goldenbraid",
		"--input", filepath.Join("..", "snapgene", "testdata"),
		"--output", out,
	}))
	for _, path := range []string{
		filepath.Join(out, gbexport.FastaDir, stem+".fasta"),
		filepath.Join(out, gbexport.GenBankDir, stem+".gbk"),
		filepath.Join(out, gbexport.SVGDir, stem+".svg"),
	} {
		_, err := os.Stat(path)
		require.NoError(t, err, path)
	}
}

func TestRunExportReturnsPipelineError(t *testing.T) {
	app := &cli.App{Name: "goldenbraid", Flags: gbcli.ExportFlags(), Action: gbcli.RunExport}
	err := app.Run([]string{"goldenbraid", "--input", t.TempDir(), "--output", t.TempDir()})
	require.ErrorIs(t, err, gbexport.ErrNoDNAFiles)
}
