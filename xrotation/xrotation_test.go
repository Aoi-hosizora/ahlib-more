package xrotation

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		_, err := New("")
		xtesting.NotNil(t, err)
		_, err = New("test.log")
		xtesting.Nil(t, err)
		_, err = New("test.log%")
		xtesting.NotNil(t, err)
		_, err = New("[x-]", WithForceNewFile(true))
		xtesting.NotNil(t, err)
		_, err = New("test.log", WithRotationMaxAge(1), WithRotationMaxCount(0))
		xtesting.Nil(t, err)
		_, err = New("test.log", WithRotationMaxAge(0), WithRotationMaxCount(1))
		xtesting.Nil(t, err)
		_, err = New("test.log", WithRotationMaxAge(1), WithRotationMaxCount(1))
		xtesting.NotNil(t, err)
	})

	t.Run("values", func(t *testing.T) {
		rl, err := New("test.log")
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.namePattern, "test.log")
		xtesting.Equal(t, rl.globPattern, "test.log")
		xtesting.Equal(t, rl.option.symlinkFilename, "")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), xtime.LocationDuration(time.Local))
		xtesting.Equal(t, rl.option.forceNewFile, false)
		xtesting.Equal(t, rl.option.rotationTime, 24*time.Hour)
		xtesting.Equal(t, rl.option.rotationSize, int64(0))
		xtesting.Equal(t, rl.option.rotationMaxAge, 7*24*time.Hour)
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(0))

		rl, err = New(
			"x", WithSymlinkFilename(""), WithSymlinkFilename("test.curr.log"), WithClock(nil),
			WithForceNewFile(false), WithRotationTime(-1), WithRotationSize(-1), WithRotationMaxAge(-1), WithRotationMaxCount(-1),
		)
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.namePattern, "x")
		xtesting.Equal(t, rl.globPattern, "x")
		xtesting.Equal(t, rl.option.symlinkFilename, "test.curr.log")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), xtime.LocationDuration(time.Local))
		xtesting.Equal(t, rl.option.forceNewFile, false)
		xtesting.Equal(t, rl.option.rotationTime, 24*time.Hour)
		xtesting.Equal(t, rl.option.rotationSize, int64(0))
		xtesting.Equal(t, rl.option.rotationMaxAge, 7*24*time.Hour)
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(0))

		rl, err = New(
			"test-%Y%m%d.log", WithSymlinkFilename(""), WithClock(xtime.UTC), WithForceNewFile(true),
			WithRotationTime(time.Hour), WithRotationSize(256), WithRotationMaxAge(time.Hour*15*24), WithRotationMaxCount(0),
		)
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.namePattern, "test-%Y%m%d.log")
		xtesting.Equal(t, rl.globPattern, "test-*.log")
		xtesting.Equal(t, rl.option.symlinkFilename, "")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), time.Duration(0))
		xtesting.Equal(t, rl.option.forceNewFile, true)
		xtesting.Equal(t, rl.option.rotationTime, time.Hour)
		xtesting.Equal(t, rl.option.rotationSize, int64(256))
		xtesting.Equal(t, rl.option.rotationMaxAge, 15*24*time.Hour)
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(0))
	})
}

func removeLoggers() {
	matches, _ := filepath.Glob("*.log*")
	for _, match := range matches {
		err := os.Remove(match)
		if err != nil {
			log.Printf("removeLoggers: os.Remove failed on `%s`", match)
		}
	}
	err := os.RemoveAll("./_test")
	if err != nil {
		log.Printf("removeLoggers: os.Remove failed on `%s`", "./_test")
	}
}

func testFileContent(t *testing.T, name, content string) {
	xtesting.SetExtraSkip(1)
	defer xtesting.SetExtraSkip(0)
	bs, err := ioutil.ReadFile(name)
	xtesting.Nil(t, err)
	if err == nil {
		xtesting.Equal(t, string(bs), content)
	}
}

func testFileExistence(t *testing.T, name string, exist bool) {
	xtesting.SetExtraSkip(1)
	defer xtesting.SetExtraSkip(0)
	_, err := os.Stat(name)
	if exist {
		xtesting.Nil(t, err)
	} else {
		xtesting.True(t, os.IsNotExist(err))
	}
}

