package stockcenter

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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
	sgr := bufio.NewScanner(r)
	for sgr.Scan() {
		record := strings.Split(sgr.Text(), "\t")
		if len(record) != 2 {
			return l, fmt.Errorf("does not expected record in line %s", sgr.Text())
		}
		if v, ok := m.Get(record[0]); ok {
			s := v.([]string)
			m.Put(record[0], append(s, record[1]))
			continue
		}
		m.Put(record[0], []string{record[1]})
	}
	if err := sgr.Err(); err != nil {
		return l, fmt.Errorf("error in scanning %s", err)
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
