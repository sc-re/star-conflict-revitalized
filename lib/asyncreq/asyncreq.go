package asyncreq

import (
	"encoding/binary"
	"log"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

func HandleAsyncReq(hdr *protocol.Header, body []byte, seq uint16, conn net.Conn, uid uint64) {
	actype := types.AsyncReqType(binary.BigEndian.Uint16(body))
	switch actype {
	case types.AC_WELCOME_MSG:
		handle_ac_welcome_msg(body, seq, hdr.Sequence, conn)
	case types.AC_MOTD:
		handle_ac_motd(body, seq, hdr.Sequence, conn)
	case types.AC_SERVER_INFO:
		handle_ac_server_info(body, seq, hdr.Sequence, conn)
	case types.AC_UNIVERSE_GET:
		handle_ac_universe_get(body, seq, hdr.Sequence, conn)
	case types.AC_ZONES_LUA_ACTIVE_EVENTS_UPDATE:
		handle_ac_zones_lua_active_events_update(body, seq, hdr.Sequence, conn)
	default:
		log.Printf("Unhandled AsyncReq of type: %s", actype)
	}
}
