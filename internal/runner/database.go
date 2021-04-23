package runner

import (
	"os"
	"path/filepath"

	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/sh"
)

const (
	gitURL   = "https://github.com/dictyBase/modware-import.git"
	cloneDir = "modware-import"
)

// Build builds the binary for modware-import project
func Build() error {
	if err := buildSetup(); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// Cleandb deletes data from arangodb database
func CleanDB(db string) error {
	if err := env.ArangoEnvs(); err != nil {
		return err
	}
	return sh.Run(
		"./importer",
		"--log-level",
		"info",
		"--is-secure",
		"arangodb",
		"delete",
		"-d",
		db,
	)
}

func buildSetup() error {
	modfile := filepath.Join(cloneDir, "go.mod")
	if _, err := os.Stat(modfile); os.IsNotExist(err) {
		if err := cloneSource("develop", cloneDir); err != nil {
			return nil
		}
	}
	return os.Chdir(cloneDir)
}

func cloneSource(branch, dir string) error {
	_, err := git.PlainClone(
		dir,
		false,
		&git.CloneOptions{
			URL:           gitURL,
			SingleBranch:  true,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
		})
	return err
}
