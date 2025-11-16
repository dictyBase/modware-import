package stockcenter

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	fperrors "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoaderContext is the initial empty context
type LoaderContext struct{}

// LoaderConfig holds viper configuration and CSV reader
type LoaderConfig struct {
	LoaderContext
	Viper  *viper.Viper
	Reader *csv.Reader
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

// GoldenBraidProcessingResult holds aggregate processing statistics
type GoldenBraidProcessingResult struct {
	Successes  []string
	Errors     error
	ErrorCount int
}

// openCSVReader opens a CSV file from config and returns a reader wrapped in IOEither
func openCSVReader(config LoaderConfig) IOE.IOEither[error, *csv.Reader] {
	return IOE.TryCatchError(func() (*csv.Reader, error) {
		inputPath := config.Viper.GetString("input")
		file, err := os.Open(config.Viper.GetString("input"))
		if err != nil {
			return nil, fmt.Errorf(
				"failed to open CSV file %s: %w",
				inputPath,
				err,
			)
		}

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = 7
		reader.TrimLeadingSpace = true

		// Skip header row
		if _, err := reader.Read(); err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to read CSV header: %w", err)
		}

		return reader, nil
	})
}

// streamAndProcessRecords reads CSV records one by one and processes them
func streamAndProcessRecords(
	config LoaderConfig,
) IOE.IOEither[error, []PlasmidProcessingResult] {
	return IOE.TryCatchError(func() ([]PlasmidProcessingResult, error) {
		userEmail := config.Viper.GetString("user-email")
		plasmidCVTerm := config.Viper.GetString("plasmid-cvterm")
		results := []PlasmidProcessingResult{}

		for {
			record, err := config.Reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("CSV read error: %w", err)
			}

			// Process single record (pure Either pipeline - integrated)
			result := F.Pipe4(
				record,
				E.FromPredicate(
					source.HasValidRecordLength,
					source.RecordLengthError,
				),
				E.Map[error](source.BuildPlasmid(
					userEmail,
					plasmidCVTerm,
				)),
				E.Chain(source.ValidatePlasmid),
				E.Fold(createErrorResult, processPlasmid),
			)

			results = append(results, result)
		}

		return results, nil
	})
}

// processPlasmid processes a single validated plasmid by loading it to the API
func processPlasmid(p *source.GoldenBraidPlasmid) PlasmidProcessingResult {
	// Execute IOEither to load plasmid to API
	result := loadPlasmidToAPI(p)()

	return E.Fold(
		func(err error) PlasmidProcessingResult {
			return PlasmidProcessingResult{
				PlasmidID: p.Name,
				Error:     O.Some(err),
			}
		},
		func(id string) PlasmidProcessingResult {
			return PlasmidProcessingResult{
				PlasmidID: id,
				Error:     O.None[error](),
			}
		},
	)(result)
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
					Genes:     O.GetOrElse(F.Constant([]string{}))(p.Genes),
					Publications: O.GetOrElse(
						F.Constant([]string{}),
					)(
						p.Publications,
					),
				},
			},
		}

		response, err := client.CreatePlasmid(context.Background(), plasmidData)
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

// aggregateResults aggregates processing results using A.Reduce
func aggregateResults(
	results []PlasmidProcessingResult,
) GoldenBraidProcessingResult {
	return F.Pipe1(
		results,
		A.Reduce(
			func(acc GoldenBraidProcessingResult, r PlasmidProcessingResult) GoldenBraidProcessingResult {
				return F.Pipe1(
					r.Error,
					O.Fold(
						// No error - success
						func() GoldenBraidProcessingResult {
							return GoldenBraidProcessingResult{
								Successes:  append(acc.Successes, r.PlasmidID),
								Errors:     acc.Errors,
								ErrorCount: acc.ErrorCount,
							}
						},
						// Error - accumulate
						func(err error) GoldenBraidProcessingResult {
							return GoldenBraidProcessingResult{
								Successes:  acc.Successes,
								Errors:     errors.Join(acc.Errors, err),
								ErrorCount: acc.ErrorCount + 1,
							}
						},
					),
				)
			},
			GoldenBraidProcessingResult{}, // Initial empty result
		),
	)
}

// toEither executes an IOEither to get an Either result
func toEither[A any](ioe IOE.IOEither[error, A]) E.Either[error, A] {
	return ioe()
}

// SetReader is a curried setter that adds CSV reader to config
var SetReader = F.Curry2(
	func(reader *csv.Reader, config LoaderConfig) LoaderConfig {
		return LoaderConfig{
			LoaderContext: config.LoaderContext,
			Viper:         config.Viper,
			Reader:        reader,
		}
	},
)

// LoadGoldenBraid is the main entry point for loading GoldenBraid CSV data
func LoadGoldenBraid(cmd *cobra.Command, args []string) error {
	errLogger := log.New(
		os.Stderr,
		"[ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	elog := E.Logger[error, GoldenBraidProcessingResult](errLogger, infoLogger)

	return F.Pipe9(
		IOE.Do[error](LoaderConfig{
			Viper:  viper.GetViper(),
			Reader: nil,
		}),
		IOE.ChainFirst(
			IOE.LogJSON[LoaderConfig](
				"Starting GoldenBraid loading with config:\n%s",
			),
		),
		IOE.Bind(SetReader, openCSVReader),
		IOE.ChainFirst(
			IOE.LogJSON[LoaderConfig](
				"CSV reader opened, processing records...\n%s",
			),
		),
		IOE.Chain(streamAndProcessRecords),
		IOE.Map[error](aggregateResults),
		IOE.ChainFirst(
			IOE.LogJSON[GoldenBraidProcessingResult](
				"Processing complete:\n%s",
			),
		),
		toEither[GoldenBraidProcessingResult],
		elog("GoldenBraid loading result"),
		E.Fold(
			fperrors.IdentityError,
			func(_ GoldenBraidProcessingResult) error { return nil },
		),
	)
}
