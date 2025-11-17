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

// Context types for CreateStrain workflow
type createStrainContext struct {
	req *stock.NewStrain
	db  *pebble.DB
}

type withGeneratedStrainID struct {
	createStrainContext
	generatedID string
}

type withStrainTimestamps struct {
	withGeneratedStrainID
	createdAt time.Time
	updatedAt time.Time
}

type withBuiltStrain struct {
	withStrainTimestamps
	strain *stock.Strain
}

type withSerializedStrain struct {
	withBuiltStrain
	protoBytes []byte
	jsonBytes  []byte
}

type withStrainBatch struct {
	withSerializedStrain
	batch *pebble.Batch
}

// Curried setters for context enrichment
var (
	setGeneratedStrainID = F.Curry2(
		func(stockID string, ctx createStrainContext) withGeneratedStrainID {
			return withGeneratedStrainID{
				createStrainContext: ctx,
				generatedID:         stockID,
			}
		},
	)

	setStrainTimestamps = F.Curry2(
		func(timestamps T.Tuple2[time.Time, time.Time], ctx withGeneratedStrainID) withStrainTimestamps {
			return withStrainTimestamps{
				withGeneratedStrainID: ctx,
				createdAt:             timestamps.F1,
				updatedAt:             timestamps.F2,
			}
		},
	)

	setBuiltStrain = F.Curry2(
		func(strain *stock.Strain, ctx withStrainTimestamps) withBuiltStrain {
			return withBuiltStrain{
				withStrainTimestamps: ctx,
				strain:               strain,
			}
		},
	)

	setSerializedStrain = F.Curry2(
		func(data T.Tuple2[[]byte, []byte], ctx withBuiltStrain) withSerializedStrain {
			return withSerializedStrain{
				withBuiltStrain: ctx,
				protoBytes:      data.F1,
				jsonBytes:       data.F2,
			}
		},
	)

	setStrainBatch = F.Curry2(
		func(batch *pebble.Batch, ctx withSerializedStrain) withStrainBatch {
			return withStrainBatch{
				withSerializedStrain: ctx,
				batch:                batch,
			}
		},
	)
)

// generateStrainID generates the next sequential strain ID
func generateStrainID(ctx createStrainContext) IOE.IOEither[error, string] {
	return IOE.TryCatchError(func() (string, error) {
		counterKey := []byte(strainCounterKey)
		data, closer, err := ctx.db.Get(counterKey)

		var nextID int64 = 1

		if err != nil {
			if !errors.Is(err, pebble.ErrNotFound) {
				return "", fmt.Errorf("failed to read strain counter: %w", err)
			}
		} else {
			nextID = decodeCounter(data) + 1
			closer.Close()
		}

		return formatStockID("DBS", nextID), nil
	})
}

// generateStrainTimestamps creates current timestamps
func generateStrainTimestamps(
	_ withGeneratedStrainID,
) IOE.IOEither[error, T.Tuple2[time.Time, time.Time]] {
	return func() E.Either[error, T.Tuple2[time.Time, time.Time]] {
		now := time.Now()
		return E.Right[error](T.MakeTuple2(now, now))
	}
}

// buildStrainFromRequest constructs a strain from the request
func buildStrainFromRequest(ctx withStrainTimestamps) *stock.Strain {
	req := ctx.req
	attrs := req.Data.Attributes

	return &stock.Strain{
		Data: &stock.Strain_Data{
			Type: "strain",
			Id:   ctx.generatedID,
			Attributes: &stock.StrainAttributes{
				CreatedBy:           attrs.CreatedBy,
				UpdatedBy:           attrs.UpdatedBy,
				CreatedAt:           nowTimestamp(),
				UpdatedAt:           nowTimestamp(),
				Depositor:           attrs.Depositor,
				Summary:             attrs.Summary,
				EditableSummary:     attrs.EditableSummary,
				Species:             attrs.Species,
				Label:               attrs.Label,
				Names:               attrs.Names,
				DictyStrainProperty: attrs.DictyStrainProperty,
				Genes:               attrs.Genes,
				Dbxrefs:             attrs.Dbxrefs,
				Publications:        attrs.Publications,
				Plasmid:             attrs.Plasmid,
				Parent:              attrs.Parent,
			},
		},
	}
}

