package slev

import (
	"errors"
	"math/big"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	BASE62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	HEX    = "0123456789abcdef"
)

type slevID struct {
	Base62 string
	Hex    string
	Bytes  []byte
	Int    big.Int
	UUID   string
	rand   big.Int
}

func (f *slevID) Timestamp() int64 {
	t := new(big.Int)
	t.Rsh(&f.Int, 80)
	return t.Int64()
}

func NewID() (string, error) {
	now := time.Now()
	timestamp := now.UnixMilli()

	bigTimestamp := big.NewInt(timestamp)

	//Generate cryptographically strong pseudo-random between 0 - max
	p := make([]byte, 10) // 80bits
	rand.Read(p)

	randomPart := new(big.Int)
	randomPart.SetBytes(p)

	timestampPart := new(big.Int)
	timestampPart.Lsh(bigTimestamp, 80)

	//Mix with random number
	randomAndTimestampCombined := timestampPart.Or(timestampPart, randomPart)

	v62 := big.NewInt(0)
	v62.SetBytes(randomAndTimestampCombined.Bytes())

	return BigIntToASCII(v62, BASE62, 22)
}

func FromBytes(fromBytes []byte) (*slevID, error) {
	if len(fromBytes) != 16 {
		return nil, errors.New("FromBytes should be len 16")
	}

	bigTimestamp := new(big.Int)
	bigTimestamp.SetBytes(fromBytes[0:6])

	randomBytes := fromBytes[6:16]

	randomPart := new(big.Int)
	randomPart.SetBytes(randomBytes)

	timestampPart := big.NewInt(0)
	timestampPart.SetBytes(bigTimestamp.Bytes())
	timestampPart.Lsh(bigTimestamp, 80)

	//Mix with random number
	randomAndTimestampCombined := timestampPart.Or(timestampPart, randomPart)

	v62 := big.NewInt(0)
	v62.SetBytes(randomAndTimestampCombined.Bytes())

	vHex := big.NewInt(0)
	vHex.SetBytes(randomAndTimestampCombined.Bytes())

	b62, b62Err := BigIntToASCII(v62, BASE62, 22)
	hex, hexErr := BigIntToASCII(vHex, HEX, 32)

	if b62Err != nil {
		return nil, errors.New("conversion to BASE62 failed")
	}
	if hexErr != nil {
		return nil, errors.New("conversion to HEX failed")
	}

	u, UUIDErr := uuid.FromBytes(randomAndTimestampCombined.Bytes())

	if UUIDErr != nil {
		return nil, errors.New("UUID generation failed")
	}

	f := slevID{
		Base62: b62,
		Hex:    hex,
		Bytes:  fromBytes,
		UUID:   u.String(),
		Int:    *randomAndTimestampCombined,
		rand:   *randomPart,
	}

	return &f, nil
}

func FromBase62(b62 string) (*slevID, error) {
	bigInt := ASCIIToBigInt(b62, BASE62)
	return FromBytes(bigInt.Bytes())
}

type Values interface {
	Timestamp() int64
	Random() *big.Int
}

// Struct is not exported
type valuesParam struct {
	ts int64
	r  *big.Int
}

func (v *valuesParam) Timestamp() int64 {
	return v.ts
}

func (v *valuesParam) Random() *big.Int {
	return v.r
}
