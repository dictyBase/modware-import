package gpadstats

import (
	"context"
	"io"
	"os"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioeither/file"
	"github.com/nao1215/filesql"
)

func closeConfig(config StatsLoaderConfig) {
	if config.DB != nil {
		_ = config.DB.Close()
	}
	if config.Reader != nil {
		if c, ok := config.Reader.(io.Closer); ok {
			_ = c.Close()
		}
	}
}

func builder(config StatsLoaderConfig) IOE.IOEither[error, StatsLoaderConfig] {
	return IOE.TryCatchError(func() (StatsLoaderConfig, error) {
		ctx := context.Background()
		validatedBuilder, err := filesql.NewBuilder().
			AddReader(config.Reader, "gpad", filesql.FileTypeTSV).
			Build(ctx)
		if err != nil {
			closeConfig(config)
			return config, err
		}
		db, err := validatedBuilder.Open(ctx)
		if err != nil {
			closeConfig(config)
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
			return StatsLoaderConfig{Reader: f}
		}),
	)
}

func openUrlResources(url string) IOE.IOEither[error, StatsLoaderConfig] {
	return F.Pipe3(
		httpGet(url),
		IOE.Chain(gzipReader),
		IOE.Chain(transformStream),
		IOE.Map[error](func(r io.Reader) StatsLoaderConfig {
			return StatsLoaderConfig{Reader: r}
		}),
	)
}

func releaseResources(
	config StatsLoaderConfig,
	_ E.Either[error, GeneCountStats],
) IOE.IOEither[error, any] {
	return IOE.TryCatchError(func() (any, error) {
		closeConfig(config)
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

func builderFromUrl(url string) IOE.IOEither[error, StatsLoaderConfig] {
	return F.Pipe1(
		openUrlResources(url),
		IOE.Chain(builder),
	)
}