// serializeStrainData serializes strain to protobuf and JSON
func serializeStrainData(
	ctx withBuiltStrain,
) IOE.IOEither[error, T.Tuple2[[]byte, []byte]] {
	return IOE.TryCatchError(func() (T.Tuple2[[]byte, []byte], error) {
		protoBytes, err := serializeStrain(ctx.strain)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		jsonIndex := buildStrainJSONIndex(ctx.strain)
		jsonBytes, err := serializeJSON(jsonIndex)
		if err != nil {
			return T.MakeTuple2[[]byte, []byte](nil, nil), err
		}

		return T.MakeTuple2(protoBytes, jsonBytes), nil
	})
}

// buildStrainBatchWrite creates pebble batch with all keys
func buildStrainBatchWrite(ctx withSerializedStrain) IOE.IOEither[error, *pebble.Batch] {
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
			[]byte("strain"),
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to set type: %w", err)
		}

		// Update counter with the numeric ID from generatedID
		counterValue := parseStockIDNumber(ctx.generatedID, stockIDPrefixLen)
		if err := batch.Set(
			keys.strainCounterKey(),
			encodeCounter(counterValue),
			pebble.Sync,
		); err != nil {
			return nil, fmt.Errorf("failed to update counter: %w", err)
		}

		// Write depositor reverse index
		if ctx.strain.Data.Attributes.Depositor != "" {
			if err := batch.Set(
				keys.depositorKey(ctx.strain.Data.Attributes.Depositor, ctx.generatedID),
				[]byte(""),
				pebble.Sync,
			); err != nil {
				return nil, fmt.Errorf("failed to set depositor index: %w", err)
			}
		}

		// Write species reverse index
		if ctx.strain.Data.Attributes.Species != "" {
			if err := batch.Set(
				keys.speciesKey(ctx.strain.Data.Attributes.Species, ctx.generatedID),
				[]byte(""),
				pebble.Sync,
			); err != nil {
				return nil, fmt.Errorf("failed to set species index: %w", err)
			}
		}

		return batch, nil
	})
}

// commitStrainBatch commits the pebble batch
func commitStrainBatch(ctx withStrainBatch) IOE.IOEither[error, withStrainBatch] {
	return IOE.TryCatchError(func() (withStrainBatch, error) {
		if err := ctx.batch.Commit(pebble.Sync); err != nil {
			return ctx, fmt.Errorf("failed to commit batch: %w", err)
		}
		return ctx, nil
	})
}

// extractCreatedStrain extracts the created strain from context
func extractCreatedStrain(ctx withStrainBatch) *stock.Strain {
	return ctx.strain
}

// CreateStrain creates a new strain using the fp-go pipeline
func (storage *Storage) CreateStrain(
	req *stock.NewStrain,
) IOE.IOEither[error, *stock.Strain] {
	return func() E.Either[error, *stock.Strain] {
		// Lock for atomic ID generation
		storage.mu.Lock()
		defer storage.mu.Unlock()

		return F.Pipe7(
			IOE.Of[error](createStrainContext{req: req, db: storage.db}),
			IOE.Bind(setGeneratedStrainID, generateStrainID),
			IOE.Bind(setStrainTimestamps, generateStrainTimestamps),
			IOE.Let[error](setBuiltStrain, buildStrainFromRequest),
			IOE.Bind(setSerializedStrain, serializeStrainData),
			IOE.Bind(setStrainBatch, buildStrainBatchWrite),
			IOE.Chain(commitStrainBatch),
			IOE.Map[error](extractCreatedStrain),
		)()
	}
}
