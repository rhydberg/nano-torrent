package handshake

import (
	"fmt"
	"io"
)

// import "github.com/rhydberg/gotorrent/torrentfile"

type Handshake struct {
	Length        byte
	Pstr          string
	ReservedBytes [8]byte
	InfoHash      [20]byte
	PeerID        [20]byte
}

func New(infohash, peerID [20]byte) *Handshake {
	pstr := "BitTorrent protocol"
	handshake := Handshake{
		Length:        byte(len(pstr)),
		Pstr:          pstr,
		ReservedBytes: [8]byte{},
		InfoHash:      infohash,
		PeerID:        peerID,
	}

	return &handshake
}

func (h *Handshake) Serialize() []byte {
	var buf []byte
	buf = append(buf, h.Length)
	buf = append(buf, []byte(h.Pstr)...)
	buf = append(buf, h.ReservedBytes[:]...)
	buf = append(buf, h.InfoHash[:]...)
	buf = append(buf, h.PeerID[:]...)

	return buf
}

func ReadHandshake(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	pStrLen := int(lengthBuf[0])
	if pStrLen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	pStrBuf := make([]byte, pStrLen)
	_, err = io.ReadFull(r, pStrBuf)
	if err != nil {
		return nil, err
	}

	reservedBytesBuf := make([]byte, 8)
	_, err = io.ReadFull(r, reservedBytesBuf)
	if err != nil {
		return nil, err
	}

	infoHashBuf := make([]byte, 20)
	_, err = io.ReadFull(r, infoHashBuf)
	if err != nil {
		return nil, err
	}

	peerIDBuf := make([]byte, 20)
	_, err = io.ReadFull(r, peerIDBuf)
	if err != nil {
		return nil, err
	}

	handshake := Handshake{
		Length:        lengthBuf[0],
		Pstr:          string(pStrBuf),
		ReservedBytes: [8]byte(reservedBytesBuf),
		InfoHash:      [20]byte(infoHashBuf),
		PeerID:        [20]byte(peerIDBuf),
	}

	return &handshake, nil
}
