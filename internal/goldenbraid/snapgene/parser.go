// Package snapgene decodes SnapGene DNA packets into poly records.
package snapgene

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	ER "github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	PR "github.com/IBM/fp-go/v2/pair"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/bebop/poly/io/genbank"
	MO "github.com/dictyBase/fp-go-loom/matchopt"
	predbytes "github.com/dictyBase/fp-go-loom/predicate/bytes"
	"github.com/dictyBase/fp-go-loom/strutils"

	PL "github.com/dictyBase/modware-import/internal/goldenbraid/polylens"
	"github.com/dictyBase/modware-import/internal/goldenbraid/seqmodel"
)

//go:generate go run github.com/IBM/fp-go/v2@v2.13.2 lens --dir . --filename gen_lens.go

// LocusNameMax limits the GenBank locus identifier in characters.
const LocusNameMax = 16

// ParseInput describes one source construct.
// fp-go:Lens
type ParseInput struct {
	SourceFile string
	Name       string
	Data       []byte
}

// Packet contains a checked TLV body. Body references the input bytes.
// fp-go:Lens
type Packet struct {
	ID   byte
	Body []byte
}

// ParseResult contains the shared record and recoverable metadata warnings.
// fp-go:Lens
type ParseResult struct {
	Record   genbank.Genbank
	Warnings []error
}

// fp-go:Lens
type parseState struct {
	Input    ParseInput
	Packets  []Packet
	Record   genbank.Genbank
	Features []genbank.Feature
	Primers  []genbank.Feature
	Notes    []notesXML
	Warnings []error
}

var (
	// ErrInvalidFeature identifies an annotation that could not be decoded.
	ErrInvalidFeature = errors.New("invalid SnapGene feature")
	// ErrInvalidPrimer identifies a primer that could not be decoded.
	ErrInvalidPrimer = errors.New("invalid SnapGene primer")
	// ErrInvalidNotes identifies a malformed optional Notes packet.
	ErrInvalidNotes = errors.New("invalid SnapGene notes")
	inputL          = MakeParseInputLenses()
	packetL         = MakePacketLenses()
	resultL         = MakeParseResultLenses()
	stateL          = MakeparseStateLenses()
	stateData       = F.Pipe1(stateL.Input, L.Compose[parseState](inputL.Data))
	stateName       = F.Pipe1(stateL.Input, L.Compose[parseState](inputL.Name))
	recordCircular  = F.Pipe2(
		PL.GenbankMeta,
		L.Compose[genbank.Genbank](PL.MetaLocus),
		L.Compose[genbank.Genbank](PL.LocusCircular),
	)
	metaLocusName   = F.Pipe1(PL.MetaLocus, L.Compose[genbank.Meta](PL.LocusName))
	metaLocusLength = F.Pipe1(PL.MetaLocus, L.Compose[genbank.Meta](PL.LocusSequenceLength))
	metaMolecule    = F.Pipe1(PL.MetaLocus, L.Compose[genbank.Meta](PL.LocusMoleculeType))
	stateMetaName   = F.Pipe2(
		stateL.Record,
		L.Compose[parseState](PL.GenbankMeta),
		L.Compose[parseState](PL.MetaName),
	)
	stateLocusName = F.Pipe2(
		stateL.Record,
		L.Compose[parseState](PL.GenbankMeta),
		L.Compose[parseState](metaLocusName),
	)
	stateLocusLength = F.Pipe2(
		stateL.Record,
		L.Compose[parseState](PL.GenbankMeta),
		L.Compose[parseState](metaLocusLength),
	)
	stateMolecule = F.Pipe2(
		stateL.Record,
		L.Compose[parseState](PL.GenbankMeta),
		L.Compose[parseState](metaMolecule),
	)
	stateDefinition = F.Pipe2(
		stateL.Record,
		L.Compose[parseState](PL.GenbankMeta),
		L.Compose[parseState](PL.MetaDefinition),
	)
	stateComment = F.Pipe2(
		stateL.Record,
		L.Compose[parseState](PL.GenbankMeta),
		L.Compose[parseState](PL.MetaOther),
	)
	stateFeatures = F.Pipe1(stateL.Record, L.Compose[parseState](PL.GenbankFeatures))
	stateSequence = F.Pipe1(stateL.Record, L.Compose[parseState](PL.GenbankSequence))

	notesType        = notesL.Type.Get
	notesUUID        = notesL.UUID.Get
	notesCreated     = notesL.Created.Get
	notesModified    = notesL.Modified.Get
	notesTransformed = notesL.Transformed.Get
)

