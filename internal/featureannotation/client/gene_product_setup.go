package client

import (
	"fmt"

	"github.com/dictyBase/arangomanager"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
)

// GeneProductCliSetup sets up the CLI environment for gene product updater
func GeneProductCliSetup(cltx *cli.Context) error {
	// Setup main database connection and gRPC client
	if err := CliSetup(cltx); err != nil {
		return err
	}

	// Setup legacy database connection
	connParams := &arangomanager.ConnectParams{
		User:     cltx.String("arangodb-user"),
		Pass:     cltx.String("arangodb-pass"),
		Database: cltx.String("legacy-database"),
		Host:     cltx.String("arangodb-host"),
		Port:     cltx.Int("arangodb-port"),
		Istls:    cltx.Bool("is-secure"),
	}

	// Create a new session and database connection for the legacy DB
	// Note: NewSessionDb returns a session and a database object.
	legacySess, legacyDbh, err := arangomanager.NewSessionDb(connParams)
	if err != nil {
		return fmt.Errorf(
			"error connecting to legacy database %s: %w",
			cltx.String("legacy-database"),
			err,
		)
	}

	registry.SetLegacyArangodbConnection(legacyDbh)
	registry.SetLegacyArangoSession(legacySess)
	return nil
}

// CleanupGeneProductResources cleans up resources used by gene product operations
func CleanupGeneProductResources() error {
	// Note: ArangoDB sessions are managed by the arangomanager library
	// and don't require explicit cleanup in this implementation
	return nil
}
