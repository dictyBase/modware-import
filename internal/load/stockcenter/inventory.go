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
	"github.com/spf13/cobra"
)

const (
	invOntology    = "strain_inventory"
	locTag         = "location"
	pubCommentTag  = "public comment"
	privCommentTag = "private comment"
	storedAsTag    = "stored as"
	vialCountTag   = "number of vials"
	vialColorTag   = "color"
	storedOnTag    = "storage date"
)

func LoadInv(cmd *cobra.Command, args []string) error {
	ir := stockcenter.NewCsvStrainInventoryReader(registry.GetReader(regs.INV_READER))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	invcount := 0
	invMap := make(map[string][]*stockcenter.StrainInventory)
	for ir.Next() {
		inv, err := ir.Value()
		if err != nil {
			return fmt.Errorf(
				"error in loading inventory for strain %s",
				err,
			)
		}
		if len(inv.PhysicalLocation) == 0 || len(inv.VialColor) == 0 {
			logger.Warnf("skipped the record %s", inv.RecordLine)
			continue
		}
		if v, ok := invMap[inv.StrainId]; ok {
			invMap[inv.StrainId] = append(v, inv)
		} else {
			invMap[inv.StrainId] = []*stockcenter.StrainInventory{inv}
		}
	}
	gc, err := client.ListAnnotationGroups(
		context.Background(),
		&pb.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;tag==%s;ontology==%s;value==%s",
				inv.StrainId,
				locTag,
				invOntology,
				inv.PhysicalLocation,
			),
		})
	if err != nil {
		if grpc.Code(err) != codes.NotFound { // error in lookup
			return err
		}
		// no inventory found, create or find all individual annotations and inventory
		if err := handleInventory(client, inv); err != nil {
			return err
		}
		logger.Debugf("created inventory for %s strain", inv.StrainId)
		invcount++
		continue
	}
	if len(gc.Data) > 1 {
		return fmt.Errorf(
			"data constraint issue, got multiple inventory groups for %s %s %s %s",
			inv.StrainId,
			locTag,
			invOntology,
			inv.PhysicalLocation,
		)
	}
	// delete inventory
	_, err = client.DeleteAnnotationGroup(
		context.Background(),
		&pb.GroupEntryId{GroupId: gc.Data[0].Group.GroupId},
	)
	if err != nil {
		return err
	}
	//create or find all individual annotations and inventory
	if err := handleInventory(client, inv); err != nil {
		return err
	}
	logger.Debugf("flush and loaded inventory for %s strain", inv.StrainId)
	invcount++
	logger.Infof("created %d inventory", invcount)
	return nil
}

func handleInventory(client pb.TaggedAnnotationServiceClient, inv *stockcenter.StrainInventory) error {
	var ids []string
	invMap := map[string]string{
		locTag:         inv.PhysicalLocation,
		storedAsTag:    inv.StoredAs,
		vialCountTag:   inv.VialsCount,
		vialColorTag:   inv.VialColor,
		privCommentTag: inv.PrivateComment,
		pubCommentTag:  inv.PublicComment,
		storedOnTag:    inv.StoredOn.Format(time.RFC3339Nano),
	}
	for t, v := range invMap {
		if len(v) > 0 {
			t, err := findOrCreateAnno(client, t, inv.StrainId, invOntology, v)
			if err != nil {
				return err
			}
			ids = append(ids, t.Data.Id)
		}
	}
	_, err := client.CreateAnnotationGroup(context.Background(), &pb.AnnotationIdList{Ids: ids})
	return err
}
