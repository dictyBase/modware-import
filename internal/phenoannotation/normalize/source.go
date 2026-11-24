package normalize

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	P "github.com/IBM/fp-go/v2/predicate"
	Str "github.com/IBM/fp-go/v2/string"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predfs "github.com/dictyBase/fp-go-loom/predicate/fs"
	predstrings "github.com/dictyBase/fp-go-loom/predicate/strings"
)

var (
	hasCreatePrefix = F.Bind2nd(strings.HasPrefix, "Phenotype_annotation_")
	hasXlsxSuffix   = predstrings.HasSuffix(".xlsx")
	createName      = F.Pipe1(hasCreatePrefix, P.And(hasXlsxSuffix))
	createEntry     = F.Pipe1(createName, P.ContraMap(fs.DirEntry.Name))
	createFile      = F.Pipe1(
		predfs.IsRegularInfo,
		P.And(F.Pipe1(createName, P.ContraMap(fs.FileInfo.Name))),
	)

	hasInput   = F.Flow2(stateInput.Get, Str.IsNonEmpty)
	hasFolder  = F.Flow2(stateFolder.Get, Str.IsNonEmpty)
	hasOut     = F.Flow2(stateOut.Get, nonblank)
	inputOnly  = F.Pipe1(hasInput, P.And(P.Not(hasFolder)))
	folderOnly = F.Pipe1(hasFolder, P.And(P.Not(hasInput)))
	validSpec  = F.Pipe2(inputOnly, P.Or(folderOnly), P.And(hasOut))

	// Valid checks the option combination before any file access.
	Valid = E.FromPredicate(validSpec, F.Constant1[State](errSpec))

	// liveState fails when the run context is done.
	liveState = E.FromError(contextError)

	// selectSource reads one explicit file or lists one folder.
	selectSource = F.Pipe1(hasInput, P.Fold(folderFiles, explicitFile))
)

func contextError(s State) error { return stateL.Context.Get(s).Err() }

func notCreateFile(info fs.FileInfo) error {
	return fmt.Errorf("not a regular create workbook: %s", info.Name())
}

func noWorkbooks(dir string, _ []string) error {
	return fmt.Errorf("no create workbooks in %s", dir)
}

// workbookPath returns the path of a listed workbook in dir.
func workbookPath(dir string, info fs.FileInfo) string {
	return filepath.Join(dir, info.Name())
}

// lstatCreateFile checks one path without following a symbolic link.
func lstatCreateFile(path string) IOE.IOEither[error, fs.FileInfo] {
	return F.Pipe2(
		path,
		IOE.Eitherize1(os.Lstat),
		IOE.ChainEitherK(E.FromPredicate(createFile, notCreateFile)),
	)
}

// explicitFile accepts one regular create workbook.
func explicitFile(s State) IOE.IOEither[error, []string] {
	return F.Pipe4(
		s,
		stateInput.Get,
		IOE.Of[error, string],
		IOE.ChainFirst(lstatCreateFile),
		IOE.Map[error](A.Of[string]),
	)
}

// folderFiles lists create workbooks directly under the folder. ReadDir
// sorts by name. DirEntry.Info returns Lstat data, so a symbolic link is
// not a regular file.
func folderFiles(s State) IOE.IOEither[error, []string] {
	dir := stateFolder.Get(s)

	return F.Pipe6(
		dir,
		IOE.Eitherize1(os.ReadDir),
		IOE.Map[error](A.Filter(createEntry)),
		IOE.Chain(IOE.TraverseArraySeq(IOE.Eitherize1(fs.DirEntry.Info))),
		IOE.Map[error](A.Filter(predfs.IsRegularInfo)),
		IOE.Map[error](A.Map(F.Bind1st(workbookPath, dir))),
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
		IOE.Chain(selectSource),
	)
}
