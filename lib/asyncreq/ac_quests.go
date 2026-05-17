package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_quests_req struct{}

func (req *ac_quests_req) parse(body []byte) error {
	return nil
}

func (req *ac_quests_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 1000))
	bw.WriteBeUint16(uint16(types.AC_QUESTS))
	bw.WriteBeUint32(0)      // idk
	bw.WriteBeUint32(493056) // idk
	bw.WriteBeUint32(534360) // idk
	bw.WriteBool(false)      // idk
	bw.WriteBool(false)      // idk
	bw.WriteBool(false)      // idk
	bw.WriteBool(true)       // idk
	bw.BwWriteByte(0)        // daily quest count
	for range 0 {
		bw.WriteBeUint32(0) // quest id
		bw.WriteBool(true)  // state
	}
	bw.BwWriteByte(0)        // quest count
	bw.BwWriteByte(0)        // quest desc count
	bw.WriteBeUint16(0)      // a count
	bw.WriteBeUint32(455727) // idk
	bw.WriteBeUint16(0)      // b count
	bw.BwWriteByte(0xff)

	return bw.ReturnSlice()
}

func handle_ac_quests(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_quests_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_quests_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
