//package stockcenter is the data source for stockcenter and related data
package stockcenter

import (
	"bufio"
	"io"
	"strings"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
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

type tsvGenotypeReader struct {
	*tsource.TsvReader
}

//NewTsvGenotypeReader is to get an instance of GenotypeReader
func NewTsvGenotypeReader(r io.Reader) GenotypeReader {
	tr := bufio.NewScanner(r)
	return &tsvGenotypeReader{&tsource.TsvReader{Reader: tr}}
}

//Value gets a new StockOrder instance
func (gr *tsvGenotypeReader) Value() (*Genotype, error) {
	g := new(Genotype)
	if gr.Err != nil {
		return g, gr.Err
	}
	g.Id = gr.Record[1]
	g.StrainId = gr.Record[0]
	g.Genotype = strings.Replace(gr.Record[2], ", ", ",", -1)
	return g, nil
}
