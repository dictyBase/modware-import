package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Constants for gene product processing
const (
	GeneProductTag = "gene product"
)

// LegacyTime handles Oracle date format "DD-MON-YY" from legacy database
type LegacyTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for Oracle date format
func (lt *LegacyTime) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	// Handle empty or null dates
	if dateStr == "" || strings.ToLower(dateStr) == "null" {
		lt.Time = time.Time{}
		return nil
	}

	// Parse Oracle date format "DD-MON-YY" (e.g., "08-APR-03")
	parsedTime, err := time.Parse("02-Jan-06", dateStr)
	if err != nil {
		return fmt.Errorf("failed to parse legacy date '%s': %w", dateStr, err)
	}

	lt.Time = parsedTime
	return nil
}

// GeneInfo holds gene information from ArangoDB
type GeneInfo struct {
	Name      string `json:"name"`
	GeneID    string `json:"gene_id"`
	FeatureID int64  `json:"feature_id"`
	CreatedBy string `json:"created_by"`
}

// GeneProductResult holds gene product query result
type GeneProductResult struct {
	GeneProduct string     `json:"gene_product"`
	CreatedBy   string     `json:"created_by"`
	CreatedOn   LegacyTime `json:"created_on"`
}

// ProcessedGeneProduct holds processed gene with product
type ProcessedGeneProduct struct {
	GeneID      string
	GeneName    string
	GeneProduct string
	CreatedBy   string
	CreatedOn   time.Time
}

// BatchGeneProductJob holds a slice of gene products for batch processing
type BatchGeneProductJob struct {
	GeneProducts []ProcessedGeneProduct
}

// BatchGeneProductResult holds the result of batch gene product processing
type BatchGeneProductResult struct {
	GeneID         string
	Success        bool
	Message        string
	Error          error
	ProcessedCount int
	SkippedCount   int
}

// GeneProductMetrics holds processing metrics
type GeneProductMetrics struct {
	TotalProcessed              int64
	SuccessCount                int64
	ErrorCount                  int64
	SkippedCount                int64
	StartTime                   time.Time
	mu                          sync.RWMutex
	TotalFetchedFromArango      int64
	AllArangoDocsFetched        bool
	JobsSubmittedToLegacyPool   int64
	JobsCompletedFromLegacyPool int64
	JobsSubmittedToGrpcPool     int64
	JobsCompletedFromGrpcPool   int64
}

// IsComplete checks if all processing is finished
func (m *GeneProductMetrics) IsComplete() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	primary := m.isPrimaryComplete()
	fallback := m.isFallbackComplete()

	// Log completion status for debugging
	if primary || fallback {
		logger := registry.GetLogger()
		logger.WithFields(map[string]any{
			"primary_complete":     primary,
			"fallback_complete":    fallback,
			"all_arango_fetched":   m.AllArangoDocsFetched,
			"legacy_submitted":     m.JobsSubmittedToLegacyPool,
			"legacy_completed":     m.JobsCompletedFromLegacyPool,
			"grpc_submitted":       m.JobsSubmittedToGrpcPool,
			"grpc_completed":       m.JobsCompletedFromGrpcPool,
			"total_processed":      m.TotalProcessed,
			"total_fetched_arango": m.TotalFetchedFromArango,
			"skipped_count":        m.SkippedCount,
		}).Debug("Completion conditions met")
	}

	return primary || fallback
}

// isPrimaryComplete checks the primary completion conditions
func (m *GeneProductMetrics) isPrimaryComplete() bool {
	allArangoFetched := m.AllArangoDocsFetched
	legacyPoolDrained := m.JobsCompletedFromLegacyPool >= m.JobsSubmittedToLegacyPool
	grpcPoolDrained := m.JobsCompletedFromGrpcPool >= m.JobsSubmittedToGrpcPool
	allProcessed := (m.TotalFetchedFromArango == 0) ||
		(m.TotalProcessed >= (m.TotalFetchedFromArango - m.SkippedCount))

	return allArangoFetched && legacyPoolDrained && grpcPoolDrained &&
		allProcessed
}

// isFallbackComplete checks fallback completion conditions for edge cases
func (m *GeneProductMetrics) isFallbackComplete() bool {
	elapsed := time.Since(m.StartTime)
	hasSignificantWork := m.TotalFetchedFromArango > 0
	hasProcessedSomething := (m.SuccessCount + m.ErrorCount) > 0
	poolsDrained := m.JobsCompletedFromLegacyPool >= m.JobsSubmittedToLegacyPool &&
		m.JobsCompletedFromGrpcPool >= m.JobsSubmittedToGrpcPool

	// More lenient fallback conditions to prevent hanging
	basicCompletion := m.AllArangoDocsFetched && poolsDrained
	timeoutCompletion := elapsed > 2*time.Minute && hasProcessedSomething &&
		poolsDrained

	// If we have processed some work and pools are drained, we're likely done
	workCompletion := hasSignificantWork && hasProcessedSomething &&
		poolsDrained &&
		elapsed > 30*time.Second

	return basicCompletion || timeoutCompletion || workCompletion
}

// GeneProductAppConfig holds configuration
type GeneProductAppConfig struct {
	Ctx              context.Context
	LegacyDatabase   string
	NumLegacyWorkers int
	NumGrpcWorkers   int
	Logger           *logrus.Entry
	Metrics          *GeneProductMetrics
}

type bridgeArangoToLegacyPoolParams struct {
	wg         *sync.WaitGroup
	ctx        context.Context
	genesChan  <-chan GeneInfo
	legacyPool *concurrent.Pool[GeneInfo, []ProcessedGeneProduct]
	metrics    *GeneProductMetrics
	logger     *logrus.Entry
}

