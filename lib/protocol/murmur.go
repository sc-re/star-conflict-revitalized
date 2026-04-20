package protocol

import (
	"encoding/binary"
)

const (
	tgpSeed = uint32(0x1337533d)
	tgpM    = uint32(0x5bd1e995)
)

func TgpChecksum(bodyLen uint32, msgType uint16, seq uint32, body []byte) uint16 {
	hdr := make([]byte, 12)
	binary.LittleEndian.PutUint32(hdr[0:], bodyLen)
	binary.LittleEndian.PutUint16(hdr[4:], msgType)
	binary.LittleEndian.PutUint16(hdr[6:], uint16(seq>>16))
	binary.LittleEndian.PutUint16(hdr[8:], uint16(seq&0xffff))

	h := uint32(12) ^ tgpSeed
	for _, data := range [][]byte{hdr, body} {
		i := 0
		for i+4 <= len(data) {
			k := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
			k *= tgpM
			k ^= k >> 24
			k *= tgpM
			h = h*tgpM ^ k
			i += 4
		}
		rem := len(data) - i
		if rem >= 3 {
			h ^= uint32(data[i+2]) << 16
		}
		if rem >= 2 {
			h ^= uint32(data[i+1]) << 8
		}
		if rem >= 1 {
			h = (h ^ uint32(data[i])) * tgpM
		}
	}
	h ^= h >> 13
	h *= tgpM
	h ^= h >> 15
	return uint16(h)
}
