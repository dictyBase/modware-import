package normalize

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	PR "github.com/IBM/fp-go/v2/pair"
	"github.com/xuri/excelize/v2"
)

var outputNames = []string{"phenotype.csv", "new-strain.csv", "rejects.csv"}

func check(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func force[A any](ioe IOE.IOEither[error, A]) (A, error) { return E.UnwrapError(ioe()) }

// dataRow builds a row in the extended layout with the production aligner.
func dataRow(cells ...string) SourceRow {
	return toSourceRow(newLine("Phenotype_annotation_T.xlsx", LayoutExtended, firstDataRow, cells))
}

func fixture(t *testing.T, path string, rows ...[]string) {
	t.Helper()
	f := excelize.NewFile()
	for i, r := range rows {
		if len(r) == 0 {
			continue
		}
		axis, err := excelize.CoordinatesToCellName(1, i+1)
		check(t, err)
		check(t, f.SetSheetRow("Sheet1", axis, &r))
	}
	check(t, f.SaveAs(path))
	check(t, f.Close())
}

func readCSV(t *testing.T, path string) [][]string {
	t.Helper()
	f, err := os.Open(path)
	check(t, err)
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	check(t, err)
	check(t, f.Close())
	return records
}

func codes(o RowOutcome) []Code {
	var out []Code
	for _, d := range o.Issues {
		out = append(out, d.Code)
	}
	return out
}

func emptyDir(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	check(t, err)
	if len(entries) != 0 {
		t.Fatal(entries)
	}
}

func baseState(out string) State {
	var s State
	s.Context = context.Background()
	s.Options.Out = out
	s.Writer = WriteFile
	return s
}

func invoke(t *testing.T, bin string, want int, args ...string) string {
	t.Helper()
	return invokeIn(t, "", bin, want, args...)
}

func invokeIn(t *testing.T, dir, bin string, want int, args ...string) string {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	buf, err := cmd.CombinedOutput()
	code := 0
	if exit, ok := errors.AsType[*exec.ExitError](err); ok {
		code = exit.ExitCode()
	} else {
		check(t, err)
	}
	if code != want {
		t.Fatalf("exit=%d want=%d: %s", code, want, buf)
	}
	return string(buf)
}

// TestCSVColumns keeps each header and its field function the same width.
func TestCSVColumns(t *testing.T) {
	if len(recordHeader) != len(recordFields(Record{})) {
		t.Fatal("record header and fields differ")
	}
	if len(diagnosticHeader) != len(diagnosticFields(Diagnostic{})) {
		t.Fatal("diagnostic header and fields differ")
	}
}

func TestHeader(t *testing.T) {
	alias := slices.Clone(header10)
	alias[6] = " Envirionment ID "
	layout, err := E.UnwrapError(decodeHeader(alias))
	check(t, err)
	if layout != LayoutShort {
		t.Fatal(layout)
	}
	alias[6] = alias[7]
	if _, err = E.UnwrapError(decodeHeader(alias)); !errors.Is(err, errBadLayout) {
		t.Fatal("duplicate header accepted")
	}
}

func TestAlignShortRow(t *testing.T) {
	r := toSourceRow(newLine("a.xlsx", LayoutShort, 4, []string{
		" DBS0000001 ", "DDPHENO:0000001", "", "", "", "", "", "", "1", "a@b.org", "", "x",
	}))
	if r.Origin.Row != 4 || r.Raw.Strain != " DBS0000001 " || r.Clean.Strain != "DBS0000001" {
		t.Fatal(r)
	}
	if r.Raw.Descriptor != "" || r.Clean.Curator != "a@b.org" || r.Extra != "x" {
		t.Fatal(r)
	}
}

func TestValidTerms(t *testing.T) {
	o := classify(dataRow(
		"DBS0000001", "", "DDPHENO:0000001", "", "", "DDASSAY:0000057", "",
		"DDENVIR:0000026", "", "36982997", "", "a@b.org",
	))
	if o.Kind != KindExisting || len(o.Issues) != 0 {
		t.Fatal(o)
	}
	if o.Record.AssayID != "DDASSAY:0000057" || o.Record.EnvID != "DDENVIR:0000026" {
		t.Fatal(o.Record)
	}
}

func TestRowFaults(t *testing.T) {
	cells := []string{
		"DBS0000001", "", "DDPHENO:0000001", "", "", "", "assay label", "",
		"environment label", "PMID:0", "", "Display Name <curator@example.org>", "outside",
	}
	o := classify(dataRow(cells...))
	want := []Code{
		CodeLabelWithoutID, CodeLabelWithoutID, CodeInvalidPMID,
		CodeInvalidCurator, CodeExtraCell,
	}
	if o.Kind != KindRejected || !slices.Equal(codes(o), want) {
		t.Fatal(o)
	}
	cells = cells[:12]
	cells[5], cells[7], cells[9], cells[11] = "BAD:0000001", "DDDENVIR:xyz", "1", "c@a.org,c@b.org"
	want = []Code{CodeInvalidAssayID, CodeInvalidEnvID, CodeInvalidCurator}
	if got := codes(classify(dataRow(cells...))); !slices.Equal(got, want) {
		t.Fatal(got)
	}
	got := codes(classify(dataRow(
		"", "", "DDPHENO:0000001", "", "", "", "", "", "", "1", "", "a@b.org",
	)))
	if !slices.Equal(got, []Code{CodeMissingDescriptor}) {
		t.Fatal(got)
	}
}

func TestAccounting(t *testing.T) {
	cells := []string{
		"", "same", "DDPHENO:0000001", "", "", "DDSTRAINCHAR:0000002", "", "",
		"", "1", "", "x@y.org",
	}
	fresh := dataRow(cells...)
	cells[0] = "DBS0000001"
	known := dataRow(cells...)
	var b Batch
	for _, r := range []SourceRow{fresh, fresh, known, {}} {
		b = joinBatch(b, outcomeBatch(classify(r)))
	}
	got := []int{len(b.Existing), len(b.New), b.Blank, b.Warnings, b.Rejected}
	if !slices.Equal(got, []int{1, 2, 1, 3, 0}) {
		t.Fatal(b)
	}
	if b.New[0].Descriptor != b.New[1].Descriptor || b.New[0].AssayID != "" {
		t.Fatal(b.New)
	}
}

func TestSelectFiles(t *testing.T) {
	dir := t.TempDir()
	s := stateFolder.Set(dir)(baseState(""))
	if _, err := force(SelectFiles(s)); err == nil {
		t.Fatal("empty selection accepted")
	}
	good := filepath.Join(dir, "Phenotype_annotation_B.xlsx")
	fixture(t, good, header10)
	update := filepath.Join(dir, "Phenotypes_update_JUN_11_2023.xlsx")
	check(t, os.WriteFile(update, nil, 0o600))
	check(t, os.WriteFile(filepath.Join(dir, "~$Phenotype_annotation_B.xlsx"), nil, 0o600))
	link := filepath.Join(dir, "Phenotype_annotation_A.xlsx")
	check(t, os.Symlink(good, link))
	files, err := force(SelectFiles(s))
	check(t, err)
	if !slices.Equal(files, []string{good}) {
		t.Fatal(files)
	}
}

func TestValidOptions(t *testing.T) {
	if E.IsLeft(Valid(stateFolder.Set("src")(baseState("out")))) {
		t.Fatal("valid options rejected")
	}
	for _, folder := range []string{"", "   ", "\t\u00a0"} {
		if !errors.Is(leftOf(Valid(stateFolder.Set(folder)(baseState("out")))), errSpec) {
			t.Fatalf("blank folder %q accepted", folder)
		}
	}
	if !errors.Is(leftOf(Valid(stateFolder.Set("src")(baseState(" ")))), errSpec) {
		t.Fatal("blank out accepted")
	}
}

func leftOf(e E.Either[error, State]) error {
	_, err := E.UnwrapError(e)
	return err
}

func TestWorkbookFailures(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Phenotype_annotation_A.xlsx")
	s := baseState("")
	check(t, os.WriteFile(path, []byte("not a zip archive"), 0o600))
	if _, err := force(s.readWorkbook(path)); err == nil {
		t.Fatal("corrupt archive accepted")
	}
	fixture(t, path)
	if _, err := force(s.readWorkbook(path)); !errors.Is(err, errNoHeader) {
		t.Fatal("missing header accepted")
	}
	f := excelize.NewFile()
	_, err := f.NewSheet("Curated")
	check(t, err)
	check(t, f.SetSheetRow("Curated", "A1", &header10))
	row := []string{"DBS0000001"}
	check(t, f.SetSheetRow("Curated", "A4", &row))
	check(t, f.SaveAs(path))
	check(t, f.Close())
	// Without --sheet the first worksheet (the empty Sheet1) is read.
	if _, err = force(s.readWorkbook(path)); !errors.Is(err, errNoHeader) {
		t.Fatal("first worksheet not used", err)
	}
	if _, err = force(stateSheet.Set("absent")(s).readWorkbook(path)); err == nil {
		t.Fatal("missing sheet accepted")
	}
	b, err := force(stateSheet.Set("Curated")(s).readWorkbook(path))
	check(t, err)
	if b.Blank != 2 || b.Rejected != 1 || b.Issues[0].Origin.Row != 4 {
		t.Fatal(b)
	}
}

func failWriter(Output) IOE.IOEither[error, string] { return IOE.Left[string](io.ErrClosedPipe) }

func TestPublish(t *testing.T) {
	d := t.TempDir()
	dst := filepath.Join(d, "out")
	_, err := force(Publish(stateL.Writer.Set(failWriter)(baseState(dst))))
	if !errors.Is(err, io.ErrClosedPipe) {
		t.Fatal(err)
	}
	emptyDir(t, d)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = force(Publish(stateL.Context.Set(ctx)(baseState(dst))))
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	emptyDir(t, d)
	_, err = force(Publish(baseState(dst)))
	check(t, err)
	for _, name := range outputNames {
		if len(readCSV(t, filepath.Join(dst, name))) != 1 {
			t.Fatal(name)
		}
	}
	if _, err = force(Publish(baseState(dst))); err == nil {
		t.Fatal("replaced an existing destination")
	}
}

func TestCancelDuringPublish(t *testing.T) {
	d := t.TempDir()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cancelAfter := func(o Output) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			defer cancel()
			return WriteFile(o)()
		}
	}
	s := baseState(filepath.Join(d, "out"))
	s.Context = ctx
	s.Writer = cancelAfter
	if _, err := force(Publish(s)); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	emptyDir(t, d)
}

