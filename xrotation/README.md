# xrotation

## Dependencies

+ github.com/Aoi-hosizora/ahlib

## Documents

### Types

+ `type Clock interface`
+ `type Option func`
+ `type RotationLogger struct`

### Variables

+ `var UTC Clock`
+ `var Local Clock`

### Constants

+ None

### Functions

+ `func WithFilenamePattern(p string) Option`
+ `func WithSymlinkFilename(f string) Option`
+ `func WithClock(c Clock) Option`
+ `func WithForceNewFile(b bool) Option`
+ `func WithRotationTime(t time.Duration) Option`
+ `func WithRotationSize(size int64) Option`
+ `func WithRotationMaxAge(age time.Duration) Option`
+ `func WithRotationMaxCount(count int32) Option`
+ `func New(options ...Option) (*RotationLogger, error)`

### Methods

+ `func (r *RotationLogger) Write(p []byte) (n int, err error)`
+ `func (r *RotationLogger) Rotate() error`
+ `func (r *RotationLogger) CurrentFilename() string`
+ `func (r *RotationLogger) Close() error`
