package xlogger

import (
	"log"
	"os"
)

// Std is the standard logger creates by log.New.
var Std = log.New(os.Stderr, "", log.LstdFlags)

var _ StdLogger = (*log.Logger)(nil)

// StdLogger describes how log.Logger works like, includes Print, Panic, Fatal series methods.
type StdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
}
