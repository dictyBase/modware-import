package stockcenter

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	RE "github.com/IBM/fp-go/readerioeither"
	R "github.com/IBM/fp-go/record"
	S "github.com/IBM/fp-go/semigroup"
	pb "github.com/dictyBase/go-genproto/dictybaseapis/annotation"
	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/fputil"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	plasmidSearchLimit = 1

	inventoryJoinQuery = `
		SELECT DISTINCT g."Plasmid Name", gi.Location
		FROM goldenbraid g
		JOIN goldenbraid_inventory gi
			ON gi.Name = g."Plasmid Name" OR gi.Name = g.Synonym`

	inventoryCountQuery = `
		SELECT COUNT(*) FROM (
			SELECT DISTINCT g."Plasmid Name", gi.Location
			FROM goldenbraid g
			JOIN goldenbraid_inventory gi
				ON gi.Name = g."Plasmid Name" OR gi.Name = g.Synonym
		) AS distinct_rows`
)

// Deps holds dependencies for inventory processing
type Deps struct {
	StockClient      stock.StockServiceClient
	AnnotationClient pb.TaggedAnnotationServiceClient
	Logger           *slog.Logger
}

// Processing is a ReaderIOEither that injects Deps
type Processing[A any] = RE.ReaderIOEither[Deps, error, A]

// PipelineContext represents the complete state of inventory processing
type PipelineContext struct {
	PlasmidName      string
	Location         string
	NameContext      *PlasmidNameContext
	PlasmidID        string
	Inventory        *pb.TaggedAnnotationGroupCollection
	InventoryDeleted bool
	AnnotationIDs    []string
	GroupCreated     bool
}

// NewPipelineContext creates context from inventory record
func NewPipelineContext(record InventoryRecord) PipelineContext {
	return PipelineContext{
		PlasmidName: record.PlasmidName,
		Location:    record.Location,
	}
}

type AnnotationContext struct {
	Client    pb.TaggedAnnotationServiceClient
	PlasmidID string
	Tag       string
	Value     string
}

type InventoryRecord struct {
	PlasmidName string
	Location    string
}

type InventoryProcessingSummary struct {
	SuccessCount  int
	ErrorCount    int
	ExpectedCount int
	Errors        []string
}

type PlasmidNameContext struct {
	Name string
	Resp *stock.PlasmidCollection
}

// setPlasmidContext updates context with API lookup result
func setPlasmidContext(
	nameCtx PlasmidNameContext,
) func(PipelineContext) PipelineContext {
	return func(ctx PipelineContext) PipelineContext {
		ctx.NameContext = &nameCtx
		return ctx
	}
}

// setPlasmidID updates context with validated plasmid ID
func setPlasmidID(id string) func(PipelineContext) PipelineContext {
	return func(ctx PipelineContext) PipelineContext {
		ctx.PlasmidID = id
		return ctx
	}
}

// setInventory updates context with existing inventory
func setInventory(
	inv *pb.TaggedAnnotationGroupCollection,
) func(PipelineContext) PipelineContext {
	return func(ctx PipelineContext) PipelineContext {
		ctx.Inventory = inv
		return ctx
	}
}

// setInventoryDeleted marks inventory as deleted
func setInventoryDeleted(
	_ *pb.TaggedAnnotationGroupCollection,
) func(PipelineContext) PipelineContext {
	return func(ctx PipelineContext) PipelineContext {
		ctx.InventoryDeleted = true
		return ctx
	}
}

// setAnnotationIDs updates context with created annotation IDs
func setAnnotationIDs(ids []string) func(PipelineContext) PipelineContext {
	return func(ctx PipelineContext) PipelineContext {
		ctx.AnnotationIDs = ids
		return ctx
	}
}

// setGroupCreated marks annotation group as created
func setGroupCreated(_ []string) func(PipelineContext) PipelineContext {
	return func(ctx PipelineContext) PipelineContext {
		ctx.GroupCreated = true
		return ctx
	}
}

