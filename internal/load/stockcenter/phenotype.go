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
	phenoMap, err := processPhenotype(&processPhenoArgs{
		pr:     pr,
		logger: logger,
		client: client,
	})
	if err != nil {
		return err
	}
	count := 0
	for id, phenoSlice := range phenoMap {
		found := true
		gc, err := getPhenotype(&getPhenoArgs{
			id:       id,
			client:   client,
			ontology: regs.PhenoOntology,
		})
		if err != nil {
			if status.Code(err) != codes.NotFound { // error in lookup
				return fmt.Errorf("error in getting phenotype of %s %s", id, err)
			}
			found = false
			logger.WithFields(logrus.Fields{
				"type":  "phenotype",
				"stock": "strain",
				"event": "get",
				"id":    id,
			}).Debugf("no phenotype")
		}
		if found {
			logger.WithFields(logrus.Fields{
				"type":  "phenotype",
				"stock": "strain",
				"event": "get",
				"id":    id,
			}).Debugf("retrieved phenotype")
			if err := delAnnotationGroup(client, gc); err != nil {
				return err
			}
			logger.WithFields(logrus.Fields{
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
		logger.WithFields(logrus.Fields{
			"type":  "phenotype",
			"stock": "strain",
			"event": "create",
			"id":    id,
			"count": len(phenoSlice),
		}).Debugf("created phenotypes")
		count += len(phenoSlice)
	}
	logger.WithFields(logrus.Fields{
		"type":  "phenotype",
		"stock": "strains",
		"event": "load",
		"count": count,
	}).Infof("loaded phenotypes")
	return nil
}

func organizePhenoAnno(pheno *stockcenter.Phenotype) map[string][]string {
	return map[string][]string{
		regs.PhenoOntology: {pheno.Observation, regs.EmptyValue},
		regs.AssayOntology: {pheno.Assay, regs.EmptyValue},
		regs.EnvOntology:   {pheno.Environment, regs.EmptyValue},
	}
}

func organizeMorePhenoAnno(pheno *stockcenter.Phenotype) [][]string {
	return [][]string{
		{regs.LiteratureTag, pheno.LiteratureID},
		{regs.NoteTag, pheno.Note},
	}
}

func createPhenotype(args *strainPhenoArgs) error {
	for i, pheno := range args.phenoSlice {
		var ids []string
	FIRST:
		for _, dataSlice := range organizeMorePhenoAnno(pheno) {
			if len(dataSlice[1]) == 0 {
				continue FIRST
			}
			anno, err := createAnnoWithRank(&createAnnoArgs{
				ontology: regs.DICTY_ANNO_ONTOLOGY,
				client:   args.client,
				tag:      dataSlice[0],
				value:    dataSlice[1],
				id:       args.id,
				rank:     i,
			})
			if err != nil {
				return err
			}
			ids = append(ids, anno.Data.Id)

		}
	SECOND:
		for onto, dataSlice := range organizePhenoAnno(pheno) {
			if len(dataSlice[0]) == 0 {
				continue SECOND
			}
			anno, err := createAnnoWithRank(&createAnnoArgs{
				client:   args.client,
				tag:      dataSlice[0],
				value:    dataSlice[1],
				id:       args.id,
				ontology: onto,
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
	return nil
}

func processPhenotype(args *processPhenoArgs) (map[string][]*stockcenter.Phenotype, error) {
	phenoMap := make(map[string][]*stockcenter.Phenotype)
	readCount := 0
	for args.pr.Next() {
		pheno, err := args.pr.Value()
		if err != nil {
			return phenoMap, fmt.Errorf("error in loading strain phenotype %s", err)
		}
		phStatus, err := validateAnnoTag(&validateTagArgs{
			client:   args.client,
			logger:   args.logger,
			tag:      pheno.Observation,
			ontology: regs.PhenoOntology,
			id:       pheno.StrainID,
			stock:    "strain",
			loader:   "phenotype",
		})
		if err != nil {
			return phenoMap, err
		}
		if !phStatus {
			continue
		}
		if len(pheno.Assay) > 0 {
			status, err := validateAnnoTag(&validateTagArgs{
				client:   args.client,
				logger:   args.logger,
				tag:      pheno.Assay,
				ontology: regs.AssayOntology,
				id:       pheno.StrainID,
				stock:    "strain",
				loader:   "phenotype",
			})
			if err != nil {
				return phenoMap, err
			}
			if !status {
				continue
			}
		}
		if phslice, ok := phenoMap[pheno.StrainID]; ok {
			phenoMap[pheno.StrainID] = append(phslice, pheno)
		} else {
			phenoMap[pheno.StrainID] = []*stockcenter.Phenotype{pheno}
		}
		readCount++
	}
	args.logger.WithFields(logrus.Fields{
		"type":  "phenotype",
		"stock": "strains",
		"event": "read",
		"count": readCount,
	}).Infof("read all record")
	return phenoMap, nil
}

func getPhenotype(args *getPhenoArgs) (*pb.TaggedAnnotationGroupCollection, error) {
	return args.client.ListAnnotationGroups(
		context.Background(),
		&pb.ListGroupParameters{
			Filter: fmt.Sprintf("entry_id==%s;ontology==%s", args.id, args.ontology),
			Limit:  100,
		})
}
