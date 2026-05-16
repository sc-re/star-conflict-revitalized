package bitwriter

import (
	"math"
	"slices"
)

type Writer struct {
	buf    []byte
	offset int64
}

func NewWriter(buf []byte) *Writer {
	return &Writer{
		buf: buf,
	}
}

func (bw *Writer) WriteBool(v bool) {
	byteOffset := bw.offset / 8
	bitOffset := bw.offset % 8
	if bitOffset == 0 {
		bw.buf = append(bw.buf, 0)
	}
	if v {
		b := byte(1) << (7 - bitOffset)
		bw.buf[byteOffset] |= b
	}
	bw.offset += 1
}

func (bw *Writer) BwWriteByte(b byte) {
	byteOffset := bw.offset / 8
	bitOffset := bw.offset % 8
	if cap(bw.buf) == len(bw.buf) {
		bw.buf = slices.Grow(bw.buf, len(bw.buf)+1)
	}
	if bitOffset == 0 {
		bw.buf = append(bw.buf, b)
	} else {
		bw.buf[byteOffset] |= b >> byte(bitOffset)
		bw.buf = append(bw.buf, b<<byte(8-bitOffset))
	}
	bw.offset += 8
}

func (bw *Writer) WriteBeUint64(u uint64) {
	bw.BwWriteByte(byte(u >> 0x38))
	bw.BwWriteByte(byte(u >> 0x30))
	bw.BwWriteByte(byte(u >> 0x28))
	bw.BwWriteByte(byte(u >> 0x20))
	bw.BwWriteByte(byte(u >> 0x18))
	bw.BwWriteByte(byte(u >> 0x10))
	bw.BwWriteByte(byte(u >> 0x8))
	bw.BwWriteByte(byte(u))
}

func (bw *Writer) WriteBeInt32(u int32) {
	bw.BwWriteByte(byte(u >> 0x18))
	bw.BwWriteByte(byte(u >> 0x10))
	bw.BwWriteByte(byte(u >> 0x8))
	bw.BwWriteByte(byte(u))
}

func (bw *Writer) WriteBeUint32(u uint32) {
	bw.BwWriteByte(byte(u >> 0x18))
	bw.BwWriteByte(byte(u >> 0x10))
	bw.BwWriteByte(byte(u >> 0x8))
	bw.BwWriteByte(byte(u))
}

func (bw *Writer) WriteBeUint16(u uint16) {
	bw.BwWriteByte(byte(u >> 8))
	bw.BwWriteByte(byte(u))
}

func (bw *Writer) WriteFloat32(f float32) {
	bw.WriteBeUint32(math.Float32bits(f))
}

func (bw *Writer) WriteString(s string) {
	for _, b := range []byte(s) {
		bw.BwWriteByte(b)
	}
}

func (bw *Writer) WriteCString(s string) {
	bw.WriteString(s)
	bw.BwWriteByte(0x0)
}

func (bw *Writer) ReturnSlice() []byte {
	return bw.buf
}
