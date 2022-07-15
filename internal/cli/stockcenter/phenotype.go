package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PhenoCmd is for loading stockcenter phenotype data
var PhenoCmd = &cobra.Command{
	Use:   "phenotype",
	Short: "load stockcenter phenotype data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadPheno,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := SetAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("phenotype-input"), regsc.PhenoReader); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	PhenoCmd.Flags().StringP(
		"phenotype-input",
		"i",
		"",
		"tsv file with strain data",
	)
	viper.BindPFlags(PhenoCmd.Flags())
}
