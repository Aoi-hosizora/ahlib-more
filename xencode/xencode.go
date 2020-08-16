package xencode

import (
	"encoding/base32"
	"encoding/base64"
)

func Base32Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

func Base32Decode(data string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(data)
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}