package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	fanno "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SynonymData holds gene synonym information from ArangoDB
type SynonymData struct {
	Name     string   `json:"name"`
	GeneID   string   `json:"gene_id"`
	Synonyms []string `json:"synonyms"`
}

// GrpcSynonymResult holds the result of a gRPC update operation for synonyms.
type GrpcSynonymResult struct {
	GeneID  string
	Success bool
	Message string
	Error   error
}

// SynonymMetrics holds processing metrics for the synonym loader.
type SynonymMetrics struct {
	TotalProcessed            int64
	SuccessCount              int64
	NotFoundCount             int64
	ErrorCount                int64
	StartTime                 time.Time
	mu                        sync.RWMutex
	TotalFetchedFromArango    int64
	AllArangoDocsFetched      bool
	JobsSubmittedToGrpcPool   int64
	JobsCompletedFromGrpcPool int64
}

// IsComplete checks if all processing is finished.
func (m *SynonymMetrics) IsComplete() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.isPrimaryComplete() || m.isFallbackComplete()
}

func (m *SynonymMetrics) isPrimaryComplete() bool {
	allArangoFetched := m.AllArangoDocsFetched
	grpcPoolDrained := m.JobsCompletedFromGrpcPool >= m.JobsSubmittedToGrpcPool
	allProcessed := (m.TotalFetchedFromArango == 0) ||
		(m.TotalProcessed >= (m.TotalFetchedFromArango))

	return allArangoFetched && grpcPoolDrained && allProcessed
}

func (m *SynonymMetrics) isFallbackComplete() bool {
	elapsed := time.Since(m.StartTime)
	hasSignificantWork := m.TotalFetchedFromArango > 0
	hasProcessedSomething := (m.SuccessCount + m.ErrorCount + m.NotFoundCount) > 0
	poolsDrained := m.JobsCompletedFromGrpcPool >= m.JobsSubmittedToGrpcPool

	return m.AllArangoDocsFetched &&
		poolsDrained &&
		hasSignificantWork &&
		hasProcessedSomething &&
		elapsed > 2*time.Minute
}

// SynonymAppConfig holds configuration for the synonym loader application.
type SynonymAppConfig struct {
	Ctx            context.Context
	NumGrpcWorkers int
	Logger         *logrus.Entry
	Metrics        *SynonymMetrics
}

type bridgeSynonymsToGrpcPoolParams struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	synsChan <-chan SynonymData
	grpcPool *concurrent.Pool[SynonymData, GrpcSynonymResult]
	metrics  *SynonymMetrics
	logger   *logrus.Entry
}

type handleSynonymGrpcResultsParams struct {
	wg       *sync.WaitGroup
	ctx      context.Context
	grpcPool *concurrent.Pool[SynonymData, GrpcSynonymResult]
	metrics  *SynonymMetrics
	logger   *logrus.Entry
}

type reportSynonymProgressParams struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	metrics *SynonymMetrics
	logger  *logrus.Entry
}

// querySynonymsParams holds parameters for querying synonyms.
type querySynonymsParams struct {
	config     SynonymAppConfig
	synsChan   chan<- SynonymData
	mainCancel context.CancelFunc
}

// RunSynonymLoader is the main entry point for the synonym loader.
func RunSynonymLoader(cltx *cli.Context) error {
	logger := registry.GetLogger()
	config := newSynonymConfigFromCliContext(cltx, logger)

	logger.Debug("Starting synonym loader...")
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()
	config.Ctx = mainCtx

	setupSignalHandling(mainCancel, logger)

	// Channel for synonyms from ArangoDB
	synonymsFromQueryChan := make(
		chan SynonymData,
		config.NumGrpcWorkers,
	)

	wg := sync.WaitGroup{}
	// Goroutine 1: ArangoDB Synonym Querier
	wg.Add(1)
	go querySynonyms(&wg, &querySynonymsParams{
		config:     config,
		synsChan:   synonymsFromQueryChan,
		mainCancel: mainCancel,
	})

	// Setup Pool
	grpcUpdatePool := setupSynonymGrpcUpdatePool(config, mainCtx)

	// Bridge from ArangoDB to gRPC Pool
	wg.Add(1)
	go bridgeSynonymsToGrpcPool(&bridgeSynonymsToGrpcPoolParams{
		wg:       &wg,
		ctx:      mainCtx,
		synsChan: synonymsFromQueryChan,
		grpcPool: grpcUpdatePool,
		metrics:  config.Metrics,
		logger:   logger,
	})

	// Handle gRPC Results
	wg.Add(1)
	go handleSynonymGrpcResults(&handleSynonymGrpcResultsParams{
		wg:       &wg,
		ctx:      mainCtx,
		grpcPool: grpcUpdatePool,
		metrics:  config.Metrics,
		logger:   logger,
	})

	// Progress Reporter
	wg.Add(1)
	go reportSynonymProgress(&reportSynonymProgressParams{
		wg:      &wg,
		ctx:     mainCtx,
		metrics: config.Metrics,
		logger:  logger,
	})

	logger.Debug("Waiting for all goroutines to complete...")
	wg.Wait()
	logger.Info("Synonym loader finished")
	return nil
}

