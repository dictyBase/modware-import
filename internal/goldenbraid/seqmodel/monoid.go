package seqmodel

import (
	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
)

var (
	summaryL = MakeExportSummaryLenses()
	// ExportSummaryMonoid adds counts and concatenates errors in source order.
	ExportSummaryMonoid = M.MakeMonoid(concatSummary, ExportSummary{})
	succeeded           = F.Pipe3(
		ExportSummary{},
		summaryL.Total.Set(1),
		summaryL.Succeeded.Set(1),
		F.Constant1[string, ExportSummary],
	)
)

func concatSummary(left, right ExportSummary) ExportSummary {
	total := summaryL.Total.Get(left) + summaryL.Total.Get(right)
	succeededCount := summaryL.Succeeded.Get(left) + summaryL.Succeeded.Get(right)
	failedCount := summaryL.Failed.Get(left) + summaryL.Failed.Get(right)
	leftErrors := summaryL.Errors.Get(left)
	rightErrors := summaryL.Errors.Get(right)
	errors := A.ArrayConcatAll(leftErrors, rightErrors)
	return F.Pipe4(ExportSummary{},
		summaryL.Total.Set(total),
		summaryL.Succeeded.Set(succeededCount),
		summaryL.Failed.Set(failedCount),
		summaryL.Errors.Set(errors))
}

func failed(err error) ExportSummary {
	return F.Pipe3(
		ExportSummary{},
		summaryL.Total.Set(1),
		summaryL.Failed.Set(1),
		summaryL.Errors.Set([]error{err}),
	)
}

// ReduceExportResults folds each attempted export into an ordered summary.
func ReduceExportResults(results []E.Either[error, string]) ExportSummary {
	return F.Pipe2(results, A.Map(E.Fold(failed, succeeded)), M.ConcatAll(ExportSummaryMonoid))
}
