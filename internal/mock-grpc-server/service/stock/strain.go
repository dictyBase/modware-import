package stock

import (
	"context"
	"fmt"
	"strings"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	T "github.com/IBM/fp-go/tuple"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/mock-grpc-server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ==================== HELPER FUNCTIONS ====================

// isValidEmail checks if a string is a valid email format
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ==================== TYPE DEFINITIONS ====================

type (
	StrainResult    = T.Tuple2[*stock.Strain, error]
	StrainEither    = E.Either[error, *stock.Strain]
	StrainIO        = IOE.IOEither[error, *stock.Strain]
	StrainConverter = func(StrainIO) StrainResult
)

// ==================== GET STRAIN ====================

type getStrainParams struct {
	ctx     context.Context
	request *stock.StockId
	storage storage.StockStorage
}

type withValidatedStrainID struct {
	getStrainParams
	validatedID string
}

type withStrainDoc struct {
	withValidatedStrainID
	strain *stock.Strain
}

var (
	setValidatedStrainID = F.Curry2(
		func(stockID string, params getStrainParams) withValidatedStrainID {
			return withValidatedStrainID{
				getStrainParams: params,
				validatedID:     stockID,
			}
		},
	)

	setStrainDoc = F.Curry2(
		func(strain *stock.Strain, params withValidatedStrainID) withStrainDoc {
			return withStrainDoc{
				withValidatedStrainID: params,
				strain:                strain,
			}
		},
	)
)

func validateStrainIDRequest(params getStrainParams) IOE.IOEither[error, string] {
	return F.Pipe1(
		validationToEither(params),
		IOE.FromEither[error, string],
	)
}

func validationToEither(params getStrainParams) E.Either[error, string] {
	if err := validateRequest(params.request); err != nil {
		return E.Left[string](fmt.Errorf("invalid request parameters: %w", err))
	}
	return F.Pipe1(
		extractRequestID(params),
		validateIDFormat,
	)
}

func validateRequest(req *stock.StockId) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}
	return nil
}

func extractRequestID(params getStrainParams) string {
	return params.request.Id
}

func validateIDFormat(stockID string) E.Either[error, string] {
	if stockID == "" {
		return E.Left[string](fmt.Errorf("stock ID cannot be empty"))
	}
	return E.Right[error](stockID)
}

func retrieveStrainFromStorage(params withValidatedStrainID) IOE.IOEither[error, *stock.Strain] {
	return F.Pipe1(
		params.storage.GetStrain(params.validatedID),
		IOE.MapLeft[*stock.Strain](enrichStrainError(params.validatedID)),
	)
}

func enrichStrainError(stockID string) func(error) error {
	return func(err error) error {
		return fmt.Errorf("failed to retrieve strain %s: %w", stockID, err)
	}
}

func extractStrainResponse(params withStrainDoc) *stock.Strain {
	return params.strain
}

func toStrainResult(ctx context.Context) StrainConverter {
	return func(io StrainIO) StrainResult {
		either := io()
		return F.Pipe1(
			either,
			E.Fold(
				createErrorResult(ctx),
				createSuccessResult,
			),
		)
	}
}

func createErrorResult(ctx context.Context) func(error) StrainResult {
	return func(err error) StrainResult {
		return T.MakeTuple2(
			&stock.Strain{},
			errorToGRPCStatus(ctx)(err),
		)
	}
}

func createSuccessResult(strain *stock.Strain) StrainResult {
	return T.MakeTuple2[*stock.Strain, error](strain, nil)
}

func errorToGRPCStatus(_ context.Context) func(error) error {
	return func(err error) error {
		if err == nil {
			return nil
		}

		errMsg := err.Error()

		switch {
		case contains(errMsg, "not found"):
			return status.Error(codes.NotFound, errMsg)
		case contains(errMsg, "invalid"):
			return status.Error(codes.InvalidArgument, errMsg)
		case contains(errMsg, "already exists"):
			return status.Error(codes.AlreadyExists, errMsg)
		default:
			return status.Error(codes.Internal, errMsg)
		}
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || findSubstring(str, substr))
}

func findSubstring(str, substr string) bool {
	for idx := 0; idx <= len(str)-len(substr); idx++ {
		if str[idx:idx+len(substr)] == substr {
			return true
		}
	}
	return false
}

func getStrain(params getStrainParams) StrainResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedStrainID, validateStrainIDRequest),
		IOE.Bind(setStrainDoc, retrieveStrainFromStorage),
		IOE.Map[error](extractStrainResponse),
		toStrainResult(params.ctx),
	)
	return result
}

// ==================== CREATE STRAIN ====================

type createStrainParams struct {
	ctx     context.Context
	request *stock.NewStrain
	storage storage.StockStorage
	config  *ServiceConfig
}

type withValidatedNewStrain struct {
	createStrainParams
	validated *stock.NewStrain
}