type badCloser struct{}

func (badCloser) Close() error { return io.ErrClosedPipe }

func failUse(badCloser) IOE.IOEither[error, int] { return IOE.Left[int](io.ErrUnexpectedEOF) }

func TestReleaseJoinsErrors(t *testing.T) {
	_, err := force(IOE.Bracket(IOE.Of[error](badCloser{}), failUse, closeJoined[badCloser, int]))
	if !errors.Is(err, io.ErrUnexpectedEOF) || !errors.Is(err, io.ErrClosedPipe) {
		t.Fatal(err)
	}
	dir := t.TempDir()
	_, err = force(IOE.Bracket(IOE.Of[error](dir), failDir, removeJoined[int]))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal(err)
	}
	if _, statErr := os.Stat(dir); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatal(statErr)
	}
}

func failDir(string) IOE.IOEither[error, int] { return IOE.Left[int](io.ErrUnexpectedEOF) }

func TestWriteFileExclusive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.csv")
	out := PR.MakePair(path, [][]string{{"a"}})
	if _, err := force(WriteFile(out)); err != nil {
		t.Fatal(err)
	}
	if _, err := force(WriteFile(out)); !errors.Is(err, os.ErrExist) {
		t.Fatal(err)
	}
}

func buildBinary(t *testing.T, dir string) string {
	t.Helper()
	bin := filepath.Join(dir, "phenoannotation")
	out, err := exec.Command("go", "build", "-o", bin, "../../../cmd/phenoannotation").
		CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v: %s", err, out)
	}
	return bin
}

