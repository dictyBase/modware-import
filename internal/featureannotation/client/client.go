package client

import (
	"fmt"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CliSetup(cltx *cli.Context) error {
	conn, err := grpc.NewClient(
		fmt.Sprintf(
			"%s:%s",
			cltx.String("feature-annotation-grpc-host"),
			cltx.String("feature-annotation-grpc-port"),
		),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}...,
	)
	if err != nil {
		return fmt.Errorf(
			"error in connecting to content grpc api server %s",
			err,
		)
	}
	registry.SetFeatureAnnotationAPIClient(
		feature.NewFeatureAnnotationServiceClient(conn),
	)
	return nil
}