type withCreatedStrain struct {
	withValidatedNewStrain
	created *stock.Strain
}

var (
	setValidatedNewStrain = F.Curry2(
		func(req *stock.NewStrain, params createStrainParams) withValidatedNewStrain {
			return withValidatedNewStrain{
				createStrainParams: params,
				validated:          req,
			}
		},
	)

	setCreatedStrain = F.Curry2(
		func(strain *stock.Strain, params withValidatedNewStrain) withCreatedStrain {
			return withCreatedStrain{
				withValidatedNewStrain: params,
				created:                strain,
			}
		},
	)
)

func validateNewStrainRequest(params createStrainParams) IOE.IOEither[error, *stock.NewStrain] {
	return func() E.Either[error, *stock.NewStrain] {
		if params.request == nil {
			return E.Left[*stock.NewStrain](fmt.Errorf("request cannot be nil"))
		}

		// Validate using protobuf validation rules
		if err := params.request.Validate(); err != nil {
			return E.Left[*stock.NewStrain](fmt.Errorf("validation failed: %w", err))
		}

		// Additional email validation
		if params.request.Data != nil && params.request.Data.Attributes != nil {
			attrs := params.request.Data.Attributes
			if !isValidEmail(attrs.CreatedBy) {
				return E.Left[*stock.NewStrain](fmt.Errorf("invalid email format for created_by"))
			}
			if !isValidEmail(attrs.UpdatedBy) {
				return E.Left[*stock.NewStrain](fmt.Errorf("invalid email format for updated_by"))
			}

			// Apply default ontology term if not provided
			if attrs.DictyStrainProperty == "" {
				attrs.DictyStrainProperty = params.config.StrainTerm
			}
		}

		return E.Right[error](params.request)
	}
}

func createStrainInStorage(params withValidatedNewStrain) IOE.IOEither[error, *stock.Strain] {
	return params.storage.CreateStrain(params.validated)
}

func extractCreatedStrainResponse(params withCreatedStrain) *stock.Strain {
	return params.created
}

func createStrain(params createStrainParams) StrainResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedNewStrain, validateNewStrainRequest),
		IOE.Bind(setCreatedStrain, createStrainInStorage),
		IOE.Map[error](extractCreatedStrainResponse),
		toStrainResult(params.ctx),
	)
	return result
}

// ==================== UPDATE STRAIN ====================

type updateStrainParams struct {
	ctx     context.Context
	request *stock.StrainUpdate
	storage storage.StockStorage
}

type withValidatedStrainUpdate struct {
	updateStrainParams
	validated *stock.StrainUpdate
}

type withUpdatedStrain struct {
	withValidatedStrainUpdate
	updated *stock.Strain
}

var (
	setValidatedStrainUpdate = F.Curry2(
		func(req *stock.StrainUpdate, params updateStrainParams) withValidatedStrainUpdate {
			return withValidatedStrainUpdate{
				updateStrainParams: params,
				validated:          req,
			}
		},
	)

	setUpdatedStrain = F.Curry2(
		func(strain *stock.Strain, params withValidatedStrainUpdate) withUpdatedStrain {
			return withUpdatedStrain{
				withValidatedStrainUpdate: params,
				updated:                   strain,
			}
		},
	)
)

func validateStrainUpdateRequest(
	params updateStrainParams,
) IOE.IOEither[error, *stock.StrainUpdate] {
	return func() E.Either[error, *stock.StrainUpdate] {
		if params.request == nil {
			return E.Left[*stock.StrainUpdate](fmt.Errorf("request cannot be nil"))
		}
		return E.Right[error](params.request)
	}
}

func updateStrainInStorage(params withValidatedStrainUpdate) IOE.IOEither[error, *stock.Strain] {
	return params.storage.UpdateStrain(params.validated)
}

func extractUpdatedStrainResponse(params withUpdatedStrain) *stock.Strain {
	return params.updated
}

func updateStrain(params updateStrainParams) StrainResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedStrainUpdate, validateStrainUpdateRequest),
		IOE.Bind(setUpdatedStrain, updateStrainInStorage),
		IOE.Map[error](extractUpdatedStrainResponse),
		toStrainResult(params.ctx),
	)
	return result
}

// ==================== LOAD STRAIN ====================

type loadStrainParams struct {
	ctx     context.Context
	request *stock.ExistingStrain
	storage storage.StockStorage
}

type withValidatedExistingStrain struct {
	loadStrainParams
	validated *stock.ExistingStrain
	stockID   string
}

type withLoadedStrain struct {
	withValidatedExistingStrain
	loaded *stock.Strain
}

