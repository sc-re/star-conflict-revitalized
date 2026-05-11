package asyncreq

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_welcome_msg_req struct {
	lang string
}

func (req *ac_welcome_msg_req) parse(body []byte) error {
	r := bytes.NewReader(body)
	err := *new(error)
	req.lang, err = protocol.ReadCString(r)
	if err != nil {
		return err
	}
	return nil
}

func (req *ac_welcome_msg_req) response(welcomeMsg string) []byte {
	resp := make([]byte, len(welcomeMsg)+4)
	binary.BigEndian.PutUint16(resp[0:], uint16(types.AC_WELCOME_MSG))
	resp[2] = 0x0
	copy(resp[3:], []byte(welcomeMsg))
	return resp
}

func handle_ac_welcome_msg(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_welcome_msg_req{}
	welcomeMsg := ""
	if err := req.parse(body[2:]); err != nil {
		// conn.Write(Failure)
	}
	log.Printf("Req: ac_welcome_msg[%v]", req)
	switch req.lang {
	case "en":
		welcomeMsg = "Welcome <br><br>Here"
	case "ru":
		welcomeMsg = "Привет"
	default:
		welcomeMsg = "Wecome <br><br>Unset language"
	}
	resp := req.response(welcomeMsg)
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