const processNote = "first\t\"quoted\"\nlast"

// writeSources writes one short and one extended workbook and an update file.
func writeSources(t *testing.T, src string) (string, string) {
	t.Helper()
	a := filepath.Join(src, "Phenotype_annotation_A.xlsx")
	z := filepath.Join(src, "Phenotype_annotation_Z.xlsx")
	valid := []string{
		"DBS0000001", "DDPHENO:0000001", "", processNote, "DDSTRAINCHAR:0000002", "",
		"DDDENVIR:0000001", "", "PMID:123", " a@b.org\u00a0",
	}
	fixture(t, z, header10, valid, nil, valid)
	fresh := []string{"", "same", "DDPHENO:0000001", "", "", "", "", "", "", "42", "", "c@x.org"}
	bad := slices.Clone(fresh)
	bad[0], bad[10], bad[11] = "DBS1", "x", ""
	fixture(t, a, header12, fresh, fresh, bad)
	update := filepath.Join(src, "Phenotypes_update_JUN_11_2023.xlsx")
	check(t, os.WriteFile(update, []byte("ignore"), 0o600))
	return a, z
}

func checkClean(t *testing.T, dir string) {
	t.Helper()
	p := readCSV(t, filepath.Join(dir, "phenotype.csv"))
	header := recordHeader
	want := []string{
		"Phenotype_annotation_Z.xlsx", "2", "DBS0000001", "", "DDPHENO:0000001", "",
		"DDENVIR:0000001", processNote, "123", "a@b.org",
	}
	if len(p) != 3 || !slices.Equal(p[0], header) || !slices.Equal(p[1], want) || p[2][1] != "4" {
		t.Fatal(p)
	}
}

