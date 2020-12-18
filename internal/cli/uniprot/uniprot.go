package uniprot

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// UniprotCmd represents the base subcommand for grouping all
// subcommands for mapping Uniprot IDs
var UniprotCmd = &cobra.Command{
	Use:   "uniprot",
	Short: "subcommand for uniprot id mapping",
}

func init() {
	UniprotCmd.AddCommand(UniprotMappingCmd)
	redisFlags()
	viper.BindPFlags(UniprotCmd.PersistentFlags())
}

func redisFlags() {
	UniprotCmd.PersistentFlags().String(
		"redis-master-service-host",
		"",
		"grpc host address for redis service",
	)
	viper.BindEnv("redis-master-service-host", "REDIS_MASTER_SERVICE_HOST")
	UniprotCmd.PersistentFlags().String(
		"redis-master-service-port",
		"",
		"grpc port for redis service",
	)
	viper.BindEnv("redis-master-service-port", "REDIS_MASTER_SERVICE_PORT")
}
