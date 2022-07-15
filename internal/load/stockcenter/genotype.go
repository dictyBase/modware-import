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

func LoadGeno(cmd *cobra.Command, args []string) error {
	gr := stockcenter.NewTsvGenotypeReader(registry.GetReader(regs.GenoReader))
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	count := 0
	uct := 0
	rct := 0
	nct := 0
	for gr.Next() {
		geno, err := gr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in loading genotype for strain %s",
				err,
			)
		}
		st, err := NewOrReloadGeno(
			client,
			&genoArgs{
				tag:      genoTag,
				id:       geno.StrainId,
				user:     regs.DefaultUser,
				ontology: regs.DictyAnnoOntology,
				value:    geno.Genotype,
			})
		if err != nil {
			return err
		}
		switch st {
		case Created:
			logger.Debugf("created genotype %s for strain %s", geno.Genotype, geno.StrainId)
			nct++
		case Updated:
			logger.Debugf("updated genotype %s for strain %s", geno.Genotype, geno.StrainId)
			uct++
		case Read:
			logger.Debugf("skipped genotype %s for strain %s", geno.Genotype, geno.StrainId)
			rct++
		}
		count += 1
	}
	logger.WithFields(
		logrus.Fields{
			"type":    "genotype",
			"stock":   "strain",
			"event":   "load",
			"count":   count,
			"read":    rct,
			"created": nct,
			"updated": uct,
		}).Infof("loaded strain genotypes")
	return nil
}

func NewOrReloadGeno(client pb.TaggedAnnotationServiceClient, p *genoArgs) (Status, error) {
	ta, err := client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      p.tag,
			EntryId:  p.id,
			Ontology: p.ontology,
		})
	switch {
	case err == nil: // exists, so check and update
		if p.value == ta.Data.Attributes.Value {
			return Read, err
		}
		data := &pb.TaggedAnnotationUpdate_Data{
			Id: ta.Data.Id,
			Attributes: &pb.TaggedAnnotationUpdateAttributes{
				Value:         p.value,
				EditableValue: p.value,
				CreatedBy:     ta.Data.Attributes.CreatedBy,
			},
		}
		_, err = client.UpdateAnnotation(
			context.Background(),
			&pb.TaggedAnnotationUpdate{Data: data},
		)
		return Updated, err
	case status.Code(err) == codes.NotFound: // create a new one
		data := &pb.NewTaggedAnnotation_Data{
			Attributes: &pb.NewTaggedAnnotationAttributes{
				Value:     p.value,
				CreatedBy: p.user,
				Tag:       p.tag,
				EntryId:   p.id,
				Ontology:  p.ontology,
			},
		}
		_, err = client.CreateAnnotation(
			context.Background(),
			&pb.NewTaggedAnnotation{Data: data},
		)
		return Created, err
	}
	return Nop, fmt.Errorf(
		"error in finding annotation %s for id %s %s",
		p.tag,
		p.id,
		err,
	)
}
