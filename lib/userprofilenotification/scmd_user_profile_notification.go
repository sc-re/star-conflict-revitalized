package userprofilenotification

import (
	"encoding/binary"
	"log"
	"net"
	"starconflict/lib/protocol"
	"starconflict/lib/types"
)

// Here we see the notification varian of AC_USER_PROFILE_GET
// Should we ever send anything else but UPN_ONLINE_STATE for
// a uid that is not the current user? Nope. But can we? Yes!
// Why? We will never know.

type UserState byte

const (
	USER_STATE_NOT_AVAILABLE UserState = iota
	USER_STATE_ONLINE
	USER_STATE_OFFLINE
)

type UserProfileNotification byte

const (
	UPN_ONLINE_STATE UserProfileNotification = iota
	UPN_ACHIEVEMENT_UPDATE
	UPN_ACHIEVEMENT_UNLOCK
	upn_we_dont_know
	UPN_AVATAR_UNLOCK
	UPN_TITLE_UNLOCK
	UPN_CLEARANCE_SCORE_UPDATE
	upn_skill_issue
)

func SendUserProfileNotificationOnlineState(conn net.Conn, uid uint64, state UserState) {
	log.Printf("Resp: SendUserProfileNotificaitonOnlineState[uid=%v, state=%v]", uid, state)
	resp := make([]byte, 0, 10)
	resp = append(resp, byte(UPN_ONLINE_STATE))
	resp = binary.BigEndian.AppendUint64(resp, uid)
	resp = append(resp, byte(state))
	conn.Write(protocol.MakeMessage(types.SCMD_USER_PROFILE_NOTIFICATION, 0, 0, resp))
}

func send_user_profile_notification() {
	upn_type := byte(0)
	uid := 1438647
	flag := byte(0)
	_, _, _ = flag, uid, upn_type

}
