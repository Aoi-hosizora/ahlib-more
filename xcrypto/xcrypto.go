package xcrypto

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/sha3"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/fnv"
)

// ====
// hash
// ====

// FNV32 uses fnv32 to hash string to uint32.
func FNV32(text string) uint32 {
	algorithm := fnv.New32() // hash/fnv
	return Uint32Hasher(algorithm, text)
}

// FNV32a uses fnv32a to hash string to uint32.
func FNV32a(text string) uint32 {
	algorithm := fnv.New32a() // hash/fnv
	return Uint32Hasher(algorithm, text)
}

// FNV64 uses fnv64 to hash string to uint64.
func FNV64(text string) uint64 {
	algorithm := fnv.New64() // hash/fnv
	return Uint64Hasher(algorithm, text)
}

// FNV64a uses fnv64a to hash string to uint64.
func FNV64a(text string) uint64 {
	algorithm := fnv.New64a() // hash/fnv
	return Uint64Hasher(algorithm, text)
}

// CRC32 uses crc32 to hash string to uint32.
func CRC32(text string) uint32 {
	algorithm := crc32.NewIEEE() // hash/crc32
	return Uint32Hasher(algorithm, text)
}

// ADLER32 uses adler32 to hash string to uint32.
func ADLER32(text string) uint32 {
	algorithm := adler32.New() // hash/adler32
	return Uint32Hasher(algorithm, text)
}

// Reference: https://crypto.stackexchange.com/questions/68307/what-is-the-difference-between-sha-3-and-sha-256

// MD4 uses md4 to hash string.
func MD4(text string) string {
	algorithm := md4.New() // x/crypto/md4
	return StringHasher(algorithm, text)
}

// MD5 uses md5 to hash string.
func MD5(text string) string {
	algorithm := md5.New() // crypto/md5
	return StringHasher(algorithm, text)
}

// SHA1 uses sha-1 to hash string.
func SHA1(text string) string {
	algorithm := sha1.New() // crypto/sha1
	return StringHasher(algorithm, text)
}

// SHA224 uses sha2-224 to hash string.
func SHA224(text string) string {
	algorithm := sha256.New224() // crypto/sha256
	return StringHasher(algorithm, text)
}

// SHA256 uses sha2-256 to hash string.
func SHA256(text string) string {
	algorithm := sha256.New() // crypto/sha256
	return StringHasher(algorithm, text)
}

// SHA384 uses sha2-384 to hash string.
func SHA384(text string) string {
	algorithm := sha512.New384() // crypto/sha512
	return StringHasher(algorithm, text)
}

// SHA512 uses sha2-512 to hash string.
func SHA512(text string) string {
	algorithm := sha512.New() // crypto/sha512
	return StringHasher(algorithm, text)
}

// SHA512_224 uses sha2-512/224 to hash string.
func SHA512_224(text string) string {
	algorithm := sha512.New512_224() // crypto/sha512
	return StringHasher(algorithm, text)
}

// SHA512_256 uses sha2-512/256 to hash string.
func SHA512_256(text string) string {
	algorithm := sha512.New512_256() // crypto/sha512
	return StringHasher(algorithm, text)
}

// SHA3_224 uses sha3-224 to hash string.
func SHA3_224(text string) string {
	algorithm := sha3.New224() // x/crypto/sha3
	return StringHasher(algorithm, text)
}

// SHA3_256 uses sha3-256 to hash string.
func SHA3_256(text string) string {
	algorithm := sha3.New256() // x/crypto/sha3
	return StringHasher(algorithm, text)
}

// SHA3_384 uses sha3-384 to hash string.
func SHA3_384(text string) string {
	algorithm := sha3.New384() // x/crypto/sha3
	return StringHasher(algorithm, text)
}

// SHA3_512 uses sha3-512 to hash string.
func SHA3_512(text string) string {
	algorithm := sha3.New512() // x/crypto/sha3
	return StringHasher(algorithm, text)
}

// Uint32Hasher uses hash.Hash32 to encode string to uint32.
func Uint32Hasher(algorithm hash.Hash32, text string) uint32 {
	_, _ = algorithm.Write(xstring.FastStob(text))
	return algorithm.Sum32()
}

// Uint64Hasher uses hash.Hash64 to encode string to uint64.
func Uint64Hasher(algorithm hash.Hash64, text string) uint64 {
	_, _ = algorithm.Write(xstring.FastStob(text))
	return algorithm.Sum64()
}

