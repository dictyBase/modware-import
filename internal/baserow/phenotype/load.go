package phenotype

import (
	"time"

	"github.com/dictyBase/modware-import/internal/baserow/database"
	"github.com/dictyBase/modware-import/internal/datasource/xls/phenotype"
	"github.com/sirupsen/logrus"
)

const ConcurrentPhenoLoader = 10

type AssignedBy struct {
	Id int `json:"id"`
}

type PhenotypePayload struct {
	StrainDescriptor string       `json:"strain_descriptor"`
	Id               []int        `json:"phenotype_id"`
	Reference        string       `json:"reference"`
	CreatedOn        time.Time    `json:"created_on"`
	AssignedBy       []AssignedBy `json:"assigned_by"`
	StrainId         string       `json:"strain_id,omitempty"`
	AssayId          []int        `json:"assay_id,omitempty"`
	EnvironmentId    []int        `json:"environment_id,omitempty"`
	Deleted          bool         `json:"deleted,omitempty"`
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
