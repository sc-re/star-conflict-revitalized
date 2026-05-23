package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

type ac_get_userdata_req struct{}

type marksElement struct {
	IsFrame    float32
	IsHpShield float32
	IsName     bool
	IsDist     float32
	IsSpecIco  bool
	IsRank     float32
}

// Maybe we should skip (de-)serialization and just store/read-as-is from/into the database
// XXX: Check if we can just send an empty initial response and have the client set this with ac_set_userdata
type ac_get_userdata_resp struct {
	HelpShown struct {
		Quests        float32
		Hangar        float32
		Hud           float32
		ShipSelection float32
		Skills        float32
		ShipsTree     float32
		Shop          float32
		Equipment     float32
	}
	BoostExpiredMsgShown float32
	HudEditor            struct {
		Configs   struct{} //... nested and nestend and moar nesting
		CurConfig float32
	}
	HudOptions struct {
		OnCircleMarksMode   string
		ScalableHealthBars  bool
		GameCursorCrossType string
		ShowGameplayInfo    float32
		GameCursorColor     string
		MarksElements       struct {
			brief_ally marksElement
			lock_ally  marksElement
			full_ally  marksElement
			brief      marksElement
			locks      marksElement
		}
	}
	ShowPremium   float32
	UcidToShow    uint64
	RaceWndShown  float32
	CurrencyShown struct {
		Free_synergy  bool
		Gold          bool
		Credits       bool
		Loaylity      float32
		Tokens        bool
		Event_credits float32
	}
	AutoRepair                float32
	UserBindings              string // XML Document
	TutorialSkipped           float32
	TutorialCompleted         float32
	Ver                       float32
	PremiumExpiredMsgShown    bool
	LastImportantActionId     uint64
	UnseenDefs                map[string]string // Only observed as empty
	HangarShwon               float32
	Region                    float32
	SpecialRaceShipsShown     float32
	Ny2018PromoShown          float32
	Ny2018PromoShown_v2       float32
	GoldAmmoMsgShown          float32
	RolesShown                map[string]float32 // role = 1-10
	AutoRefill                float32
	UiTutorialCompletedStages struct {
		ShipUpgrade   float32 `variantdict:"SHIP_UPGRADE"`
		GameTutorial  float32 `variantdict:"GAME_TUTORIAL"`
		ItemUpgrade   float32 `variantdict:"ITEM_UPGRADE"`
		RealFight     float32 `variantdict:"REAL_FIGHT"`
		PracticeFight float32 `variantdict:"PRACTICE_FIGHT"`
		ShipBuy       float32 `variantdict:"SHIP_BUY"`
		ItemEquip     float32 `variantdict:"ITEM_EQUIP"`
	}
	DeleteOvercapResources float32
}

func (req *ac_get_userdata_req) parse(body []byte) error {
	return nil
}

func (req *ac_get_userdata_req) response() []byte {
	return ac_get_userdata_default()
}

func Send_ac_get_userdata(conn net.Conn) {
	req := ac_get_userdata_req{}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, 0, 0, resp))
}

func handle_ac_get_userdata(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_get_userdata_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_get_userdata_req", "error", err)
	}
	resp := req.response()
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
