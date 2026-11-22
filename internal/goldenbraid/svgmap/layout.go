package svgmap

import (
	"math"

	A "github.com/IBM/fp-go/v2/array"
	En "github.com/IBM/fp-go/v2/endomorphism"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	ORD "github.com/IBM/fp-go/v2/ord"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/bebop/poly/io/genbank"
	predord "github.com/dictyBase/fp-go-loom/predicate/ord"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
)

type (
	indexedFeature = PR.Pair[int, genbank.Feature]
	midpointWalk   = PR.Pair[float64, float64]
	upperLabels    = PR.Pair[float64, []featureLayout]
)

var (
	spanOrder     = F.Pipe2(predord.IntOrd, ORD.ContraMap(featureSpan), ORD.Reverse[featureLayout])
	sourceOrder   = F.Pipe1(predord.IntOrd, ORD.ContraMap(featureL.ID.Get))
	yOrder        = F.Pipe1(predord.Float64Ord, ORD.ContraMap(featureL.Y.Get))
	invalidRadius = F.Pipe2(
		featureHalfWidthPx,
		ORD.Leq(predord.Float64Ord),
		P.ContraMap(featureL.Radius.Get),
	)
	sideChoice     = P.Fold(F.Constant1[float64](rightSide), F.Constant1[float64](leftSide))
	stateSize      = F.Pipe1(stateL.Config, L.Compose[renderState](configL.Size))
	hasLabel       = F.Flow2(labelText, S.IsNonEmpty)
	stateMaxLabels = F.Pipe1(stateL.Config, L.Compose[renderState](configL.MaxLabels))
	labelSides     = []string{leftSide, rightSide}
	// An occupied ring equal to the candidate moves the candidate one ring inward.
	ringBump   = P.Fold(F.Identity[int], N.Add(1))
	intEquals  = EQ.Equals(predord.IntEq)
	indexEdits = []func(indexedFeature) En.Endomorphism[featureLayout]{
		F.Flow2(PR.Head[int, genbank.Feature], featureL.ID.Set),
		F.Flow2(PR.Tail[int, genbank.Feature], featureL.Feature.Set),
		F.Flow4(
			PR.Tail[int, genbank.Feature],
			PL.FeatureLocation.Get,
			locationLeaves,
			featureL.Leaves.Set,
		),
	}
	// Each position field derives from the state and the packed feature alone.
	positionEdits = []func(pathContext) En.Endomorphism[featureLayout]{
		F.Flow2(PR.Paired(featureRadius), featureL.Radius.Set),
		F.Flow2(PR.Paired(featureAngle), featureL.Angle.Set),
		F.Flow2(PR.Paired(featureLabelY), featureL.Y.Set),
		F.Flow2(PR.Paired(featureSide), featureL.Side.Set),
		F.Flow2(PR.Paired(featureLabelX), featureL.X.Set),
	}
	midpointFold = F.Bind1st(A.Reduce[genbank.Location, midpointWalk], midpointStep)
	upperFold    = F.Bind1st(A.Reduce[featureLayout, upperLabels], pushUpper)
	nextLowerY   = F.Flow4(
		A.Last[featureLayout],
		O.Map(featureL.Y.Get),
		O.GetOrElse(F.Constant(0.0)),
		N.Add[float64](MinLabelGapPx),
	)
)

func locationLeaves(location genbank.Location) []genbank.Location {
	return F.Pipe3(
		location,
		O.FromPredicate(PL.LocationJoin.Get),
		O.Map(F.Flow2(PL.LocationSubLocations.Get, A.Chain(locationLeaves))),
		O.GetOrElse(F.Constant([]genbank.Location{location})),
	)
}

func indexFeature(index int, feature genbank.Feature) featureLayout {
	return F.Pipe2(
		indexEdits,
		A.Flap[En.Endomorphism[featureLayout]](PR.MakePair(index, feature)),
		En.Reduce[featureLayout],
	)
}

func visibleFeature(s renderState, item featureLayout) bool {
	cfg := stateL.Config.Get(s)
	feature := featureL.Feature.Get(item)
	kind := PL.FeatureType.Get(feature)
	return configL.ShowPrimers.Get(cfg) || kind != "primer_bind"
}

func layoutFeatures(s renderState) []featureLayout {
	return F.Pipe7(
		s,
		stateL.Record.Get,
		PL.GenbankFeatures.Get,
		A.MapWithIndex(indexFeature),
		A.Filter(F.Bind1st(visibleFeature, s)),
		A.Sort(spanOrder),
		A.Reduce(placeFeature, A.Empty[featureLayout]()),
		A.Map(F.Bind1st(positionFeature, s)),
	)
}

