package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	source "github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadPlasmid(cmd *cobra.Command, args []string) error {
	al, err := initAnnotatorLookup()
	if err != nil {
		return err
	}
	pl, err := initPubLookup()
	if err != nil {
		return err
	}
	gl, err := initGeneLookup()
	if err != nil {
		return err
	}

	sr := source.NewCsvPlasmidReader(
		registry.GetReader(regs.PlasmidReader),
		al, pl, gl,
	)

	client := regs.GetStockAPIClient()
	logger := registry.GetLogger()

	count := 0
	for sr.Next() {
		if err := processNextPlasmid(sr, client, logger); err != nil {
			return err
		}
		count++
	}
	logLoadStats(logger, count)

	return nil
}

func initAnnotatorLookup() (source.StockAnnotatorLookup, error) {
	al, err := source.NewStockAnnotatorLookup(
		registry.GetReader(regs.PlasmidAnnotatorReader),
	)
	if err != nil {
		return nil, fmt.Errorf("error in opening annotation source %s", err)
	}
	return al, nil
}

func initPubLookup() (source.StockPubLookup, error) {
	pl, err := source.NewStockPubLookup(
		registry.GetReader(regs.PlasmidPubReader),
	)
	if err != nil {
		return nil, fmt.Errorf("error in opening publication source %s", err)
	}
	return pl, nil
}

func initGeneLookup() (source.StockGeneLookup, error) {
	gl, err := source.NewStockGeneLookp(
		registry.GetReader(regs.PlasmidGeneReader),
	)
	if err != nil {
		return nil, fmt.Errorf("error in opening gene source %s", err)
	}
	return gl, nil
}

func processNextPlasmid(
	sr source.PlasmidReader,
	client pb.StockServiceClient,
	logger *logrus.Entry,
) error {
	plasmid, err := sr.Value()
	if err != nil {
		logger.Errorf("error in reading plasmid value from datasource %s", err)
		return nil // Continuing to next record, not returning an error to stop the loop
	}
	if len(plasmid.User) == 0 {
		logger.Errorf(
			"plasmid %s does not have a user assignment, skipping the load",
			plasmid.Id,
		)
		return nil // Continuing to next record
	}

	exists, err := plasmidExists(client, plasmid.Id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return createPlasmid(client, logger, plasmid)
		}
		return fmt.Errorf("error in finding plasmid %s %s", plasmid.Id, err)
	}

	if exists {
		return updatePlasmid(client, logger, plasmid)
	}

	return nil
}

func plasmidExists(client pb.StockServiceClient, id string) (bool, error) {
	_, err := client.GetPlasmid(context.Background(), &pb.StockId{Id: id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func createPlasmid(
	client pb.StockServiceClient,
	logger *logrus.Entry,
	plasmid *source.Plasmid,
) error {
	attr := populateExistingPlasmidAttributes(logger, plasmid)
	_, err := client.LoadPlasmid(
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
	logger.Debugf("created plasmid %s", plasmid.Id)
	return nil
}

func updatePlasmid(
	client pb.StockServiceClient,
	logger *logrus.Entry,
	plasmid *source.Plasmid,
) error {
	attr := populatePlasmidUpdateAttributes(logger, plasmid)
	_, err := client.UpdatePlasmid(
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
	logger.Debugf("updated plasmid %s", plasmid.Id)
	return nil
}

func populateExistingPlasmidAttributes(
	logger *logrus.Entry,
	plasmid *source.Plasmid,
) *pb.ExistingPlasmidAttributes {
	attr := &pb.ExistingPlasmidAttributes{
		CreatedAt: TimestampProto(plasmid.CreatedOn),
		UpdatedAt: TimestampProto(plasmid.UpdatedOn),
		CreatedBy: plasmid.User,
		Summary:   plasmid.Summary,
		Name:      plasmid.Name,
	}
	checkPublicationsAndGenes(logger, plasmid, attr)
	return attr
}

func populatePlasmidUpdateAttributes(
	logger *logrus.Entry,
	plasmid *source.Plasmid,
) *pb.PlasmidUpdateAttributes {
	attr := &pb.PlasmidUpdateAttributes{
		UpdatedBy: plasmid.User,
		Summary:   plasmid.Summary,
		Name:      plasmid.Name,
	}
	checkPublicationsAndGenes(logger, plasmid, attr)
	return attr
}

func checkPublicationsAndGenes(
	logger *logrus.Entry,
	plasmid *source.Plasmid,
	attr interface{},
) {
	switch a := attr.(type) {
	case *pb.ExistingPlasmidAttributes:
		if len(plasmid.Publications) > 0 {
			a.Publications = plasmid.Publications
		} else {
			logger.Warnf("plasmid %s has no publication entry", plasmid.Id)
		}
		if len(plasmid.Genes) > 0 {
			a.Genes = plasmid.Genes
		}
	case *pb.PlasmidUpdateAttributes:
		if len(plasmid.Publications) > 0 {
			a.Publications = plasmid.Publications
		} else {
			logger.Warnf("plasmid %s has no publication entry", plasmid.Id)
		}
		if len(plasmid.Genes) > 0 {
			a.Genes = plasmid.Genes
		}
	}
}

func logLoadStats(logger *logrus.Entry, count int) {
	logger.WithFields(
		logrus.Fields{
			"type":  "annotations",
			"stock": "plasmid",
			"event": "load",
			"count": count,
		}).Infof("loaded plasmid annotations")
}
