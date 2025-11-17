package stock

import (
	"context"
	"fmt"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	T "github.com/IBM/fp-go/tuple"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ==================== TYPE DEFINITIONS ====================

type (
	PlasmidResult    = T.Tuple2[*stock.Plasmid, error]
	PlasmidEither    = E.Either[error, *stock.Plasmid]
	PlasmidIO        = IOE.IOEither[error, *stock.Plasmid]
	PlasmidConverter = func(PlasmidIO) PlasmidResult
)

// ==================== GET PLASMID ====================

type getPlasmidParams struct {
	ctx     context.Context
	request *stock.StockId
	storage storage.StockStorage
}

type withValidatedPlasmidID struct {
	getPlasmidParams
	validatedID string
}

type withPlasmidDoc struct {
	withValidatedPlasmidID
	plasmid *stock.Plasmid
}

var (
	setValidatedPlasmidID = F.Curry2(
		func(stockID string, params getPlasmidParams) withValidatedPlasmidID {
			return withValidatedPlasmidID{
				getPlasmidParams: params,
				validatedID:      stockID,
			}
		},
	)

	setPlasmidDoc = F.Curry2(
		func(plasmid *stock.Plasmid, params withValidatedPlasmidID) withPlasmidDoc {
			return withPlasmidDoc{
				withValidatedPlasmidID: params,
				plasmid:                plasmid,
			}
		},
	)
)

func validatePlasmidIDRequest(params getPlasmidParams) IOE.IOEither[error, string] {
	return F.Pipe1(
		plasmidValidationToEither(params),
		IOE.FromEither[error, string],
	)
}

func plasmidValidationToEither(params getPlasmidParams) E.Either[error, string] {
	if err := validateRequest(params.request); err != nil {
		return E.Left[string](fmt.Errorf("invalid request parameters: %w", err))
	}
	return validateIDFormat(params.request.Id)
}

func retrievePlasmidFromStorage(params withValidatedPlasmidID) IOE.IOEither[error, *stock.Plasmid] {
	return F.Pipe1(
		params.storage.GetPlasmid(params.validatedID),
		IOE.MapLeft[*stock.Plasmid](enrichPlasmidError(params.validatedID)),
	)
}

func enrichPlasmidError(stockID string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("failed to retrieve plasmid %s: %w", stockID, err)
	}
}

func extractPlasmidResponse(params withPlasmidDoc) *stock.Plasmid {
	return params.plasmid
}

func toPlasmidResult(ctx context.Context) PlasmidConverter {
	return func(io PlasmidIO) PlasmidResult {
		either := io()
		return F.Pipe1(
			either,
			E.Fold(
				createPlasmidErrorResult(ctx),
				createPlasmidSuccessResult,
			),
		)
	}
}

func createPlasmidErrorResult(ctx context.Context) func(error) PlasmidResult {
	return func(err error) PlasmidResult {
		return T.MakeTuple2(
			&stock.Plasmid{},
			errorToGRPCStatus(ctx)(err),
		)
	}
}

func createPlasmidSuccessResult(plasmid *stock.Plasmid) PlasmidResult {
	return T.MakeTuple2[*stock.Plasmid, error](plasmid, nil)
}

func getPlasmid(params getPlasmidParams) PlasmidResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedPlasmidID, validatePlasmidIDRequest),
		IOE.Bind(setPlasmidDoc, retrievePlasmidFromStorage),
		IOE.Map[error](extractPlasmidResponse),
		toPlasmidResult(params.ctx),
	)
	return result
}

// ==================== CREATE PLASMID ====================

type createPlasmidParams struct {
	ctx     context.Context
	request *stock.NewPlasmid
	storage storage.StockStorage
	config  *ServiceConfig
}

type withValidatedNewPlasmid struct {
	createPlasmidParams
	validated *stock.NewPlasmid
}

type withCreatedPlasmid struct {
	withValidatedNewPlasmid
	created *stock.Plasmid
}

var (
	setValidatedNewPlasmid = F.Curry2(
		func(req *stock.NewPlasmid, params createPlasmidParams) withValidatedNewPlasmid {
			return withValidatedNewPlasmid{
				createPlasmidParams: params,
				validated:           req,
			}
		},
	)

	setCreatedPlasmid = F.Curry2(
		func(plasmid *stock.Plasmid, params withValidatedNewPlasmid) withCreatedPlasmid {
			return withCreatedPlasmid{
				withValidatedNewPlasmid: params,
				created:                 plasmid,
			}
		},
	)
)

