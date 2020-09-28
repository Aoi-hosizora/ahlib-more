package xlogger

import (
	"testing"
)

func TestStdLogger(t *testing.T) {
	StdLogger.Println("")
	// Println arg list ends with redundant newline
	// StdLogger.Println("\n")
	StdLogger.Println("test")
}

// noinspection GoPrintFunctions
func TestIStdLogger(t *testing.T) {
	var l IStdLogger = StdLogger
	l.Println("")
	l.Println("test")
}
