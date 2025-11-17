package pebble

import (
	"errors"
	"fmt"
	"time"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	T "github.com/IBM/fp-go/tuple"
	"github.com/cockroachdb/pebble"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
)

// ==================== CREATE PLASMID ====================

type createPlasmidContext struct {
	req *stock.NewPlasmid
	db  *pebble.DB
}

type withGeneratedPlasmidID struct {
	createPlasmidContext
	generatedID string
}

type withPlasmidTimestamps struct {
	withGeneratedPlasmidID
	createdAt time.Time
	updatedAt time.Time
}

type withBuiltPlasmid struct {
	withPlasmidTimestamps
	plasmid *stock.Plasmid
}

type withSerializedPlasmid struct {
	withBuiltPlasmid
	protoBytes []byte
	jsonBytes  []byte
}

type withPlasmidBatch struct {
	withSerializedPlasmid
	batch *pebble.Batch
}

var (
	setGeneratedPlasmidID = F.Curry2(
		func(stockID string, ctx createPlasmidContext) withGeneratedPlasmidID {
			return withGeneratedPlasmidID{
				createPlasmidContext: ctx,
				generatedID:          stockID,
			}
		},
	)

	setPlasmidTimestamps = F.Curry2(
		func(timestamps T.Tuple2[time.Time, time.Time], ctx withGeneratedPlasmidID) withPlasmidTimestamps {
			return withPlasmidTimestamps{
				withGeneratedPlasmidID: ctx,
				createdAt:              timestamps.F1,
				updatedAt:              timestamps.F2,
			}
		},
	)

	setBuiltPlasmid = F.Curry2(
		func(plasmid *stock.Plasmid, ctx withPlasmidTimestamps) withBuiltPlasmid {
			return withBuiltPlasmid{
				withPlasmidTimestamps: ctx,
				plasmid:               plasmid,
			}
		},
	)

	setSerializedPlasmid = F.Curry2(
		func(data T.Tuple2[[]byte, []byte], ctx withBuiltPlasmid) withSerializedPlasmid {
			return withSerializedPlasmid{
				withBuiltPlasmid: ctx,
				protoBytes:       data.F1,
				jsonBytes:        data.F2,
			}
		},
	)

	setPlasmidBatch = F.Curry2(
		func(batch *pebble.Batch, ctx withSerializedPlasmid) withPlasmidBatch {
			return withPlasmidBatch{
				withSerializedPlasmid: ctx,
				batch:                 batch,
			}
		},
	)
)

func generatePlasmidID(ctx createPlasmidContext) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		counterKey := []byte(plasmidCounterKey)
		data, closer, err := ctx.db.Get(counterKey)

		var nextID int64 = 1

		if err != nil {
			if !errors.Is(err, pebble.ErrNotFound) {
				return "", fmt.Errorf("failed to read plasmid counter: %w", err)
			}
		} else {
			nextID = decodeCounter(data) + 1
			closer.Close()
		}

		return formatStockID("DBP", nextID), nil
	})
}

func generatePlasmidTimestamps(
	_ withGeneratedPlasmidID,
) IOE.IOEither[error, T.Tuple2[time.Time, time.Time]] {
	return func() E.Either[error, T.Tuple2[time.Time, time.Time]] {
		now := time.Now()
		return E.Right[error](T.MakeTuple2(now, now))
	}
}

func buildPlasmidFromRequest(ctx withPlasmidTimestamps) *stock.Plasmid {
	req := ctx.req
	attrs := req.Data.Attributes

	return &stock.Plasmid{
		Data: &stock.Plasmid_Data{
			Type: "plasmid",
			Id:   ctx.generatedID,
			Attributes: &stock.PlasmidAttributes{
				CreatedBy:            attrs.CreatedBy,
				UpdatedBy:            attrs.UpdatedBy,
				CreatedAt:            nowTimestamp(),
				UpdatedAt:            nowTimestamp(),
				Depositor:            attrs.Depositor,
				Summary:              attrs.Summary,
				EditableSummary:      attrs.EditableSummary,
				Name:                 attrs.Name,
				ImageMap:             attrs.ImageMap,
				Sequence:             attrs.Sequence,
				DictyPlasmidProperty: attrs.DictyPlasmidProperty,
			},
		},
	}
}

func serializePlasmidData(
	ctx withBuiltPlasmid,
) IOE.IOEither[error, T.Tuple2[[]byte, []byte]] {
	return IOE.TryCatchError(func() (T.Tuple2[[]byte, []byte], error) {
		protoBytes, err := serializePlasmid(ctx.plasmid)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		jsonIndex := buildPlasmidJSONIndex(ctx.plasmid)
		jsonBytes, err := serializeJSON(jsonIndex)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		return T.MakeTuple2(protoBytes, jsonBytes), nil
	})
}

