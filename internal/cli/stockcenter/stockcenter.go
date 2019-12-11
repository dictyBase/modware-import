package stockcenter

import (
	"fmt"
	"os"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	FOLDER = "folder"
	BUCKET = "bucket"
)

// StockCenterCmd represents the base subcommand for grouping all
// subcommands for loading stockcenter related data
var StockCenterCmd = &cobra.Command{
	Use:   "stockcenter",
	Short: "subcommand for stockcenter data loading",
}

func init() {
	StockCenterCmd.AddCommand(
		OrderCmd,
		PhenoCmd,
		GenoCmd,
		StrainPropCmd,
		StrainCmd,
		StrainCharCmd,
		PlasmidCmd,
		StrainSynCmd,
		StrainInvCmd,
		PlasmidInvCmd,
		ReadFileCmd,
	)
	s3BucketFlags()
	stockAPIFlags()
	annoAPIFlags()
	viper.BindPFlags(StockCenterCmd.PersistentFlags())
}

func s3BucketFlags() {
	StockCenterCmd.PersistentFlags().String(
		"s3-bucket",
		"dictybase",
		"S3 bucket for input files",
	)
	StockCenterCmd.PersistentFlags().String(
		"s3-bucket-path",
		"import/stockcenter",
		"path inside S3 bucket for input files",
	)
}

func stockAPIFlags() {
	StockCenterCmd.PersistentFlags().String(
		"stock-grpc-host",
		"stock-api",
		"grpc host address for stock service",
	)
	viper.BindEnv("stock-grpc-host", "STOCK_API_SERVICE_HOST")
	StockCenterCmd.PersistentFlags().String(
		"stock-grpc-port",
		"",
		"grpc port for stock service",
	)
	viper.BindEnv("stock-grpc-port", "STOCK_API_SERVICE_PORT")
}

func annoAPIFlags() {
	StockCenterCmd.PersistentFlags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	StockCenterCmd.PersistentFlags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
}

func setAnnoAPIClient() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", viper.GetString("annotation-grpc-host"), viper.GetString("annotation-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("error in connecting to annotation grpc api server %s", err)
	}
	regsc.SetAnnotationAPIClient(
		annotation.NewTaggedAnnotationServiceClient(conn),
	)
	return nil
}

func setReader(input, key string) error {
	switch viper.GetString("input-source") {
	case FOLDER:
		r, err := os.Open(input)
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", input, err)
		}
		registry.SetReader(key, r)
	case BUCKET:
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				input,
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				input,
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(key, ar)
	default:
		return fmt.Errorf("error input source %s not supported", input)
	}
	return nil
}
