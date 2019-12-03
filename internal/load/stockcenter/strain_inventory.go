package stockcenter

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadStrainInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewTsvStrainInventoryReader(registry.GetReader(regs.InvReader))
	logger := registry.GetLogger()
	invMap, err := cacheInvByStrainId(ir, logger)
	if err != nil {
		return err
	}
	client := regs.GetAnnotationAPIClient()
	invCount := 0
	for id, invSlice := range invMap {
		found := true
		gc, err := getInventory(id, client, regs.StrainInvOnto)
		if err != nil {
			if status.Code(err) != codes.NotFound { // error in lookup
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
			if err := delExistingInventory(client, gc); err != nil {
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
		invCount += len(invSlice)
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
					"id":     inv.StrainID,
				}).Warnf("skipped the record")
			continue
		}
		if v, ok := invMap[inv.StrainID]; ok {
			invMap[inv.StrainID] = append(v, inv)
		} else {
			invMap[inv.StrainID] = []*stockcenter.StrainInventory{inv}
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
			regs.InvLocationTag:    inv.PhysicalLocation,
			regs.InvStoredAsTag:    inv.StoredAs,
			regs.InvVialCountTag:   inv.VialsCount,
			regs.InvVialColorTag:   inv.VialColor,
			regs.InvPrivCommentTag: inv.PrivateComment,
			regs.InvPubCommentTag:  inv.PublicComment,
			regs.InvStorageDateTag: inv.StoredOn.Format(time.RFC3339Nano),
		}
		if !inv.StoredOn.IsZero() {
			m[regs.InvStorageDateTag] = inv.StoredOn.Format(time.RFC3339Nano)
		}
	INNER:
		for tag, value := range m {
			if len(value) == 0 {
				continue INNER
			}
			anno, err := createAnnoWithRank(args.client, tag, inv.StrainID, regs.StrainInvOnto, value, i)
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
		return createAnno(
			args.client, regs.StrainInvOnto, args.id,
			regs.StrainInvOnto, regs.InvExistValue,
		)
	}
	return nil
}
