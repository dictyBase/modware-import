This file is a merged representation of a subset of the codebase, containing specifically included files, combined into a single document by Repomix.
The content has been processed where comments have been removed.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Only files matching these patterns are included: **/*.go
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Code comments have been removed from supported file types
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
cmd/
  modware-annotation/
    main.go
internal/
  app/
    server/
      server_feature.go
      server.go
    service/
      delete_service.go
      feature_annotation_test_helpers.go
      feature_annotation_test.go
      feature_annotation.go
      read_service.go
      service.go
      write_service.go
    validate/
      validate.go
  collection/
    collection_test.go
    collection.go
  message/
    nats/
      feature_annotation.go
      nats.go
    message.go
  model/
    feature_annotation.go
    model.go
    organism.go
  repository/
    arangodb/
      annotation_delete_test.go
      annotation_delete.go
      annotation_read_test.go
      annotation_read.go
      annotation_statement_test_helpers.go
      annotation_statement_test.go
      annotation_statement.go
      annotation_write_test.go
      annotation_write.go
      arangodb_test.go
      arangodb.go
      feature_annotation_helpers.go
      feature_annotation_pipeline.go
      feature_annotation_test_helpers.go
      feature_annotation_test.go
      feature_annotation.go
      feature_statement.go
      field.go
      list_filter_statement.go
      ontology.go
      organism_test_helpers.go
      organism_test.go
      organism.go
      pairwise.go
      parameters.go
      statement.go
    error.go
    feature_annotation.go
    organism.go
    repository.go
```

# Files

## File: internal/app/server/server.go
```go
package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
	"github.com/dictyBase/modware-annotation/internal/app/service"
	"github.com/dictyBase/modware-annotation/internal/message"
	"github.com/dictyBase/modware-annotation/internal/message/nats"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/dictyBase/modware-annotation/internal/repository/arangodb"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	gnats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	errCode  = 2
	waitTime = 2
)

type serverParams struct {
	repo repository.TaggedAnnotationRepository
	msg  message.Publisher
}

func RunServer(clt *cli.Context) error {
	spn, err := repoAndNatsConn(clt)
	if err != nil {
		return cli.NewExitError(err.Error(), errCode)
	}
	grpcS := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(getLogger(clt)),
		),
	)
	srv, err := service.NewAnnotationService(
		&service.Params{
			Repository: spn.repo,
			Publisher:  spn.msg,
			Group:      "groups",
			Options:    getGrpcOpt(),
		})
	if err != nil {
		return cli.NewExitError(err.Error(), errCode)
	}
	annotation.RegisterTaggedAnnotationServiceServer(grpcS, srv)
	reflection.Register(grpcS)

	endP := fmt.Sprintf(":%s", clt.String("port"))
	lis, err := net.Listen("tcp", endP)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("failed to listen %s", err), errCode,
		)
	}
	log.Printf("starting grpc server on %s", endP)
	if err := grpcS.Serve(lis); err != nil {
		return cli.NewExitError(err.Error(), errCode)
	}

	return nil
}

func getLogger(clt *cli.Context) *logrus.Entry {
	log := logrus.New()
	log.Out = os.Stderr
	switch clt.GlobalString("log-format") {
	case "text":
		log.Formatter = &logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	case "json":
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	}
	l := clt.GlobalString("log-level")
	switch l {
	case "debug":
		log.Level = logrus.DebugLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "error":
		log.Level = logrus.ErrorLevel
	case "fatal":
		log.Level = logrus.FatalLevel
	case "panic":
		log.Level = logrus.PanicLevel
	}

	return logrus.NewEntry(log)
}

func allParams(
	clt *cli.Context,
) (*manager.ConnectParams, *arangodb.CollectionParams, *ontoarango.CollectionParams) {
	arPort, _ := strconv.Atoi(clt.String("arangodb-port"))

	return &manager.ConnectParams{
			User:     clt.String("arangodb-user"),
			Pass:     clt.String("arangodb-pass"),
			Database: clt.String("arangodb-database"),
			Host:     clt.String("arangodb-host"),
			Port:     arPort,
			Istls:    clt.Bool("is-secure"),
		}, &arangodb.CollectionParams{
			Annotation:   clt.String("anno-collection"),
			AnnoTerm:     clt.String("annoterm-collection"),
			AnnoVersion:  clt.String("annover-collection"),
			AnnoTagGraph: clt.String("annoterm-graph"),
			AnnoVerGraph: clt.String("annover-graph"),
			AnnoGroup:    clt.String("annogroup-collection"),
			AnnoIndexes:  clt.StringSlice("annotation-index-fields"),
		}, &ontoarango.CollectionParams{
			GraphInfo:    clt.String("cv-collection"),
			OboGraph:     clt.String("obograph"),
			Relationship: clt.String("rel-collection"),
			Term:         clt.String("term-collection"),
		}
}

func getGrpcOpt() []aphgrpc.Option {
	return []aphgrpc.Option{
		aphgrpc.TopicsOption(map[string]string{
			"annotationCreate": "AnnotationService.Create",
			"annotationDelete": "AnnotationService.Delete",
			"annotationUpdate": "AnnotationService.Update",
		}),
	}
}

func repoAndNatsConn(clt *cli.Context) (*serverParams, error) {
	anrepo, err := arangodb.NewTaggedAnnotationRepo(allParams(clt))
	if err != nil {
		return &serverParams{},
			fmt.Errorf(
				"cannot connect to arangodb annotation repository %s",
				err,
			)
	}
	msp, err := nats.NewPublisher(
		clt.String("nats-host"), clt.String("nats-port"),
		gnats.MaxReconnects(-1), gnats.ReconnectWait(waitTime*time.Second),
	)
	if err != nil {
		return &serverParams{},
			fmt.Errorf("cannot connect to messaging server %s", err)
	}

	return &serverParams{
		repo: anrepo,
		msg:  msp,
	}, nil
}
```

## File: internal/app/service/delete_service.go
```go
package service

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/repository"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func (s *AnnotationService) DeleteAnnotationGroup(ctx context.Context, r *annotation.GroupEntryId) (*empty.Empty, error) {
	if err := protovalidate.Validate(r); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	if err := s.repo.RemoveAnnotationGroup(r.GroupId); err != nil {
		return nil, aphgrpc.HandleDeleteError(ctx, err)
	}

	return &empty.Empty{}, nil
}

func (s *AnnotationService) DeleteAnnotation(ctx context.Context, r *annotation.DeleteAnnotationRequest) (*empty.Empty, error) {
	if err := protovalidate.Validate(r); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	if err := s.repo.RemoveAnnotation(r.Id, r.Purge); err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleDeleteError(ctx, err)
	}

	return &empty.Empty{}, nil
}
```

## File: internal/app/service/service.go
```go
package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/arangomanager/query"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-genproto/dictybaseapis/api/upload"
	"github.com/dictyBase/go-obograph/storage"
	"github.com/dictyBase/modware-annotation/internal/message"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/dictyBase/modware-annotation/internal/repository/arangodb"
	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
)

const dividerVal = 1000000

type oboStreamHandler struct {
	writer *io.PipeWriter
	stream annotation.TaggedAnnotationService_OboJSONFileUploadServer
}


func (oh *oboStreamHandler) Write() error {
	defer oh.writer.Close()
	for {
		req, err := oh.stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("error in handling stream %s", err)
		}
		if _, err := oh.writer.Write(req.Content); err != nil {
			return fmt.Errorf(
				"error in writing the content from request %s",
				err,
			)
		}
	}

	return nil
}



type AnnotationService struct {
	*aphgrpc.Service
	repo      repository.TaggedAnnotationRepository
	publisher message.Publisher
	group     string
	annotation.UnimplementedTaggedAnnotationServiceServer
}


type Params struct {
	Repository repository.TaggedAnnotationRepository `validate:"required"`
	Publisher  message.Publisher                     `validate:"required"`
	Options    []aphgrpc.Option                      `validate:"required"`
	Group      string                                `validate:"required"`
}

func defaultOptions() *aphgrpc.ServiceOptions {
	return &aphgrpc.ServiceOptions{Resource: "annotations"}
}


func NewAnnotationService(srvP *Params) (*AnnotationService, error) {
	if err := validator.New().Struct(srvP); err != nil {
		return &AnnotationService{}, fmt.Errorf(
			"error in validating struct %s",
			err,
		)
	}
	so := defaultOptions()
	for _, optfn := range srvP.Options {
		optfn(so)
	}
	srv := &aphgrpc.Service{}
	aphgrpc.AssignFieldsToStructs(so, srv)

	return &AnnotationService{
		Service:   srv,
		repo:      srvP.Repository,
		publisher: srvP.Publisher,
		group:     srvP.Group,
	}, nil
}

func (s *AnnotationService) GetGroupResourceName() string {
	return s.group
}


func (s *AnnotationService) OboJSONFileUpload(
	stream annotation.TaggedAnnotationService_OboJSONFileUploadServer,
) error {
	in, out := io.Pipe()
	grp := new(errgroup.Group)
	defer in.Close()
	oh := &oboStreamHandler{writer: out, stream: stream}
	grp.Go(oh.Write)
	info, err := s.repo.LoadOboJSON(in)
	if err != nil {
		return aphgrpc.HandleGenericError(
			context.Background(),
			fmt.Errorf("error with loading obo %s", err),
		)
	}
	if err := grp.Wait(); err != nil {
		return aphgrpc.HandleGenericError(
			context.Background(),
			fmt.Errorf("error in waiting for the write to finish %s", err),
		)
	}

	err = stream.SendAndClose(&upload.FileUploadResponse{
		Status: uploadResponse(info),
		Msg:    "obojson file is uploaded",
	})
	if err != nil {
		return fmt.Errorf("error in closing the stream %s", err)
	}

	return nil
}

func uploadResponse(
	info *storage.UploadInformation,
) upload.FileUploadResponse_Status {
	if info.IsCreated {
		return upload.FileUploadResponse_CREATED
	}

	return upload.FileUploadResponse_UPDATED
}



func genNextCursorVal(t time.Time) int64 {
	return t.UnixNano() / dividerVal
}

func getAnnoAttributes(
	annom *model.AnnoDoc,
) *annotation.TaggedAnnotationAttributes {
	return &annotation.TaggedAnnotationAttributes{
		Value:         annom.Value,
		EditableValue: annom.EditableValue,
		CreatedBy:     annom.CreatedBy,
		CreatedAt:     aphgrpc.TimestampProto(annom.CreatedAt),
		Version:       annom.Version,
		EntryId:       annom.EnrtyId,
		Rank:          annom.Rank,
		IsObsolete:    annom.IsObsolete,
		Tag:           annom.Tag,
		Ontology:      annom.Ontology,
	}
}

func filterStrToQuery(filter string) (string, error) {
	var empty string
	if len(filter) == 0 {
		return empty, nil
	}
	p, err := query.ParseFilterString(filter)
	if err != nil {
		return empty, fmt.Errorf("error in parsing filter string")
	}
	q, err := query.GenQualifiedAQLFilterStatement(arangodb.FilterMap(), p)
	if err != nil {
		return empty, fmt.Errorf("error in generating aql statement")
	}

	return q, nil
}
```

## File: internal/app/service/write_service.go
```go
package service

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/repository"
)

func (s *AnnotationService) UpdateAnnotation(
	ctx context.Context,
	rta *annotation.TaggedAnnotationUpdate,
) (*annotation.TaggedAnnotation, error) {
	if err := protovalidate.Validate(rta); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	tga := &annotation.TaggedAnnotation{}
	mde, err := s.repo.EditAnnotation(rta)
	if err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}
	if mde.NotFound {
		return nil, aphgrpc.HandleNotFoundError(ctx, err)
	}
	tga.Data = s.getAnnoData(mde)
	err = s.publisher.Publish(s.Topics["annotationUpdate"], tga)
	if err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}

	return tga, nil
}

func (s *AnnotationService) CreateAnnotation(
	ctx context.Context,
	rta *annotation.NewTaggedAnnotation,
) (*annotation.TaggedAnnotation, error) {
	if err := protovalidate.Validate(rta); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	tga := &annotation.TaggedAnnotation{}
	m, err := s.repo.AddAnnotation(rta)
	if err != nil {
		return nil, aphgrpc.HandleInsertError(ctx, err)
	}
	tga.Data = s.getAnnoData(m)
	err = s.publisher.Publish(s.Topics["annotationCreate"], tga)
	if err != nil {
		return nil, aphgrpc.HandleInsertError(ctx, err)
	}

	return tga, nil
}

func (s *AnnotationService) AddToAnnotationGroup(
	ctx context.Context, rta *annotation.AnnotationGroupId,
) (*annotation.TaggedAnnotationGroup, error) {
	if err := protovalidate.Validate(rta); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	mga, err := s.repo.AppendToAnnotationGroup(rta.GroupId, rta.Id)
	if err != nil {
		if repository.IsGroupNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}

	return s.getGroup(mga), nil
}

func (s *AnnotationService) CreateAnnotationGroup(
	ctx context.Context, rta *annotation.AnnotationIdList,
) (*annotation.TaggedAnnotationGroup, error) {
	if err := protovalidate.Validate(rta); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	mga, err := s.repo.AddAnnotationGroup(rta.Ids...)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleInsertError(ctx, err)
	}

	return s.getGroup(mga), nil
}
```

## File: internal/app/validate/validate.go
```go
package validate

import (
	"fmt"

	"github.com/urfave/cli"
)

const errNo = 2

func ServerArgs(clt *cli.Context) error {
	for _, param := range []string{
		"arangodb-pass",
		"arangodb-database",
		"arangodb-user",
		"nats-host",
		"nats-port",
	} {
		if len(clt.String(param)) == 0 {
			return cli.NewExitError(
				fmt.Sprintf("argument %s is missing", param),
				errNo,
			)
		}
	}

	return nil
}
```

## File: internal/collection/collection_test.go
```go
package collection

import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapSeq(t *testing.T) {
	t.Parallel()
	assert := require.New(t)


	t.Run("convert ints to strings", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3, 4, 5}
		resultSeq := MapSeq(slices.Values(input), strconv.Itoa)
		expected := []string{"1", "2", "3", "4", "5"}
		assert.ElementsMatch(expected, slices.Collect(resultSeq))
	})


	t.Run("double and convert to string", func(t *testing.T) {
		t.Parallel()
		input := []int{10, 20, 30}
		resultSeq := MapSeq(slices.Values(input), func(n int) string {
			return strconv.Itoa(n * 2)
		})
		expected := []string{"20", "40", "60"}
		assert.ElementsMatch(expected, slices.Collect(resultSeq))
	})


	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		input := []int{}
		resultSeq := MapSeq(slices.Values(input), strconv.Itoa)
		assert.ElementsMatch(input, slices.Collect(resultSeq))
	})


	t.Run("conditional transformation", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3, 4, 5}
		resultSeq := MapSeq(slices.Values(input), func(n int) string {
			if n%2 == 0 {
				return "even"
			}

			return "odd"
		})
		expected := []string{"odd", "even", "odd", "even", "odd"}
		assert.ElementsMatch(expected, slices.Collect(resultSeq))
	})
}

func TestPartition(t *testing.T) {
	t.Parallel()
	assert := require.New(t)


	t.Run("partition integers by even/odd", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		isEven := func(n int) bool { return n%2 == 0 }
		evens, odds := Partition(input, isEven)
		assert.ElementsMatch([]int{2, 4, 6, 8, 10}, evens)
		assert.ElementsMatch([]int{1, 3, 5, 7, 9}, odds)
	})


	t.Run("partition strings by length", func(t *testing.T) {
		t.Parallel()
		input := []string{"a", "bb", "ccc", "dddd", "eeeee"}
		isShort := func(s string) bool { return len(s) <= 2 }

		short, long := Partition(input, isShort)

		assert.ElementsMatch([]string{"a", "bb"}, short)
		assert.ElementsMatch([]string{"ccc", "dddd", "eeeee"}, long)
	})


	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		input := []int{}
		isEven := func(n int) bool { return n%2 == 0 }

		evens, odds := Partition(input, isEven)

		assert.Empty(evens)
		assert.Empty(odds)
	})


	t.Run("all elements satisfy predicate", func(t *testing.T) {
		t.Parallel()
		input := []int{2, 4, 6, 8, 10}
		isEven := func(n int) bool { return n%2 == 0 }

		evens, odds := Partition(input, isEven)

		assert.ElementsMatch(input, evens)
		assert.Empty(odds)
	})


	t.Run("no elements satisfy predicate", func(t *testing.T) {
		t.Parallel()
		input := []int{1, 3, 5, 7, 9}
		isEven := func(n int) bool { return n%2 == 0 }

		evens, odds := Partition(input, isEven)

		assert.Empty(evens)
		assert.ElementsMatch(input, odds)
	})
}
```

## File: internal/message/nats/feature_annotation.go
```go
package nats

import (
	"encoding/json"
	"fmt"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/message"
	gnats "github.com/nats-io/nats.go"
)

type featureAnnotationPublisher struct {
	conn *gnats.Conn
}

func NewFeatureAnnotationPublisher(
	host, port string,
	options ...gnats.Option,
) (message.FeatureAnnotationPublisher, error) {
	nconn, err := gnats.Connect(
		fmt.Sprintf("nats://%s:%s", host, port),
		options...)
	if err != nil {
		return &featureAnnotationPublisher{}, fmt.Errorf(
			"error in connecting to nats server %s",
			err,
		)
	}

	return &featureAnnotationPublisher{conn: nconn}, nil
}

func (fnp *featureAnnotationPublisher) Publish(
	subj string,
	fann *feature.FeatureAnnotation,
) error {
	data, err := json.Marshal(fann)
	if err != nil {
		return fmt.Errorf("error in marshaling feature annotation %s", err)
	}
	if err := fnp.conn.Publish(subj, data); err != nil {
		return fmt.Errorf("error in publishing through nats %s", err)
	}

	return nil
}

func (fnp *featureAnnotationPublisher) Close() error {
	fnp.conn.Close()

	return nil
}
```

## File: internal/message/nats/nats.go
```go
package nats

import (
	"encoding/json"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/message"
	gnats "github.com/nats-io/nats.go"
)

type natsPublisher struct {
	conn *gnats.Conn
}

func NewPublisher(
	host, port string,
	options ...gnats.Option,
) (message.Publisher, error) {
	nconn, err := gnats.Connect(
		fmt.Sprintf("nats://%s:%s", host, port),
		options...)
	if err != nil {
		return &natsPublisher{}, fmt.Errorf(
			"error in connecting to nats server %s",
			err,
		)
	}

	return &natsPublisher{conn: nconn}, nil
}

func (n *natsPublisher) Publish(
	subj string,
	ann *annotation.TaggedAnnotation,
) error {
	data, err := json.Marshal(ann)
	if err != nil {
		return fmt.Errorf("error in marshaling annotation %s", err)
	}
	if err := n.conn.Publish(subj, data); err != nil {
		return fmt.Errorf("error in publishing through nats %s", err)
	}

	return nil
}

func (n *natsPublisher) Close() error {
	n.conn.Close()

	return nil
}
```

## File: internal/message/message.go
```go
package message

import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
)


type Publisher interface {

	Publish(subject string, ann *annotation.TaggedAnnotation) error

	Close() error
}


type FeatureAnnotationPublisher interface {

	Publish(subject string, ann *feature.FeatureAnnotation) error

	Close() error
}
```

## File: internal/model/model.go
```go
package model

import (
	"errors"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"
)

type UploadStatus int

const (
	Created UploadStatus = iota
	Updated
	Failed
)

type AnnoTag struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	IsObsolete bool   `json:"is_obsolete"`
	Ontology   string `json:"ontology"`
}

type AnnoDoc struct {
	driver.DocumentMeta
	Value         string    `json:"value"`
	EditableValue string    `json:"editable_value"`
	CreatedBy     string    `json:"created_by"`
	EnrtyId       string    `json:"entry_id"`
	Rank          int64     `json:"rank"`
	IsObsolete    bool      `json:"is_obsolete"`
	Version       int64     `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	Ontology      string    `json:"ontology,omitempty"`
	Tag           string    `json:"tag,omitempty"`
	CvtId         string    `json:"cvtid,omitempty"`
	NotFound      bool
}

type AnnoGroup struct {
	AnnoDocs  []*AnnoDoc `json:"annotations"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	GroupId   string     `json:"group_id"`
}

type DbGroup struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Group     []string  `json:"group"`
	GroupId   string    `json:"_key,omitempty"`
}

func UniqueModel[T comparable](slice []T) []T {
	result := make([]T, 0)
	seen := make(map[T]bool)
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

func DocToIDs(ml []*AnnoDoc) []string {
	str := make([]string, 0)
	for _, m := range ml {
		str = append(str, m.Key)
	}

	return str
}

func ConvToModel(i interface{}) (*AnnoDoc, error) {
	cmap, isok := i.(map[string]interface{})
	if !isok {
		return &AnnoDoc{}, errors.New("error in typecasting")
	}
	adoc := &AnnoDoc{
		Value:         cmap["value"].(string),
		EditableValue: cmap["editable_value"].(string),
		CreatedBy:     cmap["created_by"].(string),
		EnrtyId:       cmap["entry_id"].(string),
		Rank:          int64(cmap["rank"].(float64)),
		IsObsolete:    cmap["is_obsolete"].(bool),
		Version:       int64(cmap["version"].(float64)),
	}
	dstr, isok := cmap["created_at"].(string)
	if !isok {
		return &AnnoDoc{}, errors.New("error in typecasting")
	}
	t, err := time.Parse(time.RFC3339, dstr)
	if err != nil {
		return adoc, fmt.Errorf("error in parsing time %s", err)
	}
	adoc.CreatedAt = t
	adoc.DocumentMeta.Key = cmap["_key"].(string)
	adoc.DocumentMeta.Rev = cmap["_rev"].(string)

	return adoc, nil
}
```

## File: internal/model/organism.go
```go
package model

import (
	"time"

	driver "github.com/arangodb/go-driver"
)

type OrganismDoc struct {
	driver.DocumentMeta
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
	Abbreviation string    `json:"abbreviation,omitempty"`
	CommonName   string    `json:"common_name,omitempty"`
	Species      string    `json:"species"`
	Genus        string    `json:"genus"`
	NotFound     bool
}


func Schema() []byte {
	return []byte(`{
        "type": "object",
        "properties": {
            "created_at": {
                "type": "string",
		"format": "date-time"
            },
            "updated_at": {
                "type": "string",
		"format": "date-time"
            },
            "created_by": {
                "type": "string",
                "format": "email"
            },
            "updated_by": {
                "type": "string",
                "format": "email"
            },
            "abbreviation": {
                "type": "string"
            },
            "common_name": {
                "type": "string"
            },
            "species": {
                "type": "string"
            },
            "genus": {
                "type": "string"
            }
        },
	"required": ["genus", "species"]
    }`)
}
```

## File: internal/repository/arangodb/annotation_delete_test.go
```go
package arangodb

import (
	"testing"

	"github.com/dictyBase/modware-annotation/internal/model"
)

func TestRemoveFromAnnotationGroup(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsList(9)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	ids := testModelMaptoID(mla, model2IdCallback)
	g, err := anrepo.AddAnnotationGroup(ids...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	ega, err := anrepo.RemoveFromAnnotationGroup(g.GroupId, ids[:5]...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.ElementsMatch(
		testModelMaptoID(g.AnnoDocs, model2IdCallback),
		ids,
		"should match no of documents",
	)
	assert.ElementsMatch(
		ids[5:],
		testModelMaptoID(ega.AnnoDocs, model2IdCallback),
		"expected identical annotation identifiers after removing from the group",
	)
}

func TestRemoveAnnotationGroup(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsList(7)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	ids := testModelMaptoID(mla, model2IdCallback)
	g, err := anrepo.AddAnnotationGroup(ids...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	err = anrepo.RemoveAnnotationGroup(g.GroupId)
	assert.NoErrorf(err, "expect no error, received %s", err)
	err = anrepo.RemoveAnnotationGroup(g.GroupId)
	assert.Errorf(err, "should return error")
	assert.Contains(
		err.Error(),
		"removing group",
		"should contain removing group phrase",
	)
}

func TestRemoveAnnotation(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	nta := newTestTaggedAnnotation()
	m, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	nta2 := newTestTaggedAnnotationWithParams("curation", "DDB_G0287317")
	mt2, err := anrepo.AddAnnotation(nta2)
	assert.NoErrorf(err, "expect no error, received %s", err)
	err = anrepo.RemoveAnnotation(m.Key, true)
	assert.NoErrorf(err, "expect no error, received %s", err)
	err = anrepo.RemoveAnnotation(mt2.Key, false)
	assert.NoErrorf(err, "expect no error, received %s", err)
	err = anrepo.RemoveAnnotation(mt2.Key, false)
	assert.Errorf(err, "should return error")
	assert.Contains(err.Error(), "obsolete", "should contain obsolete message")
}
```

## File: internal/repository/arangodb/annotation_delete.go
```go
package arangodb

import (
	"context"
	"errors"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)

func (ar *arangorepository) RemoveAnnotation(id string, purge bool) error {
	manno := &model.AnnoDoc{}
	_, err := ar.anno.annot.ReadDocument(context.Background(), id, manno)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return &repository.AnnoNotFoundError{Id: id}
		}

		return fmt.Errorf("error in reading document %s", err)
	}
	if manno.IsObsolete {
		return fmt.Errorf(
			"annotation with id %s has already been obsolete",
			manno.Key,
		)
	}
	if purge {
		if _, err := ar.anno.annot.RemoveDocument(context.Background(), manno.Key); err != nil {
			return fmt.Errorf(
				"unable to purge annotation with id %s %s",
				manno.Key,
				err,
			)
		}

		return nil
	}
	_, err = ar.anno.annot.UpdateDocument(
		context.Background(),
		manno.Key,
		map[string]interface{}{"is_obsolete": true},
	)
	if err != nil {
		return fmt.Errorf(
			"unable to remove annotation with id %s %s",
			manno.Key,
			err,
		)
	}

	return nil
}


func (ar *arangorepository) RemoveFromAnnotationGroup(
	groupID string,
	idslice ...string,
) (*model.AnnoGroup, error) {
	manno := &model.AnnoGroup{}
	if len(idslice) <= 1 {
		return manno, errors.New(
			"need at least more than one entry to form a group",
		)
	}

	isok, err := ar.anno.annog.DocumentExists(
		context.Background(), groupID,
	)
	if err != nil {
		return manno, fmt.Errorf(
			"error in checking for existence of group identifier %s %s",
			groupID, err,
		)
	}
	if !isok {
		return manno, &repository.GroupNotFoundError{Id: groupID}
	}

	dbg := &model.DbGroup{}
	_, err = ar.anno.annog.ReadDocument(
		context.Background(),
		groupID, dbg,
	)
	if err != nil {
		return manno, fmt.Errorf("error in retrieving the group %s", err)
	}
	nids := collection.RemoveStringItems(dbg.Group, idslice...)
	mla, err := ar.getAllAnnotations(nids...)
	if err != nil {
		return manno, err
	}

	res, err := ar.database.DoRun(
		annGroupUpd,
		map[string]interface{}{
			"@anno_group_collection": ar.anno.annog.Name(),
			"key":                    groupID,
			"group":                  nids,
		})
	if err != nil {
		return manno, fmt.Errorf(
			"error in removing group members with id %s %s",
			groupID, err,
		)
	}
	ndbg := &model.DbGroup{}
	if err := res.Read(ndbg); err != nil {
		return manno, fmt.Errorf("error in reading data into struct %s", err)
	}
	manno.CreatedAt = ndbg.CreatedAt
	manno.UpdatedAt = ndbg.UpdatedAt
	manno.GroupId = ndbg.GroupId
	manno.AnnoDocs = mla

	return manno, nil
}
```

## File: internal/repository/arangodb/annotation_read_test.go
```go
package arangodb

import (
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/stretchr/testify/require"
)

const (
	filterOne   = `entry_id==DDB_G0286429;tag==private note;ontology==dicty_annotation`
	filterTwo   = `entry_id==DDB_G0294491;tag==name description;ontology==dicty_annotation`
	filterThree = `entry_id==jumbo`
)


func TestListAnnoFilter(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)


	tal := newTestTaggedAnnotationsListForFiltering(20)
	for _, anno := range tal {
		_, err := anrepo.AddAnnotation(anno)
		assert.NoErrorf(
			err,
			"setup: expect no error adding annotation, received %s",
			err,
		)
	}

	var mla, ml2, ml4 []*model.AnnoDoc

	t.Run("FilterOneFirstPage", func(t *testing.T) {
		mla = testListAnnoFilterOneFirstPage(t, assert, anrepo)
	})

	t.Run("FilterOneSecondPage", func(t *testing.T) {
		ml2 = testListAnnoFilterOneSecondPage(t, assert, anrepo, mla)
	})

	t.Run("FilterOneThirdPage", func(t *testing.T) {
		testListAnnoFilterOneThirdPage(t, assert, anrepo, ml2)
	})

	t.Run("FilterTwoFirstPage", func(t *testing.T) {
		ml4 = testListAnnoFilterTwoFirstPage(t, assert, anrepo)
	})

	t.Run("FilterTwoSecondPage", func(t *testing.T) {
		testListAnnoFilterTwoSecondPage(t, assert, anrepo, ml4)
	})

	t.Run("FilterNotFound", func(t *testing.T) {
		testListAnnoFilterNotFound(t, assert, anrepo)
	})
}

func TestGetAnnotationByID(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	nta := newTestTaggedAnnotation()
	mann, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	nta2 := newTestTaggedAnnotationWithParams("curation", "DDB_G0287317")
	ml2, err := anrepo.AddAnnotation(nta2)
	assert.NoErrorf(err, "expect no error, received %s", err)
	eim, err := anrepo.GetAnnotationByID(mann.Key)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Equal(mann.EnrtyId, eim.EnrtyId, "should match entry identifier")
	assert.Equal(mann.Ontology, eim.Ontology, "should match ontology")
	assert.Equal(mann.Tag, eim.Tag, "should match tag")
	assert.Equal(mann.Key, eim.Key, "should match the identifier")
	assert.Equal(mann.Value, eim.Value, "should match the value")
	assert.True(
		mann.CreatedAt.Equal(eim.CreatedAt),
		"should match created time of annotation",
	)
	assert.Equal(mann.Rank, eim.Rank, "should match rank")

	em2, err := anrepo.GetAnnotationByID(ml2.Key)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Equal(ml2.EnrtyId, em2.EnrtyId, "should match entry identifier")

	nie, err := anrepo.GetAnnotationByID("9999999")
	assert.Errorf(err, "expected %s error, received nothing", err)
	assert.True(
		repository.IsAnnotationNotFound(err),
		"entry should not exist",
	)
	assert.True(nie.NotFound, "entry should not exist")
}

func TestGetAnnotationByEntry(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	nta := newTestTaggedAnnotation()
	_, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	nta2 := newTestTaggedAnnotationWithParams("curation", "DDB_G0287317")
	_, err = anrepo.AddAnnotation(nta2)
	assert.NoErrorf(err, "expect no error, received %s", err)
	mae, err := anrepo.GetAnnotationByEntry(&annotation.EntryAnnotationRequest{
		Tag:      nta.Data.Attributes.Tag,
		EntryId:  nta.Data.Attributes.EntryId,
		Ontology: nta.Data.Attributes.Ontology,
	})
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Equal(int64(0), mae.Rank, "should match rank 0")
	assert.Equal(
		mae.EnrtyId,
		nta.Data.Attributes.EntryId,
		"should match the entry id",
	)

	ml2, err := anrepo.GetAnnotationByEntry(&annotation.EntryAnnotationRequest{
		Tag:      nta2.Data.Attributes.Tag,
		EntryId:  nta2.Data.Attributes.EntryId,
		Ontology: nta2.Data.Attributes.Ontology,
	})
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Equal(
		ml2.EnrtyId,
		nta2.Data.Attributes.EntryId,
		"should match the entry id",
	)
	assert.Equal(ml2.Tag, nta2.Data.Attributes.Tag, "should match the tag")

	emt, err := anrepo.GetAnnotationByEntry(&annotation.EntryAnnotationRequest{
		Tag:      nta2.Data.Attributes.Tag,
		Ontology: nta2.Data.Attributes.Ontology,
		EntryId:  "DDB_G0277853",
	})
	assert.Errorf(err, "expect %s error, received nothing", err)
	assert.True(
		repository.IsAnnotationNotFound(err),
		"the entry should not exist",
	)
	assert.True(emt.NotFound, "the entry should not exist")
}


func TestAddAnnotation(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)

	firstAnno := newTestAnnoWithTagAndOnto("dicty_annotation", "curator")
	t.Run("SuccessFirst", func(t *testing.T) {
		testAddAnnotationSuccess(t, assert, anrepo, firstAnno)
	})
	t.Run("Duplicate", func(t *testing.T) {
		testAddAnnotationDuplicate(t, assert, anrepo, firstAnno)
	})
	t.Run("NonExistentTag", func(t *testing.T) {
		testAddAnnotationNonExistentTag(t, assert, anrepo, firstAnno)
	})
	t.Run("NonExistentOntology", func(t *testing.T) {
		testAddAnnotationNonExistentOntology(t, assert, anrepo)
	})
	t.Run("SuccessSecond", func(t *testing.T) {
		testAddAnnotationSuccessSecond(t, assert, anrepo)
	})

	t.Run("SuccessThird", func(t *testing.T) {
		testAddAnnotationSuccessThird(t, assert, anrepo)
	})
}

func TestGetAnnotationGroup(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsList(4)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	ids := testModelMaptoID(mla, model2IdCallback)
	g, err := anrepo.AddAnnotationGroup(ids...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	eg, err := anrepo.GetAnnotationGroup(g.GroupId)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.ElementsMatch(
		testModelMaptoID(g.AnnoDocs, model2IdCallback),
		testModelMaptoID(eg.AnnoDocs, model2IdCallback),
		"expected identical annotation identifiers in the list",
	)
}

func TestListAnnGrFilter(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsListForFiltering(20)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	j := 5
	for i := 0; j <= len(mla); i += 5 {
		ids := testModelMaptoID(mla[i:j], model2IdCallback)
		_, err := anrepo.AddAnnotationGroup(ids...)
		assert.NoErrorf(err, "expect no error, received %s", err)
		j += 5
	}
	filterOne := `FILTER ann.entry_id == 'DDB_G0286429'
				  AND cvt.label == 'private note'
				  AND cv.metadata.namespace == 'dicty_annotation'
	`
	egl, err := anrepo.ListAnnotationGroup(0, 10, filterOne)
	assert.NoErrorf(err, "expect no error, received %s", err)
	testGroupMember(t, egl, 2, 0, "sidd@gmail.com")
	filterTwo := `FILTER ann.entry_id == 'DDB_G0294491'
				  AND cvt.label == 'name description'
				  AND cv.metadata.namespace == 'dicty_annotation'
	`
	egl2, err := anrepo.ListAnnotationGroup(0, 10, filterTwo)
	assert.NoErrorf(err, "expect no error, received %s", err)
	testGroupMember(t, egl2, 2, 1, "basu@gmail.com")
	filterThree := `FILTER cv.metadata.namespace == 'dicty_annotation'`
	egl3, err := anrepo.ListAnnotationGroup(0, 2, filterThree)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(egl3, 2, "should have two groups")
	for _, g := range egl3 {
		assert.Len(g.AnnoDocs, 5, "should have 5 annotations in each group")
	}
	egl4, err := anrepo.ListAnnotationGroup(
		toTimestamp(egl3[len(egl3)-1].CreatedAt),
		4,
		filterThree,
	)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(egl4, 3, "should have three groups")
	for _, g := range egl4 {
		assert.Len(g.AnnoDocs, 5, "should have 5 annotations in each group")
	}
	_, err = anrepo.ListAnnotationGroup(0, 4, "FILTER ann.entry_id == 'jumbo'")
	assert.Error(err, "expect error")
	assert.True(
		repository.IsAnnotationGroupListNotFound(err),
		"expect no annotation group to be found",
	)
}

func TestListAnnotationGroup(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsList(60)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	j := 5
	for i := 0; j <= len(mla); i += 5 {
		ids := testModelMaptoID(mla[i:j], model2IdCallback)
		_, err := anrepo.AddAnnotationGroup(ids...)
		assert.NoErrorf(err, "expect no error, received %s", err)
		j += 5
	}
	egl, err := anrepo.ListAnnotationGroup(0, 4, "")
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(egl, 4, "should have 4 groups")
	for _, g := range egl {
		assert.Len(g.AnnoDocs, 5, "should have 5 annotations in each group")
	}
	egl2, err := anrepo.ListAnnotationGroup(
		toTimestamp(egl[len(egl)-1].CreatedAt),
		6,
		"",
	)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(egl2, 6, "should have 6 groups")
	for _, g := range egl2 {
		assert.Len(g.AnnoDocs, 5, "should have 5 annotations in each group")
	}
	assert.Exactly(
		egl[len(egl)-1],
		egl2[0],
		"should have identical model objects",
	)
	egl3, err := anrepo.ListAnnotationGroup(
		toTimestamp(egl2[len(egl2)-1].CreatedAt),
		6,
		"",
	)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(egl3, 4, "should have 4 groups")
	for _, g := range egl3 {
		assert.Len(g.AnnoDocs, 5, "should have 5 annotations in each group")
	}
	assert.Exactly(
		egl2[len(egl2)-1],
		egl3[0],
		"should have identical model objects",
	)
}

func TestGetAnnotationTag(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	for _, tag := range tags[:6] {
		m, err := anrepo.GetAnnotationTag(tag, "dicty_annotation")
		assert.NoErrorf(err, "expect no error from fetching %s tag", tag)
		assert.Equal(m.Name, tag, "should match tag name")
		assert.Equal("dicty_annotation", m.Ontology, "should match ontology")
		assert.Falsef(m.IsObsolete, "tag %s should not be obsolete", tag)
	}
	_, err := anrepo.GetAnnotationTag("yadayada", "dicty_annotation")
	assert.Error(err, "expect error from non-existent tag")
	assert.True(
		repository.IsAnnoTagNotFound(err),
		"should be an error for non-existent tag",
	)
}

func testListAnnoFilterOneFirstPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) []*model.AnnoDoc {
	t.Helper()
	mla, err := anrepo.ListAnnotations(
		&repository.ListAnnotationsParams{Limit: 4, Filter: filterOne},
	)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(mla, 5, "should have 5 annotations")
	for _, m := range mla {
		assert.Equal("sidd@gmail.com", m.CreatedBy, "should match created by")
		assert.Equal(m.Tag, tags[0], "should match the tag")
		assert.Equal(m.EnrtyId, ddbg[0], "should match the entry id")
	}
	testModelListSort(t, mla)

	return mla
}

func testListAnnoFilterOneSecondPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	prevResult []*model.AnnoDoc,
) []*model.AnnoDoc {
	t.Helper()
	assert.NotEmpty(prevResult, "previous result should not be empty")
	ml2, err := anrepo.ListAnnotations(&repository.ListAnnotationsParams{
		Cursor: toTimestamp(prevResult[len(prevResult)-1].CreatedAt),
		Limit:  4, Filter: filterOne,
	})
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(ml2, 5, "should have five annotations")
	assert.Exactly(
		prevResult[len(prevResult)-1],
		ml2[0],
		"should have identical model objects",
	)
	testModelListSort(t, ml2)

	return ml2
}

func testListAnnoFilterOneThirdPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	prevResult []*model.AnnoDoc,
) {
	t.Helper()
	assert.NotEmpty(prevResult, "previous result should not be empty")
	ml3, err := anrepo.ListAnnotations(&repository.ListAnnotationsParams{
		Cursor: toTimestamp(prevResult[len(prevResult)-1].CreatedAt),
		Limit:  4, Filter: filterOne,
	})
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(ml3, 2, "should have two annotations")
	assert.Exactly(
		prevResult[len(prevResult)-1],
		ml3[0],
		"should have identical model objects",
	)
	testModelListSort(t, ml3)
}

func testListAnnoFilterTwoFirstPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) []*model.AnnoDoc {
	t.Helper()
	ml4, err := anrepo.ListAnnotations(
		&repository.ListAnnotationsParams{Limit: 6, Filter: filterTwo},
	)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(ml4, 7, "should have 7 annotations")
	for _, m := range ml4 {
		assert.Equal("basu@gmail.com", m.CreatedBy, "should match created by")
		assert.Equal(m.Tag, tags[1], "should match the tag")
		assert.Equal(m.EnrtyId, ddbg[1], "should match the entry id")
	}
	testModelListSort(t, ml4)

	return ml4
}

func testListAnnoFilterTwoSecondPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	prevResult []*model.AnnoDoc,
) {
	t.Helper()
	assert.NotEmpty(prevResult, "previous result should not be empty")
	ml5, err := anrepo.ListAnnotations(&repository.ListAnnotationsParams{
		Cursor: toTimestamp(prevResult[len(prevResult)-1].CreatedAt),
		Limit:  4, Filter: filterTwo,
	})
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Len(ml5, 4, "should have four annotations")
	assert.Exactly(
		prevResult[len(prevResult)-1],
		ml5[0],
		"should have identical model objects",
	)
	testModelListSort(t, ml5)
}

func testListAnnoFilterNotFound(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) {
	t.Helper()
	_, err := anrepo.ListAnnotations(
		&repository.ListAnnotationsParams{Limit: 4, Filter: filterThree},
	)
	assert.Error(err, "expect error")
	assert.True(
		repository.IsAnnotationListNotFound(err),
		"expect no annotation list found",
	)
}


func testAddAnnotationSuccess(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	nta *annotation.NewTaggedAnnotation,
) {
	t.Helper()
	mann, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.False(
		mann.IsObsolete,
		"new tagged annotation should not be obsolete",
	)
	assert.Equal(
		nta.Data.Attributes.Value,
		mann.Value,
		"should match the value",
	)
	assert.Equal(
		nta.Data.Attributes.CreatedBy,
		mann.CreatedBy,
		"should match created_by",
	)
	assert.Equal(
		nta.Data.Attributes.EntryId,
		mann.EnrtyId,
		"should match entry identifier",
	)
	assert.Equal(nta.Data.Attributes.Rank, mann.Rank, "should match the rank")
	assert.Equal(
		nta.Data.Attributes.Ontology,
		mann.Ontology,
		"should match ontology name",
	)
	assert.Equal(
		nta.Data.Attributes.Tag,
		mann.Tag,
		"should match the ontology tag",
	)
}

func testAddAnnotationDuplicate(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	nta *annotation.NewTaggedAnnotation,
) {
	t.Helper()
	_, err := anrepo.AddAnnotation(nta)
	assert.Error(err, "expect error for existing annotation")
	assert.Regexp(
		"already exists",
		err.Error(), "error should have existence of annotation",
	)
}

func testAddAnnotationNonExistentTag(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	nta *annotation.NewTaggedAnnotation,
) {
	t.Helper()

	ntaCopy := &annotation.NewTaggedAnnotation{
		Data: &annotation.NewTaggedAnnotation_Data{
			Attributes: &annotation.NewTaggedAnnotationAttributes{
				Value:     nta.Data.Attributes.Value,
				CreatedBy: nta.Data.Attributes.CreatedBy,
				EntryId:   nta.Data.Attributes.EntryId,
				Rank:      nta.Data.Attributes.Rank,
				Ontology:  nta.Data.Attributes.Ontology,
				Tag:       "respiration",
			},
		},
	}
	_, err := anrepo.AddAnnotation(ntaCopy)
	assert.Error(err, "expect error in case of non-existent ontology and tag")
	assert.Regexp(
		"respiration",
		err.Error(), "error should contain the non-existent tag name",
	)
}

func testAddAnnotationNonExistentOntology(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) {
	t.Helper()
	nta := newTestAnnoWithTagAndOnto("caboose", "description")
	_, err := anrepo.AddAnnotation(nta)
	assert.Error(err, "expect error in case of non-existent ontology and tag")
	assert.Regexp(
		"caboose",
		err.Error(), "error should contain the non-existent ontology",
	)
}

func testAddAnnotationSuccessSecond(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) {
	t.Helper()
	nta := newTestAnnoWithTagAndOnto("dicty_annotation", "summary")
	mann2, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.False(
		mann2.IsObsolete,
		"new tagged annotation should not be obsolete",
	)
	assert.Equal(
		nta.Data.Attributes.Value,
		mann2.Value,
		"should match the value",
	)
	assert.Equal(
		nta.Data.Attributes.CreatedBy,
		mann2.CreatedBy,
		"should match created_by",
	)
	assert.Equal(
		nta.Data.Attributes.EntryId,
		mann2.EnrtyId,
		"should match entry identifier",
	)
	assert.Equal(nta.Data.Attributes.Rank, mann2.Rank, "should match the rank")
	assert.Equal(
		nta.Data.Attributes.Ontology,
		mann2.Ontology,
		"should match ontology name",
	)

	assert.Equal("description", mann2.Tag, "should match the ontology tag")
}

func testAddAnnotationSuccessThird(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) {
	t.Helper()
	nta := newTestAnnoWithTagAndOnto(
		"dicty_annotation",
		"decreased 3',5'-cyclic-GMP phosphodiesterase activity",
	)
	annm, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Equal(
		nta.Data.Attributes.Ontology,
		annm.Ontology,
		"should match ontology name",
	)
	assert.Equal(nta.Data.Attributes.Tag, annm.Tag, "should match the tag")
}
```

## File: internal/repository/arangodb/annotation_statement_test.go
```go
package arangodb

import (
	"errors"
	"testing"

	"github.com/dictyBase/arangomanager/query"
	"github.com/stretchr/testify/require"
)

func TestStatementTemplate(t *testing.T) {
	t.Parallel()
	t.Run("both filters cases", func(t *testing.T) {
		t.Parallel()
		testBothFiltersStatementTemplate(t)
	})
	t.Run("first filter cases", func(t *testing.T) {
		t.Parallel()
		testFirstFilterStatementTemplate(t)
	})
	t.Run("second filter cases", func(t *testing.T) {
		t.Parallel()
		testSecondFilterStatementTemplate(t)
	})
	t.Run("invalid case", func(t *testing.T) {
		t.Parallel()
		testInvalidStatementTemplate(t)
	})
}

func TestFormatKey(t *testing.T) {
	t.Parallel()

	t.Run("both filters with cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		key := formatKey(BothFilters, true)
		assert.Equal("bothtrue", key, "should format key correctly")
	})

	t.Run("first filter without cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		key := formatKey(FirstFilter, false)
		assert.Equal("firstfalse", key, "should format key correctly")
	})

	t.Run("second filter with cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		key := formatKey(SecondFilter, true)
		assert.Equal("secondtrue", key, "should format key correctly")
	})
}

func TestDetermineStatementType(t *testing.T) {
	t.Parallel()

	t.Run("both filters", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		firstSet := []*query.Filter{
			{Field: "ann.value", Value: "test"},
		}
		secondSet := []*query.Filter{
			{Field: "cvt.label", Value: "test"},
		}
		stmtType, ok := determineStatementType(firstSet, secondSet)
		assert.True(ok, "should return true")
		assert.Equal(BothFilters, stmtType, "should return BothFilters type")
	})

	t.Run("first filter only", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		firstSet := []*query.Filter{
			{Field: "ann.value", Value: "test"},
		}
		secondSet := []*query.Filter{}
		stmtType, ok := determineStatementType(firstSet, secondSet)
		assert.True(ok, "should return true")
		assert.Equal(FirstFilter, stmtType, "should return FirstFilter type")
	})

	t.Run("second filter only", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		firstSet := []*query.Filter{}
		secondSet := []*query.Filter{
			{Field: "cvt.label", Value: "test"},
		}
		stmtType, ok := determineStatementType(firstSet, secondSet)
		assert.True(ok, "should return true")
		assert.Equal(SecondFilter, stmtType, "should return SecondFilter type")
	})

	t.Run("no filters", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		firstSet := []*query.Filter{}
		secondSet := []*query.Filter{}
		stmtType, ok := determineStatementType(firstSet, secondSet)
		assert.False(ok, "should return false")
		assert.Empty(stmtType, "should return empty statement type")
	})
}

func TestFilterAndPartitionFunc(t *testing.T) {
	t.Parallel()
	fmap := FilterMap()
	t.Run("with valid filters", func(t *testing.T) {
		t.Parallel()
		testValidFilters(t, fmap)
	})

	t.Run("with only annotation filters", func(t *testing.T) {
		t.Parallel()
		testOnlyAnnotationFilters(t, fmap)
	})

	t.Run("with only cvterm filters", func(t *testing.T) {
		t.Parallel()
		testOnlyCvtermFilters(t, fmap)
	})

	t.Run("with invalid filters", func(t *testing.T) {
		t.Parallel()
		testInvalidFilters(t, fmap)
	})

	t.Run("with mixed valid and invalid filters", func(t *testing.T) {
		t.Parallel()
		testMixedFilters(t, fmap)
	})

	t.Run("with existing error", func(t *testing.T) {
		t.Parallel()
		testExistingError(t)
	})
}

func TestParseFiltersFunc(t *testing.T) {
	t.Parallel()
	t.Run("success case", testParseFiltersFuncSuccess)
	t.Run(
		"failure case - invalid filter string",
		testParseFiltersFuncFailureInvalidString,
	)
	t.Run(
		"edge case - empty filter string",
		testParseFiltersFuncEdgeEmptyString,
	)
	t.Run("existing error case", testParseFiltersFuncExistingError)
}

func TestGenFilterStatement(t *testing.T) {
	t.Parallel()
	fmap := FilterMap()

	t.Run("success - single filter", func(t *testing.T) {
		t.Parallel()
		testGenFilterStatementSuccessSingle(t, fmap)
	})

	t.Run("success - multiple filters", func(t *testing.T) {
		t.Parallel()
		testGenFilterStatementSuccessMultiple(t, fmap)
	})

	t.Run("error - invalid filter field", func(t *testing.T) {
		t.Parallel()
		testGenFilterStatementErrorInvalidField(t, fmap)
	})

	t.Run("edge case - empty filters slice", func(t *testing.T) {
		t.Parallel()
		testGenFilterStatementEdgeEmpty(t, fmap)
	})
}

func TestBuildAQLStatementErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("with existing error", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		originalErr := errors.New("pre-existing error")
		ctx := FilterContext{Err: originalErr}
		result := buildAQLStatement(ctx)
		assert.Equal(
			originalErr,
			result.Err,
			"should return the existing error",
		)
		assert.Empty(
			result.Statement,
			"statement should be empty on existing error",
		)
	})

	t.Run("with invalid statement type", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{Type: StatementType("invalid")}
		result := buildAQLStatement(ctx)
		assert.Error(
			result.Err,
			"should return error for invalid statement type",
		)
		assert.Contains(
			result.Err.Error(),
			"no matching template found",
			"error message should indicate unsupported type",
		)
		assert.Empty(result.Statement, "statement should be empty")
	})
}

func TestBuildAQLStatementBothFilters(t *testing.T) {
	t.Parallel()
	fmap := FilterMap()

	t.Run("without cursor", func(t *testing.T) {
		t.Parallel()
		testBothFiltersWithoutCursor(t, fmap)
	})

	t.Run("with cursor", func(t *testing.T) {
		t.Parallel()
		testBothFiltersWithCursor(t, fmap)
	})
}

func TestBuildAQLStatementFirstFilter(t *testing.T) {
	t.Parallel()
	fmap := FilterMap()
	t.Run("without cursor", func(t *testing.T) {
		t.Parallel()
		testBuildAQLStatementFirstFilterWithoutCursor(t, fmap)
	})

	t.Run("with cursor", func(t *testing.T) {
		t.Parallel()
		testBuildAQLStatementFirstFilterWithCursor(t, fmap)
	})
}

func TestBuildAQLStatementSecondFilter(t *testing.T) {
	t.Parallel()
	fmap := FilterMap()
	t.Run("without cursor", func(t *testing.T) {
		t.Parallel()
		testBuildAQLStatementSecondFilterWithoutCursor(t, fmap)
	})
	t.Run("with cursor", func(t *testing.T) {
		t.Parallel()
		testBuildAQLStatementSecondFilterWithCursor(t, fmap)
	})
}

func TestBuildAQLStatementFilterGenerationErrors(t *testing.T) {
	t.Parallel()
	badFilterMap := map[string]string{"valid": "good"}

	t.Run("error generating both filters", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      BothFilters,
			HasCursor: false,
			FilterMap: badFilterMap,
			FirstSet:  []*query.Filter{createTestFilter("invalid", "val1")},
			SecondSet: []*query.Filter{createTestFilter("cvt.label", "tag1")},
		}
		result := buildAQLStatement(ctx)
		assert.Error(result.Err, "should return error from filter generation")
		assert.Contains(
			result.Err.Error(),
			"error generating annotation filter",
			"error message should indicate annotation filter error",
		)
	})

	t.Run("error generating first filter", func(t *testing.T) {
		assert := require.New(t)
		ctx := FilterContext{
			Type:      FirstFilter,
			HasCursor: false,
			FilterMap: badFilterMap,
			FirstSet:  []*query.Filter{createTestFilter("invalid", "val1")},
			SecondSet: []*query.Filter{},
		}
		result := buildAQLStatement(ctx)
		assert.Error(result.Err, "should return error from filter generation")
		assert.Contains(
			result.Err.Error(),
			"error generating annotation filter",
			"error message should indicate annotation filter error",
		)
	})

	t.Run("error generating second filter", func(t *testing.T) {
		assert := require.New(t)
		ctx := FilterContext{
			Type:      SecondFilter,
			HasCursor: false,
			FilterMap: badFilterMap,
			FirstSet:  []*query.Filter{},
			SecondSet: []*query.Filter{createTestFilter("invalid", "tag1")},
		}
		result := buildAQLStatement(ctx)
		assert.Error(result.Err, "should return error from filter generation")
		assert.Contains(
			result.Err.Error(),
			"error generating cvterm filter",
			"error message should indicate cvterm filter error",
		)
	})
}

func TestGetListAnnoStatement(t *testing.T) {
	t.Parallel()
	t.Run("basic cases", testGetListAnnoStatementBasicCases)
	t.Run("valid filters", testGetListAnnoStatementValidFilters)
	t.Run("tag filters", testGetListAnnoStatementTagFilters)
}
```

## File: internal/repository/arangodb/annotation_statement.go
```go
package arangodb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dictyBase/arangomanager/query"
	"github.com/dictyBase/modware-annotation/internal/collection"
)

const (

	BothFilters StatementType = "both"

	FirstFilter StatementType = "first"

	SecondFilter StatementType = "second"
)


type StatementType string



type PickStatementResult struct {
	Statement string
	Err       error
}


type FilterContext struct {
	FilterString string
	HasCursor    bool
	Filters      []*query.Filter
	FirstSet     []*query.Filter
	SecondSet    []*query.Filter
	FilterMap    map[string]string
	Type         StatementType
	Err          error
}


var templateMap = map[string]string{
	formatKey(BothFilters, true):   annCvtListFilterWithCursorQ,
	formatKey(BothFilters, false):  annCvtListFilterQ,
	formatKey(FirstFilter, true):   annExclusiveListFilterWithCursorQ,
	formatKey(FirstFilter, false):  annExclusiveListFilterQ,
	formatKey(SecondFilter, true):  cvtExclusiveListFilterWithCursorQ,
	formatKey(SecondFilter, false): cvtExclusiveListFilterQ,
}


func formatKey(statementType StatementType, hasCursor bool) string {
	return fmt.Sprintf("%s%v", string(statementType), hasCursor)
}


func statementTemplate(ctx FilterContext) (string, bool) {
	val, ok := templateMap[formatKey(ctx.Type, ctx.HasCursor)]

	return val, ok
}



func buildAQLStatement(ctx FilterContext) PickStatementResult {
	if ctx.Err != nil {
		return PickStatementResult{Err: ctx.Err}
	}

	var result PickStatementResult
	template, ok := statementTemplate(ctx)
	if !ok {
		result.Err = fmt.Errorf(
			"no matching template found for statement type %s with cursor=%v",
			ctx.Type,
			ctx.HasCursor,
		)

		return result
	}

	switch ctx.Type {
	case BothFilters:
		result = buildBothFiltersStatement(
			template,
			ctx.FilterMap,
			ctx.FirstSet,
			ctx.SecondSet,
		)
	case FirstFilter:
		result = buildFirstFilterStatement(
			template,
			ctx.FilterMap,
			ctx.FirstSet,
		)
	case SecondFilter:
		result = buildSecondFilterStatement(
			template,
			ctx.FilterMap,
			ctx.SecondSet,
		)
	default:
		result.Err = errors.New("unsupported statement type")
	}

	return result
}



func genFilterStatement(
	filterMap map[string]string,
	filters []*query.Filter,
	filterType string,
) (string, error) {
	filter, err := query.GenQualifiedAQLFilterStatement(filterMap, filters)
	if err != nil {
		return "", fmt.Errorf("error generating %s filter: %w", filterType, err)
	}

	return filter, nil
}



func buildBothFiltersStatement(
	template string,
	filterMap map[string]string,
	firstSet, secondSet []*query.Filter,
) PickStatementResult {
	var result PickStatementResult

	afilter, err := genFilterStatement(filterMap, firstSet, "annotation")
	if err != nil {
		result.Err = err

		return result
	}

	cfilter, err := genFilterStatement(filterMap, secondSet, "cvterm")
	if err != nil {
		result.Err = err

		return result
	}

	result.Statement = fmt.Sprintf(template, afilter, cfilter)

	return result
}



func buildFirstFilterStatement(
	template string,
	filterMap map[string]string,
	filters []*query.Filter,
) PickStatementResult {
	var result PickStatementResult

	afilter, err := genFilterStatement(filterMap, filters, "annotation")
	if err != nil {
		result.Err = err

		return result
	}

	result.Statement = fmt.Sprintf(template, afilter)

	return result
}



func buildSecondFilterStatement(
	template string,
	filterMap map[string]string,
	filters []*query.Filter,
) PickStatementResult {
	var result PickStatementResult

	cfilter, err := genFilterStatement(filterMap, filters, "cvterm")
	if err != nil {
		result.Err = err

		return result
	}

	result.Statement = fmt.Sprintf(template, cfilter)

	return result
}



func getListAnnoStatement(fstr string, cursor int64) PickStatementResult {
	if len(fstr) == 0 {
		return PickStatementResult{
			Err: errors.New("empty filter string"),
		}
	}

	return collection.Pipe4(
		FilterContext{
			FilterString: fstr,
			HasCursor:    cursor != 0,
			FilterMap:    FilterMap(),
		},
		parseFiltersFunc,
		filterAndPartitionFunc,
		determineStatementTypeFunc,
		buildAQLStatement,
	)
}


func parseFiltersFunc(ctx FilterContext) FilterContext {
	if ctx.Err != nil {
		return ctx
	}
	filters, err := query.ParseFilterString(ctx.FilterString)
	if err != nil {
		ctx.Err = fmt.Errorf(
			"error parsing filter string %q: %w",
			ctx.FilterString,
			err,
		)

		return ctx
	}
	ctx.Filters = filters

	return ctx
}


func filterAndPartitionFunc(ctx FilterContext) FilterContext {
	if ctx.Err != nil {
		return ctx
	}
	var validFilters []*query.Filter
	var firstSet []*query.Filter
	var secondSet []*query.Filter

	for _, qfl := range ctx.Filters {
		if mappedField, ok := ctx.FilterMap[qfl.Field]; ok {

			mappedFilter := &query.Filter{
				Field:    qfl.Field,
				Value:    qfl.Value,
				Operator: qfl.Operator,
			}
			if len(qfl.Logic) != 0 {
				mappedFilter.Logic = qfl.Logic
			}
			validFilters = append(validFilters, mappedFilter)

			if strings.HasPrefix(mappedField, "ann.") {
				firstSet = append(firstSet, mappedFilter)
			} else {
				secondSet = append(secondSet, mappedFilter)
			}
		}
	}

	if collection.IsEmpty(validFilters) {
		ctx.Err = fmt.Errorf(
			"no valid filters found in filter string %q",
			ctx.FilterString,
		)

		return ctx
	}


	ctx.Filters = unsetLogicIfSingleFilter(validFilters)
	ctx.FirstSet = unsetLogicIfSingleFilter(firstSet)
	ctx.SecondSet = unsetLogicIfSingleFilter(secondSet)

	return ctx
}




func unsetLogicIfSingleFilter(filters []*query.Filter) []*query.Filter {
	if len(filters) != 1 {
		return filters
	}

	newFilters := make([]*query.Filter, 1)
	newFilters[0] = &query.Filter{
		Field:    filters[0].Field,
		Value:    filters[0].Value,
		Operator: filters[0].Operator,
		Logic:    "", // Unset logic
	}

	return newFilters
}

// determineStatementTypeFunc returns a function for determining statement type in a pipeline.
func determineStatementTypeFunc(ctx FilterContext) FilterContext {
	if ctx.Err != nil {
		return ctx
	}
	statementType, ok := determineStatementType(ctx.FirstSet, ctx.SecondSet)
	if !ok {
		ctx.Err = errors.New(
			"no valid filters found after parsing",
		)

		return ctx
	}
	ctx.Type = statementType
	ctx.Filters = append(
		ctx.Filters,
		&query.Filter{Field: (string(statementType))},
	)

	return ctx
}



func determineStatementType(
	first, second []*query.Filter,
) (StatementType, bool) {
	switch {
	case !collection.IsEmpty(first) && !collection.IsEmpty(second):
		return BothFilters, true
	case !collection.IsEmpty(first) && collection.IsEmpty(second):
		return FirstFilter, true
	case collection.IsEmpty(first) && !collection.IsEmpty(second):
		return SecondFilter, true
	default:
		return "", false
	}
}
```

## File: internal/repository/arangodb/annotation_write_test.go
```go
package arangodb

import (
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
)

func TestEditAnnotation(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	nta := newTestTaggedAnnotation()
	mda, err := anrepo.AddAnnotation(nta)
	assert.NoErrorf(err, "expect no error, received %s", err)
	uan := &annotation.TaggedAnnotationUpdate{
		Data: &annotation.TaggedAnnotationUpdate_Data{
			Type: "annotations",
			Id:   mda.Key,
			Attributes: &annotation.TaggedAnnotationUpdateAttributes{
				Value:         "updated gene description",
				EditableValue: "updated gene description",
				CreatedBy:     "basu@gmail.com",
			},
		},
	}
	um, err := anrepo.EditAnnotation(uan)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Equal(mda.Version+1, um.Version, "version should be incremented by 1")
	assert.NotEqual(uan.Data.Id, um.Key, "identifier should not match")
	assert.Equal(uan.Data.Attributes.Value, um.Value, "should matches the value")
	assert.Equal(uan.Data.Attributes.CreatedBy, um.CreatedBy, "should matches created by")
}

func TestAddAnnotationGroup(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsList(8)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	ids := testModelMaptoID(mla, model2IdCallback)
	g, err := anrepo.AddAnnotationGroup(ids...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.Lenf(g.AnnoDocs, len(ids), "should have %d annotations", len(ids))
}

func TestAppendToAnntationGroup(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	tal := newTestTaggedAnnotationsList(7)
	mla := make([]*model.AnnoDoc, 0)
	for _, ann := range tal {
		m, err := anrepo.AddAnnotation(ann)
		assert.NoErrorf(err, "expect no error, received %s", err)
		mla = append(mla, m)
	}
	ids := testModelMaptoID(mla[:4], model2IdCallback)
	g, err := anrepo.AddAnnotationGroup(ids...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	nids := testModelMaptoID(mla[4:], model2IdCallback)
	eg, err := anrepo.AppendToAnnotationGroup(g.GroupId, nids...)
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.ElementsMatch(
		testModelMaptoID(eg.AnnoDocs, model2IdCallback),
		append(ids, nids...),
		"expected identical annotation identifiers after appending to the group",
	)
}
```

## File: internal/repository/arangodb/annotation_write.go
```go
package arangodb

import (
	"context"
	"errors"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)

const maxTransactionSize = 10000

func (ar *arangorepository) AddAnnotation(na *annotation.NewTaggedAnnotation) (*model.AnnoDoc, error) {
	mann := &model.AnnoDoc{}
	attr := na.Data.Attributes

	cvtid, err := ar.termID(attr.Ontology, attr.Tag)
	if err != nil {
		return mann, err
	}

	tag, err := ar.termName(cvtid)
	if err != nil {
		return mann, err
	}

	if err := ar.existAnno(attr, tag); err != nil {
		return mann, err
	}

	return ar.createAnno(
		&createParams{
			attr: attr,
			id:   cvtid,
			tag:  tag,
		},
	)
}

func (ar *arangorepository) EditAnnotation(uat *annotation.TaggedAnnotationUpdate) (*model.AnnoDoc, error) {
	mann := &model.AnnoDoc{}
	attr := uat.Data.Attributes
	rgt, err := ar.database.Get(
		fmt.Sprintf(
			annGetQ,
			ar.anno.annot.Name(),
			ar.anno.annotg.Name(),
			ar.onto.Cv.Name(),
			uat.Data.Id,
		),
	)
	if err != nil {
		return mann, fmt.Errorf("error in fetching id %s", err)
	}
	if rgt.IsEmpty() {
		mann.NotFound = true

		return mann, &repository.AnnoNotFoundError{Id: uat.Data.Id}
	}
	if err := rgt.Read(mann); err != nil {
		return mann, fmt.Errorf("error in reading to struct %s", err)
	}

	bindParams := []interface{}{
		ar.anno.annot.Name(),
		ar.anno.term.Name(),
		ar.anno.ver.Name(),
		attr.Value,
		attr.EditableValue,
		attr.CreatedBy,
		mann.EnrtyId,
		mann.Rank,
		mann.Version + 1,
		mann.CvtId,
		mann.ID.String(),
	}
	dbh := ar.database.Handler()
	idt, err := dbh.Transaction(
		context.Background(),
		annVerInstFn,
		&driver.TransactionOptions{
			WriteCollections: []string{
				ar.anno.annot.Name(),
				ar.anno.term.Name(),
				ar.anno.ver.Name(),
			},
			Params:             bindParams,
			MaxTransactionSize: maxTransactionSize,
		})
	if err != nil {
		return mann, fmt.Errorf("error in running transaction %s", err)
	}
	umd, err := model.ConvToModel(idt)
	if err != nil {
		return umd, fmt.Errorf("error in converting model struct %s", err)
	}
	umd.Ontology = mann.Ontology
	umd.Tag = mann.Tag

	return umd, nil
}


func (ar *arangorepository) AddAnnotationGroup(idslice ...string) (*model.AnnoGroup, error) {
	grp := &model.AnnoGroup{}
	if len(idslice) <= 1 {
		return grp, errors.New("need at least more than one entry to form a group")
	}

	if err := DocumentsExists(ar.anno.annot, idslice...); err != nil {
		return grp, err
	}

	mla, err := ar.getAllAnnotations(idslice...)
	if err != nil {
		return grp, err
	}
	dbg := &model.DbGroup{}
	rdn, err := ar.database.DoRun(
		annGroupInst,
		map[string]interface{}{
			"@anno_group_collection": ar.anno.annog.Name(),
			"group":                  idslice,
		},
	)
	if err != nil {
		return grp, fmt.Errorf("error in creating group %s", err)
	}
	if err := rdn.Read(dbg); err != nil {
		return grp, fmt.Errorf("error in reading to struct %s", err)
	}
	grp.CreatedAt = dbg.CreatedAt
	grp.UpdatedAt = dbg.UpdatedAt
	grp.GroupId = dbg.GroupId
	grp.AnnoDocs = mla

	return grp, nil
}


func (ar *arangorepository) RemoveAnnotationGroup(groupID string) error {
	_, err := ar.anno.annog.RemoveDocument(
		context.Background(),
		groupID,
	)
	if err != nil {
		return fmt.Errorf("error in removing group with id %s %s", groupID, err)
	}

	return nil
}


func (ar *arangorepository) AppendToAnnotationGroup(groupID string, idslice ...string) (*model.AnnoGroup, error) {
	grp := &model.AnnoGroup{}
	if len(idslice) <= 1 {
		return grp, errors.New("need at least more than one entry to form a group")
	}

	gml, err := ar.groupID2Annotations(groupID)
	if err != nil {
		return grp, err
	}

	ml, err := ar.getAllAnnotations(idslice...)
	if err != nil {
		return grp, err
	}

	aml := model.UniqueModel(append(gml, ml...))

	rdn, err := ar.database.DoRun(
		annGroupUpd,
		map[string]interface{}{
			"@anno_group_collection": ar.anno.annog.Name(),
			"key":                    groupID,
			"group":                  model.DocToIDs(aml),
		},
	)
	if err != nil {
		return grp, fmt.Errorf("error in updating group with id %s %s", groupID, err)
	}
	dbg := &model.DbGroup{}
	if err := rdn.Read(dbg); err != nil {
		return grp, fmt.Errorf("error in reading to struct %s", err)
	}
	grp.CreatedAt = dbg.CreatedAt
	grp.UpdatedAt = dbg.UpdatedAt
	grp.GroupId = dbg.GroupId
	grp.AnnoDocs = aml

	return grp, nil
}

func (ar *arangorepository) createAnno(params *createParams) (*model.AnnoDoc, error) {
	mann := &model.AnnoDoc{}
	attr := params.attr
	rins, err := ar.database.DoRun(
		annInst, map[string]interface{}{
			"@anno_collection":    ar.anno.annot.Name(),
			"@anno_cv_collection": ar.anno.term.Name(),
			"editable_value":      attr.EditableValue,
			"created_by":          attr.CreatedBy,
			"entry_id":            attr.EntryId,
			"rank":                attr.Rank,
			"value":               attr.Value,
			"to":                  params.id,
			"version":             1,
		})
	if err != nil {
		return mann, fmt.Errorf("error in running query %s", err)
	}
	if rins.IsEmpty() {
		return mann, fmt.Errorf("error in returning newly created document")
	}
	if err := rins.Read(mann); err != nil {
		return mann, fmt.Errorf("error in reading to struct %s", err)
	}
	mann.Tag = params.tag
	mann.Ontology = attr.Ontology

	return mann, nil
}
```

## File: internal/repository/arangodb/arangodb_test.go
```go
package arangodb

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dictyBase/arangomanager/testarango"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-obograph/graph"
	ontostorage "github.com/dictyBase/go-obograph/storage"
	araobo "github.com/dictyBase/go-obograph/storage/arangodb"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tags = []string{
	"private note",
	"name description",
	"name",
	"curator note",
	"description",
	"public note",
	"status",
	"curation",
	"product",
	"gene product",
	"curation status",
	"curator",
	"note",
}

var ddbg = []string{"DDB_G0286429", "DDB_G0294491"}

func toTimestamp(t time.Time) int64 {
	return t.UnixNano() / 1000000
}

func getOntoParams() *araobo.CollectionParams {
	return &araobo.CollectionParams{
		GraphInfo:    "cv",
		OboGraph:     "obograph",
		Relationship: "cvterm_relationship",
		Term:         "cvterm",
	}
}

func getCollectionParams() *CollectionParams {
	return &CollectionParams{
		Annotation:   "annotation",
		AnnoTerm:     "annotation_cvterm",
		AnnoVersion:  "annotation_version",
		AnnoTagGraph: "annotation_tag",
		AnnoVerGraph: "annotation_history",
		AnnoGroup:    "annotation_group",
		AnnoIndexes:  []string{"entry_id"},
	}
}

func loadData(tra *testarango.TestArango) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to get current dir %s", err)
	}
	res, err := os.Open(
		filepath.Join(
			filepath.Dir(dir), "testdata", "dicty_annotation.json",
		),
	)
	if err != nil {
		return fmt.Errorf("error in open file %s", err)
	}
	defer res.Close()
	gra, err := graph.BuildGraph(res)
	if err != nil {
		return fmt.Errorf("error in building graph %s", err)
	}
	dsr, err := araobo.NewDataSource(
		&araobo.ConnectParams{
			User:     tra.User,
			Pass:     tra.Pass,
			Host:     tra.Host,
			Database: tra.Database,
			Port:     tra.Port,
			Istls:    tra.Istls,
		}, getOntoParams(),
	)
	if err != nil {
		return fmt.Errorf("error in creating datasource %s", err)
	}
	if dsr.ExistsOboGraph(gra) {
		return errors.New("dicty_annotation already exist, needs a cleanp")
	}

	return saveExistentTestGraph(dsr, gra)
}

func saveExistentTestGraph(
	dsr ontostorage.DataSource,
	gra graph.OboGraph,
) error {
	if err := dsr.SaveOboGraphInfo(gra); err != nil {
		return fmt.Errorf("error in saving graph %s", err)
	}
	if _, err := dsr.SaveTerms(gra); err != nil {
		return fmt.Errorf("error in saving terms %s", err)
	}
	if _, err := dsr.SaveRelationships(gra); err != nil {
		return fmt.Errorf("error in saving relationships %s", err)
	}

	return nil
}

func newTestAnnoWithTagAndOnto(
	onto, tag string,
) *annotation.NewTaggedAnnotation {
	return &annotation.NewTaggedAnnotation{
		Data: &annotation.NewTaggedAnnotation_Data{
			Type: "annotations",
			Attributes: &annotation.NewTaggedAnnotationAttributes{
				Value:         "developmentally regulated gene",
				EditableValue: "developmentally regulated gene",
				CreatedBy:     "siddbasu@gmail.com",
				Tag:           tag,
				Ontology:      onto,
				EntryId:       "DDB_G0267474",
				Rank:          0,
			},
		},
	}
}

func newTestTaggedAnnotationWithParams(
	tag, entryID string,
) *annotation.NewTaggedAnnotation {
	return &annotation.NewTaggedAnnotation{
		Data: &annotation.NewTaggedAnnotation_Data{
			Type: "annotations",
			Attributes: &annotation.NewTaggedAnnotationAttributes{
				Value:         "developmentally regulated gene",
				EditableValue: "developmentally regulated gene",
				CreatedBy:     "siddbasu@gmail.com",
				Tag:           tag,
				Ontology:      "dicty_annotation",
				EntryId:       entryID,
				Rank:          0,
			},
		},
	}
}

func newTestTaggedAnnotation() *annotation.NewTaggedAnnotation {
	return &annotation.NewTaggedAnnotation{
		Data: &annotation.NewTaggedAnnotation_Data{
			Type: "annotations",
			Attributes: &annotation.NewTaggedAnnotationAttributes{
				Value:         "developmentally regulated gene",
				EditableValue: "developmentally regulated gene",
				CreatedBy:     "siddbasu@gmail.com",
				Tag:           "description",
				Ontology:      "dicty_annotation",
				EntryId:       "DDB_G0267474",
				Rank:          0,
			},
		},
	}
}

func newTestTaggedAnnotationsListForFiltering(
	num int,
) []*annotation.NewTaggedAnnotation {
	var nal []*annotation.NewTaggedAnnotation
	value := fmt.Sprintf("cool gene %s", tags[0])
	for zcount := 0; zcount < num/2; zcount++ {
		nal = append(nal, &annotation.NewTaggedAnnotation{
			Data: &annotation.NewTaggedAnnotation_Data{
				Type: "annotations",
				Attributes: &annotation.NewTaggedAnnotationAttributes{
					Value:         value,
					EditableValue: value,
					CreatedBy:     "sidd@gmail.com",
					Tag:           tags[0],
					Ontology:      "dicty_annotation",
					EntryId:       ddbg[0],
					Rank:          int64(zcount),
				},
			},
		})
	}
	value = fmt.Sprintf("cool gene %s", tags[1])
	for ycount := num / 2; ycount < num; ycount++ {
		nal = append(nal, &annotation.NewTaggedAnnotation{
			Data: &annotation.NewTaggedAnnotation_Data{
				Type: "annotations",
				Attributes: &annotation.NewTaggedAnnotationAttributes{
					Value:         value,
					EditableValue: value,
					CreatedBy:     "basu@gmail.com",
					Tag:           tags[1],
					Ontology:      "dicty_annotation",
					EntryId:       ddbg[1],
					Rank:          int64(ycount),
				},
			},
		})
	}

	return nal
}

func newTestTaggedAnnotationsList(num int) []*annotation.NewTaggedAnnotation {
	nal := make([]*annotation.NewTaggedAnnotation, 0)
	rsrc := rand.New(rand.NewSource(time.Now().UnixNano()))
	geneIDMax := 800000
	geneIDMin := 300000
	for i := 0; i < num; i++ {
		value := fmt.Sprintf("cool gene %s", tags[rsrc.Intn(len(tags)-1)])
		nal = append(nal, &annotation.NewTaggedAnnotation{
			Data: &annotation.NewTaggedAnnotation_Data{
				Type: "annotations",
				Attributes: &annotation.NewTaggedAnnotationAttributes{
					Value:         value,
					EditableValue: value,
					CreatedBy:     "siddbasu@gmail.com",
					Tag:           tags[rsrc.Intn(len(tags)-1)],
					Ontology:      "dicty_annotation",
					EntryId: fmt.Sprintf(
						"DDB_G0%d",
						rsrc.Intn(geneIDMax-geneIDMin)+geneIDMin,
					),
					Rank: 0,
				},
			},
		})
	}

	return nal
}

func setUp(
	t *testing.T,
) (*require.Assertions, repository.TaggedAnnotationRepository) {
	t.Helper()
	tra, err := testarango.NewTestArangoFromEnv(true)
	if err != nil {
		t.Fatalf("unable to construct new TestArango instance %s", err)
	}
	assert := require.New(t)
	repo, err := NewTaggedAnnotationRepo(
		GetConnectParamsFromDB(tra),
		getCollectionParams(),
		getOntoParams(),
	)
	assert.NoErrorf(
		err,
		"expect no error connecting to annotation repository, received %s",
		err,
	)
	err = loadData(tra)
	assert.NoError(err, "expect no error from loading ontology")

	return assert, repo
}

func tearDown(repo repository.TaggedAnnotationRepository) {
	_ = repo.Dbh().Drop()
}

func TestLoadOboJSON(t *testing.T) {
	t.Parallel()
	assert, anrepo := setUp(t)
	defer tearDown(anrepo)
	fh, err := oboReader()
	assert.NoErrorf(err, "expect no error, received %s", err)
	defer fh.Close()
	info, err := anrepo.LoadOboJSON(bufio.NewReader(fh))
	assert.NoErrorf(err, "expect no error, received %s", err)
	assert.True(info.IsCreated, "should match created status")
}

func oboReader() (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return &os.File{}, fmt.Errorf("unable to get current dir %s", err)
	}

	fhr, err := os.Open(
		filepath.Join(
			filepath.Dir(dir), "testdata", "dicty_phenotypes.json",
		),
	)
	if err != nil {
		return fhr, fmt.Errorf("error in opening file %s", err)
	}

	return fhr, nil
}

func testModelListSort(t *testing.T, m []*model.AnnoDoc) {
	t.Helper()
	assert := require.New(t)
	it, err := NewModelAnnoDocPairWiseIterator(m)
	assert.NoErrorf(err, "expect no error, received %s", err)
	for it.NextModelAnnoDocPair() {
		cm, nm := it.ModelAnnoDocPair()
		assert.Truef(
			nm.CreatedAt.Before(cm.CreatedAt),
			"date %s should be before %s",
			nm.CreatedAt.String(),
			cm.CreatedAt.String(),
		)
	}
}

func testGroupMember(
	t *testing.T,
	gla []*model.AnnoGroup,
	count, idx int,
	email string,
) {
	t.Helper()
	assert := assert.New(t)
	assert.Lenf(gla, count, "should have %d groups", count)
	for _, g := range gla {
		assert.Len(g.AnnoDocs, 5, "should have 5 annotations in each group")
		for _, gdoc := range g.AnnoDocs {
			assert.Equalf(gdoc.Tag, tags[idx], "should have %d as the tag", idx)
			assert.Equalf(
				gdoc.CreatedBy,
				email,
				"should be created by %s",
				email,
			)
			assert.Equal(
				"dicty_annotation",
				gdoc.Ontology,
				"should have dicty_annotation ontology",
			)
			assert.Equalf(
				gdoc.EnrtyId,
				ddbg[idx],
				"should have %d as entry id",
				idx,
			)
		}
	}
}

func testModelMaptoID(
	am []*model.AnnoDoc,
	fn func(m *model.AnnoDoc) string,
) []string {
	str := make([]string, 0)
	for _, m := range am {
		str = append(str, fn(m))
	}

	return str
}

func model2IdCallback(mod *model.AnnoDoc) string {
	return mod.Key
}
```

## File: internal/repository/arangodb/arangodb.go
```go
package arangodb

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
	repo "github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/go-playground/validator/v10"
)

type annoc struct {
	annot  driver.Collection
	term   driver.Collection
	ver    driver.Collection
	annog  driver.Collection
	verg   driver.Graph
	annotg driver.Graph
}

type arangorepository struct {
	sess     *manager.Session
	database *manager.Database
	anno     *annoc
	onto     *ontoarango.OntoCollection
}

func NewTaggedAnnotationRepo(
	connP *manager.ConnectParams, collP *CollectionParams, ontoP *ontoarango.CollectionParams,
) (repo.TaggedAnnotationRepository, error) {
	arp := &arangorepository{}
	if err := validator.New().Struct(collP); err != nil {
		return arp, fmt.Errorf("error in validation %s", err)
	}
	sess, dbh, err := manager.NewSessionDb(connP)
	if err != nil {
		return arp, fmt.Errorf("error in creating new session %s", err)
	}
	ontoc, err := ontoarango.CreateCollection(dbh, ontoP)
	if err != nil {
		return arp, fmt.Errorf("error in creating ontology collection %s", err)
	}
	annoc, err := setAnnotationCollection(dbh, ontoc, collP)

	return &arangorepository{
		sess:     sess,
		database: dbh,
		onto:     ontoc,
		anno:     annoc,
	}, err
}

func setAnnotationCollection(dbh *manager.Database, onto *ontoarango.OntoCollection, collP *CollectionParams) (*annoc, error) {
	annoc, err := setDocumentCollection(dbh, collP)
	if err != nil {
		return annoc, fmt.Errorf("error in creating document collection %s", err)
	}
	verg, err := dbh.FindOrCreateGraph(
		collP.AnnoVerGraph,
		[]driver.EdgeDefinition{
			{
				Collection: annoc.ver.Name(),
				From:       []string{annoc.annot.Name()},
				To:         []string{annoc.annot.Name()},
			},
		},
	)
	if err != nil {
		return annoc, fmt.Errorf("error in creating graph %s", err)
	}
	annotg, err := dbh.FindOrCreateGraph(
		collP.AnnoTagGraph,
		[]driver.EdgeDefinition{
			{
				Collection: annoc.term.Name(),
				From:       []string{annoc.annot.Name()},
				To:         []string{onto.Term.Name()},
			},
		},
	)
	if err != nil {
		return annoc, fmt.Errorf("error in creating graph %s", err)
	}
	annoc.verg = verg
	annoc.annotg = annotg
	_, _, err = dbh.EnsurePersistentIndex(
		annoc.annot.Name(),
		collP.AnnoIndexes,
		&driver.EnsurePersistentIndexOptions{
			InBackground: true,
		},
	)

	return annoc, err
}

func setDocumentCollection(dbh *manager.Database, collP *CollectionParams) (*annoc, error) {
	anns := &annoc{}
	anno, err := dbh.FindOrCreateCollection(
		collP.Annotation,
		&driver.CreateCollectionOptions{},
	)
	if err != nil {
		return anns, fmt.Errorf("error in finding or creating collection %s", err)
	}
	annogrp, err := dbh.FindOrCreateCollection(
		collP.AnnoGroup,
		&driver.CreateCollectionOptions{},
	)
	if err != nil {
		return anns, fmt.Errorf("error in finding or creating collection %s", err)
	}
	annocvt, err := dbh.FindOrCreateCollection(
		collP.AnnoTerm,
		&driver.CreateCollectionOptions{Type: driver.CollectionTypeEdge},
	)
	if err != nil {
		return anns, fmt.Errorf("error in finding or creating collection %s", err)
	}
	annov, err := dbh.FindOrCreateCollection(
		collP.AnnoVersion,
		&driver.CreateCollectionOptions{Type: driver.CollectionTypeEdge},
	)

	return &annoc{
		annot: anno,
		annog: annogrp,
		term:  annocvt,
		ver:   annov,
	}, err
}



func (ar *arangorepository) Clear() error {
	if err := ar.ClearAnnotations(); err != nil {
		return err
	}
	for _, c := range []driver.Collection{
		ar.onto.Term, ar.onto.Cv, ar.onto.Rel,
	} {
		if err := c.Truncate(context.Background()); err != nil {
			return fmt.Errorf("error in truncating %s", err)
		}
	}

	err := ar.onto.Obog.Remove(context.Background())
	if err != nil {
		return fmt.Errorf("error in removing graph %s", err)
	}

	return nil
}


func (ar *arangorepository) ClearAnnotations() error {
	for _, c := range []driver.Collection{
		ar.anno.annot, ar.anno.ver, ar.anno.term, ar.anno.annog,
	} {
		if err := c.Truncate(context.Background()); err != nil {
			return fmt.Errorf("error in truncating %s", err)
		}
	}
	for _, grph := range []driver.Graph{
		ar.anno.verg,
		ar.anno.annotg,
	} {
		arangoDb := ar.database.Handler()
		isok, err := arangoDb.GraphExists(context.Background(), grph.Name())
		if err != nil {
			return fmt.Errorf("error in checking existence of graph %s", err)
		}
		if !isok {
			continue
		}
		if err := grph.Remove(context.Background()); err != nil {
			return fmt.Errorf("error in removing graph %s", err)
		}
	}

	return nil
}

func DocumentsExists(c driver.Collection, ids ...string) error {
	for _, kdi := range ids {
		ok, err := c.DocumentExists(context.Background(), kdi)
		if err != nil {
			return fmt.Errorf("error in checking for existence of identifier %s %s", kdi, err)
		}
		if !ok {
			return &repo.AnnoNotFoundError{Id: kdi}
		}
	}

	return nil
}

func (ar *arangorepository) Dbh() *manager.Database {
	return ar.database
}
```

## File: internal/repository/arangodb/field.go
```go
package arangodb


func FilterMap() map[string]string {
	return map[string]string{
		"entry_id":   "ann.entry_id",
		"value":      "ann.value",
		"created_by": "ann.created_by",
		"version":    "ann.version",
		"rank":       "ann.rank",
		"tag":        "cvt.label",
		"ontology":   "cv.metadata.namespace",
	}
}
```

## File: internal/repository/arangodb/list_filter_statement.go
```go
package arangodb

const (
	cvtExclusiveListFilterQ = `
		LET cvtlist = (
		    FOR cvt IN @@cvterm_collection
			FOR cv IN @@cv_collection
		        	FILTER cvt.graph_id == cv._id
				FILTER cvt.deprecated == false
				%s
				RETURN { cv: cv, cvterm: cvt }
		)

		FOR row IN cvtlist
		    FOR entry IN 1..1 INBOUND row.cvterm GRAPH @anno_cvterm_graph
			    FILTER entry.is_obsolete == false
		            LIMIT @limit
		            RETURN MERGE(entry, {
				tag: row.cvterm.label,
				ontology: row.cv.metadata.namespace
			    })
	`

	annExclusiveListFilterQ = `
		LET annentries = (
		    FOR ann IN @@anno_collection
			%s
		        FILTER ann.is_obsolete == false
			SORT ann.created_at DESC
		        RETURN ann
		)

		FOR entry IN annentries
		    FOR cvt IN 1..1 OUTBOUND entry GRAPH @anno_cvterm_graph
		        FOR cv IN @@cv_collection
		            FILTER cvt.graph_id == cv._id
			    FILTER cvt.deprecated == false
		            LIMIT @limit
		            RETURN MERGE(entry, {
				tag: cvt.label,
				ontology: cv.metadata.namespace
			    })
	`
	annCvtListFilterQ = `
		LET annentries = (
		    FOR ann IN @@anno_collection
			%s
		        FILTER ann.is_obsolete == false
			SORT ann.created_at DESC
		        RETURN ann
		)

		FOR entry IN annentries
		    FOR cvt IN 1..1 OUTBOUND entry GRAPH @anno_cvterm_graph
		        FOR cv IN @@cv_collection
		            FILTER cvt.graph_id == cv._id
			    FILTER cvt.deprecated == false
			    %s
		            LIMIT @limit
		            RETURN MERGE(entry, {
				tag: cvt.label,
				ontology: cv.metadata.namespace
			    })
	`

	cvtExclusiveListFilterWithCursorQ = `
		LET cvtlist = (
		    FOR cvt IN @@cvterm_collection
			FOR cv IN @@cv_collection
		        	FILTER cvt.graph_id == cv._id
				FILTER cvt.deprecated == false
				%s
				RETURN { cv: cv, cvterm: cvt }
		)

		FOR row IN cvtlist
		    FOR entry IN 1..1 INBOUND row.cvterm GRAPH @anno_cvterm_graph
			    FILTER entry.is_obsolete == false
			    FILTER entry.created_at <= DATE_ISO8601(@cursor)
			    SORT entry.created_at DESC
		            LIMIT @limit
		            RETURN MERGE(entry, {
				tag: row.cvterm.label,
				ontology: row.cv.metadata.namespace
			    })
	`

	annExclusiveListFilterWithCursorQ = `
		LET annentries = (
		    FOR ann IN @@anno_collection
			%s
		        FILTER ann.is_obsolete == false
			FILTER ann.created_at <= DATE_ISO8601(@cursor)
			SORT ann.created_at DESC
		        RETURN ann
		)

		FOR entry IN annentries
		    FOR cvt IN 1..1 OUTBOUND entry GRAPH @anno_cvterm_graph
		        FOR cv IN @@cv_collection
		            FILTER cvt.graph_id == cv._id
			    FILTER cvt.deprecated == false
		            LIMIT @limit
		            RETURN MERGE(entry, {
				tag: cvt.label,
				ontology: cv.metadata.namespace
			    })
	`

	annCvtListFilterWithCursorQ = `
		LET annentries = (
		    FOR ann IN @@anno_collection
			%s
		        FILTER ann.is_obsolete == false
			FILTER ann.created_at <= DATE_ISO8601(@cursor)
			SORT ann.created_at DESC
		        RETURN ann
		)

		FOR entry IN annentries
		    FOR cvt IN 1..1 OUTBOUND entry GRAPH @anno_cvterm_graph
		        FOR cv IN @@cv_collection
		            FILTER cvt.graph_id == cv._id
			    FILTER cvt.deprecated == false
			    %s
		            LIMIT @limit
		            RETURN MERGE(entry, {
				tag: cvt.label,
				ontology: cv.metadata.namespace
			    })
	`
)
```

## File: internal/repository/arangodb/ontology.go
```go
package arangodb

import (
	"fmt"
	"io"

	"github.com/dictyBase/go-obograph/storage"
	ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
)

func (ar *arangorepository) LoadOboJSON(rde io.Reader) (*storage.UploadInformation, error) {
	dsb, err := ontoarango.NewDataSourceFromDb(ar.database, &ontoarango.CollectionParams{
		OboGraph:     ar.onto.Obog.Name(),
		GraphInfo:    ar.onto.Cv.Name(),
		Relationship: ar.onto.Rel.Name(),
		Term:         ar.onto.Term.Name(),
	})
	if err != nil {
		return &storage.UploadInformation{}, fmt.Errorf("error in creating new data source %s", err)
	}

	info, err := storage.LoadOboJSONFromDataSource(rde, dsb)
	if err != nil {
		return &storage.UploadInformation{}, fmt.Errorf("error in uploading JSON %s", err)
	}

	return info, nil
}

func (ar *arangorepository) termID(onto, term string) (string, error) {
	var tid string
	row, err := ar.database.GetRow(annExistTagQ, map[string]interface{}{
		"@cv_collection":     ar.onto.Cv.Name(),
		"@cvterm_collection": ar.onto.Term.Name(),
		"ontology":           onto,
		"tag":                term,
	})
	if err != nil {
		return tid, fmt.Errorf("error in running obograph retrieving query %s", err)
	}
	if row.IsEmpty() {
		return tid, fmt.Errorf("ontology %s and tag %s does not exist", onto, term)
	}
	if err := row.Read(&tid); err != nil {
		return tid, fmt.Errorf("error in retrieving obograph id %s", err)
	}

	return tid, nil
}

func (ar *arangorepository) termName(tid string) (string, error) {
	var name string
	cvtr, err := ar.database.GetRow(cvtID2LblQ, map[string]interface{}{
		"@cvterm_collection": ar.onto.Term.Name(),
		"id":                 tid,
	})
	if err != nil {
		return name,
			fmt.Errorf("error in running tag retrieving query %s", err)
	}
	if cvtr.IsEmpty() {
		return name, fmt.Errorf("cvterm id %s does not exist", tid)
	}
	if err := cvtr.Read(&name); err != nil {
		return name, fmt.Errorf("error in retrieving tag %s", err)
	}

	return name, nil
}
```

## File: internal/repository/arangodb/organism_test_helpers.go
```go
package arangodb

import (
	"testing"
	"time"

	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/arangomanager/testarango"
	"github.com/dictyBase/go-genproto/dictybaseapis/organism"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type validateOrganismParams struct {
	assertions *require.Assertions
	got        *model.OrganismDoc
	key        string
	baseOrg    *organism.NewOrganism
}

type updateOrganismParams struct {
	t      *testing.T
	asrt   *require.Assertions
	repo   repository.OrganismRepository
	id     string
	params *organism.OrganismUpdate
}

func getFullUpdateParams() *organism.OrganismUpdate {
	return &organism.OrganismUpdate{
		UpdatedBy: "another@email.com",
		Attributes: &organism.OrganismAttributes{
			Species: "purpureum", Genus: "Dictyostelium",
			CommonName: "purple slime mold", Abbreviation: "dpur",
		},
	}
}

func getPartialUpdateParams() *organism.OrganismUpdate {
	return &organism.OrganismUpdate{
		UpdatedBy: "partial@email.com",
		Attributes: &organism.OrganismAttributes{
			CommonName: "new common name",
		},
	}
}

func getNotFoundUpdateParams() *organism.OrganismUpdate {
	return &organism.OrganismUpdate{
		Id: "non_existent_id", UpdatedBy: "mock@email.com",
		Attributes: &organism.OrganismAttributes{Species: "new species"},
	}
}

func updateOrganism(params updateOrganismParams) *model.OrganismDoc {
	params.t.Helper()
	params.params.Id = params.id
	updated, err := params.repo.EditOrganism(params.params)
	params.asrt.NoError(err, "expected no error updating organism")

	return updated
}

func validateFullUpdate(
	asrt *require.Assertions,
	updated *model.OrganismDoc,
	originalKey string,
) {
	params := getFullUpdateParams()
	asrt.Equal(params.Attributes.Species, updated.Species)
	asrt.Equal(params.Attributes.Genus, updated.Genus)
	asrt.Equal(params.Attributes.CommonName, updated.CommonName)
	asrt.Equal(params.Attributes.Abbreviation, updated.Abbreviation)
	asrt.Equal(params.UpdatedBy, updated.UpdatedBy)
	asrt.Equal(originalKey, updated.Key)
}

func validatePartialUpdate(
	asrt *require.Assertions,
	updated *model.OrganismDoc,
) {
	asrt.Equal("new common name", updated.CommonName)
	asrt.Equal("purpureum", updated.Species)
	asrt.Equal("Dictyostelium", updated.Genus)
	asrt.Equal("dpur", updated.Abbreviation)
	asrt.Equal("partial@email.com", updated.UpdatedBy)
}

func validateOrganism(params validateOrganismParams) {
	params.assertions.Equal(
		params.key,
		params.got.Key,
		"should have matching keys",
	)
	params.assertions.Equal(
		params.baseOrg.Attributes.Species,
		params.got.Species,
		"should have matching species",
	)
	params.assertions.Equal(
		params.baseOrg.Attributes.Genus,
		params.got.Genus,
		"should have matching genus",
	)
	params.assertions.Equal(
		params.baseOrg.Attributes.CommonName,
		params.got.CommonName,
		"should have matching common name",
	)
	params.assertions.Equal(
		params.baseOrg.Attributes.Abbreviation,
		params.got.Abbreviation,
		"should have matching abbreviation",
	)
	params.assertions.Equal(
		params.baseOrg.CreatedBy,
		params.got.CreatedBy,
		"should have matching creator",
	)
}

func setupTestOrganism(
	t *testing.T,
	asrt *require.Assertions,
	repo repository.OrganismRepository,
) *model.OrganismDoc {
	t.Helper()
	baseOrg := &organism.NewOrganism{
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &organism.OrganismAttributes{
			Species: "discoideum", Genus: "Dictyostelium",
			CommonName: "slime mold", Abbreviation: "ddis",
		},
	}
	added, err := repo.AddOrganism(baseOrg)
	asrt.NoError(err, "expected no error adding test organism")

	return added
}

func getOrganismTestCases(baseOrg *organism.NewOrganism) []struct {
	name     string
	attrs    *organism.OrganismAttributes
	wantErr  bool
	validate func(*require.Assertions, *model.OrganismDoc)
} {
	return []struct {
		name     string
		attrs    *organism.OrganismAttributes
		wantErr  bool
		validate func(*require.Assertions, *model.OrganismDoc)
	}{
		{
			name: "success",
			attrs: &organism.OrganismAttributes{
				Species:      "discoideum",
				Genus:        "Dictyostelium",
				CommonName:   "slime mold",
				Abbreviation: "ddis",
			},
			validate: func(asrt *require.Assertions, mdl *model.OrganismDoc) {
				asrt.Equal("discoideum", mdl.Species)
				asrt.Equal("Dictyostelium", mdl.Genus)
				asrt.Equal("slime mold", mdl.CommonName)
				asrt.Equal("ddis", mdl.Abbreviation)
				asrt.Equal(baseOrg.CreatedBy, mdl.CreatedBy)
				asrt.False(mdl.NotFound)
			},
		}, {
			name: "minimal",
			attrs: &organism.OrganismAttributes{
				Species: "aurelia",
				Genus:   "Polysphondylium",
			},
			validate: func(asrt *require.Assertions, mdl *model.OrganismDoc) {
				asrt.Equal("aurelia", mdl.Species)
				asrt.Equal("Polysphondylium", mdl.Genus)
				asrt.Empty(mdl.CommonName)
				asrt.Empty(mdl.Abbreviation)
				asrt.Equal(baseOrg.CreatedBy, mdl.CreatedBy)
				asrt.False(mdl.NotFound)
			},
		},
	}
}

func setUpOrganismTest(
	t *testing.T,
) (*require.Assertions, repository.OrganismRepository) {
	t.Helper()
	tra, err := testarango.NewTestArangoFromEnv(true)
	require.NoError(t, err, "unable to construct new TestArango instance")
	assert := require.New(t)
	repo, err := NewOrganismRepo(
		GetConnectParamsFromDB(tra),
		&OrganismCollectionParams{Organism: "organism"},
	)
	assert.NoErrorf(
		err,
		"expect no error connecting to organism repository, received %s",
		err,
	)

	return assert, repo
}

func getTestOrganisms() []*organism.NewOrganism {
	return []*organism.NewOrganism{
		{
			CreatedBy: "mock@email.com",
			CreatedAt: timestamppb.New(time.Now()),
			Attributes: &organism.OrganismAttributes{
				Species: "discoideum",
				Genus:   "Dictyostelium",
			},
		},
		{
			CreatedBy: "mock@email.com",
			CreatedAt: timestamppb.New(time.Now()),
			Attributes: &organism.OrganismAttributes{
				Species: "purpureum",
				Genus:   "Dictyostelium",
			},
		},
		{
			CreatedBy: "mock@email.com",
			CreatedAt: timestamppb.New(time.Now()),
			Attributes: &organism.OrganismAttributes{
				Species: "fasciculatum",
				Genus:   "Polysphondylium",
			},
		},
	}
}



func GetConnectParamsFromDB(tra *testarango.TestArango) *manager.ConnectParams {
	return &manager.ConnectParams{
		User:     tra.User,
		Pass:     tra.Pass,
		Database: tra.Database,
		Host:     tra.Host,
		Port:     tra.Port,
		Istls:    false,
	}
}
```

## File: internal/repository/arangodb/organism_test.go
```go
package arangodb

import (
	"fmt"
	"testing"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/organism"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAddOrganism(t *testing.T) {
	t.Parallel()
	baseOrg := &organism.NewOrganism{
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
	}
	for _, tcs := range getOrganismTestCases(baseOrg) {
		t.Run(tcs.name, func(t *testing.T) {
			t.Parallel()
			asrt, repo := setUpOrganismTest(t)
			t.Cleanup(func() { _ = repo.Dbh().Drop() })
			baseOrg.Attributes = tcs.attrs
			doc, err := repo.AddOrganism(baseOrg)
			if tcs.wantErr {
				asrt.Error(err)

				return
			}
			asrt.NoError(err)
			tcs.validate(asrt, doc)
		})
	}
}

func TestAddDuplicateOrganism(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })


	baseOrg := &organism.NewOrganism{
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &organism.OrganismAttributes{
			Species: "discoideum",
			Genus:   "Dictyostelium",
		},
	}


	_, err := repo.AddOrganism(baseOrg)
	asrt.NoError(err, "expected no error adding first organism")


	_, err = repo.AddOrganism(baseOrg)
	asrt.Error(err, "expected error when adding duplicate organism")
	asrt.Contains(
		err.Error(),
		"organism Dictyostelium discoideum already exists",
		"expected duplicate organism error message",
	)
}


func TestGetOrganism(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })
	baseOrg := &organism.NewOrganism{
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &organism.OrganismAttributes{
			Species:      "discoideum",
			Genus:        "Dictyostelium",
			CommonName:   "slime mold",
			Abbreviation: "ddis",
		},
	}
	added, err := repo.AddOrganism(baseOrg)
	asrt.NoError(err, "expected no error adding test organism")

	t.Run("success", func(t *testing.T) {

		got, err := repo.GetOrganism(added.Key)
		asrt.NoError(err, "expected no error getting organism")
		validateOrganism(validateOrganismParams{
			assertions: asrt,
			got:        got,
			key:        added.Key,
			baseOrg:    baseOrg,
		})
	})

	t.Run("not found", func(t *testing.T) {

		_, err := repo.GetOrganism("non_existent_id")
		asrt.Error(err, "expected error for non-existent organism")
		asrt.True(
			repository.IsOrganismNotFound(err),
			"should be organism not found error",
		)
	})
}


func TestGetOrganismByName(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })


	baseOrg := &organism.NewOrganism{
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &organism.OrganismAttributes{
			Species:      "discoideum",
			Genus:        "Dictyostelium",
			CommonName:   "slime mold",
			Abbreviation: "ddis",
		},
	}
	added, err := repo.AddOrganism(baseOrg)
	asrt.NoError(err, "expected no error adding test organism")

	t.Run("success", func(t *testing.T) {

		got, err := repo.GetOrganismByName(
			baseOrg.Attributes.Genus,
			baseOrg.Attributes.Species,
		)
		asrt.NoError(err, "expected no error getting organism by name")
		validateOrganism(validateOrganismParams{
			assertions: asrt,
			got:        got,
			key:        added.Key,
			baseOrg:    baseOrg,
		})
	})

	t.Run("not found", func(t *testing.T) {

		_, err := repo.GetOrganismByName("NonExistent", "Species")
		asrt.Error(err, "expected error for non-existent organism")
		asrt.True(
			repository.IsOrganismNotFound(err),
			"should be organism not found error",
		)
		asrt.Contains(
			err.Error(),
			"NonExistent Species",
			"error should contain the non-existent organism name",
		)
	})
}


func TestEditOrganism(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })
	added := setupTestOrganism(t, asrt, repo)

	t.Run("success - full update", func(t *testing.T) {

		updated := updateOrganism(updateOrganismParams{
			t:      t,
			asrt:   asrt,
			repo:   repo,
			id:     added.Key,
			params: getFullUpdateParams(),
		})
		validateFullUpdate(asrt, updated, added.Key)
	})

	t.Run("success - partial update", func(t *testing.T) {

		updated := updateOrganism(updateOrganismParams{
			t:      t,
			asrt:   asrt,
			repo:   repo,
			id:     added.Key,
			params: getPartialUpdateParams(),
		})
		validatePartialUpdate(asrt, updated)
	})

	t.Run("not found", func(t *testing.T) {

		_, err := repo.EditOrganism(getNotFoundUpdateParams())
		asrt.Error(err, "expected error for non-existent organism")
		asrt.True(repository.IsOrganismNotFound(err))
	})
}


func TestRemoveOrganism(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })
	added := setupTestOrganism(t, asrt, repo)

	t.Run("success", func(t *testing.T) {

		err := repo.RemoveOrganism(added.Key)
		asrt.NoError(err, "expected no error removing organism")


		_, err = repo.GetOrganism(added.Key)
		asrt.Error(err, "expected error getting removed organism")
		asrt.True(
			repository.IsOrganismNotFound(err),
			"should be organism not found error",
		)
	})

	t.Run("not found", func(t *testing.T) {

		err := repo.RemoveOrganism("non_existent_id")
		asrt.Error(err, "expected error removing non-existent organism")
		asrt.True(
			repository.IsOrganismNotFound(err),
			"should be organism not found error",
		)
	})
}


func TestListOrganisms(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })
	organisms := getTestOrganisms()
	for _, org := range organisms {
		_, err := repo.AddOrganism(org)
		asrt.NoError(err, "expected no error adding test organism")
	}
	t.Run("success", func(t *testing.T) {
		olist, err := repo.ListOrganisms()
		asrt.NoError(err, "expected no error listing organisms")
		asrt.Len(
			olist,
			len(organisms),
			"should return correct number of organisms",
		)
		expectedOrgs := make(map[string]*organism.NewOrganism)
		for _, org := range organisms {
			key := fmt.Sprintf(
				"%s_%s",
				org.Attributes.Genus,
				org.Attributes.Species,
			)
			expectedOrgs[key] = org
		}
		for _, org := range olist {
			key := fmt.Sprintf("%s_%s", org.Genus, org.Species)
			expected, ok := expectedOrgs[key]
			asrt.True(
				ok,
				"should find matching organism for %s %s",
				org.Genus,
				org.Species,
			)
			asrt.Equal(
				expected.Attributes.Species,
				org.Species,
				"should have matching species",
			)
			asrt.Equal(
				expected.Attributes.Genus,
				org.Genus,
				"should have matching genus",
			)
			asrt.Equal(
				expected.CreatedBy,
				org.CreatedBy,
				"should have matching creator",
			)
			asrt.False(org.NotFound, "organism should exist")
		}
	})
}

func TestClearOrganisms(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpOrganismTest(t)
	t.Cleanup(func() { _ = repo.Dbh().Drop() })


	organisms := getTestOrganisms()
	for _, org := range organisms {
		_, err := repo.AddOrganism(org)
		asrt.NoError(err, "expected no error adding test organism")
	}


	list, err := repo.ListOrganisms()
	asrt.NoError(err, "expected no error listing organisms")
	asrt.Len(list, len(organisms), "should have correct number of organisms")


	err = repo.ClearOrganisms()
	asrt.NoError(err, "expected no error clearing organisms")

	_, err = repo.GetOrganismByName("Dictyostelium", "discoideum")
	asrt.Error(err, "expected error getting cleared organism")
	asrt.True(
		repository.IsOrganismNotFound(err),
		"should be organism not found error",
	)
	_, err = repo.GetOrganismByName("Polysphondylium", "fasciculatum")
	asrt.Error(err, "expected error getting cleared organism")
	asrt.True(
		repository.IsOrganismNotFound(err),
		"should be organism not found error",
	)
}
```

## File: internal/repository/arangodb/organism.go
```go
package arangodb

import (
	"context"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	dorg "github.com/dictyBase/go-genproto/dictybaseapis/organism"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/go-playground/validator/v10"
)

type organismRepo struct {
	sess     *manager.Session
	database *manager.Database
	organism driver.Collection
}



func NewOrganismRepo(
	connP *manager.ConnectParams,
	collP *OrganismCollectionParams,
) (repository.OrganismRepository, error) {

	if err := validator.New().Struct(collP); err != nil {
		return nil, fmt.Errorf("error in validation %s", err)
	}


	sess, dbh, err := manager.NewSessionDb(connP)
	if err != nil {
		return nil, fmt.Errorf("error in creating new session %s", err)
	}


	schemaOpt := &driver.CollectionSchemaOptions{
		Level:   driver.CollectionSchemaLevelStrict,
		Message: "organism schema validation failed",
		Type:    "json",
	}
	if err := schemaOpt.LoadRule(model.Schema()); err != nil {
		return nil, fmt.Errorf("error in loading schema %s", err)
	}
	orgColl, err := dbh.FindOrCreateCollection(
		collP.Organism,
		&driver.CreateCollectionOptions{
			Schema: schemaOpt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error in finding or creating organism collection %s",
			err,
		)
	}

	return &organismRepo{
		sess:     sess,
		database: dbh,
		organism: orgColl,
	}, nil
}




func (org *organismRepo) GetOrganism(id string) (*model.OrganismDoc, error) {
	doc := &model.OrganismDoc{}
	meta, err := org.organism.ReadDocument(context.Background(), id, doc)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return nil, &repository.OrganismNotFoundError{ID: id}
		}

		return nil, fmt.Errorf("error reading organism document: %w", err)
	}
	doc.DocumentMeta = meta

	return doc, nil
}





func (org *organismRepo) GetOrganismByName(
	genus, species string,
) (*model.OrganismDoc, error) {
	res, err := org.database.GetRow(orgGetByNameQ,
		map[string]interface{}{
			"@collection": org.organism.Name(),
			"genus":       genus,
			"species":     species,
		})
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	if res.IsEmpty() {
		return nil, &repository.OrganismNotFoundError{
			ID: fmt.Sprintf("%s %s", genus, species),
		}
	}

	doc := &model.OrganismDoc{}
	if err := res.Read(doc); err != nil {
		return nil, fmt.Errorf("error reading document: %w", err)
	}

	return doc, nil
}

func (org *organismRepo) AddOrganism(
	doc *dorg.NewOrganism,
) (*model.OrganismDoc, error) {

	existing, err := org.GetOrganismByName(
		doc.Attributes.Genus,
		doc.Attributes.Species,
	)
	if err != nil {

		if !repository.IsOrganismNotFound(err) {
			return nil, fmt.Errorf("error checking existing organism: %w", err)
		}
	}

	if existing != nil {
		return nil, fmt.Errorf(
			"organism %s %s already exists",
			doc.Attributes.Genus,
			doc.Attributes.Species,
		)
	}


	orgDoc := &model.OrganismDoc{
		CreatedAt:    doc.CreatedAt.AsTime(),
		UpdatedAt:    doc.CreatedAt.AsTime(),
		CreatedBy:    doc.CreatedBy,
		UpdatedBy:    doc.CreatedBy,
		Abbreviation: doc.Attributes.Abbreviation,
		CommonName:   doc.Attributes.CommonName,
		Species:      doc.Attributes.Species,
		Genus:        doc.Attributes.Genus,
	}


	meta, err := org.organism.CreateDocument(context.Background(), orgDoc)
	if err != nil {
		return nil, fmt.Errorf("error creating organism document: %w", err)
	}


	orgDoc.DocumentMeta = meta

	return orgDoc, nil
}

func (org *organismRepo) EditOrganism(
	doc *dorg.OrganismUpdate,
) (*model.OrganismDoc, error) {
	orgDoc := &model.OrganismDoc{}
	_, err := org.organism.ReadDocument(context.Background(), doc.Id, orgDoc)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return nil, &repository.OrganismNotFoundError{ID: doc.Id}
		}

		return nil, fmt.Errorf("error reading organism document: %w", err)
	}

	update := map[string]interface{}{
		"updated_at": time.Now(),
		"updated_by": doc.UpdatedBy,
	}
	attr := doc.Attributes
	if attr.Abbreviation != "" {
		update["abbreviation"] = attr.Abbreviation
	}
	if attr.CommonName != "" {
		update["common_name"] = attr.CommonName
	}
	if attr.Species != "" {
		update["species"] = attr.Species
	}
	if attr.Genus != "" {
		update["genus"] = attr.Genus
	}
	meta, err := org.organism.UpdateDocument(
		context.Background(),
		doc.Id,
		update,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating organism document: %w", err)
	}


	updatedDoc := &model.OrganismDoc{}
	_, err = org.organism.ReadDocument(
		context.Background(),
		doc.Id,
		updatedDoc,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading updated organism document: %w",
			err,
		)
	}
	updatedDoc.DocumentMeta = meta

	return updatedDoc, nil
}

func (org *organismRepo) RemoveOrganism(oid string) error {

	_, err := org.organism.ReadDocument(context.Background(), oid, nil)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return &repository.OrganismNotFoundError{ID: oid}
		}

		return fmt.Errorf("error checking organism existence: %w", err)
	}


	_, err = org.organism.RemoveDocument(context.Background(), oid)
	if err != nil {
		return fmt.Errorf("error removing organism document: %w", err)
	}

	return nil
}

func (org *organismRepo) ListOrganisms() ([]*model.OrganismDoc, error) {
	cursor, err := org.database.SearchRows(
		orgListQ,
		map[string]interface{}{
			"@collection": org.organism.Name(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error executing organism list query: %w", err)
	}
	defer cursor.Close()
	if cursor.IsEmpty() {
		return nil, &repository.ListNotFoundError{}
	}

	var organisms []*model.OrganismDoc
	for cursor.Scan() {
		omodel := &model.OrganismDoc{}
		if err := cursor.Read(omodel); err != nil {
			return nil, fmt.Errorf("error reading organism document: %w", err)
		}
		organisms = append(organisms, omodel)
	}

	return organisms, nil
}

func (org *organismRepo) ClearOrganisms() error {
	if err := org.organism.Truncate(context.Background()); err != nil {
		return fmt.Errorf("error clearing organisms collection: %w", err)
	}

	return nil
}

func (org *organismRepo) Dbh() *manager.Database {
	return org.database
}
```

## File: internal/repository/arangodb/pairwise.go
```go
package arangodb

import (
	"errors"

	"github.com/dictyBase/modware-annotation/internal/model"
)


type StringPairWiseIterator struct {
	slice []string

	firstIdx int

	secondIdx int

	lastIdx int

	firstPair bool
}



func NewStringPairWiseIterator(mdl []string) (StringPairWiseIterator, error) {
	if len(mdl) <= 1 {
		return StringPairWiseIterator{}, errors.New("not enough element to fetch pairs")
	}

	return StringPairWiseIterator{
		slice:     mdl,
		firstIdx:  0,
		secondIdx: 1,
		lastIdx:   len(mdl) - 1,
		firstPair: true,
	}, nil
}




func (p *StringPairWiseIterator) NextStringPair() bool {
	if p.firstPair {
		p.firstPair = false

		return true
	}
	if p.secondIdx == p.lastIdx {
		return false
	}
	p.firstIdx++
	p.secondIdx++

	return true
}


func (p *StringPairWiseIterator) StringPair() (string, string) {
	return p.slice[p.firstIdx], p.slice[p.secondIdx]
}


type ModelAnnoDocPairWiseIterator struct {
	slice []*model.AnnoDoc

	firstIdx int

	secondIdx int

	lastIdx int

	firstPair bool
}



func NewModelAnnoDocPairWiseIterator(mdl []*model.AnnoDoc) (ModelAnnoDocPairWiseIterator, error) {
	if len(mdl) <= 1 {
		return ModelAnnoDocPairWiseIterator{}, errors.New("not enough element to fetch pairs")
	}

	return ModelAnnoDocPairWiseIterator{
		slice:     mdl,
		firstIdx:  0,
		secondIdx: 1,
		lastIdx:   len(mdl) - 1,
		firstPair: true,
	}, nil
}




func (p *ModelAnnoDocPairWiseIterator) NextModelAnnoDocPair() bool {
	if p.firstPair {
		p.firstPair = false

		return true
	}
	if p.secondIdx == p.lastIdx {
		return false
	}
	p.firstIdx++
	p.secondIdx++

	return true
}


func (p *ModelAnnoDocPairWiseIterator) ModelAnnoDocPair() (*model.AnnoDoc, *model.AnnoDoc) {
	return p.slice[p.firstIdx], p.slice[p.secondIdx]
}
```

## File: internal/repository/organism.go
```go
package repository

import (
	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-genproto/dictybaseapis/organism"
	"github.com/dictyBase/modware-annotation/internal/model"
)



type OrganismRepository interface {

	GetOrganism(id string) (*model.OrganismDoc, error)

	GetOrganismByName(genus, species string) (*model.OrganismDoc, error)

	AddOrganism(doc *organism.NewOrganism) (*model.OrganismDoc, error)

	EditOrganism(doc *organism.OrganismUpdate) (*model.OrganismDoc, error)

	RemoveOrganism(id string) error

	ListOrganisms() ([]*model.OrganismDoc, error)

	ClearOrganisms() error

	Dbh() *manager.Database
}
```

## File: internal/repository/repository.go
```go
package repository

import (
	"io"

	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-obograph/storage"
	"github.com/dictyBase/modware-annotation/internal/model"
)


type ListAnnotationsParams struct {
	Cursor int64
	Limit  int64  `validate:"required"`
	Filter string `validate:"required"`
}



type TaggedAnnotationRepository interface {

	GetAnnotationByID(id string) (*model.AnnoDoc, error)
	GetAnnotationByEntry(
		req *annotation.EntryAnnotationRequest,
	) (*model.AnnoDoc, error)
	AddAnnotation(na *annotation.NewTaggedAnnotation) (*model.AnnoDoc, error)
	EditAnnotation(
		ua *annotation.TaggedAnnotationUpdate,
	) (*model.AnnoDoc, error)
	RemoveAnnotation(id string, purge bool) error


	ListAnnotations(*ListAnnotationsParams) ([]*model.AnnoDoc, error)
	ClearAnnotations() error
	Clear() error

	AddAnnotationGroup(idslice ...string) (*model.AnnoGroup, error)

	GetAnnotationGroup(groupID string) (*model.AnnoGroup, error)

	AppendToAnnotationGroup(
		groupID string,
		idslice ...string,
	) (*model.AnnoGroup, error)

	RemoveAnnotationGroup(groupID string) error

	RemoveFromAnnotationGroup(
		groupID string,
		idslice ...string,
	) (*model.AnnoGroup, error)


	ListAnnotationGroup(
		cursor, limit int64,
		filter string,
	) ([]*model.AnnoGroup, error)

	GetAnnotationTag(name, ontology string) (*model.AnnoTag, error)
	Dbh() *manager.Database
	LoadOboJSON(r io.Reader) (*storage.UploadInformation, error)
}
```

## File: cmd/modware-annotation/main.go
```go
package main

import (
	"log"
	"os"

	apiflag "github.com/dictyBase/aphgrpc"
	arangoflag "github.com/dictyBase/arangomanager/command/flag"
	oboaction "github.com/dictyBase/go-obograph/command/action"
	oboflag "github.com/dictyBase/go-obograph/command/flag"
	obovalidate "github.com/dictyBase/go-obograph/command/validate"
	"github.com/dictyBase/modware-annotation/internal/app/server"
	"github.com/dictyBase/modware-annotation/internal/app/validate"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "modware-annotation"
	app.Usage = "cli for modware-annotation microservice"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-format",
			Usage: "format of the logging out, either of json or text.",
			Value: "json",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "log level for the application",
			Value: "error",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "start-server",
			Usage:  "starts the modware-annotation microservice with grpc backends",
			Action: server.RunServer,
			Before: validate.ServerArgs,
			Flags:  getServerFlags(),
		},
		{
			Name:   "load-ontologies",
			Usage:  "load one or more ontologies in obograph json format",
			Action: oboaction.LoadOntologies,
			Before: obovalidate.OntologyArgs,
			Flags:  oboflag.OntologyFlags(),
		},
		{
			Name:   "start-feature-server",
			Usage:  "starts the feature annotation grpc server",
			Action: server.RunFeatureServer,
			Flags:  getFeatureServerFlags(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error in running command %s", err)
	}
}

func getServerFlags() []cli.Flag {
	flg := []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Usage: "tcp port at which the server will be available",
			Value: "9560",
		},
	}
	flg = append(flg, annoCollFlags()...)
	flg = append(flg, ontoCollFlags()...)
	flg = append(flg, arangoflag.ArangoFlags()...)
	flg = append(flg, []cli.Flag{
		cli.StringFlag{
			Name:   "arangodb-database, db",
			EnvVar: "ARANGODB_DATABASE",
			Usage:  "arangodb database name",
			Value:  "annotation",
		},
		cli.StringFlag{
			Name:  "organism-collection",
			Usage: "arangodb collection for storing organisms",
			Value: "organism",
		},
	}...)

	return append(flg, apiflag.NatsFlag()...)
}

func getFeatureServerFlags() []cli.Flag {
	flg := []cli.Flag{
		cli.StringFlag{
			Name:  "port",
			Usage: "tcp port at which the feature server will be available",
			Value: "9570",
		},
		cli.StringFlag{
			Name:   "arangodb-database, db",
			EnvVar: "ARANGODB_DATABASE",
			Usage:  "arangodb database name",
			Value:  "annofeature",
		},
		cli.StringFlag{
			Name:  "feature-collection",
			Usage: "arangodb collection for storing feature annotations",
			Value: "feature",
		},
		cli.StringFlag{
			Name:  "pub-collection",
			Usage: "arangodb collection for storing publications linked to features",
			Value: "publication",
		},
		cli.StringFlag{
			Name:  "edge-collection",
			Usage: "arangodb edge collection linking features and publications",
			Value: "feature_publication",
		},
		cli.StringFlag{
			Name:  "feature-graph",
			Usage: "arangodb graph name connecting features and publications",
			Value: "feature_pub_graph",
		},
	}
	flg = append(flg, arangoflag.ArangoFlags()...)

	return append(flg, apiflag.NatsFlag()...)
}

func ontoCollFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "term-collection",
			Usage: "arangodb collection for storing ontoloy terms",
			Value: "cvterm",
		},
		cli.StringFlag{
			Name:  "rel-collection",
			Usage: "arangodb collection for storing cvterm relationships",
			Value: "cvterm_relationship",
		},
		cli.StringFlag{
			Name:  "cv-collection",
			Usage: "arangodb collection for storing ontology information",
			Value: "cv",
		},
		cli.StringFlag{
			Name:  "obograph",
			Usage: "arangodb named graph for managing ontology graph",
			Value: "obograph",
		},
		cli.StringSliceFlag{
			Name:  "term-index-fields",
			Usage: "fields to have persistent indexes in ontology term collection",
			Value: &cli.StringSlice{"label"},
		},
	}
}

func annoCollFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "anno-collection",
			Usage: "arangodb collection for storing annotations",
			Value: "annotation",
		},
		cli.StringFlag{
			Name:  "annoterm-collection",
			Usage: "arangodb edge collection for storing links between annotation and ontology term",
			Value: "annotation_cvterm",
		},
		cli.StringFlag{
			Name:  "annover-collection",
			Usage: "arangodb edge collection to link different versions of annotation",
			Value: "annotation_version",
		},
		cli.StringFlag{
			Name:  "annogroup-collection",
			Usage: "arangodb collection for storing annotation group",
			Value: "annotation_group",
		},
		cli.StringFlag{
			Name:  "annoterm-graph",
			Usage: "arangodb named graph for managing relations between annotation and ontology term",
			Value: "annotation_tag",
		},
		cli.StringFlag{
			Name:  "annover-graph",
			Usage: "arangodb named graph for managing relations betweens different versions of annotation",
			Value: "annotation_history",
		},
		cli.StringSliceFlag{
			Name:  "annotation-index-fields",
			Usage: "fields to have persistent indexes in annotation collection",
			Value: &cli.StringSlice{"entry_id"},
		},
	}
}
```

## File: internal/app/server/server_feature.go
```go
package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/dictyBase/aphgrpc"
	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/app/service"
	"github.com/dictyBase/modware-annotation/internal/message"
	"github.com/dictyBase/modware-annotation/internal/message/nats"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/dictyBase/modware-annotation/internal/repository/arangodb"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	gnats "github.com/nats-io/nats.go"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type featureServerParams struct {
	repo repository.FeatureAnnotationRepository
	msg  message.FeatureAnnotationPublisher
}

func RunFeatureServer(clt *cli.Context) error {
	spn, err := featureRepoAndNatsConn(clt)
	if err != nil {
		return cli.NewExitError(err.Error(), errCode)
	}
	grpcS := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(getLogger(clt)),
		),
	)
	srv, err := service.NewFeatureAnnotationService(
		&service.FeatureParams{
			Repository: spn.repo,
			Publisher:  spn.msg,
			Options:    getFeatureGrpcOpt(),
		})
	if err != nil {
		return cli.NewExitError(err.Error(), errCode)
	}
	feature_annotation.RegisterFeatureAnnotationServiceServer(grpcS, srv)
	reflection.Register(grpcS)

	endP := fmt.Sprintf(":%s", clt.String("port"))
	lis, err := net.Listen("tcp", endP)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("failed to listen %s", err), errCode,
		)
	}
	log.Printf("starting feature annotation grpc server on %s", endP)
	if err := grpcS.Serve(lis); err != nil {
		return cli.NewExitError(err.Error(), errCode)
	}

	return nil
}

func allFeatureParams(
	clt *cli.Context,
) (*manager.ConnectParams, *arangodb.FeatureCollectionParams) {
	arPort, _ := strconv.Atoi(clt.String("arangodb-port"))

	return &manager.ConnectParams{
			User:     clt.String("arangodb-user"),
			Pass:     clt.String("arangodb-pass"),
			Database: clt.String("arangodb-database"),
			Host:     clt.String("arangodb-host"),
			Port:     arPort,
			Istls:    clt.Bool("is-secure"),
		}, &arangodb.FeatureCollectionParams{
			Feature: clt.String("feature-collection"),
			Pub:     clt.String("pub-collection"),
			Edge:    clt.String("edge-collection"),
			Graph:   clt.String("feature-graph"),
		}
}

func getFeatureGrpcOpt() []aphgrpc.Option {
	return []aphgrpc.Option{
		aphgrpc.TopicsOption(map[string]string{
			"featureAnnotationCreate": "FeatureAnnotationService.Create",
			"featureAnnotationUpdate": "FeatureAnnotationService.Update",
		}),
	}
}

func featureRepoAndNatsConn(clt *cli.Context) (*featureServerParams, error) {
	connectParams, collParams := allFeatureParams(clt)
	frepo, err := arangodb.NewFeatureAnnoRepo(connectParams, collParams)
	if err != nil {
		return &featureServerParams{},
			fmt.Errorf(
				"cannot connect to arangodb feature repository %s",
				err,
			)
	}
	msp, err := nats.NewFeatureAnnotationPublisher(
		clt.String("nats-host"), clt.String("nats-port"),
		gnats.MaxReconnects(-1), gnats.ReconnectWait(waitTime*time.Second),
	)
	if err != nil {
		return &featureServerParams{},
			fmt.Errorf("cannot connect to messaging server %s", err)
	}

	return &featureServerParams{
		repo: frepo,
		msg:  msp,
	}, nil
}
```

## File: internal/repository/arangodb/annotation_read.go
```go
package arangodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/go-playground/validator/v10"
)


var validate = validator.New(validator.WithRequiredStructEnabled())

func (ar *arangorepository) GetAnnotationByID(
	annoid string,
) (*model.AnnoDoc, error) {
	model := &model.AnnoDoc{}
	res, err := ar.database.Get(
		fmt.Sprintf(
			annGetQ,
			ar.anno.annot.Name(),
			ar.anno.annotg.Name(),
			ar.onto.Cv.Name(),
			annoid,
		),
	)
	if err != nil {
		return model, fmt.Errorf("error in fetching id %s", err)
	}
	if res.IsEmpty() {
		model.NotFound = true

		return model, &repository.AnnoNotFoundError{Id: annoid}
	}
	if err := res.Read(model); err != nil {
		return model, fmt.Errorf("error in reading data to structure %s", err)
	}

	return model, nil
}

func (ar *arangorepository) GetAnnotationByEntry(
	req *annotation.EntryAnnotationRequest,
) (*model.AnnoDoc, error) {
	mann := &model.AnnoDoc{}
	res, err := ar.database.Get(
		fmt.Sprintf(
			annGetByEntryQ,
			ar.anno.annot.Name(),
			ar.anno.annotg.Name(),
			ar.onto.Cv.Name(),
			req.EntryId,
			req.Rank,
			req.IsObsolete,
			req.Tag,
			req.Ontology,
		),
	)
	if err != nil {
		return mann, fmt.Errorf("error in fetching id %s", err)
	}
	if res.IsEmpty() {
		mann.NotFound = true

		return mann, &repository.AnnoNotFoundError{Id: req.EntryId}
	}
	if err := res.Read(mann); err != nil {
		return mann, fmt.Errorf("error in reading data to structure %s", err)
	}

	return mann, nil
}

func (ar *arangorepository) ListAnnotations(
	params *repository.ListAnnotationsParams,
) ([]*model.AnnoDoc, error) {
	if err := validate.Struct(params); err != nil {
		return nil, fmt.Errorf("error in valdating parameters %w", err)
	}
	annoModel := make([]*model.AnnoDoc, 0)
	bindVars := map[string]interface{}{
		"@anno_collection":  ar.anno.annot.Name(),
		"@cv_collection":    ar.onto.Cv.Name(),
		"anno_cvterm_graph": ar.anno.annotg.Name(),
		"limit":             params.Limit + 1,
	}
	if params.Cursor != 0 {
		bindVars["cursor"] = params.Cursor
	}
	result := getListAnnoStatement(params.Filter, params.Cursor)
	if result.Err != nil {
		return nil, result.Err
	}
	res, err := ar.database.SearchRows(result.Statement, bindVars)
	if err != nil {
		return annoModel, fmt.Errorf("error in searching rows %s", err)
	}
	if res.IsEmpty() {
		return annoModel, &repository.AnnoListNotFoundError{}
	}
	for res.Scan() {
		amodel := &model.AnnoDoc{}
		if err := res.Read(amodel); err != nil {
			return annoModel, fmt.Errorf(
				"error in reading data to structure %s",
				err,
			)
		}
		annoModel = append(annoModel, amodel)
	}

	return annoModel, nil
}


func (ar *arangorepository) GetAnnotationGroup(
	groupID string,
) (*model.AnnoGroup, error) {
	grp := &model.AnnoGroup{}
	ann, err := ar.groupID2Annotations(groupID)
	if err != nil {
		return grp, err
	}

	dbg := &model.DbGroup{}
	_, err = ar.anno.annog.ReadDocument(
		context.Background(),
		groupID,
		dbg,
	)
	if err != nil {
		return grp, fmt.Errorf("error in retrieving the group %s", err)
	}
	grp.CreatedAt = dbg.CreatedAt
	grp.UpdatedAt = dbg.UpdatedAt
	grp.GroupId = dbg.GroupId
	grp.AnnoDocs = ann

	return grp, nil
}



func (ar *arangorepository) ListAnnotationGroup(
	cursor, limit int64,
	fstr string,
) ([]*model.AnnoGroup, error) {
	var agrp []*model.AnnoGroup
	var stmt string
	if len(fstr) > 0 {

		stmt = fmt.Sprintf(annGroupListFilterQ,
			ar.anno.annot.Name(), ar.anno.annotg.Name(), ar.onto.Cv.Name(),
			fstr, ar.anno.annog.Name(), ar.anno.annot.Name(),
			ar.anno.annotg.Name(), ar.onto.Cv.Name(),
			limit,
		)
		if cursor != 0 {
			stmt = fmt.Sprintf(annGroupListFilterWithCursorQ,
				ar.anno.annot.Name(), ar.anno.annotg.Name(),
				ar.onto.Cv.Name(), fstr,
				ar.anno.annog.Name(), ar.anno.annot.Name(),
				ar.anno.annotg.Name(), ar.onto.Cv.Name(),
				cursor, limit,
			)
		}
	} else {

		stmt = fmt.Sprintf(annGroupListQ,
			ar.anno.annog.Name(), ar.anno.annot.Name(),
			ar.anno.annotg.Name(), ar.onto.Cv.Name(),
			limit,
		)
		if cursor != 0 {
			stmt = fmt.Sprintf(annGroupListWithCursorQ,
				ar.anno.annog.Name(), ar.anno.annot.Name(),
				ar.anno.annotg.Name(), ar.onto.Cv.Name(),
				cursor, limit,
			)
		}
	}
	res, err := ar.database.Search(stmt)
	if err != nil {
		return agrp, fmt.Errorf("error in searching rows %s", err)
	}
	if res.IsEmpty() {
		return agrp, &repository.AnnoGroupListNotFoundError{}
	}
	for res.Scan() {
		amodel := &model.AnnoGroup{}
		if err := res.Read(amodel); err != nil {
			return agrp, fmt.Errorf(
				"error in reading data to structure %s",
				err,
			)
		}
		agrp = append(agrp, amodel)
	}

	return agrp, nil
}


func (ar *arangorepository) GetAnnotationTag(
	tag, ontology string,
) (*model.AnnoTag, error) {
	annoModel := new(model.AnnoTag)
	res, err := ar.database.GetRow(
		tagGetQ,
		map[string]interface{}{
			"@cvterm_collection": ar.onto.Term.Name(),
			"@cv_collection":     ar.onto.Cv.Name(),
			"ontology":           ontology,
			"tag":                tag,
		})
	if err != nil {
		return annoModel, fmt.Errorf("error in running tag query %s", err)
	}
	if res.IsEmpty() {
		return annoModel, &repository.AnnoTagNotFoundError{Tag: tag}
	}
	if err := res.Read(annoModel); err != nil {
		return annoModel,
			fmt.Errorf(
				"error in retrieving tag %s in ontology %s %s",
				tag, ontology, err,
			)
	}

	return annoModel, nil
}

func (ar *arangorepository) existAnno(
	attr *annotation.NewTaggedAnnotationAttributes,
	tag string,
) error {
	count, err := ar.database.CountWithParams(annExistQ, map[string]interface{}{
		"@anno_collection":  ar.anno.annot.Name(),
		"@cv_collection":    ar.onto.Cv.Name(),
		"anno_cvterm_graph": ar.anno.annotg.Name(),
		"entry_id":          attr.EntryId,
		"rank":              attr.Rank,
		"ontology":          attr.Ontology,
		"tag":               tag,
	})
	if err != nil {
		return fmt.Errorf("error in count query %s", err)
	}
	if count > 0 {
		return errors.New("error in creating, annotation already exists")
	}

	return nil
}

func (ar *arangorepository) groupID2Annotations(
	groupID string,
) ([]*model.AnnoDoc, error) {
	var annoModel []*model.AnnoDoc

	isOk, err := ar.anno.annog.DocumentExists(context.Background(), groupID)
	if err != nil {
		return annoModel,
			fmt.Errorf(
				"error in checking for existence of group identifier %s %s",
				groupID,
				err,
			)
	}
	if !isOk {
		return annoModel, &repository.GroupNotFoundError{Id: groupID}
	}

	dbg := &model.DbGroup{}
	_, err = ar.anno.annog.ReadDocument(
		context.Background(),
		groupID, dbg,
	)
	if err != nil {
		return annoModel, fmt.Errorf("error in retrieving the group %s", err)
	}

	return ar.getAllAnnotations(dbg.Group...)
}

func (ar *arangorepository) getAllAnnotations(
	ids ...string,
) ([]*model.AnnoDoc, error) {
	annoModel := make([]*model.AnnoDoc, 0)
	for _, k := range ids {
		res, err := ar.database.Get(
			fmt.Sprintf(
				annGetQ, ar.anno.annot.Name(),
				ar.anno.annotg.Name(), ar.onto.Cv.Name(), k,
			),
		)
		if err != nil {
			return annoModel, fmt.Errorf("error in fetching id %s", err)
		}
		amodel := &model.AnnoDoc{}
		if err := res.Read(amodel); err != nil {
			return annoModel, fmt.Errorf(
				"error in reading data to structure %s",
				err,
			)
		}
		annoModel = append(annoModel, amodel)
	}

	return annoModel, nil
}
```

## File: internal/repository/arangodb/annotation_statement_test_helpers.go
```go
package arangodb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dictyBase/arangomanager/query"
	"github.com/stretchr/testify/require"
)

func createTestFilterWithSemicolonLogic(field, value string) *query.Filter {
	return &query.Filter{
		Field:    field,
		Value:    value,
		Operator: "==",
		Logic:    ";",
	}
}

func createTestFilter(field, value string) *query.Filter {
	return &query.Filter{
		Field:    field,
		Value:    value,
		Operator: "==",
	}
}

func testBothFiltersWithCursor(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      BothFilters,
		HasCursor: true,
		FilterMap: filterMap,
		FirstSet: []*query.Filter{
			createTestFilterWithSemicolonLogic("value", "val1"),
			createTestFilter("entry_id", "DBS01234"),
		},
		SecondSet: []*query.Filter{createTestFilter("tag", "tag1")},
	}
	result := buildAQLStatement(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.NotEmpty(result.Statement, "statement should not be empty")
	assert.Contains(
		result.Statement,
		"FILTER ann.value",
		"should contain first filter",
	)
	assert.Contains(
		result.Statement,
		"AND ann.entry_id",
		"should contain second filter with AND",
	)
	assert.Contains(
		result.Statement,
		"FILTER cvt.label",
		"should contain second filter",
	)
	assert.Contains(
		result.Statement,
		"DATE_ISO8601(@cursor)",
		"should contain cursor logic",
	)
}

func testBothFiltersWithoutCursor(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      BothFilters,
		HasCursor: false,
		FilterMap: filterMap,
		FirstSet: []*query.Filter{
			createTestFilterWithSemicolonLogic("value", "val1"),
			createTestFilter("entry_id", "DBS01234"),
		},
		SecondSet: []*query.Filter{createTestFilter("tag", "tag1")},
	}
	result := buildAQLStatement(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.NotEmpty(result.Statement, "statement should not be empty")
	assert.Contains(
		result.Statement,
		"FILTER ann.value",
		"should contain first filter",
	)
	assert.Contains(
		result.Statement,
		"AND ann.entry_id",
		"should contain second filter with AND",
	)
	assert.Contains(
		result.Statement,
		"FILTER cvt.label",
		"should contain second filter",
	)
	assert.NotContains(
		result.Statement,
		"DATE_ISO8601",
		"should not contain cursor logic",
	)
	assert.Contains(
		result.Statement,
		"FOR ann IN @@anno_collection",
		"should use annCvtListFilterQ base",
	)
}

func testBuildAQLStatementFirstFilterWithoutCursor(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      FirstFilter,
		HasCursor: false,
		FilterMap: filterMap,
		FirstSet: []*query.Filter{
			createTestFilterWithSemicolonLogic("value", "val1"),
			createTestFilter("entry_id", "DBS01234"),
		},
		SecondSet: []*query.Filter{},
	}
	result := buildAQLStatement(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.NotEmpty(result.Statement, "statement should not be empty")
	assert.Contains(
		result.Statement,
		"FILTER ann.value",
		"should contain first filter",
	)
	assert.Contains(
		result.Statement,
		"AND ann.entry_id",
		"should contain second filter with AND",
	)


}

func testBuildAQLStatementFirstFilterWithCursor(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      FirstFilter,
		HasCursor: true,
		FilterMap: filterMap,
		FirstSet: []*query.Filter{
			createTestFilterWithSemicolonLogic("value", "val1"),
			createTestFilter("entry_id", "DBS01234"),
		},
		SecondSet: []*query.Filter{},
	}
	result := buildAQLStatement(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.NotEmpty(result.Statement, "statement should not be empty")
	assert.Contains(
		result.Statement,
		"FILTER ann.value",
		"should contain first filter",
	)
	assert.Contains(
		result.Statement,
		"AND ann.entry_id",
		"should contain second filter with AND",
	)
	assert.Contains(
		result.Statement,
		"DATE_ISO8601(@cursor)",
		"should contain cursor logic",
	)
}

func testBuildAQLStatementSecondFilterWithoutCursor(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      SecondFilter,
		HasCursor: false,
		FilterMap: filterMap,
		FirstSet:  []*query.Filter{},
		SecondSet: []*query.Filter{
			createTestFilterWithSemicolonLogic("tag", "private note"),
			createTestFilter("ontology", "dicty_annotation"),
		},
	}
	result := buildAQLStatement(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.NotEmpty(result.Statement, "statement should not be empty")
	assert.Contains(
		result.Statement,
		"FILTER cvt.label",
		"should contain second filter",
	)
	assert.Contains(
		result.Statement,
		"AND cv.metadata.namespace",
		"should contain AND logic",
	)
}

func testBuildAQLStatementSecondFilterWithCursor(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      SecondFilter,
		HasCursor: true,
		FilterMap: filterMap,
		FirstSet:  []*query.Filter{},
		SecondSet: []*query.Filter{
			createTestFilterWithSemicolonLogic("tag", "private note"),
			createTestFilter("ontology", "dicty_annotation"),
		},
	}
	result := buildAQLStatement(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.NotEmpty(result.Statement, "statement should not be empty")
	assert.Contains(
		result.Statement,
		"FILTER cvt.label",
		"should contain second filter",
	)
	assert.Contains(
		result.Statement,
		"AND cv.metadata.namespace",
		"should contain AND logic",
	)
	assert.Contains(
		result.Statement,
		"DATE_ISO8601(@cursor)",
		"should contain cursor logic",
	)
}

func testBothFiltersStatementTemplate(t *testing.T) {
	t.Helper()
	t.Run("with cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      BothFilters,
			HasCursor: true,
		}
		template, ok := statementTemplate(ctx)
		assert.True(ok, "should find template")
		assert.Equal(
			annCvtListFilterWithCursorQ,
			template,
			"should return correct template",
		)
	})

	t.Run("without cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      BothFilters,
			HasCursor: false,
		}
		template, ok := statementTemplate(ctx)
		assert.True(ok, "should find template")
		assert.Equal(
			annCvtListFilterQ,
			template,
			"should return correct template",
		)
	})
}

func testFirstFilterStatementTemplate(t *testing.T) {
	t.Helper()
	t.Run("with cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      FirstFilter,
			HasCursor: true,
		}
		template, ok := statementTemplate(ctx)
		assert.True(ok, "should find template")
		assert.Equal(
			annExclusiveListFilterWithCursorQ,
			template,
			"should return correct template",
		)
	})

	t.Run("without cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      FirstFilter,
			HasCursor: false,
		}
		template, ok := statementTemplate(ctx)
		assert.True(ok, "should find template")
		assert.Equal(
			annExclusiveListFilterQ,
			template,
			"should return correct template",
		)
	})
}

func testSecondFilterStatementTemplate(t *testing.T) {
	t.Helper()
	t.Run("with cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      SecondFilter,
			HasCursor: true,
		}
		template, ok := statementTemplate(ctx)
		assert.True(ok, "should find template")
		assert.Equal(
			cvtExclusiveListFilterWithCursorQ,
			template,
			"should return correct template",
		)
	})

	t.Run("without cursor", func(t *testing.T) {
		t.Parallel()
		assert := require.New(t)
		ctx := FilterContext{
			Type:      SecondFilter,
			HasCursor: false,
		}
		template, ok := statementTemplate(ctx)
		assert.True(ok, "should find template")
		assert.Equal(
			cvtExclusiveListFilterQ,
			template,
			"should return correct template",
		)
	})
}

func testInvalidStatementTemplate(t *testing.T) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		Type:      StatementType("invalid"),
		HasCursor: false,
	}
	template, ok := statementTemplate(ctx)
	assert.False(ok, "should not find template")
	assert.Empty(template, "should return empty template")
}

func testValidFilters(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)
	ctx := FilterContext{
		FilterMap: filterMap,
		Filters: []*query.Filter{
			{Field: "entry_id", Value: "DBS123"},
			{Field: "tag", Value: "gene"},
		},
	}

	result := filterAndPartitionFunc(ctx)
	assert.NoError(result.Err, "should not have error")
	assert.Len(result.Filters, 2, "should have 2 valid filters")
	assert.Len(result.FirstSet, 1, "should have 1 filter in first set")
	assert.Len(result.SecondSet, 1, "should have 1 filter in second set")
	assert.Equal(
		"entry_id",
		result.FirstSet[0].Field,
		"first set should contain annotation filter",
	)
	assert.Equal(
		"tag",
		result.SecondSet[0].Field,
		"second set should contain cvterm filter",
	)
}

func testOnlyAnnotationFilters(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)

	ctx := FilterContext{
		FilterMap: filterMap,
		Filters: []*query.Filter{
			{Field: "entry_id", Value: "DBS123"},
			{Field: "value", Value: "test"},
		},
	}

	result := filterAndPartitionFunc(ctx)

	assert.NoError(result.Err, "should not have error")
	assert.Len(result.Filters, 2, "should have 2 valid filters")
	assert.Len(result.FirstSet, 2, "should have 2 filters in first set")
	assert.Empty(result.SecondSet, "should have 0 filters in second set")
}

func testOnlyCvtermFilters(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)

	ctx := FilterContext{
		FilterMap: filterMap,
		Filters: []*query.Filter{
			{Field: "tag", Value: "gene"},
			{Field: "ontology", Value: "GO"},
		},
	}

	result := filterAndPartitionFunc(ctx)

	assert.NoError(result.Err, "should not have error")
	assert.Len(result.Filters, 2, "should have 2 valid filters")
	assert.Empty(result.FirstSet, "should have 0 filters in first set")
	assert.Len(result.SecondSet, 2, "should have 2 filters in second set")
}

func testInvalidFilters(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)

	ctx := FilterContext{
		FilterMap: filterMap,
		Filters: []*query.Filter{
			{Field: "invalid_field", Value: "test"},
		},
	}

	result := filterAndPartitionFunc(ctx)

	assert.Error(result.Err, "should have error")
	assert.Contains(
		result.Err.Error(),
		"no valid filters found",
		"should have appropriate error message",
	)
}

func testMixedFilters(t *testing.T, filterMap map[string]string) {
	t.Helper()
	assert := require.New(t)

	ctx := FilterContext{
		FilterMap: filterMap,
		Filters: []*query.Filter{
			{Field: "entry_id", Value: "DBS123"},
			{Field: "invalid_field", Value: "test"},
		},
	}

	result := filterAndPartitionFunc(ctx)

	assert.NoError(result.Err, "should not have error")
	assert.Len(result.Filters, 1, "should have 1 valid filter")
	assert.Len(result.FirstSet, 1, "should have 1 filter in first set")
	assert.Empty(result.SecondSet, "should have 0 filters in second set")
}

func testExistingError(t *testing.T) {
	t.Helper()
	assert := require.New(t)

	existingErr := errors.New("existing error")
	ctx := FilterContext{
		Err: existingErr,
	}

	result := filterAndPartitionFunc(ctx)

	assert.Equal(existingErr, result.Err, "should preserve existing error")
}

func testParseFiltersFuncSuccess(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	ctx := FilterContext{
		FilterString: `tag==private note;ontology==dicty_annotation`,
	}
	result := parseFiltersFunc(ctx)
	assert.NoError(result.Err, "should not return error")
	assert.Len(result.Filters, 2, "should parse two filters")
	assert.Equal(
		"tag",
		result.Filters[0].Field,
		"first filter field should be tag",
	)
	assert.Equal(
		"private note",
		result.Filters[0].Value,
		"first filter value should be private note",
	)
	assert.Equal(";", result.Filters[0].Logic, "the logic should match")
	assert.Equal(
		"ontology",
		result.Filters[1].Field,
		"second filter field should be ontology",
	)
	assert.Equal(
		"dicty_annotation",
		result.Filters[1].Value,
		"second filter value should be dicty_annotation",
	)
}

func testParseFiltersFuncFailureInvalidString(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	ctx := FilterContext{
		FilterString: "tag=gene;",
	}
	result := parseFiltersFunc(ctx)
	assert.NoError(
		result.Err,
		"should not return error for this specific invalid format",
	)
	assert.Empty(
		result.Filters,
		"filters should be empty for this specific invalid format",
	)
}

func testParseFiltersFuncEdgeEmptyString(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	ctx := FilterContext{
		FilterString: "",
	}
	result := parseFiltersFunc(ctx)
	// Assuming query.ParseFilterString returns empty slice and no error for empty string
	assert.NoError(result.Err, "should not return error for empty string")
	assert.Empty(result.Filters, "should return empty slice for empty string")
}

func testParseFiltersFuncExistingError(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	originalErr := errors.New("previous error")
	ctx := FilterContext{
		FilterString: "tag==gene",
		Err:          originalErr,
	}
	result := parseFiltersFunc(ctx)

	assert.Equal(originalErr, result.Err, "should preserve the original error")
	assert.Nil(
		result.Filters,
		"filters should not be parsed when an error already exists",
	)
}

func testGenFilterStatementSuccessSingle(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	filters := []*query.Filter{
		{Field: "tag", Operator: "==", Value: "gene"},
	}
	expected := "FILTER cvt.label == 'gene'"
	filterType := "cvterm"

	stmt, err := genFilterStatement(filterMap, filters, filterType)
	assert.NoError(err, "should not return error for valid filter")
	assert.Equal(
		expected,
		stmt,
		"should generate correct AQL filter statement",
	)
}

func testGenFilterStatementSuccessMultiple(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	filters := []*query.Filter{
		{
			Field:    "entry_id",
			Operator: "==", Value: "DBS01234", Logic: ";",
		},
		{Field: "value", Operator: "!=", Value: "more test"},
	}


	expectedPart1 := "FILTER ann.entry_id == 'DBS01234'"
	expectedPart2 := "ann.value != 'more test'"
	filterType := "annotation"

	stmt, err := genFilterStatement(filterMap, filters, filterType)
	assert.NoError(
		err,
		"should not return error for multiple valid filters",
	)
	assert.Contains(stmt, expectedPart1, "should contain entry_id filter")
	assert.Contains(stmt, expectedPart2, "should contain value filter")
	assert.Contains(stmt, "AND", "should contain AND operator")

	filters2 := []*query.Filter{
		{
			Field:    "tag",
			Operator: "==", Value: "private note", Logic: ";",
		},
		{Field: "ontology", Operator: "==", Value: "dicty_annotation"},
	}
	_, err = genFilterStatement(filterMap, filters2, filterType)
	assert.NoError(
		err,
		"should not return error for multiple valid filters",
	)
}

func testGenFilterStatementErrorInvalidField(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	filters := []*query.Filter{
		{Field: "invalid_field", Operator: "==", Value: "some_value"},
	}
	filterType := "annotation"

	_, err := genFilterStatement(filterMap, filters, filterType)
	assert.Error(err, "should return error for invalid filter field")
	assert.Contains(
		err.Error(),
		fmt.Sprintf("error generating %s filter", filterType),
		"error message should include filter type",
	)
}

func testGenFilterStatementEdgeEmpty(
	t *testing.T,
	filterMap map[string]string,
) {
	t.Helper()
	assert := require.New(t)
	filters := []*query.Filter{}
	filterType := "cvterm"


	_, err := genFilterStatement(filterMap, filters, filterType)
	assert.Error(err, "should return error for empty filter slice")
}

func testGetListAnnoStatementBasicCases(t *testing.T) {
	t.Parallel()
	t.Run("empty filter string", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement("", 0)
		assert.Error(result.Err, "should return error for empty filter string")
		assert.Equal(
			"empty filter string",
			result.Err.Error(),
			"should have specific error message",
		)
		assert.Empty(result.Statement, "statement should be empty")
	})

	t.Run("invalid filter", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement("invalid_field==test", 0)
		assert.Error(result.Err, "should return error for invalid filter")
		assert.Contains(
			result.Err.Error(),
			"no valid filters found",
			"should have appropriate error message",
		)
		assert.Empty(result.Statement, "statement should be empty")
	})

	t.Run("malformed filter string", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement("value=test", 0)
		assert.Error(
			result.Err,
			"should return error for malformed filter string",
		)
		assert.Empty(result.Statement, "statement should be empty")
	})
}


func testFilterStatement(
	t *testing.T,
	filterString, expectedFilter, filterDescription string,
) {
	t.Helper()
	t.Run("without cursor", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement(filterString, 0)
		assert.NoError(result.Err, "should not return error for filter")
		assert.NotEmpty(result.Statement, "statement should not be empty")
		assert.Contains(
			result.Statement,
			expectedFilter,
			"should contain "+filterDescription,
		)
		assert.NotContains(
			result.Statement,
			"DATE_ISO8601",
			"should not contain cursor logic",
		)
	})

	t.Run("with cursor", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement(filterString, 12345)
		assert.NoError(
			result.Err,
			"should not return error for filter with cursor",
		)
		assert.NotEmpty(result.Statement, "statement should not be empty")
		assert.Contains(
			result.Statement,
			expectedFilter,
			"should contain "+filterDescription,
		)
		assert.Contains(
			result.Statement,
			"DATE_ISO8601(@cursor)",
			"should contain cursor logic",
		)
	})
}

func testGetListAnnoStatementValidFilters(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("valid filter without cursor", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement("value==test", 0)
		assert.NoError(result.Err, "should not return error for valid filter")
		assert.NotEmpty(result.Statement, "statement should not be empty")
		assert.Contains(
			result.Statement,
			"FILTER ann.value",
			"should contain annotation filter",
		)
		assert.NotContains(
			result.Statement,
			"DATE_ISO8601",
			"should not contain cursor logic",
		)
	})

	t.Run("another valid filter without cursor", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement(
			`tag==private note;ontology==dicty_annotation`,
			0,
		)
		assert.NoError(result.Err, "should not return error for valid filter")
		assert.NotEmpty(result.Statement, "statement should not be empty")
	})

	t.Run("valid filter with cursor", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement("value==test", 12345)
		assert.NoError(
			result.Err,
			"should not return error for valid filter with cursor",
		)
		assert.NotEmpty(result.Statement, "statement should not be empty")
		assert.Contains(
			result.Statement,
			"FILTER ann.value",
			"should contain annotation filter",
		)
		assert.Contains(
			result.Statement,
			"DATE_ISO8601(@cursor)",
			"should contain cursor logic",
		)
	})
}

func testGetListAnnoStatementTagFilters(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("skipping integration test")
	}


	testFilterStatement(t, "tag==gene", "FILTER cvt.label", "cvterm filter")

	t.Run("multiple filters", func(t *testing.T) {
		assert := require.New(t)
		result := getListAnnoStatement("value==test;tag==gene", 0)
		assert.NoError(
			result.Err,
			"should not return error for multiple filters",
		)
		assert.NotEmpty(result.Statement, "statement should not be empty")
		assert.Contains(
			result.Statement,
			"FILTER ann.value",
			"should contain annotation filter",
		)
		assert.Contains(
			result.Statement,
			"FILTER cvt.label",
			"should contain cvterm filter",
		)
	})
}
```

## File: internal/app/service/read_service.go
```go
package service

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)

var LIMIT int64 = 10

func (srv *AnnotationService) GetAnnotation(
	ctx context.Context,
	req *annotation.AnnotationId,
) (*annotation.TaggedAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	tna := &annotation.TaggedAnnotation{}
	mid, err := srv.repo.GetAnnotationByID(req.Id)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}
	if mid.NotFound {
		return nil, aphgrpc.HandleNotFoundError(ctx, err)
	}
	tna.Data = srv.getAnnoData(mid)

	return tna, nil
}

func (srv *AnnotationService) GetEntryAnnotation(
	ctx context.Context, rea *annotation.EntryAnnotationRequest,
) (*annotation.TaggedAnnotation, error) {
	if err := protovalidate.Validate(rea); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	tna := &annotation.TaggedAnnotation{}
	mne, err := srv.repo.GetAnnotationByEntry(rea)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}
	tna.Data = srv.getAnnoData(mne)

	return tna, nil
}

func (srv *AnnotationService) GetAnnotationGroup(
	ctx context.Context, rid *annotation.GroupEntryId,
) (*annotation.TaggedAnnotationGroup, error) {
	if err := protovalidate.Validate(rid); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	mga, err := srv.repo.GetAnnotationGroup(rid.GroupId)
	if err != nil {
		if repository.IsGroupNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}

	return srv.getGroup(mga), nil
}

func (srv *AnnotationService) ListAnnotationGroups(
	ctx context.Context, rgp *annotation.ListGroupParameters,
) (*annotation.TaggedAnnotationGroupCollection, error) {
	if err := protovalidate.Validate(rgp); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}

	searchLimit := LIMIT
	if rgp.Limit > 0 {
		searchLimit = rgp.Limit
	}
	astmt, err := filterStrToQuery(rgp.Filter)
	if err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	mgc, err := srv.repo.ListAnnotationGroup(rgp.Cursor, searchLimit, astmt)
	if err != nil {
		if repository.IsAnnotationGroupListNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}
	gcdata := make([]*annotation.TaggedAnnotationGroupCollection_Data, 0)
	for _, mgs := range mgc {
		var gdata []*annotation.TaggedAnnotationGroup_Data
		for _, m := range mgs.AnnoDocs {
			gdata = append(gdata, srv.getAnnoGroupData(m))
		}
		gcdata = append(
			gcdata,
			&annotation.TaggedAnnotationGroupCollection_Data{
				Type: srv.GetGroupResourceName(),
				Group: &annotation.TaggedAnnotationGroup{
					Data:      gdata,
					GroupId:   mgs.GroupId,
					CreatedAt: aphgrpc.TimestampProto(mgs.CreatedAt),
					UpdatedAt: aphgrpc.TimestampProto(mgs.UpdatedAt),
				},
			},
		)
	}
	if len(gcdata) < int(searchLimit)-2 {
		return &annotation.TaggedAnnotationGroupCollection{
			Data: gcdata,
			Meta: &annotation.Meta{Limit: rgp.Limit},
		}, nil
	}

	return &annotation.TaggedAnnotationGroupCollection{
		Data: gcdata[:len(gcdata)-1],
		Meta: &annotation.Meta{
			Limit:      searchLimit,
			NextCursor: genNextCursorVal(mgc[len(mgc)-1].CreatedAt),
		},
	}, nil
}

func (srv *AnnotationService) ListAnnotations(
	ctx context.Context, ral *annotation.ListParameters,
) (*annotation.TaggedAnnotationCollection, error) {
	if err := protovalidate.Validate(ral); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	tac := &annotation.TaggedAnnotationCollection{}

	searchLimit := LIMIT
	if ral.Limit > 0 {
		searchLimit = ral.Limit
	}
	mlc, err := srv.repo.ListAnnotations(
		&repository.ListAnnotationsParams{
			Cursor: ral.Cursor,
			Limit:  searchLimit,
			Filter: ral.Filter,
		})
	if err != nil {
		if repository.IsAnnotationListNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}
	tcdata := collection.Map(mlc, srv.modelToCollectionData)
	if len(tcdata) < int(searchLimit)-2 {
		tac.Data = tcdata
		tac.Meta = &annotation.Meta{Limit: ral.Limit}

		return tac, nil
	}
	tac.Data = tcdata[:len(tcdata)-1]
	tac.Meta = &annotation.Meta{
		Limit:      searchLimit,
		NextCursor: genNextCursorVal(mlc[len(mlc)-1].CreatedAt),
	}

	return tac, nil
}

func (srv *AnnotationService) GetAnnotationTag(
	ctx context.Context, rta *annotation.TagRequest,
) (*annotation.AnnotationTag, error) {
	if err := protovalidate.Validate(rta); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	tag := &annotation.AnnotationTag{}
	mta, err := srv.repo.GetAnnotationTag(rta.Name, rta.Ontology)
	if err != nil {
		if repository.IsAnnoTagNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}
	tag.Id = mta.ID
	tag.Name = mta.Name
	tag.Ontology = mta.Ontology
	tag.IsObsolete = mta.IsObsolete

	return tag, nil
}


func (srv *AnnotationService) modelToCollectionData(
	m *model.AnnoDoc,
) *annotation.TaggedAnnotationCollection_Data {
	return &annotation.TaggedAnnotationCollection_Data{
		Type:       srv.GetResourceName(),
		Id:         m.Key,
		Attributes: getAnnoAttributes(m),
	}
}

func (srv *AnnotationService) getGroup(
	mga *model.AnnoGroup,
) *annotation.TaggedAnnotationGroup {
	gta := &annotation.TaggedAnnotationGroup{}
	gta.Data = srv.getGroupData(mga)
	gta.GroupId = mga.GroupId
	gta.CreatedAt = aphgrpc.TimestampProto(mga.CreatedAt)
	gta.UpdatedAt = aphgrpc.TimestampProto(mga.UpdatedAt)

	return gta
}

func (srv *AnnotationService) getAnnoGroupData(
	m *model.AnnoDoc,
) *annotation.TaggedAnnotationGroup_Data {
	return &annotation.TaggedAnnotationGroup_Data{
		Type:       srv.GetGroupResourceName(),
		Id:         m.Key,
		Attributes: getAnnoAttributes(m),
	}
}

func (srv *AnnotationService) getAnnoData(
	m *model.AnnoDoc,
) *annotation.TaggedAnnotation_Data {
	return &annotation.TaggedAnnotation_Data{
		Type:       srv.GetGroupResourceName(),
		Id:         m.Key,
		Attributes: getAnnoAttributes(m),
	}
}

func (srv *AnnotationService) getGroupData(
	mga *model.AnnoGroup,
) []*annotation.TaggedAnnotationGroup_Data {
	gdata := make([]*annotation.TaggedAnnotationGroup_Data, 0)
	for _, m := range mga.AnnoDocs {
		gdata = append(gdata, srv.getAnnoGroupData(m))
	}

	return gdata
}
```

## File: internal/model/feature_annotation.go
```go
package model

import (
	"encoding/json"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"
)


type DbLinkDoc struct {
	PrimaryId string `json:"primary_id"`
	Database  string `json:"database"`
	Version   int64  `json:"version"`
	LinkType  string `json:"linktype,omitempty"`
	URL       string `json:"url,omitempty"`
	Label     string `json:"label,omitempty"`
}


type TagPropertyDoc struct {
	Tag       string    `json:"tag"`
	Value     string    `json:"value"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



type FeatureAnnotationDoc struct {
	driver.DocumentMeta
	Type         string           `json:"feature_type,omitempty"`
	AnnoId       string           `json:"feature_id"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	CreatedBy    string           `json:"created_by"`
	UpdatedBy    string           `json:"updated_by"`
	Name         string           `json:"name"`
	Synonyms     []string         `json:"synonyms,omitempty"`
	Publications []string         `json:"publications,omitempty"`
	Pubmed       []string         `json:"pubmed,omitempty"`
	DbLinks      []DbLinkDoc      `json:"dblinks,omitempty"`
	Properties   []TagPropertyDoc `json:"properties,omitempty"`
	IsObsolete   bool             `json:"is_obsolete"`
	NotFound     bool             `json:"-"`
}



func PubSchema() ([]byte, error) {
	schema := `{
        		"type": "object",
        		"properties": {
            		"id": { "type": "string" },
            		"created_at": { "type": "string", "format": "date-time" },
            		"updated_at": { "type": "string", "format": "date-time" }
        	},
        	"required": ["id", "created_at", "updated_at"],
        	"additionalProperties": false
    	}`

	return []byte(schema), nil
}


func FeatureAnnotationSchema() ([]byte, error) {
	baseSchema := `{
        "type": "object",
        "properties": %s,
        "required": ["feature_id", "name", "created_at", "created_by", "updated_at", "updated_by"],
        "additionalProperties": true
    }`

	properties := map[string]interface{}{
		"feature_type": map[string]string{"type": "string"},
		"feature_id":   map[string]string{"type": "string"},
		"created_at": map[string]string{
			"type":   "string",
			"format": "date-time",
		},
		"updated_at": map[string]string{
			"type":   "string",
			"format": "date-time",
		},
		"created_by": map[string]string{"type": "string", "format": "email"},
		"updated_by": map[string]string{"type": "string", "format": "email"},
		"name":       map[string]string{"type": "string"},
		"synonyms": map[string]interface{}{
			"type":  "array",
			"items": map[string]string{"type": "string"},
		},
		"dblinks":     getDbLinksSchema(),
		"properties":  getPropertiesSchema(),
		"is_obsolete": map[string]string{"type": "boolean"},
	}


	propsJSON, err := json.Marshal(properties)
	if err != nil {
		return []byte(""),
			fmt.Errorf(
				"failed to marshal feature annotation schema: %v",
				err,
			)
	}

	return []byte(fmt.Sprintf(baseSchema, string(propsJSON))), nil
}

func getDbLinksSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"primary_id": map[string]string{"type": "string"},
				"version":    map[string]string{"type": "integer"},
				"database":   map[string]string{"type": "string"},
				"linktype":   map[string]string{"type": "string"},
				"url":        map[string]string{"type": "string"},
				"label":      map[string]string{"type": "string"},
			},
			"required": []string{"primary_id", "database", "version"},
		},
	}
}

func getPropertiesSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"tag":   map[string]string{"type": "string"},
				"value": map[string]string{"type": "string"},
				"created_by": map[string]string{
					"type":   "string",
					"format": "email",
				},
				"updated_by": map[string]string{
					"type":   "string",
					"format": "email",
				},
				"created_at": map[string]string{
					"type":   "string",
					"format": "date-time",
				},
				"updated_at": map[string]string{
					"type":   "string",
					"format": "date-time",
				},
			},
			"required": []string{
				"tag",
				"value",
				"created_by",
				"created_at",
				"updated_by",
				"updated_at",
			},
		},
	}
}
```

## File: internal/repository/arangodb/parameters.go
```go
package arangodb

import "github.com/dictyBase/go-genproto/dictybaseapis/annotation"

type createParams struct {
	attr *annotation.NewTaggedAnnotationAttributes
	id   string
	tag  string
}



type OrganismCollectionParams struct {

	Organism string `validate:"required"`
}



type CollectionParams struct {

	Annotation string `validate:"required"`

	AnnoGroup string `validate:"required"`


	AnnoTerm string `validate:"required"`


	AnnoVersion string `validate:"required"`


	AnnoTagGraph string `validate:"required"`


	AnnoVerGraph string `validate:"required"`


	AnnoIndexes []string `validate:"required"`
}



type FeatureCollectionParams struct {

	Feature string `validate:"required"`
	Pub     string `validate:"required"`
	Edge    string `validate:"required"`

	Graph string `validate:"required"`
}
```

## File: internal/app/service/feature_annotation.go
```go
package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/message"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeatureAnnotationService struct {
	*aphgrpc.Service
	repo      repository.FeatureAnnotationRepository
	publisher message.FeatureAnnotationPublisher
	feature.UnimplementedFeatureAnnotationServiceServer
}

type FeatureParams struct {
	Repository repository.FeatureAnnotationRepository `validate:"required"`
	Publisher  message.FeatureAnnotationPublisher     `validate:"required"`
	Options    []aphgrpc.Option
}

func featureAnnoDefaultOptions() *aphgrpc.ServiceOptions {
	return &aphgrpc.ServiceOptions{
		Resource: "feature_annotations",
		Topics: map[string]string{
			"featureAnnotationCreate": "FeatureAnnotationCreated",
			"featureAnnotationUpdate": "FeatureAnnotationUpdated",
		},
	}
}

func NewFeatureAnnotationService(
	params *FeatureParams,
) (*FeatureAnnotationService, error) {
	if err := validator.New().Struct(params); err != nil {
		return &FeatureAnnotationService{}, fmt.Errorf(
			"error in validating params %s",
			err,
		)
	}
	svcOpt := featureAnnoDefaultOptions()
	for _, optfn := range params.Options {
		optfn(svcOpt)
	}
	srv := &aphgrpc.Service{}
	aphgrpc.AssignFieldsToStructs(svcOpt, srv)

	return &FeatureAnnotationService{
		Service:   srv,
		repo:      params.Repository,
		publisher: params.Publisher,
	}, nil
}

func (srv *FeatureAnnotationService) GetFeatureAnnotation(
	ctx context.Context,
	req *feature.FeatureAnnotationId,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	feat, err := srv.repo.GetFeatureAnnotation(req.Id)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleGetError(ctx, err)
	}

	return convertToProto(feat), nil
}


func (srv *FeatureAnnotationService) GetFeatureAnnotationByName(
	ctx context.Context,
	req *feature.FeatureName,
) (*feature.FeatureAnnotation, error) {

	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}

	feat, err := srv.repo.GetFeatureAnnotationByName(req.Name)
	if err != nil {

		if repository.IsAnnotationNotFound(err) ||
			repository.IsFeatureNameNotFound(err) {

			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}
		return nil, aphgrpc.HandleGetError(ctx, err)
	}

	return convertToProto(feat), nil
}

func (srv *FeatureAnnotationService) CreateFeatureAnnotation(
	ctx context.Context,
	req *feature.NewFeatureAnnotation,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	feat, err := srv.repo.AddFeatureAnnotation(req)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint violated") {
			return nil, aphgrpc.HandleExistError(ctx, err)
		}

		return nil, aphgrpc.HandleInsertError(ctx, err)
	}
	featProto := convertToProto(feat)
	if err := srv.publisher.Publish(srv.Topics["featureAnnotationCreate"], featProto); err != nil {
		return featProto, aphgrpc.HandleInsertError(ctx, err)
	}

	return featProto, nil
}

func (srv *FeatureAnnotationService) UpdateFeatureAnnotation(
	ctx context.Context,
	req *feature.FeatureAnnotationUpdate,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	feat, err := srv.repo.EditFeatureAnnotation(req)
	if err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}
	featProto := convertToProto(feat)
	if err := srv.publisher.Publish(
		srv.Topics["featureAnnotationUpdate"], featProto,
	); err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}

	return featProto, nil
}

func (srv *FeatureAnnotationService) DeleteFeatureAnnotation(
	ctx context.Context,
	req *feature.DeleteFeatureAnnotationRequest,
) (*emptypb.Empty, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	err := srv.repo.RemoveFeatureAnnotation(req.Id, req.Purge)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}

		return nil, aphgrpc.HandleDeleteError(ctx, err)
	}

	return &emptypb.Empty{}, nil
}

func (srv *FeatureAnnotationService) AddTag(
	ctx context.Context,
	req *feature.AddTagRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	feat, err := srv.repo.AddTag(req)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}
	featProto := convertToProto(feat)
	if err := srv.publisher.Publish(
		srv.Topics["featureAnnotationUpdate"],
		featProto,
	); err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}

	return featProto, nil
}

func (srv *FeatureAnnotationService) UpdateTag(
	ctx context.Context,
	req *feature.UpdateTagRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	feat, err := srv.repo.UpdateTag(req)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}
	featProto := convertToProto(feat)
	if err := srv.publisher.Publish(
		srv.Topics["featureAnnotationUpdate"],
		featProto,
	); err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}

	return featProto, nil
}

func (srv *FeatureAnnotationService) RemoveTag(
	ctx context.Context,
	req *feature.RemoveTagRequest,
) (*feature.FeatureAnnotation, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}

	err := srv.repo.RemoveTag(req)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}
		return nil, aphgrpc.HandleDeleteError(ctx, err)
	}


	feat, err := srv.repo.GetFeatureAnnotation(req.Id)
	if err != nil {


		if repository.IsAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(
				ctx,
				fmt.Errorf(
					"annotation not found after tag removal: %w",
					err,
				),
			)
		}
		return nil, aphgrpc.HandleGetError(
			ctx,
			fmt.Errorf(
				"failed to fetch annotation after tag removal: %w",
				err,
			),
		)
	}


	featProto := convertToProto(feat)
	if err := srv.publisher.Publish(
		srv.Topics["featureAnnotationUpdate"],
		featProto,
	); err != nil {
		return nil, aphgrpc.HandleUpdateError(ctx, err)
	}

	return featProto, nil
}

func (srv *FeatureAnnotationService) ListFeatureAnnotationsByPubmedId(
	ctx context.Context,
	req *feature.PubmedId,
) (*feature.FeatureAnnotationCollection, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}

	feats, err := srv.repo.ListByPublicationId(req.Id, "pubmed")
	if err != nil {
		if repository.IsPublicationAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}
		return nil, aphgrpc.HandleGetError(ctx, err)
	}

	return &feature.FeatureAnnotationCollection{
		Data: collection.Map(feats, convertToProto),
	}, nil
}

func (srv *FeatureAnnotationService) ListFeatureAnnotationsByDOI(
	ctx context.Context,
	req *feature.DOI,
) (*feature.FeatureAnnotationCollection, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, aphgrpc.HandleInvalidParamError(ctx, err)
	}
	feats, err := srv.repo.ListByPublicationId(req.Id, "doi")
	if err != nil {
		if repository.IsPublicationAnnotationNotFound(err) {
			return nil, aphgrpc.HandleNotFoundError(ctx, err)
		}
		return nil, aphgrpc.HandleGetError(ctx, err)
	}

	return &feature.FeatureAnnotationCollection{
		Data: collection.Map(feats, convertToProto),
	}, nil
}

func convertToProto(
	feat *model.FeatureAnnotationDoc,
) *feature.FeatureAnnotation {
	attrs := &feature.FeatureAnnotationAttributes{Name: feat.Name}


	attrs.Synonyms = feat.Synonyms
	attrs.Publications = feat.Publications
	attrs.Pubmed = feat.Pubmed
	attrs.Dblinks = collection.Map(feat.DbLinks, convertDbLink)
	attrs.Properties = collection.Map(feat.Properties, convertProperty)

	return &feature.FeatureAnnotation{
		Type:       "feature_annotations",
		Id:         feat.AnnoId,
		CreatedBy:  feat.CreatedBy,
		UpdatedBy:  feat.UpdatedBy,
		CreatedAt:  timestamppb.New(feat.CreatedAt),
		UpdatedAt:  timestamppb.New(feat.UpdatedAt),
		IsObsolete: feat.IsObsolete,
		Attributes: attrs,
	}
}

func convertDbLink(link model.DbLinkDoc) *feature.DbLink {
	return &feature.DbLink{
		Database:  link.Database,
		PrimaryId: link.PrimaryId,
		Version:   link.Version,
		Linktype:  link.LinkType,
		Url:       link.URL,
		Label:     link.Label,
	}
}

func convertProperty(prop model.TagPropertyDoc) *feature.TagProperty {
	return &feature.TagProperty{
		Tag:       prop.Tag,
		Value:     prop.Value,
		CreatedBy: prop.CreatedBy,
		UpdatedBy: prop.UpdatedBy,
		CreatedAt: timestamppb.New(prop.CreatedAt),
		UpdatedAt: timestamppb.New(prop.UpdatedAt),
	}
}
```

## File: internal/collection/collection.go
```go
package collection

import (
	"cmp"
	"iter"
	"slices"
)



func Map[T1, T2 any](slc []T1, fnc func(T1) T2) []T2 {
	ret := make([]T2, 0)
	for _, elem := range slc {
		ret = append(ret, fnc(elem))
	}

	return ret
}




func CurriedMap[T1, T2 any](fnc func(T1) T2) func([]T1) []T2 {
	return func(slc []T1) []T2 {
		return Map(slc, fnc)
	}
}



func Include[T cmp.Ordered](slice []T, element T) bool {
	if !slices.IsSorted(slice) {
		slices.Sort(slice)
	}
	_, found := slices.BinarySearch(slice, element)

	return found
}



func RemoveStringItems(slice []string, items ...string) []string {
	str := make([]string, 0)
	for _, val := range slice {
		if !Include(items, val) {
			str = append(str, val)
		}
	}

	return str
}



func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}



func CurriedFilter[T any](predicate func(T) bool) func([]T) []T {
	return func(slice []T) []T {
		return Filter(slice, predicate)
	}
}



func MapSeq[T1, T2 any](seq iter.Seq[T1], fn func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}





func PartitionTuple2[T any](
	slice []T,
	predicate func(T) bool,
) Tuple2[[]T, []T] {
	trueSlice := make([]T, 0)
	falseSlice := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}

	return NewTuple2(trueSlice, falseSlice)
}




func CurriedPartitionTuple2[T any](
	predicate func(T) bool,
) func([]T) Tuple2[[]T, []T] {
	return func(slice []T) Tuple2[[]T, []T] {
		return PartitionTuple2(slice, predicate)
	}
}





func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	trueSlice := make([]T, 0)
	falseSlice := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}

	return trueSlice, falseSlice
}




func CurriedPartition[T any](predicate func(T) bool) func([]T) ([]T, []T) {
	return func(slice []T) ([]T, []T) {
		return Partition(slice, predicate)
	}
}





func Pipe2[T1, T2, T3 any](tup T1, f1 func(T1) T2, fn2 func(T2) T3) T3 {
	return fn2(f1(tup))
}





func Pipe3[T1, T2, T3, T4 any](
	initial T1,
	f1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
) T4 {
	return fn3(fn2(f1(initial)))
}





func Pipe4[T1, T2, T3, T4, T5 any](
	initial T1,
	f1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
) T5 {
	return fn4(fn3(fn2(f1(initial))))
}





func Pipe5[T1, T2, T3, T4, T5, T6 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
) T6 {
	return fn5(fn4(fn3(fn2(fn1(initial)))))
}





func Pipe6[T1, T2, T3, T4, T5, T6, T7 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
	fn6 func(T6) T7,
) T7 {
	return fn6(fn5(fn4(fn3(fn2(fn1(initial))))))
}





func Pipe7[T1, T2, T3, T4, T5, T6, T7, T8 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
	fn6 func(T6) T7,
	fn7 func(T7) T8,
) T8 {
	return fn7(fn6(fn5(fn4(fn3(fn2(fn1(initial)))))))
}





func Pipe8[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
	fn6 func(T6) T7,
	fn7 func(T7) T8,
	fn8 func(T8) T9,
) T9 {
	return fn8(fn7(fn6(fn5(fn4(fn3(fn2(fn1(initial))))))))
}



type Tuple2[T1, T2 any] struct {
	First  T1
	Second T2
}


func NewTuple2[T1, T2 any](first T1, second T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{
		First:  first,
		Second: second,
	}
}




func SliceToTuple2[T1, T2 any](slice []any) Tuple2[T1, T2] {
	var first T1
	var second T2
	if len(slice) > 0 {
		if val, ok := slice[0].(T1); ok {
			first = val
		}
	}
	if len(slice) > 1 {
		if val, ok := slice[1].(T2); ok {
			second = val
		}
	}

	return NewTuple2(first, second)
}




func TFold[A, B, R any](
	tup Tuple2[A, B],
	folder func(Tuple2[A, B]) R,
) R {
	return folder(tup)
}



func CurriedTFold[A, B, R any](
	folder func(Tuple2[A, B]) R,
) func(Tuple2[A, B]) R {
	return func(tup Tuple2[A, B]) R {
		return TFold(tup, folder)
	}
}



func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}



func Sorted[T cmp.Ordered](slice []T) []T {
	sortedSlice := make([]T, len(slice))
	copy(sortedSlice, slice)
	slices.Sort(sortedSlice)

	return sortedSlice
}
```

## File: internal/repository/arangodb/feature_annotation_pipeline.go
```go
package arangodb

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)


type editState struct {
	fann        *featureAnnoRepo
	doc         *feature.FeatureAnnotationUpdate
	txr         *manager.TransactionHandler
	origDoc     *model.FeatureAnnotationDoc
	updatedDoc  *model.FeatureAnnotationDoc
	updateQuery string
	Err         error
}


type repoInitState struct {
	connP       *manager.ConnectParams
	collP       *FeatureCollectionParams
	sess        *manager.Session
	dbh         *manager.Database
	featureColl driver.Collection
	pubColl     driver.Collection
	edgeColl    driver.Collection
	graph       driver.Graph
	Err         error
}


func stepCreateSession(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.sess, state.dbh, state.Err = createSession(state.connP)

	return state
}


func stepCreateFeatureCollection(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.featureColl, state.Err = createFeatureCollection(
		state.dbh,
		state.collP,
	)

	return state
}


func stepCreatePubCollection(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.pubColl, state.Err = createPubCollection(state.dbh, state.collP)

	return state
}


func stepCreateEdgeCollection(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.edgeColl, state.Err = createEdgeCollection(state.dbh, state.collP)

	return state
}


func stepCreateFeatureIndices(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.Err = createFeatureIndices(state.dbh, state.featureColl)

	return state
}


func stepCreatePubIndices(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.Err = createPubIndices(state.dbh, state.pubColl)

	return state
}


func stepCreateGraph(state *repoInitState) *repoInitState {
	if state.Err != nil {
		return state
	}
	state.graph, state.Err = createFeaturePubGraph(
		state.dbh,
		state.collP.Graph,
		state.featureColl,
		state.pubColl,
		state.edgeColl,
	)

	return state
}


func stepFetchOriginalDoc(state *editState) *editState {
	if state.Err != nil {
		return state
	}

	var err error
	state.origDoc, err = state.fann.GetFeatureAnnotation(state.doc.Id)
	if err != nil {
		if repository.IsAnnotationNotFound(err) {
			state.Err = err
		} else {
			state.Err = fmt.Errorf(
				"error fetching original document: %w",
				err,
			)
		}
	}

	return state
}


func stepBeginTransaction(state *editState) *editState {
	if state.Err != nil {
		return state
	}

	state.txr, state.Err = state.fann.database.BeginTransaction(
		context.Background(),
		&manager.TransactionOptions{
			WriteCollections: []string{
				state.fann.feature.Name(),
				state.fann.pub.Name(),
				state.fann.edge.Name(),
			},
		})

	if state.Err != nil {
		state.Err = fmt.Errorf("error beginning transaction: %w", state.Err)
	}

	return state
}


func stepUpdateDocFields(state *editState) *editState {
	if state.Err != nil {
		return state
	}


	updateBasicFields(state.origDoc, state.doc)
	if state.doc.Attributes != nil {
		updateAttributes(state.origDoc, state.doc.Attributes)
	}


	state.updateQuery = fmt.Sprintf(
		"UPDATE @doc WITH @data IN %s RETURN NEW",
		state.fann.feature.Name(),
	)

	return state
}


func stepExecuteUpdate(state *editState) *editState {
	if state.Err != nil {
		return state
	}

	result, err := state.txr.DoRun(
		state.updateQuery,
		map[string]interface{}{
			"doc":  state.origDoc.Key,
			"data": state.origDoc,
		},
	)
	if err != nil {
		state.Err = fmt.Errorf("error updating feature annotation: %w", err)

		return state
	}


	state.updatedDoc = &model.FeatureAnnotationDoc{}
	if err := result.Read(state.updatedDoc); err != nil {
		state.Err = fmt.Errorf("error reading updated document: %w", err)
	}

	return state
}


func stepHandlePublications(state *editState) *editState {
	if state.Err != nil {
		return state
	}


	if state.doc.Attributes == nil {
		return state
	}


	if !collection.IsEmpty(state.doc.Attributes.Pubmed) {
		if err := state.fann.processPublicationType(
			state.txr,
			state.updatedDoc,
			state.doc.Attributes.Pubmed,
			"pubmed",
		); err != nil {
			state.Err = fmt.Errorf(
				"error processing pubmed publications: %w",
				err,
			)

			return state
		}
	}


	if !collection.IsEmpty(state.doc.Attributes.Publications) {
		if err := state.fann.processPublicationType(
			state.txr,
			state.updatedDoc,
			state.doc.Attributes.Publications,
			"doi",
		); err != nil {
			state.Err = fmt.Errorf("error processing DOI publications: %w", err)

			return state
		}
	}

	return state
}


func stepCommitTransaction(state *editState) *editState {

	if state.Err == nil {
		if err := state.txr.Commit(); err != nil {
			state.Err = fmt.Errorf("error committing transaction: %w", err)
		}

		return state
	}

	if state.txr == nil {
		return state
	}


	abortErr := state.txr.Abort()
	if abortErr != nil {
		state.Err = fmt.Errorf(
			"%v, also failed to abort transaction: %w",
			state.Err,
			abortErr,
		)
	}

	return state
}



type featureAnnotationUpdateValidator struct {
	ID        string `validate:"required"       json:"id"`
	UpdatedBy string `validate:"required,email" json:"updated_by"`
}


func stepValidateInput(state *editState) *editState {
	if state.Err != nil {
		return state
	}


	if err := validate.Struct(&featureAnnotationUpdateValidator{
		ID:        state.doc.Id,
		UpdatedBy: state.doc.UpdatedBy,
	}); err != nil {
		state.Err = fmt.Errorf("invalid feature annotation update: %w", err)

		return state
	}

	return state
}



func stepRefreshDocumentState(state *editState) *editState {
	if state.Err != nil {
		return state
	}


	updatedDoc, err := state.fann.GetFeatureAnnotation(state.doc.Id)
	if err != nil {
		state.Err = fmt.Errorf("error retrieving final document state: %w", err)

		return state
	}


	state.updatedDoc = updatedDoc

	return state
}
```

## File: internal/repository/arangodb/statement.go
```go
package arangodb

const (
	orgListQ = `FOR org IN @@collection RETURN org`
	tagGetQ  = `
		FOR cv IN @@cv_collection
			FOR cvt IN @@cvterm_collection
				FILTER cv.metadata.namespace == @ontology
				FILTER cvt.label == @tag
				FILTER cvt.graph_id == cv._id
				LIMIT 1
				RETURN {
					id: cvt.id,
					name: cvt.label,
					is_obsolete: cvt.deprecated,
					ontology: cv.metadata.namespace
				}
	`
	cvtID2LblQ = `
		FOR cvt IN @@cvterm_collection
			FILTER cvt._id == @id
			RETURN cvt.label
	`
	annExistTagQ = `
		FOR cv IN @@cv_collection
			FOR cvt IN @@cvterm_collection
				FILTER cv.metadata.namespace == @ontology
				FILTER cvt.label == @tag || @tag IN cvt.metadata.synonyms[*].value
				FILTER cvt.graph_id == cv._id
				FILTER cvt.deprecated == false
				RETURN cvt._id
	`
	annExistQ = `
		FOR ann IN @@anno_collection
			FOR v IN 1..1 OUTBOUND ann GRAPH @anno_cvterm_graph
				FOR cv IN @@cv_collection
					FILTER ann.entry_id == @entry_id
					FILTER ann.rank == @rank
					FILTER ann.is_obsolete == false
					FILTER v.label == @tag
					FILTER v.deprecated == false
					FILTER v.graph_id == cv._id
					FILTER cv.metadata.namespace == @ontology
					RETURN ann
	`
	annInst = `
		LET n = (
			INSERT {
					value: @value,
					editable_value: @editable_value,
					created_by: @created_by,
					entry_id: @entry_id,
					rank: @rank,
					is_obsolete: false,
					version: @version,
					created_at: DATE_ISO8601(DATE_NOW())
				   } IN @@anno_collection RETURN NEW
		)
		INSERT { _from: n[0]._id, _to: @to } IN @@anno_cv_collection
		RETURN n[0]
	`

	annListWithCursorQ = `
		FOR cvt IN @@cvt_collection
			FOR ann IN 1..1 INBOUND cvt GRAPH @anno_cvterm_graph
				FOR cv IN @@cv_collection
					FILTER ann.is_obsolete == false
					FILTER cvt.graph_id == cv._id
					FILTER ann.created_at <= DATE_ISO8601(@cursor)
					SORT ann.created_at DESC
					LIMIT @limit
						RETURN MERGE(
							ann,
							{ tag: cvt.label, ontology: cv.metadata.namespace }
						)
	`
	annListFilterWithCursorQ = `
		FOR cvt IN @@cvt_collection
			FOR ann IN 1..1 INBOUND cvt GRAPH @anno_cvterm_graph
				FOR cv IN @@cv_collection
					FILTER ann.is_obsolete == false
					FILTER cvt.graph_id == cv._id
					FILTER ann.created_at <= DATE_ISO8601(@cursor)
					%s
					SORT ann.created_at DESC
					LIMIT @limit
						RETURN MERGE(
							ann,
							{ tag: cvt.label, ontology: cv.metadata.namespace }
						)
	`
	annGroupInst = `
		INSERT {
				created_at: DATE_ISO8601(DATE_NOW()),
				updated_at: DATE_ISO8601(DATE_NOW()),
				group: @group
			   } IN @@anno_group_collection RETURN NEW
	`
	annGroupUpd = `
		UPDATE { _key: @key }
			WITH {
					updated_at: DATE_ISO8601(DATE_NOW()),
					group: @group
				 } IN @@anno_group_collection RETURN NEW
	`
	annGroupListFilterWithCursorQ = `
		LET filterannos = (
			FOR ann IN %s
				FOR cvt IN 1..1 OUTBOUND ann GRAPH '%s'
					FOR cv IN %s
						FILTER ann.is_obsolete == false
						FILTER cvt.graph_id == cv._id
						%s
						RETURN ann._key
		)
		FOR ag in %s
			LET annotations = (
				FOR aid in ag.group
					FOR ann IN %s
						FOR cvt IN 1..1 OUTBOUND ann GRAPH '%s'
							FOR cv IN %s
								FILTER aid == ann._key
								FILTER cvt.graph_id == cv._id
								RETURN MERGE(
									ann,
									{ tag: cvt.label, ontology: cv.metadata.namespace }
								)
			)
			FILTER ag.group ANY IN filterannos
			FILTER ag.created_at <= DATE_ISO8601(%d)
			SORT ag.created_at DESC
			LIMIT %d
			RETURN {
				created_at: ag.created_at,
				updated_at: ag.updated_at,
				group_id: ag._key,
				annotations: annotations
			}
	`
	annGroupListQ = `
		FOR ag IN %s
			LET annotations = (
				FOR aid in ag.group
					FOR ann IN %s
						FOR cvt IN 1..1 OUTBOUND ann GRAPH '%s'
							FOR cv IN %s
								FILTER aid == ann._key
								FILTER cvt.graph_id == cv._id
								RETURN MERGE(
									ann,
									{ tag: cvt.label, ontology: cv.metadata.namespace }
								)
			)
			SORT ag.created_at DESC
			LIMIT %d
			RETURN {
				created_at: ag.created_at,
				updated_at: ag.updated_at,
				group_id: ag._key,
				annotations: annotations
			}
	`
	annGroupListWithCursorQ = `
		FOR ag IN %s
			LET annotations = (
				FOR aid in ag.group
					FOR ann IN %s
						FOR cvt IN 1..1 OUTBOUND ann GRAPH '%s'
							FOR cv IN %s
								FILTER aid == ann._key
								FILTER cvt.graph_id == cv._id
								RETURN MERGE(
									ann,
									{ tag: cvt.label, ontology: cv.metadata.namespace }
								)
			)
			FILTER ag.created_at <= DATE_ISO8601(%d)
			SORT ag.created_at DESC
			LIMIT %d
			RETURN {
				created_at: ag.created_at,
				updated_at: ag.updated_at,
				group_id: ag._key,
				annotations: annotations
			}
	`
	annVerInstFn = `
		function (params) {
			var db = require('@arangodb').db
			var d = new Date(Date.now())
			var annoc = db._collection(params[0])
			var n = annoc.save({
				value: params[3],
				editable_value: params[4],
				created_by: params[5],
				entry_id: params[6],
				rank: params[7],
				is_obsolete: false,
				version: params[8],
				created_at: d.toISOString()
			}, { returnNew: true})
			annoc.update(params[10],{ is_obsolete: true })
			db._collection(params[1]).save({
				_from: n._id,
				_to: params[9]
			})
			db._collection(params[2]).save({
				_from: params[10],
				_to: n._id
			})
			return n.new
		}
	`
	annGetQ = `
		FOR ann IN %s
			FOR v IN 1..1 OUTBOUND ann GRAPH '%s'
		        FOR cv IN %s
					FILTER ann._key == '%s'
					FILTER v.graph_id == cv._id
					LIMIT 1
					RETURN MERGE(
						ann,
						{ ontology: cv.metadata.namespace, tag: v.label, cvtid: v._id}
					)
	`
	annGetByEntryQ = `
		FOR ann IN %s
			FOR v IN 1..1 OUTBOUND ann GRAPH '%s'
				FOR cv IN %s
					FILTER ann.entry_id == '%s'
					FILTER ann.rank == %d
					FILTER ann.is_obsolete == %t
					FILTER v.label == '%s'
					FILTER v.graph_id == cv._id
					FILTER cv.metadata.namespace == '%s'
					SORT ann.version DESC
					LIMIT 1
					RETURN MERGE(ann, { ontology: cv.metadata.namespace, tag: v.label })
	`
	orgGetByNameQ = `
		FOR org IN @@collection
			FILTER org.genus == @genus
			FILTER org.species == @species
			LIMIT 1
			RETURN org
	`
	annGroupListFilterQ = `
      LET filterannos = (
		  FOR ann IN %s
		      FOR cvt IN 1..1 OUTBOUND ann GRAPH '%s'
			  FOR cv IN %s
			      FILTER ann.is_obsolete == false
			      FILTER cvt.graph_id == cv._id
			      %s
			      RETURN ann._key
	      )
      FOR ag in %s
          LET annotations = (
              FOR aid in ag.group
                  FOR ann IN %s
                      FOR cvt IN 1..1 OUTBOUND ann GRAPH '%s'
                          FOR cv IN %s
                              FILTER aid == ann._key
                              FILTER cvt.graph_id == cv._id
                              RETURN MERGE(
                                  ann,
                                  { tag: cvt.label, ontology: cv.metadata.namespace }
                              )
          )
          FILTER ag.group ANY IN filterannos
          SORT ag.created_at DESC
          LIMIT %d
          RETURN {
              created_at: ag.created_at,
              updated_at: ag.updated_at,
              group_id: ag._key,
              annotations: annotations
	}
	`
)
```

## File: internal/repository/feature_annotation.go
```go
package repository

import (
	manager "github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
)



type FeatureAnnotationRepository interface {

	GetFeatureAnnotation(id string) (*model.FeatureAnnotationDoc, error)

	GetFeatureAnnotationByName(name string) (*model.FeatureAnnotationDoc, error)

	AddFeatureAnnotation(
		doc *feature.NewFeatureAnnotation,
	) (*model.FeatureAnnotationDoc, error)

	EditFeatureAnnotation(
		doc *feature.FeatureAnnotationUpdate,
	) (*model.FeatureAnnotationDoc, error)

	RemoveFeatureAnnotation(id string, purge bool) error

	ClearFeatureAnnotations() error

	Dbh() *manager.Database

	AddTag(req *feature.AddTagRequest) (*model.FeatureAnnotationDoc, error)
	UpdateTag(
		req *feature.UpdateTagRequest,
	) (*model.FeatureAnnotationDoc, error)
	RemoveTag(
		req *feature.RemoveTagRequest,
	) error

	ListByPublicationId(id string, source string) ([]*model.FeatureAnnotationDoc, error)
}
```

## File: internal/repository/error.go
```go
package repository

import (
	"fmt"
)

type AnnoNotFoundError struct {
	Id string
}

func (ae *AnnoNotFoundError) Error() string {
	return fmt.Sprintf("annotation id %s not found", ae.Id)
}

type GroupNotFoundError struct {
	Id string
}

func (ge *GroupNotFoundError) Error() string {
	return fmt.Sprintf("group id %s not found", ge.Id)
}


type FeatureNameNotFoundError struct {
	Name string
}


func (fnf *FeatureNameNotFoundError) Error() string {
	return fmt.Sprintf("feature annotation with name %s not found", fnf.Name)
}

func IsAnnotationNotFound(err error) bool {
	if _, ok := err.(*AnnoNotFoundError); ok {
		return true
	}

	return false
}


func IsFeatureNameNotFound(err error) bool {
	_, ok := err.(*FeatureNameNotFoundError)
	return ok
}

type PublicationAnnotationNotFoundError struct {
	ID     string
	Source string
}

func (panf *PublicationAnnotationNotFoundError) Error() string {
	return fmt.Sprintf(
		"no annotations found for publication ID %s with source %s",
		panf.ID,
		panf.Source,
	)
}

func IsPublicationAnnotationNotFound(err error) bool {
	_, ok := err.(*PublicationAnnotationNotFoundError)

	return ok
}

func IsGroupNotFound(err error) bool {
	if _, ok := err.(*GroupNotFoundError); ok {
		return true
	}

	return false
}

type AnnoListNotFoundError struct{}

func (al *AnnoListNotFoundError) Error() string {
	return "annotation list not found"
}

func IsAnnotationListNotFound(err error) bool {
	if _, ok := err.(*AnnoListNotFoundError); ok {
		return true
	}

	return false
}

type AnnoGroupListNotFoundError struct{}

func (agl *AnnoGroupListNotFoundError) Error() string {
	return "annotation group list not found"
}

func IsAnnotationGroupListNotFound(err error) bool {
	if _, ok := err.(*AnnoGroupListNotFoundError); ok {
		return true
	}

	return false
}

type AnnoTagNotFoundError struct {
	Tag string
}

func (at *AnnoTagNotFoundError) Error() string {
	return fmt.Sprintf("annotation tag %s not found", at.Tag)
}

func IsAnnoTagNotFound(err error) bool {
	if _, ok := err.(*AnnoTagNotFoundError); ok {
		return true
	}

	return false
}

type OrganismNotFoundError struct {
	ID string
}

func (onf *OrganismNotFoundError) Error() string {
	return fmt.Sprintf("organism id %s not found", onf.ID)
}

func IsOrganismNotFound(err error) bool {
	if _, ok := err.(*OrganismNotFoundError); ok {
		return true
	}

	return false
}

type ListNotFoundError struct{}

func (lnf *ListNotFoundError) Error() string {
	return "list not found"
}

func IsListNotFound(err error) bool {
	if _, ok := err.(*ListNotFoundError); ok {
		return true
	}

	return false
}
```

## File: internal/app/service/feature_annotation_test.go
```go
package service

import (
	"context"
	"testing"
)

func TestCreateFeatureAnnotation(t *testing.T) {
	t.Parallel()
	client, assert := setup(t)
	ctx := context.Background()
	params := &testParams{
		t:      t,
		ctx:    ctx,
		client: client,
		assert: assert,
	}
	testCreateValidFeature(params)
	testCreateMissingFields(params)
	testCreateDuplicateFeature(params)
}

func TestGetFeatureAnnotation(t *testing.T) {
	t.Parallel()
	client, assert := setup(t)
	ctx := context.Background()
	params := &testParams{
		t:      t,
		ctx:    ctx,
		client: client,
		assert: assert,
	}
	testGetExistingFeature(params)
	testGetNonExistentFeature(params)
	testGetFeatureWithInvalidID(params)
}

func TestGetFeatureAnnotationByName(t *testing.T) {
	t.Parallel()
	client, assert := setup(t)
	ctx := context.Background()
	params := &testParams{
		t:      t,
		ctx:    ctx,
		client: client,
		assert: assert,
	}
	testGetExistingFeatureByName(params)
	testGetNonExistentFeatureByName(params)
	testGetFeatureWithEmptyName(params)
}

func TestUpdateFeatureAnnotation(t *testing.T) {
	t.Parallel()
	client, assert := setup(t)
	ctx := context.Background()
	params := &testParams{
		t:      t,
		ctx:    ctx,
		client: client,
		assert: assert,
	}
	testUpdateExistingFeature(params)
	testUpdateNonExistentFeature(params)
	testUpdateWithInvalidData(params)
}

func TestListFeatureAnnotationsByPubmedId(t *testing.T) {
	t.Parallel()
	client, assert := setup(t)
	ctx := context.Background()
	params := &testParams{
		t:      t,
		ctx:    ctx,
		client: client,
		assert: assert,
	}
	testListByPubmedIdValid(params)
	testListByPubmedIdNotFound(params)
	testListByPubmedIdInvalid(params)
}

func TestListFeatureAnnotationsByDOI(t *testing.T) {
	t.Parallel()
	client, assert := setup(t)
	ctx := context.Background()
	params := &testParams{
		t:      t,
		ctx:    ctx,
		client: client,
		assert: assert,
	}
	testListByDOIValid(params)
	testListByDOINotFound(params)
	testListByDOIInvalid(params)
}
```

## File: internal/repository/arangodb/feature_annotation_test_helpers.go
```go
package arangodb

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dictyBase/arangomanager/testarango"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type validateDbLinksParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        []model.DbLinkDoc
	expected   []*feature.DbLink
}

type validatePropertiesParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        []model.TagPropertyDoc
	expected   []*feature.TagProperty
}

type removeFeatureTestCase struct {
	name    string
	purge   bool
	wantErr bool
}

type validateFeatureAnnotationParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        *model.FeatureAnnotationDoc
	base       *feature.NewFeatureAnnotation
}

type validateCompleteFeatureParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        *model.FeatureAnnotationDoc
	expected   *feature.NewFeatureAnnotation
}



type testListByPublicationIdSuccessParams struct {
	t                    *testing.T
	pubID                string
	source               string
	featureIDFieldPrefix string
	setPubFunc           func(*feature.FeatureAnnotationAttributes, []string)
	unrelatedPubID       string
	errorMsgSuffix       string
}

type featFn func() *feature.NewFeatureAnnotation

func getBaseFeatureDoc() *feature.NewFeatureAnnotation {
	return &feature.NewFeatureAnnotation{
		Id:        "DDB_G0000001",
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "base_feature",
		},
	}
}

func getCombinedFeatureDoc(
	baseFn featFn,
	advFn featFn,
) *feature.NewFeatureAnnotation {
	baseDoc := baseFn()
	feat := advFn()
	baseDoc.Attributes = feat.Attributes
	baseDoc.Id = feat.Id

	return baseDoc
}

func getFullFeatureDoc() *feature.NewFeatureAnnotation {
	return &feature.NewFeatureAnnotation{
		Id:        "DDB_G0285425",
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "original name",
			Synonyms: []string{"syn1", "syn2"},
		},
	}
}

func getCompleteFeatureDoc() *feature.NewFeatureAnnotation {
	return &feature.NewFeatureAnnotation{
		Id:        "DDB_G0285425",
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:         "gene name",
			Synonyms:     []string{"synonym1", "synonym2"},
			Publications: []string{"pub1", "pub2"},
			Pubmed:       []string{"123", "456"},
			Dblinks: []*feature.DbLink{
				{
					PrimaryId: "DDB_G0285425",
					Database:  "dictyBase",
					Version:   1,
					Linktype:  "gene",
					Url:       "http://dictybase.org/gene/DDB_G0285425",
					Label:     "gene page",
				},
			},
			Properties: []*feature.TagProperty{
				{
					Tag:       "description",
					Value:     "test gene",
					CreatedBy: "creator3@email.com",
					UpdatedBy: "updater@email.com",
					CreatedAt: timestamppb.New(time.Now()),
					UpdatedAt: timestamppb.New(time.Now()),
				},
			},
		},
	}
}

func getMultiPropertyTestCase() *feature.NewFeatureAnnotation {
	return &feature.NewFeatureAnnotation{
		Id: "DDB_G0285426",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:   "sgene",
			Pubmed: []string{"123456", "456234"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "description",
					Value:     "test description",
					CreatedBy: "creator1@email.com",
					CreatedAt: timestamppb.New(time.Now()),
					UpdatedAt: timestamppb.New(time.Now()),
				},
				{
					Tag:       "note",
					Value:     "test note",
					CreatedBy: "creator2@email.com",
					UpdatedBy: "updater@email.com",
					CreatedAt: timestamppb.New(time.Now()),
				},
				{
					Tag:       "status",
					Value:     "active",
					CreatedBy: "creator3@email.com",
					UpdatedBy: "updater@email.com",
					CreatedAt: timestamppb.New(time.Now()),
					UpdatedAt: timestamppb.New(time.Now()),
				},
			},
		},
	}
}

func getBasicTestCases() []*feature.NewFeatureAnnotation {
	return []*feature.NewFeatureAnnotation{
		{
			Attributes: &feature.FeatureAnnotationAttributes{
				Name: "required fields gene",
			},
			Id: "DDB_G0285428",
		},
		{
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:       "no properties gene",
				Properties: []*feature.TagProperty{},
			},
			Id: "DDB_G0285429",
		},
	}
}

func setUpFeatureTest(
	t *testing.T,
) (*require.Assertions, repository.FeatureAnnotationRepository) {
	t.Helper()
	tra, err := testarango.NewTestArangoFromEnv(true)
	if err != nil {
		t.Fatalf("unable to construct new TestArango instance %s", err)
	}
	assert := require.New(t)
	repo, err := NewFeatureAnnoRepo(
		GetConnectParamsFromDB(tra),
		&FeatureCollectionParams{
			Feature: "feature_test",
			Pub:     "pub_test",
			Edge:    "feature_pub_test",
			Graph:   "feature_graph",
		},
	)
	assert.NoErrorf(
		err,
		"expect no error connecting to feature repository, received %s",
		err,
	)

	return assert, repo
}

func validateProperties(params validatePropertiesParams) {
	params.t.Helper()
	params.assertions.Equal(
		len(params.expected),
		len(params.got),
		"should have same number of properties",
	)
	for idx, prop := range params.expected {
		params.assertions.Equal(
			prop.Tag,
			params.got[idx].Tag,
			"should have matching tag",
		)
		params.assertions.Equal(
			prop.Value,
			params.got[idx].Value,
			"should have matching value",
		)
		params.assertions.Equal(
			prop.CreatedBy,
			params.got[idx].CreatedBy,
			"should have matching creator",
		)
		if prop.UpdatedBy != "" {
			params.assertions.Equal(
				prop.UpdatedBy,
				params.got[idx].UpdatedBy,
				"should have matching updater",
			)
		} else {
			params.assertions.Equal(
				params.got[idx].UpdatedBy,
				params.got[idx].CreatedBy,
				"should match creator and updater",
			)
		}
	}
}




func validateDbLinks(params validateDbLinksParams) {
	params.t.Helper()
	params.assertions.Equal(
		len(params.expected),
		len(params.got),
		"should have same number of dblinks",
	)
	for idx, link := range params.expected {
		params.assertions.Equal(
			link.PrimaryId,
			params.got[idx].PrimaryId,
			"should have matching primary ID",
		)
		params.assertions.Equal(
			link.Database,
			params.got[idx].Database,
			"should have matching database",
		)
		params.assertions.Equal(
			link.Version,
			params.got[idx].Version,
			"should have matching version",
		)
		params.assertions.Equal(
			link.Linktype,
			params.got[idx].LinkType,
			"should have matching link type",
		)
		params.assertions.Equal(
			link.Url,
			params.got[idx].URL,
			"should have matching URL",
		)
		params.assertions.Equal(
			link.Label,
			params.got[idx].Label,
			"should have matching label",
		)
	}
}

func validateBasicFields(params validateFeatureAnnotationParams) {
	params.t.Helper()
	params.assertions.Regexp(
		`^DDB_G\d+`,
		params.got.AnnoId,
		"should have matching IDs",
	)
	params.assertions.Equal(
		params.base.CreatedBy,
		params.got.CreatedBy,
		"should have matching creator",
	)
	params.assertions.Regexp(
		`^[a-zA-Z0-9\s-]*$`,
		params.got.Name,
		"should have matching name",
	)
	params.assertions.Equal(
		params.base.CreatedAt.AsTime(),
		params.got.CreatedAt,
		"should have matching created date",
	)
	params.assertions.Equal(
		params.got.CreatedAt,
		params.got.UpdatedAt,
		"should have matching created and updated at",
	)
}



func validateCompleteFeatureAnnotation(params validateCompleteFeatureParams) {
	params.t.Helper()


	validateBasicFields(validateFeatureAnnotationParams{
		t:          params.t,
		assertions: params.assertions,
		got:        params.got,
		base:       params.expected,
	})


	validateDbLinks(validateDbLinksParams{
		t:          params.t,
		assertions: params.assertions,
		got:        params.got.DbLinks,
		expected:   params.expected.Attributes.Dblinks,
	})


	validateProperties(validatePropertiesParams{
		t:          params.t,
		assertions: params.assertions,
		got:        params.got.Properties,
		expected:   params.expected.Attributes.Properties,
	})


	params.assertions.ElementsMatch(
		params.expected.Attributes.Pubmed,
		params.got.Pubmed,
		"should match pubmed ids",
	)


	params.assertions.ElementsMatch(
		params.expected.Attributes.Publications,
		params.got.Publications,
		"should match publications",
	)
}

func sortTagProperties(a, b model.TagPropertyDoc) int {
	return strings.Compare(
		strings.ToLower(a.Tag),
		strings.ToLower(b.Tag),
	)
}

func getRemoveTestCases() []removeFeatureTestCase {
	return []removeFeatureTestCase{
		{
			name:    "should soft delete feature annotation",
			purge:   false,
			wantErr: false,
		},
		{
			name:    "should purge feature annotation",
			purge:   true,
			wantErr: false,
		},
		{
			name:    "should return error for non-existent ID",
			purge:   false,
			wantErr: true,
		},
	}
}

func assertListByPublicationResults(
	t *testing.T,
	asrt *require.Assertions,
	results []*model.FeatureAnnotationDoc,
	added1 *model.FeatureAnnotationDoc,
	added2 *model.FeatureAnnotationDoc,
) {
	t.Helper()
	asrt.Len(results, 2, "Should retrieve exactly 2 feature annotations")
	retrievedIDs := collection.Map(
		results,
		func(doc *model.FeatureAnnotationDoc) string {
			return doc.AnnoId
		},
	)
	expectedIDs := []string{added1.AnnoId, added2.AnnoId}
	slices.Sort(retrievedIDs)
	slices.Sort(expectedIDs)
	asrt.Equal(
		expectedIDs,
		retrievedIDs,
		"Retrieved feature IDs should match the linked ones",
	)
}

func testListByPublicationIdSuccess(
	params *testListByPublicationIdSuccessParams,
) {
	params.t.Helper()
	asrt, repo := setUpFeatureTest(params.t)
	params.t.Cleanup(cleanupDB(repo))


	feat1 := getFullFeatureDoc()
	feat1.Id = params.featureIDFieldPrefix + "1"
	params.setPubFunc(feat1.Attributes, []string{params.pubID})
	added1, err := repo.AddFeatureAnnotation(feat1)
	asrt.NoError(err, "Failed to add feature 1")

	feat2 := getFullFeatureDoc()
	feat2.Id = params.featureIDFieldPrefix + "2"
	params.setPubFunc(feat2.Attributes, []string{params.pubID})
	added2, err := repo.AddFeatureAnnotation(feat2)
	asrt.NoError(err, "Failed to add feature 2")


	feat3 := getFullFeatureDoc()
	feat3.Id = params.featureIDFieldPrefix + "3"
	params.setPubFunc(
		feat3.Attributes,
		[]string{params.unrelatedPubID},
	)
	_, err = repo.AddFeatureAnnotation(feat3)
	asrt.NoError(err, "Failed to add unrelated feature 3")


	results, err := repo.ListByPublicationId(params.pubID, params.source)
	asrt.NoError(err, "Expected no error retrieving by "+params.errorMsgSuffix)


	assertListByPublicationResults(params.t, asrt, results, added1, added2)
}

func cleanupDB(repo repository.FeatureAnnotationRepository) func() {
	return func() {
		_ = repo.Dbh().Drop()
	}
}
```

## File: internal/app/service/feature_annotation_test_helpers.go
```go
package service

import (
	"context"
	"net"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/dictyBase/arangomanager/testarango"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/repository/arangodb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)


type assertGrpcErrorParams struct {
	assert               *require.Assertions
	err                  error
	expectedCode         codes.Code
	expectedMsgSubstring string
}

type testParams struct {
	t      *testing.T
	ctx    context.Context
	client feature.FeatureAnnotationServiceClient
	assert *require.Assertions
}

type MockMessage struct{}

func (msn *MockMessage) Publish(
	subject string,
	feat *feature.FeatureAnnotation,
) error {
	return nil
}

func (msn *MockMessage) Close() error {
	return nil
}



func sortTagPropertiesByTag(a, b *feature.TagProperty) int {
	return strings.Compare(
		strings.ToLower(a.Tag),
		strings.ToLower(b.Tag),
	)
}


func extractTagAndValue(prop *feature.TagProperty) *feature.TagProperty {
	return &feature.TagProperty{
		Tag:   prop.Tag,
		Value: prop.Value,
	}
}

func setup(
	t *testing.T,
) (feature.FeatureAnnotationServiceClient, *require.Assertions) {
	t.Helper()
	assert := require.New(t)
	tra, err := testarango.NewTestArangoFromEnv(true)
	assert.NoError(err, "expect no error from creating an arangodb instance")
	repo, err := arangodb.NewFeatureAnnoRepo(
		arangodb.GetConnectParamsFromDB(tra),
		&arangodb.FeatureCollectionParams{
			Feature: "feature_test",
			Pub:     "pub_test",
			Edge:    "feature_pub_test",
			Graph:   "feature_pub_graph_test",
		},
	)
	assert.NoErrorf(
		err,
		"expect no error connecting to annotation repository, received %s",
		err,
	)

	svc, err := NewFeatureAnnotationService(&FeatureParams{
		Repository: repo,
		Publisher:  &MockMessage{},
	})
	assert.NoError(err)


	server := grpc.NewServer()
	feature.RegisterFeatureAnnotationServiceServer(server, svc)
	lis := bufconn.Listen(1024 * 1024)
	go func() {
		if err := server.Serve(lis); err != nil {
			t.Logf("Server exited with error: %v", err)
			os.Exit(1)
		}
	}()
	dialer := func(context.Context, string) (net.Conn, error) {
		conn, err := lis.Dial()
		assert.NoError(err, "expect no error from creating listener")

		return conn, nil
	}
	resolver.SetDefaultScheme("passthrough")

	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
	)
	assert.NoError(err)
	t.Cleanup(func() {
		_ = repo.Dbh().Drop()
		conn.Close()
		lis.Close()
		server.Stop()
	})

	return feature.NewFeatureAnnotationServiceClient(conn), assert
}


func newTestFeature() *feature.NewFeatureAnnotation {
	return &feature.NewFeatureAnnotation{
		Id:        "DDB_G0285425",
		CreatedBy: "testuser@dictybase.org",
		CreatedAt: timestamppb.Now(),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "Test Feature",
			Synonyms: []string{"test1", "test2"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "description",
					Value:     "Test description",
					CreatedBy: "testuser@dictybase.org",
				},
				{
					Tag:       "note",
					Value:     "Test note",
					CreatedBy: "testuser@dictybase.org",
				},
			},

		},
	}
}

func testCreateValidFeature(params *testParams) {
	params.t.Helper()
	params.t.Run("CreateValidFeatureAnnotation", func(t *testing.T) {



		req := newTestFeature()

		req.CreatedAt = timestamppb.Now()

		resp, err := params.client.CreateFeatureAnnotation(params.ctx, req)
		params.assert.NoError(err)
		params.assert.Equal(req.Id, resp.Id)
		params.assert.Equal(req.CreatedBy, resp.CreatedBy)
		params.assert.Equal(req.Attributes.Name, resp.Attributes.Name)
		params.assert.Equal(req.Attributes.Synonyms, resp.Attributes.Synonyms)


		params.assert.Len(resp.Attributes.Properties, 2)
		slices.SortFunc(req.Attributes.Properties, sortTagPropertiesByTag)
		slices.SortFunc(resp.Attributes.Properties, sortTagPropertiesByTag)
		params.assert.ElementsMatch(
			collection.Map(req.Attributes.Properties, extractTagAndValue),
			collection.Map(resp.Attributes.Properties, extractTagAndValue),
			"should have matching properties",
		)
	})
}


func testListByPublicationHelper(
	params *testParams,
	publicationType string,
	publicationID string,
	featureID1 string,
	featureID2 string,
	featureNamePrefix string,
) {
	params.t.Helper()
	feat1 := &feature.NewFeatureAnnotation{
		Id:        featureID1,
		CreatedBy: "testuser@dictybase.org",
		CreatedAt: timestamppb.Now(),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: featureNamePrefix,
		},
	}
	feat2 := &feature.NewFeatureAnnotation{
		Id:        featureID2,
		CreatedBy: "testuser@dictybase.org",
		CreatedAt: timestamppb.Now(),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: featureNamePrefix,
		},
	}

	switch publicationType {
	case "doi":
		feat1.Attributes.Publications = []string{publicationID}
		feat2.Attributes.Publications = []string{publicationID}
	case "pubmed":
		feat1.Attributes.Pubmed = []string{publicationID}
		feat2.Attributes.Pubmed = []string{publicationID}
	default:
		params.t.Fatalf("invalid publication type: %s", publicationType)
	}

	_, err := params.client.CreateFeatureAnnotation(params.ctx, feat1)
	params.assert.NoError(err)
	_, err = params.client.CreateFeatureAnnotation(params.ctx, feat2)
	params.assert.NoError(err)

	var resp *feature.FeatureAnnotationCollection

	switch publicationType {
	case "doi":
		req := &feature.DOI{Id: publicationID}
		resp, err = params.client.ListFeatureAnnotationsByDOI(params.ctx, req)
	case "pubmed":
		req := &feature.PubmedId{Id: publicationID}
		resp, err = params.client.ListFeatureAnnotationsByPubmedId(
			params.ctx,
			req,
		)
	}

	params.assert.NoError(err)
	params.assert.Len(resp.Data, 2)

	foundIDs := []string{
		resp.Data[0].Id,
		resp.Data[1].Id,
	}
	params.assert.Contains(foundIDs, feat1.Id)
	params.assert.Contains(foundIDs, feat2.Id)
}

func testListByDOIValid(params *testParams) {
	params.t.Helper()
	params.t.Run("ListByDOIValid", func(t *testing.T) {
		t.Parallel()
		testListByPublicationHelper(
			params,
			"doi",
			"10.1234/j.abcd.2023.01.001",
			"DDB_G0285430",
			"DDB_G0285431",
			"Feature DOI",
		)
	})
}

func testListByDOINotFound(params *testParams) {
	params.t.Helper()
	params.t.Run("ListByDOINotFound", func(t *testing.T) {
		t.Parallel()
		req := &feature.DOI{Id: "10.9999/non.existent.doi"}
		_, err := params.client.ListFeatureAnnotationsByDOI(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.NotFound, sts.Code())
	})
}

func testListByDOIInvalid(params *testParams) {
	params.t.Helper()
	params.t.Run("ListByDOIInvalid", func(t *testing.T) {
		t.Parallel()
		req := &feature.DOI{Id: ""} // Invalid (empty) DOI
		_, err := params.client.ListFeatureAnnotationsByDOI(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.InvalidArgument, sts.Code())
	})
}

func testListByPubmedIdValid(params *testParams) {
	params.t.Helper()
	params.t.Run("ListByPubmedIdValid", func(t *testing.T) {
		t.Parallel()
		testListByPublicationHelper(
			params,
			"pubmed",
			"12345678",
			"DDB_G0285428",
			"DDB_G0285429",
			"Feature Pubmed",
		)
	})
}

func testListByPubmedIdNotFound(params *testParams) {
	params.t.Helper()
	params.t.Run("ListByPubmedIdNotFound", func(t *testing.T) {
		t.Parallel()
		req := &feature.PubmedId{Id: "99999999"}
		_, err := params.client.ListFeatureAnnotationsByPubmedId(
			params.ctx,
			req,
		)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.NotFound, sts.Code())
	})
}

func testListByPubmedIdInvalid(params *testParams) {
	params.t.Helper()
	params.t.Run("ListByPubmedIdInvalid", func(t *testing.T) {
		t.Parallel()
		req := &feature.PubmedId{Id: ""} // Invalid (empty) pubmed ID
		_, err := params.client.ListFeatureAnnotationsByPubmedId(
			params.ctx,
			req,
		)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.InvalidArgument, sts.Code())
	})
}

func testCreateMissingFields(params *testParams) {
	params.t.Helper()
	params.t.Run("CreateFailsMissingRequiredFields", func(t *testing.T) {
		t.Parallel()
		req := &feature.NewFeatureAnnotation{
			Attributes: &feature.FeatureAnnotationAttributes{
				Name: "Invalid Feature",
			},
		}
		_, err := params.client.CreateFeatureAnnotation(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.InvalidArgument, sts.Code())
	})
}

func testCreateDuplicateFeature(params *testParams) {
	params.t.Helper()
	params.t.Run("CreateFailsDuplicateFeatureId", func(t *testing.T) {
		t.Parallel()
		req := &feature.NewFeatureAnnotation{
			Id:        "DDB_G02854297",
			CreatedBy: "testuser@dictybase.org",
			CreatedAt: timestamppb.Now(),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name: "Duplicate Feature",
			},
		}
		_, firstErr := params.client.CreateFeatureAnnotation(params.ctx, req)
		params.assert.NoError(firstErr)
		_, dupErr := params.client.CreateFeatureAnnotation(params.ctx, req)
		params.assert.Error(dupErr)
		sts, ok := status.FromError(dupErr)
		params.assert.True(ok)
		params.assert.Equal(codes.AlreadyExists, sts.Code())
	})
}

func testGetExistingFeature(params *testParams) {
	params.t.Helper()
	params.t.Run("GetExistingFeatureAnnotation", func(t *testing.T) {
		t.Parallel()

		createReq := &feature.NewFeatureAnnotation{
			Id:        "DDB_G0285426",
			CreatedBy: "testuser@dictybase.org",
			CreatedAt: timestamppb.Now(),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:     "Test Feature",
				Synonyms: []string{"test1", "test2"},
			},
		}
		_, err := params.client.CreateFeatureAnnotation(params.ctx, createReq)
		params.assert.NoError(err)


		getReq := &feature.FeatureAnnotationId{
			Id: "DDB_G0285426",
		}
		resp, err := params.client.GetFeatureAnnotation(params.ctx, getReq)
		params.assert.NoError(err)
		params.assert.Equal(createReq.Id, resp.Id)
		params.assert.Equal(createReq.CreatedBy, resp.CreatedBy)
		params.assert.Equal(createReq.Attributes.Name, resp.Attributes.Name)
		params.assert.Equal(
			createReq.Attributes.Synonyms,
			resp.Attributes.Synonyms,
		)
	})
}

func testGetNonExistentFeature(params *testParams) {
	params.t.Helper()
	params.t.Run("GetNonExistentFeatureAnnotation", func(t *testing.T) {
		t.Parallel()
		req := &feature.FeatureAnnotationId{
			Id: "DDB_G0000000",
		}
		_, err := params.client.GetFeatureAnnotation(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.NotFound, sts.Code())
	})
}

func testGetFeatureWithInvalidID(params *testParams) {
	params.t.Helper()
	params.t.Run("GetFeatureAnnotationWithInvalidID", func(t *testing.T) {
		t.Parallel()
		req := &feature.FeatureAnnotationId{
			Id: "", // Empty ID
		}
		_, err := params.client.GetFeatureAnnotation(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.InvalidArgument, sts.Code())
	})
}

func testUpdateExistingFeature(params *testParams) {
	params.t.Helper()
	params.t.Run("UpdateExistingFeatureAnnotation", func(t *testing.T) {
		t.Parallel()

		createReq := &feature.NewFeatureAnnotation{
			Id:        "DDB_G0285427",
			CreatedBy: "testuser@dictybase.org",
			CreatedAt: timestamppb.Now(),
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:     "Original Feature",
				Synonyms: []string{"orig1", "orig2"},
			},
		}
		_, err := params.client.CreateFeatureAnnotation(params.ctx, createReq)
		params.assert.NoError(err)


		updateReq := &feature.FeatureAnnotationUpdate{
			Id:        "DDB_G0285427",
			UpdatedBy: "anotheruser@dictybase.org",
			Attributes: &feature.FeatureAnnotationAttributes{
				Name:     "Updated Feature",
				Synonyms: []string{"new1", "new2"},
			},
		}
		resp, err := params.client.UpdateFeatureAnnotation(
			params.ctx,
			updateReq,
		)
		params.assert.NoError(err)
		params.assert.Equal(updateReq.Id, resp.Id)
		params.assert.Equal(updateReq.UpdatedBy, resp.UpdatedBy)
		params.assert.Equal(updateReq.Attributes.Name, resp.Attributes.Name)
		params.assert.ElementsMatch(
			slices.Concat(
				createReq.Attributes.Synonyms,
				updateReq.Attributes.Synonyms,
			),
			resp.Attributes.Synonyms,
		)
		params.assert.Equal(createReq.CreatedBy, resp.CreatedBy)
	})
}

func testUpdateNonExistentFeature(params *testParams) {
	params.t.Helper()
	params.t.Run("UpdateNonExistentFeatureAnnotation", func(t *testing.T) {
		t.Parallel()
		req := &feature.FeatureAnnotationUpdate{
			Id:        "DDB_G0000000",
			UpdatedBy: "testuser@dictybase.org",
			Attributes: &feature.FeatureAnnotationAttributes{
				Name: "Non-existent Feature",
			},
		}
		_, err := params.client.UpdateFeatureAnnotation(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		t.Log(sts.Code().String())
		params.assert.Equal(codes.Internal, sts.Code())
	})
}

func testUpdateWithInvalidData(params *testParams) {
	params.t.Helper()
	params.t.Run("UpdateFeatureAnnotationWithInvalidData", func(t *testing.T) {
		t.Parallel()
		req := &feature.FeatureAnnotationUpdate{
			Id: "", // Empty ID
			Attributes: &feature.FeatureAnnotationAttributes{
				Name: "Invalid Feature",
			},
		}
		_, err := params.client.UpdateFeatureAnnotation(params.ctx, req)
		params.assert.Error(err)
		sts, ok := status.FromError(err)
		params.assert.True(ok)
		params.assert.Equal(codes.InvalidArgument, sts.Code())
	})
}

func testGetExistingFeatureByName(params *testParams) {
	testFeatureData := newTestFeature()
	testCreateValidFeature(
		params,
	)


	featureName := testFeatureData.Attributes.Name
	req := &feature.FeatureName{Name: featureName}
	gotFeat, err := params.client.GetFeatureAnnotationByName(params.ctx, req)

	params.assert.NoError(
		err,
		"should retrieve existing feature by name without error",
	)
	params.assert.Equal(
		featureName,
		gotFeat.Attributes.Name,
		"retrieved feature name should match the known name",
	)
	params.assert.Equal(
		testFeatureData.Id,
		gotFeat.Id,
		"retrieved feature entry_id should match",
	)
	slices.SortFunc(
		testFeatureData.Attributes.Properties,
		sortTagPropertiesByTag,
	)
	slices.SortFunc(gotFeat.Attributes.Properties, sortTagPropertiesByTag)
	params.assert.ElementsMatch(
		collection.Map(
			testFeatureData.Attributes.Properties,
			extractTagAndValue,
		),
		collection.Map(gotFeat.Attributes.Properties, extractTagAndValue),
		"should have matching properties",
	)
}

func testGetNonExistentFeatureByName(params *testParams) {
	nonExistentName := "this_feature_does_not_exist_12345"
	req := &feature.FeatureName{Name: nonExistentName}
	_, err := params.client.GetFeatureAnnotationByName(params.ctx, req)

	params.assert.Error(
		err,
		"should return an error for non-existent feature name",
	)
	assertGrpcError(assertGrpcErrorParams{
		assert:               params.assert,
		err:                  err,
		expectedCode:         codes.NotFound,
		expectedMsgSubstring: "not found",
	})
}



func assertGrpcError(params assertGrpcErrorParams) {
	params.assert.Error(params.err, "expected a gRPC error")
	sts, ok := status.FromError(params.err)
	params.assert.True(ok, "error should be a gRPC status error")
	params.assert.Equal(
		params.expectedCode,
		sts.Code(),
		"expected gRPC code %s, but got %s",
		params.expectedCode,
		sts.Code(),
	)
	if params.expectedMsgSubstring != "" {
		params.assert.Contains(
			strings.ToLower(sts.Message()), // Case-insensitive check
			strings.ToLower(params.expectedMsgSubstring),
			"expected gRPC error message to contain '%s', but got '%s'",
			params.expectedMsgSubstring,
			sts.Message(),
		)
	}
}

func testGetFeatureWithEmptyName(params *testParams) {
	req := &feature.FeatureName{Name: ""} // Empty name
	_, err := params.client.GetFeatureAnnotationByName(params.ctx, req)

	params.assert.Error(err, "should return an error for empty feature name")
	assertGrpcError(assertGrpcErrorParams{
		assert:               params.assert,
		err:                  err,
		expectedCode:         codes.InvalidArgument,
		expectedMsgSubstring: "validation",
	})
}
```

## File: internal/repository/arangodb/feature_annotation_helpers.go
```go
package arangodb

import (
	"fmt"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
)


type CreateIndexArgs struct {
	Dbh          *manager.Database
	Coll         driver.Collection
	Fields       []string
	UniqueFields []string
	ErrPrefix    string
}

func createSession(
	connP *manager.ConnectParams,
) (*manager.Session, *manager.Database, error) {
	sess, dbh, err := manager.NewSessionDb(connP)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to create database session: %w",
			err,
		)
	}

	return sess, dbh, nil
}

func createFeatureCollection(
	dbh *manager.Database,
	collP *FeatureCollectionParams,
) (driver.Collection, error) {
	schema, err := model.FeatureAnnotationSchema()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate feature annotation schema: %w",
			err,
		)
	}

	schemaOpt := &driver.CollectionSchemaOptions{
		Level:   driver.CollectionSchemaLevelModerate,
		Message: "Feature annotation validation failed",
		Type:    "json",
	}
	if err := schemaOpt.LoadRule(schema); err != nil {
		return nil, fmt.Errorf("error in loading schema %s", err)
	}
	coll, err := dbh.FindOrCreateCollection(
		collP.Feature,
		&driver.CreateCollectionOptions{
			Schema: schemaOpt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create/find feature collection: %w",
			err,
		)
	}

	return coll, nil
}


func createIndices(args *CreateIndexArgs) error {

	for _, field := range args.UniqueFields {
		_, _, err := args.Dbh.EnsurePersistentIndex(
			args.Coll.Name(),
			[]string{field},
			&driver.EnsurePersistentIndexOptions{
				InBackground: true,
				Unique:       true,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to create unique %s index for %s: %w",
				field,
				args.ErrPrefix,
				err,
			)
		}
	}


	for _, field := range args.Fields {
		_, _, err := args.Dbh.EnsurePersistentIndex(
			args.Coll.Name(),
			[]string{field},
			&driver.EnsurePersistentIndexOptions{
				InBackground: true,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to create %s index for %s: %w",
				field,
				args.ErrPrefix,
				err,
			)
		}
	}

	return nil
}

func createFeatureIndices(dbh *manager.Database, coll driver.Collection) error {
	return createIndices(&CreateIndexArgs{
		Dbh:          dbh,
		Coll:         coll,
		Fields:       []string{"name"},
		UniqueFields: []string{"feature_id"},
		ErrPrefix:    "feature collection",
	})
}

func createPubIndices(dbh *manager.Database, coll driver.Collection) error {
	return createIndices(&CreateIndexArgs{
		Dbh:          dbh,
		Coll:         coll,
		UniqueFields: []string{"id"},
		ErrPrefix:    "pub collection",
	})
}

func createPubCollection(
	dbh *manager.Database,
	collP *FeatureCollectionParams,
) (driver.Collection, error) {
	schema, err := model.PubSchema()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate pub schema: %w",
			err,
		)
	}

	schemaOpt := &driver.CollectionSchemaOptions{
		Level:   driver.CollectionSchemaLevelModerate,
		Message: "Pub validation failed",
		Type:    "json",
	}
	if err := schemaOpt.LoadRule(schema); err != nil {
		return nil, fmt.Errorf("error in loading pub schema %s", err)
	}
	coll, err := dbh.FindOrCreateCollection(
		collP.Pub,
		&driver.CreateCollectionOptions{
			Schema: schemaOpt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create/find pub collection: %w",
			err,
		)
	}

	return coll, nil
}

func createEdgeCollection(
	dbh *manager.Database,
	collP *FeatureCollectionParams,
) (driver.Collection, error) {

	coll, err := dbh.FindOrCreateCollection(
		collP.Edge,
		&driver.CreateCollectionOptions{
			Type: driver.CollectionTypeEdge,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create/find edge collection: %w",
			err,
		)
	}

	return coll, nil
}

func createFeaturePubGraph(
	dbh *manager.Database,
	graphName string,
	featureColl driver.Collection,
	pubColl driver.Collection,
	edgeColl driver.Collection,
) (driver.Graph, error) {
	grph, err := dbh.FindOrCreateGraph(
		graphName,
		[]driver.EdgeDefinition{
			{
				Collection: edgeColl.Name(),
				From:       []string{featureColl.Name()},
				To:         []string{pubColl.Name()},
			},
		})
	if err != nil {
		return grph, fmt.Errorf(
			"failed to create/find graph %s: %w",
			graphName,
			err,
		)
	}

	return grph, nil
}

func updateBasicFields(
	faDoc *model.FeatureAnnotationDoc,
	doc *feature.FeatureAnnotationUpdate,
) {
	faDoc.UpdatedBy = doc.UpdatedBy
	if faDoc.IsObsolete != doc.IsObsolete {
		faDoc.IsObsolete = doc.IsObsolete
	}
}

func updateAttributes(
	mdoc *model.FeatureAnnotationDoc,
	attrs *feature.FeatureAnnotationAttributes,
) {
	mdoc.Name = attrs.Name
	mdoc.Synonyms = append(mdoc.Synonyms, attrs.Synonyms...)
	mdoc.DbLinks = append(
		mdoc.DbLinks,
		collection.Map(attrs.Dblinks, convertDbLink)...)
	mdoc.Properties = append(
		mdoc.Properties,
		collection.Map(attrs.Properties, convertProperty)...)
}

func convertDbLink(link *feature.DbLink) model.DbLinkDoc {
	return model.DbLinkDoc{
		PrimaryId: link.PrimaryId,
		Database:  link.Database,
		Version:   link.Version,
		LinkType:  link.Linktype,
		URL:       link.Url,
		Label:     link.Label,
	}
}

func convertProperty(prop *feature.TagProperty) model.TagPropertyDoc {
	mprop := model.TagPropertyDoc{
		Tag:       prop.Tag,
		Value:     prop.Value,
		CreatedBy: prop.CreatedBy,
		UpdatedBy: prop.CreatedBy,
		CreatedAt: prop.CreatedAt.AsTime(),
		UpdatedAt: prop.CreatedAt.AsTime(),
	}
	if len(prop.UpdatedBy) != 0 {
		mprop.UpdatedBy = prop.UpdatedBy
	}
	if prop.UpdatedAt != nil {
		mprop.UpdatedAt = prop.UpdatedAt.AsTime()
	}

	return mprop
}

func setOptionalFields(
	doc *feature.NewFeatureAnnotation,
	faDoc *model.FeatureAnnotationDoc,
) *model.FeatureAnnotationDoc {
	faDoc.Synonyms = doc.Attributes.Synonyms
	faDoc.DbLinks = collection.Map(doc.Attributes.Dblinks, convertDbLink)
	faDoc.Properties = collection.Map(
		doc.Attributes.Properties,
		convertProperty,
	)

	return faDoc
}
```

## File: internal/repository/arangodb/feature_statement.go
```go
package arangodb

const (





	featurePubEdgeQ = `
		FOR pub_key IN @pub_keys
			INSERT {
				_from: @feature_key,
				_to: pub_key,
				source: @source
			} IN @@edge_collection
	`






	pubUpsertQ = `
		LET allKeys = (
			FOR id_val IN @ids
				UPSERT { id: id_val }
				INSERT {
					id: id_val,
					created_at: DATE_ISO8601(DATE_NOW()),
					updated_at: DATE_ISO8601(DATE_NOW())
				}
				UPDATE { updated_at: DATE_ISO8601(DATE_NOW()) }
				IN @@collection
				RETURN NEW._id
		)
		RETURN allKeys
	`

	featureExistQ = `
		FOR f IN @@collection
			FILTER f.feature_id == @id
			LIMIT 1
			RETURN f._key
	`
	featurePurgeQ = `
        FOR f IN @@collection
            FILTER f.feature_id == @id
            LIMIT 1
            REMOVE f IN @@collection
    `
	featureObsoleteQ = `
        FOR f IN @@collection
            FILTER f.feature_id == @id
	    UPDATE f WITH { is_obsolete: true } IN @@collection
    `
	featureGetByIdQ = `
	FOR f IN @@collection
	    FILTER f.feature_id == @id
	    FILTER f.is_obsolete == false
	    LIMIT 1
	    LET pubmed = (
		FOR v,e IN 1..1 OUTBOUND f GRAPH @graph
		FILTER e.source == 'pubmed'
		RETURN v.id
	    )
	    LET doi = (
		FOR v,e IN 1..1 OUTBOUND f GRAPH @graph
		FILTER e.source == 'doi'
		RETURN v.id
	    )
	    RETURN MERGE(f, {pubmed: pubmed, publications: doi})
    `

	featAnnoGetByNameQ = `
	FOR f IN @@collection
	    FILTER f.name == @name
	    FILTER f.is_obsolete == false
	    LIMIT 1
	    LET pubmed = (
		FOR v,e IN 1..1 OUTBOUND f GRAPH @graph
		FILTER e.source == 'pubmed'
		RETURN v.id
	    )
	    LET doi = (
		FOR v,e IN 1..1 OUTBOUND f GRAPH @graph
		FILTER e.source == 'doi'
		RETURN v.id
	    )
	    RETURN MERGE(f, {pubmed: pubmed, publications: doi})
	`

	featureByPublicationIdQ = `
	FOR pub IN @@collection
    		FOR v,e IN 1..1 INBOUND pub GRAPH @graph
        	FILTER pub.id == @id
        	FILTER e.source == @source
        	FILTER v.is_obsolete == false
        	RETURN v
   `
)
```

## File: internal/repository/arangodb/feature_annotation_test.go
```go
package arangodb

import (
	"slices"
	"testing"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetFeatureAnnotation(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "expected no error adding test feature annotation")

	got, err := repo.GetFeatureAnnotation(added.AnnoId)
	asrt.NoError(err, "expected no error getting feature annotation")


	validateCompleteFeatureAnnotation(validateCompleteFeatureParams{
		t:          t,
		assertions: asrt,
		got:        got,
		expected:   feat,
	})


	asrt.NotEmpty(
		got.Pubmed,
		"should have pubmed ids in result for this test case")
	asrt.NotEmpty(
		got.Publications,
		"should have publications in result for this test case",
	)

	_, err = repo.GetFeatureAnnotation("non_existent_id")
	asrt.Error(err, "expected error for non-existent feature annotation")
	asrt.True(
		repository.IsAnnotationNotFound(err),
		"should be annotation not found error",
	)
}

func TestGetFeatureAnnotationByName(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(
		err,
		"expected no error adding test feature annotation for name lookup",
	)


	gotByName, err := repo.GetFeatureAnnotationByName(feat.Attributes.Name)
	asrt.NoError(
		err,
		"expected no error getting feature annotation by name",
	)

	asrt.Equal(
		added.AnnoId,
		gotByName.AnnoId,
		"retrieved document ID should match added document ID",
	)
	asrt.Equal(
		feat.Attributes.Name,
		gotByName.Name,
		"retrieved document name should match",
	)

	validateCompleteFeatureAnnotation(validateCompleteFeatureParams{
		t:          t,
		assertions: asrt,
		got:        gotByName,
		expected:   feat,
	})


	nonExistentName := "non_existent_feature_name"
	_, err = repo.GetFeatureAnnotationByName(nonExistentName)
	asrt.Error(
		err,
		"expected error for non-existent feature annotation name",
	)
	asrt.True(
		repository.IsFeatureNameNotFound(err),
		"should be FeatureNameNotFoundError",
	)
}

func TestAddFeatureAnnotationBasic(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))
	baseFeat := &feature.NewFeatureAnnotation{
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
	}
	for _, nfeat := range getBasicTestCases() {
		baseFeat.Attributes = nfeat.Attributes
		baseFeat.Id = nfeat.Id
		doc, err := repo.AddFeatureAnnotation(baseFeat)
		asrt.NoError(err, "expected no error adding feature annotation")
		validateBasicFields(validateFeatureAnnotationParams{
			t:          t,
			assertions: asrt,
			got:        doc,
			base:       baseFeat,
		})
	}
}

func TestAddFeatureAnnotationFull(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	feat := getCompleteFeatureDoc()
	doc, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "expected no error adding feature annotation")


	validateCompleteFeatureAnnotation(validateCompleteFeatureParams{
		t:          t,
		assertions: asrt,
		got:        doc,
		expected:   feat,
	})


	asrt.NotEmpty(
		doc.Pubmed,
		"should have pubmed ids in result for this test case")
	asrt.NotEmpty(
		doc.Publications,
		"should have publications in result for this test case",
	)
}

func TestAddFeatureAnnotationMultiProperty(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))
	feat := getCombinedFeatureDoc(getBaseFeatureDoc, getMultiPropertyTestCase)
	doc, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "expected no error adding feature annotation")
	validateBasicFields(validateFeatureAnnotationParams{
		t:          t,
		assertions: asrt,
		got:        doc,
		base:       feat,
	})
	validateProperties(validatePropertiesParams{
		t:          t,
		assertions: asrt,
		got:        doc.Properties,
		expected:   feat.Attributes.Properties,
	})
}

func TestAddDuplicateFeatureAnnotation(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := &feature.NewFeatureAnnotation{
		Id:        "DDB_G0285425",
		CreatedBy: "mock@email.com",
		CreatedAt: timestamppb.New(time.Now()),
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "gene name",
		},
	}


	_, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "expected no error adding first feature annotation")


	_, err = repo.AddFeatureAnnotation(feat)
	asrt.Error(err, "expected error when adding duplicate feature annotation")
	asrt.Contains(
		err.Error(),
		"unique constraint violated",
		"expected duplicate error message",
	)
}

func TestRemoveFeatureAnnotation(t *testing.T) {
	t.Parallel()
	for _, testCase := range getRemoveTestCases() {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			var identifier string
			assert, repo := setUpFeatureTest(t)
			t.Cleanup(func() { _ = repo.Dbh().Drop() })
			if testCase.wantErr {
				identifier = "non_existent_id"
			} else {
				doc, err := repo.AddFeatureAnnotation(getFullFeatureDoc())
				assert.NoError(err, "expected no error adding test feature annotation")
				identifier = doc.AnnoId
			}
			err := repo.RemoveFeatureAnnotation(identifier, testCase.purge)
			if testCase.wantErr {
				assert.Error(
					err,
					"expected error removing non-existent feature annotation",
				)

				return
			}
			assert.NoError(err, "expected no error removing feature annotation")
			_, err = repo.GetFeatureAnnotation(identifier)
			assert.Error(
				err,
				"expected error getting removed feature annotation",
			)
			assert.True(
				repository.IsAnnotationNotFound(err),
				"should be annotation not found error",
			)
		})
	}
}

func TestUpdateExistingFeatureAnnotation(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))
	added, err := repo.AddFeatureAnnotation(getCompleteFeatureDoc())
	asrt.NoError(err, "expected no error adding initial feature annotation")

	update := &feature.FeatureAnnotationUpdate{
		Id:        added.AnnoId,
		UpdatedBy: "updater@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "updated name",
			Synonyms: []string{"new_syn1", "new_syn2"},
		},
	}

	doc, err := repo.EditFeatureAnnotation(update)
	asrt.NoError(err)
	asrt.Equal(update.UpdatedBy, doc.UpdatedBy)
	asrt.Equal("updated name", doc.Name)

	expectedSynonyms := slices.Concat(
		added.Synonyms,
		update.Attributes.Synonyms,
	)
	slices.Sort(expectedSynonyms)
	slices.Sort(doc.Synonyms)
	asrt.ElementsMatch(
		expectedSynonyms,
		doc.Synonyms,
		"should have combined synonyms",
	)
}

func TestUpdateNonExistentFeatureAnnotation(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	update := &feature.FeatureAnnotationUpdate{
		Id:        "non_existent_id",
		UpdatedBy: "updater@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "will not update",
		},
	}

	_, err := repo.EditFeatureAnnotation(update)
	asrt.Error(err)
	asrt.True(repository.IsAnnotationNotFound(err))
}

func TestAddPropertiesToExistingFeature(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))
	added, err := repo.AddFeatureAnnotation(getCompleteFeatureDoc())
	asrt.NoError(err, "expected no error adding initial feature annotation")

	update := &feature.FeatureAnnotationUpdate{
		Id:        added.AnnoId,
		UpdatedBy: "updater@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Properties: []*feature.TagProperty{
				{
					Tag:       "description",
					Value:     "updated description",
					CreatedBy: "creator3@email.com",
					UpdatedBy: "updater@email.com",
					CreatedAt: timestamppb.New(time.Now()),
					UpdatedAt: timestamppb.New(time.Now()),
				},
			},
		},
	}

	doc, err := repo.EditFeatureAnnotation(update)
	asrt.NoError(err)
	asrt.Len(doc.Properties, 2)


	expectedProperties := slices.Concat(
		added.Properties,
		collection.Map(update.Attributes.Properties, convertProperty),
	)
	slices.SortFunc(expectedProperties, sortTagProperties)
	slices.SortFunc(doc.Properties, sortTagProperties)
	asrt.ElementsMatch(
		expectedProperties,
		doc.Properties,
		"should have combined properties",
	)
}

func TestUpdatePublications_AppendDOI(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	initialDoc := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(initialDoc)
	asrt.NoError(err, "expected no error adding initial feature annotation")
	asrt.NotEmpty(added.Publications, "Initial document should have DOIs")

	newDOIs := []string{"doi:10.1000/new1", "doi:10.1000/new2"}
	update := &feature.FeatureAnnotationUpdate{
		Id:        added.AnnoId,
		UpdatedBy: "doi_updater@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Publications: newDOIs,
		},
	}

	doc, err := repo.EditFeatureAnnotation(update)
	asrt.NoError(err)


	expectedDOIs := slices.Concat(added.Publications, newDOIs)
	slices.Sort(expectedDOIs)
	slices.Sort(doc.Publications)
	asrt.ElementsMatch(
		expectedDOIs,
		doc.Publications,
		"Publications (DOIs) should contain both initial and newly added DOIs",
	)
	asrt.Equal(update.UpdatedBy, doc.UpdatedBy, "UpdatedBy should be updated")
}

func TestUpdatePublications_AppendPubmed(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	initialDoc := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(initialDoc)
	asrt.NoError(err, "expected no error adding initial feature annotation")
	asrt.NotEmpty(added.Pubmed, "Initial document should have Pubmed IDs")

	newPubmedIDs := []string{"76543", "4839439"}
	update := &feature.FeatureAnnotationUpdate{
		Id:        added.AnnoId,
		UpdatedBy: "pubmed_updater@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Pubmed: newPubmedIDs,
		},
	}

	doc, err := repo.EditFeatureAnnotation(update)
	asrt.NoError(err)


	expectedPubmedIDs := slices.Concat(added.Pubmed, newPubmedIDs)
	slices.Sort(expectedPubmedIDs)
	slices.Sort(doc.Pubmed)
	asrt.ElementsMatch(
		expectedPubmedIDs,
		doc.Pubmed,
		"Pubmed IDs should contain both initial and newly added IDs",
	)
	asrt.Equal(update.UpdatedBy, doc.UpdatedBy, "UpdatedBy should be updated")

	expectedDOIs := added.Publications
	slices.Sort(expectedDOIs)
	slices.Sort(doc.Publications)
	asrt.ElementsMatch(
		expectedDOIs,
		doc.Publications,
		"DOIs should remain unchanged",
	)
}

func TestUpdatePublications_AddInitial(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	initialDoc := getFullFeatureDoc()
	added, err := repo.AddFeatureAnnotation(initialDoc)
	asrt.NoError(err, "expected no error adding initial feature annotation")
	asrt.Empty(added.Publications, "Initial document should have no DOIs")
	asrt.Empty(added.Pubmed, "Initial document should have no Pubmed IDs")

	newDOIs := []string{"doi:10.1000/new1", "doi:10.1000/new2"}
	newPubmedIDs := []string{"2039439", "934833"}
	update := &feature.FeatureAnnotationUpdate{
		Id:        added.AnnoId,
		UpdatedBy: "initial_pub_adder@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Publications: newDOIs,
			Pubmed:       newPubmedIDs,
		},
	}

	doc, err := repo.EditFeatureAnnotation(update)
	asrt.NoError(err)
	asrt.ElementsMatch(
		collection.Sorted(newDOIs),
		collection.Sorted(doc.Publications),
		"DOIs should be added",
	)
	asrt.ElementsMatch(
		collection.Sorted(newPubmedIDs),
		collection.Sorted(doc.Pubmed),
		"Pubmed IDs should be added",
	)
	asrt.Equal(update.UpdatedBy, doc.UpdatedBy, "UpdatedBy should be updated")
}



func TestUpdatePublications_Simultaneous(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	initialDoc := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(initialDoc)
	asrt.NoError(err, "expected no error adding initial feature annotation")
	asrt.NotEmpty(added.Publications, "Initial document should have DOIs")
	asrt.NotEmpty(added.Pubmed, "Initial document should have Pubmed IDs")

	newDOIs := []string{"doi:10.1000/new1", "doi:10.1000/new2"}
	newPubmedIDs := []string{"2039439", "934833"}
	update := &feature.FeatureAnnotationUpdate{
		Id:        added.AnnoId,
		UpdatedBy: "simul_updater@email.com",
		Attributes: &feature.FeatureAnnotationAttributes{
			Publications: newDOIs,
			Pubmed:       newPubmedIDs,
		},
	}

	doc, err := repo.EditFeatureAnnotation(update)
	asrt.NoError(err)


	expectedDOIs := slices.Concat(added.Publications, newDOIs)
	expectedPubmedIDs := slices.Concat(added.Pubmed, newPubmedIDs)


	slices.Sort(expectedDOIs)
	slices.Sort(doc.Publications)
	slices.Sort(expectedPubmedIDs)
	slices.Sort(doc.Pubmed)

	asrt.ElementsMatch(
		expectedDOIs,
		doc.Publications,
		"DOIs should contain both initial and newly added IDs",
	)
	asrt.ElementsMatch(
		expectedPubmedIDs,
		doc.Pubmed,
		"Pubmed IDs should contain both initial and newly added IDs",
	)
	asrt.Equal(update.UpdatedBy, doc.UpdatedBy, "UpdatedBy should be updated")
}

func TestUpdateFeatureAnnotation_InvalidInput(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	added, err := repo.AddFeatureAnnotation(getBaseFeatureDoc())
	asrt.NoError(err, "expected no error adding base feature annotation")


	update := &feature.FeatureAnnotationUpdate{
		Id: added.AnnoId,

		Attributes: &feature.FeatureAnnotationAttributes{
			Name: "updated name",
		},
	}

	_, err = repo.EditFeatureAnnotation(update)




	asrt.Error(err, "Expected an error due to missing UpdatedBy field")



}

func TestAddTagToExistingFeature(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "should successfully add test feature")


	tagReq := &feature.AddTagRequest{
		Id: added.AnnoId,
		Tag: &feature.TagPropertyCreate{
			Tag:       "test_tag",
			Value:     "test_value",
			CreatedBy: "tester@example.org",
		},
	}


	updated, err := repo.AddTag(tagReq)
	asrt.NoError(err, "should successfully add tag")
	asrt.Len(
		updated.Properties,
		len(feat.Attributes.Properties)+1,
		"should have one more tag",
	)


	found := false
	for _, prop := range updated.Properties {
		if prop.Tag == tagReq.Tag.Tag {
			found = true
			asrt.Equal(tagReq.Tag.Value, prop.Value, "should match tag value")
			asrt.Equal(
				tagReq.Tag.CreatedBy,
				prop.CreatedBy,
				"should match created by",
			)
			asrt.False(
				prop.CreatedAt.IsZero(),
				"should have creation timestamp",
			)
			asrt.False(prop.UpdatedAt.IsZero(), "should have update timestamp")
		}
	}
	asrt.True(found, "should find added tag")


	asrt.Equal(added.Name, updated.Name, "name should remain unchanged")
	asrt.Equal(
		added.CreatedBy,
		updated.CreatedBy,
		"created_by should remain unchanged",
	)
}

func TestAddTagToNonExistentFeature(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	_, err := repo.AddTag(&feature.AddTagRequest{
		Id: "DDB_G0000000",
		Tag: &feature.TagPropertyCreate{
			Tag:       "test_tag",
			Value:     "test_value",
			CreatedBy: "tester@example.org",
		},
	})

	asrt.Error(err, "should return error for non-existent feature")
	asrt.True(repository.IsAnnotationNotFound(err), "should be not found error")
}

func TestUpdateExistingTag(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "should create test feature")

	tagReq := &feature.AddTagRequest{
		Id: added.AnnoId,
		Tag: &feature.TagPropertyCreate{
			Tag:       "update_test",
			Value:     "initial",
			CreatedBy: "tester@example.org",
		},
	}
	tagged, err := repo.AddTag(tagReq)
	asrt.NoError(err, "should add initial tag")


	updateReq := &feature.UpdateTagRequest{
		Id: added.AnnoId,
		Tag: &feature.TagPropertyUpdate{
			Tag:       "update_test",
			Value:     "updated",
			UpdatedBy: "updater@example.org",
		},
	}


	updated, err := repo.UpdateTag(updateReq)
	asrt.NoError(err, "should successfully update tag")


	var found bool
	for _, prop := range updated.Properties {
		if prop.Tag == "update_test" {
			found = true
			asrt.Equal("updated", prop.Value, "should update value")
			asrt.Equal(
				"updater@example.org",
				prop.UpdatedBy,
				"should update modifier",
			)
			asrt.Equal(
				"tester@example.org",
				prop.CreatedBy,
				"should preserve creator",
			)
			asrt.False(prop.UpdatedAt.IsZero(), "should set update timestamp")
		}
	}
	asrt.True(found, "should find updated tag")
	asrt.Equal(tagged.Name, updated.Name, "should preserve feature name")
}

func TestUpdateNonExistentTag(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "should create test feature")


	_, err = repo.UpdateTag(&feature.UpdateTagRequest{
		Id: added.AnnoId,
		Tag: &feature.TagPropertyUpdate{
			Tag:       "ghost_tag",
			Value:     "new_value",
			UpdatedBy: "tester@example.org",
		},
	})

	asrt.Error(err, "should return error for missing tag")
}

func TestRemoveTag(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "should create base feature")


	tagReq := &feature.AddTagRequest{
		Id: added.AnnoId,
		Tag: &feature.TagPropertyCreate{
			Tag:       "remove_me",
			Value:     "temp_value",
			CreatedBy: "tester@example.org",
		},
	}
	tagged, err := repo.AddTag(tagReq)
	asrt.NoError(err, "should add test tag")
	asrt.Len(
		tagged.Properties,
		len(feat.Attributes.Properties)+1,
		"should have initial tag",
	)


	err = repo.RemoveTag(&feature.RemoveTagRequest{
		Id:  added.AnnoId,
		Tag: "remove_me",
	})
	asrt.NoError(err, "should successfully remove tag")


	updated, err := repo.GetFeatureAnnotation(added.AnnoId)
	asrt.NoError(err, "should fetch updated document")

	var found bool
	for _, prop := range updated.Properties {
		if prop.Tag == "remove_me" {
			found = true
		}
	}
	asrt.False(found, "removed tag should not exist in properties")
	asrt.Len(
		updated.Properties,
		len(tagged.Properties)-1,
		"should reduce properties count by 1",
	)
	asrt.Equal(added.Name, updated.Name, "should preserve feature name")
	asrt.Equal(added.CreatedBy, updated.CreatedBy, "should preserve created_by")
}

func TestRemoveNonExistentTag(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))


	feat := getCompleteFeatureDoc()
	added, err := repo.AddFeatureAnnotation(feat)
	asrt.NoError(err, "should create test feature")


	err = repo.RemoveTag(&feature.RemoveTagRequest{
		Id:  added.AnnoId,
		Tag: "ghost_tag",
	})
	asrt.Error(err, "should return error for missing tag")
}

func TestListByPublicationId_SuccessPubmed(t *testing.T) {
	t.Parallel()
	testListByPublicationIdSuccess(&testListByPublicationIdSuccessParams{
		t:                    t,
		pubID:                "PMID:12345",
		source:               pubmedSource,
		featureIDFieldPrefix: "DDB_G000000",
		setPubFunc:           func(attrs *feature.FeatureAnnotationAttributes, ids []string) { attrs.Pubmed = ids },
		unrelatedPubID:       "PMID:67890",
		errorMsgSuffix:       "Pubmed ID",
	})
}

func TestListByPublicationId_SuccessDOI(t *testing.T) {
	t.Parallel()
	testListByPublicationIdSuccess(&testListByPublicationIdSuccessParams{
		t:                    t,
		pubID:                "doi:10.1234/journal.1",
		source:               doiSource,
		featureIDFieldPrefix: "DDB_G000001",
		setPubFunc:           func(attrs *feature.FeatureAnnotationAttributes, ids []string) { attrs.Publications = ids },
		unrelatedPubID:       "doi:10.5678/journal.2",
		errorMsgSuffix:       "DOI",
	})
}

func TestListByPublicationId_NotFoundIncorrectID(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	pubID := "PMID:11111"
	source := "pubmed"
	feat1 := getFullFeatureDoc()
	feat1.Id = "DDB_G0000021"
	feat1.Attributes.Pubmed = []string{pubID}
	_, err := repo.AddFeatureAnnotation(feat1)
	asrt.NoError(err, "Failed to add feature 1")

	nonExistentPubID := "PMID:99999"
	_, err = repo.ListByPublicationId(nonExistentPubID, source)
	asrt.Error(err, "Expected an error for non-existent publication ID")
	asrt.True(
		repository.IsPublicationAnnotationNotFound(err),
		"Error should be PublicationAnnotationNotFoundError",
	)
}

func TestListByPublicationId_NotFoundIncorrectSource(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	pubID := "PMID:22222"
	incorrectSource := "doi"
	feat1 := getFullFeatureDoc()
	feat1.Id = "DDB_G0000031"
	feat1.Attributes.Pubmed = []string{"pubID"}
	_, err := repo.AddFeatureAnnotation(feat1)
	asrt.NoError(err, "Failed to add feature 1")

	_, err = repo.ListByPublicationId(pubID, incorrectSource)
	asrt.Error(err, "Expected an error for incorrect source")
	asrt.True(
		repository.IsPublicationAnnotationNotFound(err),
		"Error should be PublicationAnnotationNotFoundError",
	)
}

func TestListByPublicationId_NotFoundObsolete(t *testing.T) {
	t.Parallel()
	asrt, repo := setUpFeatureTest(t)
	t.Cleanup(cleanupDB(repo))

	pubID := "PMID:44444"
	source := "pubmed"

	feat1 := getFullFeatureDoc()
	feat1.Id = "DDB_G0000041"
	feat1.Attributes.Pubmed = []string{pubID}
	feat1.IsObsolete = true
	_, err := repo.AddFeatureAnnotation(feat1)
	asrt.NoError(err, "Failed to add feature 1")


	_, err = repo.ListByPublicationId(pubID, source)
	asrt.Error(
		err,
		"Expected an error as only obsolete annotations are linked",
	)
	asrt.True(
		repository.IsPublicationAnnotationNotFound(err),
		"Error should be PublicationAnnotationNotFoundError",
	)
}
```

## File: internal/repository/arangodb/feature_annotation.go
```go
package arangodb

import (
	"context"
	"fmt"
	"slices"
	"time"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)

const (
	pubmedSource = "pubmed"
	doiSource    = "doi"
)

type featureAnnoRepo struct {
	sess     *manager.Session
	database *manager.Database
	feature  driver.Collection
	pub      driver.Collection
	edge     driver.Collection
	featPub  driver.Graph
}


func NewFeatureAnnoRepo(
	connP *manager.ConnectParams,
	collP *FeatureCollectionParams,
) (repository.FeatureAnnotationRepository, error) {
	if err := validate.Struct(collP); err != nil {
		return nil, fmt.Errorf(
			"invalid feature collection parameters: %w", err,
		)
	}


	finalState := collection.Pipe7(
		&repoInitState{
			connP: connP,
			collP: collP,
		},
		stepCreateSession,
		stepCreateFeatureCollection,
		stepCreatePubCollection,
		stepCreateEdgeCollection,
		stepCreateFeatureIndices,
		stepCreatePubIndices,
		stepCreateGraph,
	)


	if finalState.Err != nil {
		return nil, fmt.Errorf(
			"error during repository initialization: %w",
			finalState.Err,
		)
	}

	return &featureAnnoRepo{
		sess:     finalState.sess,
		database: finalState.dbh,
		feature:  finalState.featureColl,
		pub:      finalState.pubColl,
		edge:     finalState.edgeColl,
		featPub:  finalState.graph,
	}, nil
}


func (fann *featureAnnoRepo) GetFeatureAnnotation(
	fid string,
) (*model.FeatureAnnotationDoc, error) {
	res, err := fann.database.GetRow(
		featureGetByIdQ,
		map[string]interface{}{
			"@collection": fann.feature.Name(),
			"graph":       fann.featPub.Name(),
			"id":          fid,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	if res.IsEmpty() {
		return nil, &repository.AnnoNotFoundError{Id: fid}
	}
	doc := &model.FeatureAnnotationDoc{}
	if err := res.Read(doc); err != nil {
		return nil, fmt.Errorf("error reading document: %w", err)
	}

	return doc, nil
}


func (fann *featureAnnoRepo) GetFeatureAnnotationByName(
	name string,
) (*model.FeatureAnnotationDoc, error) {
	res, err := fann.database.GetRow(
		featAnnoGetByNameQ,
		map[string]interface{}{
			"@collection": fann.feature.Name(),
			"graph":       fann.featPub.Name(),
			"name":        name,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error executing query for name %s: %w",
			name,
			err,
		)
	}
	if res.IsEmpty() {

		return nil, &repository.FeatureNameNotFoundError{Name: name}
	}
	doc := &model.FeatureAnnotationDoc{}
	if err := res.Read(doc); err != nil {
		return nil, fmt.Errorf(
			"error reading document for name %s: %w",
			name,
			err,
		)
	}

	return doc, nil
}


func (fann *featureAnnoRepo) AddFeatureAnnotation(
	doc *feature.NewFeatureAnnotation,
) (*model.FeatureAnnotationDoc, error) {

	faDoc := createFeatureAnnotationDoc(doc)


	txOptions := &manager.TransactionOptions{
		WriteCollections: []string{
			fann.feature.Name(),
			fann.pub.Name(),
			fann.edge.Name(),
		},
	}


	txr, err := fann.database.BeginTransaction(
		context.Background(),
		txOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}


	newDoc, err := fann.storeFeatureAnnotation(txr, faDoc)
	if err != nil {
		if abortErr := txr.Abort(); abortErr != nil {
			return nil, fmt.Errorf(
				"error in aborting transaction after %v: %w",
				err,
				abortErr,
			)
		}

		return nil, err
	}


	if err := fann.handlePublications(txr, doc, newDoc); err != nil {
		if abortErr := txr.Abort(); abortErr != nil {
			return nil, fmt.Errorf(
				"error in aborting transaction after %v: %w",
				err,
				abortErr,
			)
		}

		return nil, err
	}


	if err := txr.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return newDoc, nil
}




func (fann *featureAnnoRepo) storeFeatureAnnotation(
	txr *manager.TransactionHandler,
	faDoc *model.FeatureAnnotationDoc,
) (*model.FeatureAnnotationDoc, error) {
	result, err := txr.DoRun(
		fmt.Sprintf("INSERT @doc INTO %s RETURN NEW", fann.feature.Name()),
		map[string]interface{}{
			"doc": faDoc,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error storing feature annotation: %w", err)
	}


	newDoc := &model.FeatureAnnotationDoc{}
	if err := result.Read(newDoc); err != nil {
		return nil, fmt.Errorf("error reading result: %w", err)
	}

	return newDoc, nil
}


func (fann *featureAnnoRepo) handlePublications(
	txr *manager.TransactionHandler,
	doc *feature.NewFeatureAnnotation,
	newDoc *model.FeatureAnnotationDoc,
) error {

	if !collection.IsEmpty(doc.Attributes.Pubmed) {
		if err := fann.processPublicationType(
			txr,
			newDoc,
			doc.Attributes.Pubmed,
			pubmedSource,
		); err != nil {
			return err
		}
	}


	if !collection.IsEmpty(doc.Attributes.Publications) {
		if err := fann.processPublicationType(
			txr,
			newDoc,
			doc.Attributes.Publications,
			doiSource,
		); err != nil {
			return err
		}
	}

	return nil
}


func (fann *featureAnnoRepo) processPublicationType(
	txr *manager.TransactionHandler,
	newDoc *model.FeatureAnnotationDoc,
	pubIDs []string,
	sourceType string,
) error {

	pubKeys, err := fann.upsertPublicationsTx(txr, pubIDs)
	if err != nil {
		return err
	}


	err = fann.createPublicationEdgesTx(
		txr,
		newDoc.ID.String(),
		pubKeys,
		sourceType,
	)
	if err != nil {
		return err
	}


	if sourceType == pubmedSource {
		newDoc.Pubmed = pubIDs
	} else {
		newDoc.Publications = pubIDs
	}

	return nil
}














func (fann *featureAnnoRepo) EditFeatureAnnotation(
	doc *feature.FeatureAnnotationUpdate,
) (*model.FeatureAnnotationDoc, error) {

	finalState := collection.Pipe8(
		&editState{
			fann: fann,
			doc:  doc,
		},
		stepValidateInput,
		stepFetchOriginalDoc,
		stepBeginTransaction,
		stepUpdateDocFields,
		stepExecuteUpdate,
		stepHandlePublications,
		stepCommitTransaction,
		stepRefreshDocumentState,
	)


	if finalState.Err != nil {
		return nil, finalState.Err
	}

	return finalState.updatedDoc, nil
}


func (fann *featureAnnoRepo) ListFeatureAnnotations() ([]*model.FeatureAnnotationDoc, error) {
	return nil, fmt.Errorf("not implemented")
}



func (fann *featureAnnoRepo) ListByPublicationId(
	publicationId string,
	source string,
) ([]*model.FeatureAnnotationDoc, error) {
	binds := map[string]interface{}{
		"@collection": fann.pub.Name(),
		"graph":       fann.featPub.Name(),
		"id":          publicationId,
		"source":      source,
	}

	resultSet, err := fann.database.SearchRows(featureByPublicationIdQ, binds)
	if err != nil {
		return nil, fmt.Errorf(
			"error querying for feature annotations by publication ID %s and source %s: %w",
			publicationId,
			source,
			err,
		)
	}

	if resultSet.IsEmpty() {
		return nil, &repository.PublicationAnnotationNotFoundError{
			ID: publicationId, Source: source,
		}
	}

	var docs []*model.FeatureAnnotationDoc
	for resultSet.Scan() {
		var doc model.FeatureAnnotationDoc
		if err := resultSet.Read(&doc); err != nil {
			return nil, fmt.Errorf(
				"error reading feature annotation document: %w",
				err,
			)
		}
		docs = append(docs, &doc)
	}

	return docs, nil
}

func (fann *featureAnnoRepo) RemoveFeatureAnnotation(
	fid string,
	purge bool,
) error {
	bindVars := map[string]interface{}{
		"@collection": fann.feature.Name(),
		"id":          fid,
	}


	existRes, err := fann.database.GetRow(featureExistQ, bindVars)
	if err != nil {
		return fmt.Errorf("error checking document existence: %w", err)
	}
	if existRes.IsEmpty() {
		return &repository.AnnoNotFoundError{Id: fid}
	}

	if purge {
		err := fann.database.Do(featurePurgeQ, bindVars)
		if err != nil {
			return fmt.Errorf("error executing purge query: %w", err)
		}

		return nil
	}

	err = fann.database.Do(featureObsoleteQ, bindVars)
	if err != nil {
		return fmt.Errorf("error executing obsolete query: %w", err)
	}

	return nil
}


func (fann *featureAnnoRepo) ClearFeatureAnnotations() error {
	if err := fann.feature.Truncate(context.Background()); err != nil {
		return fmt.Errorf(
			"error clearing feature annotations collection: %w",
			err,
		)
	}

	return nil
}

func (fann *featureAnnoRepo) AddTag(
	req *feature.AddTagRequest,
) (*model.FeatureAnnotationDoc, error) {

	doc, err := fann.GetFeatureAnnotation(req.Id)
	if err != nil {
		return nil, err
	}


	newTags := doc.Properties
	newTags = append(newTags, model.TagPropertyDoc{
		Tag:       req.Tag.Tag,
		Value:     req.Tag.Value,
		CreatedBy: req.Tag.CreatedBy,
		CreatedAt: time.Now(),
		UpdatedBy: req.Tag.CreatedBy,
		UpdatedAt: time.Now(),
	})
	newDoc := &model.FeatureAnnotationDoc{}
	ctx := driver.WithReturnNew(context.Background(), newDoc)
	meta, err := fann.feature.UpdateDocument(
		ctx,
		doc.Key,
		map[string]interface{}{"properties": newTags},
	)
	if err != nil {
		return nil, fmt.Errorf("error adding tag: %w", err)
	}
	newDoc.DocumentMeta = meta

	return newDoc, nil
}

func (fann *featureAnnoRepo) UpdateTag(
	req *feature.UpdateTagRequest,
) (*model.FeatureAnnotationDoc, error) {
	doc, err := fann.GetFeatureAnnotation(req.Id)
	if err != nil {
		return nil, err
	}


	idx := slices.IndexFunc(doc.Properties, func(p model.TagPropertyDoc) bool {
		return p.Tag == req.Tag.Tag
	})
	if idx == -1 {
		return nil, fmt.Errorf("tag %s not found", req.Tag.Tag)
	}


	newProps := make([]model.TagPropertyDoc, len(doc.Properties))
	copy(newProps, doc.Properties)
	newProps[idx] = model.TagPropertyDoc{
		Tag:       req.Tag.Tag,
		Value:     req.Tag.Value,
		CreatedBy: doc.Properties[idx].CreatedBy,
		CreatedAt: doc.Properties[idx].CreatedAt,
		UpdatedBy: req.Tag.UpdatedBy,
		UpdatedAt: time.Now(),
	}


	newDoc := &model.FeatureAnnotationDoc{}
	ctx := driver.WithReturnNew(context.Background(), newDoc)
	meta, err := fann.feature.UpdateDocument(
		ctx,
		doc.Key,
		map[string]interface{}{"properties": newProps},
	)
	if err != nil {
		return nil, fmt.Errorf("error updating tag: %w", err)
	}
	newDoc.DocumentMeta = meta

	return newDoc, nil
}

func (fann *featureAnnoRepo) RemoveTag(
	req *feature.RemoveTagRequest,
) error {
	doc, err := fann.GetFeatureAnnotation(req.Id)
	if err != nil {
		return err
	}

	idx := slices.IndexFunc(doc.Properties, func(p model.TagPropertyDoc) bool {
		return p.Tag == req.Tag
	})
	if idx == -1 {
		return fmt.Errorf("tag %s not found", req.Tag)
	}


	newProps := slices.Delete(doc.Properties, idx, idx+1)

	_, err = fann.feature.UpdateDocument(
		context.Background(),
		doc.Key,
		map[string]interface{}{"properties": newProps},
	)
	if err != nil {
		return fmt.Errorf("error removing tag: %w", err)
	}

	return nil
}


func (fann *featureAnnoRepo) Dbh() *manager.Database {
	return fann.database
}



func (fann *featureAnnoRepo) upsertPublicationsTx(
	txr *manager.TransactionHandler,
	ids []string,
) ([]string, error) {
	result, err := txr.DoRun(
		pubUpsertQ,
		map[string]interface{}{
			"ids":         ids,
			"@collection": fann.pub.Name(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"error upserting publications within transaction: %w",
			err,
		)
	}
	pubKeys := make([]string, 0)
	err = result.Read(&pubKeys)
	if err != nil {
		return nil, fmt.Errorf(
			"error reading publication keys from transaction: %w",
			err,
		)
	}

	return pubKeys, nil
}


func (fann *featureAnnoRepo) createPublicationEdgesTx(
	txr *manager.TransactionHandler,
	featureKey string,
	pubKeys []string,
	source string,
) error {
	err := txr.Do(
		featurePubEdgeQ,
		map[string]interface{}{
			"feature_key":      featureKey,
			"pub_keys":         pubKeys,
			"source":           source,
			"@edge_collection": fann.edge.Name(),
		},
	)
	if err != nil {
		return fmt.Errorf(
			"error creating feature-publication edges within transaction: %w",
			err,
		)
	}

	return nil
}

func createFeatureAnnotationDoc(
	doc *feature.NewFeatureAnnotation,
) *model.FeatureAnnotationDoc {
	faDoc := &model.FeatureAnnotationDoc{
		AnnoId:     doc.Id,
		Name:       doc.Attributes.Name,
		CreatedAt:  doc.CreatedAt.AsTime(),
		UpdatedAt:  doc.CreatedAt.AsTime(),
		CreatedBy:  doc.CreatedBy,
		UpdatedBy:  doc.CreatedBy,
		IsObsolete: doc.IsObsolete,
	}
	if doc.UpdatedAt.IsValid() {
		faDoc.UpdatedAt = doc.UpdatedAt.AsTime()
	}
	if len(doc.UpdatedBy) > 0 {
		faDoc.UpdatedBy = doc.UpdatedBy
	}


	setOptionalFields(doc, faDoc)

	return faDoc
}
```
