package main

import (
	"errors"
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
	Left            int
	Downloaded      int
	TotalDataLength int
	Files           File `bencode:"files"`
}

type File []struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func main() {
	file, _ := ReadFile()
	torrentFile, _ := DecodeTorrent(file)
	torrentFile.Info.TotalLength()
	GetInfoFromTracker(torrentFile)
}

func DecodeTorrent(b []byte) (*TorrentFile, error) {
	var torrentFile TorrentFile
	err := bencode.DecodeBytes(b, &torrentFile)
	if err != nil {
		return &torrentFile, errors.New("can't decode a data")
	}

	return &torrentFile, err
}

func (info *InfoData) TotalLength() {
	total := 0
	for _, file := range info.Files {
		total += file.Length
	}
	info.Left = total - info.Downloaded
}
