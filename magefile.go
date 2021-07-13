// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	// mage:import
	"github.com/dictyBase/modware-import/internal/runner"
	// mage:import stock
	_ "github.com/dictyBase/modware-import/internal/runner/stock"
	stock "github.com/dictyBase/modware-import/internal/runner/stock"

	// mage:import ontology
	_ "github.com/dictyBase/modware-import/internal/runner/ontology"
	onto "github.com/dictyBase/modware-import/internal/runner/ontology"

	// mage:import data
	_ "github.com/dictyBase/modware-import/internal/runner/data"
	data "github.com/dictyBase/modware-import/internal/runner/data"
)

var dbs = []string{"stock", "annotation", "order"}

// CleanAll deletes all data from stock,order and annotation databases
func CleanAllDB() error {
	for _, db := range dbs {
		mg.Deps(mg.F(runner.CleanDB, db))
	}
	return nil
}

// LoadData will do the following steps...
// i)  Refresh ontology from github repository to s3 storage.
//	   Requires a git ref.
// ii) Load ontology in database
// iii)  Refresh data from github repository to s3 storage.
//	   Requires a git ref.
// iv) Load data
func LoadData(gitref string) error {
	mg.SerialDeps(
		mg.F(onto.Load, "obojson", gitref),
		mg.F(data.Refresh, gitref),
		stock.LoadStrain,
		stock.LoadPlasmid,
	)
	return nil
}
