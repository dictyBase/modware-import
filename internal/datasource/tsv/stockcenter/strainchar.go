package stockcenter

import (
	"bufio"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
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

type tsvCharacterReader struct {
	*tsource.TsvReader
}

//NewTsvCharacterReader is to get an instance of character reader
func NewTsvCharacterReader(r io.Reader) CharacterReader {
	tr := bufio.NewScanner(r)
	return &tsvCharacterReader{&tsource.TsvReader{Reader: tr}}
}

//Value gets a new Characteristics instance
func (cr *tsvCharacterReader) Value() (*Characteristics, error) {
	c := new(Characteristics)
	if cr.Err != nil {
		return c, cr.Err
	}
	c.Id = cr.Record[0]
	c.Character = cr.Record[1]
	return c, nil
}
