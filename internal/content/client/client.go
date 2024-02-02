package client

import (
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	"github.com/dictyBase/modware-import/internal/datasource/s3"
	regsc "github.com/dictyBase/modware-import/internal/registry"
	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CliSetup(cltx *cli.Context) error {
	if err := SetS3Client(cltx); err != nil {
		return cli.Exit(err.Error(), 2)
	}
	if err := SetContentAPIClient(cltx); err != nil {
		return cli.Exit(err.Error(), 2)
	}

	return nil
}

func SetS3Client(cltx *cli.Context) error {
	client, err := s3.NewCliS3Client(cltx)
	if err != nil {
		return fmt.Errorf("error in getting instance of s3 client %s", err)
	}
	regsc.SetS3Client(client)

	return nil
}

func SetContentAPIClient(cltx *cli.Context) error {
	conn, err := grpc.Dial(
		fmt.Sprintf(
			"%s:%s",
			cltx.String("content-grpc-host"),
			cltx.String("content-grpc-port"),
		),
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		}...,
	)
	if err != nil {
		return fmt.Errorf(
			"error in connecting to content grpc api server %s",
			err,
		)
	}
	registry.SetContentAPIClient(content.NewContentServiceClient(conn))
	return nil
}
