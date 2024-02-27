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
