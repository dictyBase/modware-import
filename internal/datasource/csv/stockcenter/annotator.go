package stockcenter

import (
	"encoding/csv"
	"errors"
	"io"
	"time"

	"github.com/emirpasic/gods/maps/hashmap"
)

const annoDateLayout = "2006-01-02 15:04:05"

var annMap = map[string]string{
	"jf":             "jf31@columbia.edu",
	"CGM_DDB_JAKOB":  "jf31@columbia.edu",
	"CGM_DDB_PASC":   "pgaudet@northwestern.edu",
	"CGM_DDB_STEPHY": "jf31@columbia.edu",
	"ah":             "jf31@columbia.edu",
	"sm":             "jf31@columbia.edu",
	"CGM_DDB_MARC":   "m-vincelli@northwestern.edu",
	"CGM_DDB_PFEY":   "pfey@northwestern.edu",
	"CGM_DDB_BOBD":   "robert-dodson@northwestern.edu",
	"CGM_DDB_KPIL":   "kpilchar@northwestern.edu",
	"CGM_DDB":        "dictybase@northwestern.edu",
	"CGM_DDB_KERRY":  "ksheppard@northwestern.edu",
}

// StockAnnotatorLookup is an interface for retrieving stock annotator
type StockAnnotatorLookup interface {
	StockAnnotator(id string) (string, time.Time, time.Time, bool)
}

type saLookup struct {
	smap *hashmap.Map
}

func NewStockAnnotatorLookup(r io.Reader) (StockAnnotatorLookup, error) {
	l := new(saLookup)
	m := hashmap.New()
	uar := csv.NewReader(r)
	uar.FieldsPerRecord = -1
	for {
		record, err := uar.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return l, errors.New("error in reading annotator file")
		}
		var usr string
		if _, ok := annMap[record[1]]; !ok {
			usr = annMap["jf"]
		} else {
			usr = annMap[record[1]]
		}
		m.Put(
			record[0],
			[]string{usr, record[2], record[3]},
		)
	}
	l.smap = m
	return l, nil
}

func (l *saLookup) StockAnnotator(id string) (string, time.Time, time.Time, bool) {
	var c, u time.Time
	v, ok := l.smap.Get(id)
	if !ok {
		return "", c, u, ok
	}
	record, _ := v.([]string)
	createdOn, err := time.Parse(annoDateLayout, record[1])
	if err != nil {
		return "", c, u, false
	}
	updatedOn, err := time.Parse(annoDateLayout, record[2])
	if err != nil {
		return "", c, u, false
	}
	return record[0], createdOn, updatedOn, true
}
