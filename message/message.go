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

func Read(r io.Reader) (*Message, error){
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


//Fixed length, used to request a block of pieces. The payload contains integer values specifying the index, begin location and length.
func CreateRequestMessage(index, begin, length int) *Message{
	requestLength:=13
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	message:=Message{
		Length: uint32(requestLength),
		ID: MsgRequest,
		Payload: payload,
	}

	return &message
}

//The 'have' message's payload is a single number, the index which that downloader just completed and checked the hash of.
func CreateHaveMessage(index int) *Message{
	haveLength := 5
	payload := make ([]byte, 4)
	binary.BigEndian.PutUint32(payload[:], uint32(index))

	message:=Message{
		Length: uint32(haveLength),
		ID: MsgHave,
		Payload: payload,
	}

	return &message
}





