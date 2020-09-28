package xcharset

import (
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

// DetectResult contains all the information that charset detector gives.
type DetectResult struct {
	// IANA name of the detected charset.
	Charset string

	// IANA name of the detected language. It may be empty for some charsets.
	Language string

	// The confidence of the result. Scale from 1 to 100.
	Confidence int
}

// newDetectResultFromChardet builds a DetectResult from chardet.Result.
func newDetectResultFromChardet(r *chardet.Result) *DetectResult {
	return &DetectResult{
		Charset:    r.Charset,
		Language:   r.Language,
		Confidence: r.Confidence,
	}
}

// DetectCharsetBest returns the Result with highest Confidence.
func DetectCharsetBest(bs []byte) (*DetectResult, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(bs)
	if err != nil {
		return nil, err
	}
	return newDetectResultFromChardet(result), nil
}

// DetectCharsetBest returns all Results which have non-zero Confidence. The Results are sorted by Confidence in descending order.
func DetectCharsetAll(bs []byte) ([]*DetectResult, error) {
	detector := chardet.NewTextDetector()
	results, err := detector.DetectAll(bs)
	if err != nil {
		return nil, err
	}

	out := make([]*DetectResult, len(results))
	for idx := range results {
		out[idx] = newDetectResultFromChardet(&results[idx])
	}
	return out, nil
}

// EncodeString encodes a string in a specific encoding.
func EncodeString(encode encoding.Encoding, src string) (string, error) {
	dest, _, err := transform.String(encode.NewEncoder(), src)
	return dest, err
}

// DecodeString decodes a string in a specific encoding.
func DecodeString(encode encoding.Encoding, src string) (string, error) {
	dest, _, err := transform.String(encode.NewDecoder(), src)
	return dest, err
}

// EncodeBytes encodes a bytes in a specific encoding.
func EncodeBytes(encode encoding.Encoding, src []byte) ([]byte, error) {
	dest, _, err := transform.Bytes(encode.NewEncoder(), src)
	return dest, err
}

// DecodeBytes decodes a bytes in a specific encoding.
func DecodeBytes(encode encoding.Encoding, src []byte) ([]byte, error) {
	dest, _, err := transform.Bytes(encode.NewDecoder(), src)
	return dest, err
}
