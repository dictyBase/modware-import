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
	phenoOntology     = "Dicty Phenotypes"
	envOntology       = "Dicty Environment"
	assayOntology     = "Dictyostellium Assay"
	dictyAnnoOntology = "dicty_annotation"
	user              = "dictybase@northwestern.edu"
)

func LoadPheno(cmd *cobra.Command, args []string) error {
	pr := stockcenter.NewPhenotypeReader(regs.PHENO_READER)
	client := regs.GetAnnotationAPIClient()
	logger := registry.GetLogger()
	for pr.Next() {
		pheno, err := pr.Value()
		if err != nil {
			return fmt.Errorf(
				"error in loading phenotype for strain %s",
				err,
			)
		}
	}
	exObsAnno, err := client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      pheno.Observation,
			EntryId:  pheno.Id,
			Ontology: phenoOntology,
		})
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			// create new entry
			nObsAnno, err := client.CreateAnnotation(
				context.Background(),
				&pb.NewTaggedAnnotation{
					Data: &pb.NewTaggedAnnotation_Data{
						Type: "annotation",
						Attributes: &pb.NewTaggedAnnotationAttributes{
							Value:     "none",
							CreatedBy: user,
							Tag:       pheno.Observation,
							EntryId:   pheno.Id,
							Ontology:  phenoOntology,
						},
					},
				},
			)
			if err != nil {
				return fmt.Errorf(
					"error in creating phenotype observation %s for id %s %s",
					pheno.Observation,
					pheno.Id,
					err,
				)
			}
		} else {
			return fmt.Errorf(
				"error in finding phenotype observation %s for id %s %s",
				pheno.Observation,
				pheno.Id,
				err,
			)
		}
	}

	return nil
}
