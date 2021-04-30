package ontology

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Refresh gets all obojson files from source github repository
// and stores in S3(minio) for later upload
func Refresh() error {
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	mg.Deps(runner.Build)
	s := runner.TermSpinner("Refreshing objson files ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		runner.LogLevel,
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"ontology",
		"refresh",
	)
}

// Load loads all obograph-json formatted ontologies
func Load() error {
	if err := env.ArangoEnvs(); err != nil {
		return err
	}
	mg.Deps(Refresh)
	s := runner.TermSpinner("loading obojson ontology files ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", runner.Command),
		"--log-level",
		runner.LogLevel,
		"--access-key",
		env.MinioAccessKey(),
		"--secret-key",
		env.MinioSecretKey(),
		"ontology",
		"load",
	)
}
