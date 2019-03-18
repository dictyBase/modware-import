//package stockcenter is the data source for stockcenter and related data
package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//Characteristics is the container for strain characteristics
type Characteristics struct {
	Id        string
	Character string
}

//CharacterReader is the defined interface for reading the data
type CharacterReader interface {
	datasource.IteratorWithoutValue
	Value() (*Characteristics, error)
}

type csvCharacterReader struct {
	*csource.CsvReader
}

//NewCharacterReader is to get an instance of character reader
func NewCsvCharacterReader(r io.Reader) CharacterReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvCharacterReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StockOrder instance
func (cr *csvCharacterReader) Value() (*Characteristics, error) {
	c := new(Characteristics)
	if cr.Err != nil {
		return c, cr.Err
	}
	c.Id = cr.Record[0]
	c.Character = cr.Record[1]
	return c, nil
}
