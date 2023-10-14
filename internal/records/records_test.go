package records_test

import (
	"encoding/binary"
	"hash/crc32"
	"reflect"
	"testing"

	"github.com/ksahli/bitcask/internal/records"
)

var (
	key   = []byte("test_key")
	value = []byte("test_value")
)

func TestEncode(t *testing.T) {
	encoder := binary.BigEndian

	record := records.New(key, value)

	payload := records.Encode(record, encoder)

	checksum := binary.BigEndian.Uint32(payload[:4])
	if checksum != record.Checksum() {
		t.Errorf("expecting checksum to be %d, got %d", checksum, record.Checksum())
	}

	keysz := binary.BigEndian.Uint32(payload[4:8])

	decoded := records.New(payload[16:16+keysz], payload[16+keysz:])
	if !reflect.DeepEqual(decoded, record) {
		t.Errorf("expecting decoded record to be %v got %v", decoded, record)
	}
}

func TestDecode(t *testing.T) {
	data := append(key, value...)
	checksum := crc32.Checksum(data , crc32.IEEETable)

	metadata := make([]byte, 16)
	binary.BigEndian.PutUint32(metadata[:4], checksum)

	keysz, valuesz := len(key), len(value)
	binary.BigEndian.PutUint32(metadata[4:8], uint32(keysz))
	binary.BigEndian.PutUint64(metadata[8:16], uint64(valuesz))

	payload := append(metadata, data...)

	record, err := records.Decode(payload, binary.BigEndian)
	if err != nil {
		t.Error(err)
	}

	decoded := records.New(key, value)
	if !reflect.DeepEqual(decoded, record) {
		t.Errorf("expecting decoded record to be %v got %v", decoded, record)
	}
	
}
