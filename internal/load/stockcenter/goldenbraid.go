package stockcenter

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	fperrors "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	File "github.com/IBM/fp-go/ioeither/file"
	N "github.com/IBM/fp-go/number"
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
	Logger *slog.Logger
}

const gbMaxErrorMessages = 5

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

// Context for UPSERT pipeline
type GoldenBraidContext struct {
	Plasmid   *source.GoldenBraidPlasmid
	UserEmail string
	PlasmidCV string
}

// Result of UPSERT operation
type GoldenBraidUpsertResult struct {
	PlasmidID string
	Created   bool // true=created, false=updated
	Error     error
}

// GoldenBraidProcessingResult holds aggregate processing statistics
type GoldenBraidProcessingResult struct {
	CreatedCount int      // New plasmids created
	UpdatedCount int      // Existing plasmids updated
	Successes    []string // All plasmid IDs (created + updated)
	Errors       []error
	ErrorCount   int
}

var (
	gbIntSumSemigroup      = N.SemigroupSum[int]()
	gbStringArraySemigroup = S.MakeSemigroup(func(a, b []string) []string {
		return A.ArrayConcatAll(a, b)
	})
	gbErrorsSemigroup = S.MakeSemigroup(func(a, b []error) []error {
		return F.Pipe1(
			A.ArrayConcatAll(a, b),
			A.Slice[error](0, gbMaxErrorMessages),
		)
	})
)

func GoldenBraidSummarySemigroup() S.Semigroup[GoldenBraidProcessingResult] {
	return S.MakeSemigroup(
		func(a, b GoldenBraidProcessingResult) GoldenBraidProcessingResult {
			return GoldenBraidProcessingResult{
				CreatedCount: gbIntSumSemigroup.Concat(
					a.CreatedCount,
					b.CreatedCount,
				),
				UpdatedCount: gbIntSumSemigroup.Concat(
					a.UpdatedCount,
					b.UpdatedCount,
				),
				Successes: gbStringArraySemigroup.Concat(
					a.Successes,
					b.Successes,
				),
				ErrorCount: gbIntSumSemigroup.Concat(
					a.ErrorCount,
					b.ErrorCount,
				),
				Errors: gbErrorsSemigroup.Concat(
					a.Errors,
					b.Errors,
				),
			}
		},
	)
}

func goldenBraidUpsertResultToSummary(
	result GoldenBraidUpsertResult,
) GoldenBraidProcessingResult {
	return F.Pipe1(
		result,
		F.Ternary(
			// Check if error
			func(r GoldenBraidUpsertResult) bool { return r.Error != nil },
			// Error branch
			func(r GoldenBraidUpsertResult) GoldenBraidProcessingResult {
				return GoldenBraidProcessingResult{
					ErrorCount: 1,
					Errors:     []error{r.Error},
				}
			},
			// Success branch - nested ternary for created vs updated
			func(r GoldenBraidUpsertResult) GoldenBraidProcessingResult {
				return F.Pipe1(
					r,
					F.Ternary(
						func(r GoldenBraidUpsertResult) bool { return r.Created },
						// Created branch
						func(r GoldenBraidUpsertResult) GoldenBraidProcessingResult {
							return GoldenBraidProcessingResult{
								CreatedCount: 1,
								Successes:    []string{r.PlasmidID},
							}
						},
						// Updated branch
						func(r GoldenBraidUpsertResult) GoldenBraidProcessingResult {
							return GoldenBraidProcessingResult{
								UpdatedCount: 1,
								Successes:    []string{r.PlasmidID},
							}
						},
					),
				)
			},
		),
	)
}

// Build NewPlasmid request from context
func buildNewPlasmidRequest(ctx GoldenBraidContext) *stock.NewPlasmid {
	return &stock.NewPlasmid{
		Data: &stock.NewPlasmid_Data{
			Type: "plasmid",
			Attributes: &stock.NewPlasmidAttributes{
				Name:      ctx.Plasmid.Name,
				CreatedBy: ctx.UserEmail,
				UpdatedBy: ctx.UserEmail,
				Summary:   ctx.Plasmid.Summary,
				Genes: O.GetOrElse(
					F.Constant([]string{}),
				)(
					ctx.Plasmid.Genes,
				),
				Publications: O.GetOrElse(
					F.Constant([]string{}),
				)(
					ctx.Plasmid.Publications,
				),
				DictyPlasmidProperty: ctx.PlasmidCV,
			},
		},
	}
}

