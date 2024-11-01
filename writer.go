package dotio

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

/*
Dot IO Packet format
packet data length + packet id + packet data
*/

/*
```` packet writer
*/

type Outgoing struct {
	Conn net.Conn
}

func (og *Outgoing) Write() OutgoingBuffer {
	return OutgoingBuffer{
		buffer: new(bytes.Buffer),
		packet: new(bytes.Buffer),
		conn:   og.Conn,
	}
}

type OutgoingBuffer struct {
	conn   net.Conn
	buffer *bytes.Buffer
	packet *bytes.Buffer
}

func (og *OutgoingBuffer) Sent(id []byte) error {
	og.packet.Write(id)
	for _, buffer := range og.buffer.Bytes() {
		og.packet.WriteByte(buffer)
	}

	packetLength := WriteInt32(int32(og.packet.Len()))

	_, writePacketLenE := og.conn.Write(packetLength)
	if writePacketLenE != nil {
		log.Printf("ERROR: %s\n", writePacketLenE)
		return writePacketLenE
	}

	_, writePacketDataE := og.conn.Write(og.packet.Bytes())
	if writePacketDataE != nil {
		log.Printf("ERROR: %s\n", writePacketDataE)
		return writePacketDataE
	}

	return nil
}

// 4
func (og *OutgoingBuffer) WriteInt32(number int32) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// 4
func WriteInt32(number int32) []byte {
	tempByte := new(bytes.Buffer)

	binary.Write(tempByte, binary.BigEndian, number)

	return tempByte.Bytes()
}

// 4
func (og *OutgoingBuffer) WriteFloat32(number float32) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// 8
func (og *OutgoingBuffer) WriteInt64(number int64) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// 8
func (og *OutgoingBuffer) WriteFloat64(number float64) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// length
func (og *OutgoingBuffer) WriteString(str string) {
	stringLength := WriteInt32(int32(len(str)))

	og.buffer.Write(stringLength)
	og.buffer.Write([]byte(str))
}

// 1
func (og *OutgoingBuffer) WriteBoolean(boolean bool) {
	if boolean {
		binary.Write(og.buffer, binary.BigEndian, int8(1))
	} else {
		binary.Write(og.buffer, binary.BigEndian, int8(0))
	}
}
