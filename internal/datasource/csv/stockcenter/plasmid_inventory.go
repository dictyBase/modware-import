package stockcenter

import (
	"encoding/csv"
	"io"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

//PlasmidInventory is the container for plasmid inventory
type PlasmidInventory struct {
	Id               string
	PrivateComment   string
	StoredOn         time.Time
	PhysicalLocation string
	StoredAs         string
	ObtainedAs       string
}

//PlasmidInventoryReader is the defined interface for reading the data
type PlasmidInventoryReader interface {
	datasource.IteratorWithoutValue
	Value() (*PlasmidInventory, error)
}

type csvPlasmidInventoryReader struct {
	*csource.CsvReader
}

//NewCsvPlasmidInventoryReader is to get an instance of PlasmidInventoryReader
func NewCsvPlasmidInventoryReader(r io.Reader) PlasmidInventoryReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvPlasmidInventoryReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StrainInventory instance
func (pir *csvPlasmidInventoryReader) Value() (*PlasmidInventory, error) {
	inv := new(PlasmidInventory)
	if pir.Err != nil {
		return inv, pir.Err
	}
	inv.Id = pir.Record[0]
	inv.PhysicalLocation = pir.Record[1]
	if len(pir.Record[2]) > 0 {
		inv.ObtainedAs = pir.Record[2]
	}
	if len(pir.Record[3]) > 0 {
		inv.StoredAs = pir.Record[3]
	}
	if len(pir.Record[4]) > 0 {
		storedOn, err := time.Parse(registry.STOCK_DATE_LAYOUT, pir.Record[4])
		if err != nil {
			return inv, err
		}
		inv.StoredOn = storedOn
	}
	if len(pir.Record[5]) > 0 {
		inv.PrivateComment = pir.Record[5]
	}
	return inv, nil
}
