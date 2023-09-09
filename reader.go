package main

import (
	"errors"
	"os"
)

func ReadFile() ([]byte, error) {
	file, err := os.ReadFile("Starfield.torrent")
	if err != nil {
		return nil, errors.New("can't read a file")
	}
	return file, err
}
