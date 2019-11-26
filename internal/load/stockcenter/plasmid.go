package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
)

func LoadPlasmid(cmd *cobra.Command, args []string) error {
	al, err := source.NewStockAnnotatorLookup(registry.GetReader(regs.PLASMID_ANNOTATOR_READER))
	if err != nil {
		return fmt.Errorf("error in opening annotation source %s", err)
	}
	pl, err := source.NewStockPubLookup(registry.GetReader(regs.PLASMID_PUB_READER))
	if err != nil {
		return fmt.Errorf("error in opening publication source %s", err)
	}
	gl, err := source.NewStockGeneLookp(registry.GetReader(regs.PLASMID_GENE_READER))
	if err != nil {
		return fmt.Errorf("error in opening gene source %s", err)
	}
	sr := source.NewCsvPlasmidReader(
		registry.GetReader(regs.PLASMID_READER),
		al, pl, gl,
	)
	logger := registry.GetLogger()
	client := regs.GetStockAPIClient()
	for sr.Next() {
		plasmid, err := sr.Value()
		if err != nil {
			logger.Errorf("error in reading plasmid value from datasource %s", err)
			continue
		}
		if len(plasmid.User) == 0 {
			logger.Errorf("plasmid %s does not have a user assignment, skipping the load", plasmid.Id)
			continue
		}
		_, err = client.GetPlasmid(context.Background(), &pb.StockId{Id: plasmid.Id})
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				// create new strain entry
				attr := &pb.ExistingPlasmidAttributes{
					CreatedAt: TimestampProto(plasmid.CreatedOn),
					UpdatedAt: TimestampProto(plasmid.UpdatedOn),
					CreatedBy: plasmid.User,
					Summary:   plasmid.Summary,
					Name:      plasmid.Name,
				}
				if len(plasmid.Publications) > 0 {
					attr.Publications = plasmid.Publications
				} else {
					logger.Warnf("plasmid %s has no publication entry", plasmid.Id)
				}
				if len(plasmid.Genes) > 0 {
					attr.Genes = plasmid.Genes
				}
				npl, err := client.LoadPlasmid(
					context.Background(),
					&pb.ExistingPlasmid{
						Data: &pb.ExistingPlasmid_Data{
							Type:       "plasmid",
							Id:         plasmid.Id,
							Attributes: attr,
						},
					})
				if err != nil {
					return fmt.Errorf("error in creating plasmid %s %s", plasmid.Id, err)
				}
				logger.Infof("created plasmid %s", npl.Data.Id)
				continue
			}
			return fmt.Errorf("error in finding plasmid %s %s", plasmid.Id, err)
		}
		// update plasmid
		attr := &pb.PlasmidUpdateAttributes{
			UpdatedBy: plasmid.User,
			Summary:   plasmid.Summary,
			Name:      plasmid.Name,
		}
		if len(plasmid.Publications) > 0 {
			attr.Publications = plasmid.Publications
		} else {
			logger.Warnf("plasmid %s has no publication entry", plasmid.Id)
		}
		if len(plasmid.Genes) > 0 {
			attr.Genes = plasmid.Genes
		}
		upl, err := client.UpdatePlasmid(
			context.Background(),
			&pb.PlasmidUpdate{
				Data: &pb.PlasmidUpdate_Data{
					Type:       "plasmid",
					Id:         plasmid.Id,
					Attributes: attr,
				},
			})
		if err != nil {
			return fmt.Errorf("error in updating plasmid %s %s", plasmid.Id, err)
		}
		logger.Infof("updated plasmid %s", upl.Data.Id)
	}
	return nil
}
