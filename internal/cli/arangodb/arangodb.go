package arangodb

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ArangodbCmd is the subcommand for managing arangodb database
var ArangodbCmd = &cobra.Command{
	Use:   "arangodb",
	Short: "subcommand for arangodb database management",
}

func init() {
	ArangodbCmd.AddCommand(DeleteCmd)
	arangodbFlags()
	viper.BindPFlags(ArangodbCmd.PersistentFlags())
}

func arangodbFlags() {
	ArangodbCmd.PersistentFlags().StringP(
		"arangodb-pass",
		"p",
		"",
		"arangodb database password",
	)
	ArangodbCmd.PersistentFlags().StringP(
		"arangodb-user",
		"u",
		"",
		"arangodb database user",
	)
	ArangodbCmd.PersistentFlags().StringP(
		"arangodb-host",
		"H",
		"arangodb",
		"arangodb database host",
	)
	ArangodbCmd.PersistentFlags().IntP(
		"arangodb-port",
		"port",
		8529,
		"arangodb database port",
	)
	ArangodbCmd.PersistentFlags().Bool(
		"is-secure",
		true,
		"flag for secured or unsecured arangodb endpoint",
	)
	viper.BindEnv("arangodb-pass", "ARANGODB_PASS")
	viper.BindEnv("arangodb-user", "ARANGODB_USER")
	viper.BindEnv("arangodb-database", "ARANGODB_DATABASE")
	viper.BindEnv("arangodb-host", "ARANGODB_SERVICE_HOST")
	viper.BindEnv("arangodb-port", "ARANGODB_SERVICE_PORT")
}
