package ontology

import (
	"github.com/dictyBase/modware-import/internal/cli/arangodb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// OntologyCmd is the subcommand for managing ontology
var OntologyCmd = &cobra.Command{
	Use:   "ontology",
	Short: "subcommand for ontology management",
}

func init() {
	arangodb.ArangodbFlags()
	extraArangodbFlags()
	viper.BindPFlags(arangodb.ArangodbCmd.PersistentFlags())
	viper.BindPFlags(OntologyCmd.PersistentFlags())
}

func extraArangodbFlags() {
	OntologyCmd.PersistentFlags().String(
		"arangodb-database",
		"",
		"arangodb database name",
	)
	viper.BindEnv("arangodb-database", "ARANGODB_DATABASE")
}
