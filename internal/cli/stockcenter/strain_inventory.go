package stockcenter

import (
	"fmt"
	"os"

	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainInvCmd is for loading strain inventory data
var StrainInvCmd = &cobra.Command{
	Use:   "strain-inventory",
	Short: "load strain inventory data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadStrainInv,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setAnnoAPIClient(); err != nil {
			return err
		}
		if err := setStrainInvReader(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	StrainInvCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	StrainInvCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	StrainInvCmd.Flags().StringP(
		"strain-inventory-input",
		"i",
		"",
		"tsv file with inventory data",
	)
	viper.BindPFlags(StrainInvCmd.Flags())
}

func setStrainInvReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		pr, err := os.Open(viper.GetString("strain-inventory-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("plasmid-inventory-input"), err)
		}
		registry.SetReader(regsc.InvReader, pr)
	case BUCKET:
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("strain-inventory-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-strain-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.InvReader, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}
