package strain

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	TestFile  = "strain_annotation.xlsx"
	TestSheet = "Strain_Annotations"
)

func strainTestFile() (string, error) {
	var empty string
	dir, err := os.Getwd()
	if err != nil {
		return empty, fmt.Errorf(
			"error in getting current working directory %s",
			err,
		)
	}

	return filepath.Join(dir, "../../../../testdata/", TestFile), nil
}

func TestNewStrainAnnotationReader(t *testing.T) {
	assert := require.New(t)
	strainFile, err := strainTestFile()
	assert.NoError(err)
	reader, err := NewStrainAnnotationReader(strainFile, TestSheet, time.Now())
	assert.NoError(err)
	assert.NotNil(reader)
}

func TestValue(t *testing.T) {
	assert := require.New(t)
	strainFile, err := strainTestFile()
	assert.NoError(err)
	t.Run("should have first row with expected values", func(t *testing.T) {
		reader, err := NewStrainAnnotationReader(
			strainFile,
			TestSheet,
			time.Now(),
		)
		assert.NoError(err)
		assert.True(reader.Next())
		anno, err := reader.Value()
		assert.NoError(err)
		assert.False(anno.IsEmpty())
		assert.Equal("far1-", anno.Descriptor())
		assert.Regexp("far1 null mutant", anno.Summary())
		assert.Equal(
			"DDSTRAINCHAR:0000003,DDSTRAINCHAR:0000013",
			anno.Characteristic(),
		)
		assert.Equal("DDGENMOD:0000003", anno.GeneticModification())
		assert.Equal("DDMUMET:0000010", anno.MutagenesisMethod())
		assert.Equal("DBS0236486", anno.ParentId())
		assert.Equal("DDB_G0281211", anno.Genes())
		assert.Equal("Dictyostelium discoideum", anno.Species())
		assert.Equal("PMID:35916164", anno.Reference())
		assert.Equal("robert-dodson@northwestern.edu", anno.AssignedBy())
		assert.False(anno.HasName())
		assert.False(anno.HasPlasmid())
		assert.True(anno.HasGenotype())
		assert.Equal("axeA1,axeB1,axeC1,far1-,[bsRcas],bsR", anno.Genotype())
		assert.False(anno.HasDepositor())
	})
	t.Run("should detect empty rows", func(t *testing.T) {
		reader, err := NewStrainAnnotationReader(
			strainFile,
			TestSheet,
			time.Now(),
		)
		assert.NoError(err)
		assert.True(reader.Next())
		anno, err := reader.Value()
		assert.NoError(err)
		assert.False(anno.IsEmpty())

		assert.True(reader.Next())
		anno2, err := reader.Value()
		assert.NoError(err)
		assert.False(anno2.IsEmpty())

		assert.True(reader.Next())
		anno3, err := reader.Value()
		assert.NoError(err)
		assert.True(anno3.IsEmpty())
	})
}
