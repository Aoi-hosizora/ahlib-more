package xpflag

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	saved := _cmd
	defer func() { _cmd = saved }()
	saved2 := os.Args
	defer func() { os.Args = saved2 }()
	saved3 := _osStderr
	defer func() { _osStderr = saved3 }()
	saved4 := _osExit
	defer func() { _osExit = saved4 }()

	var pHelp *bool
	var pConfig *string
	define := func(includeError bool) {
		_cmd = pflag.NewFlagSet("application", pflag.ContinueOnError) // avoid to influence global _cmd
		_cmd.Usage = func() { DefaultUsage(_cmd) }
		if includeError {
			pHelp = Cmd().BoolP("help", "h", false, "show help message")
		}
		pConfig = Cmd().StringP("config", "c", "./config.json", "config file path")
	}
	sb := &strings.Builder{}
	_osStderr = sb
	exitCode := -1
	_osExit = func(code int) {
		exitCode = code
	}

	// 1. PrintUsage
	xtesting.EmptyCollection(t, strings.TrimSpace(Cmd().FlagUsages()))
	PrintUsage()
	define(false)
	xtesting.Equal(t, len(strings.Split(strings.TrimSpace(Cmd().FlagUsages()), "\n")), 1)
	PrintUsage()
	define(true)
	xtesting.Equal(t, len(strings.Split(strings.TrimSpace(Cmd().FlagUsages()), "\n")), 2)
	PrintUsage()

	// => test normal parse
	for _, tc := range []struct {
		giveArgs  []string
		wantError bool
		checkFn   func()
	}{
		{[]string{}, false, func() {
			xtesting.Equal(t, *pHelp, false)
			xtesting.Equal(t, *pConfig, "./config.json")
		}},
		{[]string{"-x"}, true, nil},
		{[]string{"-xxx"}, true, nil},
		{[]string{"--xxx"}, true, nil},
		{[]string{"-h", "-1"}, true, nil},
		{[]string{"-x", "-c"}, true, nil},
		{[]string{"-h", "123"}, false, func() {
			xtesting.Equal(t, *pHelp, true)
			xtesting.Equal(t, *pConfig, "./config.json")
		}},
		{[]string{"--help", "-h", "-c", "xxx"}, false, func() {
			xtesting.Equal(t, *pHelp, true)
			xtesting.Equal(t, *pConfig, "xxx")
		}},
		{[]string{"-c", "--help"}, false, func() {
			xtesting.Equal(t, *pHelp, false)
			xtesting.Equal(t, *pConfig, "--help")
		}},
	} {
		t.Run(strings.Join(tc.giveArgs, " "), func(t *testing.T) {
			os.Args = append([]string{"application"}, tc.giveArgs...)
			define(true)

			// 2. Parse
			err := Parse()
			xtesting.Equal(t, err != nil, tc.wantError)
			if err == nil && tc.checkFn != nil {
				tc.checkFn()
			}

			// 3. MustParse
			sb.Reset()
			exitCode = -1
			MustParse()
			if !tc.wantError {
				xtesting.Equal(t, sb.String(), "")
				xtesting.Equal(t, exitCode, -1) // success
			} else {
				log.Println(err)
				firstLine := strings.Split(sb.String(), "\n")[0]
				xtesting.Equal(t, firstLine, fmt.Sprintf("Error: %v", err))
				xtesting.Equal(t, exitCode, 2)
			}
		})
	}

	// => test parse without help flag
	for _, tc := range []struct {
		giveArgs      []string
		wantHelpError bool
	}{
		{[]string{""}, false},
		{[]string{"-h"}, true},
		{[]string{"--help"}, true},
		{[]string{"-c -h"}, false},
		{[]string{"-h -c"}, true},
	} {
		t.Run(strings.Join(tc.giveArgs, " "), func(t *testing.T) {
			os.Args = append([]string{"application"}, tc.giveArgs...)
			define(false)

			// 4. Parse (for help)
			err := Parse()
			xtesting.Equal(t, errors.Is(err, pflag.ErrHelp), tc.wantHelpError)
			log.Println(err) // may be "pflag: help requested"

			// 5. MustParse (for help)
			sb.Reset()
			exitCode = -1
			MustParse()
			if !tc.wantHelpError {
				xtesting.Equal(t, exitCode, -1) // success
			} else {
				xtesting.Equal(t, exitCode, 0)
				xtesting.Equal(t, strings.Contains(sb.String(), pflag.ErrHelp.Error()), false)
			}
		})
	}
}
