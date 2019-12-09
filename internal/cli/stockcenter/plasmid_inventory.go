package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"

	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PlasmidInvCmd is for loading plasmid inventory data
var PlasmidInvCmd = &cobra.Command{
	Use:   "plasmid-inventory",
	Short: "load plasmid inventory data",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadPlasmidInv,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("plasmid-inventory-input"), regsc.InvReader); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	PlasmidInvCmd.Flags().StringP(
		"plasmid-inventory-input",
		"i",
		"",
		"tsv file with inventory data",
	)
	viper.BindPFlags(PlasmidInvCmd.Flags())
}
