package stockcenter

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
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

// SetStockClient creates stock gRPC connection and registers client
func SetStockClient(c *cli.Context) IOE.IOEither[error, *grpc.ClientConn] {
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

// SetAnnotationClient creates annotation gRPC connection and registers client
func SetAnnotationClient(c *cli.Context) IOE.IOEither[error, *grpc.ClientConn] {
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

// SetClients sets up both stock and annotation clients
func SetClients(c *cli.Context) error {
	return F.Pipe2(
		IOE.SequenceArraySeq([]IOE.IOEither[error, *grpc.ClientConn]{
			SetStockClient(c),
			SetAnnotationClient(c),
		}),
		fputil.ToEither[error, []*grpc.ClientConn],
		E.Fold(
			F.Identity[error],
			F.Constant1[[]*grpc.ClientConn, error](nil),
		),
	)
}
