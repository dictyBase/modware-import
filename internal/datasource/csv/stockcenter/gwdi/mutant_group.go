package gwdi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type GWDIMutantReader interface {
	Next() bool
	Value() (*GWDIStrain, error)
}

type groupItr struct {
	itr iterator.Iterator
}

func NewGWDIMutantIterator(itr iterator.Iterator) GWDIMutantReader {
	return &groupItr{itr: itr}
}

func (g *groupItr) Next() bool {
	return g.itr.Next()
}

func (g *groupItr) Value() (*GWDIStrain, error) {
	strain := &GWDIStrain{}
	if err := json.Unmarshal(g.itr.Value(), strain); err != nil {
		return strain, fmt.Errorf("error in decoding value for strain group %s", err)
	}
	return strain, nil
}

type annoFn func(r []string) *GWDIStrain

// GWDI is for managing gwdi data
type GWDI struct {
	listCache  ListCacher
	mapper     IdMapper
	reader     io.Reader
	annoMapper map[string]annoFn
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
	g.annoMapper = map[string]annoFn{
		"NA_multiple":              multiple_na_annotation,
		"NA_single":                single_na_annotation,
		"intergenic_down_multiple": intergenic_multiple_down_annotation,
		"intergenic_up_multiple":   intergenic_multiple_up_annotation,
		"intergenic_both_multiple": intergenic_multiple_both_annotation,
		"intergenic_none_multiple": intergenic_multiple_no_gene_annotation,
		"intergenic_up_single":     intergenic_single_up_annotation,
		"intergenic_down_single":   intergenic_single_down_annotation,
		"intergenic_both_single":   intergenic_single_both_annotation,
		"intergenic_none_single":   intergenic_single_no_gene_annotation,
		"intragenic_single":        intragenic_single_annotation,
		"intragenic_multiple":      intragenic_multiple_annotation,
	}
	return g, nil
}

func (g *GWDI) MutantReader(group string) GWDIMutantReader {
	return NewGWDIMutantIterator(g.listCache.IterateByPrefix([]byte(group)))
}

func (g *GWDI) AllGroups() ([]string, error) {
	var all []string
	prefixes, err := g.listCache.CommonPrefixes()
	if err != nil {
		return all, err
	}
	for _, p := range prefixes {
		all = append(all, string(p))
	}
	return all, nil
}

func (g *GWDI) AnnotateMutant() error {
	if err := g.DedupId(); err != nil {
		return err
	}
	return g.GroupMutant()
}

func (g *GWDI) GroupMutant() error {
	cache := g.listCache
	cache.StartBatch()
	m := g.mapper
	itr := m.Iterate()
	for itr.Next() {
		r := strings.Split(string(itr.Value()), "\t")
		group := inferGroup(r)
		fn, ok := g.annoMapper[group]
		if !ok {
			return fmt.Errorf("unexpected group %s", group)
		}
		value, err := json.Marshal(fn(r))
		if err != nil {
			return fmt.Errorf("error in encoding strain %s", err)
		}
		cache.AppendToBatch([]byte(group), value)
	}
	itr.Release()
	if err := cache.CommitBatch(); err != nil {
		return fmt.Errorf("error in writing to cache %s", err)
	}
	return nil
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
		ok, err := m.Exist([]byte(record[0]))
		if err != nil {
			return err
		}
		if ok {
		INNER:
			for _, p := range chars {
				id := fmt.Sprintf("%s%s", record[0], string(p))
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
		if err := m.Put([]byte(record[0]), []byte(rstr)); err != nil {
			return err
		}
	}
	return nil
}

func inferGroup(r []string) string {
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
