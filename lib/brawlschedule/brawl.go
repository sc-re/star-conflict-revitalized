package brawlschedule

import (
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type intPair struct {
	a uint32
	b uint32
}

type brawl_schedule_entry struct {
	schedule []intPair
	name     string
}

func ScmdBrawlSchedule() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 250))
	BwScmdBrawlSchedule(bw)
	return bw.ReturnSlice()
}

// Original Server currently just return 28 times the same data
func BwScmdBrawlSchedule(bw *bitwriter.Writer) {
	entry := brawl_schedule_entry{
		schedule: []intPair{{1, 2}, {9, 10}, {15, 16}, {17, 18}, {19, 20}, {21, 22}},
		name:     "gtdm",
	}
	for range 28 {
		bw.WriteBeUint32(uint32(len(entry.schedule)))

		for _, s := range entry.schedule {
			bw.WriteBeUint32(s.a)
			bw.WriteBeUint32(s.b)
		}
		bw.WriteCString(entry.name)
	}
}

func SendScmdBrawlSchedule(conn net.Conn) {
	resp := ScmdBrawlSchedule()
	/*
		if err != nil {
			slog.Error("Failed to create ScmdBrawlSchedule", "error", err)
			return
		}*/
	conn.Write(protocol.MakeMessage(types.SCMD_BRAWL_SCHEDULE, 0, 0, resp))
}
