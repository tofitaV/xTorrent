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
	"strings"
)

const (
	Udp  string = "udp"
	Http        = "http"
)

func GetInfoFromTracker(torrent TorrentFile) {

	peerId, _ := CreatePeerId()
	tHash, _ := calculateTorrentFileHash(torrent)

	mainUrl := GetProtocol(Http, torrent)[0]

	baseURL := mainUrl
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

func GetProtocol(p string, t TorrentFile) []string {
	var httpAnnounces []string
	for _, announceGroup := range t.AnnounceList {
		for _, announce := range announceGroup {
			if strings.Contains(announce, p) {
				httpAnnounces = append(httpAnnounces, announce)
			}
		}
	}
	return httpAnnounces
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
