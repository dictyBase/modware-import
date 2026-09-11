package mock

import (
	"fmt"
	"time"

	feature "github.com/dictyBase/go-genproto/dictybaseapis/feature_annotation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	randomGeneIDModulo = 10000000 // Modulo value for generating random gene IDs
	journalIDModulo    = 100000   // Modulo value for journal IDs
	basePubmedID       = 10000000 // Base value for PubMed IDs
	pubmedIDModulo     = 90000000 // Modulo value for PubMed IDs
)

// String constants for mock annotation attributes.
const (
	geneType         = "gene"
	functionTag      = "function"
	testCreatorEmail = "test@dictybase.org"
	curatorEmail     = "curator@dictybase.org"
	adminEmail       = "admin@dictybase.org"
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
		Type: geneType,
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
					Tag:       functionTag,
					Value:     "cytoskeleton organization",
					CreatedBy: testCreatorEmail,
					CreatedAt: now,
				},
				{
					Tag:       "location",
					Value:     "cytoplasm",
					CreatedBy: testCreatorEmail,
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
		CreatedBy:  testCreatorEmail,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createMyoBAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: geneType,
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
					Tag:       functionTag,
					Value:     "motor activity",
					CreatedBy: curatorEmail,
					CreatedAt: now,
				},
				{
					Tag:       "pathway",
					Value:     "cell motility",
					CreatedBy: curatorEmail,
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
		CreatedBy:  curatorEmail,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createPakAAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: geneType,
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
					Tag:       functionTag,
					Value:     "protein serine/threonine kinase activity",
					CreatedBy: adminEmail,
					CreatedAt: now,
				},
				{
					Tag:       "regulation",
					Value:     "positive regulation of cell migration",
					CreatedBy: adminEmail,
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  adminEmail,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createRasGAnnotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: geneType,
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
					Tag:       functionTag,
					Value:     "GTPase activity",
					CreatedBy: testCreatorEmail,
					CreatedAt: now,
				},
				{
					Tag:       "domain",
					Value:     "Ras family",
					CreatedBy: testCreatorEmail,
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  testCreatorEmail,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

func createDiscoidin1Annotation(now *timestamppb.Timestamp) *feature.FeatureAnnotation {
	return &feature.FeatureAnnotation{
		Type: geneType,
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
					Tag:       functionTag,
					Value:     "carbohydrate binding",
					CreatedBy: curatorEmail,
					CreatedAt: now,
				},
				{
					Tag:       "expression",
					Value:     "developmentally regulated",
					CreatedBy: curatorEmail,
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  curatorEmail,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

// GenerateRandomFeatureAnnotation creates a single random feature annotation for testing
func GenerateRandomFeatureAnnotation() *feature.FeatureAnnotation {
	id := fmt.Sprintf("DDB_G%07d", time.Now().UnixNano()%randomGeneIDModulo)
	geneNames := []string{"geneA", "geneB", "geneC", "testGene", "mockGene"}
	functions := []string{
		"protein binding",
		"catalytic activity",
		"transcription factor",
		"enzyme activity",
	}

	now := timestamppb.New(time.Now())

	return &feature.FeatureAnnotation{
		Type: geneType,
		Id:   id,
		Attributes: &feature.FeatureAnnotationAttributes{
			Name: geneNames[time.Now().UnixNano()%int64(len(geneNames))],
			Publications: []string{
				fmt.Sprintf("10.1000/journal.%d", time.Now().UnixNano()%journalIDModulo),
			},
			Pubmed: []string{
				fmt.Sprintf("%d", basePubmedID+time.Now().UnixNano()%pubmedIDModulo),
			},
			Properties: []*feature.TagProperty{
				{
					Tag:       functionTag,
					Value:     functions[time.Now().UnixNano()%int64(len(functions))],
					CreatedBy: testCreatorEmail,
					CreatedAt: now,
				},
			},
		},
		CreatedBy:  testCreatorEmail,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsObsolete: false,
	}
}

// ValidEmails returns a list of valid email addresses for testing
func ValidEmails() []string {
	return []string{
		testCreatorEmail,
		curatorEmail,
		adminEmail,
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
