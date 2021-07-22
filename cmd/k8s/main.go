package main

import (
	"fmt"
	"os"

	"github.com/dictyBase/modware-import/internal/k8s"
	"github.com/dictyBase/modware-import/internal/k8s/app"
	"github.com/kris-nova/naml"
)

func main() {
	cleanerApp, err := app.NewDBCleaner(
		&k8s.AppParams{
			Name:        "db-cleaner",
			Description: "app to flush all data in arangodb databases",
			Namespace:   "dictybase",
			Fragment:    "flush",
		},
		k8s.NewImageSpec("dictybase/modware-import", "develop", "IfNotPresent"),
		"debug",
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	uniportApp, err := app.NewUniprotLoader(
		&k8s.AppParams{
			Name:        "uniprot-loader",
			Description: "app to load uniprot data",
			Namespace:   "dictybase",
			Fragment:    "import",
		},
		k8s.NewImageSpec("dictybase/modware-import", "develop", "IfNotPresent"),
		"debug",
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	naml.Register(cleanerApp)
	naml.Register(uniportApp)
	if err := naml.RunCommandLine(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