func checkFolder(t *testing.T, first, second string) {
	t.Helper()
	for _, name := range outputNames {
		x, err := os.ReadFile(filepath.Join(first, name))
		check(t, err)
		y, err := os.ReadFile(filepath.Join(second, name))
		check(t, err)
		if string(x) != string(y) {
			t.Fatal(name)
		}
	}
	n := readCSV(t, filepath.Join(first, "new-strain.csv"))
	if len(n) != 3 || n[1][3] != n[2][3] || n[2][1] != "3" {
		t.Fatal(n)
	}
	var got []string
	for _, d := range readCSV(t, filepath.Join(first, "rejects.csv"))[1:] {
		got = append(got, d[3])
	}
	want := []string{
		"invalid_strain_id", "deleted_in_create", "invalid_assigned_by",
		"unsupported_strain_characteristic", "corrected_env_prefix",
		"unsupported_strain_characteristic", "corrected_env_prefix",
	}
	if !slices.Equal(got, want) {
		t.Fatal(got)
	}
}

func TestProcess(t *testing.T) {
	d := t.TempDir()
	src := filepath.Join(d, "src")
	check(t, os.Mkdir(src, 0o700))
	bin := buildBinary(t, d)
	a, z := writeSources(t, src)

	only := filepath.Join(d, "only")
	check(t, os.Mkdir(only, 0o700))
	check(t, os.Rename(z, filepath.Join(only, filepath.Base(z))))
	clean := filepath.Join(d, "clean")
	invoke(t, bin, 0, "normalize", "--folder", only, "--out", clean)
	checkClean(t, clean)
	check(t, os.Rename(filepath.Join(only, filepath.Base(z)), z))

	// Without --out the batch goes to pheno_normalize in the work dir.
	work := filepath.Join(d, "work")
	check(t, os.Mkdir(work, 0o700))
	invokeIn(t, work, bin, 2, "normalize", "--folder", src)
	checkFolder(t, filepath.Join(work, "pheno_normalize"), filepath.Join(work, "pheno_normalize"))

	first, second := filepath.Join(d, "first"), filepath.Join(d, "second")
	for _, dest := range []string{first, second} {
		msg := invoke(t, bin, 2, "normalize", "--folder", src, "--out", dest)
		if !strings.Contains(msg, "files=2 blank=1 existing=2 new=2 rejected=1 warnings=4") {
			t.Fatal(msg)
		}
	}
	checkFolder(t, first, second)

	unused := filepath.Join(d, "unused")
	invoke(t, bin, 1, "normalize", "--folder", src, "--out", first)
	invoke(t, bin, 1, "normalize", "--out", unused)
	invoke(t, bin, 1, "normalize", "--folder", "  ", "--out", unused)
	invoke(t, bin, 1, "normalize", "--folder", src, "--sheet", "absent", "--out", unused)
	invoke(t, bin, 1, "normalize", "--input", a, "--folder", src, "--out", unused)
	invoke(t, bin, 1, "normalize", "--typo")
	invoke(t, bin, 1)
	if _, err := os.Stat(unused); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}
