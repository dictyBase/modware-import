package stockcenter

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type strainPhenoArgs struct {
	id         string
	client     pb.TaggedAnnotationServiceClient
	phenoSlice []*stockcenter.Phenotype
}

type strainInvArgs struct {
	id       string
	client   pb.TaggedAnnotationServiceClient
	invSlice []*stockcenter.StrainInventory
	found    bool
}

type plasmidInvArgs struct {
	id       string
	client   pb.TaggedAnnotationServiceClient
	invSlice []*stockcenter.PlasmidInventory
	found    bool
}

func createAnnoWithRank(client pb.TaggedAnnotationServiceClient, tag, id, ontology, value string, rank int) (*pb.TaggedAnnotation, error) {
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
					Rank:      int64(rank),
				},
			},
		},
	)
	if err != nil {
		return ta, fmt.Errorf(
			"error in creating annotation %s for id %s with rank %d %s",
			tag,
			id,
			rank,
			err,
		)
	}
	return ta, nil
}

func createAnno(client pb.TaggedAnnotationServiceClient, tag, id, ontology, value string) error {
	_, err := client.CreateAnnotation(
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
	if err != nil {
		return fmt.Errorf(
			"error in creating annotation %s for id %s %s",
			tag,
			id,
			err,
		)
	}
	return nil
}

func findOrCreateAnnoWithRank(client pb.TaggedAnnotationServiceClient, tag, id, ontology, value string, rank int) (*pb.TaggedAnnotation, bool, error) {
	create := false
	ta, err := client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      tag,
			EntryId:  id,
			Ontology: ontology,
			Rank:     int64(rank),
		})
	switch {
	case err == nil:
		return ta, create, nil
	case status.Code(err) == codes.NotFound:
		_, err := client.CreateAnnotation(
			context.Background(),
			&pb.NewTaggedAnnotation{
				Data: &pb.NewTaggedAnnotation_Data{
					Attributes: &pb.NewTaggedAnnotationAttributes{
						Value:     value,
						CreatedBy: regs.DEFAULT_USER,
						Tag:       tag,
						EntryId:   id,
						Ontology:  ontology,
						Rank:      int64(rank),
					},
				},
			},
		)
		if err != nil {
			return ta, create, fmt.Errorf(
				"error in creating annotation %s for id %s %s",
				tag,
				id,
				err,
			)
		}
		create = true
	}
	return ta, create, err
}

func findOrCreateAnnoWithStatus(client pb.TaggedAnnotationServiceClient, tag, id, ontology, value string) (bool, error) {
	create := false
	_, err := client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      tag,
			EntryId:  id,
			Ontology: ontology,
		})
	switch {
	case err == nil:
		return create, nil
	case status.Code(err) == codes.NotFound:
		_, err := client.CreateAnnotation(
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
		if err != nil {
			return create, fmt.Errorf(
				"error in finding annotation %s for id %s %s",
				tag,
				id,
				err,
			)
		}
		create = true
	}
	return create, nil
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
	case status.Code(err) == codes.NotFound:
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

func getInventory(id string, client pb.TaggedAnnotationServiceClient, onto string) (*pb.TaggedAnnotationGroupCollection, error) {
	return client.ListAnnotationGroups(
		context.Background(),
		&pb.ListGroupParameters{
			Filter: fmt.Sprintf(
				"entry_id==%s;tag==%s;ontology==%s",
				id, regs.InvLocationTag, onto,
			),
		})
}

func delAnnotationGroup(client pb.TaggedAnnotationServiceClient, gc *pb.TaggedAnnotationGroupCollection) error {
	for _, gcd := range gc.Data {
		// remove annotations group
		_, err := client.DeleteAnnotationGroup(
			context.Background(),
			&pb.GroupEntryId{GroupId: gcd.Group.GroupId},
		)
		if err != nil {
			return fmt.Errorf("error in deleting annotation group %s %s", gcd.Group.GroupId, err)
		}
		// remove all annotations
		for _, gd := range gcd.Group.Data {
			_, err := client.DeleteAnnotation(
				context.Background(),
				&pb.DeleteAnnotationRequest{Id: gd.Id, Purge: true},
			)
			if err != nil {
				return fmt.Errorf("error in deleting annotation %s %s", gd.Id, err)
			}
		}
	}
	return nil
}

func TimestampProto(t time.Time) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(t)
	return ts
}
