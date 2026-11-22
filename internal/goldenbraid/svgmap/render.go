package svgmap

import (
	"encoding/xml"
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/bebop/poly/io/genbank"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
)

// fp-go:Lens
type svgDocument struct {
	XMLName     xml.Name      `xml:"svg"`
	Namespace   string        `xml:"xmlns,attr"`
	Role        string        `xml:"role,attr"`
	ViewBox     string        `xml:"viewBox,attr"`
	Title       string        `xml:"title"`
	Description string        `xml:"desc"`
	Backbone    svgCircle     `xml:"circle"`
	Paths       []svgPath     `xml:"path"`
	Leaders     []svgPolyline `xml:"polyline"`
	Text        []svgText     `xml:"text"`
}

// fp-go:Lens
type svgCircle struct {
	CX     float64 `xml:"cx,attr"`
	CY     float64 `xml:"cy,attr"`
	Radius float64 `xml:"r,attr"`
	Fill   string  `xml:"fill,attr"`
	Stroke string  `xml:"stroke,attr"`
}

// fp-go:Lens
type svgPath struct {
	ID   string `xml:"data-feature-id,attr"`
	Ring int    `xml:"data-ring,attr"`
	D    string `xml:"d,attr"`
	Fill string `xml:"fill,attr"`
}

// fp-go:Lens
type svgPolyline struct {
	Points string `xml:"points,attr"`
	Fill   string `xml:"fill,attr"`
	Stroke string `xml:"stroke,attr"`
}

// fp-go:Lens
type svgText struct {
	ID       string  `xml:"data-feature-id,attr,omitempty"`
	Side     string  `xml:"data-side,attr,omitempty"`
	X        float64 `xml:"x,attr"`
	Y        float64 `xml:"y,attr"`
	Length   float64 `xml:"textLength,attr,omitempty"`
	Adjust   string  `xml:"lengthAdjust,attr,omitempty"`
	FontSize float64 `xml:"font-size,attr"`
	Anchor   string  `xml:"text-anchor,attr,omitempty"`
	Value    string  `xml:",chardata"`
}

// fp-go:Lens
type arcShape struct {
	Center  float64
	Radius  float64
	Start   float64
	End     float64
	Reverse bool
}

type (
	pathContext = PR.Pair[renderState, featureLayout]
	leafContext = PR.Pair[pathContext, genbank.Location]
)

