// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	// mage:import
	_ "github.com/dictyBase/modware-import/internal/runner"
	// mage:import stock
	"github.com/dictyBase/modware-import/internal/runner/stock"
	// mage:import annotation
	"github.com/dictyBase/modware-import/internal/runner/annotation"
	// mage:import order
	"github.com/dictyBase/modware-import/internal/runner/order"
)

// CleanAll deletes all data from stock,order and annotation databases
func CleanAll() error {
	mg.Deps(stock.Clean, annotation.Clean, order.Clean)
	return nil
}
