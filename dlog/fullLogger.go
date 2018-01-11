package dlog

// A FullLogger supports, in addition to Logger's functions,
// the ability to set and get a log level, and create and write
// logs directly to file.
type FullLogger interface {
	Logger
	GetLogLevel() Level
	FileWrite(in ...interface{})
	SetDebugFilter(filter string)
	SetDebugLevel(dL Level)
	CreateLogFile()
}

var fullOakLogger FullLogger

// GetLogLevel returns the log level of the fullOakLogger, or
// NONE if there is not fullOakLogger.
var GetLogLevel = func() Level {
	return NONE
}

// FileWrite logs by writing to file (if possible) but does
// not log to console as well.
// This is a NOP if fullOakLogger is not set by SetLogger.
var FileWrite = func(...interface{}) {}

// SetDebugFilter defines a string that all logs should be
// checked against-- if the log message does not contain
// the input string the log will not log.
// This is a NOP if fullOakLogger is not set by SetLogger.
var SetDebugFilter = func(string) {}

// SetDebugLevel sets the log level of the fullOakLogger.
// This is a NOP if fullOakLogger is not set by SetLogger.
var SetDebugLevel = func(Level) {}

// CreateLogFile creates a file for logs to be written to
// by log functions.
// This is a NOP if fullOakLogger is not set by SetLogger.
var CreateLogFile = func() {}
