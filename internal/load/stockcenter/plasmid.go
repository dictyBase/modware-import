package stockcenter

import (
	"context"
	"errors"
	"fmt"

	E "github.com/IBM/fp-go/either"
	fperrors "github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	P "github.com/IBM/fp-go/predicate"
	SG "github.com/IBM/fp-go/semigroup"
	S "github.com/IBM/fp-go/string"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var notEmpty = P.Not(S.IsEmpty)

// PlasmidEnv contains all dependencies for plasmid loading
type PlasmidEnv struct {
	Logger    *logrus.Entry
	APIClient pb.StockServiceClient
}

// ProcessingResult accumulates successes and errors
type ProcessingResult struct {
	Successes  []string // Successfully processed plasmid IDs
	Errors     error    // All validation/processing errors
	ErrorCount int      // Number of errors (for accurate statistics)
}

// ProcessingResultSemigroup for combining results
var ProcessingResultSemigroup = SG.MakeSemigroup(
	func(a, b ProcessingResult) ProcessingResult {
		return ProcessingResult{
			Successes:  append(a.Successes, b.Successes...),
			Errors:     errors.Join(a.Errors, b.Errors),
			ErrorCount: a.ErrorCount + b.ErrorCount,
		}
	},
)

// LookupTuple holds all lookup instances
type LookupTuple struct {
	Annotator  source.StockAnnotatorLookup
	PubLookup  source.StockPubLookup
	GeneLookup source.StockGeneLookup
}

// Context types for Do/Bind pattern
type LookupContext struct{}

type WithAnnotator struct {
	LookupContext
	Annotator source.StockAnnotatorLookup
}

type WithPubLookup struct {
	WithAnnotator
	PubLookup source.StockPubLookup
}

type WithGeneLookup struct {
	WithPubLookup
	GeneLookup source.StockGeneLookup
}

// Curried setters for Do/Bind pattern
var (
	SetAnnotator = F.Curry2(
		func(annotator source.StockAnnotatorLookup, ctx LookupContext) WithAnnotator {
			return WithAnnotator{
				LookupContext: ctx,
				Annotator:     annotator,
			}
		},
	)

	SetPubLookup = F.Curry2(
		func(pubLookup source.StockPubLookup, ctx WithAnnotator) WithPubLookup {
			return WithPubLookup{
				WithAnnotator: ctx,
				PubLookup:     pubLookup,
			}
		},
	)

	SetGeneLookup = F.Curry2(
		func(geneLookup source.StockGeneLookup, ctx WithPubLookup) WithGeneLookup {
			return WithGeneLookup{
				WithPubLookup: ctx,
				GeneLookup:    geneLookup,
			}
		},
	)
)

// initAnnotatorLookup creates annotator lookup as IOEither
func initAnnotatorLookup(
	ctx LookupContext,
) IOE.IOEither[error, source.StockAnnotatorLookup] {
	return IOE.TryCatchError(func() (source.StockAnnotatorLookup, error) {
		lookup, err := source.NewStockAnnotatorLookup(
			registry.GetReader(regs.PlasmidAnnotatorReader),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"error in opening annotation source: %w",
				err,
			)
		}
		return lookup, nil
	})
}

// initPubLookup creates publication lookup as IOEither
func initPubLookup(
	ctx WithAnnotator,
) IOE.IOEither[error, source.StockPubLookup] {
	return IOE.TryCatchError(func() (source.StockPubLookup, error) {
		lookup, err := source.NewStockPubLookup(
			registry.GetReader(regs.PlasmidPubReader),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"error in opening publication source: %w",
				err,
			)
		}
		return lookup, nil
	})
}

// initGeneLookup creates gene lookup as IOEither
func initGeneLookup(
	ctx WithPubLookup,
) IOE.IOEither[error, source.StockGeneLookup] {
	return IOE.TryCatchError(func() (source.StockGeneLookup, error) {
		lookup, err := source.NewStockGeneLookp(
			registry.GetReader(regs.PlasmidGeneReader),
		)
		if err != nil {
			return nil, fmt.Errorf("error in opening gene source: %w", err)
		}
		return lookup, nil
	})
}

