package xlogrus

import (
	"bytes"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xcolor"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"sync"
	"time"
)

// SimpleFormatter represents a simple formatter for logrus.Logger, it only formats level, time, caller and message information
// in colored, and the logrus.Fields data is not be logged like logrus.TextFormatter does.
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

// Format renders a single log entry, this method implements logrus.Formatter.
func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	s.initOnce(entry)

	buf := &bytes.Buffer{}
	if entry.Buffer != nil {
		buf = entry.Buffer
	}

	// 1. time
	timeFormat := time.RFC3339 // default format
	if s.TimestampFormat != "" {
		timeFormat = s.TimestampFormat
	}
	now := entry.Time.Format(timeFormat)

	// 2. caller
	caller := ""
	if entry.HasCaller() {
		var funcVal, fileVal string
		if s.RuntimeCaller != nil {
			funcVal, fileVal = s.RuntimeCaller(entry.Caller)
		} else {
			funcVal = fmt.Sprintf("%s()", entry.Caller.Function)
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}
		sp := strings.Builder{}
		if fileVal != "" {
			sp.WriteByte(' ')
			sp.WriteString(fileVal)
		}
		if funcVal != "" {
			sp.WriteByte(' ')
			sp.WriteString(funcVal)
		}
		caller = sp.String() // runtime caller string
	}

	// 3. message
	level := strings.ToUpper(entry.Level.String()[0:4])
	message := strings.TrimSuffix(entry.Message, "\n")

	// write
	if s.DisableColor {
		_, _ = fmt.Fprintf(buf, "%s [%s]%s %s", level, now, caller, message)
	} else {
		levelColor := int(s.levelColor(entry.Level))
		_, _ = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m [%s]%s %s", levelColor, level, now, caller, message)
	}
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

// levelColor returns the color code from logrus.Level.
func (s *SimpleFormatter) levelColor(level logrus.Level) uint8 {
	switch level {
	case logrus.InfoLevel:
		return xcolor.Blue.Code()
	case logrus.WarnLevel:
		return xcolor.Yellow.Code()
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return xcolor.Red.Code()
	default: // debug, trace
		return xcolor.White.Code()
	}
}
