package snapgene

import (
	"encoding/xml"
	"errors"
	"regexp"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/bebop/poly/io/genbank"
	eitherparse "github.com/dictyBase/fp-go-loom/either/parse"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predarrays "github.com/dictyBase/fp-go-loom/predicate/array"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
)

// fp-go:Lens
type (
	featuresXML struct {
		XMLName xml.Name     `xml:"Features"`
		Items   []featureXML `xml:"Feature"`
	}
	// fp-go:Lens
	featureXML struct {
		Name      string       `xml:"name,attr"`
		Type      string       `xml:"type,attr"`
		Direction string       `xml:"directionality,attr"`
		Segments  []segmentXML `xml:"Segment"`
	}
	// fp-go:Lens
	segmentXML struct {
		Range string `xml:"range,attr"`
		Type  string `xml:"type,attr"`
		Color string `xml:"color,attr"`
	}
	// fp-go:Lens
	primersXML struct {
		XMLName xml.Name    `xml:"Primers"`
		Items   []primerXML `xml:"Primer"`
	}
	// fp-go:Lens
	primerXML struct {
		Name  string    `xml:"name,attr"`
		Sites []siteXML `xml:"BindingSite"`
	}
	// fp-go:Lens
	siteXML struct {
		Location   string `xml:"location,attr"`
		Strand     string `xml:"boundStrand,attr"`
		Simplified string `xml:"simplified,attr"`
	}
	// fp-go:Lens
	notesXML struct {
		XMLName     xml.Name `xml:"Notes"`
		Type        string   `xml:"Type"`
		UUID        string   `xml:"UUID"`
		Created     string   `xml:"Created"`
		Modified    string   `xml:"LastModified"`
		Transformed string   `xml:"TransformedInto"`
	}
	// fp-go:Lens
	rangeBounds struct {
		Start  int
		End    int
		Length int
	}
	// fp-go:Lens
	featureState struct {
		XML        featureXML
		Record     genbank.Genbank
		Leaves     []genbank.Location
		Color      O.Option[seqmodel.Color]
		Location   genbank.Location
		Complement bool
	}
	// fp-go:Lens
	siteState struct {
		XML        siteXML
		Bounds     rangeBounds
		Complement bool
	}
)

var (
	featuresL             = MakefeaturesXMLLenses()
	featureL              = MakefeatureXMLLenses()
	segmentL              = MakesegmentXMLLenses()
	primersL              = MakeprimersXMLLenses()
	primerL               = MakeprimerXMLLenses()
	siteL                 = MakesiteXMLLenses()
	notesL                = MakenotesXMLLenses()
	boundsL               = MakerangeBoundsLenses()
	featureStateL         = MakefeatureStateLenses()
	siteStateL            = MakesiteStateLenses()
	rangePattern          = regexp.MustCompile(`^[0-9]+-[0-9]+$`)
	directionPattern      = regexp.MustCompile(`^[0-3]?$`)
	strandPattern         = regexp.MustCompile(`^[01]$`)
	errInvalidRange       = errors.New("invalid SnapGene range")
	errInvalidDirection   = errors.New("invalid SnapGene direction")
	errInvalidType        = errors.New("empty SnapGene feature type")
	standardSegment       = F.Flow2(segmentL.Type.Get, predord.StrEq("standard"))
	fullSite              = F.Flow2(siteL.Simplified.Get, predord.NotEqualStr("1"))
	featureLocationChoice = P.Fold(joinLocation, firstLocation)
	featureRangeChoice    = P.Fold(normalFeatureRange, wrappedFeatureRange)
	intGt                 = ORD.Gt(predord.IntOrd)
	isWrappedRange        = F.Flow2(startPastEnd, intGt(0))
	complementDirection   = MO.Match(
		false,
		[]O.Kleisli[string, bool]{MO.Const(reverseFeature, true)},
	)
	complementStrand = MO.Match(
		false,
		[]O.Kleisli[string, bool]{MO.Const(reversePrimer, true)},
	)
	featureDirection = F.Pipe1(featureStateL.XML, L.Compose[featureState](featureL.Direction))
	featureName      = F.Pipe1(featureStateL.XML, L.Compose[featureState](featureL.Name))
	featureType      = F.Pipe1(featureStateL.XML, L.Compose[featureState](featureL.Type))
	siteStrand       = F.Pipe1(siteStateL.XML, L.Compose[siteState](siteL.Strand))
	reverseFeature   = predord.StrEq("2")
	reversePrimer    = predord.StrEq("1")
	emptyColor       = predord.StrEq("")
	noColor          = predord.StrEq("noColor")
	parseColor       = MO.Match(
		E.Right[error](O.None[seqmodel.Color]()),
		[]O.Kleisli[string, E.Either[error, O.Option[seqmodel.Color]]]{
			MO.Const(emptyColor, E.Right[error](O.None[seqmodel.Color]())),
			MO.Const(noColor, E.Right[error](O.None[seqmodel.Color]())),
			MO.Default(F.Flow2(seqmodel.ParseColor, E.Map[error](O.Some[seqmodel.Color]))),
		},
	)
)

