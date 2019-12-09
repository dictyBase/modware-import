package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GenoCmd is for loading stockcenter genotype data
var GenoCmd = &cobra.Command{
	Use:   "genotype",
	Short: "load stockcenter genotype data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadGeno,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("genotype-input"), regsc.GENO_READER); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	GenoCmd.Flags().StringP(
		"genotype-input",
		"i",
		"",
		"csv file with genotype data",
	)
	viper.BindPFlags(GenoCmd.Flags())
}
