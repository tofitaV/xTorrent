package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	Udp  string = "udp"
	Http        = "http"
)

func GetInfoFromTracker(torrent *TorrentFile) {

	peerId, _ := CreatePeerId()

	baseURL := torrent.Announce
	queryParams := url.Values{
		"info_hash":  []string{string(torrent.InfoHash[:])},
		"peer_id":    []string{peerId},
		"uploaded":   []string{"0"},
		"downloaded": []string{strconv.Itoa(torrent.FileData.Downloaded)},
		"left":       []string{strconv.Itoa(torrent.FileData.Left)},
		"port":       []string{"6969"},
		"compact":    []string{"1"},
	}

	fullURL := baseURL + "?" + queryParams.Encode()

	client := &http.Client{}
	req, err := http.NewRequest(
		"GET", fullURL, nil,
	)
	fmt.Println(req.URL)
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

func GetProtocol(p string, t *TorrentFile) []string {
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
