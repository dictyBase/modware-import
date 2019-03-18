package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//StrainProp is the container for strain properties
type StrainProp struct {
	Id       string
	Property string
	Value    string
}

//StrainPropReader is the defined interface for reading the data
type StrainPropReader interface {
	datasource.IteratorWithoutValue
	Value() (*StrainProp, error)
}

type csvStrainPropReader struct {
	*csource.CsvReader
}

//NewStrainPropReader is to get an instance of StrainPropReader
func NewStrainPropReader(r io.Reader) StrainPropReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainPropReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StrainProp instance
func (spr *csvStrainPropReader) Value() (*StrainProp, error) {
	prop := new(StrainProp)
	if spr.Err != nil {
		return prop, spr.Err
	}
	prop.Id = spr.Record[0]
	prop.Property = spr.Record[1]
	prop.Value = spr.Record[2]
	return prop, nil
}
