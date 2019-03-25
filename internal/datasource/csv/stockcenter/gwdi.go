package stockcenter

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
)

var chrMap = map[string]string{
	"DDB0166986": "DDB0166986",
	"DDB0203883": "DDB0203883",
	"DDB0185230": "DDB0185230",
	"DDB0187328": "DDB0187328",
	"DDB0189305": "DDB0189305",
	"DDB0237465": "R",
	"DDB0169550": "M",
	"DDB0215018": "2F",
	"DDB0220052": "BF",
	"DDB0215151": "3F",
	"DDB0232428": "chr1",
	"DDB0232429": "chr2",
	"DDB0232430": "chr3",
	"DDB0232431": "chr4",
	"DDB0232432": "chr5",
	"DDB0232433": "chr6",
}

//GWDIStrain is the container for GWDI strain
type GWDIStrain struct {
	Label    string
	Name     string
	Summary  string
	GeneId   string
	Genotype string
}

//GWDIStrainReader is the defined interface for reading the data
type GWDIStrainReader interface {
	datasource.IteratorWithoutValue
	Value() (*GWDIStrain, error)
}

type csvGWDIStraineader struct {
	*csource.CsvReader
}

//NewStockPropReader is to get an instance of StrainPropReader
func NewGWDIStrainReader(r io.Reader) GWDIStrainReader {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.Comma = '\t'
	return &csvGWDIStraineader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new StockProp instance
func (g *csvGWDIStraineader) Value() (*GWDIStrain, error) {
	gst := new(GWDIStrain)
	if g.Err != nil {
		return gst, g.Err
	}
	gene := g.Record[7]
	if strings.HasPrefix(gene, "DDB_G") {
		gst.Label = g.Record[0]
		gst.GeneId = gene
	} else if g.Record[7] == "none" {
		gst.Label = g.Record[0]
	} else {
		gst.Label = fmt.Sprintf("%s-", gene)
		gst.GeneId = gene
	}
	gst.Name = g.Record[0]
	var summ strings.Builder
	fmt.Fprintf(
		&summ,
		"Genome Wide Dictyostelium Insertion bank(GWDI) %s mutatnt;",
		gst.Label,
	)
	fmt.Fprintf(
		&summ,
		"insertion at position %s, %s;",
		strings.Replace(g.Record[2], ",", "", -1),
		chrMap[g.Record[1]],
	)
	if _, err := summ.WriteString("used enzyme: SphI"); err != nil {
		return gst, err
	}
	gst.Summary = summ.String()
	return gst, nil
}
