package stockcenter

import (
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GoldenbraidCmd is for loading GoldenBraid plasmid CSV data directly
var GoldenbraidCmd = &cobra.Command{
	Use:   "goldenbraid",
	Short: "load GoldenBraid plasmid CSV data directly",
	Long: `Load GoldenBraid plasmid CSV data directly into the stock center API.

This command reads a GoldenBraid CSV file, parses plasmid data, validates it,
and loads it directly to the stock center without requiring intermediate files.`,
	Args:    cobra.NoArgs,
	RunE:    loader.LoadGoldenBraid,
	PreRunE: setGoldenbraidPreRun,
}

func setGoldenbraidPreRun(_ *cobra.Command, _ []string) error {
	// Set up API client (reuse existing strain API client setup)
	return SetStockAPIClient()
}

func init() {
	GoldenbraidCmd.Flags().StringP(
		"input", "i", "",
		"GoldenBraid CSV file path (local path or filename in bucket)",
	)
	GoldenbraidCmd.MarkFlagRequired("input")

	GoldenbraidCmd.Flags().StringP(
		"user-email", "u", "",
		"Email of the user loading the data",
	)
	GoldenbraidCmd.MarkFlagRequired("user-email")

	GoldenbraidCmd.Flags().StringP(
		"plasmid-cvterm", "c", "",
		"Plasmid ontology term (e.g., 'GB vector')",
	)
	GoldenbraidCmd.MarkFlagRequired("plasmid-cvterm")

	viper.BindPFlags(GoldenbraidCmd.Flags())
}
