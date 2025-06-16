package cli

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/sirupsen/logrus"
)

// bridgeArangoToLegacyPool transfers genes from ArangoDB to legacy query pool
func bridgeArangoToLegacyPool(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	genesChan <-chan GeneInfo,
	legacyQueryPool *concurrent.Pool[GeneInfo, ProcessedGeneProduct],
	metrics *GeneProductMetrics,
	logger *logrus.Entry,
) {
	defer wg.Done()
	defer legacyQueryPool.Close()

	logger.Info("Starting Arango-to-Legacy-Pool bridge...")
	for {
		select {
		case <-mainCtx.Done():
			return
		case gene, ok := <-genesChan:
			if !ok {
				return
			}
			legacyQueryPool.Submit(gene)
			metrics.mu.Lock()
			metrics.JobsSubmittedToLegacyPool++
			metrics.mu.Unlock()
			logger.WithFields(logrus.Fields{
				"gene_id": gene.GeneID,
				"stage":   "submitted_to_legacy_pool",
			}).Debug("Gene submitted for legacy query")
		}
	}
}

// bridgeLegacyToGrpcPool transfers processed genes to gRPC update pool
func bridgeLegacyToGrpcPool(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	legacyQueryPool *concurrent.Pool[GeneInfo, ProcessedGeneProduct],
	grpcUpdatePool *concurrent.Pool[ProcessedGeneProduct, GrpcUpdateResult],
	metrics *GeneProductMetrics,
	logger *logrus.Entry,
) {
	defer wg.Done()
	defer grpcUpdatePool.Close()

	logger.Debug("Starting Legacy-Pool-to-gRPC-Pool bridge...")
	for {
		select {
		case <-mainCtx.Done():
			logger.Debug("Legacy-to-gRPC bridge: context done, exiting.")
			return
		case result, ok := <-legacyQueryPool.Results():
			if !ok {
				logger.Debug("Legacy query pool results channel closed.")
				return
			}
			metrics.mu.Lock()
			metrics.JobsCompletedFromLegacyPool++
			metrics.mu.Unlock()

			if result.Error != nil {
				logger.WithFields(logrus.Fields{
					"job_id": result.JobID,
					"error":  result.Error,
				}).Error("Legacy query failed")
				continue
			}

			// Skip if no gene product
			if result.Output.GeneProduct == "" {
				metrics.mu.Lock()
				metrics.SkippedCount++
				metrics.mu.Unlock()
				logger.Debugf(
					"No gene product for gene %s, skipping gRPC update",
					result.Output.GeneID,
				)
				continue
			}

			grpcUpdatePool.Submit(result.Output)
			metrics.mu.Lock()
			metrics.JobsSubmittedToGrpcPool++
			metrics.mu.Unlock()

			logger.WithFields(logrus.Fields{
				"gene_id": result.Output.GeneID,
				"stage":   "submitted_to_grpc_pool",
			}).Debug("Gene product submitted for gRPC update")

		case err, ok := <-legacyQueryPool.Errors():
			if !ok {
				logger.Debug("Legacy query pool errors channel closed.")
				return
			}
			logger.Errorf("Async error from legacy query pool: %v", err)
		case <-time.After(5 * time.Second): // Added timeout as per plan, though not strictly in bridge logic
			logger.Debug("Timeout waiting for legacy query result - continuing")
			continue
		}
	}
}

// handleGeneProductGrpcResults handles results from gRPC update pool
func handleGeneProductGrpcResults(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	grpcUpdatePool *concurrent.Pool[ProcessedGeneProduct, GrpcUpdateResult],
	metrics *GeneProductMetrics,
	logger *logrus.Entry,
) {
	defer wg.Done()
	logger.Debug("Starting gRPC results handler...")

	for {
		select {
		case <-mainCtx.Done():
			return
		case result, ok := <-grpcUpdatePool.Results():
			if !ok {
				logger.Debug("gRPC update pool results channel closed.")
				return
			}

			metrics.mu.Lock()
			metrics.TotalProcessed++
			metrics.JobsCompletedFromGrpcPool++
			if result.Error != nil {
				metrics.ErrorCount++
				metrics.mu.Unlock()
				logger.Errorf(
					"Error updating gene %s: %v",
					result.Output.GeneID,
					result.Error,
				)
			} else {
				metrics.SuccessCount++
				metrics.mu.Unlock()
				logger.Infof(
					"Successfully processed gene %s: %s",
					result.Output.GeneID,
					result.Output.Message,
				)
			}
		case err, ok := <-grpcUpdatePool.Errors():
			if !ok {
				logger.Debug("gRPC update pool errors channel closed.")
				return
			}
			metrics.mu.Lock()
			metrics.ErrorCount++
			metrics.mu.Unlock()
			logger.Errorf("Async error from gRPC update pool: %v", err)
		}
	}
}

// reportGeneProductProgress reports processing progress
func reportGeneProductProgress(
	wg *sync.WaitGroup,
	ctx context.Context,
	metrics *GeneProductMetrics,
	logger *logrus.Entry,
) {
	defer wg.Done()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logCurrentMetrics := func(message string) {
		metrics.mu.RLock() // Changed to RLock for reading metrics
		defer metrics.mu.RUnlock()
		elapsed := time.Since(metrics.StartTime)
		rate := 0.0
		if elapsed.Seconds() > 0 {
			rate = float64(metrics.TotalProcessed) / elapsed.Seconds()
		}
		logger.WithFields(logrus.Fields{
			"read_from_db":     metrics.TotalFetchedFromArango,
			"total_processed":  metrics.TotalProcessed,
			"success_count":    metrics.SuccessCount,
			"error_count":      metrics.ErrorCount,
			"skipped_count":    metrics.SkippedCount,
			"processing_rate":  fmt.Sprintf("%.2f genes/sec", rate),
			"elapsed_time":     elapsed.String(),
			"legacy_submitted": metrics.JobsSubmittedToLegacyPool,
			"legacy_completed": metrics.JobsCompletedFromLegacyPool,
			"grpc_submitted":   metrics.JobsSubmittedToGrpcPool,
			"grpc_completed":   metrics.JobsCompletedFromGrpcPool,
		}).Info(message)
	}

	for {
		select {
		case <-ctx.Done():
			logCurrentMetrics("Final gene product processing report")
			return
		case <-ticker.C:
			logCurrentMetrics("Gene product processing progress")

			metrics.mu.RLock()
			// Condition to stop reporter if all work is done
			if metrics.AllArangoDocsFetched &&
				metrics.JobsCompletedFromLegacyPool >= metrics.JobsSubmittedToLegacyPool && // Ensure legacy pool is drained
				metrics.JobsCompletedFromGrpcPool >= metrics.JobsSubmittedToGrpcPool &&
				(metrics.TotalFetchedFromArango == 0 || // Handle case where no docs were fetched
					metrics.TotalProcessed >= (metrics.TotalFetchedFromArango-metrics.SkippedCount)) { // All fetched (minus skipped) are processed
				logger.Info(
					"All gene product updates appear completed. Stopping progress reporter.",
				)
				metrics.mu.RUnlock()
				return
			}
			metrics.mu.RUnlock()
		}
	}
}
