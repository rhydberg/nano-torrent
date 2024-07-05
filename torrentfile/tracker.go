package torrentfile

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/rhydberg/gotorrent/peers"
)

type  trackerResponse struct{ 
	Interval int `bencode:"interval"`
	Peers string `bencode:"peers"`
}

func (tf TorrentFile) buildTrackerURL() (string, error) {
	base, err := url.Parse(tf.Announce)
	if err != nil {
		return "", err
	}

	port := 6667
	fmt.Println(string(tf.InfoHash[:]))
	fmt.Println(string(tf.PeerID[:]))
	params := url.Values{
		"info_hash":  []string{string(tf.InfoHash[:])},
		"peer_id":    []string{string(tf.PeerID[:])},
		"port":       []string{strconv.Itoa(port)},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tf.Info.Length)},
	}

	// fmt.Printf("%v", params.Encode())

	base.RawQuery = params.Encode()
	// fmt.Printf("%v\n", base.String())
	return base.String(), nil

}

func (tf TorrentFile) RequestPeers() ([]peers.Peer, error) {
	url, err:= tf.buildTrackerURL()
	if err!=nil{
		log.Fatal("Could not build tracker URL")

	}

	client := &http.Client{Timeout: 15*time.Second}
	response, err:= client.Get(url)
	if err!=nil{
		log.Fatal("Could not get a response", err)
	}

	defer response.Body.Close()

	trResp := trackerResponse{}
	err = bencode.Unmarshal(response.Body, &trResp)
	if err!=nil{
		log.Fatal("Could not unmarshal requenst ", err)
	}

	return peers.GetPeers([]byte(trResp.Peers))

	
}
