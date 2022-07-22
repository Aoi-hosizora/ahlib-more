package xrotation

import (
	"errors"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xerror"
	"github.com/Aoi-hosizora/ahlib/xtime"
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
	symlinkFilename string
	nowClock        xtime.Clock
	forceNewFile    bool

	rotationTime     time.Duration
	rotationSize     int64
	rotationMaxAge   time.Duration
	rotationMaxCount int32
}

// Option represents an option type for RotationLogger's option, can be created by WithXXX functions.
type Option func(*loggerOptions)

// WithSymlinkFilename creates an Option to specify symlink filename for RotationLogger, defaults to empty, and means not to create symlink.
func WithSymlinkFilename(f string) Option {
	return func(o *loggerOptions) {
		o.symlinkFilename = f
	}
}

// WithClock creates an Option to specify a xtime.Clock for RotationLogger, defaults to xtime.Local.
func WithClock(c xtime.Clock) Option {
	return func(o *loggerOptions) {
		o.nowClock = c
	}
}

// WithForceNewFile creates an Option to let RotationLogger create a new file when write initially, defaults to false.
func WithForceNewFile(b bool) Option {
	return func(o *loggerOptions) {
		o.forceNewFile = b
	}
}

// WithRotationTime creates an Option to specify a rotation time for RotationLogger, defaults to 24 hours.
func WithRotationTime(t time.Duration) Option {
	return func(o *loggerOptions) {
		if t < 0 {
			t = 0
		}
		o.rotationTime = t
	}
}

// WithRotationSize creates an Option to specify a rotation size for RotationLogger, defaults to no limit.
func WithRotationSize(size int64) Option {
	return func(o *loggerOptions) {
		if size < 0 {
			size = 0
		}
		o.rotationSize = size
	}
}

// WithRotationMaxAge creates an Option to specify rotation loggers' max age for RotationLogger, defaults to 7 days if maxCount is not set.
// Note that maxAge and maxCount cannot be set at the same time.
func WithRotationMaxAge(age time.Duration) Option {
	return func(o *loggerOptions) {
		if age < 0 {
			age = 0
		}
		o.rotationMaxAge = age
	}
}

// WithRotationMaxCount creates an Option to specify rotation loggers' max count for RotationLogger, defaults to no limits, and it cannot less
// than one. Note that maxAge and maxCount cannot be set at the same time.
func WithRotationMaxCount(count int32) Option {
	return func(o *loggerOptions) {
		if count < 0 {
			count = 0
		}
		o.rotationMaxCount = count
	}
}

// RotationLogger represents a rotation logger, which will gets automatically rotated when new file created. Some codes and interfaces are referred
// from https://github.com/lestrrat-go/file-rotatelogs.
type RotationLogger struct {
	option      *loggerOptions
	namePattern string
	globPattern string

	mu             sync.RWMutex
	currFile       *os.File
	currBasename   string
	currGeneration uint32
	currFilename   string
}

var _ io.WriteCloser = (*RotationLogger)(nil)

var (
	errEmptyFilenamePattern   = errors.New("xrotation: empty filename pattern is not allowed")
	errRotationMaxAgeMaxCount = errors.New("xrotation: rotation max age and max count can not be set at the same time")
)

const (
	errInvalidFilenamePattern = "xrotation: filename pattern `%s` is invalid: %w"
)

// New creates a RotationLogger with given filename pattern (in C-style / strftime) and Option-s, returns error if you give invalid options.
//
// Example:
// 	rl, err := New(
// 		"console.%Y%m%d.log",
// 		WithSymlinkFilename("console.current.log"),
// 		WithClock(xtime.UTC),
// 		WithForceNewFile(false),
// 		WithRotationSize(20*1024*1024),     // 20M
// 		WithRotationTime(24*time.Hour),     // 1d
// 		WithRotationMaxAge(7*24*time.Hour), // 7d
// 	)
func New(pattern string, options ...Option) (*RotationLogger, error) {
	opt := &loggerOptions{}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}
	if opt.nowClock == nil {
		opt.nowClock = xtime.Local
	}
	if opt.rotationTime == 0 {
		opt.rotationTime = 24 * time.Hour
	}

	// check options
	if pattern == "" {
		return nil, errEmptyFilenamePattern
	}
	if opt.rotationMaxAge > 0 && opt.rotationMaxCount > 0 {
		return nil, errRotationMaxAgeMaxCount
	}
	if opt.rotationMaxAge == 0 && opt.rotationMaxCount == 0 {
		opt.rotationMaxAge = 7 * 24 * time.Hour
	}

	// check filename pattern
	_, err := xtime.StrftimeInString(pattern, time.Now())
	if err != nil {
		return nil, fmt.Errorf(errInvalidFilenamePattern, pattern, err)
	}
	globPattern := xtime.StrftimeToGlobPattern(pattern)
	_, err = filepath.Match(globPattern, "")
	if err != nil {
		return nil, fmt.Errorf(errInvalidFilenamePattern, pattern, err)
	}

	logger := &RotationLogger{option: opt, namePattern: pattern, globPattern: globPattern}
	return logger, nil
}

