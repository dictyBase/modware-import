package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GwdiCmd is for loading gwdi strains
var GwdiCmd = &cobra.Command{
	Use:   "gwdi",
	Short: "load gwdi strains",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadGwdi,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := SetStrainAPIClient(); err != nil {
			return err
		}
		if err := SetAnnoAPIClient(); err != nil {
			return err
		}
		if err := setReader(viper.GetString("gwdi-input"), regsc.GwdiReader); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	GwdiCmd.Flags().StringP(
		"gwdi-input",
		"i",
		"",
		"csv file with gwdi strains data",
	)
	GwdiCmd.Flags().Bool(
		"gwdi-prune",
		true,
		"prune all gwdi strain records",
	)
	GwdiCmd.Flags().IntP(
		"concurrency",
		"c",
		15,
		"No of concurrent workers",
	)
	viper.BindPFlags(GwdiCmd.Flags())
}
