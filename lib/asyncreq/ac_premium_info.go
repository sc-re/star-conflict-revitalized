package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_premium_info_req struct{}

func (req *ac_premium_info_req) parse(body []byte) error {
	return nil
}

func (req *ac_premium_info_req) response() []byte {
	ret := make([]byte, 0, 10)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_PREMIUM_INFO))
	ret = binary.BigEndian.AppendUint64(ret, 1767806947605) // Timestamp for premium
	return ret
}

func Send_ac_premium_info(conn net.Conn) {
	req := ac_premium_info_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_premium_info(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_premium_info_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_premium_info_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
