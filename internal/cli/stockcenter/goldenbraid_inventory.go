package stockcenter

import (
	"log/slog"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/urfave/cli/v2"
)

func getSlogHandler(c *cli.Context) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if c.String("log-level") == "debug" {
		opts.Level = slog.LevelDebug
	}
	if c.String("log-format") == "text" {
		return slog.NewTextHandler(os.Stderr, opts)
	}
	return slog.NewJSONHandler(os.Stderr, opts)
}

func LoadGoldenBraidInventory(cmd *cli.Context) error {
	handler := getSlogHandler(cmd)
	slogger := slog.New(handler)
	elog := E.Logger[error, stockcenter.InventoryProcessingSummary](
		slog.NewLogLogger(handler, slog.LevelInfo),
		slog.NewLogLogger(handler, slog.LevelError),
	)

	return F.Pipe2(
		IOE.Bracket(
			stockcenter.InventoryBuilderFromFile(cmd, slogger),
			stockcenter.ProcessInventory,
			stockcenter.ReleaseInventoryResources,
		),
		fputil.ToEither[error, stockcenter.InventoryProcessingSummary],
		E.Fold(
			func(err error) error {
				_ = elog(
					"GoldenBraid inventory loading error",
				)(
					E.Left[stockcenter.InventoryProcessingSummary](err),
				)
				return err
			},
			func(summary stockcenter.InventoryProcessingSummary) error {
				_ = elog("GoldenBraid inventory loading success")(E.Right[error](summary))
				slogger.Info(
					"GoldenBraid inventory loading summary",
					"successes", summary.SuccessCount,
					"errors", summary.ErrorCount,
				)
				return nil
			},
		),
	)
}
