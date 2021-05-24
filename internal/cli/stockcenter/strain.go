package stockcenter

import (
	"encoding/csv"
	"fmt"
	"os"

	"google.golang.org/grpc"

	"github.com/dictyBase/go-genproto/dictybaseapis/stock"
	loader "github.com/dictyBase/modware-import/internal/load/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/minio/minio-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StrainCmd is for loading stockcenter strain data
var StrainCmd = &cobra.Command{
	Use:     "strain",
	Short:   "load stockcenter strain data",
	Args:    cobra.NoArgs,
	RunE:    loader.LoadStrain,
	PreRunE: setStrainPreRun,
}

func setStrainPreRun(cmd *cobra.Command, args []string) error {
	if err := setStrainAPIClient(); err != nil {
		return err
	}
	if err := setStrainInputReader(); err != nil {
		return err
	}
	return nil
}

func setStrainAPIClient() error {
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

func setStrainInputReader() error {
	switch viper.GetString("input-source") {
	case FOLDER:
		ar, err := os.Open(viper.GetString("strain-annotator-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strain-annotator-input"), err)
		}
		pr, err := os.Open(viper.GetString("strain-pub-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strain-pub-input"), err)
		}
		gr, err := os.Open(viper.GetString("strain-gene-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strain-gene-input"), err)
		}
		sr, err := os.Open(viper.GetString("strain-input"))
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", viper.GetString("strain-input"), err)
		}
		registry.SetReader(regsc.STRAIN_ANNOTATOR_READER, ar)
		registry.SetReader(regsc.STRAIN_PUB_READER, pr)
		registry.SetReader(regsc.STRAIN_GENE_READER, gr)
		registry.SetReader(regsc.STRAIN_READER, sr)
	case BUCKET:
		objpath := fmt.Sprintf("%s/%s",
			viper.GetString("s3-bucket-path"), viper.GetString("strain-annotator-input"),
		)
		registry.GetLogger().Infof("fetch object %s", objpath)
		ar, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			objpath,
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strain-annotator-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		csvr := csv.NewReader(ar)
		records, err := csvr.ReadAll()
		if err != nil {
			return fmt.Errorf("cannot read all data directly %s", err)
		}
		registry.GetLogger().Infof("read %d lines from annotator", len(records))
		gr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("strain-gene-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strain-gene-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		pr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("strain-pub-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strain-pub-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		sr, err := registry.GetS3Client().GetObject(
			viper.GetString("s3-bucket"),
			fmt.Sprintf(
				"%s/%s",
				viper.GetString("s3-bucket-path"),
				viper.GetString("strain-input"),
			),
			minio.GetObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf(
				"error in getting file %s from bucket %s %s",
				viper.GetString("strain-input"),
				viper.GetString("s3-bucket-path"),
				err,
			)
		}
		//registry.SetReader(regsc.STRAIN_ANNOTATOR_READER, ar)
		registry.SetReader("spongebob", ar)
		registry.SetReader(regsc.STRAIN_PUB_READER, pr)
		registry.SetReader(regsc.STRAIN_GENE_READER, gr)
		registry.SetReader(regsc.STRAIN_READER, sr)
	default:
		return fmt.Errorf("error input source %s not supported", viper.GetString("input-source"))
	}
	return nil
}

func init() {
	StrainCmd.Flags().StringP(
		"strain-annotator-input",
		"a",
		"",
		"csv file that provides mapping among strain identifier, annotator and annotation timestamp",
	)
	StrainCmd.Flags().StringP(
		"strain-gene-input",
		"g",
		"",
		"csv file that maps strains to gene identifiers",
	)
	StrainCmd.Flags().StringP(
		"strain-pub-input",
		"p",
		"",
		"csv file that maps strains to publication identifiers",
	)
	StrainCmd.Flags().StringP(
		"strain-input",
		"i",
		"",
		"csv file with strain data",
	)
	viper.BindPFlags(StrainCmd.Flags())
}
