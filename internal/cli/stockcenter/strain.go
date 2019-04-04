package stockcenter

import (
	"fmt"
	"os"

	"google.golang.org/grpc"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainCmd is for loading stockcenter strain data
var StrainCmd = &cobra.Command{
	Use:     "strain",
	Short:   "load stockcenter strain data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrain,
	PreRunE: setStrainPreRun,
}

func setStrainPreRun(cmd *cobra.Command, args []string) error {
	if err := setStrainAPIClient(); err != nil {
		return err
	}
	if err := setStrainInputReader(); err != nil {
		return err
	}
	return nil
}

func setStrainAPIClient() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", viper.GetString("stock-grpc-host"), viper.GetString("stock-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("error in connecting to stock grpc api server %s", err)
	}
	regsc.SetStockAPIClient(
		stock.NewStockServiceClient(conn),
	)
	return nil
}

func setStrainInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		ar, err := os.Open(viper.GetString("strain-annotator-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strain-annotator-input"), err)
		}
		registry.SetReader(regsc.STRAIN_ANNOTATOR_READER, ar)
		sr, err := os.Open(viper.GetString("input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("input"), err)
		}
		registry.SetReader(regsc.STRAIN_READER, sr)
	case "bucket":
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket-path"),
			viper.GetString("strain-annotator-input"),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strain-annotator-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.STRAIN_ANNOTATOR_READER, ar)
		sr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket-path"),
			viper.GetString("input"),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.STRAIN_READER, sr)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	StrainCmd.Flags().String(
		"stock-grpc-host",
		"stock-api",
		"grpc host address for stock service",
	)
	viper.BindEnv("stock-grpc-host", "STOCK_API_SERVICE_HOST")
	StrainCmd.Flags().String(
		"stock-grpc-port",
		"",
		"grpc port for stock service",
	)
	viper.BindEnv("stock-grpc-port", "STOCK_API_SERVICE_PORT")
	StrainCmd.MarkFlagRequired("stock-grpc-port")
	StrainCmd.Flags().StringP(
		"strain-annotator-input",
		"a",
		"",
		"csv file that provides mapping among strain identifier, annotator and annotation timestamp",
	)
	StrainCmd.Flags().StringP(
		"input",
		"i",
		"",
		"csv file with strain data",
	)
	viper.BindPFlags(StrainCmd.Flags())
}
