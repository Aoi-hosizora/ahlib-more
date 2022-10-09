package xpflag

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"io"
	"os"
)

// _cmd is a global pflag.FlagSet.
var _cmd = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

func init() {
	_cmd.Usage = func() { DefaultUsage(_cmd) }
}

// Cmd returns the global pflag.FlagSet. This should be used if want to use Parse, MustParse or PrintUsage.
func Cmd() *pflag.FlagSet {
	return _cmd
}

// DefaultUsage is the default usage function, which prints the default usage string to os.Stderr.
func DefaultUsage(cmd *pflag.FlagSet) {
	fmt.Fprintf(_osStderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(_osStderr, "%s", cmd.FlagUsages())
}

// PrintUsage prints the default usage string to os.Stderr.
func PrintUsage() {
	if _cmd.Usage != nil {
		_cmd.Usage()
	}
}

// Parse parses flags passed to the program and returns error, note that this function never panic or exit the program.
func Parse() error {
	return _cmd.Parse(os.Args[1:]) // may return "pflag.ErrHelp"
}

// _osStderr always equals to os.Stderr except when testing.
var _osStderr io.Writer = os.Stderr

// _osExit always equals to os.Exit except when testing.
var _osExit = os.Exit

// MustParse parses flags passed to the program, prints error message and exits the program when failed to parse.
func MustParse() {
	err := Parse()
	if err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			_osExit(0)
			return
		}
		fmt.Fprintf(_osStderr, "Error: %v\n", err)
		fmt.Fprintf(_osStderr, "Run '%s --help' for usage.\n", os.Args[0])
		_osExit(2)
	}
}
