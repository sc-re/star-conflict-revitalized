package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"io"

	"net"
	"slices"

	"starconflict/lib/bitwriter"
	"starconflict/lib/protocol"
	"starconflict/lib/types"

	"github.com/jmoiron/sqlx"
	"github.com/monnand/dhkx"
)

type authMethod uint8

const (
	AM_UNKNOWN authMethod = iota
	AM_INTERNAL
	AM_YUPLAY
	AM_STEAM
	AM_ARC
	AM_FACEBOOK
	AM_VKONTAKTE
	AM_FAKE
	AM_BOT
	AM_NUM
)

type authRequest struct {
	version    string
	hash       uint32
	authMethod authMethod
	email      string
	gBmodP     *dhkx.DHKey
	cipherText [256]byte
	password   []byte
	platform   string
	machineId  string
}

type user struct {
	Uid      uint64
	Nickname string
	Zone     byte
	Password string
	Mail     string
}

func readAuthMethod(reader *bytes.Reader) (authMethod, error) {
	u, err := protocol.ReadUint8(reader)
	return authMethod(u), err
}

func readgBmodp(reader *bytes.Reader) (*dhkx.DHKey, error) {
	buf := make([]byte, 128)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	key := dhkx.NewPublicKey(buf)
	return key, nil
}

func (auth *authRequest) parseCCmdAuthRequest(body []byte) error {
	r := bytes.NewReader(body)
	err := *new(error)
	auth.version, err = protocol.ReadCString(r)
	if err != nil {
		return err
	}
	auth.hash, err = protocol.ReadUint32(r)
	if err != nil {
		return err
	}
	auth.authMethod, err = readAuthMethod(r)
	if err != nil {
		return err
	}
	if auth.authMethod != AM_INTERNAL {
		return errors.New("Unhandled authenticationMethod")
	}
	auth.email, err = protocol.ReadCString(r)
	if err != nil {
		return err
	}
	if _, err := r.ReadByte(); err != nil {
		return err
	}
	auth.gBmodP, err = readgBmodp(r)
	if err != nil {
		return err
	}
	if _, err := io.ReadFull(r, auth.cipherText[:]); err != nil {
		return err
	}

	return nil
}

func (auth *authRequest) decryptPassword(group *dhkx.DHGroup, key *dhkx.DHKey) error {
	buf := make([]byte, 256)
	sharedKey, err := group.ComputeKey(auth.gBmodP, key)
	if err != nil {
		return err
	}
	iv := sha256.Sum256([]byte(auth.email))
	block, err := aes.NewCipher(sharedKey.Bytes()[:16])
	if err != nil {
		return err
	}
	mode := cipher.NewCBCDecrypter(block, iv[:16])
	mode.CryptBlocks(buf, auth.cipherText[:])
	off := 0
	for k, v := range buf {
		if v == 0x0 {
			off = k
			break
		}
	}
	auth.password = buf[:off]
	return nil
}

func (auth *authRequest) checkHash() bool {
	// XXX: implement
	_ = auth.hash
	return true
}

func (auth *authRequest) checkClientVersion() bool {
	return auth.version == "1.14.09.166666"
}

func Authenticate(body []byte, db *sqlx.DB, key *dhkx.DHKey, group *dhkx.DHGroup) (bool, uint64, types.MasterServerDisconnectReason, error) {
	auth := authRequest{}
	err := auth.parseCCmdAuthRequest(body)
	if err != nil {
		return false, 0, types.MasterServerDisconnectReason(0xff), err
	}
	err = auth.decryptPassword(group, key)
	user := user{}
	if err := db.Get(&user, "SELECT * FROM user WHERE mail=$1", auth.email); err != nil {
		return false, 0, types.DR_INVALID_LOGIN, err
	}
	valid := checkPassword(auth.password, user.Password)
	if !valid {
		return false, 0, types.DR_INVALID_LOGIN, err
	}
	valid = auth.checkClientVersion()
	if !valid {
		return false, 0, types.DR_BAD_CLIENT_VERSION, nil
	}
	valid = auth.checkHash()
	if !valid {
		return false, 0, types.DR_BAD_CLIENT_VERSION, nil
	}
	return true, user.Uid, 0, nil
}

func CreateClientChallenge() (*dhkx.DHKey, *dhkx.DHGroup) {
	// #LastMillenium, everything created after Y200 sucks
	// #1024bitsOughtToBeEnoughForEveryone
	// rfc2409 Second Oakley Group
	group, _ := dhkx.GetGroup(2)
	priv, _ := group.GeneratePrivateKey(nil)
	return priv, group
}

func SendAuthAck(conn net.Conn, db *sqlx.DB, seq uint16, seqRet uint16, uid uint64) error {
	// TODO: Actually generate some tokens and store them somewhere
	idkToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	token := uint64(10)

	var user user
	if err := db.Get(&user, "SELECT nickname,uid,zone FROM user WHERE uid=$1", uid); err != nil {
		return err
	}

	bw := bitwriter.NewWriter(make([]byte, 0, 86+len(user.Nickname)))
	bw.WriteBeUint64(user.Uid)
	bw.WriteBeUint64(token)
	bw.WriteCString(user.Nickname)
	bw.WriteCString(idkToken)
	bw.WriteBool(false)
	bw.WriteBeUint32(uint32(user.Zone))
	_, err := conn.Write(protocol.MakeMessage(types.SCMD_AUTH_ACK, seq, seqRet, bw.ReturnSlice()))
	return err
}

func SendChallenge(conn net.Conn, key *dhkx.DHKey, seq uint16) error {
	pub := key.Bytes()
	length := []byte{byte(len(pub))}
	body := slices.Concat(length, pub)
	_, err := conn.Write(protocol.MakeMessage(types.SCMD_AUTH_REQ, seq, 0, body))
	return err
}
