package xlogrus

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib-more/xrotation"
	"github.com/sirupsen/logrus"
	"time"
)

// rotationHookOptions is a type of RotationHook's option, each field can be set by RotationHookOption function type.
type rotationHookOptions struct {
	rotationLogger *xrotation.RotationLogger
	formatter      logrus.Formatter
	level          logrus.Level
}

// RotationHookOption represents an option type for RotationHook's option, can be created by WithXXX functions.
type RotationHookOption func(*rotationHookOptions)

// WithRotationLogger creates an RotationHookOption to specific xrotation.RotationLogger for rotation, it is a required option.
func WithRotationLogger(rl *xrotation.RotationLogger) RotationHookOption {
	return func(o *rotationHookOptions) {
		o.rotationLogger = rl
	}
}

// WithFormatter creates an RotationHookOption to specific logrus.Formatter for rotation, defaults to logrus.JSONFormatter with time.RFC3339.
func WithFormatter(formatter logrus.Formatter) RotationHookOption {
	return func(o *rotationHookOptions) {
		o.formatter = formatter
	}
}

// WithLevel creates an RotationHookOption to specific lowest logrus.Level for rotation, defaults to logrus.WarnLevel.
func WithLevel(level logrus.Level) RotationHookOption {
	return func(o *rotationHookOptions) {
		if level < logrus.PanicLevel || level > logrus.TraceLevel {
			level = logrus.PanicLevel // -> 0
		}
		o.level = level
	}
}

// RotationHook represents a logrus.Hook for xrotation.RotationLogger, which will gets automatically rotated when new file created.
type RotationHook struct {
	option *rotationHookOptions
}

var _ logrus.Hook = (*RotationHook)(nil)

var (
	errNilRotationLogger = errors.New("xlogrus: nil xrotation.RotationLogger")
)

// NewRotationHook creates a RotationHook with given RotationHookOption-s.
//
// Example:
// 	l := logrus.New()
// 	l.SetLevel(logrus.TraceLevel)
// 	rl, _ := xrotation.New(xrotation.WithFilenamePattern("logger.%Y%m%d.log"), ...)
// 	hook, _ := NewRotationHook(
// 		WithRotationLogger(rl),
// 		WithFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}),
// 		WithLevel(logrus.WarnLevel),
// 	)
// 	l.AddHook(hook)
func NewRotationHook(options ...RotationHookOption) (*RotationHook, error) {
	opt := &rotationHookOptions{level: logrus.Level(4294967295)}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.rotationLogger == nil {
		return nil, errNilRotationLogger
	}
	if opt.formatter == nil {
		opt.formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339}
	}
	if opt.level == 4294967295 {
		opt.level = logrus.WarnLevel
	}
	return &RotationHook{option: opt}, nil
}

// Levels implements the logrus.Hook interface.
func (r *RotationHook) Levels() []logrus.Level {
	return logrus.AllLevels[:r.option.level+1]
}

// Fire writes logrus.Entry data to given xrotation.RotationLogger, and implements the logrus.Hook interface.
func (r *RotationHook) Fire(entry *logrus.Entry) error {
	bs, err := r.option.formatter.Format(entry)
	if err == nil {
		_, _ = r.option.rotationLogger.Write(bs)
	}
	return nil // ignore error
}
