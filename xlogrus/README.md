# xlogrus

## Dependencies

+ xrotation
+ github.com/Aoi-hosizora/ahlib
+ github.com/sirupsen/logrus (NOTE: do not update to v1.9.0)

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

+ `func WithTimestampFormat(f string) SimpleFormatterOption`
+ `func WithTimeLocation(loc *time.Location) SimpleFormatterOption`
+ `func WithDisableColor(disable bool) SimpleFormatterOption`
+ `func WithCallerFormatter(formatter func(*runtime.Frame) (filename string, funcname string)) SimpleFormatterOption`
+ `func WithLevelFormatter(formatter func(logrus.Level) string) SimpleFormatterOption`
+ `func WithMessageFormatter(formatter func(level, time, caller, message string) string) SimpleFormatterOption`
+ `func NewSimpleFormatter(options ...SimpleFormatterOption) *SimpleFormatter`
+ `func RFC3339JsonFormatter() *logrus.JSONFormatter`
+ `func RFC3339ColoredTextFormatter() *logrus.TextFormatter`
+ `func WithRotateFormatter(formatter logrus.Formatter) RotationHookOption`
+ `func WithRotateLevel(level logrus.Level) RotationHookOption`
+ `func NewRotationHook(rotation *xrotation.RotationLogger, options ...RotationHookOption) *RotationHook`

### Methods

+ `func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error)`
+ `func (r *RotationHook) Levels() []logrus.Level`
+ `func (r *RotationHook) Fire(entry *logrus.Entry) error`
