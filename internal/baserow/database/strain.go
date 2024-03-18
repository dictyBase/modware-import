package database

type StrainTableManager struct {
	*TableManager
}

func (strn *StrainTableManager) FieldNames() []string {
	return []string{
		"annotation_id",
		"strain_id",
		"strain_names",
		"strain_descriptor",
		"strain_summary",
		"systematic_name",
		"strain_characteristics_id",
		"strain_characteristics",
		"genetic_modification_id",
		"genetic_modification",
		"mutagenesis_method_id",
		"mutagenesis_method",
		"plasmid",
		"parent_strain_id",
		"associated_genes",
		"genotype",
		"depositor",
		"species",
		"reference",
		"assigned_by",
		"created_on",
		"last_modified",
	}
}

func (strn *StrainTableManager) LinkFieldChangeSpecs(
	idMaps map[string]int,
) map[string]map[string]interface{} {
	paramsMap := make(map[string]map[string]interface{})
	paramsMap["strain_characteristics_id"] = map[string]interface{}{
		"name":              "strain_characteristics_id",
		"type":              "link_row",
		"link_row_table_id": idMaps["strainchar-ontology-table"],
	}
	paramsMap["genetic_modification_id"] = map[string]interface{}{
		"name":              "genetic_modification_id",
		"type":              "link_row",
		"link_row_table_id": idMaps["genetic-mod-ontology-table"],
	}
	paramsMap["mutagenesis_method_id"] = map[string]interface{}{
		"name":              "mutagenesis_method_id",
		"type":              "link_row",
		"link_row_table_id": idMaps["mutagenesis-method-ontology-table"],
	}

	return paramsMap
}

func (strn *StrainTableManager) FieldChangeSpecs() map[string]map[string]interface{} {
	paramsMap := make(map[string]map[string]interface{})
	paramsMap["annotation_id"] = map[string]interface{}{
		"name": "annotation_id",
		"type": "uuid",
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
		"type":              "date",
		"date_format":       "US",
		"date_include_time": "true",
		"date_time_format":  "12",
		"date_show_tzinfo":  "true",
	}
	paramsMap["assigned_by"] = map[string]interface{}{
		"name": "assigned_by",
		"type": "multiple_collaborators",
	}
	paramsMap["strain_characteristics"] = map[string]interface{}{
		"name":    "strain_characteristics",
		"type":    "formula",
		"formula": "lookup('strain_characteristics_id','name')",
	}
	paramsMap["mutagenesis_method"] = map[string]interface{}{
		"name":    "mutagenesis_method",
		"type":    "formula",
		"formula": "lookup('mutagenesis_method_id','name')",
	}
	paramsMap["genetic_modification"] = map[string]interface{}{
		"name":    "genetic_modification",
		"type":    "formula",
		"formula": "lookup('genetic_modification_id','name')",
	}
	paramsMap["depositor"] = map[string]interface{}{
		"name": "depositor",
		"type": "email",
	}

	return paramsMap
}
