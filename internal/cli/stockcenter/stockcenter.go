package stockcenter

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StockCenterCmd represents the base subcommand for grouping all
// subcommands for loading stockcenter related data
var StockCenterCmd = &cobra.Command{
	Use:   "stockcenter",
	Short: "subcommand for stockcenter data loading",
}

func init() {
	StockCenterCmd.AddCommand(OrderCmd)
	StockCenterCmd.PersistentFlags().String(
		"s3-bucket-path",
		"dictybase/import/stockcenter",
		"S3 bucket path where all stockcenter data will be kept",
	)
	viper.BindPFlags(StockCenterCmd.PersistentFlags())
}
