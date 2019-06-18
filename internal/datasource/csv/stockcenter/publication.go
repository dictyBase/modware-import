package stockcenter

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/emirpasic/gods/maps/hashmap"
)

//StockPubLookup is an interface for retrieving publication
//linked to a stock record
type StockPubLookup interface {
	//StockPub looks up a stock identifier and returns a slice
	//with a list of publication identifiers
	StockPub(id string) []string
}

type saPubLookup struct {
	smap *hashmap.Map
}

//NewStockPubLookup returns an StockPubLookup implementing struct
func NewStockPubLookup(r io.Reader) (StockPubLookup, error) {
	l := new(saPubLookup)
	m := hashmap.New()
	spr := csv.NewReader(r)
	spr.FieldsPerRecord = -1
	for {
		record, err := spr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return l, err
		}
		if strings.HasPrefix(record[1], "d") {
			continue
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

func (sl *saPubLookup) StockPub(id string) []string {
	if v, ok := sl.smap.Get(id); ok {
		s := v.([]string)
		return s
	}
	return []string{""}
}
