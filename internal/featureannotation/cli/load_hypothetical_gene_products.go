package cli

import (
	"context"
	"fmt"
	"strings"

	P "github.com/IBM/fp-go/predicate"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/IBM/fp-go/ioeither/file"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// HypotheticalProteinProduct is the product name for hypothetical proteins
	HypotheticalProteinProduct = "conserved hypothetical protein"
)

// Point-free helper predicates
var (
	// stripBOM removes Unicode byte order mark from beginning of string
	stripBOM = func(s string) string {
		return strings.TrimPrefix(s, "\ufeff")
	}

	// isNotEmpty checks if trimmed string is not empty
	isNotEmpty = F.Flow2(
		strings.TrimSpace,
		S.IsNonEmpty,
	)

	// doesNotStartWithHash checks if trimmed string doesn't start with #
	doesNotStartWithHash = F.Flow2(
		strings.TrimSpace,
		func(s string) bool { return !strings.HasPrefix(s, "#") },
	)

	// isValidGeneLine combines predicates using monoid
	isValidGeneLine = P.SemigroupAll[string]().Concat(
		isNotEmpty,
		doesNotStartWithHash,
	)

	parseLine = F.Flow3(
		strings.TrimSpace,
		stripBOM,
		O.FromPredicate(isValidGeneLine),
	)

	// processGeneFile processes an entire file of gene lines
	// Returns only valid, non-comment gene IDs in uppercase
	processGeneFile = F.Flow3(
		func(cnt []byte) string { return string(cnt) },
		func(content string) []string { return strings.Split(content, "\n") },
		A.FilterMap(parseLine),
	)
)

// LoadHypotheticalGeneProductsParams holds parameters for loading hypothetical gene products
type LoadHypotheticalGeneProductsParams struct {
	User   string
	Logger *logrus.Entry
	Client feature.FeatureAnnotationServiceClient
}

// ProcessingConfig holds immutable configuration for processing gene IDs
type ProcessingConfig struct {
	GeneIDs []string
	User    string
	Logger  *logrus.Entry
	Client  feature.FeatureAnnotationServiceClient
}

// HypotheticalGeneProcessingResult represents the result of processing gene IDs
type HypotheticalGeneProcessingResult struct {
	Processed int
	Skipped   int
	Errors    int
}

// readGeneIDsFromFile reads gene IDs from a text file using IOEither
func readGeneIDsFromFile(filePath string) IOE.IOEither[error, []string] {
	return F.Pipe2(
		file.ReadFile(filePath),
		IOE.Map[error](processGeneFile),
		IOE.MapLeft[[]string](func(err error) error {
			return fmt.Errorf(
				"failed to read gene IDs from file %s: %w",
				filePath,
				err,
			)
		}),
	)
}

// createProcessingConfig creates an immutable ProcessingConfig from params and gene IDs
// This is a pure function that doesn't mutate anything
var createProcessingConfig = F.Curry2(
	func(
		params *LoadHypotheticalGeneProductsParams,
		geneIDs []string,
	) ProcessingConfig {
		return ProcessingConfig{
			GeneIDs: geneIDs,
			User:    params.User,
			Logger:  params.Logger,
			Client:  params.Client,
		}
	},
)

// processGeneID processes a single gene ID using the feature annotation API
// Uses F.Curry2 for proper currying and accepts ProcessingConfig
var processGeneID = F.Curry2(
	func(config ProcessingConfig, geneID string) E.Either[error, string] {
		ctx := context.Background()

		// Try to get existing annotation
		existing, err := config.Client.GetFeatureAnnotation(
			ctx,
			&pb.FeatureAnnotationId{Id: geneID},
		)
		if err != nil {
			return handleNewGeneProduct(config)(geneID)(err)
		}

		return handleExistingGeneProduct(config)(geneID)(existing)
	},
)

// hasHypotheticalProduct checks if a tag property is the hypothetical protein product
var hasHypotheticalProduct = func(tag *pb.TagProperty) bool {
	return tag.Tag == "product" && tag.Value == HypotheticalProteinProduct
}

