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
	timestampFormat  string
	useUTCTime       bool
	disableColor     bool
	callerFormatter  func(*runtime.Frame) (filename string, funcname string)
	levelFormatter   func(logrus.Level) string
	messageFormatter func(level, time, caller, message string) string
}

// SimpleFormatterOption represents an option type for SimpleFormatter's option, can be created by WithXXX functions.
type SimpleFormatterOption func(*simpleFormatterOptions)

// WithTimestampFormat creates an SimpleFormatterOption to specific timestamp format, defaults to time.RFC3339.
func WithTimestampFormat(f string) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.timestampFormat = f
	}
}

// WithUseUTCTime creates an SimpleFormatterOption to specific use the time.UTC layout or not, defaults to false, and means defaults to use time.Local layout.
func WithUseUTCTime(use bool) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.useUTCTime = use
	}
}

// WithDisableColor creates an SimpleFormatterOption to disable the colored format, defaults to false, and means defaults to enable colored format.
func WithDisableColor(disable bool) SimpleFormatterOption {
	return func(c *simpleFormatterOptions) {
		c.disableColor = disable
	}
}

// WithCallerFormatter creates an SimpleFormatterOption to specific the caller's runtime.Frame formatter, defaults to use filename without path and function's shortname.
func WithCallerFormatter(formatter func(*runtime.Frame) (filename string, funcname string)) SimpleFormatterOption {
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

// WithMessageFormatter creates an SimpleFormatterOption to specific the logger formatter.
//
// The default format logs like:
// 	WARN [2021-08-29T05:56:25+08:00] test
// 	INFO [2021-08-29T05:56:25+08:00] filename.go:123 funcname() > test
func WithMessageFormatter(formatter func(level, time, caller, message string) string) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.messageFormatter = formatter
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
// 		WithUseUTCTime(true),
// 		WithDisableColor(false),
// 		WithCallerFormatter(func(*runtime.Frame) (string, string) { return "", "" }),
// 		WithLevelFormatter(func(l logrus.Level) string { return strings.ToUpper(l.String())[:1] }),
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
func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s.initOnce(entry)

	// 1. time and message
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
		var filename, funcName string
		if f := s.option.callerFormatter; f != nil {
			filename, funcName = f(entry.Caller)
		} else {
			_, filename = filepath.Split(entry.Caller.File)
			filename = fmt.Sprintf("%s:%d", filename, entry.Caller.Line)
			_, funcName = filepath.Split(entry.Caller.Function)
			funcName = fmt.Sprintf("%s()", funcName)
		}
		parts := make([]string, 0, 2)
		if filename != "" {
			parts = append(parts, filename)
		}
		if funcName != "" {
			parts = append(parts, funcName)
		}
		caller = strings.Join(parts, " ") // "filename:line funcname()"
	}

	// 3. format
	level := ""
	if f := s.option.levelFormatter; f != nil {
		level = f(entry.Level)
	}
	if level == "" {
		level = strings.ToUpper(entry.Level.String()[0:4])
	}
	if !s.option.disableColor {
		level = s.levelColor(entry.Level).Sprintf(level)
		if caller != "" {
			caller = xcolor.Faint.Sprint(caller)
		}
	}
	formatted := ""
	if f := s.option.messageFormatter; f != nil {
		formatted = f(level, now, caller, message)
	}
	if formatted == "" {
		if caller == "" {
			formatted = fmt.Sprintf("%s [%s] %s\n", level, now, message)
		} else {
			formatted = fmt.Sprintf("%s [%s] %s > %s\n", level, now, caller, message)
		}
	}

	// *. write to buffer
	buf := &bytes.Buffer{}
	if entry.Buffer != nil {
		buf = entry.Buffer
	}
	_, _ = buf.WriteString(formatted)
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

// RFC3339JsonFormatter returns a logrus.JSONFormatter with time.RFC3339Nano timestamp format and custom logrus.FieldMap.
func RFC3339JsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "entries",
		PrettyPrint:       false,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
			logrus.FieldKeyFunc:  "function",
			logrus.FieldKeyFile:  "file",
		},
	}
}

// RFC3339ColoredTextFormatter returns a logrus.TextFormatter with time.RFC3339 timestamp format and force color and quote.
func RFC3339ColoredTextFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		ForceColors:      true,
		ForceQuote:       true,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  time.RFC3339,
		DisableSorting:   false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			_, funcname := filepath.Split(frame.Function)
			_, filename := filepath.Split(frame.File)
			funcname = fmt.Sprintf("%s()", funcname)
			filename = fmt.Sprintf(" %s:%d", filename, frame.Line)
			return funcname, filename
		},
	}
}