func validateNewPlasmidRequest(params createPlasmidParams) IOE.IOEither[error, *stock.NewPlasmid] {
	return func() E.Either[error, *stock.NewPlasmid] {
		if params.request == nil {
			return E.Left[*stock.NewPlasmid](fmt.Errorf("request cannot be nil"))
		}

		// Validate using protobuf validation rules
		if err := params.request.Validate(); err != nil {
			return E.Left[*stock.NewPlasmid](fmt.Errorf("validation failed: %w", err))
		}

		// Apply default ontology term if not provided
		if params.request.Data != nil && params.request.Data.Attributes != nil {
			attrs := params.request.Data.Attributes
			if attrs.DictyPlasmidProperty == "" {
				attrs.DictyPlasmidProperty = params.config.PlasmidTerm
			}
		}

		return E.Right[error](params.request)
	}
}

func createPlasmidInStorage(params withValidatedNewPlasmid) IOE.IOEither[error, *stock.Plasmid] {
	return params.storage.CreatePlasmid(params.validated)
}

func extractCreatedPlasmidResponse(params withCreatedPlasmid) *stock.Plasmid {
	return params.created
}

func createPlasmid(params createPlasmidParams) PlasmidResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedNewPlasmid, validateNewPlasmidRequest),
		IOE.Bind(setCreatedPlasmid, createPlasmidInStorage),
		IOE.Map[error](extractCreatedPlasmidResponse),
		toPlasmidResult(params.ctx),
	)
	return result
}

// ==================== UPDATE PLASMID ====================

type updatePlasmidParams struct {
	ctx     context.Context
	request *stock.PlasmidUpdate
	storage storage.StockStorage
}

type withValidatedPlasmidUpdate struct {
	updatePlasmidParams
	validated *stock.PlasmidUpdate
}

type withUpdatedPlasmid struct {
	withValidatedPlasmidUpdate
	updated *stock.Plasmid
}

var (
	setValidatedPlasmidUpdate = F.Curry2(
		func(req *stock.PlasmidUpdate, params updatePlasmidParams) withValidatedPlasmidUpdate {
			return withValidatedPlasmidUpdate{
				updatePlasmidParams: params,
				validated:           req,
			}
		},
	)

	setUpdatedPlasmid = F.Curry2(
		func(plasmid *stock.Plasmid, params withValidatedPlasmidUpdate) withUpdatedPlasmid {
			return withUpdatedPlasmid{
				withValidatedPlasmidUpdate: params,
				updated:                    plasmid,
			}
		},
	)
)

func validatePlasmidUpdateRequest(
	params updatePlasmidParams,
) IOE.IOEither[error, *stock.PlasmidUpdate] {
	return func() E.Either[error, *stock.PlasmidUpdate] {
		if params.request == nil {
			return E.Left[*stock.PlasmidUpdate](fmt.Errorf("request cannot be nil"))
		}
		return E.Right[error](params.request)
	}
}

func updatePlasmidInStorage(params withValidatedPlasmidUpdate) IOE.IOEither[error, *stock.Plasmid] {
	return params.storage.UpdatePlasmid(params.validated)
}

func extractUpdatedPlasmidResponse(params withUpdatedPlasmid) *stock.Plasmid {
	return params.updated
}

func updatePlasmid(params updatePlasmidParams) PlasmidResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedPlasmidUpdate, validatePlasmidUpdateRequest),
		IOE.Bind(setUpdatedPlasmid, updatePlasmidInStorage),
		IOE.Map[error](extractUpdatedPlasmidResponse),
		toPlasmidResult(params.ctx),
	)
	return result
}

// ==================== LOAD PLASMID ====================

type loadPlasmidParams struct {
	ctx     context.Context
	request *stock.ExistingPlasmid
	storage storage.StockStorage
}

type withValidatedExistingPlasmid struct {
	loadPlasmidParams
	validated *stock.ExistingPlasmid
	stockID   string
}

type withLoadedPlasmid struct {
	withValidatedExistingPlasmid
	loaded *stock.Plasmid
}

var (
	setValidatedExistingPlasmid = F.Curry2(
		func(data T.Tuple2[*stock.ExistingPlasmid, string], params loadPlasmidParams) withValidatedExistingPlasmid {
			return withValidatedExistingPlasmid{
				loadPlasmidParams: params,
				validated:         data.F1,
				stockID:           data.F2,
			}
		},
	)

	setLoadedPlasmid = F.Curry2(
		func(plasmid *stock.Plasmid, params withValidatedExistingPlasmid) withLoadedPlasmid {
			return withLoadedPlasmid{
				withValidatedExistingPlasmid: params,
				loaded:                       plasmid,
			}
		},
	)
)

