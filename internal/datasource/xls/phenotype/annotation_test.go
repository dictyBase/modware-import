package phenotype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssayID(t *testing.T) {
	annotation := PhenotypeAnnotation{assayID: "A123"}
	assert.Equal(t, "A123", annotation.AssayID())
}

func TestHasAssayID(t *testing.T) {
	annotationWithID := PhenotypeAnnotation{assayID: "DDASSAY_4893439"}
	annotationWithoutID := PhenotypeAnnotation{}
	assert.True(t, annotationWithID.HasAssayID())
	assert.False(t, annotationWithoutID.HasAssayID())
}

func TestEnvironmentID(t *testing.T) {
	annotation := PhenotypeAnnotation{environmentID: "E123"}
	assert.Equal(t, "E123", annotation.EnvironmentID())
}

func TestHasEnvironmentID(t *testing.T) {
	annotationWithID := PhenotypeAnnotation{environmentID: "E123"}
	annotationWithoutID := PhenotypeAnnotation{}
	assert.True(t, annotationWithID.HasEnvironmentID())
	assert.False(t, annotationWithoutID.HasEnvironmentID())
}

func TestNotes(t *testing.T) {
	annotation := PhenotypeAnnotation{notes: "This is a note."}
	assert.Equal(t, "This is a note.", annotation.Notes())
}

func TestHasNotes(t *testing.T) {
	annotationWithNotes := PhenotypeAnnotation{notes: "This is a note."}
	annotationWithoutNotes := PhenotypeAnnotation{}
	assert.True(t, annotationWithNotes.HasNotes())
	assert.False(t, annotationWithoutNotes.HasNotes())
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

func TestHasStrainID(t *testing.T) {
	annotationWithID := PhenotypeAnnotation{strainID: "S123"}
	annotationWithoutID := PhenotypeAnnotation{}
	assert.True(t, annotationWithID.HasStrainID())
	assert.False(t, annotationWithoutID.HasStrainID())
}

func TestPhenotypeID(t *testing.T) {
	annotation := PhenotypeAnnotation{phenotypeID: "P123"}
	assert.Equal(t, "P123", annotation.PhenotypeID())
}

func TestStrainID(t *testing.T) {
	annotation := PhenotypeAnnotation{strainID: "S123"}
	assert.Equal(t, "S123", annotation.StrainID())
}

func TestStrainDescriptor(t *testing.T) {
	tests := []struct {
		name             string
		strainDescriptor string
		want             string
	}{
		{"Valid strain descriptor", "ABCD-1234", "ABCD-1234"},
		{"Empty strain descriptor", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pha := PhenotypeAnnotation{
				strainDescriptor: tt.strainDescriptor,
			}
			assert.Equal(t, tt.want, pha.StrainDescriptor())
		})
	}
}

func TestHasStrainDescriptor(t *testing.T) {
	annotationWithDescriptor := PhenotypeAnnotation{strainDescriptor: "ABC123"}
	annotationWithoutDescriptor := PhenotypeAnnotation{}
	assert.True(t, annotationWithDescriptor.HasStrainDescriptor())
	assert.False(t, annotationWithoutDescriptor.HasStrainDescriptor())
}
