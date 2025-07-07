package mock

import (
	"fmt"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GenerateFeatureAnnotations creates a set of realistic mock feature annotations
func GenerateFeatureAnnotations() []*feature.FeatureAnnotation {
	now := timestamppb.New(time.Now())

	return []*feature.FeatureAnnotation{
		createActAAnnotation(now),
		createMyoBAnnotation(now),
		createPakAAnnotation(now),
		createRasGAnnotation(now),
		createDiscoidin1Annotation(now),
	}
}

func createActAAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: "gene",
		Id:   "DDB_G0267398",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "actA",
			Synonyms: []string{"actin", "act1"},
			Publications: []string{
				"10.1016/j.cell.2023.001234",
				"10.1038/nature.2023.5678",
			},
			Pubmed: []string{"12345678", "87654321"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "function",
					Value:     "cytoskeleton organization",
					CreatedBy: "test@dictybase.org",
					CreatedAt: now,
				},
				{
					Tag:       "location",
					Value:     "cytoplasm",
					CreatedBy: "test@dictybase.org",
					CreatedAt: now,
				},
			},
			Dblinks: []*feature.DbLink{
				{
					PrimaryId: "UNIPROT:P13363",
					Database:  "UniProt",
					Linktype:  "protein",
					Url:       "https://www.uniprot.org/uniprot/P13363",
				},
			},
		},
		CreatedBy:  "test@dictybase.org",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createMyoBAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: "gene",
		Id:   "DDB_G0275199",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "myoB",
			Synonyms: []string{"myosin II heavy chain B", "myo2"},
			Publications: []string{
				"10.1074/jbc.2023.298765",
			},
			Pubmed: []string{"11223344"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "function",
					Value:     "motor activity",
					CreatedBy: "curator@dictybase.org",
					CreatedAt: now,
				},
				{
					Tag:       "pathway",
					Value:     "cell motility",
					CreatedBy: "curator@dictybase.org",
					CreatedAt: now,
				},
			},
			Dblinks: []*feature.DbLink{
				{
					PrimaryId: "UNIPROT:Q54G86",
					Database:  "UniProt",
					Linktype:  "protein",
					Url:       "https://www.uniprot.org/uniprot/Q54G86",
				},
			},
		},
		CreatedBy:  "curator@dictybase.org",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createPakAAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: "gene",
		Id:   "DDB_G0282525",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "pakA",
			Synonyms: []string{"p21-activated kinase A", "pak1"},
			Publications: []string{
				"10.1083/jcb.2023.202301056",
				"10.1242/dev.2023.200123",
			},
			Pubmed: []string{"33445566", "77889900"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "function",
					Value:     "protein serine/threonine kinase activity",
					CreatedBy: "admin@dictybase.org",
					CreatedAt: now,
				},
				{
					Tag:       "regulation",
					Value:     "positive regulation of cell migration",
					CreatedBy: "admin@dictybase.org",
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  "admin@dictybase.org",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createRasGAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: "gene",
		Id:   "DDB_G0283471",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "rasG",
			Synonyms: []string{"ras protein G", "ras-like GTPase"},
			Publications: []string{
				"10.1016/j.devcel.2023.04.015",
			},
			Pubmed: []string{"55667788"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "function",
					Value:     "GTPase activity",
					CreatedBy: "test@dictybase.org",
					CreatedAt: now,
				},
				{
					Tag:       "domain",
					Value:     "Ras family",
					CreatedBy: "test@dictybase.org",
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  "test@dictybase.org",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createDiscoidin1Annotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: "gene",
		Id:   "DDB_G0291234",
		Attributes: &feature.FeatureAnnotationAttributes{
			Name:     "discoidin1",
			Synonyms: []string{"disc1", "lectin"},
			Publications: []string{
				"10.1371/journal.pone.2023.0123456",
			},
			Pubmed: []string{"99887766"},
			Properties: []*feature.TagProperty{
				{
					Tag:       "function",
					Value:     "carbohydrate binding",
					CreatedBy: "curator@dictybase.org",
					CreatedAt: now,
				},
				{
					Tag:       "expression",
					Value:     "developmentally regulated",
					CreatedBy: "curator@dictybase.org",
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  "curator@dictybase.org",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

// GenerateRandomFeatureAnnotation creates a single random feature annotation for testing
func GenerateRandomFeatureAnnotation() *feature.FeatureAnnotation {
	id := fmt.Sprintf("DDB_G%07d", time.Now().UnixNano()%10000000)
	geneNames := []string{"geneA", "geneB", "geneC", "testGene", "mockGene"}
	functions := []string{"protein binding", "catalytic activity", "transcription factor", "enzyme activity"}

	now := timestamppb.New(time.Now())

	return &feature.FeatureAnnotation{
		Type: "gene",
		Id:   id,
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: geneNames[time.Now().UnixNano()%int64(len(geneNames))],
			Publications: []string{
				fmt.Sprintf("10.1000/journal.%d", time.Now().UnixNano()%100000),
			},
			Pubmed: []string{
				fmt.Sprintf("%d", 10000000+time.Now().UnixNano()%90000000),
			},
			Properties: []*feature.TagProperty{
				{
					Tag:       "function",
					Value:     functions[time.Now().UnixNano()%int64(len(functions))],
					CreatedBy: "test@dictybase.org",
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  "test@dictybase.org",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

// ValidEmails returns a list of valid email addresses for testing
func ValidEmails() []string {
	return []string{
		"test@dictybase.org",
		"curator@dictybase.org",
		"admin@dictybase.org",
		"user@dictybase.org",
		"researcher@dictybase.org",
	}
}

// ValidDOIs returns a list of valid DOI patterns for testing
func ValidDOIs() []string {
	return []string{
		"10.1016/j.cell.2023.001234",
		"10.1038/nature.2023.5678",
		"10.1074/jbc.2023.298765",
		"10.1083/jcb.2023.202301056",
		"10.1242/dev.2023.200123",
		"10.1371/journal.pone.2023.0123456",
	}
}

// ValidPubmedIDs returns a list of valid PubMed IDs for testing
func ValidPubmedIDs() []string {
	return []string{
		"12345678",
		"87654321",
		"11223344",
		"33445566",
		"77889900",
		"55667788",
		"99887766",
	}
}
