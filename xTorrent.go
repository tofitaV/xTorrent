package main

import (
	"errors"
	"github.com/zeebo/bencode"
)

type TorrentFile struct {
	Publisher         string
	PublisherUrl      string
	Announce          string         `bencode:"announce"`
	AnnounceList      [][]string     `bencode:"announce-list"`
	AzureusProperties map[string]int `bencode:"azureus_properties"`
	Comment           string         `bencode:"comment"`
	CreatedBy         string         `bencode:"created by"`
	Info              Info           `bencode:"info"`
	FileData          FileData
}

type FileData struct {
	Left            int
	Downloaded      int
	TotalDataLength int
}

type Info struct {
	Name        string
	PieceLength int
	Pieces      string
	Files       File `bencode:"files"`
}

type File []struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func main() {
	file, _ := ReadFile()
	torrentFile, _ := DecodeTorrent(file)
	torrentFile.TotalLength()
	GetInfoFromTracker(torrentFile)
}

func DecodeTorrent(b []byte) (*TorrentFile, error) {
	var torrentFile TorrentFile
	err := bencode.DecodeBytes(b, &torrentFile)
	if err != nil {
		return &torrentFile, errors.New("can't decode a data")
	}

	var topMap map[string]interface{}
	err = bencode.DecodeBytes(b, &topMap)
	if err != nil {
		return &torrentFile, errors.New("can't decode a data")
	}

	torrentFile.Info.Name = topMap["info"].(map[string]interface{})["name"].(string)
	torrentFile.Info.PieceLength = int(topMap["info"].(map[string]interface{})["piece length"].(int64))
	torrentFile.Info.Pieces = topMap["info"].(map[string]interface{})["pieces"].(string)
	torrentFile.Publisher = topMap["publisher"].(string)
	torrentFile.PublisherUrl = topMap["publisher-url"].(string)

	return &torrentFile, err
}

func (t *TorrentFile) TotalLength() {
	total := 0
	for _, file := range t.Info.Files {
		total += file.Length
	}
	t.FileData.Left = total - t.FileData.Downloaded
}
