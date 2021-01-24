package xlogrus

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSimpleFormatter(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	sb := &strings.Builder{}
	l.SetOutput(sb)

	for _, tc := range []struct {
		giveFmt      *SimpleFormatter
		giveReporter bool
		giveFn       func()
		want1        string
		want2        string
	}{
		{&SimpleFormatter{}, false,
			func() { l.Trace("test") }, "\x1b[37mTRAC\x1b[0m [", "] test\n"},
		{&SimpleFormatter{TimestampFormat: time.RFC3339}, false,
			func() { l.Error("test") }, "\x1b[31mERRO\x1b[0m [", "] test\n"},
		{&SimpleFormatter{RuntimeCaller: func(*runtime.Frame) (string, string) { return "", "" }}, true,
			func() { l.Warn("test") }, "\x1b[33mWARN\x1b[0m [", "] test\n"},
		{&SimpleFormatter{RuntimeCaller: func(*runtime.Frame) (string, string) { return "fn()", "a.go:1" }}, true,
			func() { l.Info("test") }, "\x1b[34mINFO\x1b[0m [", "] a.go:1 fn() test\n"},
		{&SimpleFormatter{DisableColor: true}, false,
			func() { l.Debug("test") }, "DEBU [", "] test\n"},
		{&SimpleFormatter{RuntimeCaller: nil, DisableColor: true}, true,
			func() { l.Debug("test") }, "", ""}, // ignore test
	} {
		l.SetFormatter(tc.giveFmt)
		l.SetReportCaller(tc.giveReporter)
		sb.Reset()
		tc.giveFn()
		output := sb.String()

		if !tc.giveReporter || tc.giveFmt.RuntimeCaller != nil {
			len1 := len(tc.want1)
			len2 := len(tc.want2)
			xtesting.Equal(t, output[:len1], tc.want1)
			xtesting.Equal(t, output[len(output)-len2:], tc.want2)
		}
	}
}

func TestRotateFileHook(t *testing.T) {
	for _, tc := range []struct {
		giveCfg   *RotateFileConfig
		wantPanic bool
	}{
		{nil, true},
		{&RotateFileConfig{Filename: ""}, true},
		{&RotateFileConfig{Filename: "log"}, false},
		{&RotateFileConfig{Filename: "log", Level: 20}, true},
		{&RotateFileConfig{Filename: "log", Level: logrus.WarnLevel}, false},
		{&RotateFileConfig{Filename: "log", Formatter: nil}, false},
		{&RotateFileConfig{Filename: "log", Formatter: &logrus.JSONFormatter{}}, false},
	} {
		if tc.wantPanic {
			xtesting.Panic(t, func() { NewRotateFileHook(tc.giveCfg) })
		} else {
			NewRotateFileHook(tc.giveCfg)
		}
	}

	hook := NewRotateFileHook(&RotateFileConfig{
		Filename:  "console.log",
		Level:     logrus.WarnLevel,
		Formatter: &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
		MaxAge:    30,
		MaxSize:   100,
		LocalTime: false,
		Compress:  true,
	})
	xtesting.Equal(t, hook.Levels(), []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel})
	xtesting.Nil(t, hook.Fire(logrus.WithField("key", "value")))
}

func TestRotateLogHook(t *testing.T) {
	for _, tc := range []struct {
		giveCfg   *RotateLogConfig
		wantPanic bool
	}{
		{nil, true},
		{&RotateLogConfig{Filename: ""}, true},
		{&RotateLogConfig{Filename: "log"}, false},
		{&RotateLogConfig{Filename: "log", Level: 20}, true},
		{&RotateLogConfig{Filename: "log", Level: logrus.WarnLevel}, false},
		{&RotateLogConfig{Filename: "log", FilenameTimePart: "%"}, true},
		{&RotateLogConfig{Filename: "log", FilenameTimePart: ".%Y%m%d.log"}, false},
		{&RotateLogConfig{Filename: "log", LocalTime: true}, false},
		{&RotateLogConfig{Filename: "log", LocalTime: false}, false},
		{&RotateLogConfig{Filename: "log", ForceNewFile: true}, false},
		{&RotateLogConfig{Filename: "log", ForceNewFile: false}, false},
		{&RotateLogConfig{Filename: "log", Formatter: nil}, false},
		{&RotateLogConfig{Filename: "log", Formatter: &logrus.JSONFormatter{}}, false},
	} {
		if tc.wantPanic {
			xtesting.Panic(t, func() { NewRotateLogHook(tc.giveCfg) })
		} else {
			NewRotateLogHook(tc.giveCfg)
		}
	}

	hook := NewRotateLogHook(&RotateLogConfig{
		Filename:         "console",
		FilenameTimePart: ".%Y%m%d.log",
		LinkFileName:     "console.curr.log",
		Level:            logrus.WarnLevel,
		Formatter:        &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
		MaxAge:           time.Hour * 24 * 30,
		MaxSize:          100,
		RotationTime:     time.Hour * 24,
		LocalTime:        false,
		ForceNewFile:     false,
	})
	xtesting.Equal(t, hook.Levels(), []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel})
	xtesting.Nil(t, hook.Fire(logrus.WithField("key", "value")))
}
