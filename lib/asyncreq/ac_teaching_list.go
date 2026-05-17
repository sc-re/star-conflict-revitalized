package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_teaching_list_req struct{}

func (req *ac_teaching_list_req) parse(body []byte) error {
	return nil
}

func (req *ac_teaching_list_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	bw.WriteBeUint16(uint16(types.AC_TEACHING_LIST))
	for range 6 {
		bw.WriteBeUint32(0)
	}
	bw.WriteBool(true)
	bw.WriteBool(true)
	return bw.ReturnSlice()
}

func handle_ac_teaching_list(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_teaching_list_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_teaching_list_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
