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
	func(loader *PhenotypeLoader, pheno *phenotype.PhenotypeAnnotation) *PhenotypeLoader {
		payload := &PhenotypePayload{
			Reference: pheno.Reference(),
			CreatedOn: pheno.CreatedOn(),
		}
		if pheno.HasStrainId() {
			payload.StrainId = pheno.StrainId()
		}
		if pheno.HasStrainDescriptor() {
			payload.StrainDescriptor = pheno.StrainDescriptor()
		}
		loader.Payload = payload
		return loader
	},
)

var assayIdHandler = F.Curry2(
	func(assayId int, loader *PhenotypeLoader) *PhenotypeLoader {
		if assayId != 0 {
			loader.Payload.AssayId = []int{assayId}
		}
		return loader
	})

var envIdHandler = F.Curry2(
	func(envId int, loader *PhenotypeLoader) *PhenotypeLoader {
		if envId != 0 {
			loader.Payload.EnvironmentId = []int{envId}
		}
		return loader
	})

var phenoIdHandler = F.Curry2(
	func(phenoId int, loader *PhenotypeLoader) *PhenotypeLoader {
		loader.Payload.Id = []int{phenoId}
		return loader
	})

var assignedByIdHandler = F.Curry2(
	func(aid int, loader *PhenotypeLoader) *PhenotypeLoader {
		if aid != 0 {
			loader.Payload.AssignedBy = []common.AssignedBy{{Id: aid}}
		}
		return loader
	})

func environmentId(loader *PhenotypeLoader) E.Either[error, int] {
	if !loader.Annotation.HasEnvironmentId() {
		return E.Right[error](0)
	}
	envid, err := loader.TableManager.SearchRows(
		F.Pipe2(
			loader.Annotation.EnvironmentId(),
			common.ProcessOntologyTermId,
			common.ProcessEnvOntologyTerm,
		),
		loader.OntologyTableMap["env-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}
	return E.Right[error](envid)
}

func assayId(loader *PhenotypeLoader) E.Either[error, int] {
	if !loader.Annotation.HasAssayId() {
		return E.Right[error](0)
	}
	asid, err := loader.TableManager.SearchRows(
		common.ProcessOntologyTermId(loader.Annotation.AssayId()),
		loader.OntologyTableMap["assay-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}
	return E.Right[error](asid)
}

func assignedById(loader *PhenotypeLoader) E.Either[error, int] {
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

func phenotypeId(loader *PhenotypeLoader) E.Either[error, int] {
	phid, err := loader.TableManager.SearchRows(
		common.ProcessOntologyTermId(loader.Annotation.PhenotypeId()),
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
		Msg: fmt.Sprintf("created phenotype with annotation id %s", res.AnnoId),
	}
}
