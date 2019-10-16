package stockcenter

import (
	"context"
	"fmt"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func createAnno(client pb.TaggedAnnotationServiceClient, tag, id, ontology, value string) (*pb.TaggedAnnotation, error) {
	ta, err := client.CreateAnnotation(
		context.Background(),
		&pb.NewTaggedAnnotation{
			Data: &pb.NewTaggedAnnotation_Data{
				Attributes: &pb.NewTaggedAnnotationAttributes{
					Value:     value,
					CreatedBy: regs.DEFAULT_USER,
					Tag:       tag,
					EntryId:   id,
					Ontology:  ontology,
				},
			},
		},
	)
	return ta, fmt.Errorf(
		"error in creating annotation %s for id %s %s",
		tag,
		id,
		err,
	)
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
						CreatedBy: regs.DEFAULT_USER,
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
