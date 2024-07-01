package main

import (
	// "fmt"
	"github.com/rhydberg/gotorrent/torrentfile"
	"log"
)

func main() {
	path := "debian-12.5.0-amd64-netinst.iso.torrent"

	tf, err := torrentfile.GetTorrentFile(path)

	if err != nil {
		log.Fatalf("Error parsing torrent file")
	}

	// url, _ := tf.buildTrackerURL()
	// fmt.Printf("%v", url)
	tf.RequestPeers()
}
