This file is a merged representation of a subset of the codebase, containing specifically included files, combined into a single document by Repomix.
The content has been processed where comments have been removed, content has been compressed (code blocks are separated by ⋮---- delimiter).

# File Summary

## Purpose
This file contains a packed representation of a subset of the repository's contents that is considered the most important context.
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
- Content has been compressed - code blocks are separated by ⋮---- delimiter
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
⋮----
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
⋮----
"fmt"
"log"
"net"
"os"
"strconv"
"time"
⋮----
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
⋮----
const (
	errCode  = 2
	waitTime = 2
)
⋮----
type serverParams struct {
	repo repository.TaggedAnnotationRepository
	msg  message.Publisher
}
⋮----
func RunServer(clt *cli.Context) error
⋮----
func getLogger(clt *cli.Context) *logrus.Entry
⋮----
func allParams(
	clt *cli.Context,
) (*manager.ConnectParams, *arangodb.CollectionParams, *ontoarango.CollectionParams)
⋮----
func getGrpcOpt() []aphgrpc.Option
⋮----
func repoAndNatsConn(clt *cli.Context) (*serverParams, error)
```

## File: internal/app/service/delete_service.go
```go
package service
⋮----
import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/repository"
	empty "google.golang.org/protobuf/types/known/emptypb"
)
⋮----
"context"
⋮----
"github.com/bufbuild/protovalidate-go"
"github.com/dictyBase/aphgrpc"
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/repository"
empty "google.golang.org/protobuf/types/known/emptypb"
⋮----
func (s *AnnotationService) DeleteAnnotationGroup(ctx context.Context, r *annotation.GroupEntryId) (*empty.Empty, error)
⋮----
func (s *AnnotationService) DeleteAnnotation(ctx context.Context, r *annotation.DeleteAnnotationRequest) (*empty.Empty, error)
```

## File: internal/app/service/read_service.go
```go
package service
⋮----
import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)
⋮----
"context"
⋮----
"github.com/bufbuild/protovalidate-go"
"github.com/dictyBase/aphgrpc"
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
⋮----
var LIMIT int64 = 10
⋮----
func (srv *AnnotationService) GetAnnotation(
	ctx context.Context,
	req *annotation.AnnotationId,
) (*annotation.TaggedAnnotation, error)
⋮----
func (srv *AnnotationService) GetEntryAnnotation(
	ctx context.Context, rea *annotation.EntryAnnotationRequest,
) (*annotation.TaggedAnnotation, error)
⋮----
func (srv *AnnotationService) GetAnnotationGroup(
	ctx context.Context, rid *annotation.GroupEntryId,
) (*annotation.TaggedAnnotationGroup, error)
⋮----
func (srv *AnnotationService) ListAnnotationGroups(
	ctx context.Context, rgp *annotation.ListGroupParameters,
) (*annotation.TaggedAnnotationGroupCollection, error)
⋮----
var gdata []*annotation.TaggedAnnotationGroup_Data
⋮----
func (srv *AnnotationService) ListAnnotations(
	ctx context.Context, ral *annotation.ListParameters,
) (*annotation.TaggedAnnotationCollection, error)
⋮----
func (srv *AnnotationService) GetAnnotationTag(
	ctx context.Context, rta *annotation.TagRequest,
) (*annotation.AnnotationTag, error)
⋮----
func (srv *AnnotationService) modelToCollectionData(
	m *model.AnnoDoc,
) *annotation.TaggedAnnotationCollection_Data
⋮----
func (srv *AnnotationService) getGroup(
	mga *model.AnnoGroup,
) *annotation.TaggedAnnotationGroup
⋮----
func (srv *AnnotationService) getAnnoGroupData(
	m *model.AnnoDoc,
) *annotation.TaggedAnnotationGroup_Data
⋮----
func (srv *AnnotationService) getAnnoData(
	m *model.AnnoDoc,
) *annotation.TaggedAnnotation_Data
⋮----
func (srv *AnnotationService) getGroupData(
	mga *model.AnnoGroup,
) []*annotation.TaggedAnnotationGroup_Data
```

## File: internal/app/service/service.go
```go
package service
⋮----
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
⋮----
"context"
"fmt"
"io"
"time"
⋮----
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
⋮----
const dividerVal = 1000000
⋮----
type oboStreamHandler struct {
	writer *io.PipeWriter
	stream annotation.TaggedAnnotationService_OboJSONFileUploadServer
}
⋮----
func (oh *oboStreamHandler) Write() error
⋮----
type AnnotationService struct {
	*aphgrpc.Service
	repo      repository.TaggedAnnotationRepository
	publisher message.Publisher
	group     string
	annotation.UnimplementedTaggedAnnotationServiceServer
}
⋮----
type Params struct {
	Repository repository.TaggedAnnotationRepository `validate:"required"`
	Publisher  message.Publisher                     `validate:"required"`
	Options    []aphgrpc.Option                      `validate:"required"`
	Group      string                                `validate:"required"`
}
⋮----
func defaultOptions() *aphgrpc.ServiceOptions
⋮----
func NewAnnotationService(srvP *Params) (*AnnotationService, error)
⋮----
func (s *AnnotationService) GetGroupResourceName() string
⋮----
func (s *AnnotationService) OboJSONFileUpload(
	stream annotation.TaggedAnnotationService_OboJSONFileUploadServer,
) error
⋮----
func uploadResponse(
	info *storage.UploadInformation,
) upload.FileUploadResponse_Status
⋮----
func genNextCursorVal(t time.Time) int64
⋮----
func getAnnoAttributes(
	annom *model.AnnoDoc,
) *annotation.TaggedAnnotationAttributes
⋮----
func filterStrToQuery(filter string) (string, error)
⋮----
var empty string
```

## File: internal/app/service/write_service.go
```go
package service
⋮----
import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/dictyBase/aphgrpc"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/repository"
)
⋮----
"context"
⋮----
"github.com/bufbuild/protovalidate-go"
"github.com/dictyBase/aphgrpc"
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/repository"
⋮----
func (s *AnnotationService) UpdateAnnotation(
	ctx context.Context,
	rta *annotation.TaggedAnnotationUpdate,
) (*annotation.TaggedAnnotation, error)
⋮----
func (s *AnnotationService) CreateAnnotation(
	ctx context.Context,
	rta *annotation.NewTaggedAnnotation,
) (*annotation.TaggedAnnotation, error)
⋮----
func (s *AnnotationService) AddToAnnotationGroup(
	ctx context.Context, rta *annotation.AnnotationGroupId,
) (*annotation.TaggedAnnotationGroup, error)
⋮----
func (s *AnnotationService) CreateAnnotationGroup(
	ctx context.Context, rta *annotation.AnnotationIdList,
) (*annotation.TaggedAnnotationGroup, error)
```

## File: internal/app/validate/validate.go
```go
package validate
⋮----
import (
	"fmt"

	"github.com/urfave/cli"
)
⋮----
"fmt"
⋮----
"github.com/urfave/cli"
⋮----
const errNo = 2
⋮----
func ServerArgs(clt *cli.Context) error
```

## File: internal/message/nats/feature_annotation.go
```go
package nats
⋮----
import (
	"encoding/json"
	"fmt"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/message"
	gnats "github.com/nats-io/nats.go"
)
⋮----
"encoding/json"
"fmt"
⋮----
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/message"
gnats "github.com/nats-io/nats.go"
⋮----
type featureAnnotationPublisher struct {
	conn *gnats.Conn
}
⋮----
func NewFeatureAnnotationPublisher(
	host, port string,
	options ...gnats.Option,
) (message.FeatureAnnotationPublisher, error)
⋮----
func (fnp *featureAnnotationPublisher) Publish(
	subj string,
	fann *feature.FeatureAnnotation,
) error
⋮----
func (fnp *featureAnnotationPublisher) Close() error
```

## File: internal/message/nats/nats.go
```go
package nats
⋮----
import (
	"encoding/json"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/message"
	gnats "github.com/nats-io/nats.go"
)
⋮----
"encoding/json"
"fmt"
⋮----
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/message"
gnats "github.com/nats-io/nats.go"
⋮----
type natsPublisher struct {
	conn *gnats.Conn
}
⋮----
func NewPublisher(
	host, port string,
	options ...gnats.Option,
) (message.Publisher, error)
⋮----
func (n *natsPublisher) Publish(
	subj string,
	ann *annotation.TaggedAnnotation,
) error
⋮----
func (n *natsPublisher) Close() error
```

## File: internal/message/message.go
```go
package message
⋮----
import (
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
)
⋮----
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
⋮----
type Publisher interface {

	Publish(subject string, ann *annotation.TaggedAnnotation) error

	Close() error
}
⋮----
type FeatureAnnotationPublisher interface {

	Publish(subject string, ann *feature.FeatureAnnotation) error

	Close() error
}
```

## File: internal/model/feature_annotation.go
```go
package model
⋮----
import (
	"encoding/json"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"
)
⋮----
"encoding/json"
"fmt"
"time"
⋮----
driver "github.com/arangodb/go-driver"
⋮----
type DbLinkDoc struct {
	PrimaryId string `json:"primary_id"`
	Database  string `json:"database"`
	Version   int64  `json:"version"`
	LinkType  string `json:"linktype,omitempty"`
	URL       string `json:"url,omitempty"`
	Label     string `json:"label,omitempty"`
}
⋮----
type TagPropertyDoc struct {
	Tag       string    `json:"tag"`
	Value     string    `json:"value"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
⋮----
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
⋮----
func PubSchema() ([]byte, error)
⋮----
func FeatureAnnotationSchema() ([]byte, error)
⋮----
func getDbLinksSchema() map[string]interface
⋮----
func getPropertiesSchema() map[string]interface
```

## File: internal/model/model.go
```go
package model
⋮----
import (
	"errors"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"
)
⋮----
"errors"
"fmt"
"time"
⋮----
driver "github.com/arangodb/go-driver"
⋮----
type UploadStatus int
⋮----
const (
	Created UploadStatus = iota
	Updated
	Failed
)
⋮----
type AnnoTag struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	IsObsolete bool   `json:"is_obsolete"`
	Ontology   string `json:"ontology"`
}
⋮----
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
⋮----
type AnnoGroup struct {
	AnnoDocs  []*AnnoDoc `json:"annotations"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	GroupId   string     `json:"group_id"`
}
⋮----
type DbGroup struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Group     []string  `json:"group"`
	GroupId   string    `json:"_key,omitempty"`
}
⋮----
func UniqueModel[T comparable](slice []T) []T
⋮----
func DocToIDs(ml []*AnnoDoc) []string
⋮----
func ConvToModel(i interface
```

## File: internal/model/organism.go
```go
package model
⋮----
import (
	"time"

	driver "github.com/arangodb/go-driver"
)
⋮----
"time"
⋮----
driver "github.com/arangodb/go-driver"
⋮----
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
⋮----
func Schema() []byte
```

## File: internal/repository/arangodb/annotation_delete_test.go
```go
package arangodb
⋮----
import (
	"testing"

	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
"testing"
⋮----
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
func TestRemoveFromAnnotationGroup(t *testing.T)
⋮----
func TestRemoveAnnotationGroup(t *testing.T)
⋮----
func TestRemoveAnnotation(t *testing.T)
```

## File: internal/repository/arangodb/annotation_delete.go
```go
package arangodb
⋮----
import (
	"context"
	"errors"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)
⋮----
"context"
"errors"
"fmt"
⋮----
driver "github.com/arangodb/go-driver"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
⋮----
func (ar *arangorepository) RemoveAnnotation(id string, purge bool) error
⋮----
func (ar *arangorepository) RemoveFromAnnotationGroup(
	groupID string,
	idslice ...string,
) (*model.AnnoGroup, error)
```

## File: internal/repository/arangodb/annotation_read_test.go
```go
package arangodb
⋮----
import (
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/stretchr/testify/require"
)
⋮----
"testing"
⋮----
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/stretchr/testify/require"
⋮----
const (
	filterOne   = `entry_id==DDB_G0286429;tag==private note;ontology==dicty_annotation`
	filterTwo   = `entry_id==DDB_G0294491;tag==name description;ontology==dicty_annotation`
	filterThree = `entry_id==jumbo`
)
⋮----
func TestListAnnoFilter(t *testing.T)
⋮----
var mla, ml2, ml4 []*model.AnnoDoc
⋮----
func TestGetAnnotationByID(t *testing.T)
⋮----
func TestGetAnnotationByEntry(t *testing.T)
⋮----
func TestAddAnnotation(t *testing.T)
⋮----
func TestGetAnnotationGroup(t *testing.T)
⋮----
func TestListAnnGrFilter(t *testing.T)
⋮----
func TestListAnnotationGroup(t *testing.T)
⋮----
func TestGetAnnotationTag(t *testing.T)
⋮----
func testListAnnoFilterOneFirstPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) []*model.AnnoDoc
⋮----
func testListAnnoFilterOneSecondPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	prevResult []*model.AnnoDoc,
) []*model.AnnoDoc
⋮----
func testListAnnoFilterOneThirdPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	prevResult []*model.AnnoDoc,
)
⋮----
func testListAnnoFilterTwoFirstPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
) []*model.AnnoDoc
⋮----
func testListAnnoFilterTwoSecondPage(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	prevResult []*model.AnnoDoc,
)
⋮----
func testListAnnoFilterNotFound(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
)
⋮----
func testAddAnnotationSuccess(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	nta *annotation.NewTaggedAnnotation,
)
⋮----
func testAddAnnotationDuplicate(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	nta *annotation.NewTaggedAnnotation,
)
⋮----
func testAddAnnotationNonExistentTag(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
	nta *annotation.NewTaggedAnnotation,
)
⋮----
func testAddAnnotationNonExistentOntology(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
)
⋮----
func testAddAnnotationSuccessSecond(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
)
⋮----
func testAddAnnotationSuccessThird(
	t *testing.T,
	assert *require.Assertions,
	anrepo repository.TaggedAnnotationRepository,
)
```

## File: internal/repository/arangodb/annotation_statement_test_helpers.go
```go
package arangodb
⋮----
import (
	"errors"
	"fmt"
	"testing"

	"github.com/dictyBase/arangomanager/query"
	"github.com/stretchr/testify/require"
)
⋮----
"errors"
"fmt"
"testing"
⋮----
"github.com/dictyBase/arangomanager/query"
"github.com/stretchr/testify/require"
⋮----
func createTestFilterWithSemicolonLogic(field, value string) *query.Filter
⋮----
func createTestFilter(field, value string) *query.Filter
⋮----
func testBothFiltersWithCursor(t *testing.T, filterMap map[string]string)
⋮----
func testBothFiltersWithoutCursor(t *testing.T, filterMap map[string]string)
⋮----
func testBuildAQLStatementFirstFilterWithoutCursor(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testBuildAQLStatementFirstFilterWithCursor(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testBuildAQLStatementSecondFilterWithoutCursor(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testBuildAQLStatementSecondFilterWithCursor(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testBothFiltersStatementTemplate(t *testing.T)
⋮----
func testFirstFilterStatementTemplate(t *testing.T)
⋮----
func testSecondFilterStatementTemplate(t *testing.T)
⋮----
func testInvalidStatementTemplate(t *testing.T)
⋮----
func testValidFilters(t *testing.T, filterMap map[string]string)
⋮----
func testOnlyAnnotationFilters(t *testing.T, filterMap map[string]string)
⋮----
func testOnlyCvtermFilters(t *testing.T, filterMap map[string]string)
⋮----
func testInvalidFilters(t *testing.T, filterMap map[string]string)
⋮----
func testMixedFilters(t *testing.T, filterMap map[string]string)
⋮----
func testExistingError(t *testing.T)
⋮----
func testParseFiltersFuncSuccess(t *testing.T)
⋮----
func testParseFiltersFuncFailureInvalidString(t *testing.T)
⋮----
func testParseFiltersFuncEdgeEmptyString(t *testing.T)
⋮----
// Assuming query.ParseFilterString returns empty slice and no error for empty string
⋮----
func testParseFiltersFuncExistingError(t *testing.T)
⋮----
func testGenFilterStatementSuccessSingle(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testGenFilterStatementSuccessMultiple(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testGenFilterStatementErrorInvalidField(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testGenFilterStatementEdgeEmpty(
	t *testing.T,
	filterMap map[string]string,
)
⋮----
func testGetListAnnoStatementBasicCases(t *testing.T)
⋮----
func testFilterStatement(
	t *testing.T,
	filterString, expectedFilter, filterDescription string,
)
⋮----
func testGetListAnnoStatementValidFilters(t *testing.T)
⋮----
func testGetListAnnoStatementTagFilters(t *testing.T)
```

## File: internal/repository/arangodb/annotation_statement_test.go
```go
package arangodb
⋮----
import (
	"errors"
	"testing"

	"github.com/dictyBase/arangomanager/query"
	"github.com/stretchr/testify/require"
)
⋮----
"errors"
"testing"
⋮----
"github.com/dictyBase/arangomanager/query"
"github.com/stretchr/testify/require"
⋮----
func TestStatementTemplate(t *testing.T)
⋮----
func TestFormatKey(t *testing.T)
⋮----
func TestDetermineStatementType(t *testing.T)
⋮----
func TestFilterAndPartitionFunc(t *testing.T)
⋮----
func TestParseFiltersFunc(t *testing.T)
⋮----
func TestGenFilterStatement(t *testing.T)
⋮----
func TestBuildAQLStatementErrorHandling(t *testing.T)
⋮----
func TestBuildAQLStatementBothFilters(t *testing.T)
⋮----
func TestBuildAQLStatementFirstFilter(t *testing.T)
⋮----
func TestBuildAQLStatementSecondFilter(t *testing.T)
⋮----
func TestBuildAQLStatementFilterGenerationErrors(t *testing.T)
⋮----
func TestGetListAnnoStatement(t *testing.T)
```

## File: internal/repository/arangodb/annotation_statement.go
```go
package arangodb
⋮----
import (
	"errors"
	"fmt"
	"strings"

	"github.com/dictyBase/arangomanager/query"
	"github.com/dictyBase/modware-annotation/internal/collection"
)
⋮----
"errors"
"fmt"
"strings"
⋮----
"github.com/dictyBase/arangomanager/query"
"github.com/dictyBase/modware-annotation/internal/collection"
⋮----
const (

	BothFilters StatementType = "both"

	FirstFilter StatementType = "first"

	SecondFilter StatementType = "second"
)
⋮----
type StatementType string
⋮----
type PickStatementResult struct {
	Statement string
	Err       error
}
⋮----
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
⋮----
var templateMap = map[string]string{
	formatKey(BothFilters, true):   annCvtListFilterWithCursorQ,
	formatKey(BothFilters, false):  annCvtListFilterQ,
	formatKey(FirstFilter, true):   annExclusiveListFilterWithCursorQ,
	formatKey(FirstFilter, false):  annExclusiveListFilterQ,
	formatKey(SecondFilter, true):  cvtExclusiveListFilterWithCursorQ,
	formatKey(SecondFilter, false): cvtExclusiveListFilterQ,
}
⋮----
func formatKey(statementType StatementType, hasCursor bool) string
⋮----
func statementTemplate(ctx FilterContext) (string, bool)
⋮----
func buildAQLStatement(ctx FilterContext) PickStatementResult
⋮----
var result PickStatementResult
⋮----
func genFilterStatement(
	filterMap map[string]string,
	filters []*query.Filter,
	filterType string,
) (string, error)
⋮----
func buildBothFiltersStatement(
	template string,
	filterMap map[string]string,
	firstSet, secondSet []*query.Filter,
) PickStatementResult
⋮----
func buildFirstFilterStatement(
	template string,
	filterMap map[string]string,
	filters []*query.Filter,
) PickStatementResult
⋮----
func buildSecondFilterStatement(
	template string,
	filterMap map[string]string,
	filters []*query.Filter,
) PickStatementResult
⋮----
func getListAnnoStatement(fstr string, cursor int64) PickStatementResult
⋮----
func parseFiltersFunc(ctx FilterContext) FilterContext
⋮----
func filterAndPartitionFunc(ctx FilterContext) FilterContext
⋮----
var validFilters []*query.Filter
var firstSet []*query.Filter
var secondSet []*query.Filter
⋮----
func unsetLogicIfSingleFilter(filters []*query.Filter) []*query.Filter
⋮----
Logic:    "", // Unset logic
⋮----
// determineStatementTypeFunc returns a function for determining statement type in a pipeline.
func determineStatementTypeFunc(ctx FilterContext) FilterContext
⋮----
func determineStatementType(
	first, second []*query.Filter,
) (StatementType, bool)
```

## File: internal/repository/arangodb/annotation_write_test.go
```go
package arangodb
⋮----
import (
	"testing"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
"testing"
⋮----
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
func TestEditAnnotation(t *testing.T)
⋮----
func TestAddAnnotationGroup(t *testing.T)
⋮----
func TestAppendToAnntationGroup(t *testing.T)
```

## File: internal/repository/arangodb/annotation_write.go
```go
package arangodb
⋮----
import (
	"context"
	"errors"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
)
⋮----
"context"
"errors"
"fmt"
⋮----
driver "github.com/arangodb/go-driver"
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
⋮----
const maxTransactionSize = 10000
⋮----
func (ar *arangorepository) AddAnnotation(na *annotation.NewTaggedAnnotation) (*model.AnnoDoc, error)
⋮----
func (ar *arangorepository) EditAnnotation(uat *annotation.TaggedAnnotationUpdate) (*model.AnnoDoc, error)
⋮----
func (ar *arangorepository) AddAnnotationGroup(idslice ...string) (*model.AnnoGroup, error)
⋮----
func (ar *arangorepository) RemoveAnnotationGroup(groupID string) error
⋮----
func (ar *arangorepository) AppendToAnnotationGroup(groupID string, idslice ...string) (*model.AnnoGroup, error)
⋮----
func (ar *arangorepository) createAnno(params *createParams) (*model.AnnoDoc, error)
```

## File: internal/repository/arangodb/arangodb_test.go
```go
package arangodb
⋮----
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
⋮----
"bufio"
"errors"
"fmt"
"math/rand"
"os"
"path/filepath"
"testing"
"time"
⋮----
"github.com/dictyBase/arangomanager/testarango"
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/go-obograph/graph"
ontostorage "github.com/dictyBase/go-obograph/storage"
araobo "github.com/dictyBase/go-obograph/storage/arangodb"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
⋮----
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
⋮----
var ddbg = []string{"DDB_G0286429", "DDB_G0294491"}
⋮----
func toTimestamp(t time.Time) int64
⋮----
func getOntoParams() *araobo.CollectionParams
⋮----
func getCollectionParams() *CollectionParams
⋮----
func loadData(tra *testarango.TestArango) error
⋮----
func saveExistentTestGraph(
	dsr ontostorage.DataSource,
	gra graph.OboGraph,
) error
⋮----
func newTestAnnoWithTagAndOnto(
	onto, tag string,
) *annotation.NewTaggedAnnotation
⋮----
func newTestTaggedAnnotationWithParams(
	tag, entryID string,
) *annotation.NewTaggedAnnotation
⋮----
func newTestTaggedAnnotation() *annotation.NewTaggedAnnotation
⋮----
func newTestTaggedAnnotationsListForFiltering(
	num int,
) []*annotation.NewTaggedAnnotation
⋮----
var nal []*annotation.NewTaggedAnnotation
⋮----
func newTestTaggedAnnotationsList(num int) []*annotation.NewTaggedAnnotation
⋮----
func setUp(
	t *testing.T,
) (*require.Assertions, repository.TaggedAnnotationRepository)
⋮----
func tearDown(repo repository.TaggedAnnotationRepository)
⋮----
func TestLoadOboJSON(t *testing.T)
⋮----
func oboReader() (*os.File, error)
⋮----
func testModelListSort(t *testing.T, m []*model.AnnoDoc)
⋮----
func testGroupMember(
	t *testing.T,
	gla []*model.AnnoGroup,
	count, idx int,
	email string,
)
⋮----
func testModelMaptoID(
	am []*model.AnnoDoc,
	fn func(m *model.AnnoDoc) string,
) []string
⋮----
func model2IdCallback(mod *model.AnnoDoc) string
```

## File: internal/repository/arangodb/arangodb.go
```go
package arangodb
⋮----
import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
	repo "github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/go-playground/validator/v10"
)
⋮----
"context"
"fmt"
⋮----
driver "github.com/arangodb/go-driver"
manager "github.com/dictyBase/arangomanager"
ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
repo "github.com/dictyBase/modware-annotation/internal/repository"
"github.com/go-playground/validator/v10"
⋮----
type annoc struct {
	annot  driver.Collection
	term   driver.Collection
	ver    driver.Collection
	annog  driver.Collection
	verg   driver.Graph
	annotg driver.Graph
}
⋮----
type arangorepository struct {
	sess     *manager.Session
	database *manager.Database
	anno     *annoc
	onto     *ontoarango.OntoCollection
}
⋮----
func NewTaggedAnnotationRepo(
	connP *manager.ConnectParams, collP *CollectionParams, ontoP *ontoarango.CollectionParams,
) (repo.TaggedAnnotationRepository, error)
⋮----
func setAnnotationCollection(dbh *manager.Database, onto *ontoarango.OntoCollection, collP *CollectionParams) (*annoc, error)
⋮----
func setDocumentCollection(dbh *manager.Database, collP *CollectionParams) (*annoc, error)
⋮----
func (ar *arangorepository) Clear() error
⋮----
func (ar *arangorepository) ClearAnnotations() error
⋮----
func DocumentsExists(c driver.Collection, ids ...string) error
⋮----
func (ar *arangorepository) Dbh() *manager.Database
```

## File: internal/repository/arangodb/field.go
```go
package arangodb
⋮----
func FilterMap() map[string]string
```

## File: internal/repository/arangodb/list_filter_statement.go
```go
package arangodb
⋮----
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
```

## File: internal/repository/arangodb/ontology.go
```go
package arangodb
⋮----
import (
	"fmt"
	"io"

	"github.com/dictyBase/go-obograph/storage"
	ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
)
⋮----
"fmt"
"io"
⋮----
"github.com/dictyBase/go-obograph/storage"
ontoarango "github.com/dictyBase/go-obograph/storage/arangodb"
⋮----
func (ar *arangorepository) LoadOboJSON(rde io.Reader) (*storage.UploadInformation, error)
⋮----
func (ar *arangorepository) termID(onto, term string) (string, error)
⋮----
var tid string
⋮----
func (ar *arangorepository) termName(tid string) (string, error)
⋮----
var name string
```

## File: internal/repository/arangodb/organism_test_helpers.go
```go
package arangodb
⋮----
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
⋮----
"testing"
"time"
⋮----
manager "github.com/dictyBase/arangomanager"
"github.com/dictyBase/arangomanager/testarango"
"github.com/dictyBase/go-genproto/dictybaseapis/organism"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/stretchr/testify/require"
"google.golang.org/protobuf/types/known/timestamppb"
⋮----
type validateOrganismParams struct {
	assertions *require.Assertions
	got        *model.OrganismDoc
	key        string
	baseOrg    *organism.NewOrganism
}
⋮----
type updateOrganismParams struct {
	t      *testing.T
	asrt   *require.Assertions
	repo   repository.OrganismRepository
	id     string
	params *organism.OrganismUpdate
}
⋮----
func getFullUpdateParams() *organism.OrganismUpdate
⋮----
func getPartialUpdateParams() *organism.OrganismUpdate
⋮----
func getNotFoundUpdateParams() *organism.OrganismUpdate
⋮----
func updateOrganism(params updateOrganismParams) *model.OrganismDoc
⋮----
func validateFullUpdate(
	asrt *require.Assertions,
	updated *model.OrganismDoc,
	originalKey string,
)
⋮----
func validatePartialUpdate(
	asrt *require.Assertions,
	updated *model.OrganismDoc,
)
⋮----
func validateOrganism(params validateOrganismParams)
⋮----
func setupTestOrganism(
	t *testing.T,
	asrt *require.Assertions,
	repo repository.OrganismRepository,
) *model.OrganismDoc
⋮----
func getOrganismTestCases(baseOrg *organism.NewOrganism) []struct
⋮----
func setUpOrganismTest(
	t *testing.T,
) (*require.Assertions, repository.OrganismRepository)
⋮----
func getTestOrganisms() []*organism.NewOrganism
⋮----
func GetConnectParamsFromDB(tra *testarango.TestArango) *manager.ConnectParams
```

## File: internal/repository/arangodb/organism.go
```go
package arangodb
⋮----
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
⋮----
"context"
"fmt"
"time"
⋮----
driver "github.com/arangodb/go-driver"
manager "github.com/dictyBase/arangomanager"
dorg "github.com/dictyBase/go-genproto/dictybaseapis/organism"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/go-playground/validator/v10"
⋮----
type organismRepo struct {
	sess     *manager.Session
	database *manager.Database
	organism driver.Collection
}
⋮----
func NewOrganismRepo(
	connP *manager.ConnectParams,
	collP *OrganismCollectionParams,
) (repository.OrganismRepository, error)
⋮----
func (org *organismRepo) GetOrganism(id string) (*model.OrganismDoc, error)
⋮----
func (org *organismRepo) GetOrganismByName(
	genus, species string,
) (*model.OrganismDoc, error)
⋮----
func (org *organismRepo) AddOrganism(
	doc *dorg.NewOrganism,
) (*model.OrganismDoc, error)
⋮----
func (org *organismRepo) EditOrganism(
	doc *dorg.OrganismUpdate,
) (*model.OrganismDoc, error)
⋮----
func (org *organismRepo) RemoveOrganism(oid string) error
⋮----
func (org *organismRepo) ListOrganisms() ([]*model.OrganismDoc, error)
⋮----
var organisms []*model.OrganismDoc
⋮----
func (org *organismRepo) ClearOrganisms() error
⋮----
func (org *organismRepo) Dbh() *manager.Database
```

## File: internal/repository/arangodb/pairwise.go
```go
package arangodb
⋮----
import (
	"errors"

	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
"errors"
⋮----
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
type StringPairWiseIterator struct {
	slice []string

	firstIdx int

	secondIdx int

	lastIdx int

	firstPair bool
}
⋮----
func NewStringPairWiseIterator(mdl []string) (StringPairWiseIterator, error)
⋮----
func (p *StringPairWiseIterator) NextStringPair() bool
⋮----
func (p *StringPairWiseIterator) StringPair() (string, string)
⋮----
type ModelAnnoDocPairWiseIterator struct {
	slice []*model.AnnoDoc

	firstIdx int

	secondIdx int

	lastIdx int

	firstPair bool
}
⋮----
func NewModelAnnoDocPairWiseIterator(mdl []*model.AnnoDoc) (ModelAnnoDocPairWiseIterator, error)
⋮----
func (p *ModelAnnoDocPairWiseIterator) NextModelAnnoDocPair() bool
⋮----
func (p *ModelAnnoDocPairWiseIterator) ModelAnnoDocPair() (*model.AnnoDoc, *model.AnnoDoc)
```

## File: internal/repository/arangodb/parameters.go
```go
package arangodb
⋮----
import "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
⋮----
type createParams struct {
	attr *annotation.NewTaggedAnnotationAttributes
	id   string
	tag  string
}
⋮----
type OrganismCollectionParams struct {

	Organism string `validate:"required"`
}
⋮----
type CollectionParams struct {

	Annotation string `validate:"required"`

	AnnoGroup string `validate:"required"`


	AnnoTerm string `validate:"required"`


	AnnoVersion string `validate:"required"`


	AnnoTagGraph string `validate:"required"`


	AnnoVerGraph string `validate:"required"`


	AnnoIndexes []string `validate:"required"`
}
⋮----
type FeatureCollectionParams struct {

	Feature string `validate:"required"`
	Pub     string `validate:"required"`
	Edge    string `validate:"required"`

	Graph string `validate:"required"`
}
```

## File: internal/repository/organism.go
```go
package repository
⋮----
import (
	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-genproto/dictybaseapis/organism"
	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
manager "github.com/dictyBase/arangomanager"
"github.com/dictyBase/go-genproto/dictybaseapis/organism"
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
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
⋮----
import (
	"io"

	manager "github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/go-obograph/storage"
	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
"io"
⋮----
manager "github.com/dictyBase/arangomanager"
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/go-obograph/storage"
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
type ListAnnotationsParams struct {
	Cursor int64
	Limit  int64  `validate:"required"`
	Filter string `validate:"required"`
}
⋮----
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
⋮----
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
⋮----
"log"
"os"
⋮----
apiflag "github.com/dictyBase/aphgrpc"
arangoflag "github.com/dictyBase/arangomanager/command/flag"
oboaction "github.com/dictyBase/go-obograph/command/action"
oboflag "github.com/dictyBase/go-obograph/command/flag"
obovalidate "github.com/dictyBase/go-obograph/command/validate"
"github.com/dictyBase/modware-annotation/internal/app/server"
"github.com/dictyBase/modware-annotation/internal/app/validate"
"github.com/urfave/cli"
⋮----
func main()
⋮----
func getServerFlags() []cli.Flag
⋮----
func getFeatureServerFlags() []cli.Flag
⋮----
func ontoCollFlags() []cli.Flag
⋮----
func annoCollFlags() []cli.Flag
```

## File: internal/app/server/server_feature.go
```go
package server
⋮----
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
⋮----
"fmt"
"log"
"net"
"strconv"
"time"
⋮----
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
⋮----
type featureServerParams struct {
	repo repository.FeatureAnnotationRepository
	msg  message.FeatureAnnotationPublisher
}
⋮----
func RunFeatureServer(clt *cli.Context) error
⋮----
func allFeatureParams(
	clt *cli.Context,
) (*manager.ConnectParams, *arangodb.FeatureCollectionParams)
⋮----
func getFeatureGrpcOpt() []aphgrpc.Option
⋮----
func featureRepoAndNatsConn(clt *cli.Context) (*featureServerParams, error)
```

## File: internal/collection/collection_test.go
```go
package collection
⋮----
import (
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)
⋮----
"slices"
"strconv"
"testing"
⋮----
"github.com/stretchr/testify/require"
⋮----
func TestMapSeq(t *testing.T)
⋮----
func TestPartition(t *testing.T)
⋮----
func TestFind(t *testing.T)
⋮----
func testFindIntegerFound(t *testing.T)
⋮----
func testFindIntegerNotFound(t *testing.T)
⋮----
func testFindStringMultipleMatches(t *testing.T)
⋮----
func testFindStruct(t *testing.T)
⋮----
type person struct {
		name string
		age  int
	}
⋮----
func testFindEmptySlice(t *testing.T)
⋮----
func testFindNilSlice(t *testing.T)
⋮----
var slice []int
```

## File: internal/repository/arangodb/annotation_read.go
```go
package arangodb
⋮----
import (
	"context"
	"errors"
	"fmt"

	"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/go-playground/validator/v10"
)
⋮----
"context"
"errors"
"fmt"
⋮----
"github.com/dictyBase/go-genproto/dictybaseapis/annotation"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/go-playground/validator/v10"
⋮----
var validate = validator.New(validator.WithRequiredStructEnabled())
⋮----
func (ar *arangorepository) GetAnnotationByID(
	annoid string,
) (*model.AnnoDoc, error)
⋮----
func (ar *arangorepository) GetAnnotationByEntry(
	req *annotation.EntryAnnotationRequest,
) (*model.AnnoDoc, error)
⋮----
func (ar *arangorepository) ListAnnotations(
	params *repository.ListAnnotationsParams,
) ([]*model.AnnoDoc, error)
⋮----
func (ar *arangorepository) GetAnnotationGroup(
	groupID string,
) (*model.AnnoGroup, error)
⋮----
func (ar *arangorepository) ListAnnotationGroup(
	cursor, limit int64,
	fstr string,
) ([]*model.AnnoGroup, error)
⋮----
var agrp []*model.AnnoGroup
var stmt string
⋮----
func (ar *arangorepository) GetAnnotationTag(
	tag, ontology string,
) (*model.AnnoTag, error)
⋮----
func (ar *arangorepository) existAnno(
	attr *annotation.NewTaggedAnnotationAttributes,
	tag string,
) error
⋮----
func (ar *arangorepository) groupID2Annotations(
	groupID string,
) ([]*model.AnnoDoc, error)
⋮----
var annoModel []*model.AnnoDoc
⋮----
func (ar *arangorepository) getAllAnnotations(
	ids ...string,
) ([]*model.AnnoDoc, error)
```

## File: internal/repository/arangodb/feature_annotation_helpers.go
```go
package arangodb
⋮----
import (
	"fmt"

	driver "github.com/arangodb/go-driver"
	manager "github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
"fmt"
⋮----
driver "github.com/arangodb/go-driver"
manager "github.com/dictyBase/arangomanager"
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
type CreateIndexArgs struct {
	Dbh          *manager.Database
	Coll         driver.Collection
	Fields       []string
	UniqueFields []string
	ErrPrefix    string
}
⋮----
func createSession(
	connP *manager.ConnectParams,
) (*manager.Session, *manager.Database, error)
⋮----
func createFeatureCollection(
	dbh *manager.Database,
	collP *FeatureCollectionParams,
) (driver.Collection, error)
⋮----
func createIndices(args *CreateIndexArgs) error
⋮----
func createFeatureIndices(dbh *manager.Database, coll driver.Collection) error
⋮----
func createPubIndices(dbh *manager.Database, coll driver.Collection) error
⋮----
func createPubCollection(
	dbh *manager.Database,
	collP *FeatureCollectionParams,
) (driver.Collection, error)
⋮----
func createEdgeCollection(
	dbh *manager.Database,
	collP *FeatureCollectionParams,
) (driver.Collection, error)
⋮----
func createFeaturePubGraph(
	dbh *manager.Database,
	graphName string,
	featureColl driver.Collection,
	pubColl driver.Collection,
	edgeColl driver.Collection,
) (driver.Graph, error)
⋮----
func updateBasicFields(
	faDoc *model.FeatureAnnotationDoc,
	doc *feature.FeatureAnnotationUpdate,
)
⋮----
func updateAttributes(
	mdoc *model.FeatureAnnotationDoc,
	attrs *feature.FeatureAnnotationAttributes,
)
⋮----
func convertDbLink(link *feature.DbLink) model.DbLinkDoc
⋮----
func convertProperty(prop *feature.TagProperty) model.TagPropertyDoc
⋮----
func setOptionalFields(
	doc *feature.NewFeatureAnnotation,
	faDoc *model.FeatureAnnotationDoc,
) *model.FeatureAnnotationDoc
```

## File: internal/repository/arangodb/feature_annotation_pipeline.go
```go
package arangodb
⋮----
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
⋮----
"context"
"fmt"
⋮----
driver "github.com/arangodb/go-driver"
manager "github.com/dictyBase/arangomanager"
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
⋮----
type editState struct {
	fann        *featureAnnoRepo
	doc         *feature.FeatureAnnotationUpdate
	txr         *manager.TransactionHandler
	origDoc     *model.FeatureAnnotationDoc
	updatedDoc  *model.FeatureAnnotationDoc
	updateQuery string
	Err         error
}
⋮----
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
⋮----
func stepCreateSession(state *repoInitState) *repoInitState
⋮----
func stepCreateFeatureCollection(state *repoInitState) *repoInitState
⋮----
func stepCreatePubCollection(state *repoInitState) *repoInitState
⋮----
func stepCreateEdgeCollection(state *repoInitState) *repoInitState
⋮----
func stepCreateFeatureIndices(state *repoInitState) *repoInitState
⋮----
func stepCreatePubIndices(state *repoInitState) *repoInitState
⋮----
func stepCreateGraph(state *repoInitState) *repoInitState
⋮----
func stepFetchOriginalDoc(state *editState) *editState
⋮----
var err error
⋮----
func stepBeginTransaction(state *editState) *editState
⋮----
func stepUpdateDocFields(state *editState) *editState
⋮----
func stepExecuteUpdate(state *editState) *editState
⋮----
func stepHandlePublications(state *editState) *editState
⋮----
func stepCommitTransaction(state *editState) *editState
⋮----
type featureAnnotationUpdateValidator struct {
	ID        string `validate:"required"       json:"id"`
	UpdatedBy string `validate:"required,email" json:"updated_by"`
}
⋮----
func stepValidateInput(state *editState) *editState
⋮----
func stepRefreshDocumentState(state *editState) *editState
```

## File: internal/repository/arangodb/organism_test.go
```go
package arangodb
⋮----
import (
	"fmt"
	"testing"
	"time"

	"github.com/dictyBase/go-genproto/dictybaseapis/organism"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)
⋮----
"fmt"
"testing"
"time"
⋮----
"github.com/dictyBase/go-genproto/dictybaseapis/organism"
"github.com/dictyBase/modware-annotation/internal/repository"
"google.golang.org/protobuf/types/known/timestamppb"
⋮----
func TestAddOrganism(t *testing.T)
⋮----
func TestAddDuplicateOrganism(t *testing.T)
⋮----
func TestGetOrganism(t *testing.T)
⋮----
func TestGetOrganismByName(t *testing.T)
⋮----
func TestEditOrganism(t *testing.T)
⋮----
func TestRemoveOrganism(t *testing.T)
⋮----
func TestListOrganisms(t *testing.T)
⋮----
func TestClearOrganisms(t *testing.T)
```

## File: internal/repository/arangodb/statement.go
```go
package arangodb
⋮----
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
```

## File: internal/app/service/feature_annotation.go
```go
package service
⋮----
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
⋮----
"context"
"fmt"
"strings"
⋮----
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
⋮----
type FeatureAnnotationService struct {
	*aphgrpc.Service
	repo      repository.FeatureAnnotationRepository
	publisher message.FeatureAnnotationPublisher
	feature.UnimplementedFeatureAnnotationServiceServer
}
⋮----
type FeatureParams struct {
	Repository repository.FeatureAnnotationRepository `validate:"required"`
	Publisher  message.FeatureAnnotationPublisher     `validate:"required"`
	Options    []aphgrpc.Option
}
⋮----
func featureAnnoDefaultOptions() *aphgrpc.ServiceOptions
⋮----
func NewFeatureAnnotationService(
	params *FeatureParams,
) (*FeatureAnnotationService, error)
⋮----
func (srv *FeatureAnnotationService) GetFeatureAnnotation(
	ctx context.Context,
	req *feature.FeatureAnnotationId,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) GetFeatureAnnotationByName(
	ctx context.Context,
	req *feature.FeatureName,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) CreateFeatureAnnotation(
	ctx context.Context,
	req *feature.NewFeatureAnnotation,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) UpdateFeatureAnnotation(
	ctx context.Context,
	req *feature.FeatureAnnotationUpdate,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) DeleteFeatureAnnotation(
	ctx context.Context,
	req *feature.DeleteFeatureAnnotationRequest,
) (*emptypb.Empty, error)
⋮----
func (srv *FeatureAnnotationService) AddTag(
	ctx context.Context,
	req *feature.AddTagRequest,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) UpdateTag(
	ctx context.Context,
	req *feature.UpdateTagRequest,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) RemoveTag(
	ctx context.Context,
	req *feature.RemoveTagRequest,
) (*feature.FeatureAnnotation, error)
⋮----
func (srv *FeatureAnnotationService) ListFeatureAnnotationsByPubmedId(
	ctx context.Context,
	req *feature.PubmedId,
) (*feature.FeatureAnnotationCollection, error)
⋮----
func (srv *FeatureAnnotationService) ListFeatureAnnotationsByDOI(
	ctx context.Context,
	req *feature.DOI,
) (*feature.FeatureAnnotationCollection, error)
⋮----
func convertToProto(
	feat *model.FeatureAnnotationDoc,
) *feature.FeatureAnnotation
⋮----
func convertDbLink(link model.DbLinkDoc) *feature.DbLink
⋮----
func convertProperty(prop model.TagPropertyDoc) *feature.TagProperty
```

## File: internal/collection/collection.go
```go
package collection
⋮----
import (
	"cmp"
	"iter"
	"slices"
)
⋮----
"cmp"
"iter"
"slices"
⋮----
func Find[T any](slice []T, predicate func(T) bool) (*T, bool)
⋮----
func Map[T1, T2 any](slc []T1, fnc func(T1) T2) []T2
⋮----
func CurriedMap[T1, T2 any](fnc func(T1) T2) func([]T1) []T2
⋮----
func Include[T cmp.Ordered](slice []T, element T) bool
⋮----
func RemoveStringItems(slice []string, items ...string) []string
⋮----
func Filter[T any](slice []T, predicate func(T) bool) []T
⋮----
func CurriedFilter[T any](predicate func(T) bool) func([]T) []T
⋮----
func MapSeq[T1, T2 any](seq iter.Seq[T1], fn func(T1) T2) iter.Seq[T2]
⋮----
func PartitionTuple2[T any](
	slice []T,
	predicate func(T) bool,
) Tuple2[[]T, []T]
⋮----
func CurriedPartitionTuple2[T any](
	predicate func(T) bool,
) func([]T) Tuple2[[]T, []T]
⋮----
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T)
⋮----
func CurriedPartition[T any](predicate func(T) bool) func([]T) ([]T, []T)
⋮----
func Pipe2[T1, T2, T3 any](tup T1, f1 func(T1) T2, fn2 func(T2) T3) T3
⋮----
func Pipe3[T1, T2, T3, T4 any](
	initial T1,
	f1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
) T4
⋮----
func Pipe4[T1, T2, T3, T4, T5 any](
	initial T1,
	f1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
) T5
⋮----
func Pipe5[T1, T2, T3, T4, T5, T6 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
) T6
⋮----
func Pipe6[T1, T2, T3, T4, T5, T6, T7 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
	fn6 func(T6) T7,
) T7
⋮----
func Pipe7[T1, T2, T3, T4, T5, T6, T7, T8 any](
	initial T1,
	fn1 func(T1) T2,
	fn2 func(T2) T3,
	fn3 func(T3) T4,
	fn4 func(T4) T5,
	fn5 func(T5) T6,
	fn6 func(T6) T7,
	fn7 func(T7) T8,
) T8
⋮----
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
) T9
⋮----
type Tuple2[T1, T2 any] struct {
	First  T1
	Second T2
}
⋮----
func NewTuple2[T1, T2 any](first T1, second T2) Tuple2[T1, T2]
⋮----
func SliceToTuple2[T1, T2 any](slice []any) Tuple2[T1, T2]
⋮----
var first T1
var second T2
⋮----
func TFold[A, B, R any](
	tup Tuple2[A, B],
	folder func(Tuple2[A, B]) R,
) R
⋮----
func CurriedTFold[A, B, R any](
	folder func(Tuple2[A, B]) R,
) func(Tuple2[A, B]) R
⋮----
func IsEmpty[T any](slice []T) bool
⋮----
func Sorted[T cmp.Ordered](slice []T) []T
```

## File: internal/repository/feature_annotation.go
```go
package repository
⋮----
import (
	manager "github.com/dictyBase/arangomanager"
	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/model"
)
⋮----
manager "github.com/dictyBase/arangomanager"
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/model"
⋮----
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
⋮----
import (
	"fmt"
)
⋮----
"fmt"
⋮----
type AnnoNotFoundError struct {
	Id string
}
⋮----
func (ae *AnnoNotFoundError) Error() string
⋮----
type GroupNotFoundError struct {
	Id string
}
⋮----
type FeatureNameNotFoundError struct {
	Name string
}
⋮----
func IsAnnotationNotFound(err error) bool
⋮----
func IsFeatureNameNotFound(err error) bool
⋮----
type PublicationAnnotationNotFoundError struct {
	ID     string
	Source string
}
⋮----
func IsPublicationAnnotationNotFound(err error) bool
⋮----
func IsGroupNotFound(err error) bool
⋮----
type AnnoListNotFoundError struct{}
⋮----
func IsAnnotationListNotFound(err error) bool
⋮----
type AnnoGroupListNotFoundError struct{}
⋮----
func IsAnnotationGroupListNotFound(err error) bool
⋮----
type AnnoTagNotFoundError struct {
	Tag string
}
⋮----
func IsAnnoTagNotFound(err error) bool
⋮----
type OrganismNotFoundError struct {
	ID string
}
⋮----
func IsOrganismNotFound(err error) bool
⋮----
type ListNotFoundError struct{}
⋮----
func IsListNotFound(err error) bool
```

## File: internal/app/service/feature_annotation_test.go
```go
package service
⋮----
import (
	"context"
	"testing"
)
⋮----
"context"
"testing"
⋮----
func TestCreateFeatureAnnotation(t *testing.T)
⋮----
func TestGetFeatureAnnotation(t *testing.T)
⋮----
func TestGetFeatureAnnotationByName(t *testing.T)
⋮----
func TestUpdateFeatureAnnotation(t *testing.T)
⋮----
func TestListFeatureAnnotationsByPubmedId(t *testing.T)
⋮----
func TestListFeatureAnnotationsByDOI(t *testing.T)
```

## File: internal/repository/arangodb/feature_annotation_test_helpers.go
```go
package arangodb
⋮----
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
⋮----
"slices"
"strings"
"testing"
"time"
⋮----
"github.com/dictyBase/arangomanager/testarango"
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/stretchr/testify/require"
"google.golang.org/protobuf/types/known/timestamppb"
⋮----
type validateDbLinksParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        []model.DbLinkDoc
	expected   []*feature.DbLink
}
⋮----
type validatePropertiesParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        []model.TagPropertyDoc
	expected   []*feature.TagProperty
}
⋮----
type removeFeatureTestCase struct {
	name    string
	purge   bool
	wantErr bool
}
⋮----
type validateFeatureAnnotationParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        *model.FeatureAnnotationDoc
	base       *feature.NewFeatureAnnotation
}
⋮----
type validateCompleteFeatureParams struct {
	t          *testing.T
	assertions *require.Assertions
	got        *model.FeatureAnnotationDoc
	expected   *feature.NewFeatureAnnotation
}
⋮----
type testListByPublicationIdSuccessParams struct {
	t                    *testing.T
	pubID                string
	source               string
	featureIDFieldPrefix string
	setPubFunc           func(*feature.FeatureAnnotationAttributes, []string)
	unrelatedPubID       string
	errorMsgSuffix       string
}
⋮----
type featFn func() *feature.NewFeatureAnnotation
⋮----
func getBaseFeatureDoc() *feature.NewFeatureAnnotation
⋮----
func getCombinedFeatureDoc(
	baseFn featFn,
	advFn featFn,
) *feature.NewFeatureAnnotation
⋮----
func getFullFeatureDoc() *feature.NewFeatureAnnotation
⋮----
func getCompleteFeatureDoc() *feature.NewFeatureAnnotation
⋮----
func getMultiPropertyTestCase() *feature.NewFeatureAnnotation
⋮----
func getBasicTestCases() []*feature.NewFeatureAnnotation
⋮----
func setUpFeatureTest(
	t *testing.T,
) (*require.Assertions, repository.FeatureAnnotationRepository)
⋮----
func validateProperties(params validatePropertiesParams)
⋮----
func validateDbLinks(params validateDbLinksParams)
⋮----
func validateBasicFields(params validateFeatureAnnotationParams)
⋮----
func validateCompleteFeatureAnnotation(params validateCompleteFeatureParams)
⋮----
func sortTagProperties(a, b model.TagPropertyDoc) int
⋮----
func getRemoveTestCases() []removeFeatureTestCase
⋮----
func assertListByPublicationResults(
	t *testing.T,
	asrt *require.Assertions,
	results []*model.FeatureAnnotationDoc,
	added1 *model.FeatureAnnotationDoc,
	added2 *model.FeatureAnnotationDoc,
)
⋮----
func testListByPublicationIdSuccess(
	params *testListByPublicationIdSuccessParams,
)
⋮----
func cleanupDB(repo repository.FeatureAnnotationRepository) func()
```

## File: internal/repository/arangodb/feature_statement.go
```go
package arangodb
⋮----
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
```

## File: internal/app/service/feature_annotation_test_helpers.go
```go
package service
⋮----
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
⋮----
"context"
"net"
"os"
"slices"
"strings"
"testing"
⋮----
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
⋮----
type assertGrpcErrorParams struct {
	assert               *require.Assertions
	err                  error
	expectedCode         codes.Code
	expectedMsgSubstring string
}
⋮----
type testParams struct {
	t      *testing.T
	ctx    context.Context
	client feature.FeatureAnnotationServiceClient
	assert *require.Assertions
}
⋮----
type MockMessage struct{}
⋮----
func (msn *MockMessage) Publish(
	subject string,
	feat *feature.FeatureAnnotation,
) error
⋮----
func (msn *MockMessage) Close() error
⋮----
func sortTagPropertiesByTag(a, b *feature.TagProperty) int
⋮----
func extractTagAndValue(prop *feature.TagProperty) *feature.TagProperty
⋮----
func setup(
	t *testing.T,
) (feature.FeatureAnnotationServiceClient, *require.Assertions)
⋮----
func newTestFeature() *feature.NewFeatureAnnotation
⋮----
func testCreateValidFeature(params *testParams)
⋮----
func testListByPublicationHelper(
	params *testParams,
	publicationType string,
	publicationID string,
	featureID1 string,
	featureID2 string,
	featureNamePrefix string,
)
⋮----
var resp *feature.FeatureAnnotationCollection
⋮----
func testListByDOIValid(params *testParams)
⋮----
func testListByDOINotFound(params *testParams)
⋮----
func testListByDOIInvalid(params *testParams)
⋮----
req := &feature.DOI{Id: ""} // Invalid (empty) DOI
⋮----
func testListByPubmedIdValid(params *testParams)
⋮----
func testListByPubmedIdNotFound(params *testParams)
⋮----
func testListByPubmedIdInvalid(params *testParams)
⋮----
req := &feature.PubmedId{Id: ""} // Invalid (empty) pubmed ID
⋮----
func testCreateMissingFields(params *testParams)
⋮----
func testCreateDuplicateFeature(params *testParams)
⋮----
func testGetExistingFeature(params *testParams)
⋮----
func testGetNonExistentFeature(params *testParams)
⋮----
func testGetFeatureWithInvalidID(params *testParams)
⋮----
Id: "", // Empty ID
⋮----
func testUpdateExistingFeature(params *testParams)
⋮----
func testUpdateNonExistentFeature(params *testParams)
⋮----
func testUpdateWithInvalidData(params *testParams)
⋮----
func testGetExistingFeatureByName(params *testParams)
⋮----
func testGetNonExistentFeatureByName(params *testParams)
⋮----
func assertGrpcError(params assertGrpcErrorParams)
⋮----
strings.ToLower(sts.Message()), // Case-insensitive check
⋮----
func testGetFeatureWithEmptyName(params *testParams)
⋮----
req := &feature.FeatureName{Name: ""} // Empty name
```

## File: internal/repository/arangodb/feature_annotation_test.go
```go
package arangodb
⋮----
import (
	"slices"
	"testing"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"github.com/dictyBase/modware-annotation/internal/collection"
	"github.com/dictyBase/modware-annotation/internal/model"
	"github.com/dictyBase/modware-annotation/internal/repository"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)
⋮----
"slices"
"testing"
"time"
⋮----
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
"github.com/stretchr/testify/require"
"google.golang.org/protobuf/types/known/timestamppb"
⋮----
func TestGetFeatureAnnotation(t *testing.T)
⋮----
func TestGetFeatureAnnotationByName(t *testing.T)
⋮----
func TestAddFeatureAnnotationBasic(t *testing.T)
⋮----
func TestAddFeatureAnnotationFull(t *testing.T)
⋮----
func TestAddFeatureAnnotationMultiProperty(t *testing.T)
⋮----
func TestAddDuplicateFeatureAnnotation(t *testing.T)
⋮----
func TestRemoveFeatureAnnotation(t *testing.T)
⋮----
var identifier string
⋮----
func TestUpdateExistingFeatureAnnotation(t *testing.T)
⋮----
func TestUpdateNonExistentFeatureAnnotation(t *testing.T)
⋮----
func TestAddPropertiesToExistingFeature(t *testing.T)
⋮----
func TestUpdatePublications_AppendDOI(t *testing.T)
⋮----
func TestUpdatePublications_AppendPubmed(t *testing.T)
⋮----
func TestUpdatePublications_AddInitial(t *testing.T)
⋮----
func TestUpdatePublications_Simultaneous(t *testing.T)
⋮----
func TestUpdateFeatureAnnotation_InvalidInput(t *testing.T)
⋮----
func TestAddTag_WithDefaultTimestamp(t *testing.T)
⋮----
func TestAddTag_WithProvidedTimestamp(t *testing.T)
⋮----
func TestAddTagToNonExistentFeature(t *testing.T)
⋮----
func seedAnnotationWithTags(
	t *testing.T,
	repo repository.FeatureAnnotationRepository,
) *model.FeatureAnnotationDoc
⋮----
func TestUpdateTag_SuccessDefaultTimestamp(t *testing.T)
⋮----
func TestUpdateTag_SuccessExplicitTimestamp(t *testing.T)
⋮----
func TestUpdateTag_FailNonExistentTag(t *testing.T)
⋮----
func TestUpdateTag_FailNonExistentFeature(t *testing.T)
⋮----
var nfErr *repository.AnnoNotFoundError
⋮----
func TestRemoveTag(t *testing.T)
⋮----
var found bool
⋮----
func TestRemoveNonExistentTag(t *testing.T)
⋮----
func TestListByPublicationId_SuccessPubmed(t *testing.T)
⋮----
func TestListByPublicationId_SuccessDOI(t *testing.T)
⋮----
func TestListByPublicationId_NotFoundIncorrectID(t *testing.T)
⋮----
func TestListByPublicationId_NotFoundIncorrectSource(t *testing.T)
⋮----
func TestListByPublicationId_NotFoundObsolete(t *testing.T)
```

## File: internal/repository/arangodb/feature_annotation.go
```go
package arangodb
⋮----
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
⋮----
"context"
"fmt"
"slices"
"time"
⋮----
driver "github.com/arangodb/go-driver"
manager "github.com/dictyBase/arangomanager"
feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
"github.com/dictyBase/modware-annotation/internal/collection"
"github.com/dictyBase/modware-annotation/internal/model"
"github.com/dictyBase/modware-annotation/internal/repository"
⋮----
const (
	pubmedSource = "pubmed"
	doiSource    = "doi"
)
⋮----
type featureAnnoRepo struct {
	sess     *manager.Session
	database *manager.Database
	feature  driver.Collection
	pub      driver.Collection
	edge     driver.Collection
	featPub  driver.Graph
}
⋮----
func NewFeatureAnnoRepo(
	connP *manager.ConnectParams,
	collP *FeatureCollectionParams,
) (repository.FeatureAnnotationRepository, error)
⋮----
func (fann *featureAnnoRepo) GetFeatureAnnotation(
	fid string,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) GetFeatureAnnotationByName(
	name string,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) AddFeatureAnnotation(
	doc *feature.NewFeatureAnnotation,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) storeFeatureAnnotation(
	txr *manager.TransactionHandler,
	faDoc *model.FeatureAnnotationDoc,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) handlePublications(
	txr *manager.TransactionHandler,
	doc *feature.NewFeatureAnnotation,
	newDoc *model.FeatureAnnotationDoc,
) error
⋮----
func (fann *featureAnnoRepo) processPublicationType(
	txr *manager.TransactionHandler,
	newDoc *model.FeatureAnnotationDoc,
	pubIDs []string,
	sourceType string,
) error
⋮----
func (fann *featureAnnoRepo) EditFeatureAnnotation(
	doc *feature.FeatureAnnotationUpdate,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) ListFeatureAnnotations() ([]*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) ListByPublicationId(
	publicationId string,
	source string,
) ([]*model.FeatureAnnotationDoc, error)
⋮----
var docs []*model.FeatureAnnotationDoc
⋮----
var doc model.FeatureAnnotationDoc
⋮----
func (fann *featureAnnoRepo) RemoveFeatureAnnotation(
	fid string,
	purge bool,
) error
⋮----
func (fann *featureAnnoRepo) ClearFeatureAnnotations() error
⋮----
func (fann *featureAnnoRepo) AddTag(
	req *feature.AddTagRequest,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) UpdateTag(
	req *feature.UpdateTagRequest,
) (*model.FeatureAnnotationDoc, error)
⋮----
func (fann *featureAnnoRepo) RemoveTag(
	req *feature.RemoveTagRequest,
) error
⋮----
func (fann *featureAnnoRepo) Dbh() *manager.Database
⋮----
func (fann *featureAnnoRepo) upsertPublicationsTx(
	txr *manager.TransactionHandler,
	ids []string,
) ([]string, error)
⋮----
func (fann *featureAnnoRepo) createPublicationEdgesTx(
	txr *manager.TransactionHandler,
	featureKey string,
	pubKeys []string,
	source string,
) error
⋮----
func findTagPredicate(tag string) func(p model.TagPropertyDoc) bool
⋮----
func createFeatureAnnotationDoc(
	doc *feature.NewFeatureAnnotation,
) *model.FeatureAnnotationDoc
```
