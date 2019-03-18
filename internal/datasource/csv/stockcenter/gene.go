package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//StrainGene is the container for strain and gene data
type StrainGene struct {
	Id     string
	GeneId string
}

//StrainGeneReader is the defined interface for reading the strain and gene data
type StrainGeneReader interface {
	datasource.IteratorWithoutValue
	Value() (*StrainGene, error)
}

type csvStrainGeneReader struct {
	*csource.CsvReader
}

//NewCsvStrainGeneReader is to get an instance of StrainGeneReader
func NewCsvStrainGeneReader(r io.Reader) StrainGeneReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainGeneReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StrainGene instance
func (sgr *csvStrainGeneReader) Value() (*StrainGene, error) {
	sg := new(StrainGene)
	if sgr.Err != nil {
		return sg, sgr.Err
	}
	sg.Id = sgr.Record[0]
	sg.GeneId = sgr.Record[1]
	return sg, nil
}