// Build PlasmidUpdate request from context and existing plasmid
func buildUpdatePlasmidRequest(
	ctx GoldenBraidContext,
	existing *stock.Plasmid,
) *stock.PlasmidUpdate {
	return &stock.PlasmidUpdate{
		Data: &stock.PlasmidUpdate_Data{
			Id:   existing.Data.Id,
			Type: "plasmid",
			Attributes: &stock.PlasmidUpdateAttributes{
				Name:      ctx.Plasmid.Name,
				UpdatedBy: ctx.UserEmail,
				Summary:   ctx.Plasmid.Summary,
				Genes: O.GetOrElse(
					F.Constant([]string{}),
				)(
					ctx.Plasmid.Genes,
				),
				Publications: O.GetOrElse(
					F.Constant([]string{}),
				)(
					ctx.Plasmid.Publications,
				),
				DictyPlasmidProperty: ctx.PlasmidCV,
			},
		},
	}
}

func createNewPlasmid(
	ctx GoldenBraidContext,
) func(O.Option[*stock.Plasmid]) IOE.IOEither[error, GoldenBraidUpsertResult] {
	return func(_ O.Option[*stock.Plasmid]) IOE.IOEither[error, GoldenBraidUpsertResult] {
		return F.Pipe2(
			IOE.TryCatchError(func() (*stock.Plasmid, error) {
				return regsc.GetStockAPIClient().CreatePlasmid(
					context.Background(),
					buildNewPlasmidRequest(ctx),
				)
			}),
			IOE.MapLeft[*stock.Plasmid](func(err error) error {
				return fmt.Errorf(
					"failed to create plasmid %s: %w",
					ctx.Plasmid.Name,
					err,
				)
			}),
			IOE.Map[error](
				func(created *stock.Plasmid) GoldenBraidUpsertResult {
					return GoldenBraidUpsertResult{
						PlasmidID: created.Data.Id,
						Created:   true,
						Error:     nil,
					}
				},
			),
		)
	}
}

func updateExistingPlasmid(
	ctx GoldenBraidContext,
) func(O.Option[*stock.Plasmid]) IOE.IOEither[error, GoldenBraidUpsertResult] {
	return func(opt O.Option[*stock.Plasmid]) IOE.IOEither[error, GoldenBraidUpsertResult] {
		existing := O.GetOrElse(func() *stock.Plasmid {
			panic("impossible: Option was None in update branch")
		})(opt)
		return F.Pipe2(
			IOE.TryCatchError(func() (*stock.Plasmid, error) {
				return regsc.GetStockAPIClient().UpdatePlasmid(
					context.Background(),
					buildUpdatePlasmidRequest(ctx, existing),
				)
			}),
			IOE.MapLeft[*stock.Plasmid](func(err error) error {
				return fmt.Errorf(
					"failed to update plasmid %s: %w",
					existing.Data.Attributes.Name,
					err,
				)
			}),
			IOE.Map[error](
				func(updated *stock.Plasmid) GoldenBraidUpsertResult {
					return GoldenBraidUpsertResult{
						PlasmidID: updated.Data.Id,
						Created:   false,
						Error:     nil,
					}
				},
			),
		)
	}
}

// processPlasmidWithUpsert processes a single validated plasmid by upserting it to the API
func processPlasmidWithUpsert(
	userEmail string,
	plasmidCVTerm string,
) func(*source.GoldenBraidPlasmid) E.Either[error, GoldenBraidUpsertResult] {
	return func(plasmid *source.GoldenBraidPlasmid) E.Either[error, GoldenBraidUpsertResult] {
		ctx := GoldenBraidContext{
			Plasmid:   plasmid,
			UserEmail: userEmail,
			PlasmidCV: plasmidCVTerm,
		}
		return F.Pipe3(
			fetchPlasmidByName(ctx.Plasmid.Name),
			IOE.Chain(F.Ternary(
				O.IsSome[*stock.Plasmid],
				updateExistingPlasmid(ctx),
				createNewPlasmid(ctx),
			)),
			fputil.ToEither[error, GoldenBraidUpsertResult],
			E.MapLeft[GoldenBraidUpsertResult](func(err error) error {
				return fmt.Errorf(
					"failed to upsert plasmid %s: %w",
					plasmid.Name,
					err,
				)
			}),
		)
	}
}

