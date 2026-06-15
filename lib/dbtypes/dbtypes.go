package dbtypes

type Accounts struct {
	Uid                          uint64     `bson:"uid,omitempty"`
	Puid                         uint64     `bson:"puid,omitempty"`
	YuplayID                     int64      `bson:"yuplayID,omitempty"`
	SteamID                      int64      `bson:"steamID,omitempty"`
	YupPurchases                 struct{}   `bson:"yupPurchases"`
	SteamYupPurchases            struct{}   `bson:"steamYupPurchases"`
	Username                     string     `bson:"username,omitempty"`
	Password                     string     `bson:"password,omitempty"`
	NickName                     string     `bson:"nickName,omitempty"`
	AdmRole                      string     `bson:"admRole,omitempty"`
	Credits                      uint64     `bson:"credits,omitempty"`
	GoldCredits                  uint64     `bson:"goldCredits,omitempty"`
	TokenCredits                 uint64     `bson:"tokenCredits,omitempty"`
	VesselExpPool                uint32     `bson:"vesselExpPool"`
	Karma                        int32      `bson:"karma,omitempty"`
	AccountRank                  int32      `bson:"accountRank,omitempty"`
	MilitaryRank                 int32      `bson:"militaryRank,omitempty"`
	Elo                          int32      `bson:"elo,omitempty"`
	PremiumTill                  uint64     `bson:"premiumTill,omitempty"`
	Promo                        bool       `bson:"promo,omitempty"`
	Userdata                     []byte     `bson:"userdata,omitempty"`
	RaceRep                      struct{}   `bson:"raceRep,omitempty"`
	FactionRep                   struct{}   `bson:"factionRep,omitempty"`
	Talents                      []struct{} `bson:"talents,omitempty"`
	TalentsSets                  []struct{} `bson:"talentsSets,omitempty"`
	Items                        []struct{} `bson:"items,omitempty"`
	AutogenItems                 []struct{} `bson:"autogenItems,omitempty"`
	Vessels                      []struct{} `bson:"vessels,omitempty"`
	BSlots                       []struct{} `bson:"bSlots,omitempty"`
	Auras                        []struct{} `bson:"auras,omitempty"`
	Title                        string     `bson:"title,omitempty"`
	AcquiredTitles               []struct{} `bson:"acquiredTitles,omitempty"`
	Avatar                       string     `bson:"avatar,omitempty"`
	AcquiredAvatars              []struct{} `bson:"acquiredAvatars,omitempty"`
	Motto                        string     `bson:"motto,omitempty"`
	AcquiredMottos               []string   `bson:"acquiredMottos,omitempty"`
	Referrals                    []struct{} `bson:"referrals"`
	Referrer                     int64      `bson:"referrar"` // Uid of referrer
	UnlimPve_missionLevels       int64      `bson:"unlimPve_missionLeves,omitempty"`
	UnlimPve_playerLevel1        int64      `bson:"unlimPve_playerLevel1,omitempty"`
	UnlimPve_playerLevel2        int64      `bson:"unlimPve_PlayerLevel2,omitempty"`
	UnlimPve_PlayerBuffsDisabled bool       `bson:"unlimPve_PlayerBuffsDisabled,omitempty"`
	WavePve_availableWave        int64      `bson:"wavePve_availableWave,omitempty"`
	PlayerLuaData                struct{}   `bson:"playerLuaData,omitempty"`
	OutgoingTransactions         struct{}   `bson:"outgoingTransactions,omitempty"`
	IngoingTransactions          struct{}   `bson:"ingoingTransactions,omitempty"`
	SpaceStation                 int64      `bson:"spaceStation,omitempty"`
	ver                          int64      `bson:"ver,omitempty"`
}

type Counters struct {
	C    int64 `bson:"c"`
	Time int64 `bson:"time"`
}

type AccountStats struct {
	Uid uint64 `bson:"uid"`
}