func TestSimpleWrite(t *testing.T) {
	t.Run("strftime", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("logger.%Y%m%d.log", WithRotationTime(time.Hour*24), WithClock(clock))
		xtesting.Nil(t, rl.Close()) // <- r.currFile == nil

		rl, _ = New("logger.%Y%m%d.log", WithRotationTime(time.Hour*24), WithClock(clock))
		_, err := fmt.Fprintf(rl, "hello world 1") // <- create a new file, logger.01.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.20010101.log", true)
		testFileContent(t, "logger.20010101.log", "hello world 1")

		*pNow = time.Date(2002, 2, 2, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		_, err = fmt.Fprintf(rl, "hello world 2") // <- create a new file, logger.02.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.20010101.log", true)
		testFileExistence(t, "logger.20020202.log", true)
		testFileContent(t, "logger.20010101.log", "hello world 1")
		testFileContent(t, "logger.20020202.log", "hello world 2")

		_, err = fmt.Fprintf(rl, "hello world 3") // <- use logger.02.log
		testFileContent(t, "logger.20020202.log", "hello world 2hello world 3")
		xtesting.Nil(t, rl.Close()) // normal close
		xtesting.Nil(t, rl.currFile)
		xtesting.Zero(t, rl.currBasename)
		xtesting.Zero(t, rl.currGeneration)
		xtesting.Zero(t, rl.currFilename)
		_, err = fmt.Fprintf(rl, "hello world 4") // <- normally use logger.02.log
		testFileContent(t, "logger.20020202.log", "hello world 2hello world 3hello world 4")
		xtesting.Nil(t, rl.Close()) // normal close
	})

	t.Run("reopen and ForceNewFile", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, _ := New("logger.log", WithRotationTime(time.Hour*24), WithRotationSize(15), WithForceNewFile(false))
		_, err := fmt.Fprintf(rl, "hello world 1\n") // <- 14, create logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileContent(t, "logger.log", "hello world 1\n")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationTime(time.Hour*24), WithRotationSize(15), WithForceNewFile(false))
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\n")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationTime(time.Hour*24), WithRotationSize(15), WithForceNewFile(false))
		_, err = fmt.Fprintf(rl, "hello world 3\n") // <- create logger.log_1
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log_1", true)
		testFileExistence(t, "logger.log_2", false)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\n")
		testFileContent(t, "logger.log_1", "hello world 3\n")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationTime(time.Hour*24), WithRotationSize(15), WithForceNewFile(true))
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- create logger.log_2
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log_2", true)
		testFileExistence(t, "logger.log_3", false)
		testFileContent(t, "logger.log_1", "hello world 3\n")
		testFileContent(t, "logger.log_2", "hello world 4\n")
		xtesting.Nil(t, rl.Close())
	})
}

