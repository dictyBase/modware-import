//package stockcenter is the data source for stockcenter and related data
package stockcenter

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//Genotype is the container for strain genotype
type Genotype struct {
	Id       string
	StrainId string
	Genotype string
}

//GenotypeReader is the defined interface for reading the data
type GenotypeReader interface {
	datasource.IteratorWithoutValue
	Value() (*Genotype, error)
}

type csvGenotypeReader struct {
	*csource.CsvReader
}

//NewCsvGenotypeReader is to get an instance of GenotypeReader
func NewCsvGenotypeReader(r io.Reader) GenotypeReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvGenotypeReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StockOrder instance
func (gr *csvGenotypeReader) Value() (*Genotype, error) {
	g := new(Genotype)
	if gr.Err != nil {
		return g, gr.Err
	}
	g.Id = gr.Record[1]
	g.StrainId = gr.Record[0]
	g.Genotype = strings.Replace(gr.Record[2], ", ", ",", -1)
	return g, nil
}
