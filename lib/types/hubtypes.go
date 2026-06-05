//go:generate stringer -type=HubCommandType

package types

type ServerType uint32

const (
	ST_INVALID ServerType = iota
	ST_LOAD_BALANCER
	ST_SHARD
	ST_MATCHMAKER
	ST_DEDICATED_SERVER
	ST_HUB
	ST_ADMIN_CONSOLE
	ST_GAME_REWARD_DISPENCER
	ST_AUTHENTICATOR
	ST_CHAT_SERVER
	ST_STAT_SERVER
	ST_CONFIG_SERVER
	ST_UNIVERSE
	ST_QUEST_SERVER
	ST_ZONE_COORDINATOR
	ST_ZONE_INSTANCE
	ST_LOBBY_SERVER
	ST_LEADERBOARD_SERVER
	ST_LEADERBOARD_SERVERv2
	ST_CLAN_SERVER
	ST_TOURNMANET
	ST_EARLY_MATCHMAKER
	ST_SCRIPTS_SERVER
	ST_ADVENTURE_SERVER
)

type HubCommandType uint16

const (
	INITIAL_REGISTER             HubCommandType = 0x00
	CL_HEARTBEAT                 HubCommandType = 0x03
	CL_ADDR_ADVERTISE            HubCommandType = 0x0b
	CC_CLOUD_CONFIG              HubCommandType = 0xab
	CL_REGISTER                  HubCommandType = 0xad
	CL_IDK                       HubCommandType = 0xac
	CC_ZI_PLAYER_CONNECTED       HubCommandType = 0x35 // Also gets acked and returned by Dedicated Server
	CC_REGISTER_PLAYER_LIST      HubCommandType = 0x36 // ditto
	CC_PLAYER_ACK                HubCommandType = 0x37
	CC_GAMEPLAY_PROPS_UPDATE     HubCommandType = 0x38
	CL_CLIENT_CAPABILITYIES      HubCommandType = 0x39
	CC_HEARTBEAT                 HubCommandType = 0x3b
	CC_RANK_UPDATE               HubCommandType = 0xb8
	CC_PROFILE_FIELD             HubCommandType = 0xb9
	CC_PROFILE_BIT               HubCommandType = 0xba
	CC_INVENTORY_BIT             HubCommandType = 0xbb
	CC_LOADOUT_UPDATE            HubCommandType = 0xbc
	CC_ZONE_REGION_UPDATE        HubCommandType = 0xc7
	CC_ZI_KICK_PLAYER_FROM_ZONES HubCommandType = 0xcb
	CC_DS_MULTIPLE_LOGIN         HubCommandType = 0xcc
	CC_ZI_REMOVE_PLAYER          HubCommandType = 0xcf
	CC_PLAYER_UPDATE             HubCommandType = 0xd0
	CC_ZI_REGISTER_PLAYER        HubCommandType = 0xd3
	CC_PROFILE_RELOAD            HubCommandType = 0xd4
	CC_PLAYER_STATE_LIST         HubCommandType = 0xd7
	CC_GAMEPLAY_RPC              HubCommandType = 0xd8
	CC_PLAYER_DOCK_ACTION        HubCommandType = 0xd9
	CC_PLAYER_DEPART             HubCommandType = 0xda
	CC_LOADOUT_UPDATE_5x         HubCommandType = 0xdb
	CC_LOOT_CONFIG               HubCommandType = 0xf2
	CC_GAMEPLAY_CONFIG           HubCommandType = 0xf4
	CC_PLAYER_NOTIFICATION       HubCommandType = 0xf8
	CC_PLAYER_REMOVED            HubCommandType = 0xfb
)
