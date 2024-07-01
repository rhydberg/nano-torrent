package handshake

import "io"

// import "github.com/rhydberg/gotorrent/torrentfile"

type Handshake struct {
	Length   byte
	Pstr 	string
	ReservedBytes [8]byte
	InfoHash [20]byte
	PeerID [20]byte
}

func New(infohash, peerID [20]byte) *Handshake{
	pstr:="BitTorrent protocol"
	handshake := Handshake{
		Length: byte(len(pstr)),
		Pstr: pstr,
		ReservedBytes: [8]byte{},
		InfoHash: infohash,
		PeerID: peerID,
	}

	return &handshake
}

func (h *Handshake) Serialize() []byte{
	var buf []byte
	buf = append(buf, h.Length)
	buf = append(buf, []byte(h.Pstr)...)
	buf = append(buf, h.ReservedBytes[:]...)
	buf = append(buf, h.InfoHash[:]...)
	buf = append(buf, h.PeerID[:]...)

	return buf
}

func ReadHandshake(r io.Reader) (*Handshake, error){
	
}