package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dictyBase/aphgrpc"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadStrain(cmd *cobra.Command, args []string) error {
	pl, err := source.NewStockPubLookup(
		registry.GetReader(regs.StrainPubReader),
	)
	if err != nil {
		return fmt.Errorf("error in opening publication source %s", err)
	}
	gl, err := source.NewStockGeneLookp(
		registry.GetReader(regs.StrainGeneReader),
	)
	if err != nil {
		return fmt.Errorf("error in opening gene source %s", err)
	}
	al, err := source.NewStockAnnotatorLookup(
		registry.GetReader(regs.StrainAnnotatorReader),
	)
	if err != nil {
		return fmt.Errorf("error in opening annotation source %s", err)
	}
	sr := source.NewCsvStrainReader(
		registry.GetReader(regs.StrainReader),
		al,
		pl,
		gl,
	)
	logger := registry.GetLogger()
	client := regs.GetStockAPIClient()
	count := 0
	for sr.Next() {
		strain, err := sr.Value()
		if err != nil {
			logger.Errorf(
				"error in reading strain value from datasource %s",
				err,
			)
			continue
		}
		if len(strain.User) == 0 {
			logger.Errorf(
				"strain %s does not have a user assignment, skipping the load",
				strain.Id,
			)
			continue
		}
		_, err = client.GetStrain(
			context.Background(),
			&pb.StockId{Id: strain.Id},
		)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// create new strain entry
				attr := &pb.ExistingStrainAttributes{
					CreatedAt: aphgrpc.TimestampProto(strain.CreatedOn),
					UpdatedAt: aphgrpc.TimestampProto(strain.UpdatedOn),
					CreatedBy: strain.User,
					UpdatedBy: strain.User,
					Summary:   strain.Summary,
					Species:   strain.Species,
					Label:     strain.Descriptor,
				}
				if len(strain.Publications) > 0 {
					attr.Publications = strain.Publications
				} else {
					logger.Warnf("strain %s has no publication entry", strain.Id)
				}
				if len(strain.Genes) > 0 {
					attr.Genes = strain.Genes
				}
				nstr, err := client.LoadStrain(
					context.Background(),
					&pb.ExistingStrain{
						Data: &pb.ExistingStrain_Data{
							Type:       "strain",
							Id:         strain.Id,
							Attributes: attr,
						},
					})
				if err != nil {
					return fmt.Errorf(
						"error in creating strain %s %s",
						strain.Id,
						err,
					)
				}
				logger.Debugf("created strain %s", nstr.Data.Id)
				count += 1
				continue
			}
			return fmt.Errorf("error in finding strain %s %s", strain.Id, err)
		}
		// update strains
		attr := &pb.StrainUpdateAttributes{
			UpdatedBy: strain.User,
			Summary:   strain.Summary,
			Label:     strain.Descriptor,
		}
		if len(strain.Publications) > 0 {
			attr.Publications = strain.Publications
		} else {
			logger.Warnf("strain %s has no publication entry", strain.Id)
		}
		if len(strain.Genes) > 0 {
			attr.Genes = strain.Genes
		}
		ustr, err := client.UpdateStrain(
			context.Background(),
			&pb.StrainUpdate{
				Data: &pb.StrainUpdate_Data{
					Type:       "strain",
					Id:         strain.Id,
					Attributes: attr,
				},
			})
		if err != nil {
			return fmt.Errorf("error in updating strain %s %s", strain.Id, err)
		}
		logger.Debugf("updated strain %s", ustr.Data.Id)
		count += 1
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "annotations",
			"stock": "strains",
			"event": "load",
			"count": count,
		}).Infof("loaded strain annotations")
	return nil
}
