//package stockcenter is the data source for stockcenter and related data
package stockcenter

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

const orderDateLayout = "2006-01-02 15:04:05"

//StockOrder is the container for order data
type StockOrder struct {
	CreatedAt time.Time
	User      string
	Items     []string
}

//StockOrderReader is the defined interface for reading the data
type StockOrderReader interface {
	datasource.IteratorWithoutValue
	Value() (*StockOrder, error)
}

type csvOrderReader struct {
	*csource.CsvReader
}

//NewCsvStockOrderReader is to get an instance of order reader
func NewCsvStockOrderReader(r io.Reader) StockOrderReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	return &csvOrderReader{&csource.CsvReader{Reader: r}}
}

//Value gets a new StockOrder instance
func (or *csvOrderReader) Value() (*StockOrder, error) {
	so := new(StockOrder)
	if or.Err != nil {
		return so, or.Err
	}
	created, err := time.Parse(orderDateLayout, or.Record[0])
	if err != nil {
		return so, fmt.Errorf("error in parsing order data %s", err)
	}
	so.CreatedAt = created
	so.User = or.Record[1]
	so.Items = or.Record[2]
	return so, nil
}
