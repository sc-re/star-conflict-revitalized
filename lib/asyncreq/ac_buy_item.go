package asyncreq

import (
	"log/slog"

	"starconflict/lib/bitreader"
	"starconflict/lib/bitwriter"
	"starconflict/lib/cmdstore"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
	"starconflict/lib/types/store"
	_ "starconflict/lib/variantdict"
)

type ac_buy_item_req struct {
	storeItemId uint32
	amout       uint32 // max 0x3fff
	creditsType store.CreditsType
	hasDiscount bool
	mode        uint32
}

type ac_buy_item_resp struct {
	status store.StoreBuyResultType
}

func (req *ac_buy_item_req) parse(body []byte) error {
	br := bitreader.NewReader(body)
	var err error
	req.storeItemId, err = br.ReadBeUint32()
	if err != nil {
		return err
	}
	req.amout, err = br.ReadBeUint32()
	if err != nil {
		return err
	}
	req.hasDiscount, err = br.ReadBool()
	if err != nil {
		return err
	}
	if v, err := br.ReadByte(); err != nil {
		return err
	} else {
		req.creditsType = store.CreditsType(v)
	}
	return nil
}

func (req *ac_buy_item_req) response(resp *ac_buy_item_resp) []byte {
	br := bitwriter.NewWriter(make([]byte, 0, 30))
	br.WriteBeUint16(uint16(types.AC_BUY_ITEM))
	br.WriteBeUint32(req.storeItemId)
	br.BwWriteByte(byte(byte(resp.status)))
	br.WriteBeUint64(0) // Item Id/Vessel ID/whatever
	br.WriteCString("") // Def_name
	br.WriteBeUint32(0) // Quantitiy
	br.WriteBool(false) // idk
	br.WriteBeUint64(0) // idk
	return br.ReturnSlice()
}

// XXX: We need abstraction over credit handling to send notifications after value changes
func (req *ac_buy_item_req) handle_motto(session *session.Session, item cmdstore.StoreItem) []byte {
	resp := ac_buy_item_resp{
		status: store.StoreBuyResultError,
	}
	if req.creditsType != store.CreditsTypeGold {
		resp.status = store.StoreBuyResultInvalidData
		return req.response(&resp)
	}
	tx, err := session.Db.Begin()
	if err != nil {
		slog.Error("Failed to start Database transaction")
		return req.response(&resp)
	}
	defer tx.Rollback()
	if result, err := tx.Exec("UPDATE credits SET gold = gold - $1 WHERE uid=$2 AND gold >= $1", item.PremiumPrice, session.Uid); err != nil {
		slog.Error("Failed to update database", "err", err)
		resp.status = store.StoreBuyResultError
		return req.response(&resp)
	} else {
		if rows, err := result.RowsAffected(); err != nil {
			slog.Error("Failed to check affected rows", "err", err)
			resp.status = store.StoreBuyResultError
			return req.response(&resp)
		} else if rows != 1 {
			slog.Error("Failed to substract gold", "err", err)
			resp.status = store.StoreBuyResultCantAfford
			return req.response(&resp)
		}
	}
	if _, err := tx.Exec("INSERT INTO mottos (motto, uid) VALUES ($1, $2)", item.ItemName, session.Uid); err != nil {
		slog.Error("failed to insert Motto", "err", err)
		resp.status = store.StoreBuyResultError
		return req.response(&resp)
	}
	if err := tx.Commit(); err != nil {
		slog.Error("Failed to commit transaction", "err", err)
		resp.status = store.StoreBuyResultError
		return req.response(&resp)
	}

	//session.Conn.Write(protocol.MakeMessage(types.SCMD_NOTIFICATION, 0, 0, variantdict.Marshal(struct{})))
	resp.status = store.StoreBuyResultOk

	return req.response(&resp)
}

func handle_ac_buy_item(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_buy_item_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_buy_item_req", "error", err)
	}
	item := cmdstore.StoreItem{}
	if err := session.Db.Get(&item, "SELECT * FROM store WHERE StoreItemId=$1", req.storeItemId); err != nil {
		slog.Error("Failed to get item from database", "StoreItemId", req.storeItemId)
	}

	resp := req.response(&ac_buy_item_resp{status: store.StoreBuyResultSteamDeniedTransaction})

	switch item.ItemType {
	case cmdstore.ItemTypeVessel:
	case cmdstore.ItemTypeModule:
	case cmdstore.ItemTypeRocket:
	case cmdstore.ItemTypeDrug:
	case cmdstore.ItemTypeBundle:
	case cmdstore.ItemTypeResource:
	case cmdstore.ItemTypeBlueprint:
	case cmdstore.ItemTypeAvatar:
	case cmdstore.ItemTypeMotto:
		resp = req.handle_motto(session, item)
	case cmdstore.ItemTypeJunk:
	}

	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