// newSynonymConfigFromCliContext creates config from CLI context.
func newSynonymConfigFromCliContext(
	cltx *cli.Context,
	logger *logrus.Entry,
) SynonymAppConfig {
	return SynonymAppConfig{
		NumGrpcWorkers: cltx.Int("grpc-workers"),
		Logger:         logger,
		Metrics: &SynonymMetrics{
			StartTime: time.Now(),
		},
	}
}

// setupSynonymGrpcUpdatePool sets up gRPC update pool for synonyms.
func setupSynonymGrpcUpdatePool(
	config SynonymAppConfig,
	mainCtx context.Context,
) *concurrent.Pool[SynonymData, GrpcSynonymResult] {
	pool := concurrent.NewPool(
		grpcSynonymWorkerFunc(
			config,
			registry.GetFeatureAnnotationAPIClient(),
		),
		concurrent.WithWorkers[SynonymData, GrpcSynonymResult](
			config.NumGrpcWorkers,
		),
		concurrent.WithContext[SynonymData, GrpcSynonymResult](mainCtx),
		concurrent.WithBufferSize[SynonymData, GrpcSynonymResult](
			config.NumGrpcWorkers*2,
		),
	)
	pool.Start()
	return pool
}

// querySynonyms queries ArangoDB for synonyms.
func querySynonyms(
	wg *sync.WaitGroup,
	params *querySynonymsParams,
) {
	defer wg.Done()
	defer close(params.synsChan)

	logger := params.config.Logger
	dbh := registry.GetArangodbConnection()
	logger.Info("Starting synonym query...")

	cursor, err := dbh.Search(ListSynonyms)
	if err != nil {
		logger.Errorf("Error querying synonyms: %v", err)
		params.mainCancel()
		return
	}
	defer cursor.Close()

	if cursor.IsEmpty() {
		logger.Warn("No synonyms found.")
		params.config.Metrics.mu.Lock()
		params.config.Metrics.TotalFetchedFromArango = 0
		params.config.Metrics.AllArangoDocsFetched = true
		params.config.Metrics.mu.Unlock()
		return
	}

	var synCount int64
	for cursor.Scan() {
		select {
		case <-params.config.Ctx.Done():
			logger.Info(
				`Synonym querier stopping due to context cancellation.`,
			)
			return
		default:
			synData := SynonymData{}
			if err := cursor.Read(&synData); err != nil {
				logger.Errorf("Failed to read synonym data: %v", err)
				params.mainCancel()
				return
			}

			select {
			case params.synsChan <- synData:
				synCount++
			case <-params.config.Ctx.Done():
				logger.Info("Synonym querier stopping during send.")
				return
			}
		}
	}

	logger.Infof("Successfully fetched %d synonym records.", synCount)

	params.config.Metrics.mu.Lock()
	params.config.Metrics.TotalFetchedFromArango = synCount
	params.config.Metrics.AllArangoDocsFetched = true
	params.config.Metrics.mu.Unlock()
}

// bridgeSynonymsToGrpcPool transfers synonyms to gRPC update pool.
func bridgeSynonymsToGrpcPool(params *bridgeSynonymsToGrpcPoolParams) {
	defer params.wg.Done()
	defer params.grpcPool.Close()

	params.logger.Info("Starting Synonyms-to-gRPC-Pool bridge...")
	for {
		select {
		case <-params.ctx.Done():
			return
		case syn, ok := <-params.synsChan:
			if !ok {
				return
			}
			params.grpcPool.Submit(syn)
			params.metrics.mu.Lock()
			params.metrics.JobsSubmittedToGrpcPool++
			params.metrics.mu.Unlock()
			params.logger.WithFields(logrus.Fields{
				"gene_id": syn.GeneID,
				"stage":   "submitted_to_grpc_pool",
			}).Debug("Synonym data submitted for gRPC update")
		}
	}
}

