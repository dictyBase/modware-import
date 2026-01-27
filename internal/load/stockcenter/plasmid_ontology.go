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
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
	File "github.com/IBM/fp-go/ioeither/file"
	O "github.com/IBM/fp-go/option"
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
	PlasmidID string
	Term      string
	Created   bool
	Processed bool
	Error     O.Option[error]
}

type KeywordProcessingSummary struct {
	Created    []string `json:"created"`
	Existing   []string `json:"existing"`
	Skipped    int      `json:"skipped"`
	ErrorCount int      `json:"error_count"`
	Errors     []string `json:"errors"`
}

var (
	errSkipRecord = errors.New("skip record")

	keywordIsNotBlankRecord = A.Any(func(s string) bool {
		return len(strings.TrimSpace(s)) > 0
	})
)

const keywordRecordLength = 3

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
		"created", len(summary.Created),
		"existing", len(summary.Existing),
		"skipped", summary.Skipped,
		"errors", summary.ErrorCount,
	)
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
		keywordFn := keywordFilterProperty(prop)
		summary := KeywordProcessingSummary{}
		for {
			record, err := config.Reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return summary, fmt.Errorf("tsv read error: %w", err)
			}

			result := F.Pipe6(
				record,
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
				E.Chain(keywordFn),
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
			return errSkipRecord
		},
	)
}

func handleKeywordPipelineError(err error) KeywordProcessingResult {
	if errors.Is(err, errSkipRecord) {
		return KeywordProcessingResult{
			Processed: false,
			Error:     O.None[error](),
		}
	}
	return KeywordProcessingResult{
		Processed: false,
		Error:     O.Some(err),
	}
}

func keywordHasValidRecordLength(record []string) bool {
	return len(record) == keywordRecordLength
}

func keywordRecordLengthError(record []string) error {
	return fmt.Errorf(
		"invalid record: expected exactly %d columns, got %d (%s)",
		keywordRecordLength,
		len(record),
		strings.Join(record, "\t"),
	)
}

func parseKeywordRecord(record []string) KeywordRecord {
	return KeywordRecord{
		PlasmidID: strings.TrimSpace(record[0]),
		Property:  strings.ToLower(strings.TrimSpace(record[1])),
		Term:      strings.TrimSpace(record[2]),
	}
}

func associateKeywordTerm(
	record KeywordRecord,
) E.Either[error, KeywordProcessingResult] {
	fn := IOE.TryCatchError(func() (KeywordProcessingResult, error) {
		client := regsc.GetStockAPIClient()
		plasmid, err := client.GetPlasmid(
			context.Background(),
			&stockpb.StockId{Id: record.PlasmidID},
		)
		if err != nil {
			return KeywordProcessingResult{}, fmt.Errorf(
				"failed to fetch plasmid %s: %w",
				record.PlasmidID,
				err,
			)
		}
		existing := plasmid.GetData().GetAttributes().GetDictyPlasmidProperty()
		if strings.EqualFold(existing, record.Term) {
			return keywordCreateSuccessResult(
				record.PlasmidID,
				record.Term,
				false,
			), nil
		}

		_, err = client.UpdatePlasmid(
			context.Background(),
			&stockpb.PlasmidUpdate{
				Data: &stockpb.PlasmidUpdate_Data{
					Type: "plasmid",
					Id:   record.PlasmidID,
					Attributes: &stockpb.PlasmidUpdateAttributes{
						UpdatedBy:            regsc.DefaultUser,
						DictyPlasmidProperty: record.Term,
					},
				},
			},
		)
		if err != nil {
			return KeywordProcessingResult{}, fmt.Errorf(
				"failed to update plasmid %s: %w",
				record.PlasmidID,
				err,
			)
		}

		return keywordCreateSuccessResult(
			record.PlasmidID,
			record.Term,
			true,
		), nil
	})
	return fputil.ToEither(fn)
}

func keywordCreateSuccessResult(
	id, term string,
	created bool,
) KeywordProcessingResult {
	return KeywordProcessingResult{
		PlasmidID: id,
		Term:      term,
		Created:   created,
		Processed: true,
		Error:     O.None[error](),
	}
}

func SummarySemigroup() S.Semigroup[KeywordProcessingSummary] {
	return S.MakeSemigroup(
		func(a, b KeywordProcessingSummary) KeywordProcessingSummary {
			return KeywordProcessingSummary{
				Created:    append(a.Created, b.Created...),
				Existing:   append(a.Existing, b.Existing...),
				Skipped:    a.Skipped + b.Skipped,
				ErrorCount: a.ErrorCount + b.ErrorCount,
				Errors:     append(a.Errors, b.Errors...),
			}
		},
	)
}

func resultToSummary(
	result KeywordProcessingResult,
) KeywordProcessingSummary {
	return F.Pipe1(
		result.Error,
		O.Fold(
			func() KeywordProcessingSummary {
				switch {
				case !result.Processed:
					return KeywordProcessingSummary{Skipped: 1}
				case result.Created:
					return KeywordProcessingSummary{
						Created: []string{keywordFormatAssociation(result)},
					}
				default:
					return KeywordProcessingSummary{
						Existing: []string{keywordFormatAssociation(result)},
					}
				}
			},
			func(err error) KeywordProcessingSummary {
				return KeywordProcessingSummary{
					ErrorCount: 1,
					Errors:     []string{err.Error()},
				}
			},
		),
	)
}

func keywordFormatAssociation(result KeywordProcessingResult) string {
	return F.Pipe2(
		result.Term,
		O.FromPredicate(func(s string) bool {
			return len(strings.TrimSpace(s)) > 0
		}),
		O.Fold(
			func() string { return result.PlasmidID },
			func(term string) string {
				return fmt.Sprintf(
					"%s:%s",
					result.PlasmidID,
					term,
				)
			},
		),
	)
}

func keywordSkipRecordError(_ []string) error {
	return errSkipRecord
}