func convertFeature(record genbank.Genbank, feature featureXML) E.Either[error, genbank.Feature] {
	return F.Pipe11(
		featureState{},
		featureStateL.XML.Set(feature),
		featureStateL.Record.Set(record),
		E.Of[error, featureState],
		E.Chain(
			E.FromPredicate(validFeatureDirection, F.Constant1[featureState](errInvalidDirection)),
		),
		E.Chain(E.FromPredicate(validFeatureType, F.Constant1[featureState](errInvalidType))),
		E.Bind(featureStateL.Leaves.Set, featureLeaves),
		E.Chain(
			E.FromPredicate(
				hasFeatureLeaves,
				F.Constant1[featureState](seqmodel.ErrSegmentOutOfBounds),
			),
		),
		E.Bind(featureStateL.Color.Set, featureColor),
		E.Let[error](
			featureStateL.Complement.Set,
			F.Flow2(featureDirection.Get, complementDirection),
		),
		E.Let[error](featureStateL.Location.Set, featureLocation),
		E.Map[error](buildFeature),
	)
}

func validFeatureDirection(s featureState) bool {
	return F.Pipe3(s, featureStateL.XML.Get, featureL.Direction.Get, directionPattern.MatchString)
}

func validFeatureType(s featureState) bool {
	return F.Pipe3(s, featureStateL.XML.Get, featureL.Type.Get, S.IsNonEmpty)
}

func featureLeaves(s featureState) E.Either[error, []genbank.Location] {
	return F.Pipe5(
		s,
		featureStateL.XML.Get,
		featureL.Segments.Get,
		A.Filter(F.Flow2(segmentL.Type.Get, predord.NotEqualStr("gap"))),
		E.TraverseArray(
			F.Flow3(
				F.Bind1st(PR.MakePair[featureState, segmentXML], s),
				PR.MapHead[segmentXML](featureStateL.Record.Get),
				PR.Paired(segmentLocations),
			),
		),
		E.Map[error](A.Flatten[genbank.Location]),
	)
}

func segmentLocations(
	record genbank.Genbank,
	segment segmentXML,
) E.Either[error, []genbank.Location] {
	return F.Pipe3(
		segment,
		segmentL.Range.Get,
		parseRange,
		E.Chain(F.Bind1st(featureRange, record)),
	)
}

func parseRange(span string) E.Either[error, rangeBounds] {
	return F.Pipe3(
		span,
		E.FromPredicate(rangePattern.MatchString, F.Constant1[string](errInvalidRange)),
		E.Map[error](splitRange),
		E.Chain(parseBounds),
	)
}

func splitRange(span string) []string { return strings.Split(span, "-") }

func parseBounds(parts []string) E.Either[error, rangeBounds] {
	return F.Pipe2(parts, E.TraverseArray(eitherparse.ParseInt), E.Map[error](boundsFromEndpoints))
}

func boundsFromEndpoints(parts []int) rangeBounds {
	return F.Pipe2(rangeBounds{}, boundsL.Start.Set(parts[0]), boundsL.End.Set(parts[1]))
}

func featureRange(record genbank.Genbank, bounds rangeBounds) E.Either[error, []genbank.Location] {
	return F.Pipe2(
		bounds,
		E.FromPredicate(
			F.Bind1st(validFeatureRange, record),
			F.Constant1[rangeBounds](seqmodel.ErrSegmentOutOfBounds),
		),
		E.Map[error](F.Bind1st(splitFeatureRange, record)),
	)
}

