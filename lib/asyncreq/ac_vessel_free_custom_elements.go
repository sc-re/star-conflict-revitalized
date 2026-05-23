package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_vessel_free_custom_elements_req struct{}

func (req *ac_vessel_free_custom_elements_req) parse(body []byte) error {
	return nil
}

func (req *ac_vessel_free_custom_elements_req) response() []byte {
	ret := make([]byte, 0, 100)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_VESSEL_FREE_CUSTOM_ELEMENTS))
	customElements := []string{"ny_2018_03"}
	ret = binary.BigEndian.AppendUint32(ret, uint32(len(customElements)))
	for _, v := range customElements {
		ret = append(ret, []byte(v)...)
		ret = append(ret, 0x00)
	}
	return ret
}

func Send_ac_vessel_free_custom_elements(conn net.Conn) {
	req := ac_vessel_free_custom_elements_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_vessel_free_custom_elements(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_vessel_free_custom_elements_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_vessel_free_custom_elements_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
