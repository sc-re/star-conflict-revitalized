package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_player_credentials_req struct{}

func (req *ac_player_credentials_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_credentials_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 40))
	bw.WriteBeUint16(uint16(types.AC_PLAYER_CREDENTIALS))
	bw.WriteCString("Test") // username
	bw.BwWriteByte(1)       // Idk
	bw.WriteBeUint64(0)     // Steam id
	bw.WriteBeUint64(10)    // Some other external(?) Account id
	bw.WriteBeUint32(10)    // some account level thingie?
	bw.WriteBool(false)     // idk
	bw.WriteBeUint32(0)     // Would be a variantdict, but always empty
	return bw.ReturnSlice()
}

func Send_ac_player_credentials(conn net.Conn) {
	req := ac_player_credentials_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_credentials(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_player_credentials_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_credentials_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
