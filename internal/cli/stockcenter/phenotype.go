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

// PhenoCmd is for loading stockcenter phenotype data
var PhenoCmd = &cobra.Command{
	Use:     "phenotype",
	Short:   "load stockcenter phenotype data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadPheno,
	PreRunE: setPhenoPreRun,
}

func setPhenoPreRun(cmd *cobra.Command, args []string) error {
	if err := setAnnoAPIClient(); err != nil {
		return err
	}
	if err := setPhenoInputReader(); err != nil {
		return err
	}
	return nil
}

func setPhenoInputReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		pr, err := os.Open(viper.GetString("phenotype-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("phenotype-input"), err)
		}
		registry.SetReader(regsc.PHENO_READER, pr)
	case BUCKET:
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket-path"),
			viper.GetString("phenotype-input"),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("phenotype-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.PHENO_READER, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	PhenoCmd.Flags().String(
		"annotation-grpc-host",
		"annotation-api",
		"grpc host address for annotation service",
	)
	viper.BindEnv("annotation-grpc-host", "ANNOTATION_API_SERVICE_HOST")
	PhenoCmd.Flags().String(
		"annotation-grpc-port",
		"",
		"grpc port for annotation service",
	)
	viper.BindEnv("annotation-grpc-port", "ANNOTATION_API_SERVICE_PORT")
	PhenoCmd.Flags().StringP(
		"phenotype-input",
		"i",
		"",
		"csv file with strain data",
	)
	viper.BindPFlags(PhenoCmd.Flags())
}
