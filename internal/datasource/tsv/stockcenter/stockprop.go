package stockcenter

import (
	"bufio"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv"
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

type tsvStockPropReader struct {
	*tsource.TsvReader
}

//NewTsvStockPropReader is to get an instance of StockPropReader
func NewTsvStockPropReader(r io.Reader) StockPropReader {
	tr := bufio.NewScanner(r)
	return &tsvStockPropReader{&tsource.TsvReader{Reader: tr}}
}

//Value gets a new StockProp instance
func (spr *tsvStockPropReader) Value() (*StockProp, error) {
	prop := new(StockProp)
	if spr.Err != nil {
		return prop, spr.Err
	}
	prop.Id = spr.Record[0]
	prop.Property = spr.Record[1]
	prop.Value = spr.Record[2]
	return prop, nil
}
