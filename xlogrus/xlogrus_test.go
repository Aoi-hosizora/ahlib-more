package xlogrus

import (
	"github.com/Aoi-hosizora/ahlib-more/xrotation"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestSimpleFormatter(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	sb := &strings.Builder{}
	l.SetOutput(io.MultiWriter(os.Stdout, sb))

	for _, tc := range []struct {
		giveFmt    *SimpleFormatter
		giveReport bool
		giveFn     func(l *logrus.Entry)
		want       string
	}{
		{NewSimpleFormatter(), false, func(l *logrus.Entry) { l.Trace("test") },
			"\x1b[37mTRAC\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithTimestampFormat("")), false, func(l *logrus.Entry) { l.Debug("test") },
			"\x1b[37mDEBU\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithTimestampFormat("2006-01-02 15:04:05")), false, func(l *logrus.Entry) { l.Info("test") },
			"\x1b[34mINFO\x1b[0m [2021-12-29 18:16:13] test\n"},
		{NewSimpleFormatter(), true, func(l *logrus.Entry) { l.Warn("test") },
			"\x1b[33mWARN\x1b[0m [2021-12-29T18:16:13+08:00] xlogrus_test.go:34 xlogrus.TestSimpleFormatter.func4() > test\n"},
		{NewSimpleFormatter(WithCallerFormatter(func(*runtime.Frame) (string, string) { return "", "" })), true, func(l *logrus.Entry) { l.Error("test") },
			"\x1b[31mERRO\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithCallerFormatter(func(*runtime.Frame) (string, string) { return "fn()", "" })), true, func(l *logrus.Entry) { l.Panic("test") },
			"\x1b[31mPANI\x1b[0m [2021-12-29T18:16:13+08:00] fn() > test\n"},
		{NewSimpleFormatter(WithCallerFormatter(func(*runtime.Frame) (string, string) { return "fn()", "a.go:1" })), true, func(l *logrus.Entry) { l.Trace("test") },
			"\x1b[37mTRAC\x1b[0m [2021-12-29T18:16:13+08:00] a.go:1 fn() > test\n"},
		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return l.String() })), false, func(l *logrus.Entry) { l.Debug("test") },
			"\x1b[37mdebug\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return strings.ToUpper(l.String())[:1] })), false, func(l *logrus.Entry) { l.Info("test") },
			"\x1b[34mI\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return "" })), false, func(l *logrus.Entry) { l.Warn("test") },
			"\x1b[33mWARN\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithDisableColor(true)), false, func(l *logrus.Entry) { l.Error("test") },
			"ERRO [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithUseUTCTime(true)), false, func(l *logrus.Entry) { l.Panic("test") },
			"\x1b[31mPANI\x1b[0m [2021-12-29T10:16:13Z] test\n"},
	} {
		l.SetFormatter(tc.giveFmt)
		l.SetReportCaller(tc.giveReport)
		sb.Reset()

		entry := l.WithTime(time.Date(2021, 12, 29, 18, 16, 13, 0, time.Local)).WithField("key", "value")
		func() {
			defer func() { recover() }()
			tc.giveFn(entry)
		}()
		xtesting.Equal(t, sb.String(), tc.want)
	}
}

func TestRotationHook(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339, ForceColors: true})
	l.SetOutput(os.Stdout)

	type Os = []RotationHookOption
	rl, err := xrotation.New(xrotation.WithFilenamePattern("logger.log"))
	xtesting.Nil(t, err)
	rl2, err := xrotation.New(xrotation.WithFilenamePattern("console.log"))
	xtesting.Nil(t, err)

	for _, tc := range []struct {
		giveOptions  Os
		giveFn       func(l *logrus.Entry)
		wantErr      bool
		wantLevels   []logrus.Level
		wantFilename string
		wantContent  string
	}{
		{Os{}, nil, true, nil, "", ""},
		{Os{WithRotationLogger(nil), WithLevel(logrus.TraceLevel)}, nil, true, nil, "", ""},
		{Os{WithRotationLogger(rl), WithLevel(logrus.TraceLevel)}, func(l *logrus.Entry) { l.Trace("test") },
			false, logrus.AllLevels, "logger.log", "{\"key\":\"value\",\"level\":\"trace\",\"msg\":\"test\",\"time\":\"2021-12-29T18:16:13+08:00\"}\n"},
		{Os{WithRotationLogger(rl2), WithLevel(logrus.DebugLevel), WithFormatter(&logrus.JSONFormatter{})}, func(l *logrus.Entry) { l.Debug("test") },
			false, logrus.AllLevels[:logrus.DebugLevel+1], "console.log", "{\"key\":\"value\",\"level\":\"debug\",\"msg\":\"test\",\"time\":\"2021-12-29T18:16:13+08:00\"}\n"},
		{Os{WithRotationLogger(rl), WithFormatter(&logrus.TextFormatter{})}, func(l *logrus.Entry) { l.Info("test") },
			false, logrus.AllLevels[:logrus.WarnLevel+1], "", ""},
		{Os{WithRotationLogger(rl), WithFormatter(&logrus.TextFormatter{DisableColors: true})}, func(l *logrus.Entry) { l.Warn("test") },
			false, logrus.AllLevels[:logrus.WarnLevel+1], "logger.log", "time=\"2021-12-29T18:16:13+08:00\" level=warning msg=test key=value\n"},
		{Os{WithRotationLogger(rl2), WithLevel(999999), WithFormatter(NewSimpleFormatter(WithDisableColor(true)))}, func(l *logrus.Entry) { l.Panic("test") },
			false, []logrus.Level{logrus.PanicLevel}, "console.log", "PANI [2021-12-29T18:16:13+08:00] test\n"},
	} {
		hook, err := NewRotationHook(tc.giveOptions...)
		xtesting.Equal(t, err != nil, tc.wantErr)
		if err != nil {
			continue
		}
		l.Hooks = logrus.LevelHooks{}
		l.AddHook(hook)
		xtesting.Equal(t, hook.Levels(), tc.wantLevels)
		if tc.wantFilename != "" {
			_ = os.Remove(tc.wantFilename)
		}

		entry := l.WithTime(time.Date(2021, 12, 29, 18, 16, 13, 0, time.Local)).WithField("key", "value")
		func() {
			defer func() { recover() }()
			tc.giveFn(entry)
		}()
		_, err = os.Stat(tc.wantFilename)
		if tc.wantFilename == "" {
			xtesting.True(t, os.IsNotExist(err))
		} else {
			xtesting.Nil(t, err)
			bs, err := ioutil.ReadFile(tc.wantFilename)
			xtesting.Nil(t, err)
			xtesting.Equal(t, string(bs), tc.wantContent)
		}

		if o := hook.option.rotationLogger; o != nil {
			_ = o.Close()
		}
		if tc.wantFilename != "" {
			_ = os.Remove(tc.wantFilename)
		}
	}
}
