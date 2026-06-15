package asyncreq

import (
	"context"
	"encoding/binary"
	"log/slog"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ac_set_userdata_req struct {
	userdata []byte
}

func (req *ac_set_userdata_req) parse(body []byte) error {
	req.userdata = body
	return nil
}

func (req *ac_set_userdata_req) response(ver uint32) []byte {
	ret := make([]byte, 0, 6)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_SET_USERDATA))
	// version should probably be the on in the requst,
	// but it seems like the clinet doesn't care what we return
	ret = binary.BigEndian.AppendUint32(ret, ver)
	return ret
}

func handle_ac_set_userdata(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_set_userdata_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_set_userdata_req", "error", err)
	}
	// XXX: This is a free data storage API, no size limit, no data validation, just pure storage
	if _, err := session.Db.Database("cosmosim").Collection("accounts").UpdateOne(context.TODO(), bson.M{"uid": session.Uid}, bson.M{"$set": bson.M{"userdata": req.userdata}}); err != nil {
		slog.Error("Failed to save userdata", "error", err)
		return
	}

	resp := req.response(1)
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
