package protocol

import (
	"encoding/binary"
)

func MakePacket(msgType uint16, seq uint32, body []byte) []byte {
	hdr := make([]byte, 12)
	binary.BigEndian.PutUint32(hdr[0:], uint32(len(body)))
	binary.BigEndian.PutUint16(hdr[4:], msgType)
	binary.BigEndian.PutUint32(hdr[6:], seq)
	binary.BigEndian.PutUint16(hdr[10:], TgpChecksum(uint32(len(body)), msgType, seq, body))
	return append(hdr, body...)
}
