package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_battle_slots_req struct{}

func (req *ac_battle_slots_req) parse(body []byte) error {
	return nil
}

func (req *ac_battle_slots_req) response() []byte {
	ret := make([]byte, 0, 46)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_BATTLE_SLOTS))
	ret = binary.BigEndian.AppendUint32(ret, 4) // Slot count
	for range 4 {
		ret = binary.BigEndian.AppendUint64(ret, 0) // Vessel ID
	}
	ret = binary.BigEndian.AppendUint64(ret, 0x0) // idk
	ret = binary.BigEndian.AppendUint64(ret, 0xc)
	return ret
}

func Send_ac_battle_slots(conn net.Conn) {
	req := ac_battle_slots_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_battle_slots(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_battle_slots_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_battle_slots_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
