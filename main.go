package main

func main() {
	var filePath = "Starfield.torrent"
	file, _ := ReadFile(filePath)
	DownloadFiles(file)
}
