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
		_, err := New()
		xtesting.NotNil(t, err)
		_, err = New(WithFilenamePattern("test.log"))
		xtesting.Nil(t, err)
		_, err = New(WithFilenamePattern("test.log%"))
		xtesting.NotNil(t, err)
		_, err = New(WithFilenamePattern("[x-]"), WithForceNewFile(true))
		xtesting.NotNil(t, err)
		_, err = New(WithFilenamePattern("test.log"), WithRotationMaxAge(1), WithRotationMaxCount(1))
		xtesting.Nil(t, err)
		_, err = New(WithFilenamePattern("test.log"), WithRotationMaxAge(1), WithRotationMaxCount(2))
		xtesting.NotNil(t, err)
	})

	t.Run("values", func(t *testing.T) {
		rl, err := New(
			WithFilenamePattern(""), WithFilenamePattern("test.log"), WithSymlinkFilename("test.curr.log"), WithClock(UTC),
			WithForceNewFile(true), WithRotationTime(1), WithRotationSize(1), WithRotationMaxAge(0), WithRotationMaxCount(-1),
		)
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.option.filenamePattern, "test.log")
		xtesting.Equal(t, rl.option.symlinkFilename, "test.curr.log")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), time.Duration(0))
		xtesting.Equal(t, rl.option.forceNewFile, true)
		xtesting.Equal(t, rl.option.rotationTime, time.Duration(1))
		xtesting.Equal(t, rl.option.rotationSize, int64(1))
		xtesting.Equal(t, rl.option.rotationMaxAge, 7*24*time.Hour)
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(0))

		rl, err = New(
			WithFilenamePattern("test.log"), WithSymlinkFilename("test.curr.log"), WithSymlinkFilename(""), WithClock(nil),
			WithForceNewFile(false), WithRotationTime(-1), WithRotationSize(-1), WithRotationMaxAge(-1), WithRotationMaxCount(2),
		)
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.option.filenamePattern, "test.log")
		xtesting.Equal(t, rl.option.symlinkFilename, "")
		xtesting.Equal(t, xtime.LocationDuration(rl.option.nowClock.Now().Location()), xtime.LocationDuration(time.Local))
		xtesting.Equal(t, rl.option.forceNewFile, false)
		xtesting.Equal(t, rl.option.rotationTime, 24*time.Hour)
		xtesting.Equal(t, rl.option.rotationSize, int64(0))
		xtesting.Equal(t, rl.option.rotationMaxAge, time.Duration(0))
		xtesting.Equal(t, rl.option.rotationMaxCount, int32(2))
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