// Write implements the io.Writer interface, it writes given bytes to file, and does rotation when a new file is created.
func (r *RotationLogger) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	writer, err := r.getRotatedWriter(false) // in some cases, it is no need to do rotation
	if err != nil {
		return 0, err
	}
	return writer.Write(p)
}

// Rotate rotates the logger files first manually, returns error when new file is unavailable to get, or rotate failed.
func (r *RotationLogger) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.getRotatedWriter(true) // rotation will be done in all cases
	return err
}

// CurrentFilename returns the current file name that the RotationLogger is writing to.
func (r *RotationLogger) CurrentFilename() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.currFilename
}

// Close implements the io.Closer interface, it closes the opened file, you can also call Write later because the closed file will be opened again.
func (r *RotationLogger) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.currFile == nil {
		return nil
	}
	_ = r.currFile.Close()

	// initialize all the states
	r.currFile = nil
	r.currBasename = ""
	r.currGeneration = 0
	r.currFilename = ""
	return nil
}

// ===================
// core implementation
// ===================

// These unexported variables are only used for testing.
var (
	_t_testHookMkdir   func()
	_t_testHookSymlink [3]func() string
)

const (
	errCreateDirectory  = "xrotation: failed to create directory `%s`: %w"
	errOpenOrCreateFile = "xrotation: failed to open or create file `%s`: %w"
	warnCreateSymlink   = "xrotation warning: failed to create symlink for `%s`: %v"
	warnDoRotation      = "xrotation warning: failed to rotate: [%v]"
	errDoRotation       = "xrotation: failed to rotate: [%w]"
)

