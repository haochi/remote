package remote

import (
	"encoding/binary"
	"github.com/gansidui/gotcp"
	"io"
	"net"
)

type Protocol struct {
}

func (this *Protocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	var (
		idBytes     []byte = make([]byte, idSize)
		lengthBytes []byte = make([]byte, messageLengthSize)
		length      int
	)

	// read id
	if _, err := io.ReadFull(conn, idBytes); err != nil {
		return nil, err
	}

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}
	length = int(binary.BigEndian.Uint32(lengthBytes[:messageLengthSize]))

	buff := make([]byte, idSize+messageLengthSize+length)

	// set id
	copy(buff[idStartIndex:idEndIndex], idBytes)

	// set length
	copy(buff[messageLengthStartIndex:messageLengthEndIndex], lengthBytes)

	// copy body
	if _, err := io.ReadFull(conn, buff[messageStartIndex:]); err != nil {
		return nil, err
	}

	return newPacketWithSize(buff), nil
}
