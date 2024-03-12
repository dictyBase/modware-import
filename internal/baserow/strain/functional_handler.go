package strain

import (
	"fmt"
	"net/http"
	"strings"

	J "github.com/IBM/fp-go/json"

	F "github.com/IBM/fp-go/function"

	H "github.com/IBM/fp-go/context/readerioeither/http"
	E "github.com/IBM/fp-go/either"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/dictyBase/modware-import/internal/datasource/xls/strain"
)

var (
	strainCreateHTTP = H.ReadJSON[strainCreateResp](
		H.MakeClient(http.DefaultClient),
	)
)

var mutagenesisIdHandler = F.Curry2(
	func(mutId int, loader *StrainLoader) *StrainLoader {
		if mutId != 0 {
			loader.Payload.GeneticModificationId = []int{mutId}
		}

		return loader
	})

var charIdsHandler = F.Curry2(
	func(charIds []int, loader *StrainLoader) *StrainLoader {
		if len(charIds) != 0 {
			loader.Payload.GeneticModificationId = charIds
		}

		return loader
	})

var genModIdHandler = F.Curry2(
	func(genmodId int, loader *StrainLoader) *StrainLoader {
		if genmodId != 0 {
			loader.Payload.GeneticModificationId = []int{genmodId}
		}

		return loader
	})

var initialPayload = F.Curry2(
	func(loader *StrainLoader, strn *strain.StrainAnnotation) *StrainLoader {
		payload := &StrainPayload{
			Descriptor: strn.Descriptor(),
			Reference:  strn.Reference(),
			Species:    strn.Species(),
			Summary:    strn.Summary(),
		}
		if strn.HasName() {
			payload.Names = strn.Name()
		}
		if strn.HasSystematicName() {
			payload.SystematicName = strn.SystematicName()
		}
		if strn.HasPlasmid() {
			payload.Plasmid = strn.Plasmid()
		}
		if strn.HasParentId() {
			payload.ParentId = strn.ParentId()
		}
		if strn.HasGenes() {
			payload.Genes = strn.Genes()
		}
		if strn.HasGenotype() {
			payload.Genotype = strn.Genotype()
		}
		if strn.HasDepositor() {
			payload.Depositor = strn.Depositor()
		}
		loader.Payload = payload

		return loader
	})

type strainCreateResp struct {
	AnnoId string `json:"annotation_id"`
}

func onStrainCreateFeedbackSuccess(
	res strainCreateResp,
) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf("created strain with annotation id %s", res.AnnoId),
	}
}

func onStrainCreateFeedbackError(err error) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{Err: err}
}

func processOntologyTermId(val string) string {
	return strings.Replace(val, ":", "_", 1)
}

func mutagenesisId(loader *StrainLoader) E.Either[error, int] {
	mutId, err := loader.TableManager.SearchRows(
		processOntologyTermId(loader.Annotation.MutagenesisMethod()),
		loader.OntologyTableMap["mutagenesis-method-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}

	return E.Right[error](mutId)
}

func genmodId(loader *StrainLoader) E.Either[error, int] {
	mutId, err := loader.TableManager.SearchRows(
		processOntologyTermId(loader.Annotation.GeneticModification()),
		loader.OntologyTableMap["genetic-mod-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}

	return E.Right[error](mutId)
}

func characteristicIds(loader *StrainLoader) E.Either[error, []int] {
	charIds := make([]int, 0)
	for _, charac := range strings.Split(loader.Annotation.Characteristic(), ",") {
		cid, err := loader.TableManager.SearchRows(
			processOntologyTermId(charac),
			loader.OntologyTableMap["strainchar-ontology-table"],
		)
		if err != nil {
			return E.Left[[]int](err)
		}
		charIds = append(charIds, cid)
	}

	return E.Right[error](charIds)
}

func loaderToPayload(ldr *StrainLoader) *StrainPayload { return ldr.Payload }

func marshalPayload(payload *StrainPayload) E.Either[error, []byte] {
	return F.Pipe1(payload, J.Marshal)
}
