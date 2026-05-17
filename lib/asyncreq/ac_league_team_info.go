package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_league_team_info_req struct{}

func (req *ac_league_team_info_req) parse(body []byte) error {
	return nil
}

func (req *ac_league_team_info_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 56))
	bw.WriteBeUint16(uint16(types.AC_LEAGUE_TEAM_INFO))
	bw.BwWriteByte(0)   // status
	bw.WriteBeUint64(0) // team id
	bw.WriteCString("") // team name
	bw.WriteCString("") // short name
	bw.WriteBeUint64(0) // team captain uid
	bw.BwWriteByte(0)   // member count
	// write members
	bw.BwWriteByte(0) // some other count
	// write idk what
	bw.WriteFloat32(1000.0) // Team rating
	bw.WriteBeUint32(0)     // idk
	bw.WriteBeUint32(0)     // idk
	bw.WriteBeUint32(0)     // idk
	bw.WriteBeUint64(0)     // idk
	bw.WriteBeUint64(0)     // idk
	bw.WriteBool(false)     // idk
	return bw.ReturnSlice()
}

func handle_ac_league_team_info(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_league_team_info_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_league_team_info_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
