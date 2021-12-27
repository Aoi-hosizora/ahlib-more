package xrotation

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		_, err := New()
		xtesting.NotNil(t, err)
		_, err = New(WithFilenamePattern("test.log"))
		xtesting.Nil(t, err)
		_, err = New(WithFilenamePattern("test.log%"))
		xtesting.NotNil(t, err)
		_, err = New(WithFilenamePattern("[x-]"), WithForceNewFile(true))
		xtesting.NotNil(t, err)
		_, err = New(WithFilenamePattern("test.log"), WithRotationMaxAge(1), WithRotationMaxCount(1))
		xtesting.NotNil(t, err)
	})

	t.Run("values", func(t *testing.T) {
		rl, err := New(
			WithFilenamePattern(""), WithFilenamePattern("test.log"), WithSymlinkFilename("test.curr.log"), WithClock(UTC),
			WithUseSingleFile(false), WithForceNewFile(true),
			WithRotationTime(1), WithRotationSize(1), WithRotationMaxAge(0), WithRotationMaxCount(-1),
		)
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.option.filenamePattern, "test.log")
		xtesting.Equal(t, rl.option.symlinkFilename, "test.curr.log")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), time.Duration(0))
		xtesting.Equal(t, rl.option.useSingleFile, false)
		xtesting.Equal(t, rl.option.forceNewFile, true)
		xtesting.Equal(t, rl.option.rotationTime, time.Duration(1))
		xtesting.Equal(t, rl.option.rotationSize, int64(1))
		xtesting.Equal(t, rl.option.rotationMaxAge, 7*24*time.Hour)
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(0))

		rl, err = New(
			WithFilenamePattern("test.log"), WithSymlinkFilename("test.curr.log"), WithSymlinkFilename(""), WithClock(nil),
			WithUseSingleFile(true), WithForceNewFile(false),
			WithRotationTime(-1), WithRotationSize(-1), WithRotationMaxAge(-1), WithRotationMaxCount(1),
		)
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.option.filenamePattern, "test.log")
		xtesting.Equal(t, rl.option.symlinkFilename, "")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), xtime.LocationDuration(time.Local))
		xtesting.Equal(t, rl.option.useSingleFile, true)
		xtesting.Equal(t, rl.option.forceNewFile, false)
		xtesting.Equal(t, rl.option.rotationTime, 24*time.Hour)
		xtesting.Equal(t, rl.option.rotationSize, int64(0))
		xtesting.Equal(t, rl.option.rotationMaxAge, time.Duration(0))
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(1))
	})
}

func removeLoggers() {
	matches, _ := filepath.Glob("*.log*")
	for _, match := range matches {
		_ = os.Remove(match)
	}
}

func testFileContent(t *testing.T, name, content string) {
	bs, err := ioutil.ReadFile(name)
	xtesting.Nil(t, err)
	if err == nil {
		xtesting.Equal(t, string(bs), content, " | name: %s, content: %s", name, strings.ReplaceAll(content, "\n", "\\n"))
	}
}

func testFileExistence(t *testing.T, name string, exist bool) {
	_, err := os.Stat(name)
	if exist {
		xtesting.Nil(t, err, " | name: %s, exist: %t", name, exist)
	} else {
		xtesting.True(t, os.IsNotExist(err), " | name: %s, exist: %t", name, exist)
	}
}

func TestWriter(t *testing.T) {
	t.Run("force new file", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		f, err := os.Create("logger.log")
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		rl, err := New(WithFilenamePattern("logger.log"), WithForceNewFile(true)) // cannot del logger.log
		xtesting.NotNil(t, err)
		testFileExistence(t, "logger.log", true)

		xtesting.Nil(t, f.Close())
		rl, err = New(WithFilenamePattern("logger.log"), WithForceNewFile(true))
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", false)
		_, err = fmt.Fprintf(rl, "hello world") // <- create new file, logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileContent(t, "logger.log", "hello world")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		xtesting.Nil(t, rl.Close())

		rl, err = New(WithFilenamePattern("logger.log"))
		xtesting.Nil(t, err)
		xtesting.Nil(t, rl.Close()) // r.currFile == nil
	})

	t.Run("some write demos", func(t *testing.T) {
	})

	t.Run("rotate directly", func(t *testing.T) {

	})
}

func TestRotate(t *testing.T) {

}
