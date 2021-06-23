package data

import "github.com/spf13/cobra"

const (
	GroupTag = "data-file-group"
)

// DataCmd manages data files
var DataCmd = &cobra.Command{
	Use:   "data",
	Short: "subcommand for managing data files",
}

func init() {
	DataCmd.AddCommand(RefreshCmd)
}
