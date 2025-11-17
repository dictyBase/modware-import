package phenotype

import (
	"fmt"

	"github.com/dictyBase/modware-import/internal/baserow/common"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"

	"github.com/dictyBase/modware-import/internal/datasource/xls/phenotype"
)

var initialPayload = F.Curry2(
	func(loader *PhenotypeLoader, pheno *phenotype.Annotation) *PhenotypeLoader {
		payload := &PhenotypePayload{
			Reference: pheno.Reference(),
			CreatedOn: pheno.CreatedOn(),
		}
		if pheno.HasStrainID() {
			payload.StrainID = pheno.StrainID()
		}
		if pheno.HasStrainDescriptor() {
			payload.StrainDescriptor = pheno.StrainDescriptor()
		}
		loader.Payload = payload
		return loader
	},
)

var assayIDHandler = F.Curry2(
	func(assayID int, loader *PhenotypeLoader) *PhenotypeLoader {
		if assayID != 0 {
			loader.Payload.AssayID = []int{assayID}
		}
		return loader
	})

var envIDHandler = F.Curry2(
	func(envID int, loader *PhenotypeLoader) *PhenotypeLoader {
		if envID != 0 {
			loader.Payload.EnvironmentID = []int{envID}
		}
		return loader
	})

var phenoIDHandler = F.Curry2(
	func(phenoID int, loader *PhenotypeLoader) *PhenotypeLoader {
		loader.Payload.ID = []int{phenoID}
		return loader
	})

var assignedByIDHandler = F.Curry2(
	func(aid int, loader *PhenotypeLoader) *PhenotypeLoader {
		if aid != 0 {
			loader.Payload.AssignedBy = []common.AssignedBy{{ID: aid}}
		}
		return loader
	})

func environmentID(loader *PhenotypeLoader) E.Either[error, int] {
	if !loader.Annotation.HasEnvironmentID() {
		return E.Right[error](0)
	}
	envid, err := loader.TableManager.SearchRows(
		F.Pipe2(
			loader.Annotation.EnvironmentID(),
			common.ProcessOntologyTermID,
			common.ProcessEnvOntologyTerm,
		),
		loader.OntologyTableMap["env-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}
	return E.Right[error](envid)
}

func assayID(loader *PhenotypeLoader) E.Either[error, int] {
	if !loader.Annotation.HasAssayID() {
		return E.Right[error](0)
	}
	asid, err := loader.TableManager.SearchRows(
		common.ProcessOntologyTermID(loader.Annotation.AssayID()),
		loader.OntologyTableMap["assay-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}
	return E.Right[error](asid)
}

func assignedByID(loader *PhenotypeLoader) E.Either[error, int] {
	ok, aid, err := loader.WorkspaceManager.SearchWorkspaceUser(
		loader.Workspace, loader.Annotation.AssignedBy(),
	)
	if err != nil {
		return E.Left[int](err)
	}
	if !ok {
		return E.Right[error](0)
	}

	return E.Right[error](aid)
}

func phenotypeID(loader *PhenotypeLoader) E.Either[error, int] {
	phid, err := loader.TableManager.SearchRows(
		common.ProcessOntologyTermID(loader.Annotation.PhenotypeID()),
		loader.OntologyTableMap["phenotype-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}
	return E.Right[error](phid)
}

func onPhenoCreateFeedbackSuccess(
	res common.CreateResp,
) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf("created phenotype with annotation id %s", res.AnnoID),
	}
}
