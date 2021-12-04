package xlogrus

import (
	"github.com/ah-forklib/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

// RotateFileConfig represents RotateFileHook's config.
type RotateFileConfig struct {
	// Filename represents the log filename, with filepath, filename and extension, required.
	Filename string

	// Level represents the lowest log level, defaults to logrus.PanicLevel.
	Level logrus.Level

	// Formatter represents the logger formatter, defaults to logrus.JSONFormatter.
	Formatter logrus.Formatter

	// MaxAge represents the max day counts of the file, defaults not to remove old logs.
	MaxAge int

	// MaxSize represents the max size in MB of the file, defaults to 100MB.
	MaxSize int

	// LocalTime represents the switcher for local or UTC time, defaults to use UTC time.
	LocalTime bool

	// Compress represents the switcher for compression, defaults not to perform compression.
	Compress bool
}

// RotateFileHook represents a logrus hook for writing logs into a single file.
//
// Example:
// 	hook := NewRotateFileHook(&RotateFileConfig{
// 		Filename:  "console.log",
// 		Level:     logrus.WarnLevel,
// 		Formatter: &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
// 		MaxAge:    30,
// 		MaxSize:   100,
// 		LocalTime: false,
// 		Compress:  true,
// 	})
// 	logger.AddHook(hook)
type RotateFileHook struct {
	// config is the rotate config.
	config *RotateFileConfig

	// writer is the io.Writer for log file rotation.
	writer io.Writer
}

const (
	panicNilConfig     = "xlogrus: nil config"
	panicEmptyFilename = "xlogrus: empty filename for rotation"
	panicInvalidLevel  = "xlogrus: invalid level"
)

// NewRotateFileHook creates a RotateFileHook as logrus.Hook with RotateFileConfig.
func NewRotateFileHook(config *RotateFileConfig) logrus.Hook {
	if config == nil {
		panic(panicNilConfig)
	}
	if config.Filename == "" {
		panic(panicEmptyFilename)
	}
	if config.Level < logrus.PanicLevel || config.Level > logrus.TraceLevel {
		panic(panicInvalidLevel)
	}
	if config.Formatter == nil {
		config.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339}
	}

	writer := &lumberjack.Logger{
		Filename:  config.Filename,
		MaxSize:   config.MaxSize,
		MaxAge:    config.MaxAge,
		LocalTime: config.LocalTime,
		Compress:  config.Compress,
	}

	return &RotateFileHook{config: config, writer: writer}
}

func (r *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.config.Level+1]
}

// Fire writes logrus.Entry data to io.Writer, this implements logrus.Hook.
func (r *RotateFileHook) Fire(entry *logrus.Entry) error {
	b, _ := r.config.Formatter.Format(entry) // ignore error
	_, _ = r.writer.Write(b)
	return nil
}
