// Package strain provides structures and functions to work with
// biological strain annotations. It defines the StrainAnnotation
// type which encapsulates various attributes related to a strain
// such as its descriptor, species, genetic modifications, and
// other metadata. The package allows checking the presence of
// specific attributes and retrieving their values.
package strain

type StrainAnnotation struct {
	descriptor          string `validate:"required"`
	species             string `validate:"required"`
	assignedBy          string `validate:"required"`
	reference           string `validate:"required_with=descriptor"`
	summary             string `validate:"required_with=descriptor"`
	characteristic      string `validate:"required_with=descriptor"`
	geneticModification string `validate:"required_with=characteristic"`
	mutagenesisMethod   string `validate:"required_with=geneticModification"`
	id                  string
	name                string
	systematicName      string
	plasmid             string
	parentId            string
	genes               string
	genotype            string
	depositor           string
}

func (strain *StrainAnnotation) HasId() bool {
	return len(strain.id) > 0
}

func (strain *StrainAnnotation) Id() string {
	return strain.id
}

func (strain *StrainAnnotation) HasName() bool {
	return len(strain.name) > 0
}

func (strain *StrainAnnotation) Name() string {
	return strain.name
}

func (strain *StrainAnnotation) HasSystematicName() bool {
	return len(strain.systematicName) > 0
}

func (strain *StrainAnnotation) HasPlasmid() bool {
	return len(strain.plasmid) > 0
}

func (strain *StrainAnnotation) Plasmid() string {
	return strain.plasmid
}

func (strain *StrainAnnotation) HasParentId() bool {
	return len(strain.parentId) > 0
}

func (strain *StrainAnnotation) ParentId() string {
	return strain.parentId
}

func (strain *StrainAnnotation) Descriptor() string {
	return strain.descriptor
}

func (strain *StrainAnnotation) Species() string {
	return strain.species
}

func (strain *StrainAnnotation) AssignedBy() string {
	return strain.assignedBy
}

func (strain *StrainAnnotation) Reference() string {
	return strain.reference
}

func (strain *StrainAnnotation) Summary() string {
	return strain.summary
}

func (strain *StrainAnnotation) Characteristic() string {
	return strain.characteristic
}

func (strain *StrainAnnotation) GeneticModification() string {
	return strain.geneticModification
}

func (strain *StrainAnnotation) MutagenesisMethod() string {
	return strain.mutagenesisMethod
}

func (strain *StrainAnnotation) HasGenes() bool {
	return len(strain.genes) > 0
}

func (strain *StrainAnnotation) Genes() string {
	return strain.genes
}

func (strain *StrainAnnotation) HasGenotype() bool {
	return len(strain.genotype) > 0
}

func (strain *StrainAnnotation) Genotype() string {
	return strain.genotype
}

func (strain *StrainAnnotation) HasDepositor() bool {
	return len(strain.depositor) > 0
}

func (strain *StrainAnnotation) Depositor() string {
	return strain.depositor
}
