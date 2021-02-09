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
		"NA_multiple":              multipleNaAnnotation,
		"NA_single":                singleNaAnnotation,
		"intergenic_down_multiple": intergenicMultipleDownAnnotation,
		"intergenic_up_multiple":   intergenicMultipleUpAnnotation,
		"intergenic_both_multiple": intergenicMultipleBothAnnotation,
		"intergenic_none_multiple": intergenicMultipleNoGeneAnnotation,
		"intergenic_up_single":     intergenicSingleUpAnnotation,
		"intergenic_down_single":   intergenicSingleDownAnnotation,
		"intergenic_both_single":   intergenicSingleBothAnnotation,
		"intergenic_none_single":   intergenicSingleNoGeneAnnotation,
		"intragenic_single":        intragenicSingleAnnotation,
		"intragenic_multiple":      intragenicMultipleAnnotation,
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
		id, err := g.syntheticId(record[0])
		if err != nil {
			return err
		}
		record[0] = id
		rstr := strings.Join(record, "\t")
		if err := m.Put([]byte(record[0]), []byte(rstr)); err != nil {
			return err
		}
	}
	return nil
}

func (g *GWDI) syntheticId(r string) (string, error) {
	ok, err := g.mapper.Exist([]byte(r))
	if err != nil {
		return r, err
	}
	if !ok {
		return r, nil
	}
	id, err := g.selectId(r)
	if err != nil {
		return r, err
	}
	return id, nil
}

func (g *GWDI) selectId(r string) (string, error) {
	chars := "abcdefghijklmnopqrst"
	var id string
	for _, p := range chars {
		id = fmt.Sprintf("%s%s", r, string(p))
		ok, err := g.mapper.Exist([]byte(id))
		if err != nil {
			return id, err
		}
		if !ok {
			break
		}
	}
	return id, nil
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
		switch {
		case r[7] == "none":
			if r[4] == "1" {
				group = fmt.Sprintf("%s_none_single", r[6])
			} else {
				group = fmt.Sprintf("%s_none_multiple", r[6])
			}
		case r[4] == "1":
			group = fmt.Sprintf("%s_single", r[6])
		default:
			group = fmt.Sprintf("%s_multiple", r[6])
		}
	}
	return group
}
