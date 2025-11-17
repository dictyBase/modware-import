package cli

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	geneProductBridgeTimeoutSeconds = 10 // Timeout for gene product bridge operations
	progressReportIntervalSeconds   = 30 // Interval for progress reporting
)

// handleLegacyResultsChannelClosure handles the closure of the legacy pool results channel
func handleLegacyResultsChannelClosure(errorsChan <-chan error, logger *logrus.Entry) {
	logger.Debug("Legacy query pool results channel closed.")
	// Check if errors channel is also closed
	select {
	case err, errOk := <-errorsChan:
		if !errOk {
			logger.Debug("Legacy query pool errors channel also closed, exiting bridge.")
		} else {
			logger.Errorf("Final async error from legacy query pool: %v", err)
		}
	default:
		// No error to process, continue
	}
}

// handleLegacyErrorsChannelClosure handles the closure of the legacy pool errors channel
func handleLegacyErrorsChannelClosure(params *bridgeLegacyToGrpcPoolParams) bool {
	params.logger.Debug("Legacy query pool errors channel closed.")
	return true // Always exit when errors channel closes
}

// handleGrpcResultsChannelClosure handles closure of the gRPC results channel
func handleGrpcResultsChannelClosure(errorsChan <-chan error, logger *logrus.Entry) {
	logger.Debug("Batch gRPC update pool results channel closed.")
	// Check if errors channel is also closed
	select {
	case err, errOk := <-errorsChan:
		if !errOk {
			logger.Debug("Batch gRPC update pool errors channel also closed, exiting handler.")
		} else {
			logger.Errorf("Final async error from batch gRPC update pool: %v", err)
		}
	default:
		// No error to process, continue
	}
}

// handleGrpcErrorsChannelClosure handles closure of the gRPC errors channel
func handleGrpcErrorsChannelClosure(params *handleGeneProductGrpcResultsParams) bool {
	params.logger.Debug("Batch gRPC update pool errors channel closed.")
	return true
}

// processLegacyPoolResult processes a successful result from the legacy pool
func processLegacyPoolResult(
	resultOutput []ProcessedGeneProduct,
	resultError error,
	jobID string,
	params *bridgeLegacyToGrpcPoolParams,
) {
	params.metrics.mu.Lock()
	params.metrics.JobsCompletedFromLegacyPool++
	params.metrics.mu.Unlock()

	if resultError != nil {
		params.logger.WithFields(logrus.Fields{
			"job_id": jobID,
			"error":  resultError,
		}).Error("Legacy query failed")
		return
	}

	// Process each gene product in the slice
	processGeneProductSlice(resultOutput, params)
}

