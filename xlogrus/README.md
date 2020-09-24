# xlogrus

### Dependencies

+ github.com/sirupsen/logrus
+ gopkg.in/natefinch/lumberjack.v2
+ github.com/lestrrat-go/file-rotatelogs
+ xcolor

### Functions

#### Formatter

+ `type CustomFormatter struct {}`
+ `(f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error)`

#### Rotate Hook

+ `type RotateFileConfig struct {}`
+ `type RotateFileHook struct {}`
+ `NewRotateFileHook(config *RotateFileConfig) logrus.Hook`
+ `type RotateLogConfig struct {}`
+ `type RotateLogHook struct {}`
+ `NewRotateLogHook(config *RotateLogConfig) logrus.Hook`
