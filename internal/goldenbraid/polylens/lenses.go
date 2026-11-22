// Package polylens provides immutable lenses for external poly records.
package polylens

import (
	"maps"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/bebop/poly/io/fasta"
	"github.com/bebop/poly/io/genbank"
)

var (
	// GenbankMeta focuses on Genbank.Meta.
	GenbankMeta = L.MakeLens(
		func(s genbank.Genbank) genbank.Meta { return s.Meta },
		func(s genbank.Genbank, value genbank.Meta) genbank.Genbank {
			value.Other = maps.Clone(value.Other)
			s.Meta = value
			return s
		},
	)
	// GenbankFeatures focuses on Genbank.Features.
	GenbankFeatures = L.MakeLens(
		func(s genbank.Genbank) []genbank.Feature { return s.Features },
		func(s genbank.Genbank, value []genbank.Feature) genbank.Genbank {
			s.Features = cloneFeatures(value)
			return s
		},
	)
	// GenbankSequence focuses on Genbank.Sequence.
	GenbankSequence = L.MakeLens(
		func(s genbank.Genbank) string { return s.Sequence },
		func(s genbank.Genbank, value string) genbank.Genbank { s.Sequence = value; return s },
	)
	// MetaName focuses on Meta.Name.
	MetaName = L.MakeLens(
		func(s genbank.Meta) string { return s.Name },
		func(s genbank.Meta, value string) genbank.Meta { s.Name = value; return s },
	)
	// MetaDefinition focuses on Meta.Definition.
	MetaDefinition = L.MakeLens(
		func(s genbank.Meta) string { return s.Definition },
		func(s genbank.Meta, value string) genbank.Meta { s.Definition = value; return s },
	)
	// MetaOther focuses on Meta.Other.
	MetaOther = L.MakeLens(
		func(s genbank.Meta) map[string]string { return s.Other },
		func(s genbank.Meta, value map[string]string) genbank.Meta { s.Other = maps.Clone(value); return s },
	)
	// MetaLocus focuses on Meta.Locus.
	MetaLocus = L.MakeLens(
		func(s genbank.Meta) genbank.Locus { return s.Locus },
		func(s genbank.Meta, value genbank.Locus) genbank.Meta { s.Locus = value; return s },
	)
	// LocusName focuses on Locus.Name.
	LocusName = L.MakeLens(
		func(s genbank.Locus) string { return s.Name },
		func(s genbank.Locus, value string) genbank.Locus { s.Name = value; return s },
	)
	// LocusSequenceLength focuses on Locus.SequenceLength.
	LocusSequenceLength = L.MakeLens(
		func(s genbank.Locus) string { return s.SequenceLength },
		func(s genbank.Locus, value string) genbank.Locus { s.SequenceLength = value; return s },
	)
	// LocusMoleculeType focuses on Locus.MoleculeType.
	LocusMoleculeType = L.MakeLens(
		func(s genbank.Locus) string { return s.MoleculeType },
		func(s genbank.Locus, value string) genbank.Locus { s.MoleculeType = value; return s },
	)
	// LocusGenbankDivision focuses on Locus.GenbankDivision.
	LocusGenbankDivision = L.MakeLens(
		func(s genbank.Locus) string { return s.GenbankDivision },
		func(s genbank.Locus, value string) genbank.Locus { s.GenbankDivision = value; return s },
	)
	// LocusModificationDate focuses on Locus.ModificationDate.
	LocusModificationDate = L.MakeLens(
		func(s genbank.Locus) string { return s.ModificationDate },
		func(s genbank.Locus, value string) genbank.Locus { s.ModificationDate = value; return s },
	)
	// LocusSequenceCoding focuses on Locus.SequenceCoding.
	LocusSequenceCoding = L.MakeLens(
		func(s genbank.Locus) string { return s.SequenceCoding },
		func(s genbank.Locus, value string) genbank.Locus { s.SequenceCoding = value; return s },
	)
	// LocusCircular focuses on Locus.Circular.
	LocusCircular = L.MakeLens(
		func(s genbank.Locus) bool { return s.Circular },
		func(s genbank.Locus, value bool) genbank.Locus { s.Circular = value; return s },
	)
	// FeatureType focuses on Feature.Type.
	FeatureType = L.MakeLens(
		func(s genbank.Feature) string { return s.Type },
		func(s genbank.Feature, value string) genbank.Feature { s.Type = value; return s },
	)
	// FeatureDescription focuses on Feature.Description.
	FeatureDescription = L.MakeLens(
		func(s genbank.Feature) string { return s.Description },
		func(s genbank.Feature, value string) genbank.Feature { s.Description = value; return s },
	)
	// FeatureAttributes focuses on Feature.Attributes.
	FeatureAttributes = L.MakeLens(
		func(s genbank.Feature) map[string]string { return s.Attributes },
		func(s genbank.Feature, value map[string]string) genbank.Feature {
			s.Attributes = maps.Clone(value)
			return s
		},
	)
	// FeatureLocation focuses on Feature.Location.
	FeatureLocation = L.MakeLens(
		func(s genbank.Feature) genbank.Location { return s.Location },
		func(s genbank.Feature, value genbank.Location) genbank.Feature {
			s.Location = cloneLocation(value)
			return s
		},
	)
	// LocationStart focuses on Location.Start.
	LocationStart = L.MakeLens(
		func(s genbank.Location) int { return s.Start },
		func(s genbank.Location, value int) genbank.Location { s.Start = value; return s },
	)
	// LocationEnd focuses on Location.End.
	LocationEnd = L.MakeLens(
		func(s genbank.Location) int { return s.End },
		func(s genbank.Location, value int) genbank.Location { s.End = value; return s },
	)
	// LocationComplement focuses on Location.Complement.
	LocationComplement = L.MakeLens(
		func(s genbank.Location) bool { return s.Complement },
		func(s genbank.Location, value bool) genbank.Location { s.Complement = value; return s },
	)
	// LocationJoin focuses on Location.Join.
	LocationJoin = L.MakeLens(
		func(s genbank.Location) bool { return s.Join },
		func(s genbank.Location, value bool) genbank.Location { s.Join = value; return s },
	)
	// LocationFivePrimePartial focuses on Location.FivePrimePartial.
	LocationFivePrimePartial = L.MakeLens(
		func(s genbank.Location) bool { return s.FivePrimePartial },
		func(s genbank.Location, value bool) genbank.Location { s.FivePrimePartial = value; return s },
	)
	// LocationThreePrimePartial focuses on Location.ThreePrimePartial.
	LocationThreePrimePartial = L.MakeLens(
		func(s genbank.Location) bool { return s.ThreePrimePartial },
		func(s genbank.Location, value bool) genbank.Location { s.ThreePrimePartial = value; return s },
	)
	// LocationGbkLocationString focuses on Location.GbkLocationString.
	LocationGbkLocationString = L.MakeLens(
		func(s genbank.Location) string { return s.GbkLocationString },
		func(s genbank.Location, value string) genbank.Location { s.GbkLocationString = value; return s },
	)
	// LocationSubLocations focuses on Location.SubLocations.
	LocationSubLocations = L.MakeLens(
		func(s genbank.Location) []genbank.Location { return s.SubLocations },
		func(s genbank.Location, value []genbank.Location) genbank.Location {
			s.SubLocations = cloneLocations(value)
			return s
		},
	)
	// FastaName focuses on Fasta.Name.
	FastaName = L.MakeLens(
		func(s fasta.Fasta) string { return s.Name },
		func(s fasta.Fasta, value string) fasta.Fasta { s.Name = value; return s },
	)
	// FastaSequence focuses on Fasta.Sequence.
	FastaSequence = L.MakeLens(
		func(s fasta.Fasta) string { return s.Sequence },
		func(s fasta.Fasta, value string) fasta.Fasta { s.Sequence = value; return s },
	)
)

// Raw focus breaks recursive initialization: the public setter clones children.
var rawSubLocations = L.MakeLens(
	func(s genbank.Location) []genbank.Location { return s.SubLocations },
	func(s genbank.Location, value []genbank.Location) genbank.Location { s.SubLocations = value; return s },
)
var cloneFeatures = A.Map(cloneFeature)

func cloneFeature(f genbank.Feature) genbank.Feature {
	return F.Pipe2(f,
		FeatureAttributes.Modify(F.Identity[map[string]string]),
		FeatureLocation.Modify(F.Identity[genbank.Location]))
}

func cloneLocation(l genbank.Location) genbank.Location {
	return F.Pipe1(l, rawSubLocations.Modify(cloneLocations))
}

func cloneLocations(locations []genbank.Location) []genbank.Location {
	return F.Pipe1(locations, A.Map(cloneLocation))
}
