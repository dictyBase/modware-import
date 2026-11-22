package seqmodel

import (
	"encoding/hex"
	"regexp"
	"strings"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/string"
)

var (
	colorPattern = regexp.MustCompile(`\A#[0-9A-Fa-f]{6}\z`)
	dnaPattern   = regexp.MustCompile(`\A[ACGTacgt]+\z`)
	colorL       = MakeColorLenses()
)

// ValidateSequence validates ASCII DNA and normalizes its case.
func ValidateSequence(sequence string) E.Either[error, string] {
	return F.Pipe3(sequence,
		E.FromPredicate(S.IsNonEmpty, emptySequence),
		E.Chain(E.FromPredicate(dnaPattern.MatchString, invalidDNABase)),
		E.Map[error](strings.ToUpper))
}

// ParseColor accepts exactly six hexadecimal digits after a hash.
func ParseColor(s string) E.Either[error, Color] {
	return F.Pipe4(s,
		E.FromPredicate(colorPattern.MatchString, invalidColor),
		E.Map[error](F.Bind2nd(strings.TrimPrefix, "#")),
		E.Chain(E.Eitherize1(hex.DecodeString)),
		E.Map[error](colorFromBytes))
}

func colorFromBytes(b []byte) Color {
	return F.Pipe3(Color{}, colorL.R.Set(b[0]), colorL.G.Set(b[1]), colorL.B.Set(b[2]))
}

// Hex returns a lower-case CSS hex color with leading zeros.
func (c Color) Hex() string {
	return F.Pipe2(
		[]byte{colorL.R.Get(c), colorL.G.Get(c), colorL.B.Get(c)},
		hex.EncodeToString,
		S.Prepend("#"),
	)
}
