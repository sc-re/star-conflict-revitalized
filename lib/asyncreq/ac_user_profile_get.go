package asyncreq

import (
	"log/slog"
	"starconflict/lib/bitreader"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
)

type userProfileFieldMask uint32

const (
	upfState            userProfileFieldMask = 1
	upfClanId           userProfileFieldMask = 1 << 1
	upfGeneralStats     userProfileFieldMask = 1 << 2
	upfVesselsRankStats userProfileFieldMask = 1 << 3
	upfAchievements     userProfileFieldMask = 1 << 4
	upfMedals           userProfileFieldMask = 1 << 5
	upfTitles           userProfileFieldMask = 1 << 6
	upfAvatars          userProfileFieldMask = 1 << 7
	upfMottos           userProfileFieldMask = 1 << 8
	upfAtlas            userProfileFieldMask = 1 << 9
)

type userProfileGeneralStat uint8

const (
	upgsKarma userProfileGeneralStat = iota
	upgsFirstLogin
	upgsPlayedTime
	upgsBattleTime
	upgsMaxSessionLength
	upgsNumBattles
	upgsNumVictories
	upgsWinStreak
	upgsNumPveBattles
	upgsRating
	upgsNumPlayersKilled
	upgsMaxPlayersKilledInOneBattle
	upgsNumAssists
	upgsMaxAssistsInOneBattle
	upgsTotalDamageDone
	upgsTotalHealingDone
	upgsVitalPointDamageDone
	upgsAccountRank
	upgsVesselsTotal
	upgsVesselsPremium
	upgsVesselsAtMaxLvl
	upgsThumbsUpAlly
	upgsThumbsUpEnemy
	upgsLeagueTeamId
	upgsLeagueTeamMajor
	upgsLeagueTeamRank
	upgsMilitaryRank
	upgsMilitaryExp
	upgsFactionRepEmpire
	upgsFactionRepFederation
	upgsFactionRepJericho
	upgsFactionRepEnclave
	upgsFactionRepCyber2
)

type userState byte

const (
	userStateNotAvailable userState = iota
	userStateOffline
	userStateOnline
)

type profileRequest struct {
	uid  uint64
	flag uint32
}

type ac_user_profile_get_req struct {
	profilesRequets []profileRequest
}

func upfStateWriter(bw *bitwriter.Writer, uid uint64) {
	bw.BwWriteByte(byte(userStateOnline)) // 0
	bw.WriteBeUint64(0)                   // timestamp, state last changed
}

func upfClanIdWriter(bw *bitwriter.Writer, uid uint64) {
	bw.WriteBeUint64(0) // clan id
}

func upfGeneralStatsWriter(bw *bitwriter.Writer, uid uint64) {
	for range 33 {
		bw.WriteBeUint64(0) // userProfileGeneralStat
	}
}

func upfVesselsRankStatsWriter(bw *bitwriter.Writer, uid uint64) {
	slog.Error("Nobody has implemented upfVesselRankStatsWriter yet")
}

func upfAchievementsWriter(bw *bitwriter.Writer, uid uint64) {
	slog.Error("Nobody has implemented upfAchievementsWriter yet")
}

func upfMedalsWriter(bw *bitwriter.Writer, uid uint64) {
	for range 62 {
		bw.WriteBeUint32(0)
	}
}

func upfTitlesWriter(bw *bitwriter.Writer, uid uint64) {
	bw.WriteBeUint16(0) // active title
	for range 0x180 {
		bw.WriteBool(false) // true if a certain title is owned by the account
	}
}

func upfAvatarsWriter(bw *bitwriter.Writer, uid uint64) {
	bw.WriteCString("Avatar_73")
	bw.WriteBeUint16(1) // Unlocked Avatar Count
	bw.WriteCString("Avatar_73")
}

func upfMottosWriter(bw *bitwriter.Writer, uid uint64, session *session.Session) {
	mottos := []string{}
	if err := session.Db.Select(&mottos, "SELECT motto FROM mottos WHERE uid=$1", uid); err != nil {
		mottos = []string{}
	}
	activeMotto := ""
	session.Db.Get(&activeMotto, "SELECT activeMotto FROM user WHERE uid=$1", uid)

	bw.WriteCString(activeMotto)
	bw.WriteBeUint16(uint16(len(mottos))) // Unlocked Motto Count
	for _, motto := range mottos {
		bw.WriteCString(motto)
	}
}

