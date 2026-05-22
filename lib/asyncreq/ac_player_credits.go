package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_player_credits_req struct {
	flag uint16
}

func (req *ac_player_credits_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_credits_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 68))
	bw.WriteBeUint16(uint16(types.AC_PLAYER_CREDITS))
	bw.WriteBeUint16(req.flag)
	bw.WriteBeUint64(0) // credits
	if req.flag&0x02 != 0 {
		bw.WriteBeUint64(1234) // gold
	}
	if req.flag&0x04 != 0 {
		bw.WriteBeUint64(3456) // Iridium
	}
	if req.flag&0x08 != 0 {
		bw.WriteBeUint64(4) // Xenochips
	}
	if req.flag&0x10 != 0 {
		bw.WriteBeUint64(1778418304383) // Premium timestamp?
	}
	if req.flag&0x20 != 0 {
		bw.WriteBeUint64(0) // idk, vid
	}
	if req.flag&0x40 != 0 {
		bw.WriteBeUint32(573123) // Free Synergy
	}
	// idk, some ressources?
	if req.flag&0x80 != 0 {
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
	}
	return bw.ReturnSlice()
}

func Send_ac_player_credits(conn net.Conn) {
	req := ac_player_credits_req{flag: 222}

	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_credits(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_player_credits_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_credits_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
