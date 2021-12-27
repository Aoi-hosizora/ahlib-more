package xcrypto

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"testing"
)

func TestUintHasher(t *testing.T) {
	for _, tc := range []struct {
		giveFn   func(string) uint32
		giveText string
		want     uint32
	}{
		{FNV32, "", 0x811c9dc5},
		{FNV32, "test", 0xbc2c0be9},
		{FNV32, "hello world", 0x548da96f},
		{FNV32, "测试 テス тест", 0x15a513dc},
		{FNV32a, "", 0x811c9dc5},
		{FNV32a, "test", 0xafd071e5},
		{FNV32a, "hello world", 0xd58b3fa7},
		{FNV32a, "测试 テス тест", 0x1e0ace72},
		{CRC32, "", 0x0},
		{CRC32, "test", 0xd87f7e0c},
		{CRC32, "hello world", 0xd4a1185},
		{CRC32, "测试 テス тест", 0xa699b2d6},
		{ADLER32, "", 0x1},
		{ADLER32, "test", 0x45d01c1},
		{ADLER32, "hello world", 0x1a0b045d},
		{ADLER32, "测试 テス тест", 0xa43d0e1a},
	} {
		xtesting.Equal(t, tc.giveFn(tc.giveText), tc.want)
	}

	for _, tc := range []struct {
		giveFn   func(string) uint64
		giveText string
		want     uint64
	}{
		{FNV64, "", 0xcbf29ce484222325},
		{FNV64, "test", 0x8c093f7e9fccbf69},
		{FNV64, "hello world", 0x7dcf62cdb1910e6f},
		{FNV64, "测试 テス тест", 0xefa05d5a0bc1da7c},
		{FNV64a, "", 0xcbf29ce484222325},
		{FNV64a, "test", 0xf9e6e6ef197c2b25},
		{FNV64a, "hello world", 0x779a65e7023cd2e7},
		{FNV64a, "测试 テス тест", 0xa8009ce94a3ad872},
		{CRC64, "", 0x0},
		{CRC64, "test", 0x287c72c850000000},
		{CRC64, "hello world", 0xb9cf3f572ad9ac3e},
		{CRC64, "测试 テス тест", 0xe16038d3f4fca746},
	} {
		xtesting.Equal(t, tc.giveFn(tc.giveText), tc.want)
	}
}

