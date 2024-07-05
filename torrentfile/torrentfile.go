package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"

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
	PieceHashes  [][20]byte
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

func (tf TorrentFile) getPieceHashes() ([][20]byte, error){
	lenHash := 20
	buf:=[]byte(tf.Info.Pieces)
	if len(buf)%lenHash !=0{
		return nil, fmt.Errorf("pieces not a multiple of 20")
	}

	numHashes := len(buf)/lenHash

	pieceHashes := make([][20]byte, numHashes)

	for i:=0;i<numHashes;i++{
		offset := i*lenHash
		copy(pieceHashes[i][:], buf[offset:offset+lenHash])
	}

	return pieceHashes, nil
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



	tf.PieceHashes, err = tf.getPieceHashes()

	if err!=nil{
		log.Fatal("could not split into hashes, ",err)
	}

	tf.PeerID, err = tf.getPeerID()

	if err != nil {
		log.Fatal("Could not generate Peer ID")
	}
	return tf, nil
}
