package xlogrus

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"testing"
	"time"
)

func TestCustomFormatter(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)

	// no caller, with color
	l.SetReportCaller(false)
	l.SetFormatter(&CustomFormatter{
		RuntimeCaller:   nil,
		TimestampFormat: "",
	})
	l.Trace("test")
	l.Info("test")

	// default caller, with color
	l.SetReportCaller(true)
	l.SetFormatter(&CustomFormatter{
		DisableColor:    false,
		RuntimeCaller:   nil,
		TimestampFormat: time.RubyDate,
	})
	l.Warn("test")
	l.Error("test")

	// no color, with caller
	l.SetReportCaller(true)
	l.SetFormatter(&CustomFormatter{
		DisableColor:  true,
		RuntimeCaller: func(*runtime.Frame) (string, string) { return "package.funcname", "filename:1" },
	})
	l.Trace("test")
	l.Info("test")
	l.Warn("test")
	l.Error("test")
}

type fakeFormatter struct{}

func (f *fakeFormatter) Format(*logrus.Entry) ([]byte, error) {
	return []byte("fake"), fmt.Errorf("fake error")
}

func TestRotateFileHook(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	cfg := &RotateFileConfig{
		Filename:  "./logs/file.log",
		Level:     logrus.TraceLevel,
		Formatter: &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}
	l.AddHook(NewRotateFileHook(cfg))
	l.SetFormatter(&CustomFormatter{DisableColor: true})

	for i := 0; i < 5; i++ {
		l.Infof("test at %s", time.Now().Format(time.RFC3339Nano))
		time.Sleep(time.Millisecond * 100)
	}

	l = logrus.New()
	l.SetLevel(logrus.TraceLevel)
	cfg.Formatter = &fakeFormatter{}
	l.AddHook(NewRotateFileHook(cfg))
	l.Error("!!!!!!")
}

func TestRotateLogHook(t *testing.T) {
	func() {
		defer func() { recover() }()
		_ = NewRotateLogHook(&RotateLogConfig{Filename: "%"})
	}()

	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	cfg := &RotateLogConfig{
		MaxAge:       15 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
		LocalTime:    true,
		ForceNewFile: true,
		Filepath:     "./logs/",
		Filename:     "console",
		Level:        logrus.TraceLevel,
		Formatter:    &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}
	l.AddHook(NewRotateLogHook(cfg))
	l.SetFormatter(&CustomFormatter{DisableColor: true})

	for i := 0; i < 5; i++ {
		l.Infof("test at %s", time.Now().Format(time.RFC3339Nano))
		time.Sleep(time.Millisecond * 100)
	}

	l = logrus.New()
	l.SetLevel(logrus.TraceLevel)
	cfg.Formatter = &fakeFormatter{}
	l.AddHook(NewRotateLogHook(cfg))
	l.Error("!!!!!!")
}