// fetchPlasmidByName uses ListPlasmids with filter to find plasmid by name
func fetchPlasmidByName(
	name string,
) IOE.IOEither[error, O.Option[*stock.Plasmid]] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
			return regsc.GetStockAPIClient().ListPlasmids(
				context.Background(),
				&stock.StockParameters{
					Filter: fmt.Sprintf("plasmid_name===%s", name),
					Limit:  1,
				},
			)
		}),
		IOE.MapLeft[*stock.PlasmidCollection](func(err error) error {
			return fmt.Errorf("failed to fetch plasmid %s: %w", name, err)
		}),
		IOE.Map[error](collectionToOption),
	)
}

// collectionToOption converts PlasmidCollection to Option[Plasmid]
func collectionToOption(
	collection *stock.PlasmidCollection,
) O.Option[*stock.Plasmid] {
	return F.Pipe2(
		collection.Data,
		A.Head[*stock.PlasmidCollection_Data], // Returns Option - None if empty
		O.Map(convertCollectionDataToPlasmid), // Transform if Some
	)
}

// convertCollectionDataToPlasmid converts collection data to full Plasmid message
func convertCollectionDataToPlasmid(
	data *stock.PlasmidCollection_Data,
) *stock.Plasmid {
	return &stock.Plasmid{
		Data: &stock.Plasmid_Data{
			Type:       data.Type,
			Id:         data.Id,
			Attributes: data.Attributes,
		},
	}
}

// openReader opens a CSV file from config and returns a reader wrapped in IOEither
func openReader(
	config LoaderConfig,
) IOE.IOEither[error, LoaderConfig] {
	return F.Pipe1(
		config,
		F.Switch(
			func(cfg LoaderConfig) string {
				source, _ := cfg.Cmd.Flags().GetString("input-source")
				return source
			},
			map[string]func(LoaderConfig) IOE.IOEither[error, LoaderConfig]{
				"folder": openGoldenBraidFileReader,
				"bucket": openGoldenBraidS3Reader,
			},
			defaultGoldenBraidReader,
		),
	)
}

func openGoldenBraidFileReader(config LoaderConfig) IOE.IOEither[error, LoaderConfig] {
	inputPath, _ := config.Cmd.Flags().GetString("input")

	return F.Pipe2(
		File.Open(inputPath),
		IOE.MapLeft[*os.File](func(err error) error {
			return fmt.Errorf("failed to open CSV file %s: %w", inputPath, err)
		}),
		IOE.Chain(func(file *os.File) IOE.IOEither[error, LoaderConfig] {
			reader := csv.NewReader(file)
			reader.FieldsPerRecord = source.GoldenBraidFieldCount
			reader.TrimLeadingSpace = true

			// Read and skip header - returns IOE
			return F.Pipe2(
				IOE.TryCatchError(func() ([]string, error) {
					return reader.Read()
				}),
				IOE.MapLeft[[]string](func(err error) error {
					file.Close() // Clean up on error
					return fmt.Errorf("failed to read CSV header: %w", err)
				}),
				IOE.Map[error](func(_ []string) LoaderConfig {
					config.Reader = reader
					config.Closer = file
					return config
				}),
			)
		}),
	)
}

