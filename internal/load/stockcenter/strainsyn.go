package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
)

const (
	synTag = "synonym"
)

func LoadStrainSynProp(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewCsvStockPropReader(registry.GetReader(regs.STRAINPROP_READER))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	pcount := 0
	synMap := make(map[string][]*stockcenter.StockProp)
	for pr.Next() {
		prop, err := pr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in reading property for strain %s",
				err,
			)
		}
		if prop.Property != synTag {
			continue
		}
		// cache all synonyms
		if _, ok := synMap[prop.Id]; !ok {
			synMap[prop.Id] = []*stockcenter.StockProp{prop}
		} else {
			synMap[prop.Id] = append(synMap[prop.Id], prop)
		}
	}
	// load all the synonyms
	for entryId, props := range synMap {
		tac, err := client.ListAnnotations(
			context.Background(),
			&annotation.ListParameters{
				Limit: 20,
				Filter: fmt.Sprintf(
					"entry_id==%s;tag==%s;ontology==%s",
					entryId, synTag, regs.DICTY_ANNO_ONTOLOGY,
				)})
		if err != nil {
			if grpc.Code(err) != codes.NotFound {
				return fmt.Errorf("error in listing synonyms for %s %s", entryId, err)
			}
		} else {
			// remove synonyms
			for _, ta := range tac.Data {
				_, err := client.DeleteAnnotation(
					context.Background(),
					&annotation.DeleteAnnotationRequest{
						EntryId: ta.Id,
						Purge:   true,
					})
				if err != nil {
					return fmt.Errorf("unable to remove synonyms for %s %s", entryId, err)
				}
			}
		}
		// reload all synonyms
		for i, p := range props {
			_, err := client.CreateAnnotation(
				context.Background(),
				&annotation.NewTaggedAnnotation{
					Data: &annotation.NewTaggedAnnotation_Data{
						Attributes: &annotation.NewTaggedAnnotationAttributes{
							Value:     p.Value,
							CreatedBy: regs.DEFAULT_USER,
							Tag:       synTag,
							EntryId:   entryId,
							Ontology:  regs.DICTY_ANNO_ONTOLOGY,
							Rank:      int64(i),
						},
					},
				})
			if err != nil {
				return fmt.Errorf("unable to load synonym %s for %s %s", p.Value, entryId, err)
			}
		}
		logger.Debugf("loaded all strain synonyms for %s", entryId)
		pcount++
	}
	logger.Infof("loaded synonyms for %d strains", pcount)
	return nil
}
