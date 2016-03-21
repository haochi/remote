package remote

import "encoding/binary"

const idStartIndex = 0
const idSize = 2
const idEndIndex = idStartIndex + idSize
const messageLengthStartIndex = idEndIndex
const messageLengthSize = 4
const messageLengthEndIndex = messageLengthStartIndex + messageLengthSize
const messageStartIndex = messageLengthEndIndex

type Packet struct {
	// [id:2B][message length:4B][message:...]
	buff []byte
}

func (this *Packet) Serialize() []byte {
	return this.buff
}

func (this *Packet) getId() int {
	return int(binary.BigEndian.Uint16(this.buff[idStartIndex:idEndIndex]))
}

func (this *Packet) getLength() int {
	return int(binary.BigEndian.Uint32(this.buff[messageLengthStartIndex:messageLengthEndIndex]))
}

func (this *Packet) getBody() []byte {
	return this.buff[messageStartIndex:]
}

func newPacketWithSize(buff []byte) *Packet {
	p := &Packet{}
	p.buff = buff
	return p
}

func newPacket(id int, buff []byte) *Packet {
	p := &Packet{}
	p.buff = make([]byte, idSize+messageLengthSize+len(buff))
	binary.BigEndian.PutUint16(p.buff[idStartIndex:idEndIndex], uint16(id))
	binary.BigEndian.PutUint32(p.buff[messageLengthStartIndex:messageLengthEndIndex], uint32(len(buff)))
	copy(p.buff[messageStartIndex:], buff)
	return p
}
