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
// i)   Refresh and load ontology
// ii)  Refresh and load data in s3 storage
// iii) Load data
func LoadData() error {
	mg.SerialDeps(
		mg.F(onto.Load, "obojson"),
		data.Refresh,
		stock.LoadStrain,
		stock.LoadPlasmid,
	)
	return nil
}
