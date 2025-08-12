package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	fanno "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gene struct {
	FeatureID int    `json:"feature_id"`
	GeneID    string `json:"gene_id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

type GeneWithPubmed struct {
	Gene
	Pubmeds    []string
	Skip       bool
	SkipReason string
}

type GrpcAnnotationResult struct {
	GeneID  string
	Success bool
	Skipped bool
	Message string
	Error   error
}

type FeatureAnnotationMetrics struct {
	TotalProcessed              int64
	SuccessCount                int64
	ErrorCount                  int64
	StartTime                   time.Time
	mu                          sync.RWMutex
	TotalFetchedFromArango      int64
	AllArangoDocsFetched        bool
	JobsSubmittedToPubmedPool   int64
	JobsCompletedFromPubmedPool int64
	JobsSubmittedToGrpcPool     int64
	JobsCompletedFromGrpcPool   int64
}

func (m *FeatureAnnotationMetrics) IsComplete() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	allArangoFetched := m.AllArangoDocsFetched
	pubmedPoolDrained := m.JobsCompletedFromPubmedPool >= m.JobsSubmittedToPubmedPool
	grpcPoolDrained := m.JobsCompletedFromGrpcPool >= m.JobsSubmittedToGrpcPool
	allProcessed := (m.TotalFetchedFromArango == 0) ||
		(m.TotalProcessed >= m.TotalFetchedFromArango)

	return allArangoFetched && pubmedPoolDrained && grpcPoolDrained &&
		allProcessed
}

type FeatureAnnotationAppConfig struct {
	Ctx              context.Context
	NumPubmedWorkers int
	NumGrpcWorkers   int
	Logger           *logrus.Entry
	Metrics          *FeatureAnnotationMetrics
}

func newFeatureAnnotationConfigFromCliContext(
	cltx *cli.Context,
	logger *logrus.Entry,
) FeatureAnnotationAppConfig {
	return FeatureAnnotationAppConfig{
		NumPubmedWorkers: cltx.Int("pubmed-workers"),
		NumGrpcWorkers:   cltx.Int("grpc-workers"),
		Logger:           logger,
		Metrics: &FeatureAnnotationMetrics{
			StartTime: time.Now(),
		},
	}
}

func RunFeatureAnnotationLoader(cltx *cli.Context) error {
	logger := registry.GetLogger()
	config := newFeatureAnnotationConfigFromCliContext(cltx, logger)

	logger.Debug("Starting feature annotation loader...")
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()
	config.Ctx = mainCtx

	setupSignalHandling(mainCancel, logger)

	genesFromQueryChan := make(chan Gene, config.NumPubmedWorkers)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go queryActiveGenesForAnnotation(&wg, &queryActiveGenesForAnnotationParams{
		config:     config,
		genesChan:  genesFromQueryChan,
		mainCancel: mainCancel,
	})

	client := registry.GetFeatureAnnotationAPIClient()
	pubmedFetchPool := setupPubmedFetchPool(config, mainCtx, client)
	annotationCreatePool := setupAnnotationCreatePool(config, mainCtx, client)

	wg.Add(1)
	go func() {
		defer wg.Done()
		bridgeGenesToPubmedPool(&bridgeGenesToPubmedPoolParams{
			ctx:        mainCtx,
			genesChan:  genesFromQueryChan,
			pubmedPool: pubmedFetchPool,
			metrics:    config.Metrics,
			logger:     logger,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bridgePubmedToGrpcPool(&bridgePubmedToGrpcPoolParams{
			ctx:        mainCtx,
			pubmedPool: pubmedFetchPool,
			grpcPool:   annotationCreatePool,
			metrics:    config.Metrics,
			logger:     logger,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		handleAnnotationGrpcResults(&handleAnnotationGrpcResultsParams{
			ctx:      mainCtx,
			grpcPool: annotationCreatePool,
			metrics:  config.Metrics,
			logger:   logger,
		})
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reportAnnotationProgress(&reportAnnotationProgressParams{
			ctx:        mainCtx,
			metrics:    config.Metrics,
			logger:     logger,
			mainCancel: mainCancel,
		})
	}()

	logger.Debug("Waiting for all goroutines to complete...")
	wg.Wait()
	logger.Info("Feature annotation loader finished")
	return nil
}

type queryActiveGenesForAnnotationParams struct {
	config     FeatureAnnotationAppConfig
	genesChan  chan<- Gene
	mainCancel context.CancelFunc
}

func queryActiveGenesForAnnotation(
	wg *sync.WaitGroup,
	params *queryActiveGenesForAnnotationParams,
) {
	defer wg.Done()
	defer close(params.genesChan)

	logger := params.config.Logger
	dbh := registry.GetArangodbConnection()
	logger.Info("Starting active genes query...")

	cursor, err := dbh.Search(ListActiveGenesQ)
	if err != nil {
		logger.Errorf("Error querying active genes: %v", err)
		params.mainCancel()
		return
	}
	defer cursor.Close()

	if cursor.IsEmpty() {
		logger.Warn("No active genes found.")
		params.config.Metrics.mu.Lock()
		params.config.Metrics.TotalFetchedFromArango = 0
		params.config.Metrics.AllArangoDocsFetched = true
		params.config.Metrics.mu.Unlock()
		return
	}

	var geneCount int64
	for cursor.Scan() {
		select {
		case <-params.config.Ctx.Done():
			logger.Info(
				`Active genes querier stopping due to context cancellation.`,
			)
			return
		default:
			gene := Gene{}
			if err := cursor.Read(&gene); err != nil {
				logger.Errorf("Failed to read gene: %v", err)
				params.mainCancel()
				return
			}

			select {
			case params.genesChan <- gene:
				geneCount++
			case <-params.config.Ctx.Done():
				logger.Info("Active genes querier stopping during send.")
				return
			}
		}
	}

	logger.Infof("Successfully fetched %d active genes.", geneCount)

	params.config.Metrics.mu.Lock()
	params.config.Metrics.TotalFetchedFromArango = geneCount
	params.config.Metrics.AllArangoDocsFetched = true
	params.config.Metrics.mu.Unlock()
}

func setupPubmedFetchPool(
	config FeatureAnnotationAppConfig,
	mainCtx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
) *concurrent.Pool[Gene, GeneWithPubmed] {
	pool := concurrent.NewPool(
		pubmedFetcherWorkerFunc(config, grpcClient),
		concurrent.WithWorkers[Gene, GeneWithPubmed](config.NumPubmedWorkers),
		concurrent.WithContext[Gene, GeneWithPubmed](mainCtx),
		concurrent.WithBufferSize[Gene, GeneWithPubmed](
			config.NumPubmedWorkers*2,
		),
	)
	pool.Start()
	return pool
}

func setupAnnotationCreatePool(
	config FeatureAnnotationAppConfig,
	mainCtx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
) *concurrent.Pool[GeneWithPubmed, GrpcAnnotationResult] {
	pool := concurrent.NewPool(
		annotationCreatorWorkerFunc(grpcClient),
		concurrent.WithWorkers[GeneWithPubmed, GrpcAnnotationResult](
			config.NumGrpcWorkers,
		),
		concurrent.WithContext[GeneWithPubmed, GrpcAnnotationResult](mainCtx),
		concurrent.WithBufferSize[GeneWithPubmed, GrpcAnnotationResult](
			config.NumGrpcWorkers*2,
		),
	)
	pool.Start()
	return pool
}

type bridgeGenesToPubmedPoolParams struct {
	ctx        context.Context
	genesChan  <-chan Gene
	pubmedPool *concurrent.Pool[Gene, GeneWithPubmed]
	metrics    *FeatureAnnotationMetrics
	logger     *logrus.Entry
}

func bridgeGenesToPubmedPool(params *bridgeGenesToPubmedPoolParams) {
	defer params.pubmedPool.Close()

	params.logger.Info("Starting Genes-to-Pubmed-Pool bridge...")
	for {
		select {
		case <-params.ctx.Done():
			return
		case gene, ok := <-params.genesChan:
			if !ok {
				return
			}
			params.pubmedPool.Submit(gene)
			params.metrics.mu.Lock()
			params.metrics.JobsSubmittedToPubmedPool++
			params.metrics.mu.Unlock()
			params.logger.WithFields(logrus.Fields{
				"gene_id": gene.GeneID,
				"stage":   "submitted_to_pubmed_pool",
			}).Debug("Gene submitted for pubmed fetching")
		}
	}
}

type bridgePubmedToGrpcPoolParams struct {
	ctx        context.Context
	pubmedPool *concurrent.Pool[Gene, GeneWithPubmed]
	grpcPool   *concurrent.Pool[GeneWithPubmed, GrpcAnnotationResult]
	metrics    *FeatureAnnotationMetrics
	logger     *logrus.Entry
}

func bridgePubmedToGrpcPool(params *bridgePubmedToGrpcPoolParams) {
	defer params.grpcPool.Close()
	params.logger.Debug("Starting Pubmed-Pool-to-gRPC-Pool bridge...")
	for {
		select {
		case <-params.ctx.Done():
			params.logger.Debug("Pubmed-to-gRPC bridge: context done, exiting.")
			return
		case result, ok := <-params.pubmedPool.Results():
			if !ok {
				params.logger.Debug("Pubmed pool results channel closed.")
				return
			}
			params.metrics.mu.Lock()
			params.metrics.JobsCompletedFromPubmedPool++
			params.metrics.mu.Unlock()

			if result.Error != nil {
				params.logger.WithFields(logrus.Fields{
					"job_id": result.JobID,
					"error":  result.Error,
				}).Error("Pubmed fetching failed")
				continue
			}
			params.grpcPool.Submit(result.Output)
			params.metrics.mu.Lock()
			params.metrics.JobsSubmittedToGrpcPool++
			params.metrics.mu.Unlock()
			params.logger.WithFields(logrus.Fields{
				"gene_id": result.Output.GeneID,
				"stage":   "submitted_to_grpc_pool",
			}).Debug("Gene with pubmeds submitted for gRPC creation")
		case err, ok := <-params.pubmedPool.Errors():
			if !ok {
				params.logger.Debug("Pubmed pool errors channel closed.")
				return
			}
			params.logger.Errorf("Async error from pubmed pool: %v", err)
		}
	}
}

type handleAnnotationGrpcResultsParams struct {
	ctx      context.Context
	grpcPool *concurrent.Pool[GeneWithPubmed, GrpcAnnotationResult]
	metrics  *FeatureAnnotationMetrics
	logger   *logrus.Entry
}

func handleAnnotationGrpcResults(params *handleAnnotationGrpcResultsParams) {
	params.logger.Debug("Starting gRPC results handler...")

	for {
		select {
		case <-params.ctx.Done():
			return
		case result, ok := <-params.grpcPool.Results():
			if !ok {
				params.logger.Debug("gRPC create pool results channel closed.")
				return
			}

			params.metrics.mu.Lock()
			params.metrics.TotalProcessed++
			params.metrics.JobsCompletedFromGrpcPool++
			if result.Error != nil {
				params.metrics.ErrorCount++
				params.metrics.mu.Unlock()
				params.logger.Errorf(
					"Error creating annotation for gene %s: %v",
					result.Output.GeneID,
					result.Error,
				)
			} else {
				params.metrics.SuccessCount++
				params.metrics.mu.Unlock()
				if result.Output.Skipped {
					params.logger.Infof(
						"Annotation processing for gene %s: %s",
						result.Output.GeneID,
						result.Output.Message,
					)
				} else {
					params.logger.Infof(
						"Successfully created annotation for gene %s: %s",
						result.Output.GeneID,
						result.Output.Message,
					)
				}
			}
		case err, ok := <-params.grpcPool.Errors():
			if !ok {
				params.logger.Debug("gRPC create pool errors channel closed.")
				return
			}
			params.metrics.mu.Lock()
			params.metrics.ErrorCount++
			params.metrics.mu.Unlock()
			params.logger.Errorf("Async error from gRPC create pool: %v", err)
		}
	}
}

type reportAnnotationProgressParams struct {
	ctx        context.Context
	metrics    *FeatureAnnotationMetrics
	logger     *logrus.Entry
	mainCancel context.CancelFunc
}

func reportAnnotationProgress(params *reportAnnotationProgressParams) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logCurrentMetrics := func(message string) {
		params.metrics.mu.RLock()
		defer params.metrics.mu.RUnlock()
		elapsed := time.Since(params.metrics.StartTime)
		rate := 0.0
		if elapsed.Seconds() > 0 {
			rate = float64(params.metrics.TotalProcessed) / elapsed.Seconds()
		}
		params.logger.WithFields(logrus.Fields{
			"read_from_db":     params.metrics.TotalFetchedFromArango,
			"total_processed":  params.metrics.TotalProcessed,
			"success_count":    params.metrics.SuccessCount,
			"error_count":      params.metrics.ErrorCount,
			"processing_rate":  fmt.Sprintf("%.2f genes/sec", rate),
			"elapsed_time":     elapsed.String(),
			"pubmed_submitted": params.metrics.JobsSubmittedToPubmedPool,
			"pubmed_completed": params.metrics.JobsCompletedFromPubmedPool,
			"grpc_submitted":   params.metrics.JobsSubmittedToGrpcPool,
			"grpc_completed":   params.metrics.JobsCompletedFromGrpcPool,
		}).Info(message)
	}

	for {
		select {
		case <-params.ctx.Done():
			logCurrentMetrics("Final annotation processing report")
			return
		case <-ticker.C:
			logCurrentMetrics("Annotation processing progress")
			if params.metrics.IsComplete() {
				params.logger.Info(
					"All annotation updates appear completed. Stopping progress reporter.",
				)
				params.mainCancel()
				return
			}
		}
	}
}

func pubmedFetcherWorkerFunc(
	config FeatureAnnotationAppConfig,
	grpcClient fanno.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[Gene, GeneWithPubmed] {
	return func(
		ctx context.Context,
		job concurrent.Job[Gene],
	) (GeneWithPubmed, error) {
		gene := job.Payload
		logger := config.Logger

		exists, err := checkAnnotationExists(ctx, grpcClient, gene.GeneID)
		if err != nil {
			return GeneWithPubmed{}, fmt.Errorf(
				"error checking for feature annotation for gene %s: %w",
				gene.GeneID,
				err,
			)
		}
		if exists {
			return GeneWithPubmed{
				Gene:       gene,
				Skip:       true,
				SkipReason: "annotation already exists",
			}, nil
		}

		dbh := registry.GetArangodbConnection()
		pubmedResult, err := dbh.SearchRows(
			ListPubmedsByFeature,
			map[string]any{
				"feature_id": gene.FeatureID,
			},
		)
		if err != nil {
			return GeneWithPubmed{}, fmt.Errorf(
				"error querying PubMed IDs for feature %d: %w",
				gene.FeatureID,
				err,
			)
		}
		defer pubmedResult.Close()

		var pubmedIDs []string
		for pubmedResult.Scan() {
			var pubmed string
			if err := pubmedResult.Read(&pubmed); err != nil {
				return GeneWithPubmed{}, fmt.Errorf(
					"error reading PubMed ID for feature %d: %w",
					gene.FeatureID,
					err,
				)
			}
			pubmedIDs = append(pubmedIDs, pubmed)
		}

		if len(pubmedIDs) == 0 {
			return GeneWithPubmed{
				Gene:       gene,
				Skip:       true,
				SkipReason: "no pubmed entries found",
			}, nil
		}

		logger.Debugf("Feature %d has PubMed reference: %v",
			gene.FeatureID,
			pubmedIDs,
		)
		return GeneWithPubmed{Gene: gene, Pubmeds: pubmedIDs}, nil
	}
}

func annotationCreatorWorkerFunc(
	grpcClient fanno.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[GeneWithPubmed, GrpcAnnotationResult] {
	return func(
		ctx context.Context,
		job concurrent.Job[GeneWithPubmed],
	) (GrpcAnnotationResult, error) {
		geneWithPubmed := job.Payload
		result := GrpcAnnotationResult{
			GeneID:  geneWithPubmed.GeneID,
			Success: false,
		}

		if geneWithPubmed.Skip {
			result.Success = true
			result.Skipped = true
			result.Message = fmt.Sprintf(
				"Skipped, %s",
				geneWithPubmed.SkipReason,
			)
			return result, nil
		}

		res, err := createNewAnnotation(ctx, grpcClient, geneWithPubmed)
		if err != nil {
			result.Error = err
			result.Message = fmt.Sprintf(
				"failed to create feature annotation for gene %s: %v",
				geneWithPubmed.GeneID,
				err,
			)
			return result, err
		}

		result.Success = true
		result.Message = fmt.Sprintf(
			"Created new feature annotation record %s for feature name %s",
			res.Attributes.Name,
			res.Id,
		)
		return result, nil
	}
}

// checkAnnotationExists checks if a feature annotation already exists.
// It returns true if it exists, false if it does not exist (NotFound),
// and an error for any other grpc status.
func checkAnnotationExists(
	ctx context.Context,
	client fanno.FeatureAnnotationServiceClient,
	geneID string,
) (bool, error) {
	_, err := client.GetFeatureAnnotation(
		ctx,
		&fanno.FeatureAnnotationId{Id: geneID},
	)
	if err == nil {
		return true, nil // Exists
	}
	if status.Code(err) == codes.NotFound {
		return false, nil // Does not exist
	}
	return false, err // Other error
}

func createNewAnnotation(
	ctx context.Context,
	grpcClient fanno.FeatureAnnotationServiceClient,
	geneWithPubmed GeneWithPubmed,
) (*fanno.FeatureAnnotation, error) {
	return grpcClient.CreateFeatureAnnotation(
		ctx,
		&fanno.NewFeatureAnnotation{
			Type:       "gene",
			Id:         geneWithPubmed.GeneID,
			IsObsolete: false,
			CreatedAt:  timestamppb.Now(),
			UpdatedAt:  timestamppb.Now(),
			CreatedBy: resolveCreatorFromCreatedBy(
				geneWithPubmed.CreatedBy,
			),
			Attributes: &fanno.FeatureAnnotationAttributes{
				Name:   geneWithPubmed.Name,
				Pubmed: geneWithPubmed.Pubmeds,
			},
		},
	)
}
