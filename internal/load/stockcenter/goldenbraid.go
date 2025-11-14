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
	F "github.com/IBM/fp-go/function"
	fperrors "github.com/IBM/fp-go/errors"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PlasmidProcessingResult represents the result of processing a single plasmid
type PlasmidProcessingResult struct {
	PlasmidID string
	Error     O.Option[error]
}

// ProcessingResult holds aggregate processing statistics
type ProcessingResult struct {
	Successes  []string
	Errors     error
	ErrorCount int
}

// openCSVReader opens a CSV file and returns a reader wrapped in IOEither
func openCSVReader(path string) IOE.IOEither[error, *csv.Reader] {
	return IOE.TryCatchError(func() (*csv.Reader, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open CSV file %s: %w", path, err)
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
	userEmail string,
	plasmidCVTerm string,
) func(*csv.Reader) IOE.IOEither[error, []PlasmidProcessingResult] {
	return func(reader *csv.Reader) IOE.IOEither[error, []PlasmidProcessingResult] {
		return IOE.TryCatchError(func() ([]PlasmidProcessingResult, error) {
			results := []PlasmidProcessingResult{}

			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, fmt.Errorf("CSV read error: %w", err)
				}

				// Process single record (pure Either pipeline)
				result := F.Pipe3(
					source.ParseRecord(record, userEmail, plasmidCVTerm),
					E.Chain(source.ValidatePlasmid),
					E.Fold(
						func(e error) PlasmidProcessingResult {
							return PlasmidProcessingResult{
								PlasmidID: "",
								Error:     O.Some(e),
							}
						},
						func(p *source.GoldenBraidPlasmid) PlasmidProcessingResult {
							return processPlasmid(p)
						},
					),
				)

				results = append(results, result)
			}

			return results, nil
		})
	}
}

// processPlasmid processes a single validated plasmid by loading it to the API
func processPlasmid(p *source.GoldenBraidPlasmid) PlasmidProcessingResult {
	// Execute IOEither to load plasmid to API
	result := F.Pipe1(
		loadPlasmidToAPI(p),
		IOE.Fold(
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
		),
	)()

	return result
}

// loadPlasmidToAPI loads a plasmid to the stock center API
func loadPlasmidToAPI(p *source.GoldenBraidPlasmid) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		client := registry.GetStockClient()

		plasmidData := &stock.StockPlasmid{
			Data: &stock.StockPlasmid_Data{
				Type: "plasmid",
				Attributes: &stock.PlasmidAttributes{
					Name:      p.Name,
					Summary:   p.Summary,
					CreatedBy: p.User,
					UpdatedBy: p.User,
					Plasmid:   p.PlasmidType,
					// Convert Option fields to slices
					Genes:        O.GetOrElse(F.Constant([]string{}))(p.Genes),
					Publications: O.GetOrElse(F.Constant([]string{}))(p.Publications),
				},
			},
		}

		response, err := client.CreatePlasmid(context.Background(), plasmidData)
		if err != nil {
			return "", fmt.Errorf("failed to create plasmid %s: %w", p.Name, err)
		}

		return response.Data.Id, nil
	})
}

// aggregateResults aggregates processing results using A.Reduce
func aggregateResults(results []PlasmidProcessingResult) ProcessingResult {
	return F.Pipe1(
		results,
		A.Reduce(
			func(acc ProcessingResult, r PlasmidProcessingResult) ProcessingResult {
				return F.Pipe1(
					r.Error,
					O.Fold(
						// No error - success
						func() ProcessingResult {
							return ProcessingResult{
								Successes:  append(acc.Successes, r.PlasmidID),
								Errors:     acc.Errors,
								ErrorCount: acc.ErrorCount,
							}
						},
						// Error - accumulate
						func(err error) ProcessingResult {
							return ProcessingResult{
								Successes:  acc.Successes,
								Errors:     errors.Join(acc.Errors, err),
								ErrorCount: acc.ErrorCount + 1,
							}
						},
					),
				)
			},
			ProcessingResult{}, // Initial empty result
		),
	)
}

// toEither executes an IOEither to get an Either result
func toEither[A any](ioe IOE.IOEither[error, A]) E.Either[error, A] {
	return ioe()
}

// LoadGoldenBraid is the main entry point for loading GoldenBraid CSV data
func LoadGoldenBraid(cmd *cobra.Command, args []string) error {
	userEmail := viper.GetString("user-email")
	plasmidCVTerm := viper.GetString("plasmid-cvterm")
	inputPath := viper.GetString("input")

	// Configure loggers
	errLogger := log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	elog := E.Logger[error, ProcessingResult](errLogger, infoLogger)

	return F.Pipe3(
		F.Pipe7(
			IOE.Of[error](inputPath),
			IOE.ChainFirst(IOE.LogJSON[string]("Starting GoldenBraid loading:\n%s")),
			IOE.Chain(openCSVReader),
			IOE.Chain(streamAndProcessRecords(userEmail, plasmidCVTerm)),
			IOE.Map[error](aggregateResults),
			IOE.ChainFirst(IOE.LogJSON[ProcessingResult]("Processing results:\n%s")),
		),
		toEither[ProcessingResult],
		elog("GoldenBraid loading result"),
		E.Fold(
			fperrors.IdentityError,
			func(_ ProcessingResult) error { return nil },
		),
	)
}
