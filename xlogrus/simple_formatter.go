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
	timeLocation     *time.Location
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

// WithTimeLocation creates an SimpleFormatterOption to specific the time.Location for entry's time, defaults to time.Local.
func WithTimeLocation(loc *time.Location) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.timeLocation = loc
	}
}

// WithDisableColor creates an SimpleFormatterOption to disable the colored format, defaults to false, and means defaults to enable colored format.
func WithDisableColor(disable bool) SimpleFormatterOption {
	return func(o *simpleFormatterOptions) {
		o.disableColor = disable
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
// 		WithTimeLocation(time.UTC),
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
	if opt.timeLocation == nil {
		opt.timeLocation = time.Local
	}
	if opt.callerFormatter == nil {
		opt.callerFormatter = func(frame *runtime.Frame) (filename string, funcname string) {
			_, filename = filepath.Split(frame.File)
			_, funcname = filepath.Split(frame.Function)
			filename = fmt.Sprintf("%s:%d", filename, frame.Line)
			funcname = fmt.Sprintf("%s()", funcname)
			return filename, funcname
		}
	}
	if opt.levelFormatter == nil {
		opt.levelFormatter = func(level logrus.Level) string {
			return strings.ToUpper(level.String()[0:4])
		}
	}
	if opt.messageFormatter == nil {
		opt.messageFormatter = func(level, time, caller, message string) string {
			if caller == "" {
				return fmt.Sprintf("%s [%s] %s\n", level, time, message)
			}
			return fmt.Sprintf("%s [%s] %s > %s\n", level, time, caller, message)
		}
	}
	return &SimpleFormatter{option: opt}
}

// Format formats a single logrus.Entry, this method implements logrus.Formatter interface.
func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s.terminalOnce.Do(func() {
		// 0. initialize the terminal io.Writer for color supported
		if entry.Logger != nil && !s.option.disableColor {
			xcolor.InitTerminal(entry.Logger.Out)
		}
	})

	// 1. level and time
	level := s.option.levelFormatter(entry.Level)
	now := entry.Time.In(s.option.timeLocation).Format(s.option.timestampFormat)

	// 2. runtime caller
	caller := ""
	if entry.HasCaller() {
		filename, funcname := s.option.callerFormatter(entry.Caller)
		parts := make([]string, 0, 2)
		if filename != "" {
			parts = append(parts, filename)
		}
		if funcname != "" {
			parts = append(parts, funcname)
		}
		caller = strings.Join(parts, " ") // "filename:line funcname()"
	}

	// 3. color and message
	if !s.option.disableColor {
		level = s.levelColor(entry.Level).Sprintf(level)
		if caller != "" {
			caller = s.callerStyle().Sprint(caller)
		}
	}
	message := strings.TrimSuffix(entry.Message, "\n")
	message = s.option.messageFormatter(level, now, caller, message)

	// *. write to buffer
	buf := entry.Buffer
	if entry.Buffer == nil {
		// nil only when fire hooks
		buf = &bytes.Buffer{}
	}
	_, _ = buf.WriteString(message)
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

// callerStyle returns the xcolor.Style for caller.
func (s *SimpleFormatter) callerStyle() xcolor.Style {
	return xcolor.Faint // you may not see this style in Windows's cmd
}

// RFC3339JsonFormatter returns a logrus.JSONFormatter with time.RFC3339Nano timestamp format, "entries" data key and custom logrus.FieldMap.
func RFC3339JsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339Nano,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		PrettyPrint:       false,
		DataKey:           "entries",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
			logrus.FieldKeyFile:  "file",
			logrus.FieldKeyFunc:  "function",
			// @@time -> @level -> @message -> entries -> file -> function
		},
	}
}

// RFC3339ColoredTextFormatter returns a logrus.TextFormatter with time.RFC3339 timestamp format, full timestamp format, force color and quote.
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
			filename = fmt.Sprintf("%s:%d", filename, frame.Line)
			return funcname, filename
		},
	}
}
