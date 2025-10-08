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

// GeneProcessingAction represents the action taken for a gene
type GeneProcessingAction int

const (
	// GeneCreated indicates a new annotation was created
	GeneCreated GeneProcessingAction = iota
	// GeneUpdated indicates a product tag was added to existing annotation
	GeneUpdated
	// GeneSkipped indicates the annotation already had the product tag
	GeneSkipped
)

// GeneProcessingResult represents successful processing of a single gene
type GeneProcessingResult struct {
	GeneID string
	Action GeneProcessingAction
}

// ProcessingStats holds aggregate processing statistics
type ProcessingStats struct {
	Created int
	Updated int
	Skipped int
	Total   int
}

// GeneProcessingContext holds initial processing context
type GeneProcessingContext struct {
	Config ProcessingConfig
	GeneID string
}

// WithAnnotation holds context with fetched annotation
type WithAnnotation struct {
	GeneProcessingContext
	Annotation O.Option[*pb.FeatureAnnotation]
}

// WithAction holds context with determined action
type WithAction struct {
	WithAnnotation
	Action GeneProcessingAction
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

// setAnnotation is a setter for Do notation that adds annotation to context
var setAnnotation = F.Curry2(
	func(ann O.Option[*pb.FeatureAnnotation], ctx GeneProcessingContext) WithAnnotation {
		return WithAnnotation{
			GeneProcessingContext: ctx,
			Annotation:            ann,
		}
	},
)

// setAction is a setter for Do notation that adds action to context
var setAction = F.Curry2(
	func(action GeneProcessingAction, ctx WithAnnotation) WithAction {
		return WithAction{
			WithAnnotation: ctx,
			Action:         action,
		}
	},
)

// isNotFoundError checks if error is gRPC NotFound
var isNotFoundError = func(err error) bool {
	return status.Code(err) == codes.NotFound
}

// isHypotheticalProductTag checks if a tag property is the hypothetical protein product
var isHypotheticalProductTag = func(tag *pb.TagProperty) bool {
	return tag.Tag == "product" && tag.Value == HypotheticalProteinProduct
}

// convertGrpcResultToOption converts gRPC (result, error) tuple to (Option, error)
// Handles three cases: success, NotFound, and other errors
var convertGrpcResultToOption = func(
	geneID string,
	result *pb.FeatureAnnotation,
	err error,
) (O.Option[*pb.FeatureAnnotation], error) {
	if err == nil {
		return O.Some(result), nil
	}

	if isNotFoundError(err) {
		return O.None[*pb.FeatureAnnotation](), nil
	}

	return O.None[*pb.FeatureAnnotation](),
		fmt.Errorf("failed to fetch annotation for %s: %w", geneID, err)
}

// fetchAnnotationIO fetches annotation using TryCatchError
// Returns IOEither[error, Option[Annotation]]
var fetchAnnotationIO = F.Curry2(
	func(config ProcessingConfig, geneID string) IOE.IOEither[error, O.Option[*pb.FeatureAnnotation]] {
		return IOE.TryCatchError(
			func() (O.Option[*pb.FeatureAnnotation], error) {
				ctx := context.Background()
				result, err := config.Client.GetFeatureAnnotation(
					ctx, &pb.FeatureAnnotationId{Id: geneID},
				)
				return convertGrpcResultToOption(geneID, result, err)
			},
		)
	},
)

// createAnnotationWithProductIO creates annotation with hypothetical product
var createAnnotationWithProductIO = F.Curry2(
	func(config ProcessingConfig, geneID string) IOE.IOEither[error, GeneProcessingAction] {
		return IOE.TryCatchError(func() (GeneProcessingAction, error) {
			ctx := context.Background()
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

			_, err := config.Client.CreateFeatureAnnotation(ctx, nfa)
			if err != nil {
				return 0,
					fmt.Errorf(
						"failed to create annotation for %s: %w",
						geneID,
						err,
					)
			}

			return GeneCreated, nil
		})
	},
)

// addProductTagIO adds product tag to existing annotation
var addProductTagIO = F.Curry2(
	func(config ProcessingConfig, geneID string) IOE.IOEither[error, GeneProcessingAction] {
		return IOE.TryCatchError(func() (GeneProcessingAction, error) {
			ctx := context.Background()
			_, err := config.Client.AddTag(ctx, &pb.AddTagRequest{
				Id: geneID,
				Tag: &pb.TagPropertyCreate{
					Tag:       "product",
					Value:     HypotheticalProteinProduct,
					CreatedBy: config.User,
					CreatedAt: timestamppb.Now(),
				},
			})
			if err != nil {
				return 0,
					fmt.Errorf("failed to add tag for %s: %w", geneID, err)
			}

			return GeneUpdated, nil
		})
	},
)

// returnSkippedAction returns GeneSkipped action (used in Fold)
var returnSkippedAction = func(*pb.TagProperty) IOE.IOEither[error, GeneProcessingAction] {
	return IOE.Of[error](GeneSkipped)
}

// handleProductNotFound handles product not found case
var handleProductNotFound = F.Curry2(
	func(config ProcessingConfig, geneID string) func() IOE.IOEither[error, GeneProcessingAction] {
		return func() IOE.IOEither[error, GeneProcessingAction] {
			return addProductTagIO(config)(geneID)
		}
	},
)

// getActionForExistingAnnotation checks if annotation has hypothetical product and returns action
var getActionForExistingAnnotation = F.Curry2(
	func(config ProcessingConfig, geneID string) func(*pb.FeatureAnnotation) IOE.IOEither[error, GeneProcessingAction] {
		return func(ann *pb.FeatureAnnotation) IOE.IOEither[error, GeneProcessingAction] {
			return F.Pipe2(
				ann.Attributes.Properties,
				A.FindFirst(isHypotheticalProductTag),
				O.Fold(
					handleProductNotFound(config)(geneID),
					returnSkippedAction,
				),
			)
		}
	},
)

// handleAnnotationNotFound handles annotation not found case
var handleAnnotationNotFound = F.Curry2(
	func(config ProcessingConfig, geneID string) func() IOE.IOEither[error, GeneProcessingAction] {
		return func() IOE.IOEither[error, GeneProcessingAction] {
			return createAnnotationWithProductIO(config)(geneID)
		}
	},
)

// ensureHypotheticalProductIO ensures hypothetical product exists
var ensureHypotheticalProductIO = F.Curry2(
	func(config ProcessingConfig, ctx WithAnnotation) IOE.IOEither[error, GeneProcessingAction] {
		return F.Pipe1(
			ctx.Annotation,
			O.Fold(
				handleAnnotationNotFound(config)(ctx.GeneID),
				getActionForExistingAnnotation(config)(ctx.GeneID),
			),
		)
	},
)

// fetchAnnotationForContext fetches annotation for context
var fetchAnnotationForContext = func(ctx GeneProcessingContext) IOE.IOEither[error, O.Option[*pb.FeatureAnnotation]] {
	return fetchAnnotationIO(ctx.Config)(ctx.GeneID)
}

// extractGeneResult extracts final result from context
var extractGeneResult = func(ctx WithAction) GeneProcessingResult {
	return GeneProcessingResult{
		GeneID: ctx.GeneID,
		Action: ctx.Action,
	}
}

// processAllGenes processes all genes using TraverseArraySeq
func processAllGenes(
	cfg ProcessingConfig,
) IOE.IOEither[error, []GeneProcessingResult] {
	return F.Pipe1(
		cfg.GeneIDs,
		IOE.TraverseArraySeq(
			func(geneID string) IOE.IOEither[error, GeneProcessingResult] {
				return F.Pipe3(
					IOE.Of[error](GeneProcessingContext{
						Config: cfg,
						GeneID: geneID,
					}),
					IOE.Bind(setAnnotation, fetchAnnotationForContext),
					IOE.Bind(setAction, ensureHypotheticalProductIO(cfg)),
					IOE.Map[error](extractGeneResult),
				)
			},
		),
	)
}

// aggregateResults aggregates processing results into stats
var aggregateResults = func(results []GeneProcessingResult) ProcessingStats {
	reducer := func(stats ProcessingStats, result GeneProcessingResult) ProcessingStats {
		switch result.Action {
		case GeneCreated:
			return ProcessingStats{
				Total:   stats.Total,
				Created: stats.Created + 1,
				Updated: stats.Updated,
				Skipped: stats.Skipped,
			}
		case GeneUpdated:
			return ProcessingStats{
				Total:   stats.Total,
				Created: stats.Created,
				Updated: stats.Updated + 1,
				Skipped: stats.Skipped,
			}
		case GeneSkipped:
			return ProcessingStats{
				Total:   stats.Total,
				Created: stats.Created,
				Updated: stats.Updated,
				Skipped: stats.Skipped + 1,
			}
		default:
			// This switch is exhaustive for all GeneProcessingAction values
			// If we reach here, it indicates a programming error
			return stats
		}
	}

	return F.Pipe1(
		results,
		A.Reduce(reducer, ProcessingStats{Total: len(results)}),
	)
}

// handleProcessingError logs error and returns it
var handleProcessingError = F.Curry2(
	func(logger *logrus.Entry, err error) error {
		logger.Errorf("failed to load hypothetical gene products: %v", err)
		return err
	},
)

// logAndCheckStats logs final statistics
var logAndCheckStats = F.Curry2(
	func(logger *logrus.Entry, stats ProcessingStats) error {
		logger.Infof(
			"finished loading hypothetical gene products. Created: %d, Updated: %d, Skipped: %d, Total: %d",
			stats.Created,
			stats.Updated,
			stats.Skipped,
			stats.Total,
		)
		return nil
	},
)

// LoadHypotheticalGeneProducts is the main action handler for the command
func LoadHypotheticalGeneProducts(c *cli.Context) error {
	logger := registry.GetLogger()
	params := &LoadHypotheticalGeneProductsParams{
		Client: registry.GetFeatureAnnotationAPIClient(),
		User:   c.String("user"),
		Logger: logger,
	}

	logger.Infof(
		"loading hypothetical gene products from %s",
		c.String("input"),
	)

	// Execute the pure functional pipeline - no mutation, consistent IOEither context
	return F.Pipe2(
		F.Pipe5(
			c.String("input"),
			readGeneIDsFromFile,
			IOE.Map[error](createProcessingConfig(params)),
			IOE.Map[error](func(cfg ProcessingConfig) ProcessingConfig {
				logger.Infof("found %d gene IDs to process", len(cfg.GeneIDs))
				return cfg
			}),
			IOE.Chain(processAllGenes),
			IOE.Map[error](aggregateResults),
		),
		toEither[ProcessingStats],
		E.Fold(
			handleProcessingError(logger),
			logAndCheckStats(logger),
		),
	)
}

// toEither executes an IOEither to get an Either result
func toEither[A any](ioe IOE.IOEither[error, A]) E.Either[error, A] {
	return ioe()
}
