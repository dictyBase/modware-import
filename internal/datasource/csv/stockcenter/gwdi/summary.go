package gwdi

import (
	"fmt"
	"strings"
)

func summInterMultipleUpDown(orientation string) string {
	var b strings.Builder
	b.WriteString(summInterUpDown(orientation))
	b.WriteString(" This stock contains %s individual mutants")
	return b.String()
}

func summInterMultipleBoth() string {
	var b strings.Builder
	b.WriteString(summInterSingleBoth())
	b.WriteString(" This stock contains %s individual mutants")
	return b.String()
}

func summInterSingleBoth() string {
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) intergenic mutant,")
	b.WriteString(" insertion is within 500 bp of start codons")
	b.WriteString(" of two neighboring genes oriented in opposite direction;")
	b.WriteString(" potentially affected genes are %s and %s,")
	b.WriteString(" insertion at position %s, %s")
	b.WriteString(" %s at genomic sites; %s orientation.")
	return b.String()
}

func summInterUpDown(orientation string) string {
	strand := "Crick"
	if orientation == "downstream" {
		strand = "Watson"
	}
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) intergenic mutant,")
	b.WriteString(" insertion is within 500 bp of start codon;")
	b.WriteString(" nearest gene %s is ")
	b.WriteString(
		fmt.Sprintf(
			"%s of the insertions site (%s strand)",
			orientation, strand,
		))
	b.WriteString(" insertion at position %s, %s")
	b.WriteString(" %s at genomic sites; %s orientation.")
	return b.String()
}

func summInterNoGeneMultiple() string {
	var b strings.Builder
	b.WriteString(summInterNoGeneSingle())
	b.WriteString(" This stock contains %s individual mutants")
	return b.String()
}

func summInterNoGeneSingle() string {
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) intergenic mutant,")
	b.WriteString(" not near a known coding region;")
	b.WriteString(" insertion at position %s, %s,")
	b.WriteString(" %s at genomic sites; %s orientation.")
	return b.String()
}

func summIntraMultiple() string {
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) mutants,")
	b.WriteString(" %s intragenic insertions;")
	b.WriteString(" insertion at position %s, %s,")
	b.WriteString(" %s at genomic sites; %s orientation;")
	b.WriteString(" this stock contains %s individual mutants")
	return b.String()
}

func summaryIntraSingle() string {
	var b strings.Builder
	b.WriteString("Genome Wide Dictyostelium Insertion bank (GWDI) %s mutant;")
	b.WriteString(" insertion at position %s, %s,")
	b.WriteString(" %s at genomic sites; %s orientation.")
	return b.String()
}

func summaryNAMultiple() string {
	var b strings.Builder
	b.WriteString(summaryIntraSingle())
	b.WriteString(" this stock contains %s individual mutants")
	return b.String()
}
