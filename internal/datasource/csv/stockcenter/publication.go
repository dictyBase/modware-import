package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//StockPub is the container for stock and its associated publication
type StockPub struct {
	Id    string
	PubId string
}

//StockPubReader is the defined interface for reading the data
type StockPubReader interface {
	datasource.IteratorWithoutValue
	Value() (*StockPub, error)
}

type csvStockPubReader struct {
	*csource.CsvReader
}

//NewStockPubReader is to get an instance of StockPubReader
func NewStockPubReader(r io.Reader) StockPubReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStockPubReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StockProp instance
func (spr *csvStockPubReader) Value() (*StockPub, error) {
	pub := new(StockPub)
	if spr.Err != nil {
		return pub, spr.Err
	}
	pub.Id = spr.Record[0]
	pub.PubId = spr.Record[1]
	return pub, nil
}
