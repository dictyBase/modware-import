package stockcenter

import (
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadGwdi(cmd *cobra.Command, args []string) error {
	gr := stockcenter.NewGWDIStrainReader(registry.GetReader(regs.GWDI_READER))
	stclient := regs.GetStockAPIClient()
	logger := registry.GetLogger().WithFields(logrus.Fields{
		"type":  "gwdi",
		"stock": "strain",
	})
	count := 0
	for gr.Next() {
		gwdi, err := gr.Value()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"event": "read",
			}).Errorf("gwdi datasource error %s", err)
			continue
		}
		attr := &pb.NewStrainAttributes{
			CreatedBy: regs.DEFAULT_USER,
			UpdatedBy: regs.DEFAULT_USER,
			Summary:   gwdi.Summary,
		}
	}
	return nil
}
