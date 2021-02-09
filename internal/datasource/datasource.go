//package datasource is for managing data from different source

package datasource

// IteratorWithoutValue is an interface for iteration only
type IteratorWithoutValue interface {
	Next() bool
}
