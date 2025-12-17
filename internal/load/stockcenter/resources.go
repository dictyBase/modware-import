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

type InventoryLoaderIOE = IOE.IOEither[error, InventoryLoaderConfig]

func closeInventoryConfig(config InventoryLoaderConfig) {
	if config.DB != nil {
		_ = config.DB.Close()
	}
	if config.Reader != nil {
		_ = config.Reader.Close()
	}
}

func inventoryBuilder(
	config InventoryLoaderConfig,
) InventoryLoaderIOE {
	ctx := context.Background()
	return F.Pipe2(
		IOE.TryCatchError(func() (*sql.DB, error) {
			validatedBuilder, err := filesql.NewBuilder().
				AddReader(
					config.Reader,
					"inventory",
					filesql.FileTypeCSV,
				).
				Build(ctx)
			if err != nil {
				return nil, err
			}
			return validatedBuilder.Open(ctx)
		}),
		IOE.MapLeft[*sql.DB](func(err error) error {
			return fmt.Errorf(
				"failed to build inventory database: %w",
				err,
			)
		}),
		IOE.Map[error](func(db *sql.DB) InventoryLoaderConfig {
			config.DB = db
			return config
		}),
	)
}

func getInventorySource(cfg InventoryLoaderConfig) string {
	return cfg.Cmd.String("input-source")
}

func invLoaderFromFile(
	cfg InventoryLoaderConfig,
) InventoryLoaderIOE {
	inputPath := cfg.Cmd.String("input")
	return F.Pipe2(
		IOE.TryCatchError(func() (io.ReadCloser, error) {
			return os.Open(inputPath)
		}),
		IOE.MapLeft[io.ReadCloser](func(err error) error {
			return fmt.Errorf(
				"failed to open CSV file %s: %w",
				inputPath,
				err,
			)
		}),
		IOE.Map[error](
			func(reader io.ReadCloser) InventoryLoaderConfig {
				cfg.Reader = reader
				return cfg
			},
		),
	)
}

func invLoaderFromS3Bucket(
	cfg InventoryLoaderConfig,
) InventoryLoaderIOE {
	objectPath := fmt.Sprintf(
		"%s/%s",
		cfg.Cmd.String("s3-bucket-path"),
		cfg.Cmd.String("input"),
	)
	return F.Pipe2(
		IOE.TryCatchError(func() (io.ReadCloser, error) {
			return registry.GetS3Client().GetObject(
				cfg.Cmd.String("s3-bucket"),
				objectPath,
				minio.GetObjectOptions{},
			)
		}),
		IOE.MapLeft[io.ReadCloser](func(err error) error {
			return fmt.Errorf(
				"failed to open s3 file %s: %w",
				objectPath,
				err,
			)
		}),
		IOE.Map[error](
			func(reader io.ReadCloser) InventoryLoaderConfig {
				cfg.Reader = reader
				return cfg
			},
		),
	)
}

func invLoaderMap() map[string]func(InventoryLoaderConfig) InventoryLoaderIOE {
	return map[string]func(InventoryLoaderConfig) InventoryLoaderIOE{
		"folder": invLoaderFromFile,
		"bucket": invLoaderFromS3Bucket,
	}
}

func defaultInvLoader(
	cfg InventoryLoaderConfig,
) InventoryLoaderIOE {
	return IOE.Left[InventoryLoaderConfig](
		fmt.Errorf(
			"unsupported input source %s",
			cfg.Cmd.String("input-source"),
		),
	)
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
) InventoryLoaderIOE {
	return F.Pipe2(
		IOE.Do[error](InventoryLoaderConfig{
			Cmd:    cmd,
			Logger: logger,
		}),
		IOE.Chain(F.Switch(
			getInventorySource,
			invLoaderMap(),
			defaultInvLoader,
		)),
		IOE.Chain(inventoryBuilder),
	)
}
