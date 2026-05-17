package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_friends_send_request_req struct{}

func (req *ac_friends_send_request_req) parse(body []byte) error {
	return nil
}

func (req *ac_friends_send_request_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	bw.WriteBeUint16(uint16(types.AC_FRIENDS_SEND_REQUEST))
	bw.BwWriteByte(0) // Friends Count
	bw.BwWriteByte(0) // Incoming Friends Requests Count
	bw.BwWriteByte(0) // Outoing Friends Requests Count
	bw.BwWriteByte(0) // Ignored Count
	bw.BwWriteByte(0) // Watchlist Count
	bw.BwWriteByte(0) // idk count
	bw.BwWriteByte(0) // idk2 count
	return bw.ReturnSlice()
}

func handle_ac_friends_send_request(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_friends_send_request_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_friends_send_request_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