var (
	setValidatedExistingStrain = F.Curry2(
		func(data T.Tuple2[*stock.ExistingStrain, string], params loadStrainParams) withValidatedExistingStrain {
			return withValidatedExistingStrain{
				loadStrainParams: params,
				validated:        data.F1,
				stockID:          data.F2,
			}
		},
	)

	setLoadedStrain = F.Curry2(
		func(strain *stock.Strain, params withValidatedExistingStrain) withLoadedStrain {
			return withLoadedStrain{
				withValidatedExistingStrain: params,
				loaded:                      strain,
			}
		},
	)
)

func validateExistingStrainRequest(
	params loadStrainParams,
) IOE.IOEither[error, T.Tuple2[*stock.ExistingStrain, string]] {
	return func() E.Either[error, T.Tuple2[*stock.ExistingStrain, string]] {
		if params.request == nil {
			return E.Left[T.Tuple2[*stock.ExistingStrain, string]](
				fmt.Errorf("request cannot be nil"),
			)
		}
		if params.request.Data == nil || params.request.Data.Id == "" {
			return E.Left[T.Tuple2[*stock.ExistingStrain, string]](
				fmt.Errorf("stock ID cannot be empty"),
			)
		}
		return E.Right[error](T.MakeTuple2(params.request, params.request.Data.Id))
	}
}

func loadStrainIntoStorage(params withValidatedExistingStrain) IOE.IOEither[error, *stock.Strain] {
	return params.storage.LoadStrain(params.stockID, params.validated)
}

func extractLoadedStrainResponse(params withLoadedStrain) *stock.Strain {
	return params.loaded
}

func loadStrain(params loadStrainParams) StrainResult {
	result := F.Pipe4(
		IOE.Of[error](params),
		IOE.Bind(setValidatedExistingStrain, validateExistingStrainRequest),
		IOE.Bind(setLoadedStrain, loadStrainIntoStorage),
		IOE.Map[error](extractLoadedStrainResponse),
		toStrainResult(params.ctx),
	)
	return result
}

// ==================== LIST STRAINS ====================

type listStrainsParams struct {
	ctx     context.Context
	params  *stock.StockParameters
	storage storage.StockStorage
}

type (
	StrainCollectionResult    = T.Tuple2[*stock.StrainCollection, error]
	StrainCollectionEither    = E.Either[error, *stock.StrainCollection]
	StrainCollectionIO        = IOE.IOEither[error, *stock.StrainCollection]
	StrainCollectionConverter = func(StrainCollectionIO) StrainCollectionResult
)

func retrieveStrainCollection(
	params listStrainsParams,
) IOE.IOEither[error, *stock.StrainCollection] {
	return params.storage.ListStrains(params.params)
}

func toStrainCollectionResult(ctx context.Context) StrainCollectionConverter {
	return func(io StrainCollectionIO) StrainCollectionResult {
		either := io()
		return F.Pipe1(
			either,
			E.Fold(
				func(err error) StrainCollectionResult {
					return T.MakeTuple2(
						&stock.StrainCollection{},
						errorToGRPCStatus(ctx)(err),
					)
				},
				func(collection *stock.StrainCollection) StrainCollectionResult {
					return T.MakeTuple2[*stock.StrainCollection, error](collection, nil)
				},
			),
		)
	}
}

func listStrains(params listStrainsParams) StrainCollectionResult {
	result := F.Pipe2(
		IOE.Of[error](params),
		IOE.Chain(retrieveStrainCollection),
		toStrainCollectionResult(params.ctx),
	)
	return result
}

// ==================== LIST STRAINS BY IDS ====================

type listStrainsByIDsParams struct {
	ctx     context.Context
	request *stock.StockIdList
	storage storage.StockStorage
}

type (
	StrainListResult    = T.Tuple2[*stock.StrainList, error]
	StrainListEither    = E.Either[error, *stock.StrainList]
	StrainListIO        = IOE.IOEither[error, *stock.StrainList]
	StrainListConverter = func(StrainListIO) StrainListResult
)

func retrieveStrainsByIDs(params listStrainsByIDsParams) IOE.IOEither[error, *stock.StrainList] {
	return params.storage.ListStrainsByIDs(params.request.Id)
}

func toStrainListResult(ctx context.Context) StrainListConverter {
	return func(io StrainListIO) StrainListResult {
		either := io()
		return F.Pipe1(
			either,
			E.Fold(
				func(err error) StrainListResult {
					return T.MakeTuple2(
						&stock.StrainList{},
						errorToGRPCStatus(ctx)(err),
					)
				},
				func(list *stock.StrainList) StrainListResult {
					return T.MakeTuple2[*stock.StrainList, error](list, nil)
				},
			),
		)
	}
}

func listStrainsByIDs(params listStrainsByIDsParams) StrainListResult {
	result := F.Pipe2(
		IOE.Of[error](params),
		IOE.Chain(retrieveStrainsByIDs),
		toStrainListResult(params.ctx),
	)
	return result
}
