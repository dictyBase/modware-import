package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/registry"
)

// queryActiveGenesParams holds parameters for querying active genes
type queryActiveGenesParams struct {
	config     GeneProductAppConfig
	genesChan  chan<- GeneInfo
	mainCancel context.CancelFunc
}

// queryActiveGenes queries ArangoDB for active genes
func queryActiveGenes(
	wg *sync.WaitGroup,
	params *queryActiveGenesParams,
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
			gene := GeneInfo{}
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
