package xhash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"hash/fnv"
)

func FNV32(text string) uint32 {
	algorithm := fnv.New32()
	return uint32Hasher(algorithm, text)
}

func FNV32a(text string) uint32 {
	algorithm := fnv.New32a()
	return uint32Hasher(algorithm, text)
}

func FNV64(text string) uint64 {
	algorithm := fnv.New64()
	return uint64Hasher(algorithm, text)
}

func FNV64a(text string) uint64 {
	algorithm := fnv.New64a()
	return uint64Hasher(algorithm, text)
}

func MD5(text string) string {
	algorithm := md5.New()
	return stringHasher(algorithm, text)
}

func SHA1(text string) string {
	algorithm := sha1.New()
	return stringHasher(algorithm, text)
}

func SHA256(text string) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, text)
}

func SHA512(text string) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, text)
}

func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func uint32Hasher(algorithm hash.Hash32, text string) uint32 {
	_, _ = algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

func uint64Hasher(algorithm hash.Hash64, text string) uint64 {
	_, _ = algorithm.Write([]byte(text))
	return algorithm.Sum64()
}
