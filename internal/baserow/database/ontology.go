package database

type OntologyTableManager struct {
	*TableManager
}

func (ont *OntologyTableManager) FieldNames() []string {
	return []string{"term_id", "name", "is_obsolete"}
}

func (ont *OntologyTableManager) FieldDefs() []map[string]interface{} {
	return []map[string]interface{}{
		{"name": "name", "type": "text"},
		{"name": "term_id", "type": "text"},
		{"name": "is_obsolete", "type": "boolean"},
	}
}

func (ont *OntologyTableManager) TabledateParams() map[string]map[string]interface{} {
	paramsMap := make(map[string]map[string]interface{})
	paramsMap["is_obsolete"] = map[string]interface{}{
		"name": "is_obsolete",
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

	return paramsMap
}
