# xlogrus

## Dependencies

+ github.com/sirupsen/logrus
+ github.com/ah-forklib/lumberjack (gopkg.in/natefinch/lumberjack.v2)
+ github.com/ah-forklib/rotatelogs (github.com/lestrrat-go/file-rotatelogs)
+ github.com/ah-forklib/strftime (github.com/lestrrat-go/strftime)

## Documents

### Functions

### Types

+ `type SimpleFormatter struct`
+ `type RotateFileConfig struct`
+ `type RotateFileHook struct`
+ `type RotateLogConfig struct`
+ `type RotateLogHook struct`

### Variables

+ None

### Constants

+ None

### Functions

+ `func NewRotateFileHook(config *RotateFileConfig) logrus.Hook`
+ `func NewRotateLogHook(config *RotateLogConfig) logrus.Hook`

### Methods

+ `func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error)`
+ `func (r *RotateFileHook) Fire(entry *logrus.Entry) error`
+ `func (r *RotateLogHook) Fire(entry *logrus.Entry) error`
