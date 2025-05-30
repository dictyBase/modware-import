package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/sirupsen/logrus"
)

// ArangoResultDoc represents the structure of a document from ArangoDB.
type ArangoResultDoc struct {
	ID    string           `json:"id"` // This is dbx.accession, likely the feature_id
	Props []ArangoProperty `json:"props"`
}

// ArangoProperty represents a single property object from ArangoDB.
type ArangoProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProcessedGeneData holds the gene ID and its list of HTML-stripped property values.
type ProcessedGeneData struct {
	GeneID            string
	StrippedPropsText []StrippedProperty
}

// StrippedProperty holds the original property name and its stripped text.
type StrippedProperty struct {
	OriginalName string
	StrippedText string
}

// GrpcUpdateResult holds the result of a gRPC update operation.
type GrpcUpdateResult struct {
	GeneID  string
	Success bool
	Message string
	Error   error
}

// bridgeArangoToHTMLPool transfers documents from the ArangoDB channel to the
// HTML processing pool.
func bridgeArangoToHTMLPool(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	arangoDocsFromQueryChan <-chan ArangoResultDoc,
	htmlProcessingPool *concurrent.Pool[ArangoResultDoc, ProcessedGeneData],
	logger *logrus.Entry,
) {
	defer wg.Done()
	defer htmlProcessingPool.Close()

	logger.Info("Starting Arango-to-HTML-Pool bridge...")
	for {
		select {
		case <-mainCtx.Done():
			return
		case arangoDoc, ok := <-arangoDocsFromQueryChan:
			if !ok {
				return
			}
			htmlProcessingPool.Submit(arangoDoc)
		}
	}
}

// bridgeHTMLToGrpcPool transfers processed data from the HTML pool to the gRPC
// update pool.
func bridgeHTMLToGrpcPool(
	wg *sync.WaitGroup,
	mainCtx context.Context,
	htmlProcessingPool *concurrent.Pool[ArangoResultDoc, ProcessedGeneData],
	grpcUpdatePool *concurrent.Pool[ProcessedGeneData, GrpcUpdateResult],
	logger *logrus.Entry,
) {
	defer wg.Done()
	defer grpcUpdatePool.Close()
	logger.Debug("Starting HTML-Pool-to-gRPC-Pool bridge...")
	for {
		select {
		case <-mainCtx.Done():
			return
		case result, ok := <-htmlProcessingPool.Results():
			if !ok {
				return
			}
			if result.Error != nil {
				logger.Errorf(
					"Error from HTML processing pool for job %s: %v",
					result.JobID,
					result.Error,
				)
				return
			}
			grpcUpdatePool.Submit(result.Output)
		case err, ok := <-htmlProcessingPool.Errors():
			if !ok {
				return
			}
			logger.Errorf("Async error from HTML processing pool: %v", err)
		}
	}
}
