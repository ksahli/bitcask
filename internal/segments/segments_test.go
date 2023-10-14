package segments_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ksahli/bitcask/internal/records"
	"github.com/ksahli/bitcask/internal/segments"
)

func TestNew(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "test.data")
	file, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}

	payload := []byte("payload")
	if _, err := file.Write(payload); err != nil {
		t.Error(err)
	}

	segment, err := segments.New(file)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedOffset := int64(len(payload))
	if segment.Offset() != expectedOffset {
		t.Errorf("want offset to be %d got %d'", expectedOffset, segment.Offset())
	}
}

func TestAppend(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "test.data")
	file, err := os.Create(path)
	if err != nil {
		t.Error(err)
	}

	segment, err := segments.New(file)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	key, value := []byte("key"), []byte("value")

	{
		record := records.New(key, value)
		offset, err := segment.Append(record)
		if err != nil {
			t.Error(err)
		}

		if offset != 0 {
			t.Errorf("expecting offset to be 0, got %d", offset)
		}
	}

	{
		record := records.New(key, value)
		offset, err := segment.Append(record)
		if err != nil {
			t.Error(err)
		}

		if offset != 4096 {
			t.Errorf("expecting offset to be 4096, got %d", offset)
		}
	}
}
