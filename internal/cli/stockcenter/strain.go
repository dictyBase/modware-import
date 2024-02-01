package stockcenter

import (
	"fmt"
	"io"
	"os"

	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainCmd is for loading stockcenter strain data
var StrainCmd = &cobra.Command{
	Use:     "strain",
	Short:   "load stockcenter strain data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrain,
	PreRunE: setStrainPreRun,
}

func setStrainPreRun(cmd *cobra.Command, args []string) error {
	if err := SetStrainAPIClient(); err != nil {
		return err
	}
	if err := setStrainInputReader(); err != nil {
		return err
	}
	return nil
}

func init() {
	StrainCmd.Flags().StringP(
		"strain-annotator-input",
		"a",
		"",
		"csv file that provides mapping among strain identifier, annotator and annotation timestamp",
	)
	StrainCmd.Flags().StringP(
		"strain-gene-input",
		"g",
		"",
		"csv file that maps strains to gene identifiers",
	)
	StrainCmd.Flags().StringP(
		"strain-pub-input",
		"p",
		"",
		"csv file that maps strains to publication identifiers",
	)
	StrainCmd.Flags().StringP(
		"strain-input",
		"i",
		"",
		"csv file with strain data",
	)
	viper.BindPFlags(StrainCmd.Flags())
}

func setStrainInputReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		if err := openLocalFiles(); err != nil {
			return err
		}
	case BUCKET:
		if err := getS3Files(); err != nil {
			return err
		}
	default:
		return fmt.Errorf(
			"error input source %s not supported",
			viper.GetString("input-source"),
		)
	}

	return nil
}

func openLocalFiles() error {
	ar, err := openFile(viper.GetString("strain-annotator-input"))
	if err != nil {
		return err
	}
	pr, err := openFile(viper.GetString("strain-pub-input"))
	if err != nil {
		return err
	}
	gr, err := openFile(viper.GetString("strain-gene-input"))
	if err != nil {
		return err
	}
	sr, err := openFile(viper.GetString("strain-input"))
	if err != nil {
		return err
	}
	setRegistryReaders(ar, pr, gr, sr)
	return nil
}

func getS3Files() error {
	ar, err := getS3Object(viper.GetString("strain-annotator-input"))
	if err != nil {
		return err
	}
	pr, err := getS3Object(viper.GetString("strain-pub-input"))
	if err != nil {
		return err
	}
	gr, err := getS3Object(viper.GetString("strain-gene-input"))
	if err != nil {
		return err
	}
	sr, err := getS3Object(viper.GetString("strain-input"))
	if err != nil {
		return err
	}
	setRegistryReaders(ar, pr, gr, sr)
	return nil
}

func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error in opening file %s %s", filePath, err)
	}
	return file, nil
}

func getS3Object(filePath string) (*minio.Object, error) {
	object, err := registry.GetS3Client().GetObject(
		viper.GetString("s3-bucket"),
		fmt.Sprintf("%s/%s", viper.GetString("s3-bucket-path"), filePath),
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error in getting file %s from bucket %s %s",
			filePath,
			viper.GetString("s3-bucket-path"),
			err,
		)
	}
	return object, nil
}

func setRegistryReaders(ar, pr, gr, sr io.Reader) {
	registry.SetReader(regsc.StrainAnnotatorReader, ar)
	registry.SetReader(regsc.StrainPubReader, pr)
	registry.SetReader(regsc.StrainGeneReader, gr)
	registry.SetReader(regsc.StrainReader, sr)
}