func InventorySummarySemigroup() S.Semigroup[InventoryProcessingSummary] {
	return S.MakeSemigroup(
		func(a, b InventoryProcessingSummary) InventoryProcessingSummary {
			return InventoryProcessingSummary{
				SuccessCount:  a.SuccessCount + b.SuccessCount,
				ErrorCount:    a.ErrorCount + b.ErrorCount,
				ExpectedCount: max(a.ExpectedCount, b.ExpectedCount),
				Errors:        append(a.Errors, b.Errors...),
			}
		},
	)
}

// listPlasmidsRE calls API and maps result to PlasmidNameContext using Reader
func listPlasmidsRE(
	ctx PipelineContext,
) Processing[PlasmidNameContext] {
	return F.Pipe1(
		RE.Ask[Deps, error](),
		RE.Chain(func(deps Deps) Processing[PlasmidNameContext] {
			return RE.FromIOEither[Deps](
				F.Pipe2(
					IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
						return deps.StockClient.ListPlasmids(
							context.Background(),
							&stock.StockParameters{
								Filter: fmt.Sprintf(
									"plasmid_name===%s",
									ctx.PlasmidName,
								),
								Limit: plasmidSearchLimit,
							},
						)
					}),
					IOE.MapLeft[*stock.PlasmidCollection](
						func(err error) error {
							return fmt.Errorf(
								"error listing plasmids for name %s: %w",
								ctx.PlasmidName,
								err,
							)
						},
					),
					IOE.Map[error](
						func(resp *stock.PlasmidCollection) PlasmidNameContext {
							return PlasmidNameContext{
								Name: ctx.PlasmidName,
								Resp: resp,
							}
						},
					),
				),
			)
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

// requireNameContext converts nullable pointer to Either using fp-go patterns
func requireNameContext(
	ctx PipelineContext,
) E.Either[error, PlasmidNameContext] {
	return F.Pipe2(
		O.FromNillable(ctx.NameContext),
		O.Map(func(nc *PlasmidNameContext) PlasmidNameContext { return *nc }),
		E.FromOption[PlasmidNameContext](func() error {
			return fmt.Errorf(
				"nameContext not populated for plasmid %s",
				ctx.PlasmidName,
			)
		}),
	)
}

// validatePlasmidResponseRE validates plasmid response and extracts ID using pure Either logic
func validatePlasmidResponseRE(
	ctx PipelineContext,
) Processing[string] {
	return RE.FromIOEither[Deps](
		IOE.FromEither(
			F.Pipe2(
				requireNameContext(ctx),
				E.Chain(
					func(nc PlasmidNameContext) E.Either[error, PlasmidNameContext] {
						return E.FromPredicate(
							hasSinglePlasmid,
							plasmidCountError,
						)(
							nc,
						)
					},
				),
				E.Map[error](firstPlasmidID),
			),
		),
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

// createInventoryAnnotationsRE creates annotation IDs for inventory attributes from context using Reader
func createInventoryAnnotationsRE(
	ctx PipelineContext,
) Processing[[]string] {
	return F.Pipe1(
		RE.Ask[Deps, error](),
		RE.Chain(func(deps Deps) Processing[[]string] {
			return RE.FromIOEither[Deps](
				F.Pipe2(
					map[string]string{
						regsc.InvLocationTag: ctx.Location,
						regsc.InvStorageDateTag: time.Now().
							Format(time.RFC3339Nano),
					},
					R.Collect(
						func(tag string, value string) IOE.IOEither[error, string] {
							return F.Pipe3(
								AnnotationContext{
									Client:    deps.AnnotationClient,
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
				),
			)
		}),
	)
}

// getInventoryRE retrieves existing inventory for plasmid from context
func getInventoryRE(
	ctx PipelineContext,
) Processing[*pb.TaggedAnnotationGroupCollection] {
	filter := fmt.Sprintf(
		"entry_id==%s;tag==%s;ontology==%s",
		ctx.PlasmidID,
		regsc.InvLocationTag,
		regsc.PlasmidInvOntO,
	)
	return F.Pipe1(
		RE.Ask[Deps, error](),
		RE.Chain(
			func(deps Deps) Processing[*pb.TaggedAnnotationGroupCollection] {
				return RE.FromIOEither[Deps](
					F.Pipe1(
						IOE.TryCatchError(
							func() (*pb.TaggedAnnotationGroupCollection, error) {
								return deps.AnnotationClient.ListAnnotationGroups(
									context.Background(),
									&pb.ListGroupParameters{Filter: filter},
								)
							},
						),
						IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](
							func(err error) error {
								return fmt.Errorf(
									"error checking inventory: %w",
									err,
								)
							},
						),
					),
				)
			},
		),
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

// deleteInventoryIfExistsRE conditionally deletes existing inventory from context using Reader
func deleteInventoryIfExistsRE(
	ctx PipelineContext,
) Processing[*pb.TaggedAnnotationGroupCollection] {
	return F.Pipe1(
		RE.Ask[Deps, error](),
		RE.Chain(
			func(deps Deps) Processing[*pb.TaggedAnnotationGroupCollection] {
				return RE.FromIOEither[Deps](
					F.Pipe3(
						// Step 1: Extract annotation IDs and delete them
						F.Pipe2(
							ctx.Inventory,
							extractAnnotationIDs,
							A.Map(deleteAnnotationIOE(deps.AnnotationClient)),
						),
						IOE.SequenceArray[error, string],
						// Step 2: Extract group IDs and delete them
						IOE.Chain(
							func(_ []string) IOE.IOEither[error, []string] {
								return F.Pipe3(
									ctx.Inventory,
									extractGroupIDs,
									A.Map(
										deleteGroupIOE(deps.AnnotationClient),
									),
									IOE.SequenceArray[error, string],
								)
							},
						),
						// Step 3: Return original collection
						IOE.Map[error](
							func(_ []string) *pb.TaggedAnnotationGroupCollection {
								return ctx.Inventory
							},
						),
					),
				)
			},
		),
	)
}

// createAnnotationGroupRE creates annotation group from context using Reader
func createAnnotationGroupRE(
	ctx PipelineContext,
) Processing[[]string] {
	return F.Pipe1(
		RE.Ask[Deps, error](),
		RE.Chain(func(deps Deps) Processing[[]string] {
			return RE.FromIOEither[Deps](
				F.Pipe1(
					IOE.TryCatchError(func() ([]string, error) {
						_, err := deps.AnnotationClient.CreateAnnotationGroup(
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
				),
			)
		}),
	)
}

// markInventoryExistenceRE marks inventory as existing from context using Reader
func markInventoryExistenceRE(
	ctx PipelineContext,
) Processing[string] {
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
	return F.Pipe1(
		RE.Ask[Deps, error](),
		RE.Chain(func(deps Deps) Processing[string] {
			return RE.FromIOEither[Deps](
				F.Pipe2(
					// 1. Perform the IO operation
					IOE.TryCatchError(func() (*pb.TaggedAnnotation, error) {
						return deps.AnnotationClient.CreateAnnotation(
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
				),
			)
		}),
	)
}

// queryInventoryCount runs the COUNT query against the joined tables
func queryInventoryCount(db *sql.DB) IOE.IOEither[error, int] {
	return IOE.TryCatchError(func() (int, error) {
		var count int
		err := db.QueryRow(inventoryCountQuery).Scan(&count)
		return count, err
	})
}

// queryInventoryRows runs the JOIN query and returns the result rows
func queryInventoryRows(db *sql.DB) IOE.IOEither[error, *sql.Rows] {
	return F.Pipe1(
		IOE.TryCatchError(func() (*sql.Rows, error) {
			return db.Query(inventoryJoinQuery)
		}),
		IOE.MapLeft[*sql.Rows](func(err error) error {
			return fmt.Errorf("error querying inventory: %w", err)
		}),
	)
}

// stampExpectedCount returns a function that immutably sets ExpectedCount on a summary
func stampExpectedCount(
	summary InventoryProcessingSummary,
) func(int) InventoryProcessingSummary {
	return func(expected int) InventoryProcessingSummary {
		return InventoryProcessingSummary{
			SuccessCount:  summary.SuccessCount,
			ErrorCount:    summary.ErrorCount,
			ExpectedCount: expected,
			Errors:        summary.Errors,
		}
	}
}

func ProcessInventory(
	deps Deps,
	config InventoryLoaderConfig,
) IOE.IOEither[error, InventoryProcessingSummary] {
	return F.Pipe2(
		queryInventoryRows(config.DB),
		IOE.Chain(
			func(rows *sql.Rows) IOE.IOEither[error, InventoryProcessingSummary] {
				return foldRowsWithSemigroup(deps, rows)
			},
		),
		IOE.Chain(
			func(summary InventoryProcessingSummary) IOE.IOEither[error, InventoryProcessingSummary] {
				return F.Pipe1(
					queryInventoryCount(config.DB),
					IOE.Map[error](stampExpectedCount(summary)),
				)
			},
		),
	)
}

// foldRowsWithSemigroup streams over rows one at a time, folding with Semigroup
// This processes rows lazily without loading all into memory
func foldRowsWithSemigroup(
	deps Deps,
	rows *sql.Rows,
) IOE.IOEither[error, InventoryProcessingSummary] {
	return IOE.TryCatchError(func() (InventoryProcessingSummary, error) {
		defer rows.Close()

		// Initial accumulator (identity for Semigroup)
		summary := InventoryProcessingSummary{}
		semigroup := InventorySummarySemigroup()

		// Stream over rows, folding one at a time
		for rows.Next() {
			var record InventoryRecord
			if err := rows.Scan(
				&record.PlasmidName,
				&record.Location,
			); err != nil {
				scanError := InventoryProcessingSummary{
					ErrorCount: 1,
					Errors: []string{fmt.Sprintf(
						"scan error: %v",
						err,
					)},
				}
				summary = semigroup.Concat(summary, scanError)
				continue
			}
			rowSummary := processRowToSummary(deps, record)
			summary = semigroup.Concat(summary, rowSummary)
		}
		// Check for iteration errors
		if err := rows.Err(); err != nil {
			return summary, fmt.Errorf("error iterating rows: %w", err)
		}
		return summary, nil
	})
}

// onProcessSuccess handles successful processing
func onProcessSuccess(_ string) InventoryProcessingSummary {
	return InventoryProcessingSummary{
		SuccessCount: 1,
	}
}

// onProcessError handles processing failure
var onProcessError = F.Curry2(
	func(name string, err error) InventoryProcessingSummary {
		return InventoryProcessingSummary{
			ErrorCount: 1,
			Errors: []string{
				fmt.Sprintf(
					"%s: %v",
					name,
					err,
				),
			},
		}
	},
)

// processRowToSummary processes one row and returns a summary with injected dependencies
func processRowToSummary(
	deps Deps,
	record InventoryRecord,
) InventoryProcessingSummary {
	return F.Pipe9(
		RE.Of[Deps, error](NewPipelineContext(record)),
		RE.Bind(setPlasmidContext, listPlasmidsRE),
		RE.Bind(setPlasmidID, validatePlasmidResponseRE),
		RE.Bind(setInventory, getInventoryRE),
		RE.Bind(setInventoryDeleted, deleteInventoryIfExistsRE),
		RE.Bind(setAnnotationIDs, createInventoryAnnotationsRE),
		RE.Bind(setGroupCreated, createAnnotationGroupRE),
		RE.Chain(markInventoryExistenceRE),
		func(r Processing[string]) E.Either[error, string] {
			return fputil.ToEither(r(deps))
		},
		E.Fold(
			onProcessError(record.PlasmidName),
			onProcessSuccess,
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
