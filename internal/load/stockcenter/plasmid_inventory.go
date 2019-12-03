package stockcenter

import (
	"context"
	"fmt"
	"time"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func LoadPlasmidInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewTsvPlasmidInventoryReader(registry.GetReader(regs.INV_READER))
	logger := registry.GetLogger()
	invMap, err := cacheInvByPlasmidId(ir, logger)
	if err != nil {
		return err
	}
	client := regs.GetAnnotationAPIClient()
	invCount := 0
	for id, invSlice := range invMap {
		found := true
		gc, err := getInventory(id, client, regs.PLASMID_INV_ONTO)
		if err != nil {
			if grpc.Code(err) != codes.NotFound { // error in lookup
				return err
			}
			found = false
			logger.WithFields(
				logrus.Fields{
					"type":  "inventory",
					"stock": "plasmid",
					"event": "get",
					"id":    id,
				}).Debugf("no inventories")
		}
		if found {
			logger.WithFields(
				logrus.Fields{
					"type":  "inventory",
					"stock": "plasmid",
					"event": "get",
					"id":    id,
				}).Debugf("retrieved inventories")
			if err := delExistingInventory(id, client, gc); err != nil {
				return err
			}
			logger.WithFields(
				logrus.Fields{
					"type":  "inventory",
					"stock": "plasmid",
					"event": "delete",
					"id":    id,
				}).Debugf("deleted inventories")
		}
		if err := createPlasmidInventory(id, client, invSlice, found); err != nil {
			return err
		}
		logger.WithFields(
			logrus.Fields{
				"type":  "inventory",
				"stock": "plasmid",
				"event": "create",
				"id":    id,
				"count": len(invSlice),
			}).Debugf("created inventories")
		invCount = invCount + len(invSlice)
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": "plasmids",
			"event": "load",
			"count": invCount,
		}).Infof("loaded inventories")
	return nil
}

func cacheInvByPlasmidId(ir stockcenter.PlasmidInventoryReader, logger *logrus.Entry) (map[string][]*stockcenter.PlasmidInventory, error) {
	invMap := make(map[string][]*stockcenter.PlasmidInventory)
	for ir.Next() {
		inv, err := ir.Value()
		if err != nil {
			return invMap, fmt.Errorf(
				"error in loading inventory for plasmid %s",
				err,
			)
		}
		if len(inv.PhysicalLocation) == 0 {
			logger.WithFields(
				logrus.Fields{
					"type":   "inventory",
					"stock":  "plasmid",
					"event":  "skip record",
					"output": inv.RecordLine,
				}).Warnf("skipped the record")
			continue
		}
		if invSlice, ok := invMap[inv.PlasmidId]; ok {
			invMap[inv.PlasmidId] = append(invSlice, inv)
			continue
		}
		invMap[inv.PlasmidId] = []*stockcenter.PlasmidInventory{inv}
	}
	return invMap, nil
}

func createPlasmidInventory(id string, client pb.TaggedAnnotationServiceClient, invSlice []*stockcenter.PlasmidInventory, found bool) error {
	for i, inv := range invSlice {
		var ids []string
		m := map[string]string{
			regs.INV_LOCATION_TAG:     inv.PhysicalLocation,
			regs.INV_STORED_AS_TAG:    inv.StoredAs,
			regs.INV_PRIV_COMMENT_TAG: inv.PrivateComment,
			regs.INV_OBTAINED_AS_TAG:  inv.ObtainedAs,
			regs.PLASMID_INV_TAG:      regs.INV_EXIST_VALUE,
		}
		if !inv.StoredOn.IsZero() {
			m[regs.INV_STORAGE_DATE_TAG] = inv.StoredOn.Format(time.RFC3339Nano)
		}
	INNER:
		for tag, value := range m {
			if len(value) == 0 {
				continue INNER
			}
			anno, err := createAnnoWithRank(client, tag, inv.PlasmidId, regs.PLASMID_INV_ONTO, value, i)
			if err != nil {
				return err
			}
			ids = append(ids, anno.Data.Id)
		}
		_, err := client.CreateAnnotationGroup(context.Background(), &pb.AnnotationIdList{Ids: ids})
		if err != nil {
			return err
		}
	}
	// create presence of inventory annotation
	if !found {
		_, err := createAnno(
			client, regs.PLASMID_INV_ONTO, invSlice[0].PlasmidId,
			regs.PLASMID_INV_ONTO, regs.INV_EXIST_VALUE,
		)
		if err != nil {
			return err
		}

	}
	return nil
}
