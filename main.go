package main

import (
	// "fmt"
	"fmt"
	"log"

	// "github.com/rhydberg/gotorrent/client"
	// "github.com/rhydberg/gotorrent/handshake"
	"github.com/rhydberg/gotorrent/torrentfile"
)



func fn(s []int){
	s = append(s[:0], 2,3)
}




func main() {
	path := "debian-12.5.0-amd64-netinst.iso.torrent"

	tf, err := torrentfile.GetTorrentFile(path)

	if err != nil {
		log.Fatalf("Error parsing torrent file")
	}

	fmt.Printf("%+v", tf.PieceHashes)

	// url, _ := tf.buildTrackerURL()
	// fmt.Printf("%v", url)
	p, err:= tf.RequestPeers()

	// for _,peer := range(p){
	// 	_=handshake.New(tf.InfoHash, tf.PeerID)
	// 	_,_=client.New(peer, tf.InfoHash, tf.PeerID)
	// }


}
