// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	// mage:import
	"github.com/dictyBase/modware-import/internal/runner"
	// mage:import stock
	_ "github.com/dictyBase/modware-import/internal/runner/stock"
	// mage:import ontology
	_ "github.com/dictyBase/modware-import/internal/runner/ontology"
)

var dbs = []string{"stock", "annotation", "order"}

// CleanAll deletes all data from stock,order and annotation databases
func CleanAllDB(branch string) error {
	mg.Deps(mg.F(runner.BuildBranch, branch))
	for _, db := range dbs {
		mg.Deps(mg.F(runner.CleanDB(db)))
	}
	return nil
}
