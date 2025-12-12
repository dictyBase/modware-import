package gpadstats

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/IBM/fp-go/ioeither/file"
	fputil "github.com/dictyBase/modware-import/internal/fputil"
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

func Run(cltx *cli.Context) error {
	return F.Pipe2(
		IOE.Bracket(
			F.Pipe2(IOE.Do[error](StatsLoaderConfig{
				Path: cltx.String("file"),
			}),
				IOE.Chain(openResources),
				IOE.Chain(builder),
			),
			computeStats,
			releaseResources,
		),
		fputil.ToEither[error, int],
		E.Fold(
			func(err error) error {
				return cli.Exit(fmt.Sprintf("Error: %v", err), 1)
			},
			func(count int) error {
				fmt.Printf("Unique Gene Count: %d\n", count)
				return nil
			},
		),
	)
}