// StringHasher uses hash.Hash to encode string to string.
func StringHasher(algorithm hash.Hash, text string) string {
	_, _ = algorithm.Write(xstring.FastStob(text))
	return HexEncodeToString(algorithm.Sum(nil))
}

// ===============
// encode & decode
// ===============

// HexEncodeToBytes encodes bytes to hex bytes.
func HexEncodeToBytes(data []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(data))) // encoding/hex
	hex.Encode(dst, data)
	return dst
}

// HexEncodeToString encodes bytes to hex string.
func HexEncodeToString(data []byte) string {
	return xstring.FastBtos(HexEncodeToBytes(data))
}

// HexDecodeFromBytes decodes bytes from hex bytes.
func HexDecodeFromBytes(data []byte) ([]byte, error) {
	buf := make([]byte, hex.DecodedLen(len(data))) // encoding/hex
	n, err := hex.Decode(buf, data)
	return buf[:n], err
}

// HexDecodeFromString decodes bytes from hex string.
func HexDecodeFromString(data string) ([]byte, error) {
	return HexDecodeFromBytes(xstring.FastStob(data))
}

// Base32EncodeToBytes encodes bytes to base32 bytes.
func Base32EncodeToBytes(data []byte) []byte {
	enc := base32.StdEncoding // encoding/base32
	buf := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(buf, data)
	return buf
}

// Base32EncodeToString encodes bytes to base32 string.
func Base32EncodeToString(data []byte) string {
	return xstring.FastBtos(Base32EncodeToBytes(data))
}

// Base32DecodeFromBytes decodes bytes from base32 bytes.
func Base32DecodeFromBytes(data []byte) ([]byte, error) {
	enc := base32.StdEncoding // encoding/base32
	buf := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(buf, data)
	return buf[:n], err
}

// Base32DecodeFromString decodes bytes from base32 string.
func Base32DecodeFromString(data string) ([]byte, error) {
	return Base32DecodeFromBytes(xstring.FastStob(data))
}

// Base64EncodeToBytes encodes bytes to base64 bytes.
func Base64EncodeToBytes(data []byte) []byte {
	enc := base64.StdEncoding // encoding/base64
	buf := make([]byte, enc.EncodedLen(len(data)))
	enc.Encode(buf, data)
	return buf
}

// Base64EncodeToString encodes bytes to base64 string.
func Base64EncodeToString(data []byte) string {
	return xstring.FastBtos(Base64EncodeToBytes(data))
}

// Base64DecodeFromBytes decodes bytes from base64 bytes.
func Base64DecodeFromBytes(data []byte) ([]byte, error) {
	enc := base64.StdEncoding // encoding/base64
	buf := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(buf, data)
	return buf[:n], err
}

// Base64DecodeFromString decodes bytes from base64 string.
func Base64DecodeFromString(data string) ([]byte, error) {
	return Base64DecodeFromBytes(xstring.FastStob(data))
}

// ====
// pkcs
// ====

const (
	panicBlockSize = "xcrypto: blockSize must larger then 0"
)

// PKCS5Padding uses PKCS#5 and PKCS#7 to pad data to block aligned bytes.
func PKCS5Padding(data []byte, blockSize int) []byte {
	if blockSize <= 0 {
		panic(panicBlockSize)
	}

	padLen := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padText...)
}

// PKCS5Trimming uses PKCS#5 and PKCS#7 to trim data from block aligned bytes.
func PKCS5Trimming(data []byte) []byte {
	length := len(data)
	padLen := int(data[length-1])
	return data[:length-padLen]
}

// ======
// bcrypt
// ======

const (
	BcryptMinCost     int = 4  // The bcrypt minimum allowable cost.
	BcryptMaxCost     int = 31 // The bcrypt maximum allowable cost.
	BcryptDefaultCost int = 10 // The bcrypt default cost, and this will actually be set if a cost is below BcryptMinCost.
)

// BcryptEncrypt uses bcrypt to encrypt password using given cost.
func BcryptEncrypt(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost) // x/crypto/bcrypt
}

// BcryptEncryptWithDefaultCost uses bcrypt to encrypt password using BcryptDefaultCost.
func BcryptEncryptWithDefaultCost(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, BcryptDefaultCost)
}

// BcryptCompare compares hashed encrypted password and given password.
func BcryptCompare(password, encrypted []byte) (ok bool, err error) {
	err = bcrypt.CompareHashAndPassword(encrypted, password)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}
