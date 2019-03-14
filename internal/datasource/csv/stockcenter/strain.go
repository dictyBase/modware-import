package stockcenter

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//Strain is the container for strain data
type Strain struct {
	Id         string
	Descriptor string
	Summary    string
	Species    string
	User       string
	CreatedOn  time.Time
	UpdatedOn  time.Time
}

//StrainReader is the defined interface for reading the strain data
type StrainReader interface {
	datasource.IteratorWithoutValue
	Value() (*Strain, error)
}

type csvStrainReader struct {
	*csource.CsvReader
	lookup StrainAnnotatorLookup
}

//NewCsvStrainReader is to get an instance of strain reader
func NewCsvStrainReader(r io.Reader, al StrainAnnotatorLookup) StrainReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainReader{
		CsvReader: &csource.CsvReader{Reader: cr},
		lookup:    al,
	}
}

//Value gets a new Strain instance
func (sr *csvStrainReader) Value() (*Strain, error) {
	s := new(Strain)
	if sr.Err != nil {
		return s, sr.Err
	}
	s.Id = sr.Record[0]
	s.Descriptor = sr.Record[1]
	s.Species = sr.Record[2]
	s.Summary = sr.Record[3]
	user, c, u, ok := sr.lookup.StrainAnnotator(sr.Record[0])
	if ok {
		s.User = user
		s.CreatedOn = c
		s.UpdatedOn = u
	}
	return s, nil
}
