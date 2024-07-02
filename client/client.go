package client

import (
	"bytes"
	"fmt"
	"net"
	"time"
	"github.com/rhydberg/gotorrent/bitfield"
	"github.com/rhydberg/gotorrent/handshake"
	"github.com/rhydberg/gotorrent/message"
	"github.com/rhydberg/gotorrent/peers"
)


type Client struct{
	Conn net.Conn
	Choked bool
	Bitfield bitfield.Bitfield
	peer peers.Peer
	infoHash [20]byte
	peerID [20]byte	 
}

func New(peer peers.Peer, infoHash, peerID [20]byte) (*Client, error){
	conn, err:= net.DialTimeout("tcp", peer.String(), 3*time.Second)
	if err!=nil{
		return nil, err
	}

	_, err = doHandshake(conn, infoHash, peerID)
	if err!=nil{
		conn.Close()
		return nil, err
	}

	bf, err := receiveBitfield(conn) //'bitfield' is only ever sent as the first message.
	if err!=nil{
		conn.Close()
		return nil, err
	}

	client:= Client{
		Conn: conn,
		Choked: true,
		Bitfield: bf,
		peer: peer,
		infoHash: infoHash,
		peerID: peerID,

	}
	return &client, nil
}

func receiveBitfield(conn net.Conn) (bitfield.Bitfield, error){
	conn.SetDeadline(time.Now().Add(5*time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err:= message.Read(conn)
	if err!=nil{
		return nil, err
	}
	if msg == nil {
		err := fmt.Errorf("expected bitfield but got %s", msg)
		return nil, err
	}
	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("expected bitfield but got ID %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil

}

func doHandshake(conn net.Conn, infoHash, peerID [20]byte)(*handshake.Handshake, error){
	conn.SetDeadline(time.Now().Add(3*time.Second))
	defer conn.SetDeadline(time.Time{})

	request := handshake.New(infoHash, peerID)
	serializedRequest:=request.Serialize()
	_, err := conn.Write(serializedRequest)
	if err!=nil{
		return nil, err
	}

	response, err:= handshake.ReadHandshake(conn)
	if err!=nil{
		return nil ,err
	}
	if !bytes.Equal(response.InfoHash[:], infoHash[:]){
		return nil, fmt.Errorf("Expected infohash %x but got %x", response.InfoHash, infoHash)
	}

	return response, nil
}

