package gbexport

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"path"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	P "github.com/IBM/fp-go/v2/predicate"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqformat"
	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
	"github.com/dictyBase/modware-import/internal/goldenbraid/svgmap"
)

var (
	jobL         = MakefileJobLenses()
	parseInputL  = snapgene.MakeParseInputLenses()
	parseResultL = snapgene.MakeParseResultLenses()
	renderInputL = svgmap.MakeRenderInputLenses()

	jobCtx      = F.Pipe1(jobL.Run, L.Compose[fileJob](runCtx))
	jobLogger   = F.Pipe1(jobL.Run, L.Compose[fileJob](runLogger))
	jobPrimers  = F.Pipe1(jobL.Run, L.Compose[fileJob](runPrimers))
	jobOutBlob  = F.Pipe1(jobL.Run, L.Compose[fileJob](runL.OutputBucket))
	jobPrefix   = F.Pipe1(jobL.Run, L.Compose[fileJob](runOutputKey))
	jobKey      = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Key))
	jobStem     = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Stem))
	jobMember   = F.Pipe1(jobL.Source, L.Compose[fileJob](sourceL.Member))
	jobRecord   = F.Pipe1(jobL.Result, L.Compose[fileJob](parseResultL.Record))
	jobWarnings = F.Pipe1(jobL.Result, L.Compose[fileJob](parseResultL.Warnings))

	isArchiveMember = F.Flow2(jobMember.Get, P.IsNonZero[*zip.File]())
	// readSource reads an archive member from memory and any other
	// source from the input bucket.
	readSource = F.Pipe1(isArchiveMember, P.Fold(readObject, readArchiveMember))
)

// convertSource reads one source, parses it, encodes FASTA, GenBank,
// and SVG artifacts, and writes them sequentially.
func convertSource(job fileJob) IOE.IOEither[error, struct{}] {
	return F.Pipe9(
		job,
		IOE.Of[error, fileJob],
		IOE.Bind(jobL.Data.Set, readSource),
		IOE.ChainEitherK(requireLiveContext),
		IOE.ChainEitherK(parseJob),
		IOE.ChainFirstIOK[error](fileJob.logWarnings),
		IOE.ChainEitherK(encodeJob),
		IOE.Chain(writeArtifacts),
		IOE.MapTo[error, fileJob](struct{}{}),
		IOE.MapLeft[struct{}](F.Bind1st(wrapSourceError, job)),
	)
}

func requireLiveContext(job fileJob) E.Either[error, fileJob] {
	ctx := jobCtx.Get(job)
	return F.Pipe2(
		context.Cause(ctx),
		E.FromPredicate(P.IsZero[error](), F.Identity[error]),
		E.MapTo[error, error](job),
	)
}

// readObject reads the source key of job from the input bucket.
func readObject(job fileJob) IOE.IOEither[error, []byte] {
	run := jobL.Run.Get(job)
	key := jobKey.Get(job)
	return readBucketObject(run, key)
}

// parseJob parses the SnapGene data of job and stores the result on job.
func parseJob(job fileJob) E.Either[error, fileJob] {
	return F.Pipe2(job, E.Of[error, fileJob], E.Bind(jobL.Result.Set, parseSnapGene))
}

func parseSnapGene(job fileJob) E.Either[error, snapgene.ParseResult] {
	key := jobKey.Get(job)
	stem := jobStem.Get(job)
	data := jobL.Data.Get(job)
	return F.Pipe4(
		snapgene.ParseInput{},
		parseInputL.SourceFile.Set(key),
		parseInputL.Name.Set(stem),
		parseInputL.Data.Set(data),
		snapgene.Parse,
	)
}

// encodeJob encodes FASTA, GenBank, and SVG bytes before any write.
func encodeJob(job fileJob) E.Either[error, fileJob] {
	return F.Pipe4(
		job,
		E.Of[error, fileJob],
		E.Bind(jobL.FASTA.Set, F.Flow2(jobRecord.Get, seqformat.ToFASTABytes)),
		E.Bind(jobL.GenBank.Set, F.Flow2(jobRecord.Get, seqformat.ToGenBankBytes)),
		E.Bind(jobL.SVG.Set, renderSVG),
	)
}

func renderSVG(job fileJob) E.Either[error, []byte] {
	record := jobRecord.Get(job)
	showPrimers := jobPrimers.Get(job)
	opt := svgmap.WithShowPrimers(showPrimers)
	opts := A.Of(opt)
	return F.Pipe3(
		svgmap.RenderInput{},
		renderInputL.Record.Set(record),
		renderInputL.Options.Set(opts),
		svgmap.Render,
	)
}

// writeArtifacts writes FASTA, GenBank, and SVG in order and stops on the
// first write error.
func writeArtifacts(job fileJob) IOE.IOEither[error, fileJob] {
	return F.Pipe4(
		job,
		IOE.Of[error, fileJob],
		IOE.ChainFirst(writeFASTA),
		IOE.ChainFirst(writeGenBank),
		IOE.ChainFirst(writeSVG),
	)
}

func writeFASTA(job fileJob) IOE.IOEither[error, struct{}] {
	data := jobL.FASTA.Get(job)
	return writeOutput(job, FastaDir, fastaExt, data)
}

func writeGenBank(job fileJob) IOE.IOEither[error, struct{}] {
	data := jobL.GenBank.Get(job)
	return writeOutput(job, GenBankDir, genBankExt, data)
}

func writeSVG(job fileJob) IOE.IOEither[error, struct{}] {
	data := jobL.SVG.Get(job)
	return writeOutput(job, SVGDir, svgExt, data)
}

// writeOutput writes one artifact object to the output bucket.
func writeOutput(
	job fileJob,
	folder string,
	ext string,
	data []byte,
) IOE.IOEither[error, struct{}] {
	bucket, ctx := jobOutBlob.Get(job), jobCtx.Get(job)
	prefix, stem := jobPrefix.Get(job), jobStem.Get(job)
	key := path.Join(prefix, folder, stem+ext)
	return F.Pipe1(
		IOE.TryCatchError(
			func() (struct{}, error) {
				return struct{}{}, bucket.WriteAll(ctx, key, data, nil)
			},
		),
		IOE.MapLeft[struct{}](ER.OnError("write %q", key)),
	)
}

func (job fileJob) logWarnings() IO.IO[[]IO.Void] {
	return F.Pipe2(job, jobWarnings.Get, IO.TraverseArraySeq(job.logWarning))
}

func (job fileJob) logWarning(warning error) IO.IO[IO.Void] {
	return IO.FromImpure(
		func() {
			jobLogger.Get(job).Warn("parse warning", "key", jobKey.Get(job), "error", warning)
		},
	)
}

// wrapSourceError names the source key of job in err.
func wrapSourceError(job fileJob, err error) error {
	return fmt.Errorf("convert %q: %w", jobKey.Get(job), err)
}

// readArchiveMember reads the archive member of job.
func readArchiveMember(job fileJob) IOE.IOEither[error, []byte] {
	return F.Pipe2(job, jobMember.Get, readMember)
}

// readMember opens member, reads all of it, and closes it on success
// and on failure.
func readMember(member *zip.File) IOE.IOEither[error, []byte] {
	return F.Pipe2(
		readAll,
		IOE.WithResource[[]byte](openMember(member), closeResource[io.ReadCloser]),
		IOE.MapLeft[[]byte](ER.OnError("read archive member %q", memberName(member))),
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
