package ontology

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Refresh gets all obojson formatted files from source github repository
// and stores in S3(minio) for later upload. The ontology group and git branch have to be supplied.
func Refresh(group, branch string) error {
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	mg.Deps(mg.F(runner.BuildBranch, branch))
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
		"--group",
		group,
	)
}

// Load loads all obograph-json formatted ontologies
// The ontology group and git branch have to be supplied.
func Load(group, branch string) error {
	if err := env.ArangoEnvs(); err != nil {
		return err
	}
	mg.Deps(mg.F(Refresh, group, branch))
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
		"--group",
		group,
	)
}
