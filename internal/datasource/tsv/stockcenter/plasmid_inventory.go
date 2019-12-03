package stockcenter

import (
	"bufio"
	"io"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
)

//PlasmidInventory is the container for plasmid inventory
type PlasmidInventory struct {
	PlasmidID        string
	PrivateComment   string
	StoredOn         time.Time
	PhysicalLocation string
	StoredAs         string
	ObtainedAs       string
	RecordLine       string
}

//PlasmidInventoryReader is the defined interface for reading the data
type PlasmidInventoryReader interface {
	datasource.IteratorWithoutValue
	Value() (*PlasmidInventory, error)
}

type tsvPlasmidInventoryReader struct {
	*tsource.TsvReader
}

//NewTsvPlasmidInventoryReader is to get an instance of PlasmidInventoryReader
func NewTsvPlasmidInventoryReader(r io.Reader) PlasmidInventoryReader {
	tr := bufio.NewScanner(r)
	return &tsvPlasmidInventoryReader{&tsource.TsvReader{Reader: tr}}
}

//Value gets a new StrainInventory instance
func (pir *tsvPlasmidInventoryReader) Value() (*PlasmidInventory, error) {
	inv := new(PlasmidInventory)
	if pir.Err != nil {
		return inv, pir.Err
	}
	inv.PlasmidID = pir.Record[0]
	inv.PhysicalLocation = pir.Record[1]
	if len(pir.Record[2]) > 0 {
		inv.ObtainedAs = pir.Record[2]
	}
	if len(pir.Record[3]) > 0 {
		inv.StoredAs = pir.Record[3]
	}
	if len(pir.Record[4]) > 0 {
		storedOn, err := parseInvDate(pir.Record[4])
		if err != nil {
			return inv, err
		}
		inv.StoredOn = storedOn
	}
	if len(pir.Record[5]) > 0 {
		inv.PrivateComment = pir.Record[5]
	}
	inv.RecordLine = strings.Join(pir.Record, "\t")
	return inv, nil
}
