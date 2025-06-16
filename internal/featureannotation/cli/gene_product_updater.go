package cli

import (
	"context"
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

// GeneInfo holds gene information from ArangoDB
type GeneInfo struct {
	Name      string `json:"name"`
	GeneID    string `json:"gene_id"`
	FeatureID int64  `json:"feature_id"`
	CreatedBy string `json:"created_by"`
}

// GeneProductResult holds gene product query result
type GeneProductResult struct {
	GeneProduct string `json:"gene_product"`
	CreatedBy   string `json:"created_by"`
}

// ProcessedGeneProduct holds processed gene with product
type ProcessedGeneProduct struct {
	GeneID      string
	GeneName    string
	GeneProduct string
	CreatedBy   string
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

// GeneProductAppConfig holds configuration
type GeneProductAppConfig struct {
	Ctx              context.Context
	LegacyDatabase   string
	NumLegacyWorkers int
	NumGrpcWorkers   int
	Logger           *logrus.Entry
	Metrics          *GeneProductMetrics
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
) *concurrent.Pool[GeneInfo, ProcessedGeneProduct] {
	pool := concurrent.NewPool(
		legacyDBQueryWorkerFunc(config),
		concurrent.WithWorkers[GeneInfo, ProcessedGeneProduct](
			config.NumLegacyWorkers,
		),
		concurrent.WithContext[GeneInfo, ProcessedGeneProduct](mainCtx),
		concurrent.WithBufferSize[GeneInfo, ProcessedGeneProduct](
			config.NumLegacyWorkers*2,
		),
	)
	pool.Start()
	return pool
}

func setupGrpcUpdatePool(
	config GeneProductAppConfig,
	mainCtx context.Context,
) *concurrent.Pool[ProcessedGeneProduct, GrpcUpdateResult] {
	pool := concurrent.NewPool(
		geneProductGrpcWorkerFunc(
			config,
			registry.GetFeatureAnnotationAPIClient(),
		),
		concurrent.WithWorkers[ProcessedGeneProduct, GrpcUpdateResult](
			config.NumGrpcWorkers,
		),
		concurrent.WithContext[ProcessedGeneProduct, GrpcUpdateResult](mainCtx),
		concurrent.WithBufferSize[ProcessedGeneProduct, GrpcUpdateResult](
			config.NumGrpcWorkers*2,
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
	grpcUpdatePool := setupGrpcUpdatePool(config, mainCtx)

	// Bridge from ArangoDB to Legacy Query Pool
	wg.Add(1)
	go bridgeArangoToLegacyPool(
		&wg,
		mainCtx,
		genesFromQueryChan,
		legacyQueryPool,
		config.Metrics,
		logger,
	)

	// Bridge from Legacy Query Pool to gRPC Pool
	wg.Add(1)
	go bridgeLegacyToGrpcPool(
		&wg,
		mainCtx,
		legacyQueryPool,
		grpcUpdatePool,
		config.Metrics,
		logger,
	)

	// Handle gRPC Results
	wg.Add(1)
	go handleGeneProductGrpcResults(
		&wg,
		mainCtx,
		grpcUpdatePool,
		config.Metrics,
		logger,
	)

	// Progress Reporter
	wg.Add(1)
	go reportGeneProductProgress(&wg, mainCtx, config.Metrics, logger)

	logger.Debug("Waiting for all goroutines to complete...")
	wg.Wait()
	logger.Info("Gene product updater finished")
	return nil
}
