package xlogrus

import (
	"fmt"
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

	d := time.Date(2021, 12, 29, 18, 16, 13, 0, time.FixedZone("", 8*60*60))
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
		{NewSimpleFormatter(WithUseUTCTime(true)), false, func(l *logrus.Entry) { l.Warn("test") },
			"\x1b[33mWARN\x1b[0m [2021-12-29T10:16:13Z] test\n"},
		{NewSimpleFormatter(WithDisableColor(false)), false, func(l *logrus.Entry) { l.Error("test") },
			"\x1b[31mERRO\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithTimestampFormat("2006-01-02"), WithDisableColor(true), WithUseUTCTime(false)), false, func(l *logrus.Entry) { l.Panic("test") },
			"PANI [2021-12-29] test\n"},

		{NewSimpleFormatter(), true, func(l *logrus.Entry) { l.Trace("test") },
			"\x1b[37mTRAC\x1b[0m [2021-12-29T18:16:13+08:00] \x1b[2mxlogrus_test.go:43 xlogrus.TestSimpleFormatter.func7()\x1b[0m > test\n"},
		{NewSimpleFormatter(WithCallerFormatter(func(*runtime.Frame) (string, string) { return "", "" })), true, func(l *logrus.Entry) { l.Debug("test") },
			"\x1b[37mDEBU\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithCallerFormatter(func(*runtime.Frame) (string, string) { return "", "fn()" }), WithDisableColor(true)), true, func(l *logrus.Entry) { l.Info("test") },
			"INFO [2021-12-29T18:16:13+08:00] fn() > test\n"},
		{NewSimpleFormatter(WithCallerFormatter(func(*runtime.Frame) (string, string) { return "aaa.go:123", "fff()" })), true, func(l *logrus.Entry) { l.Warn("test") },
			"\x1b[33mWARN\x1b[0m [2021-12-29T18:16:13+08:00] \x1b[2maaa.go:123 fff()\x1b[0m > test\n"},

		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return l.String() })), false, func(l *logrus.Entry) { l.Trace("test") },
			"\x1b[37mtrace\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return "" })), false, func(l *logrus.Entry) { l.Debug("test") },
			"\x1b[37mDEBU\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return l.String()[:1] })), false, func(l *logrus.Entry) { l.Info("test") },
			"\x1b[34mi\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithLevelFormatter(func(l logrus.Level) string { return strings.ToUpper(l.String())[:1] }), WithDisableColor(true)), false, func(l *logrus.Entry) { l.Error("test") },
			"E [2021-12-29T18:16:13+08:00] test\n"},

		{NewSimpleFormatter(WithMessageFormatter(func(l, t, c, m string) string { return fmt.Sprintf("%s %s: %s\n", t, l, m) }), WithTimestampFormat("2006-01-02 15:04:05")), false, func(l *logrus.Entry) { l.Trace("test") },
			"2021-12-29 18:16:13 \x1b[37mTRAC\x1b[0m: test\n"},
		{NewSimpleFormatter(WithMessageFormatter(func(l, t, c, m string) string { return "" })), false, func(l *logrus.Entry) { l.Debug("test") },
			"\x1b[37mDEBU\x1b[0m [2021-12-29T18:16:13+08:00] test\n"},
		{NewSimpleFormatter(WithMessageFormatter(func(l, t, c, m string) string { return fmt.Sprintf("[%s] [%s] %s\n", t, l, m) }), WithDisableColor(true)), false, func(l *logrus.Entry) { l.Info("test") },
			"[2021-12-29T18:16:13+08:00] [INFO] test\n"},
		{NewSimpleFormatter(WithMessageFormatter(func(l, t, c, m string) string { return fmt.Sprintf("%s[%s] %s\n", l, t, m) }), WithLevelFormatter(func(l logrus.Level) string { return strings.ToUpper(l.String())[:1] })), false, func(l *logrus.Entry) { l.Warn("test") },
			"\x1b[33mW\x1b[0m[2021-12-29T18:16:13+08:00] test\n"},
	} {
		l.SetFormatter(tc.giveFmt)
		l.SetReportCaller(tc.giveReport)
		sb.Reset()

		entry := l.WithTime(d).WithField("key", "value")
		func() {
			defer func() { recover() }()
			tc.giveFn(entry)
		}()
		xtesting.Equal(t, sb.String(), tc.want)
	}
}

