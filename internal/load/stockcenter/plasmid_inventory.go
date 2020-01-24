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

func LoadPlasmidInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewTsvPlasmidInventoryReader(registry.GetReader(regs.InvReader))
	logger := registry.GetLogger()
	invMap, err := cacheInvByPlasmidId(ir, logger)
	if err != nil {
		return err
	}
	logger.Debugf("cached %d plasmid inventories", len(invMap))
	client := regs.GetAnnotationAPIClient()
	invCount := 0
	for id, invSlice := range invMap {
		gc, err := getInventory(id, client, regs.PlasmidInvOntO)
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
		err = createPlasmidInventory(&plasmidInvArgs{
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
				"stock": "plasmid",
				"event": "create",
				"id":    id,
				"count": len(invSlice),
			}).Debugf("created inventories")
		invCount += len(invSlice)
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
					"event":  "skip",
					"output": inv.RecordLine,
				}).Warnf("skipped the record")
			continue
		}
		if invSlice, ok := invMap[inv.PlasmidID]; ok {
			invMap[inv.PlasmidID] = append(invSlice, inv)
			continue
		}
		invMap[inv.PlasmidID] = []*stockcenter.PlasmidInventory{inv}
	}
	return invMap, nil
}

func organizePlasmidInvAnno(inv *stockcenter.PlasmidInventory) map[string]string {
	return map[string]string{
		regs.InvLocationTag:    inv.PhysicalLocation,
		regs.InvStoredAsTag:    inv.StoredAs,
		regs.InvPrivCommentTag: inv.PrivateComment,
		regs.InvObtainedAsTag:  inv.ObtainedAs,
	}
}

func createPlasmidInventory(args *plasmidInvArgs) error {
	for i, inv := range args.invSlice {
		var ids []string
		m := organizePlasmidInvAnno(inv)
		if !inv.StoredOn.IsZero() {
			m[regs.InvStorageDateTag] = inv.StoredOn.Format(time.RFC3339Nano)
		}
	INNER:
		for tag, value := range m {
			if len(value) == 0 {
				continue INNER
			}
			anno, err := createAnnoWithRank(&createAnnoArgs{
				ontology: regs.PlasmidInvOntO,
				client:   args.client,
				id:       inv.PlasmidID,
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
		return createAnno(
			args.client, regs.PlasmidInvTag, args.id,
			regs.PlasmidInvOntO, regs.InvExistValue,
		)
	}
	return nil
}
