package message

import (
	"encoding/binary"
	"io"
)


type messageID uint8

const (
	// MsgChoke chokes the receiver
	MsgChoke messageID = 0
	// MsgUnchoke unchokes the receiver
	MsgUnchoke messageID = 1
	// MsgInterested expresses interest in receiving data
	MsgInterested messageID = 2
	// MsgNotInterested expresses disinterest in receiving data
	MsgNotInterested messageID = 3
	// MsgHave alerts the receiver that the sender has downloaded a piece
	MsgHave messageID = 4
	// MsgBitfield encodes which pieces that the sender has downloaded
	MsgBitfield messageID = 5
	// MsgRequest requests a block of data from the receiver
	MsgRequest messageID = 6
	// MsgPiece delivers a block of data to fulfill a request
	MsgPiece messageID = 7
	// MsgCancel cancels a request
	MsgCancel messageID = 8
)

type Message struct {
	Length uint32
	ID messageID
	Payload []byte
}

func (m *Message) Serialize() []byte {
	if m == nil{
		emptyMessage := make([]byte, 4)
		return emptyMessage 
	}

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf[:], m.Length)

	buf = append(buf, byte(m.ID))

	buf = append(buf, m.Payload...)

	return buf
}

func ReadMessage(r io.Reader) (*Message, error){
	lengthBuf := make([]byte, 4)
	_,err := io.ReadFull(r, lengthBuf)
	if err!=nil{
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	if length ==0 { // keep-alive
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_,err= io.ReadFull(r,messageBuf)
	if err!=nil{
		return nil,err
	}

	message:= Message{
		Length: length,
		ID: messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &message, nil
}

