package tsv

import (
	"bufio"
	"strings"
)

// Reader is to manage record from tsv file
type Reader struct {
	Reader *bufio.Scanner
	Record []string
	Err    error
}

// Next read the next tsv record
func (r *Reader) Next() bool {
	resp := r.Reader.Scan()
	if err := r.Reader.Err(); err != nil {
		r.Err = err
		return false
	}
	if !resp {
		return resp
	}
	r.Record = strings.Split(r.Reader.Text(), "\t")
	return resp
}
