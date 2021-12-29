# xlogrus

## Dependencies

+ github.com/Aoi-hosizora/ahlib
+ github.com/sirupsen/logrus

## Documents

### Types

+ `type SimpleFormatterOption func`
+ `type SimpleFormatter struct`
+ `type RotationHookOption func`
+ `type RotationHook struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func WithTimestampFormat(format string) SimpleFormatterOption`
+ `func WithCallerFormatter(formatter func(*runtime.Frame) (function string, file string)) SimpleFormatterOption`
+ `func WithLevelFormatter(formatter func(logrus.Level) string) SimpleFormatterOption`
+ `func WithDisableColor(disable bool) SimpleFormatterOption`
+ `func WithUseUTCTime(use bool) SimpleFormatterOption`
+ `func NewSimpleFormatter(options ...SimpleFormatterOption) *SimpleFormatter`
+ `func WithRotationLogger(rl *xrotation.RotationLogger) RotationHookOption`
+ `func WithFormatter(formatter logrus.Formatter) RotationHookOption`
+ `func WithLevel(level logrus.Level) RotationHookOption`
+ `func NewRotationHook(options ...RotationHookOption) (logrus.Hook, error)`

### Methods

+ `func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error)`
+ `func (r *RotationHook) Levels() []logrus.Level`
+ `func (r *RotationHook) Fire(entry *logrus.Entry) error`
