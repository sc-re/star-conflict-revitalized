package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type ac_get_craft_resources_req struct{}

func (req *ac_get_craft_resources_req) parse(body []byte) error {
	return nil
}

func (req *ac_get_craft_resources_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 500))
	bw.WriteBeUint16(uint16(types.AC_GET_CRAFT_RESOURCES))
	resources := map[string]int32{
		"ow_Impure_iridium": 321,
	}
	variantdict.BwMarshal(bw, resources)
	return bw.ReturnSlice()
}

func Send_ac_get_craft_resources(conn net.Conn) {
	req := ac_get_craft_resources_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_get_craft_resources(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_get_craft_resources_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_get_craft_resources_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
