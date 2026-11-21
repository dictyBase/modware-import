package snapgene

import (
	"encoding/binary"
	"errors"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	En "github.com/IBM/fp-go/v2/endomorphism"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	PR "github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
)

// MaxPacketBytes bounds each packet body, not the file size.
const MaxPacketBytes uint32 = 32 << 20

const (
	packetHeaderBytes      = 5
	cookieBytes            = 14
	cookieSignature        = "SnapGene"
	dnaID             byte = 0
	featuresID        byte = 10
	primersID         byte = 5
	notesID           byte = 6
	cookieID          byte = 9
	circularFlag      byte = 1
)

var (
	// ErrInvalidMagicCookie identifies a missing or invalid first cookie packet.
	ErrInvalidMagicCookie = errors.New("invalid SnapGene magic cookie")
	// ErrPacketTooLarge identifies a body that exceeds MaxPacketBytes.
	ErrPacketTooLarge = errors.New("SnapGene packet too large")
	// ErrTruncatedPacket identifies incomplete packet headers or bodies.
	ErrTruncatedPacket = errors.New("truncated SnapGene packet")
	// ErrMissingDNAPacket identifies a construct with no DNA packet.
	ErrMissingDNAPacket = errors.New("missing SnapGene DNA packet")
	scanL               = MakescanStateLenses()
	byteEq              = EQ.FromStrictEquals[byte]()
	predordByteDNA      = F.Pipe1(dnaID, EQ.Equals(byteEq))
	isFeaturesID        = F.Pipe1(featuresID, EQ.Equals(byteEq))
	isPrimersID         = F.Pipe1(primersID, EQ.Equals(byteEq))
	isNotesID           = F.Pipe1(notesID, EQ.Equals(byteEq))
	atHeader            = P.Fold(E.Of[error, scanState], scanHeader)
	// retainPayload appends only cookie, DNA, and annotation packets.
	retainPayload = F.Pipe1(
		isRetainedPacket,
		P.Fold(F.Identity[scanState], appendPacket),
	)
	scanData = F.Flip(
		F.Bind1st(A.ReduceWithIndex[byte, E.Either[error, scanState]], scanOffset),
	)
)

// fp-go:Lens
type scanState struct {
	Data    []byte
	Next    int
	Size    uint32
	Packets []Packet
}

// Packets checks all frames and retains only cookie, DNA, and annotation payloads.
// Unknown and history packets still receive full length validation.
func Packets(data []byte) E.Either[error, []Packet] {
	return F.Pipe5(
		scanState{},
		scanL.Data.Set(data),
		E.Of[error, scanState],
		scanData(data),
		E.Map[error](
			scanL.Packets.Get,
		),
		E.Chain(E.FromPredicate(validCookie, F.Constant1[[]Packet](ErrInvalidMagicCookie))),
	)
}

func scanOffset(index int, acc E.Either[error, scanState], _ byte) E.Either[error, scanState] {
	return F.Pipe1(acc, E.Chain(F.Bind1st(scanAt, index)))
}

func scanAt(index int, s scanState) E.Either[error, scanState] {
	return F.Pipe1(s, atHeader(F.Bind1st(isNextHeader, index)))
}

func isNextHeader(index int, s scanState) bool { return index == scanL.Next.Get(s) }

func scanHeader(s scanState) E.Either[error, scanState] {
	return F.Pipe7(
		s,
		E.Of[error, scanState],
		E.Chain(E.FromPredicate(headerFits, F.Constant1[scanState](ErrTruncatedPacket))),
		E.Chain(
			E.FromPredicate(firstPacketIsCookie, F.Constant1[scanState](ErrInvalidMagicCookie)),
		),
		E.Let[error](scanL.Size.Set, packetLength),
		E.Chain(E.FromPredicate(packetSizeFits, F.Constant1[scanState](ErrPacketTooLarge))),
		E.Chain(E.FromPredicate(bodyFits, F.Constant1[scanState](ErrTruncatedPacket))),
		E.Map[error](keepPacket),
	)
}

func headerFits(s scanState) bool {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	return len(data)-next >= packetHeaderBytes
}

func firstPacketIsCookie(s scanState) bool {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	return next != 0 || data[next] == cookieID
}

func packetLength(s scanState) uint32 {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	return binary.BigEndian.Uint32(data[next+1 : next+packetHeaderBytes])
}

func packetSizeFits(s scanState) bool { return scanL.Size.Get(s) <= MaxPacketBytes }

func bodyFits(s scanState) bool {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	size := scanL.Size.Get(s)
	return int64(len(data)-next-packetHeaderBytes) >= int64(size)
}

func keepPacket(s scanState) scanState {
	return F.Pipe2(
		PR.MakePair(s, s),
		PR.BiMap(retainPayload, F.Flow2(nextPacketOffset, scanL.Next.Set)),
		PR.Merge(F.Identity[En.Endomorphism[scanState]]),
	)
}

func nextPacketOffset(s scanState) int {
	return scanL.Next.Get(s) + packetHeaderBytes + int(scanL.Size.Get(s))
}

func isRetainedPacket(s scanState) bool {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	id := data[next]
	return id == cookieID || id == dnaID || id == featuresID || id == primersID || id == notesID
}

func appendPacket(s scanState) scanState {
	return F.Pipe4(s, packetAt, A.Push[Packet], scanL.Packets.Modify, En.Read(s))
}

func packetAt(s scanState) Packet {
	return F.Pipe2(packetEdits, A.Ap[En.Endomorphism[Packet]]([]scanState{s}), En.Reduce[Packet])
}

func packetTag(s scanState) byte {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	return data[next]
}

func packetBody(s scanState) []byte {
	data := scanL.Data.Get(s)
	next := scanL.Next.Get(s)
	size := scanL.Size.Get(s)
	return data[next+packetHeaderBytes : next+packetHeaderBytes+int(size)]
}

func validCookie(packets []Packet) bool {
	return F.Pipe2(packets, A.Head[Packet], O.Exists(cookieMatches))
}

func cookieMatches(packet Packet) bool {
	body := packetL.Body.Get(packet)
	return packetL.ID.Get(packet) == cookieID && len(body) == cookieBytes &&
		string(body[:len(cookieSignature)]) == cookieSignature
}

var packetEdits = []func(scanState) En.Endomorphism[Packet]{
	F.Flow2(packetTag, packetL.ID.Set),
	F.Flow2(packetBody, packetL.Body.Set),
}
