package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
	"strconv"
)

type LeaderBoardQueryType float32

const (
	LeaderBoardQueryRange LeaderBoardQueryType = iota
	LeaderBoardQueryEntityList
	LeaderBoardQuerySize
	LeaderBoardQueryNum
)

type leaderboardEntry struct {
	e uint64  // uid
	p int32   // position
	s float32 // value
}

type ac_leaderboard_get_range_resp []leaderboardEntry
type ac_leaderboard_get_entity_list_resp []leaderboardEntry
type ac_leaderboard_get_size_resp []int32
type ac_leaderboard_get_num_resp []int32

type ac_leaderboard_get_req struct {
	Lb   string // Leaderboard
	Qt   LeaderBoardQueryType
	To   float32  // LQT_RANGE
	From float32  // LQT_RANGE
	List []string // LQT_ENTITY_LIST
}

func (req *ac_leaderboard_get_req) parse(body []byte) error {
	if err := variantdict.Unmarshal(body, req); err != nil {
		return err
	}
	slog.Debug("Parses request", "req", req)
	return nil
}

func (req *ac_leaderboard_get_req) response() []byte {
	ret := make([]byte, 0, 7)
	ret = binary.BigEndian.AppendUint16(ret, uint16(types.AC_LEADERBOARD_GET))
	ret = append(ret, 4)
	ret = binary.BigEndian.AppendUint64(ret, 0)
	return ret
}

func getLeaderboardSize() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 13))
	bw.WriteBeUint16(uint16(types.AC_LEADERBOARD_GET))
	bw.BwWriteByte(0) // status
	resp := ac_leaderboard_get_size_resp{748}
	if err := variantdict.BwMarshal(bw, resp); err != nil {
		slog.Error("Failed to marshal response", "err", err)
		return nil
	}
	return bw.ReturnSlice()
}

func (req *ac_leaderboard_get_req) getEntities() []byte {
	entities := ac_leaderboard_get_entity_list_resp{}
	for _, uid := range req.List {
		u, _ := strconv.Atoi(uid)
		entities = append(entities, leaderboardEntry{
			e: uint64(u),
			p: 0,
			s: 500.0,
		})
	}
	bw := bitwriter.NewWriter(make([]byte, 0, 50))
	bw.WriteBeUint16(uint16(types.AC_LEADERBOARD_GET))
	bw.BwWriteByte(0) // status
	variantdict.BwMarshal(bw, entities)
	return bw.ReturnSlice()
}

func handle_ac_leaderboard_get(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_leaderboard_get_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_leaderboard_get_req", "error", err)
	}
	switch req.Qt {
	case LeaderBoardQueryRange:
	case LeaderBoardQueryEntityList:
		slog.Debug("Got leaderboard entity list request", "leaderboard", req.Lb, "uids", req.List)
		resp := req.getEntities()
		session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
	case LeaderBoardQuerySize:
		slog.Debug("Got Leaderboard size request", "leaderboard", req.Lb)
		resp := getLeaderboardSize()
		session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
	case LeaderBoardQueryNum:
	default:
		resp := req.response()
		session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
		slog.Error("Got unknown Leaderboard Query: %v", "QueryType", req.Qt)
	}
}
