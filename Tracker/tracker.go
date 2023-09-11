package Tracker

import (
	"awesomeProject/Peers"
	"awesomeProject/Torrent"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jackpal/bencode-go"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func GetPeers(t *Torrent.TorrentFile) ([]Peers.Peer, string, error) {

	peerId, _ := CreatePeerId()

	baseURL := t.Announce
	queryParams := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{peerId},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(t.Info.Length)},
		"port":       []string{"6969"},
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
		return nil, "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	trackerResponse := TrackerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return nil, "", err
	}

	peers, err := Peers.ParsePeers([]byte(trackerResponse.Peers))
	return peers, peerId, err
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
