package gpadstats

import (
	"fmt"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/urfave/cli/v2"
)

func toEither[ERR, A any](io IOE.IOEither[ERR, A]) E.Either[ERR, A] {
	return io()
}

func onStatsError(err error) T.Tuple2[GeneCountStats, error] {
	return T.MakeTuple2(GeneCountStats{}, err)
}

func onStatsSuccess(stats GeneCountStats) T.Tuple2[GeneCountStats, error] {
	return T.MakeTuple2(stats, (error)(nil))
}

func Run(cltx *cli.Context) error {
	output := F.Pipe2(
		IOE.Bracket(
			builderFromFile(cltx.String("file")),
			queryCounts,
			releaseResources,
		),
		toEither[error, GeneCountStats],
		E.Fold(
			onStatsError,
			onStatsSuccess,
		),
	)

	if output.F2 != nil {
		return cli.Exit(fmt.Sprintf("Error: %v", output.F2), 1)
	}
	fmt.Printf("Unique Gene Count: %d\n", output.F1.Count)
	fmt.Printf("Unique Gene with ECO code Count: %d\n", output.F1.EcoCount)
	return nil
}
