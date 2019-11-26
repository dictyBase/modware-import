package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/spf13/cobra"
)

// StrainInvCmd is for loading strain inventory data
var StrainInvCmd = &cobra.Command{
	Use:     "strain-inventory",
	Short:   "load strain inventory data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrainInv,
	PreRunE: setInvPreRun,
}

func init() {
	initInvCmd(StrainInvCmd)
}
