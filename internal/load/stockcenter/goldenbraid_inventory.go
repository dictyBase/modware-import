package stockcenter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	S "github.com/IBM/fp-go/semigroup"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

const (
	plasmidSearchLimit = 2
)

type InventoryRecord struct {
	PlasmidName string
	Location    string
}

type InventoryProcessingSummary struct {
	SuccessCount int
	ErrorCount   int
	Errors       []string
}

type PlasmidNameContext struct {
	Name string
	Resp *stock.PlasmidCollection
}

func InventorySummarySemigroup() S.Semigroup[InventoryProcessingSummary] {
	return S.MakeSemigroup(
		func(a, b InventoryProcessingSummary) InventoryProcessingSummary {
			return InventoryProcessingSummary{
				SuccessCount: a.SuccessCount + b.SuccessCount,
				ErrorCount:   a.ErrorCount + b.ErrorCount,
				Errors:       append(a.Errors, b.Errors...),
			}
		},
	)
}

// listPlasmidsIOE calls API and maps result to PlasmidNameContext
func listPlasmidsIOE(name string) IOE.IOEither[error, PlasmidNameContext] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
			return regsc.GetStockAPIClient().
				ListPlasmids(context.Background(),
					&stock.StockParameters{
						Filter: fmt.Sprintf("name==%s", name),
						Limit:  plasmidSearchLimit,
					})
		}),
		IOE.MapLeft[*stock.PlasmidCollection](func(err error) error {
			return fmt.Errorf(
				"error listing plasmids for name %s: %w",
				name,
				err,
			)
		}),
		IOE.Map[error](func(resp *stock.PlasmidCollection) PlasmidNameContext {
			return PlasmidNameContext{
				Name: name,
				Resp: resp,
			}
		}),
	)
}

func hasSinglePlasmid(ctx PlasmidNameContext) bool {
	return len(ctx.Resp.Data) == 1
}

func plasmidCountError(ctx PlasmidNameContext) error {
	return fmt.Errorf(
		"expected 1 plasmid for %s, found %d",
		ctx.Name,
		len(ctx.Resp.Data),
	)
}

func firstPlasmidID(ctx PlasmidNameContext) string {
	return ctx.Resp.Data[0].Id
}

// validatePlasmidResponse uses Either for pure validation
func validatePlasmidResponse(
	ctx PlasmidNameContext,
) E.Either[error, string] {
	return F.Pipe2(
		ctx,
		E.FromPredicate(
			hasSinglePlasmid,
			plasmidCountError,
		),
		E.Map[error](firstPlasmidID),
	)
}

// ProcessRowCtx is the initial context for processRow pipeline
type ProcessRowCtx struct {
	PlasmidName string
	Location    string
}

// WithPlasmidContext adds the plasmid name context from API
type WithPlasmidContext struct {
	ProcessRowCtx
	NameContext PlasmidNameContext
}

// WithPlasmidID adds the validated plasmid ID
type WithPlasmidID struct {
	WithPlasmidContext
	PlasmidID string
}

// WithInventory adds the existing inventory collection
type WithInventory struct {
	WithPlasmidID
	Inventory *pb.TaggedAnnotationGroupCollection
}

// WithDeletedInventory represents state after deletion check
type WithDeletedInventory struct {
	WithInventory
}

// WithAnnotationIDs adds the created annotation IDs
type WithAnnotationIDs struct {
	WithDeletedInventory
	AnnotationIDs []string
}

// WithAnnotationGroup represents state after annotation group creation
type WithAnnotationGroup struct {
	WithAnnotationIDs
}

// SetPlasmidContext sets the plasmid name context
var SetPlasmidContext = F.Curry2(func(ctx PlasmidNameContext, s ProcessRowCtx) WithPlasmidContext {
	return WithPlasmidContext{ProcessRowCtx: s, NameContext: ctx}
})

// SetPlasmidID sets the validated plasmid ID
var SetPlasmidID = F.Curry2(func(id string, s WithPlasmidContext) WithPlasmidID {
	return WithPlasmidID{WithPlasmidContext: s, PlasmidID: id}
})

// SetInventory sets the existing inventory collection
var SetInventory = F.Curry2(
	func(inv *pb.TaggedAnnotationGroupCollection, s WithPlasmidID) WithInventory {
		return WithInventory{WithPlasmidID: s, Inventory: inv}
	},
)

// SetDeletedInventory marks inventory as deleted (or confirmed absent)
var SetDeletedInventory = F.Curry2(
	func(_ *pb.TaggedAnnotationGroupCollection, s WithInventory) WithDeletedInventory {
		return WithDeletedInventory{WithInventory: s}
	},
)

// SetAnnotationIDs sets the created annotation IDs
var SetAnnotationIDs = F.Curry2(func(ids []string, s WithDeletedInventory) WithAnnotationIDs {
	return WithAnnotationIDs{WithDeletedInventory: s, AnnotationIDs: ids}
})

