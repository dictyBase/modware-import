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
)

func LoadStrainInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewTsvStrainInventoryReader(registry.GetReader(regs.InvReader))
	logger := registry.GetLogger().WithFields(logrus.Fields{
		"type":  "inventory",
		"stock": "strain",
	})
	invMap, err := cacheInvByStrainId(ir, logger)
	if err != nil {
		return err
	}
	client := regs.GetAnnotationAPIClient()
	invCount := 0
	for id, invSlice := range invMap {
		gc, err := getInventory(id, client, regs.StrainInvOnto)
		found, err := handleAnnoRetrieval(&annoParams{
			id:     id,
			gc:     gc,
			err:    err,
			client: client,
			logger: logger,
			loader: "inventory",
		})
		if err != nil {
			return err
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
		logger.WithFields(logrus.Fields{
			"event": "create",
			"id":    id,
			"count": len(invSlice),
		}).Debug("created inventories")
		invCount += len(invSlice)
	}
	logger.WithFields(logrus.Fields{
		"event": "load",
		"count": invCount,
	}).Info("loaded inventories")
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
			logger.WithFields(logrus.Fields{
				"event":  "skip record",
				"output": inv.RecordLine,
				"id":     inv.StrainID,
			}).Warn("skipped the record")
			continue
		}
		if v, ok := invMap[inv.StrainID]; ok {
			invMap[inv.StrainID] = append(v, inv)
		} else {
			invMap[inv.StrainID] = []*stockcenter.StrainInventory{inv}
		}
		readCount += 1
	}
	logger.WithFields(logrus.Fields{
		"event": "read",
		"count": readCount,
	}).Debug("read all record")
	return invMap, nil
}

func organizeStrainInvAnno(inv *stockcenter.StrainInventory) map[string]string {
	return map[string]string{
		regs.InvLocationTag:    inv.PhysicalLocation,
		regs.InvStoredAsTag:    inv.StoredAs,
		regs.InvVialCountTag:   inv.VialsCount,
		regs.InvVialColorTag:   inv.VialColor,
		regs.InvPrivCommentTag: inv.PrivateComment,
		regs.InvPubCommentTag:  inv.PublicComment,
		regs.InvStorageDateTag: inv.StoredOn.Format(time.RFC3339Nano),
	}
}

func createStrainInventory(args *strainInvArgs) error {
	for i, inv := range args.invSlice {
		var ids []string
		m := organizeStrainInvAnno(inv)
		if !inv.StoredOn.IsZero() {
			m[regs.InvStorageDateTag] = inv.StoredOn.Format(time.RFC3339Nano)
		}
	INNER:
		for tag, value := range m {
			if len(value) == 0 {
				continue INNER
			}
			anno, err := createAnnoWithRank(&createAnnoArgs{
				ontology: regs.StrainInvOnto,
				client:   args.client,
				id:       inv.StrainID,
				value:    value,
				tag:      tag,
				rank:     i,
			})
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
		return createAnno(&createAnnoArgs{
			client:   args.client,
			tag:      regs.StrainInvTag,
			id:       args.id,
			ontology: regs.StrainInvOnto,
			value:    regs.InvExistValue,
			user:     regs.DEFAULT_USER,
		})
	}
	return nil
}
