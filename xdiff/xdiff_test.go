package main

import (
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"testing"
)

func TestDiff(t *testing.T) {
	ud := difflib.UnifiedDiff{
		A:        difflib.SplitLines("0 a\n1 b\n2 c\n\n\n3 d \n4 e\n5 f\n6 g"),
		B:        difflib.SplitLines("1 b\n\n2 c\n\n3 d\n5 f\n7 h\n6 g\n8 ii"),
		FromFile: "A",
		FromDate: "date_a",
		ToFile:   "B",
		ToDate:   "date_b",
		Context:  2,
	}
	text, _ := difflib.GetUnifiedDiffString(ud)
	// fmt.Println(text)
	// fmt.Print("=============================================\n\n")
	text = UnifiedDiffString(&ud)
	fmt.Println(text)
	fmt.Print("=============================================\n\n")

	cd := difflib.ContextDiff(ud)
	text, _ = difflib.GetContextDiffString(cd)
	// fmt.Println(text)
	// fmt.Print("=============================================\n\n")
	text = ContextDiffString(&cd)
	fmt.Println(text)
	fmt.Print("=============================================\n\n")
}
