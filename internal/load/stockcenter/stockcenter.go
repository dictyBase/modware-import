package stockcenter

import (
	"context"
	"fmt"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/sirupsen/logrus"
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

func getInventory(id string, client pb.TaggedAnnotationServiceClient, onto, stock string, logger *logrus.Entry) (*pb.TaggedAnnotationGroupCollection, error) {
	gc, err := client.ListAnnotationGroups(
		context.Background(),
		&pb.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;tag==%s;ontology==%s",
				id, regs.INV_LOCATION_TAG, onto,
			),
		})
	if err != nil {
		if grpc.Code(err) != codes.NotFound { // error in lookup
			return gc, err
		}
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": stock,
			"event": "get",
			"id":    id,
		}).Debugf("retrieved inventories")

	return gc, nil
}

func delExistingInventory(id string, client pb.TaggedAnnotationServiceClient, stock string, gc *pb.TaggedAnnotationGroupCollection, logger *logrus.Entry) error {
	for _, gcd := range gc.Data {
		// remove annotations group
		_, err := client.DeleteAnnotationGroup(
			context.Background(),
			&pb.GroupEntryId{GroupId: gcd.Group.GroupId},
		)
		if err != nil {
			return err
		}
		// remove all annotations
		for _, gd := range gcd.Group.Data {
			_, err := client.DeleteAnnotation(
				context.Background(),
				&pb.DeleteAnnotationRequest{Id: gd.Id, Purge: true},
			)
			if err != nil {
				return err
			}
		}
	}
	logger.WithFields(
		logrus.Fields{
			"type":  "inventory",
			"stock": stock,
			"event": "delete",
			"id":    id,
		}).Debugf("deleted inventories")
	return nil
}
