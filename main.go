package main

import (
	"log"

	"github.com/ksahli/bitcask/internal/datafile"
)

func main() {

	df, err := datafile.Open("test.data")
	if err != nil {
		log.Fatal(err)
	}

	defer df.Close()

	offset, size, err := df.Write(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(offset)
	log.Println(size)
}
