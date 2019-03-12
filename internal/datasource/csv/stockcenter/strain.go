package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//Strain is the container for strain data
type Strain struct {
	Id         string
	Descriptor string
	Summary    string
	Species    string
}

//StrainReader is the defined interface for reading the strain data
type StrainReader interface {
	datasource.IteratorWithoutValue
	Value() (*Strain, error)
}

type csvStrainReader struct {
	*csource.CsvReader
}

//NewCsvStrainReader is to get an instance of strain reader
func NewCsvStrainReader(r io.Reader) StrainReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainReader{&csource.CsvReader{Reader: cr}}
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
	return s, nil
}
