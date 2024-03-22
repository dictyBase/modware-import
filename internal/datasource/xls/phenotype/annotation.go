// Package phenotype defines the structure and associated methods for handling
// phenotype annotations. Each annotation is represented by the PhenotypeAnnotation
// struct, which includes details like strain ID, phenotype ID, assay ID, environment ID,
// strain descriptor, notes, reference, assigned by, and flags for deletion and emptiness.
// The struct's methods provide accessors and checkers for the various fields.
package phenotype

import "time"

// PhenotypeAnnotation represents annotations related to a phenotype.
type PhenotypeAnnotation struct {
	strainId      string
	phenotypeId   string `validate:"required_with=strainDescriptor"`
	assayId       string
	environmentId string
	notes         string
	reference     string `validate:"required_with=phenotypeId"`
	assignedBy    string `validate:"required_with=phenotypeId"`
	deleted       bool
	empty         bool
	createdOn     time.Time `validate:"required"`
}

func (pha *PhenotypeAnnotation) CreatedOn() time.Time {
	return pha.createdOn
}

// AssayId returns the assay ID associated with the phenotype annotation.
func (pha *PhenotypeAnnotation) AssayId() string {
	return pha.assayId
}

// HasAssayId checks whether an assay ID is associated with the phenotype annotation.
// It returns true if the assay ID is not set.
func (pha *PhenotypeAnnotation) HasAssayId() bool {
	return len(pha.assayId) == 0
}

// HasEnvironmentId checks whether an environment ID is associated with the phenotype annotation.
// It returns true if the environment ID is not set.
func (pha *PhenotypeAnnotation) HasEnvironmentId() bool {
	return len(pha.environmentId) == 0
}

// EnvironmentId returns the environment ID associated with the phenotype annotation.
func (pha *PhenotypeAnnotation) EnvironmentId() string {
	return pha.environmentId
}

func (pha *PhenotypeAnnotation) HasNotes() bool {
	return len(pha.notes) == 0
}

// Notes returns any notes associated with the phenotype annotation.
func (pha *PhenotypeAnnotation) Notes() string {
	return pha.notes
}

// Reference returns the reference associated with the phenotype annotation.
func (pha *PhenotypeAnnotation) Reference() string {
	return pha.reference
}

// AssignedBy returns the identifier of the entity that assigned the phenotype annotation.
func (pha *PhenotypeAnnotation) AssignedBy() string {
	return pha.assignedBy
}

// IsEmpty checks if the phenotype annotation is marked as empty.
// It returns true if the annotation is considered empty.
func (pha *PhenotypeAnnotation) IsEmpty() bool {
	return pha.empty
}

// HasStrainId checks whether a strain ID is associated with the phenotype annotation.
// It returns true if the strain ID is not set.
func (pha *PhenotypeAnnotation) HasStrainId() bool {
	return len(pha.strainId) == 0
}

// PhenotypeId returns the phenotype ID associated with the phenotype annotation.
func (pha *PhenotypeAnnotation) PhenotypeId() string {
	return pha.phenotypeId
}

// StrainId returns the strain ID associated with the phenotype annotation.
func (pha *PhenotypeAnnotation) StrainId() string {
	return pha.strainId
}