var (
	float64Gt   = ORD.Gt(predord.Float64Ord)
	intGeq      = ORD.Geq(predord.IntOrd)
	hasName     = F.Flow3(stateL.Record.Get, fullName.Get, S.IsNonEmpty)
	hasSequence = F.Flow2(sequenceLength, P.IsNonZero[int]())
	// validConfig accepts a size wide enough for the label margins and
	// a label column, and a label limit that is not negative.
	validConfig = F.Pipe2(
		F.Flow2(marginRoom, float64Gt(0)),
		P.And(F.Flow2(labelColumnWidth, float64Gt(0))),
		P.And(F.Flow2(stateMaxLabels.Get, intGeq(0))),
	)
	documentL = MakesvgDocumentLenses()
	circleL   = MakesvgCircleLenses()
	pathL     = MakesvgPathLenses()
	leaderL   = MakesvgPolylineLenses()
	textL     = MakesvgTextLenses()
	arcL      = MakearcShapeLenses()
	fullName  = F.Pipe1(PL.GenbankMeta, L.Compose[genbank.Genbank](PL.MetaName))
	stateName = F.Pipe1(stateL.Record, L.Compose[renderState](fullName))
	circular  = F.Pipe2(
		PL.GenbankMeta,
		L.Compose[genbank.Genbank](PL.MetaLocus),
		L.Compose[genbank.Genbank](PL.LocusCircular),
	)
	largeArcChoice  = P.Fold(F.Constant1[float64](0), F.Constant1[float64](1))
	directionChoice = P.Fold(F.Constant1[bool](1.0), F.Constant1[bool](-1.0))
	renderSeed      = F.Flow4(
		PR.Of[RenderInput],
		PR.BiMap(
			F.Flow2(inputL.Record.Get, stateL.Record.Set),
			F.Flow5(
				inputL.Options.Get,
				A.Reverse[Option],
				M.ConcatAll(En.Monoid[Config]()),
				En.Read(defaultConfig),
				stateL.Config.Set,
			),
		),
		PR.Paired(En.Monoid[renderState]().Concat),
		En.Build[renderState],
	)
	titleText = F.Flow3(
		PR.Of[renderState],
		PR.BiMap(stateName.Get, stateL.Length.Get),
		PR.Paired(formatTitle),
	)
	descriptionText = F.Flow3(
		PR.Of[renderState],
		PR.BiMap(
			F.Flow3(stateL.Features.Get, A.Filter(hasLabel), A.Size[featureLayout]),
			F.Flow2(stateL.Labels.Get, A.Size[featureLayout]),
		),
		PR.Paired(formatDescription),
	)
	documentEdits = []func(renderState) En.Endomorphism[svgDocument]{
		F.Constant1[renderState](documentL.Namespace.Set("http://www.w3.org/2000/svg")),
		F.Constant1[renderState](documentL.Role.Set("img")),
		F.Flow3(stateSize.Get, S.Format[int]("0 0 %[1]d %[1]d"), documentL.ViewBox.Set),
		F.Flow2(titleText, documentL.Title.Set),
		F.Flow2(descriptionText, documentL.Description.Set),
		F.Flow2(backbone, documentL.Backbone.Set),
		F.Flow4(
			PR.Of[renderState],
			PR.MapTail[renderState](F.Flow2(stateL.Features.Get, A.Sort(sourceOrder))),
			PR.Paired(featureListPaths),
			documentL.Paths.Set,
		),
		F.Flow4(
			PR.Of[renderState],
			PR.MapTail[renderState](stateL.Labels.Get),
			PR.Paired(labelLeaders),
			documentL.Leaders.Set,
		),
		F.Flow4(
			PR.Of[renderState],
			PR.BiMap(
				F.Flow3(
					PR.Of[renderState],
					PR.MapTail[renderState](stateL.Labels.Get),
					PR.Paired(labelTexts),
				),
				centerLabel,
			),
			PR.Merge(A.Push[svgText]),
			documentL.Text.Set,
		),
	}
	backboneEdits = []func(renderState) En.Endomorphism[svgCircle]{
		F.Flow2(centerPx, circleL.CX.Set),
		F.Flow2(centerPx, circleL.CY.Set),
		F.Flow2(backboneRadius, circleL.Radius.Set),
		F.Constant1[renderState](circleL.Fill.Set("none")),
		F.Constant1[renderState](circleL.Stroke.Set(backboneStroke)),
	}
	leafState = F.Flow2(PR.Head[pathContext, genbank.Location], PR.Head[renderState, featureLayout])
	leafItem  = F.Flow2(PR.Head[pathContext, genbank.Location], PR.Tail[renderState, featureLayout])
	pathEdits = []func(leafContext) En.Endomorphism[svgPath]{
		F.Flow4(leafItem, featureL.ID.Get, strconv.Itoa, pathL.ID.Set),
		F.Flow3(leafItem, featureL.Ring.Get, pathL.Ring.Set),
		F.Flow2(leafPath, pathL.D.Set),
		F.Flow3(leafItem, featureFill, pathL.Fill.Set),
	}
	arcEdits = []func(leafContext) En.Endomorphism[arcShape]{
		F.Flow3(leafState, centerPx, arcL.Center.Set),
		F.Flow3(leafItem, featureL.Radius.Get, arcL.Radius.Set),
		F.Flow2(F.Bind1st(leafAngle, PL.LocationStart.Get), arcL.Start.Set),
		F.Flow2(F.Bind1st(leafAngle, PL.LocationEnd.Get), arcL.End.Set),
		F.Flow5(
			leafItem,
			featureL.Feature.Get,
			PL.FeatureLocation.Get,
			PL.LocationComplement.Get,
			arcL.Reverse.Set,
		),
	}
	leaderEdits = []func(pathContext) En.Endomorphism[svgPolyline]{
		F.Flow2(PR.Paired(leaderPoints), leaderL.Points.Set),
		F.Constant1[pathContext](leaderL.Fill.Set("none")),
		F.Constant1[pathContext](leaderL.Stroke.Set(labelStroke)),
	}
	labelEdits = []func(pathContext) En.Endomorphism[svgText]{
		F.Flow4(PR.Tail[renderState, featureLayout], featureL.ID.Get, strconv.Itoa, textL.ID.Set),
		F.Flow3(PR.Tail[renderState, featureLayout], featureL.Side.Get, textL.Side.Set),
		F.Flow3(PR.Tail[renderState, featureLayout], featureL.X.Get, textL.X.Set),
		F.Flow3(PR.Tail[renderState, featureLayout], featureL.Y.Get, textL.Y.Set),
		F.Flow2(PR.Paired(labelWidth), textL.Length.Set),
		F.Constant1[pathContext](textL.Adjust.Set("spacingAndGlyphs")),
		F.Constant1[pathContext](textL.FontSize.Set(labelFontSizePx)),
		F.Constant1[pathContext](textL.Anchor.Set("start")),
		F.Flow3(PR.Tail[renderState, featureLayout], labelText, textL.Value.Set),
	}
	centerEdits = []func(renderState) En.Endomorphism[svgText]{
		F.Flow2(centerPx, textL.X.Set),
		F.Flow2(centerPx, textL.Y.Set),
		F.Constant1[renderState](textL.FontSize.Set(labelFontSizePx)),
		F.Constant1[renderState](textL.Anchor.Set("middle")),
		F.Flow3(stateL.Length.Get, S.Format[int]("%d bp"), textL.Value.Set),
	}
	featureFill = F.Flow5(
		featureL.Feature.Get,
		PL.FeatureAttributes.Get,
		R.Lookup[string]("ApEinfo_fwdcolor"),
		O.Filter(S.IsNonEmpty),
		O.GetOrElse(F.Constant(defaultFill)),
	)
)

