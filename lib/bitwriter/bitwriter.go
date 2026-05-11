package bitwriter

import (
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

func (bw *Writer) WriteBool(v bool) error {
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
	return nil
}

func (bw *Writer) WriteByte(b byte) error {
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
	return nil
}

func (bw *Writer) WriteBeUint64(u uint64) error {
	bw.WriteByte(byte(u >> 0x38))
	bw.WriteByte(byte(u >> 0x30))
	bw.WriteByte(byte(u >> 0x28))
	bw.WriteByte(byte(u >> 0x20))
	bw.WriteByte(byte(u >> 0x18))
	bw.WriteByte(byte(u >> 0x10))
	bw.WriteByte(byte(u >> 0x8))
	bw.WriteByte(byte(u))
	return nil
}

func (bw *Writer) WriteBeInt32(u int32) error {
	bw.WriteByte(byte(u >> 0x18))
	bw.WriteByte(byte(u >> 0x10))
	bw.WriteByte(byte(u >> 0x8))
	bw.WriteByte(byte(u))
	return nil
}

func (bw *Writer) WriteBeUint32(u uint32) error {
	bw.WriteByte(byte(u >> 0x18))
	bw.WriteByte(byte(u >> 0x10))
	bw.WriteByte(byte(u >> 0x8))
	bw.WriteByte(byte(u))
	return nil
}

func (bw *Writer) WriteBeUint16(u uint16) error {
	bw.WriteByte(byte(u >> 8))
	bw.WriteByte(byte(u))
	return nil
}

func (bw *Writer) WriteString(s string) error {
	for _, b := range []byte(s) {
		if err := bw.WriteByte(b); err != nil {
			return err
		}
	}
	return nil
}

func (bw *Writer) WriteCString(s string) error {
	if err := bw.WriteString(s); err != nil {
		return err
	}
	if err := bw.WriteByte(0x0); err != nil {
		return err
	}
	return nil
}

func (bw *Writer) ReturnSlice() []byte {
	return bw.buf
}
