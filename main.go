package main

func main() {
	file, _ := ReadFile()
	torrentFile, _ := DecodeTorrent(file)
	torrentFile.TotalLength()
	torrentFile.GetHashes()
	GetInfoFromTracker(torrentFile)
}
