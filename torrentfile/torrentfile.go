package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	// "fmt"
	"log"
	"os"

	"github.com/jackpal/bencode-go"
)

type torrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type TorrentFile struct {
	Announce     string `bencode:"announce"`
	Comment      string `bencode:"comment"`
	CreationDate uint32 `bencode:"creation date"`
	InfoHash     [20]byte
	Info         torrentInfo `bencode:"info"`
	PeerID       [20]byte
}

func (tf TorrentFile) getInfoHash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, tf.Info)
	if err != nil {
		return [20]byte{}, err
	}

	hash := sha1.Sum(buf.Bytes())
	return hash, nil
}

func (tf TorrentFile) getPeerID() ([20]byte, error) {
	var pid [20]byte
	_, err := rand.Read(pid[:])
	// fmt.Printf("PID IS %v\n", []string{string(pid[:])})
	if err != nil {
		return pid, err
	}

	return pid, nil
}

func GetTorrentFile(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	tf := TorrentFile{}
	err = bencode.Unmarshal(file, &tf)
	if err != nil {
		return TorrentFile{}, err
	}

	tf.InfoHash, err = tf.getInfoHash()
	if err != nil {
		log.Fatal("Unable to get info hash")
	}

	tf.PeerID, err = tf.getPeerID()

	if err != nil {
		log.Fatal("Could not generate Peer ID")
	}
	return tf, nil
}
