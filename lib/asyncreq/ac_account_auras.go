package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_account_auras_req struct{}

func (req *ac_account_auras_req) parse(body []byte) error {
	return nil
}

func (req *ac_account_auras_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 574))
	bw.WriteBeUint16(uint16(types.AC_ACCOUNT_AURAS))
	bw.WriteBool(true) // idk
	bw.BwWriteByte(0)  // Aura count
	return bw.ReturnSlice()
}

func Send_ac_account_auras(conn net.Conn) {
	req := ac_account_auras_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_account_auras(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_account_auras_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_account_auras_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
