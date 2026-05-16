package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"starconflict/lib/types"
)

type Header struct {
	Length         uint32
	Sequence       uint16
	ReturnSequence uint16
	CommandType    types.MessageType
	Checksum       uint16
	Special        bool
}

func MakeMessage(cmdType types.MessageType, seq uint16, seqRet uint16, body []byte) []byte {
	hdr := make([]byte, 12)
	binary.BigEndian.PutUint32(hdr[0:], uint32(len(body)))
	binary.BigEndian.PutUint16(hdr[4:], seq)
	binary.BigEndian.PutUint16(hdr[6:], seqRet)
	binary.BigEndian.PutUint16(hdr[8:], uint16(cmdType))
	binary.BigEndian.PutUint16(hdr[10:], TgpChecksum(uint32(len(body)), seq, seqRet, uint16(cmdType), body))
	return append(hdr, body...)
}

func ParseNextMessage(conn net.Conn) (*Header, []byte, error) {
	hdrb := make([]byte, 12)
	hdr := Header{}
	var body []byte
	if _, err := io.ReadFull(conn, hdrb); err != nil {
		return nil, nil, fmt.Errorf("read hdr: %w", err)
	}

	hdr.Length = binary.BigEndian.Uint32(hdrb[0:])
	if hdr.Length > 0xfffffc {
		// Handle Keep Alive BS
		hdr.Special = true
		return &hdr, body, nil
	}

	hdr.Special = false
	hdr.Sequence = binary.BigEndian.Uint16(hdrb[4:])
	hdr.ReturnSequence = binary.BigEndian.Uint16(hdrb[6:])
	hdr.CommandType = types.MessageType(binary.BigEndian.Uint16(hdrb[8:]))
	hdr.Checksum = binary.BigEndian.Uint16(hdrb[10:])
	body = make([]byte, hdr.Length)

	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, nil, fmt.Errorf("read body: %w", err)
	}
	return &hdr, body, nil
}
