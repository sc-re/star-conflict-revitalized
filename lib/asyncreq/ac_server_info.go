package asyncreq

import (
	"encoding/binary"
	"log"
	"math"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"time"
)

type ac_server_info_req struct{}

func (req *ac_server_info_req) response(sandboxAccess uint32, mm_enable_pve_raids uint8, mm_enable_league uint8, mm_enable_coop_vs_ai uint8) []byte {
	resp := make([]byte, 22)
	binary.BigEndian.PutUint16(resp[0:], uint16(types.AC_SERVER_INFO))
	timestamp := float64(time.Now().Unix())
	binary.LittleEndian.PutUint64(resp[2:], math.Float64bits(timestamp))
	// TODO: Figure out what this does
	binary.LittleEndian.PutUint32(resp[10:], 0x0bd18549)
	binary.LittleEndian.PutUint32(resp[14:], sandboxAccess)
	resp[18] = mm_enable_pve_raids
	resp[19] = mm_enable_league
	resp[20] = mm_enable_coop_vs_ai

	return resp
}

func handle_ac_server_info(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_server_info_req{}
	if len(body) > 2 {
		// err
	}
	log.Printf("Req: ac_server_info[%v]", req)
	resp := req.response(4, 0, 0, 0)
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
