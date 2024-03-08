package strain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasId(t *testing.T) {
	s := StrainAnnotation{id: "123"}
	assert.True(t, s.HasId(), "Expected HasId to return true when id is set")
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasId(),
		"Expected HasId to return false when id is not set",
	)
}

func TestId(t *testing.T) {
	s := StrainAnnotation{id: "123"}
	assert.Equal(t, "123", s.Id(), "Expected Id to return the correct id value")
}

func TestHasName(t *testing.T) {
	s := StrainAnnotation{name: "Strain Name"}
	assert.True(
		t,
		s.HasName(),
		"Expected HasName to return true when name is set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasName(),
		"Expected HasName to return false when name is not set",
	)
}

func TestName(t *testing.T) {
	s := StrainAnnotation{name: "Name"}
	assert.Equal(
		t,
		"Name",
		s.Name(),
		"Expected Name to return the correct name value",
	)
}

func TestHasSystematicName(t *testing.T) {
	s := StrainAnnotation{systematicName: "Sys Name"}
	assert.True(
		t,
		s.HasSystematicName(),
		"Expected HasSystematicName to return true when systematicName is set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasSystematicName(),
		"Expected HasSystematicName to return false when systematicName is not set",
	)
}

func TestHasPlasmid(t *testing.T) {
	s := StrainAnnotation{plasmid: "Plasmid Info"}
	assert.True(
		t,
		s.HasPlasmid(),
		"Expected HasPlasmid to return true when plasmid is set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasPlasmid(),
		"Expected HasPlasmid to return false when plasmid is not set",
	)
}

func TestPlasmid(t *testing.T) {
	s := StrainAnnotation{plasmid: "Plasmid Info"}
	assert.Equal(
		t,
		"Plasmid Info",
		s.Plasmid(),
		"Expected Plasmid to return the correct plasmid value",
	)
}

func TestHasParentId(t *testing.T) {
	s := StrainAnnotation{parentId: "Parent123"}
	assert.True(
		t,
		s.HasParentId(),
		"Expected HasParentId to return true when parentId is set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasParentId(),
		"Expected HasParentId to return false when parentId is not set",
	)
}

func TestParentId(t *testing.T) {
	s := StrainAnnotation{parentId: "Parent123"}
	assert.Equal(
		t,
		"Parent123",
		s.ParentId(),
		"Expected ParentId to return the correct parentId value",
	)
}

func TestDescriptor(t *testing.T) {
	s := StrainAnnotation{descriptor: "Descriptor Info"}
	assert.Equal(
		t,
		"Descriptor Info",
		s.Descriptor(),
		"Expected Descriptor to return the correct descriptor value",
	)
}

func TestSpecies(t *testing.T) {
	s := StrainAnnotation{species: "Species Name"}
	assert.Equal(
		t,
		"Species Name",
		s.Species(),
		"Expected Species to return the correct species value",
	)
}

func TestAssignedBy(t *testing.T) {
	s := StrainAnnotation{assignedBy: "Assigned By Entity"}
	assert.Equal(
		t,
		"Assigned By Entity",
		s.AssignedBy(),
		"Expected AssignedBy to return the correct assignedBy value",
	)
}

func TestReference(t *testing.T) {
	s := StrainAnnotation{reference: "Reference Info"}
	assert.Equal(
		t,
		"Reference Info",
		s.Reference(),
		"Expected Reference to return the correct reference value",
	)
}

func TestSummary(t *testing.T) {
	s := StrainAnnotation{summary: "Summary Info"}
	assert.Equal(
		t,
		"Summary Info",
		s.Summary(),
		"Expected Summary to return the correct summary value",
	)
}

func TestCharacteristic(t *testing.T) {
	s := StrainAnnotation{characteristic: "Characteristic Info"}
	assert.Equal(
		t,
		"Characteristic Info",
		s.Characteristic(),
		"Expected Characteristic to return the correct characteristic value",
	)
}

func TestGeneticModification(t *testing.T) {
	s := StrainAnnotation{
		geneticModification: "Genetic Modification Info",
	}
	assert.Equal(
		t,
		"Genetic Modification Info",
		s.GeneticModification(),
		"Expected GeneticModification to return the correct geneticModification value",
	)
}

func TestMutagenesisMethod(t *testing.T) {
	s := StrainAnnotation{mutagenesisMethod: "Mutagenesis Method Info"}
	assert.Equal(
		t,
		"Mutagenesis Method Info",
		s.MutagenesisMethod(),
		"Expected MutagenesisMethod to return the correct mutagenesisMethod value",
	)
}

func TestHasGenes(t *testing.T) {
	s := StrainAnnotation{genes: "Gene1, Gene2"}
	assert.True(
		t,
		s.HasGenes(),
		"Expected HasGenes to return true when genes are set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasGenes(),
		"Expected HasGenes to return false when genes are not set",
	)
}

func TestGenes(t *testing.T) {
	s := StrainAnnotation{genes: "Gene1, Gene2"}
	assert.Equal(
		t,
		"Gene1, Gene2",
		s.Genes(),
		"Expected Genes to return the correct genes value",
	)
}

func TestHasGenotype(t *testing.T) {
	s := StrainAnnotation{genotype: "Genotype Info"}
	assert.True(
		t,
		s.HasGenotype(),
		"Expected HasGenotype to return true when genotype is set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasGenotype(),
		"Expected HasGenotype to return false when genotype is not set",
	)
}

func TestGenotype(t *testing.T) {
	s := StrainAnnotation{genotype: "Genotype Info"}
	assert.Equal(
		t,
		"Genotype Info",
		s.Genotype(),
		"Expected Genotype to return the correct genotype value",
	)
}

func TestHasDepositor(t *testing.T) {
	s := StrainAnnotation{depositor: "Depositor Info"}
	assert.True(
		t,
		s.HasDepositor(),
		"Expected HasDepositor to return true when depositor is set",
	)
	s = StrainAnnotation{}
	assert.False(
		t,
		s.HasDepositor(),
		"Expected HasDepositor to return false when depositor is not set",
	)
}

func TestDepositor(t *testing.T) {
	s := StrainAnnotation{depositor: "Depositor Info"}
	assert.Equal(
		t,
		"Depositor Info",
		s.Depositor(),
		"Expected Depositor to return the correct depositor value",
	)
}
