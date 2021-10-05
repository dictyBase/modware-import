package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainSynCmd is for loading stockcenter strain synonyms
var StrainSynCmd = &cobra.Command{
	Use:   "strainsyn",
	Short: "load stockcenter strain synonyms",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadStrainSynProp,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := SetAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("strainsyn-input"), regsc.StrainSynReader); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	StrainSynCmd.Flags().StringP(
		"strainsyn-input",
		"i",
		"",
		"csv file with strain synonyms data",
	)
	viper.BindPFlags(StrainSynCmd.Flags())
}
