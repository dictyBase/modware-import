package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//StockProp is the container for stock properties
type StockProp struct {
	Id       string
	Property string
	Value    string
}

//StockPropReader is the defined interface for reading the data
type StockPropReader interface {
	datasource.IteratorWithoutValue
	Value() (*StockProp, error)
}

type csvStockPropReader struct {
	*csource.CsvReader
}

//NewCsvStockPropReader is to get an instance of StockPropReader
func NewCsvStockPropReader(r io.Reader) StockPropReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStockPropReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StockProp instance
func (spr *csvStockPropReader) Value() (*StockProp, error) {
	prop := new(StockProp)
	if spr.Err != nil {
		return prop, spr.Err
	}
	prop.Id = spr.Record[0]
	prop.Property = spr.Record[1]
	prop.Value = spr.Record[2]
	return prop, nil
}
