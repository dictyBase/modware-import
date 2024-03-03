package phenotype

type PhenotypeAnnotation struct {
	strainId         string
	phenotypeId      string `validate:"required_with=strainDescriptor"`
	assayId          string
	environmentId    string
	strainDescriptor string `validate:"required"`
	notes            string
	reference        string `validate:"required_with=phenotypeId"`
	assignedBy       string `validate:"required_with=phenotypeId"`
	deleted          bool
	empty            bool
}

func (pha *PhenotypeAnnotation) AssayId() string {
	return pha.assayId
}

func (pha *PhenotypeAnnotation) HasAssayId() bool {
	return len(pha.assayId) == 0
}

func (pha *PhenotypeAnnotation) HasEnvironmentId() bool {
	return len(pha.environmentId) == 0
}

func (pha *PhenotypeAnnotation) EnvironmentId() string {
	return pha.environmentId
}

func (pha *PhenotypeAnnotation) Notes() string {
	return pha.notes
}

func (pha *PhenotypeAnnotation) Reference() string {
	return pha.reference
}

func (pha *PhenotypeAnnotation) AssignedBy() string {
	return pha.assignedBy
}

func (pha *PhenotypeAnnotation) IsEmpty() bool {
	return pha.empty
}

func (pha *PhenotypeAnnotation) HasStrainId() bool {
	return len(pha.strainId) == 0
}

func (pha *PhenotypeAnnotation) PhenotypeId() string {
	return pha.phenotypeId
}

func (pha *PhenotypeAnnotation) StrainDescriptor() string {
	return pha.strainDescriptor
}

func (pha *PhenotypeAnnotation) StrainId() string {
	return pha.strainId
}
