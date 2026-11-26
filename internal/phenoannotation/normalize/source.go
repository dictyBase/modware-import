package normalize

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	FILE "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predfs "github.com/dictyBase/fp-go-loom/predicate/fs"
	predstrings "github.com/dictyBase/fp-go-loom/predicate/strings"
)

var (
	hasCreatePrefix = F.Bind2nd(strings.HasPrefix, "Phenotype_annotation_")
	hasXlsxSuffix   = predstrings.HasSuffix(".xlsx")
	createName      = F.Pipe1(hasCreatePrefix, P.And(hasXlsxSuffix))
	createEntry     = F.Pipe1(createName, P.ContraMap(fs.DirEntry.Name))

	// liveState fails when the run context is done.
	liveState = E.FromError(contextError)
)

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

// folderFiles lists create workbooks directly under the folder. ReadDir
// sorts by name. DirEntry.Info returns Lstat data, so a symbolic link is
// not a regular file.
func folderFiles(s State) IOE.IOEither[error, []string] {
	dir := stateFolder.Get(s)

	return F.Pipe7(
		dir,
		IOE.Eitherize1(os.ReadDir),
		IOE.Map[error](A.Filter(createEntry)),
		IOE.Chain(IOE.TraverseArraySeq(IOE.Eitherize1(fs.DirEntry.Info))),
		IOE.Map[error](A.Filter(predfs.IsRegularInfo)),
		IOE.Map[error](A.Map(F.Flow2(fs.FileInfo.Name, FILE.Join))),
		IOE.Map[error](A.Ap[string](A.Of(dir))),
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
