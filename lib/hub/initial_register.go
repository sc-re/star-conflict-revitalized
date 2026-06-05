package hub

import (
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type player struct {
	Uid uint64
}

type players struct {
	Players []player
}

func HandleInitialRegister(conn net.Conn, hdr protocol.HubHeader) {
	hdr.DstId = 0x5800000058000000
	hdr.SrcId = uint64(types.ST_DEDICATED_SERVER)
	hdr.ReturnSequence = hdr.Sequence
	hdr.Sequence = 0
	bw := bitwriter.NewWriter(make([]byte, 0, 20))
	bw.WriteBool(false)
	cvars := map[string]string{}
	cvars["db_mongoHosts"] = "127.0.0.27:2727"
	cvars["db_dedicatedDbCollectionSuffix"] = ""
	cvars["db_stats_name"] = "stats_cosmosim"
	cvars["db_name"] = "cosmosim"
	variantdict.BwMarshal(bw, cvars)

	conn.Write(protocol.MakeHubMessage(hdr, bw.ReturnSlice()))

	hdr.ReturnSequence = 0
	hdr.Sequence = 3
	hdr.CommandType = types.CC_CLOUD_CONFIG
	bw = bitwriter.NewWriter(make([]byte, 0, 20))
	bw.WriteBeInt32(1)
	bw.WriteBeUint16(0)
	bw.BwWriteByte(0)
	conn.Write(protocol.MakeHubMessage(hdr, bw.ReturnSlice()))

}
