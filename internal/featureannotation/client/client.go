package client

import (
	"fmt"
	"os" // Add this import

	"github.com/arangodb/go-driver" // Add this import
	"github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CliSetup(cltx *cli.Context) error {
	// Setup ArangoDB session
	tls := cltx.Bool("is-secure")
	session, db, err := arangomanager.NewSessionDb(
		&arangomanager.ConnectParams{
			User:     cltx.String("arangodb-user"),
			Pass:     cltx.String("arangodb-pass"),
			Database: cltx.String("arangodb-database"),
			Host:     cltx.String("arangodb-host"),
			Port:     cltx.Int("arangodb-port"),
			Istls:    tls,
		},
	)
	if err != nil {
		return fmt.Errorf("error connecting to arangodb: %s", err)
	}
	registry.SetArangoSession(session)
	registry.SetArangodbConnection(db)

	// Setup gRPC client
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:%s",
			cltx.String("feature-annotation-grpc-host"),
			cltx.String("feature-annotation-grpc-port"),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf(
			"error connecting to feature annotation service: %s",
			err,
		)
	}
	registry.SetFeatureAnnotationAPIClient(
		feature.NewFeatureAnnotationServiceClient(conn),
	)
	return nil
}

func setupArangoDBConnection(
	cltx *cli.Context,
) (*arangomanager.Database, error) {
	session, database, err := arangomanager.NewSessionDb(
		&arangomanager.ConnectParams{
			User:     cltx.String("arangodb-user"),
			Pass:     cltx.String("arangodb-pass"),
			Database: cltx.String("arangodb-database"),
			Host:     cltx.String("arangodb-host"),
			Port:     cltx.Int("arangodb-port"),
			Istls:    cltx.Bool("is-secure"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to arangodb: %w", err)
	}
	registry.SetArangoSession(session)
	registry.SetArangodbConnection(database)
	return database, nil
}

func verifyArangoDBCollection(
	cltx *cli.Context,
	database *arangomanager.Database,
) error {
	collName := cltx.String("collection")
	if _, err := database.Collection(collName); err != nil {
		// Check if the error indicates the collection doesn't exist
		if driver.IsNotFoundGeneral(err) {
			return fmt.Errorf("collection %s does not exist", collName)
		}
		return fmt.Errorf("error checking collection %s: %w", collName, err)
	}
	return nil
}

func verifyAndRegisterCSVFile(cltx *cli.Context) error {
	csvFilePath := cltx.String("csv-file")
	if _, err := os.Stat(csvFilePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("csv file %s does not exist", csvFilePath)
		}
		return fmt.Errorf("error checking csv file %s: %w", csvFilePath, err)
	}
	registry.SetCSVFilePath(csvFilePath)
	return nil
}

// CSVArangodbCliSetup sets up the ArangoDB connection and CSV file for loading.
func CSVArangodbCliSetup(cltx *cli.Context) error {
	database, err := setupArangoDBConnection(cltx)
	if err != nil {
		return err
	}
	if err := verifyArangoDBCollection(cltx, database); err != nil {
		return err
	}
	return verifyAndRegisterCSVFile(cltx)
}
