package Client

import (
	"awesomeProject/Peers"
	"fmt"
	"net"
	"time"
)

type Torrent struct {
	Peers       []Peers.Peer
	PeerID      string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	Announce    string
}

type pieceResult struct {
	index int
	buf   []byte
}

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

func (t *Torrent) calculateBoundsForPiece(index int) (begin int, end int) {
	begin = index * t.PieceLength
	end = begin + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return begin, end
}
func (t *Torrent) calculatePieceSize(index int) int {
	begin, end := t.calculateBoundsForPiece(index)
	return end - begin
}

func (t *Torrent) Start() {

	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)

	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		workQueue <- &pieceWork{index, hash, length}
	}

	for _, peer := range t.Peers {
		go StartDownloadFile(peer, workQueue, results)
	}

}

func StartDownloadFile(peer Peers.Peer, workQueue chan *pieceWork, pieceResult chan *pieceResult) {
	conn, err := net.DialTimeout("tcp", Peers.ToString(peer), 2*time.Second)
	if err != nil {
		fmt.Printf("Error connecting to the peer: %s, \n%s", peer.IP, err)
	}

	fmt.Println("Connected to", peer.IP)

	//err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return
	}

	defer func(conn net.Conn, t time.Time) {
		err := conn.SetDeadline(t)
		if err != nil {

		}
	}(conn, time.Time{})

	defer conn.Close()
}
