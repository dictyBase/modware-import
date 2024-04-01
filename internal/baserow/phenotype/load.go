package phenotype

import (
	"bytes"
	"context"
	"fmt"
	"time"

	R "github.com/IBM/fp-go/context/readerioeither"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"golang.org/x/exp/slices"

	"github.com/dictyBase/modware-import/internal/baserow/common"
	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/dictyBase/modware-import/internal/concurrent"
	"github.com/dictyBase/modware-import/internal/datasource/xls/phenotype"
	"github.com/sirupsen/logrus"
)

const ConcurrentPhenoLoader = 10

type PhenotypePayload struct {
	Id               []int               `json:"phenotype_id"`
	Reference        string              `json:"reference"`
	CreatedOn        time.Time           `json:"created_on"`
	AssignedBy       []common.AssignedBy `json:"assigned_by"`
	StrainId         string              `json:"strain_id,omitempty"`
	StrainDescriptor string              `json:"strain_descriptor,omitempty"`
	AssayId          []int               `json:"assay_id,omitempty"`
	EnvironmentId    []int               `json:"environment_id,omitempty"`
	Deleted          bool                `json:"deleted,omitempty"`
}

type PhenotypeLoader struct {
	Workspace        string
	Host             string
	Token            string
	TableId          int
	OntologyTableMap map[string]int
	Payload          *PhenotypePayload
	Logger           *logrus.Entry
	TableManager     *database.TableManager
	Annotation       *phenotype.PhenotypeAnnotation
	WorkspaceManager *database.WorkspaceManager
}

type PhenotypeLoaderProperties struct {
	Workspace        string
	Host             string
	Token            string
	TableId          int
	OntologyTableMap map[string]int
	Payload          *PhenotypePayload
	Logger           *logrus.Entry
	TableManager     *database.TableManager
	Annotation       *phenotype.PhenotypeAnnotation
	WorkspaceManager *database.WorkspaceManager
}

// NewPhenotypeLoader creates a new instance of PhenotypeLoader with the provided properties.
func NewPhenotypeLoader(props *PhenotypeLoaderProperties) *PhenotypeLoader {
	return &PhenotypeLoader{
		Workspace:        props.Workspace,
		Host:             props.Host,
		Token:            props.Token,
		TableId:          props.TableId,
		Logger:           props.Logger,
		OntologyTableMap: props.OntologyTableMap,
		TableManager:     props.TableManager,
		WorkspaceManager: props.WorkspaceManager,
	}
}

// Load processes the phenotype data from the given reader and loads them into the database.
// It concurrently processes and loads data up to a set limit before waiting and continuing with the next batch.
func (loader *PhenotypeLoader) Load(
	reader *phenotype.PhenotypeAnnotationReader,
) error {
	tasks := make(
		[]concurrent.TaskWrapper[*phenotype.PhenotypeAnnotation, string],
		0,
		ConcurrentPhenoLoader,
	)
	for reader.Next() {
		pheno, err := reader.Value()
		if err != nil {
			return err
		}
		if pheno.IsEmpty() {
			continue
		}
		tasks = append(
			tasks,
			concurrent.TaskWrapper[*phenotype.PhenotypeAnnotation, string]{
				TaskFunc: loader.addPhenotypeRow,
				Input:    pheno,
			},
		)
		if len(tasks) == ConcurrentPhenoLoader {
			loader.Logger.Debug("going to load phenotypes")
			if err := concurrent.RunTasks(tasks, loader.Logger); err != nil {
				return err
			}
			tasks = slices.Delete(tasks, 0, len(tasks))
		}
	}
	loader.Logger.Infof("going to load remaining %d phenotypes", len(tasks))
	if err := concurrent.RunTasks(tasks, loader.Logger); err != nil {
		return err
	}

	return nil
}

func (loader *PhenotypeLoader) addPheno(
	pheno *phenotype.PhenotypeAnnotation,
) E.Either[error, *PhenotypeLoader] {
	newLoader := NewPhenotypeLoader(&PhenotypeLoaderProperties{
		Host:             loader.Host,
		Token:            loader.Token,
		Workspace:        loader.Workspace,
		TableId:          loader.TableId,
		Logger:           loader.Logger,
		OntologyTableMap: loader.OntologyTableMap,
		TableManager:     loader.TableManager,
		WorkspaceManager: loader.WorkspaceManager,
	})
	newLoader.Annotation = pheno
	return E.Right[error](newLoader)
}

func (loader *PhenotypeLoader) addPhenotypeRow(
	pheno *phenotype.PhenotypeAnnotation,
) (string, error) {
	var empty string
	content := F.Pipe8(
		E.Do[error](pheno),
		E.Bind(initialPayload, loader.addPheno),
		E.Bind(phenoIdHandler, phenotypeId),
		E.Bind(assayIdHandler, assayId),
		E.Bind(envIdHandler, environmentId),
		E.Bind(assignedByIdHandler, assignedById),
		E.Map[error, *PhenotypeLoader](loaderToPayload),
		E.Chain[error, *PhenotypePayload](common.MarshalPayload),
		E.Fold(httpapi.OnJSONPayloadError, httpapi.OnJSONPayloadSuccess),
	)
	if content.Error != nil {
		return empty, content.Error
	}
	resp := F.Pipe3(
		loader.createPhenotypeURL(),
		httpapi.MakeHTTPRequest("POST", bytes.NewBuffer(content.Payload)),
		R.Map(httpapi.SetHeaderWithJWT(loader.Token)),
		common.CreateHTTP,
	)(context.Background())
	output := F.Pipe1(
		resp(),
		E.Fold(common.OnCreateFeedbackError, onPhenoCreateFeedbackSuccess),
	)
	return output.Msg, output.Err
}

func loaderToPayload(loader *PhenotypeLoader) *PhenotypePayload {
	return loader.Payload
}

func (loader *PhenotypeLoader) createPhenotypeURL() string {
	return fmt.Sprintf(
		"https://%s/api/database/rows/table/%d/?user_field_names=true",
		loader.Host,
		loader.TableId,
	)
}
