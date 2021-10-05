package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainInvCmd is for loading strain inventory data
var StrainInvCmd = &cobra.Command{
	Use:   "strain-inventory",
	Short: "load strain inventory data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadStrainInv,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := SetAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("strain-inventory-input"), regsc.InvReader); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	StrainInvCmd.Flags().StringP(
		"strain-inventory-input",
		"i",
		"",
		"tsv file with inventory data",
	)
	viper.BindPFlags(StrainInvCmd.Flags())
}
