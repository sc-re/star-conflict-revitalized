package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type ac_get_blueprints_req struct{}

func (req *ac_get_blueprints_req) parse(body []byte) error {
	return nil
}

func (req *ac_get_blueprints_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	bw.WriteBeUint16(uint16(types.AC_GET_BLUEPRINTS))
	bw.WriteBool(true)
	blueprints := map[string]int32{
		"BP_LaserPenetrate": 1,
	}
	variantdict.BwMarshal(bw, blueprints)
	return bw.ReturnSlice()
}

func Send_ac_get_blueprints(conn net.Conn) {
	req := ac_get_blueprints_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_get_blueprints(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_get_blueprints_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_get_blueprints_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
