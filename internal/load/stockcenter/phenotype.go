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
	// group all phenotypes by their ids
	stPheno := make(map[string][]*stockcenter.Phenotype)
	for pr.Next() {
		pheno, err := pr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in loading phenotype for strain %s",
				err,
			)
		}
		stPheno[pheno.StrainId] = append(stPheno[pheno.StrainId], pheno)
	}
	phcount := 0
	for sid, phgrp := range stPheno {
		gc, err := client.ListAnnotationGroups(
			context.Background(),
			&pb.ListGroupParameters{
				Limit:  20,
				Filter: fmt.Sprintf("entry_id==%s", sid),
			})
		if err != nil {
			if grpc.Code(err) == codes.NotFound {
				logger.Infof("no phenotype group found for id %s", sid)
				continue
			}
			return err
		}
		// fetch and remove all phenotype(groups) for every strains
		for _, g := range gc.Data {
			_, err := client.DeleteAnnotationGroup(
				context.Background(),
				&pb.GroupEntryId{GroupId: g.Group.GroupId},
			)
			if err != nil {
				return err
			}
			logger.Debugf("removing group %s for strain id %s", g.Group.GroupId, sid)
		}
		for _, ph := range phgrp {
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
				ta, err := findOrCreateAnno(client, ph.Assay, sid, regs.DICTY_ANNO_ONTOLOGY, "novalue")
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
			if err != nil {
				return err
			}
			phcount++
		}
		logger.Debugf("created %d phenotypes for %s strain", len(phgrp), sid)
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