func upfAtlasWriter(bw *bitwriter.Writer, uid uint64) {
	bw.WriteBeInt32(1234) // Experience Points
	bw.WriteBool(true)    // module progress dict, true = no dict
	bw.WriteBool(true)    // vessel progess dict, true = no dict
}

func readVarUint(br *bitreader.Reader) (uint32, error) {
	b, err := br.ReadBool()
	if err != nil {
		return 0, err
	}
	if !b {
		i, err := br.ReadByte()
		return uint32(i), err
	}
	b, err = br.ReadBool()
	if err != nil {
		return 0, err
	}
	if !b {
		i, err := br.ReadBeUint16()
		return uint32(i), err
	}
	return br.ReadBeUint32()
}

func (req *ac_user_profile_get_req) parse(body []byte) error {
	br := bitreader.NewReader(body)
	count, err := br.ReadBeUint32()
	if err != nil {
		return err
	}
	for range count {
		uid, err := br.ReadBeUint64()
		if err != nil {
			return err
		}
		_ = uid
		flag, err := readVarUint(br)
		if err != nil {
			return err
		}
		req.profilesRequets = append(req.profilesRequets, profileRequest{uid, flag})
	}
	return nil
}

func (req *ac_user_profile_get_req) response(session *session.Session) []byte {
	bw := bitwriter.NewWriter(make([]byte, 0, 2500))
	bw.WriteBeUint16(uint16(types.AC_USER_PROFILE_GET))
	bw.WriteBeUint16(uint16(len(req.profilesRequets))) // maybe only fill in in the end so we can drop uids if they don't exist, throw errors, idk?

	for _, v := range req.profilesRequets {
		bw.WriteBeUint64(v.uid)
		bw.WriteBool(true)
		bw.WriteBool(true)
		bw.WriteBeUint32(v.flag) // XXX: Check if there is a good reason to not just always send as 32bit int
		if v.flag&uint32(upfState) != 0 {
			slog.Debug("Writing upfState", "uid", v.uid)
			upfStateWriter(bw, v.uid)
		}
		if v.flag&uint32(upfClanId) != 0 {
			slog.Debug("Writing upfClanId", "uid", v.uid)
			upfClanIdWriter(bw, v.uid)
		}
		if v.flag&uint32(upfGeneralStats) != 0 {
			slog.Debug("Writing upfGeneralStats", "uid", v.uid)
			upfGeneralStatsWriter(bw, v.uid)
		}
		if v.flag&uint32(upfVesselsRankStats) != 0 {
			slog.Debug("Writing upfVesselsRankStats", "uid", v.uid)
			upfVesselsRankStatsWriter(bw, v.uid)
		}
		if v.flag&uint32(upfAchievements) != 0 {
			slog.Debug("Writing upfAchievements", "uid", v.uid)
			upfAchievementsWriter(bw, v.uid)
		}
		if v.flag&uint32(upfMedals) != 0 {
			slog.Debug("Writing upfMedals", "uid", v.uid)
			upfMedalsWriter(bw, v.uid)
		}
		if v.flag&uint32(upfTitles) != 0 {
			slog.Debug("Writing upfTitles", "uid", v.uid)
			upfTitlesWriter(bw, v.uid)
		}
		if v.flag&uint32(upfAvatars) != 0 {
			slog.Debug("Writing upfAvatars", "uid", v.uid)
			upfAvatarsWriter(bw, v.uid)
		}
		if v.flag&uint32(upfMottos) != 0 {
			slog.Debug("Writing upfMottos", "uid", v.uid)
			upfMottosWriter(bw, v.uid, session)
		}
		if v.flag&uint32(upfAtlas) != 0 {
			slog.Debug("Writing upfAtlas", "uid", v.uid)
			upfAtlasWriter(bw, v.uid)
		}
	}

	return bw.ReturnSlice()
}

func handle_ac_user_profile_get(body []byte, seq uint16, seqRet uint16, session *session.Session) {
	req := ac_user_profile_get_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_user_profile_get_req", "error", err)
	}
	slog.Debug("Fetching profile data for profiles", "len", len(req.profilesRequets), "request", req)
	resp := req.response(session)
	session.Conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
