package database

type PhenotypeTableManager struct {
	*TableManager
}

func (pheno *PhenotypeTableManager) FieldNames() []string {
	return []string{
		"strain_id",
		"strain_descriptor",
		"phenotype_term",
		"assay_term",
		"environment_term",
		"reference",
		"deleted",
		"assigned_by",
		"created_on",
		"last_modified",
	}
}
