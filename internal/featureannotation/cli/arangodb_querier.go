package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/registry"
)

// queryArangoParams holds the parameters for the queryArango function.
type queryArangoParams struct {
	ctx            context.Context
	config         AppConfig // Contains AQLQuery, Logger, and Metrics
	arangoDocsChan chan<- ArangoResultDoc
	mainCancel     context.CancelFunc
}

// queryArango is responsible for querying ArangoDB and sending documents to a channel.
// It calls params.mainCancel if critical errors occur or the context is done.
// It also updates params.config.Metrics with the total count of fetched documents
// and a flag indicating when all documents have been fetched.
func queryArango(wg *sync.WaitGroup, params *queryArangoParams) {
	defer wg.Done()
	defer close(
		params.arangoDocsChan,
	) // Ensure channel is always closed on exit

	logger := params.config.Logger // Use logger from AppConfig
	dbh := registry.GetArangodbConnection()
	logger.Info("Starting ArangoDB querier...")

	cursor, err := dbh.Search(params.config.AQLQuery)
	if err != nil {
		logger.Errorf("Error in ArangoDB query: %v", err)
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	defer cursor.Close()

	if cursor.IsEmpty() {
		logger.Warn("ArangoDB querier finished (no data).")
		params.config.Metrics.mu.Lock()
		params.config.Metrics.TotalFetchedFromArango = 0
		params.config.Metrics.AllArangoDocsFetched = true
		params.config.Metrics.mu.Unlock()
		// No need to call mainCancel if empty, let pipeline finish gracefully
		return
	}

	var docCount int64 // Use int64 to match metric type
	for cursor.Scan() {
		select {
		case <-params.ctx.Done():
			logger.Info(
				"ArangoDB querier stopping due to context cancellation during scan.",
			)
			// AllArangoDocsFetched will remain false, progress
			// reporter will stop via ctx.Done()
			return
		default:
			var doc ArangoResultDoc
			if errRead := cursor.Read(&doc); errRead != nil {
				logger.Errorf(
					"Failed to read document from ArangoDB cursor: %v",
					errRead,
				)
				params.mainCancel() // Signal other goroutines to stop
				return
			}
			// Send the document to the channel, handling context cancellation
			select {
			case params.arangoDocsChan <- doc:
				docCount++
			case <-params.ctx.Done():
				logger.Info(
					"ArangoDB querier stopping due to context cancellation during send.",
				)
				// AllArangoDocsFetched will remain false, progress reporter will stop via ctx.Done()
				return
			}
		}
	}

	logger.Infof(
		"Successfully fetched and sent %d documents from ArangoDB.",
		docCount,
	)

	// Update the total count of documents that will be processed by the pipeline.
	params.config.Metrics.mu.Lock()
	params.config.Metrics.TotalFetchedFromArango = docCount
	params.config.Metrics.mu.Unlock()

	// All documents have been fetched and sent successfully
	params.config.Metrics.mu.Lock()
	params.config.Metrics.AllArangoDocsFetched = true
	params.config.Metrics.mu.Unlock()
}
