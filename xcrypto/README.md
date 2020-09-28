# xcrypto

### Dependencies

+ xtesting*

### Functions

+ `FNV32(text string) uint32`
+ `FNV32a(text string) uint32`
+ `FNV64(text string) uint64`
+ `FNV64a(text string) uint64`
+ `MD5(text string) string`
+ `SHA1(text string) string`
+ `SHA256(text string) string`
+ `SHA512(text string) string`
+ `StringHasher(algorithm hash.Hash, text string) string`
+ `Uint32Hasher(algorithm hash.Hash32, text string) uint32`
+ `Uint64Hasher(algorithm hash.Hash64, text string) uint64`
+ `Base32Encode(data []byte) string`
+ `Base32Decode(data string) ([]byte, error)`
+ `Base64Encode(data []byte) string`
+ `Base64Decode(data string) ([]byte, error)`