// Parse validates framing and DNA, then decodes optional annotations and notes.
func Parse(in ParseInput) E.Either[error, ParseResult] {
	return F.Pipe14(
		parseState{},
		stateL.Input.Set(in),
		E.Of[error, parseState],
		E.Chain(
			E.FromPredicate(
				F.Flow2(stateName.Get, S.IsNonEmpty),
				F.Constant1[parseState](seqmodel.ErrEmptyName),
			),
		),
		E.Bind(stateL.Packets.Set, F.Flow2(stateData.Get, Packets)),
		E.Bind(
			stateL.Record.Set,
			F.Flow5(
				stateL.Packets.Get,
				A.Filter(F.Flow2(packetL.ID.Get, predordByteDNA)),
				A.Head[Packet],
				E.FromOption[Packet](F.Constant(ErrMissingDNAPacket)),
				E.Chain(decodeDNA),
			),
		),
		E.Map[error](decodeAnnotations),
		E.Let[error](stateMetaName.Set, stateName.Get),
		E.Let[error](stateLocusName.Set, F.Flow3(stateName.Get, S.ToRunes, clipLocus)),
		E.Let[error](stateLocusLength.Set, F.Flow3(stateSequence.Get, S.Size, strconv.Itoa)),
		E.Let[error](stateMolecule.Set, F.Constant1[parseState]("DNA")),
		E.Let[error](stateDefinition.Set, constructStateDefinition),
		E.Let[error](stateComment.Set, F.Flow2(stateL.Notes.Get, commentText)),
		E.Let[error](stateFeatures.Set, annotationFeatures),
		E.Map[error](buildResult),
	)
}

func decodeDNA(packet Packet) E.Either[error, genbank.Genbank] {
	return F.Pipe3(
		packet,
		packetL.Body.Get,
		E.FromPredicate(predbytes.HasPositiveLen, F.Constant1[[]byte](seqmodel.ErrEmptySequence)),
		E.Chain(validateBases),
	)
}

func validateBases(body []byte) E.Either[error, genbank.Genbank] {
	return F.Pipe3(
		PR.MakePair(body, body),
		PR.BiMap(isCircularBody, F.Flow2(sequenceBases, seqmodel.ValidateSequence)),
		PR.Sequence(E.Map[error, string, PR.Pair[bool, string]]),
		E.Map[error](PR.Paired(sequenceRecord)),
	)
}

func isCircularBody(body []byte) bool { return body[0]&circularFlag != 0 }

func sequenceBases(body []byte) string { return string(body[1:]) }

func sequenceRecord(circular bool, sequence string) genbank.Genbank {
	return F.Pipe2(
		genbank.Genbank{},
		recordCircular.Set(circular),
		PL.GenbankSequence.Set(sequence),
	)
}

func decodeAnnotations(s parseState) parseState {
	return F.Pipe2(s, stateL.Packets.Get, A.Reduce(decodePacket, s))
}

// decodePacket adds a features, primers, or notes packet to s and
// skips every other packet.
func decodePacket(s parseState, packet Packet) parseState {
	return F.Pipe1(packet, MO.Match(s, []O.Kleisli[Packet, parseState]{
		MO.Case(isFeaturesPacket, F.Bind1st(readFeatures, s)),
		MO.Case(isPrimersPacket, F.Bind1st(readPrimers, s)),
		MO.Case(isNotesPacket, F.Bind1st(readNotes, s)),
	}))
}

func readFeatures(s parseState, packet Packet) parseState {
	return F.Pipe5(
		packet,
		packetL.Body.Get,
		decodeXML[featuresXML],
		E.Map[error](featuresL.Items.Get),
		E.Map[error](A.Map(F.Flow3(
			F.Bind1st(PR.MakePair[parseState, featureXML], s),
			PR.MapHead[featureXML](stateL.Record.Get), PR.Paired(convertFeature)))),
		E.Fold(F.Bind1st(addFeaturePacketWarning, s), A.Reduce(collectFeature, s)),
	)
}

func readPrimers(s parseState, packet Packet) parseState {
	return F.Pipe5(
		packet,
		packetL.Body.Get,
		decodeXML[primersXML],
		E.Map[error](primersL.Items.Get),
		E.Map[error](A.Map(F.Flow3(
			F.Bind1st(PR.MakePair[parseState, primerXML], s),
			PR.MapHead[primerXML](stateL.Record.Get), PR.Paired(convertPrimer)))),
		E.Fold(F.Bind1st(addPrimerPacketWarning, s), A.Reduce(collectPrimer, s)),
	)
}

func readNotes(s parseState, packet Packet) parseState {
	return F.Pipe4(
		packet,
		packetL.Body.Get,
		decodeXML[notesXML],
		E.MapLeft[notesXML](ER.OnError("%w", ErrInvalidNotes)),
		E.Fold(F.Bind1st(addWarning, s), F.Bind1st(addNotes, s)),
	)
}

func decodeXML[T any](data []byte) E.Either[error, T] {
	return F.Pipe2(data, E.Eitherize1(unmarshalXML[T]), E.Map[error](F.Deref[T]))
}

// unmarshalXML is the raw decode call. It returns the pointer that Unmarshal fills.
func unmarshalXML[T any](data []byte) (*T, error) {
	value := new(T)
	return value, xml.Unmarshal(data, value)
}

func addFeaturePacketWarning(s parseState, err error) parseState {
	return F.Pipe2(err, ER.OnError("%w", ErrInvalidFeature), F.Bind1st(addWarning, s))
}

func addPrimerPacketWarning(s parseState, err error) parseState {
	return F.Pipe2(err, ER.OnError("%w", ErrInvalidPrimer), F.Bind1st(addWarning, s))
}

