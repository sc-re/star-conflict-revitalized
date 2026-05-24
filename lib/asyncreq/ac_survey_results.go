package asyncreq

import (
	"log/slog"
	"starconflict/lib/bitreader"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type ac_survey_results_req struct {
	lang string
}

func (req *ac_survey_results_req) parse(body []byte) error {
	br := bitreader.NewReader(body)
	var err error
	req.lang, err = br.ReadCString()
	return err
}

// Probably same format as ac_survey_get_new
func (req *ac_survey_results_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 7))
	bw.WriteBeUint16(uint16(types.AC_SURVEY_RESULTS))
	bw.BwWriteByte(0)
	survey := map[string]string{}
	if err := variantdict.BwMarshal(bw, survey); err != nil {
		slog.Error("Failed to marshal ac_survey_results", "err", err)
	}
	return bw.ReturnSlice()
}

func handle_ac_survey_results(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_survey_results_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_survey_results_req", "error", err)
	}
	resp := req.response()
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
