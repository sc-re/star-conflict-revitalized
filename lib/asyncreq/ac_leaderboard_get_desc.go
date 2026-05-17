package asyncreq

import (
	"log/slog"
	"net"
	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
	"starconflict/lib/variantdict"
)

type ac_leaderboard_get_desc_req struct{}

type rewardParams struct {
	StringId     string
	TokenCredits float32
	Tid          float32
}

type reward struct {
	Type   int32
	Params rewardParams
	PosMin int32
	PosMax int32
}

type leaderboardEntries struct {
	Name            string
	EntityType      int32
	Dir             int32
	ExpiresAt       uint64
	RenewalInterval int32
	ExpAction       int32
	LastDecay       uint64
	DecayInterval   int32
	DecayPower      float32
	Rewards         map[string]reward
}

func (req *ac_leaderboard_get_desc_req) parse(body []byte) error {
	return nil
}

func defaultLeaderboardDescs() []leaderboardEntries {
	entries := []leaderboardEntries{
		{Name: "player_prestige", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_pvp_effectiveness", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.10000000149011612, Rewards: map[string]reward{}},
		{Name: "player_pve_effectiveness", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.10000000149011612, Rewards: map[string]reward{}},
		{Name: "player_battles", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_wins", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_kills", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_karma_good", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_karma_bad", EntityType: 0, Dir: 1, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_goals", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_hunger", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_pressure", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_time_pve_raid", EntityType: 0, Dir: 1, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_time_pve_raid_waterharvest_t3", EntityType: 0, Dir: 1, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_time_pve_raid_waterharvest_t4", EntityType: 0, Dir: 1, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_time_pve_raid_waterharvest_t5", EntityType: 0, Dir: 1, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_time_pve_planetoidraid", EntityType: 0, Dir: 1, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "unlim_pve_level_planet_war_waves_T1", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_wavePve_wave", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "aprilcontrolgwin", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_kills_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779321600000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop1", TokenCredits: 400.0}, PosMin: 1, PosMax: 1}, "1": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop2", TokenCredits: 200.0}, PosMin: 2, PosMax: 2}, "2": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop3", TokenCredits: 100.0}, PosMin: 3, PosMax: 3}}},
		{Name: "player_pvp_eff_max_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779368400000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop1", TokenCredits: 400.0}, PosMin: 1, PosMax: 1}, "1": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop2", TokenCredits: 200.0}, PosMin: 2, PosMax: 2}, "2": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop3", TokenCredits: 100.0}, PosMin: 3, PosMax: 3}}},
		{Name: "player_pve_eff_max_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779368400000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop1", TokenCredits: 400.0}, PosMin: 1, PosMax: 1}, "1": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop2", TokenCredits: 200.0}, PosMin: 2, PosMax: 2}, "2": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop3", TokenCredits: 100.0}, PosMin: 3, PosMax: 3}}},
		{Name: "player_pve_win_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779321600000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop1", TokenCredits: 400.0}, PosMin: 1, PosMax: 1}, "1": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop2", TokenCredits: 200.0}, PosMin: 2, PosMax: 2}, "2": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop3", TokenCredits: 100.0}, PosMin: 3, PosMax: 3}}},
		{Name: "player_eff_rating_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779321600000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop1eff", Tid: 31.0}, PosMin: 1, PosMax: 3}, "1": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop2eff", Tid: 32.0}, PosMin: 4, PosMax: 20}, "2": {Type: 0, Params: rewardParams{StringId: "LeaderboardTop3eff", Tid: 33.0}, PosMin: 21, PosMax: 50}}},
		{Name: "cm_pvp_effectiveness", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.10000000149011612, Rewards: map[string]reward{}},
		{Name: "cm_pve_effectiveness", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.10000000149011612, Rewards: map[string]reward{}},
		{Name: "cm_pressure", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.10000000149011612, Rewards: map[string]reward{}},
		{Name: "cm_military_rank", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "cm_war_damage", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "cm_war_losses", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "CLS_PVP_RATING", EntityType: 1, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778972400000, DecayInterval: 60, DecayPower: 0.009999999776482582, Rewards: map[string]reward{}},
		{Name: "CLS_PVE_RATING", EntityType: 1, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778972400000, DecayInterval: 60, DecayPower: 0.009999999776482582, Rewards: map[string]reward{}},
		{Name: "CLS_WAR_DAMAGE", EntityType: 1, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "CLS_ZONE_CAPTURE", EntityType: 1, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "CLS_WAR_RATING", EntityType: 1, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "CLS_WAR_RATING_MAX", EntityType: 1, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "league_rating", EntityType: 4, Dir: 0, ExpiresAt: 1782345600000, RenewalInterval: 129600, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTopLgs1", Tid: 201.0}, PosMin: 1, PosMax: 8}}},
		{Name: "major_league_rating", EntityType: 0, Dir: 0, ExpiresAt: 1782345600000, RenewalInterval: 129600, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{"0": {Type: 0, Params: rewardParams{StringId: "LeaderboardTopLgs1", Tid: 201.0}, PosMin: 1, PosMax: 8}}},
		{Name: "event_faction_war_empire", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "event_faction_war_federation", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "event_faction_war_jericho", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "event_faction_war_pirate", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_zm_win_humans", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_zm_win_zombies", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "autogen_dm", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "three_teams_dm_kills", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_cd_win", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "one_against_all_boss_win", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "one_against_all_hunter_eff", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_kills_marathon", EntityType: 0, Dir: 0, ExpiresAt: 18446744073709551615, RenewalInterval: 0, ExpAction: 2, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "unlim_pve_level_magnificent_seven", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "unlim_pve_level_stealth", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "unlim_pve_level_rescue_pirates_base", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "cm_effectiveness", EntityType: 2, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.009999999776482582, Rewards: map[string]reward{}},
		{Name: "player_effectiveness_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779321600000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_battles_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779321600000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_wins_weekly", EntityType: 0, Dir: 0, ExpiresAt: 1779321600000, RenewalInterval: 10080, ExpAction: 0, LastDecay: 0, DecayInterval: 0, DecayPower: 0.0, Rewards: map[string]reward{}},
		{Name: "player_effectiveness", EntityType: 0, Dir: 0, ExpiresAt: 0, RenewalInterval: 0, ExpAction: 0, LastDecay: 1778889600000, DecayInterval: 1440, DecayPower: 0.009999999776482582, Rewards: map[string]reward{}},
	}
	return entries
}

func (req *ac_leaderboard_get_desc_req) response() ([]byte, error) {
	bw := bitwriter.NewWriter(make([]byte, 0, 12000))
	bw.WriteBeUint16(uint16(types.AC_LEADERBOARD_GET_DESCS))
	leaderboards := defaultLeaderboardDescs()
	bw.WriteBeUint32(uint32(len(leaderboards)))
	for _, v := range leaderboards {
		if err := variantdict.BwMarshal(bw, v); err != nil {
			return nil, err
		}
	}
	return bw.ReturnSlice(), nil
}

func handle_ac_leaderboard_get_desc(body []byte, seq uint16, seqRet uint16, conn net.Conn) {
	req := ac_leaderboard_get_desc_req{}
	if err := req.parse(body[2:]); err != nil {
		slog.Error("Failed to parse ac_leaderboard_get_desc_req", "error", err)
	}
	resp, err := req.response()
	if err != nil {
		slog.Error("Failed to create ac_leaderboard_get_desc", "error", err)
		// XXX: Send a proper status message
		return
	}
	conn.Write(protocol.MakeMessage(types.CSCMD_ASYNC_REQ, seq, seqRet, resp))
}