// Render encodes a nonempty circular record without file or bucket effects.
func Render(input RenderInput) E.Either[error, []byte] {
	return F.Pipe13(
		input,
		renderSeed,
		E.Of[error, renderState],
		E.Chain(E.FromPredicate(hasName, F.Constant1[renderState](seqmodel.ErrEmptyName))),
		E.Chain(E.FromPredicate(hasSequence, F.Constant1[renderState](seqmodel.ErrEmptySequence))),
		E.Chain(E.FromPredicate(isCircular, F.Constant1[renderState](ErrNotCircular))),
		E.Chain(E.FromPredicate(validConfig, F.Constant1[renderState](ErrInvalidConfig))),
		E.Let[error](stateL.Length.Set, sequenceLength),
		E.Chain(
			E.FromPredicate(
				validLocations,
				F.Constant1[renderState](seqmodel.ErrSegmentOutOfBounds),
			),
		),
		E.Let[error](stateL.Features.Set, layoutFeatures),
		E.Chain(E.FromPredicate(validRings, F.Constant1[renderState](ErrInvalidConfig))),
		E.Let[error](stateL.Labels.Set, layoutLabels),
		E.Map[error](buildDocument),
		E.Chain(marshalDocument),
	)
}

func isCircular(s renderState) bool {
	return F.Pipe2(s, stateL.Record.Get, circular.Get)
}

func sequenceLength(s renderState) int {
	record := stateL.Record.Get(s)
	sequence := PL.GenbankSequence.Get(record)
	return len(sequence)
}

// marginRoom is the map size left after both label margins.
func marginRoom(s renderState) float64 { return sizePx(s) - LabelMarginPx*half }

// labelColumnWidth is the width left for label text on one side of
// the map.
func labelColumnWidth(s renderState) float64 {
	size := sizePx(s)
	return size/half - size*radiusFraction - labelLeaderGapPx - LabelMarginPx
}

func validLocations(s renderState) bool {
	return F.Pipe5(
		s,
		stateL.Record.Get,
		PL.GenbankFeatures.Get,
		A.MapWithIndex(indexFeature),
		A.Filter(F.Bind1st(visibleFeature, s)),
		P.Not(A.Any(F.Bind1st(invalidFeature, s))),
	)
}

func invalidFeature(s renderState, item featureLayout) bool {
	leaves := featureL.Leaves.Get(item)
	return len(leaves) == 0 || F.Pipe1(leaves, A.Any(F.Bind1st(invalidLeaf, s)))
}

func invalidLeaf(s renderState, leaf genbank.Location) bool {
	start := PL.LocationStart.Get(leaf)
	end := PL.LocationEnd.Get(leaf)
	return start < 0 || start >= end || end > stateL.Length.Get(s)
}

func validRings(s renderState) bool {
	return F.Pipe2(
		s,
		stateL.Features.Get,
		P.Not(A.Any(invalidRadius)),
	)
}

func buildDocument(s renderState) svgDocument {
	return F.Pipe2(documentEdits, A.Flap[En.Endomorphism[svgDocument]](s), En.Reduce[svgDocument])
}

func formatTitle(name string, length int) string {
	return fmt.Sprintf("%s (%d bp)", name, length)
}

func formatDescription(labeled, visible int) string {
	return fmt.Sprintf(
		"Circular sequence map. %d feature labels omitted by label limits or viewport capacity.",
		labeled-visible,
	)
}

func featureListPaths(s renderState, features []featureLayout) []svgPath {
	return F.Pipe1(
		features,
		A.Chain(F.Flow4(
			F.Bind1st(PR.MakePair[renderState, featureLayout], s),
			PR.Of[pathContext],
			PR.MapTail[pathContext](
				F.Flow2(PR.Tail[renderState, featureLayout], featureL.Leaves.Get),
			),
			PR.Paired(leafPaths),
		)),
	)
}

func leafPaths(context pathContext, leaves []genbank.Location) []svgPath {
	return F.Pipe1(leaves, A.Map(F.Bind1st(makePath, context)))
}

func labelLeaders(s renderState, labels []featureLayout) []svgPolyline {
	return F.Pipe1(labels, A.Map(F.Bind1st(labelLeader, s)))
}

