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

//NewPlasmidGenbankReader is to get an instance of PlasmidGenbankReader
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

//PlasmidGene is the container for plasmid and gene identifier links
type PlasmidGene struct {
	Id     string
	GeneId string
}

//PlasmidGene is the defined interface for reading the data
type PlasmidGenbankReader interface {
	datasource.IteratorWithoutValue
	Value() (*PlasmidGene, error)
}

type csvPlasmidGeneReader struct {
	*csource.CsvReader
}

//NewPlasmidGeneReader is to get an instance of PlasmidGeneReader
func NewPlasmidGeneReader(r io.Reader) PlasmidGeneReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidGeneReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new PlasmidGene instance
func (pger *csvPlasmidGeneReader) Value() (*PlasmidGene, error) {
	gene := new(PlasmidGene)
	if pger.Err != nil {
		return gene, pger.Err
	}
	gene.Id = pger.Record[0]
	gene.GeneId = pger.Record[1]
	return gene, nil
}
