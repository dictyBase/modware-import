package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/spf13/cobra"
)

// PlasmidInvCmd is for loading plasmid inventory data
var PlasmidInvCmd = &cobra.Command{
	Use:     "plasmid-inventory",
	Short:   "load plasmid inventory data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadPlasmidInv,
	PreRunE: setInvPreRun,
}

func init() {
	initInvCmd(PlasmidInvCmd)
}
