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
	// GeneCreated indicates a new annotation was created
	GeneCreated GeneProcessingAction = iota
	// GeneUpdated indicates a product tag was added to existing annotation
	GeneUpdated
	// GeneSkipped indicates the annotation already had the product tag
	GeneSkipped
)

// Point-free helper predicates
var (
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

	// createProcessingConfig creates an immutable ProcessingConfig from params and
	// gene IDs.
	createProcessingConfig = F.Curry2(
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

	// setAnnotation is a setter for Do notation that adds annotation to
	// context
	setAnnotation = F.Curry2(
		func(
			ann O.Option[*pb.FeatureAnnotation],
			ctx GeneProcessingContext,
		) WithAnnotation {
			return WithAnnotation{
				GeneProcessingContext: ctx,
				Annotation:            ann,
			}
		},
	)

	// setAction is a setter for Do notation that adds action to context
	setAction = F.Curry2(
		func(
			action GeneProcessingAction,
			ctx WithAnnotation,
		) WithAction {
			return WithAction{
				WithAnnotation: ctx,
				Action:         action,
			}
		},
	)
	checkAndHandleProductTag = F.Curry2(
		func(
			ctx WithAnnotation,
			ann *pb.FeatureAnnotation,
		) IOE.IOEither[error, GeneProcessingAction] {
			return F.Pipe2(
				ann.Attributes.Properties,
				A.FindFirst(isHypotheticalProductTag),
				O.Fold(
					addProductTag(ctx.GeneProcessingContext),
					returnSkippedAction,
				),
			)
		})
)

// GeneProcessingAction represents the action taken for a gene
type GeneProcessingAction int

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

// createGeneAnnotation creates a new gene annotation with hypothetical product
func createGeneAnnotation(
	config ProcessingConfig,
	geneID string,
) IOE.IOEither[error, GeneProcessingAction] {
	return IOE.TryCatchError(func() (GeneProcessingAction, error) {
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

		_, err := config.Client.CreateFeatureAnnotation(
			context.Background(),
			nfa,
		)
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
}

// processAnnotation decides whether to create, update or skip a gene annotation
func processAnnotation(
	ctx WithAnnotation,
) IOE.IOEither[error, GeneProcessingAction] {
	return F.Pipe1(
		ctx.Annotation,
		O.Fold(
			handleAnnotationNotFound(ctx),
			checkAndHandleProductTag(ctx),
		),
	)
}

// handleAnnotationNotFound returns a function that creates a gene annotation.
func handleAnnotationNotFound(
	ctx WithAnnotation,
) func() IOE.IOEither[error, GeneProcessingAction] {
	return func() IOE.IOEither[error, GeneProcessingAction] {
		return createGeneAnnotation(ctx.Config, ctx.GeneID)
	}
}

// addProductTag returns a function that adds a product tag to a gene.
func addProductTag(
	gctx GeneProcessingContext,
) func() IOE.IOEither[error, GeneProcessingAction] {
	return func() IOE.IOEither[error, GeneProcessingAction] {
		return addProductTagIO(gctx)
	}
}

// fetchAnnotationForContext fetches annotation for context
func fetchAnnotationForContext(
	ctx GeneProcessingContext,
) IOE.IOEither[error, O.Option[*pb.FeatureAnnotation]] {
	return IOE.TryCatchError(
		func() (O.Option[*pb.FeatureAnnotation], error) {
			result, err := ctx.Config.Client.
				GetFeatureAnnotation(
					context.Background(),
					&pb.FeatureAnnotationId{
						Id: ctx.GeneID,
					},
				)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					return O.None[*pb.FeatureAnnotation](), nil
				}
				return O.None[*pb.FeatureAnnotation](), err
			}
			return O.Some(result), nil
		},
	)
}

// extractGeneResult extracts final result from context
func extractGeneResult(ctx WithAction) GeneProcessingResult {
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
					IOE.Bind(setAction, processAnnotation),
					IOE.Map[error](extractGeneResult),
				)
			},
		),
	)
}

// aggregateResults aggregates processing results into stats
func aggregateResults(results []GeneProcessingResult) ProcessingStats {
	reducer := func(
		stats ProcessingStats,
		result GeneProcessingResult,
	) ProcessingStats {
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

// stripBOM removes Unicode byte order mark from beginning of string
func stripBOM(s string) string {
	return strings.TrimPrefix(s, "\ufeff")
}

// isHypotheticalProductTag checks if a tag property is the hypothetical protein
// product
func isHypotheticalProductTag(tag *pb.TagProperty) bool {
	return tag.Tag == "product" && tag.Value == HypotheticalProteinProduct
}

// addProductTagIO adds product tag to existing annotation
func addProductTagIO(
	ctx GeneProcessingContext,
) IOE.IOEither[error, GeneProcessingAction] {
	return IOE.TryCatchError(func() (GeneProcessingAction, error) {
		gctx := context.Background()
		_, err := ctx.Config.Client.AddTag(gctx, &pb.AddTagRequest{
			Id: ctx.GeneID,
			Tag: &pb.TagPropertyCreate{
				Tag:       "product",
				Value:     HypotheticalProteinProduct,
				CreatedBy: ctx.Config.User,
				CreatedAt: timestamppb.Now(),
			},
		})
		if err != nil {
			return 0,
				fmt.Errorf(
					"failed to add tag for %s: %w",
					ctx.GeneID,
					err,
				)
		}

		return GeneUpdated, nil
	})
}

// returnSkippedAction returns GeneSkipped action (used in Fold)
func returnSkippedAction(
	*pb.TagProperty,
) IOE.IOEither[error, GeneProcessingAction] {
	return IOE.Of[error](GeneSkipped)
}

// toEither executes an IOEither to get an Either result
func toEither[A any](ioe IOE.IOEither[error, A]) E.Either[error, A] {
	return ioe()
}
