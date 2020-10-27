package xlogrus

import (
	"github.com/ah-forklib/rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"path"
	"time"
)

// RotateLogHook's config
type RotateLogConfig struct {
	MaxAge       time.Duration // default to one week
	RotationTime time.Duration // default to one day
	LocalTime    bool          // default to false (UTC)
	ForceNewFile bool          // force to create a new file to record log

	Filepath  string           // log filepath
	Filename  string           // log filename, without extension
	Level     logrus.Level     // log level
	Formatter logrus.Formatter // text formatter
}

// Write log into files (split logs to files automatically)
type RotateLogHook struct {
	config    *RotateLogConfig
	logWriter io.Writer
}

func NewRotateLogHook(config *RotateLogConfig) logrus.Hook {
	fileName := path.Join(config.Filepath, config.Filename)

	options := []rotatelogs.Option{
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(config.MaxAge),
		rotatelogs.WithRotationTime(config.RotationTime),
	}
	options = append(options, rotatelogs.WithClock(rotatelogs.UTC))
	if config.LocalTime {
		options = append(options, rotatelogs.WithClock(rotatelogs.Local))
	}
	if config.ForceNewFile {
		options = append(options, rotatelogs.ForceNewFile())
	}

	writer, err := rotatelogs.New(fileName+".%Y%m%d.log", options...)
	if err != nil {
		panic(err)
	}

	return &RotateLogHook{
		config:    config,
		logWriter: writer,
	}
}

func (r *RotateLogHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.config.Level+1]
}

func (r *RotateLogHook) Fire(entry *logrus.Entry) error {
	b, err := r.config.Formatter.Format(entry)
	if err != nil {
		return err // unreachable
	}

	_, _ = r.logWriter.Write(b) // lock
	return nil
}
