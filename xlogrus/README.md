# xlogrus

## Functions

### Logrus Formatter

+ `type CustomFormatter struct {}`
+ `(f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error)`

### Logrus Rotate Hook

+ `type RotateFileConfig struct {}`
+ `type RotateFileHook struct {}`
+ `NewRotateFileHook(config *RotateFileConfig) logrus.Hook`
+ `type RotateLogConfig struct {}`
+ `type RotateLogHook struct {}`
+ `NewRotateLogHook(config *RotateLogConfig) logrus.Hook`
