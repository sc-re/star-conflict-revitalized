package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
)

type ac_clan_request_profile_req struct {
	uid uint64
}

func (req *ac_clan_request_profile_req) parse(body []byte) error {
	req.uid = binary.BigEndian.Uint64(body)
	return nil
}

func (req *ac_clan_request_profile_req) response() []byte {
	ret := make([]byte, 0, 18)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_CLAN_REQUEST_PROFILE))
	ret = binary.BigEndian.AppendUint64(ret, req.uid)
	ret = binary.BigEndian.AppendUint64(ret, 0) // XXX: clan id
	return ret
}

func handle_ac_clan_request_profile(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_clan_request_profile_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_clan_request_profile_req", "error", err)
	}
	resp := req.response()
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
