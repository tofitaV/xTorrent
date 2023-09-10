package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func ParsePeers(peersByte []byte) ([]Peer, error) {
	const peerSize = 6
	numPeers := len(peersByte) / peerSize
	if len(peersByte)%peerSize != 0 {
		err := fmt.Errorf("Received malformed peers")
		return nil, err
	}
	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peersByte[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersByte[offset+4 : offset+6]))
	}
	return peers, nil
}