func featureSpan(item featureLayout) int {
	return F.Pipe2(item, featureL.Leaves.Get, A.Reduce(addLeafLength, 0))
}

func addLeafLength(total int, leaf genbank.Location) int {
	return total + PL.LocationEnd.Get(leaf) - PL.LocationStart.Get(leaf)
}

func placeFeature(placed []featureLayout, item featureLayout) []featureLayout {
	return F.Pipe2(item, F.Bind1st(findRing, placed), F.Bind1st(A.Append[featureLayout], placed))
}

// findRing picks the first ring, from the current one, that no overlapping feature uses.
// Ascending occupied rings move the candidate past each ring that it meets.
func findRing(placed []featureLayout, item featureLayout) featureLayout {
	return F.Pipe6(
		placed,
		A.Filter(F.Bind1st(sharesBases, item)),
		A.Map(featureL.Ring.Get),
		A.Sort(predord.IntOrd),
		A.Reduce(nextFreeRing, featureL.Ring.Get(item)),
		featureL.Ring.Set,
		En.Read(item),
	)
}

func nextFreeRing(candidate, occupied int) int {
	return F.Pipe3(occupied, intEquals, ringBump, En.Read(candidate))
}

func sharesBases(item, placed featureLayout) bool {
	return F.Pipe2(
		PR.MakePair(item, placed),
		PR.BiMap(featureL.Leaves.Get, featureL.Leaves.Get),
		PR.Paired(leavesOverlap),
	)
}

func leavesOverlap(first, second []genbank.Location) bool {
	return F.Pipe1(first, A.Any(F.Bind1st(overlapsAny, second)))
}

func overlapsAny(leaves []genbank.Location, leaf genbank.Location) bool {
	return F.Pipe1(leaves, A.Any(F.Bind1st(intervalsOverlap, leaf)))
}

func intervalsOverlap(first, second genbank.Location) bool {
	return PL.LocationStart.Get(first) < PL.LocationEnd.Get(second) &&
		PL.LocationStart.Get(second) < PL.LocationEnd.Get(first)
}

func positionFeature(s renderState, item featureLayout) featureLayout {
	return F.Pipe3(
		positionEdits,
		A.Flap[En.Endomorphism[featureLayout]](PR.MakePair(s, item)),
		M.ConcatAll(En.Monoid[featureLayout]()),
		En.Read(item),
	)
}

func sizePx(s renderState) float64 { return float64(stateSize.Get(s)) }

func lengthPx(s renderState) float64 { return float64(stateL.Length.Get(s)) }

// angleAt maps a base position to a clockwise angle that starts at 12 o'clock.
func angleAt(position, length float64) float64 {
	return fullTurn*math.Pi*position/length - math.Pi/half
}

func featureRadius(s renderState, item featureLayout) float64 {
	size := sizePx(s)
	return size*radiusFraction - float64(featureL.Ring.Get(item))*ringStepPx
}

func featureAngle(s renderState, item featureLayout) float64 {
	midpoint := featureMidpoint(item)
	length := lengthPx(s)
	return angleAt(midpoint, length)
}

func featureLabelY(s renderState, item featureLayout) float64 {
	size := sizePx(s)
	radius := featureRadius(s, item)
	angle := featureAngle(s, item)
	return size/half + radius*math.Sin(angle)
}

func featureSide(s renderState, item featureLayout) string {
	angle := featureAngle(s, item)
	return labelSide(angle)
}

func featureLabelX(s renderState, item featureLayout) float64 {
	size := sizePx(s)
	side := featureSide(s, item)
	return labelX(size, side)
}

// featureMidpoint walks occupied bases in source order, not the intervening gaps.
func featureMidpoint(item featureLayout) float64 {
	return F.Pipe4(
		item,
		PR.Of[featureLayout],
		PR.BiMap(featureL.Leaves.Get, midpointSeed),
		PR.Merge(midpointFold),
		PR.Tail[float64, float64],
	)
}

func midpointSeed(item featureLayout) midpointWalk {
	span := float64(featureSpan(item))
	return PR.MakePair(span/half, 0.0)
}

func midpointStep(acc midpointWalk, leaf genbank.Location) midpointWalk {
	return F.Pipe3(
		acc,
		O.FromPredicate(F.Flow2(PR.Head[float64, float64], N.MoreThan(0.0))),
		O.Map(F.Bind1st(advanceMidpoint, leaf)),
		O.GetOrElse(F.Constant(acc)),
	)
}

