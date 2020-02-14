package stockcenter

import (
	"context"
	"fmt"

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
	annclient := regs.GetAnnotationAPIClient()
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
		strain, err := createGwdi(stclient, gwdi)
		if err != nil {
			return fmt.Errorf("error in creating new gwdi strain record  %s", err)
		}
		err = createAnno(&createAnnoArgs{
			user:     regs.DEFAULT_USER,
			id:       strain.Data.Id,
			client:   annclient,
			ontology: regs.DICTY_ANNO_ONTOLOGY,
			tag:      genoTag,
			value:    gwdi.Genotype,
		})
		if err != nil {
			return fmt.Errorf("cannot create genotype of gwdi strain %s %s", strain.Data.Id, err)
		}
		for _, char := range gwdi.Characters {
			err = createAnno(&createAnnoArgs{
				user:     regs.DEFAULT_USER,
				id:       strain.Data.Id,
				client:   annclient,
				ontology: strainCharOnto,
				tag:      char,
				value:    val,
			})
			if err != nil {
				return fmt.Errorf(
					"cannot create characteristic %s of gwdi strain %s %s",
					char, strain.Data.Id, err,
				)
			}
		}
		for onto, prop := range gwdi.Properties {
			err = createAnno(&createAnnoArgs{
				user:     regs.DEFAULT_USER,
				id:       strain.Data.Id,
				client:   annclient,
				ontology: onto,
				tag:      prop.Property,
				value:    prop.Value,
			})
			if err != nil {
				return fmt.Errorf(
					"cannot create property %s of gwdi strain %s %s",
					prop.Property, strain.Data.Id, err,
				)
			}
		}
		logger.WithFields(logrus.Fields{
			"event": "create",
			"id":    strain.Data.Id,
		}).Debug("new gwdi strain record")
		count++
	}
	logger.WithFields(logrus.Fields{
		"event": "load",
		"count": count,
	}).Info("all gwdi records")
	return nil
}

func createGwdi(client pb.StockServiceClient, gwdi *stockcenter.GWDIStrain) (*pb.Strain, error) {
	attr := &pb.NewStrainAttributes{
		CreatedBy:    regs.DEFAULT_USER,
		UpdatedBy:    regs.DEFAULT_USER,
		Summary:      gwdi.Summary,
		Genes:        gwdi.Genes,
		Depositor:    gwdi.Depositor,
		Label:        gwdi.Label,
		Species:      gwdi.Species,
		Plasmid:      gwdi.Plasmid,
		Parent:       gwdi.Parent,
		Publications: []string{gwdi.Publication},
		Names:        []string{gwdi.Name},
	}
	return client.CreateStrain(
		context.Background(),
		&pb.NewStrain{
			Data: &pb.NewStrain_Data{
				Type:       "strain",
				Attributes: attr,
			},
		},
	)
}
