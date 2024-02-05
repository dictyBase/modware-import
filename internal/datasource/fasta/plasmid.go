package fasta

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase/modware-import/internal/datasource"
	"github.com/dictybase-playground/gobio/seqio"
)

type PlasmidSeq struct {
	Id  string
	Seq string
}

type PlasmidSeqReader interface {
	datasource.IteratorWithoutValue
	Value() (*PlasmidSeq, error)
}

type dirSeqReader struct {
	fasFiles []fs.FileInfo
	idx      int
	dir      string
	currFile string
}

func NewPlasmidSeqReader(dir string) (PlasmidSeqReader, error) {
	info, err := os.ReadDir(dir)
	if err != nil {
		return &dirSeqReader{}, fmt.Errorf(
			"error in reading dir %s %s",
			dir,
			err,
		)
	}
	var files []fs.FileInfo
	for _, entry := range info {
		if entry.IsDir() {
			continue
		}
		finfo, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if strings.HasSuffix(entry.Name(), "fasta") {
			files = append(files, finfo)
		}
	}
	return &dirSeqReader{fasFiles: files, idx: 0, dir: dir}, nil
}

func (seqr *dirSeqReader) Next() bool {
	if seqr.idx == len(seqr.fasFiles) {
		return false
	}
	seqr.currFile = filepath.Join(seqr.dir, seqr.fasFiles[seqr.idx].Name())
	seqr.idx += 1
	return true
}

func (seqr *dirSeqReader) Value() (*PlasmidSeq, error) {
	pseq := new(PlasmidSeq)
	reader, err := seqio.NewFastaReader(seqr.currFile)
	if err != nil {
		return pseq, err
	}
	if !reader.HasEntry() {
		if err := reader.Err(); err != nil {
			return pseq, err
		}
		return pseq, fmt.Errorf(
			"error in reading fasta entry from file %s",
			seqr.currFile,
		)
	}
	pseq.Id = strings.Split(filepath.Base(seqr.currFile), ".")[0]
	pseq.Seq = string(reader.NextEntry().Sequence)
	return pseq, nil
}
