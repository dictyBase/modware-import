package stockcenter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/sirupsen/logrus"

	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateAnnoTag(args *validateTagArgs) (bool, error) {
	tag, err := args.client.GetAnnotationTag(
		context.Background(),
		&pb.TagRequest{Name: args.tag, Ontology: args.ontology},
	)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			args.logger.WithFields(
				logrus.Fields{
					"type":     args.loader,
					"stock":    args.stock,
					"tag":      args.tag,
					"ontology": args.ontology,
					"id":       args.id,
					"event":    "non-existent tag",
				}).Warn("tag does not exist")
			return false, nil
		}
		return false, fmt.Errorf("error in tag lookup %s", err)
	}
	if tag.IsObsolete {
		args.logger.WithFields(
			logrus.Fields{
				"type":     args.loader,
				"stock":    args.stock,
				"tag":      args.tag,
				"ontology": args.ontology,
				"id":       args.id,
				"event":    "obsolete tag",
			}).Warn("tag is obsolete")
		return false, nil
	}
	return true, nil
}

func createAnnoWithRank(args *createAnnoArgs) (*pb.TaggedAnnotation, error) {
	ta, err := args.client.CreateAnnotation(
		context.Background(),
		&pb.NewTaggedAnnotation{
			Data: &pb.NewTaggedAnnotation_Data{
				Attributes: &pb.NewTaggedAnnotationAttributes{
					Value:     args.value,
					CreatedBy: regs.DEFAULT_USER,
					Tag:       args.tag,
					EntryId:   args.id,
					Ontology:  args.ontology,
					Rank:      int64(args.rank),
				},
			},
		},
	)
	if err != nil {
		return ta, fmt.Errorf(
			"error in creating annotation %s for id %s with rank %d %s",
			args.tag,
			args.id,
			args.rank,
			err,
		)
	}
	return ta, nil
}

func createAnno(args *createAnnoArgs) error {
	fmt.Printf("args %+v\n", args)
	_, err := args.client.CreateAnnotation(
		context.Background(),
		&pb.NewTaggedAnnotation{
			Data: &pb.NewTaggedAnnotation_Data{
				Attributes: &pb.NewTaggedAnnotationAttributes{
					Value:     args.value,
					CreatedBy: args.user,
					Tag:       args.tag,
					EntryId:   args.id,
					Ontology:  args.ontology,
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf(
			"error in creating annotation %s for id %s with ontology %s %s",
			args.tag,
			args.id,
			args.ontology,
			err,
		)
	}
	return nil
}

func findOrCreateAnnoWithStatus(args *createAnnoArgs) (bool, error) {
	create := false
	var errVal error
	_, err := args.client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      args.tag,
			EntryId:  args.id,
			Ontology: args.ontology,
		})
	switch {
	case err == nil:
		errVal = nil
	case status.Code(err) == codes.NotFound:
		err = createAnno(&createAnnoArgs{
			value:    args.value,
			user:     regs.DEFAULT_USER,
			tag:      args.tag,
			id:       args.id,
			ontology: args.ontology,
			client:   args.client,
		})
		if err != nil {
			errVal = err
		} else {
			create = true
		}
	case err != nil:
		errVal = fmt.Errorf(
			"error in finding annotation %s for id %s %s",
			args.tag,
			args.id,
			err,
		)
	}
	return create, errVal
}

func findOrCreateAnno(args *createAnnoArgs) (*pb.TaggedAnnotation, error) {
	ta, err := args.client.GetEntryAnnotation(
		context.Background(),
		&pb.EntryAnnotationRequest{
			Tag:      args.tag,
			EntryId:  args.id,
			Ontology: args.ontology,
		})
	switch {
	case err == nil:
		return ta, nil
	case status.Code(err) == codes.NotFound:
		return args.client.CreateAnnotation(
			context.Background(),
			&pb.NewTaggedAnnotation{
				Data: &pb.NewTaggedAnnotation_Data{
					Attributes: &pb.NewTaggedAnnotationAttributes{
						Value:     args.value,
						CreatedBy: regs.DEFAULT_USER,
						Tag:       args.tag,
						EntryId:   args.id,
						Ontology:  args.ontology,
					},
				},
			},
		)
	}
	return ta, fmt.Errorf(
		"error in finding annotation %s for id %s %s",
		args.tag,
		args.id,
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

func handleAnnoRetrieval(args *annoParams) (bool, error) {
	found := true
	if args.err != nil {
		if status.Code(args.err) != codes.NotFound { // error in lookup
			return found, fmt.Errorf("error in getting %s of %s %s", args.loader, args.id, args.err)
		}
		found = false
		args.logger.WithFields(logrus.Fields{
			"event": "get",
			"id":    args.id,
		}).Debugf("no %s", args.loader)
		return found, nil
	}
	args.logger.WithFields(logrus.Fields{
		"event": "get",
		"id":    args.id,
	}).Debugf("retrieved %s", args.loader)
	if err := delAnnotationGroup(args.client, args.gc); err != nil {
		return found, err
	}
	args.logger.WithFields(logrus.Fields{
		"event": "delete",
		"id":    args.id,
	}).Debugf("deleted %s", args.loader)
	return found, nil
}

// waitForPipeline waits for results from all error channels.
// It returns early on the first error.
func waitForPipeline(errs ...<-chan error) error {
	for err := range mergeErrors(errs...) {
		if err != nil {
			return err
		}
	}
	return nil
}

// mergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
func mergeErrors(cs ...<-chan error) <-chan error {
	// We must ensure that the output channel has the capacity to
	// hold as many errors
	// as there are error channels.
	// This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls
	// wg.Done.
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan error) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	// Start a goroutine to close out once all the output goroutines
	// are done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
