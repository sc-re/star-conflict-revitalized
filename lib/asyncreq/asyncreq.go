package asyncreq

import (
	"encoding/binary"
	"log/slog"
	"starconflict/lib/protocol"
	"starconflict/lib/session"
	"starconflict/lib/types"
	"time"
)

func SendAcPlayerPush(session *session.Session) {
	time.Sleep(7 * time.Second)
	Send_ac_player_autogen_inventory(session)
	Send_ac_player_credentials(session)
	Send_ac_player_credits(session)
	Send_ac_player_inventory(session.Conn)
	Send_ac_player_stats(session.Conn)
	Send_ac_player_vessels(session.Conn)
	Send_ac_get_fb_token(session.Conn)
	Send_ac_battle_slots(session.Conn)
	Send_ac_premium_info(session.Conn)
	Send_ac_get_userdata(session)
	Send_ac_rewarded_tutorials(session.Conn)
	Send_ac_obtain_referal_key(session.Conn)
	Send_ac_account_auras(session.Conn)
	Send_ac_get_blueprints(session.Conn)
	Send_ac_get_craft_resources(session.Conn)
	Send_ac_get_visited_free_space_zones(session.Conn)
	Send_ac_vessel_free_custom_elements(session.Conn)
}

func HandleAsyncReq(hdr *protocol.Header, body []byte, seq uint16, session *session.Session) {
	actype := types.AsyncReqType(binary.BigEndian.Uint16(body))
	switch actype {
	case types.AC_WELCOME_MSG:
		handle_ac_welcome_msg(body, seq, hdr.Sequence, session.Conn)
	case types.AC_MOTD:
		handle_ac_motd(body, seq, hdr.Sequence, session.Conn)
	case types.AC_SERVER_INFO:
		handle_ac_server_info(body, seq, hdr.Sequence, session.Conn)
	case types.AC_UNIVERSE_GET:
		handle_ac_universe_get(body, seq, hdr.Sequence, session.Conn)
	case types.AC_ZONES_LUA_ACTIVE_EVENTS_UPDATE:
		handle_ac_zones_lua_active_events_update(body, seq, hdr.Sequence, session.Conn)
	case types.AC_LOAD_INITIAL_PLAYER_DATA:
		handle_ac_load_initial_player_data(body, seq, hdr.Sequence, session.Conn)
	case types.AC_LOBBY_INFO:
		handle_ac_lobby_info(body, seq, hdr.Sequence, session.Conn)
	case types.AC_CLAN_REQUEST_DESC:
		handle_ac_clan_request_desc(body, seq, hdr.Sequence, session.Conn)
	case types.AC_FRIENDS_SEND_REQUEST:
		handle_ac_friends_send_request(body, seq, hdr.Sequence, session.Conn)
	case types.AC_LEADERBOARD_GET_DESCS:
		handle_ac_leaderboard_get_desc(body, seq, hdr.Sequence, session.Conn)
	case types.AC_TEACHING_LIST:
		handle_ac_teaching_list(body, seq, hdr.Sequence, session.Conn)
	case types.AC_QUESTS:
		handle_ac_quests(body, seq, hdr.Sequence, session.Conn)
	case types.AC_ADVENTURES:
		handle_ac_adventures(body, seq, hdr.Sequence, session.Conn)
	case types.AC_SQUAD_INFO:
		handle_ac_squad_info(body, seq, hdr.Sequence, session.Conn)
	case types.AC_LEAGUE_TEAM_INFO:
		handle_ac_league_team_info(body, seq, hdr.Sequence, session.Conn)
	case types.AC_MAIL_GET:
		handle_ac_mail_get(body, seq, hdr.Sequence, session.Conn)
	case types.AC_USER_NOTES:
		handle_ac_user_notes(body, seq, hdr.Sequence, session.Conn)
	case types.AC_USER_PROFILE_GET:
		handle_ac_user_profile_get(body, seq, hdr.Sequence, session.Conn)
	case types.AC_GET_USERDATA:
		handle_ac_get_userdata(body, seq, hdr.Sequence, session)
	case types.AC_SET_USERDATA:
		handle_ac_set_userdata(body, seq, hdr.Sequence, session)
	default:
		slog.Warn("Unhandled AsyncReq", "AsyncType", actype)
	}
}
