package main

import "fmt"

// "fmt"
// "github.com/rhydberg/gotorrent/torrentfile"
// "log"

func fn(s []int){
	s = append(s[:0], 2,3)
}


func main() {
	// path := "debian-12.5.0-amd64-netinst.iso.torrent"

	// tf, err := torrentfile.GetTorrentFile(path)

	// if err != nil {
	// 	log.Fatalf("Error parsing torrent file")
	// }

	// // url, _ := tf.buildTrackerURL()
	// // fmt.Printf("%v", url)
	// tf.RequestPeers()

	slice := make([]int, 3)
	fmt.Printf("%v", slice)
	fn(slice)
	fmt.Printf("%v", slice)
}
