package ontology

import (
	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Refresh gets all obojson formatted files from source github repository
// and stores in S3(minio) for later upload.
// The ontology group,binary path and git reference(tag,branch or commit id)
// have to be supplied.
func Refresh(bin, group, gitref string) error {
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	runner.ConsoleLog("Refreshing objson files ...")
	defer runner.ConsoleLog("Done refreshing objson files ...")
	return sh.Run(
		bin, "--log-level", runner.LogLevel,
		"--access-key", env.MinioAccessKey(),
		"--secret-key", env.MinioSecretKey(),
		"--s3-bucket-path", "import/obograph-json",
		"ontology", "refresh",
		"--group", group,
		"--branch", gitref,
	)
}

// Load loads all obograph-json formatted ontologies
// The ontology group and git branch have to be supplied.
func Load(group, gitref string) error {
	bin, err := runner.LookUp()
	if err != nil {
		return err
	}
	if err := env.ArangoEnvs(); err != nil {
		return err
	}
	if err := env.ArangoDBName(); err != nil {
		return err
	}
	mg.Deps(mg.F(Refresh, bin, group, gitref))
	runner.ConsoleLog("loading obojson ontology files ...")
	defer runner.ConsoleLog("Done loading obojosn files ....")
	return sh.Run(
		bin, "ontology", "--log-level", runner.LogLevel,
		"--access-key", env.MinioAccessKey(),
		"--secret-key", env.MinioSecretKey(),
		"--s3-bucket-path", "import/obograph-json",
		"--is-secure", "load",
		"--group", group,
	)
}
