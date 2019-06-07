package stockcenter

import (
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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

func setAnnoAPIClient() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", viper.GetString("annotation-grpc-host"), viper.GetString("annotation-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("error in connecting to annotation grpc api server %s", err)
	}
	regsc.SetAnnotationAPIClient(
		annotation.NewTaggedAnnotationServiceClient(conn),
	)
	return nil
}
