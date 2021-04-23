package annotation

import (
	"github.com/dictyBase/modware-import/internal/runner"
	"github.com/magefile/mage/mg"
)

// Clean deletes all order data from arangodb database
func Clean() error {
	mg.SerialDeps(
		runner.Build,
		mg.F(runner.CleanDb, "annotation"),
	)
	return nil
}
