package gwdi

import (
	"fmt"
	"regexp"

	tsource "github.com/dictyBase/modware-import/internal/datasource/tsv/stockcenter"
	regs "github.com/dictyBase/modware-import/internal/registry/stockcenter"
)

const (
	genoTmpl = `axeA1,axeB1,axeC1,%s,[bsRcas],bsR`
)

var disrupt_rgxp = regexp.MustCompile(`^(DDB_G[0-9]{5,})`)
var suffix_rgxp = regexp.MustCompile(`(^GWDI_\d+_[A-Z]{1,2}_\d+)[a-z]{1}$`)

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
	"DDB0169550": "mitochondrial DNA",
	"DDB0215018": "unintegrated genomic DNA",
	"DDB0220052": "unintegrated genomic DNA",
	"DDB0215151": "unintegrated genomic DNA",
	"DDB0237465": "extrachromosomal ribosomal RNA",
	"DDB0232428": "chr 1",
	"DDB0232429": "chr 2",
	"DDB0232430": "chr 3",
	"DDB0232431": "chr 4",
	"DDB0232432": "chr 5",
	"DDB0232433": "chr 6",
}

//GWDIStrain is the container for GWDI strain
type GWDIStrain struct {
	Label       string                        `json:"label"`
	Name        string                        `json:"name"`
	Summary     string                        `json:"summary"`
	Genotype    string                        `json:"genotype"`
	Parent      string                        `json:"parent"`
	Plasmid     string                        `json:"plasmid,omitempty"`
	Species     string                        `json:"species"`
	Depositor   string                        `json:"depositor"`
	Publication string                        `json:"publication"`
	Characters  []string                      `json:"characters"`
	Genes       []string                      `json:"genes,omitempty"`
	Properties  map[string]*tsource.StockProp `json:"properties"`
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

func removeSuffix(id string) string {
	if suffix_rgxp.MatchString(id) {
		return suffix_rgxp.ReplaceAllString(id, "$1")
	}
	return id
}

func intragenic_mutant_annotation(r []string) *GWDIStrain {
	d := fmt.Sprintf("%s-", r[8])
	strain := defaultGWDIStrain()
	strain.Label = d
	strain.Name = r[0]
	strain.Genes = []string{r[8]}
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	return strain
}

func geneless_mutant_annotation(r []string) *GWDIStrain {
	strain := defaultGWDIStrain()
	strain.Label = removeSuffix(r[0])
	strain.Name = r[0]
	strain.Genotype = fmt.Sprintf(genoTmpl, strain.Label)
	strain.Characters[2] = "mutant"
	strain.Properties[regs.DICTY_ANNO_ONTOLOGY] = &tsource.StockProp{
		Property: "mutant type",
		Value:    "exogenous insertion",
	}
	return strain
}

func intergenic_multiple_no_gene_annotation(r []string) *GWDIStrain {
	strain := geneless_mutant_annotation(r)
	strain.Summary = fmt.Sprintf(
		summInterNoGeneMultiple(),
		r[2], chrMap[r[1]],
		insrMap[r[3]], r[5], r[4],
	)
	return strain
}

func intergenic_multiple_both_annotation(r []string) *GWDIStrain {
	strain := defaultGWDIStrain()
	m := disrupt_rgxp.FindStringSubmatch(r[7])
	d := fmt.Sprintf("[%s/%s]-", m[0], m[1])
	strain.Label = d
	strain.Name = r[0]
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	strain.Characters[2] = "mutant"
	strain.Genes = []string{m[0], m[1]}
	strain.Summary = fmt.Sprintf(
		summInterMultipleBoth(),
		m[0], m[1], r[2], chrMap[r[1]],
		insrMap[r[3]], r[5], r[4],
	)
	return strain
}

func intergenic_multiple_up_down_annotation(r []string, orientation string) *GWDIStrain {
	strain := defaultGWDIStrain()
	d := fmt.Sprintf("[%s]-", r[8])
	strain.Label = d
	strain.Name = r[0]
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	strain.Characters[2] = "mutant"
	strain.Genes = []string{r[8]}
	strain.Summary = fmt.Sprintf(
		summInterMultipleUpDown(orientation),
		r[8], r[2], chrMap[r[1]],
		insrMap[r[3]], r[5], r[4],
	)
	return strain
}

func intergenic_multiple_down_annotation(r []string) *GWDIStrain {
	return intergenic_multiple_up_down_annotation(r, "downstream")
}

func intergenic_multiple_up_annotation(r []string) *GWDIStrain {
	return intergenic_multiple_up_down_annotation(r, "upstream")
}

func intergenic_single_up_down_annotation(r []string, orientation string) *GWDIStrain {
	strain := defaultGWDIStrain()
	d := fmt.Sprintf("[%s]-", r[8])
	strain.Label = d
	strain.Name = r[0]
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	strain.Characters[2] = "mutant"
	strain.Genes = []string{r[8]}
	strain.Summary = fmt.Sprintf(
		summInterUpDown(orientation),
		r[8], r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
}

func intergenic_single_both_annotation(r []string) *GWDIStrain {
	strain := defaultGWDIStrain()
	m := disrupt_rgxp.FindStringSubmatch(r[7])
	d := fmt.Sprintf("[%s/%s]-", m[0], m[1])
	strain.Label = d
	strain.Name = r[0]
	strain.Genotype = fmt.Sprintf(genoTmpl, d)
	strain.Characters[2] = "mutant"
	strain.Genes = []string{m[0], m[1]}
	strain.Summary = fmt.Sprintf(
		summInterSingleBoth(),
		m[0], m[1], r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
}

func intergenic_single_down_annotation(r []string) *GWDIStrain {
	return intergenic_single_up_down_annotation(r, "downstream")
}

func intergenic_single_up_annotation(r []string) *GWDIStrain {
	return intergenic_single_up_down_annotation(r, "upstream")
}

func intergenic_single_no_gene_annotation(r []string) *GWDIStrain {
	strain := geneless_mutant_annotation(r)
	strain.Summary = fmt.Sprintf(
		summInterNoGeneSingle(),
		r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
}

func intragenic_multiple_annotation(r []string) *GWDIStrain {
	strain := intragenic_mutant_annotation(r)
	strain.Summary = fmt.Sprintf(
		summIntraMultiple(),
		strain.Label, r[2], chrMap[r[1]],
		insrMap[r[3]], r[5], r[4],
	)
	return strain
}

func intragenic_single_annotation(r []string) *GWDIStrain {
	strain := intragenic_mutant_annotation(r)
	strain.Summary = fmt.Sprintf(
		summaryIntraSingle(),
		strain.Label, r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
}

func single_na_annotation(r []string) *GWDIStrain {
	strain := geneless_mutant_annotation(r)
	strain.Properties[regs.DICTY_ANNO_ONTOLOGY] = &tsource.StockProp{
		Property: "mutant type",
		Value:    "endogenous insertion",
	}
	strain.Summary = fmt.Sprintf(
		summaryIntraSingle(),
		strain.Label, r[2], chrMap[r[1]],
		insrMap[r[3]], r[5],
	)
	return strain
}

func multiple_na_annotation(r []string) *GWDIStrain {
	strain := geneless_mutant_annotation(r)
	strain.Properties[regs.DICTY_ANNO_ONTOLOGY] = &tsource.StockProp{
		Property: "mutant type",
		Value:    "endogenous insertion",
	}
	strain.Summary = fmt.Sprintf(
		summaryNAMultiple(),
		strain.Label, r[2], chrMap[r[1]],
		insrMap[r[3]], r[5], r[4],
	)
	return strain
}
