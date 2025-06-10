package cli

import (
	"context"
	"sync"
	"time"

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
	metrics *ProcessingMetrics,
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
			metrics.mu.Lock()
			metrics.JobsSubmittedToHTMLPool++
			metrics.mu.Unlock()
			logger.WithFields(logrus.Fields{
				"job_id": arangoDoc.ID,
				"stage":  "submitted_to_html_pool",
			}).Debug("Job submitted for HTML processing")
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
	metrics *ProcessingMetrics,
	logger *logrus.Entry,
) {
	defer wg.Done()
	defer grpcUpdatePool.Close()
	logger.Debug("Starting HTML-Pool-to-gRPC-Pool bridge...")
	for {
		select {
		case <-mainCtx.Done():
			logger.Debug("HTML-Pool-to-gRPC-Pool bridge: main context done, exiting.")
			return
		case result, ok := <-htmlProcessingPool.Results():
			if !ok {
				logger.Debug("HTML-Pool-to-gRPC-Pool bridge: HTML processing results channel closed.")
				return
			}
			metrics.mu.Lock()
			metrics.JobsCompletedFromHTMLPool++
			metrics.mu.Unlock()
			if result.Error != nil {
				logger.WithFields(logrus.Fields{
					"job_id": result.JobID,
					"error":  result.Error,
				}).Debug("HTML processing failed for job")
				logger.Errorf(
					"HTML processing error for job %s: %v",
					result.JobID,
					result.Error,
				)
				continue
			}
			logger.WithFields(logrus.Fields{
				"job_id":  result.JobID,
				"gene_id": result.Output.GeneID,
			}).Debug("HTML processing successful for job, submitting to gRPC pool")
			grpcUpdatePool.Submit(result.Output)
			metrics.mu.Lock()
			metrics.JobsSubmittedToGrpcPool++
			metrics.mu.Unlock()
			logger.WithFields(logrus.Fields{
				"job_id":  result.JobID,
				"gene_id": result.Output.GeneID,
				"stage":   "submitted_to_grpc_pool",
			}).Debug("Job submitted for gRPC update")
		case err, ok := <-htmlProcessingPool.Errors():
			if !ok {
				logger.Debug("HTML-Pool-to-gRPC-Pool bridge: HTML processing errors channel closed.")
				return
			}
			logger.Errorf("Async error from HTML processing pool: %v", err)
		case <-time.After(5 * time.Second):
			logger.Debug("Timeout waiting for HTML processing result - continuing")
			continue
		}
	}
}
