package asyncreq

import (
	"log/slog"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
)

type ac_player_credits_req struct {
	flag uint16
}

type ac_player_credits_resp struct {
	Gold        uint64
	Iridium     uint64
	Credits     uint64
	FreeSynergy uint32
}

func (req *ac_player_credits_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_credits_req) response(session *session.Session) []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 68))
	bw.WriteBeUint16(uint16(types.AC_PLAYER_CREDITS))
	bw.WriteBeUint16(req.flag)
	credits := ac_player_credits_resp{}
	if err := session.Db.Get(&credits, "SELECT gold, iridium, credits, freeSynergy FROM credits WHERE uid=$1", session.Uid); err != nil {
		slog.Error("Failed to get Credits from database", "err", err, "uid", session.Uid)
	}

	bw.WriteBeUint64(credits.Credits)
	if req.flag&0x02 != 0 {
		bw.WriteBeUint64(credits.Gold)
	}
	if req.flag&0x04 != 0 {
		bw.WriteBeUint64(credits.Iridium)
	}
	if req.flag&0x08 != 0 {
		bw.WriteBeUint64(4) // Xenochips
	}
	if req.flag&0x10 != 0 {
		bw.WriteBeUint64(1779558934) // Premium timestamp?
	}
	if req.flag&0x20 != 0 {
		bw.WriteBeUint64(0) // idk, vid
	}
	if req.flag&0x40 != 0 {
		bw.WriteBeUint32(credits.FreeSynergy)
	}
	// idk, some ressources?
	if req.flag&0x80 != 0 {
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
		bw.WriteBeUint32(0)
	}
	return bw.ReturnSlice()
}

func Send_ac_player_credits(session *session.Session) {
	req := ac_player_credits_req{flag: 222}

	resp := req.response(session)
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_credits(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_player_credits_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_credits_req", "error", err)
	}
	resp := req.response(session)
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
