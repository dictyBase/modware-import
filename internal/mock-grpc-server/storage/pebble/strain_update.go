package pebble

import (
	"fmt"

	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	T "github.com/IBM/fp-go/tuple"
	"github.com/cockroachdb/pebble"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

// Context types for UpdateStrain workflow
type updateStrainContext struct {
	req *stock.StrainUpdate
	db  *pebble.DB
}

type withExistingStrain struct {
	updateStrainContext
	existing *stock.Strain
}

type withUpdatedStrain struct {
	withExistingStrain
	updated *stock.Strain
}

type withUpdatedSerializedStrain struct {
	withUpdatedStrain
	protoBytes []byte
	jsonBytes  []byte
}

type withUpdateStrainBatch struct {
	withUpdatedSerializedStrain
	batch *pebble.Batch
}

// Curried setters for update context enrichment
var (
	setExistingStrain = F.Curry2(
		func(strain *stock.Strain, ctx updateStrainContext) withExistingStrain {
			return withExistingStrain{
				updateStrainContext: ctx,
				existing:            strain,
			}
		},
	)

	setUpdatedStrain = F.Curry2(
		func(strain *stock.Strain, ctx withExistingStrain) withUpdatedStrain {
			return withUpdatedStrain{
				withExistingStrain: ctx,
				updated:            strain,
			}
		},
	)

	setUpdatedSerializedStrain = F.Curry2(
		func(data T.Tuple2[[]byte, []byte], ctx withUpdatedStrain) withUpdatedSerializedStrain {
			return withUpdatedSerializedStrain{
				withUpdatedStrain: ctx,
				protoBytes:        data.F1,
				jsonBytes:         data.F2,
			}
		},
	)

	setUpdateStrainBatch = F.Curry2(
		func(batch *pebble.Batch, ctx withUpdatedSerializedStrain) withUpdateStrainBatch {
			return withUpdateStrainBatch{
				withUpdatedSerializedStrain: ctx,
				batch:                       batch,
			}
		},
	)
)

// retrieveExistingStrain retrieves the existing strain to update
func (storage *Storage) retrieveExistingStrain(
	ctx updateStrainContext,
) IOE.IOEither[error, *stock.Strain] {
	return storage.GetStrain(ctx.req.Data.Id)
}

// applyStrainUpdate applies the update to the existing strain
func applyStrainUpdate(ctx withExistingStrain) *stock.Strain {
	updated := ctx.existing
	updateAttrs := ctx.req.Data.Attributes

	// Update mutable fields
	if updateAttrs.UpdatedBy != "" {
		updated.Data.Attributes.UpdatedBy = updateAttrs.UpdatedBy
	}
	updated.Data.Attributes.UpdatedAt = nowTimestamp()

	if updateAttrs.Summary != "" {
		updated.Data.Attributes.Summary = updateAttrs.Summary
	}
	if updateAttrs.EditableSummary != "" {
		updated.Data.Attributes.EditableSummary = updateAttrs.EditableSummary
	}

	return updated
}

// serializeUpdatedStrainData serializes updated strain
func serializeUpdatedStrainData(
	ctx withUpdatedStrain,
) IOE.IOEither[error, T.Tuple2[[]byte, []byte]] {
	return IOE.TryCatchError(func() (T.Tuple2[[]byte, []byte], error) {
		protoBytes, err := serializeStrain(ctx.updated)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		jsonIndex := buildStrainJSONIndex(ctx.updated)
		jsonBytes, err := serializeJSON(jsonIndex)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		return T.MakeTuple2(protoBytes, jsonBytes), nil
	})
}

// buildUpdateStrainBatch creates batch for update operation
func buildUpdateStrainBatch(
	ctx withUpdatedSerializedStrain,
) IOE.IOEither[error, *pebble.Batch] {
	return IOE.TryCatchError(func() (*pebble.Batch, error) {
		batch := ctx.db.NewBatch()
		keys := newKeyBuilder()
		stockID := ctx.req.Data.Id

		// Update stock document
		if err := batch.Set(
			keys.stockKey(stockID),
			ctx.protoBytes,
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}

		// Update JSON index
		if err := batch.Set(
			keys.indexKey(stockID),
			ctx.jsonBytes,
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to update index: %w", err)
		}

		return batch, nil
	})
}

// commitUpdateStrainBatch commits the update batch
func commitUpdateStrainBatch(
	ctx withUpdateStrainBatch,
) IOE.IOEither[error, withUpdateStrainBatch] {
	return IOE.TryCatchError(func() (withUpdateStrainBatch, error) {
		if err := ctx.batch.Commit(pebble.Sync); err != nil {
			return ctx, fmt.Errorf("failed to commit update batch: %w", err)
		}
		return ctx, nil
	})
}

// extractUpdatedStrain extracts the updated strain
func extractUpdatedStrain(ctx withUpdateStrainBatch) *stock.Strain {
	return ctx.updated
}

// UpdateStrain updates an existing strain
func (storage *Storage) UpdateStrain(
	req *stock.StrainUpdate,
) IOE.IOEither[error, *stock.Strain] {
	return F.Pipe6(
		IOE.Of[error](updateStrainContext{req: req, db: storage.db}),
		IOE.Bind(setExistingStrain, storage.retrieveExistingStrain),
		IOE.Let[error](setUpdatedStrain, applyStrainUpdate),
		IOE.Bind(setUpdatedSerializedStrain, serializeUpdatedStrainData),
		IOE.Bind(setUpdateStrainBatch, buildUpdateStrainBatch),
		IOE.Chain(commitUpdateStrainBatch),
		IOE.Map[error](extractUpdatedStrain),
	)
}

// LoadStrain loads an existing strain with a specific ID
func (storage *Storage) LoadStrain(
	stockID string,
	req *stock.ExistingStrain,
) IOE.IOEither[error, *stock.Strain] {
	return F.Pipe7(
		IOE.Of[error](createStrainContext{
			req: convertExistingToNewStrain(req),
			db:  storage.db,
		}),
		IOE.Let[error](setGeneratedStrainID, func(_ createStrainContext) string {
			return stockID
		}),
		IOE.Bind(setStrainTimestamps, generateStrainTimestamps),
		IOE.Let[error](setBuiltStrain, buildStrainFromRequest),
		IOE.Bind(setSerializedStrain, serializeStrainData),
		IOE.Bind(setStrainBatch, buildStrainBatchWrite),
		IOE.Chain(commitStrainBatch),
		IOE.Map[error](extractCreatedStrain),
	)
}

// convertExistingToNewStrain converts ExistingStrain to NewStrain for loading
func convertExistingToNewStrain(
	existing *stock.ExistingStrain,
) *stock.NewStrain {
	attrs := existing.Data.Attributes
	return &stock.NewStrain{
		Data: &stock.NewStrain_Data{
			Type: "strain",
			Attributes: &stock.NewStrainAttributes{
				CreatedBy:           attrs.CreatedBy,
				UpdatedBy:           attrs.UpdatedBy,
				Summary:             attrs.Summary,
				EditableSummary:     attrs.EditableSummary,
				Depositor:           attrs.Depositor,
				Genes:               attrs.Genes,
				Dbxrefs:             attrs.Dbxrefs,
				Publications:        attrs.Publications,
				Label:               attrs.Label,
				Species:             attrs.Species,
				Plasmid:             attrs.Plasmid,
				Parent:              attrs.Parent,
				Names:               attrs.Names,
				DictyStrainProperty: attrs.DictyStrainProperty,
			},
		},
	}
}
