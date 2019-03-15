package stockcenter

import (
	"fmt"
	"os"

	"google.golang.org/grpc"

	"github.com/dictyBase/go-genproto/dictybaseapis/order"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// OrderCmd is for loading stockcenter order data
var OrderCmd = &cobra.Command{
	Use:   "order",
	Short: "load stockcenter order data",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("i am order")
	},
	PreRunE: setPreRun,
}

func setPreRun(cmd *cobra.Command, args []string) error {
	if err := setAPIClient(); err != nil {
		return err
	}
	if err := setInputReader(); err != nil {
		return err
	}
	return nil
}

func setAPIClient() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", viper.GetString("order-grpc-host"), viper.GetString("order-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("error in connecting to order grpc api server %s", err)
	}
	regsc.SetOrderAPIClient(
		order.NewOrderServiceClient(conn),
	)
	return nil
}

func setInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		pr, err := os.Open(viper.GetString("plasmid-map-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("plasmid-map-input"), err)
		}
		registry.SetReader(regsc.PLASMID_ID_MAP_READER, pr)
		or, err := os.Open(viper.GetString("input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("input"), err)
		}
		registry.SetReader(regsc.ORDER_READER, or)
	case "bucket":
		pr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket-path"),
			viper.GetString("plasmid-map-input"),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-map-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.PLASMID_ID_MAP_READER, pr)
		or, err := registry.GetS3Client().GetObject(
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
		registry.SetReader(regsc.ORDER_READER, or)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	OrderCmd.Flags().String(
		"order-grpc-host",
		"order-api",
		"grpc host address for order service",
	)
	viper.BindEnv("order-grpc-host", "ORDER_API_SERVICE_HOST")
	OrderCmd.Flags().String(
		"order-grpc-port",
		"",
		"grpc port for order service",
	)
	viper.BindEnv("order-grpc-port", "ORDER_API_SERVICE_PORT")
	OrderCmd.MarkFlagRequired("order-grpc-port")
	OrderCmd.Flags().StringP(
		"plasmid-map-input",
		"",
		"",
		"csv file that provides mapping between plamid name and identifier in their first two columns",
	)
	OrderCmd.Flags().StringP(
		"input",
		"i",
		"",
		"csv file with order data",
	)
	viper.BindPFlags(OrderCmd.Flags())
}