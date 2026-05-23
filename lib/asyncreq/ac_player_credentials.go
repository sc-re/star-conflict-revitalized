package asyncreq

import (
	"log/slog"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
)

type ac_player_credentials_req struct{}

type ac_player_credentials_resp struct {
	username string
}

func (req *ac_player_credentials_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_credentials_req) response(session *session.Session) []byte {
	var username string
	bw := bitwriter.NewWriter(make([]byte, 0, 40))
	bw.WriteBeUint16(uint16(types.AC_PLAYER_CREDENTIALS))
	if err := session.Db.Get(&username, "SELECT nickname FROM user WHERE uid=$1", session.Uid); err != nil {
		slog.Error("Failed to get nickname from Database", "err", err, "uid", session.Uid)
		// Err: Disconnect?
	}
	bw.WriteCString(username) // username
	bw.BwWriteByte(0)         // Idk
	bw.WriteBeUint64(0)       // Steam id
	bw.WriteBeUint64(0)       // Some other external(?) Account id
	bw.WriteBeUint32(0)       // some account level thingie?
	bw.WriteBool(false)       // idk
	bw.WriteBeUint32(0)       // Would be a variantdict, but always empty
	return bw.ReturnSlice()
}

func Send_ac_player_credentials(session *session.Session) {
	req := ac_player_credentials_req{}
	resp := req.response(session)
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_credentials(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_player_credentials_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_credentials_req", "error", err)
	}
	resp := req.response(session)
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