// SetAnnotationGroup marks annotation group as created
var SetAnnotationGroup = F.Curry2(func(_ []string, s WithAnnotationIDs) WithAnnotationGroup {
	return WithAnnotationGroup{WithAnnotationIDs: s}
})

// SetFinalID sets the final plasmid ID result
var SetFinalID = F.Curry2(func(id string, _ WithAnnotationGroup) string {
	return id
})

// getPlasmidContext retrieves plasmid name context from API
func getPlasmidContext(ctx ProcessRowCtx) IOE.IOEither[error, PlasmidNameContext] {
	return listPlasmidsIOE(ctx.PlasmidName)
}

// validateAndExtractID validates plasmid response and extracts ID
func validateAndExtractID(ctx WithPlasmidContext) IOE.IOEither[error, string] {
	return IOE.FromEither(validatePlasmidResponse(ctx.NameContext))
}

// createInventoryAnnotations creates annotation IDs for inventory attributes from context
func createInventoryAnnotations(ctx WithDeletedInventory) IOE.IOEither[error, []string] {
	client := regsc.GetAnnotationAPIClient()
	attributes := map[string]string{
		regsc.InvLocationTag:    ctx.Location,
		regsc.InvStorageDateTag: time.Now().Format(time.RFC3339Nano),
	}

	return createAnnotationsForAttributes(ctx.PlasmidID, attributes, client)
}

// getInventoryIO retrieves existing inventory for plasmid from context
func getInventoryIO(ctx WithPlasmidID) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
	client := regsc.GetAnnotationAPIClient()

	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.TaggedAnnotationGroupCollection, error) {
			return client.ListAnnotationGroups(
				context.Background(),
				&pb.ListGroupParameters{
					Filter: fmt.Sprintf(
						"entry_id==%s;tag==%s;ontology==%s",
						ctx.PlasmidID,
						regsc.InvLocationTag,
						regsc.PlasmidInvOntO,
					),
				},
			)
		}),
		IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](func(err error) error {
			return fmt.Errorf("error checking inventory: %w", err)
		}),
	)
}

// deleteInventoryIfExists conditionally deletes existing inventory from context
func deleteInventoryIfExists(
	ctx WithInventory,
) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
	if len(ctx.Inventory.Data) == 0 {
		// Nothing to delete - return success
		return IOE.Of[error](ctx.Inventory)
	}

	client := regsc.GetAnnotationAPIClient()
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.TaggedAnnotationGroupCollection, error) {
			if err := delAnnotationGroup(client, ctx.Inventory); err != nil {
				return nil, err
			}
			return ctx.Inventory, nil
		}),
		IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](func(err error) error {
			return fmt.Errorf("error deleting existing inventory: %w", err)
		}),
	)
}

// createAnnotationsForAttributes creates annotations for all attributes
// This replaces the imperative for loop with functional map
func createAnnotationsForAttributes(
	plasmidID string,
	attributes map[string]string,
	client pb.TaggedAnnotationServiceClient,
) IOE.IOEither[error, []string] {
	// Convert map to slice of key-value pairs
	type attrPair struct{ tag, value string }
	pairs := make([]attrPair, 0, len(attributes))
	for tag, value := range attributes {
		pairs = append(pairs, attrPair{tag, value})
	}

	// Use Traverse pattern: []A → IOEither[error, []B]
	return IOE.TraverseArraySeq(
		func(pair attrPair) IOE.IOEither[error, string] {
			return F.Pipe2(
				IOE.TryCatchError(func() (*pb.TaggedAnnotation, error) {
					return createAnnoWithRank(&createAnnoArgs{
						ontology: regsc.PlasmidInvOntO,
						client:   client,
						id:       plasmidID,
						value:    pair.value,
						tag:      pair.tag,
						rank:     0,
					})
				}),
				IOE.MapLeft[*pb.TaggedAnnotation](func(err error) error {
					return fmt.Errorf(
						"failed to create annotation for %s: %w",
						pair.tag,
						err,
					)
				}),
				IOE.Map[error](func(anno *pb.TaggedAnnotation) string {
					if anno == nil || anno.Data == nil {
						return ""
					}
					return anno.Data.Id
				}),
			)
		},
	)(
		pairs,
	)
}

// createAnnotationGroupIO creates annotation group from context
func createAnnotationGroupIO(ctx WithAnnotationIDs) IOE.IOEither[error, []string] {
	client := regsc.GetAnnotationAPIClient()

	return F.Pipe1(
		IOE.TryCatchError(func() ([]string, error) {
			_, err := client.CreateAnnotationGroup(
				context.Background(),
				&pb.AnnotationIdList{Ids: ctx.AnnotationIDs},
			)
			return ctx.AnnotationIDs, err
		}),
		IOE.MapLeft[[]string](func(err error) error {
			return fmt.Errorf("failed to create annotation group: %w", err)
		}),
	)
}

