package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_adventures_req struct{}

func (req *ac_adventures_req) parse(body []byte) error {
	return nil
}

func (req *ac_adventures_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 4))
	bw.WriteBeUint16(uint16(types.AC_ADVENTURES))
	bw.WriteBeUint16(0) // Adventure Count
	return bw.ReturnSlice()
}

func handle_ac_adventures(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_adventures_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_adventures_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
