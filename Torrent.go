package main

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	InfoHash          [20]byte
	PieceHashes       [][20]byte
	Publisher         string         `bencode:"publisher"`
	PublisherUrl      string         `bencode:"publisher-url"`
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
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Files       File   `bencode:"files"`
}

type File []struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func DecodeTorrent(b []byte) (*TorrentFile, error) {
	var torrentFile TorrentFile
	err := bencode.Unmarshal(bytes.NewReader(b), &torrentFile)
	if err != nil {
		return &torrentFile, errors.New("can't decode a data")
	}

	torrentFile.TotalLength()
	torrentFile.GetHashes()
	return &torrentFile, err
}

func DownloadFiles(filePath []byte) {
	torrentFile, err := DecodeTorrent(filePath)
	if err != nil {
		errors.New("can't get hash")
	}
	peers, _ := GetPeers(torrentFile)
	fmt.Println(peers)
}

func (t *TorrentFile) TotalLength() {
	total := 0
	for _, file := range t.Info.Files {
		total += file.Length
	}
	t.FileData.Left = total - t.FileData.Downloaded
}

func (i *Info) Hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *TorrentFile) SplitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Info.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (t *TorrentFile) GetHashes() {
	infoHash, err := t.Info.Hash()
	if err != nil {
		errors.New("can't get hash")
	}
	pieceHashes, err := t.SplitPieceHashes()
	if err != nil {
		errors.New("can't splite hashes")
	}
	t.InfoHash = infoHash
	t.PieceHashes = pieceHashes
}
