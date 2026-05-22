package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_player_vessels_req struct{}

func (req *ac_player_vessels_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_vessels_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 2000))
	bw.WriteBeUint16(uint16(types.AC_PLAYER_VESSELS))
	bw.WriteBeUint32(0)  // Vessel Count
	bw.WriteFloat32(0.0) // idk
	bw.WriteFloat32(0.0) // Fleet Strength
	return bw.ReturnSlice()
}

func Send_ac_player_vessels(conn net.Conn) {
	req := ac_player_vessels_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_vessels(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_player_vessels_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_vessels_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
