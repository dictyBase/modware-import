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

// PlasmidInvCmd is for loading plasmid inventory data
var PlasmidInvCmd = &cobra.Command{
	Use:   "plasmid-inventory",
	Short: "load plasmid inventory data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadPlasmidInv,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setAnnoAPIClient(); err != nil {
			return err
		}
		if err := setPlasmidInvReader(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	PlasmidInvCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	PlasmidInvCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	PlasmidInvCmd.Flags().StringP(
		"plasmid-inventory-input",
		"i",
		"",
		"tsv file with inventory data",
	)
	viper.BindPFlags(PlasmidInvCmd.Flags())
}

func setPlasmidInvReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		pr, err := os.Open(viper.GetString("plasmid-inventory-input"))
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
				viper.GetString("plasmid-inventory-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-inventory-input"),
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