// initAllLookups initializes all lookups using Do/Bind pattern
func initAllLookups() IOE.IOEither[error, LookupTuple] {
	return F.Pipe5(
		IOE.Do[error](LookupContext{}),
		IOE.Bind(SetAnnotator, initAnnotatorLookup),
		IOE.Bind(SetPubLookup, initPubLookup),
		IOE.Bind(SetGeneLookup, initGeneLookup),
		IOE.Map[error](func(ctx WithGeneLookup) LookupTuple {
			return LookupTuple{
				Annotator:  ctx.Annotator,
				PubLookup:  ctx.PubLookup,
				GeneLookup: ctx.GeneLookup,
			}
		}),
		IOE.MapLeft[LookupTuple](
			fperrors.OnError("failed to initialize lookups"),
		),
	)
}

// createPlasmidReader creates CSV reader with lookups
func createPlasmidReader(
	lookups LookupTuple,
) IOE.IOEither[error, source.PlasmidReader] {
	return IOE.TryCatchError(func() (source.PlasmidReader, error) {
		reader := source.NewCsvPlasmidReader(
			registry.GetReader(regs.PlasmidReader),
			lookups.Annotator,
			lookups.PubLookup,
			lookups.GeneLookup,
		)
		return reader, nil
	})
}

// hasNextRecord checks if more records exist
func hasNextRecord(reader source.PlasmidReader) IOE.IOEither[error, bool] {
	return IOE.Of[error](reader.Next())
}

// readPlasmidRecord reads one record as IOEither using FromEither
func readPlasmidRecord(
	reader source.PlasmidReader,
) IOE.IOEither[error, *source.Plasmid] {
	return IOE.FromEither(reader.Value())
}

func hasValidPlasmidUser(plasmid *source.Plasmid) bool {
	return len(plasmid.User) == 0
}

func plasmidUserError(plasmid *source.Plasmid) error {
	return fmt.Errorf(
		"field User: user assignment required for plasmid %s",
		plasmid.Id,
	)
}

func hasValidPlasmidId(plasmid *source.Plasmid) bool {
	return notEmpty(plasmid.Id)
}

func plasmidIdError(plasmid *source.Plasmid) error {
	return fmt.Errorf("field Id: required")
}

func hasValidPlasmidName(plasmid *source.Plasmid) bool {
	return notEmpty(plasmid.Name)
}

func plasmidNameError(plasmid *source.Plasmid) error {
	return fmt.Errorf("field Name: required")
}

// validatePlasmid performs all validation checks
func validatePlasmid(
	plasmid *source.Plasmid,
) E.Either[error, *source.Plasmid] {
	return F.Pipe4(
		E.Of[error](plasmid),
		E.Chain(E.FromPredicate(hasValidPlasmidUser, plasmidUserError)),
		E.Chain(E.FromPredicate(hasValidPlasmidId, plasmidIdError)),
		E.Chain(E.FromPredicate(hasValidPlasmidName, plasmidNameError)),
		E.MapLeft[*source.Plasmid](func(err error) error {
			return fmt.Errorf("plasmid %s: %w", plasmid.Id, err)
		}),
	)
}

// checkPlasmidExists wraps gRPC GetPlasmid call in IOEither
func checkPlasmidExists(env PlasmidEnv, id string) IOE.IOEither[error, bool] {
	return IOE.TryCatchError(func() (bool, error) {
		_, err := env.APIClient.GetPlasmid(
			context.Background(),
			&pb.StockId{Id: id},
		)
		if err != nil && status.Code(err) == codes.NotFound {
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf(
				"error checking plasmid existence: %w",
				err,
			)
		}
		return true, nil
	})
}

// createPlasmidAPI wraps LoadPlasmid API call
func createPlasmidAPI(
	env PlasmidEnv,
	plasmid *source.Plasmid,
) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		attr := populateExistingPlasmidAttributes(env.Logger, plasmid)
		_, err := env.APIClient.LoadPlasmid(
			context.Background(),
			&pb.ExistingPlasmid{
				Data: &pb.ExistingPlasmid_Data{
					Type:       "plasmid",
					Id:         plasmid.Id,
					Attributes: attr,
				},
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"failed to create plasmid %s: %w",
				plasmid.Id,
				err,
			)
		}
		env.Logger.Debugf("created plasmid %s", plasmid.Id)
		return plasmid.Id, nil
	})
}

