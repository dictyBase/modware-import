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
	N "github.com/IBM/fp-go/number"
	S "github.com/IBM/fp-go/semigroup"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/datasource/s3"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/logger"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/urfave/cli/v2"
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

	// Field-level semigroups for composable aggregation
	intSumSemigroup = N.SemigroupSum[int]()
	errorsSemigroup = S.MakeSemigroup(func(a, b []error) []error {
		return F.Pipe1(
			A.ArrayConcatAll(a, b),
			A.Slice[error](0, maxErrorMessages),
		)
	})
)

type plasmidOntologyLoadResult struct {
	Summary KeywordProcessingSummary
	Error   error
}

type KeywordLoaderContext struct{}

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

func configureTSVReader(reader io.Reader) *csv.Reader {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	return csvReader
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
			return KeywordProcessingSummary{
				CreatedCount: intSumSemigroup.Concat(
					a.CreatedCount,
					b.CreatedCount,
				),
				ExistingCount: intSumSemigroup.Concat(
					a.ExistingCount,
					b.ExistingCount,
				),
				ErrorCount: intSumSemigroup.Concat(
					a.ErrorCount,
					b.ErrorCount,
				),
				Errors: errorsSemigroup.Concat(
					a.Errors,
					b.Errors,
				),
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

// ===== CLI Context Support =====

type KeywordLoaderCliConfig struct {
	KeywordLoaderContext
	Cmd    *cli.Context
	Reader *csv.Reader
	Closer io.Closer
}

type KeywordReaderCliIOE = IOE.IOEither[error, KeywordLoaderCliConfig]

func openKeywordReaderCli(config KeywordLoaderCliConfig) KeywordReaderCliIOE {
	return F.Pipe1(
		config,
		F.Switch(
			func(cfg KeywordLoaderCliConfig) string {
				return cfg.Cmd.String("input-source")
			},
			map[string]func(KeywordLoaderCliConfig) KeywordReaderCliIOE{
				"folder": keywordReaderFromFileCli,
				"bucket": keywordReaderFromS3BucketCli,
			},
			defaultKeywordReaderCli,
		),
	)
}

func keywordReaderFromFileCli(cfg KeywordLoaderCliConfig) KeywordReaderCliIOE {
	inputPath := cfg.Cmd.String("input")
	return F.Pipe2(
		File.Open(inputPath),
		IOE.MapLeft[*os.File](func(err error) error {
			return fmt.Errorf("failed to open TSV file %s: %w", inputPath, err)
		}),
		IOE.Map[error](func(file *os.File) KeywordLoaderCliConfig {
			cfg.Reader = configureTSVReader(file)
			cfg.Closer = file
			return cfg
		}),
	)
}

// S3ReaderContext holds state accumulated through the S3 reader pipeline
type S3ReaderContext struct {
	Cmd        *cli.Context
	Bucket     string
	ObjectPath string
	Client     *minio.Client
	Reader     io.ReadCloser
}

// setS3Client returns a function that sets the S3 client in the context
func setS3Client(client *minio.Client) func(S3ReaderContext) S3ReaderContext {
	return func(ctx S3ReaderContext) S3ReaderContext {
		ctx.Client = client
		return ctx
	}
}

// setS3Reader returns a function that sets the reader in the context
func setS3Reader(reader io.ReadCloser) func(S3ReaderContext) S3ReaderContext {
	return func(ctx S3ReaderContext) S3ReaderContext {
		ctx.Reader = reader
		return ctx
	}
}

// createS3ClientIOE creates the S3 client from CLI context
func createS3ClientIOE(ctx S3ReaderContext) IOE.IOEither[error, *minio.Client] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*minio.Client, error) {
			return s3.NewCliS3Client(ctx.Cmd)
		}),
		IOE.MapLeft[*minio.Client](func(err error) error {
			return fmt.Errorf("failed to create S3 client: %w", err)
		}),
	)
}

// statS3ObjectIOE validates object exists and credentials are valid
func statS3ObjectIOE(ctx S3ReaderContext) IOE.IOEither[error, minio.ObjectInfo] {
	return F.Pipe1(
		IOE.TryCatchError(func() (minio.ObjectInfo, error) {
			return ctx.Client.StatObject(
				ctx.Bucket,
				ctx.ObjectPath,
				minio.StatObjectOptions{},
			)
		}),
		IOE.MapLeft[minio.ObjectInfo](func(err error) error {
			return fmt.Errorf("failed to stat s3 object %s: %w", ctx.ObjectPath, err)
		}),
	)
}

