package stockcenter

import (
	"bufio"
	"io"
	"strings"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
)

//Phenotype is the container for strain phenotype
type Phenotype struct {
	StrainId     string
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

type tsvPhenotypeReader struct {
	*tsource.TsvReader
}

//NewPhenotypeReader is to get an instance of PhenotypeReader
func NewPhenotypeReader(r io.Reader) PhenotypeReader {
	cr := bufio.NewScanner(r)
	return &tsvPhenotypeReader{&tsource.TsvReader{Reader: cr}}
}

//Value gets a new Phenotype instance
func (cp *tsvPhenotypeReader) Value() (*Phenotype, error) {
	p := new(Phenotype)
	if cp.Err != nil {
		return p, cp.Err
	}
	p.StrainId = cp.Record[0]
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
