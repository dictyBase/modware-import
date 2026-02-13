package stockcenter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	N "github.com/IBM/fp-go/number"
	S "github.com/IBM/fp-go/semigroup"

	stockpb "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/fputil"
	"github.com/dictyBase/modware-import/internal/logger"
	regsc "github.com/dictyBase/modware-import/internal/registry/stockcenter"
	"github.com/urfave/cli/v2"
)

const maxErrorSamples = 5

// OntologyUpdateStats tracks update operation metrics
type OntologyUpdateStats struct {
	ProcessedCount int       `json:"processed_count"` // Total plasmids examined
	UpdatedCount   int       `json:"updated_count"`   // Successfully updated
	ErrorCount     int       `json:"error_count"`     // Failed updates
	Errors         []error   `json:"-"`               // Error samples (first 5)
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
}

// UpdateContext - Context struct for updatePlasmidTerm (point-free, no lambda wrapper)
type UpdateContext struct {
	Client  stockpb.StockServiceClient
	Term    string
	Plasmid *stockpb.Plasmid
}

// ProcessContext - Context struct for processPlasmid (point-free, no lambda wrapper)
type ProcessContext struct {
	Client  stockpb.StockServiceClient
	Term    string
	Plasmid *stockpb.PlasmidCollection_Data
}

// Reusable from existing code
var (
	intSumSemigroup = N.SemigroupSum[int]()
	errorsSemigroup = S.MakeSemigroup(func(a, b []error) []error {
		return F.Pipe1(
			A.ArrayConcatAll(a, b),
			A.Slice[error](0, maxErrorSamples),
		)
	})

	// countPlasmids counts plasmids in a slice (point-free)
	countPlasmids = func(plasmids []*stockpb.PlasmidCollection_Data) int {
		return len(plasmids)
	}

	// createInitialStats creates initial stats with count (point-free)
	createInitialStats = func(count int) OntologyUpdateStats {
		return OntologyUpdateStats{ProcessedCount: count}
	}

	updateErrorSummary = func(err error) OntologyUpdateStats {
		return OntologyUpdateStats{
			ErrorCount: 1,
			Errors:     []error{err},
		}
	}

	updateSuccessSummary = func(_ *stockpb.Plasmid) OntologyUpdateStats {
		return OntologyUpdateStats{
			UpdatedCount: 1,
		}
	}
)

// StatsSemigroup creates a semigroup for combining stats
func StatsSemigroup() S.Semigroup[OntologyUpdateStats] {
	return S.MakeSemigroup(
		func(a, b OntologyUpdateStats) OntologyUpdateStats {
			return OntologyUpdateStats{
				ProcessedCount: intSumSemigroup.Concat(
					a.ProcessedCount,
					b.ProcessedCount,
				),
				UpdatedCount: intSumSemigroup.Concat(
					a.UpdatedCount,
					b.UpdatedCount,
				),
				ErrorCount: intSumSemigroup.Concat(
					a.ErrorCount,
					b.ErrorCount,
				),
				Errors:    errorsSemigroup.Concat(a.Errors, b.Errors),
				StartTime: a.StartTime,
				EndTime:   b.EndTime,
			}
		},
	)
}

// updatePlasmidTerm updates a plasmid's property to the target term
// Point-free: takes UpdateContext directly, no lambda wrapper
func updatePlasmidTerm(
	ctx ProcessContext,
) IOE.IOEither[error, *stockpb.Plasmid] {
	input := &stockpb.PlasmidUpdate{
		Data: &stockpb.PlasmidUpdate_Data{
			Type: "plasmid",
			Id:   ctx.Plasmid.Id,
			Attributes: &stockpb.PlasmidUpdateAttributes{
				UpdatedBy:            regsc.DefaultUser,
				DictyPlasmidProperty: ctx.Term,
			},
		},
	}
	return F.Pipe1(
		IOE.TryCatchError(func() (*stockpb.Plasmid, error) {
			return ctx.Client.UpdatePlasmid(
				context.Background(),
				input,
			)
		}),
		IOE.MapLeft[*stockpb.Plasmid](func(err error) error {
			return fmt.Errorf(
				"failed to update plasmid %s: %w",
				ctx.Plasmid.Id,
				err,
			)
		}),
	)
}

func processBatch(
	plasmids []*stockpb.PlasmidCollection_Data,
	client stockpb.StockServiceClient,
	term string,
) OntologyUpdateStats {
	return F.Pipe5(
		plasmids,
		A.Map(func(pl *stockpb.PlasmidCollection_Data) ProcessContext {
			return ProcessContext{
				Client:  client,
				Term:    term,
				Plasmid: pl,
			}
		}),
		A.Map(updatePlasmidTerm), // Process each context (point-free!)
		A.Map(fputil.ToEither[error, *stockpb.Plasmid]),
		A.Map(E.Fold(updateErrorSummary, updateSuccessSummary)),
		A.Reduce(StatsSemigroup().Concat,
			F.Pipe2(plasmids, countPlasmids, createInitialStats),
		),
	)
}

// LoadPlasmidOntologyCli is the CLI entry point for ontology updates
func LoadPlasmidOntologyCli(cmd *cli.Context) error {
	handler := logger.GetCliSlogHandler(cmd)
	slogger := slog.New(handler)
	client := regsc.GetStockAPIClient()

	term := cmd.String("ontology-term")
	batchSize := cmd.Int("batch-size")

	cursor := int64(0)
	totalStats := OntologyUpdateStats{
		StartTime: time.Now(),
	}

	slogger.Info("Starting plasmid ontology update",
		"target_term", term,
		"batch_size", batchSize,
	)

	// Pagination loop
	filter := fmt.Sprintf("tag!=%s;tag!=GB vector", term)
	for {
		params := &stockpb.StockParameters{
			Limit:  int64(batchSize),
			Cursor: cursor,
			Filter: filter,
		}
		resp, err := client.ListPlasmids(context.Background(), params)
		if err != nil {
			return fmt.Errorf(
				"failed to list plasmids at cursor %d: %w",
				cursor,
				err,
			)
		}

		runningStats := processBatch(resp.Data, client, term)
		runningStats.StartTime = totalStats.StartTime
		runningStats.EndTime = time.Now()
		totalStats = StatsSemigroup().Concat(totalStats, runningStats)

		// Pagination: check for next page
		if resp.Meta.NextCursor == 0 {
			break
		}
		cursor = resp.Meta.NextCursor
	}

	totalStats.EndTime = time.Now()
	return handleOntologyOutput(totalStats, slogger)
}

// handleOntologyOutput logs final results and returns error if too many failures
func handleOntologyOutput(
	stats OntologyUpdateStats,
	slogger *slog.Logger,
) error {
	if stats.ErrorCount > 0 {
		joinedErrors := errors.Join(stats.Errors...)
		slogger.Log(
			context.Background(),
			slog.LevelError,
			"errors encountered",
			"top 5 errors", joinedErrors,
		)
		return joinedErrors
	}
	slogger.Log(context.Background(), slog.LevelInfo,
		"Plasmid ontology update complete",
		"processed", stats.ProcessedCount,
		"updated", stats.UpdatedCount,
		"errors", stats.ErrorCount,
		"duration_ms", stats.EndTime.Sub(stats.StartTime).Milliseconds(),
	)
	return nil
}
