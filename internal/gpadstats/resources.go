package gpadstats

import (
	"context"
	"os"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioeither/file"
	"github.com/nao1215/filesql"
)

func builder(config StatsLoaderConfig) IOE.IOEither[error, StatsLoaderConfig] {
	return IOE.TryCatchError(func() (StatsLoaderConfig, error) {
		ctx := context.Background()
		f := config.File
		validatedBuilder, err := filesql.NewBuilder().
			AddReader(f, "gpad", filesql.FileTypeTSV).
			Build(ctx)
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

func releaseResources(
	config StatsLoaderConfig,
	_ E.Either[error, GeneCountStats],
) IOE.IOEither[error, any] {
	return IOE.TryCatchError(func() (any, error) {
		config.DB.Close()
		config.File.Close()
		return nil, nil
	})
}

func builderFromFile(file string) IOE.IOEither[error, StatsLoaderConfig] {
	return F.Pipe2(
		IOE.Do[error](StatsLoaderConfig{
			Path: file,
		}),
		IOE.Chain(openResources),
		IOE.Chain(builder),
	)
}
