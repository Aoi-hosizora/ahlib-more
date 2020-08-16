package xencode

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"testing"
)

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
