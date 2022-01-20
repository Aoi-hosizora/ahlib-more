package main

import (
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/pmezard/go-difflib/difflib"
	"strings"
)

// UnifiedDiffString aligns and colorizes the difflib.UnifiedDiff's compare result string returned from difflib.GetUnifiedDiffString.
func UnifiedDiffString(ud *difflib.UnifiedDiff) string {
	s, _ := difflib.GetUnifiedDiffString(*ud) // ignore error, actually it will not return error
	return colorize(s, ud.Eol, true)
}

// ContextDiffString aligns and colorizes the difflib.ContextDiff's compare result string returned from difflib.GetContextDiffString.
func ContextDiffString(cd *difflib.ContextDiff) string {
	s, _ := difflib.GetContextDiffString(*cd) // ignore error
	return colorize(s, cd.Eol, false)
}

// colorize aligns and colorizes the given string, notes that alignment is only for unified diff, that is to add a space between prefix and content.
func colorize(s, eol string, unified bool) string {
	if len(eol) == 0 {
		eol = "\n"
	}
	lines := strings.Split(s, eol)
	outLines := make([]string, 0, len(lines))
	for _, line := range lines {
		minusColor := xcolor.BrightRed
		plusColor := xcolor.BrightGreen
		exclamationColor := xcolor.BrightYellow
		normalColor := xcolor.BrightWhite

		// 1. skip
		if len(line) == 0 || strings.HasPrefix(line, "*** ") || strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ") ||
			strings.HasPrefix(line, "***************") || strings.HasPrefix(line, "@@ -") {
			outLines = append(outLines, line)
			continue
		}

		// 2. colorize
		if strings.HasPrefix(line, "-") {
			if len(line) == 1 {
				line = "- "
			} else if !unified {
				line = "-" + line[1:]
			} else {
				line = "- " + line[1:]
			}
			line = minusColor.Sprint(line)
		} else if strings.HasPrefix(line, "+") {
			if len(line) == 1 {
				line = "+ "
			} else if !unified {
				line = "+" + line[1:]
			} else {
				line = "+ " + line[1:]
			}
			line = plusColor.Sprint(line)
		} else if strings.HasPrefix(line, "!") {
			line = "!" + line[1:]
			line = exclamationColor.Sprint(line)
		} else if strings.HasPrefix(line, " ") {
			if len(line) == 1 {
				line = "  "
			} else if !unified {
				line = " " + line[1:]
			} else {
				line = "  " + line[1:]
			}
			line = normalColor.Sprint(line)
		}
		outLines = append(outLines, line)
	}
	return strings.Join(outLines, eol)
}
