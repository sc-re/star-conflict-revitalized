package asyncreq

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_motd_req struct {
	lang string
}

func (req *ac_motd_req) parse(body []byte) error {
	r := bytes.NewReader(body)
	err := *new(error)
	req.lang, err = protocol.ReadCString(r)
	if err != nil {
		return err
	}
	return nil
}

func (req *ac_motd_req) response(motd string) []byte {
	// TODO: Actually observe a motd that isn't empty
	_ = motd
	resp := make([]byte, 3)
	binary.BigEndian.PutUint16(resp[0:], uint16(types.AC_MOTD))
	resp[2] = 0x04
	return resp
}

func handle_ac_motd(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_motd_req{}
	motd := ""
	if err := req.parse(body[2:]); err != nil {
		// conn.Write(Failure)
	}
	log.Printf("Req: ac_motd[%v]", req)
	resp := req.response(motd)
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
