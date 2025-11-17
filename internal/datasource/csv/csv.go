// package csv is for managing data from csv file

package csv

import (
	"encoding/csv"
	"io"
)

// Reader is to manage record from csv file
type Reader struct {
	Reader *csv.Reader
	Record []string
	Err    error
}

// Next read the next csv record
func (r *Reader) Next() bool {
	r.Err = nil
	record, err := r.Reader.Read()
	if err == io.EOF {
		return false
	}
	if err != nil {
		r.Err = err
		return true
	}
	if len(r.Record) > 1 {
		r.Record = nil
	}
	r.Record = append(r.Record, record...)
	return true
}
