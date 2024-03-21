package phenotype

import (
	"github.com/dictyBase/modware-import/internal/baserow/common"
	"fmt"
	"net/http"
	"strings"

	"github.com/dictyBase/modware-import/internal/baserow/httpapi"

	H "github.com/IBM/fp-go/context/readerioeither/http"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	J "github.com/IBM/fp-go/json"

	"github.com/dictyBase/modware-import/internal/datasource/xls/phenotype"
)

var (
	phenoCreateHTTP = H.ReadJSON[phenoCreateResp](
		H.MakeClient(http.DefaultClient),
	)
)

var initialPayload = F.Curry2(
	func(loader *PhenotypeLoader, pheno *phenotype.PhenotypeAnnotation) *PhenotypeLoader {
		payload := &PhenotypePayload{
			StrainDescriptor: pheno.StrainDescriptor(),
			Reference:        pheno.Reference(),
			CreatedOn:        pheno.CreatedOn(),
		}
		if pheno.HasStrainId() {
			payload.StrainId = pheno.StrainId()
		}
		loader.Payload = payload
		return loader
	},
)

var assayIdHandler = F.Curry2(
	func(assayId []int, loader *PhenotypeLoader) *PhenotypeLoader {
		if len(assayId) != 0 {
			loader.Payload.AssayId = assayId
		}
		return loader
	})

var envIdHandler = F.Curry2(
	func(envId []int, loader *PhenotypeLoader) *PhenotypeLoader {
		loader.Payload.EnvironmentId = envId
		return loader
	})

var phenoIdHandler = F.Curry2(
	func(phenoId []int, loader *PhenotypeLoader) *PhenotypeLoader {
		loader.Payload.Id = phenoId
		return loader
	})

type phenoCreateResp struct {
	AnnoId string `json:"annotation_id"`
}

func environmentId(loader *PhenotypeLoader) E.Either[error, []int] {
	if !loader.Annotation.HasEnvironmentId() {
		return E.Right[error]([]int{0})
	}
	envid, err := loader.TableManager.SearchRows(
		processOntologyTermId(loader.Annotation.EnvironmentId()),
		loader.OntologyTableMap["env-ontology-table"],
	)
	if err != nil {
		return E.Left[[]int](err)
	}
	return E.Right[error]([]int{envid})
}

func assayId(loader *PhenotypeLoader) E.Either[error, []int] {
	if !loader.Annotation.HasAssayId() {
		return E.Right[error]([]int{0})
	}
	asid, err := loader.TableManager.SearchRows(
		processOntologyTermId(loader.Annotation.AssayId()),
		loader.OntologyTableMap["assay-ontology-table"],
	)
	if err != nil {
		return E.Left[[]int](err)
	}
	return E.Right[error]([]int{asid})
}

func phenotypeId(loader *PhenotypeLoader) E.Either[error, []int] {
	phid, err := loader.TableManager.SearchRows(
		processOntologyTermId(loader.Annotation.PhenotypeId()),
		loader.OntologyTableMap["phenotype-ontology-table"],
	)
	if err != nil {
		return E.Left[[]int](err)
	}
	return E.Right[error]([]int{phid})
}


// Use common.MarshalPayload instead

func onPhenoCreateFeedbackSuccess(
	res phenoCreateResp,
) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf("created phenotype with annotation id %s", res.AnnoId),
	}
}

func onPhenoCreateFeedbackError(err error) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{Err: err}
}
