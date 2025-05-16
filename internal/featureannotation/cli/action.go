package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func LoadCSVToArangodb(cltx *cli.Context) error {
	logger := registry.GetLogger()
	result := collection.Pipe5(
		cltx,
		setupPipeline,
		setupFileProcessing,
		validateHeaders,
		processCSVRecords,
		finalizeBatchProcessing,
	)

	// Check the final result for errors
	if result.Error != nil {
		logger.Errorf("Pipeline failed: %s", result.Error)
		return cli.Exit(
			fmt.Sprintf("Pipeline execution failed: %s", result.Error.Error()),
			2,
		)
	}

	// If successful, log the final count (already logged in
	// finalizeBatchProcessing, but good for confirmation)
	logger.Infof(
		"Pipeline completed successfully. Total documents updated: %d",
		result.TotalUpdated,
	)

	return nil
}

func LoadFeatureAnnotation(cltx *cli.Context) error {
	logger := registry.GetLogger()
	dbh := registry.GetArangodbConnection()
	client := registry.GetFeatureAnnotationAPIClient()

	result, err := dbh.Search(ListActiveGenesQ)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	defer result.Close()

	if result.IsEmpty() {
		logger.Error("No PubMed references found in feature annotations")
		return nil
	}

	for result.Scan() {
		entry := &Gene{}
		if err := result.Read(&entry); err != nil {
			return cli.Exit(
				fmt.Sprintf("error reading query result: %s", err),
				2,
			)
		}
		// Call helper function to process the entry
		params := &processGeneEntryParams{
			entry:  entry,
			dbh:    dbh,
			client: client,
			logger: logger,
		}
		if err := processGeneEntry(params); err != nil {
			// Exit if the helper function encounters an error
			return cli.Exit(err.Error(), 2)
		}
	}

	return nil
}
