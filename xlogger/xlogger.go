package xlogger

import (
	"log"
	_ "unsafe"
)

//go:linkname std log.std

// std represents the standard internal logger creates by log.New, that is `log.New(os.Stderr, "", LstdFlags)`.
var std *log.Logger

// Std returns the standard logger creates by log.New, and it equals to log.Default.
func Std() *log.Logger {
	return std
}

var _ StdLogger = (*log.Logger)(nil)

// StdLogger describes how log.Logger works like, includes Print, Panic, Fatal series methods.
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
}
