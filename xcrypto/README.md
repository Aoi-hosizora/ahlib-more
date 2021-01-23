# xcrypto

## Dependencies

+ golang.org/x/crypto
+ xtesting*

## Documents

### Types

+ None

### Variables

+ None

### Constants

+ `const BcryptMinCost int`
+ `const BcryptMaxCost int`
+ `const BcryptDefaultCost int`

### Functions

+ `func FNV32(text string) uint32`
+ `func FNV32a(text string) uint32`
+ `func FNV64(text string) uint64`
+ `func FNV64a(text string) uint64`
+ `func CRC32(text string) uint32`
+ `func ADLER32(text string) uint32`
+ `func MD4(text string) string`
+ `func MD5(text string) string`
+ `func SHA1(text string) string`
+ `func SHA224(text string) string`
+ `func SHA256(text string) string`
+ `func SHA384(text string) string`
+ `func SHA512(text string) string`
+ `func SHA512_224(text string) string`
+ `func SHA512_256(text string) string`
+ `func SHA3_224(text string) string`
+ `func SHA3_256(text string) string`
+ `func SHA3_384(text string) string`
+ `func SHA3_512(text string) string`
+ `func Uint32Hasher(algorithm hash.Hash32, text string) uint32`
+ `func Uint64Hasher(algorithm hash.Hash64, text string) uint64`
+ `func StringHasher(algorithm hash.Hash, text string) string`
+ `func HexEncodeToBytes(data []byte) []byte`
+ `func HexEncodeToString(data []byte) string`
+ `func HexDecodeFromBytes(data []byte) ([]byte, error)`
+ `func HexDecodeFromString(data string) ([]byte, error)`
+ `func Base32EncodeToBytes(data []byte) []byte`
+ `func Base32EncodeToString(data []byte) string`
+ `func Base32DecodeFromBytes(data []byte) ([]byte, error)`
+ `func Base32DecodeFromString(data string) ([]byte, error)`
+ `func Base64EncodeToBytes(data []byte) []byte`
+ `func Base64EncodeToString(data []byte) string`
+ `func Base64DecodeFromBytes(data []byte) ([]byte, error)`
+ `func Base64DecodeFromString(data string) ([]byte, error)`
+ `func PKCS5Padding(data []byte, blockSize int) []byte`
+ `func PKCS5Trimming(data []byte) []byte`
+ `func BcryptEncrypt(password []byte, cost int) ([]byte, error)`
+ `func BcryptEncryptWithDefaultCost(password []byte) ([]byte, error)`
+ `func BcryptCompare(password, encrypted []byte) (ok bool, err error)`

### Methods

+ None
