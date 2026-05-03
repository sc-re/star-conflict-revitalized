package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

func ReadCString(reader *bytes.Reader) (string, error) {
	buf := new(strings.Builder)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0x00 {
			break
		}
		if err := buf.WriteByte(b); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

func ReadUint8(reader *bytes.Reader) (uint8, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

func ReadUint32(reader *bytes.Reader) (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf), nil
}
