package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"

	"golang.org/x/crypto/argon2"
)

type phc struct {
	id              string
	version         int
	memoryCost      uint32
	timeCost        uint32
	parallelismCost uint8
	keyLength       uint32
	salt            []byte
	hash            []byte
}

// rfc9106 section-7.3
const (
	defaultTimeCost        = 1
	defaultMemoryCost      = 1024 * 1024 * 2 // 2GiB
	defaultParallelismCost = 4
	defaultSaltLength      = 16
	defaultKeyLength       = 32
)

var (
	formatString = "$%s$v=%d$m=%d,t=%d,p=%d,l=%d$%s$%s"
)

func HashPassword(password []byte) (string, error) {
	args := phc{
		id:              "argon2id",
		version:         argon2.Version,
		memoryCost:      defaultMemoryCost,
		timeCost:        defaultTimeCost,
		parallelismCost: defaultParallelismCost,
		keyLength:       defaultKeyLength,
	}
	if err := args.generateSalt(defaultSaltLength); err != nil {
		return "", err
	}

	args.hash = argon2.IDKey(password, args.salt, args.timeCost, args.memoryCost, args.parallelismCost, args.keyLength)
	return args.phcEncode(), nil
}

func checkPassword(password []byte, hash string) bool {
	args := phc{}
	if err := args.phcDecode(hash); err != nil {
		slog.Error("Failed to decode phc string", "error", err)
		return false
	}
	key := argon2.IDKey(password, args.salt, args.timeCost, args.memoryCost, args.parallelismCost, args.keyLength)
	return subtle.ConstantTimeCompare(key, args.hash) == 1
}

func (phc *phc) generateSalt(saltLength uint) error {
	phc.salt = make([]byte, saltLength)
	if _, err := rand.Read(phc.salt); err != nil {
		return err
	}
	return nil
}

func (phc *phc) phcEncode() string {
	salt := base64.RawStdEncoding.EncodeToString(phc.salt)
	hash := base64.RawStdEncoding.EncodeToString(phc.hash)

	ret := fmt.Sprintf(formatString, phc.id, phc.version, phc.memoryCost, phc.timeCost, phc.parallelismCost, phc.keyLength, salt, hash)
	return ret
}

func (phc *phc) phcDecode(phcString string) error {
	fields := strings.Split(phcString, "$")
	if len(fields) != 6 {
		return fmt.Errorf("invalid format of phcString")
	}
	phc.id = fields[1]
	if _, err := fmt.Sscanf(fields[2], "v=%d", &phc.version); err != nil {
		return err
	}
	if _, err := fmt.Sscanf(fields[3], "m=%d,t=%d,p=%d,l=%d", &phc.memoryCost, &phc.timeCost, &phc.parallelismCost, &phc.keyLength); err != nil {
		return err
	}
	err := *new(error)
	phc.salt, err = base64.RawStdEncoding.Strict().DecodeString(fields[4])
	if err != nil {
		return err
	}
	phc.hash, err = base64.RawStdEncoding.Strict().DecodeString(fields[5])
	if err != nil {
		return err
	}
	return nil
}
