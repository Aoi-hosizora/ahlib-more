package xlogger

import (
	"log"
	_ "unsafe"
)

//go:linkname Std log.std

// Std represents the standard logger creates by log.New, that equals to `log.New(os.Stderr, "", log.LstdFlags)`.
var Std *log.Logger

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
