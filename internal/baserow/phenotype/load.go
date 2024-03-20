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
	StrainId         string       `json:"strain_id"`
	StrainDescriptor string       `json:"strain_descriptor"`
	Id               []int        `json:"phenotype_id"`
	AssayId          []int        `json:"assay_id"`
	EnvironmentId    []int        `json:"environment_id"`
	Reference        string       `json:"reference"`
	Deleted          bool         `json:"deleted"`
	CreatedOn        time.Time    `json:"created_on"`
	AssignedBy       []AssignedBy `json:"assigned_by"`
}

type PhenotypeLoader struct {
	Workspace        string
	Host             string
	Token            string
	TableId          int
	Logger           *logrus.Entry
	OntologyTableMap map[string]int
	TableManager     *database.TableManager
	Payload          *PhenotypePayload
	Annotation       *phenotype.PhenotypeAnnotation
	WorkspaceManager *database.WorkspaceManager
}

type PhenotypeLoaderProperties struct {
	Workspace        string
	Host             string
	Token            string
	TableId          int
	Logger           *logrus.Entry
	OntologyTableMap map[string]int
	TableManager     *database.TableManager
	Payload          *PhenotypePayload
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