// updatePlasmidAPI wraps UpdatePlasmid API call
func updatePlasmidAPI(
	env PlasmidEnv,
	plasmid *source.Plasmid,
) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		attr := populatePlasmidUpdateAttributes(env.Logger, plasmid)
		_, err := env.APIClient.UpdatePlasmid(
			context.Background(),
			&pb.PlasmidUpdate{
				Data: &pb.PlasmidUpdate_Data{
					Type:       "plasmid",
					Id:         plasmid.Id,
					Attributes: attr,
				},
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"failed to update plasmid %s: %w",
				plasmid.Id,
				err,
			)
		}
		env.Logger.Debugf("updated plasmid %s", plasmid.Id)
		return plasmid.Id, nil
	})
}

// mergeResult is a curried function that merges new result into accumulator
var mergeResult = F.Curry2(
	func(acc ProcessingResult, newResult ProcessingResult) ProcessingResult {
		return ProcessingResultSemigroup.Concat(acc, newResult)
	},
)

// validateAndProcessPlasmid performs validation and API operations
func validateAndProcessPlasmid(
	env PlasmidEnv,
	plasmid *source.Plasmid,
) IOE.IOEither[error, ProcessingResult] {
	return IOE.TryCatchError(func() (ProcessingResult, error) {
		validation := validatePlasmid(plasmid)
		// Handle validation result
		result := E.Fold(
			func(err error) ProcessingResult {
				// Log validation error
				env.Logger.Errorf("%v", err)
				return ProcessingResult{
					Successes:  []string{},
					Errors:     err,
					ErrorCount: 1,
				}
			},
			func(validPlasmid *source.Plasmid) ProcessingResult {
				// Process the plasmid through API
				processEither := F.Pipe2(
					checkPlasmidExists(env, validPlasmid.Id),
					IOE.Chain(func(exists bool) IOE.IOEither[error, string] {
						if exists {
							return updatePlasmidAPI(env, validPlasmid)
						}
						return createPlasmidAPI(env, validPlasmid)
					}),
					IOE.Map[error](func(id string) ProcessingResult {
						// Success case
						return ProcessingResult{
							Successes:  []string{id},
							Errors:     nil,
							ErrorCount: 0,
						}
					}),
				)()

				// Convert API errors to ProcessingResult
				return E.GetOrElse(func(err error) ProcessingResult {
					env.Logger.Errorf(
						"error processing plasmid %s: %v",
						validPlasmid.Id,
						err,
					)
					return ProcessingResult{
						Successes: []string{},
						Errors: fmt.Errorf(
							"plasmid %s: %w",
							validPlasmid.Id,
							err,
						),
						ErrorCount: 1,
					}
				})(processEither)
			},
		)(validation)

		return result, nil
	})
}

// processSingleRecord processes one record and recurses
func processSingleRecord(
	env PlasmidEnv,
	reader source.PlasmidReader,
	acc ProcessingResult,
) IOE.IOEither[error, ProcessingResult] {
	return F.Pipe1(
		readPlasmidRecord(reader),
		IOE.Chain(
			func(plasmid *source.Plasmid) IOE.IOEither[error, ProcessingResult] {
				return F.Pipe1(
					validateAndProcessPlasmid(env, plasmid),
					IOE.Chain(
						func(result ProcessingResult) IOE.IOEither[error, ProcessingResult] {
							newAcc := mergeResult(acc)(result)
							return processNextRecord(env, reader, newAcc)
						},
					),
				)
			},
		),
	)
}

// processNextRecord is the tail-recursive processor
func processNextRecord(
	env PlasmidEnv,
	reader source.PlasmidReader,
	acc ProcessingResult,
) IOE.IOEither[error, ProcessingResult] {
	return F.Pipe2(
		hasNextRecord(reader),
		IOE.Chain(func(hasNext bool) IOE.IOEither[error, ProcessingResult] {
			if !hasNext {
				// Base case: no more records, return accumulated result
				return IOE.Of[error](acc)
			}
			// Recursive case: process one record and continue
			return processSingleRecord(env, reader, acc)
		}),
		IOE.MapLeft[ProcessingResult](
			fperrors.OnError("stream processing failed"),
		),
	)
}

// streamProcessRecords processes CSV records recursively
func streamProcessRecords(
	env PlasmidEnv,
	reader source.PlasmidReader,
) IOE.IOEither[error, ProcessingResult] {
	return processNextRecord(env, reader, ProcessingResult{
		Successes:  []string{},
		Errors:     nil,
		ErrorCount: 0,
	})
}

