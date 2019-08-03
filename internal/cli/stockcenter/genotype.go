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

// GenoCmd is for loading stockcenter genotype data
var GenoCmd = &cobra.Command{
	Use:     "genotype",
	Short:   "load stockcenter genotype data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadGeno,
	PreRunE: setGenoPreRun,
}

func setGenoPreRun(cmd *cobra.Command, args []string) error {
	if err := setAnnoAPIClient(); err != nil {
		return err
	}
	if err := setGenoInputReader(); err != nil {
		return err
	}
	return nil
}

func setGenoInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		pr, err := os.Open(viper.GetString("genotype-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("genotype-input"), err)
		}
		registry.SetReader(regsc.GENO_READER, pr)
	case "bucket":
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("genotype-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("genotype-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.GENO_READER, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	GenoCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	GenoCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	GenoCmd.Flags().StringP(
		"genotype-input",
		"i",
		"",
		"csv file with genotype data",
	)
	viper.BindPFlags(GenoCmd.Flags())
}
