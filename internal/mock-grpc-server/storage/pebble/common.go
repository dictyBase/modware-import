package pebble

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	stock "github.com/dictyBase/go-genproto/dictybaseapis/stock"
	"github.com/dictyBase/modware-import/internal/config"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// jsonIndex represents the JSON structure for filtering
type jsonIndex map[string]interface{}

// buildStrainJSONIndex creates a JSON index from a strain for filtering
func buildStrainJSONIndex(strain *stock.Strain) jsonIndex {
	if strain == nil || strain.Data == nil {
		return jsonIndex{}
	}

	attrs := strain.Data.Attributes
	if attrs == nil {
		attrs = &stock.StrainAttributes{}
	}

	index := jsonIndex{
		"id":        strain.Data.Id,
		"type":      strain.Data.Type,
		"depositor": attrs.Depositor,
		"summary":   attrs.Summary,
		"species":   attrs.Species,
		"names":     attrs.Names,
	}

	if attrs.CreatedAt != nil {
		index["created_at"] = attrs.CreatedAt.AsTime().Format(time.RFC3339)
	}
	if attrs.UpdatedAt != nil {
		index["updated_at"] = attrs.UpdatedAt.AsTime().Format(time.RFC3339)
	}

	return index
}

// buildPlasmidJSONIndex creates a JSON index from a plasmid for filtering
func buildPlasmidJSONIndex(plasmid *stock.Plasmid) jsonIndex {
	if plasmid == nil || plasmid.Data == nil {
		return jsonIndex{}
	}

	attrs := plasmid.Data.Attributes
	if attrs == nil {
		attrs = &stock.PlasmidAttributes{}
	}

	index := jsonIndex{
		"id":        plasmid.Data.Id,
		"type":      plasmid.Data.Type,
		"depositor": attrs.Depositor,
		"summary":   attrs.Summary,
		"name":      attrs.Name,
	}

	if attrs.CreatedAt != nil {
		index["created_at"] = attrs.CreatedAt.AsTime().Format(time.RFC3339)
	}
	if attrs.UpdatedAt != nil {
		index["updated_at"] = attrs.UpdatedAt.AsTime().Format(time.RFC3339)
	}

	return index
}

// serializeStrain serializes a strain to protobuf bytes
func serializeStrain(strain *stock.Strain) ([]byte, error) {
	protoBytes, err := proto.Marshal(strain)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal strain: %w", err)
	}
	return protoBytes, nil
}

// deserializeStrain deserializes protobuf bytes to a strain
func deserializeStrain(data []byte) (*stock.Strain, error) {
	var strain stock.Strain
	if err := proto.Unmarshal(data, &strain); err != nil {
		return nil, fmt.Errorf("failed to unmarshal strain: %w", err)
	}
	return &strain, nil
}

// serializePlasmid serializes a plasmid to protobuf bytes
func serializePlasmid(plasmid *stock.Plasmid) ([]byte, error) {
	protoBytes, err := proto.Marshal(plasmid)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal plasmid: %w", err)
	}
	return protoBytes, nil
}

// deserializePlasmid deserializes protobuf bytes to a plasmid
func deserializePlasmid(data []byte) (*stock.Plasmid, error) {
	var plasmid stock.Plasmid
	if err := proto.Unmarshal(data, &plasmid); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plasmid: %w", err)
	}
	return &plasmid, nil
}

// serializeJSON serializes a value to JSON bytes
func serializeJSON(value interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return jsonBytes, nil
}

// encodeCounter encodes an int64 counter value to bytes
func encodeCounter(value int64) []byte {
	if value < 0 {
		value = 0
	}
	buf := make([]byte, config.DefaultBufferSize)
	binary.BigEndian.PutUint64(buf, uint64(value)) // #nosec G115 -- validated non-negative
	return buf
}

// decodeCounter decodes bytes to an int64 counter value
func decodeCounter(data []byte) int64 {
	if len(data) != config.DefaultBufferSize {
		return 0
	}
	uValue := binary.BigEndian.Uint64(data)
	// Prevent overflow: uint64 values > math.MaxInt64 would overflow when converted to int64
	if uValue > (1<<63 - 1) {
		return 0
	}
	return int64(uValue) // #nosec G115 -- validated within int64 range
}

// nowTimestamp returns the current time as a protobuf timestamp
func nowTimestamp() *timestamppb.Timestamp {
	return timestamppb.New(time.Now())
}

// formatStockID formats a stock ID with the given prefix and number
func formatStockID(prefix string, number int64) string {
	return fmt.Sprintf("%s%07d", prefix, number)
}

// parseStockIDNumber extracts the numeric part from a stock ID (e.g., "DBS0000002" → 2)
func parseStockIDNumber(stockID string, prefixLen int) int64 {
	if len(stockID) <= prefixLen {
		return 0
	}
	numStr := stockID[prefixLen:]
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
