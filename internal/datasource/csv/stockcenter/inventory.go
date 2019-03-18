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

//StrainInventory is the container for strain inventory
type StrainInventory struct {
	Id               string
	PrivateComment   string
	PublicComment    string
	StoredOn         time.Time
	StoredType       string
	VialsCount       int64
	VialColor        string
	PhysicalLocation string
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
	storedOn, err := time.Parse(invDateLayout, sir.Record[5])
	if err != nil {
		return inv, err
	}
	inv.Id = sir.Record[0]
	inv.PhysicalLocation = sir.Record[1]
	inv.VialColor = sir.Record[2]
	vc, _ := strconv.ParseInt(sir.Record[3], 10, 64)
	inv.VialsCount = vc
	inv.StoredType = sir.Record[4]
	inv.StoredOn = storedOn
	inv.PublicComment = sir.Record[6]
	if len(sir.Record) >= 8 {
		inv.PrivateComment = sir.Record[7]
	}
	return inv, nil
}
