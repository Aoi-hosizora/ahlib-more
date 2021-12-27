package xrotation

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/ah-forklib/strftime"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type loggerOptions struct {
	filenamePattern string
	symlinkFilename string
	nowClock        Clock
	useSingleFile   bool
	forceNewFile    bool

	rotationTime     time.Duration
	rotationSize     int64
	rotationMaxAge   time.Duration
	rotationMaxCount int32
}

type (
	Clock   interface{ Now() time.Time }
	clockFn func() time.Time
)

func (c clockFn) Now() time.Time {
	return c()
}

var (
	_ Clock = (*clockFn)(nil)

	UTC   = clockFn(func() time.Time { return time.Now().UTC() })
	Local = clockFn(time.Now)
)

type Option func(*loggerOptions)

func WithFilenamePattern(p string) Option {
	return func(o *loggerOptions) {
		o.filenamePattern = p
	}
}

func WithSymlinkFilename(f string) Option {
	return func(o *loggerOptions) {
		o.symlinkFilename = f
	}
}

func WithClock(c Clock) Option {
	return func(o *loggerOptions) {
		if c != nil {
			o.nowClock = c
		}
	}
}

func WithUseSingleFile(b bool) Option {
	return func(o *loggerOptions) {
		o.useSingleFile = b
	}
}

func WithForceNewFile(b bool) Option {
	return func(o *loggerOptions) {
		o.forceNewFile = b
	}
}

func WithRotationTime(t time.Duration) Option {
	return func(o *loggerOptions) {
		if t > 0 {
			o.rotationTime = t
		}
	}
}

func WithRotationSize(size int64) Option {
	return func(o *loggerOptions) {
		if size > 0 {
			o.rotationSize = size
		}
	}
}

func WithRotationMaxAge(age time.Duration) Option {
	return func(o *loggerOptions) {
		if age > 0 {
			o.rotationMaxAge = age
		}
	}
}

func WithRotationMaxCount(count int32) Option {
	return func(o *loggerOptions) {
		if count > 0 {
			o.rotationMaxCount = count
		}
	}
}

type RotationLogger struct {
	option         *loggerOptions
	namePattern    *strftime.Strftime
	globPattern    string
	mu             sync.RWMutex
	currFile       *os.File
	currBasename   string
	currGeneration uint32
	currFilename   string
}

func New(options ...Option) (*RotationLogger, error) {
	opt := &loggerOptions{nowClock: Local, rotationTime: 24 * time.Hour}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.filenamePattern == "" {
		return nil, errors.New("empty filename pattern is not allowed")
	}
	if opt.rotationMaxAge > 0 && opt.rotationMaxCount > 0 {
		return nil, errors.New("rotation max age and max count can not be set at the same time")
	}
	if opt.rotationMaxAge == 0 && opt.rotationMaxCount == 0 {
		opt.rotationMaxAge = 7 * 24 * time.Hour
	}

	namePattern, err := xtime.NewStrftime(opt.filenamePattern)
	if err != nil {
		return nil, fmt.Errorf("filename pattern `%s` is invalid: %w", opt.filenamePattern, err)
	}
	globPattern := xtime.ToGlobPattern(opt.filenamePattern)

	if opt.forceNewFile {
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to match glob pattern `%s`: %w", globPattern, err)
		}
		if !opt.useSingleFile {
			moreMatches, _ := filepath.Glob(globPattern + "_*")
			matches = append(matches, moreMatches...)
		}
		for _, match := range matches {
			err := os.Remove(match)
			if err != nil {
				return nil, fmt.Errorf("failed to remove matched file `%s`: %w", match, err)
			}
		}
	}

	logger := &RotationLogger{option: opt, namePattern: namePattern, globPattern: globPattern}
	return logger, nil
}

func (r *RotationLogger) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	writer, err := r.getRotatedWriter(false)
	if err != nil {
		return 0, err
	}
	return writer.Write(p)
}

func (r *RotationLogger) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.getRotatedWriter(true)
	return err
}

func (r *RotationLogger) CurrentFilename() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.currFilename
}

func (r *RotationLogger) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.currFile == nil {
		return nil
	}
	_ = r.currFile.Close()
	r.currFile = nil
	return nil
}

// ===================
// core implementation
// ===================

