package content

import (
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/content"
	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
