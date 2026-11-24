package normalize

import (
	A "github.com/IBM/fp-go/v2/array"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/predicate"
)

const summaryTemplate = "files={{len .Files}} blank={{.Batch.Blank}} " +
	"existing={{len .Batch.Existing}} new={{len .Batch.New}} " +
	"rejected={{.Batch.Rejected}} warnings={{.Batch.Warnings}}"

// kindIs matches an outcome of the given kind.
func kindIs(kind RowKind) P.Predicate[RowOutcome] {
	return F.Pipe2(kind, EQ.Equals(kindEq), P.ContraMap(outcomeL.Kind.Get))
}

// recordIf returns the record of o when o has the given kind.
func recordIf(o RowOutcome, kind RowKind) []Record {
	return F.Pipe3(
		o,
		O.FromPredicate(kindIs(kind)),
		O.Map(outcomeL.Record.Get),
		O.Fold(A.Empty[Record], A.Of[Record]),
	)
}

// countIf returns 1 when o has the given kind, else 0.
func countIf(o RowOutcome, kind RowKind) int {
	return F.Pipe2(
		o,
		O.FromPredicate(kindIs(kind)),
		O.Fold(F.Constant(0), F.Constant1[RowOutcome](1)),
	)
}

// outcomeBatch turns one classified row into a batch of one row. A
// warning of a rejected row stays in the issues. A diagnostic is never
// counted as a rejected row.
func outcomeBatch(o RowOutcome) Batch {
	issues := outcomeL.Issues.Get(o)

	return F.Pipe6(
		Batch{},
		batchL.Existing.Set(recordIf(o, KindExisting)),
		batchL.New.Set(recordIf(o, KindNew)),
		batchL.Issues.Set(issues),
		batchL.Warnings.Set(F.Pipe2(issues, A.Filter(isWarning), A.Size[Diagnostic])),
		batchL.Blank.Set(countIf(o, KindBlank)),
		batchL.Rejected.Set(countIf(o, KindRejected)),
	)
}

// joinBatch appends batch b after batch a. It keeps the row order.
func joinBatch(a, b Batch) Batch {
	return F.Pipe6(
		a,
		batchL.Existing.Set(A.ArrayConcatAll(batchL.Existing.Get(a), batchL.Existing.Get(b))),
		batchL.New.Set(A.ArrayConcatAll(batchL.New.Get(a), batchL.New.Get(b))),
		batchL.Issues.Set(A.ArrayConcatAll(batchL.Issues.Get(a), batchL.Issues.Get(b))),
		batchL.Warnings.Set(batchL.Warnings.Get(a)+batchL.Warnings.Get(b)),
		batchL.Blank.Set(batchL.Blank.Get(a)+batchL.Blank.Get(b)),
		batchL.Rejected.Set(batchL.Rejected.Get(a)+batchL.Rejected.Get(b)),
	)
}

// addRow is the row pipeline: align the cells, classify the row, and
// add it to the batch.
func addRow(b Batch, l line) Batch {
	return F.Pipe4(l, toSourceRow, classify, outcomeBatch, F.Bind1st(joinBatch, b))
}

// PrintSummary prints one summary line and returns the state.
func PrintSummary(s State) IO.IO[State] {
	return F.Pipe1(s, IO.PrintGo[State](summaryTemplate))
}
