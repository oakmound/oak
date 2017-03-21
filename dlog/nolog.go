//+build nolog

// Package dlog provides logging functions with
// caller file and line information,
// logging levels and level and text filters.
package dlog

// Logging levels
const (
	NONE = iota
	ERROR
	WARN
	INFO
	VERBOSE
)

func SetDebugFilter(filter string) {
}

func SetDebugLevel(dL int) {
}

func CreateLogFile() {
}

func Error(in ...interface{}) {
}

func Warn(in ...interface{}) {
}

func Info(in ...interface{}) {
}

func Verb(in ...interface{}) {
}

func SetStringDebugLevel(debugL string) {
}
