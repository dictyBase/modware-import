package main

import (
	"fmt"
	"os"

	"github.com/dictyBase/modware-import/internal/cli/stockcenter"
	"github.com/dictyBase/modware-import/internal/k8s/wait"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "k8s",
		Usage: "Kubernetes utility commands",
		Commands: []*cli.Command{
			{
				Name:   "wait-job",
				Usage:  "Wait for a Kubernetes job to complete, detecting stuck pods early",
				Action: wait.JobAction,
				Flags: append(
					stockcenter.LoggingFlags(),
					&cli.StringFlag{
						Name:     "name",
						Usage:    "Job name to wait for",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "namespace",
						Usage: "Kubernetes namespace",
						Value: "dev",
					},
					&cli.StringFlag{
						Name:  "timeout",
						Usage: "Maximum wait duration (e.g. 60s, 5m)",
						Value: "60s",
					},
					&cli.StringFlag{
						Name:    "kubeconfig",
						Usage:   "Path to kubeconfig file",
						EnvVars: []string{"KUBECONFIG"},
					},
				),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