func TestWrite(t *testing.T) {
	t.Run("some write demos", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, _ := New("logger.log", WithRotationSize(15))
		xtesting.Equal(t, rl.CurrentFilename(), "")
		_, err := fmt.Fprintf(rl, "hello world 1\n") // <- 14, create logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "hello world 1\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")

		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")

		_, err = fmt.Fprintf(rl, "hello world 3\n") // <- create: logger.log_1
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log_1", true)
		testFileExistence(t, "logger.log_2", false)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\n")
		testFileContent(t, "logger.log_1", "hello world 3\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log_1")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationSize(29))
		xtesting.Equal(t, rl.CurrentFilename(), "")
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- use logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", true)
		testFileExistence(t, "logger.log_2", false)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\nhello world 4\n")
		testFileContent(t, "logger.log_1", "hello world 3\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")

		_, err = fmt.Fprintf(rl, "hello world 5\n") // <- create: logger.log_2
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log_2", true)
		testFileExistence(t, "logger.log_3", false)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\nhello world 4\n")
		testFileContent(t, "logger.log_1", "hello world 3\n")
		testFileContent(t, "logger.log_2", "hello world 5\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log_2")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationSize(29))
		xtesting.Equal(t, rl.CurrentFilename(), "")
		_, err = fmt.Fprintf(rl, "hello world 6\n") // <- create: logger.log_3
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log_3", true)
		testFileExistence(t, "logger.log_4", false)
		testFileContent(t, "logger.log_2", "hello world 5\n")
		testFileContent(t, "logger.log_3", "hello world 6\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log_3")
		xtesting.Nil(t, rl.Close())
	})

	t.Run("simple rotate", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, _ := New("logger.log", WithRotationSize(15))
		xtesting.Nil(t, rl.Rotate()) // <- create logger.log
		testFileExistence(t, "logger.log", true)
		testFileContent(t, "logger.log", "")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		xtesting.Nil(t, rl.Rotate()) // <- use logger.log
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationSize(15))
		xtesting.Nil(t, rl.Rotate()) // <- use logger.log
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		_, err := fmt.Fprintf(rl, "hello world 1\n") // <- write to logger.log
		xtesting.Nil(t, err)
		testFileContent(t, "logger.log", "hello world 1\n")
		xtesting.Nil(t, rl.Rotate()) // <- use logger.log
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "hello world 1\n")
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- write to logger.log
		xtesting.Nil(t, err)
		testFileContent(t, "logger.log", "hello world 1\nhello world 2\n")
		xtesting.Nil(t, rl.Rotate()) // <- create logger_1.log
		testFileExistence(t, "logger.log_1", true)
		testFileExistence(t, "logger.log_2", false)
		testFileContent(t, "logger.log_1", "") // <- empty logger.log_1
		xtesting.Nil(t, rl.Close())

		rl, _ = New("logger.log", WithRotationSize(15))
		xtesting.Nil(t, rl.Rotate()) // <- create logger.log_2
		testFileExistence(t, "logger.log_1", true)
		testFileExistence(t, "logger.log_2", true)
		testFileContent(t, "logger.log_2", "")
		_, err = fmt.Fprintf(rl, "hello world 3\n") // <- write to logger.log_2
		testFileContent(t, "logger.log_1", "")
		testFileContent(t, "logger.log_2", "hello world 3\n")
		xtesting.Nil(t, rl.Close())
	})
}

