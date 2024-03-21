package strain

import (
	"fmt"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/baserow/common"

	F "github.com/IBM/fp-go/function"

	E "github.com/IBM/fp-go/either"
	"github.com/dictyBase/modware-import/internal/baserow/httpapi"
	"github.com/dictyBase/modware-import/internal/datasource/xls/strain"
)

var assignedByIdHandler = F.Curry2(
	func(aid int, loader *StrainLoader) *StrainLoader {
		if aid != 0 {
			loader.Payload.AssignedBy = []AssignedBy{{Id: aid}}
		}
		return loader
	})

var mutagenesisIdHandler = F.Curry2(
	func(mutId int, loader *StrainLoader) *StrainLoader {
		if mutId != 0 {
			loader.Payload.MutagenesisMethodId = []int{mutId}
		}

		return loader
	})

var charIdsHandler = F.Curry2(
	func(charIds []int, loader *StrainLoader) *StrainLoader {
		if len(charIds) != 0 {
			loader.Payload.StrainCharacteristicsId = charIds
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

var creationTimeHandler = F.Curry2(
	func(createdOn time.Time, loader *StrainLoader) *StrainLoader {
		loader.Payload.CreatedOn = createdOn
		return loader
	})
var creationTime = F.Curry2(
	func(createdOn time.Time, loader *StrainLoader) E.Either[error, time.Time] {
		return E.Right[error](createdOn)
	})

func onStrainCreateFeedbackSuccess(
	res common.CreateResp,
) httpapi.ResponseFeedback {
	return httpapi.ResponseFeedback{
		Msg: fmt.Sprintf("created strain with annotation id %s", res.AnnoId),
	}
}

func assignedById(loader *StrainLoader) E.Either[error, int] {
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

func mutagenesisId(loader *StrainLoader) E.Either[error, int] {
	mutId, err := loader.TableManager.SearchRows(
		common.ProcessOntologyTermId(loader.Annotation.MutagenesisMethod()),
		loader.OntologyTableMap["mutagenesis-method-ontology-table"],
	)
	if err != nil {
		return E.Left[int](err)
	}

	return E.Right[error](mutId)
}

func genmodId(loader *StrainLoader) E.Either[error, int] {
	mutId, err := loader.TableManager.SearchRows(
		common.ProcessOntologyTermId(loader.Annotation.GeneticModification()),
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
			common.ProcessOntologyTermId(charac),
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
