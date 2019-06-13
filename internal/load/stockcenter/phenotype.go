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
	user          = "dictybase@northwestern.edu"
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
			if err := handlePhenotype(client, pheno, pheno.StrainId); err != nil {
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
		if err := handlePhenotype(client, pheno, pheno.StrainId); err != nil {
			return err
		}
		logger.Debugf("flush and loaded phenotypes for %s strain", pheno.StrainId)
		phcount++
	}
	logger.Infof("created %d phenotypes", phcount)
	return nil
}

func findOrCreateAnno(client pb.TaggedAnnotationServiceClient, tag, id, ontology, value string) (*pb.TaggedAnnotation, error) {
	ta, err := client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      tag,
			EntryId:  id,
			Ontology: ontology,
		})
	switch {
	case err == nil:
		return ta, nil
	case grpc.Code(err) == codes.NotFound:
		return client.CreateAnnotation(
			context.Background(),
			&pb.NewTaggedAnnotation{
				Data: &pb.NewTaggedAnnotation_Data{
					Attributes: &pb.NewTaggedAnnotationAttributes{
						Value:     value,
						CreatedBy: user,
						Tag:       tag,
						EntryId:   id,
						Ontology:  ontology,
					},
				},
			},
		)

	}
	return ta, fmt.Errorf(
		"error in finding annotation %s for id %s %s",
		tag,
		id,
		err,
	)
}

func handlePhenotype(client pb.TaggedAnnotationServiceClient, ph *stockcenter.Phenotype, sid string) error {
	var ids []string
	to, err := findOrCreateAnno(client, ph.Observation, sid, phenoOntology, "novalue")
	if err != nil {
		return err
	}
	ids = append(ids, to.Data.Id)
	tl, err := findOrCreateAnno(client, literatureTag, sid, regs.DICTY_ANNO_ONTOLOGY, ph.LiteratureId)
	if err != nil {
		return err
	}
	ids = append(ids, tl.Data.Id)
	if len(ph.Note) > 1 {
		tn, err := findOrCreateAnno(client, noteTag, sid, regs.DICTY_ANNO_ONTOLOGY, ph.Note)
		if err != nil {
			return err
		}
		ids = append(ids, tn.Data.Id)
	}
	if len(ph.Assay) > 1 {
		ta, err := findOrCreateAnno(client, ph.Assay, sid, assayOntology, "novalue")
		if err != nil {
			return err
		}
		ids = append(ids, ta.Data.Id)
	}
	if len(ph.Environment) > 1 {
		te, err := findOrCreateAnno(client, ph.Environment, sid, envOntology, "novalue")
		if err != nil {
			return err
		}
		ids = append(ids, te.Data.Id)
	}
	_, err = client.CreateAnnotationGroup(context.Background(), &pb.AnnotationIdList{Ids: ids})
	return err
}
