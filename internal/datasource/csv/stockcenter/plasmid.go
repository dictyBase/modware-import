package stockcenter

import (
	"encoding/csv"
	"io"
	"time"

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

//Plasmid is the container for plasmid data
type Plasmid struct {
	Id           string
	Summary      string
	User         string
	CreatedOn    time.Time
	UpdatedOn    time.Time
	Name         string
	Publications []string
	Genes        []string
}

//PlasmidReader is the defined interface for reading the plasmid data
type PlasmidReader interface {
	datasource.IteratorWithoutValue
	Value() (*Plasmid, error)
}

type csvPlasmidReader struct {
	*csource.CsvReader
	alookup StockAnnotatorLookup
	plookup StockPubLookup
	glookup StockGeneLookup
}

//NewCsvPlasmidReader is to get an instance of PlasmidReader instance
func NewCsvPlasmidReader(r io.Reader, al StockAnnotatorLookup, pl StockPubLookup, gl StockGeneLookup) PlasmidReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidReader{
		CsvReader: &csource.CsvReader{Reader: cr},
		alookup:   al,
		plookup:   pl,
		glookup:   gl,
	}
}

//Value gets a new Plasmid instance
func (plr *csvPlasmidReader) Value() (*Plasmid, error) {
	p := new(Plasmid)
	if plr.Err != nil {
		return p, plr.Err
	}
	p.Id = plr.Record[0]
	p.Name = plr.Record[1]
	p.Summary = plr.Record[2]
	user, c, u, ok := plr.alookup.StockAnnotator(plr.Record[0])
	if ok {
		p.User = user
		p.CreatedOn = c
		p.UpdatedOn = u
	}
	p.Publications = append(p.Publications, plr.plookup.StockPub(p.Id)...)
	p.Genes = append(p.Genes, plr.glookup.StockGene(p.Id)...)
	return p, nil
}
