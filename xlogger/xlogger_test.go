package xlogger

import (
	"testing"
)

func TestStdLogger(t *testing.T) {
	StdLogger.Println("")
	StdLogger.Println("test")
}
