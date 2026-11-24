package normalize

import (
	"encoding/csv"
	"os"
	"path/filepath"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	PR "github.com/IBM/fp-go/v2/pair"
)

const (
	filePerm       = 0o600
	createFlags    = os.O_CREATE | os.O_EXCL | os.O_WRONLY
	stagingPattern = ".pheno-*"
)

// batchFiles returns the three output files of a batch in dir.
func batchFiles(dir string, b Batch) []Output {
	return []Output{
		PR.MakePair(filepath.Join(dir, "phenotype.csv"), recordRows(batchL.Existing.Get(b))),
		PR.MakePair(filepath.Join(dir, "new-strain.csv"), recordRows(batchL.New.Get(b))),
		PR.MakePair(filepath.Join(dir, "rejects.csv"), diagnosticRows(batchL.Issues.Get(b))),
	}
}

// createExclusive creates path and fails when path exists.
func createExclusive(path string) IOE.IOEither[error, *os.File] {
	return IOE.TryCatchError(func() (*os.File, error) {
		return os.OpenFile(path, createFlags, filePerm)
	})
}

// writeCSV writes all rows to f. WriteAll flushes the writer and returns
// the flush error.
func writeCSV(f *os.File, rows [][]string) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) {
		return struct{}{}, csv.NewWriter(f).WriteAll(rows)
	})
}

// WriteFile creates the file of out exclusively, writes its rows as CSV,
// and closes the file. A close error joins a write error.
func WriteFile(out Output) IOE.IOEither[error, string] {
	return F.Pipe1(
		IOE.Bracket(
			createExclusive(PR.Head(out)),
			F.Bind2nd(writeCSV, PR.Tail(out)),
			closeJoined[*os.File, struct{}],
		),
		IOE.MapTo[error, struct{}](PR.Head(out)),
	)
}

// removeAll removes the staging directory. After a successful rename the
// directory does not exist, and RemoveAll returns nil.
func removeAll(dir string) IOE.IOEither[error, struct{}] {
	return IOE.TryCatchError(func() (struct{}, error) { return struct{}{}, os.RemoveAll(dir) })
}

// removeJoined is the release step of IOE.Bracket for the staging
// directory. A removal error never hides the error of the commit.
func removeJoined[A any](dir string, used E.Either[error, A]) IOE.IOEither[error, struct{}] {
	return F.Pipe2(dir, removeAll, IOE.MapLeft[struct{}](F.Bind1st(joinUseError[A], used)))
}

// stagingDir checks the context and creates a new directory next to the
// destination.
func stagingDir(s State) IOE.IOEither[error, string] {
	return F.Pipe4(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Map[error](F.Flow3(stateOut.Get, filepath.Clean, filepath.Dir)),
		IOE.Chain(F.Bind2nd(IOE.Eitherize2(os.MkdirTemp), stagingPattern)),
	)
}

// writeOutput checks the context, then writes one file with the writer of
// the state.
func (s State) writeOutput(out Output) IOE.IOEither[error, string] {
	return F.Pipe3(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Chain(F.Constant1[State](F.Pipe1(out, stateL.Writer.Get(s)))),
	)
}

// renameInto checks the context and renames dir to the destination. It
// never replaces an existing destination. The rename is the only commit
// point.
func renameInto(dir string, s State) IOE.IOEither[error, State] {
	return F.Pipe5(
		s,
		IOE.Of[error, State],
		IOE.ChainEitherK(liveState),
		IOE.Map[error](F.Flow2(stateOut.Get, filepath.Clean)),
		IOE.Chain(F.Bind1st(IOE.Eitherize2(renameNew), dir)),
		IOE.MapTo[error, string](s),
	)
}

// commit writes the three files of the batch into dir and renames dir to
// the destination.
func (s State) commit(dir string) IOE.IOEither[error, State] {
	return F.Pipe4(
		s,
		stateL.Batch.Get,
		F.Bind1st(batchFiles, dir),
		IOE.TraverseArraySeq(s.writeOutput),
		IOE.Chain(F.Constant1[[]string](renameInto(dir, s))),
	)
}

// Publish writes the batch into a new sibling directory and renames it to
// the destination. It removes the staging directory on every path.
func Publish(s State) IOE.IOEither[error, State] {
	return IOE.Bracket(stagingDir(s), s.commit, removeJoined[State])
}
