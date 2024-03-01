package database

type PhenotypeTableManager struct {
	*TableManager
}

func (pheno *PhenotypeTableManager) FieldNames() []string {
	return []string{
		"annotation_id",
		"strain_id",
		"strain_descriptor",
		"phenotype_id",
		"phenotype_term",
		"assay_id",
		"assay_term",
		"environment_id",
		"environment_term",
		"reference",
		"deleted",
		"assigned_by",
		"created_on",
		"last_modified",
	}
}

func (pheno *PhenotypeTableManager) FieldChangeSpecs() map[string]map[string]interface{} {
	paramsMap := make(map[string]map[string]interface{})
	paramsMap["annotation_id"] = map[string]interface{}{
		"name": "annotation_id",
		"type": "uuid",
	}
	paramsMap["deleted"] = map[string]interface{}{
		"name": "deleted",
		"type": "boolean",
	}
	paramsMap["last_modified"] = map[string]interface{}{
		"name":              "last_modified",
		"type":              "created_on",
		"date_format":       "US",
		"date_include_time": "true",
		"date_time_format":  "12",
		"date_show_tzinfo":  "true",
	}
	paramsMap["created_on"] = map[string]interface{}{
		"name":              "created_on",
		"type":              "created_on",
		"date_format":       "US",
		"date_include_time": "true",
		"date_time_format":  "12",
		"date_show_tzinfo":  "true",
	}
	paramsMap["assigned_by"] = map[string]interface{}{
		"name": "assigned_by",
		"type": "multiple_collaborators",
	}
	paramsMap["phenotype_term"] = map[string]interface{}{
		"name":    "phenotype_term",
		"type":    "formula",
		"formula": "lookup('phenotype_id','name')",
	}
	paramsMap["environment_term"] = map[string]interface{}{
		"name":    "environment_term",
		"type":    "formula",
		"formula": "lookup('environment_id','name')",
	}
	paramsMap["assay_term"] = map[string]interface{}{
		"name":    "assay_term",
		"type":    "formula",
		"formula": "lookup('assay_id','name')",
	}

	return paramsMap
}