// markInventoryExistence marks inventory as existing from context
func markInventoryExistence(ctx WithAnnotationGroup) IOE.IOEither[error, string] {
	client := regsc.GetAnnotationAPIClient()

	return F.Pipe1(
		IOE.TryCatchError(func() (string, error) {
			err := createAnno(&createAnnoArgs{
				client:   client,
				tag:      regsc.PlasmidInvTag,
				id:       ctx.PlasmidID,
				ontology: regsc.PlasmidInvOntO,
				value:    regsc.InvExistValue,
			})
			return ctx.PlasmidID, err
		}),
		IOE.MapLeft[string](func(err error) error {
			return fmt.Errorf("failed to mark inventory existence: %w", err)
		}),
	)
}

// processRow processes a single inventory record using Do/Bind pattern
func processRow(record InventoryRecord) IOE.IOEither[error, string] {
	return F.Pipe8(
		IOE.Do[error](ProcessRowCtx(record)),
		IOE.Bind(SetPlasmidContext, getPlasmidContext),
		IOE.Bind(SetPlasmidID, validateAndExtractID),
		IOE.Bind(SetInventory, getInventoryIO),
		IOE.Bind(SetDeletedInventory, deleteInventoryIfExists),
		IOE.Bind(SetAnnotationIDs, createInventoryAnnotations),
		IOE.Bind(SetAnnotationGroup, createAnnotationGroupIO),
		IOE.Bind(SetFinalID, markInventoryExistence),
		IOE.Map[error](F.Identity[string]),
	)
}

func ProcessInventory(
	config InventoryLoaderConfig,
) IOE.IOEither[error, InventoryProcessingSummary] {
	return F.Pipe1(
		// Step 1: Query rows (returns *sql.Rows)
		queryInventoryRows(config.DB),
		// Step 2: Fold over rows using Semigroup (streaming, no materialization)
		IOE.Chain(
			func(rows *sql.Rows) IOE.IOEither[error, InventoryProcessingSummary] {
				return foldRowsWithSemigroup(rows, config.Logger)
			},
		),
	)
}

// queryInventoryRows queries the database and returns sql.Rows
func queryInventoryRows(db *sql.DB) IOE.IOEither[error, *sql.Rows] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*sql.Rows, error) {
			return db.Query("SELECT Name, Location FROM inventory")
		}),
		IOE.MapLeft[*sql.Rows](func(err error) error {
			return fmt.Errorf("error querying inventory: %w", err)
		}),
	)
}

// foldRowsWithSemigroup streams over rows one at a time, folding with Semigroup
// This processes rows lazily without loading all into memory
func foldRowsWithSemigroup(
	rows *sql.Rows,
	logger *slog.Logger,
) IOE.IOEither[error, InventoryProcessingSummary] {
	return IOE.TryCatchError(func() (InventoryProcessingSummary, error) {
		defer rows.Close()

		// Initial accumulator (identity for Semigroup)
		summary := InventoryProcessingSummary{}
		semigroup := InventorySummarySemigroup()

		// Stream over rows, folding one at a time
		for rows.Next() {
			var name, location string

			// Scan row
			if err := rows.Scan(&name, &location); err != nil {
				// Create error summary for scan failure
				scanError := InventoryProcessingSummary{
					ErrorCount: 1,
					Errors:     []string{fmt.Sprintf("scan error: %v", err)},
				}
				// Accumulate using Semigroup
				summary = semigroup.Concat(summary, scanError)
				logger.Error("failed to scan row", "error", err)
				continue
			}

			record := InventoryRecord{
				PlasmidName: name,
				Location:    location,
			}

			// Process single row to get its summary
			rowSummary := processRowToSummary(record, logger)

			// Accumulate using Semigroup (fold step)
			summary = semigroup.Concat(summary, rowSummary)
		}

		// Check for iteration errors
		if err := rows.Err(); err != nil {
			return summary, fmt.Errorf("error iterating rows: %w", err)
		}

		return summary, nil
	})
}

// processRowToSummary processes one row and returns a summary
// Pure function that doesn't perform I/O (I/O happens in processRow)
func processRowToSummary(
	record InventoryRecord,
	logger *slog.Logger,
) InventoryProcessingSummary {
	// Execute the IOEither to get Either result
	result := processRow(record)()

	// Fold the Either into a summary
	return E.Fold(
		// Error case
		func(err error) InventoryProcessingSummary {
			logger.Error(
				"failed to process inventory",
				"plasmid", record.PlasmidName,
				"error", err,
			)
			return InventoryProcessingSummary{
				ErrorCount: 1,
				Errors: []string{
					fmt.Sprintf("%s: %v", record.PlasmidName, err),
				},
			}
		},
		// Success case
		func(_ string) InventoryProcessingSummary {
			logger.Info(
				"processed inventory",
				"plasmid", record.PlasmidName,
				"location", record.Location,
			)
			return InventoryProcessingSummary{
				SuccessCount: 1,
			}
		},
	)(result)
}
