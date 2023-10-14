package records

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

type Record struct {
	key, value []byte
}

func (r Record) Key() []byte {
	return r.key
}

func (r Record) Value() []byte {
	return r.value
}

func (r Record) Bytes() []byte {
	data := append(r.key, r.value...)
	return data
}

func (r Record) Checksum() uint32 {
	checksum := crc32.Checksum(r.Bytes() , crc32.IEEETable)
	return checksum
}

func New(key, value []byte) Record {
	record := Record{
		key:      key,
		value:    value,
	}
	return record
}

func Encode(record Record, encoder binary.ByteOrder) []byte {
	metadata := make([]byte, 16)

	checksum := record.Checksum()
	encoder.PutUint32(metadata[:4], checksum)

	keysz, valuesz := len(record.Key()), len(record.Value())
	encoder.PutUint32(metadata[4:8], uint32(keysz))
	encoder.PutUint64(metadata[8:16], uint64(valuesz))

	data := record.Bytes()

	return append(metadata, data...)
}

func Decode(payload []byte, decoder binary.ByteOrder) (Record, error) {
	checksum := decoder.Uint32(payload[:4])

	keysz := decoder.Uint32(payload[4:8])

	key, value := payload[16:16+keysz], payload[16+keysz:]
	record := New(key, value) 

	if record.Checksum() != checksum {
		return Record{}, errors.New("invalid record")
	}

	return record, nil
}