func buildPlasmidBatchWrite(ctx withSerializedPlasmid) IOE.IOEither[error, *pebble.Batch] {
	return IOE.TryCatchError(func() (*pebble.Batch, error) {
		batch := ctx.db.NewBatch()
		keys := newKeyBuilder()

		// Write stock document
		if err := batch.Set(
			keys.stockKey(ctx.generatedID),
			ctx.protoBytes,
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to set stock: %w", err)
		}

		// Write JSON index
		if err := batch.Set(
			keys.indexKey(ctx.generatedID),
			ctx.jsonBytes,
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to set index: %w", err)
		}

		// Write type classification
		if err := batch.Set(
			keys.typeKey(ctx.generatedID),
			[]byte("plasmid"),
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to set type: %w", err)
		}

		// Update counter
		counterValue := decodeCounter([]byte{}) + 1
		if err := batch.Set(
			keys.plasmidCounterKey(),
			encodeCounter(counterValue),
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to update counter: %w", err)
		}

		// Write depositor reverse index
		if ctx.plasmid.Data.Attributes.Depositor != "" {
			if err := batch.Set(
				keys.depositorKey(ctx.plasmid.Data.Attributes.Depositor, ctx.generatedID),
				[]byte(""),
				pebble.Sync,
			); err != nil {
				return nil, fmt.Errorf("failed to set depositor index: %w", err)
			}
		}

		return batch, nil
	})
}

func commitPlasmidBatch(ctx withPlasmidBatch) IOE.IOEither[error, withPlasmidBatch] {
	return IOE.TryCatchError(func() (withPlasmidBatch, error) {
		if err := ctx.batch.Commit(pebble.Sync); err != nil {
			return ctx, fmt.Errorf("failed to commit batch: %w", err)
		}
		return ctx, nil
	})
}

func extractCreatedPlasmid(ctx withPlasmidBatch) *stock.Plasmid {
	return ctx.plasmid
}

func (storage *pebbleStorage) CreatePlasmid(
	req *stock.NewPlasmid,
) IOE.IOEither[error, *stock.Plasmid] {
	return F.Pipe7(
		IOE.Of[error](createPlasmidContext{req: req, db: storage.db}),
		IOE.Bind(setGeneratedPlasmidID, generatePlasmidID),
		IOE.Bind(setPlasmidTimestamps, generatePlasmidTimestamps),
		IOE.Let[error](setBuiltPlasmid, buildPlasmidFromRequest),
		IOE.Bind(setSerializedPlasmid, serializePlasmidData),
		IOE.Bind(setPlasmidBatch, buildPlasmidBatchWrite),
		IOE.Chain(commitPlasmidBatch),
		IOE.Map[error](extractCreatedPlasmid),
	)
}

// ==================== UPDATE PLASMID ====================

type updatePlasmidContext struct {
	req *stock.PlasmidUpdate
	db  *pebble.DB
}

type withExistingPlasmid struct {
	updatePlasmidContext
	existing *stock.Plasmid
}

type withUpdatedPlasmid struct {
	withExistingPlasmid
	updated *stock.Plasmid
}

type withUpdatedSerializedPlasmid struct {
	withUpdatedPlasmid
	protoBytes []byte
	jsonBytes  []byte
}

type withUpdatePlasmidBatch struct {
	withUpdatedSerializedPlasmid
	batch *pebble.Batch
}

var (
	setExistingPlasmid = F.Curry2(
		func(plasmid *stock.Plasmid, ctx updatePlasmidContext) withExistingPlasmid {
			return withExistingPlasmid{
				updatePlasmidContext: ctx,
				existing:             plasmid,
			}
		},
	)

	setUpdatedPlasmid = F.Curry2(
		func(plasmid *stock.Plasmid, ctx withExistingPlasmid) withUpdatedPlasmid {
			return withUpdatedPlasmid{
				withExistingPlasmid: ctx,
				updated:             plasmid,
			}
		},
	)

	setUpdatedSerializedPlasmid = F.Curry2(
		func(data T.Tuple2[[]byte, []byte], ctx withUpdatedPlasmid) withUpdatedSerializedPlasmid {
			return withUpdatedSerializedPlasmid{
				withUpdatedPlasmid: ctx,
				protoBytes:         data.F1,
				jsonBytes:          data.F2,
			}
		},
	)

	setUpdatePlasmidBatch = F.Curry2(
		func(batch *pebble.Batch, ctx withUpdatedSerializedPlasmid) withUpdatePlasmidBatch {
			return withUpdatePlasmidBatch{
				withUpdatedSerializedPlasmid: ctx,
				batch:                        batch,
			}
		},
	)
)

