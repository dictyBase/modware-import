// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	// mage:import stock
	"github.com/dictyBase/modware-import/internal/runner/stock"
	// mage:import annotation
	"github.com/dictyBase/modware-import/internal/runner/annotation"
	// mage:import order
	"github.com/dictyBase/modware-import/internal/runner/order"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Build builds the binary for modware-import project
func Build() error {
	fmt.Println("building ......")
	return sh.Run("go", "build", "-o", "importer", "cmd/import/main.go")
}

// CleanAll deletes all data from stock,order and annotation databases
func CleanAll() error {
	mg.Deps(stock.Clean, annotation.Clean, order.Clean)
	return nil
}