type bridgeLegacyToGrpcPoolParams struct {
	wg            *sync.WaitGroup
	ctx           context.Context
	legacyPool    *concurrent.Pool[GeneInfo, []ProcessedGeneProduct]
	batchGrpcPool *concurrent.Pool[BatchGeneProductJob, BatchGeneProductResult]
	metrics       *GeneProductMetrics
	logger        *logrus.Entry
}

type handleGeneProductGrpcResultsParams struct {
	wg            *sync.WaitGroup
	ctx           context.Context
	batchGrpcPool *concurrent.Pool[BatchGeneProductJob, BatchGeneProductResult]
	metrics       *GeneProductMetrics
	logger        *logrus.Entry
}

type reportGeneProductProgressParams struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	metrics *GeneProductMetrics
	logger  *logrus.Entry
}

// newGeneProductConfigFromCliContext creates config from CLI context
func newGeneProductConfigFromCliContext(
	cltx *cli.Context,
	logger *logrus.Entry,
) GeneProductAppConfig {
	return GeneProductAppConfig{
		LegacyDatabase:   cltx.String("legacy-database"),
		NumLegacyWorkers: cltx.Int("legacy-workers"),
		NumGrpcWorkers:   cltx.Int("grpc-workers"),
		Logger:           logger,
		Metrics: &GeneProductMetrics{
			StartTime: time.Now(),
		},
	}
}

func setupLegacyQueryPool(
	config GeneProductAppConfig,
	mainCtx context.Context,
) *concurrent.Pool[GeneInfo, []ProcessedGeneProduct] {
	pool := concurrent.NewPool(
		legacyDBQueryWorkerFunc(config),
		concurrent.WithWorkers[GeneInfo, []ProcessedGeneProduct](
			config.NumLegacyWorkers,
		),
		concurrent.WithContext[GeneInfo, []ProcessedGeneProduct](mainCtx),
		concurrent.WithBufferSize[GeneInfo, []ProcessedGeneProduct](
			config.NumLegacyWorkers*bufferSizeMultiplier,
		),
	)
	pool.Start()
	return pool
}

// setupBatchGrpcUpdatePool sets up batch gRPC update pool for gene products
func setupBatchGrpcUpdatePool(
	config GeneProductAppConfig,
	mainCtx context.Context,
) *concurrent.Pool[BatchGeneProductJob, BatchGeneProductResult] {
	pool := concurrent.NewPool(
		batchGeneProductGrpcWorkerFunc(
			config,
			registry.GetFeatureAnnotationAPIClient(),
		),
		concurrent.WithWorkers[BatchGeneProductJob, BatchGeneProductResult](
			config.NumGrpcWorkers,
		),
		concurrent.WithContext[BatchGeneProductJob, BatchGeneProductResult](
			mainCtx,
		),
		concurrent.WithBufferSize[BatchGeneProductJob, BatchGeneProductResult](
			config.NumGrpcWorkers*bufferSizeMultiplier,
		),
	)
	pool.Start()
	return pool
}

// RunGeneProductUpdater is the main entry point for the gene product updater
func RunGeneProductUpdater(cltx *cli.Context) error {
	logger := registry.GetLogger()
	config := newGeneProductConfigFromCliContext(cltx, logger)

	logger.Debug("Starting gene product updater...")
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()
	config.Ctx = mainCtx

	setupSignalHandling(
		mainCancel,
		logger,
	)
	// Channel for genes from ArangoDB
	genesFromQueryChan := make(
		chan GeneInfo,
		config.NumLegacyWorkers,
	)

	wg := sync.WaitGroup{}
	// Goroutine 1: ArangoDB Active Genes Querier
	wg.Add(1)
	go queryActiveGenes(&wg, &queryActiveGenesParams{
		config:     config,
		genesChan:  genesFromQueryChan,
		mainCancel: mainCancel,
	})

	// Setup Pools
	legacyQueryPool := setupLegacyQueryPool(config, mainCtx)
	batchGrpcPool := setupBatchGrpcUpdatePool(config, mainCtx)

	// Bridge from ArangoDB to Legacy Query Pool
	wg.Add(1)
	go bridgeArangoToLegacyPool(&bridgeArangoToLegacyPoolParams{
		wg:         &wg,
		ctx:        mainCtx,
		genesChan:  genesFromQueryChan,
		legacyPool: legacyQueryPool,
		metrics:    config.Metrics,
		logger:     logger,
	})

	// Bridge from Legacy Query Pool to gRPC Pool
	wg.Add(1)
	go bridgeLegacyToGrpcPool(&bridgeLegacyToGrpcPoolParams{
		wg:            &wg,
		ctx:           mainCtx,
		legacyPool:    legacyQueryPool,
		batchGrpcPool: batchGrpcPool,
		metrics:       config.Metrics,
		logger:        logger,
	})

	// Handle gRPC Results
	wg.Add(1)
	go handleGeneProductGrpcResults(&handleGeneProductGrpcResultsParams{
		wg:            &wg,
		ctx:           mainCtx,
		batchGrpcPool: batchGrpcPool,
		metrics:       config.Metrics,
		logger:        logger,
	})

	// Progress Reporter
	wg.Add(1)
	go reportGeneProductProgress(&reportGeneProductProgressParams{
		wg:      &wg,
		ctx:     mainCtx,
		metrics: config.Metrics,
		logger:  logger,
	})

	logger.Debug("Waiting for all goroutines to complete...")
	wg.Wait()
	logger.Info("Gene product updater finished")
	return nil
}
