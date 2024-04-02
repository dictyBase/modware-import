package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	ucli "github.com/urfave/cli/v2"

	E "github.com/IBM/fp-go/either"

	A "github.com/IBM/fp-go/array"
	S "github.com/IBM/fp-go/string"

	O "github.com/IBM/fp-go/option"

	F "github.com/IBM/fp-go/function"
)

var (
	Split = F.Curry2(
		F.Bind2nd[string, string, []string],
	)(strings.Split)
	flagNamesHandler = flagNames()
)

type SliceWithError[T any] struct {
	Slice []T
	Error error
}

type flagNamesReturn func([]ucli.Flag) []string

func onErrorWithSlice[T any](err error) SliceWithError[T] {
	return SliceWithError[T]{Error: err}
}

func onSuccessWithSlice[T any](slice []T) SliceWithError[T] {
	return SliceWithError[T]{Slice: slice}
}

func noDir(rec fs.DirEntry) bool { return !rec.IsDir() }

func parsePhenoFileName(file string) (time.Time, error) {
	output := F.Pipe7(
		file,
		filepath.Base,
		Split("."),
		A.Head,
		O.GetOrElse(F.Constant("")),
		Split("_"),
		A.SliceRight[string](2),
		S.Join(":"),
	)
	if len(output) == 0 {
		return time.Time{}, fmt.Errorf("error in parsing file name %s", file)
	}
	return time.Parse("Jan:02:2006", output)
}

func isPhenoAnnoFile(
	rec fs.DirEntry,
) bool {
	return F.Pipe1(rec.Name(), S.Includes("annotation"))
}

func listPhenoFiles(folder string) ([]string, error) {
	output := F.Pipe2(
		E.TryCatchError(os.ReadDir(folder)),
		E.Map[error](func(files []fs.DirEntry) []string {
			return F.Pipe3(
				files,
				A.Filter(noDir),
				A.Filter(isPhenoAnnoFile),
				A.Map(F.Curry2(fullPath)(folder)),
			)
		}),
		E.Fold[error, []string](onErrorWithSlice, onSuccessWithSlice),
	)
	return output.Slice, output.Error
}

func parseStrainFileName(file string) (time.Time, error) {
	output := F.Pipe7(
		file,
		filepath.Base,
		Split("."),
		A.Head,
		O.GetOrElse(F.Constant("")),
		Split("-"),
		A.SliceRight[string](3),
		S.Join(":"),
	)
	if len(output) == 0 {
		return time.Time{}, fmt.Errorf("error in parsing file name %s", file)
	}
	return time.Parse("Jan:02:2006", output)
}

func listStrainFiles(folder string) ([]string, error) {
	output := F.Pipe2(
		E.TryCatchError(os.ReadDir(folder)),
		E.Map[error](func(files []fs.DirEntry) []string {
			return F.Pipe3(
				files,
				A.Filter(noDir),
				A.Filter(isStrainAnnoFile),
				A.Map(F.Curry2(fullPath)(folder)),
			)
		}),
		E.Fold[error, []string](onErrorWithSlice, onSuccessWithSlice),
	)
	return output.Slice, output.Error
}

func isStrainAnnoFile(rec fs.DirEntry) bool {
	return F.Pipe1(rec.Name(), S.Includes("PMID"))
}

func flagNames() flagNamesReturn {
	return F.Flow2(
		A.Map(flgName),
		A.Map(firstName),
	)
}

func flgName(flg ucli.Flag) []string {
	return flg.Names()
}

func firstName(names []string) string {
	return names[0]
}

func fullPath(folder string, rec fs.DirEntry) string {
	return filepath.Join(folder, rec.Name())
}
