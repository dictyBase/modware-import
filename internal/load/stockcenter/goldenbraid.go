package stockcenter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"

	E "github.com/IBM/fp-go/either"
	fperrors "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/semigroup"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoaderContext is the initial empty context
type LoaderContext struct{}

// LoaderConfig holds viper configuration and CSV reader
type LoaderConfig struct {
	LoaderContext
	Cmd    *cobra.Command
	Reader *csv.Reader
	Closer io.Closer
	Viper  *viper.Viper
}

type GoldenBraidReaderResource struct {
	Reader *csv.Reader
	Closer io.Closer
}

// PlasmidProcessingResult represents the result of processing a single plasmid
type PlasmidProcessingResult struct {
	PlasmidID string
	Error     O.Option[error]
}

// createErrorResult creates a PlasmidProcessingResult for an error
func createErrorResult(e error) PlasmidProcessingResult {
	return PlasmidProcessingResult{
		PlasmidID: "",
		Error:     O.Some(e),
	}
}

func createProcessingResult(id string) PlasmidProcessingResult {
	return PlasmidProcessingResult{
		PlasmidID: id,
		Error:     O.None[error](),
	}
}

func logGoldenBraidStep[T any](
	logger *slog.Logger,
	msg string,
) func(T) IOE.IOEither[error, T] {
	return func(data T) IOE.IOEither[error, T] {
		return IOE.TryCatchError(func() (T, error) {
			logger.Info(msg, "payload", data)
			return data, nil
		})
	}
}

// GoldenBraidProcessingResult holds aggregate processing statistics
type GoldenBraidProcessingResult struct {
	Successes  []string
	Errors     []string
	ErrorCount int
}

func GoldenBraidSummarySemigroup() S.Semigroup[GoldenBraidProcessingResult] {
	return S.MakeSemigroup(
		func(a, b GoldenBraidProcessingResult) GoldenBraidProcessingResult {
			return GoldenBraidProcessingResult{
				Successes:  append(a.Successes, b.Successes...),
				ErrorCount: a.ErrorCount + b.ErrorCount,
				Errors:     append(a.Errors, b.Errors...),
			}
		},
	)
}

// openReader opens a CSV file from config and returns a reader wrapped in IOEither
func openReader(
	config LoaderConfig,
) IOE.IOEither[error, GoldenBraidReaderResource] {
	return IOE.TryCatchError(func() (GoldenBraidReaderResource, error) {
		source, _ := config.Cmd.Flags().GetString("input-source")
		switch source {
		case "folder":
			return openGoldenBraidFileReader(config)
		case "bucket":
			return openGoldenBraidS3Reader(config)
		default:
			return GoldenBraidReaderResource{}, fmt.Errorf(
				"unsupported input source %s",
				source,
			)
		}
	})
}

func openGoldenBraidFileReader(config LoaderConfig) (GoldenBraidReaderResource, error) {
	inputPath, _ := config.Cmd.Flags().GetString("input")
	file, err := os.Open(inputPath)
	if err != nil {
		return GoldenBraidReaderResource{}, fmt.Errorf(
			"failed to open CSV file %s: %w",
			inputPath,
			err,
		)
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = source.GoldenBraidFieldCount
	reader.TrimLeadingSpace = true

	// Skip header row
	if _, err := reader.Read(); err != nil {
		file.Close()
		return GoldenBraidReaderResource{}, fmt.Errorf(
			"failed to read CSV header: %w",
			err,
		)
	}

	return GoldenBraidReaderResource{Reader: reader, Closer: file}, nil
}

func openGoldenBraidS3Reader(config LoaderConfig) (GoldenBraidReaderResource, error) {
	bucket, _ := config.Cmd.Flags().GetString("s3-bucket")
	path, _ := config.Cmd.Flags().GetString("s3-bucket-path")
	file, _ := config.Cmd.Flags().GetString("input")
	reader, err := registry.GetS3Client().GetObject(
		bucket,
		fmt.Sprintf("%s/%s", path, file),
		minio.GetObjectOptions{},
	)
	if err != nil {
		return GoldenBraidReaderResource{}, fmt.Errorf(
			"failed to open s3 file %s/%s: %w",
			path,
			file,
			err,
		)
	}
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = source.GoldenBraidFieldCount
	csvReader.TrimLeadingSpace = true

	// Skip header row
	if _, err := csvReader.Read(); err != nil {
		reader.Close()
		return GoldenBraidReaderResource{}, fmt.Errorf(
			"failed to read CSV header: %w",
			err,
		)
	}
	return GoldenBraidReaderResource{Reader: csvReader, Closer: reader}, nil
}

// SetReader is a curried setter that adds CSV reader to config
var SetReader = F.Curry2(
	func(resource GoldenBraidReaderResource, config LoaderConfig) LoaderConfig {
		config.Reader = resource.Reader
		config.Closer = resource.Closer
		return config
	},
)

// streamAndProcessRecords reads CSV records one by one and processes them
func streamAndProcessRecords(
	config LoaderConfig,
) IOE.IOEither[error, GoldenBraidProcessingResult] {
	return IOE.TryCatchError(func() (GoldenBraidProcessingResult, error) {
		if config.Closer != nil {
			defer config.Closer.Close()
		}
		userEmail := config.Viper.GetString("user-email")
		plasmidCVTerm := config.Viper.GetString("plasmid-cvterm")
		summary := GoldenBraidProcessingResult{}

		for {
			record, err := config.Reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return summary, fmt.Errorf("CSV read error: %w", err)
			}

			// Process single record (pure Either pipeline - integrated)
			result := F.Pipe7(
				record,
				E.FromPredicate(
					source.HasValidRecordLength,
					source.RecordLengthError,
				),
				E.Map[error](source.BuildPlasmid(
					userEmail,
					plasmidCVTerm,
				)),
				E.Chain(E.FromPredicate(
					source.HasValidName,
					source.NameError,
				)),
				E.Chain(
					E.FromPredicate(
						source.HasValidSummary,
						source.SummaryError,
					),
				),
				E.Chain(E.FromPredicate(
					source.HasValidUser,
					source.UserError,
				)),
				E.Chain(processPlasmid),
				E.Fold(
					createErrorResult,
					createProcessingResult,
				),
			)

			summary = GoldenBraidSummarySemigroup().Concat(summary, goldenBraidResultToSummary(result))
		}

		return summary, nil
	})
}

