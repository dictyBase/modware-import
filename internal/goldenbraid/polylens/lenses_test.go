package polylens_test

import (
	"testing"

	"github.com/bebop/poly/io/genbank"
	"github.com/stretchr/testify/require"

	"github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
)

func TestFeatureLocationOwnsNestedSlices(t *testing.T) {
	original := genbank.Location{SubLocations: []genbank.Location{
		{Start: 10, SubLocations: []genbank.Location{{Start: 20}}},
	}}
	feature := polylens.FeatureLocation.Set(original)(genbank.Feature{})
	feature.Location.SubLocations[0].Start = 99
	feature.Location.SubLocations[0].SubLocations[0].Start = 99
	require.Equal(t, 10, original.SubLocations[0].Start)
	require.Equal(t, 20, original.SubLocations[0].SubLocations[0].Start)
}

func TestGenbankMetaOwnsOtherMap(t *testing.T) {
	original := genbank.Meta{Other: map[string]string{"comment": "original"}}
	record := polylens.GenbankMeta.Set(original)(genbank.Genbank{})
	record.Meta.Other["comment"] = "changed"
	require.Equal(t, "original", original.Other["comment"])
}

func TestGenbankFeaturesOwnQualifiersAndLocations(t *testing.T) {
	original := []genbank.Feature{{
		Attributes: map[string]string{"label": "original"},
		Location:   genbank.Location{SubLocations: []genbank.Location{{Start: 10}}},
	}}
	record := polylens.GenbankFeatures.Set(original)(genbank.Genbank{})
	record.Features[0].Attributes["label"] = "changed"
	record.Features[0].Location.SubLocations[0].Start = 99
	record.Features[0].Type = "CDS"
	require.Equal(t, "original", original[0].Attributes["label"])
	require.Equal(t, 10, original[0].Location.SubLocations[0].Start)
	require.Empty(t, original[0].Type)
}

func TestQualifierAndSubLocationSettersOwnValues(t *testing.T) {
	attrs := map[string]string{"label": "original"}
	feature := polylens.FeatureAttributes.Set(attrs)(genbank.Feature{})
	feature.Attributes["label"] = "changed"
	require.Equal(t, "original", attrs["label"])
	other := map[string]string{"comment": "original"}
	meta := polylens.MetaOther.Set(other)(genbank.Meta{})
	meta.Other["comment"] = "changed"
	require.Equal(t, "original", other["comment"])
	locations := []genbank.Location{{SubLocations: []genbank.Location{{Start: 10}}}}
	location := polylens.LocationSubLocations.Set(locations)(genbank.Location{})
	location.SubLocations[0].SubLocations[0].Start = 99
	require.Equal(t, 10, locations[0].SubLocations[0].Start)
}