func labelTexts(s renderState, labels []featureLayout) []svgText {
	return F.Pipe1(labels, A.Map(F.Bind1st(featureLabel, s)))
}

func centerPx(s renderState) float64 { return sizePx(s) / half }

func backboneRadius(s renderState) float64 { return sizePx(s) * radiusFraction }

func backbone(s renderState) svgCircle {
	return F.Pipe2(backboneEdits, A.Flap[En.Endomorphism[svgCircle]](s), En.Reduce[svgCircle])
}

func makePath(context pathContext, leaf genbank.Location) svgPath {
	return F.Pipe2(
		pathEdits,
		A.Flap[En.Endomorphism[svgPath]](PR.MakePair(context, leaf)),
		En.Reduce[svgPath],
	)
}

func leafPath(lc leafContext) string {
	return F.Pipe3(arcEdits, A.Flap[En.Endomorphism[arcShape]](lc), En.Reduce[arcShape], arcPath)
}

func leafAngle(position func(genbank.Location) int, lc leafContext) float64 {
	base := float64(position(PR.Tail(lc)))
	s := leafState(lc)
	length := lengthPx(s)
	return angleAt(base, length)
}

func arcPath(shape arcShape) string {
	center := arcL.Center.Get(shape)
	radius := arcL.Radius.Get(shape)
	start := arcL.Start.Get(shape)
	end := arcL.End.Get(shape)
	reverse := arcL.Reverse.Get(shape)
	sign := direction(reverse)
	sweep := (1 + int(sign)) / int(half)
	otherSweep := 1 - sweep
	span := end - start
	head := math.Min(ArrowAngleDeg*math.Pi/degreesPerHalfTurn, span/half)
	origin := (start + end - sign*span) / half
	tip := (start + end + sign*span) / half
	neck := tip - sign*head
	outerRadius := radius + featureHalfWidthPx
	innerRadius := radius - featureHalfWidthPx
	outerStart := polar(center, outerRadius, origin)
	outerEnd := polar(center, outerRadius, neck)
	tipPoint := polar(center, radius, tip)
	innerEnd := polar(center, innerRadius, neck)
	innerStart := polar(center, innerRadius, origin)
	large := largeArcFlag(span - head)
	return fmt.Sprintf(
		"M %s A %.3f %.3f 0 %d %d %s L %s L %s A %.3f %.3f 0 %d %d %s Z",
		outerStart,
		outerRadius,
		outerRadius,
		large,
		sweep,
		outerEnd,
		tipPoint,
		innerEnd,
		innerRadius,
		innerRadius,
		large,
		otherSweep,
		innerStart,
	)
}

func direction(reverse bool) float64 {
	return F.Pipe1(reverse, directionChoice(F.Identity[bool]))
}

func largeArcFlag(angle float64) int {
	return F.Pipe1(angle, largeArcChoice(N.MoreThan(math.Pi)))
}

func polar(center, radius, angle float64) string {
	x := center + radius*math.Cos(angle)
	y := center + radius*math.Sin(angle)
	return fmt.Sprintf("%.3f,%.3f", x, y)
}

func labelLeader(s renderState, item featureLayout) svgPolyline {
	return F.Pipe2(
		leaderEdits,
		A.Flap[En.Endomorphism[svgPolyline]](PR.MakePair(s, item)),
		En.Reduce[svgPolyline],
	)
}

func leaderPoints(s renderState, item featureLayout) string {
	size := sizePx(s)
	origin := polar(size/half, featureL.Radius.Get(item), featureL.Angle.Get(item))
	x := featureL.X.Get(item)
	y := featureL.Y.Get(item)
	side := featureL.Side.Get(item)
	width := labelWidth(s, item)
	sign := direction(side == leftSide)
	endpoint := x + width*(1-sign)/half
	return fmt.Sprintf("%s %.3f,%.3f", origin, endpoint, y)
}

func featureLabel(s renderState, item featureLayout) svgText {
	return F.Pipe2(
		labelEdits,
		A.Flap[En.Endomorphism[svgText]](PR.MakePair(s, item)),
		En.Reduce[svgText],
	)
}

func labelWidth(s renderState, item featureLayout) float64 {
	text := labelText(item)
	estimated := float64(utf8.RuneCountInString(text)) * glyphWidthPx
	return math.Min(labelColumnWidth(s), estimated)
}

func centerLabel(s renderState) svgText {
	return F.Pipe2(centerEdits, A.Flap[En.Endomorphism[svgText]](s), En.Reduce[svgText])
}

func marshalDocument(doc svgDocument) E.Either[error, []byte] {
	return F.Pipe2(
		any(doc),
		E.Eitherize1(xml.Marshal),
		E.MapLeft[[]byte](ER.OnError("encode SVG map")),
	)
}
