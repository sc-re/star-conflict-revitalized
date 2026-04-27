package protocol

import (
	"encoding/binary"
)

func MakePacket(cmdType uint16, seq uint16, seqRet uint16, body []byte) []byte {
	hdr := make([]byte, 12)
	binary.BigEndian.PutUint32(hdr[0:], uint32(len(body)))
	binary.BigEndian.PutUint16(hdr[4:], seq)
	binary.BigEndian.PutUint16(hdr[6:], seqRet)
	binary.BigEndian.PutUint16(hdr[8:], cmdType)
	binary.BigEndian.PutUint16(hdr[10:], TgpChecksum(uint32(len(body)), seq, seqRet, cmdType, body))
	return append(hdr, body...)
}
