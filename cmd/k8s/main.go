package main

import (
	"fmt"
	"os"

	"github.com/dictyBase/modware-import/internal/k8s"
	"github.com/dictyBase/modware-import/internal/k8s/app"
	"github.com/kris-nova/naml"
)

func main() {
	cleanerApp := app.NewDBCleaner(
		&k8s.AppParams{
			Name:        "db",
			Description: "app to flush all data in arangodb databases",
			Namespace:   "dictybase",
			Fragment:    "cleaner",
		},
		k8s.NewImageSpec("dictybase/modware-import", "develop", "IfNotPresent"),
		"debug",
	)
	uniportApp := app.NewUniprotLoader(
		&k8s.AppParams{
			Name:        "uniprot",
			Description: "app to load uniprot data",
			Namespace:   "dictybase",
			Fragment:    "loader",
		},
		k8s.NewImageSpec("dictybase/modware-import", "develop", "IfNotPresent"),
		"debug",
	)
	naml.Register(cleanerApp)
	naml.Register(uniportApp)
	if err := naml.RunCommandLine(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
