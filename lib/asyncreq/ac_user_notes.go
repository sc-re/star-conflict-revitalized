package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_user_notes_req struct{}

func (req *ac_user_notes_req) parse(body []byte) error {
	return nil
}

func (req *ac_user_notes_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 6))
	bw.WriteBeUint16(uint16(types.AC_USER_NOTES))
	bw.WriteBeUint32(0) // No user notes
	return bw.ReturnSlice()
}

func handle_ac_user_notes(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_user_notes_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_user_notes_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