func (storage *pebbleStorage) retrieveExistingPlasmid(
	ctx updatePlasmidContext,
) IOE.IOEither[error, *stock.Plasmid] {
	return storage.GetPlasmid(ctx.req.Data.Id)
}

func applyPlasmidUpdate(ctx withExistingPlasmid) *stock.Plasmid {
	updated := ctx.existing
	updateAttrs := ctx.req.Data.Attributes

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

func serializeUpdatedPlasmidData(
	ctx withUpdatedPlasmid,
) IOE.IOEither[error, T.Tuple2[[]byte, []byte]] {
	return IOE.TryCatchError(func() (T.Tuple2[[]byte, []byte], error) {
		protoBytes, err := serializePlasmid(ctx.updated)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		jsonIndex := buildPlasmidJSONIndex(ctx.updated)
		jsonBytes, err := serializeJSON(jsonIndex)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		return T.MakeTuple2(protoBytes, jsonBytes), nil
	})
}

func buildUpdatePlasmidBatch(
	ctx withUpdatedSerializedPlasmid,
) IOE.IOEither[error, *pebble.Batch] {
	return IOE.TryCatchError(func() (*pebble.Batch, error) {
		batch := ctx.db.NewBatch()
		keys := newKeyBuilder()
		stockID := ctx.req.Data.Id

		if err := batch.Set(
			keys.stockKey(stockID),
			ctx.protoBytes,
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}

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

func commitUpdatePlasmidBatch(
	ctx withUpdatePlasmidBatch,
) IOE.IOEither[error, withUpdatePlasmidBatch] {
	return IOE.TryCatchError(func() (withUpdatePlasmidBatch, error) {
		if err := ctx.batch.Commit(pebble.Sync); err != nil {
			return ctx, fmt.Errorf("failed to commit update batch: %w", err)
		}
		return ctx, nil
	})
}

func extractUpdatedPlasmid(ctx withUpdatePlasmidBatch) *stock.Plasmid {
	return ctx.updated
}

func (storage *pebbleStorage) UpdatePlasmid(
	req *stock.PlasmidUpdate,
) IOE.IOEither[error, *stock.Plasmid] {
	return F.Pipe6(
		IOE.Of[error](updatePlasmidContext{req: req, db: storage.db}),
		IOE.Bind(setExistingPlasmid, storage.retrieveExistingPlasmid),
		IOE.Let[error](setUpdatedPlasmid, applyPlasmidUpdate),
		IOE.Bind(setUpdatedSerializedPlasmid, serializeUpdatedPlasmidData),
		IOE.Bind(setUpdatePlasmidBatch, buildUpdatePlasmidBatch),
		IOE.Chain(commitUpdatePlasmidBatch),
		IOE.Map[error](extractUpdatedPlasmid),
	)
}

// ==================== LOAD PLASMID ====================

func (storage *pebbleStorage) LoadPlasmid(
	stockID string,
	req *stock.ExistingPlasmid,
) IOE.IOEither[error, *stock.Plasmid] {
	return F.Pipe7(
		IOE.Of[error](createPlasmidContext{
			req: convertExistingToNewPlasmid(req),
			db:  storage.db,
		}),
		IOE.Let[error](setGeneratedPlasmidID, func(_ createPlasmidContext) string {
			return stockID
		}),
		IOE.Bind(setPlasmidTimestamps, generatePlasmidTimestamps),
		IOE.Let[error](setBuiltPlasmid, buildPlasmidFromRequest),
		IOE.Bind(setSerializedPlasmid, serializePlasmidData),
		IOE.Bind(setPlasmidBatch, buildPlasmidBatchWrite),
		IOE.Chain(commitPlasmidBatch),
		IOE.Map[error](extractCreatedPlasmid),
	)
}

func convertExistingToNewPlasmid(
	existing *stock.ExistingPlasmid,
) *stock.NewPlasmid {
	attrs := existing.Data.Attributes
	return &stock.NewPlasmid{
		Data: &stock.NewPlasmid_Data{
			Type: "plasmid",
			Attributes: &stock.NewPlasmidAttributes{
				CreatedBy:            attrs.CreatedBy,
				UpdatedBy:            attrs.UpdatedBy,
				Summary:              attrs.Summary,
				EditableSummary:      attrs.EditableSummary,
				Depositor:            attrs.Depositor,
				Genes:                attrs.Genes,
				Dbxrefs:              attrs.Dbxrefs,
				Publications:         attrs.Publications,
				ImageMap:             attrs.ImageMap,
				Sequence:             attrs.Sequence,
				Name:                 attrs.Name,
				DictyPlasmidProperty: attrs.DictyPlasmidProperty,
			},
		},
	}
}