// handleSynonymGrpcResults handles results from gRPC update pool.
func handleSynonymGrpcResults(params *handleSynonymGrpcResultsParams) {
	defer params.wg.Done()
	params.logger.Debug("Starting gRPC results handler for synonyms...")

	for {
		select {
		case <-params.ctx.Done():
			return
		case result, ok := <-params.grpcPool.Results():
			if !ok {
				params.logger.Debug(
					"gRPC update pool results channel closed.",
				)
				return
			}

			params.metrics.mu.Lock()
			params.metrics.TotalProcessed++
			params.metrics.JobsCompletedFromGrpcPool++
			if result.Error != nil {
				if result.Output.Message == "not found" {
					params.metrics.NotFoundCount++
				} else {
					params.metrics.ErrorCount++
				}
				params.metrics.mu.Unlock()
				params.logger.Errorf(
					"Error updating synonyms for gene %s: %v",
					result.Output.GeneID,
					result.Error,
				)
			} else {
				params.metrics.SuccessCount++
				params.metrics.mu.Unlock()
				params.logger.Infof(
					"Successfully processed synonyms for gene %s: %s",
					result.Output.GeneID,
					result.Output.Message,
				)
			}
		case err, ok := <-params.grpcPool.Errors():
			if !ok {
				params.logger.Debug(
					"gRPC update pool errors channel closed.",
				)
				return
			}
			params.metrics.mu.Lock()
			params.metrics.ErrorCount++
			params.metrics.mu.Unlock()
			params.logger.Errorf(
				"Async error from gRPC update pool: %v",
				err,
			)
		}
	}
}

// reportSynonymProgress reports processing progress for synonyms.
func reportSynonymProgress(params *reportSynonymProgressParams) {
	defer params.wg.Done()
	ticker := time.NewTicker(30 * time.Second)
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
			"read_from_db":    params.metrics.TotalFetchedFromArango,
			"total_processed": params.metrics.TotalProcessed,
			"success_count":   params.metrics.SuccessCount,
			"not_found_count": params.metrics.NotFoundCount,
			"error_count":     params.metrics.ErrorCount,
			"processing_rate": fmt.Sprintf("%.2f genes/sec", rate),
			"elapsed_time":    elapsed.String(),
			"grpc_submitted":  params.metrics.JobsSubmittedToGrpcPool,
			"grpc_completed":  params.metrics.JobsCompletedFromGrpcPool,
		}).Info(message)
	}

	for {
		select {
		case <-params.ctx.Done():
			logCurrentMetrics("Final synonym processing report")
			return
		case <-ticker.C:
			logCurrentMetrics("Synonym processing progress")

			if params.metrics.IsComplete() {
				params.logger.Info(
					"All synonym updates appear completed. Stopping progress reporter.",
				)
				return
			}
		}
	}
}

// grpcSynonymWorkerFunc creates a worker for processing synonyms.
func grpcSynonymWorkerFunc(
	config SynonymAppConfig,
	grpcClient fanno.FeatureAnnotationServiceClient,
) concurrent.WorkerFunc[SynonymData, GrpcSynonymResult] {
	return func(
		ctx context.Context,
		job concurrent.Job[SynonymData],
	) (GrpcSynonymResult, error) {
		synData := job.Payload
		logger := config.Logger
		result := GrpcSynonymResult{
			GeneID:  synData.GeneID,
			Success: false,
		}

		// Get existing feature annotation
		featAnno, err := grpcClient.GetFeatureAnnotation(
			ctx,
			&fanno.FeatureAnnotationId{Id: synData.GeneID},
		)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				logger.Warnf(
					"feature annotation not found for gene %s",
					synData.GeneID,
				)
				result.Message = "not found"
				result.Error = err
				return result, err
			}
			result.Error = err
			result.Message = fmt.Sprintf(
				"failed to get feature annotation: %v",
				err,
			)
			return result, err
		}

		// Merge synonyms
		existingSynonyms := mapset.NewSet(
			featAnno.Attributes.Synonyms...)
		newSynonyms := collection.Filter(
			synData.Synonyms,
			func(synonym string) bool {
				return !existingSynonyms.Contains(synonym)
			},
		)
		if len(newSynonyms) == 0 {
			result.Success = true
			result.Message = "All synonyms already exist"
			logger.Debugf(
				"All synonyms for gene %s already exist, skipping",
				synData.GeneID,
			)
			return result, nil
		}
		_, err = grpcClient.UpdateFeatureAnnotation(
			ctx,
			&fanno.FeatureAnnotationUpdate{
				Id:        featAnno.Id,
				UpdatedBy: DefaultUserName,
				Attributes: &fanno.FeatureAnnotationAttributes{
					Synonyms: newSynonyms,
				},
			},
		)
		if err != nil {
			result.Error = err
			result.Message = fmt.Sprintf(
				"failed to update feature annotation: %v",
				err,
			)
			return result, err
		}

		result.Success = true
		result.Message = fmt.Sprintf(
			"Successfully added %d new synonyms",
			len(newSynonyms),
		)
		return result, nil
	}
}
