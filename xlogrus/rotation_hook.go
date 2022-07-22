package xlogrus

import (
	"github.com/Aoi-hosizora/ahlib-more/xrotation"
	"github.com/sirupsen/logrus"
	"time"
)

// rotationHookOptions is a type of RotationHook's option, each field can be set by RotationHookOption function type.
type rotationHookOptions struct {
	formatter logrus.Formatter
	level     logrus.Level
}

// RotationHookOption represents an option type for RotationHook's option, can be created by WithXXX functions.
type RotationHookOption func(*rotationHookOptions)

// WithRotateFormatter creates an RotationHookOption to specify logrus.Formatter for rotation, defaults to logrus.JSONFormatter with time.RFC3339.
func WithRotateFormatter(formatter logrus.Formatter) RotationHookOption {
	return func(o *rotationHookOptions) {
		o.formatter = formatter
	}
}

// WithRotateLevel creates an RotationHookOption to specify the lowest logrus.Level for rotation, defaults to logrus.WarnLevel.
func WithRotateLevel(level logrus.Level) RotationHookOption {
	return func(o *rotationHookOptions) {
		if level < logrus.PanicLevel || level > logrus.TraceLevel {
			level = logrus.PanicLevel // -> 0
		}
		o.level = level
	}
}

// RotationHook represents a logrus.Hook for xrotation.RotationLogger, which will gets automatically rotated when new file created.
type RotationHook struct {
	option   *rotationHookOptions
	rotation *xrotation.RotationLogger
}

var _ logrus.Hook = (*RotationHook)(nil)

const (
	panicNilRotationLogger = "xlogrus: nil xrotation.RotationLogger"
)

// NewRotationHook creates a RotationHook with given xrotation.RotationLogger and RotationHookOption-s.
//
// Example:
// 	l := logrus.New()
// 	l.SetLevel(logrus.TraceLevel)
// 	rotation, _ := xrotation.New(xrotation.WithFilenamePattern("logger.%Y%m%d.log"), ...)
// 	hook, _ := NewRotationHook(
// 		rotation,
// 		WithRotateFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}),
// 		WithRotateLevel(logrus.WarnLevel),
// 	)
// 	l.AddHook(hook)
func NewRotationHook(rotation *xrotation.RotationLogger, options ...RotationHookOption) *RotationHook {
	if rotation == nil {
		panic(panicNilRotationLogger)
	}
	opt := &rotationHookOptions{level: logrus.Level(4294967295)} // -> as default level
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.formatter == nil {
		opt.formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339}
	}
	if opt.level == 4294967295 {
		opt.level = logrus.WarnLevel
	}
	return &RotationHook{rotation: rotation, option: opt}
}

// Levels implements the logrus.Hook interface.
func (r *RotationHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.option.level+1]
}

// Fire writes logrus.Entry data to given xrotation.RotationLogger, and implements the logrus.Hook interface.
func (r *RotationHook) Fire(entry *logrus.Entry) (err error) {
	bs, err := r.option.formatter.Format(entry)
	if err == nil {
		_, _ = r.rotation.Write(bs)
	}
	return nil // ignore error
}
