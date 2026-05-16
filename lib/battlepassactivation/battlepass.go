package battlepassactivation

import (
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

func ScmdBattlePassActivation() []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 100))
	BwScmdBattlePassActivation(bw)
	return bw.ReturnSlice()
}

func BwScmdBattlePassActivation(bw *bitwriter.Writer) {
	bw.WriteBeUint64(1772168400000)
	battlePassCount := uint16(3)
	bw.WriteBeUint16(battlePassCount)
	for i := range battlePassCount {
		bw.WriteBeUint16(i)
		bw.WriteBeUint64(1780376400000)

	}
}

func SendScmdBattlePassActivation(conn net.Conn) {
	resp := ScmdBattlePassActivation()
	/*
		if err != nil {
			slog.Error("Failed to create ScmdBattlePassActivation", "error", err)
		}
	*/
	conn.Write(protocol.MakeMessage(types.SCMD_BATTLE_PASS_ACTIVATION, 0, 0, resp))

}
