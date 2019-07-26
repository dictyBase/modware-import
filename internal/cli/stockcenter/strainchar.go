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

// StrainCharCmd is for loading stockcenter strain characteristics data
var StrainCharCmd = &cobra.Command{
	Use:     "strainchar",
	Short:   "load stockcenter strain characteristics data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrainChar,
	PreRunE: setStrainCharPreRun,
}

func setStrainCharPreRun(cmd *cobra.Command, args []string) error {
	if err := setAnnoAPIClient(); err != nil {
		return err
	}
	if err := setStrainCharInputReader(); err != nil {
		return err
	}
	return nil
}

func setStrainCharInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		pr, err := os.Open(viper.GetString("strainchar-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strainchar-input"), err)
		}
		registry.SetReader(regsc.STRAINCHAR_READER, pr)
	case "bucket":
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket-path"),
			viper.GetString("strainchar-input"),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strainchar-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.STRAINCHAR_READER, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	StrainCharCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	StrainCharCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	StrainCharCmd.Flags().StringP(
		"strainchar-input",
		"i",
		"",
		"csv file with strain characteristics data",
	)
	viper.BindPFlags(StrainCharCmd.Flags())
}