func validFeatureRange(record genbank.Genbank, bounds rangeBounds) bool {
	sequence := PL.GenbankSequence.Get(record)
	circular := recordCircular.Get(record)
	start := boundsL.Start.Get(bounds)
	end := boundsL.End.Get(bounds)
	return start >= 1 && start <= len(sequence) && end >= 1 && end <= len(sequence) &&
		(start <= end || circular)
}

func splitFeatureRange(record genbank.Genbank, bounds rangeBounds) []genbank.Location {
	return F.Pipe3(PR.MakePair(record, bounds),
		PR.MapHead[rangeBounds](F.Flow3(PL.GenbankSequence.Get, S.Size, boundsL.Length.Set)),
		PR.Merge(En.Read[rangeBounds]), featureRangeChoice(isWrappedRange))
}

// startPastEnd is how far the range start lies after its end. It is
// positive for a range that crosses the origin.
func startPastEnd(bounds rangeBounds) int {
	return boundsL.Start.Get(bounds) - boundsL.End.Get(bounds)
}

func normalFeatureRange(bounds rangeBounds) []genbank.Location {
	return F.Pipe1(normalLeafEdits, A.Map(F.Bind2nd(rangeLeaf, bounds)))
}

// wrappedFeatureRange splits an origin-crossing range at base 1 in source order.
func wrappedFeatureRange(bounds rangeBounds) []genbank.Location {
	return F.Pipe1(wrappedLeafEdits, A.Map(F.Bind2nd(rangeLeaf, bounds)))
}

func rangeLeaf(edits []rangeEdit, bounds rangeBounds) genbank.Location {
	return F.Pipe2(
		edits,
		A.Ap[En.Endomorphism[genbank.Location]]([]rangeBounds{bounds}),
		En.Reduce[genbank.Location],
	)
}

func rangeStart(bounds rangeBounds) int { return boundsL.Start.Get(bounds) - 1 }

func hasFeatureLeaves(s featureState) bool {
	return F.Pipe2(s, featureStateL.Leaves.Get, predarrays.IsNonEmpty[genbank.Location]())
}

func featureColor(s featureState) E.Either[error, O.Option[seqmodel.Color]] {
	return F.Pipe7(
		s,
		featureStateL.XML.Get,
		featureL.Segments.Get,
		A.Filter(standardSegment),
		A.Head[segmentXML],
		O.Map(segmentL.Color.Get),
		O.GetOrElse(F.Constant("")),
		parseColor,
	)
}

func featureLocation(s featureState) genbank.Location {
	return F.Pipe2(
		PR.MakePair(s, s),
		PR.BiMap(
			F.Flow2(
				featureStateL.Leaves.Get,
				featureLocationChoice(predarrays.LenEq[genbank.Location](1)),
			),
			F.Flow2(featureStateL.Complement.Get, PL.LocationComplement.Set),
		),
		PR.Merge(F.Identity[En.Endomorphism[genbank.Location]]),
	)
}

func firstLocation(locations []genbank.Location) genbank.Location { return locations[0] }

func joinLocation(locations []genbank.Location) genbank.Location {
	return F.Pipe2(
		genbank.Location{},
		PL.LocationJoin.Set(true),
		PL.LocationSubLocations.Set(locations),
	)
}

func buildFeature(s featureState) genbank.Feature {
	return F.Pipe2(
		featureEdits,
		A.Ap[En.Endomorphism[genbank.Feature]]([]featureState{s}),
		En.Reduce[genbank.Feature],
	)
}

func convertPrimer(record genbank.Genbank, primer primerXML) E.Either[error, []genbank.Feature] {
	return F.Pipe3(
		PR.MakePair(primer, primer),
		PR.BiMap(
			primerL.Name.Get,
			F.Flow3(
				primerL.Sites.Get,
				A.Filter(fullSite),
				E.TraverseArray(F.Bind1st(convertSite, record)),
			),
		),
		PR.Sequence(E.Map[error, []genbank.Location, PR.Pair[string, []genbank.Location]]),
		E.Map[error](
			F.Flow2(
				PR.MapHead[[]genbank.Location](F.Curry2(primerFeature)),
				PR.Merge(F.Flip(A.Map[genbank.Location, genbank.Feature])),
			),
		),
	)
}

