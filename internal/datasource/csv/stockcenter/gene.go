package stockcenter

import (
	"encoding/csv"
	"io"

	"github.com/emirpasic/gods/maps/hashmap"
)

//StrainGeneLookup is an interface for retrieving gene mapped to
//a strain
type StrainGeneLookup interface {
	//StrainGene looks up a strain identifier and returns a slice
	//with a list of gene identifiers
	StrainGene(id string) []string
}

type saGeneLookup struct {
	smap *hashmap.Map
}

//NewStrainGeneLookp returns an StrainGeneLookup implementing struct
func NewStrainGeneLookp(r io.Reader) (StrainGeneLookup, error) {
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

func (sl *saGeneLookup) StrainGene(id string) []string {
	if v, ok := sl.smap.Get(id); ok {
		s := v.([]string)
		return s
	}
	return []string{""}
}
