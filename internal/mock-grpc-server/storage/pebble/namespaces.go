package pebble

import (
	"fmt"
)

// Namespace prefixes for organizing different data types in Pebble
const (
	// Stock documents - protobuf serialized
	stockPrefix = "stock:"

	// JSON indices - for filtering/queries
	indexPrefix = "index:"

	// Type edges - stock classification (strain or plasmid)
	typePrefix = "type:"

	// Parent relationships - strain hierarchy
	parentPrefix = "parent:"

	// Ontology terms - stock classification
	termPrefix = "term:"

	// ID counters - sequential ID generation
	counterPrefix = "counter:"

	// Reverse indices - for common queries
	depositorPrefix = "depositor:"
	speciesPrefix   = "species:"
)

// Counter keys
const (
	strainCounterKey  = counterPrefix + "strain"
	plasmidCounterKey = counterPrefix + "plasmid"
)

// keyBuilder provides methods for constructing namespaced keys
type keyBuilder struct{}

// newKeyBuilder creates a new key builder instance
func newKeyBuilder() keyBuilder {
	return keyBuilder{}
}

// stockKey returns the key for a stock document
func (keyBuilder) stockKey(stockID string) []byte {
	return []byte(stockPrefix + stockID)
}

// indexKey returns the key for a JSON index
func (keyBuilder) indexKey(stockID string) []byte {
	return []byte(indexPrefix + stockID)
}

// typeKey returns the key for stock type classification
func (keyBuilder) typeKey(stockID string) []byte {
	return []byte(typePrefix + stockID)
}

// parentKey returns the key for parent strain relationship
func (keyBuilder) parentKey(strainID string) []byte {
	return []byte(parentPrefix + strainID)
}

// termKey returns the key for ontology term
func (keyBuilder) termKey(stockID string) []byte {
	return []byte(termPrefix + stockID)
}

// depositorKey returns the reverse index key for depositor
func (keyBuilder) depositorKey(depositor, stockID string) []byte {
	return []byte(fmt.Sprintf("%s%s:%s", depositorPrefix, depositor, stockID))
}

// speciesKey returns the reverse index key for species
func (keyBuilder) speciesKey(species, stockID string) []byte {
	return []byte(fmt.Sprintf("%s%s:%s", speciesPrefix, species, stockID))
}

// strainCounterKey returns the key for strain ID counter
func (keyBuilder) strainCounterKey() []byte {
	return []byte(strainCounterKey)
}

// plasmidCounterKey returns the key for plasmid ID counter
func (keyBuilder) plasmidCounterKey() []byte {
	return []byte(plasmidCounterKey)
}
