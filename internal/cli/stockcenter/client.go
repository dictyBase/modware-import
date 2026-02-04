package stockcenter

import (
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/datasource/s3"
	"github.com/dictyBase/modware-import/internal/fputil"
	registry "github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
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

// SetStockClient creates stock gRPC connection and registers client
func SetStockClient(cltx *cli.Context) IOE.IOEither[error, *grpc.ClientConn] {
	return F.Pipe2(
		GRPCClientConfig{
			Host:        cltx.String("stock-grpc-host"),
			Port:        cltx.String("stock-grpc-port"),
			ServiceName: "stock grpc server",
		},
		createGRPCConnection,
		IOE.Map[error](func(conn *grpc.ClientConn) *grpc.ClientConn {
			regsc.SetStockAPIClient(stock.NewStockServiceClient(conn))
			return conn
		}),
	)
}

// SetStockClientWrapper wraps SetStockClient for use as a Before hook
// Returns error directly instead of IOEither for urfave/cli compatibility
func SetStockClientWrapper(cltx *cli.Context) error {
	return F.Pipe2(
		SetStockClient(cltx),
		fputil.ToEither[error, *grpc.ClientConn],
		E.Fold(
			F.Identity[error],
			F.Constant1[*grpc.ClientConn, error](nil),
		),
	)
}

// SetAnnotationClient creates annotation gRPC connection and registers client
func SetAnnotationClient(
	cltx *cli.Context,
) IOE.IOEither[error, *grpc.ClientConn] {
	return F.Pipe2(
		GRPCClientConfig{
			Host:        cltx.String("annotation-grpc-host"),
			Port:        cltx.String("annotation-grpc-port"),
			ServiceName: "annotation grpc server",
		},
		createGRPCConnection,
		IOE.Map[error](func(conn *grpc.ClientConn) *grpc.ClientConn {
			regsc.SetAnnotationAPIClient(
				annotation.NewTaggedAnnotationServiceClient(conn),
			)
			return conn
		}),
	)
}

// SetClients sets up both stock and annotation clients
func SetClients(cltx *cli.Context) error {
	return F.Pipe2(
		IOE.SequenceArraySeq([]IOE.IOEither[error, *grpc.ClientConn]{
			SetStockClient(cltx),
			SetAnnotationClient(cltx),
		}),
		fputil.ToEither[error, []*grpc.ClientConn],
		E.Fold(
			F.Identity[error],
			F.Constant1[[]*grpc.ClientConn, error](nil),
		),
	)
}

// SetS3Client initializes the S3 client
func SetS3Client(cltx *cli.Context) IOE.IOEither[error, *minio.Client] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*minio.Client, error) {
			return s3.NewCliS3Client(cltx)
		}),
		IOE.MapLeft[*minio.Client](func(err error) error {
			return fmt.Errorf("error in getting instance of s3 client %w", err)
		}),
		IOE.Map[error](func(client *minio.Client) *minio.Client {
			registry.SetS3Client(client)
			return client
		}),
	)
}

// SetStockAndS3Clients sets up both stock client and S3 client
// Uses SequenceArraySeq for independent operations that should both succeed
func SetStockAndS3Clients(cltx *cli.Context) error {
	return F.Pipe2(
		IOE.SequenceArraySeq([]IOE.IOEither[error, any]{
			IOE.Map[error](func(conn *grpc.ClientConn) any { return conn })(SetStockClient(cltx)),
			IOE.Map[error](func(client *minio.Client) any { return client })(SetS3Client(cltx)),
		}),
		fputil.ToEither[error, []any],
		E.Fold(
			F.Identity[error],
			F.Constant1[[]any, error](nil),
		),
	)
}
