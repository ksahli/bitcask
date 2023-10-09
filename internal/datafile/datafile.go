package datafile

import (
	"os"
)

const (
	flags       = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	permissions = 0600
)

type Datafile struct {
	file   *os.File
	offset int64
}

func (df *Datafile) Write(payload []byte) (int64, int64, error) {
	n, err := df.file.Write(payload)
	if err != nil {
		return -1, -1, err
	}

	size := int64(n)

	offset := df.offset
	df.offset = df.offset + size
	return offset, size, nil
}

func (df *Datafile) Close() error {
	return df.file.Close()
}

func Open(path string) (*Datafile, error) {
	file, err := os.OpenFile(path, flags, permissions)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := info.Size()

	datafile := Datafile{
		file:   file,
		offset: size,
	}

	return &datafile, nil
}
