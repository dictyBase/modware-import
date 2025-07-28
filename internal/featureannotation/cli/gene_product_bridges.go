package cli

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// bridgeArangoToLegacyPool transfers genes from ArangoDB to legacy query pool
func bridgeArangoToLegacyPool(params *bridgeArangoToLegacyPoolParams) {
	defer params.wg.Done()
	defer params.legacyPool.Close()

	params.logger.Info("Starting Arango-to-Legacy-Pool bridge...")
	for {
		select {
		case <-params.ctx.Done():
			return
		case gene, ok := <-params.genesChan:
			if !ok {
				return
			}
			params.legacyPool.Submit(gene)
			params.metrics.mu.Lock()
			params.metrics.JobsSubmittedToLegacyPool++
			params.metrics.mu.Unlock()
			params.logger.WithFields(logrus.Fields{
				"gene_id": gene.GeneID,
				"stage":   "submitted_to_legacy_pool",
			}).Debug("Gene submitted for legacy query")
		}
	}
}

// bridgeLegacyToGrpcPool transfers processed genes to gRPC update pool
func bridgeLegacyToGrpcPool(params *bridgeLegacyToGrpcPoolParams) {
	defer params.wg.Done()
	defer params.batchGrpcPool.Close()

	params.logger.Debug("Starting Legacy-Pool-to-gRPC-Pool bridge...")
	for {
		select {
		case <-params.ctx.Done():
			params.logger.Debug("Legacy-to-gRPC bridge: context done, exiting.")
			return
		case result, ok := <-params.legacyPool.Results():
			if !ok {
				params.logger.Debug("Legacy query pool results channel closed.")
				return
			}
			params.metrics.mu.Lock()
			params.metrics.JobsCompletedFromLegacyPool++
			params.metrics.mu.Unlock()

			if result.Error != nil {
				params.logger.WithFields(logrus.Fields{
					"job_id": result.JobID,
					"error":  result.Error,
				}).Error("Legacy query failed")
				continue
			}

			// Process each gene product in the slice
			processGeneProductSlice(result.Output, params)

		case err, ok := <-params.legacyPool.Errors():
			if !ok {
				params.logger.Debug("Legacy query pool errors channel closed.")
				return
			}
			params.logger.Errorf("Async error from legacy query pool: %v", err)
		case <-time.After(5 * time.Second): // Added timeout as per plan, though not strictly in bridge logic
			params.logger.Debug(
				"Timeout waiting for legacy query result - continuing",
			)
			continue
		}
	}
}

// processGeneProductSlice processes a slice of gene products and submits them
// to gRPC pool as a batch job.
func processGeneProductSlice(
	geneProducts []ProcessedGeneProduct,
	params *bridgeLegacyToGrpcPoolParams,
) {
	// If no gene products, skip
	if len(geneProducts) == 0 {
		params.logger.Debug("No gene products to process")
		return
	}

	// Get the gene ID from the first product (all should have the same gene ID)
	geneID := geneProducts[0].GeneID

	// Create batch job and submit to batch gRPC pool
	batchJob := BatchGeneProductJob{
		GeneProducts: geneProducts,
	}

	params.batchGrpcPool.Submit(batchJob)
	params.metrics.mu.Lock()
	params.metrics.JobsSubmittedToGrpcPool++
	params.metrics.mu.Unlock()

	params.logger.WithFields(logrus.Fields{
		"gene_id":       geneID,
		"product_count": len(geneProducts),
		"stage":         "submitted_to_batch_grpc_pool",
	}).Debug("Gene products batch submitted for gRPC update")
}

// handleGeneProductGrpcResults handles results from batch gRPC update pool
func handleGeneProductGrpcResults(params *handleGeneProductGrpcResultsParams) {
	defer params.wg.Done()
	params.logger.Debug("Starting batch gRPC results handler...")

	for {
		select {
		case <-params.ctx.Done():
			return
		case result, ok := <-params.batchGrpcPool.Results():
			if !ok {
				params.logger.Debug(
					"Batch gRPC update pool results channel closed.",
				)
				return
			}

			params.metrics.mu.Lock()
			params.metrics.TotalProcessed += int64(result.Output.ProcessedCount)
			params.metrics.SkippedCount += int64(result.Output.SkippedCount)
			params.metrics.JobsCompletedFromGrpcPool++
			if result.Error != nil {
				params.metrics.ErrorCount++
				params.metrics.mu.Unlock()
				params.logger.Errorf(
					"Error updating gene %s: %v",
					result.Output.GeneID,
					result.Error,
				)
			} else {
				params.metrics.SuccessCount++
				params.metrics.mu.Unlock()
				params.logger.WithFields(logrus.Fields{
					"gene_id":         result.Output.GeneID,
					"processed_count": result.Output.ProcessedCount,
					"skipped_count":   result.Output.SkippedCount,
				}).Infof(
					"Successfully processed gene %s: %s",
					result.Output.GeneID,
					result.Output.Message,
				)
			}
		case err, ok := <-params.batchGrpcPool.Errors():
			if !ok {
				params.logger.Debug(
					"Batch gRPC update pool errors channel closed.",
				)
				return
			}
			params.metrics.mu.Lock()
			params.metrics.ErrorCount++
			params.metrics.mu.Unlock()
			params.logger.Errorf(
				"Async error from batch gRPC update pool: %v",
				err,
			)
		}
	}
}

// reportGeneProductProgress reports processing progress
func reportGeneProductProgress(params *reportGeneProductProgressParams) {
	defer params.wg.Done()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logCurrentMetrics := func(message string) {
		params.metrics.mu.RLock() // Changed to RLock for reading metrics
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
			"skipped_count":    params.metrics.SkippedCount,
			"processing_rate":  fmt.Sprintf("%.2f genes/sec", rate),
			"elapsed_time":     elapsed.String(),
			"legacy_submitted": params.metrics.JobsSubmittedToLegacyPool,
			"legacy_completed": params.metrics.JobsCompletedFromLegacyPool,
			"grpc_submitted":   params.metrics.JobsSubmittedToGrpcPool,
			"grpc_completed":   params.metrics.JobsCompletedFromGrpcPool,
		}).Info(message)
	}

	for {
		select {
		case <-params.ctx.Done():
			logCurrentMetrics("Final gene product processing report")
			return
		case <-ticker.C:
			logCurrentMetrics("Gene product processing progress")

			if params.metrics.IsComplete() {
				params.logger.Info(
					"All gene product updates appear completed. Stopping progress reporter.",
				)
				return
			}
		}
	}
}
