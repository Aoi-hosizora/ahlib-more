# xcharset

### Dependencies

+ golang.org/x/text
+ github.com/saintfish
+ xtesting*

### Functions

+ `type DetectResult struct {}`
+ `DetectCharsetBest(bs []byte) (*DetectResult, error)`
+ `DetectCharsetAll(bs []byte) ([]*DetectResult, error)`
+ `EncodeString(encode encoding.Encoding, src string) (string, error)`
+ `DecodeString(encode encoding.Encoding, src string) (string, error)`
+ `EncodeBytes(encode encoding.Encoding, src []byte) ([]byte, error)`
+ `DecodeBytes(encode encoding.Encoding, src []byte) ([]byte, error)`
+ `TrimBomString(str string) string`
+ `TrimBomBytes(bs []byte) []byte`
+ `GetEncoding(iana string) (encode encoding.Encoding, existed bool)`
