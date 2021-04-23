// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	// mage:import stock
	"github.com/dictyBase/modware-import/internal/runner/stock"
	// mage:import annotation
	"github.com/dictyBase/modware-import/internal/runner/annotation"
	// mage:import order
	"github.com/dictyBase/modware-import/internal/runner/order"
)

const (
	gitURL = "https://github.com/dictyBase/modware-import.git"
)

// Build builds the binary for modware-import project
func Build() error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		mg.F(CloneSource, "develop")
	}
	fmt.Println("building ......")
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// CleanAll deletes all data from stock,order and annotation databases
func CleanAll() error {
	mg.Deps(stock.Clean, annotation.Clean, order.Clean)
	return nil
}

// CloneSource get the source code from the default git repository
func CloneSource(branch string) error {
	if len(branch) == 0 {
		branch = "develop"
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = git.PlainClone(
		dir,
		false,
		&git.CloneOptions{
			URL:           gitURL,
			SingleBranch:  true,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
		})
	return err
}
