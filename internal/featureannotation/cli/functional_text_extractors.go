package cli

import (
	"strings"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
	"github.com/PuerkitoBio/goquery"
)

// CellTextExtractor represents a function that extracts text from a single cell
type CellTextExtractor func(*goquery.Selection) O.Option[string]

// normalizeText processes text by replacing newlines with spaces
func normalizeText(text string) string {
	return strings.ReplaceAll(text, "\n", " ")
}

// textExtractor creates an Option-returning extractor with selector and normalization
func textExtractor(selector func(*goquery.Selection) string) CellTextExtractor {
	return func(cell *goquery.Selection) O.Option[string] {
		text := strings.TrimSpace(selector(cell))
		if text != "" {
			return O.Some(normalizeText(text))
		}
		return O.None[string]()
	}
}

// extractH2Text extracts text from h2 elements within a cell
var extractH2Text = textExtractor(func(c *goquery.Selection) string {
	return c.Find("h2").Text()
})

// extractDirectText extracts direct text content from a cell
var extractDirectText = textExtractor(func(c *goquery.Selection) string {
	return c.Text()
})

// directTextAlternative creates an alternative function for direct text
// extraction
func directTextAlternative(cell *goquery.Selection) func() O.Option[string] {
	return func() O.Option[string] {
		return extractDirectText(cell)
	}
}

// extractTextFromCellAtIndex extracts text from a cell at a given index
func extractTextFromCellAtIndex(
	cells *goquery.Selection,
) func(int) O.Option[string] {
	return func(index int) O.Option[string] {
		return extractTextFromSingleCell(cells.Eq(index))
	}
}

// createTupleWithIndex creates a tuple combining text with a given index
func createTupleWithIndex(index int) func(string) T.Tuple2[string, int] {
	return func(text string) T.Tuple2[string, int] {
		return T.MakeTuple2(text, index)
	}
}

// extractTextWithIndexFromCell extracts text with index from a cell at a given
// index
func extractTextWithIndexFromCell(
	cells *goquery.Selection,
) func(int) O.Option[T.Tuple2[string, int]] {
	return func(index int) O.Option[T.Tuple2[string, int]] {
		return F.Pipe1(
			extractTextFromSingleCell(cells.Eq(index)),
			O.Map(createTupleWithIndex(index)),
		)
	}
}

// extractTextFromSingleCell tries multiple extraction strategies using pipe
// with Option alternative
func extractTextFromSingleCell(cell *goquery.Selection) O.Option[string] {
	return F.Pipe2(
		cell,
		extractH2Text,
		O.Alt(directTextAlternative(cell)),
	)
}

// createCellRange generates a slice of cell indices for the specified range
func createCellRange(startCol, endCol, maxCells int) []int {
	actualEnd := min(endCol, maxCells-1)
	if startCol > actualEnd {
		return []int{}
	}
	cellIndices := make([]int, 0, actualEnd-startCol+1)
	for i := startCol; i <= actualEnd; i++ {
		cellIndices = append(cellIndices, i)
	}
	return cellIndices
}

// extractTextFromColumns searches for non-empty text in the specified column
// range It first checks for h2 elements within td cells, then falls back to
// direct text content
func extractTextFromColumns(
	cells *goquery.Selection,
	startCol, endCol int,
) string {
	return F.Pipe3(
		createCellRange(startCol, endCol, cells.Length()),
		A.FilterMap(extractTextFromCellAtIndex(cells)),
		A.Head[string],
		O.GetOrElse(F.Constant("")),
	)
}

// extractTextFromColumnsWithIndex searches for non-empty text in the specified
// column range and returns both the text and the column index where it was
// found
func extractTextFromColumnsWithIndex(
	cells *goquery.Selection,
	startCol, endCol int,
) (string, int) {
	result := F.Pipe3(
		createCellRange(startCol, endCol, cells.Length()),
		A.FilterMap(extractTextWithIndexFromCell(cells)),
		A.Head[T.Tuple2[string, int]],
		O.GetOrElse(F.Constant(T.MakeTuple2("", -1))),
	)
	return result.F1, result.F2
}
