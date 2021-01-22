package xcrypto

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"log"
	"testing"
)

func TestFNV32(t *testing.T) {
	data := FNV32("Raed Shomali")
	xtesting.Equal(t, data, uint32(0x194b953a))
}

func TestFNV32a(t *testing.T) {
	data := FNV32a("Raed Shomali")
	xtesting.Equal(t, data, uint32(0xdd36df08))
}

func TestFNV64(t *testing.T) {
	data := FNV64("Raed Shomali")
	xtesting.Equal(t, data, uint64(0xc03d55b5ff7722da))
}

func TestFNV64a(t *testing.T) {
	data := FNV64a("Raed Shomali")
	xtesting.Equal(t, data, uint64(0xf7fc847f1e6b4148))
}

func TestMD5(t *testing.T) {
	data := MD5("Raed Shomali")
	xtesting.Equal(t, data, "c313bc3b48fcfed9abc733429665b105")
}

func TestSHA1(t *testing.T) {
	data := SHA1("Raed Shomali")
	xtesting.Equal(t, data, "e0d66f6f09de72942e83289cc994b3c721ab34c5")
}

func TestSHA256(t *testing.T) {
	data := SHA256("Raed Shomali")
	xtesting.Equal(t, data, "75894b9be21065a833e57bfe4440b375fc216f120a965243c9be8b2dc36709c2")
}

func TestSHA512(t *testing.T) {
	data := SHA512("Raed Shomali")
	xtesting.Equal(t, data, "406e8d495140187a8b09893c30d054cf385ad7359855db0d2e0386c7189ac1c4667a4816d1b63a19f3d8ccdcbace7861ec4cc6ff5e2a1659c8f4360bda699b42")
}

func TestBase32Encode(t *testing.T) {
	data := Base32Encode([]byte("Raed Shomali"))
	xtesting.Equal(t, data, "KJQWKZBAKNUG63LBNRUQ====")
}

func TestBase32Decode(t *testing.T) {
	data, err := Base32Decode("KJQWKZBAKNUG63LBNRUQ====")
	xtesting.Equal(t, err, nil)
	xtesting.Equal(t, string(data), "Raed Shomali")
}

func TestBase64Encode(t *testing.T) {
	data := Base64Encode([]byte("Raed Shomali"))
	xtesting.Equal(t, data, "UmFlZCBTaG9tYWxp")
}

func TestBase64Decode(t *testing.T) {
	data, err := Base64Decode("UmFlZCBTaG9tYWxp")
	xtesting.Equal(t, err, nil)
	xtesting.Equal(t, string(data), "Raed Shomali")
}

func TestPassword(t *testing.T) {
	password := []byte("123")

	for _, cost := range []int{1, MinCost, 7, DefaultCost, 13} { // 1 4 7 10 13
		encrypted, err := Encrypt(password, cost)
		xtesting.Nil(t, err)
		log.Println(cost, ":", string(encrypted))
		check, err := Check(password, encrypted) // true nil
		xtesting.Nil(t, err)
		xtesting.True(t, check)
	}

	encrypted, err := Encrypt(password, MaxCost+100)
	xtesting.NotNil(t, err)

	encrypted, err = EncryptWithDefaultCost(password)
	xtesting.Nil(t, err)
	log.Println(DefaultCost, ":", string(encrypted))
	check, err := Check(password, encrypted) // true nil
	xtesting.Nil(t, err)
	xtesting.True(t, check)

	check, err = Check(password, []byte{}) // false, err
	xtesting.NotNil(t, err)
	xtesting.False(t, check)

	check, err = Check([]byte("miss"), encrypted) // false, nil
	xtesting.Nil(t, err)
	xtesting.False(t, check)
}
