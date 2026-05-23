package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_obtain_referral_key_req struct{}

func (req *ac_obtain_referral_key_req) parse(body []byte) error {
	return nil
}

func (req *ac_obtain_referral_key_req) response() []byte {
	ret := make([]byte, 0, 36)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_OBTAIN_REFERRAL_KEY))
	ret = append(ret, []byte{0x9b, 0x99, 0x9c, 0x18, 0x30,
		0xb0, 0x9a, 0x98, 0x9a, 0x1b, 0x1c, 0x98, 0xb1,
		0x1a, 0x1a, 0xb3, 0x1a, 0x98, 0xb2, 0x1a, 0x32,
		0x9a, 0x19, 0x1b, 0x31, 0x98, 0x98, 0x1a, 0x19,
		0x1c, 0x31, 0x1c, 0x00, 0x00}...)
	return ret
}

func Send_ac_obtain_referal_key(conn net.Conn) {
	req := ac_obtain_referral_key_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_obtain_referral_key(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_obtain_referral_key_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_obtain_referral_key_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