func TestRotate(t *testing.T) {
	t.Run("max age", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("logger.%Y%m%d.log", WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxAge(3*24*time.Hour), WithClock(clock))

		*pNow = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		_, _ = fmt.Fprintf(rl, "hello world 4\n")
		_, _ = fmt.Fprintf(rl, "hello world 5\n") // <- logger.01.log_2
		_ = os.Chtimes("logger.20010101.log", now, now)
		_ = os.Chtimes("logger.20010101.log_1", now, now)
		_ = os.Chtimes("logger.20010101.log_2", now, now)
		*pNow = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 6\n") // <- logger.02.log
		_, _ = fmt.Fprintf(rl, "hello world 7\n")
		_, _ = fmt.Fprintf(rl, "hello world 8\n") // <- logger.02.log_1
		_ = os.Chtimes("logger.20010102.log", now, now)
		_ = os.Chtimes("logger.20010102.log_1", now, now)
		*pNow = xtime.SetDay(now, 3)              // 3d
		_, _ = fmt.Fprintf(rl, "hello world 9\n") // <- logger.03.log
		_ = os.Chtimes("logger.20010103.log", now, now)
		testFileExistence(t, "logger.20010101.log", true)
		testFileExistence(t, "logger.20010101.log_1", true)
		testFileExistence(t, "logger.20010101.log_2", true)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)

		*pNow = xtime.SetHour(now, 1)
		*pNow = xtime.SetDay(now, 4) // 4d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)
		*pNow = xtime.SetDay(now, 5)         // 5d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", true)
		*pNow = xtime.SetDay(now, 6) // 6d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.03.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", false)
		_ = rl.Close()
	})

	t.Run("max count", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 1, 1, 1, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("logger.%Y%m%d.log", WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxCount(3), WithClock(clock))

		*pNow = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		_, _ = fmt.Fprintf(rl, "hello world 4\n")
		_, _ = fmt.Fprintf(rl, "hello world 5\n") // <- logger.01.log_2
		*pNow = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 6\n") // <- logger.02.log
		_, _ = fmt.Fprintf(rl, "hello world 7\n")
		_, _ = fmt.Fprintf(rl, "hello world 8\n") // <- logger.02.log_1
		*pNow = xtime.SetDay(now, 3)              // 3d
		_, _ = fmt.Fprintf(rl, "hello world 9\n") // <- logger.03.log
		testFileExistence(t, "logger.20010101.log", true)
		testFileExistence(t, "logger.20010101.log_1", true)
		testFileExistence(t, "logger.20010101.log_2", true)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)

		*pNow = xtime.SetHour(now, 1)
		*pNow = xtime.SetDay(now, 4) // 4d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)
		*pNow = xtime.SetDay(now, 5)         // 5d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", true)
		*pNow = xtime.SetDay(now, 6) // 6d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.03.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", false)

		f1, err := os.OpenFile("logger.20010104.log", os.O_APPEND, 0644)
		xtesting.Nil(t, err)
		f2, err := os.OpenFile("logger.20010105.log", os.O_APPEND, 0644)
		xtesting.Nil(t, err)
		*pNow = xtime.SetDay(now, 7) // 7d1h
		_, err = fmt.Fprintf(rl, "")
		xtesting.Nil(t, err)            // <- Warning: failed to remove logger.04.log
		xtesting.NotNil(t, rl.Rotate()) // <- need to rotate, error: failed to remove logger.04.log
		testFileExistence(t, "logger.20010104.log", true)
		*pNow = xtime.SetDay(now, 8)    // 8d1h
		xtesting.NotNil(t, rl.Rotate()) // <- error: failed to remove logger.04.log, logger.05.log
		testFileExistence(t, "logger.20010104.log", true)
		testFileExistence(t, "logger.20010105.log", true)
		xtesting.Nil(t, f1.Close())
		xtesting.Nil(t, f2.Close())
		*pNow = xtime.SetDay(now, 9) // 9d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.04.log and logger.05.log
		testFileExistence(t, "logger.20010104.log", false)
		testFileExistence(t, "logger.20010105.log", false)
		_ = rl.Close()
	})
}

