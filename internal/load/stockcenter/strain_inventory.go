package stockcenter

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadStrainInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewTsvStrainInventoryReader(registry.GetReader(regs.INV_READER))
	logger := registry.GetLogger()
	invMap, err := cacheInvByStrainId(ir, logger)
	if err != nil {
		return err
	}
	client := regs.GetAnnotationAPIClient()
	invCount := 0
	for id, invSlice := range invMap {
		found := true
		gc, err := getInventory(id, client, regs.STRAIN_INV_ONTO)
		if err != nil {
			if grpc.Code(err) != codes.NotFound { // error in lookup
				return err
			}
			found = false
			logger.WithFields(
				logrus.Fields{
					"type":  "inventory",
					"stock": "strain",
					"event": "get",
					"id":    id,
				}).Debugf("no inventories")
		}
		if found { // remove if inventory exists
			logger.WithFields(
				logrus.Fields{
					"type":  "inventory",
					"stock": "strain",
					"event": "get",
					"id":    id,
				}).Debugf("retrieved inventories")
			if err := delExistingInventory(id, client, gc); err != nil {
				return err
			}
			logger.WithFields(
				logrus.Fields{
					"type":  "inventory",
					"stock": "strain",
					"event": "delete",
					"id":    id,
				}).Debugf("deleted inventories")
		}
		err = createStrainInventory(&strainInvArgs{
			id:       id,
			client:   client,
			invSlice: invSlice,
			found:    found,
		})
		if err != nil {
			return err
		}

		logger.WithFields(
			logrus.Fields{
				"type":  "inventory",
				"stock": "strain",
				"event": "create",
				"id":    id,
				"count": len(invSlice),
			}).Debugf("created inventories")
		invCount = invCount + len(invSlice)
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

func cacheInvByStrainId(ir stockcenter.StrainInventoryReader, logger *logrus.Entry) (map[string][]*stockcenter.StrainInventory, error) {
	invMap := make(map[string][]*stockcenter.StrainInventory)
	readCount := 0
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
					"id":     inv.StrainId,
				}).Warnf("skipped the record")
			continue
		}
		if v, ok := invMap[inv.StrainId]; ok {
			invMap[inv.StrainId] = append(v, inv)
		} else {
			invMap[inv.StrainId] = []*stockcenter.StrainInventory{inv}
		}
		readCount += 1
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": "strains",
			"event": "read",
			"count": readCount,
		}).Infof("read all record")
	return invMap, nil
}

func createStrainInventory(args *strainInvArgs) error {
	for i, inv := range args.invSlice {
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
		if !inv.StoredOn.IsZero() {
			m[regs.INV_STORAGE_DATE_TAG] = inv.StoredOn.Format(time.RFC3339Nano)
		}
	INNER:
		for tag, value := range m {
			if len(value) == 0 {
				continue INNER
			}
			anno, err := createAnnoWithRank(args.client, tag, inv.StrainId, regs.STRAIN_INV_ONTO, value, i)
			if err != nil {
				return err
			}
			ids = append(ids, anno.Data.Id)
		}
		_, err := args.client.CreateAnnotationGroup(context.Background(), &pb.AnnotationIdList{Ids: ids})
		if err != nil {
			return err
		}
	}
	// create presence of inventory annotation
	if !args.found {
		_, err := createAnno(
			args.client, regs.STRAIN_INV_ONTO, args.id,
			regs.STRAIN_INV_ONTO, regs.INV_EXIST_VALUE,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
