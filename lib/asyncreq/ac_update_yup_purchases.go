package asyncreq

import (
	"fmt"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

// We probably only ever need ac_update_yup_purchases if we want to emulate Gaijin Purchases/DLC behaviour
// The server usually sends the response unsolicited for any account that has a DLC

type ac_update_yup_purchases_req struct{}

func (req *ac_update_yup_purchases_req) response() ([]byte, error) {
	return nil, fmt.Errorf("Response for ac_update_yup_purchases not implemented")
}

func handle_ac_update_yup_purchases(body []byte, seq uint16, seqRet uint16, conn net.Conn) error {
	req := ac_update_yup_purchases_req{}
	if len(body) > 0 {
		return fmt.Errorf("Malformed ac_update_yp_purchases request")
	}
	resp, err := req.response()
	if err != nil {
		return err
	}
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
	return nil
}