func openGoldenBraidS3Reader(config LoaderConfig) IOE.IOEither[error, LoaderConfig] {
	bucket, _ := config.Cmd.Flags().GetString("s3-bucket")
	bucketPath, _ := config.Cmd.Flags().GetString("s3-bucket-path")
	file, _ := config.Cmd.Flags().GetString("input")
	objectPath := path.Join(bucketPath, file)

	return F.Pipe2(
		IOE.TryCatchError(func() (io.ReadCloser, error) {
			return registry.GetS3Client().GetObject(
				bucket,
				objectPath,
				minio.GetObjectOptions{},
			)
		}),
		IOE.MapLeft[io.ReadCloser](func(err error) error {
			return fmt.Errorf("failed to open s3 file %s: %w", objectPath, err)
		}),
		IOE.Chain(func(reader io.ReadCloser) IOE.IOEither[error, LoaderConfig] {
			csvReader := csv.NewReader(reader)
			csvReader.FieldsPerRecord = source.GoldenBraidFieldCount
			csvReader.TrimLeadingSpace = true

			// Read and skip header - returns IOE
			return F.Pipe2(
				IOE.TryCatchError(func() ([]string, error) {
					return csvReader.Read()
				}),
				IOE.MapLeft[[]string](func(err error) error {
					reader.Close() // Clean up on error
					return fmt.Errorf("failed to read CSV header: %w", err)
				}),
				IOE.Map[error](func(_ []string) LoaderConfig {
					config.Reader = csvReader
					config.Closer = reader
					return config
				}),
			)
		}),
	)
}

func defaultGoldenBraidReader(cfg LoaderConfig) IOE.IOEither[error, LoaderConfig] {
	source, _ := cfg.Cmd.Flags().GetString("input-source")
	return IOE.Left[LoaderConfig](
		fmt.Errorf("unsupported input source %s", source),
	)
}

// streamAndProcessRecords reads CSV records one by one and processes them
func streamAndProcessRecords(
	config LoaderConfig,
) IOE.IOEither[error, GoldenBraidProcessingResult] {
	return IOE.TryCatchError(func() (GoldenBraidProcessingResult, error) {
		userEmail := config.Viper.GetString("user-email")
		plasmidCVTerm := config.Viper.GetString("plasmid-cvterm")
		summary := GoldenBraidProcessingResult{}

		for {
			record, err := config.Reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				config.Logger.Error(
					"CSV read error, partial results",
					"successes", len(summary.Successes),
					"errors", summary.ErrorCount,
					"error", err,
				)
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
				E.Chain(processPlasmidWithUpsert(userEmail, plasmidCVTerm)),
				E.Fold(
					func(err error) GoldenBraidUpsertResult {
						return GoldenBraidUpsertResult{Error: err}
					},
					F.Identity[GoldenBraidUpsertResult],
				),
			)

			summary = GoldenBraidSummarySemigroup().Concat(
				summary,
				goldenBraidUpsertResultToSummary(result),
			)
		}

		return summary, nil
	})
}

func releaseGoldenBraidReader(
	config LoaderConfig,
	_ E.Either[error, GoldenBraidProcessingResult],
) IOE.IOEither[error, any] {
	return IOE.TryCatchError(func() (any, error) {
		if config.Closer != nil {
			return nil, config.Closer.Close()
		}
		return nil, nil
	})
}

func useGoldenBraidReader(
	config LoaderConfig,
) IOE.IOEither[error, GoldenBraidProcessingResult] {
	return F.Pipe2(
		config,
		logGoldenBraidStep[LoaderConfig](
			config.Logger,
			"CSV reader opened, processing records",
		),
		IOE.Chain(streamAndProcessRecords),
	)
}

// LoadGoldenBraid is the main entry point for loading GoldenBraid CSV data
func LoadGoldenBraid(cmd *cobra.Command, _ []string) error {
	handler := logger.GetSlogHandler(cmd)
	slogger := slog.New(handler)
	elog := E.Logger[error, GoldenBraidProcessingResult](
		slog.NewLogLogger(handler, slog.LevelInfo),
		slog.NewLogLogger(handler, slog.LevelError),
	)

	return F.Pipe6(
		IOE.Do[error](LoaderConfig{
			Viper:  viper.GetViper(),
			Cmd:    cmd,
			Reader: nil,
			Logger: slogger,
		}),
		IOE.ChainFirst(
			logGoldenBraidStep[LoaderConfig](
				slogger,
				"Starting GoldenBraid loading",
			),
		),
		IOE.Chain(
			func(config LoaderConfig) IOE.IOEither[error, GoldenBraidProcessingResult] {
				return IOE.Bracket(
					openReader(config),
					useGoldenBraidReader,
					releaseGoldenBraidReader,
				)
			},
		),
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
					"created", summary.CreatedCount,
					"updated", summary.UpdatedCount,
					"total_successes", len(summary.Successes),
					"errors", summary.ErrorCount,
				)
				return nil
			},
		),
	)
}
