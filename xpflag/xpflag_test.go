package xpflag

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/spf13/pflag"
	"log"
	"os"
	"strings"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	field := xreflect.FieldValueOf(pflag.CommandLine, "errorHandling")
	errorHandling := xreflect.GetUnexportedField(field).Interface().(pflag.ErrorHandling)
	xtesting.Equal(t, errorHandling, pflag.ExitOnError)
	xtesting.Nil(t, Parse())
	xtesting.Equal(t, errorHandling, pflag.ExitOnError)
}

func TestParse(t *testing.T) {
	saved := pflag.CommandLine
	defer func() { pflag.CommandLine = saved }()
	saved2 := os.Args
	defer func() { os.Args = saved2 }()
	saved3 := osExit
	defer func() { osExit = saved3 }()
	saved4 := osStderr
	defer func() { osStderr = saved4 }()
	saved5 := pflag.Usage
	defer func() { pflag.Usage = saved5 }()

	var pHelp *bool
	var pConfig *string
	define := func() {
		pflag.CommandLine = pflag.NewFlagSet("application", pflag.ContinueOnError)
		pflag.Usage = func() {}
		pHelp = pflag.CommandLine.BoolP("help", "h", false, "show help message")
		pConfig = pflag.CommandLine.StringP("config", "c", "./config.json", "config file path")
	}

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
			define()
			err := Parse()
			xtesting.Equal(t, err != nil, tc.wantError)
			if err == nil && tc.checkFn != nil {
				tc.checkFn()
			}

			sb := &strings.Builder{}
			osStderr = sb
			exitCode := 0
			osExit = func(code int) {
				exitCode = code
			}
			ParseDefault()
			if err == nil {
				xtesting.Equal(t, sb.String(), "")
				xtesting.Equal(t, exitCode, 0)
			} else {
				log.Println(err)
				xtesting.Equal(t, sb.String(), fmt.Sprintf("Error, %v\n", err))
				xtesting.Equal(t, exitCode, 2)
			}
		})
	}
}
