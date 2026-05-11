package asyncreq

import (
	"encoding/binary"
	"log"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_zones_lua_active_events_update_req struct{}

func (req *ac_zones_lua_active_events_update_req) response() []byte {
	resp := make([]byte, 3)
	binary.BigEndian.PutUint16(resp[0:], uint16(types.AC_ZONES_LUA_ACTIVE_EVENTS_UPDATE))
	resp[2] = 0x0
	return resp
}

func handle_ac_zones_lua_active_events_update(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_zones_lua_active_events_update_req{}
	// TODO: actually parse and do something with body
	log.Printf("Req: ac_zones_lua_active_events_update[%v]", req)
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
