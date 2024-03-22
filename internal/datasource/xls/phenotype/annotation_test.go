package phenotype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssayId(t *testing.T) {
	annotation := PhenotypeAnnotation{assayId: "A123"}
	assert.Equal(t, "A123", annotation.AssayId())
}

func TestHasAssayId(t *testing.T) {
	annotationWithId := PhenotypeAnnotation{assayId: "A123"}
	annotationWithoutId := PhenotypeAnnotation{}
	assert.False(t, annotationWithId.HasAssayId())
	assert.True(t, annotationWithoutId.HasAssayId())
}

func TestEnvironmentId(t *testing.T) {
	annotation := PhenotypeAnnotation{environmentId: "E123"}
	assert.Equal(t, "E123", annotation.EnvironmentId())
}

func TestHasEnvironmentId(t *testing.T) {
	annotationWithId := PhenotypeAnnotation{environmentId: "E123"}
	annotationWithoutId := PhenotypeAnnotation{}
	assert.False(t, annotationWithId.HasEnvironmentId())
	assert.True(t, annotationWithoutId.HasEnvironmentId())
}

func TestNotes(t *testing.T) {
	annotation := PhenotypeAnnotation{notes: "This is a note."}
	assert.Equal(t, "This is a note.", annotation.Notes())
}

func TestHasNotes(t *testing.T) {
	annotationWithNotes := PhenotypeAnnotation{notes: "This is a note."}
	annotationWithoutNotes := PhenotypeAnnotation{}
	assert.False(t, annotationWithNotes.HasNotes())
	assert.True(t, annotationWithoutNotes.HasNotes())
}

func TestReference(t *testing.T) {
	annotation := PhenotypeAnnotation{reference: "Ref123"}
	assert.Equal(t, "Ref123", annotation.Reference())
}

func TestAssignedBy(t *testing.T) {
	annotation := PhenotypeAnnotation{assignedBy: "User123"}
	assert.Equal(t, "User123", annotation.AssignedBy())
}

func TestIsEmpty(t *testing.T) {
	emptyAnnotation := PhenotypeAnnotation{empty: true}
	nonEmptyAnnotation := PhenotypeAnnotation{empty: false}
	assert.True(t, emptyAnnotation.IsEmpty())
	assert.False(t, nonEmptyAnnotation.IsEmpty())
}

func TestHasStrainId(t *testing.T) {
	annotationWithId := PhenotypeAnnotation{strainId: "S123"}
	annotationWithoutId := PhenotypeAnnotation{}
	assert.False(t, annotationWithId.HasStrainId())
	assert.True(t, annotationWithoutId.HasStrainId())
}

func TestPhenotypeId(t *testing.T) {
	annotation := PhenotypeAnnotation{phenotypeId: "P123"}
	assert.Equal(t, "P123", annotation.PhenotypeId())
}

func TestStrainId(t *testing.T) {
	annotation := PhenotypeAnnotation{strainId: "S123"}
	assert.Equal(t, "S123", annotation.StrainId())
}
