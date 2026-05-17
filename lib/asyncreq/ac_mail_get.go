package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_mail_get_req struct{}

func (req *ac_mail_get_req) parse(body []byte) error {
	return nil
}

func (req *ac_mail_get_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	bw.WriteBeUint16(uint16(types.AC_MAIL_GET))
	bw.BwWriteByte(0)   // status
	bw.WriteBool(false) // keep existing mails (if this message does replace all mails or appends)
	bw.WriteBeUint16(0) // mail count
	return bw.ReturnSlice()
}

func handle_ac_mail_get(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_mail_get_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_mail_get_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
