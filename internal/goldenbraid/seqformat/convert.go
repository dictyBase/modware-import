// Package seqformat encodes shared poly records as GenBank and FASTA bytes.
package seqformat

import (
	"bytes"
	"errors"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	PR "github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/bebop/poly/io/fasta"
	"github.com/bebop/poly/io/genbank"
	"github.com/dictyBase/fp-go-loom/strutils"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
)

const (
	featuresHeading      = "FEATURES             Location/Qualifiers\n"
	originHeading        = "ORIGIN\n"
	emptyFeatureBoundary = featuresHeading + originHeading
	featureTypeFormat    = "     %-16s"
	qualifierKeyFormat   = "                     /%s=\""
	lineEnd              = "\n"
	qualifierEnd         = "\"" + lineEnd
)

var (
	// ErrMissingFeatureBoundary identifies an unsupported GenBank builder layout.
	ErrMissingFeatureBoundary = errors.New("GenBank FEATURES/ORIGIN boundary is missing")
	fullName                  = F.Pipe1(PL.GenbankMeta, L.Compose[genbank.Genbank](PL.MetaName))
	qualifierQuotes           = strings.NewReplacer("\"", "\"\"")
	featureBoundary           = []byte(emptyFeatureBoundary)
	concatPair                = PR.Paired(S.Monoid.Concat)
	// The pair holds the encoded features block, then the built GenBank bytes.
	hasFeatureBoundary = F.Flow2(
		PR.Tail[[]byte, []byte],
		F.Bind2nd(bytes.Contains, featureBoundary),
	)
	missingFeatureBoundary = F.Constant1[PR.Pair[[]byte, []byte]](ErrMissingFeatureBoundary)
	buildWithoutFeatures   = F.Flow3(
		PL.GenbankFeatures.Set(nil),
		E.Eitherize1(genbank.Build),
		E.MapLeft[[]byte](ER.OnError("build GenBank")),
	)
	featuresBlock = F.Flow6(
		PL.GenbankFeatures.Get,
		A.Map(featureText),
		A.Fold(S.Monoid),
		S.Prepend(featuresHeading),
		S.Append(originHeading),
		S.ToBytes,
	)
	// One feature is its type column, location line, then sorted qualifier lines.
	featureText = F.Pipe1(
		[]func(genbank.Feature) string{
			F.Flow2(PL.FeatureType.Get, S.Format[string](featureTypeFormat)),
			F.Flow3(PL.FeatureLocation.Get, genbank.BuildLocationString, S.Append(lineEnd)),
			featureQualifiers,
		},
		strutils.ConcatFields[genbank.Feature](S.Monoid),
	)
	featureQualifiers = F.Flow3(
		PL.FeatureAttributes.Get,
		F.Pipe1(qualifierText, R.CollectOrd[string, string](S.Ord)),
		A.Fold(S.Monoid),
	)
	// A qualifier line doubles the embedded quotes of its value.
	qualifierText = PR.Unpaired(F.Flow2(
		PR.BiMap(
			S.Format[string](qualifierKeyFormat),
			F.Flow2(qualifierQuotes.Replace, S.Append(qualifierEnd)),
		),
		concatPair,
	))
	fastaRecord = F.Flow4(
		PR.Of[genbank.Genbank],
		PR.BiMap(
			F.Flow2(fullName.Get, PL.FastaName.Set),
			F.Flow2(PL.GenbankSequence.Get, PL.FastaSequence.Set),
		),
		PR.Paired(En.Monoid[fasta.Fasta]().Concat),
		En.Build[fasta.Fasta],
	)
)

// ToGenBankBytes retains poly headers and bases, with deterministic features.
func ToGenBankBytes(record genbank.Genbank) E.Either[error, []byte] {
	return F.Pipe5(
		record,
		PR.Of[genbank.Genbank],
		PR.BiMap(featuresBlock, buildWithoutFeatures),
		PR.Sequence(E.Map[error, []byte, PR.Pair[[]byte, []byte]]),
		E.Chain(E.FromPredicate(hasFeatureBoundary, missingFeatureBoundary)),
		E.Map[error](PR.Paired(replaceFeatureBoundary)),
	)
}

// replaceFeatureBoundary puts the encoded block in place of the empty FEATURES heading.
func replaceFeatureBoundary(block, data []byte) []byte {
	return bytes.Replace(data, featureBoundary, block, 1)
}

// ToFASTABytes uses the full construct name rather than its short locus name.
func ToFASTABytes(record genbank.Genbank) E.Either[error, []byte] {
	return F.Pipe4(
		record,
		fastaRecord,
		A.Of[fasta.Fasta],
		E.Eitherize1(fasta.Build),
		E.MapLeft[[]byte](ER.OnError("build FASTA")),
	)
}
