//+build nolog

package dlog

// Logging levels
const (
	NONE = iota
	ERROR
	WARN
	INFO
	VERBOSE
)

// The nolog file serves to remove all logging functionality.
// this is in case logging is suspected to cause performance
// issues, i.e. in a final release, without having to strip
// code of calls to logging functions.

// In practice, logging doesn't appear to affect performance terribly.

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

func FileWrite(in ...interface{}) {}
