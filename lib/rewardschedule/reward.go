package rewardschedule

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type rewardScheduleEntry struct {
	GoldReward float32
	HourBegin  float32
	HourEnd    float32
}

func ScmdRewardSchedule() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	if err := BwScmdRewardSchedule(bw); err != nil {
		return nil, err
	}
	return bw.ReturnSlice(), nil
}

func BwScmdRewardSchedule(bw *bitwriter.Writer) error {
	rewards := map[string]rewardScheduleEntry{
		"0": {GoldReward: 10.0, HourBegin: 13.0, HourEnd: 24.0},
		"1": {GoldReward: 10.0, HourBegin: 0.0, HourEnd: 1.0},
	}
	err := variantdict.BwMarshal(bw, rewards)
	return err
}

func SendScmdRewardSchedule(conn net.Conn) {
	resp, err := ScmdRewardSchedule()
	if err != nil {
		slog.Error("Failed to creeate ScmdRewardSchedule", "error", err)
		return
	}
	conn.Write(protocol.MakeMessage(types.SCMD_BRAWL_SCHEDULE, 0, 0, resp))
}
