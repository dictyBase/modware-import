// Package gbexport exports GoldenBraid assemblies as FASTA, GenBank, and SVG.
package gbexport

import (
	"archive/zip"
	"context"
	"errors"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"

	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
	"github.com/dictyBase/modware-import/internal/goldenbraid/snapgene"
)

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

const (
	// FastaDir holds FASTA artifacts under the output root.
	FastaDir = "fasta"
	// GenBankDir holds GenBank artifacts under the output root.
	GenBankDir = "genbank"
	// SVGDir holds circular map artifacts under the output root.
	SVGDir = "svg"

	fastaExt          = ".fasta"
	genBankExt        = ".gbk"
	svgExt            = ".svg"
	dnaExt            = ".dna"
	zipExt            = ".zip"
	typeAuto          = "auto"
	typeDNA           = "dna"
	typeZip           = "zip"
	typeDir           = "dir"
	schemeFile        = "file"
	schemeS3          = "s3"
	keySeparator      = "/"
	macMetadataDir    = "__MACOSX"
	appleDoublePrefix = "._"
	maxFolderDepth    = 1
)

var (
	// ErrNoDNAFiles identifies an input without selectable .dna files.
	ErrNoDNAFiles = errors.New("no SnapGene .dna files found")
	// ErrPartialExportFailure identifies a batch with at least one failed file.
	ErrPartialExportFailure = errors.New("one or more assemblies failed to export")
	// ErrDuplicateStem identifies a later source whose stem is already reserved.
	ErrDuplicateStem = errors.New("duplicate assembly stem")
	// ErrUnknownInputType identifies an input type other than auto, dna, zip, or dir.
	ErrUnknownInputType = errors.New("unknown input type")
	// ErrUnsupportedScheme identifies a location scheme without a bucket driver.
	ErrUnsupportedScheme = errors.New("unsupported location scheme")
)

// Params holds the export locations and the storage settings.
// fp-go:Lens
type Params struct {
	Input       string
	InputType   string
	Output      string
	S3Endpoint  string
	S3Region    string
	AccessKey   string
	SecretKey   string
	ShowPrimers bool
}

// State carries the dependencies and the batch summary of one export run.
// fp-go:Lens
type State struct {
	Ctx     context.Context
	Logger  *slog.Logger
	Params  Params
	Summary seqmodel.ExportSummary
}

// FileBundle holds the three encoded artifacts of one construct.
// fp-go:Lens
type FileBundle struct {
	Stem    string
	FASTA   []byte
	GenBank []byte
	SVG     []byte
}

// location is a bucket root (a local directory or an S3 bucket) and a key.
// fp-go:Lens
type location struct {
	Scheme string
	Kind   string
	Bucket string
	Key    string
}

// source is one selected construct: a bucket object or an archive member.
// fp-go:Lens
type source struct {
	Key       string
	Stem      string
	Member    *zip.File
	Duplicate bool
}

// runState holds the resolved locations and, inside the resource scopes,
// the open buckets. Pipeline returns only State, never a bucket handle.
// fp-go:Lens
type runState struct {
	State        State
	Input        location
	Output       location
	Sources      []source
	InputBucket  *blob.Bucket
	OutputBucket *blob.Bucket
}

// bucketTarget describes how to open one side of the export.
// fp-go:Lens
type bucketTarget struct {
	Ctx     context.Context
	Params  Params
	Scheme  string
	Name    string
	Options fileblob.Options
	Client  *s3.Client
}

// objectRef addresses one object or one list prefix in an open bucket.
// fp-go:Lens
type objectRef struct {
	Bucket *blob.Bucket
	Ctx    context.Context
	Key    string
}

// fp-go:Lens
type fileJob struct {
	Ctx         context.Context
	Logger      *slog.Logger
	ShowPrimers bool
	Input       *blob.Bucket
	Output      *blob.Bucket
	Prefix      string
	Source      source
	Data        []byte
	Result      snapgene.ParseResult
	Bundle      FileBundle
}

// fp-go:Lens
type artifact struct {
	Key  string
	Data []byte
}
