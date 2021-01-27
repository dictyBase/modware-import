package uniprot

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/registry"
	r "github.com/go-redis/redis/v7"
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

func setRedisClient() error {
	client := r.NewClient(&r.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			viper.GetString("redis-master-service-host"),
			viper.GetString("redis-master-service-port"),
		),
	})
	err := client.Ping().Err()
	if err != nil {
		return fmt.Errorf("error pinging redis %s", err)
	}
	registry.SetRedisClient(client)
	return nil
}
