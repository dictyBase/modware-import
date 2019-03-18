package stockcenter

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

const invDateLayout = "02-JAN-06"

//Inventory is the container for strain inventory
type Inventory struct {
	Id               string
	PrivateComment   string
	PublicComment    string
	StoredOn         time.Time
	StoredType       string
	VialsCount       int64
	VialColor        string
	PhysicalLocation string
}

//InventoryReader is the defined interface for reading the data
type InventoryReader interface {
	datasource.IteratorWithoutValue
	Value() (*Characteristics, error)
}

type csvInventoryReader struct {
	*csource.CsvReader
}

//NewCharacterReader is to get an instance of InventoryReader
func NewCsvInventoryReader(r io.Reader) InventoryReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvCharacterReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new Inventory instance
func (ci *csvInventoryReader) Value() (*Inventory, error) {
	inv := new(Inventory)
	if ci.Err != nil {
		return inv, ci.Err
	}
	storedOn, err := time.Parse(invDateLayout, ci.Record[5])
	if err != nil {
		return inv, err
	}
	inv.Id = ci.Record[0]
	inv.PhysicalLocation = ci.Record[1]
	inv.VialColor = ci.Record[2]
	vc, _ := strconv.ParseInt(ci.Record[3], 10, 64)
	inv.VialsCount = vc
	inv.StoredType = ci.Record[4]
	inv.StoredOn = storedOn
	inv.PublicComment = ci.Record[6]
	if len(ci.Record) >= 8 {
		inv.PrivateComment = ci.Record[7]
	}
	return inv, nil
}