// getRotatedWriter does: check whether it needs to create new file, create a unique-filename file, generate symlink and do rotation.
func (r *RotationLogger) getRotatedWriter(rotateManually bool) (io.Writer, error) {
	// check whether it needs to create new file
	createNewFile := false
	generation := r.currGeneration
	basename, _ := xtime.StrftimeInString(r.namePattern, xtime.TruncateTime(r.option.nowClock.Now(), r.option.rotationTime))
	if r.currFilename == "" { // invoke initially
		fi, err := os.Stat(basename)
		if existed := !os.IsNotExist(err); !existed || r.option.forceNewFile || (r.option.rotationSize > 0 && fi.Size() >= r.option.rotationSize) {
			createNewFile = true // 4.
			if existed {
				generation = 1
			} else {
				generation = 0
			}
		} else {
			createNewFile = false // 3.
		}
	} else if basename != r.currBasename { // new basename
		createNewFile = true // 2.
		generation = 0
	} else { // check whether file exceeds rotation size
		fi, err := os.Stat(r.currFilename)
		if err == nil && r.option.rotationSize > 0 && fi.Size() >= r.option.rotationSize {
			createNewFile = true // 2.
			generation++
		}
	}

	// cases the following code deals with:
	// 1.1. !createNewFile && currFile != nil && !rotateManually => return directly (happens in most cases)
	// 1.2. !createNewFile && currFile != nil && rotateManually  => close the file, open it again, check symlink and do rotate (happens when calling Rotate())
	// 2.   createNewFile  && currFile != nil                    => create a new file with basename or basename_x (happens when rotation basename changes or file exceeds rotation size)
	// 3.   !createNewFile && currFile == nil                    => open the old file with basename (happens when the first time call this method, with file exists)
	// 4.   createNewFile  && currFile == nil                    => same with 2 (happens when the first time call this method, with file not exists, or forceNewFile, or file size exceeds)
	filename := basename
	if !createNewFile && r.currFile != nil {
		if !rotateManually {
			// also don't check symlink and do rotation
			return r.currFile, nil
		}
		filename = r.currFilename
		// close first, later it will be reopened
		_ = r.currFile.Close()
		r.currFile = nil
	}

	// generate a non-conflict filename
	if createNewFile {
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
	}

	// open or create the file
	if createNewFile {
		dirname := filepath.Dir(filename)
		if _, err := os.Stat(dirname); os.IsNotExist(err) {
			if _t_testHookMkdir != nil { // only used when testing
				_t_testHookMkdir()
			}
			err := os.MkdirAll(dirname, 0755) // drwxr-xr-x
			if err != nil {
				return nil, fmt.Errorf(errCreateDirectory, dirname, err)
			}
		}
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) // -rwxr--r--
	if err != nil {
		return nil, fmt.Errorf(errOpenOrCreateFile, filename, err)
	}

	// generate a symlink and do rotation
	if r.option.symlinkFilename != "" {
		// only when need to create symlink links to current filename
		err := createSymlink(filename, r.option.symlinkFilename)
		if err != nil {
			// Windows: "A required privilege is not held by the client"
			log.Printf(warnCreateSymlink, filename, err) // ignore symlink error
		}
	}
	if createNewFile || rotateManually {
		// only when need to create a new file or rotate manually
		err := doRotation(r.globPattern, r.option.nowClock.Now(), r.option.rotationMaxAge, r.option.rotationMaxCount) // errors returned from os.Remove
		if err != nil {
			if !rotateManually {
				log.Printf(warnDoRotation, err) // ignore rotation error
			} else {
				_ = file.Close()
				return nil, fmt.Errorf(errDoRotation, err)
			}
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
		if _t_testHookSymlink[0] != nil { // only used when testing
			_t_testHookSymlink[0]()
		}
		err := os.MkdirAll(linkDirname, 0755)
		if err != nil {
			// hint: no need for "xrotation: " prefix
			return fmt.Errorf("failed to create directory `%s`: %w", linkDirname, err)
		}
	}

	// check the relative path of destination
	destinationPath, _ := filepath.Abs(filename)
	linkDirnamePath, _ := filepath.Abs(linkDirname)
	if _t_testHookSymlink[1] != nil {
		linkDirnamePath = _t_testHookSymlink[1]()
	}
	destination, err := filepath.Rel(linkDirnamePath, destinationPath)
	if err != nil {
		return fmt.Errorf("failed to evaluate the relative path from `%s` to `%s`: %w", destinationPath, linkDirnamePath, err)
	}

	// make symlink and rename to the link file
	tempLinkname := filename + "_symlink"
	if _, err := os.Stat(tempLinkname); err == nil {
		_ = os.Remove(tempLinkname)
	}
	if _t_testHookSymlink[2] != nil {
		_t_testHookSymlink[2]()
	}
	err = os.Symlink(destination, tempLinkname)
	if err != nil {
		return fmt.Errorf("failed to create symlink `%s`: %w", tempLinkname, err)
	}
	err = os.Rename(tempLinkname, linkname)
	if err != nil {
		return fmt.Errorf("failed to rename symlink `%s` to `%s`: %w", tempLinkname, linkname, err)
	}
	return nil
}

// doRotation does the real rotation work, this will rotate for loggers' max age or for loggers' max count, and remove all unlinked files.
func doRotation(globPattern string, now time.Time, maxAge time.Duration, maxCount int32) error {
	// get matches by glob pattern
	matches, _ := filepath.Glob(globPattern) // error is always nil if in safe manner, here ignore it
	unlinkFiles := make([]string, 0)

	// I) rotate for max age
	if maxAge > 0 {
		cutoffDuration := now.Add(-1 * maxAge)
		for _, match := range matches {
			fi, err := os.Lstat(match)
			if err != nil || (fi.Mode()&os.ModeSymlink) == os.ModeSymlink {
				continue
			}
			if fi.ModTime().Before(cutoffDuration) {
				unlinkFiles = append(unlinkFiles, match)
			}
		}
	}

	// II) rotate for max count
	if count := int(maxCount); count > 0 {
		type nameTimeTuple struct {
			name string
			mod  time.Time
		}
		pairs := make([]nameTimeTuple, 0, len(matches))
		for _, match := range matches {
			fi, err := os.Lstat(match)
			if err != nil || (fi.Mode()&os.ModeSymlink) == os.ModeSymlink {
				continue
			}
			pairs = append(pairs, nameTimeTuple{match, fi.ModTime()})
		}
		if len(pairs) > count {
			sort.Slice(pairs, func(i, j int) bool { return pairs[i].mod.Before(pairs[j].mod) })
			for _, fi := range pairs[:len(pairs)-count] {
				unlinkFiles = append(unlinkFiles, fi.name)
			}
		}
	}

	// expand unlinkFiles for file with "xxx_*" name
	if len(unlinkFiles) == 0 {
		return nil
	}
	moreMatches, _ := filepath.Glob(globPattern + "_*") // also ignore error
	if len(moreMatches) > 0 {
		more := make([]string, 0)
		for _, match := range moreMatches {
			for _, path := range unlinkFiles {
				if strings.HasPrefix(match, path) {
					more = append(more, match)
				}
			}
		}
		unlinkFiles = append(unlinkFiles, more...)
	}

	// remove unlinked files
	errs := make([]error, 0)
	for _, path := range unlinkFiles {
		err := os.Remove(path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return xerror.Combine(errs...)
}
