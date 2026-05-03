package asyncreq

import (
	"encoding/binary"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_universe_get struct{}

func (req *ac_universe_get) response() []byte {
	resp := make([]byte, 5)
	binary.BigEndian.PutUint16(resp[0:], uint16(types.AC_UNIVERSE_GET))
	resp[2] = 0x2
	resp[3] = 0x0
	resp[4] = 0x0 // array length
	return resp
}

func handle_ac_universe_get(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_universe_get{}
	if len(body) > 2 {
		// err
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
