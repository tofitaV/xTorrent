package main

import (
	"fmt"
	"github.com/zeebo/bencode"
)

type TorrentFile struct {
	Name              string         `bencode:"name"`
	PieceLength       int            `bencode:"piece length"`
	Pieces            string         `bencode:"pieces"`
	Announce          string         `bencode:"announce"`
	AnnounceList      [][]string     `bencode:"announce-list"`
	AzureusProperties map[string]int `bencode:"azureus_properties"`
	Comment           string         `bencode:"comment"`
	CreatedBy         string         `bencode:"created by"`
	Info              InfoData       `bencode:"info"`
}

type InfoData struct {
	Length int  `bencode:"length"`
	Files  File `bencode:"files"`
}

type File []struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func main() {
	file, _ := ReadFile()

	var torrentFile TorrentFile
	err := bencode.DecodeBytes(file, &torrentFile)
	if err != nil {
		return
	}
	fmt.Println()

	GetInfoFromTracker(torrentFile)
}