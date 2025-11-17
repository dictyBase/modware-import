// Package strain provides structures and functions to work with
// biological strain annotations. It defines the Annotation
// type which encapsulates various attributes related to a strain
// such as its descriptor, species, genetic modifications, and
// other metadata. The package allows checking the presence of
// specific attributes and retrieving their values.
package strain

type Annotation struct {
	descriptor          string `validate:"required"`
	species             string `validate:"required"`
	assignedBy          string `validate:"required"`
	reference           string `validate:"required_with=descriptor"`
	summary             string `validate:"required_with=descriptor"`
	characteristic      string `validate:"required_with=descriptor"`
	geneticModification string `validate:"required_with=characteristic"`
	mutagenesisMethod   string
	id                  string
	name                string
	systematicName      string
	plasmid             string
	parentID            string
	genes               string
	genotype            string
	depositor           string
	empty               bool
}

func (strain *Annotation) IsEmpty() bool {
	return strain.empty
}

func (strain *Annotation) HasID() bool {
	return len(strain.id) > 0
}

func (strain *Annotation) ID() string {
	return strain.id
}

func (strain *Annotation) HasName() bool {
	return len(strain.name) > 0
}

func (strain *Annotation) Name() string {
	return strain.name
}

func (strain *Annotation) HasSystematicName() bool {
	return len(strain.systematicName) > 0
}

func (strain *Annotation) SystematicName() string {
	return strain.systematicName
}

func (strain *Annotation) HasPlasmid() bool {
	return len(strain.plasmid) > 0
}

func (strain *Annotation) Plasmid() string {
	return strain.plasmid
}

func (strain *Annotation) HasParentID() bool {
	return len(strain.parentID) > 0
}

func (strain *Annotation) ParentID() string {
	return strain.parentID
}

func (strain *Annotation) Descriptor() string {
	return strain.descriptor
}

func (strain *Annotation) Species() string {
	return strain.species
}

func (strain *Annotation) AssignedBy() string {
	return strain.assignedBy
}

func (strain *Annotation) Reference() string {
	return strain.reference
}

func (strain *Annotation) Summary() string {
	return strain.summary
}

func (strain *Annotation) Characteristic() string {
	return strain.characteristic
}

func (strain *Annotation) GeneticModification() string {
	return strain.geneticModification
}

func (strain *Annotation) MutagenesisMethod() string {
	return strain.mutagenesisMethod
}

func (strain *Annotation) HasGenes() bool {
	return len(strain.genes) > 0
}

func (strain *Annotation) Genes() string {
	return strain.genes
}

func (strain *Annotation) HasGenotype() bool {
	return len(strain.genotype) > 0
}

func (strain *Annotation) Genotype() string {
	return strain.genotype
}

func (strain *Annotation) HasDepositor() bool {
	return len(strain.depositor) > 0
}

func (strain *Annotation) Depositor() string {
	return strain.depositor
}
