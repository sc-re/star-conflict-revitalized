package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
)

type ac_player_autogen_inventory_req struct{}

func (req *ac_player_autogen_inventory_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_autogen_inventory_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 3000))
	bw.WriteBeUint16(uint16(types.AC_PLAYER_AUTOGEN_INVENTORY))
	bw.WriteBeUint32(0)  // Item count
	bw.BwWriteByte(1)    // Inventory level
	bw.WriteBeUint32(40) // Inventory size
	return bw.ReturnSlice()
}

func Send_ac_player_autogen_inventory(session *session.Session) {
	req := ac_player_autogen_inventory_req{}
	resp := req.response()
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_autogen_inventory(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_player_autogen_inventory_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_autogen_inventory_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
