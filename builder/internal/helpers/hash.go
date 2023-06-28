package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"strconv"
	"strings"

	"github.com/packwiz/packwiz/curseforge/murmur2"
)

type HashStringer interface {
	hash.Hash
	HashToString([]byte) string
}

type HexStringer struct {
	hash.Hash
}

type Number32As64Stringer struct {
	hash.Hash
}

type Number64Stringer struct {
	hash.Hash
}

func (HexStringer) HashToString(data []byte) string {
	return hex.EncodeToString(data)
}

func (Number32As64Stringer) HashToString(data []byte) string {
	return strconv.FormatUint(uint64(binary.BigEndian.Uint32(data)), 10)
}

func (Number64Stringer) HashToString(data []byte) string {
	return strconv.FormatUint(binary.BigEndian.Uint64(data), 10)
}

func GetHashImpl(hashType string) (HashStringer, error) {
	switch strings.ToLower(hashType) {
	case "sha1":
		return HexStringer{sha1.New()}, nil
	case "sha256":
		return HexStringer{sha256.New()}, nil
	case "sha512":
		return HexStringer{sha512.New()}, nil
	case "md5":
		return HexStringer{md5.New()}, nil
	case "murmur2":
		return Number32As64Stringer{murmur2.New()}, nil
	case "length-bytes":
		return Number64Stringer{&LengthHasher{}}, nil
	}

	return nil, fmt.Errorf("hash implementation %s not found", hashType)
}
