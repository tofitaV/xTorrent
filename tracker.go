package main

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/zeebo/bencode"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetInfoFromTracker(torrent TorrentFile) {

	peerId, _ := CreatePeerId()
	tHash, _ := calculateTorrentFileHash(torrent)

	baseURL := torrent.Announce
	queryParams := url.Values{
		"info_hash":  []string{tHash},
		"peer_id":    []string{peerId},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{"0"},
		"port":       []string{"6889"},
		"compact":    []string{"1"},
	}

	fullURL := baseURL + "?" + queryParams.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fullURL, nil,
	)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading", err)
		return
	}

	fmt.Println(string(body))
}

func CreatePeerId() (string, error) {
	randomPartLength := 6

	randomPart := make([]byte, randomPartLength)
	_, err := rand.Read(randomPart)
	if err != nil {
		return "", err
	}
	prefix := "-LT1100-"

	peerID := prefix + hex.EncodeToString(randomPart)

	return peerID, nil
}

func calculateTorrentFileHash(torrent TorrentFile) (string, error) {
	infoData, err := bencode.EncodeBytes(torrent.Info)
	if err != nil {
		return "", err
	}

	hash := sha1.Sum(infoData)

	infoHash := fmt.Sprintf("%x", hash)

	return infoHash, nil
}
