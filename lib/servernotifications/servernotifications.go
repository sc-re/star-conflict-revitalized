package servernotifications

import "github.com/samber/mo"

type SN_ATLAS_INIT struct {
	AtlasModulesNum int32
}

type SN_LOBBY_GROUP_CREATED struct {
	ownerUid uint64 `variantdict:"owner_uid"`
}

type SN_PLAYER_BOUGHT_PREMIUM_VESSEL struct {
	Def      string
	Uid      uint64
	NickName string
}

type SN_PLAYER_CRAFTED_VESSEL struct {
	Def      string
	Uid      uint64
	NickName string
}

type SN_VESSELS_AUTO_REFILLED struct {
	MunitionTransfered int32
	MunitionPurchased  int32
	Credits            int32
	GoldCredits        int32
	TokenCredits       int32
	EventCredits       int32
}

type SN_LIKE_ADDED struct {
	Giver    int32 // XXX: int32 has to be wrong?
	Receiver int32 // ditto, or maybe both 32 and 64 works? but why signed?
}

type SN_VK_GROUP_REWARD_PROMO struct {
	GoldCredits int32
	Credits     int32
	PremiumTime int32
	aura        string
}

type SN_GAME_REWARD struct {
	GameMode             int32
	IsLooser             int32
	IncAntibotCounter    bool
	UnlimPveMissionLevel int32  // probably only in pve
	PveMissionName       string // ditto
}

type SN_MAINTENANCE_COUNTDOWN struct {
	MinutesLeft int32
}

type SN_PREMIUM_ACCESS struct {
	PremiumAccess int32
}

type SN_NEW_LETTER struct {
	MailId uint64
}

type SN_REFEREE_ADDED struct {
	Uid uint64
}

type SN_REFEREE_BONUS_GOLD struct {
	Referee   uint64
	BonusGold int32
}

type SN_ADMIN_TASK_RESULT struct {
	Success int32
}

//Notification(SN_DAILY_LOGIN,
// {'tier': i32(2), 'auras': bag({}),
//  'bundles': bag(
//                 {'Bundle_Iridium_Overseer_T5': bag({'0': bag({'aura': str('Bundle_Daily_Pvp_open'), 'type': i32(5)}),
//                                                     '1': bag({'aura': str('Bundle_Daily_Raid_open'), 'type': i32(5)}),
//                                                     '2': bag({'aura': str('Bundle_Daily_Pve_open'), 'type': i32(5)}),
//                                                     '3': bag({'aura': str('Bundle_Iridium_Overseer_T4_1_open'), 'type': i32(5)}),
//                                                     '4': bag({'aura': str('Bundle_Iridium_Overseer_T4_2_open'), 'type': i32(5)}),
//                                                     '5': bag({'aura': str('Bundle_Iridium_Overseer_T5_5_open'), 'type': i32(5)})
//                                                    })
//                 }
//                )
// }
//)
// gamedata/shared/login_rewards.lua

type SN_DAILY_LOGIN struct {
	Tier    int32 // 2
	Auras   struct{}
	Bundles struct{}
}

type corpStruct struct {
	Tag  string
	Name string
}

type SN_ZONE_OWNER_CHANGED struct {
	Zid int32 // Zone ID
	Old corpStruct
	New corpStruct
}

// 1: 2.0, zone_name, level_transit_emp_rock_frontier_e_ow
// 2: 0.0, param, 40.0
type GameEventParam struct {
	EventType float32 `variantdict:"type"`
	Name      string
	Value     mo.Either[float32, string]
}

type SN_GAME_EVENT_STATE_CHANGED struct {
	MessageType int32  // 2?
	TextId      string // ow_PIRATE_EXPANSION_periodic
	Params      []GameEventParam
}

type SN_SHIP_QUEST_STARTED struct {
	QuestId int32
	TextId  int32
	Vids    []uint64 // Vehicle IDs (upto 6)
}
