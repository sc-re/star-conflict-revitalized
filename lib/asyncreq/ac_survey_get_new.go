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

type ac_survey_get_new_req struct {
	lang string
}

func (req *ac_survey_get_new_req) parse(body []byte) error {
	br := bitreader.NewReader(body)
	var err error
	req.lang, err = br.ReadCString()
	return err
}

type answer struct {
	Text      string
	RecordIdx int32
}

// XXX: Mostly correct, but not fully, needs investigating
type survey struct {
	Sid        uint64
	Question   string
	IsResults  bool
	Multiple   int32
	Answers    map[string]answer
	Results    map[string]uint64
	TotalVoted uint64
}

func (req *ac_survey_get_new_req) response() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 7))
	bw.WriteBeUint16(uint16(types.AC_SURVEY_GET_NEW))
	bw.BwWriteByte(0)
	surveys := map[string]string{}
	if err := variantdict.BwMarshal(bw, surveys); err != nil {
		slog.Error("Failed to marshal ac_survey_get_new", "err", err)
	}

	return bw.ReturnSlice()
}

func handle_ac_survey_get_new(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_survey_get_new_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_survey_get_new_req", "error", err)
	}
	resp := req.response()
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
