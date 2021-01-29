package gwdi

import (
	"fmt"
	"regexp"
	"strings"

	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

const (
	genoTmpl = `axeA1,axeB1,axeC1,%s,[bsRcas],bsR`
)

var grxp = regexp.MustCompile(`(DDB_G[0-9]{5,}){2,}`)

var insrMap = map[string]string{
	"G1": "GATC",
	"G2": "GATC",
	"C4": "CATG",
	"C6": "CATG",
	"C8": "CATG",
}

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
	"DDB0232428": "chr 1",
	"DDB0232429": "chr 2",
	"DDB0232430": "chr 3",
	"DDB0232431": "chr 4",
	"DDB0232432": "chr 5",
	"DDB0232433": "chr 6",
}

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

func summIntraMultiple() string {
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) mutants,")
	b.WriteString(" %s intragenic insertions;")
	b.WriteString(" insertion at position %s, %s,")
	b.WriteString(" at %s genomic sites; %s orientation;")
	b.WriteString(" this stock contains %s individual mutants")
	return b.String()
}

func summaryIntraSingle() string {
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) %s mutant;")
	b.WriteString(" insertion at position %s, %s,")
	b.WriteString(" at %s genomic sites; %s orientation.")
	return b.String()
}

func defaultGWDIStrain() *GWDIStrain {
	return &GWDIStrain{
		Parent:      "DBS0351471",
		Plasmid:     "Blasticidin S resistance cassette",
		Depositor:   "baldwinAJ@cardiff.ac.uk",
		Species:     "Dictyostelium discoideum",
		Publication: "10.1101/582072",
		Characters: []string{
			"blasticidin resistant",
			"axenic",
			"null mutant",
		},
		Properties: map[string]*tsource.StockProp{
			regs.DICTY_ANNO_ONTOLOGY: {
				Property: "mutant type",
				Value:    "endogenous insertion",
			},
			regs.DICTY_MUTAGENESIS_ONTOLOGY: {
				Property: "mutagenesis method",
				Value:    "Restriction Enzyme-Mediated Integration",
			},
		},
	}
}

func intragenic_multiple_annotation(r []string) *GWDIStrain {
	d := fmt.Sprintf("%s-", r[8])
	strain := defaultGWDIStrain()
	strain.Label = d
	strain.Name = r[0]
	strain.Genes = []string{r[8]}
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	strain.Summary = fmt.Sprintf(
		summIntraMultiple(),
		d, r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
}

func intragenic_single_annotation(r []string) *GWDIStrain {
	d := fmt.Sprintf("%s-", r[8])
	strain := defaultGWDIStrain()
	strain.Label = d
	strain.Name = r[0]
	strain.Genes = []string{r[8]}
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	strain.Summary = fmt.Sprintf(
		summaryIntraSingle(),
		d, r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
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
