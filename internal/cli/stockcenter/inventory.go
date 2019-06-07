package stockcenter

import (
	"fmt"
	"os"

	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InvCmd is for loading stockcenter inventory data
var InvCmd = &cobra.Command{
	Use:     "inventory",
	Short:   "load stockcenter inventory data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadInv,
	PreRunE: setInvPreRun,
}

func setInvPreRun(cmd *cobra.Command, args []string) error {
	if err := setAnnoAPIClient(); err != nil {
		return err
	}
	if err := setInvInputReader(); err != nil {
		return err
	}
	return nil
}

func setInvInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		pr, err := os.Open(viper.GetString("input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("input"), err)
		}
		registry.SetReader(regsc.INV_READER, pr)
	case "bucket":
		ar, err := registry.GetS3Client().GetObject(
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
		registry.SetReader(regsc.INV_READER, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	InvCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	InvCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	InvCmd.MarkFlagRequired("annotation-grpc-port")
	InvCmd.Flags().StringP(
		"input",
		"i",
		"",
		"csv file with inventory data",
	)
	viper.BindPFlags(InvCmd.Flags())
}