func TestInDifferentDir(t *testing.T) {
	t.Run("some write demos", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, _ := New("./_test/_test/logger.log", WithRotationSize(15))
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- create logger.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n") // <- use logger.log
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- create: logger.log_1
		testFileExistence(t, "./_test/_test/logger.log", true)
		testFileExistence(t, "./_test/_test/logger.log_1", true)
		testFileExistence(t, "./_test/_test/logger.log_2", false)
		testFileContent(t, "./_test/_test/logger.log", "hello world 1\nhello world 2\n")
		testFileContent(t, "./_test/_test/logger.log_1", "hello world 3\n")
		xtesting.Equal(t, rl.CurrentFilename(), "./_test/_test/logger.log_1")
		xtesting.Nil(t, rl.Close())

		rl, _ = New("./_test/_test/logger.log", WithRotationSize(29))
		_, _ = fmt.Fprintf(rl, "hello world 4\n") // <- use logger.log
		_, _ = fmt.Fprintf(rl, "hello world 5\n") // <- create: logger.log_2
		testFileExistence(t, "./_test/_test/logger.log", true)
		testFileExistence(t, "./_test/_test/logger.log_1", true)
		testFileExistence(t, "./_test/_test/logger.log_2", true)
		testFileExistence(t, "./_test/_test/logger.log_3", false)
		testFileContent(t, "./_test/_test/logger.log", "hello world 1\nhello world 2\nhello world 4\n")
		testFileContent(t, "./_test/_test/logger.log_1", "hello world 3\n")
		testFileContent(t, "./_test/_test/logger.log_2", "hello world 5\n")
		xtesting.Equal(t, rl.CurrentFilename(), "./_test/_test/logger.log_2")
		xtesting.Nil(t, rl.Close())
	})

	t.Run("max age", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("./_test/_test/logger.%Y%m%d.log", WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxAge(2*24*time.Hour), WithClock(clock))

		*pNow = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		_ = os.Chtimes("./_test/_test/logger.20010101.log", now, now)
		_ = os.Chtimes("./_test/_test/logger.20010101.log_1", now, now)
		*pNow = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 4\n") // <- logger.02.log
		_ = os.Chtimes("./_test/_test/logger.20010102.log", now, now)
		testFileExistence(t, "./_test/_test/logger.20010101.log", true)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", true)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)

		*pNow = xtime.SetHour(now, 1)
		*pNow = xtime.SetDay(now, 3) // 3d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)
		*pNow = xtime.SetDay(now, 4)         // 4d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", false)
		_ = rl.Close()
	})

	t.Run("max count", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 1, 1, 1, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("./_test/_test/logger.%Y%m%d.log", WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxCount(2), WithClock(clock))

		*pNow = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		*pNow = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 4\n") // <- logger.02.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", true)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", true)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)

		*pNow = xtime.SetHour(now, 1)
		*pNow = xtime.SetDay(now, 3) // 3d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)
		*pNow = xtime.SetDay(now, 4)         // 4d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", false)
		_ = rl.Close()
	})

	t.Run("cover errors", func(t *testing.T) {
		removeLoggers()
		now := time.Date(2001, 1, 1, 1, 1, 1, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)

		rl, _ := New("./_test/logger.%Y%m%d.log", WithRotationTime(time.Hour*24), WithClock(clock))
		*pNow = xtime.SetDay(now, 1)
		_t_testHookMkdir = func() {
			f, err := os.OpenFile("_test", os.O_CREATE, 0644)
			xtesting.Nil(t, err)
			xtesting.Nil(t, f.Close())
		}
		xtesting.NotNil(t, rl.Rotate()) // Rotate: MkdirAll failed
		removeLoggers()
		_t_testHookMkdir = func() {
			err := os.MkdirAll("./_test/logger.20010101.log", 0755)
			xtesting.Nil(t, err)
		}
		_, err := fmt.Fprintf(rl, "test OpenFile")
		xtesting.NotNil(t, err) // OpenFile failed
		_ = rl.Close()
		_t_testHookMkdir = nil
		removeLoggers()
	})
}

