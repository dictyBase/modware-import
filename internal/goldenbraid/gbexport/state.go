// Package gbexport exports GoldenBraid assemblies as FASTA, GenBank, and SVG.
package gbexport

import (
	"archive/zip"
	"context"
	"errors"
	"log/slog"

	"gocloud.dev/blob"

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
	typeDNA           = "dna"
	typeZip           = "zip"
	typeDir           = "dir"
	schemeFile        = "file"
	schemeS3          = "s3"
	schemeSeparator   = "://"
	keySeparator      = "/"
	macMetadataDir    = "__MACOSX"
	appleDoublePrefix = "._"
	maxFolderDepth    = 1
)

var (
	// ErrNoDNAFiles identifies an input without selectable .dna files.
	ErrNoDNAFiles = errors.New("no SnapGene .dna files found")
	// ErrDuplicateStem identifies two sources with the same file stem. Their
	// artifacts would overwrite each other.
	ErrDuplicateStem = errors.New("duplicate assembly stem")
	// ErrUnknownInputType identifies an input kind other than dna, zip, or dir.
	ErrUnknownInputType = errors.New("unknown input type")
	// ErrUnsupportedScheme identifies a location scheme without a bucket driver.
	ErrUnsupportedScheme = errors.New("unsupported location scheme")
)

// Params holds the export locations and the storage settings.
// fp-go:Lens
type Params struct {
	Input       string
	Output      string
	S3Endpoint  string
	S3Region    string
	AccessKey   string
	SecretKey   string
	ShowPrimers bool
}

// State carries the dependencies of one export run and, after a
// successful run, the number of exported sources.
// fp-go:Lens
type State struct {
	Ctx      context.Context
	Logger   *slog.Logger
	Params   Params
	Exported int
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
	Key    string
	Stem   string
	Member *zip.File
}

// runState carries resolved locations, open buckets, and selected sources.
// Pipeline returns only State, never a bucket handle.
// fp-go:Lens
type runState struct {
	State        State
	Input        location
	Output       location
	InputBucket  *blob.Bucket
	OutputBucket *blob.Bucket
	Sources      []source
}

// fileJob carries one source through read, parse, encode, and write.
// fp-go:Lens
type fileJob struct {
	Run     runState
	Source  source
	Data    []byte
	Result  snapgene.ParseResult
	FASTA   []byte
	GenBank []byte
	SVG     []byte
}
