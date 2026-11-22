package database

type OntologyTableManager struct {
	*TableManager
}

func (ont *OntologyTableManager) FieldNames() []string {
	return []string{
		"term_id",
		"name",
		"is_obsolete",
		"created_on",
		"last_modified",
	}
}

func (ont *OntologyTableManager) FieldDefs() []map[string]any {
	return []map[string]any{
		{"name": "name", "type": "text"},
		{"name": "term_id", "type": "text"},
		{"name": "is_obsolete", "type": "boolean"},
	}
}

func (ont *OntologyTableManager) FieldChangeSpecs() map[string]map[string]any {
	paramsMap := make(map[string]map[string]any)
	paramsMap["is_obsolete"] = map[string]any{
		"name": "is_obsolete",
		"type": "boolean",
	}
	paramsMap["last_modified"] = map[string]any{
		"name":              "last_modified",
		"type":              "created_on",
		"date_format":       "US",
		"date_include_time": "true",
		"date_time_format":  "12",
		"date_show_tzinfo":  "true",
	}
	paramsMap["created_on"] = map[string]any{
		"name":              "created_on",
		"type":              "created_on",
		"date_format":       "US",
		"date_include_time": "true",
		"date_time_format":  "12",
		"date_show_tzinfo":  "true",
	}

	return paramsMap
}
