package xcharset

import (
	"bytes"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"strings"
)

const (
	BOM  = "\xef\xbb\xbf"
	BOM2 = "\xef\xbf\xbe"
)

// TrimBomString removes BOM from a string.
func TrimBomString(str string) string {
	return strings.TrimPrefix(strings.TrimPrefix(str, BOM), BOM2)
}

// TrimBomBytes removes BOM from a bytes.
func TrimBomBytes(bs []byte) []byte {
	return bytes.TrimPrefix(bytes.TrimPrefix(bs, []byte(BOM)), []byte(BOM2))
}

const (
	IANA_UTF8    = "UTF-8"    // *
	IANA_UTF16BE = "UTF-16BE" // *
	IANA_UTF16LE = "UTF-16LE" // *
	IANA_UTF32BE = "UTF-32BE" // *
	IANA_UTF32LE = "UTF-32LE" // *

	IANA_ISO88591    = "ISO-8859-1"   // Latin-1, 1: en, da, de, es, fr, it, nl, no, pt, sv; 2: cs, hu, pl, ro
	IANA_ISO88595    = "ISO-8859-5"   // ru
	IANA_ISO88596    = "ISO-8859-6"   // ar
	IANA_ISO88597    = "ISO-8859-7"   // el
	IANA_ISO88598    = "ISO-8859-8"   // he
	IANA_ISO88598I   = "ISO-8859-8-I" // he
	IANA_ISO88599    = "ISO-8859-9"   // tr
	IANA_WINDOWS1251 = "windows-1251" // ar
	IANA_WINDOWS1256 = "windows-1256" // ar
	IANA_KOI8R       = "KOI8-R"       // ru

	IANA_SHIFTJIS  = "Shift_JIS"   // ja
	IANA_GB18030   = "GB-18030"    // zh
	IANA_EUCJP     = "EUC-JP"      // ja
	IANA_EUCKR     = "EUC-KR"      // ko
	IANA_BIG5      = "Big5"        // zh
	IANA_ISO2022JP = "ISO-2022-JP" // jp
	IANA_ISO2022KR = "ISO-2022-KR" // kr
	IANA_ISO2022CN = "ISO-2022-CN" // cn

	IANA_IBM424RTL = "IBM420_rtl" // he, ar
	IANA_IBM424LTR = "IBM420_ltr" // he, ar
)

// GetEncoding returns a encoding.Encoding from some IANA.
func GetEncoding(iana string) (encode encoding.Encoding, existed bool) {
	switch iana {
	case IANA_UTF8:
		return unicode.UTF8, true
	case IANA_UTF16BE:
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM), true
	case IANA_UTF16LE:
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM), true
	case IANA_UTF32BE:
		return utf32.UTF32(utf32.BigEndian, utf32.IgnoreBOM), true
	case IANA_UTF32LE:
		return utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM), true

	case IANA_SHIFTJIS:
		return japanese.ShiftJIS, true
	case IANA_EUCJP:
		return japanese.EUCJP, true
	case IANA_ISO2022JP:
		return japanese.ISO2022JP, true
	case IANA_GB18030:
		return simplifiedchinese.GB18030, true
	case IANA_BIG5:
		return traditionalchinese.Big5, true
	case IANA_ISO2022CN:
		// not found
	case IANA_EUCKR:
		return korean.EUCKR, true
	case IANA_ISO2022KR:
		// not found

	case IANA_ISO88591, IANA_ISO88595, IANA_ISO88596, IANA_ISO88597, IANA_ISO88598, IANA_ISO88598I, IANA_ISO88599:
		// not found
	case IANA_WINDOWS1251, IANA_WINDOWS1256, IANA_KOI8R, IANA_IBM424RTL, IANA_IBM424LTR:
		// not found
	}

	return nil, false
}
