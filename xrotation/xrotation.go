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
	"sort"
	"strings"
	"sync"
	"time"
)

// loggerOptions is a type of RotationLogger's option, each field can be set by Option function type.
type loggerOptions struct {
	filenamePattern string
	symlinkFilename string
	nowClock        Clock
	forceNewFile    bool

	rotationTime     time.Duration
	rotationSize     int64
	rotationMaxAge   time.Duration
	rotationMaxCount int32
}

type (
	// Clock represents an interface used by RotationLogger to determine the current time.
	Clock interface{ Now() time.Time }

	// clockFn is an internal type that implements Clock interface.
	clockFn func() time.Time
)

// Now implements the Clock interface.
func (c clockFn) Now() time.Time {
	return c()
}

var (
	_ Clock = (*clockFn)(nil)

	// UTC is an object satisfying the Clock interface, which returns the current time in UTC.
	UTC Clock = clockFn(func() time.Time { return time.Now().UTC() })

	// Local is an object satisfying the Clock interface, which returns the current time in the local timezone.
	Local Clock = clockFn(time.Now)
)

// Option represents an option type for loggerOptions, can be created by WithXXX functions, is used to set up RotationLogger's options.
type Option func(*loggerOptions)

// WithFilenamePattern creates an Option to specific filename pattern for RotationLogger, it is a required option.
func WithFilenamePattern(p string) Option {
	return func(o *loggerOptions) {
		o.filenamePattern = p
	}
}

// WithSymlinkFilename creates an Option to specific symlink filename for RotationLogger, defaults to empty, and means not to create symlink.
func WithSymlinkFilename(f string) Option {
	return func(o *loggerOptions) {
		o.symlinkFilename = f
	}
}

// WithClock creates an Option to specific a Clock for RotationLogger, defaults to Local.
func WithClock(c Clock) Option {
	return func(o *loggerOptions) {
		if c != nil {
			o.nowClock = c
		}
	}
}

// WithForceNewFile creates an Option to let RotationLogger delete all logger files when New is called, defaults to false.
func WithForceNewFile(b bool) Option {
	return func(o *loggerOptions) {
		o.forceNewFile = b
	}
}

// WithRotationTime creates an Option to specific a rotation time for RotationLogger, defaults to 24 hours.
func WithRotationTime(t time.Duration) Option {
	return func(o *loggerOptions) {
		if t > 0 {
			o.rotationTime = t
		}
	}
}

// WithRotationSize creates an Option to specific a rotation size for RotationLogger, defaults to no limit.
func WithRotationSize(size int64) Option {
	return func(o *loggerOptions) {
		if size > 0 {
			o.rotationSize = size
		}
	}
}

// WithRotationMaxAge creates an Option to specific rotation loggers' max age for RotationLogger, defaults to 7 days if maxCount is not set.
// Notes that maxAge and maxCount cannot be set at the same time.
func WithRotationMaxAge(age time.Duration) Option {
	return func(o *loggerOptions) {
		if age > 0 {
			o.rotationMaxAge = age
		}
	}
}

// WithRotationMaxCount creates an Option to specific rotation loggers' max count for RotationLogger, defaults to no limits, and it cannot less than 1.
// Notes that maxAge and maxCount cannot be set at the same time.
func WithRotationMaxCount(count int32) Option {
	return func(o *loggerOptions) {
		if count == 1 {
			count = 0
		}
		if count > 0 {
			o.rotationMaxCount = count
		}
	}
}

// RotationLogger represents a rotation logger, which will gets automatically rotated as you write to it.
type RotationLogger struct {
	option      *loggerOptions
	namePattern *strftime.Strftime
	globPattern string

	mu             sync.RWMutex
	currFile       *os.File
	currBasename   string
	currGeneration uint32
	currFilename   string
}

