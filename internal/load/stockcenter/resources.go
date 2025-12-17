package stockcenter

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/nao1215/filesql"
	"github.com/urfave/cli/v2"
)

type InventoryLoaderConfig struct {
	DB     *sql.DB
	Reader io.ReadCloser
	Logger *slog.Logger
	Cmd    *cli.Context
}

func closeInventoryConfig(config InventoryLoaderConfig) {
	if config.DB != nil {
		_ = config.DB.Close()
	}
	if config.Reader != nil {
		_ = config.Reader.Close()
	}
}

func inventoryBuilder(config InventoryLoaderConfig) IOE.IOEither[error, InventoryLoaderConfig] {
	return IOE.TryCatchError(func() (InventoryLoaderConfig, error) {
		ctx := context.Background()
		// Initialize filesql builder
		// Ensure filesql handles the CSV headers correctly
		validatedBuilder, err := filesql.NewBuilder().
			AddReader(config.Reader, "inventory", filesql.FileTypeCSV).
			Build(ctx)
		if err != nil {
			closeInventoryConfig(config)
			return config, err
		}
		db, err := validatedBuilder.Open(ctx)
		if err != nil {
			closeInventoryConfig(config)
			return config, err
		}
		config.DB = db
		return config, nil
	})
}

func openInventoryReader(config InventoryLoaderConfig) IOE.IOEither[error, InventoryLoaderConfig] {
	return IOE.TryCatchError(func() (InventoryLoaderConfig, error) {
		source := config.Cmd.String("input-source")
		switch source {
		case "folder":
			return openInventoryFileReader(config)
		case "bucket":
			return openInventoryS3Reader(config)
		default:
			return config, fmt.Errorf("unsupported input source %s", source)
		}
	})
}

func openInventoryFileReader(config InventoryLoaderConfig) (InventoryLoaderConfig, error) {
	inputPath := config.Cmd.String("input")
	file, err := os.Open(inputPath)
	if err != nil {
		return config, fmt.Errorf("failed to open CSV file %s: %w", inputPath, err)
	}
	config.Reader = file
	return config, nil
}

func openInventoryS3Reader(config InventoryLoaderConfig) (InventoryLoaderConfig, error) {
	bucket := config.Cmd.String("s3-bucket")
	bucketPath := config.Cmd.String("s3-bucket-path")
	file := config.Cmd.String("input")
	reader, err := registry.GetS3Client().GetObject(
		bucket,
		fmt.Sprintf("%s/%s", bucketPath, file),
		minio.GetObjectOptions{},
	)
	if err != nil {
		return config, fmt.Errorf("failed to open s3 file %s/%s: %w", bucketPath, file, err)
	}
	config.Reader = reader
	return config, nil
}

func ReleaseInventoryResources(
	config InventoryLoaderConfig,
	_ E.Either[error, InventoryProcessingSummary],
) IOE.IOEither[error, any] {
	return IOE.TryCatchError(func() (any, error) {
		closeInventoryConfig(config)
		return nil, nil
	})
}

func InventoryBuilderFromFile(
	cmd *cli.Context,
	logger *slog.Logger,
) IOE.IOEither[error, InventoryLoaderConfig] {
	return F.Pipe2(
		IOE.Do[error](InventoryLoaderConfig{
			Cmd:    cmd,
			Logger: logger,
		}),
		IOE.Chain(openInventoryReader),
		IOE.Chain(inventoryBuilder),
	)
}
