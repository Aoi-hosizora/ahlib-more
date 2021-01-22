package xlogger

import (
	"testing"
)

func TestStd(t *testing.T) {
	Std.Print("test")
	Std.Printf("%s", "test")
	Std.Println("test")
	func() {
		defer func() { recover() }()
		Std.Panic("test")
	}()
	func() {
		defer func() { recover() }()
		Std.Panicf("test")
	}()
	func() {
		defer func() { recover() }()
		Std.Panicln("test")
	}()
}

func TestStdLogger(t *testing.T) {
	var l StdLogger = Std
	l.Print("test")
	l.Printf("%s", "test")
	l.Println("test")
	func() {
		defer func() { recover() }()
		l.Panic("test")
	}()
	func() {
		defer func() { recover() }()
		l.Panicf("test")
	}()
	func() {
		defer func() { recover() }()
		l.Panicln("test")
	}()
}
