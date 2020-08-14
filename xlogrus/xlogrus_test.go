package xlogrus

import (
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestCustomFormatter(t *testing.T) {
	l := logrus.New()
	l.SetFormatter(&CustomFormatter{ForceColor: true})
	l.Error("test")
	l.Warn("test")
}

func TestRotateFileHook(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.AddHook(NewRotateFileHook(&RotateFileConfig{
		Filename:  "./logs/file.log",
		Level:     logrus.TraceLevel,
		Formatter: &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}))
	l.SetFormatter(&CustomFormatter{ForceColor: true})

	for i := 0; i < 5; i++ {
		l.Errorf("test at %s", time.Now().Format(time.RFC3339))
		time.Sleep(time.Second * 2)
	}
}

func TestRotateLogHook(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	l.AddHook(NewRotateLogHook(&RotateLogConfig{
		MaxAge:       15 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
		LocalTime:    true,
		Filepath:     "./logs/",
		Filename:     "console",
		Level:        logrus.TraceLevel,
		Formatter:    &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}))
	l.SetFormatter(&CustomFormatter{ForceColor: true})

	for i := 0; i < 5; i++ {
		l.Errorf("test at %s", time.Now().Format(time.RFC3339))
		time.Sleep(time.Second * 2)
	}
}
