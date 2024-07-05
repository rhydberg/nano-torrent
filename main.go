package main

import (
	// "fmt"
	"flag"
	// "fmt"
	"log"
	"os"

	// "github.com/rhydberg/gotorrent/client"
	// "github.com/rhydberg/gotorrent/handshake"
	"github.com/rhydberg/gotorrent/torrentfile"
)



func fn(s []int){
	s = append(s[:0], 2,3)
}




func main() {
	debian_path := "debian-12.5.0-amd64-netinst.iso.torrent"

	path:= flag.String("p", debian_path, "The path to the torrent file")
	out:= flag.String("o", "out", "the output file")
	flag.Parse()

	tf, err := torrentfile.GetTorrentFile(*path)

	if err != nil {
		log.Fatalf("Error parsing torrent file")
	}

	outFile,err := os.Create(*out)
	buf, err := tf.Download()
	if err!=nil{
		log.Fatal("Error downloading", err)
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)

	if err!=nil{
		log.Fatal("Error writing to file ", err)
	}




	// fmt.Printf("%+v", tf.PieceHashes)

	// url, _ := tf.buildTrackerURL()
	// fmt.Printf("%v", url)
	// p, err:= tf.RequestPeers()

	// for _,peer := range(p){
	// 	_=handshake.New(tf.InfoHash, tf.PeerID)
	// 	_,_=client.New(peer, tf.InfoHash, tf.PeerID)
	// }


}
