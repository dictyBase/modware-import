package gwdi

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// GWDI is for managing gwdi data
type GWDI struct {
	listCache ListCacher
	mapper    IdMapper
	reader    io.Reader
}

func NewGWDI(r io.Reader) (*GWDI, error) {
	g := &GWDI{}
	c, err := NewListCache()
	if err != nil {
		return g, err
	}
	m, err := NewIdMap()
	if err != nil {
		return g, err
	}
	g.listCache = c
	g.mapper = m
	g.reader = r
	return g, nil
}

func (g *GWDI) GroupMutant() error {
	cache := g.listCache
	cache.StartBatch()
	m := g.mapper
	itr := m.Iterate()
	for itr.Next() {
		r := strings.Split(string(itr.Value()), "\t")
		group := createGroups(r)
		cache.AppendToBatch([]byte(key), value)
	}
	itr.Release()
	if err := cache.CommitBatch(); err != nil {
		return fmt.Errorf("error in writing to cache %s", err)
	}
}

func (g *GWDI) DedupId() error {
	chars := "abcdefghijklmnopqrst"
	m := g.mapper
	r := csv.NewReader(g.reader)
	r.Comment = '#'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error in reading record %s", err)
		}
		id := record[0]
		ok, err := m.Exist([]byte(id))
		if err != nil {
			return err
		}
		if ok {
		INNER:
			for _, p := range chars {
				id = fmt.Sprintf("%s%s", id, string(p))
				ok, err := m.Exist([]byte(id))
				if err != nil {
					return err
				}
				if !ok {
					record[0] = id
					break INNER
				}
			}
		}
		rstr := strings.Join(record, "\t")
		if err := m.Put([]byte(id), []byte(rstr)); err != nil {
			return err
		}
	}
	return nil
}

func createGroups(r []string) string {
	var group string
	if r[6] == "intragenic" || r[6] == "NA" {
		if r[4] == "1" {
			group = fmt.Sprintf("%s_single", r[6])
		} else {
			group = fmt.Sprintf("%s_multiple", r[6])
		}
	} else {
		if r[7] == "none" {
			if r[4] == "1" {
				group = fmt.Sprintf("%s_none_single", r[6])
			} else {
				group = fmt.Sprintf("%s_none_multiple", r[6])
			}
		} else if r[4] == "1" {
			group = fmt.Sprintf("%s_single", r[6])
		} else {
			group = fmt.Sprintf("%s_multiple", r[6])
		}
	}
	return group
}
