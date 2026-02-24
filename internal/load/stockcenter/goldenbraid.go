package stockcenter

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
	File "github.com/IBM/fp-go/ioeither/file"
	N "github.com/IBM/fp-go/number"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/urfave/cli/v2"
)

// LoaderContext is the initial empty context
type LoaderContext struct{}

// LoaderConfig holds CLI context and CSV reader
type LoaderConfig struct {
	LoaderContext
	Cmd    *cli.Context
	Reader *csv.Reader
	Closer io.Closer
}

const gbMaxErrorMessages = 5

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
				return cfg.Cmd.String("input-source")
			},
			map[string]func(LoaderConfig) IOE.IOEither[error, LoaderConfig]{
				"folder": openGoldenBraidFileReader,
				"bucket": openGoldenBraidS3Reader,
			},
			defaultGoldenBraidReader,
		),
	)
}

func openGoldenBraidFileReader(
	config LoaderConfig,
) IOE.IOEither[error, LoaderConfig] {
	inputPath := config.Cmd.String("input")

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

func openGoldenBraidS3Reader(
	config LoaderConfig,
) IOE.IOEither[error, LoaderConfig] {
	bucket := config.Cmd.String("s3-bucket")
	bucketPath := config.Cmd.String("s3-bucket-path")
	file := config.Cmd.String("input")
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

func defaultGoldenBraidReader(
	cfg LoaderConfig,
) IOE.IOEither[error, LoaderConfig] {
	return IOE.Left[LoaderConfig](
		fmt.Errorf("unsupported input source %s", cfg.Cmd.String("input-source")),
	)
}

// streamAndProcessRecords reads CSV records one by one and processes them
func streamAndProcessRecords(
	config LoaderConfig,
) IOE.IOEither[error, GoldenBraidProcessingResult] {
	return IOE.TryCatchError(func() (GoldenBraidProcessingResult, error) {
		userEmail := config.Cmd.String("user-email")
		plasmidCVTerm := config.Cmd.String("plasmid-cvterm")
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
		IOE.Of[error](config),
		IOE.ChainFirstIOK[error](
			IO.Logf[LoaderConfig](
				"CSV reader opened, processing records: %+v",
			),
		),
		IOE.Chain(streamAndProcessRecords),
	)
}

// LoadGoldenBraidCli is the main entry point for loading GoldenBraid CSV data
func LoadGoldenBraidCli(cmd *cli.Context) error {
	handler := logger.GetCliSlogHandler(cmd)
	slogger := slog.New(handler)
	elog := E.Logger[error, GoldenBraidProcessingResult](
		slog.NewLogLogger(handler, slog.LevelInfo),
		slog.NewLogLogger(handler, slog.LevelError),
	)

	output := F.Pipe6(
		IOE.Do[error](LoaderConfig{
			Cmd:    cmd,
			Reader: nil,
		}),
		IOE.ChainFirstIOK[error](
			IO.Logf[LoaderConfig](
				"Starting GoldenBraid loading: %+v",
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
		IOE.ChainFirstIOK[error](
			IO.Logf[GoldenBraidProcessingResult](
				"Processing complete: %+v",
			),
		),
		fputil.ToEither[error, GoldenBraidProcessingResult],
		elog("GoldenBraid loading result"),
		E.Fold(
			onGoldenBraidSummaryError,
			onGoldenBraidSummarySuccess),
	)

	return handleGoldenBraidSummaryOutput(slogger, output)
}

func onGoldenBraidSummaryError(
	err error,
) T.Tuple2[GoldenBraidProcessingResult, error] {
	return T.MakeTuple2(GoldenBraidProcessingResult{}, err)
}

func onGoldenBraidSummarySuccess(
	summary GoldenBraidProcessingResult,
) T.Tuple2[GoldenBraidProcessingResult, error] {
	return T.MakeTuple2(summary, (error)(nil))
}

func handleGoldenBraidSummaryOutput(
	slogger *slog.Logger,
	output T.Tuple2[GoldenBraidProcessingResult, error],
) error {
	if output.F2 != nil {
		return output.F2
	}
	summary := output.F1
	slogger.Info(
		"GoldenBraid loading summary",
		"created", summary.CreatedCount,
		"updated", summary.UpdatedCount,
		"total_successes", len(summary.Successes),
		"errors", summary.ErrorCount,
	)

	if len(summary.Errors) > 0 {
		joinedErrors := errors.Join(summary.Errors...)
		slogger.Warn(
			"error samples (first 5)",
			"sample_count", len(summary.Errors),
			"total_errors", summary.ErrorCount,
			"messages", joinedErrors.Error(),
		)
	}

	return nil
}
