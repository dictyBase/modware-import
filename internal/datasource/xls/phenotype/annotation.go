// Package phenotype defines the structure and associated methods for handling
// phenotype annotations. Each annotation is represented by the Annotation
// struct, which includes details like strain ID, phenotype ID, assay ID, environment ID,
// strain descriptor, notes, reference, assigned by, and flags for deletion and emptiness.
// The struct's methods provide accessors and checkers for the various fields.
package phenotype

import (
	"regexp"
	"time"
)

var assayRgxp = regexp.MustCompile(`^DDASSAY_\d{6,}$`)

// Annotation represents annotations related to a phenotype.
type Annotation struct {
	strainID         string
	strainDescriptor string
	phenotypeID      string `validate:"required"`
	assayID          string
	environmentID    string
	notes            string
	reference        string `validate:"required_with=phenotypeID"`
	assignedBy       string `validate:"required_with=phenotypeID"`
	deleted          bool
	empty            bool
	createdOn        time.Time `validate:"required"`
}

func (pha *Annotation) CreatedOn() time.Time {
	return pha.createdOn
}

// AssayID returns the assay ID associated with the phenotype annotation.
func (pha *Annotation) AssayID() string {
	return pha.assayID
}

// HasAssayID checks whether an assay ID is associated with the phenotype annotation.
// It returns true if the assay ID is set.
func (pha *Annotation) HasAssayID() bool {
	if len(pha.assayID) == 0 {
		return false
	}
	return assayRgxp.MatchString(pha.assayID)
}

// HasEnvironmentID checks whether an environment ID is associated with the phenotype annotation.
// It returns true if the environment ID is set.
func (pha *Annotation) HasEnvironmentID() bool {
	return len(pha.environmentID) > 0
}

// EnvironmentID returns the environment ID associated with the phenotype annotation.
func (pha *Annotation) EnvironmentID() string {
	return pha.environmentID
}

func (pha *Annotation) HasNotes() bool {
	return len(pha.notes) > 0
}

// Notes returns any notes associated with the phenotype annotation.
func (pha *Annotation) Notes() string {
	return pha.notes
}

// Reference returns the reference associated with the phenotype annotation.
func (pha *Annotation) Reference() string {
	return pha.reference
}

// AssignedBy returns the identifier of the entity that assigned the phenotype annotation.
func (pha *Annotation) AssignedBy() string {
	return pha.assignedBy
}

// IsEmpty checks if the phenotype annotation is marked as empty.
// It returns true if the annotation is considered empty.
func (pha *Annotation) IsEmpty() bool {
	return pha.empty
}

// HasStrainID checks whether a strain ID is associated with the phenotype annotation.
// It returns true if the strain ID is set.
func (pha *Annotation) HasStrainID() bool {
	return len(pha.strainID) > 0
}

// HasStrainDescriptor checks if the Annotation instance has a strainDescriptor defined.
// It returns true if the strainDescriptor is not empty, otherwise false.
func (pha *Annotation) HasStrainDescriptor() bool {
	return len(pha.strainDescriptor) > 0
}

// StrainDescriptor returns the strain descriptor associated with the Annotation.
func (pha *Annotation) StrainDescriptor() string {
	return pha.strainDescriptor
}

// PhenotypeID returns the phenotype ID associated with the phenotype annotation.
func (pha *Annotation) PhenotypeID() string {
	return pha.phenotypeID
}

// StrainID returns the strain ID associated with the phenotype annotation.
func (pha *Annotation) StrainID() string {
	return pha.strainID
}
