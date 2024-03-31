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
	assert.True(t, annotationWithId.HasAssayId())
	assert.False(t, annotationWithoutId.HasAssayId())
}

func TestEnvironmentId(t *testing.T) {
	annotation := PhenotypeAnnotation{environmentId: "E123"}
	assert.Equal(t, "E123", annotation.EnvironmentId())
}

func TestHasEnvironmentId(t *testing.T) {
	annotationWithId := PhenotypeAnnotation{environmentId: "E123"}
	annotationWithoutId := PhenotypeAnnotation{}
	assert.True(t, annotationWithId.HasEnvironmentId())
	assert.False(t, annotationWithoutId.HasEnvironmentId())
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

func TestHasStrainId(t *testing.T) {
	annotationWithId := PhenotypeAnnotation{strainId: "S123"}
	annotationWithoutId := PhenotypeAnnotation{}
	assert.True(t, annotationWithId.HasStrainId())
	assert.False(t, annotationWithoutId.HasStrainId())
}

func TestPhenotypeId(t *testing.T) {
	annotation := PhenotypeAnnotation{phenotypeId: "P123"}
	assert.Equal(t, "P123", annotation.PhenotypeId())
}

func TestStrainId(t *testing.T) {
	annotation := PhenotypeAnnotation{strainId: "S123"}
	assert.Equal(t, "S123", annotation.StrainId())
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
