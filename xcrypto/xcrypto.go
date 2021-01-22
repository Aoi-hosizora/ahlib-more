package xcrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"hash"
	"hash/fnv"
)

func FNV32(text string) uint32 {
	algorithm := fnv.New32() // hash/fnv
	return Uint32Hasher(algorithm, text)
}

func FNV32a(text string) uint32 {
	algorithm := fnv.New32a() // hash/fnv
	return Uint32Hasher(algorithm, text)
}

func FNV64(text string) uint64 {
	algorithm := fnv.New64() // hash/fnv
	return Uint64Hasher(algorithm, text)
}

func FNV64a(text string) uint64 {
	algorithm := fnv.New64a() // hash/fnv
	return Uint64Hasher(algorithm, text)
}

func MD5(text string) string {
	algorithm := md5.New() // crypto/md5
	return StringHasher(algorithm, text)
}

func SHA1(text string) string {
	algorithm := sha1.New() // crypto/sha1
	return StringHasher(algorithm, text)
}

func SHA256(text string) string {
	algorithm := sha256.New() // crypto/sha256
	return StringHasher(algorithm, text)
}

func SHA512(text string) string {
	algorithm := sha512.New() // crypto/sha512
	return StringHasher(algorithm, text)
}

// StringHasher use hash.Hash to encode string.
func StringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Uint32Hasher use hash.Hash to encode uint32.
func Uint32Hasher(algorithm hash.Hash32, text string) uint32 {
	_, _ = algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

// Uint64Hasher use hash.Hash to encode uint64.
func Uint64Hasher(algorithm hash.Hash64, text string) uint64 {
	_, _ = algorithm.Write([]byte(text))
	return algorithm.Sum64()
}

// Base32Encode use base32.StdEncoding (standard base32 encoding) to encode in base32.
func Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data) // encoding/base32
}

// Base32Decode use base32.StdEncoding (standard base32 encoding) to decode in base32.
func Base32Decode(data string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(data) // encoding/base32
}

// Base64Encode use base64.StdEncoding (standard base32 encoding) to encode in base64.
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // encoding/base64
}

// Base64Decode use base64.StdEncoding (standard base32 encoding) to decode in base64.
func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data) // encoding/base64
}

const (
	MinCost     int = 4  // The minimum allowable cost.
	MaxCost     int = 31 // The maximum allowable cost.
	DefaultCost int = 10 // The cost that will actually be set if a cost is below MinCost.
)

// Use bcrypt with cost to encrypt password.
// If the cost given is less than MinCost, the cost will be set to DefaultCost instead.
func Encrypt(password []byte, cost int) ([]byte, error) {
	encrypted, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

// Use bcrypt with DefaultCost to encrypt password.
func EncryptWithDefaultCost(password []byte) ([]byte, error) {
	return Encrypt(password, DefaultCost)
}

// Check the password is the same.
func Check(password, encrypted []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(encrypted, password)
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}