func addWarning(s parseState, err error) parseState {
	return F.Pipe1(s, stateL.Warnings.Modify(A.Push(err)))
}

func addNotes(s parseState, notes notesXML) parseState {
	return F.Pipe1(s, stateL.Notes.Modify(A.Push(notes)))
}

func collectFeature(s parseState, item E.Either[error, genbank.Feature]) parseState {
	return F.Pipe2(
		item,
		E.MapLeft[genbank.Feature](ER.OnError("%w", ErrInvalidFeature)),
		E.Fold(F.Bind1st(addWarning, s), F.Bind1st(addFeature, s)),
	)
}

func addFeature(s parseState, feature genbank.Feature) parseState {
	return F.Pipe1(s, stateL.Features.Modify(A.Push(feature)))
}

func collectPrimer(s parseState, item E.Either[error, []genbank.Feature]) parseState {
	return F.Pipe2(
		item,
		E.MapLeft[[]genbank.Feature](ER.OnError("%w", ErrInvalidPrimer)),
		E.Fold(F.Bind1st(addWarning, s), F.Bind1st(addPrimerFeatures, s)),
	)
}

func addPrimerFeatures(s parseState, features []genbank.Feature) parseState {
	return F.Pipe1(s, stateL.Primers.Modify(A.Concat(features)))
}

func annotationFeatures(s parseState) []genbank.Feature {
	return F.Pipe2(
		PR.MakePair(s, s),
		PR.BiMap(stateL.Features.Get, stateL.Primers.Get),
		PR.Merge(A.Concat[genbank.Feature]),
	)
}

func clipLocus(runes []rune) string { return string(runes[:min(len(runes), LocusNameMax)]) }

func constructStateDefinition(s parseState) string {
	return F.Pipe2(
		PR.MakePair(s, s),
		PR.BiMap(stateName.Get, stateL.Notes.Get),
		PR.Paired(definitionText),
	)
}

func definitionText(name string, notes []notesXML) string {
	return F.Pipe1(
		latestNote(notesType, notes),
		MO.Match(
			name+".",
			[]O.Kleisli[string, string]{
				MO.Case(S.IsNonEmpty, F.Bind1st(constructDefinition, name)),
			},
		),
	)
}

func constructDefinition(name, kind string) string {
	return F.Pipe4(kind, strings.ToLower, S.Append(" construct."), S.Prepend(", "), S.Prepend(name))
}

func latestNote(get func(notesXML) string, notes []notesXML) string {
	return F.Pipe4(
		notes,
		A.Map(get),
		A.Filter(S.IsNonEmpty),
		A.Last[string],
		O.GetOrElse(F.Constant("")),
	)
}

func commentText(notes []notesXML) map[string]string {
	return F.Pipe4(
		[]func([]notesXML) string{uuidText, createdText, modifiedText, transformedText},
		A.Ap[string]([][]notesXML{notes}),
		strutils.JoinSemicolon,
		strutils.MapNonEmpty(strutils.Period),
		commentQualifier,
	)
}

func uuidText(notes []notesXML) string {
	return F.Pipe3(
		notes, F.Bind1st(latestNote, notesUUID),
		removeTrademark,
		strutils.MapNonEmpty(S.Prepend("SnapGene UUID ")),
	)
}

func createdText(notes []notesXML) string {
	return F.Pipe3(
		notes,
		F.Bind1st(latestNote, notesCreated),
		removeTrademark,
		strutils.MapNonEmpty(S.Prepend("created ")),
	)
}

func modifiedText(notes []notesXML) string {
	return F.Pipe3(
		notes,
		F.Bind1st(latestNote, notesModified),
		removeTrademark,
		strutils.MapNonEmpty(S.Prepend("modified ")),
	)
}

func transformedText(notes []notesXML) string {
	return F.Pipe3(
		notes, F.Bind1st(latestNote, notesTransformed),
		removeTrademark,
		strutils.MapNonEmpty(S.Prepend("transformed into ")),
	)
}

var (
	isFeaturesPacket = F.Flow2(packetL.ID.Get, isFeaturesID)
	isPrimersPacket  = F.Flow2(packetL.ID.Get, isPrimersID)
	isNotesPacket    = F.Flow2(packetL.ID.Get, isNotesID)
	removeTrademark  = strings.NewReplacer("™", "").Replace
	commentQualifier = MO.Match(
		map[string]string{},
		[]O.Kleisli[string, map[string]string]{MO.Case(S.IsNonEmpty, makeCommentQualifier)},
	)
)

func makeCommentQualifier(comment string) map[string]string {
	return F.Pipe1(map[string]string{}, R.UpsertAt("COMMENT", comment))
}

func buildResult(s parseState) ParseResult {
	return F.Pipe3(
		PR.MakePair(s, s),
		PR.BiMap(
			F.Flow2(stateL.Record.Get, resultL.Record.Set),
			F.Flow2(stateL.Warnings.Get, resultL.Warnings.Set),
		),
		PR.Paired(En.Monoid[ParseResult]().Concat),
		En.Build[ParseResult],
	)
}