// handleExistingGeneProduct handles the case where the gene annotation already exists
// Uses F.Curry3 for proper currying and accepts ProcessingConfig
var handleExistingGeneProduct = F.Curry3(
	func(
		config ProcessingConfig,
		geneID string,
		existing *pb.FeatureAnnotation,
	) E.Either[error, string] {
		// Check if the hypothetical protein product already exists using functional pattern
		productExists := F.Pipe2(
			existing.Attributes.Properties,
			A.FindFirst(hasHypotheticalProduct),
			O.IsSome[*pb.TagProperty],
		)

		if productExists {
			config.Logger.Debugf(
				"gene product '%s' already exists for %s, skipping",
				HypotheticalProteinProduct,
				geneID,
			)
			return E.Right[error](fmt.Sprintf("skipped:%s", geneID))
		}

		// Add the hypothetical protein product tag
		ctx := context.Background()
		updated, err := config.Client.AddTag(ctx, &pb.AddTagRequest{
			Id: geneID,
			Tag: &pb.TagPropertyCreate{
				Tag:       "product",
				Value:     HypotheticalProteinProduct,
				CreatedBy: config.User,
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			return E.Left[string](
				fmt.Errorf("failed to add tag for %s: %w", geneID, err),
			)
		}

		config.Logger.Debugf(
			"successfully added gene product for %s",
			updated.Id,
		)
		return E.Right[error](fmt.Sprintf("updated:%s", geneID))
	},
)

// handleNewGeneProduct handles the case where the gene annotation doesn't exist
// Uses F.Curry3 for proper currying and accepts ProcessingConfig
var handleNewGeneProduct = F.Curry3(
	func(
		config ProcessingConfig,
		geneID string,
		grpcErr error,
	) E.Either[error, string] {
		if status.Code(grpcErr) != codes.NotFound {
			return E.Left[string](
				fmt.Errorf(
					"error finding feature annotation for %s: %w",
					geneID,
					grpcErr,
				),
			)
		}

		// Create new feature annotation with hypothetical protein product
		nfa := &pb.NewFeatureAnnotation{
			Id:        geneID,
			CreatedBy: config.User,
			CreatedAt: timestamppb.Now(),
			Attributes: &pb.FeatureAnnotationAttributes{
				Name: geneID,
				Properties: []*pb.TagProperty{{
					Tag:       "product",
					Value:     HypotheticalProteinProduct,
					CreatedBy: config.User,
					CreatedAt: timestamppb.Now(),
				}},
			},
		}

		ctx := context.Background()
		created, err := config.Client.CreateFeatureAnnotation(ctx, nfa)
		if err != nil {
			return E.Left[string](
				fmt.Errorf(
					"failed to create annotation for %s: %w",
					geneID,
					err,
				),
			)
		}

		config.Logger.Debugf("created new feature annotation %s", created.Id)
		return E.Right[error](fmt.Sprintf("created:%s", geneID))
	},
)

// ProcessingOutcome represents the outcome of processing a single gene
type ProcessingOutcome struct {
	GeneID    string
	Status    string // "created", "updated", "skipped", or "error"
	ErrorInfo O.Option[error]
}

// parseStatusFromMessage extracts the status from a result message
var parseStatusFromMessage = func(msg string) string {
	switch {
	case strings.HasPrefix(msg, "skipped:"):
		return "skipped"
	case strings.HasPrefix(msg, "created:"):
		return "created"
	default:
		return "updated"
	}
}

// createSuccessOutcome creates a successful processing outcome from a message
var createSuccessOutcome = F.Curry2(
	func(geneID string, msg string) ProcessingOutcome {
		return ProcessingOutcome{
			GeneID:    geneID,
			Status:    parseStatusFromMessage(msg),
			ErrorInfo: O.None[error](),
		}
	},
)

// createErrorOutcome creates an error processing outcome
var createErrorOutcome = F.Curry2(
	func(geneID string, err error) ProcessingOutcome {
		return ProcessingOutcome{
			GeneID:    geneID,
			Status:    "error",
			ErrorInfo: O.Some(err),
		}
	},
)

// logGeneError logs an error for a specific gene ID
var logGeneError = F.Curry3(
	func(logger *logrus.Entry, geneID string, err error) error {
		logger.WithFields(logrus.Fields{
			"gene-id": geneID,
			"error":   err,
		}).Error("error processing gene ID")
		return err
	},
)

// processGeneWithResultIO processes a gene and returns the outcome
// Uses F.Curry2 for currying and accepts ProcessingConfig
var processGeneWithResultIO = F.Curry2(
	func(config ProcessingConfig, geneID string) ProcessingOutcome {
		processor := processGeneID(config)

		onError := F.Flow2(
			logGeneError(config.Logger)(geneID),
			createErrorOutcome(geneID),
		)

		return F.Pipe1(
			processor(geneID),
			E.Fold(
				onError,
				createSuccessOutcome(geneID),
			),
		)
	},
)

// countByStatus counts outcomes by status type
var countByStatus = F.Curry2(
	func(status string, outcomes []ProcessingOutcome) int {
		return F.Pipe2(
			outcomes,
			A.Filter(
				func(o ProcessingOutcome) bool { return o.Status == status },
			),
			func(filtered []ProcessingOutcome) int { return len(filtered) },
		)
	},
)

// logProgress logs progress for a batch of outcomes
// Uses F.Curry2 for currying and accepts ProcessingConfig
var logProgress = F.Curry2(
	func(
		config ProcessingConfig,
		total int,
	) func(int, ProcessingOutcome) ProcessingOutcome {
		return func(idx int, outcome ProcessingOutcome) ProcessingOutcome {
			// Log progress every 50 genes
			if (idx+1)%50 == 0 {
				config.Logger.Infof(
					"processed %d/%d gene IDs",
					idx+1,
					total,
				)
			}
			return outcome
		}
	},
)

// processGeneIDsIO processes all gene IDs and returns a processing result
// Returns IOEither to make I/O effects explicit in the type system
// Uses functional patterns: Map, Filter, and reduce operations
func processGeneIDsIO(
	config ProcessingConfig,
) IOE.IOEither[error, HypotheticalGeneProcessingResult] {
	return func() E.Either[error, HypotheticalGeneProcessingResult] {
		config.Logger.Infof("found %d gene IDs to process", len(config.GeneIDs))

		// Process all gene IDs functionally
		outcomes := F.Pipe2(
			config.GeneIDs,
			A.Map(processGeneWithResultIO(config)),
			A.MapWithIndex(logProgress(config)(len(config.GeneIDs))),
		)

		// Count outcomes by status using functional composition
		result := HypotheticalGeneProcessingResult{
			Processed: countByStatus(
				"created",
			)(
				outcomes,
			) + countByStatus(
				"updated",
			)(
				outcomes,
			),
			Skipped: countByStatus("skipped")(outcomes),
			Errors:  countByStatus("error")(outcomes),
		}

		return E.Right[error](result)
	}
}

// handleProcessingError logs and returns the error
var handleProcessingError = F.Curry2(
	func(logger *logrus.Entry, err error) error {
		logger.Errorf("failed to load hypothetical gene products: %v", err)
		return err
	},
)

// handleProcessingSuccess logs the result and returns error if there were failures
var handleProcessingSuccess = F.Curry2(
	func(
		logger *logrus.Entry,
		result HypotheticalGeneProcessingResult,
	) error {
		logger.Infof(
			"finished loading hypothetical gene products. Processed: %d, Skipped: %d, Errors: %d",
			result.Processed,
			result.Skipped,
			result.Errors,
		)
		if result.Errors > 0 {
			return fmt.Errorf(
				"encountered %d errors during loading",
				result.Errors,
			)
		}
		return nil
	},
)

// LoadHypotheticalGeneProducts is the main action handler for the command
func LoadHypotheticalGeneProducts(c *cli.Context) error {
	logger := registry.GetLogger()
	configProcessesor := createProcessingConfig(
		&LoadHypotheticalGeneProductsParams{
			Client: registry.GetFeatureAnnotationAPIClient(),
			User:   c.String("user"),
			Logger: logger,
		},
	)

	logger.Infof(
		"loading hypothetical gene products from %s",
		c.String("input"),
	)

	// Execute the pure functional pipeline - no mutation, consistent IOEither context
	return F.Pipe5(
		c.String("input"),
		readGeneIDsFromFile,
		IOE.Map[error](configProcessesor),
		IOE.Chain(processGeneIDsIO),
		toEither[HypotheticalGeneProcessingResult],
		E.Fold(
			handleProcessingError(logger),
			handleProcessingSuccess(logger),
		),
	)
}

// toEither executes an IOEither to get an Either result
func toEither[A any](ioe IOE.IOEither[error, A]) E.Either[error, A] {
	return ioe()
}
