package stockcenter

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
	registry "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

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

type tsvStrainInventoryReader struct {
	*tsource.TsvReader
}

//NewTsvStrainInventoryReader is to get an instance of StrainInventoryReader
func NewTsvStrainInventoryReader(r io.Reader) StrainInventoryReader {
	tr := bufio.NewScanner(r)
	return &tsvStrainInventoryReader{&tsource.TsvReader{Reader: tr}}
}

//Value gets a new StrainInventory instance
func (sir *tsvStrainInventoryReader) Value() (*StrainInventory, error) {
	inv := new(StrainInventory)
	if sir.Err != nil {
		return inv, sir.Err
	}
	inv.StrainId = sir.Record[0]
	inv.PhysicalLocation = sir.Record[1]
	inv.VialColor = sir.Record[2]
	inv.VialsCount = sir.Record[3]
	if len(sir.Record[5]) > 0 {
		inv.StoredAs = sir.Record[5]
	}
	if len(sir.Record[6]) > 0 {
		m := dateRegxp.FindStringSubmatch(sir.Record[6])
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
	inv.PrivateComment = sir.Record[7]
	if len(sir.Record) >= 9 {
		inv.PublicComment = sir.Record[8]
	}
	inv.RecordLine = strings.Join(sir.Record, "\t")
	return inv, nil
}

//NewCsvStrainInventoryReader is to get an instance of StrainInventoryReader
func NewCsvStrainInventoryReader(r io.Reader) StrainInventoryReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	cr.LazyQuotes = true
	return &csvStrainInventoryReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StrainInventory instance
func (sir *csvStrainInventoryReader) Value() (*StrainInventory, error) {
	inv := new(StrainInventory)
	if sir.Err != nil {
		return inv, sir.Err
	}
	inv.StrainId = sir.Record[0]
	inv.PhysicalLocation = sir.Record[1]
	inv.VialColor = sir.Record[2]
	inv.VialsCount = sir.Record[3]
	if len(sir.Record[5]) > 0 {
		inv.StoredAs = sir.Record[5]
	}
	if len(sir.Record[6]) > 0 {
		m := dateRegxp.FindStringSubmatch(sir.Record[6])
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
	inv.PrivateComment = sir.Record[7]
	if len(sir.Record) >= 9 {
		inv.PublicComment = sir.Record[8]
	}
	inv.RecordLine = strings.Join(sir.Record, "\t")
	return inv, nil
}

func ucFirstAllLower(s string) string {
	return fmt.Sprintf("%s%s", string(s[0]), strings.ToLower(s[1:]))
}
