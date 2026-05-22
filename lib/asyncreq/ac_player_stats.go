package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type ac_player_stats_req struct{}

type ac_player_stats_resp struct {
	GamePlayed            int32
	AntibotCounter        int32
	MachineAntibotCounter int32
	TwoStepAuthEnabled    int32 `variantdict:"2StepAuthEnabled"`
}

func (req *ac_player_stats_req) parse(body []byte) error {
	return nil
}

func (req *ac_player_stats_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 92))
	stats := ac_player_stats_resp{
		GamePlayed:            500,
		AntibotCounter:        500,
		MachineAntibotCounter: 0,
		TwoStepAuthEnabled:    0,
	}

	bw.WriteBeUint16(uint16(types.AC_PLAYER_STATS))
	variantdict.BwMarshal(bw, stats)
	return bw.ReturnSlice()
}

func Send_ac_player_stats(conn net.Conn) {
	req := ac_player_stats_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_player_stats(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_player_stats_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_player_stats_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
