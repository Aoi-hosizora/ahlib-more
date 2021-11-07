package xcharset

import (
	"github.com/saintfish/chardet"
)

// DetectResult contains the information for charset detector. See chardet.Result.
type DetectResult struct {
	// Charset represents IANA or MIME name of the detected charset.
	Charset string

	// Language represents IANA name of the detected language. It may be empty for some charsets.
	Language string

	// Confidence represents the confidence of the result. Scale from 1 to 100.
	Confidence int
}

// DetectBestCharset detects bytes and returns the charset result with the highest confidence.
func DetectBestCharset(bs []byte) (*DetectResult, bool) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(bs)
	if err != nil {
		return nil, false // empty result
	}

	return newDetectResultFromChardet(result), true
}

// DetectAllCharsets detects bytes and returns all charsets in confidence's descending order.
func DetectAllCharsets(bs []byte) ([]*DetectResult, bool) {
	detector := chardet.NewTextDetector()
	results, err := detector.DetectAll(bs)
	if err != nil {
		return nil, false // empty result
	}

	out := make([]*DetectResult, len(results))
	for idx := range results {
		out[idx] = newDetectResultFromChardet(&results[idx])
	}
	return out, true
}

// newDetectResultFromChardet creates a DetectResult from chardet.Result, note that there are some bugs in `chardet` package.
func newDetectResultFromChardet(r *chardet.Result) *DetectResult {
	charset := r.Charset
	language := r.Language

	switch charset {
	// case "ISO-8859-1":
	// 	switch language {
	// 	case "cs", "hu", "pl", "ro":
	// 		charset = "ISO-8859-2"
	// 	}
	case "GB-18030":
		charset = "GB18030"
	case "ISO-2022-JP":
		language = "ja"
		// case "ISO-2022-KR":
		// 	language = "ko"
		// case "ISO-2022-CN":
		// 	language = "cn"
	}

	return &DetectResult{Charset: charset, Language: language, Confidence: r.Confidence}
}
