package xlogrus

import (
	"github.com/ah-forklib/rotatelogs"
	"github.com/ah-forklib/strftime"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

// RotateLogConfig represents RotateLogHook's config.
type RotateLogConfig struct {
	// Filename represents the log filename without time part and extension, required.
	Filename string

	// FilenameTimePart represents time part after filename, defaults to ".%Y%m%d.log". See strftime.New.
	FilenameTimePart string

	// LinkFileName represents the symbolic link filename, defaults to "", no link will be written.
	LinkFileName string

	// Level represents the lowest log level, defaults to logrus.PanicLevel.
	Level logrus.Level

	// Formatter represents the logger formatter, defaults to logrus.JSONFormatter.
	Formatter logrus.Formatter

	// MaxAge represents the max duration of the file, defaults to one week.
	MaxAge time.Duration

	// MaxSize represents the max size in MB of the file, defaults to no limit.
	MaxSize int

	// RotationTime represents the rotation duration of the file, defaults to one day.
	RotationTime time.Duration

	// LocalTime represents the switcher for local or UTC time, defaults to use UTC time.
	LocalTime bool

	// ForceNewFile represents the switcher for forcing to save to new file, defaults to false.
	ForceNewFile bool
}

// RotateLogHook represents a logrus hook for writing logs into files splitting by time.
type RotateLogHook struct {
	// config is the rotate config.
	config *RotateLogConfig

	// writer is the io.Writer for log file rotation.
	writer io.Writer
}

const (
	panicInvalidTimePattern = "xlogrus: invalid time pattern for filename"
)

// NewRotateLogHook creates a RotateLogHook as logrus.Hook with RotateLogConfig.
//
// Example:
// 	hook := NewRotateLogHook(&RotateLogConfig{
// 		Filename:         "console",
// 		FilenameTimePart: ".%Y%m%d.log",
// 		LinkFileName:     "console.curr.log",
// 		Level:            logrus.WarnLevel,
// 		Formatter:        &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
// 		MaxAge:           time.Hour * 24 * 30,
// 		MaxSize:          100,
// 		RotationTime:     time.Hour * 24,
// 		LocalTime:        false,
// 		ForceNewFile:     false,
// 	})
// 	logger.AddHook(hook)
func NewRotateLogHook(config *RotateLogConfig) logrus.Hook {
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

	timePartName := ".%Y%m%d.log" // default time part name
	if config.FilenameTimePart != "" {
		timePartName = config.FilenameTimePart
	}
	_, err := strftime.New(timePartName)
	if err != nil {
		panic(panicInvalidTimePattern)
	}

	options := []rotatelogs.Option{
		rotatelogs.WithLinkName(config.LinkFileName),
		rotatelogs.WithMaxAge(config.MaxAge),
		rotatelogs.WithRotationSize(int64(float64(config.MaxSize) * 1024 * 1024)), // MB -> B
		rotatelogs.WithRotationTime(config.RotationTime),
	}
	if config.LocalTime {
		options = append(options, rotatelogs.WithClock(rotatelogs.Local))
	} else {
		options = append(options, rotatelogs.WithClock(rotatelogs.UTC))
	}
	if config.ForceNewFile {
		options = append(options, rotatelogs.ForceNewFile())
	}

	filename := config.Filename + timePartName
	writer, _ := rotatelogs.New(filename, options...) // no error

	return &RotateLogHook{config: config, writer: writer}
}

func (r *RotateLogHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.config.Level+1]
}

// Fire writes logrus.Entry data to io.Writer, this implements logrus.Hook.
func (r *RotateLogHook) Fire(entry *logrus.Entry) error {
	b, _ := r.config.Formatter.Format(entry) // ignore error
	_, _ = r.writer.Write(b)
	return nil
}
