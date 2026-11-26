package normalize

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	endoutils "github.com/dictyBase/fp-go-loom/endomorphism"
	"github.com/xuri/excelize/v2"
)

var (
	// lineOrigin is the source file and the row number of a line.
	lineOrigin = endoutils.BuildFrom([]func(line) En.Endomorphism[Origin]{
		F.Flow2(lineL.File.Get, originL.File.Set),
		F.Flow2(lineL.Row.Get, originL.Row.Set),
	})

	// toSourceRow keeps the raw cells and a trimmed copy. The trimmed copy
	// goes into the checks and the records. The raw copy goes into the
	// diagnostics.
	toSourceRow = endoutils.BuildFrom([]func(line) En.Endomorphism[SourceRow]{
		F.Flow2(lineOrigin, rowL.Origin.Set),
		F.Flow3(expandCells, toCells, rowL.Raw.Set),
		F.Flow4(expandCells, A.Map(strings.TrimSpace), toCells, rowL.Clean.Set),
		F.Flow2(overflow, rowL.Extra.Set),
	})

	wantsSheet = F.Pipe1(Str.IsNonEmpty, P.ContraMap(requestL.Want.Get))
	// chooseSheet uses the first worksheet when the option is empty.
	chooseSheet = F.Pipe1(wantsSheet, P.Fold(firstSheet, namedSheet))
)

// expandCells puts the cells of one line in internal column order.
func expandCells(l line) []string {
	return F.Pipe1(
		layoutColumns[lineL.Layout.Get(l)],
		A.Map(F.Bind1st(cellAt, lineL.Cells.Get(l))),
	)
}

// overflow returns the first nonblank cell after the last header column.
func overflow(l line) string {
	return F.Pipe3(
		l,
		lineL.Cells.Get,
		A.SliceRight[string](layoutWidth[lineL.Layout.Get(l)]),
		firstNonblank,
	)
}

// newLine records one data row as read.
func newLine(file string, layout Layout, row int, cells []string) line {
	return F.Pipe4(
		line{},
		lineL.File.Set(file),
		lineL.Layout.Set(layout),
		lineL.Row.Set(row),
		lineL.Cells.Set(cells),
	)
}

func missingSheet(want string) error { return fmt.Errorf("missing worksheet %q", want) }

// firstSheet returns the first worksheet in workbook order.
func firstSheet(r sheetRequest) E.Either[error, string] {
	return F.Pipe3(
		r,
		requestL.Names.Get,
		A.Head[string],
		E.FromOption[string](F.Constant(errNoSheet)),
	)
}

func namedSheet(r sheetRequest) E.Either[error, string] {
	return F.Pipe2(
		r,
		requestL.Want.Get,
		E.FromPredicate(F.Bind1st(slices.Contains[[]string], requestL.Names.Get(r)), missingSheet),
	)
}

// openWorkbook closes the partial handle that excelize can return
// together with an open error.
//
// fp-go:allow-branch excelize can return a handle together with an error
func openWorkbook(path string) (*excelize.File, error) {
	book, err := excelize.OpenFile(path)
	if err != nil && book != nil {
		return nil, errors.Join(err, book.Close())
	}

	return book, err
}

// readHeader reads row 1 and returns the layout of the worksheet.
//
// fp-go:allow-branch excelize exposes only a stateful row iterator
func readHeader(rows *excelize.Rows) (Layout, error) {
	if !rows.Next() {
		return LayoutUnknown, errors.Join(errNoHeader, rows.Error())
	}
	cells, err := rows.Columns()
	if err != nil {
		return LayoutUnknown, fmt.Errorf("row %d: %w", headerRow, err)
	}

	return E.UnwrapError(decodeHeader(cells))
}

