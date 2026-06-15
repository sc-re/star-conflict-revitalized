package asyncreq

import (
	"context"
	"log/slog"
	"starconflict/lib/bitwriter"
	"starconflict/lib/dbtypes"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"

	"go.mongodb.org/mongo-driver/v2/bson"
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
	var account dbtypes.Accounts
	if err := session.Db.Database("cosmosim").Collection("accounts").FindOne(context.TODO(), bson.M{"uid": session.Uid}).Decode(&account); err != nil {
		slog.Error("Failed to get account", "error", err)
		return nil
	}

	bw.WriteBeUint64(account.Credits)
	if req.flag&0x02 != 0 {
		bw.WriteBeUint64(account.GoldCredits)
	}
	if req.flag&0x04 != 0 {
		bw.WriteBeUint64(account.TokenCredits)
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
		bw.WriteBeUint32(account.VesselExpPool)
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
