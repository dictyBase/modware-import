package stockcenter

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/csv/stockcenter"
	"github.com/dictyBase/modware-import/internal/registry"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/spf13/cobra"
)

const (
	phenoOntology = "Dicty Phenotypes"
	envOntology   = "Dicty Environment"
	assayOntology = "Dictyostellium Assay"
	literatureTag = "literature_tag"
	noteTag       = "public note"
)

func LoadPheno(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewPhenotypeReader(registry.GetReader(regs.PHENO_READER))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	phcount := 0
	for pr.Next() {
		pheno, err := pr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in loading phenotype for strain %s",
				err,
			)
		}
		gc, err := client.ListAnnotationGroups(
			context.Background(),
			&pb.ListGroupParameters{
				Filter: fmt.Sprintf(
					"entry_id==%s;tag==%s;ontology==%s",
					pheno.StrainId,
					pheno.Observation,
					phenoOntology,
				),
			})
		if err != nil {
			if grpc.Code(err) != codes.NotFound { // error in lookup
				return err
			}
			// no phenotype group found, create or find all individual annotations and phenotype
			if err := handlePhenotype(client, pheno); err != nil {
				return err
			}
			logger.Debugf("created phenotypes for %s strain", pheno.StrainId)
			phcount++
			continue
		}
		if len(gc.Data) > 1 {
			return fmt.Errorf(
				"data constraint issue, got multiple phenotype groups for %s %s %s",
				pheno.StrainId,
				pheno.Observation,
				phenoOntology,
			)
		}
		// delete phenotype group
		_, err = client.DeleteAnnotationGroup(
			context.Background(),
			&pb.GroupEntryId{GroupId: gc.Data[0].Group.GroupId},
		)
		if err != nil {
			return err
		}
		//create or find all individual annotations and phenotype
		if err := handlePhenotype(client, pheno); err != nil {
			return err
		}
		logger.Debugf("flush and loaded phenotypes for %s strain", pheno.StrainId)
		phcount++
	}
	logger.Infof("created %d phenotypes", phcount)
	return nil
}

func handlePhenotype(client pb.TaggedAnnotationServiceClient, ph *stockcenter.Phenotype) error {
	var ids []string
	phenMap := map[string][]string{
		ph.Observation: []string{phenoOntology, "novalue"},
		ph.Assay:       []string{assayOntology, "novalue"},
		ph.Environment: []string{envOntology, "novalue"},
		literatureTag:  []string{regs.DICTY_ANNO_ONTOLOGY, ph.LiteratureId},
		noteTag:        []string{regs.DICTY_ANNO_ONTOLOGY, ph.Note},
	}
	for t, s := range phenMap {
		t, err := findOrCreateAnno(client, t, ph.StrainId, s[0], s[1])
		if err != nil {
			return err
		}
		ids = append(ids, t.Data.Id)
	}
	_, err = client.CreateAnnotationGroup(context.Background(), &pb.AnnotationIdList{Ids: ids})
	return err
}
