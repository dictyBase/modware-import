package gbexport

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"path"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqformat"
	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
	"github.com/dictyBase/modware-import/internal/goldenbraid/svgmap"
)

var (
	jobL          = MakefileJobLenses()
	artifactL     = MakeartifactLenses()
	parseInputL   = snapgene.MakeParseInputLenses()
	parseResultL  = snapgene.MakeParseResultLenses()
	renderInputL  = svgmap.MakeRenderInputLenses()
	jobKey        = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Key))
	jobStem       = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Stem))
	jobMember     = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Member))
	jobDuplicate  = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Duplicate))
	jobRecord     = F.Pipe1(jobL.Result, L.Compose[fileJob](parseResultL.Record))
	jobWarnings   = F.Pipe1(jobL.Result, L.Compose[fileJob](parseResultL.Warnings))
	jobBundleStem = F.Pipe1(jobL.Bundle, L.Compose[fileJob](bundleL.Stem))
	jobFASTA      = F.Pipe1(jobL.Bundle, L.Compose[fileJob](bundleL.FASTA))
	jobGenBank    = F.Pipe1(jobL.Bundle, L.Compose[fileJob](bundleL.GenBank))
	jobSVG        = F.Pipe1(jobL.Bundle, L.Compose[fileJob](bundleL.SVG))

	fastaSpec   = PR.MakePair(FastaDir, fastaExt)
	genBankSpec = PR.MakePair(GenBankDir, genBankExt)
	svgSpec     = PR.MakePair(SVGDir, svgExt)

	rejectDuplicate = E.FromPredicate(
		P.Not(jobDuplicate.Get),
		F.Constant1[fileJob](ErrDuplicateStem),
	)
	requireLiveContext = E.FromPredicate(
		F.Flow3(jobL.Ctx.Get, context.Cause, isNilError),
		F.Flow2(jobL.Ctx.Get, context.Cause),
	)
	readSource = F.Pipe1(
		isArchiveMember,
		P.Fold(F.Flow2(jobRef, readRef), readArchiveMember),
	)
	jobRef = F.Bind2nd(build[fileJob, objectRef], []func(fileJob) En.Endomorphism[objectRef]{
		F.Flow2(jobL.Input.Get, refL.Bucket.Set),
		F.Flow2(jobL.Ctx.Get, refL.Ctx.Set),
		F.Flow2(jobKey.Get, refL.Key.Set),
	})
	parseJob = F.Flow2(
		E.Of[error, fileJob],
		E.Bind(jobL.Result.Set, F.Flow2(
			F.Bind2nd(
				build[fileJob, snapgene.ParseInput],
				[]func(fileJob) En.Endomorphism[snapgene.ParseInput]{
					F.Flow2(jobKey.Get, parseInputL.SourceFile.Set),
					F.Flow2(jobStem.Get, parseInputL.Name.Set),
					F.Flow2(jobL.Data.Get, parseInputL.Data.Set),
				},
			),
			snapgene.Parse,
		)),
	)
	// encodeJob encodes all three artifacts before any write.
	encodeJob = F.Flow5(
		E.Of[error, fileJob],
		E.Map[error](F.Bind2nd(update[fileJob], []func(fileJob) En.Endomorphism[fileJob]{
			F.Flow2(jobStem.Get, jobBundleStem.Set),
		})),
		E.Bind(jobFASTA.Set, F.Flow2(jobRecord.Get, seqformat.ToFASTABytes)),
		E.Bind(jobGenBank.Set, F.Flow2(jobRecord.Get, seqformat.ToGenBankBytes)),
		E.Bind(jobSVG.Set, F.Flow2(
			F.Bind2nd(
				build[fileJob, svgmap.RenderInput],
				[]func(fileJob) En.Endomorphism[svgmap.RenderInput]{
					F.Flow2(jobRecord.Get, renderInputL.Record.Set),
					F.Flow4(
						jobL.ShowPrimers.Get,
						svgmap.WithShowPrimers,
						A.Of[svgmap.Option],
						renderInputL.Options.Set,
					),
				},
			),
			svgmap.Render,
		)),
	)
	// artifactBuilders fix the write order: FASTA, GenBank, then SVG.
	artifactBuilders = []func(fileJob) artifact{
		F.Bind2nd(build[fileJob, artifact], []func(fileJob) En.Endomorphism[artifact]{
			F.Flow2(F.Bind2nd(artifactKey, fastaSpec), artifactL.Key.Set),
			F.Flow2(jobFASTA.Get, artifactL.Data.Set),
		}),
		F.Bind2nd(build[fileJob, artifact], []func(fileJob) En.Endomorphism[artifact]{
			F.Flow2(F.Bind2nd(artifactKey, genBankSpec), artifactL.Key.Set),
			F.Flow2(jobGenBank.Get, artifactL.Data.Set),
		}),
		F.Bind2nd(build[fileJob, artifact], []func(fileJob) En.Endomorphism[artifact]{
			F.Flow2(F.Bind2nd(artifactKey, svgSpec), artifactL.Key.Set),
			F.Flow2(jobSVG.Get, artifactL.Data.Set),
		}),
	}
	writeArtifacts = F.Bind2nd(writeAll, artifactBuilders)
)

