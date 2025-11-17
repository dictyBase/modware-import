package stockcenter

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

// Strain is the container for strain data
type Strain struct {
	ID           string
	Descriptor   string
	Summary      string
	Species      string
	User         string
	Publications []string
	Genes        []string
	CreatedOn    time.Time
	UpdatedOn    time.Time
}

// StrainReader is the defined interface for reading the strain data
type StrainReader interface {
	datasource.IteratorWithoutValue
	Value() (*Strain, error)
}

type csvStrainReader struct {
	*csource.Reader
	lookup  StockAnnotatorLookup
	plookup StockPubLookup
	glookup StockGeneLookup
}

// NewCsvStrainReader is to get an instance of strain reader
func NewCsvStrainReader(
	r io.Reader,
	al StockAnnotatorLookup,
	pl StockPubLookup,
	gl StockGeneLookup,
) StrainReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainReader{
		Reader:  &csource.Reader{Reader: cr},
		lookup:  al,
		plookup: pl,
		glookup: gl,
	}
}

// Value gets a new Strain instance
func (sr *csvStrainReader) Value() (*Strain, error) {
	s := new(Strain)
	if sr.Err != nil {
		return s, sr.Err
	}
	s.ID = sr.Record[0]
	s.Descriptor = sr.Record[1]
	s.Species = sr.Record[2]
	s.Summary = sr.Record[3]
	user, c, u, ok := sr.lookup.StockAnnotator(sr.Record[0])
	if ok {
		s.User = user
		s.CreatedOn = c
		s.UpdatedOn = u
	}
	s.Publications = append(s.Publications, sr.plookup.StockPub(s.ID)...)
	s.Genes = append(s.Genes, sr.glookup.StockGene(s.ID)...)
	return s, nil
}
