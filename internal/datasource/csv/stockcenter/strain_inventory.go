package stockcenter

import (
	"encoding/csv"
	"io"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

const invDateLayout = "02-JAN-06"

//StrainInventory is the container for strain inventory
type StrainInventory struct {
	StrainId         string
	PrivateComment   string
	PublicComment    string
	StoredOn         time.Time
	StoredAs         string
	VialsCount       string
	VialColor        string
	PhysicalLocation string
	RecordLine       string
}

//StrainInventoryReader is the defined interface for reading the data
type StrainInventoryReader interface {
	datasource.IteratorWithoutValue
	Value() (*StrainInventory, error)
}

type csvStrainInventoryReader struct {
	*csource.CsvReader
}

//NewCsvStrainInventoryReader is to get an instance of StrainInventoryReader
func NewCsvStrainInventoryReader(r io.Reader) StrainInventoryReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainInventoryReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StrainInventory instance
func (sir *csvStrainInventoryReader) Value() (*StrainInventory, error) {
	inv := new(StrainInventory)
	if sir.Err != nil {
		return inv, sir.Err
	}
	storedOn, err := time.Parse(registry.STOCK_DATE_LAYOUT, sir.Record[5])
	if err != nil {
		return inv, err
	}
	inv.StoredOn = storedOn
	inv.StrainId = sir.Record[0]
	inv.PhysicalLocation = sir.Record[1]
	inv.VialColor = sir.Record[2]
	inv.VialsCount = sir.Record[3]
	inv.StoredAs = sir.Record[4]
	inv.PrivateComment = sir.Record[6]
	if len(sir.Record) >= 8 {
		inv.PublicComment = sir.Record[7]
	}
	inv.RecordLine = strings.Join(sir.Record, "\t")
	return inv, nil
}
