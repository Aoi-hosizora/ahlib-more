package xlogrus

import (
	"bytes"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// simpleFormatterOptions is a type of SimpleFormatter's option, each field can be set by SimpleFormatterOption function type.
type simpleFormatterOptions struct {
	timestampFormat string
	callerFormatter func(*runtime.Frame) (function string, file string)
	levelFormatter  func(logrus.Level) string
	disableColor    bool
	useUTCTime      bool
}

// SimpleFormatterOption represents an option type for SimpleFormatter's option, can be created by WithXXX functions.
type SimpleFormatterOption func(*simpleFormatterOptions)

// WithTimestampFormat creates an SimpleFormatterOption to specific timestamp format, defaults to time.RFC3339.
func WithTimestampFormat(format string) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.timestampFormat = format
	}
}

// WithCallerFormatter creates an SimpleFormatterOption to specific the caller's runtime.Frame formatter, defaults to use the function shortname and filename without path.
func WithCallerFormatter(formatter func(*runtime.Frame) (function string, file string)) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.callerFormatter = formatter
	}
}

// WithLevelFormatter creates an SimpleFormatterOption to specific the logrus.Level formatter, defaults to use the first four character in capital of the level.
func WithLevelFormatter(formatter func(logrus.Level) string) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.levelFormatter = formatter
	}
}

// WithDisableColor creates an SimpleFormatterOption to disable the colored format, defaults to false, and means defaults to enable colored format.
func WithDisableColor(disable bool) SimpleFormatterOption {
	return func(c *simpleFormatterOptions) {
		c.disableColor = disable
	}
}

// WithUseUTCTime creates an SimpleFormatterOption to specific use the time.UTC layout or not, defaults to false, and means defaults to use time.Local layout.
func WithUseUTCTime(use bool) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.useUTCTime = use
	}
}

// SimpleFormatter represents a simple formatter for logrus.Logger, it only formats level, time, caller and message information with color or without color, also notes that
// the logrus.Fields data will not be formatted unlink what logrus.TextFormatter does.
type SimpleFormatter struct {
	option       *simpleFormatterOptions
	terminalOnce sync.Once
}

var _ logrus.Formatter = (*SimpleFormatter)(nil)

// NewSimpleFormatter creates an SimpleFormatter with given SimpleFormatterOption-s.
//
// Example:
// 	l := logrus.New()
// 	l.SetLevel(logrus.TraceLevel)
// 	l.SetReportCaller(true)
// 	l.SetFormatter(NewSimpleFormatter(
// 		WithTimestampFormat("2006-01-02 15:04:05"),
// 		WithCallerFormatter(func(*runtime.Frame) (string, string) { return "", "" }), // can use to disable report caller
// 		WithLevelFormatter(func(l logrus.Level) string { return strings.ToUpper(l.String())[:1] }),
// 		WithDisableColor(false), // defaults to false
// 		WithUseUTCTime(true), // defaults to false
// 	))
func NewSimpleFormatter(options ...SimpleFormatterOption) *SimpleFormatter {
	opt := &simpleFormatterOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.timestampFormat == "" {
		opt.timestampFormat = time.RFC3339
	}
	return &SimpleFormatter{option: opt}
}

// initOnce initializes the terminal io.Writer for color supported, this method will be called only once.
func (s *SimpleFormatter) initOnce(entry *logrus.Entry) {
	s.terminalOnce.Do(func() {
		if entry.Logger != nil && !s.option.disableColor {
			xcolor.InitTerminal(entry.Logger.Out)
		}
	})
}

// Format formats a single logrus.Entry, this method implements logrus.Formatter interface.
//
// Logs like:
// 	WARN [2021-08-29T05:56:25+08:00] test
// 	INFO [2021-08-29T05:56:25+08:00] a.go:1 fn() > test
func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s.initOnce(entry)

	// 1. time and message
	level := ""
	if f := s.option.levelFormatter; f != nil {
		level = f(entry.Level)
	}
	if level == "" {
		level = strings.ToUpper(entry.Level.String()[0:4])
	}
	t := entry.Time
	if s.option.useUTCTime {
		t = t.UTC()
	} else {
		t = t.Local()
	}
	now := t.Format(s.option.timestampFormat)
	message := strings.TrimSuffix(entry.Message, "\n")

	// 2. runtime caller
	caller := ""
	if entry.HasCaller() {
		var funcName, filename string
		if f := s.option.callerFormatter; f != nil {
			funcName, filename = f(entry.Caller)
		} else {
			_, funcName = filepath.Split(entry.Caller.Function)
			_, filename = filepath.Split(entry.Caller.File)
			funcName = fmt.Sprintf("%s()", funcName)
			filename = fmt.Sprintf("%s:%d", filename, entry.Caller.Line)
		}
		parts := make([]string, 0, 3)
		if filename != "" {
			parts = append(parts, filename)
		}
		if funcName != "" {
			parts = append(parts, funcName)
		}
		if len(parts) > 0 {
			parts = append(parts, ">")
			caller = " " + strings.Join(parts, " ")
		}
	}

	// *. write to buffer
	buf := &bytes.Buffer{}
	if entry.Buffer != nil {
		buf = entry.Buffer
	}
	levelString := level
	if !s.option.disableColor {
		levelString = s.levelColor(entry.Level).Sprintf(level)
	}
	_, _ = fmt.Fprintf(buf, "%s [%s]%s %s\n", levelString, now, caller, message) // ignore error
	return buf.Bytes(), nil
}

// levelColor returns the xcolor.Color from logrus.Level.
func (s *SimpleFormatter) levelColor(level logrus.Level) xcolor.Color {
	switch level {
	case logrus.InfoLevel:
		return xcolor.Blue
	case logrus.WarnLevel:
		return xcolor.Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return xcolor.Red
	default: // debug, trace
		return xcolor.White
	}
}