func (r *RotationLogger) getRotatedWriter(rotateManually bool) (io.Writer, error) {
	createNewFile := false
	generation := r.currGeneration
	basename := r.namePattern.FormatString(xtime.TruncateTime(r.option.nowClock.Now(), r.option.rotationTime)) // <<<
	if basename != r.currBasename {
		createNewFile = true
		generation = 0
	} else if rotateManually {
		createNewFile = true
		generation++
	} else {
		fi, err := os.Stat(r.currFilename)
		if err == nil && r.option.rotationSize > 0 && fi.Size() >= r.option.rotationSize {
			createNewFile = true
			generation++
		}
	}

	if !createNewFile {
		return r.currFile, nil
	}
	if r.option.useSingleFile {
		_, err := os.Stat(basename)
		if err == nil {
			_ = r.currFile.Close()
			generation = 0
			err = os.Remove(basename)
			if err != nil {
				return nil, fmt.Errorf("failed to remove file `%s`: %w", basename, err)
			}
		}
	}

	filename := basename
	var tempName string
	for ; ; generation++ {
		if generation == 0 {
			tempName = filename
		} else {
			tempName = fmt.Sprintf("%s_%d", filename, generation)
		}
		if _, err := os.Stat(tempName); os.IsNotExist(err) {
			filename = tempName
			break
		}
	}

	dirname := filepath.Dir(filename)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err := os.MkdirAll(dirname, 0755) // drwxr-xr-x
		if err != nil {
			return nil, fmt.Errorf("failed to create directory `%s`: %w", dirname, err)
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // -rwxr--r--
	if err != nil {
		return nil, fmt.Errorf("failed to create file `%s`: %w", filename, err)
	}
	if !r.option.useSingleFile {
		err = r.doContentRotate(filename)
		if err != nil {
			err = fmt.Errorf("failed to rotate on pattern `%s`: %w", r.globPattern, err)
			if rotateManually {
				_ = file.Close()
				return nil, err
			} else {
				log.Println(err)
				// ignore error
			}
		}
	}

	_ = r.currFile.Close()
	r.currFile = file
	r.currGeneration = generation
	r.currBasename = basename
	r.currFilename = filename
	return file, nil
}

func createSymlink(filename, linkname string) error {
	tempLinkname := filename + "_symlink"
	destination := filename
	linkDirname := filepath.Dir(linkname)
	if filepath.Dir(destination) != filepath.Dir(linkname) {
		destinationPath, err := filepath.Abs(destination)
		if err != nil {
			return fmt.Errorf("failed to get the absolute path of file `%s`: %w", destination, err)
		}
		destination, err = filepath.Rel(linkDirname, destinationPath)
		if err != nil {
			return fmt.Errorf("failed to evaluate the relative path from `%s` to `%s`: %w", destinationPath, linkDirname, err)
		}
	}
	err := os.Symlink(destination, tempLinkname)
	if err != nil {
		return fmt.Errorf("failed to create new symlink `%s`: %w", tempLinkname, err)
	}

	if _, err := os.Stat(linkDirname); os.IsNotExist(err) {
		err := os.MkdirAll(linkDirname, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory `%s`: %w", linkDirname, err)
		}
	}
	err = os.Rename(tempLinkname, linkname)
	if err != nil {
		return fmt.Errorf("failed to rename symlink `%s` to `%s`: %w", tempLinkname, linkname, err)
	}
	return nil
}

func (r *RotationLogger) doContentRotate(filename string) error {
	if r.option.symlinkFilename != "" {
		err := createSymlink(filename, r.option.symlinkFilename)
		log.Println(err)
		// ignore error, especially on Windows: A required privilege is not held by the client
	}

	matches, err := filepath.Glob(r.globPattern) // TODO
	if err != nil {
		return fmt.Errorf("failed to match glob pattern `%s`: %w", r.globPattern, err)
	}

	unlinkFiles := make([]string, 0)
	cutoffDuration := r.option.nowClock.Now().Add(-1 * r.option.rotationMaxAge)
	for _, match := range matches {
		fi, err := os.Lstat(match)
		if err != nil {
			continue
		}
		if (fi.Mode() & os.ModeSymlink) == os.ModeSymlink {
			continue
		}
		if r.option.rotationMaxAge > 0 && fi.ModTime().After(cutoffDuration) {
			continue
		}
		unlinkFiles = append(unlinkFiles, match)
	}

	if count := int(r.option.rotationMaxCount); count > 0 {
		if len(unlinkFiles) > count {
			unlinkFiles = unlinkFiles[:len(unlinkFiles)-count]
		}
	}
	if len(unlinkFiles) == 0 {
		return nil
	}

	for _, path := range unlinkFiles {
		_ = os.Remove(path)
	}
	return nil
}
