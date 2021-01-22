package xcharset

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

// EncodeString encodes a string to given encoding.
func EncodeString(encoding encoding.Encoding, s string) (string, error) {
	result, _, err := transform.String(encoding.NewEncoder(), s)
	return result, err
}

// DecodeString decodes a string to given encoding.
func DecodeString(encoding encoding.Encoding, s string) (string, error) {
	result, _, err := transform.String(encoding.NewDecoder(), s)
	return result, err
}

// EncodeBytes encodes a bytes to given encoding.
func EncodeBytes(encoding encoding.Encoding, bs []byte) ([]byte, error) {
	result, _, err := transform.Bytes(encoding.NewEncoder(), bs)
	return result, err
}

// DecodeBytes decodes a bytes to given encoding.
func DecodeBytes(encoding encoding.Encoding, bs []byte) ([]byte, error) {
	result, _, err := transform.Bytes(encoding.NewDecoder(), bs)
	return result, err
}

// See https://github.com/saintfish/chardet/blob/master/detector.go and https://www.iana.org/assignments/charset-reg/charset-reg.xhtml.
const (
	IANA_UTF8    = "UTF-8"    // *
	IANA_UTF16BE = "UTF-16BE" // *
	IANA_UTF16LE = "UTF-16LE" // *
	IANA_UTF32BE = "UTF-32BE" // *
	IANA_UTF32LE = "UTF-32LE" // *

	IANA_ISO8859_1   = "ISO-8859-1"   // en, da, de, es, fr, it, nl, no, pt, sv
	IANA_ISO8859_2   = "ISO-8859-2"   // cs, hu, pl, ro
	IANA_ISO8859_5   = "ISO-8859-5"   // ru
	IANA_ISO8859_6   = "ISO-8859-6"   // ar
	IANA_ISO8859_7   = "ISO-8859-7"   // el
	IANA_ISO8859_8   = "ISO-8859-8"   // he
	IANA_ISO8859_8I  = "ISO-8859-8-I" // he
	IANA_ISO8859_9   = "ISO-8859-9"   // tr
	IANA_KOI8R       = "KOI8-R"       // ru
	IANA_KOI8U       = "KOI8-U"       // uk
	IANA_WINDOWS1251 = "windows-1251" // ar
	IANA_WINDOWS1256 = "windows-1256" // ar
	IANA_IBM424RTL   = "IBM424_rtl"   // he
	IANA_IBM424LTR   = "IBM424_ltr"   // he
	IANA_IBM420RTL   = "IBM420_rtl"   // ar
	IANA_IBM420LTR   = "IBM420_ltr"   // ar

	IANA_SHIFTJIS  = "Shift_JIS"   // ja
	IANA_GBK       = "GBK"         // zh
	IANA_GB18030   = "GB18030"     // zh
	IANA_BIG5      = "Big5"        // zh
	IANA_EUCJP     = "EUC-JP"      // ja
	IANA_EUCKR     = "EUC-KR"      // ko
	IANA_ISO2022JP = "ISO-2022-JP" // jp
	IANA_ISO2022KR = "ISO-2022-KR" // kr
	IANA_ISO2022CN = "ISO-2022-CN" // cn
)

// GetEncoding returns an encoding.Encoding from some IANA or MIME names.
func GetEncoding(iana string) (encode encoding.Encoding, exist bool) {
	switch iana {
	// utf8, utf16, utf32
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

	// single_byte
	case IANA_ISO8859_1:
		return charmap.ISO8859_1, true
	case IANA_ISO8859_2:
		return charmap.ISO8859_2, true
	case IANA_ISO8859_5:
		return charmap.ISO8859_5, true
	case IANA_ISO8859_6:
		return charmap.ISO8859_6, true
	case IANA_ISO8859_7:
		return charmap.ISO8859_7, true
	case IANA_ISO8859_8:
		return charmap.ISO8859_8, true
	case IANA_ISO8859_8I:
		return charmap.ISO8859_8I, true
	case IANA_ISO8859_9:
		return charmap.ISO8859_9, true
	case IANA_KOI8R:
		return charmap.KOI8R, true
	case IANA_KOI8U:
		return charmap.KOI8U, true
	case IANA_WINDOWS1251:
		return charmap.Windows1251, true
	case IANA_WINDOWS1256:
		return charmap.Windows1256, true
	case IANA_IBM424RTL, IANA_IBM424LTR, IANA_IBM420RTL, IANA_IBM420LTR:
		// not found

	// multi_byte
	case IANA_SHIFTJIS:
		return japanese.ShiftJIS, true
	case IANA_GBK:
		return simplifiedchinese.GBK, true
	case IANA_GB18030:
		return simplifiedchinese.GB18030, true
	case IANA_BIG5:
		return traditionalchinese.Big5, true
	case IANA_EUCJP:
		return japanese.EUCJP, true
	case IANA_EUCKR:
		return korean.EUCKR, true
	case IANA_ISO2022JP:
		return japanese.ISO2022JP, true
	case IANA_ISO2022KR, IANA_ISO2022CN:
		// not found
	}

	// not found
	return nil, false
}
