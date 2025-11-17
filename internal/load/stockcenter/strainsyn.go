package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/config"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadStrainSynProp(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewTsvStockPropReader(
		registry.GetReader(regs.StrainSynReader),
	)
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	synMap, err := readStrainSynonyms(pr)
	if err != nil {
		return err
	}
	logger.Debugf("going to load %d synonyms", len(synMap))
	pcount, err := processSynonyms(synMap, client, logger)
	if err != nil {
		return err
	}
	logProcessedSynonyms(pcount, logger)
	return nil
}

func readStrainSynonyms(
	pr stockcenter.StockPropReader,
) (map[string][]*stockcenter.StockProp, error) {
	synMap := make(map[string][]*stockcenter.StockProp)
	for pr.Next() {
		prop, err := pr.Value()
		if err != nil {
			return nil, fmt.Errorf(
				"error in reading property for strain %s",
				err,
			)
		}
		if prop.Property != synTag {
			continue
		}
		synMap[prop.Id] = append(synMap[prop.Id], prop)
	}
	return synMap, nil
}

func processSynonyms(
	synMap map[string][]*stockcenter.StockProp,
	client annotation.TaggedAnnotationServiceClient,
	logger *logrus.Entry,
) (int, error) {
	pcount := 0
	for entryID, props := range synMap {
		if err := removeExistingSynonyms(entryID, client, logger); err != nil {
			return pcount, err
		}
		if err := reloadSynonyms(entryID, props, client, logger); err != nil {
			return pcount, err
		}
		pcount++
	}
	return pcount, nil
}

func removeExistingSynonyms(
	entryID string,
	client annotation.TaggedAnnotationServiceClient,
	logger *logrus.Entry,
) error {
	tac, err := client.ListAnnotations(getContext(), &annotation.ListParameters{
		Limit:  config.DefaultCSVWorkerPoolSize,
		Filter: buildFilter(entryID),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("error in listing synonyms for %s %s", entryID, err)
	}
	if tac == nil {
		logger.Debugf("synonym %s is absent, no need to remove it", entryID)
		return nil
	}
	for _, ta := range tac.Data {
		if err := deleteAnnotation(ta.Id, client); err != nil {
			return fmt.Errorf(
				"unable to remove synonyms for %s %s",
				entryID,
				err,
			)
		}
	}
	logger.Debugf("removed %d synonyms for id %s", len(tac.Data), entryID)
	return nil
}

func deleteAnnotation(
	annotationID string,
	client annotation.TaggedAnnotationServiceClient,
) error {
	_, err := client.DeleteAnnotation(
		getContext(),
		&annotation.DeleteAnnotationRequest{
			Id:    annotationID,
			Purge: true,
		},
	)
	return err
}

func reloadSynonyms(
	entryID string,
	props []*stockcenter.StockProp,
	client annotation.TaggedAnnotationServiceClient,
	logger *logrus.Entry,
) error {
	for i, p := range props {
		if err := createAnnotation(entryID, i, p, client); err != nil {
			return fmt.Errorf(
				"unable to load synonym %s for %s %s",
				p.Value,
				entryID,
				err,
			)
		}
	}
	logger.Debugf("loaded all %d synonyms for %s", len(props), entryID)
	return nil
}

func createAnnotation(
	entryID string,
	index int,
	prop *stockcenter.StockProp,
	client annotation.TaggedAnnotationServiceClient,
) error {
	_, err := client.CreateAnnotation(
		getContext(),
		&annotation.NewTaggedAnnotation{
			Data: &annotation.NewTaggedAnnotation_Data{
				Attributes: &annotation.NewTaggedAnnotationAttributes{
					Value:     prop.Value,
					CreatedBy: regs.DefaultUser,
					Tag:       synTag,
					EntryId:   entryID,
					Ontology:  regs.DictyAnnoOntology,
					Rank:      int64(index),
				},
			},
		},
	)
	return err
}

func buildFilter(entryID string) string {
	return fmt.Sprintf(
		"entry_id===%s;tag===%s;ontology===%s",
		entryID, synTag, regs.DictyAnnoOntology,
	)
}

func getContext() context.Context {
	return context.Background()
}

func logProcessedSynonyms(pcount int, logger *logrus.Entry) {
	logger.WithFields(
		logrus.Fields{
			"type":  "synonym",
			"stock": "strain",
			"event": "load",
			"count": pcount,
		}).Infof("loaded strain synonym")
}
