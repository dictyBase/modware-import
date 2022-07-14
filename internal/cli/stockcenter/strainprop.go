package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainPropCmd is for loading stockcenter strain properties data
var StrainPropCmd = &cobra.Command{
	Use:   "strainprop",
	Short: "load stockcenter strain properties data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadStrainProp,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := SetAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("strainprop-input"), regsc.StrainpropReader); err != nil {
			return err
		}
		return nil
	},
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
