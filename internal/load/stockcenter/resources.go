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
	File "github.com/IBM/fp-go/ioeither/file"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"
	"github.com/nao1215/filesql"
	"github.com/urfave/cli/v2"
)

type InventoryLoaderConfig struct {
	DB                *sql.DB
	Reader            io.ReadCloser
	GoldenBraidReader io.ReadCloser
	Logger            *slog.Logger
	Cmd               *cli.Context
}

// DBBuildContext contains only the fields needed for building inventory database
type DBBuildContext struct {
	Reader            io.ReadCloser
	GoldenBraidReader io.ReadCloser
	Context           context.Context
	Builder           *filesql.DBBuilder
}

type InventoryLoaderIOE = IOE.IOEither[error, InventoryLoaderConfig]

func closeInventoryConfig(config InventoryLoaderConfig) {
	if config.DB != nil {
		_ = config.DB.Close()
	}
	if config.Reader != nil {
		_ = config.Reader.Close()
	}
	if config.GoldenBraidReader != nil {
		_ = config.GoldenBraidReader.Close()
	}
}

// buildInventorySQL creates and validates a filesql builder from two readers:
// goldenbraid_inventory (inventory CSV) and goldenbraid (plasmid CSV).
func buildInventorySQL(
	ctx DBBuildContext,
) IOE.IOEither[error, DBBuildContext] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*filesql.DBBuilder, error) {
			return filesql.NewBuilder().
				AddReader(ctx.Reader, "goldenbraid_inventory", filesql.FileTypeCSV).
				AddReader(ctx.GoldenBraidReader, "goldenbraid", filesql.FileTypeCSV).
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
			Reader:            config.Reader,
			GoldenBraidReader: config.GoldenBraidReader,
			Context:           context.Background(),
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
				DB:                db,
				Reader:            config.Reader,
				GoldenBraidReader: config.GoldenBraidReader,
				Logger:            config.Logger,
				Cmd:               config.Cmd,
			}
		}),
	)
}

func invLoaderFromFile(
	cfg InventoryLoaderConfig,
) InventoryLoaderIOE {
	inputPath := cfg.Cmd.String("input")
	return F.Pipe2(
		File.Open(inputPath),
		IOE.MapLeft[*os.File](func(err error) error {
			return fmt.Errorf(
				"failed to open CSV file %s: %w",
				inputPath,
				err,
			)
		}),
		IOE.Map[error](func(f *os.File) InventoryLoaderConfig {
			cfg.Reader = f
			return cfg
		}),
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
		IOE.Map[error](func(reader io.ReadCloser) InventoryLoaderConfig {
			cfg.Reader = reader
			return cfg
		}),
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

func gbLoaderFromFile(cfg InventoryLoaderConfig) InventoryLoaderIOE {
	inputPath := cfg.Cmd.String("goldenbraid-input")
	return F.Pipe2(
		File.Open(inputPath),
		IOE.MapLeft[*os.File](func(err error) error {
			return fmt.Errorf(
				"failed to open goldenbraid CSV file %s: %w",
				inputPath,
				err,
			)
		}),
		IOE.Map[error](func(f *os.File) InventoryLoaderConfig {
			cfg.GoldenBraidReader = f
			return cfg
		}),
	)
}

func gbLoaderFromS3Bucket(cfg InventoryLoaderConfig) InventoryLoaderIOE {
	objectPath := fmt.Sprintf(
		"%s/%s",
		cfg.Cmd.String("s3-bucket-path"),
		cfg.Cmd.String("goldenbraid-input"),
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
				"failed to open s3 goldenbraid file %s: %w",
				objectPath,
				err,
			)
		}),
		IOE.Map[error](func(reader io.ReadCloser) InventoryLoaderConfig {
			cfg.GoldenBraidReader = reader
			return cfg
		}),
	)
}

func ReleaseInventoryResources(
	config InventoryLoaderConfig,
	_ E.Either[error, InventoryProcessingSummary],
) IOE.IOEither[error, any] {
	return IOE.TryCatchError(func() (any, error) {
		closeInventoryConfig(config)
		return struct{}{}, nil
	})
}

func InventoryBuilderFromFile(
	cmd *cli.Context,
	logger *slog.Logger,
) InventoryLoaderIOE {
	return F.Pipe3(
		IOE.Do[error](InventoryLoaderConfig{
			Cmd:    cmd,
			Logger: logger,
		}),
		// Step 1: open inventory file → sets Reader
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
		// Step 2: open goldenbraid plasmid file → sets GoldenBraidReader
		IOE.Chain(
			F.Switch(
				func(cfg InventoryLoaderConfig) string {
					return cfg.Cmd.String("input-source")
				},
				map[string]func(InventoryLoaderConfig) InventoryLoaderIOE{
					"folder": gbLoaderFromFile,
					"bucket": gbLoaderFromS3Bucket,
				},
				defaultInvLoader,
			)),
		// Step 3: build filesql DB from both readers
		IOE.Chain(inventoryBuilder),
	)
}
