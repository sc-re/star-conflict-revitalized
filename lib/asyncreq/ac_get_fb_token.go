package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_get_fb_token_req struct{}

func (req *ac_get_fb_token_req) parse(body []byte) error {
	return nil
}

// We don't care about fb tokens, the game just excpects this message
func (req *ac_get_fb_token_req) response() []byte {
	ret := make([]byte, 20)
	binary.BigEndian.PutUint16(ret[0:], uint16(types.AC_GET_FB_TOKEN))
	return ret
}

func Send_ac_get_fb_token(conn net.Conn) {
	req := ac_get_fb_token_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_get_fb_token(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_get_fb_token_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_get_fb_token_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
