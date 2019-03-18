package stockcenter

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//Phenotype is the container for strain phenotype
type Phenotype struct {
	Id           string
	Observation  string
	Environment  string
	Assay        string
	Note         string
	LiteratureId string
}

//PhenotypeReader is the defined interface for reading the data
type PhenotypeReader interface {
	datasource.IteratorWithoutValue
	Value() (*Phenotype, error)
}

type csvPhenotypeReader struct {
	*csource.CsvReader
}

//NewPhenotypeReader is to get an instance of PhenotypeReader
func NewPhenotypeReader(r io.Reader) PhenotypeReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPhenotypeReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new Phenotype instance
func (cp *csvPhenotypeReader) Value() (*Phenotype, error) {
	p := new(Phenotype)
	if cp.Err != nil {
		return p, cp.Err
	}
	p.Id = cp.Record[0]
	p.Observation = cp.Record[1]
	if len(cp.Record[2]) > 0 {
		p.Environment = cp.Record[2]
	}
	if len(cp.Record[3]) > 0 {
		p.Assay = cp.Record[3]
	}
	p.LiteratureId = cp.Record[4]
	if !strings.HasPrefix(cp.Record[5], "[") {
		p.Note = cp.Record[5]
	}
	return p, nil
}
