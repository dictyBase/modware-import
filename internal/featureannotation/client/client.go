package client

import (
	"fmt"

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
