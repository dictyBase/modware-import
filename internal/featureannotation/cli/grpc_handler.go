package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/sirupsen/logrus"
)

// handleGrpcResults processes the results from the gRPC update pool and updates metrics.
func handleGrpcResults(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	grpcUpdatePool *concurrent.Pool[ProcessedGeneData, GrpcUpdateResult],
	metrics *ProcessingMetrics, // Add this parameter
	logger *logrus.Entry,
) {
	defer wg.Done()
	logger.Debug("Starting gRPC Update Pool results handler...")
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
					"Error from gRPC update pool for gene %s (Job ID %s): %v",
					result.Output.GeneID,
					result.JobID,
					result.Error,
				)
				// Do not return here, allow other results/errors to be processed
			} else {
				metrics.SuccessCount++
				metrics.mu.Unlock()
				logger.WithFields(logrus.Fields{
					"gene_id": result.Output.GeneID,
					"job_id":  result.JobID,
					"stage":   "grpc_update_completed",
				}).Debug("Job completed from gRPC pool")
				logger.Infof(
					"Successfully updated gene %s (Job ID %s). %s", result.Output.GeneID, result.JobID, result.Output.Message)
			}
		case err, ok := <-grpcUpdatePool.Errors():
			if !ok {
				logger.Debug("gRPC update pool errors channel closed.")
				return
			}
			metrics.mu.Lock()
			metrics.TotalProcessed++ // Assuming an error here means a job was attempted
			metrics.ErrorCount++
			metrics.mu.Unlock()
			logger.Errorf("Async error from gRPC update pool: %v", err)
			// Do not return here, allow other results/errors to be processed
		}
	}
}