func TestWriter(t *testing.T) {
	t.Run("strftime", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, err := New(WithFilenamePattern("logger.%Y%m%d.log"), WithRotationTime(time.Hour*24), WithClock(clock))
		xtesting.Nil(t, err)

		_, err = fmt.Fprintf(rl, "hello world 1") // <- create a new file
		testFileExistence(t, "logger.20010101.log", true)
		testFileContent(t, "logger.20010101.log", "hello world 1")

		*nowP = time.Date(2002, 2, 2, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		_, err = fmt.Fprintf(rl, "hello world 2") // <- create a new file
		testFileExistence(t, "logger.20010101.log", true)
		testFileExistence(t, "logger.20020202.log", true)
		testFileContent(t, "logger.20010101.log", "hello world 1")
		testFileContent(t, "logger.20020202.log", "hello world 2")

		_, err = fmt.Fprintf(rl, "hello world 3") // <- use old file
		testFileContent(t, "logger.20020202.log", "hello world 2hello world 3")
		xtesting.Nil(t, rl.Close())
	})

	t.Run("force new file", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		f1, err := os.Create("logger.log")
		xtesting.Nil(t, err)
		f2, err := os.Create("logger.log_1")
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", true)
		rl, err := New(WithFilenamePattern("logger.log"), WithForceNewFile(true)) // cannot remove logger.log and logger.log_1
		xtesting.NotNil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", true)

		xtesting.Nil(t, f1.Close())
		xtesting.Nil(t, f2.Close())
		rl, err = New(WithFilenamePattern("logger.log"), WithForceNewFile(true)) // success to remove logger.log and logger.log_1
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", false)
		testFileExistence(t, "logger.log_1", false)

		_, err = fmt.Fprintf(rl, "hello world") // <- create new file: logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileContent(t, "logger.log", "hello world")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		xtesting.Nil(t, rl.Close())

		rl, err = New(WithFilenamePattern("logger.log"))
		xtesting.Nil(t, err)
		xtesting.Nil(t, rl.Close()) // <- r.currFile == nil
	})

	t.Run("simple write demos", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, err := New(WithFilenamePattern("logger.log"), WithRotationSize(15))
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 1\n") // <- 14, create logger.log
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

		rl, err = New(WithFilenamePattern("logger.log"), WithRotationSize(15))
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- create: logger.log_2
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log_1", true)
		testFileExistence(t, "logger.log_2", true)
		testFileContent(t, "logger.log_1", "hello world 3\n")
		testFileContent(t, "logger.log_2", "hello world 4\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log_2")
		xtesting.Nil(t, rl.Close())

		rl, err = New(WithFilenamePattern("logger.log"), WithRotationSize(15), WithForceNewFile(true)) // <- delete logger.log and logger.log_*
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", false)
		testFileExistence(t, "logger.log_1", false)
		testFileExistence(t, "logger.log_2", false)
		_, err = fmt.Fprintf(rl, "hello world 5\n") // <- create logger.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 6\n") // <- create logger.log
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "hello world 5\nhello world 6\n")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		xtesting.Nil(t, rl.Close())
	})

	t.Run("simple rotate and write", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, err := New(WithFilenamePattern("logger.log"), WithRotationSize(15))
		xtesting.Nil(t, err)
		xtesting.Equal(t, rl.CurrentFilename(), "")
		xtesting.Nil(t, rl.Rotate()) // <- create logger.log
		testFileExistence(t, "logger.log", true)
		testFileContent(t, "logger.log", "")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		xtesting.Nil(t, rl.Rotate()) // <- use logger.log
		testFileExistence(t, "logger.log", true)
		testFileExistence(t, "logger.log_1", false)
		testFileContent(t, "logger.log", "")
		xtesting.Equal(t, rl.CurrentFilename(), "logger.log")
		_, err = fmt.Fprintf(rl, "hello world 1\n") // <- write to logger.log
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

		rl, err = New(WithFilenamePattern("logger.log"), WithRotationSize(15))
		xtesting.Nil(t, err)
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
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("logger.%Y%m%d.log"),
			WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxAge(3*24*time.Hour), WithClock(clock))

		*nowP = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		_, _ = fmt.Fprintf(rl, "hello world 4\n")
		_, _ = fmt.Fprintf(rl, "hello world 5\n") // <- logger.01.log_2
		_ = os.Chtimes("logger.20010101.log", now, now)
		_ = os.Chtimes("logger.20010101.log_1", now, now)
		_ = os.Chtimes("logger.20010101.log_2", now, now)
		*nowP = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 6\n") // <- logger.02.log
		_, _ = fmt.Fprintf(rl, "hello world 7\n")
		_, _ = fmt.Fprintf(rl, "hello world 8\n") // <- logger.02.log_1
		_ = os.Chtimes("logger.20010102.log", now, now)
		_ = os.Chtimes("logger.20010102.log_1", now, now)
		*nowP = xtime.SetDay(now, 3)              // 3d
		_, _ = fmt.Fprintf(rl, "hello world 9\n") // <- logger.03.log
		_ = os.Chtimes("logger.20010103.log", now, now)
		testFileExistence(t, "logger.20010101.log", true)
		testFileExistence(t, "logger.20010101.log_1", true)
		testFileExistence(t, "logger.20010101.log_2", true)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)

		*nowP = xtime.SetHour(now, 1)
		*nowP = xtime.SetDay(now, 4) // 4d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)
		*nowP = xtime.SetDay(now, 5)         // 5d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", true)
		*nowP = xtime.SetDay(now, 6) // 6d1h
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
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("logger.%Y%m%d.log"),
			WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxCount(3), WithClock(clock))

		*nowP = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		_, _ = fmt.Fprintf(rl, "hello world 4\n")
		_, _ = fmt.Fprintf(rl, "hello world 5\n") // <- logger.01.log_2
		*nowP = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 6\n") // <- logger.02.log
		_, _ = fmt.Fprintf(rl, "hello world 7\n")
		_, _ = fmt.Fprintf(rl, "hello world 8\n") // <- logger.02.log_1
		*nowP = xtime.SetDay(now, 3)              // 3d
		_, _ = fmt.Fprintf(rl, "hello world 9\n") // <- logger.03.log
		testFileExistence(t, "logger.20010101.log", true)
		testFileExistence(t, "logger.20010101.log_1", true)
		testFileExistence(t, "logger.20010101.log_2", true)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)

		*nowP = xtime.SetHour(now, 1)
		*nowP = xtime.SetDay(now, 4) // 4d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileExistence(t, "logger.20010102.log_1", true)
		testFileExistence(t, "logger.20010103.log", true)
		*nowP = xtime.SetDay(now, 5)         // 5d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", true)
		*nowP = xtime.SetDay(now, 6) // 6d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.03.log
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010101.log_2", false)
		testFileExistence(t, "logger.20010102.log", false)
		testFileExistence(t, "logger.20010102.log_1", false)
		testFileExistence(t, "logger.20010103.log", false)

		f, err := os.OpenFile("logger.20010104.log", os.O_APPEND, 0644)
		xtesting.Nil(t, err)
		*nowP = xtime.SetDay(now, 7) // 7d1h
		xtesting.Nil(t, rl.Rotate()) // <- warning: failed to remove logger.04.log
		testFileExistence(t, "logger.20010104.log", true)
		xtesting.Nil(t, f.Close())
		*nowP = xtime.SetDay(now, 8) // 8d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.04.log and logger.05.log
		testFileExistence(t, "logger.20010104.log", false)
		testFileExistence(t, "logger.20010105.log", false)
		_ = rl.Close()
	})

	t.Run("cover errors", func(t *testing.T) {
		removeLoggers()
		now := time.Date(2001, 1, 1, 1, 1, 1, 0, time.FixedZone("", 8*60*60))
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })

		// 1.
		rl, _ := New(WithFilenamePattern("logger.%Y%m%d.log"), WithRotationTime(time.Hour*24), WithRotationMaxCount(2), WithClock(clock))
		*nowP = xtime.SetDay(now, 1)
		_, _ = fmt.Fprintf(rl, "hello world") // <- create logger.01.log
		oldPattern := rl.globPattern
		rl.globPattern = "[]"
		*nowP = xtime.SetDay(now, 2)
		_, err := fmt.Fprintf(rl, "hello world") // <- create logger.02.log, with warning
		xtesting.Nil(t, err)
		*nowP = xtime.SetDay(now, 3)
		xtesting.NotNil(t, rl.Rotate()) // <- still error
		rl.globPattern = oldPattern
		*nowP = xtime.SetDay(now, 4)
		xtesting.Nil(t, rl.Rotate()) // <- success
		xtesting.Nil(t, rl.Close())
		removeLoggers()
	})
}

