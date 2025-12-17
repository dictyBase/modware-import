package stockcenter

import (
	"log/slog"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/urfave/cli/v2"
)

type inventoryLoadResult struct {
	Summary stockcenter.InventoryProcessingSummary
	Error   error
}

func onInventoryError(err error) inventoryLoadResult {
	return inventoryLoadResult{Error: err}
}

func onInventorySuccess(
	summary stockcenter.InventoryProcessingSummary,
) inventoryLoadResult {
	return inventoryLoadResult{Summary: summary}
}

func handleInventoryOutput(
	result inventoryLoadResult,
	slogger *slog.Logger,
) error {
	if result.Error != nil {
		return result.Error
	}

	slogger.Info(
		"GoldenBraid inventory loading summary",
		"successes", result.Summary.SuccessCount,
		"errors", result.Summary.ErrorCount,
	)
	return nil
}

func LoadGoldenBraidInventory(cmd *cli.Context) error {
	handler := logger.GetCliSlogHandler(cmd)
	slogger := slog.New(handler)

	output := F.Pipe2(
		IOE.Bracket(
			stockcenter.InventoryBuilderFromFile(cmd, slogger),
			stockcenter.ProcessInventory,
			stockcenter.ReleaseInventoryResources,
		),
		fputil.ToEither[error, stockcenter.InventoryProcessingSummary],
		E.Fold(onInventoryError, onInventorySuccess),
	)

	return handleInventoryOutput(output, slogger)
}
