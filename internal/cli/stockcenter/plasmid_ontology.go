package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultPropertyName = "keyword"
)

var PlasmidOntologyCmd = &cobra.Command{
	Use:   "plasmid-ontology",
	Short: "associate plasmids with ontology keywords",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadPlasmidOntology,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return SetStockAPIClient()
	},
}

func init() {
	PlasmidOntologyCmd.Flags().StringP(
		"input",
		"i",
		"",
		"tsv file with plasmid properties",
	)
	PlasmidOntologyCmd.MarkFlagRequired("input")
	PlasmidOntologyCmd.Flags().String(
		"property",
		defaultPropertyName,
		"property label to filter plasmid keyword rows",
	)
	viper.BindPFlags(PlasmidOntologyCmd.Flags())
}