func advanceMidpoint(leaf genbank.Location, acc midpointWalk) midpointWalk {
	remaining := PR.Head(acc)
	start := float64(PL.LocationStart.Get(leaf))
	end := float64(PL.LocationEnd.Get(leaf))
	consumed := math.Min(remaining, end-start)
	return PR.MakePair(remaining-consumed, start+consumed)
}

func labelSide(angle float64) string {
	return F.Pipe2(angle, math.Cos, sideChoice(N.LessThan(0.0)))
}

func labelX(size float64, side string) float64 {
	sign := direction(side == leftSide)
	rightX := size/half + size*radiusFraction + labelLeaderGapPx
	return LabelMarginPx + (sign+1)/half*(rightX-LabelMarginPx)
}

func labelText(item featureLayout) string {
	feature := featureL.Feature.Get(item)
	attributes := PL.FeatureAttributes.Get(feature)
	return attributes["label"]
}

func layoutLabels(s renderState) []featureLayout {
	return F.Pipe5(
		s,
		PR.Of[renderState],
		PR.BiMap(
			F.Flow3(stateL.Features.Get, A.Filter(hasLabel), A.Sort(sourceOrder)),
			stateMaxLabels.Get,
		),
		PR.Merge(F.Bind1st(A.Slice[featureLayout], 0)),
		F.Bind1st(relaxLabels, s),
		A.Sort(sourceOrder),
	)
}

func relaxLabels(s renderState, items []featureLayout) []featureLayout {
	return F.Pipe2(
		labelSides,
		A.Map(F.Bind1st(sideLabels, items)),
		A.Chain(F.Bind1st(relaxSide, s)),
	)
}

func sideLabels(items []featureLayout, side string) []featureLayout {
	return F.Pipe1(items, A.Filter(F.Flow2(featureL.Side.Get, predord.StrEq(side))))
}

// relaxSide keeps the per-side capacity, then spaces labels downward and upward.
func relaxSide(s renderState, items []featureLayout) []featureLayout {
	return F.Pipe7(
		PR.MakePair(items, s),
		PR.MapTail[[]featureLayout](labelCapacity),
		PR.Merge(F.Bind1st(A.Slice[featureLayout], 0)),
		A.Sort(yOrder),
		A.Reduce(pushLower, A.Empty[featureLayout]()),
		A.Reverse[featureLayout],
		F.Bind1st(spaceUpward, s),
		A.Reverse[featureLayout],
	)
}

func labelCapacity(s renderState) int {
	size := stateSize.Get(s)
	return (size-LabelMarginPx-LabelMarginPx)/MinLabelGapPx + 1
}

func spaceUpward(s renderState, items []featureLayout) []featureLayout {
	return F.Pipe3(
		PR.MakePair(items, s),
		PR.MapTail[[]featureLayout](upperSeed),
		PR.Merge(upperFold),
		PR.Tail[float64, []featureLayout],
	)
}

func upperSeed(s renderState) upperLabels {
	upper := float64(stateSize.Get(s) - LabelMarginPx + MinLabelGapPx)
	return PR.MakePair(upper, A.Empty[featureLayout]())
}

func pushLower(placed []featureLayout, item featureLayout) []featureLayout {
	return F.Pipe4(
		item,
		F.Bind1st(lowerLabelY, placed),
		featureL.Y.Set,
		En.Read(item),
		F.Bind1st(A.Append[featureLayout], placed),
	)
}

// lowerLabelY keeps the ideal position unless the previous label is too close.
func lowerLabelY(placed []featureLayout, item featureLayout) float64 {
	return F.Pipe2(
		PR.MakePair(placed, item),
		PR.BiMap(nextLowerY, featureL.Y.Get),
		PR.Merge(F.Curry2(math.Max)),
	)
}

func pushUpper(acc upperLabels, item featureLayout) upperLabels {
	return F.Pipe4(
		item,
		F.Bind1st(upperLabelY, acc),
		featureL.Y.Set,
		En.Read(item),
		F.Bind1st(stackUpper, acc),
	)
}

func upperLabelY(acc upperLabels, item featureLayout) float64 {
	ideal := featureL.Y.Get(item)
	ceiling := PR.Head(acc) - MinLabelGapPx
	return math.Min(ideal, ceiling)
}

func stackUpper(acc upperLabels, placed featureLayout) upperLabels {
	y := featureL.Y.Get(placed)
	items := A.Append(PR.Tail(acc), placed)
	return PR.MakePair(y, items)
}
