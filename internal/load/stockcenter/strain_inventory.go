package stockcenter

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadStrainInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewCsvStrainInventoryReader(registry.GetReader(regs.INV_READER))
	logger := registry.GetLogger()
	invMap, err := cacheInvByStrainId(ir, logger)
	if err != nil {
		return err
	}
	client := regs.GetAnnotationAPIClient()
	invCount := 0
	for id, invSlice := range invMap {
		gc, err := getInventory(id, client, logger)
		if err != nil {
			return err
		}
		if err := delExistingInventory(id, client, gc, logger); err != nil {
			return err
		}
		if err := createStrainInventory(id, client, invSlice, logger); err != nil {
			return err
		}
		invCount++
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": "strains",
			"event": "load",
			"count": invCount,
		}).Infof("loaded inventories")
	return nil
}

func getInventory(id string, client pb.TaggedAnnotationServiceClient, logger *logrus.Entry) (*pb.TaggedAnnotationGroupCollection, error) {
	gc, err := client.ListAnnotationGroups(
		context.Background(),
		&pb.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;tag==%s;ontology==%s",
				id, regs.INV_LOCATION_TAG, regs.STRAIN_INV_ONTO,
			),
		})
	if err != nil {
		if grpc.Code(err) != codes.NotFound { // error in lookup
			return gc, err
		}
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": "strains",
			"event": "get",
			"id":    id,
		}).Debugf("retrieved inventories")

	return gc, nil
}

func delExistingInventory(id string, client pb.TaggedAnnotationServiceClient, gc *pb.TaggedAnnotationGroupCollection, logger *logrus.Entry) error {
	for _, gcd := range gc.Data {
		// remove annotations group
		_, err := client.DeleteAnnotationGroup(
			context.Background(),
			&pb.GroupEntryId{GroupId: gcd.Group.GroupId},
		)
		if err != nil {
			return err
		}
		// remove all annotations
		for _, gd := range gcd.Group.Data {
			_, err := client.DeleteAnnotation(
				context.Background(),
				&pb.DeleteAnnotationRequest{Id: gd.Id, Purge: true},
			)
			if err != nil {
				return err
			}
		}
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": "strains",
			"event": "delete",
			"id":    id,
		}).Debugf("deleted inventories")
	return nil
}

func cacheInvByStrainId(ir stockcenter.StrainInventoryReader, logger *logrus.Entry) (map[string][]*stockcenter.StrainInventory, error) {
	invMap := make(map[string][]*stockcenter.StrainInventory)
	for ir.Next() {
		inv, err := ir.Value()
		if err != nil {
			return invMap, fmt.Errorf(
				"error in loading inventory for strain %s",
				err,
			)
		}
		if len(inv.PhysicalLocation) == 0 || len(inv.VialColor) == 0 {
			logger.WithFields(
				logrus.Fields{
					"type":   "inventory",
					"stock":  "strains",
					"event":  "skip record",
					"output": inv.RecordLine,
				}).Warnf("skipped the record")
			continue
		}
		if v, ok := invMap[inv.StrainId]; ok {
			invMap[inv.StrainId] = append(v, inv)
		} else {
			invMap[inv.StrainId] = []*stockcenter.StrainInventory{inv}
		}
	}
	return invMap, nil
}

func createStrainInventory(id string, client pb.TaggedAnnotationServiceClient, invSlice []*stockcenter.StrainInventory, logger *logrus.Entry) error {
	for _, inv := range invSlice {
		var ids []string
		m := map[string]string{
			regs.INV_LOCATION_TAG:     inv.PhysicalLocation,
			regs.INV_STORED_AS_TAG:    inv.StoredAs,
			regs.INV_VIAL_COUNT_TAG:   inv.VialsCount,
			regs.INV_VIAL_COLOR_TAG:   inv.VialColor,
			regs.INV_PRIV_COMMENT_TAG: inv.PrivateComment,
			regs.INV_PUB_COMMENT_TAG:  inv.PublicComment,
			regs.INV_STORAGE_DATE_TAG: inv.StoredOn.Format(time.RFC3339Nano),
		}
		for t, v := range m {
			if len(v) == 0 {
				continue
			}
			t, err := createAnno(client, t, inv.StrainId, regs.STRAIN_INV_ONTO, v)
			if err != nil {
				return err
			}
			ids = append(ids, t.Data.Id)
		}
		_, err := client.CreateAnnotationGroup(context.Background(), &pb.AnnotationIdList{Ids: ids})
		if err != nil {
			return err
		}
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": "strains",
			"event": "create",
			"id":    id,
		}).Debugf("created inventories")
	return nil
}