// New creates a RotationLogger with given Option-s, notes that if forceNewFile option is set to true, all related logger files will be deleted first.
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

	// parse filename pattern
	namePattern, err := xtime.NewStrftime(opt.filenamePattern)
	if err != nil {
		return nil, fmt.Errorf("filename pattern `%s` is invalid: %w", opt.filenamePattern, err)
	}
	globPattern := xtime.ToGlobPattern(opt.filenamePattern)

	// remove files first if force new file
	if opt.forceNewFile {
		matches, err := filepath.Glob(globPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to match glob pattern `%s`: %w", globPattern, err)
		}
		moreMatches, _ := filepath.Glob(globPattern + "_*")
		matches = append(matches, moreMatches...)
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

// Write implements the io.Writer interface, will get rotated writer and do rotate first.
func (r *RotationLogger) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	writer, err := r.getRotatedWriter(true)
	if err != nil {
		return 0, err
	}
	return writer.Write(p)
}

// Rotate rotates the logger first manually, returns error when new file is unavailable to get, or rotate failed.
func (r *RotationLogger) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.getRotatedWriter(false)
	return err
}

// CurrentFilename returns the current file name that the RotationLogger is writing to.
func (r *RotationLogger) CurrentFilename() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.currFilename
}

// Close implements the io.Closer interface.
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

var (
	// _t_testHookMkdir is only used for testing RotationLogger.getRotatedWriter.
	_t_testHookMkdir func()

	// _t_testHookSymlinkMkdir is only used for testing createSymlink.
	_t_testHookSymlinkMkdir func()

	// _t_testHookSymlinkMkdir2 is only used for testing createSymlink.
	_t_testHookSymlinkMkdir2 func() string

	// _t_testHookSymlinkMkdir3 is only used for testing createSymlink.
	_t_testHookSymlinkMkdir3 func()
)

