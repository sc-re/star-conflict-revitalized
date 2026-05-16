package brawlschedule

import (
	"log/slog"
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

func ScmdBrawlSchedule() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 250))
	if err := BwScmdBrawlSchedule(bw); err != nil {
		return nil, err
	}
	return bw.ReturnSlice(), nil
}

// Original Server currently just return 28 times the same data
func BwScmdBrawlSchedule(bw *bitwriter.Writer) error {
	entry := brawl_schedule_entry{
		schedule: []intPair{{1, 2}, {9, 10}, {15, 16}, {17, 18}, {19, 20}, {21, 22}},
		name:     "gtdm",
	}
	for range 28 {
		if err := bw.WriteBeUint32(uint32(len(entry.schedule))); err != nil {
			return err
		}

		for _, s := range entry.schedule {
			if err := bw.WriteBeUint32(s.a); err != nil {
				return err
			}

			if err := bw.WriteBeUint32(s.b); err != nil {
				return err
			}

		}
		if err := bw.WriteCString(entry.name); err != nil {
			return err
		}

	}
	return nil
}

func SendScmdBrawlSchedule(conn net.Conn) {
	resp, err := ScmdBrawlSchedule()
	if err != nil {
		slog.Error("Failed to create ScmdBrawlSchedule", "error", err)
		return
	}
	conn.Write(protocol.MakeMessage(types.SCMD_BRAWL_SCHEDULE, 0, 0, resp))
}
