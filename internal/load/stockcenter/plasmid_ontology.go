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
	fperrors "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/logger"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
)

type KeywordLoaderContext struct{}

type KeywordReaderResource struct {
	Reader *csv.Reader
	Closer io.Closer
}

type KeywordLoaderConfig struct {
	KeywordLoaderContext
	Cmd    *cobra.Command
	Reader *csv.Reader
	Closer io.Closer
}

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
	Err        error    `json:"-"`
}

var (
	setKeywordReader = F.Curry2(
		func(resource KeywordReaderResource, config KeywordLoaderConfig) KeywordLoaderConfig {
			config.Reader = resource.Reader
			config.Closer = resource.Closer
			return config
		},
	)
	errSkipRecord = errors.New("skip record")
)

const keywordRecordLength = 3

func LoadPlasmidOntology(cmd *cobra.Command, _ []string) error {
	handler := logger.GetSlogHandler(cmd)
	slogger := slog.New(handler)
	elog := E.Logger[error, KeywordProcessingSummary](
		slog.NewLogLogger(handler, slog.LevelInfo),
		slog.NewLogLogger(handler, slog.LevelError),
	)

	return F.Pipe8(
		IOE.Do[error](KeywordLoaderConfig{Cmd: cmd}),
		IOE.ChainFirst(
			IOE.LogJSON[KeywordLoaderConfig](
				"Starting plasmid ontology association:\n%s",
			),
		),
		IOE.Bind(setKeywordReader, openKeywordReader),
		IOE.Chain(streamKeywordRecords),
		IOE.Map[error](aggregateKeywordResults),
		IOE.ChainFirst(
			IOE.LogJSON[KeywordProcessingSummary](
				"Plasmid ontology association summary:\n%s",
			),
		),
		fputil.ToEither[error, KeywordProcessingSummary],
		elog("plasmid ontology association result"),
		E.Fold(
			fperrors.IdentityError,
			func(summary KeywordProcessingSummary) error {
				slogger.Info(
					"plasmid ontology summary",
					"created", len(summary.Created),
					"existing", len(summary.Existing),
					"skipped", summary.Skipped,
					"errors", summary.ErrorCount,
				)
				if summary.ErrorCount > 0 {
					slogger.Error(
						"plasmid ontology errors",
						"errors",
						summary.Errors,
					)
				}
				return summary.Err
			},
		),
	)
}

func openKeywordReader(
	config KeywordLoaderConfig,
) IOE.IOEither[error, KeywordReaderResource] {
	return IOE.TryCatchError(func() (KeywordReaderResource, error) {
		inputPath, _ := config.Cmd.Flags().GetString("input")
		file, err := os.Open(inputPath)
		if err != nil {
			return KeywordReaderResource{}, fmt.Errorf(
				"failed to open TSV file %s: %w",
				inputPath,
				err,
			)
		}
		reader := csv.NewReader(file)
		reader.Comma = '\t'
		reader.FieldsPerRecord = -1
		reader.TrimLeadingSpace = true
		return KeywordReaderResource{Reader: reader, Closer: file}, nil
	})
}

func streamKeywordRecords(
	config KeywordLoaderConfig,
) IOE.IOEither[error, []KeywordProcessingResult] {
	return IOE.TryCatchError(func() ([]KeywordProcessingResult, error) {
		if config.Closer != nil {
			defer config.Closer.Close()
		}

		prop, _ := config.Cmd.Flags().GetString("property")
		keywordFn := keywordFilterProperty(prop)
		results := []KeywordProcessingResult{}
		for {
			record, err := config.Reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("tsv read error: %w", err)
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

			results = append(results, result)
		}

		return results, nil
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
		"invalid record: expected at least %d columns, got %d (%s)",
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
		existing := plasmid.Data.Attributes.DictyPlasmidProperty
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

func aggregateKeywordResults(
	results []KeywordProcessingResult,
) KeywordProcessingSummary {
	return F.Pipe1(
		results,
		A.Reduce(
			func(acc KeywordProcessingSummary, result KeywordProcessingResult) KeywordProcessingSummary {
				return F.Pipe1(
					result.Error,
					O.Fold(
						func() KeywordProcessingSummary {
							if !result.Processed {
								acc.Skipped++
								return acc
							}
							association := keywordFormatAssociation(result)
							if result.Created {
								acc.Created = append(acc.Created, association)
							} else {
								acc.Existing = append(acc.Existing, association)
							}
							return acc
						},
						func(err error) KeywordProcessingSummary {
							acc.ErrorCount++
							acc.Errors = append(acc.Errors, err.Error())
							acc.Err = errors.Join(acc.Err, err)
							return acc
						},
					),
				)
			},
			KeywordProcessingSummary{},
		),
	)
}

func keywordFormatAssociation(result KeywordProcessingResult) string {
	if strings.TrimSpace(result.Term) == "" {
		return result.PlasmidID
	}
	return fmt.Sprintf("%s:%s", result.PlasmidID, result.Term)
}

func keywordIsNotBlankRecord(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return true
		}
	}
	return false
}

func keywordSkipRecordError(_ []string) error {
	return errSkipRecord
}