func TestDiffDir(t *testing.T) {
	t.Run("force new file", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		xtesting.Nil(t, os.MkdirAll("./_test/_test", 0644))
		f1, _ := os.Create("./_test/_test/logger.log")
		f2, _ := os.Create("./_test/_test/logger.log_1")
		testFileExistence(t, "./_test/_test/logger.log", true)
		testFileExistence(t, "./_test/_test/logger.log_1", true)
		rl, _ := New(WithFilenamePattern("./_test/_test/logger.log"), WithForceNewFile(true)) // cannot remove logger.log and logger.log_1
		testFileExistence(t, "./_test/_test/logger.log", true)
		testFileExistence(t, "./_test/_test/logger.log_1", true)

		xtesting.Nil(t, f1.Close())
		xtesting.Nil(t, f2.Close())
		rl, _ = New(WithFilenamePattern("./_test/_test/logger.log"), WithForceNewFile(true)) // success to remove logger.log and logger.log_1
		testFileExistence(t, "./_test/_test/logger.log", false)
		testFileExistence(t, "./_test/_test/logger.log_1", false)
		xtesting.Nil(t, rl.Close())
	})

	t.Run("simple write demos", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		rl, _ := New(WithFilenamePattern("./_test/_test/logger.log"), WithRotationSize(15))
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- 14, create logger.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n") // <- use logger.log
		testFileExistence(t, "./_test/_test/logger.log", true)
		testFileExistence(t, "./_test/_test/logger.log_1", false)
		testFileContent(t, "./_test/_test/logger.log", "hello world 1\nhello world 2\n")
		xtesting.Equal(t, rl.CurrentFilename(), "./_test/_test/logger.log")

		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- create: logger.log_1
		testFileExistence(t, "./_test/_test/logger.log_1", true)
		testFileContent(t, "./_test/_test/logger.log", "hello world 1\nhello world 2\n")
		testFileContent(t, "./_test/_test/logger.log_1", "hello world 3\n")
		xtesting.Equal(t, rl.CurrentFilename(), "./_test/_test/logger.log_1")
		xtesting.Nil(t, rl.Close())
	})

	t.Run("max age", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("./_test/_test/logger.%Y%m%d.log"),
			WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxAge(2*24*time.Hour), WithClock(clock))

		*nowP = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		_ = os.Chtimes("./_test/_test/logger.20010101.log", now, now)
		_ = os.Chtimes("./_test/_test/logger.20010101.log_1", now, now)
		*nowP = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 4\n") // <- logger.02.log
		_ = os.Chtimes("./_test/_test/logger.20010102.log", now, now)
		testFileExistence(t, "./_test/_test/logger.20010101.log", true)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", true)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)

		*nowP = xtime.SetHour(now, 1)
		*nowP = xtime.SetDay(now, 3) // 3d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)
		*nowP = xtime.SetDay(now, 4)         // 4d1h
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
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("./_test/_test/logger.%Y%m%d.log"),
			WithRotationSize(15), WithRotationTime(time.Hour*24), WithRotationMaxCount(2), WithClock(clock))

		*nowP = xtime.SetDay(now, 1)              // 1d
		_, _ = fmt.Fprintf(rl, "hello world 1\n") // <- logger.01.log
		_, _ = fmt.Fprintf(rl, "hello world 2\n")
		_, _ = fmt.Fprintf(rl, "hello world 3\n") // <- logger.01.log_1
		*nowP = xtime.SetDay(now, 2)              // 2d
		_, _ = fmt.Fprintf(rl, "hello world 4\n") // <- logger.02.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", true)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", true)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)

		*nowP = xtime.SetHour(now, 1)
		*nowP = xtime.SetDay(now, 3) // 3d1h
		xtesting.Nil(t, rl.Rotate()) // <- delete logger.01.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", true)
		*nowP = xtime.SetDay(now, 4)         // 4d1h
		_, _ = fmt.Fprintf(rl, "for rotate") // <- delete logger.02.log
		testFileExistence(t, "./_test/_test/logger.20010101.log", false)
		testFileExistence(t, "./_test/_test/logger.20010101.log_1", false)
		testFileExistence(t, "./_test/_test/logger.20010102.log", false)
		_ = rl.Close()
	})

	t.Run("cover errors", func(t *testing.T) {
		removeLoggers()
		now := time.Date(2001, 1, 1, 1, 1, 1, 0, time.FixedZone("", 8*60*60))
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })

		// 1.
		rl, _ := New(WithFilenamePattern("./_test/logger.%Y%m%d.log"), WithRotationTime(time.Hour*24), WithClock(clock))
		*nowP = xtime.SetDay(now, 1)
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
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("logger.%Y%m%d.log"), WithSymlinkFilename("logger.current.log"), WithRotationSize(15), WithRotationMaxCount(2), WithClock(clock))
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
		*nowP = xtime.SetDay(now, 2)                // 2d
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- create: logger.02.log
		testFileContent(t, "logger.20010102.log", "hello world 4\n")
		testFileContent(t, "logger.current.log", "hello world 4\n")
		*nowP = xtime.SetDay(now, 3)                // 3d
		_, err = fmt.Fprintf(rl, "hello world 5\n") // <- create: logger.03.log, need rotate
		xtesting.Nil(t, err)
		testFileExistence(t, "logger.20010101.log", false)
		testFileExistence(t, "logger.20010101.log_1", false)
		testFileExistence(t, "logger.20010102.log", true)
		testFileContent(t, "logger.20010103.log", "hello world 5\n")
		testFileContent(t, "logger.current.log", "hello world 5\n")
		_ = rl.Close()
		removeLoggers()

		rl, _ = New(WithFilenamePattern("logger.%Y%m%d.log"), WithSymlinkFilename("logger.current.log"), WithRotationSize(15), WithRotationMaxAge(time.Hour*24*2), WithClock(clock), WithForceNewFile(true))
		*nowP = xtime.SetDay(now, 1)                // 1d
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
		*nowP = xtime.SetDay(now, 2)                // 2d
		_, err = fmt.Fprintf(rl, "hello world 4\n") // <- create: logger.02.log
		_ = os.Chtimes("logger.20010102.log", now, now)
		testFileContent(t, "logger.20010102.log", "hello world 4\n")
		testFileContent(t, "logger.current.log", "hello world 4\n")
		*nowP = xtime.SetHour(now, 1)
		*nowP = xtime.SetDay(now, 3) // 3d1h
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
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("_test/logger.%Y%m%d.log"), WithSymlinkFilename("_test/_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		_, err := fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		testFileContent(t, "_test/logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileExistence(t, "_test/_test/logger.current.log", true)
		testFileContent(t, "_test/_test/logger.current.log", "hello world 1\nhello world 2\n")
		_ = rl.Close()
		removeLoggers()

		rl, _ = New(WithFilenamePattern("_test/_test/logger.%Y%m%d.log"), WithSymlinkFilename("_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		_, err = fmt.Fprintf(rl, "hello world 1\n") // <- create: logger.01.log
		xtesting.Nil(t, err)
		_, err = fmt.Fprintf(rl, "hello world 2\n") // <- use: logger.01.log
		xtesting.Nil(t, err)
		testFileContent(t, "_test/_test/logger.20010101.log", "hello world 1\nhello world 2\n")
		testFileExistence(t, "_test/logger.current.log", true)
		testFileContent(t, "_test/logger.current.log", "hello world 1\nhello world 2\n")
		_ = rl.Close()
		removeLoggers()
	})

	t.Run("cover errors", func(t *testing.T) {
		removeLoggers()
		defer removeLoggers()

		now := time.Date(2001, 1, 1, 0, 0, 0, 0, time.FixedZone("", 8*60*60))
		nowP := &now
		clock := clockFn(func() time.Time { return *nowP })
		rl, _ := New(WithFilenamePattern("logger.%Y%m%d.log"), WithSymlinkFilename("_test/logger.current.log"), WithRotationSize(15), WithClock(clock))
		*nowP = xtime.SetDay(now, 1)
		_t_testHookSymlinkMkdir = func() {
			f, err := os.OpenFile("_test", os.O_CREATE, 0644)
			xtesting.Nil(t, err)
			xtesting.Nil(t, f.Close())
		}
		xtesting.Nil(t, rl.Rotate()) // MkdirAll failed
		_t_testHookSymlinkMkdir = nil

		*nowP = xtime.SetDay(now, 2)
		_t_testHookSymlinkMkdir2 = func() string {
			return ".." // hack
		}
		xtesting.Nil(t, rl.Rotate()) // Rel failed
		_t_testHookSymlinkMkdir2 = nil

		*nowP = xtime.SetDay(now, 3)
		_t_testHookSymlinkMkdir3 = func() {
			xtesting.Nil(t, os.MkdirAll("logger.20010103.log_symlink", 0755))
		}
		xtesting.Nil(t, rl.Rotate()) // Symlink failed

		_ = os.RemoveAll("_test")
		*nowP = xtime.SetDay(now, 4)
		_t_testHookSymlinkMkdir3 = func() {
			xtesting.Nil(t, os.MkdirAll("_test/logger.current.log", 0755))
		}
		xtesting.Nil(t, rl.Rotate()) // Rename failed
		_t_testHookSymlinkMkdir3 = nil

		_ = rl.Close()
	})
}
