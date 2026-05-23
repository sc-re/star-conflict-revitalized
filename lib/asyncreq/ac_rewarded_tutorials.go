package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_rewarded_tutorials_req struct{}

func (req *ac_rewarded_tutorials_req) parse(body []byte) error {
	return nil
}

func (req *ac_rewarded_tutorials_req) response() []byte {
	ret := make([]byte, 0, 6)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_REWARDED_TUTORIALS))
	ret = append(ret, 1) // tutorial count
	ret = append(ret, 0) // list of tuturial ids
	return ret
}

func Send_ac_rewarded_tutorials(conn net.Conn) {
	req := ac_rewarded_tutorials_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_rewarded_tutorials(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_rewarded_tutorials_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_reward_tutorials_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
