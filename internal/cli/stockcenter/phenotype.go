package stockcenter

import (
	"fmt"
	"os"

	"google.golang.org/grpc"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	minio "github.com/minio/minio-go"
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
	if err := setPhenoAPIClient(); err != nil {
		return err
	}
	if err := setPhenoInputReader(); err != nil {
		return err
	}
	return nil
}

func setPhenoAPIClient() error {
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

func setPhenoInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		pr, err := os.Open(viper.GetString("input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("input"), err)
		}
		registry.SetReader(regsc.PHENO_READER, pr)
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
	PhenoCmd.MarkFlagRequired("annotation-grpc-port")
	PhenoCmd.Flags().StringP(
		"input",
		"i",
		"",
		"csv file with strain data",
	)
	viper.BindPFlags(StrainCmd.Flags())
}
