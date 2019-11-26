package stockcenter

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
	"github.com/dictyBase/modware-import/internal/regexp"
	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

//PlasmidInventory is the container for plasmid inventory
type PlasmidInventory struct {
	PlasmidId        string
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
	inv.PlasmidId = pir.Record[0]
	inv.PhysicalLocation = pir.Record[1]
	if len(pir.Record[2]) > 0 {
		inv.ObtainedAs = pir.Record[2]
	}
	if len(pir.Record[3]) > 0 {
		inv.StoredAs = pir.Record[3]
	}
	if len(pir.Record[4]) > 0 {
		m := regexp.DateRegxp.FindStringSubmatch(pir.Record[4])
		if m != nil {
			storedOn, err := time.Parse(
				registry.STOCK_DATE_LAYOUT,
				fmt.Sprintf("%s-%s-%s", m[1], ucFirstAllLower(m[2]), m[3]),
			)
			if err != nil {
				return inv, err
			}
			inv.StoredOn = storedOn
		}
	}
	if len(pir.Record[5]) > 0 {
		inv.PrivateComment = pir.Record[5]
	}
	inv.RecordLine = strings.Join(pir.Record, "\t")
	return inv, nil
}
