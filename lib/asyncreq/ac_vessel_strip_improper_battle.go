package asyncreq

import (
	"log"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"time"
)

type ac_vessel_strip_improper_battle_req struct{}

func (req *ac_vessel_strip_improper_battle_req) response() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 8))
	bw.WriteBeUint16(uint16(types.AC_VESSEL_STRIP_IMPROPER_BATTLE))
	bw.BwWriteByte(0)   // status
	bw.WriteBool(false) // has vessel
	bw.WriteBeUint32(0)
	return bw.ReturnSlice(), nil
}

// TODO: Login should trigger code to check if any ships of a given user have unfit equipment
// and update the Database as needed
func Send_ac_vessel_strip_improper_battle(conn net.Conn) {
	// TODO: Figure out what should trigger this
	// Maybe the first time after the client request AC_UNIVERSE_GET?
	time.Sleep(3 * time.Second)
	req := ac_vessel_strip_improper_battle_req{}
	resp, _ := req.response()
	log.Printf("Resp: ac_vessel_strip_improper_battle[]")
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_vessel_strip_imporer_battle(body []byte, seq uint16, seqRet uint16, conn net.Conn) error {
	return nil
}
