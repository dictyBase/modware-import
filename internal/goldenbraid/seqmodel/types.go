// Package seqmodel defines shared sequence export domain types.
package seqmodel

import "github.com/dictyBase/fp-go-loom/eithererr"

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

// Color is an RGB feature color.
// fp-go:Lens
type Color struct{ R, G, B uint8 }

// ExportSummary counts attempted exports and preserves errors in source order.
// fp-go:Lens
type ExportSummary struct {
	Total     int
	Succeeded int
	Failed    int
	Errors    []error
}

var (
	invalidColor          = eithererr.ConstErr[string]("invalid color")
	emptyName             = eithererr.ConstErr[string]("empty name")
	emptySequence         = eithererr.ConstErr[string]("empty sequence")
	invalidDNABase        = eithererr.ConstErr[string]("invalid DNA base")
	segmentOutOfBounds    = eithererr.ConstErr[string]("segment out of bounds")
	ErrInvalidColor       = invalidColor("")
	ErrEmptyName          = emptyName("")
	ErrEmptySequence      = emptySequence("")
	ErrInvalidDNABase     = invalidDNABase("")
	ErrSegmentOutOfBounds = segmentOutOfBounds("")
)
