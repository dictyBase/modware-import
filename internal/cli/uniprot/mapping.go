package uniprot

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	loader "github.com/dictyBase/modware-import/internal/load/uniprot"
)

// UniprotMappingCmd is for loading mapping of Uniprot IDs to Gene IDs.
var UniprotMappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "load uniprot id mappings",
	Args:  cobra.NoArgs,
	RunE:  loader.LoadUniprotMappings,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return setRedisClient()
	},
}

var uniprotURL = fmt.Sprintf(
	"%s/?%s&%s&%s",
	"https://www.uniprot.org/uniprot",
	"query=taxonomy:44689",
	"columns=id,database(dictyBase),genes(PREFERRED)",
	"format=tab",
)

func init() {
	UniprotMappingCmd.Flags().StringP(
		"uniprot-url",
		"u",
		uniprotURL,
		"uniprot endpoint",
	)
	viper.BindPFlags(UniprotMappingCmd.Flags())
}