// exportFile reads, parses, encodes, and writes one source. A success is
// counted only after all three writes succeed.
func exportFile(job fileJob) IOE.IOEither[error, string] {
	return F.Pipe11(
		job,
		IOE.Of[error, fileJob],
		IOE.ChainEitherK(rejectDuplicate),
		IOE.Bind(jobL.Data.Set, readSource),
		IOE.ChainEitherK(requireLiveContext),
		IOE.ChainEitherK(parseJob),
		IOE.ChainFirstIOK[error](fileJob.logWarnings),
		IOE.ChainEitherK(encodeJob),
		IOE.ChainFirst(writeArtifacts),
		IOE.Map[error](jobKey.Get),
		IOE.MapLeft[string](F.Bind1st(wrapExportError, job)),
		IOE.TapLeft[string](job.logFailure),
	)
}

func (job fileJob) logWarnings() IO.IO[[]IO.Void] {
	return F.Pipe2(job, jobWarnings.Get, IO.TraverseArraySeq(job.logWarning))
}

func (job fileJob) logWarning(warning error) IO.IO[IO.Void] {
	return IO.FromImpure(
		func() {
			jobL.Logger.Get(job).Warn("parse warning", "key", jobKey.Get(job), "error", warning)
		},
	)
}

func (job fileJob) logFailure(err error) IOE.IOEither[error, IO.Void] {
	return F.Pipe1(
		IO.FromImpure(
			func() {
				jobL.Logger.Get(job).Error("export failed", "key", jobKey.Get(job), "error", err)
			},
		),
		IOE.FromIO[error, IO.Void],
	)
}

func wrapExportError(job fileJob, err error) error {
	return fmt.Errorf("export %q: %w", jobKey.Get(job), err)
}

func isNilError(err error) bool { return err == nil }

func isArchiveMember(job fileJob) bool { return jobMember.Get(job) != nil }

func readArchiveMember(job fileJob) IOE.IOEither[error, []byte] {
	return F.Pipe4(
		job,
		jobMember.Get,
		openMember,
		F.Bind2nd(withResource[io.ReadCloser, []byte], readAll),
		IOE.MapLeft[[]byte](ER.OnError("read archive member %q", jobKey.Get(job))),
	)
}

func openMember(member *zip.File) IOE.IOEither[error, io.ReadCloser] {
	return IOE.TryCatchError(
		func() (io.ReadCloser, error) { return member.Open() },
	)
}

func readAll(reader io.ReadCloser) IOE.IOEither[error, []byte] {
	return IOE.TryCatchError(
		func() ([]byte, error) { return io.ReadAll(reader) },
	)
}

func artifactKey(job fileJob, spec PR.Pair[string, string]) string {
	prefix, folder := jobL.Prefix.Get(job), PR.Head(spec)
	name := jobStem.Get(job) + PR.Tail(spec)
	return path.Join(prefix, folder, name)
}

func writeAll(job fileJob, builders []func(fileJob) artifact) IOE.IOEither[error, []struct{}] {
	return F.Pipe2(
		builders,
		A.Flap[artifact](job),
		IOE.TraverseArraySeq(F.Bind1st(writeArtifact, job)),
	)
}

// writeArtifact replaces the object at the artifact key. fileblob writes
// a temporary file and renames it over the key.
func writeArtifact(job fileJob, item artifact) IOE.IOEither[error, struct{}] {
	return F.Pipe1(
		IOE.TryCatchError(
			func() (struct{}, error) {
				bucket, ctx := jobL.Output.Get(job), jobL.Ctx.Get(job)
				key, data := artifactL.Key.Get(item), artifactL.Data.Get(item)
				return struct{}{}, bucket.WriteAll(ctx, key, data, nil)
			},
		),
		IOE.MapLeft[struct{}](ER.OnError("write %q", artifactL.Key.Get(item))),
	)
}
