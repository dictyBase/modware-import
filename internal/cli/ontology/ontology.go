package ontology

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	GroupTag = "ontology-group"
)

// OntologyCmd is the subcommand for managing ontology
var OntologyCmd = &cobra.Command{
	Use:   "ontology",
	Short: "subcommand for ontology management",
}

func init() {
	OntologyCmd.AddCommand(LoadCmd, RefreshCmd)
	ontologyStorageFlags()
	viper.BindPFlags(OntologyCmd.PersistentFlags())
}

func ontologyStorageFlags() {
	OntologyCmd.PersistentFlags().String(
		"arangodb-database",
		"",
		"arangodb database name",
	)
	OntologyCmd.PersistentFlags().StringP(
		"arangodb-user",
		"u",
		"",
		"arangodb database user",
	)
	OntologyCmd.PersistentFlags().StringP(
		"arangodb-pass",
		"p",
		"",
		"arangodb database password",
	)
	OntologyCmd.PersistentFlags().StringP(
		"arangodb-host",
		"H",
		"arangodb",
		"arangodb database host",
	)
	OntologyCmd.PersistentFlags().Int(
		"arangodb-port",
		8529,
		"arangodb database port",
	)
	OntologyCmd.PersistentFlags().Bool(
		"is-secure",
		false,
		"flag for secured or unsecured arangodb endpoint",
	)
	viper.BindEnv("arangodb-pass", "ARANGODB_PASS")
	viper.BindEnv("arangodb-user", "ARANGODB_USER")
	viper.BindEnv("arangodb-host", "ARANGODB_SERVICE_HOST")
	viper.BindEnv("arangodb-port", "ARANGODB_SERVICE_PORT")
	viper.BindEnv("arangodb-database", "ARANGODB_DATABASE")
}
