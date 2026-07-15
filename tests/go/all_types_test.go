package gotests

import (
	"encoding/json"
	"reflect"
	"testing"

	m "github.com/bbanez/minibin/tests/dist/go"
)

func ptr[T any](value T) *T { return &value }

func TestAllTypesDocumentRoundTrip(t *testing.T) {
	document := &m.AllTypesDocument{
		Id:                  "document-001",
		Title:               ptr("A fully populated document"),
		Aliases:             []string{"primary", "backup"},
		OptionalAliases:     []*string{ptr("legacy")},
		MinimumI32:          -2000000000,
		OptionalI32:         ptr(int32(2000000000)),
		I32Values:           []int32{-1, 0, 1},
		OptionalI32Values:   []*int32{ptr(int32(-42)), ptr(int32(42))},
		MinimumI64:          -9000000000000,
		OptionalI64:         ptr(int64(9000000000000)),
		I64Values:           []int64{-9_000_000_000, 0, 9_000_000_000},
		OptionalI64Values:   []*int64{ptr(int64(-7)), ptr(int64(7))},
		MaximumU32:          4294967295,
		OptionalU32:         ptr(uint32(42)),
		U32Values:           []uint32{0, 1, 4294967295},
		OptionalU32Values:   []*uint32{ptr(uint32(9))},
		MaximumU64:          9000000000000,
		OptionalU64:         ptr(uint64(99)),
		U64Values:           []uint64{0, 1, 9000000000000},
		OptionalU64Values:   []*uint64{ptr(uint64(11))},
		Price:               -123.45,
		OptionalRatio:       ptr(float32(0.12345)),
		F32Values:           []float32{-2.5, 0, 3.125},
		OptionalF32Values:   []*float32{ptr(float32(-1.2)), ptr(float32(4.5))},
		Measurement:         -987654.1234567,
		OptionalMeasurement: ptr(0.1256),
		F64Values:           []float64{-1.234567, 0, 9.876543},
		OptionalF64Values:   []*float64{ptr(12.34)},
		Enabled:             true,
		OptionalEnabled:     ptr(false),
		Flags:               []bool{true, false, true},
		OptionalFlags:       []*bool{ptr(false), ptr(true)},
		State:               m.LIFECYCLE_IN_REVIEW,
		Priorities:          []m.Priority{m.PRIORITY_LOW, m.PRIORITY_CRITICAL},
		Payload:             []byte{0, 1, 127, 128, 255},
		OptionalPayload:     ptr([]byte{9, 8, 7}),
		Chunks:              [][]byte{{}, {1}, {2, 3, 4}},
		OptionalChunks:      []*[]byte{ptr([]byte{5, 6})},
		Origin:              m.Coordinate{Latitude: 51.5, Longitude: -0.1, AltitudeMeters: ptr(float32(35.25))},
		OptionalDestination: &m.Coordinate{Latitude: -33.8, Longitude: 151.2},
		Attachments: []*m.Attachment{{
			FileName: "report.bin",
			Content:  []byte{1, 2, 3},
			Checksum: ptr([]byte{4, 5}),
			Labels:   []*string{ptr("binary"), ptr("report")},
			Sizes:    []uint64{3, 1024},
		}},
		OptionalAuditTrail: []*m.AuditEntry{{
			Sequence:      1,
			ActorId:       -42,
			Message:       "published",
			PreviousState: m.LIFECYCLE_DRAFT,
			NextState:     m.LIFECYCLE_PUBLISHED,
			ChangedFields: []string{"state", "title"},
			Location:      &m.Coordinate{Latitude: 40.7, Longitude: -74},
		}},
	}

	unpacked, err := m.UnpackAllTypesDocument(document.Pack(), nil)
	if err != nil {
		t.Fatalf("unpack document: %v", err)
	}
	if !reflect.DeepEqual(document, unpacked) {
		t.Fatalf("round trip mismatch\nwant: %#v\n got: %#v", document, unpacked)
	}
}

func TestAllTypesDocumentBatchRoundTrip(t *testing.T) {
	document := m.AllTypesDocument{Id: "minimal", Aliases: []string{}, I32Values: []int32{}, I64Values: []int64{}, U32Values: []uint32{}, U64Values: []uint64{}, F32Values: []float32{}, F64Values: []float64{}, Flags: []bool{}, State: m.LIFECYCLE_DRAFT, Priorities: []m.Priority{}, Payload: []byte{}, Chunks: [][]byte{}, Origin: m.Coordinate{}, Attachments: []*m.Attachment{}}
	batch := &m.AllTypesDocumentBatch{BatchId: "batch-001", Documents: []*m.AllTypesDocument{&document}}

	unpacked, err := m.UnpackAllTypesDocumentBatch(batch.Pack(), nil)
	if err != nil {
		t.Fatalf("unpack batch: %v", err)
	}
	if unpacked.BatchId != batch.BatchId || len(unpacked.Documents) != 1 || unpacked.Documents[0].Id != document.Id {
		t.Fatalf("batch round trip mismatch: %#v", unpacked)
	}
}

