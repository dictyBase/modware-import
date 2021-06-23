package data

import (
	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/magefile/mage/sh"
)

// Refresh gets all dictybase data files from source github repository
// and stores in S3(minio) server.
func Refresh() error {
	bin, err := runner.LookUp()
	if err != nil {
		return err
	}
	if err := env.MinioEnvs(); err != nil {
		return err
	}
	s := runner.TermSpinner("Refreshing data files ...")
	defer s.Stop()
	s.Start()
	return sh.Run(
		bin, "--log-level", runner.LogLevel,
		"--access-key", env.MinioAccessKey(),
		"--secret-key", env.MinioSecretKey(),
		"--s3-bucket-path", "import",
		"data", "refresh",
		"--group", "migration-import-data",
	)
}
