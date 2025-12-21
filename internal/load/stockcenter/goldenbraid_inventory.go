package stockcenter

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	R "github.com/IBM/fp-go/record"
	S "github.com/IBM/fp-go/semigroup"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/fputil"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	plasmidSearchLimit = 2
)

// ProcessRowCtx is the initial context for processRow pipeline
type ProcessRowCtx struct {
	PlasmidName string
	Location    string
}

type AnnotationContext struct {
	Client    pb.TaggedAnnotationServiceClient
	PlasmidID string
	Tag       string
	Value     string
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

// SetPlasmidContext sets the plasmid name context
var SetPlasmidContext = F.Curry2(
	func(ctx PlasmidNameContext, s ProcessRowCtx) WithPlasmidContext {
		return WithPlasmidContext{ProcessRowCtx: s, NameContext: ctx}
	},
)

// SetPlasmidID sets the validated plasmid ID
var SetPlasmidID = F.Curry2(
	func(id string, s WithPlasmidContext) WithPlasmidID {
		return WithPlasmidID{WithPlasmidContext: s, PlasmidID: id}
	},
)

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
var SetAnnotationIDs = F.Curry2(
	func(ids []string, s WithDeletedInventory) WithAnnotationIDs {
		return WithAnnotationIDs{WithDeletedInventory: s, AnnotationIDs: ids}
	},
)

// SetAnnotationGroup marks annotation group as created
var SetAnnotationGroup = F.Curry2(
	func(_ []string, s WithAnnotationIDs) WithAnnotationGroup {
		return WithAnnotationGroup{WithAnnotationIDs: s}
	},
)

// SetFinalID sets the final plasmid ID result
var SetFinalID = F.Curry2(func(id string, _ WithAnnotationGroup) string {
	return id
})

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
func listPlasmidsIOE(
	ctx ProcessRowCtx,
) IOE.IOEither[error, PlasmidNameContext] {
	return F.Pipe2(
		IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
			return regsc.GetStockAPIClient().
				ListPlasmids(context.Background(),
					&stock.StockParameters{
						Filter: fmt.Sprintf("name==%s", ctx.PlasmidName),
						Limit:  plasmidSearchLimit,
					})
		}),
		IOE.MapLeft[*stock.PlasmidCollection](func(err error) error {
			return fmt.Errorf(
				"error listing plasmids for name %s: %w",
				ctx.PlasmidName,
				err,
			)
		}),
		IOE.Map[error](func(resp *stock.PlasmidCollection) PlasmidNameContext {
			return PlasmidNameContext{
				Name: ctx.PlasmidName,
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

// validatePlasmidResponse validates plasmid response and extracts ID using pure Either logic
func validatePlasmidResponse(
	ctx WithPlasmidContext,
) IOE.IOEither[error, string] {
	return F.Pipe3(
		ctx.NameContext,
		E.FromPredicate(
			hasSinglePlasmid,
			plasmidCountError,
		),
		E.Map[error](firstPlasmidID),
		IOE.FromEither,
	)
}

func createAnnotationIOE(
	ctx AnnotationContext,
) IOE.IOEither[error, *pb.TaggedAnnotation] {
	input := &pb.NewTaggedAnnotation_Data{
		Attributes: &pb.NewTaggedAnnotationAttributes{
			Value:     ctx.Value,
			CreatedBy: regsc.DefaultUser,
			Tag:       ctx.Tag,
			EntryId:   ctx.PlasmidID,
			Ontology:  regsc.PlasmidInvOntO,
			Rank:      0,
		},
	}
	return IOE.TryCatchError(func() (*pb.TaggedAnnotation, error) {
		return ctx.Client.CreateAnnotation(context.Background(),
			&pb.NewTaggedAnnotation{Data: input})
	})
}

// createInventoryAnnotations creates annotation IDs for inventory attributes from context
func createInventoryAnnotations(
	ctx WithDeletedInventory,
) IOE.IOEither[error, []string] {
	return F.Pipe2(
		map[string]string{
			regsc.InvLocationTag: ctx.Location,
			regsc.InvStorageDateTag: time.Now().
				Format(time.RFC3339Nano),
		},
		R.Collect(
			func(tag string, value string) IOE.IOEither[error, string] {
				return F.Pipe3(
					AnnotationContext{
						Client:    regsc.GetAnnotationAPIClient(),
						PlasmidID: ctx.PlasmidID,
						Tag:       tag,
						Value:     value,
					},
					createAnnotationIOE,
					IOE.MapLeft[*pb.TaggedAnnotation](
						wrapAnnotationError(tag),
					),
					IOE.Map[error](extractAnnotationID),
				)
			},
		),
		IOE.SequenceArray[error, string],
	)
}

// getInventoryIO retrieves existing inventory for plasmid from context
func getInventoryIO(
	ctx WithPlasmidID,
) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
	filter := fmt.Sprintf(
		"entry_id==%s;tag==%s;ontology==%s",
		ctx.PlasmidID,
		regsc.InvLocationTag,
		regsc.PlasmidInvOntO,
	)
	return F.Pipe1(
		IOE.TryCatchError(func() (*pb.TaggedAnnotationGroupCollection, error) {
			return regsc.GetAnnotationAPIClient().ListAnnotationGroups(
				context.Background(),
				&pb.ListGroupParameters{Filter: filter},
			)
		}),
		IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](func(err error) error {
			return fmt.Errorf("error checking inventory: %w", err)
		}),
	)
}

func toAnnoID(anno *pb.TaggedAnnotationGroup_Data) string {
	return anno.Id
}

func toGroupID(gcd *pb.TaggedAnnotationGroupCollection_Data) string {
	return gcd.Group.GroupId
}

func extractAnnotationIDs(gc *pb.TaggedAnnotationGroupCollection) []string {
	return F.Pipe1(
		gc.Data,
		A.Chain(func(gcd *pb.TaggedAnnotationGroupCollection_Data) []string {
			return F.Pipe1(
				gcd.Group.Data,
				A.Map(toAnnoID),
			)
		}),
	)
}

func extractGroupIDs(gc *pb.TaggedAnnotationGroupCollection) []string {
	return F.Pipe1(
		gc.Data,
		A.Map(toGroupID),
	)
}

func deleteAnnotationIOE(
	client pb.TaggedAnnotationServiceClient,
) func(string) IOE.IOEither[error, string] {
	return func(id string) IOE.IOEither[error, string] {
		return F.Pipe2(
			IOE.TryCatchError(func() (*emptypb.Empty, error) {
				return client.DeleteAnnotation(
					context.Background(),
					&pb.DeleteAnnotationRequest{
						Id: id, Purge: true,
					},
				)
			}),
			IOE.MapLeft[*emptypb.Empty](func(err error) error {
				return fmt.Errorf(
					"error deleting annotation %s: %w",
					id,
					err,
				)
			}),
			IOE.Map[error](func(_ *emptypb.Empty) string {
				return id
			}),
		)
	}
}

func deleteGroupIOE(
	client pb.TaggedAnnotationServiceClient,
) func(string) IOE.IOEither[error, string] {
	return func(groupID string) IOE.IOEither[error, string] {
		return F.Pipe2(
			IOE.TryCatchError(func() (*emptypb.Empty, error) {
				return client.DeleteAnnotationGroup(
					context.Background(),
					&pb.GroupEntryId{GroupId: groupID},
				)
			}),
			IOE.MapLeft[*emptypb.Empty](func(err error) error {
				return fmt.Errorf(
					"error deleting annotation group %s: %w",
					groupID,
					err,
				)
			}),
			IOE.Map[error](func(_ *emptypb.Empty) string {
				return groupID
			}),
		)
	}
}

// deleteInventoryIfExists conditionally deletes existing inventory from context
func deleteInventoryIfExists(
	ctx WithInventory,
) IOE.IOEither[error, *pb.TaggedAnnotationGroupCollection] {
	client := regsc.GetAnnotationAPIClient()

	return F.Pipe3(
		// Step 1: Extract annotation IDs and delete them
		F.Pipe2(
			ctx.Inventory,
			extractAnnotationIDs,
			A.Map(deleteAnnotationIOE(client)),
		),
		IOE.SequenceArray[error, string],
		// Step 2: Extract group IDs and delete them
		IOE.Chain(func(_ []string) IOE.IOEither[error, []string] {
			return F.Pipe3(
				ctx.Inventory,
				extractGroupIDs,
				A.Map(deleteGroupIOE(client)),
				IOE.SequenceArray[error, string],
			)
		}),
		// Step 3: Return original collection
		IOE.Map[error](func(_ []string) *pb.TaggedAnnotationGroupCollection {
			return ctx.Inventory
		}),
	)
}

// createAnnotationGroupIO creates annotation group from context
func createAnnotationGroupIO(
	ctx WithAnnotationIDs,
) IOE.IOEither[error, []string] {
	return F.Pipe1(
		IOE.TryCatchError(func() ([]string, error) {
			_, err := regsc.GetAnnotationAPIClient().
				CreateAnnotationGroup(
					context.Background(),
					&pb.AnnotationIdList{
						Ids: ctx.AnnotationIDs,
					},
				)
			return ctx.AnnotationIDs, err
		}),
		IOE.MapLeft[[]string](func(err error) error {
			return fmt.Errorf(
				"failed to create annotation group: %w",
				err,
			)
		}),
	)
}

// markInventoryExistence marks inventory as existing from context
func markInventoryExistence(
	ctx WithAnnotationGroup,
) IOE.IOEither[error, string] {
	input := &pb.NewTaggedAnnotation{
		Data: &pb.NewTaggedAnnotation_Data{
			Attributes: &pb.NewTaggedAnnotationAttributes{
				Value:     regsc.InvExistValue,
				CreatedBy: regsc.DefaultUser,
				Tag:       regsc.PlasmidInvTag,
				EntryId:   ctx.PlasmidID,
				Ontology:  regsc.PlasmidInvOntO,
			},
		},
	}
	return F.Pipe2(
		// 1. Perform the IO operation
		IOE.TryCatchError(func() (*pb.TaggedAnnotation, error) {
			return regsc.GetAnnotationAPIClient().
				CreateAnnotation(
					context.Background(),
					input,
				)
		}),
		// 2. Handle errors
		IOE.MapLeft[*pb.TaggedAnnotation](func(err error) error {
			return fmt.Errorf(
				"failed to mark inventory existence: %w",
				err,
			)
		}),
		// 3. Map success to the desired output (PlasmidID)
		IOE.Map[error](func(_ *pb.TaggedAnnotation) string {
			return ctx.PlasmidID
		}),
		// 3. Map success to the desired output (PlasmidID)
		IOE.Map[error](func(_ *pb.TaggedAnnotation) string {
			return ctx.PlasmidID
		}),
	)
}

// processRow processes a single inventory record using Do/Bind pattern
func processRow(record InventoryRecord) IOE.IOEither[error, string] {
	return F.Pipe8(
		IOE.Do[error](ProcessRowCtx(record)),
		IOE.Bind(SetPlasmidContext, listPlasmidsIOE),
		IOE.Bind(SetPlasmidID, validatePlasmidResponse),
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
	return F.Pipe2(
		IOE.TryCatchError(func() (*sql.Rows, error) {
			return config.DB.Query(
				"SELECT Name, Location FROM inventory",
			)
		}),
		IOE.MapLeft[*sql.Rows](func(err error) error {
			return fmt.Errorf("error querying inventory: %w", err)
		}),
		IOE.Chain(foldRowsWithSemigroup),
	)
}

// foldRowsWithSemigroup streams over rows one at a time, folding with Semigroup
// This processes rows lazily without loading all into memory
func foldRowsWithSemigroup(
	rows *sql.Rows,
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
				continue
			}

			record := InventoryRecord{
				PlasmidName: name,
				Location:    location,
			}

			// Process single row to get its summary
			rowSummary := processRowToSummary(record)

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
func processRowToSummary(record InventoryRecord) InventoryProcessingSummary {
	return F.Pipe3(
		record,
		processRow,
		fputil.ToEither,
		E.Fold(
			func(err error) InventoryProcessingSummary {
				return InventoryProcessingSummary{
					ErrorCount: 1,
					Errors: []string{
						fmt.Sprintf(
							"%s: %v",
							record.PlasmidName,
							err,
						),
					},
				}
			},
			func(_ string) InventoryProcessingSummary {
				return InventoryProcessingSummary{
					SuccessCount: 1,
				}
			},
		),
	)
}

// extractAnnotationID extracts ID from annotation response
func extractAnnotationID(anno *pb.TaggedAnnotation) string {
	return anno.Data.Id
}

// wrapAnnotationError creates tag-specific error wrapper
func wrapAnnotationError(tag string) func(error) error {
	return func(err error) error {
		return fmt.Errorf(
			"failed to create annotation for %s: %w",
			tag,
			err,
		)
	}
}
