package normalize

import (
	"fmt"
	"io/fs"
	"os"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	FILE "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predfs "github.com/dictyBase/fp-go-loom/predicate/fs"
)

// liveState fails when the run context is done.
var liveState = E.FromError(contextError)

func createEntry(e fs.DirEntry) bool {
	return F.Pipe2(
		e,
		fs.DirEntry.Name,
		F.Pipe1(
			Str.HasPrefix("Phenotype_annotation_"),
			P.And(Str.HasSuffix(".xlsx")),
		),
	)
}

// Valid checks the options before any file access. A folder or an output
// path of only spaces is not a path.
func Valid(s State) E.Either[error, State] {
	return F.Pipe1(
		s,
		E.FromPredicate(
			F.Pipe1(
				F.Flow2(stateFolder.Get, nonblank),
				P.And(F.Flow2(stateOut.Get, nonblank)),
			),
			F.Constant1[State](errSpec),
		),
	)
}

func contextError(s State) error { return stateL.Context.Get(s).Err() }

func noWorkbooks(dir string, _ []string) error {
	return fmt.Errorf("no create workbooks in %s", dir)
}

// workbookPaths keeps regular files and joins their names with dir.
func workbookPaths(dir string, infos []fs.FileInfo) []string {
	return F.Pipe2(
		infos,
		A.Filter(predfs.IsRegularInfo),
		A.Map(F.Flow3(fs.FileInfo.Name, FILE.Join, En.Read(dir))),
	)
}

// folderFiles lists create workbooks directly under the folder. ReadDir
// sorts by name. DirEntry.Info returns Lstat data, so a symbolic link is
// not a regular file.
func folderFiles(s State) IOE.IOEither[error, []string] {
	dir := stateFolder.Get(s)

	return F.Pipe5(
		dir,
		IOE.Eitherize1(os.ReadDir),
		IOE.Map[error](A.Filter(createEntry)),
		IOE.Chain(IOE.TraverseArraySeq(IOE.Eitherize1(fs.DirEntry.Info))),
		IOE.Map[error](F.Bind1st(workbookPaths, dir)),
		IOE.ChainEitherK(E.FromPredicate(
			predarrays.IsNonEmpty[string](),
			F.Bind1st(noWorkbooks, dir),
		)),
	)
}

// SelectFiles returns the workbook paths in a stable order.
func SelectFiles(s State) IOE.IOEither[error, []string] {
	return F.Pipe3(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Chain(folderFiles),
	)
}
