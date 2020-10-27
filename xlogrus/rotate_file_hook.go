package xlogrus

import (
	"github.com/ah-forklib/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
)

// RotateFileHook's config
type RotateFileConfig struct {
	MaxSize    int  // default to 100 megabytes
	MaxAge     int  // default not to remove old log
	MaxBackups int  // default to retain all old log files
	LocalTime  bool // default to use UTC time
	Compress   bool // default not to perform compression

	Filename  string           // log filename (with filepath, filename and extension)
	Level     logrus.Level     // log level
	Formatter logrus.Formatter // text formatter
}

// Write log into files (split logs to files manually)
type RotateFileHook struct {
	config    *RotateFileConfig
	logWriter io.Writer
}

func NewRotateFileHook(config *RotateFileConfig) logrus.Hook {
	return &RotateFileHook{
		config: config,
		logWriter: &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			LocalTime:  config.LocalTime,
			Compress:   config.Compress,
		},
	}
}

func (r *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.config.Level+1]
}

func (r *RotateFileHook) Fire(entry *logrus.Entry) error {
	b, err := r.config.Formatter.Format(entry)
	if err != nil {
		return err // unreachable
	}

	_, _ = r.logWriter.Write(b) // lock
	return nil
}
