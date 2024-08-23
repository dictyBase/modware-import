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
	"%sorganism_id:%d&fields=%s&format=%s&size=%d,`",
	"https://rest.uniprot.org/uniprotkb/search?query=",
	44689, "id,xref_dictybase", "json", 500,
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