func TestLogrusFormatters(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	sb := &strings.Builder{}
	l.SetOutput(io.MultiWriter(os.Stdout, sb))

	d := time.Date(2022, 1, 20, 14, 56, 20, 123456789, time.Local)
	for _, tc := range []struct {
		giveFmt    logrus.Formatter
		giveFunc   func(*logrus.Entry, ...interface{})
		giveFields logrus.Fields
		giveMsg    string
		giveReport bool
		want       string
	}{
		{RFC3339JsonFormatter(), (*logrus.Entry).Debug, logrus.Fields{}, "hello world 1", false,
			`{"@@time":"2022-01-20T14:56:20.123456789+08:00","@level":"debug","@message":"hello world 1","entries":{}}`},
		{RFC3339JsonFormatter(), (*logrus.Entry).Info, logrus.Fields{"key": "v"}, "hello world 2", false,
			`{"@@time":"2022-01-20T14:56:20.123456789+08:00","@level":"info","@message":"hello world 2","entries":{"key":"v"}}`},
		{RFC3339JsonFormatter(), (*logrus.Entry).Warn, logrus.Fields{"empty": ""}, "hello world 3", false,
			`{"@@time":"2022-01-20T14:56:20.123456789+08:00","@level":"warning","@message":"hello world 3","entries":{"empty":""}}`},
		{RFC3339JsonFormatter(), (*logrus.Entry).Error, logrus.Fields{"i": 123, "f": 3.14}, "hello world 4", true,
			`{"@@time":"2022-01-20T14:56:20.123456789+08:00","@level":"error","@message":"hello world 4","entries":{"f":3.14,"i":123}`},

		{RFC3339ColoredTextFormatter(), (*logrus.Entry).Debug, logrus.Fields{}, "hello world 1", false, ""},
		{RFC3339ColoredTextFormatter(), (*logrus.Entry).Info, logrus.Fields{"key": "v"}, "hello world 2", false, ""},
		{RFC3339ColoredTextFormatter(), (*logrus.Entry).Warn, logrus.Fields{"empty": ""}, "hello world 3", false, ""},
		{RFC3339ColoredTextFormatter(), (*logrus.Entry).Error, logrus.Fields{"i": 123, "f": 3.14}, "hello world 4", true, ""},
	} {
		t.Run(tc.giveMsg, func(t *testing.T) {
			l.SetFormatter(tc.giveFmt)
			l.SetReportCaller(tc.giveReport)

			sb.Reset()
			tc.giveFunc(l.WithTime(d).WithFields(tc.giveFields), tc.giveMsg)
			l.SetReportCaller(false)
			xtesting.Equal(t, sb.String()[:len(tc.want)], tc.want)
		})
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
		giveRotation *xrotation.RotationLogger
		giveOptions  Os
		giveFn       func(l *logrus.Entry)
		wantPanic    bool
		wantLevels   []logrus.Level
		wantFilename string
		wantContent  string
	}{
		{nil, Os{}, nil, true, nil, "", ""},
		{nil, Os{WithRotateLevel(logrus.TraceLevel)}, nil, true, nil, "", ""},
		{rl, Os{WithRotateLevel(logrus.TraceLevel)}, func(l *logrus.Entry) { l.Trace("test") },
			false, logrus.AllLevels, "logger.log", "{\"key\":\"value\",\"level\":\"trace\",\"msg\":\"test\",\"time\":\"2021-12-29T18:16:13+08:00\"}\n"},
		{rl2, Os{WithRotateLevel(logrus.DebugLevel), WithRotateFormatter(&logrus.JSONFormatter{})}, func(l *logrus.Entry) { l.Debug("test") },
			false, logrus.AllLevels[:logrus.DebugLevel+1], "console.log", "{\"key\":\"value\",\"level\":\"debug\",\"msg\":\"test\",\"time\":\"2021-12-29T18:16:13+08:00\"}\n"},
		{rl, Os{WithRotateFormatter(&logrus.TextFormatter{})}, func(l *logrus.Entry) { l.Info("test") },
			false, logrus.AllLevels[:logrus.WarnLevel+1], "", ""},
		{rl, Os{WithRotateFormatter(&logrus.TextFormatter{DisableColors: true})}, func(l *logrus.Entry) { l.Warn("test") },
			false, logrus.AllLevels[:logrus.WarnLevel+1], "logger.log", "time=\"2021-12-29T18:16:13+08:00\" level=warning msg=test key=value\n"},
		{rl2, Os{WithRotateLevel(999999), WithRotateFormatter(NewSimpleFormatter(WithDisableColor(true)))}, func(l *logrus.Entry) { l.Panic("test") },
			false, []logrus.Level{logrus.PanicLevel}, "console.log", "PANI [2021-12-29T18:16:13+08:00] test\n"},
	} {
		if tc.wantPanic {
			xtesting.Panic(t, func() { NewRotationHook(tc.giveRotation, tc.giveOptions...) })
			continue
		}
		hook := NewRotationHook(tc.giveRotation, tc.giveOptions...)
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
		_, err := os.Stat(tc.wantFilename)
		if tc.wantFilename == "" {
			xtesting.True(t, os.IsNotExist(err))
		} else {
			xtesting.Nil(t, err)
			bs, err := ioutil.ReadFile(tc.wantFilename)
			xtesting.Nil(t, err)
			xtesting.Equal(t, string(bs), tc.wantContent)
		}

		if o := hook.rotation; o != nil {
			_ = o.Close()
		}
		if tc.wantFilename != "" {
			_ = os.Remove(tc.wantFilename)
		}
	}
}
