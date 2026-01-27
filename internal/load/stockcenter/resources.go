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

// DBBuildContext contains only the fields needed for building inventory database
type DBBuildContext struct {
	Reader  io.ReadCloser
	Context context.Context
	Builder *filesql.DBBuilder
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

// buildInventorySQL creates and validates a filesql builder from a reader
func buildInventorySQL(
	ctx DBBuildContext,
) IOE.IOEither[error, DBBuildContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*filesql.DBBuilder, error) {
			return filesql.NewBuilder().
				AddReader(
					ctx.Reader,
					"inventory",
					filesql.FileTypeCSV,
				).
				Build(ctx.Context)
		}),
		IOE.Map[error](func(builder *filesql.DBBuilder) DBBuildContext {
			ctx.Builder = builder
			return ctx
		}),
	)
}

// openInventoryDB opens a database connection from a validated builder
func openInventoryDB(
	ctx DBBuildContext,
) IOE.IOEither[error, *sql.DB] {
	return IOE.TryCatchError(func() (*sql.DB, error) {
		return ctx.Builder.Open(ctx.Context)
	})
}

func inventoryBuilder(
	config InventoryLoaderConfig,
) InventoryLoaderIOE {
	return F.Pipe4(
		DBBuildContext{
			Reader:  config.Reader,
			Context: context.Background(),
		},
		buildInventorySQL,
		IOE.Chain(openInventoryDB),
		IOE.MapLeft[*sql.DB](func(err error) error {
			return fmt.Errorf(
				"failed to build inventory database: %w",
				err,
			)
		}),
		IOE.Map[error](func(db *sql.DB) InventoryLoaderConfig {
			return InventoryLoaderConfig{
				DB:     db,
				Reader: config.Reader,
				Logger: config.Logger,
				Cmd:    config.Cmd,
			}
		}),
	)
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
		IOE.Chain(
			F.Switch(
				func(cfg InventoryLoaderConfig) string {
					return cfg.Cmd.String("input-source")
				},
				map[string]func(InventoryLoaderConfig) InventoryLoaderIOE{
					"folder": invLoaderFromFile,
					"bucket": invLoaderFromS3Bucket,
				},
				defaultInvLoader,
			)),
		IOE.Chain(inventoryBuilder),
	)
}
