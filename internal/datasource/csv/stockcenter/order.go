//package stockcenter is the data source for stockcenter and related data
package stockcenter

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/dictyBase/modware-import/internal/datasource"
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
	r      *csv.Reader
	record []string
	err    error
}

//NewCsvStockOrderReader is to get an instance of order reader
func NewCsvStockOrderReader(r io.Reader) StockOrderReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	return &CsvOrderReader{r: cr}
}

//Next is to read next row of order data
func (csv *csvOrderReader) Next() bool {
	csv.err = nil
	record, err := csv.r.Read()
	if err == io.EOF {
		return false
	}
	if err != nil {
		csv.err = err
		return true
	}
	if len(csv.record) > 1 {
		csv.record = nil
	}
	csv.record = append(csv.record, record...)
	return true
}

//Value gets a new StockOrder instance
func (csv *csvOrderReader) Value() (*StockOrder, error) {
	so := new(StockOrder)
	if csv.err != nil {
		return so, nil
	}
	created, err := time.Parse(orderDateLayout, csv.record[0])
	if err != nil {
		return so, fmt.Errorf("error in parsing order data %s", err)
	}
	so.CreatedAt = created
	so.User = csv.record[1]
	so.Items = csv.record[2:]
	return so, nil
}
