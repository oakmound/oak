package dlog

import (
	"io"
)

// A Logger is a minimal log interface for the content oak wants to log:
// four levels of logging.
type Logger interface {
	Error(...interface{})
	Info(...interface{})
	Verb(...interface{})
	SetFilter(func(string) bool)
	GetLogLevel() Level
	SetLogLevel(l Level) error
	SetOutput(io.Writer)
}

// DefaultLogger is the Logger which all oak log messages are passed through.
var DefaultLogger Logger = NewLogger()

// ErrorCheck checks that the input is not nil, then calls Error on it if it is
// not. Otherwise it does nothing.
// Emits the input error as is for additional processing if desired.
func ErrorCheck(in error) error {
	if in != nil {
		Error(in)
	}
	return in
}

// Error will write a log if the debug level is not NONE
func Error(vs ...interface{}) {
	DefaultLogger.Error(vs...)
}

// Info will write a log if the debug level is higher than ERROR
func Info(vs ...interface{}) {
	DefaultLogger.Info(vs...)
}

// Verb will write a log if the debug level is higher than INFO
func Verb(vs ...interface{}) {
	DefaultLogger.Verb(vs...)
}

// GetLogLevel returns the set logger's log level
func GetLogLevel() Level {
	return DefaultLogger.GetLogLevel()
}

// SetFilter defines a custom filter function. Log lines that
// return false when passed to this function will not be output.
func SetFilter(filter func(string) bool) {
	DefaultLogger.SetFilter(filter)
}

// SetLogLevel sets the log level of the default logger.
func SetLogLevel(l Level) error {
	return DefaultLogger.SetLogLevel(l)
}

func SetOutput(w io.Writer) {
	DefaultLogger.SetOutput(w)
}
