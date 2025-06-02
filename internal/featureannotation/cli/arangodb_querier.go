package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/registry"
)

// queryArangoParams holds the parameters for the queryArango function.
type queryArangoParams struct {
	ctx            context.Context
	config         AppConfig // Contains AQLQuery and Logger
	arangoDocsChan chan<- ArangoResultDoc
	mainCancel     context.CancelFunc
}

// queryArango is responsible for querying ArangoDB and sending documents to a channel.
// It calls params.mainCancel if critical errors occur or the context is done.
func queryArango(wg *sync.WaitGroup, params *queryArangoParams) {
	defer wg.Done()
	defer close(params.arangoDocsChan)

	dbh := registry.GetArangodbConnection()
	params.config.Logger.Info("Starting ArangoDB querier...")
	cursor, err := dbh.Search(params.config.AQLQuery)
	if err != nil {
		params.config.Logger.Errorf("error in arangodb query %v", err)
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	defer cursor.Close()
	if cursor.IsEmpty() {
		params.config.Logger.Warn("ArangoDB querier finished (no data).")
		params.mainCancel() // Signal other goroutines to stop
		return
	}
	docCount := 0
	for cursor.Scan() {
		select {
		case <-params.ctx.Done():
			return
		default:
			var doc ArangoResultDoc
			if errRead := cursor.Read(&doc); errRead != nil {
				params.config.Logger.Errorf(
					"Failed to read document from ArangoDB cursor: %v",
					errRead,
				)
				params.mainCancel() // Signal other goroutines to stop
			}
			// Send the document to the channel, handling context
			// cancellation
			select {
			case params.arangoDocsChan <- doc:
				docCount++
			case <-params.ctx.Done():
				return
			}
		}
	}
	params.config.Logger.Infof(
		"Successfully fetched %d documents from ArangoDB.",
		docCount,
	)
}
