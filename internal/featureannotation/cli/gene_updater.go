package cli

import (
	"context"
	"sync"

	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// DefaultAQLQuery is the default query to fetch gene data from ArangoDB.
// Exported for use in flag.go
const DefaultAQLQuery = `
FOR ftype IN cvterm
 FOR feat IN feature
    FOR dbx IN dbxref
        FILTER ftype.name == 'gene'
        FILTER feat.type_id == ftype.cvterm_id
        FILTER feat.dbxref_id == dbx.dbxref_id
        LET props = (
            FOR fprop IN featureprop
                FOR cvt IN cvterm
                    FILTER cvt.name IN ['description','name description']
                    FILTER feat.feature_id == fprop.feature_id
                    FILTER fprop.type_id == cvt.cvterm_id
                    RETURN {
                        name: cvt.name,
                        value: fprop.value
                    }
        )
        FILTER LENGTH(props) > 0
        RETURN {
            id: dbx.accession,
            props: props
        }
`

// AppConfig holds all configuration for the application.
type AppConfig struct {
	AQLQuery             string
	ArangoUser           string // For authorship in gRPC updates
	NumProcessingWorkers int
	NumGrpcWorkers       int
	Logger               *logrus.Entry
}

// newAppConfigFromCliContext creates an AppConfig from CLI context and a logger.
func newAppConfigFromCliContext(
	cltx *cli.Context,
	logger *logrus.Entry,
) AppConfig {
	return AppConfig{
		AQLQuery:             cltx.String("aql-query"),
		ArangoUser:           cltx.String("arangodb-user"), // For authorship
		NumProcessingWorkers: cltx.Int("processing-workers"),
		NumGrpcWorkers:       cltx.Int("grpc-workers"),
		Logger:               logger,
	}
}

func RunGeneUpdater(cltx *cli.Context) error {
	logger := registry.GetLogger()
	config := newAppConfigFromCliContext(cltx, logger)

	logger.Debug("Starting gene updater application...")
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()
	setupSignalHandling(mainCancel, logger)

	arangoDocsFromQueryChan := make(
		chan ArangoResultDoc,
		config.NumProcessingWorkers,
	)
	wg := sync.WaitGroup{} // Goroutine 1: ArangoDB Querier
	wg.Add(1)
	go queryArango(&wg, &queryArangoParams{
		ctx:            mainCtx,
		config:         config,
		arangoDocsChan: arangoDocsFromQueryChan,
		mainCancel:     mainCancel,
	})
	// Setup HTML Processing Pool
	htmlProcessingPool := concurrent.NewPool(
		htmlProcessingWorkerFunc(config.Logger),
		concurrent.WithWorkers[ArangoResultDoc, ProcessedGeneData](
			config.NumProcessingWorkers,
		),
		concurrent.WithContext[ArangoResultDoc, ProcessedGeneData](mainCtx),
		concurrent.WithBufferSize[ArangoResultDoc, ProcessedGeneData](
			config.NumProcessingWorkers*2,
		),
	)
	htmlProcessingPool.Start()
	grpcUpdatePool := concurrent.NewPool(
		grpcUpdateWorkerFunc(
			config,
			registry.GetFeatureAnnotationAPIClient()),
		concurrent.WithWorkers[ProcessedGeneData, GrpcUpdateResult](
			config.NumGrpcWorkers,
		),
		concurrent.WithContext[ProcessedGeneData, GrpcUpdateResult](mainCtx),
		concurrent.WithBufferSize[ProcessedGeneData, GrpcUpdateResult](
			config.NumGrpcWorkers*2,
		),
	)
	grpcUpdatePool.Start()
	// Start bridge from Arango Docs to HTML Processing Pool goroutine
	wg.Add(1)
	go bridgeArangoToHTMLPool(
		&wg,
		mainCtx,
		arangoDocsFromQueryChan,
		htmlProcessingPool,
		logger,
	)
	// Start bridge from HTML Processing Results to gRPC Update Pool goroutine
	wg.Add(1)
	go bridgeHTMLToGrpcPool(
		&wg,
		mainCtx,
		htmlProcessingPool,
		grpcUpdatePool,
		logger,
	)
	// Start gRPC Update Results Handler goroutine
	wg.Add(1)
	go handleGrpcResults(&wg, mainCtx, grpcUpdatePool, logger)
	logger.Debug("Waiting for all main goroutines to complete...")
	wg.Wait()
	logger.Info("Gene updater application finished")
	return nil
}
