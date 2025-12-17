package main

import (
	"fmt"
	"log"
	"os"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/fputil"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClientConfig holds configuration for creating a gRPC client
type GRPCClientConfig struct {
	Host        string
	Port        string
	ServiceName string // For error messages
}

// createGRPCConnection creates a gRPC client connection
// Separates connection creation from client registration
func createGRPCConnection(
	config GRPCClientConfig,
) IOE.IOEither[error, *grpc.ClientConn] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*grpc.ClientConn, error) {
			return grpc.NewClient(
				fmt.Sprintf("%s:%s", config.Host, config.Port),
				grpc.WithTransportCredentials(
					insecure.NewCredentials(),
				),
			)
		}),
		IOE.MapLeft[*grpc.ClientConn](func(err error) error {
			return fmt.Errorf(
				"failed to connect to %s: %w",
				config.ServiceName,
				err,
			)
		}))
}

// registerStockClient registers stock API client and returns connection
func registerStockClient(
	conn *grpc.ClientConn,
) IOE.IOEither[error, *grpc.ClientConn] {
	return IOE.TryCatchError(func() (*grpc.ClientConn, error) {
		regsc.SetStockAPIClient(stock.NewStockServiceClient(conn))
		return conn, nil
	})
}

// registerAnnotationClient registers annotation API client and returns connection
func registerAnnotationClient(
	conn *grpc.ClientConn,
) IOE.IOEither[error, *grpc.ClientConn] {
	return IOE.TryCatchError(func() (*grpc.ClientConn, error) {
		regsc.SetAnnotationAPIClient(
			annotation.NewTaggedAnnotationServiceClient(conn),
		)
		return conn, nil
	})
}

// setStockClient creates stock gRPC connection and registers client
func setStockClient(c *cli.Context) IOE.IOEither[error, *grpc.ClientConn] {
	return F.Pipe2(
		GRPCClientConfig{
			Host:        c.String("stock-grpc-host"),
			Port:        c.String("stock-grpc-port"),
			ServiceName: "stock grpc server",
		},
		createGRPCConnection,
		IOE.Chain(registerStockClient),
	)
}

// setAnnotationClient creates annotation gRPC connection and registers client
func setAnnotationClient(c *cli.Context) IOE.IOEither[error, *grpc.ClientConn] {
	return F.Pipe2(
		GRPCClientConfig{
			Host:        c.String("annotation-grpc-host"),
			Port:        c.String("annotation-grpc-port"),
			ServiceName: "annotation grpc server",
		},
		createGRPCConnection,
		IOE.Chain(registerAnnotationClient),
	)
}

func setClients(c *cli.Context) error {
	return F.Pipe3(
		setStockClient(c),
		IOE.Chain(
			func(_ *grpc.ClientConn) IOE.IOEither[error, *grpc.ClientConn] {
				return setAnnotationClient(c)
			},
		),
		fputil.ToEither[error, *grpc.ClientConn],
		E.Fold(
			F.Identity[error],
			F.Constant1[*grpc.ClientConn, error](nil),
		),
	)
}

func main() {
	app := &cli.App{
		Name:  "goldenbraid",
		Usage: "GoldenBraid data loading tools",
		Commands: []*cli.Command{
			{
				Name:   "inventory",
				Usage:  "Load GoldenBraid inventory",
				Action: stockcenter.LoadGoldenBraidInventory,
				Before: setClients,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "Input CSV file path",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "input-source",
						Usage: "Source of input file (folder or bucket)",
						Value: "bucket",
					},
					&cli.StringFlag{
						Name:  "s3-bucket",
						Usage: "S3 bucket for input files",
						Value: "dictybase",
					},
					&cli.StringFlag{
						Name:  "s3-bucket-path",
						Usage: "Path inside S3 bucket for input files",
						Value: "import/stockcenter",
					},
					&cli.StringFlag{
						Name:  "stock-grpc-host",
						Usage: "gRPC host address for stock service",
						Value: "stock-api",
					},
					&cli.StringFlag{
						Name:  "stock-grpc-port",
						Usage: "gRPC port for stock service",
						Value: "9560",
					},
					&cli.StringFlag{
						Name:  "annotation-grpc-host",
						Usage: "gRPC host address for annotation service",
						Value: "annotation-api",
					},
					&cli.StringFlag{
						Name:  "annotation-grpc-port",
						Usage: "gRPC port for annotation service",
						Value: "9560",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
