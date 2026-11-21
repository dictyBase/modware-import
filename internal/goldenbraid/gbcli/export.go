package gbcli

import (
	"context"
	"log/slog"
	"net"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	ioeutils "github.com/dictyBase/fp-go-loom/ioeitherutils"
	"github.com/urfave/cli/v2"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/goldenbraid/gbexport"
	"github.com/dictyBase/modware-import/internal/logger"
)

const s3Scheme = "http://"

var (
	stateL     = gbexport.MakeStateLenses()
	paramsL    = gbexport.MakeParamsLenses()
	flagString = (*cli.Context).String
	paramEdits = []func(*cli.Context) En.Endomorphism[gbexport.Params]{
		F.Flow2(F.Bind2nd(flagString, "input"), paramsL.Input.Set),
		F.Flow2(F.Bind2nd(flagString, "input-type"), paramsL.InputType.Set),
		F.Flow2(F.Bind2nd(flagString, "output"), paramsL.Output.Set),
		F.Flow2(s3Endpoint, paramsL.S3Endpoint.Set),
		F.Flow2(F.Bind2nd(flagString, "s3-region"), paramsL.S3Region.Set),
		F.Flow2(F.Bind2nd(flagString, "access-key"), paramsL.AccessKey.Set),
		F.Flow2(F.Bind2nd(flagString, "secret-key"), paramsL.SecretKey.Set),
		F.Flow2(F.Bind2nd((*cli.Context).Bool, "show-primers"), paramsL.ShowPrimers.Set),
	}
	stateEdits = []func(*cli.Context) En.Endomorphism[gbexport.State]{
		F.Flow2(cliContext, stateL.Ctx.Set),
		F.Flow3(logger.GetCliSlogHandler, slog.New, stateL.Logger.Set),
		F.Flow2(paramsOf, stateL.Params.Set),
	}
)

// ExportFlags returns the export flags plus the shared S3 and logging flags.
func ExportFlags() []cli.Flag {
	return A.ArrayConcatAll(exportFlags(), stockcenter.S3Flags(), stockcenter.LoggingFlags())
}

// RunExport runs the export pipeline for one CLI invocation.
func RunExport(c *cli.Context) error {
	return F.Pipe4(
		c,
		toExportState,
		gbexport.Pipeline,
		ioeutils.ToEither[error, gbexport.State],
		E.ToError[gbexport.State],
	)
}

func toExportState(c *cli.Context) gbexport.State {
	return F.Pipe2(
		stateEdits,
		A.Flap[En.Endomorphism[gbexport.State]](c),
		En.Reduce[gbexport.State],
	)
}

func paramsOf(c *cli.Context) gbexport.Params {
	return F.Pipe2(
		paramEdits,
		A.Flap[En.Endomorphism[gbexport.Params]](c),
		En.Reduce[gbexport.Params],
	)
}

func cliContext(c *cli.Context) context.Context { return c.Context }

// s3Endpoint always uses HTTP: http://<s3-server>:<s3-server-port>.
func s3Endpoint(c *cli.Context) string {
	host, port := c.String("s3-server"), c.String("s3-server-port")
	return s3Scheme + net.JoinHostPort(host, port)
}
