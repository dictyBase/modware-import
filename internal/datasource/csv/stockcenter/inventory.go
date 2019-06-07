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
	StrainId         string
	PrivateComment   string
	PublicComment    string
	StoredOn         time.Time
	StoredAs         string
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
	inv.StoredOn = storedOn
	inv.StrainId = sir.Record[0]
	inv.PhysicalLocation = sir.Record[1]
	inv.VialColor = sir.Record[2]
	vc, _ := strconv.ParseInt(sir.Record[3], 10, 64)
	inv.VialsCount = vc
	inv.StoredAs = sir.Record[4]
	inv.PrivateComment = sir.Record[6]
	if len(sir.Record) >= 8 {
		inv.PublicComment = sir.Record[7]
	}
	return inv, nil
}

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
	storedOn, err := time.Parse(invDateLayout, pir.Record[4])
	if err != nil {
		return inv, err
	}
	inv.StoredOn = storedOn
	inv.Id = pir.Record[0]
	inv.PhysicalLocation = pir.Record[1]
	if len(pir.Record[2]) > 0 {
		inv.ObtainedAs = pir.Record[2]
	}
	if len(pir.Record[3]) > 0 {
		inv.StoredAs = pir.Record[3]
	}
	if len(pir.Record[5]) > 0 {
		inv.PrivateComment = pir.Record[5]
	}
	return inv, nil
}