func TestAllTypesDocumentAccessorsMutationAndCopy(t *testing.T) {
	document := &m.AllTypesDocument{
		Id: "before", Aliases: []string{"first"}, I32Values: []int32{1}, I64Values: []int64{2}, U32Values: []uint32{3}, U64Values: []uint64{4},
		F32Values: []float32{1.5}, F64Values: []float64{2.25}, Flags: []bool{true}, State: m.LIFECYCLE_DRAFT, Priorities: []m.Priority{m.PRIORITY_LOW},
		Payload: []byte{1}, Chunks: [][]byte{{2}}, Origin: m.Coordinate{Latitude: 1, Longitude: 2}, Attachments: []*m.Attachment{},
	}

	document.SetId("after")
	document.SetTitle(ptr("optional title"))
	document.SetAliases([]string{"updated", "aliases"})
	document.SetOptionalI32(ptr(int32(-99)))
	document.SetState(m.LIFECYCLE_ARCHIVED)
	document.SetPayload([]byte{9, 8, 7})
	document.SetOrigin(m.Coordinate{Latitude: 3, Longitude: 4, AltitudeMeters: ptr(float32(5.5))})

	if document.GetId() != "after" || *document.GetTitle() != "optional title" || !reflect.DeepEqual(document.GetAliases(), []string{"updated", "aliases"}) {
		t.Fatal("string accessor mutation was not retained")
	}
	if *document.GetOptionalI32() != -99 || document.GetState() != m.LIFECYCLE_ARCHIVED || !reflect.DeepEqual(document.GetPayload(), []byte{9, 8, 7}) {
		t.Fatal("scalar accessor mutation was not retained")
	}
	if document.GetOrigin().Latitude != 3 || document.GetOrigin().AltitudeMeters == nil || *document.GetOrigin().AltitudeMeters != 5.5 {
		t.Fatal("nested object accessor mutation was not retained")
	}

	copy := document.Copy()
	copy.SetId("copy")
	copy.SetAliases([]string{"copy-only"})
	copyOrigin := copy.GetOrigin()
	copyOrigin.Latitude = 99
	copy.SetOrigin(copyOrigin)
	if document.GetId() != "after" || !reflect.DeepEqual(document.GetAliases(), []string{"updated", "aliases"}) || document.GetOrigin().Latitude != 3 {
		t.Fatal("copy mutation changed the original document")
	}

	unpacked, err := m.UnpackAllTypesDocument(document.Pack(), nil)
	if err != nil {
		t.Fatalf("unpack mutated document: %v", err)
	}
	if unpacked.GetId() != "after" || unpacked.GetState() != m.LIFECYCLE_ARCHIVED || !reflect.DeepEqual(unpacked.GetPayload(), []byte{9, 8, 7}) || unpacked.GetOrigin().Latitude != 3 {
		t.Fatal("mutated values did not survive round trip")
	}
}

func TestAllTypesDocumentRejectsMalformedPayloads(t *testing.T) {
	testCases := [][]byte{
		{0},                      // Missing type/length byte.
		{0, 0x10},                // String missing its length byte.
		{0, 0x10, 2, 'x'},        // String length exceeds remaining data.
		{0, 0x28, 0, 0, 0, 0, 0}, // Int32 declares a nine-byte value.
		{0, 0x80},                // Bool missing its value byte.
	}

	for _, payload := range testCases {
		if _, err := m.UnpackAllTypesDocument(payload, nil); err == nil {
			t.Fatalf("expected malformed payload %v to fail", payload)
		}
	}
}

func TestAllTypesDocumentRejectsMismatchedWireType(t *testing.T) {
	// Field 0 is Id (string); encode it with the int32 wire tag instead.
	if _, err := m.UnpackAllTypesDocument([]byte{0, 0x20, 0, 1}, nil); err == nil {
		t.Fatal("expected mismatched wire type to fail")
	}
}

func TestAllTypesDocumentToJsonEscapesStrings(t *testing.T) {
	document := &m.AllTypesDocument{Id: "quote: \" and control: \x01", State: m.LIFECYCLE_DRAFT}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(document.ToJson()), &decoded); err != nil {
		t.Fatalf("ToJson returned invalid JSON: %v", err)
	}
	if decoded["id"] != document.Id {
		t.Fatalf("JSON value mismatch: %#v", decoded)
	}
}

func TestAllTypesDocumentRejectsUnknownEnumValue(t *testing.T) {
	// Field 32 is State; enum values use wire tag 10.
	payload := []byte{32, 0xA0, 7, 'm', 'i', 's', 's', 'i', 'n', 'g'}
	if _, err := m.UnpackAllTypesDocument(payload, nil); err == nil {
		t.Fatal("expected unknown enum value to fail")
	}
}

func TestAllTypesDocumentUsesEnumWireTag(t *testing.T) {
	document := &m.AllTypesDocument{State: m.LIFECYCLE_DRAFT}
	packed := document.Pack()
	found := false
	for i := 0; i+1 < len(packed); i++ {
		if packed[i] == 32 && packed[i+1]>>4 == 10 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected enum wire tag 10, got payload %v", packed)
	}
}
