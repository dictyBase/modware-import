package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainSynCmd is for loading stockcenter strain synonyms
var StrainSynCmd = &cobra.Command{
	Use:     "strainsyn",
	Short:   "load stockcenter strain synonyms",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrainSynProp,
	PreRunE: setStrainPropPreRun,
}

func init() {
	StrainSynCmd.Flags().StringP(
		"strainprop-input",
		"i",
		"",
		"csv file with strain synonyms data",
	)
	viper.BindPFlags(StrainSynCmd.Flags())
}
