package xlogger

import (
	"log"
	"os"
)

// Standard logger defined by log.New.
var StdLogger = log.New(os.Stderr, "", log.LstdFlags)

// IStdLogger defines how log.Logger works like.
type IStdLogger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}
