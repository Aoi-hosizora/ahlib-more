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

type CustomFormatter struct {
	TimestampFormat string
	RuntimeCaller   func(*runtime.Frame) (function string, file string)
	ForceColor      bool

	terminalInitOnce sync.Once
}

func (f *CustomFormatter) hasColor() bool {
	return runtime.GOOS != "windows" || f.ForceColor
}

func (f *CustomFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil && f.hasColor() {
		xcolor.InitTerminal(entry.Logger.Out)
	}
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.terminalInitOnce.Do(func() { f.init(entry) })

	b := &bytes.Buffer{}
	if entry.Buffer != nil {
		b = entry.Buffer
	}

	caller := ""
	if entry.HasCaller() {
		var funcVal, fileVal string
		if f.RuntimeCaller != nil {
			funcVal, fileVal = f.RuntimeCaller(entry.Caller)
		} else {
			funcVal = fmt.Sprintf("%s()", entry.Caller.Function)
			fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		}
		sp := strings.Builder{}
		if fileVal != "" {
			sp.WriteString(" ")
			sp.WriteString(fileVal)
		}
		if funcVal != "" {
			sp.WriteString(" ")
			sp.WriteString(funcVal)
		}
		caller = sp.String()
	}

	timeFmt := time.RFC3339
	if f.TimestampFormat != "" {
		timeFmt = f.TimestampFormat
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]
	message := strings.TrimSuffix(entry.Message, "\n")
	now := entry.Time.Format(timeFmt)

	if f.hasColor() {
		var levelColor int
		switch entry.Level {
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = 37 // gray
		case logrus.WarnLevel:
			levelColor = 33 // yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = 31 // red
		default: // info
			levelColor = 36 // blue
		}
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m [%s]%s %s", levelColor, levelText, now, caller, message)
	} else {
		_, _ = fmt.Fprintf(b, "%s [%s]%s %s", levelText, now, caller, message)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
