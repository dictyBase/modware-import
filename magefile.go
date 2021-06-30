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
)

var dbs = []string{"stock", "annotation", "order"}

// CleanAll deletes all data from stock,order and annotation databases
func CleanAllDB() error {
	for _, db := range dbs {
		mg.Deps(mg.F(runner.CleanDB, db))
	}
	return nil
}

// LoadData will clean the collections, load ontologies and then load all the data
func LoadData() error {
	mg.SerialDeps(
		CleanAllDB,
		mg.F(onto.Load, "obojson"),
		stock.LoadStrain,
	)
	return nil
}
