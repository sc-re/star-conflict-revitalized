package auth

import (
	"reflect"
	"testing"
)

func TestPassword(t *testing.T) {
	args := phc{
		id:              "argon2",
		version:         0x19,
		memoryCost:      64,
		timeCost:        128,
		parallelismCost: 2,
		keyLength:       3,
		salt:            []byte("salt"),
		hash:            []byte("hash"),
	}
	phcString := args.phcEncode()
	argsRet := phc{}
	if err := argsRet.phcDecode(phcString); err != nil {
		t.Errorf("%v", err)
	}
	if !reflect.DeepEqual(args, argsRet) {
		t.Errorf("phcDecode and phcEncode do not argree on coding args=%v, argsRet=%v, string=%v", args, argsRet, phcString)
	}
}
