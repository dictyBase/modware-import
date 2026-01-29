package stockcenter

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	Eq "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
	File "github.com/IBM/fp-go/ioeither/file"
	S "github.com/IBM/fp-go/semigroup"
	T "github.com/IBM/fp-go/tuple"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/logger"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
)

const keywordRecordLength = 3

const maxErrorMessages = 5

var (
	errSkipRecord    = errors.New("skip record")
	errWrongProperty = errors.New("wrong property type")

	// Predicates for result classification
	hasError = func(r KeywordProcessingResult) bool { return r.Error != nil }
	// isWrongPropertyError: checks ONLY for wrong property type (expected filtering, not an error)
	// These lines have a different property type and should be skipped silently
	isWrongPropertyError = func(r KeywordProcessingResult) bool {
		return errors.Is(r.Error, errWrongProperty)
	}
	isCreated = func(r KeywordProcessingResult) bool { return r.Created }

	// Summary builders
	// skipSummary: for wrong property records only - skipped silently, not counted
	skipSummary    = func(_ KeywordProcessingResult) KeywordProcessingSummary { return KeywordProcessingSummary{} }
	createdSummary = func(_ KeywordProcessingResult) KeywordProcessingSummary {
		return KeywordProcessingSummary{CreatedCount: 1}
	}
	existingSummary = func(_ KeywordProcessingResult) KeywordProcessingSummary {
		return KeywordProcessingSummary{ExistingCount: 1}
	}
	// errorSummary: for real errors AND errSkipRecord (blank/invalid lines are errors)
	errorSummary = func(r KeywordProcessingResult) KeywordProcessingSummary {
		return KeywordProcessingSummary{
			ErrorCount: 1,
			Errors:     []error{r.Error},
		}
	}

	keywordIsNotBlankRecord = A.Any(func(s string) bool {
		return len(strings.TrimSpace(s)) > 0
	})

	caseInsensitiveEq = Eq.FromEquals(strings.EqualFold)

	trimSpace         = strings.TrimSpace
	normalizeProperty = F.Flow2(strings.TrimSpace, strings.ToLower)
	SetPlasmid        = F.Curry2(
		func(plasmid *stockpb.Plasmid, ctx KeywordContext) WithPlasmid {
			return WithPlasmid{
				KeywordContext: ctx,
				Plasmid:        plasmid,
			}
		},
	)
)

type KeywordLoaderContext struct{}

type KeywordLoaderConfig struct {
	KeywordLoaderContext
	Cmd    *cobra.Command
	Reader *csv.Reader
	Closer io.Closer
}

type KeywordReaderIOE = IOE.IOEither[error, KeywordLoaderConfig]

type KeywordRecord struct {
	PlasmidID string
	Property  string
	Term      string
}

type KeywordProcessingResult struct {
	Created bool
	Error   error
}

type KeywordProcessingSummary struct {
	CreatedCount  int     `json:"created_count"`
	ExistingCount int     `json:"existing_count"`
	ErrorCount    int     `json:"error_count"`
	Errors        []error `json:"-"`
}

type KeywordContext struct {
	Record KeywordRecord
}

type WithPlasmid struct {
	KeywordContext
	Plasmid *stockpb.Plasmid
}

func LoadPlasmidOntology(cmd *cobra.Command, _ []string) error {
	handler := logger.GetSlogHandler(cmd)
	slogger := slog.New(handler)
	elog := E.Logger[error, KeywordProcessingSummary](
		slog.NewLogLogger(handler, slog.LevelInfo),
		slog.NewLogLogger(handler, slog.LevelError),
	)

	output := F.Pipe7(
		IOE.Do[error](KeywordLoaderConfig{Cmd: cmd}),
		IOE.ChainFirstIOK[error](
			IO.Logf[KeywordLoaderConfig](
				"Starting plasmid ontology association: %+v",
			),
		),
		IOE.Chain(openKeywordReader),
		IOE.Chain(processAndAggregateKeywordRecords),
		IOE.ChainFirstIOK[error](
			IO.Logf[KeywordProcessingSummary](
				"Plasmid ontology association summary: %+v",
			),
		),
		fputil.ToEither[error, KeywordProcessingSummary],
		elog("plasmid ontology association result"),
		E.Fold(onSummaryError, onSummarySuccess),
	)

	return handleSummaryOutput(slogger, output)
}

func onSummaryError(err error) T.Tuple2[KeywordProcessingSummary, error] {
	return T.MakeTuple2(KeywordProcessingSummary{}, err)
}

