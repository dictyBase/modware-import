package client

import (
	"context"
	"fmt"

	"github.com/dictyBase/modware-import/internal/registry"
	rds "github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
)

func SetRedisClient(cltx *cli.Context) error {
	client := rds.NewClient(&rds.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			cltx.String("redis-master-service-host"),
			cltx.String("redis-master-service-port"),
		),
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("error pinging redis: %w", err)
	}
	registry.SetRedisClient(client)
	return nil
}
