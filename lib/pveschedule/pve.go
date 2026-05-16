package pveschedule

import (
	"log/slog"
	"net"

	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type pveSchedule struct {
	one   string `variantdict:"1"`
	two   string `variantdict:"2"`
	three string `variantdict:"3"`
}

func ScmdPveSchedule() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	if err := BwScmdPveSchedule(bw); err != nil {
		return nil, err
	}
	return bw.ReturnSlice(), nil
}

func BwScmdPveSchedule(bw *bitwriter.Writer) error {
	schedule := pveSchedule{
		one:   "pve_magnificent_seven",
		two:   "s1347_geostation",
		three: "pve_stealth",
	}
	err := variantdict.BwMarshal(bw, schedule)
	return err
}

func SendScmdPveSchedule(conn net.Conn) {
	resp, err := ScmdPveSchedule()
	if err != nil {
		slog.Error("Failed to create ScmdPveSchedule")
		return
	}
	conn.Write(protocol.MakeMessage(types.SCMD_PVE_SCHEDULE, 0, 0, resp))
}
