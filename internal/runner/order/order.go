package order

import (
	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Build builds the binary for modware-import project
func Build() error {
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// Clean deletes all order data from arangodb database
func Clean() error {
	if err := runner.ArangoEnvs(); err != nil {
		return err
	}
	mg.Deps(Build)
	return sh.Run(
		"./importer",
		"--log-level",
		"info",
		"--is-secure",
		"arangodb",
		"delete",
		"-d",
		"order",
	)
}
