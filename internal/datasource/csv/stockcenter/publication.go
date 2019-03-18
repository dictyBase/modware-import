package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

//StrainPublication is the container for strain
type StrainPub struct {
	Id    string
	PubId string
}

//StrainPubReader is the defined interface for reading the data
type StrainPubReader interface {
	datasource.IteratorWithoutValue
	Value() (*StrainPub, error)
}

type csvStrainPubReader struct {
	*csource.CsvReader
}

//NewStrainPubReader is to get an instance of StrainPubReader
func NewStrainPubReader(r io.Reader) StrainPubReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvStrainPubReader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StrainProp instance
func (spr *csvStrainPubReader) Value() (*StrainPub, error) {
	pub := new(StrainPub)
	if spr.Err != nil {
		return pub, spr.Err
	}
	pub.Id = spr.Record[0]
	pub.PubId = spr.Record[1]
	return pub, nil
}
