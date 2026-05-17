package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_lobby_info_req struct{}

func (req *ac_lobby_info_req) parse(body []byte) error {
	return nil
}

func (req *ac_lobby_info_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	bw.WriteBeUint16(uint16(types.AC_LOBBY_INFO))
	bw.WriteBeUint64(0) // Lobby ID
	bw.WriteCString("") // Name
	bw.WriteBeUint32(0) // Idk
	bw.WriteCString("") // level_def
	bw.BwWriteByte(0)   // a
	bw.BwWriteByte(0)   // botcount
	bw.BwWriteByte(0)   // c
	bw.WriteBool(false) // friendly fire
	bw.WriteBool(false) // self fire
	bw.WriteBool(false) // idk
	bw.WriteBool(false)
	bw.WriteBool(false)
	bw.WriteBool(false)
	bw.WriteFloat32(1.0) // idk
	bw.WriteFloat32(1.0) // idk
	bw.WriteBeUint16(0)  // Ship Role Mask
	bw.WriteBeUint32(0)  // Ship Rank Mask
	bw.WriteBool(true)   // Autobalance
	bw.WriteBeUint64(0)  // idk
	bw.BwWriteByte(0)    // player count
	bw.WriteCString("")  // bot preset
	bw.WriteCString("")  // idk
	bw.WriteBool(false)  // esports
	bw.WriteBool(false)  // idk
	bw.WriteCString("")  // team1 dreadnought
	bw.WriteCString("")  // team2 dreadnought

	return bw.ReturnSlice()
}

func handle_ac_lobby_info(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_lobby_info_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_lobby_info_req", "error", err)
	}
	resp := req.response()
	slog.Debug("AC_LOBBY_INFO", "resp", resp)
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
