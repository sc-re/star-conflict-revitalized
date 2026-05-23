package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type ac_get_visited_free_space_zones_req struct{}

func (req *ac_get_visited_free_space_zones_req) parse(body []byte) error {
	return nil
}

func (req *ac_get_visited_free_space_zones_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 500))
	bw.WriteBeUint16(uint16(types.AC_GET_VISITED_FREE_SPACE_ZONES))
	visitedZones := map[string]bool{
		"1": true,
	}
	variantdict.BwMarshal(bw, visitedZones)
	return bw.ReturnSlice()
}

func Send_ac_get_visited_free_space_zones(conn net.Conn) {
	req := ac_get_visited_free_space_zones_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_get_visited_free_space_zones(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_get_visited_free_space_zones_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_get_visited_free_space_zones_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
