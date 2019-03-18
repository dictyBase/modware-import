package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//PlasmidGenbank is the container for genbank link for plasmid
type PlasmidGenbank struct {
	Id      string
	Genbank string
}

//PlasmidGenbankReader is the defined interface for reading the data
type PlasmidGenbankReader interface {
	datasource.IteratorWithoutValue
	Value() (*PlasmidGenbank, error)
}

type csvPlasmidGenbankReader struct {
	*csource.CsvReader
}

//NewPlasmidReader is to get an instance of PlasmidGenbankReader
func NewPlasmidGenbankReader(r io.Reader) PlasmidGenbankReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidGenbankReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new PlasmidGenbank instance
func (pgr *csvPlasmidGenbankReader) Value() (*PlasmidGenbank, error) {
	g := new(PlasmidGenbank)
	if pgr.Err != nil {
		return g, pgr.Err
	}
	g.Id = pgr.Record[0]
	g.Genbank = pgr.Record[1]
	return g, nil
}
