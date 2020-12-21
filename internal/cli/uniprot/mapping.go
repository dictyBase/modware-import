package uniprot

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UniprotMappingCmd is for loading mapping of Uniprot IDs to Gene IDs.
var UniprotMappingCmd = &cobra.Command{
	Use:   "mapping",
	Short: "load uniprot id mappings",
	Args:  cobra.NoArgs,
	// RunE:  loader.LoadUniprotMappings,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := setRedisClient(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	UniprotMappingCmd.Flags().StringP(
		"uniprot-url",
		"u",
		"https://www.uniprot.org/uniprot/?query=taxonomy:44689&columns=id,database(dictyBase),genes(PREFERRED)&format=tab",
		"uniprot endpoint",
	)
	viper.BindPFlags(UniprotMappingCmd.Flags())
}
