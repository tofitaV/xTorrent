package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jackpal/bencode-go"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	Udp  string = "udp"
	Http        = "http"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func GetPeers(torrent *TorrentFile) ([]Peer, error) {
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
		return nil, err
	}

	defer resp.Body.Close()

	trackerResponse := TrackerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return nil, err
	}

	peers, err := ParsePeers([]byte(trackerResponse.Peers))
	return peers, err
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