func TestStringHasher(t *testing.T) {
	for _, tc := range []struct {
		giveFn   func(string) string
		giveText string
		wantText string
	}{
		{FNV128, "", "6c62272e07bb014262b821756295c58d"},
		{FNV128, "test", "66ab2a8b6f757277b806e89c56faf339"},
		{FNV128, "hello world", "e1b1650f0631aef5566634b6c074ac1f"},
		{FNV128a, "", "6c62272e07bb014262b821756295c58d"},
		{FNV128a, "test", "69d061a9c5757277b806e99413dd99a5"},
		{FNV128a, "hello world", "6c155799fdc8eec4b91523808e7726b7"},
		{MD4, "", "31d6cfe0d16ae931b73c59d7e0c089c0"},
		{MD4, "test", "db346d691d7acc4dc2625db19f9e3f52"},
		{MD4, "hello world", "aa010fbc1d14c795d86ef98c95479d17"},
		{MD5, "", "d41d8cd98f00b204e9800998ecf8427e"},
		{MD5, "test", "098f6bcd4621d373cade4e832627b4f6"},
		{MD5, "hello world", "5eb63bbbe01eeed093cb22bb8f5acdc3"},
		{SHA1, "", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{SHA1, "test", "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"},
		{SHA1, "hello world", "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"},
		{SHA224, "", "d14a028c2a3a2bc9476102bb288234c415a2b01f828ea62ac5b3e42f"},
		{SHA224, "test", "90a3ed9e32b2aaf4c61c410eb925426119e1a9dc53d4286ade99a809"},
		{SHA224, "hello world", "2f05477fc24bb4faefd86517156dafdecec45b8ad3cf2522a563582b"},
		{SHA256, "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{SHA256, "test", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
		{SHA256, "hello world", "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
		{SHA384, "", "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b"},
		{SHA384, "test", "768412320f7b0aa5812fce428dc4706b3cae50e02a64caa16a782249bfe8efc4b7ef1ccb126255d196047dfedf17a0a9"},
		{SHA384, "hello world", "fdbd8e75a67f29f701a4e040385e2e23986303ea10239211af907fcbb83578b3e417cb71ce646efd0819dd8c088de1bd"},
		{SHA512, "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{SHA512, "test", "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"},
		{SHA512, "hello world", "309ecc489c12d6eb4cc40f50c902f2b4d0ed77ee511a7c7a9bcd3ca86d4cd86f989dd35bc5ff499670da34255b45b0cfd830e81f605dcf7dc5542e93ae9cd76f"},
		{SHA512_224, "", "6ed0dd02806fa89e25de060c19d3ac86cabb87d6a0ddd05c333b84f4"},
		{SHA512_224, "test", "06001bf08dfb17d2b54925116823be230e98b5c6c278303bc4909a8c"},
		{SHA512_224, "hello world", "22e0d52336f64a998085078b05a6e37b26f8120f43bf4db4c43a64ee"},
		{SHA512_256, "", "c672b8d1ef56ed28ab87c3622c5114069bdd3ad7b8f9737498d0c01ecef0967a"},
		{SHA512_256, "test", "3d37fe58435e0d87323dee4a2c1b339ef954de63716ee79f5747f94d974f913f"},
		{SHA512_256, "hello world", "0ac561fac838104e3f2e4ad107b4bee3e938bf15f2b15f009ccccd61a913f017"},
		{SHA3_224, "", "6b4e03423667dbb73b6e15454f0eb1abd4597f9a1b078e3f5b5a6bc7"},
		{SHA3_224, "test", "3797bf0afbbfca4a7bbba7602a2b552746876517a7f9b7ce2db0ae7b"},
		{SHA3_224, "hello world", "dfb7f18c77e928bb56faeb2da27291bd790bc1045cde45f3210bb6c5"},
		{SHA3_256, "", "a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a"},
		{SHA3_256, "test", "36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80"},
		{SHA3_256, "hello world", "644bcc7e564373040999aac89e7622f3ca71fba1d972fd94a31c3bfbf24e3938"},
		{SHA3_384, "", "0c63a75b845e4f7d01107d852e4c2485c51a50aaaa94fc61995e71bbee983a2ac3713831264adb47fb6bd1e058d5f004"},
		{SHA3_384, "test", "e516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd"},
		{SHA3_384, "hello world", "83bff28dde1b1bf5810071c6643c08e5b05bdb836effd70b403ea8ea0a634dc4997eb1053aa3593f590f9c63630dd90b"},
		{SHA3_512, "", "a69f73cca23a9ac5c8b567dc185a756e97c982164fe25859e0d1dcc1475c80a615b2123af1f5f94c11e3e9402c3ac558f500199d95b6d3e301758586281dcd26"},
		{SHA3_512, "test", "9ece086e9bac491fac5c1d1046ca11d737b92a2b2ebd93f005d7b710110c0a678288166e7fbe796883a4f2e9b3ca9f484f521d0ce464345cc1aec96779149c14"},
		{SHA3_512, "hello world", "840006653e9ac9e95117a15c915caab81662918e925de9e004f774ff82d7079a40d4d27b1b372657c61d46d470304c88c788b3a4527ad074d1dccbee5dbaa99a"},
	} {
		xtesting.Equal(t, tc.giveFn(tc.giveText), tc.wantText)
	}
}

func TestEncodeDecode(t *testing.T) {
	testStr := "test"
	testBs := []byte(testStr)
	helloWorldStr := "hello world"
	helloWorldBs := []byte(helloWorldStr)
	test2Str := "测试 テス тест"
	test2Bs := []byte(test2Str)

	for _, tc := range []struct {
		giveFn1 func([]byte) []byte
		giveFn2 func([]byte) string
		give    []byte
		want    string
	}{
		{HexEncodeToBytes, HexEncodeToString, nil, ""},
		{HexEncodeToBytes, HexEncodeToString, testBs, "74657374"},
		{HexEncodeToBytes, HexEncodeToString, helloWorldBs, "68656c6c6f20776f726c64"},
		{HexEncodeToBytes, HexEncodeToString, test2Bs, "e6b58be8af9520e38386e382b920d182d0b5d181d182"},
		{Base32EncodeToBytes, Base32EncodeToString, nil, ""},
		{Base32EncodeToBytes, Base32EncodeToString, testBs, "ORSXG5A="},
		{Base32EncodeToBytes, Base32EncodeToString, helloWorldBs, "NBSWY3DPEB3W64TMMQ======"},
		{Base32EncodeToBytes, Base32EncodeToString, test2Bs, "422YX2FPSUQOHA4G4OBLSIGRQLILLUMB2GBA===="},
		{Base64EncodeToBytes, Base64EncodeToString, nil, ""},
		{Base64EncodeToBytes, Base64EncodeToString, testBs, "dGVzdA=="},
		{Base64EncodeToBytes, Base64EncodeToString, helloWorldBs, "aGVsbG8gd29ybGQ="},
		{Base64EncodeToBytes, Base64EncodeToString, test2Bs, "5rWL6K+VIOODhuOCuSDRgtC10YHRgg=="},
	} {
		xtesting.Equal(t, string(tc.giveFn1(tc.give)), tc.want)
		xtesting.Equal(t, tc.giveFn2(tc.give), tc.want)
	}

	for _, tc := range []struct {
		giveFn1 func([]byte) ([]byte, error)
		giveFn2 func(string) ([]byte, error)
		give    string
		want    string
	}{
		{HexDecodeFromBytes, HexDecodeFromString, "", ""},
		{HexDecodeFromBytes, HexDecodeFromString, "74657374", testStr},
		{HexDecodeFromBytes, HexDecodeFromString, "68656c6c6f20776f726c64", helloWorldStr},
		{HexDecodeFromBytes, HexDecodeFromString, "e6b58be8af9520e38386e382b920d182d0b5d181d182", test2Str},
		{Base32DecodeFromBytes, Base32DecodeFromString, "", ""},
		{Base32DecodeFromBytes, Base32DecodeFromString, "ORSXG5A=", testStr},
		{Base32DecodeFromBytes, Base32DecodeFromString, "NBSWY3DPEB3W64TMMQ======", helloWorldStr},
		{Base32DecodeFromBytes, Base32DecodeFromString, "422YX2FPSUQOHA4G4OBLSIGRQLILLUMB2GBA====", test2Str},
		{Base64DecodeFromBytes, Base64DecodeFromString, "", ""},
		{Base64DecodeFromBytes, Base64DecodeFromString, "dGVzdA==", testStr},
		{Base64DecodeFromBytes, Base64DecodeFromString, "aGVsbG8gd29ybGQ=", helloWorldStr},
		{Base64DecodeFromBytes, Base64DecodeFromString, "5rWL6K+VIOODhuOCuSDRgtC10YHRgg==", test2Str},
	} {
		bs1, err := tc.giveFn1([]byte(tc.give))
		xtesting.Nil(t, err)
		xtesting.Equal(t, string(bs1), tc.want)

		bs2, err := tc.giveFn2(tc.give)
		xtesting.Nil(t, err)
		xtesting.Equal(t, string(bs2), tc.want)
	}
}

func TestPKCS5(t *testing.T) {
	for _, tc := range []struct {
		giveData    []byte
		giveSize    int
		wantAligned []byte
		wantPanic   bool
	}{
		{[]byte{}, 1, []byte{0x1}, false},
		{[]byte{}, 4, []byte{0x4, 0x4, 0x4, 0x4}, false},
		{[]byte{'t', 'e', 's', 't'}, 1, []byte{'t', 'e', 's', 't', 0x1}, false},
		{[]byte{'t', 'e', 's', 't'}, 5, []byte{'t', 'e', 's', 't', 0x1}, false},
		{[]byte{'t', 'e', 's', 't'}, 2, []byte{'t', 'e', 's', 't', 0x2, 0x2}, false},
		{[]byte{'t', 'e', 's', 't'}, 7, []byte{'t', 'e', 's', 't', 0x3, 0x3, 0x3}, false},
		{[]byte{'t', 'e', 's', 't'}, 4, []byte{'t', 'e', 's', 't', 0x4, 0x4, 0x4, 0x4}, false},
		{[]byte{'t', 'e', 's', 't', ' '}, 9, []byte{'t', 'e', 's', 't', ' ', 0x4, 0x4, 0x4, 0x4}, false},
		{[]byte{'t', 'e', 's', 't'}, 0, nil, true},
		{[]byte{'t', 'e', 's', 't'}, -1, nil, true},
	} {
		if tc.wantPanic {
			xtesting.Panic(t, func() { PKCS5Padding(tc.giveData, tc.giveSize) })
		} else {
			xtesting.Equal(t, PKCS5Padding(tc.giveData, tc.giveSize), tc.wantAligned)
			xtesting.Equal(t, PKCS5Trimming(tc.wantAligned), tc.giveData)
		}
	}
}

func TestBcrypt(t *testing.T) {
	_ = BcryptMaxCost // too slow

	for _, tc := range []struct {
		givePass   string
		giveCost   int
		useDefault bool
	}{
		{"", 0, false}, // -> BcryptMinCost
		{"test", 0, false},
		{"hello world", 0, false},
		{"", BcryptMinCost, false},
		{"test", BcryptMinCost, false},
		{"hello world", BcryptMinCost, false},
		{"", BcryptDefaultCost, false},
		{"test", BcryptDefaultCost, false},
		{"hello world", BcryptDefaultCost, false},
		{"", 0, true},
		{"test", 0, true},
		{"hello world", 0, true},
	} {
		pass := []byte(tc.givePass)
		var encrypted []byte
		if tc.useDefault {
			var err error
			encrypted, err = BcryptEncrypt(pass, tc.giveCost)
			xtesting.Nil(t, err)
		} else {
			var err error
			encrypted, err = BcryptEncryptWithDefaultCost(pass)
			xtesting.Nil(t, err)
		}

		ok, err := BcryptCompare(pass, encrypted)
		xtesting.True(t, ok)
		xtesting.Nil(t, err)

		ok, err = BcryptCompare([]byte("fake password"), encrypted)
		xtesting.False(t, ok)
		xtesting.Nil(t, err)

		ok, err = BcryptCompare(pass, []byte{})
		xtesting.False(t, ok)
		xtesting.NotNil(t, err)
	}
}
