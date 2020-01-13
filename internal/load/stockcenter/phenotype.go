package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func LoadPheno(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewPhenotypeReader(registry.GetReader(regs.PHENO_READER))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	phenoMap, err := cachePhenotype(pr, logger)
	if err != nil {
		return err
	}
	count := 0
	for id, phenoSlice := range phenoMap {
		found := true
		gc, err := getPhenotype(client, phenoSlice[0])
		if err != nil {
			if status.Code(err) != codes.NotFound { // error in lookup
				return fmt.Errorf("error in getting phenotype of %s %s", id, err)
			}
			found = false
			logger.WithFields(
				logrus.Fields{
					"type":  "phenotype",
					"stock": "strain",
					"event": "get",
					"id":    id,
				}).Debugf("no phenotype")
		}
		if found {
			logger.WithFields(
				logrus.Fields{
					"type":  "phenotype",
					"stock": "strain",
					"event": "get",
					"id":    id,
				}).Debugf("retrieved phenotype")
			if err := delAnnotationGroup(client, gc); err != nil {
				return err
			}
			logger.WithFields(
				logrus.Fields{
					"type":  "phenotype",
					"stock": "strain",
					"event": "delete",
					"id":    id,
				}).Debugf("deleted phenotype")
		}
		err = createPhenotype(&strainPhenoArgs{
			id:         id,
			client:     client,
			phenoSlice: phenoSlice,
		})
		if err != nil {
			return err
		}
		logger.WithFields(
			logrus.Fields{
				"type":  "phenotype",
				"stock": "strain",
				"event": "create",
				"id":    id,
				"count": len(phenoSlice),
			}).Debugf("created phenotypes")
		count += len(phenoSlice)
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "phenotype",
			"stock": "strains",
			"event": "load",
			"count": count,
		}).Infof("loaded phenotypes")
	return nil
}

func createPhenotype(args *strainPhenoArgs) error {
	for i, pheno := range args.phenoSlice {
		var ids []string
		m := map[string][]string{
			regs.PhenoOntology: []string{pheno.Observation, regs.EmptyValue},
			regs.AssayOntology: []string{pheno.Assay, regs.EmptyValue},
			regs.EnvOntology:   []string{pheno.Environment, regs.EmptyValue},
			regs.DICTY_ANNO_ONTOLOGY: []string{
				regs.LiteratureTag, pheno.LiteratureId,
				regs.NoteTag, pheno.Note},
		}
	INNER:
		for onto, dataSlice := range m {
			if onto == regs.DICTY_ANNO_ONTOLOGY {
				if len(dataSlice[1]) > 0 {
					anno, err := createAnnoWithRank(args.client, dataSlice[0], args.id, onto, dataSlice[1], i)
					if err != nil {
						return err
					}
					ids = append(ids, anno.Data.Id)
				}
				if len(dataSlice[3]) > 0 {
					anno, err := createAnnoWithRank(args.client, dataSlice[2], args.id, onto, dataSlice[3], i)
					if err != nil {
						return err
					}
					ids = append(ids, anno.Data.Id)
				}
				continue INNER
			}
			if len(dataSlice[1]) == 0 {
				continue INNER
			}
			anno, err := createAnnoWithRank(args.client, dataSlice[0], args.id, onto, dataSlice[1], i)
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
	return nil
}

func cachePhenotype(pr stockcenter.PhenotypeReader, logger *logrus.Entry) (map[string][]*stockcenter.Phenotype, error) {
	phenoMap := make(map[string][]*stockcenter.Phenotype)
	readCount := 0
	for pr.Next() {
		pheno, err := pr.Value()
		if err != nil {
			return phenoMap, fmt.Errorf(
				"error in loading strain phenotype %s", err,
			)
		}
		if phslice, ok := phenoMap[pheno.StrainId]; ok {
			phenoMap[pheno.StrainId] = append(phslice, pheno)
		} else {
			phenoMap[pheno.StrainId] = []*stockcenter.Phenotype{pheno}
		}
		readCount++
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "phenotype",
			"stock": "strains",
			"event": "read",
			"count": readCount,
		}).Debugf("read all record")
	return phenoMap, nil
}

func getPhenotype(client pb.TaggedAnnotationServiceClient, pheno *stockcenter.Phenotype) (*pb.TaggedAnnotationGroupCollection, error) {
	return client.ListAnnotationGroups(
		context.Background(),
		&pb.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;tag==%s;ontology==%s",
				pheno.StrainId,
				pheno.Observation,
				regs.PhenoOntology,
			),
		})
}
