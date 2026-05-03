package protocol

import (
	"starconflict/lib/types"
)

func MakeDisconnectMessage(reason types.MasterServerDisconnectReason) []byte {
	ret := make([]byte, 12)
	ret[0], ret[1], ret[2], ret[3] = 0xff, 0xff, 0xff, 0xff
	ret[8] = uint8(reason)
	return ret
}
