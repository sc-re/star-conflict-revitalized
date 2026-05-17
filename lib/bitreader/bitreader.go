package bitreader

import (
	"fmt"
	"math"
)

type Reader struct {
	buf    []byte
	offset int64
}

func NewReader(buf []byte) *Reader {
	return &Reader{
		buf: buf,
	}
}

func (bw *Reader) ReadBool() (bool, error) {
	byteOffset := bw.offset / 8
	bitOffset := bw.offset % 8
	if len(bw.buf) <= int(byteOffset) {
		return false, fmt.Errorf("EofByte, tried to read from []byte len(%v) at offset %v", len(bw.buf), byteOffset)
	}
	b := bw.buf[byteOffset] >> (7 - bitOffset)
	bw.offset += 1
	return b == 1, nil
}

func (bw *Reader) ReadByte() (byte, error) {
	var ret byte
	byteOffset := bw.offset / 8
	bitOffset := bw.offset % 8
	if bitOffset == 0 {
		if len(bw.buf) <= int(byteOffset) {
			return 0, fmt.Errorf("EofByte, tried to read from []byte len(%v) at offset %v", len(bw.buf), byteOffset)
		}
		ret = bw.buf[byteOffset]
	} else {
		if len(bw.buf)-1 <= int(byteOffset) {
			return 0, fmt.Errorf("EofByte, tried to read from []byte len(%v) at offset %v", len(bw.buf), byteOffset)
		}
		ret = bw.buf[byteOffset] << byte(bitOffset)
		ret |= bw.buf[byteOffset+1] >> byte(8-bitOffset)
	}
	bw.offset += 8
	return ret, nil
}

func (bw *Reader) ReadBeUint64() (uint64, error) {
	var ret uint64
	b, err := bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x38

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x30

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x28

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x20

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x18

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x10

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b) << 0x8

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint64(b)

	return ret, nil
}

func (bw *Reader) ReadBeInt32() (int32, error) {
	var ret int32
	b, err := bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= int32(b) << 0x18

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= int32(b) << 0x10

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= int32(b) << 0x8

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= int32(b)
	return ret, nil
}

func (bw *Reader) ReadBeUint32() (uint32, error) {
	var ret uint32
	b, err := bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint32(b) << 0x18

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint32(b) << 0x10

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint32(b) << 0x8

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint32(b)

	return ret, nil
}

func (bw *Reader) ReadBeUint16() (uint16, error) {
	var ret uint16
	b, err := bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint16(b) << 8

	b, err = bw.ReadByte()
	if err != nil {
		return 0, err
	}
	ret |= uint16(b)
	return ret, nil
}

func (bw *Reader) ReadFloat32() (float32, error) {
	i, err := bw.ReadBeUint32()
	if err != nil {
		return 0.0, err
	}
	f := math.Float32frombits(i)
	return f, nil
}

func (bw *Reader) ReadCString() (string, error) {
	var ret []byte
	for {
		b, err := bw.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0x00 {
			return string(ret), nil
		}
		ret = append(ret, b)
	}
}

func (bw *Reader) ReturnSlice() []byte {
	return bw.buf
}
