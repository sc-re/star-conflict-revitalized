package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_clan_request_desc_req struct{}

func (req *ac_clan_request_desc_req) parse(body []byte) error {
	return nil
}

func (req *ac_clan_request_desc_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	bw.WriteBeUint16(uint16(types.AC_CLAN_REQUEST_DESC))
	bw.WriteBeUint64(0)
	bw.WriteCString("") // Name
	bw.WriteCString("") // tag
	bw.WriteCString("") // motd
	bw.WriteCString("") // desc
	bw.WriteCString("") // emblem
	bw.WriteCString("") // current clan ship
	bw.WriteBeUint64(0) // creation date
	bw.WriteBeUint32(1148846080)
	bw.WriteBeUint32(0) // Counter Target
	bw.WriteBeUint32(0) // Counter Progress
	bw.WriteBeInt32(-1) // Clan Quest ID
	bw.WriteBeUint16(0) // Clan Quest Progress
	bw.BwWriteByte(0)   // Recruiting
	bw.WriteBeUint32(0) // member count
	bw.WriteBeUint32(0) // Invite count
	bw.WriteBeUint32(0) // Join Request Count
	bw.BwWriteByte(0)   // upgrade_a
	bw.BwWriteByte(0)   // upgrade_b
	// resources
	bw.WriteBeUint64(0)
	bw.WriteBeUint64(0)
	bw.WriteBeUint32(0) // StreamFlags
	// Should be a variantdict
	bw.WriteBeUint32(0)
	bw.BwWriteByte(0)
	return bw.ReturnSlice()
}

func handle_ac_clan_request_desc(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_clan_request_desc_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_clan_request_desc_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
