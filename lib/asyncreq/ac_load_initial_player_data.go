package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/battlepassactivation"
	"starconflict/lib/bitwriter"
	"starconflict/lib/brawlschedule"
	"starconflict/lib/leagueforbiddenequipment"
	"starconflict/lib/protocol"
	"starconflict/lib/pveschedule"
	"starconflict/lib/rewardschedule"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type pveModifierEntry struct {
	a    uint32
	b    uint32
	name string
}

type ac_load_initial_player_data_req struct {
	lang string
}

func (req *ac_load_initial_player_data_req) parse(body []byte) error {
	return nil
}

// XXX: There is probably some message somewhere that will also send this dict
func unlimPveMissionLevel(bw *bitwriter.Writer) error {
	level := map[string]int32{
		"planet_war_waves_T1": 25,
		"piratebay_hard":      69,
	}
	return variantdict.BwMarshal(bw, level)
}

func battlePassPlayerData(bw *bitwriter.Writer) {
	bw.WriteBeUint64(1748581200000)
	type stageDataEntry struct {
		idk                  uint16
		highestUnlockedStage uint16
	}
	stageData := []stageDataEntry{
		{idk: 0, highestUnlockedStage: 37},
		{idk: 1, highestUnlockedStage: 1},
	}
	bw.WriteBeUint16(uint16(len(stageData)))
	for _, v := range stageData {
		bw.WriteBeUint16(v.idk)
		bw.WriteBeUint16(v.highestUnlockedStage)
	}
	unlockedRewards := []string{"Bundle_Battle_Pass_2_1"}
	bw.WriteBeUint32(uint32(len(unlockedRewards)))
	for _, v := range unlockedRewards {
		bw.WriteCString(v)
	}
	type somethingTimed struct {
		timeStamp uint64
		idk       []uint16
	}
	var idk []somethingTimed
	bw.WriteBeUint32(uint32(len(idk)))
	for _, v := range idk {
		bw.WriteBeUint64(v.timeStamp)
		bw.WriteBeUint32(uint32(len(v.idk)))
		for _, v := range v.idk {
			bw.WriteBeUint16(v)
		}
	}
}

// XXX: Check what these numbers actually do
func pveModifiers(bw *bitwriter.Writer) error {
	data := []pveModifierEntry{
		{a: 0, b: 0, name: "s8256_pve_raid_planetoid"},
		{a: 2, b: 2, name: "s8256_pve_raid_planetoid"},
		{a: 4, b: 4, name: "s8256_pve_raid_planetoid"},
		{a: 6, b: 6, name: "s8256_pve_raid_planetoid"},
		{a: 8, b: 8, name: "s8256_pve_raid_planetoid"},
		{a: 10, b: 10, name: "s8256_pve_raid_planetoid"},
		{a: 12, b: 12, name: "s8256_pve_raid_planetoid"},
		{a: 14, b: 14, name: "s8256_pve_raid_planetoid"},
		{a: 16, b: 16, name: "s8256_pve_raid_planetoid"},
		{a: 18, b: 18, name: "s8256_pve_raid_planetoid"},
		{a: 20, b: 20, name: "s8256_pve_raid_planetoid"},
		{a: 22, b: 22, name: "s8256_pve_raid_planetoid"},
		{a: 24, b: 24, name: "pve_raid_waterharvest"},
		{a: 26, b: 26, name: "pve_raid_waterharvest"},
		{a: 28, b: 28, name: "pve_raid_waterharvest"},
		{a: 30, b: 30, name: "pve_raid_waterharvest"},
		{a: 32, b: 32, name: "pve_raid_waterharvest"},
		{a: 34, b: 34, name: "pve_raid_waterharvest"},
		{a: 36, b: 36, name: "pve_raid_waterharvest"},
		{a: 38, b: 38, name: "pve_raid_waterharvest"},
		{a: 40, b: 40, name: "pve_raid_waterharvest"},
		{a: 42, b: 42, name: "pve_raid_waterharvest"},
		{a: 44, b: 44, name: "pve_raid_waterharvest"},
		{a: 46, b: 46, name: "pve_raid_waterharvest"},
		{a: 48, b: 48, name: "pve_raid"},
		{a: 50, b: 50, name: "pve_raid"},
		{a: 52, b: 52, name: "pve_raid"},
		{a: 54, b: 54, name: "pve_raid"},
		{a: 56, b: 56, name: "pve_raid"},
		{a: 58, b: 58, name: "pve_raid"},
		{a: 60, b: 60, name: "pve_raid"},
		{a: 62, b: 62, name: "pve_raid"},
		{a: 64, b: 64, name: "pve_raid"},
		{a: 66, b: 66, name: "pve_raid"},
		{a: 68, b: 68, name: "pve_raid"},
		{a: 70, b: 70, name: "pve_raid"},
		{a: 72, b: 72, name: "s8256_pve_raid_planetoid"},
		{a: 74, b: 74, name: "s8256_pve_raid_planetoid"},
		{a: 76, b: 76, name: "s8256_pve_raid_planetoid"},
		{a: 78, b: 78, name: "s8256_pve_raid_planetoid"},
		{a: 80, b: 80, name: "s8256_pve_raid_planetoid"},
		{a: 82, b: 82, name: "s8256_pve_raid_planetoid"},
		{a: 84, b: 84, name: "s8256_pve_raid_planetoid"},
		{a: 86, b: 86, name: "s8256_pve_raid_planetoid"},
		{a: 88, b: 88, name: "s8256_pve_raid_planetoid"},
		{a: 90, b: 90, name: "s8256_pve_raid_planetoid"},
		{a: 92, b: 92, name: "s8256_pve_raid_planetoid"},
		{a: 94, b: 94, name: "s8256_pve_raid_planetoid"},
		{a: 96, b: 96, name: "pve_raid"},
		{a: 98, b: 98, name: "pve_raid"},
		{a: 100, b: 100, name: "pve_raid"},
		{a: 102, b: 102, name: "pve_raid"},
		{a: 104, b: 104, name: "pve_raid"},
		{a: 106, b: 106, name: "pve_raid"},
		{a: 108, b: 108, name: "pve_raid"},
		{a: 110, b: 110, name: "pve_raid"},
		{a: 112, b: 112, name: "pve_raid"},
		{a: 114, b: 114, name: "pve_raid"},
		{a: 116, b: 116, name: "pve_raid"},
		{a: 118, b: 118, name: "pve_raid"},
		{a: 120, b: 120, name: "pve_raid_waterharvest"},
		{a: 122, b: 122, name: "pve_raid_waterharvest"},
		{a: 124, b: 124, name: "pve_raid_waterharvest"},
		{a: 126, b: 126, name: "pve_raid_waterharvest"},
		{a: 128, b: 128, name: "pve_raid_waterharvest"},
		{a: 130, b: 130, name: "pve_raid_waterharvest"},
		{a: 132, b: 132, name: "pve_raid_waterharvest"},
		{a: 134, b: 134, name: "pve_raid_waterharvest"},
		{a: 136, b: 136, name: "pve_raid_waterharvest"},
		{a: 138, b: 138, name: "pve_raid_waterharvest"},
		{a: 140, b: 140, name: "pve_raid_waterharvest"},
		{a: 142, b: 142, name: "pve_raid_waterharvest"},
		{a: 144, b: 144, name: "pve_raid_waterharvest"},
		{a: 146, b: 146, name: "pve_raid"},
		{a: 148, b: 148, name: "s8256_pve_raid_planetoid"},
		{a: 150, b: 150, name: "pve_raid_waterharvest"},
		{a: 152, b: 152, name: "pve_raid"},
		{a: 154, b: 154, name: "s8256_pve_raid_planetoid"},
		{a: 156, b: 156, name: "pve_raid_waterharvest"},
		{a: 158, b: 158, name: "pve_raid"},
		{a: 160, b: 160, name: "s8256_pve_raid_planetoid"},
		{a: 162, b: 162, name: "pve_raid_waterharvest"},
		{a: 164, b: 164, name: "pve_raid"},
		{a: 166, b: 166, name: "s8256_pve_raid_planetoid"},
	}
	bw.WriteBeUint32(uint32(len(data)))
	for _, v := range data {
		bw.WriteBeUint32(v.a)
		bw.WriteBeUint32(v.b)
		bw.WriteCString(v.name)
	}
	return nil
}