func onSummarySuccess(
	summary KeywordProcessingSummary,
) T.Tuple2[KeywordProcessingSummary, error] {
	return T.MakeTuple2(summary, (error)(nil))
}

func handleSummaryOutput(
	slogger *slog.Logger,
	output T.Tuple2[KeywordProcessingSummary, error],
) error {
	if output.F2 != nil {
		return output.F2
	}
	summary := output.F1
	slogger.Info(
		"plasmid ontology summary",
		"created", summary.CreatedCount,
		"existing", summary.ExistingCount,
		"errors", summary.ErrorCount,
	)

	// Log error samples if present - use errors.Join to combine into single string
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

func configureTSVReader(reader io.Reader) *csv.Reader {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	return csvReader
}

func openKeywordReader(config KeywordLoaderConfig) KeywordReaderIOE {
	return F.Pipe1(
		config,
		F.Switch(
			func(cfg KeywordLoaderConfig) string {
				source, _ := cfg.
					Cmd.
					Flags().
					GetString("input-source")
				return source
			},
			map[string]func(KeywordLoaderConfig) KeywordReaderIOE{
				"folder": keywordReaderFromFile,
				"bucket": keywordReaderFromS3Bucket,
			},
			defaultKeywordReader,
		),
	)
}

func keywordReaderFromFile(cfg KeywordLoaderConfig) KeywordReaderIOE {
	inputPath, _ := cfg.Cmd.Flags().GetString("input")
	return F.Pipe2(
		File.Open(inputPath),
		IOE.MapLeft[*os.File](func(err error) error {
			return fmt.Errorf("failed to open TSV file %s: %w", inputPath, err)
		}),
		IOE.Map[error](func(file *os.File) KeywordLoaderConfig {
			cfg.Reader = configureTSVReader(file)
			cfg.Closer = file
			return cfg
		}),
	)
}

func keywordReaderFromS3Bucket(cfg KeywordLoaderConfig) KeywordReaderIOE {
	path, _ := cfg.Cmd.Flags().GetString("s3-bucket-path")
	file, _ := cfg.Cmd.Flags().GetString("input")
	bucket, _ := cfg.Cmd.Flags().GetString("s3-bucket")
	objectPath := fmt.Sprintf("%s/%s", path, file)

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
		IOE.Map[error](func(reader io.ReadCloser) KeywordLoaderConfig {
			cfg.Reader = configureTSVReader(reader)
			cfg.Closer = reader
			return cfg
		}),
	)
}

func defaultKeywordReader(cfg KeywordLoaderConfig) KeywordReaderIOE {
	source, _ := cfg.Cmd.Flags().GetString("input-source")
	return IOE.Left[KeywordLoaderConfig](
		fmt.Errorf("unsupported input source %s", source),
	)
}

func processAndAggregateKeywordRecords(
	config KeywordLoaderConfig,
) IOE.IOEither[error, KeywordProcessingSummary] {
	return IOE.TryCatchError(func() (KeywordProcessingSummary, error) {
		if config.Closer != nil {
			defer config.Closer.Close()
		}

		prop, _ := config.Cmd.Flags().GetString("property")
		summary := KeywordProcessingSummary{}
		for {
			keyRecord, err := config.Reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return summary, fmt.Errorf("tsv read error: %w", err)
			}

			result := F.Pipe6(
				keyRecord,
				E.FromPredicate(
					keywordIsNotBlankRecord,
					keywordSkipRecordError,
				),
				E.Chain(
					E.FromPredicate(
						keywordHasValidRecordLength,
						keywordRecordLengthError,
					),
				),
				E.Map[error](parseKeywordRecord),
				E.Chain(keywordFilterProperty(prop)),
				E.Chain(associateKeywordTerm),
				E.Fold(
					handleKeywordPipelineError,
					F.Identity[KeywordProcessingResult],
				),
			)

			summary = SummarySemigroup().Concat(
				summary,
				resultToSummary(result),
			)
		}

		return summary, nil
	})
}

func keywordFilterProperty(
	target string,
) func(KeywordRecord) E.Either[error, KeywordRecord] {
	return E.FromPredicate(
		func(r KeywordRecord) bool {
			return strings.EqualFold(r.Property, target)
		},
		func(_ KeywordRecord) error {
			return errWrongProperty
		},
	)
}

func handleKeywordPipelineError(err error) KeywordProcessingResult {
	return KeywordProcessingResult{
		Created: false,
		Error:   err,
	}
}

