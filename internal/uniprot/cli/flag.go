package cli

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/config"
	"github.com/urfave/cli/v2"
)

const (
	uniprotBaseURL  = "https://rest.uniprot.org/uniprotkb/search?query="
	uniprotFields   = "id,xref_dictybase"
	uniprotFormat   = "json"
	uniprotPageSize = 500
)

var uniprotURL = fmt.Sprintf(
	"%sorganism_id:%d&fields=%s&format=%s&size=%d",
	uniprotBaseURL,
	config.UniProtPort,
	uniprotFields,
	uniprotFormat,
	uniprotPageSize,
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
