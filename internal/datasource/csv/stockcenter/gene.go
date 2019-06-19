package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/emirpasic/gods/maps/hashmap"
)

//StockGeneLookup is an interface for retrieving gene mapped to
//a stock entry
type StockGeneLookup interface {
	//StockGene looks up a stock identifier and returns a slice
	//with a list of gene identifiers
	StockGene(id string) []string
}

type saGeneLookup struct {
	smap *hashmap.Map
}

//NewStockGeneLookp returns a struct implementing StockGeneLookup interface
func NewStockGeneLookp(r io.Reader) (StockGeneLookup, error) {
	l := new(saGeneLookup)
	m := hashmap.New()
	sgr := csv.NewReader(r)
	sgr.FieldsPerRecord = -1
	for {
		record, err := sgr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return l, err
		}
		if v, ok := m.Get(record[0]); ok {
			s := v.([]string)
			m.Put(record[0], append(s, record[1]))
			continue
		}
		m.Put(record[0], []string{record[1]})
	}
	l.smap = m
	return l, nil
}

func (sl *saGeneLookup) StockGene(id string) []string {
	if v, ok := sl.smap.Get(id); ok {
		s := v.([]string)
		return s
	}
	return []string{""}
}
