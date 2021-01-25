# xcharset

## Dependencies

+ github.com/saintfish/chardet
+ golang.org/x/text
+ xtesting*

## Documents

### Types

+ `type DetectResult struct`

### Variables

+ None

### Constants

+ `const IANA_UTF8 string`
+ `const IANA_UTF16BE string`
+ `const IANA_UTF16LE string`
+ `const IANA_UTF32BE string`
+ `const IANA_UTF32LE string`
+ `const IANA_ISO8859_1 string`
+ `const IANA_ISO8859_2 string`
+ `const IANA_ISO8859_5 string`
+ `const IANA_ISO8859_6 string`
+ `const IANA_ISO8859_7 string`
+ `const IANA_ISO8859_8 string`
+ `const IANA_ISO8859_8I string`
+ `const IANA_ISO8859_9 string`
+ `const IANA_KOI8R string`
+ `const IANA_WINDOWS1251 string`
+ `const IANA_WINDOWS1256 string`
+ `const IANA_IBM424RTL string`
+ `const IANA_IBM424LTR string`
+ `const IANA_IBM420RTL string`
+ `const IANA_IBM420LTR string`
+ `const IANA_SHIFTJIS string`
+ `const IANA_GBK string`
+ `const IANA_GB18030 string`
+ `const IANA_BIG5 string`
+ `const IANA_EUCJP string`
+ `const IANA_EUCKR string`
+ `const IANA_ISO2022JP string`
+ `const IANA_ISO2022KR string`
+ `const IANA_ISO2022CN string`

### Functions

+ `func DetectBestCharset(bs []byte) (*DetectResult, bool)`
+ `func DetectAllCharsets(bs []byte) ([]*DetectResult, bool)`
+ `func EncodeString(encoding encoding.Encoding, s string) (string, error)`
+ `func DecodeString(encoding encoding.Encoding, s string) (string, error)`
+ `func EncodeBytes(encoding encoding.Encoding, bs []byte) ([]byte, error)`
+ `func DecodeBytes(encoding encoding.Encoding, bs []byte) ([]byte, error)`
+ `func GetEncoding(iana string) (encode encoding.Encoding, exist bool)`

### Methods

+ None
