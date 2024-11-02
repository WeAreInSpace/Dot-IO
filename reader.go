package dotio

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

/*
Dot I/O Packet format
packet data length + packet id + packet data
*/

type Inbound struct {
	Conn net.Conn
}

func (ib Inbound) Read() (id int32, ibb InboundBuffer, err error) {
	tempLength := make([]byte, 4)
	_, readLenE := ib.Conn.Read(tempLength)
	if readLenE != nil {
		log.Printf("ERROR: %s", readLenE)
		return 0, InboundBuffer{buffer: new(bytes.Buffer)}, readLenE
	}

	length := ReadInt32(tempLength)

	tempPacket := make([]byte, length)

	_, readAllPacketE := ib.Conn.Read(tempPacket)
	if readAllPacketE != nil {
		log.Printf("ERROR: %s", readLenE)
		return 0, InboundBuffer{buffer: new(bytes.Buffer)}, readLenE
	}

	buffer := InboundBuffer{
		buffer: new(bytes.Buffer),
	}
	buffer.buffer.Write(tempPacket)

	id = buffer.ReadInt32()

	return id, buffer, nil
}

type InboundBuffer struct {
	buffer *bytes.Buffer
}

// 4
func (ib *InboundBuffer) ReadInt32() int32 {
	tempPacket := new(bytes.Buffer)

	var i int8
	for {
		if i == 4 {
			break
		}
		data, _ := ib.buffer.ReadByte()
		tempPacket.WriteByte(data)
		i += 1
	}

	var number int32
	binary.Read(tempPacket, binary.BigEndian, &number)
	return number
}

// 4
func (ib *InboundBuffer) ReadFloat32() float32 {
	tempPacket := new(bytes.Buffer)

	var i int8
	for {
		if i == 4 {
			break
		}
		data, _ := ib.buffer.ReadByte()
		tempPacket.WriteByte(data)
		i += 1
	}

	var number float32
	binary.Read(tempPacket, binary.BigEndian, &number)
	return number
}

func ReadInt32(rawNumber []byte) int32 {
	tempByte := new(bytes.Buffer)
	tempByte.Write(rawNumber)

	var number int32
	binary.Read(tempByte, binary.BigEndian, &number)
	return number
}

// 8
func (ib *InboundBuffer) ReadInt64() int64 {
	tempPacket := new(bytes.Buffer)

	var i int8
	for {
		if i == 8 {
			break
		}
		data, _ := ib.buffer.ReadByte()
		tempPacket.WriteByte(data)
		i += 1
	}

	var number int64
	binary.Read(tempPacket, binary.BigEndian, &number)
	return number
}

// 4
func (ib *InboundBuffer) ReadFloat64() float64 {
	tempPacket := new(bytes.Buffer)

	var i int8
	for {
		if i == 8 {
			break
		}
		data, _ := ib.buffer.ReadByte()
		tempPacket.WriteByte(data)
		i += 1
	}

	var number float64
	binary.Read(tempPacket, binary.BigEndian, &number)
	return number
}

// length
func (ib *InboundBuffer) ReadString() string {
	length := ib.ReadInt32()

	tempPacket := new(bytes.Buffer)
	var i int32
	for {
		if i == int32(length) {
			break
		}
		data, _ := ib.buffer.ReadByte()
		tempPacket.WriteByte(data)
		i += 1
	}

	return tempPacket.String()
}

// 1
func (ib *InboundBuffer) ReadBoolean() bool {
	tempPacket, _ := ib.buffer.ReadByte()

	if tempPacket == 0 {
		return false
	} else {
		return true
	}
}