// readSheet reads one worksheet row by row. Each Next advances exactly
// one physical row, so the counter is the Excel row number. Each data
// row goes through the row pipeline addRow and into the batch.
//
// fp-go:allow-loop excelize exposes only a stateful row iterator
// fp-go:allow-branch the iterator reports an error for each row
func (s State) readSheet(file string, rows *excelize.Rows) (Batch, error) {
	layout, err := readHeader(rows)
	if err != nil {
		return Batch{}, err
	}
	ctx := stateL.Context.Get(s)
	batch := Batch{}
	for row := firstDataRow; rows.Next(); row++ {
		if err := ctx.Err(); err != nil {
			return Batch{}, err
		}
		cells, err := rows.Columns()
		if err != nil {
			return Batch{}, fmt.Errorf("row %d: %w", row, err)
		}
		batch = addRow(batch, newLine(file, layout, row, cells))
	}

	return batch, rows.Error()
}

// closeEffect closes a workbook, a row iterator, or a file.
func closeEffect[R io.Closer](r R) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) { return struct{}{}, r.Close() })
}

// joinUseError joins the error of the use step, if any, with a release
// error. IOE.Bracket otherwise returns only the release error.
func joinUseError[A any](used E.Either[error, A], releaseErr error) error {
	return errors.Join(E.ToError(used), releaseErr)
}

// closeJoined is the release step of IOE.Bracket for a closer. A close
// error never hides the error of the use step.
func closeJoined[R io.Closer, A any](r R, used E.Either[error, A]) IOE.IOEither[error, struct{}] {
	return F.Pipe2(
		r,
		closeEffect[R],
		IOE.MapLeft[struct{}](F.Bind1st(joinUseError[A], used)),
	)
}

func workbookError(path string, err error) error {
	return fmt.Errorf("workbook %s: %w", path, err)
}

// pickSheet chooses a worksheet from the --sheet option and the
// worksheet names of a workbook.
func pickSheet(want string, names []string) E.Either[error, string] {
	return F.Pipe3(
		sheetRequest{},
		requestL.Want.Set(want),
		requestL.Names.Set(names),
		chooseSheet,
	)
}

// scanSheet reads one worksheet and closes the row iterator. Excelize
// keeps the workbook path in book.Path.
func (s State) scanSheet(book *excelize.File, name string) IOE.IOEither[error, Batch] {
	return IOE.Bracket(
		F.Pipe1(name, IOE.Eitherize1(book.Rows)),
		F.Bind1st(IOE.Eitherize2(s.readSheet), filepath.Base(book.Path)),
		closeJoined[*excelize.Rows, Batch],
	)
}

// readBook chooses the worksheet of an open workbook and reads it.
func (s State) readBook(book *excelize.File) IOE.IOEither[error, Batch] {
	return F.Pipe3(
		book.GetSheetList(),
		IOE.Of[error, []string],
		IOE.ChainEitherK(F.Bind1st(pickSheet, stateSheet.Get(s))),
		IOE.Chain(F.Bind1st(s.scanSheet, book)),
	)
}

// bookBatch opens the workbook at path, reads it, and closes it on every
// path.
func (s State) bookBatch(path string) IOE.IOEither[error, Batch] {
	return IOE.Bracket(
		F.Pipe1(path, IOE.Eitherize1(openWorkbook)),
		s.readBook,
		closeJoined[*excelize.File, Batch],
	)
}

// readWorkbook checks the context, reads one workbook, and names the
// workbook in an error.
func (s State) readWorkbook(path string) IOE.IOEither[error, Batch] {
	return F.Pipe4(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Chain(F.Bind2nd(State.bookBatch, path)),
		IOE.MapLeft[Batch](F.Bind1st(workbookError, path)),
	)
}

// ReadFiles reads the selected workbooks one at a time, in order, and
// joins their batches.
func ReadFiles(s State) IOE.IOEither[error, Batch] {
	return F.Pipe3(
		s,
		stateL.Files.Get,
		IOE.TraverseArraySeq(s.readWorkbook),
		IOE.Map[error](A.Reduce(joinBatch, Batch{})),
	)
}
