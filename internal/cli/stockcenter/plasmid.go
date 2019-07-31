package stockcenter

import (
	"fmt"
	"os"

	"google.golang.org/grpc"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PlasmidCmd is for loading stockcenter plasmid data
var PlasmidCmd = &cobra.Command{
	Use:     "plasmid",
	Short:   "load stockcenter plasmid data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadPlasmid,
	PreRunE: setPlasmidPreRun,
}

func setPlasmidPreRun(cmd *cobra.Command, args []string) error {
	if err := setPlasmidAPIClient(); err != nil {
		return err
	}
	if err := setPlasmidInputReader(); err != nil {
		return err
	}
	return nil
}

func setPlasmidAPIClient() error {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", viper.GetString("stock-grpc-host"), viper.GetString("stock-grpc-port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("error in connecting to stock grpc api server %s", err)
	}
	regsc.SetStockAPIClient(
		stock.NewStockServiceClient(conn),
	)
	return nil
}

func setPlasmidInputReader() error {
	switch viper.GetString("input-source") {
	case "folder":
		ar, err := os.Open(viper.GetString("plasmid-annotator-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("plasmid-annotator-input"), err)
		}
		pr, err := os.Open(viper.GetString("plasmid-pub-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("plasmid-pub-input"), err)
		}
		gr, err := os.Open(viper.GetString("plasmid-gene-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("plasmid-gene-input"), err)
		}
		sr, err := os.Open(viper.GetString("plasmid-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("plasmid-input"), err)
		}
		registry.SetReader(regsc.PLASMID_ANNOTATOR_READER, ar)
		registry.SetReader(regsc.PLASMID_PUB_READER, pr)
		registry.SetReader(regsc.PLASMID_GENE_READER, gr)
		registry.SetReader(regsc.PLASMID_READER, sr)
	case "bucket":
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf("%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("plasmid-annotator-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-annotator-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		gr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf("%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("plasmid-gene-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-gene-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		pr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf("%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("plasmid-pub-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-pub-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		sr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf("%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("plasmid-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("plasmid-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		registry.SetReader(regsc.PLASMID_ANNOTATOR_READER, ar)
		registry.SetReader(regsc.PLASMID_PUB_READER, pr)
		registry.SetReader(regsc.PLASMID_GENE_READER, gr)
		registry.SetReader(regsc.PLASMID_READER, sr)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	PlasmidCmd.Flags().String(
		"stock-grpc-host",
		"stock-api",
		"grpc host address for stock service",
	)
	viper.BindEnv("stock-grpc-host", "STOCK_API_SERVICE_HOST")
	PlasmidCmd.Flags().String(
		"stock-grpc-port",
		"",
		"grpc port for stock service",
	)
	viper.BindEnv("stock-grpc-port", "STOCK_API_SERVICE_PORT")
	PlasmidCmd.Flags().StringP(
		"plasmid-annotator-input",
		"a",
		"",
		"csv file that provides mapping among plasmid identifier, annotator and annotation timestamp",
	)
	PlasmidCmd.Flags().StringP(
		"plasmid-gene-input",
		"g",
		"",
		"csv file that maps plasmids to gene identifiers",
	)
	PlasmidCmd.Flags().StringP(
		"plasmid-pub-input",
		"p",
		"",
		"csv file that maps plasmids to publication identifiers",
	)
	PlasmidCmd.Flags().StringP(
		"plasmid-input",
		"i",
		"",
		"csv file with plasmid data",
	)
	viper.BindPFlags(PlasmidCmd.Flags())
}
