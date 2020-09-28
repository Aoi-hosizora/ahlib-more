package xcharset

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"golang.org/x/text/encoding/unicode"
	"testing"
)

func TestDetectCharsetBest(t *testing.T) {
	res, err := DetectCharsetBest([]byte("test"))
	xtesting.Nil(t, err)
	xtesting.Equal(t, res.Charset, "ISO-8859-1")
}

func TestDetectCharsetAll(t *testing.T) {
	res, err := DetectCharsetAll([]byte("test"))
	xtesting.Nil(t, err)
	xtesting.Equal(t, res[0].Charset, "ISO-8859-1")
}

func TestEncode(t *testing.T) {
	dest, err := EncodeString(unicode.UTF8, "test")
	xtesting.Nil(t, err)
	xtesting.Equal(t, dest, "test")

	dest2, err := EncodeBytes(unicode.UTF8, []byte("test"))
	xtesting.Nil(t, err)
	xtesting.Equal(t, dest2, []byte("test"))
}

func TestDecode(t *testing.T) {
	dest, err := DecodeString(unicode.UTF8, "test")
	xtesting.Nil(t, err)
	xtesting.Equal(t, dest, "test")

	dest2, err := DecodeBytes(unicode.UTF8, []byte("test"))
	xtesting.Nil(t, err)
	xtesting.Equal(t, dest2, []byte("test"))
}

func TestTrimBom(t *testing.T) {
	src := "\xef\xbb\xbftest"
	dest := TrimBomString(src)
	xtesting.Equal(t, dest, "test")

	src2 := []byte("\xef\xbb\xbftest")
	dest2 := TrimBomBytes(src2)
	xtesting.Equal(t, dest2, []byte("test"))
}

func TestGetEncoding(t *testing.T) {
	_ = GetEncoding
}
