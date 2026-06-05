package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"starconflict/lib/types"
)

type HubHeader struct {
	Length         uint32
	Sequence       uint16
	ReturnSequence uint16
	SrcId          uint64
	DstId          uint64
	CmdId          byte
	CommandType    types.HubCommandType
	Special        bool
}

func MakeHubMessage(header HubHeader, body []byte) []byte {
	hdr := make([]byte, 0, 27)
	hdr = binary.BigEndian.AppendUint32(hdr, uint32(len(body)))
	hdr = binary.BigEndian.AppendUint16(hdr, header.Sequence)
	hdr = binary.BigEndian.AppendUint16(hdr, header.ReturnSequence)
	hdr = binary.BigEndian.AppendUint64(hdr, header.SrcId)
	hdr = binary.BigEndian.AppendUint64(hdr, header.DstId)
	hdr = append(hdr, header.CmdId)
	hdr = binary.BigEndian.AppendUint16(hdr, uint16(header.CommandType))
	return append(hdr, body...)
}

func ParseNextHubMessage(conn net.Conn) (*HubHeader, []byte, error) {
	hdrb := make([]byte, 27)
	hdr := HubHeader{}
	var body []byte
	if _, err := io.ReadFull(conn, hdrb); err != nil {
		return nil, nil, fmt.Errorf("read hub hdr: %w", err)
	}
	hdr.Length = binary.BigEndian.Uint32(hdrb[0:])
	if hdr.Length > 0xfffffc {
		hdr.Special = true
		return &hdr, body, nil
	}
	hdr.Special = false
	hdr.Sequence = binary.BigEndian.Uint16(hdrb[4:])
	hdr.ReturnSequence = binary.BigEndian.Uint16(hdrb[6:])
	hdr.SrcId = binary.BigEndian.Uint64(hdrb[8:])
	hdr.DstId = binary.BigEndian.Uint64(hdrb[16:])
	hdr.CmdId = hdrb[24]
	hdr.CommandType = types.HubCommandType(binary.BigEndian.Uint16(hdrb[25:]))
	body = make([]byte, hdr.Length)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, nil, fmt.Errorf("read hub body: %w", err)
	}
	return &hdr, body, nil
}
