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
	r.ReadByte()
	auth.gBmodP, err = readgBmodp(r)
	if err != nil {
		return err
	}
	io.ReadFull(r, auth.cipherText[:])

	return nil
}

func checkPassword(password []byte, email string) (bool, error) {
	// XXX: Totally legit
	if email != "test@localhost" {
		return false, nil
	}
	if !bytes.Equal(password, []byte("test")) {
		return false, nil
	}
	return true, nil
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

func Authenticate(body []byte, key *dhkx.DHKey, group *dhkx.DHGroup) (bool, types.MasterServerDisconnectReason, error) {
	auth := authRequest{}
	err := auth.parseCCmdAuthRequest(body)
	if err != nil {
		return false, types.MasterServerDisconnectReason(0xff), err
	}
	err = auth.decryptPassword(group, key)
	valid, err := checkPassword(auth.password, auth.email)
	if !valid {
		return false, types.DR_INVALID_LOGIN, err
	}
	valid = auth.checkClientVersion()
	if !valid {
		return false, types.DR_BAD_CLIENT_VERSION, nil
	}
	valid = auth.checkHash()
	if !valid {
		return false, types.DR_BAD_CLIENT_VERSION, nil
	}
	return true, 0, nil
}

func CreateClientChallenge() (*dhkx.DHKey, *dhkx.DHGroup) {
	// #LastMillenium, everything created after Y200 sucks
	// #1024bitsOughtToBeEnoughForEveryone
	// rfc2409 Second Oakley Group
	group, _ := dhkx.GetGroup(2)
	priv, _ := group.GeneratePrivateKey(nil)
	return priv, group
}

func SendAuthAck(conn net.Conn, seq uint16, seqRet uint16) error {
	// TODO: get form DB
	nick := "Bruh"
	idkToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	uid := uint64(10)
	token := uint64(10)
	zone := uint32(1)
	// Actual code
	bw := bitwriter.NewWriter(make([]byte, 0, 86+len(nick)))
	bw.WriteBeUint64(uid)
	bw.WriteBeUint64(token)
	bw.WriteCString(nick)
	bw.WriteCString(idkToken)
	bw.WriteBool(false)
	bw.WriteBeUint32(zone)
	_, err := conn.Write(protocol.MakeMessage(types.SCMD_AUTH_ACK, seq, seqRet, bw.ReturnSlice()))
	return err
}

func SendChallenge(conn net.Conn, key *dhkx.DHKey, seq uint16) error {
	pub := key.Bytes()
	len := []byte{byte(len(pub))}
	body := slices.Concat(len, pub)
	_, err := conn.Write(protocol.MakeMessage(types.SCMD_AUTH_REQ, seq, 0, body))
	return err
}
