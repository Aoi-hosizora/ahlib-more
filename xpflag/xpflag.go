package xpflag

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/spf13/pflag"
	"io"
	"os"
	"reflect"
)

// Parse parses flag definitions from the argument list, and this function never panic or call os.Exit.
// Visit pflag.Parse for more details.
func Parse() error {
	errorHandling := xreflect.FieldValueOf(pflag.CommandLine, "errorHandling")
	saved := xreflect.GetUnexportedField(errorHandling)
	xreflect.SetUnexportedField(errorHandling, reflect.ValueOf(pflag.ContinueOnError))
	defer xreflect.SetUnexportedField(errorHandling, saved)
	return pflag.CommandLine.Parse(os.Args[1:])
}

// osStderr always equals to os.Stderr except when testing.
var osStderr io.Writer = os.Stderr

// osExit always equals to os.Exit except when testing.
var osExit = os.Exit

// ParseDefault parses flag definitions from the argument list, if it failed to parse, this function will
// print error message and default usage by pflag.Usage to os.Stderr.
func ParseDefault() {
	err := Parse()
	if err != nil {
		_, _ = fmt.Fprintf(osStderr, "Error, %v\n", err)
		pflag.Usage()
		osExit(2)
	}
}
