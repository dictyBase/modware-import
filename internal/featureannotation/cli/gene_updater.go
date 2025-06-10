package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// ProcessingMetrics holds counters for tracking progress.
type ProcessingMetrics struct {
	TotalProcessed int64
	SuccessCount   int64
	ErrorCount     int64
	StartTime      time.Time
	mu             sync.RWMutex
	// TotalFetchedFromArango stores the total number of items fetched by queryArango.
	// This field is set by queryArango once the total count is known.
	TotalFetchedFromArango int64
	// AllArangoDocsFetched is a flag set to true by queryArango after all documents
	// have been fetched and sent to the processing pipeline.
	AllArangoDocsFetched bool
	// Intermediate tracking counters for detailed pipeline monitoring
	JobsSubmittedToHTMLPool   int64
	JobsCompletedFromHTMLPool int64
	JobsSubmittedToGrpcPool   int64
	JobsCompletedFromGrpcPool int64
}

// DefaultAQLQuery is the default query to fetch gene data from ArangoDB.
// Exported for use in flag.go
const DefaultAQLQuery = `
FOR ftype IN cvterm
 FOR feat IN feature
    FOR dbx IN dbxref
        FILTER ftype.name == 'gene'
        FILTER feat.type_id == ftype.cvterm_id
        FILTER feat.dbxref_id == dbx.dbxref_id
        LET props = (
            FOR fprop IN featureprop
                FOR cvt IN cvterm
                    FILTER cvt.name IN ['description','name description']
                    FILTER feat.feature_id == fprop.feature_id
                    FILTER fprop.type_id == cvt.cvterm_id
                    RETURN {
                        name: cvt.name,
                        value: fprop.value
                    }
        )
        FILTER LENGTH(props) > 0
        RETURN {
            id: dbx.accession,
            props: props
        }
`

// AppConfig holds all configuration for the application.
type AppConfig struct {
	AQLQuery             string
	ArangoUser           string // For authorship in gRPC updates
	NumProcessingWorkers int
	NumGrpcWorkers       int
	Logger               *logrus.Entry
	Metrics              *ProcessingMetrics // Add this field
}

// newAppConfigFromCliContext creates an AppConfig from CLI context and a logger.
func newAppConfigFromCliContext(
	cltx *cli.Context,
	logger *logrus.Entry,
) AppConfig {
	return AppConfig{
		AQLQuery:             cltx.String("aql-query"),
		ArangoUser:           cltx.String("arangodb-user"), // For authorship
		NumProcessingWorkers: cltx.Int("processing-workers"),
		NumGrpcWorkers:       cltx.Int("grpc-workers"),
		Logger:               logger,
		Metrics: &ProcessingMetrics{
			StartTime: time.Now(),
		},
	}
}

// reportProgress periodically logs processing metrics.
func reportProgress(
	wg *sync.WaitGroup,
	ctx context.Context,
	metrics *ProcessingMetrics,
	logger *logrus.Entry,
) {
	defer wg.Done()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Helper function to log metrics
	logCurrentMetrics := func(message string) {
		elapsed := time.Since(metrics.StartTime)
		rate := 0.0
		if elapsed.Seconds() > 0 {
			rate = float64(metrics.TotalProcessed) / elapsed.Seconds()
		}
		logger.WithFields(logrus.Fields{
			"read_from_db":    metrics.TotalFetchedFromArango,
			"total_processed": metrics.TotalProcessed,
			"success_count":   metrics.SuccessCount,
			"error_count":     metrics.ErrorCount,
			"processing_rate": fmt.Sprintf("%.2f genes/sec", rate),
			"elapsed_time":    elapsed.String(),
			"html_submitted":  metrics.JobsSubmittedToHTMLPool,
			"html_completed":  metrics.JobsCompletedFromHTMLPool,
			"grpc_submitted":  metrics.JobsSubmittedToGrpcPool,
			"grpc_completed":  metrics.JobsCompletedFromGrpcPool,
		}).Info(message)
	}

	for {
		select {
		case <-ctx.Done():
			metrics.mu.RLock()
			logCurrentMetrics("Final processing report on cancellation")
			metrics.mu.RUnlock()
			logger.Debug(
				"Stopping progress reporter due to context cancellation.",
			)
			return
		case <-ticker.C:
			metrics.mu.RLock()
			logCurrentMetrics("Processing progress")

			// Check if all jobs have been submitted AND all submitted jobs are complete
			if metrics.AllArangoDocsFetched &&
				metrics.JobsCompletedFromGrpcPool >= metrics.JobsSubmittedToGrpcPool {
				logger.Info("All jobs completed. Stopping progress reporter.")
				metrics.mu.RUnlock() // Release lock before returning
				return
			}
			metrics.mu.RUnlock()
		}
	}
}

func RunGeneUpdater(cltx *cli.Context) error {
	logger := registry.GetLogger()
	config := newAppConfigFromCliContext(cltx, logger)

	logger.Debug("Starting gene updater application...")
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()
	setupSignalHandling(mainCancel, logger)

	arangoDocsFromQueryChan := make(
		chan ArangoResultDoc,
		config.NumProcessingWorkers,
	)
	wg := sync.WaitGroup{} // Goroutine 1: ArangoDB Querier
	wg.Add(1)
	go queryArango(&wg, &queryArangoParams{
		ctx:            mainCtx,
		config:         config,
		arangoDocsChan: arangoDocsFromQueryChan,
		mainCancel:     mainCancel,
	})
	// Setup HTML Processing Pool
	htmlProcessingPool := concurrent.NewPool(
		htmlProcessingWorkerFunc(config.Logger),
		concurrent.WithWorkers[ArangoResultDoc, ProcessedGeneData](
			config.NumProcessingWorkers,
		),
		concurrent.WithContext[ArangoResultDoc, ProcessedGeneData](mainCtx),
		concurrent.WithBufferSize[ArangoResultDoc, ProcessedGeneData](
			config.NumProcessingWorkers*2,
		),
	)
	htmlProcessingPool.Start()
	grpcUpdatePool := concurrent.NewPool(
		grpcUpdateWorkerFunc(
			config,
			registry.GetFeatureAnnotationAPIClient()),
		concurrent.WithWorkers[ProcessedGeneData, GrpcUpdateResult](
			config.NumGrpcWorkers,
		),
		concurrent.WithContext[ProcessedGeneData, GrpcUpdateResult](mainCtx),
		concurrent.WithBufferSize[ProcessedGeneData, GrpcUpdateResult](
			config.NumGrpcWorkers*2,
		),
	)
	grpcUpdatePool.Start()
	// Start bridge from Arango Docs to HTML Processing Pool goroutine
	wg.Add(1)
	go bridgeArangoToHTMLPool(
		&wg,
		mainCtx,
		arangoDocsFromQueryChan,
		htmlProcessingPool,
		config.Metrics,
		logger,
	)
	// Start bridge from HTML Processing Results to gRPC Update Pool goroutine
	wg.Add(1)
	go bridgeHTMLToGrpcPool(
		&wg,
		mainCtx,
		htmlProcessingPool,
		grpcUpdatePool,
		config.Metrics,
		logger,
	)
	// Start gRPC Update Results Handler goroutine
	wg.Add(1)
	go handleGrpcResults(&wg, mainCtx, grpcUpdatePool, config.Metrics, logger)

	// Start Progress Reporter goroutine
	wg.Add(1)
	go reportProgress(&wg, mainCtx, config.Metrics, logger)

	logger.Debug("Waiting for all main goroutines to complete...")
	wg.Wait()
	logger.Info("Gene updater application finished")
	return nil
}
