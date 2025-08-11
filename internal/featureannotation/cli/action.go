package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

func LoadCSVToArangodb(cltx *cli.Context) error {
	logger := registry.GetLogger()
	result := collection.Pipe4(
		cltx,
		setupPipeline,
		setupFileProcessing,
		validateHeaders,
		processCSVRecords,
	)

	if result.File != nil {
		defer result.File.Close()
	}

	// Check for errors
	if result.Error != nil {
		logger.Errorf("Pipeline failed: %s", result.Error)
		return cli.Exit(
			fmt.Sprintf("Pipeline execution failed: %s", result.Error.Error()),
			2,
		)
	}

	// Log success
	logger.Infof(
		"Successfully finished processing CSV for collection %s. Total documents updated: %d",
		result.Setup.CollectionName,
		result.UpdateCount,
	)

	return nil
}
