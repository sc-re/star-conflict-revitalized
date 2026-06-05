package hub

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type registerConfig struct {
	ZoneId int32
	//GameSessionId        int32 `variantdict:"gameSessionId"`
	Cid                  uint64
	CsMapIdx             int32
	CsRace               int32
	MaxVesselTier        int32
	ClanShipFit          struct{}
	ZoneBlocks           struct{}
	AutoJoinUid          uint64
	AutoJoinCid          uint64
	InstanceTransferData struct{}
}

func HandleRegister(conn net.Conn, hdr protocol.HubHeader, body []byte) {
	slog.Debug("HandleRegister", "header", hdr, "body", body)
	hdr.ReturnSequence = hdr.Sequence
	hdr.Sequence = 2
	bw := bitwriter.NewWriter(make([]byte, 0, 20))
	bw.WriteBool(true)
	config := registerConfig{
		ZoneId: 2,
	}
	variantdict.BwMarshal(bw, config)
	conn.Write(protocol.MakeHubMessage(hdr, bw.ReturnSlice()))

	hdr.ReturnSequence = 0
	hdr.Sequence = 3
	hdr.CommandType = types.CC_REGISTER_PLAYER_LIST
	bw = bitwriter.NewWriter(make([]byte, 0, 20))
	bw.WriteBeUint64(0xffffffff)
	variantdict.BwMarshal(bw, players{
		[]player{
			{Uid: 19},
			{Uid: 10},
		},
	})
	conn.Write(protocol.MakeHubMessage(hdr, bw.ReturnSlice()))
}
