package segments

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"

	"github.com/ksahli/bitcask/internal/records"
)

type File interface {
	fs.File
	io.Writer
}

// Segment is a data file
type Segment struct {
	file   File
	offset int64
}

// Offset of next append
func (s Segment) Offset() int64 {
	return s.offset
}

// Append a record
func (s *Segment) Append(record records.Record) (int64, error) {
	writer := bufio.NewWriter(s.file)

	payload := records.Encode(record, binary.BigEndian)

	if _, err := writer.Write(payload); err != nil {
		werr := fmt.Errorf("failed to write value: %w", err)
		return -1, werr
	}

	offset, increment := s.offset, writer.Size()
	if err := writer.Flush(); err != nil {
		werr := fmt.Errorf("failed to write record: %w", err)
		return -1, werr
	}

	s.offset = s.offset + int64(increment)
	return offset, nil
}

// Read a record at a given offset
func (s *Segment) Read(int64) (records.Record, error) {
	return records.Record{}, nil
}

// New segment
func New(file File) (*Segment, error) {
	info, err := file.Stat()
	if err != nil {
		werr := fmt.Errorf("failed to get file info: %w", err)
		return nil, werr
	}
	offset := info.Size()

	segment := Segment{
		file:   file,
		offset: offset,
	}
	return &segment, nil
}
