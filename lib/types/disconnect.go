package types

type MasterServerDisconnectReason uint8

const (
	DR_ACCESS_DENIED                       MasterServerDisconnectReason = 37
	DR_ACCOUNT_ATTACHMENT_CONFLICT         MasterServerDisconnectReason = 54
	DR_ARC_AUTH_FAILED                     MasterServerDisconnectReason = 58
	DR_ATTACHED_TO_ANOTHER_ARC_ACCT        MasterServerDisconnectReason = 59
	DR_ATTACHED_TO_ANOTHER_STEAM_ACCT      MasterServerDisconnectReason = 48
	DR_BAD_CLIENT_VERSION                  MasterServerDisconnectReason = 35
	DR_BANNED                              MasterServerDisconnectReason = 40
	DR_CLIENT_QUIT                         MasterServerDisconnectReason = 41
	DR_ERROR_CHANNEL_CAP                   MasterServerDisconnectReason = 4
	DR_ERROR_FAILED_TO_CONNECT             MasterServerDisconnectReason = 3
	DR_ERROR_INITIATING_CONNECTION         MasterServerDisconnectReason = 2
	DR_GAME_SESSION_ABANDONED              MasterServerDisconnectReason = 45
	DR_GAME_SESSION_FATAL                  MasterServerDisconnectReason = 44
	DR_GAME_SESSION_FINISHED               MasterServerDisconnectReason = 43
	DR_GENERIC_ERROR                       MasterServerDisconnectReason = 1
	DR_INTERNAL_ERROR                      MasterServerDisconnectReason = 32
	DR_INVALID_LOGIN                       MasterServerDisconnectReason = 36
	DR_IN_HUGE_DEBT                        MasterServerDisconnectReason = 51
	DR_KEEPALIVE_TIMEOUT                   MasterServerDisconnectReason = 50
	DR_KICK                                MasterServerDisconnectReason = 39
	DR_MULTIPLE_LOGIN                      MasterServerDisconnectReason = 38
	DR_NETWORK_FAIL                        MasterServerDisconnectReason = 33
	DR_NETWORK_FAIL_OS                     MasterServerDisconnectReason = 5
	DR_NOT_ALLOWED_FOR_AUTO_GENERATED_ACCT MasterServerDisconnectReason = 52
	DR_NO_ACCOUNT_ATTACHED_TO_STEAM_LOGIN  MasterServerDisconnectReason = 49
	DR_NO_REASON                           MasterServerDisconnectReason = 0
	DR_SERVER_SHUTDOWN                     MasterServerDisconnectReason = 42
	DR_STEAM_AUTH_FAILED                   MasterServerDisconnectReason = 47
	DR_TEMPORARY_UNAVAILABLE               MasterServerDisconnectReason = 46
	DR_TOO_MANY_CLEINTS                    MasterServerDisconnectReason = 34
	DR_TOO_MANY_REQUESTS                   MasterServerDisconnectReason = 53
	DR_YUPLAY_2STEP                        MasterServerDisconnectReason = 55
	DR_YUPLAY_2STEPERROR                   MasterServerDisconnectReason = 57
	DR_YUPLAY_FROZEN                       MasterServerDisconnectReason = 56
)

/*
MasterServerDisconnectReason = {
       },
*/
