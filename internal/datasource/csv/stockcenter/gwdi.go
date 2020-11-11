package stockcenter

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"

	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"

	"github.com/dictyBase/modware-import/internal/datasource"
	csource "github.com/dictyBase/modware-import/internal/datasource/csv"
	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
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

var grxp = regexp.MustCompile(`(DDB_G[0-9]{5,}){2,}`)

//GWDIStrain is the container for GWDI strain
type GWDIStrain struct {
	Label       string
	Name        string
	Summary     string
	Genotype    string
	Parent      string
	Plasmid     string
	Species     string
	Depositor   string
	Publication string
	Characters  []string
	Genes       []string
	Properties  map[string]*tsource.StockProp
}

//GWDIStrainReader is the defined interface for reading the data
type GWDIStrainReader interface {
	datasource.IteratorWithoutValue
	Value() (*GWDIStrain, error)
}

type csvGWDIStraineader struct {
	*csource.CsvReader
}

func defaultGWDIStrain() *GWDIStrain {
	gst := &GWDIStrain{
		Parent:      "DBS0351471",
		Plasmid:     "Blasticidin S resistance cassette",
		Depositor:   "baldwinAJ@cardiff.ac.uk",
		Species:     "Dictyostelium discoideum",
		Publication: "10.1101/582072",
	}
	gst.Characters = []string{
		"blasticidin resistant",
		"axenic",
		"null mutant",
	}
	gst.Properties = map[string]*tsource.StockProp{
		regs.DICTY_ANNO_ONTOLOGY: {
			Property: "mutant type",
			Value:    "endogenous insertion",
		},
		regs.DICTY_MUTAGENESIS_ONTOLOGY: {
			Property: "mutagenesis method",
			Value:    "Restriction Enzyme-Mediated Integration",
		},
	}
	return gst
}

//NewGWDIStrainReader is to get an instance of GWDIStrainReader
func NewGWDIStrainReader(r io.Reader) GWDIStrainReader {
	cr := csv.NewReader(r)
	cr.Comment = '#'
	return &csvGWDIStraineader{&csource.CsvReader{Reader: cr}}
}

//Value gets a new GWDIStrain instance
func (g *csvGWDIStraineader) Value() (*GWDIStrain, error) {
	gst := defaultGWDIStrain()
	if g.Err != nil {
		return gst, g.Err
	}
	gene := g.Record[7]
	switch {
	case strings.HasPrefix(gene, "DDB_G"):
		gst.Label = g.Record[0]
	case g.Record[7] == "none":
		gst.Label = g.Record[0]
	default:
		gst.Label = fmt.Sprintf("%s-", gene)
	}
	switch g.Record[6] {
	case "intragenic", "intergenic_down", "intergenic_up":
		gst.Genes = []string{gene}
	case "intergenic_both":
		gst.Genes = grxp.FindAllString(gene, -1)
	}
	gst.Name = g.Record[0]
	var summ strings.Builder
	fmt.Fprintf(
		&summ,
		"Genome Wide Dictyostelium Insertion bank(GWDI) %s mutant;",
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
