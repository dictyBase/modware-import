package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//Inventory is the container for strain inventory
type Plasmid struct {
	Id      string
	Plasmid string
}

//PlasmidReader is the defined interface for reading the data
type PlasmidReader interface {
	datasource.IteratorWithoutValue
	Value() (*Plasmid, error)
}

type csvPlasmidReader struct {
	*csource.CsvReader
}

//NewPlasmidReader is to get an instance of PlasmidReader
func NewPlasmidReader(r io.Reader) PlasmidReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new Inventory instance
func (pr *csvPlasmidReader) Value() (*Plasmid, error) {
	p := new(Plasmid)
	if pr.Err != nil {
		return p, pr.Err
	}
	p.Id = pr.Record[0]
	p.Plasmid = pr.Record[1]
	return p, nil
}