func TestSymlink(t *testing.T) {
	t.Run("same directory", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("logger.%Y%m%d.log", WithSymlinkFilename("logger.current.log"), WithRotationSize(15), WithRotationMaxCount(2), WithClock(clock))
		_, err := fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		testFileContent(t, "logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileExistence(t, "logger.current.log", true)
		fi, err := os.Lstat("logger.current.log")
		xtesting.Nil(t, err)
		xtesting.True(t, fi.Mode()&os.ModeSymlink == os.ModeSymlink)
		testFileContent(t, "logger.current.log", "hello world 1\nhello world 2\n")

		_, err = fmt.Fprintf(rl, "hello world 3\n") // <- create: logger.01.log_1
		xtesting.Nil(t, err)
		testFileContent(t, "logger.20010101.log_1", "hello world 3\n")
		testFileContent(t, "logger.current.log", "hello world 3\n")
		*pNow = xtime.SetDay(now, 2)                // 2d
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- create: logger.02.log
		testFileContent(t, "logger.20010102.log", "hello world 4\n")
		testFileContent(t, "logger.current.log", "hello world 4\n")
		*pNow = xtime.SetDay(now, 3)                // 3d
		_, err = fmt.Fprintf(rl, "hello world 5\n") // <- create: logger.03.log, need rotate
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileContent(t, "logger.20010103.log", "hello world 5\n")
		testFileContent(t, "logger.current.log", "hello world 5\n")
		_ = rl.Close()
		removeLoggers()

		rl, _ = New("logger.%Y%m%d.log", WithSymlinkFilename("logger.current.log"), WithRotationSize(15), WithRotationMaxAge(time.Hour*24*2), WithClock(clock), WithForceNewFile(true))
		*pNow = xtime.SetDay(now, 1)                // 1d
		_, err = fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		_ = os.Chtimes("logger.20010101.log", now, now)
		testFileContent(t, "logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileContent(t, "logger.current.log", "hello world 1\nhello world 2\n")

		_, err = fmt.Fprintf(rl, "hello world 3\n") // <- create: logger.01.log_1
		xtesting.Nil(t, err)
		_ = os.Chtimes("logger.20010101.log_1", now, now)
		testFileContent(t, "logger.20010101.log_1", "hello world 3\n")
		testFileContent(t, "logger.current.log", "hello world 3\n")
		*pNow = xtime.SetDay(now, 2)                // 2d
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- create: logger.02.log
		_ = os.Chtimes("logger.20010102.log", now, now)
		testFileContent(t, "logger.20010102.log", "hello world 4\n")
		testFileContent(t, "logger.current.log", "hello world 4\n")
		*pNow = xtime.SetHour(now, 1)
		*pNow = xtime.SetDay(now, 3) // 3d1h
		f, err := os.OpenFile("logger.20010103.log_symlink", os.O_CREATE, 0644)
		xtesting.Nil(t, err) // <- fake symlink, need to be deleted
		xtesting.Nil(t, f.Close())
		_, err = fmt.Fprintf(rl, "hello world 5\n") // <- create: logger.03.log, need rotate
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileContent(t, "logger.20010103.log", "hello world 5\n")
		testFileContent(t, "logger.current.log", "hello world 5\n")
		_ = rl.Close()
	})

	t.Run("different directory", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("_test/logger.%Y%m%d.log", WithSymlinkFilename("_test/_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		_, err := fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		testFileContent(t, "_test/logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileExistence(t, "_test/_test/logger.current.log", true)
		testFileContent(t, "_test/_test/logger.current.log", "hello world 1\nhello world 2\n")
		_ = rl.Close()
		removeLoggers()

		rl, _ = New("_test/_test/logger.%Y%m%d.log", WithSymlinkFilename("_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		_, err = fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		testFileContent(t, "_test/_test/logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileExistence(t, "_test/logger.current.log", true)
		testFileContent(t, "_test/logger.current.log", "hello world 1\nhello world 2\n")
		_ = rl.Close()
		removeLoggers()

		rl, _ = New("_test/logger.%Y%m%d.log", WithSymlinkFilename("_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		_, err = fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		testFileContent(t, "_test/logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileExistence(t, "_test/logger.current.log", true)
		testFileContent(t, "_test/logger.current.log", "hello world 1\nhello world 2\n")
		_ = rl.Close()
		removeLoggers()
	})

	t.Run("cover errors", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		pNow := &now
		clock := xtime.CustomClock(pNow)
		rl, _ := New("logger.%Y%m%d.log", WithSymlinkFilename("_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		*pNow = xtime.SetDay(now, 1)
		_t_testHookSymlink[0] = func() string {
			f, err := os.OpenFile("_test", os.O_CREATE, 0644)
			xtesting.Nil(t, err)
			xtesting.Nil(t, f.Close())
			return ""
		}
		xtesting.Nil(t, rl.Rotate()) // MkdirAll failed
		_t_testHookSymlink[0] = nil

		*pNow = xtime.SetDay(now, 2)
		_t_testHookSymlink[1] = func() string {
			return ".." // hack
		}
		xtesting.Nil(t, rl.Rotate()) // Rel failed
		_t_testHookSymlink[1] = nil

		*pNow = xtime.SetDay(now, 3)
		_t_testHookSymlink[2] = func() string {
			xtesting.Nil(t, os.MkdirAll("logger.20010103.log_symlink", 0755))
			return ""
		}
		xtesting.Nil(t, rl.Rotate()) // Symlink failed

		_ = os.RemoveAll("_test")
		*pNow = xtime.SetDay(now, 4)
		_t_testHookSymlink[2] = func() string {
			xtesting.Nil(t, os.MkdirAll("_test/logger.current.log", 0755))
			return ""
		}
		xtesting.Nil(t, rl.Rotate()) // Rename failed
		_t_testHookSymlink[2] = nil

		_ = rl.Close()
	})
}