func keywordHasValidRecordLength(keyRecord []string) bool {
	return len(keyRecord) == keywordRecordLength
}

func keywordRecordLengthError(keyRecord []string) error {
	return fmt.Errorf(
		"invalid record: expected exactly %d columns, got %d (%s)",
		keywordRecordLength,
		len(keyRecord),
		strings.Join(keyRecord, "\t"),
	)
}

func parseKeywordRecord(keyRecord []string) KeywordRecord {
	return KeywordRecord{
		PlasmidID: trimSpace(keyRecord[0]),
		Property:  normalizeProperty(keyRecord[1]),
		Term:      trimSpace(keyRecord[2]),
	}
}

func fetchPlasmid(ctx KeywordContext) IOE.IOEither[error, *stockpb.Plasmid] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*stockpb.Plasmid, error) {
			return regsc.GetStockAPIClient().GetPlasmid(
				context.Background(),
				&stockpb.StockId{Id: ctx.Record.PlasmidID},
			)
		}),
		IOE.MapLeft[*stockpb.Plasmid](func(err error) error {
			return fmt.Errorf(
				"failed to fetch plasmid %s: %w",
				ctx.Record.PlasmidID,
				err,
			)
		}),
	)
}

func propertyAlreadyMatches(ctx WithPlasmid) bool {
	return caseInsensitiveEq.Equals(
		ctx.Plasmid.GetData().GetAttributes().GetDictyPlasmidProperty(),
		ctx.Record.Term,
	)
}

func returnExistingResult(
	ctx WithPlasmid,
) IOE.IOEither[error, KeywordProcessingResult] {
	return IOE.Of[error](KeywordProcessingResult{Created: false})
}

func updatePlasmidProperty(
	ctx WithPlasmid,
) IOE.IOEither[error, KeywordProcessingResult] {
	input := &stockpb.PlasmidUpdate{
		Data: &stockpb.PlasmidUpdate_Data{
			Type: "plasmid",
			Id:   ctx.Record.PlasmidID,
			Attributes: &stockpb.PlasmidUpdateAttributes{
				UpdatedBy:            regsc.DefaultUser,
				DictyPlasmidProperty: ctx.Record.Term,
			},
		},
	}
	return F.Pipe2(
		IOE.TryCatchError(func() (*stockpb.Plasmid, error) {
			return regsc.GetStockAPIClient().UpdatePlasmid(
				context.Background(),
				input,
			)
		}),
		IOE.MapLeft[*stockpb.Plasmid](func(err error) error {
			return fmt.Errorf(
				"failed to update plasmid %s: %w",
				ctx.Record.PlasmidID,
				err,
			)
		}),
		IOE.Map[error](func(_ *stockpb.Plasmid) KeywordProcessingResult {
			return KeywordProcessingResult{
				Created: true,
				Error:   nil,
			}
		}),
	)
}

func associateKeywordTerm(
	keyRecord KeywordRecord,
) E.Either[error, KeywordProcessingResult] {
	return F.Pipe3(
		IOE.Of[error](KeywordContext{Record: keyRecord}),
		IOE.Bind(SetPlasmid, fetchPlasmid),
		IOE.Chain(
			F.Ternary(
				propertyAlreadyMatches,
				returnExistingResult,
				updatePlasmidProperty,
			),
		),
		fputil.ToEither,
	)
}

func SummarySemigroup() S.Semigroup[KeywordProcessingSummary] {
	return S.MakeSemigroup(
		func(a, b KeywordProcessingSummary) KeywordProcessingSummary {
			combinedErrors := append(a.Errors, b.Errors...)

			// Cap at first maxErrorMessages errors
			if len(combinedErrors) > maxErrorMessages {
				combinedErrors = combinedErrors[:maxErrorMessages]
			}

			return KeywordProcessingSummary{
				CreatedCount:  a.CreatedCount + b.CreatedCount,
				ExistingCount: a.ExistingCount + b.ExistingCount,
				ErrorCount:    a.ErrorCount + b.ErrorCount,
				Errors:        combinedErrors,
			}
		},
	)
}

func resultToSummary(
	result KeywordProcessingResult,
) KeywordProcessingSummary {
	return F.Pipe1(
		result,
		F.Ternary(
			hasError,
			F.Ternary(
				isWrongPropertyError,
				skipSummary,
				errorSummary,
			),
			F.Ternary(isCreated, createdSummary, existingSummary),
		),
	)
}

func keywordSkipRecordError(_ []string) error {
	return errSkipRecord
}