// loadPlasmidWorkflow orchestrates the entire loading process
func loadPlasmidWorkflow(env PlasmidEnv) IOE.IOEither[error, ProcessingResult] {
	return F.Pipe2(
		initAllLookups(),
		IOE.Chain(createPlasmidReader),
		IOE.Chain(
			func(reader source.PlasmidReader) IOE.IOEither[error, ProcessingResult] {
				return streamProcessRecords(env, reader)
			},
		),
	)
}

// buildPlasmidEnv constructs environment from registry
func buildPlasmidEnv() PlasmidEnv {
	return PlasmidEnv{
		Logger:    registry.GetLogger(),
		APIClient: regs.GetStockAPIClient(),
	}
}

// logFinalStats logs the final processing statistics
func logFinalStats(logger *logrus.Entry, result ProcessingResult) {
	successCount := len(result.Successes)
	errorCount := result.ErrorCount

	logger.WithFields(
		logrus.Fields{
			"type":    "annotations",
			"stock":   "plasmid",
			"event":   "load",
			"success": successCount,
			"errors":  errorCount,
			"total":   successCount + errorCount,
		},
	).Infof(
		"loaded plasmid annotations: %d succeeded, %d failed",
		successCount,
		errorCount,
	)

	// Log errors if any
	if result.Errors != nil {
		logger.Warnf("validation/processing errors: %v", result.Errors)
	}
}

// LoadPlasmid is the main entry point for plasmid loading
func LoadPlasmid(cmd *cobra.Command, args []string) error {
	env := buildPlasmidEnv()

	// Execute the workflow
	resultIO := loadPlasmidWorkflow(env)

	// Run the IOEither to get Either result
	resultEither := resultIO()

	// Handle the result
	return E.Fold(
		// Error case: workflow failed
		func(err error) error {
			env.Logger.Errorf("plasmid loading workflow failed: %v", err)
			return err
		},
		// Success case: log stats and return
		func(result ProcessingResult) error {
			logFinalStats(env.Logger, result)

			// If there were any errors, we still return nil to indicate
			// the workflow completed (errors were logged)
			return nil
		},
	)(resultEither)
}

// populateExistingPlasmidAttributes creates attributes for new plasmid
func populateExistingPlasmidAttributes(
	logger *logrus.Entry,
	plasmid *source.Plasmid,
) *pb.ExistingPlasmidAttributes {
	attr := &pb.ExistingPlasmidAttributes{
		CreatedAt: TimestampProto(plasmid.CreatedOn),
		UpdatedAt: TimestampProto(plasmid.UpdatedOn),
		CreatedBy: plasmid.User,
		Summary:   plasmid.Summary,
		Name:      plasmid.Name,
	}
	checkPublicationsAndGenes(logger, plasmid, attr)
	return attr
}

// populatePlasmidUpdateAttributes creates attributes for plasmid update
func populatePlasmidUpdateAttributes(
	logger *logrus.Entry,
	plasmid *source.Plasmid,
) *pb.PlasmidUpdateAttributes {
	attr := &pb.PlasmidUpdateAttributes{
		UpdatedBy: plasmid.User,
		Summary:   plasmid.Summary,
		Name:      plasmid.Name,
	}
	checkPublicationsAndGenes(logger, plasmid, attr)
	return attr
}

// checkPublicationsAndGenes populates publication and gene data
func checkPublicationsAndGenes(
	logger *logrus.Entry,
	plasmid *source.Plasmid,
	attr any,
) {
	switch a := attr.(type) {
	case *pb.ExistingPlasmidAttributes:
		if len(plasmid.Publications) > 0 {
			a.Publications = plasmid.Publications
		} else {
			logger.Warnf("plasmid %s has no publication entry", plasmid.Id)
		}
		if len(plasmid.Genes) > 0 {
			a.Genes = plasmid.Genes
		}
	case *pb.PlasmidUpdateAttributes:
		if len(plasmid.Publications) > 0 {
			a.Publications = plasmid.Publications
		} else {
			logger.Warnf("plasmid %s has no publication entry", plasmid.Id)
		}
		if len(plasmid.Genes) > 0 {
			a.Genes = plasmid.Genes
		}
	}
}