// getRotatedWriter does: check if need to create new file, create a name-non-conflict file, generate a symlink and do rotation.
func (r *RotationLogger) getRotatedWriter(ignoreRotateError bool) (io.Writer, error) {
	// check if need to create new file
	needCreate := false
	generation := r.currGeneration
	basename := r.namePattern.FormatString(xtime.TruncateTime(r.option.nowClock.Now(), r.option.rotationTime))
	if r.currBasename == "" || basename != r.currBasename { // initial or new basename
		needCreate = true
		generation = 0
	} else {
		fi, err := os.Stat(r.currFilename)
		if err == nil && r.option.rotationSize > 0 && fi.Size() >= r.option.rotationSize { // exceed rotation size
			needCreate = true
			generation++
		}
	}
	if !needCreate {
		return r.currFile, nil
	}

	// generate a non-conflict filename
	filename := basename
	var tempName string
	for ; ; generation++ {
		if generation == 0 {
			tempName = filename
		} else {
			tempName = fmt.Sprintf("%s_%d", filename, generation) // xxx, xxx_1, xxx_2, ...
		}
		if _, err := os.Stat(tempName); os.IsNotExist(err) {
			filename = tempName
			break
		}
	}

	// create a new file
	dirname := filepath.Dir(filename)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if _t_testHookMkdir != nil { // only used when testing
			_t_testHookMkdir()
		}
		err := os.MkdirAll(dirname, 0755) // drwxr-xr-x
		if err != nil {
			return nil, fmt.Errorf("failed to create directory `%s`: %w", dirname, err)
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // -rwxr--r--
	if err != nil {
		return nil, fmt.Errorf("failed to create file `%s`: %w", filename, err)
	}

	// generate a symlink and do rotation
	if r.option.symlinkFilename != "" {
		err := createSymlink(filename, r.option.symlinkFilename)
		if err != nil {
			log.Printf("Warning: failed to create symlink for `%s`: %v", filename, err)
			// ignore error, especially on Windows: "A required privilege is not held by the client"
		}
	}
	err = r.doRotation()
	if err != nil {
		err = fmt.Errorf("failed to rotate on pattern `%s`: %w", r.globPattern, err)
		if ignoreRotateError {
			log.Printf("Warning: %v", err) // ignore error
		} else {
			_ = file.Close()
			return nil, err
		}
	}

	if r.currFile != nil {
		_ = r.currFile.Close()
	}
	r.currFile = file
	r.currGeneration = generation
	r.currBasename = basename
	r.currFilename = filename
	return file, nil
}

// createSymlink creates a symlink file `linkname` and its destination is `filename`.
func createSymlink(filename, linkname string) error {
	// create target link file directory
	linkDirname := filepath.Dir(linkname)
	if _, err := os.Stat(linkDirname); os.IsNotExist(err) {
		if _t_testHookSymlinkMkdir != nil {
			_t_testHookSymlinkMkdir()
		}
		err := os.MkdirAll(linkDirname, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory `%s`: %w", linkDirname, err)
		}
	}

	// check the relative path of destination
	destination := filename
	if filepath.Dir(destination) != filepath.Dir(linkname) {
		destinationPath, _ := filepath.Abs(destination)
		linkDirnamePath, _ := filepath.Abs(linkDirname)
		var err error
		if _t_testHookSymlinkMkdir2 != nil {
			linkDirnamePath = _t_testHookSymlinkMkdir2()
		}
		destination, err = filepath.Rel(linkDirnamePath, destinationPath)
		if err != nil {
			return fmt.Errorf("failed to evaluate the relative path from `%s` to `%s`: %w", destinationPath, linkDirnamePath, err)
		}
	}

	// make symlink and rename to the link file
	tempLinkname := filename + "_symlink"
	if _, err := os.Stat(tempLinkname); err == nil {
		_ = os.Remove(tempLinkname)
	}
	if _t_testHookSymlinkMkdir3 != nil {
		_t_testHookSymlinkMkdir3()
	}
	err := os.Symlink(destination, tempLinkname)
	if err != nil {
		return fmt.Errorf("failed to create new symlink `%s`: %w", tempLinkname, err)
	}
	err = os.Rename(tempLinkname, linkname)
	if err != nil {
		return fmt.Errorf("failed to rename symlink `%s` to `%s`: %w", tempLinkname, linkname, err)
	}
	return nil
}

// doRotation does the real rotation work, this will rotate for loggers' max age or for loggers' max count, and remove all unlinked files.
func (r *RotationLogger) doRotation() error {
	// get matches by glob pattern
	matches, err := filepath.Glob(r.globPattern)
	if err != nil {
		return fmt.Errorf("failed to match glob pattern `%s`: %w", r.globPattern, err)
	}
	unlinkFiles := make([]string, 0)

	// I) rotate for max age
	if r.option.rotationMaxAge > 0 {
		cutoffDuration := r.option.nowClock.Now().Add(-1 * r.option.rotationMaxAge)
		for _, match := range matches {
			fi, err := os.Lstat(match)
			if err != nil || (fi.Mode()&os.ModeSymlink) == os.ModeSymlink {
				continue
			}
			if r.option.rotationMaxAge > 0 && fi.ModTime().After(cutoffDuration) {
				continue
			}
			unlinkFiles = append(unlinkFiles, match)
		}
	}

	// II) rotate for max count
	if count := int(r.option.rotationMaxCount); count > 0 {
		type nameTimeTuple struct {
			name string
			mod  time.Time
		}
		matchesFileInfos := make([]nameTimeTuple, 0)
		for _, match := range matches {
			fi, err := os.Lstat(match)
			if err != nil || (fi.Mode()&os.ModeSymlink) == os.ModeSymlink {
				continue
			}
			matchesFileInfos = append(matchesFileInfos, nameTimeTuple{match, fi.ModTime()})
		}
		if len(matchesFileInfos) > count {
			sort.Slice(matchesFileInfos, func(i, j int) bool {
				return matchesFileInfos[i].mod.Before(matchesFileInfos[j].mod)
			})
			for _, fi := range matchesFileInfos[:len(matchesFileInfos)-count] {
				unlinkFiles = append(unlinkFiles, fi.name)
			}
		}
	}

	// add moreUnlinkFiles according unlinkFiles
	if len(unlinkFiles) == 0 {
		return nil
	}
	moreMatches, _ := filepath.Glob(r.globPattern + "_*")
	moreUnlinkFiles := make([]string, 0)
	for _, path := range unlinkFiles {
		for _, match := range moreMatches {
			if strings.HasPrefix(match, path) {
				moreUnlinkFiles = append(moreUnlinkFiles, match)
			}
		}
	}

	// remove unlinked files
	for _, path := range append(unlinkFiles, moreUnlinkFiles...) {
		err := os.Remove(path)
		if err != nil {
			log.Printf("Warning: failed to remove file `%s`: %v", path, err) // ignore error
		}
	}
	return nil
}
