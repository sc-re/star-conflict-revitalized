package battlepassactivation

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

func ScmdBattlePassActivation() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	if err := BwScmdBattlePassActivation(bw); err != nil {
		return nil, err
	}
	return bw.ReturnSlice(), nil
}

func BwScmdBattlePassActivation(bw *bitwriter.Writer) error {
	if err := bw.WriteBeUint64(1772168400000); err != nil {
		return err
	}
	battlePassCount := uint16(3)
	if err := bw.WriteBeUint16(battlePassCount); err != nil {
		return err
	}
	for i := range battlePassCount {
		if err := bw.WriteBeUint16(i); err != nil {
			return err
		}
		if err := bw.WriteBeUint64(1780376400000); err != nil {
			return err
		}

	}
	return nil
}

func SendScmdBattlePassActivation(conn net.Conn) {
	resp, err := ScmdBattlePassActivation()
	if err != nil {
		slog.Error("Failed to create ScmdBattlePassActivation", "error", err)
	}
	conn.Write(protocol.MakeMessage(types.SCMD_BATTLE_PASS_ACTIVATION, 0, 0, resp))

}
