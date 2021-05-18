package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/dictyBase/modware-import/internal/runner/env"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/sh"
)

const (
	gitURL   = "https://github.com/dictyBase/modware-import.git"
	cloneDir = "modware-import"
	branch   = "develop"
	Command  = "importer"
	LogLevel = "info"
)

func TermSpinnerWithPrefixColor(prefix, color string) *spinner.Spinner {
	s := spinner.New(
		spinner.CharSets[33],
		300*time.Millisecond,
	)
	_ = s.Color("bgHiBlack", "bold", color)
	s.Prefix = fmt.Sprintf("%s  ", prefix)
	return s
}

func TermSpinner(prefix string) *spinner.Spinner {
	return TermSpinnerWithPrefixColor(prefix, "fgHiGreen")
}

// A standalone builder, it builds the binary
// after checking out the source code from the given branch
func BuildBranch(branch string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := buildSetup(cloneDir, branch); err != nil {
		return err
	}
	return buildAndClean(currDir, cloneDir)
}

// Builds the modware-import binary. It is intended to run
// from source folder
func Build() error {
	s := TermSpinner("building modware-import binary ...")
	defer s.Stop()
	s.Start()
	return sh.Run("go", "build", "-o", Command, "cmd/import/main.go")
}

// Another standalone builder, it builds the binary after
// checking out the source code from develop branch
func MagicBuild() error {
	if err := buildSetup(cloneDir, branch); err != nil {
		return err
	}
	s := TermSpinner("building modware-import binary ...")
	defer s.Stop()
	s.Start()
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// Cleandb deletes data from arangodb database
func CleanDB(db string) error {
	if err := env.ArangoEnvs(); err != nil {
		return err
	}
	s := TermSpinnerWithPrefixColor(
		fmt.Sprintf("cleaning database %s ...", db),
		"fgHiMagenta",
	)
	defer s.Stop()
	s.Start()
	return sh.Run(
		fmt.Sprintf("./%s", Command),
		"arangodb",
		"--log-level",
		LogLevel,
		"--is-secure",
		"delete",
		"--database",
		db,
	)
}

func buildAndClean(curr, dir string) error {
	if err := Build(); err != nil {
		return err
	}
	if err := os.Chdir(curr); err != nil {
		return err
	}
	dst := filepath.Join(curr, Command)
	src := filepath.Join(dir, Command)
	return sh.Copy(dst, src)
}

func buildSetup(dir, branch string) error {
	modfile := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(modfile); os.IsNotExist(err) {
		if err := cloneSource(branch, dir); err != nil {
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
