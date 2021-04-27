package ontology

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// LoadCmd obojson formmated ontologies to arangodb
var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load an obojson formatted ontologies to arangodb",
	Args:  cobra.NoArgs,
}

func init() {
	loadFlags()
	viper.BindPFlags(LoadCmd.Flags())
}

func loadFlags() {
	LoadCmd.Flags().StringSliceP(
		"obojson",
		"j",
		[]string{""},
		"input ontology files in obograph json format",
	)
	LoadCmd.Flags().String(
		"obograph",
		"obograph",
		"arangodb named graph for managing ontology graph",
	)
	LoadCmd.Flags().String(
		"cv-collection",
		"cv",
		"arangodb collection for storing ontology information",
	)
	LoadCmd.Flags().String(
		"rel-collection",
		"cvterm_relationship",
		"arangodb collection for storing cvterm relationships",
	)
	LoadCmd.Flags().String(
		"term-collection",
		"cvterm",
		"arangodb collection for storing ontoloy terms",
	)
	LoadCmd.Flags().String(
		"s3-bucket",
		"dictybase",
		"S3 bucket for input files",
	)
	LoadCmd.Flags().String(
		"s3-bucket-path",
		"import/ontology",
		"path inside S3 bucket for obojson files",
	)
}
