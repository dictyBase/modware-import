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

// PlasmidCmd is for loading stockcenter plasmid data
var PlasmidCmd = &cobra.Command{
	Use:     "plasmid",
	Short:   "load stockcenter plasmid data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadPlasmid,
	PreRunE: setPlasmidPreRun,
}

func plasmidInputMap() map[string]string {
	return map[string]string{
		"plasmid-annotator-input": regsc.PlasmidAnnotatorReader,
		"plasmid-pub-input":       regsc.PlasmidPubReader,
		"plasmid-gene-input":      regsc.PlasmidGeneReader,
		"plasmid-input":           regsc.PlasmidReader,
	}
}

func setPlasmidPreRun(cmd *cobra.Command, args []string) error {
	for _, fn := range []func() error{SetStrainAPIClient, setPlasmidInputReader} {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func setPlasmidInputReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		for k, v := range plasmidInputMap() {
			reader, err := os.Open(viper.GetString(k))
			if err != nil {
				return fmt.Errorf(
					"error in opening file %s %s", viper.GetString(k), err,
				)
			}
			registry.SetReader(v, reader)
		}
	case BUCKET:
		for k, v := range plasmidInputMap() {
			reader, err := registry.GetS3Client().GetObject(
				viper.GetString("s3-bucket"),
				fmt.Sprintf("%s/%s", viper.GetString("s3-bucket-path"), viper.GetString(k)),
				minio.GetObjectOptions{},
			)
			if err != nil {
				return fmt.Errorf(
					"error in getting file %s from bucket %s %s",
					viper.GetString("plasmid-annotator-input"),
					viper.GetString(k), err,
				)
			}
			registry.SetReader(v, reader)
		}
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	PlasmidCmd.Flags().StringP(
		"plasmid-annotator-input", "a", "",
		"csv file that provides mapping among plasmid identifier, annotator and annotation timestamp",
	)
	PlasmidCmd.Flags().StringP(
		"plasmid-gene-input", "g", "",
		"csv file that maps plasmids to gene identifiers",
	)
	PlasmidCmd.Flags().StringP("plasmid-pub-input", "p", "",
		"csv file that maps plasmids to publication identifiers",
	)
	PlasmidCmd.Flags().StringP("plasmid-input", "i", "",
		"csv file with plasmid data",
	)
	viper.BindPFlags(PlasmidCmd.Flags())
}
