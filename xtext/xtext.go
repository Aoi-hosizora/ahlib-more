package xtext

import (
	"github.com/saintfish/chardet"
	_ "golang.org/x/text"
)

// https://github.com/Aoi-hosizora/common_private_api/blob/master/src/service/http.go

type DetectResult struct {
	Charset    string
	Language   string
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
