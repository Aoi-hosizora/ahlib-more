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

	src = "\xef\xbf\xbetest"
	dest = TrimBomString(src)
	xtesting.Equal(t, dest, "test")

	src2 := []byte("\xef\xbb\xbftest")
	dest2 := TrimBomBytes(src2)
	xtesting.Equal(t, dest2, []byte("test"))

	src2 = []byte("\xef\xbf\xbetest")
	dest2 = TrimBomBytes(src2)
	xtesting.Equal(t, dest2, []byte("test"))
}

func TestGetEncoding(t *testing.T) {
	_, ok := GetEncoding(IANA_UTF8)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_UTF16BE)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_UTF16LE)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_UTF32BE)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_UTF32LE)
	xtesting.True(t, ok)

	_, ok = GetEncoding(IANA_SHIFTJIS)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_EUCJP)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_ISO2022JP)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_GB18030)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_BIG5)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_ISO2022CN)
	xtesting.False(t, ok)
	_, ok = GetEncoding(IANA_EUCKR)
	xtesting.True(t, ok)
	_, ok = GetEncoding(IANA_ISO2022KR)
	xtesting.False(t, ok)

	for _, name := range []string{
		IANA_ISO88591, IANA_ISO88595, IANA_ISO88596, IANA_ISO88597, IANA_ISO88598, IANA_ISO88598I, IANA_ISO88599,
		IANA_WINDOWS1251, IANA_WINDOWS1256, IANA_KOI8R, IANA_IBM424RTL, IANA_IBM424LTR,
	} {
		_, ok := GetEncoding(name)
		xtesting.False(t, ok)
	}

	_, ok = GetEncoding("")
	xtesting.False(t, ok)
}
