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
	Deps             Deps
	PlasmidName      string
	Location         string
	PlasmidID        string
	Inventory        *pb.TaggedAnnotationGroupCollection
	InventoryDeleted bool
	AnnotationIDs    []string
	GroupCreated     bool
}

type AnnotationContext struct {
	Client    pb.TaggedAnnotationServiceClient
	PlasmidID string
	Tag       string
	Value     string
}

type InventoryProcessingSummary struct {
	SuccessCount  int
	ErrorCount    int
	ExpectedCount int
	Errors        []string
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

// listPlasmidsIOE calls ListPlasmids and returns an Option: None if not found, Some(first) if found.
func listPlasmidsIOE(
	ctx PipelineContext,
) IOE.IOEither[error, O.Option[*stock.Plasmid]] {
	return F.Pipe4(
		IOE.TryCatchError(func() (*stock.PlasmidCollection, error) {
			return ctx.Deps.StockClient.ListPlasmids(
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
		IOE.MapLeft[*stock.PlasmidCollection](func(err error) error {
			return fmt.Errorf(
				"error listing plasmids for name %s: %w",
				ctx.PlasmidName,
				err,
			)
		}),
		IOE.Map[error](func(c *stock.PlasmidCollection) []*stock.PlasmidCollection_Data {
			return c.Data
		}),
		IOE.Map[error](A.Head[*stock.PlasmidCollection_Data]),
		IOE.Map[error](O.Map(convertCollectionDataToPlasmid)),
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

// createInventoryAnnotationsRE creates annotation IDs for inventory attributes from context
func createInventoryAnnotationsRE(
	ctx PipelineContext,
) Processing[[]string] {
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
							Client:    ctx.Deps.AnnotationClient,
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
}

// getInventoryRE retrieves existing inventory for plasmid from context
func getInventoryRE(
	ctx PipelineContext,
) Processing[*pb.TaggedAnnotationGroupCollection] {
	filter := fmt.Sprintf(
		"entry_id===%s;tag===%s;ontology===%s",
		ctx.PlasmidID,
		regsc.InvLocationTag,
		regsc.PlasmidInvOntO,
	)
	return RE.FromIOEither[Deps](
		F.Pipe1(
			IOE.TryCatchError(
				func() (*pb.TaggedAnnotationGroupCollection, error) {
					return ctx.Deps.AnnotationClient.ListAnnotationGroups(
						context.Background(),
						&pb.ListGroupParameters{Filter: filter},
					)
				},
			),
			IOE.MapLeft[*pb.TaggedAnnotationGroupCollection](
				func(err error) error {
					return fmt.Errorf("error checking inventory: %w", err)
				},
			),
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

// deleteInventoryIfExistsRE conditionally deletes existing inventory from context
func deleteInventoryIfExistsRE(
	ctx PipelineContext,
) Processing[*pb.TaggedAnnotationGroupCollection] {
	return RE.FromIOEither[Deps](
		F.Pipe3(
			F.Pipe2(
				ctx.Inventory,
				extractAnnotationIDs,
				A.Map(deleteAnnotationIOE(ctx.Deps.AnnotationClient)),
			),
			IOE.SequenceArray[error, string],
			IOE.Chain(
				func(_ []string) IOE.IOEither[error, []string] {
					return F.Pipe3(
						ctx.Inventory,
						extractGroupIDs,
						A.Map(deleteGroupIOE(ctx.Deps.AnnotationClient)),
						IOE.SequenceArray[error, string],
					)
				},
			),
			IOE.Map[error](
				func(_ []string) *pb.TaggedAnnotationGroupCollection {
					return ctx.Inventory
				},
			),
		),
	)
}

// createAnnotationGroupRE creates annotation group from context
func createAnnotationGroupRE(
	ctx PipelineContext,
) Processing[[]string] {
	return RE.FromIOEither[Deps](
		F.Pipe1(
			IOE.TryCatchError(func() ([]string, error) {
				_, err := ctx.Deps.AnnotationClient.CreateAnnotationGroup(
					context.Background(),
					&pb.AnnotationIdList{
						Ids: ctx.AnnotationIDs,
					},
				)
				return ctx.AnnotationIDs, err
			}),
			IOE.MapLeft[[]string](func(err error) error {
				return fmt.Errorf("failed to create annotation group: %w", err)
			}),
		),
	)
}

// markInventoryExistenceRE marks inventory as existing from context
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
	return RE.FromIOEither[Deps](
		F.Pipe2(
			IOE.TryCatchError(func() (*pb.TaggedAnnotation, error) {
				return ctx.Deps.AnnotationClient.CreateAnnotation(
					context.Background(),
					input,
				)
			}),
			IOE.MapLeft[*pb.TaggedAnnotation](func(err error) error {
				return fmt.Errorf("failed to mark inventory existence: %w", err)
			}),
			IOE.Map[error](func(_ *pb.TaggedAnnotation) string {
				return ctx.PlasmidID
			}),
		),
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
func foldRowsWithSemigroup(
	deps Deps,
	rows *sql.Rows,
) IOE.IOEither[error, InventoryProcessingSummary] {
	return IOE.TryCatchError(func() (InventoryProcessingSummary, error) {
		defer rows.Close()

		summary := InventoryProcessingSummary{}
		semigroup := InventorySummarySemigroup()

		for rows.Next() {
			ctx := PipelineContext{Deps: deps}
			if err := rows.Scan(
				&ctx.PlasmidName,
				&ctx.Location,
			); err != nil {
				summary = semigroup.Concat(
					summary,
					InventoryProcessingSummary{
						ErrorCount: 1,
						Errors: []string{fmt.Sprintf(
							"scan error: %v",
							err,
						)},
					})
				continue
			}
			summary = semigroup.Concat(
				summary,
				processRowToSummary(ctx),
			)
		}
		if err := rows.Err(); err != nil {
			return summary, fmt.Errorf(
				"error iterating rows: %w",
				err,
			)
		}
		return summary, nil
	})
}

// skipInventory is the O.Fold None branch: plasmid not in stock, nothing to annotate.
func skipInventory() E.Either[error, InventoryProcessingSummary] {
	return E.Right[error](InventoryProcessingSummary{})
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
				fmt.Sprintf("%s: %v", name, err),
			},
		}
	},
)

// runProcessing applies ctx to a Processing pipeline, runs the IO, and maps to summary.
func runProcessing(
	ctx PipelineContext,
) func(Processing[string]) E.Either[error, InventoryProcessingSummary] {
	return func(pipeline Processing[string]) E.Either[error, InventoryProcessingSummary] {
		return F.Pipe2(
			pipeline(ctx.Deps),
			fputil.ToEither[error, string],
			E.Map[error](onProcessSuccess),
		)
	}
}

// processFoundPlasmid runs the full inventory pipeline for a found plasmid.
// Curried: processFoundPlasmid(ctx)(plasmid) — used directly as the O.Fold Some branch.
var processFoundPlasmid = F.Curry2(
	func(ctx PipelineContext, plasmid *stock.Plasmid) E.Either[error, InventoryProcessingSummary] {
		ctx.PlasmidID = plasmid.Data.Id
		return F.Pipe1(
			F.Pipe5(
				RE.Of[Deps, error](ctx),
				RE.Bind(setInventory, getInventoryRE),
				RE.Bind(setInventoryDeleted, deleteInventoryIfExistsRE),
				RE.Bind(setAnnotationIDs, createInventoryAnnotationsRE),
				RE.Bind(setGroupCreated, createAnnotationGroupRE),
				RE.Chain(markInventoryExistenceRE),
			),
			runProcessing(ctx),
		)
	},
)

// processRowToSummary processes one inventory row and returns a summary.
// Phase 1: look up plasmid by name → Either[error, Option[*Plasmid]].
// Phase 2: O.Fold — None skips silently, Some calls processFoundPlasmid(ctx) directly.
func processRowToSummary(ctx PipelineContext) InventoryProcessingSummary {
	return F.Pipe4(
		ctx,
		listPlasmidsIOE,
		fputil.ToEither,
		E.Chain(O.Fold(skipInventory, processFoundPlasmid(ctx))),
		E.Fold(
			onProcessError(ctx.PlasmidName),
			F.Identity[InventoryProcessingSummary],
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
