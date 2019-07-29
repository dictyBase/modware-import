package stockcenter

import (
	"fmt"
	"os"

	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ReadFileCmd is for reading any data file
var ReadFileCmd = &cobra.Command{
	Use:     "readfile",
	Short:   "read data from any file",
	Args:    cobra.NoArgs,
	RunE:    LoadReadFile,
	PreRunE: setReadFilePreRun,
}

func setReadFilePreRun(cmd *cobra.Command, args []string) error {
	if err := setReadFileInputReader(); err != nil {
		return err
	}
	return nil
}

func setReadFileInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		pr, err := os.Open(viper.GetString("readfile-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("readfile-input"), err)
		}
		registry.SetReader(regsc.READFILE_READER, pr)
	case "bucket":
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket-path"),
			viper.GetString("readfile-input"),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("readfile-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.READFILE_READER, ar)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	ReadFileCmd.Flags().StringP(
		"readfile-input",
		"i",
		"",
		"input file with data",
	)
	viper.BindPFlags(ReadFileCmd.Flags())
}
