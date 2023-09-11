package main

import (
	client "awesomeProject/Client"
	"awesomeProject/Reader"
	"awesomeProject/Torrent"
	"awesomeProject/Tracker"
)

func main() {
	var filePath = "InputData/Starfield.torrent"
	file, _ := Reader.ReadFile(filePath)
	var torrentFile *Torrent.TorrentFile
	torrentFile = Torrent.DownloadFiles(file)
	var peers, peerId, _ = Tracker.GetPeers(torrentFile)

	torrent := client.Torrent{
		Peers:       peers,
		PeerID:      peerId,
		InfoHash:    torrentFile.InfoHash,
		Name:        torrentFile.Info.Name,
		Length:      torrentFile.Info.Length,
		PieceHashes: torrentFile.PieceHashes,
		PieceLength: torrentFile.Info.PieceLength,
		Announce:    torrentFile.Announce,
	}

	torrent.Start()
}
