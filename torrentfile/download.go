package torrentfile

import (
	"log"
	"time"
	"crypto/sha1"
	"bytes"
	"fmt"
	"runtime"
	"github.com/rhydberg/gotorrent/client"
	"github.com/rhydberg/gotorrent/message"
	"github.com/rhydberg/gotorrent/peers"
)

// MaxBlockSize is the largest number of bytes a request can ask for
const MaxBlockSize = 16384

// MaxBacklog is the number of unfulfilled requests a client can have in its pipeline
const MaxBacklog = 5

type pieceWork struct {
	index int
	hash [20]byte
	length int
}

type pieceResult struct{
	index int
	buf []byte
}

type pieceProgress struct {
	index int
	client *client.Client
	buf []byte
	downloaded int
	requested int
	backlog int
}

func(state *pieceProgress) readMessage() error{
	msg,err:= message.Read(state.client.Conn)

	if err!=nil{
		return err
	}

	if msg == nil{
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		state.client.Choked = false
	case message.MsgChoke:
		state.client.Choked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.client.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n, err := message.ParsePiece(state.index, state.buf, msg)
		if err != nil {
			return err
		}
		state.downloaded += n
		state.backlog--
	}
	return nil

}

func attemptDownloadPiece(c *client.Client, pw *pieceWork) ([]byte, error){
	state:= pieceProgress{
		index: pw.index,
		client: c,
		buf: make([]byte, pw.length),
	}
	c.Conn.SetDeadline((time.Now().Add(30* time.Second)))
	defer c.Conn.SetDeadline(time.Time{})

	for state.downloaded<pw.length{
		if !state.client.Choked{
			for state.backlog < MaxBacklog && state.requested < pw.length {
				blocksize:= MaxBlockSize

				if pw.length-state.requested < blocksize{
					blocksize = pw.length - state.requested
				}

				err := c.SendRequest(pw.index, state.requested, blocksize)
				if err!=nil{
					return nil, err
				}
				state.backlog++
				state.requested += blocksize
			}
		}

		err:= state.readMessage()
		if err!=nil{
			return nil, err
		}
	}
	return state.buf, nil
}

func checkIntegrity(pw *pieceWork, buf []byte) error {
	hash := sha1.Sum(buf)
	if !bytes.Equal(hash[:], pw.hash[:]) {
		return fmt.Errorf("Index %d failed integrity check", pw.index)
	}
	return nil
}

func (tf *TorrentFile) startDownloadWorker(peer peers.Peer, workQueue chan *pieceWork, results chan *pieceResult){
	client, err := client.New(peer, tf.InfoHash, tf.PeerID)
	fmt.Println("Making client for ", peer.IP)
	if err!=nil{
		log.Print("could not create client for",peer.IP, err)
		return
	}
	defer client.Conn.Close()

	log.Printf("successfully made client for %s\n", peer.IP)

	client.SendUnchoke()
	client.SendInterested()

	for pw := range workQueue {
		if !client.Bitfield.HasPiece(pw.index){
			workQueue <- pw
			continue
		}
		buf, err := attemptDownloadPiece(client, pw)
		
		if err!=nil{
			log.Printf("Piece #%d failed integrity check\n", pw.index)
			workQueue <- pw // Put piece back on the queue
			continue
		}

		err = checkIntegrity(pw, buf)
		if err != nil {
			log.Printf("Piece #%d failed integrity check\n", pw.index)
			workQueue <- pw // Put piece back on the queue
			continue
		}

		client.SendHave(pw.index)
		results <- &pieceResult{pw.index, buf}
	}


}

func (tf *TorrentFile) Download() ([]byte, error){
	workQueue := make(chan *pieceWork, len(tf.PieceHashes))
	results:= make(chan *pieceResult)

	for index, hash := range tf.PieceHashes{
		begin := index*tf.Info.PieceLength
		end := (index+1)*tf.Info.PieceLength;
		if end > tf.Info.Length{
			end = tf.Info.Length
		}

		length := end - begin
		workQueue <- &pieceWork{index, hash, length}
	}

	peers, err:= tf.RequestPeers()
	// fmt.Printf("%v",peers)

	if err!=nil{
		log.Fatal("could not get peers: ",err)
	}

	for _, peer := range peers {
		go tf.startDownloadWorker(peer, workQueue, results)

	}

	buf:=make([]byte, tf.Info.Length)
	donePieces:=0

	for donePieces<len(tf.PieceHashes){
		result := <-results
		begin := result.index*tf.Info.PieceLength
		end := (result.index+1)*tf.Info.PieceLength;
		if end > tf.Info.Length{
			end = tf.Info.Length
		}

		copy(buf[begin:end], result.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(tf.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		log.Printf("(%0.2f%%) piece #%d from %d peers\n", percent, result.index, numWorkers)


	}
	close(workQueue)
	return buf, nil

}