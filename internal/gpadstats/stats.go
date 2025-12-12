package gpadstats

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioeither/file"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/nao1215/filesql"
	"github.com/urfave/cli/v2"
)

type StatsLoaderConfig struct {
	Path string
	File *os.File
	DB   *sql.DB
}

func builder(config StatsLoaderConfig) IOE.IOEither[error, StatsLoaderConfig] {
	return IOE.TryCatchError(func() (StatsLoaderConfig, error) {
		ctx := context.Background()
		f := config.File
		builder := filesql.NewBuilder().
			AddReader(f, "gpad", filesql.FileTypeTSV)
		validatedBuilder, err := builder.Build(ctx)
		if err != nil {
			_ = f.Close()
			return config, err
		}
		db, err := validatedBuilder.Open(ctx)
		if err != nil {
			_ = f.Close()
			return config, err
		}
		config.DB = db
		return config, nil
	})
}

func openResources(
	config StatsLoaderConfig,
) IOE.IOEither[error, StatsLoaderConfig] {
	return F.Pipe1(
		file.Open(config.Path),
		IOE.Map[error](func(f *os.File) StatsLoaderConfig {
			return StatsLoaderConfig{File: f}
		}),
	)
}

func computeStats(config StatsLoaderConfig) IOE.IOEither[error, int] {
	return IOE.TryCatchError(func() (int, error) {
		var count int
		err := config.DB.QueryRow("SELECT COUNT(DISTINCT DB_Object_ID) FROM gpad").
			Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("failed to execute count query: %w", err)
		}
		return count, nil
	})
}

func releaseResources(
	config StatsLoaderConfig,
	_ E.Either[error, int],
) IOE.IOEither[error, any] {
	return IOE.TryCatchError(func() (any, error) {
		config.DB.Close()
		config.File.Close()
		return nil, nil
	})
}

func toEither[ERR, A any](io IOE.IOEither[ERR, A]) E.Either[ERR, A] {
	return io()
}

func Run(cltx *cli.Context) error {
	output := F.Pipe2(
		IOE.Bracket(
			F.Pipe2(
				IOE.Do[error](StatsLoaderConfig{
					Path: cltx.String("file"),
				}),
				IOE.Chain(openResources),
				IOE.Chain(builder),
			),
			computeStats,
			releaseResources,
		),
		toEither[error, int],
		E.Fold(
			F.Bind1st(T.MakeTuple2[int, error], 0),
			F.Bind2nd(T.MakeTuple2[int, error], nil),
		),
	)

	if output.F2 != nil {
		return cli.Exit(fmt.Sprintf("Error: %v", output.F2), 1)
	}
	fmt.Printf("Unique Gene Count: %d\n", output.F1)
	return nil
}