func convertSite(record genbank.Genbank, site siteXML) E.Either[error, genbank.Location] {
	return F.Pipe7(
		siteState{},
		siteStateL.XML.Set(site),
		E.Of[error, siteState],
		E.Chain(E.FromPredicate(validSiteStrand, F.Constant1[siteState](errInvalidDirection))),
		E.Bind(siteStateL.Bounds.Set, F.Flow3(siteStateL.XML.Get, siteL.Location.Get, parseRange)),
		E.Chain(
			E.FromPredicate(
				F.Bind1st(validSiteRange, record),
				F.Constant1[siteState](seqmodel.ErrSegmentOutOfBounds),
			),
		),
		E.Let[error](siteStateL.Complement.Set, F.Flow2(siteStrand.Get, complementStrand)),
		E.Map[error](siteLocation),
	)
}

func validSiteStrand(s siteState) bool {
	return F.Pipe3(s, siteStateL.XML.Get, siteL.Strand.Get, strandPattern.MatchString)
}

func validSiteRange(record genbank.Genbank, s siteState) bool {
	sequence := PL.GenbankSequence.Get(record)
	bounds := siteStateL.Bounds.Get(s)
	start := boundsL.Start.Get(bounds)
	end := boundsL.End.Get(bounds)
	return start >= 0 && start <= end && end < len(sequence)
}

func siteLocation(s siteState) genbank.Location {
	return F.Pipe2(
		siteEdits,
		A.Ap[En.Endomorphism[genbank.Location]]([]siteState{s}),
		En.Reduce[genbank.Location],
	)
}

func primerFeature(name string, location genbank.Location) genbank.Feature {
	return F.Pipe2(
		primerEdits,
		A.Ap[En.Endomorphism[genbank.Feature]](
			[]PR.Pair[string, genbank.Location]{PR.MakePair(name, location)},
		),
		En.Reduce[genbank.Feature],
	)
}

// rangeEdit sets one field of a half-open leaf from one-based inclusive bounds.
type rangeEdit = func(rangeBounds) En.Endomorphism[genbank.Location]

var (
	// normalLeafEdits produce one leaf [start-1, end).
	normalLeafEdits = [][]rangeEdit{{
		F.Flow2(rangeStart, PL.LocationStart.Set),
		F.Flow2(boundsL.End.Get, PL.LocationEnd.Set),
	}}
	// wrappedLeafEdits produce [start-1, length) and then [0, end).
	wrappedLeafEdits = [][]rangeEdit{
		{
			F.Flow2(rangeStart, PL.LocationStart.Set),
			F.Flow2(boundsL.Length.Get, PL.LocationEnd.Set),
		},
		{F.Flow2(boundsL.End.Get, PL.LocationEnd.Set)},
	}
	featureEdits = []func(featureState) En.Endomorphism[genbank.Feature]{
		F.Flow2(featureType.Get, PL.FeatureType.Set),
		F.Flow2(featureStateL.Location.Get, PL.FeatureLocation.Set),

		F.Flow3(
			featureName.Get,
			F.Bind1st(R.Singleton[string, string], "label"),
			PL.FeatureAttributes.Set,
		),
		// Reduce applies edits left-to-right: establish label, then add optional color.
		F.Flow2(featureStateL.Color.Get, O.Fold(
			F.Constant(F.Identity[genbank.Feature]),
			F.Flow3(
				seqmodel.Color.Hex,
				F.Bind1st(R.UpsertAt[string, string], "ApEinfo_fwdcolor"),
				PL.FeatureAttributes.Modify,
			),
		)),
	}
	siteStart = F.Pipe1(siteStateL.Bounds, L.Compose[siteState](boundsL.Start))
	siteEnd   = F.Pipe1(siteStateL.Bounds, L.Compose[siteState](boundsL.End))
	siteEdits = []func(siteState) En.Endomorphism[genbank.Location]{
		F.Flow2(siteStart.Get, PL.LocationStart.Set),
		F.Flow3(siteEnd.Get, N.Add(1), PL.LocationEnd.Set),
		F.Flow2(siteStateL.Complement.Get, PL.LocationComplement.Set),
	}
	primerEdits = []func(PR.Pair[string, genbank.Location]) En.Endomorphism[genbank.Feature]{
		F.Constant1[PR.Pair[string, genbank.Location]](PL.FeatureType.Set("primer_bind")),
		F.Flow2(PR.Tail[string, genbank.Location], PL.FeatureLocation.Set),
		F.Flow3(
			PR.Head[string, genbank.Location],
			F.Bind1st(R.Singleton[string, string], "label"),
			PL.FeatureAttributes.Set,
		),
	}
)
