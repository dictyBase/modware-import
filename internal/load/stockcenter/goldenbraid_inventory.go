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

// listPlasmids creates the API call function
func listPlasmids(name string) func() (*stock.PlasmidCollection, error) {
	return func() (*stock.PlasmidCollection, error) {
		return regsc.GetStockAPIClient().
			ListPlasmids(context.Background(),
				&stock.StockParameters{
					Filter: fmt.Sprintf(
						"name==%s",
						name,
					),
					Limit: plasmidSearchLimit,
				})
	}
}

var plasmidListError = F.Curry2(func(name string, err error) error {
	return fmt.Errorf(
		"error listing plasmids for name %s: %w",
		name,
		err,
	)
})

var toPlasmidNameContext = F.Curry2(
	func(name string, resp *stock.PlasmidCollection) PlasmidNameContext {
		return PlasmidNameContext{
			Name: name,
			Resp: resp,
		}
	},
)

// resolvePlasmidID finds the plasmid ID for a given name
func resolvePlasmidID(name string) IOE.IOEither[error, string] {
	return F.Pipe3(
		// Step 1: Call API
		IOE.TryCatchError(listPlasmids(name)),
		// Step 2: Add error context
		IOE.MapLeft[*stock.PlasmidCollection](plasmidListError(name)),
		// Step 3: Map to Context
		IOE.Map[error](toPlasmidNameContext(name)),
		// Step 4: Validate response
		IOE.ChainEitherK(validatePlasmidResponse),
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

// syncInventory handles the check-delete-create logic for inventory
func syncInventory(plasmidID, location string) IOE.IOEither[error, string] {
	client := regsc.GetAnnotationAPIClient()

	return F.Pipe2(
		// Step 1: Get existing inventory
		getInventoryIO(plasmidID, client),
		// Step 2: Delete if exists
		IOE.Chain(
			func(gc *pb.TaggedAnnotationGroupCollection) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
				return deleteInventoryIfExists(client, gc)
			},
		),
		// Step 3: Create new inventory
		IOE.Chain(
			func(*pb.TaggedAnnotationGroupCollection) IOE.IOEither[error, string] {
				return createInventoryIO(plasmidID, location, client)
			},
		),
	)
}

// getInventoryIO wraps getInventory in IOEither
func getInventoryIO(
	plasmidID string,
	client pb.TaggedAnnotationServiceClient,
) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.TaggedAnnotationGroupCollection, error) {
			return getInventory(plasmidID, client, regsc.PlasmidInvOntO)
		}),
		IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](func(err error) error {
			return fmt.Errorf("error checking inventory: %w", err)
		}),
	)
}

// deleteInventoryIfExists conditionally deletes
func deleteInventoryIfExists(
	client pb.TaggedAnnotationServiceClient,
	gc *pb.TaggedAnnotationGroupCollection,
) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
	if len(gc.Data) == 0 {
		// Nothing to delete - return success
		return IOE.Of[error](gc)
	}

	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.TaggedAnnotationGroupCollection, error) {
			if err := delAnnotationGroup(client, gc); err != nil {
				return nil, err
			}
			return gc, nil
		}),
		IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](func(err error) error {
			return fmt.Errorf("error deleting existing inventory: %w", err)
		}),
	)
}

func createInventoryIO(
	plasmidID, location string,
	client pb.TaggedAnnotationServiceClient,
) IOE.IOEither[error, string] {
	attributes := map[string]string{
		regsc.InvLocationTag:    location,
		regsc.InvStorageDateTag: time.Now().Format(time.RFC3339Nano),
	}

	return F.Pipe2(
		// Step 1: Create all annotations (using Traverse pattern)
		createAnnotationsForAttributes(plasmidID, attributes, client),
		// Step 2: Create annotation group
		IOE.Chain(func(annoIDs []string) IOE.IOEither[error, []string] {
			return createAnnotationGroupIO(client, annoIDs)
		}),
		// Step 3: Mark existence
		IOE.Chain(func(_ []string) IOE.IOEither[error, string] {
			return markInventoryExistence(plasmidID, client)
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

// createAnnotationGroupIO creates the annotation group
func createAnnotationGroupIO(
	client pb.TaggedAnnotationServiceClient,
	annoIDs []string,
) IOE.IOEither[error, []string] {
	return F.Pipe1(
		IOE.TryCatchError(func() ([]string, error) {
			_, err := client.CreateAnnotationGroup(
				context.Background(),
				&pb.AnnotationIdList{Ids: annoIDs},
			)
			return annoIDs, err
		}),
		IOE.MapLeft[[]string](func(err error) error {
			return fmt.Errorf("failed to create annotation group: %w", err)
		}),
	)
}

// markInventoryExistence marks the inventory as existing
func markInventoryExistence(
	plasmidID string,
	client pb.TaggedAnnotationServiceClient,
) IOE.IOEither[error, string] {
	return F.Pipe1(
		IOE.TryCatchError(func() (string, error) {
			err := createAnno(&createAnnoArgs{
				client:   client,
				tag:      regsc.PlasmidInvTag,
				id:       plasmidID,
				ontology: regsc.PlasmidInvOntO,
				value:    regsc.InvExistValue,
			})
			return plasmidID, err
		}),
		IOE.MapLeft[string](func(err error) error {
			return fmt.Errorf("failed to mark inventory existence: %w", err)
		}),
	)
}

func processRow(record InventoryRecord) IOE.IOEither[error, string] {
	return F.Pipe2(
		resolvePlasmidID(record.PlasmidName),
		IOE.Chain(func(id string) IOE.IOEither[error, string] {
			return syncInventory(id, record.Location)
		}),
		IOE.Map[error](func(id string) string { return id }),
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
