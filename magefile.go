// +build mage

package main

import (
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/dictyBase/modware-import/internal/runner"
	// mage:import stock
	"github.com/dictyBase/modware-import/internal/runner/stock"
	// mage:import annotation
	"github.com/dictyBase/modware-import/internal/runner/annotation"
	// mage:import order
	"github.com/dictyBase/modware-import/internal/runner/order"
)

const (
	cloneDir = "modware-import"
)

// Build builds the binary for modware-import project
func Build() error {
	modfile := filepath.Join(cloneDir, "go.mod")
	if _, err := os.Stat(modfile); os.IsNotExist(err) {
		if err := runner.CloneSource("develop", cloneDir); err != nil {
			return nil
		}
	}
	if err := os.Chdir(cloneDir); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// CleanAll deletes all data from stock,order and annotation databases
func CleanAll() error {
	mg.Deps(stock.Clean, annotation.Clean, order.Clean)
	return nil
}
