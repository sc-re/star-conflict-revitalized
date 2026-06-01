package asyncreq

import (
	"encoding/binary"
	"errors"
	"log/slog"
	"starconflict/lib/bitreader"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"

	"github.com/mattn/go-sqlite3"
)

type NickNameChangeResultType byte

const (
	NickNameChangeResultOk NickNameChangeResultType = iota
	NickNameChangeResultTimeout
	NickNameChangeResultCantAfford
	NickNameChangeResultNickTaken
	NickNameChangeResultInvalidSymbol
	NickNameChangeResultInvalidLength
	NickNameChangeResultError
)

type ac_change_player_nickname_req struct {
	name string
}

func (req *ac_change_player_nickname_req) parse(body []byte) error {
	var err error
	br := bitreader.NewReader(body)
	req.name, err = br.ReadCString()
	return err
}

func (req *ac_change_player_nickname_req) response(status NickNameChangeResultType) []byte {
	ret := make([]byte, 0, 20)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_CHANGE_PLAYER_NICKNAME))
	ret = append(ret, byte(status))
	ret = append(ret, []byte(req.name)...)
	ret = append(ret, 0)
	return ret
}

func (req *ac_change_player_nickname_req) verifyNickName() (NickNameChangeResultType, bool) {
	if len(req.name) > 16 || len(req.name) < 2 {
		return NickNameChangeResultInvalidLength, true
	}
	for _, c := range req.name {
		if ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {

		} else {
			return NickNameChangeResultInvalidSymbol, true
		}
	}
	return NickNameChangeResultOk, false
}

// XXX: Check if we might have to sent some notifications after changing the name?
func handle_ac_change_player_nickname(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_change_player_nickname_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_change_player_nickname_req", "error", err)
	}
	resp := req.response(NickNameChangeResultOk)
	if result, invalid := req.verifyNickName(); invalid {
		resp = req.response(result)
	} else if result, err := session.Db.Exec("UPDATE user SET nickname = $1 WHERE uid = $2", req.name, session.Uid); err != nil {
		resp = req.response(NickNameChangeResultError)
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				resp = req.response(NickNameChangeResultNickTaken)
			} else {
				slog.Warn("Db error", "err", err)
			}
		}
	} else {
		if rows, err := result.RowsAffected(); err != nil {
			resp = req.response(NickNameChangeResultError)
		} else if rows != 1 {
			resp = req.response(NickNameChangeResultError)
		}
	}
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
