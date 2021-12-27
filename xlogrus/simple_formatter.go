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

// SimpleFormatter represents a simple formatter for logrus.Logger, it only formats level, time, caller and message information with color
// or without color, also note that the logrus.Fields data will not be printed as logrus.TextFormatter does.
type SimpleFormatter struct {
	// TimestampFormat represents the time format, uses time.RFC3339 as default.
	TimestampFormat string

	// RuntimeCaller represents the caller prettifier, uses function and filename directly as default.
	RuntimeCaller func(*runtime.Frame) (function string, file string)

	// DisableColor represents the switcher for color, uses false (use color) as default.
	DisableColor bool

	// terminalInitOnce is the init function. See initOnce.
	terminalInitOnce sync.Once
}

// initOnce initializes the terminal for color supported, this method will be called only once.
func (s *SimpleFormatter) initOnce(entry *logrus.Entry) {
	s.terminalInitOnce.Do(func() {
		if entry.Logger != nil && !s.DisableColor {
			xcolor.InitTerminal(entry.Logger.Out)
		}
	})
}

// Format formats a single log entry, this method implements logrus.Formatter interface.
//
// Logs like:
// 	WARN [2021-08-29T05:56:25+08:00] test
// 	INFO [2021-08-29T05:56:25+08:00] a.go:1 fn() > test
func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s.initOnce(entry)

	// 1. time
	timeFormat := time.RFC3339 // default format
	if s.TimestampFormat != "" {
		timeFormat = s.TimestampFormat
	}
	now := entry.Time.Format(timeFormat)

	// 2. caller
	caller := ""
	if entry.HasCaller() {
		var funcName, filename string
		if s.RuntimeCaller != nil {
			funcName, filename = s.RuntimeCaller(entry.Caller)
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
		}
		if len(parts) > 0 {
			caller = " " + strings.Join(parts, " ")
		}
	}

	// 3. message
	level := strings.ToUpper(entry.Level.String()[0:4])
	message := strings.TrimSuffix(entry.Message, "\n")

	// *. write to buffer
	buf := &bytes.Buffer{}
	if entry.Buffer != nil {
		buf = entry.Buffer
	}
	levelString := level
	if !s.DisableColor {
		color := s.levelColor(entry.Level)
		levelString = color.Sprintf(level)
	}
	_, _ = fmt.Fprintf(buf, "%s [%s]%s %s\n", levelString, now, caller, message)
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
