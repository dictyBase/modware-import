package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var uniprotURL = fmt.Sprintf(
	"%sorganism_id:%d&fields=%s&format=%s&size=%d,`",
	"https://rest.uniprot.org/uniprotkb/search?query=",
	44689, "id,xref_dictybase", "json", 500,
)

func UniprotFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "uniprot-url",
			Aliases: []string{"u"},
			Usage:   "URL for fetching Uniprot data",
			Value:   uniprotURL,
		},
		&cli.StringFlag{
			Name:    "redis-service-host",
			Aliases: []string{"s"},
			Usage:   "Redis service host address",
			EnvVars: []string{"REDIS_SERVICE_HOST"},
		},
		&cli.StringFlag{
			Name:    "redis-service-port",
			Aliases: []string{"p"},
			Usage:   "Redis service port",
			Value:   "6379",
			EnvVars: []string{"REDIS_SERVICE_PORT"},
		},
	}
}
