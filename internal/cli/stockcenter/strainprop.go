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

// StrainPropCmd is for loading stockcenter strain properties data
var StrainPropCmd = &cobra.Command{
	Use:     "strainprop",
	Short:   "load stockcenter strain properties data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrainProp,
	PreRunE: setStrainPropPreRun,
}

func setStrainPropPreRun(cmd *cobra.Command, args []string) error {
	if err := setAnnoAPIClient(); err != nil {
		return err
	}
	if err := setStrainPropInputReader(); err != nil {
		return err
	}
	return nil
}

func setStrainPropInputReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		pr, err := os.Open(viper.GetString("strainprop-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strainprop-input"), err)
		}
		registry.SetReader(regsc.STRAINPROP_READER, pr)
	case BUCKET:
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("strainprop-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strainprop-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.STRAINPROP_READER, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	StrainPropCmd.Flags().StringP(
		"strainprop-input",
		"i",
		"",
		"csv file with strain property data",
	)
	viper.BindPFlags(StrainPropCmd.Flags())
}