// updateGrpcMetrics updates metrics and logs for a gRPC result
func updateGrpcMetrics(
	result BatchGeneProductResult,
	error error,
	params *handleGeneProductGrpcResultsParams,
) {
	params.metrics.mu.Lock()
	params.metrics.TotalProcessed += int64(result.ProcessedCount)
	params.metrics.SkippedCount += int64(result.SkippedCount)
	params.metrics.JobsCompletedFromGrpcPool++

	if error != nil {
		params.metrics.ErrorCount++
		params.metrics.mu.Unlock()
		params.logger.Errorf(
			"Error updating gene %s: %v",
			result.GeneID,
			error,
		)
	} else {
		params.metrics.SuccessCount++
		params.metrics.mu.Unlock()
		params.logger.WithFields(logrus.Fields{
			"gene_id":         result.GeneID,
			"processed_count": result.ProcessedCount,
			"skipped_count":   result.SkippedCount,
		}).Infof(
			"Successfully processed gene %s: %s",
			result.GeneID,
			result.Message,
		)
	}
}

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
	defer func() {
		params.logger.Debug("Closing batch gRPC pool...")
		params.batchGrpcPool.Close()
		params.logger.Debug("Batch gRPC pool closed")
	}()

	params.logger.Debug("Starting Legacy-Pool-to-gRPC-Pool bridge...")

	resultsChan := params.legacyPool.Results()
	errorsChan := params.legacyPool.Errors()

	for {
		select {
		case <-params.ctx.Done():
			params.logger.Debug("Legacy-to-gRPC bridge: context done, exiting.")
			return
		case result, ok := <-resultsChan:
			if !ok {
				handleLegacyResultsChannelClosure(errorsChan, params.logger)
				return
			} else {
				processLegacyPoolResult(result.Output, result.Error, result.JobID, params)
			}

		case err, ok := <-errorsChan:
			if !ok {
				if handleLegacyErrorsChannelClosure(params) {
					return
				}
			}
			params.logger.Errorf("Async error from legacy query pool: %v", err)
		case <-time.After(geneProductBridgeTimeoutSeconds * time.Second):
			// Longer timeout with completion check
			params.logger.Debug(
				"Timeout waiting for legacy query result - checking completion status",
			)
			if params.metrics.IsComplete() {
				params.logger.Debug("Metrics indicate completion, exiting legacy-to-gRPC bridge")
				return
			}
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

	resultsChan := params.batchGrpcPool.Results()
	errorsChan := params.batchGrpcPool.Errors()

	for {
		select {
		case <-params.ctx.Done():
			params.logger.Debug("gRPC results handler: context done, exiting.")
			return
		case result, ok := <-resultsChan:
			if !ok {
				handleGrpcResultsChannelClosure(errorsChan, params.logger)
				return
			} else {
				updateGrpcMetrics(result.Output, result.Error, params)
			}
		case err, ok := <-errorsChan:
			if !ok {
				if handleGrpcErrorsChannelClosure(params) {
					return
				}
			} else {
				params.metrics.mu.Lock()
				params.metrics.ErrorCount++
				params.metrics.mu.Unlock()
				params.logger.Errorf(
					"Async error from batch gRPC update pool: %v",
					err,
				)
			}
		case <-time.After(geneProductBridgeTimeoutSeconds * time.Second):
			// Check completion status during timeout
			params.logger.Debug("Timeout waiting for gRPC results - checking completion status")
			if params.metrics.IsComplete() {
				params.logger.Debug("Metrics indicate completion, exiting gRPC results handler")
				return
			}
			continue
		}
	}
}

// reportGeneProductProgress reports processing progress
func reportGeneProductProgress(params *reportGeneProductProgressParams) {
	defer params.wg.Done()
	ticker := time.NewTicker(progressReportIntervalSeconds * time.Second)
	defer ticker.Stop()

	logCurrentMetrics := func(message string) {
		params.metrics.mu.RLock() // Changed to RLock for reading metrics
		defer params.metrics.mu.RUnlock()
		elapsed := time.Since(params.metrics.StartTime)
		rate := 0.0
		if elapsed.Seconds() > 0 {
			rate = float64(params.metrics.TotalProcessed) / elapsed.Seconds()
		}

		// Check completion status for debugging
		isComplete := params.metrics.IsComplete()

		params.logger.WithFields(logrus.Fields{
			"read_from_db":       params.metrics.TotalFetchedFromArango,
			"total_processed":    params.metrics.TotalProcessed,
			"success_count":      params.metrics.SuccessCount,
			"error_count":        params.metrics.ErrorCount,
			"skipped_count":      params.metrics.SkippedCount,
			"processing_rate":    fmt.Sprintf("%.2f genes/sec", rate),
			"elapsed_time":       elapsed.String(),
			"legacy_submitted":   params.metrics.JobsSubmittedToLegacyPool,
			"legacy_completed":   params.metrics.JobsCompletedFromLegacyPool,
			"grpc_submitted":     params.metrics.JobsSubmittedToGrpcPool,
			"grpc_completed":     params.metrics.JobsCompletedFromGrpcPool,
			"all_arango_fetched": params.metrics.AllArangoDocsFetched,
			"is_complete":        isComplete,
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