// getS3ObjectReaderIOE fetches the object stream
func getS3ObjectReaderIOE(ctx S3ReaderContext) IOE.IOEither[error, io.ReadCloser] {
	return F.Pipe1(
		IOE.TryCatchError(func() (io.ReadCloser, error) {
			return ctx.Client.GetObject(
				ctx.Bucket,
				ctx.ObjectPath,
				minio.GetObjectOptions{},
			)
		}),
		IOE.MapLeft[io.ReadCloser](func(err error) error {
			return fmt.Errorf("failed to open s3 file %s: %w", ctx.ObjectPath, err)
		}),
	)
}

func keywordReaderFromS3BucketCli(cfg KeywordLoaderCliConfig) KeywordReaderCliIOE {
	path := cfg.Cmd.String("s3-bucket-path")
	file := cfg.Cmd.String("input")
	bucket := cfg.Cmd.String("s3-bucket")
	objectPath := fmt.Sprintf("%s/%s", path, file)

	initialCtx := S3ReaderContext{
		Cmd:        cfg.Cmd,
		Bucket:     bucket,
		ObjectPath: objectPath,
	}

	// Flat pipeline: Do -> Bind(client) -> ChainFirst(stat) -> Bind(reader) -> Map(config)
	return F.Pipe4(
		IOE.Do[error](initialCtx),
		IOE.Bind(setS3Client, createS3ClientIOE),
		IOE.ChainFirst(statS3ObjectIOE),
		IOE.Bind(setS3Reader, getS3ObjectReaderIOE),
		IOE.Map[error](func(ctx S3ReaderContext) KeywordLoaderCliConfig {
			cfg.Reader = configureTSVReader(ctx.Reader)
			cfg.Closer = ctx.Reader
			return cfg
		}),
	)
}

func defaultKeywordReaderCli(cfg KeywordLoaderCliConfig) KeywordReaderCliIOE {
	source := cfg.Cmd.String("input-source")
	return IOE.Left[KeywordLoaderCliConfig](
		fmt.Errorf("unsupported input source %s", source),
	)
}

func processAndAggregateKeywordRecordsCli(
	config KeywordLoaderCliConfig,
) IOE.IOEither[error, KeywordProcessingSummary] {
	return IOE.TryCatchError(func() (KeywordProcessingSummary, error) {
		if config.Closer != nil {
			defer config.Closer.Close()
		}

		prop := config.Cmd.String("property")
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

func onPlasmidOntologyError(err error) plasmidOntologyLoadResult {
	return plasmidOntologyLoadResult{Error: err}
}

func onPlasmidOntologySuccess(
	summary KeywordProcessingSummary,
) plasmidOntologyLoadResult {
	return plasmidOntologyLoadResult{Summary: summary}
}

func handlePlasmidOntologyOutput(
	result plasmidOntologyLoadResult,
	slogger *slog.Logger,
) error {
	if result.Error != nil {
		return result.Error
	}

	slogger.Info(
		"plasmid ontology summary",
		"created", result.Summary.CreatedCount,
		"existing", result.Summary.ExistingCount,
		"errors", result.Summary.ErrorCount,
	)

	// Log error samples if present
	if len(result.Summary.Errors) > 0 {
		joinedErrors := errors.Join(result.Summary.Errors...)
		slogger.Warn(
			"error samples (first 5)",
			"sample_count", len(result.Summary.Errors),
			"total_errors", result.Summary.ErrorCount,
			"messages", joinedErrors.Error(),
		)
	}

	return nil
}

func LoadPlasmidOntologyCli(cmd *cli.Context) error {
	handler := logger.GetCliSlogHandler(cmd)
	slogger := slog.New(handler)

	output := F.Pipe1(
		F.Pipe4(
			IOE.Do[error](KeywordLoaderCliConfig{Cmd: cmd}),
			IOE.Chain(openKeywordReaderCli),
			IOE.Chain(processAndAggregateKeywordRecordsCli),
			IOE.ChainFirstIOK[error](
				IO.Logf[KeywordProcessingSummary](
					"Plasmid ontology association summary: %+v",
				),
			),
			fputil.ToEither[error, KeywordProcessingSummary],
		),
		E.Fold(onPlasmidOntologyError, onPlasmidOntologySuccess),
	)

	return handlePlasmidOntologyOutput(output, slogger)
}