func goldenBraidResultToSummary(
	result PlasmidProcessingResult,
) GoldenBraidProcessingResult {
	return F.Pipe1(
		result.Error,
		O.Fold(
			func() GoldenBraidProcessingResult {
				return GoldenBraidProcessingResult{
					Successes: []string{result.PlasmidID},
				}
			},
			func(err error) GoldenBraidProcessingResult {
				return GoldenBraidProcessingResult{
					ErrorCount: 1,
					Errors:     []string{err.Error()},
				}
			},
		),
	)
}

// processPlasmid processes a single validated plasmid by loading it to the API
func processPlasmid(p *source.GoldenBraidPlasmid) E.Either[error, string] {
	return F.Pipe2(p, loadPlasmidToAPI, fputil.ToEither[error, string])
}

// loadPlasmidToAPI loads a plasmid to the stock center API
func loadPlasmidToAPI(
	p *source.GoldenBraidPlasmid,
) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		client := regsc.GetStockAPIClient()

		plasmidData := &stock.NewPlasmid{
			Data: &stock.NewPlasmid_Data{
				Type: "plasmid",
				Attributes: &stock.NewPlasmidAttributes{
					Name:      p.Name,
					Summary:   p.Summary,
					CreatedBy: p.User,
					UpdatedBy: p.User,
					Genes: O.GetOrElse(
						F.Constant([]string{}),
					)(
						p.Genes,
					),
					Publications: O.GetOrElse(
						F.Constant([]string{}),
					)(
						p.Publications,
					),
				},
			},
		}

		response, err := client.CreatePlasmid(
			context.Background(),
			plasmidData,
		)
		if err != nil {
			return "", fmt.Errorf(
				"failed to create plasmid %s: %w",
				p.Name,
				err,
			)
		}

		return response.Data.Id, nil
	})
}

// LoadGoldenBraid is the main entry point for loading GoldenBraid CSV data
func LoadGoldenBraid(cmd *cobra.Command, _ []string) error {
	handler := logger.GetSlogHandler(cmd)
	slogger := slog.New(handler)
	elog := E.Logger[error, GoldenBraidProcessingResult](
		slog.NewLogLogger(handler, slog.LevelInfo),
		slog.NewLogLogger(handler, slog.LevelError),
	)

	return F.Pipe8(
		IOE.Do[error](LoaderConfig{
			Viper:  viper.GetViper(),
			Cmd:    cmd,
			Reader: nil,
		}),
		IOE.ChainFirst(
			logGoldenBraidStep[LoaderConfig](
				slogger,
				"Starting GoldenBraid loading",
			),
		),
		IOE.Bind(SetReader, openReader),
		IOE.ChainFirst(
			logGoldenBraidStep[LoaderConfig](
				slogger,
				"CSV reader opened, processing records",
			),
		),
		IOE.Chain(streamAndProcessRecords),
		IOE.ChainFirst(
			logGoldenBraidStep[GoldenBraidProcessingResult](
				slogger,
				"Processing complete",
			),
		),
		fputil.ToEither[error, GoldenBraidProcessingResult],
		elog("GoldenBraid loading result"),
		E.Fold(
			fperrors.IdentityError,
			func(summary GoldenBraidProcessingResult) error {
				slogger.Info(
					"GoldenBraid loading summary",
					"successes", len(summary.Successes),
					"errors", summary.ErrorCount,
				)
				return nil
			},
		),
	)
}