func validateExistingPlasmidRequest(
	params loadPlasmidParams,
) IOE.IOEither[error, T.Tuple2[*stock.ExistingPlasmid, string]] {
	return func() E.Either[error, T.Tuple2[*stock.ExistingPlasmid, string]] {
		if params.request == nil {
			return E.Left[T.Tuple2[*stock.ExistingPlasmid, string]](
				fmt.Errorf("request cannot be nil"),
			)
		}
		if params.request.Data == nil || params.request.Data.Id == "" {
			return E.Left[T.Tuple2[*stock.ExistingPlasmid, string]](
				fmt.Errorf("stock ID cannot be empty"),
			)
		}
		return E.Right[error](T.MakeTuple2(params.request, params.request.Data.Id))
	}
}

func loadPlasmidIntoStorage(
	params withValidatedExistingPlasmid,
) IOE.IOEither[error, *stock.Plasmid] {
	return params.storage.LoadPlasmid(params.stockID, params.validated)
}

func extractLoadedPlasmidResponse(params withLoadedPlasmid) *stock.Plasmid {
	return params.loaded
}

func loadPlasmid(params loadPlasmidParams) PlasmidResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedExistingPlasmid, validateExistingPlasmidRequest),
		IOE.Bind(setLoadedPlasmid, loadPlasmidIntoStorage),
		IOE.Map[error](extractLoadedPlasmidResponse),
		toPlasmidResult(params.ctx),
	)
	return result
}

// ==================== LIST PLASMIDS ====================

type listPlasmidsParams struct {
	ctx     context.Context
	params  *stock.StockParameters
	storage storage.StockStorage
}

type (
	PlasmidCollectionResult    = T.Tuple2[*stock.PlasmidCollection, error]
	PlasmidCollectionEither    = E.Either[error, *stock.PlasmidCollection]
	PlasmidCollectionIO        = IOE.IOEither[error, *stock.PlasmidCollection]
	PlasmidCollectionConverter = func(PlasmidCollectionIO) PlasmidCollectionResult
)

func retrievePlasmidCollection(
	params listPlasmidsParams,
) IOE.IOEither[error, *stock.PlasmidCollection] {
	return params.storage.ListPlasmids(params.params)
}

func toPlasmidCollectionResult(ctx context.Context) PlasmidCollectionConverter {
	return func(io PlasmidCollectionIO) PlasmidCollectionResult {
		either := io()
		return F.Pipe1(
			either,
			E.Fold(
				func(err error) PlasmidCollectionResult {
					return T.MakeTuple2(
						&stock.PlasmidCollection{},
						errorToGRPCStatus(ctx)(err),
					)
				},
				func(collection *stock.PlasmidCollection) PlasmidCollectionResult {
					return T.MakeTuple2[*stock.PlasmidCollection, error](collection, nil)
				},
			),
		)
	}
}

func listPlasmids(params listPlasmidsParams) PlasmidCollectionResult {
	result := F.Pipe2(
		IOE.Of[error](params),
		IOE.Chain(retrievePlasmidCollection),
		toPlasmidCollectionResult(params.ctx),
	)
	return result
}

// ==================== REMOVE STOCK ====================

type removeStockParams struct {
	ctx     context.Context
	request *stock.StockId
	storage storage.StockStorage
}

type (
	EmptyResult    = T.Tuple2[*emptypb.Empty, error]
	EmptyEither    = E.Either[error, *emptypb.Empty]
	EmptyIO        = IOE.IOEither[error, *emptypb.Empty]
	EmptyConverter = func(EmptyIO) EmptyResult
)

func validateRemoveStockRequest(params removeStockParams) IOE.IOEither[error, string] {
	return func() E.Either[error, string] {
		if err := validateRequest(params.request); err != nil {
			return E.Left[string](fmt.Errorf("invalid request: %w", err))
		}
		return validateIDFormat(params.request.Id)
	}
}

func removeStockFromStorage(
	stockID string,
	params removeStockParams,
) IOE.IOEither[error, *emptypb.Empty] {
	return F.Pipe1(
		params.storage.RemoveStock(stockID),
		IOE.Map[error](func(struct{}) *emptypb.Empty {
			return &emptypb.Empty{}
		}),
	)
}

func toEmptyResult(ctx context.Context) EmptyConverter {
	return func(io EmptyIO) EmptyResult {
		either := io()
		return F.Pipe1(
			either,
			E.Fold(
				func(err error) EmptyResult {
					return T.MakeTuple2(
						&emptypb.Empty{},
						errorToGRPCStatus(ctx)(err),
					)
				},
				func(empty *emptypb.Empty) EmptyResult {
					return T.MakeTuple2[*emptypb.Empty, error](empty, nil)
				},
			),
		)
	}
}

func removeStock(params removeStockParams) EmptyResult {
	result := F.Pipe2(
		validateRemoveStockRequest(params),
		IOE.Chain(func(stockID string) IOE.IOEither[error, *emptypb.Empty] {
			return removeStockFromStorage(stockID, params)
		}),
		toEmptyResult(params.ctx),
	)
	return result
}
