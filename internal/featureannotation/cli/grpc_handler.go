package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/sirupsen/logrus"
)

// handleGrpcResults processes the results from the gRPC update pool.
func handleGrpcResults(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	grpcUpdatePool *concurrent.Pool[ProcessedGeneData, GrpcUpdateResult],
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
				return
			}
			if result.Error != nil {
				logger.Errorf(
					"Error from gRPC update pool for gene %s (Job ID %s): %v",
					result.Output.GeneID,
					result.JobID,
					result.Error,
				)
				return
			} else {
				logger.Infof(
					"Successfully updated gene %s (Job ID %s). %s", result.Output.GeneID, result.JobID, result.Output.Message)
			}
		case err, ok := <-grpcUpdatePool.Errors():
			if !ok {
				return
			}
			logger.Errorf("Async error from gRPC update pool: %v", err)
			return
		}
	}
}
