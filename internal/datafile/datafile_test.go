package datafile_test

import (
	"path/filepath"
	"testing"

	"github.com/ksahli/bitcask/internal/datafile"
)

func TestOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.data")

	df, err := datafile.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer df.Close()
}

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.data")

	df, err := datafile.Open(path)
	if err != nil {
		t.Error(err)
	}

	payload := []byte("test payload")
	offset, size, err := df.Write(payload)
	if err != nil {
		t.Error(err)
	}

	expectedOffset := int64(0)
	if expectedOffset != offset {
		t.Errorf("want offset %d, got %d", expectedOffset, offset)
	}

	expectedSize := int64(len(payload))
	if expectedSize != size {
		t.Errorf("expecting size %d, got %d", expectedSize, size)
	}
}