func (req *ac_load_initial_player_data_req) response() ([]byte, error) {
	resp := make([]byte, 0, 250000)
	bw := bitwriter.NewWriter(resp)
	bw.WriteBeUint16(uint16(types.AC_LOAD_INITIAL_PLAYER_DATA))

	bw.WriteBeUint64(1) // Account Revision
	bw.WriteBeUint32(1) // Format version
	bw.WriteBool(false) // idk
	bw.WriteBeUint32(1) // head account field, jk idk
	bw.WriteBool(false) // idk
	bw.WriteBeUint64(0) // moar zeros
	bw.WriteCString("")
	bw.WriteBeUint32(0) // Steam DLC count
	bw.WriteBeUint32(0) // Gaijin DLC count
	bw.WriteBeUint32(0) // Owned DLCs count
	if err := pveModifiers(bw); err != nil {
		return nil, err
	}
	// Unnamed/Empty brawl schedule thingie
	bw.WriteBeUint32(1)
	bw.WriteBeUint32(17)
	bw.WriteBeUint32(18)
	bw.WriteCString("")
	brawlschedule.BwScmdBrawlSchedule(bw)
	if err := rewardschedule.BwScmdRewardSchedule(bw); err != nil {
		return nil, err
	}
	if err := pveschedule.BwScmdPveSchedule(bw); err != nil {
		return nil, err
	}
	bw.WriteBeUint32(0)   // idk
	bw.WriteBool(false)   // idk
	bw.BwWriteByte(17)    // max unlocked vessel rank
	bw.BwWriteByte(19)    // account rank
	bw.WriteBeInt32(1480) // account exp
	bw.WriteBool(false)   // idk
	bw.WriteBeUint32(0)   // leading_advert (variantdict, send as empty for now)
	if err := unlimPveMissionLevel(bw); err != nil {
		return nil, err
	}
	// idk
	bw.WriteBeInt32(1)
	bw.WriteBeInt32(1)

	bw.WriteBool(false)
	bw.WriteBeInt32(-1)
	bw.WriteBeUint64(0)
	if err := leagueforbiddenequipment.BwScmdLeagueFrobiddenEquipment(bw); err != nil {
		return nil, err
	}
	battlepassactivation.BwScmdBattlePassActivation(bw)
	battlePassPlayerData(bw)
	bw.WriteBeInt32(50) // idk
	bw.WriteBool(true)
	if err := variantdict.BwMarshal(bw, getGamePlayMap()); err != nil {
		return nil, err
	}
	// SCMD_ZONES_WITH_DISABLED_QUESTS (65bytes)
	bw.BwWriteByte(0)
	for range 8 {
		bw.WriteBeUint64(0)
	}
	return bw.ReturnSlice(), nil
}

func handle_ac_load_initial_player_data(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_load_initial_player_data_req{}
	if err := req.parse(body[2:]); err != nil {
		// conn.Write(Failure)
	}
	slog.Info("Parsed request", "Request", req)
	resp, err := req.response()
	if err != nil {
		slog.Error("Failed to create ac_load_initial_player_data response", "error", err)
		// XXX: Disconnect the client
		return
	}
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
