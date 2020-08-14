package xlogger

import (
	"log"
	"os"
)

var StdLogger = log.New(os.Stderr, "", log.LstdFlags)
